package main

import (
    "fmt"
)


/*
insructions:
In this exercise you'll use Go's concurrency features to parallelize a web crawler.

Modify the Crawl function to fetch URLs in parallel without fetching the same URL twice.
*/

type Fetcher interface {
    // Fetch returns the body of URL and
    // a slice of URLs found on that page.
    Fetch(url string) (body string, urls []string, err error)
}

type VisitRegister interface {
  Visit(string)
  IsVisited(string) bool
  Close()
}

type visitMap chan visitData

type visitData struct {
  action visitAction
  url string
  result chan<- bool
}

type visitAction int

const (
  add visitAction = iota
  has
  end
)

func NewVisitRegister() VisitRegister {
  vm := make(visitMap)
  go vm.run()
  return vm
}

func (vm visitMap) run() {
  store := make(map[string]bool)
  for command := range vm {
    switch command.action {
      case add:
        store[command.url] = true
      case has:
        _, found := store[command.url]
        command.result <- found
      case end:
        close(vm)
    }
  }
}

func (vm visitMap) Visit(url string) {
  vm <- visitData{action: add, url: url}
}

func (vm visitMap) IsVisited(url string) bool {
  reply := make(chan bool)
  vm <- visitData{action: has, url: url, result: reply}
  return <-reply
}

func (vm visitMap) Close() {
  vm <- visitData{action: end}
}

func Crawl(url string, depth int, fetcher Fetcher) {
  visitRegister := NewVisitRegister()
  done := make(chan struct{})
  go crawl(done, url, depth, visitRegister, fetcher)
  <-done
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func crawl(done chan<- struct{}, url string, depth int, visitRegister VisitRegister, fetcher Fetcher) {
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
    body, urls, err := fetcher.Fetch(url)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("found: %s %q\n", url, body)
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
    Crawl("http://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
    body string
    urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
    if res, ok := f[url]; ok {
        return res.body, res.urls, nil
    }
    return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
    "http://golang.org/": &fakeResult{
        "The Go Programming Language",
        []string{
            "http://golang.org/pkg/",
            "http://golang.org/cmd/",
        },
    },
    "http://golang.org/pkg/": &fakeResult{
        "Packages",
        []string{
            "http://golang.org/",
            "http://golang.org/cmd/",
            "http://golang.org/pkg/fmt/",
            "http://golang.org/pkg/os/",
        },
    },
    "http://golang.org/pkg/fmt/": &fakeResult{
        "Package fmt",
        []string{
            "http://golang.org/",
            "http://golang.org/pkg/",
        },
    },
    "http://golang.org/pkg/os/": &fakeResult{
        "Package os",
        []string{
            "http://golang.org/",
            "http://golang.org/pkg/",
        },
    },
}
