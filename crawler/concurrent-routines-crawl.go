package main

import (
  "fmt"
  "github.com/dfjones/gotour/crawler/visitreg"
  "github.com/dfjones/gotour/crawler/visitreg/channel"
  "github.com/dfjones/gotour/crawler/visitreg/mutex"
  "sync/atomic"
)

// ChanRoutineCrawl begins a crawl operation using the channel based visited map
func ChanRoutineCrawl(url string, depth int, fetcher Fetcher) map[string]bool {
  return startCrawl(url, depth, channel.New(), fetcher)
}

// MutexRoutineCrawl begins a crawl operation using the mutex based visited map
func MutexRoutineCrawl(url string, depth int, fetcher Fetcher) map[string]bool {
  return startCrawl(url, depth, mutex.New(), fetcher)
}

func startRoutineCrawl(url string, depth int, visitRegister visitreg.VisitRegister, fetcher Fetcher) map[string]bool {
  done := make(chan struct{})
  startRoutines(1)
  go routinecrawl(done, url, depth, visitRegister, fetcher)
  <-done
  return visitRegister.Close()
}

var routineCount int64

func startRoutines(count int) {
  atomic.AddInt64(&routineCount, int64(count))
}

func finishRoutine(done chan<- struct{}) {
  remaining := atomic.AddInt64(&routineCount, -1)
  if remaining == 0 {
    done <- struct{}{}
  }
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func routinecrawl(done chan<- struct{}, url string, depth int, visitRegister visitreg.VisitRegister, fetcher Fetcher) {
  // Done: Fetch URLs in parallel.
  // Done: Don't fetch the same URL twice.
  // This implementation doesn't do either:
  defer finishRoutine(done)
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
  //fmt.Printf("found: %s %q\n", url, body)
  startRoutines(len(urls))
  for _, u := range urls {
    go reccrawl(done, u, depth-1, visitRegister, fetcher)
  }
}
