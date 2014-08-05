package pool

import (
  "fmt"
  "github.com/dfjones/gotour/crawler"
  "github.com/dfjones/gotour/crawler/visitreg"
  "github.com/dfjones/gotour/crawler/visitreg/channel"
  "github.com/dfjones/gotour/crawler/visitreg/mutex"
  "sync/atomic"
)

func ChanCrawl(url string, depth int, fetcher crawler.Fetcher) map[string]bool {
  return startCrawl(url, depth, channel.New(), fetcher)
}

func MutexCrawl(url string, depth int, fetcher crawler.Fetcher) map[string]bool {
  return startCrawl(url, depth, mutex.New(), fetcher)
}

func startCrawl(url string, depth int, visitRegister visitreg.VisitRegister, fetcher crawler.Fetcher) map[string] bool {
  workers := 4
  done := make(chan struct{})
  defer close(done)
  // todo: this arbitrary buffer size doesn't seem like a great idea
  workChan = make(chan *crawlParams, workers*64)
  defer close(workChan)
  for i := 0; i < workers; i++ {
    go func() {
      for p := range workChan {
        crawl(done, p.url, p.depth, visitRegister, fetcher)
      }
    }()
  }
  start(1)
  workChan <- &crawlParams{
    url,
    depth,
  }
  <-done
  return visitRegister.Close()
}

type crawlParams struct {
  url string
  depth int
}

var workChan chan *crawlParams

var routineCount int64

func start(count int) {
  atomic.AddInt64(&routineCount, int64(count))
}

func finish(done chan<- struct{}) {
  remaining := atomic.AddInt64(&routineCount, -1)
  if remaining == 0 {
    done <- struct{}{}
  }
}

func crawl(done chan<- struct{}, url string, depth int, visitRegister visitreg.VisitRegister, fetcher crawler.Fetcher) {
  defer finish(done)
  if depth <= 0 {
    return
  }
  if visitRegister.IsVisited(url) {
    return
  }
  visitRegister.Visit(url)
  _, urls, err := fetcher.Fetch(url)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Printf("found: %s\n", url)
  start(len(urls))
  for _, u := range urls {
    workChan <- &crawlParams{u, depth-1}
  }
}