package main

import (
    "fmt"
    "runtime"
)

type VisitRegister interface {
  Visit(string)
  IsVisited(string) bool
  Close() map[string]bool
}

type visitMap chan visitData

type visitData struct {
  action visitAction
  url string
  result chan<- bool
  data chan<- map[string]bool
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
        command.data <- store
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

func (vm visitMap) Close() (map[string]bool){
  data := make (chan map[string]bool)
  vm <- visitData{action: end, data: data}
  return <-data
}

func Crawl(url string, depth int, fetcher Fetcher) (map[string]bool){
  visitRegister := NewVisitRegister()
  done := make(chan struct{})
  go crawl(done, url, depth, visitRegister, fetcher)
  <-done
  return visitRegister.Close()
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
  visited := Crawl("0", 100, graph)
  fmt.Println("Visited: ", len(visited))
}
