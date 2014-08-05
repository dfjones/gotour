package recursive

import (
  "fmt"
  "github.com/dfjones/gotour/crawler"
  "github.com/dfjones/gotour/crawler/visitreg"
  "github.com/dfjones/gotour/crawler/visitreg/channel"
  "github.com/dfjones/gotour/crawler/visitreg/mutex"
)

// ChanCrawl begins a crawl operation using the channel based visited map
func ChanCrawl(url string, depth int, fetcher crawler.Fetcher) map[string]bool {
  return startCrawl(url, depth, channel.New(), fetcher)
}

// MutexCrawl begins a crawl operation using the mutex based visited map
func MutexCrawl(url string, depth int, fetcher crawler.Fetcher) map[string]bool {
  return startCrawl(url, depth, mutex.New(), fetcher)
}

func startCrawl(url string, depth int, visitRegister visitreg.VisitRegister, fetcher crawler.Fetcher) map[string]bool {
  done := make(chan struct{})
  go crawl(done, url, depth, visitRegister, fetcher)
  <-done
  return visitRegister.Close()
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func crawl(done chan<- struct{}, url string, depth int, visitRegister visitreg.VisitRegister, fetcher crawler.Fetcher) {
  defer func() {
    done <- struct{}{}
  }()
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
  childrenDone := make(chan struct{})
  for _, u := range urls {
    go crawl(childrenDone, u, depth-1, visitRegister, fetcher)
  }
  for i := 0; i < len(urls); {
    <-childrenDone
    i++
  }
}
