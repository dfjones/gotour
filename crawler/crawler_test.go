package main

import (
  "fmt"
  "runtime"
  "testing"
)

var depth = 50
var graphSize = 100000
var graph Fetcher

func init() {
  numCpus := runtime.NumCPU()
  runtime.GOMAXPROCS(runtime.NumCPU())
  fmt.Println("Running using num cpus: ", numCpus)
  graph = NewFakeFetcher(graphSize, 1, 10)
  fmt.Println("Graph built size: ", graphSize)
}

func BenchmarkChanCrawl(b *testing.B) {
  visited := ChanCrawl("0", depth, graph)
  b.StopTimer()
  b.Log("Benchmark finished. Visited: ", len(visited))
}
