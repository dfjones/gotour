package main

import (
  "fmt"
  "reflect"
  "runtime"
  "sort"
  "testing"
)

var depth = 6
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

  if !mapKeysEqual(v1, v2) {
    t.Errorf("%v != %v", v1, v2)
  }

  if !mapKeysEqual(v1, v3) {
    t.Errorf("%v != %v", v1, v3)
  }
}

func mapKeysEqual(v1, v2 visitMap) bool {
  k1 := keys(v1)
  k2 := keys(v2)
  sort.Strings(k1)
  sort.Strings(k2)
  return reflect.DeepEqual(k1, k2)
}

func keys(m visitMap) []string {
  keys := []string{}
  for k := range m {
    keys = append(keys, k)
  }
  return keys
}

func BenchmarkSerialCrawl(b *testing.B) {
  for i := 0; i < b.N; i++ {
    SerialCrawl("0", depth, graph)
  }
}

func BenchmarkRecursiveChanCrawl(b *testing.B) {
  for i := 0; i < b.N; i++ {
    ChanCrawl("0", depth, graph)
  }
}

func BenchmarkRecursiveMutexCrawl(b *testing.B) {
  for i := 0; i < b.N; i++ {
    MutexCrawl("0", depth, graph)
  }
}
