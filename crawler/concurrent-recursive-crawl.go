package main

import (
  "fmt"
  "github.com/dfjones/tour/crawler/visitreg"
  "github.com/dfjones/tour/crawler/visitreg/channel"
  "github.com/dfjones/tour/crawler/visitreg/mutex"
  "runtime"
)

// ChanCrawl begins a crawl operation using the channel based visited map
func ChanCrawl(url string, depth int, fetcher Fetcher) map[string]bool {
  return startCrawl(url, depth, channel.New(), fetcher)
}

// MutexCrawl begins a crawl operation using the mutex based visited map
func MutexCrawl(url string, depth int, fetcher Fetcher) map[string]bool {
  return startCrawl(url, depth, mutex.New(), fetcher)
}

func startCrawl(url string, depth int, visitRegister visitreg.VisitRegister, fetcher Fetcher) map[string]bool {
  done := make(chan struct{})
  go chancrawl(done, url, depth, visitRegister, fetcher)
  <-done
  return visitRegister.Close()
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func chancrawl(done chan<- struct{}, url string, depth int, visitRegister visitreg.VisitRegister, fetcher Fetcher) {
  // Done: Fetch URLs in parallel.
  // Done: Don't fetch the same URL twice.
  // This implementation doesn't do either:
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
    go chancrawl(childrenDone, u, depth-1, visitRegister, fetcher)
  }
  for i := 0; i < len(urls); {
    <-childrenDone
    i++
  }
}

func main() {
  numCpus := runtime.NumCPU()
  fmt.Println("Running using num cpus: ", numCpus)
  runtime.GOMAXPROCS(runtime.NumCPU())
  graph := NewFakeFetcher(100000, 1, 10)
  fmt.Println("Graph built")
  visited := ChanCrawl("0", 4, graph)
  fmt.Println("Visited: ", len(visited))
}
