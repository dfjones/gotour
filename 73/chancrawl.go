package main

import (
    "fmt"
    "runtime"
    "github.com/dfjones/tour/73/visitreg"
    "github.com/dfjones/tour/73/visitreg/channel"
)

func Crawl(url string, depth int, fetcher Fetcher) (map[string]bool){
  visitRegister := channel.New()
  done := make(chan struct{})
  go crawl(done, url, depth, visitRegister, fetcher)
  <-done
  return visitRegister.Close()
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func crawl(done chan<- struct{}, url string, depth int, visitRegister visitreg.VisitRegister, fetcher Fetcher) {
    // Done: Fetch URLs in parallel.
    // Done: Don't fetch the same URL twice.
    // This implementation doesn't do either:
    defer func () {
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

func main() {
  numCpus := runtime.NumCPU()
  fmt.Println("Running using num cpus: ", numCpus)
  runtime.GOMAXPROCS(runtime.NumCPU())
  graph := NewFakeFetcher(100000, 1, 10)
  fmt.Println("Graph built")
  visited := Crawl("0", 4, graph)
  fmt.Println("Visited: ", len(visited))
}
