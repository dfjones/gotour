package main

import (
  "fmt"
  "github.com/dfjones/gotour/crawler"
  "github.com/dfjones/gotour/crawler/concurrent/graph"
  "runtime"
)

func main() {
  numCpus := runtime.NumCPU()
  fmt.Println("Running using num cpus: ", numCpus)
  runtime.GOMAXPROCS(runtime.NumCPU())
  data := crawler.NewFakeFetcher(100000, 1, 10)
  fmt.Println("Graph built")
  visited := graph.ChanCrawl("0", 4, data)
  fmt.Println("Visited: ", len(visited))
}
