package main

import (
  "fmt"
  "reflect"
  "runtime"
  "testing"
)

var depth = 8
var graphSize = 100000
var graph Fetcher
var smallGraph Fetcher

func init() {
  numCpus := runtime.NumCPU()
  runtime.GOMAXPROCS(runtime.NumCPU())
  fmt.Println("Running using num cpus: ", numCpus)
  smallGraph = NewFakeFetcher(100, 1, 5)
  graph = NewFakeFetcher(graphSize, 1, 10)
  fmt.Println("Graph built size: ", graphSize)
}

func TestCrawlEquiv(t *testing.T) {
  depth := 3
  v1 := SerialCrawl("0", depth, smallGraph)
  v2 := ChanCrawl("0", depth, smallGraph)
  v3 := MutexCrawl("0", depth, smallGraph)

  if !reflect.DeepEqual(v1, v2) {
    t.Errorf("%v != %v", v1, v2)
  }

  if !reflect.DeepEqual(v1, v3) {
    t.Errorf("%v != %v", v1, v3)
  }
}

func BenchmarkSerialCrawl(b *testing.B) {
  visited := SerialCrawl("0", depth, graph)
  b.StopTimer()
  b.Log("Benchmark finished. Visited: ", len(visited))
}

func BenchmarkRecursiveChanCrawl(b *testing.B) {
  visited := ChanCrawl("0", depth, graph)
  b.StopTimer()
  b.Log("Benchmark finished. Visited: ", len(visited))
}

func BenchmarkRecursiveMutexCrawl(b *testing.B) {
  visited := MutexCrawl("0", depth, graph)
  b.StopTimer()
  b.Log("Benchmark finished. Visited: ", len(visited))
}
