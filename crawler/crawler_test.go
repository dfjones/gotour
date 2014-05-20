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

type crawlFunc func(url string, depth int, fetcher Fetcher) map[string]bool

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
  vms := []visitMap{}
  funcs := []crawlFunc{SerialCrawl, ChanCrawl, MutexCrawl,
    ChanRoutineCrawl, MutexRoutineCrawl}
  for _, f := range funcs {
    res := f("0", depth, smallGraph)
    vms = append(vms, res)
  }

  assertMapsEqual(t, vms)
}

func assertMapsEqual(t *testing.T, vms []visitMap) {
  v1 := vms[0]
  k1 := keys(v1)
  sort.Strings(k1)

  for _, v2 := range vms[1:] {
    k2 := keys(v2)
    sort.Strings(k2)
    if !reflect.DeepEqual(k1, k2) {
      t.Errorf("%v != %v", k1, k2)
    }
  }
}

func keys(m visitMap) []string {
  keys := []string{}
  for k := range m {
    keys = append(keys, k)
  }
  return keys
}

func BenchmarkSerialCrawl(b *testing.B) {
  bench(b, SerialCrawl)
}

func BenchmarkRecursiveChanCrawl(b *testing.B) {
  bench(b, ChanCrawl)
}

func BenchmarkRecursiveMutexCrawl(b *testing.B) {
  bench(b, MutexCrawl)
}

func BenchmarkChanRoutineCrawl(b *testing.B) {
  bench(b, ChanRoutineCrawl)
}

func BenchmarkMutexRoutineCrawl(b *testing.B) {
  bench(b, MutexRoutineCrawl)
}

func bench(b *testing.B, f crawlFunc) {
  var size = 0
  for i := 0; i < b.N; i++ {
    size = len(f("0", depth, graph))
  }
  b.Log("Benchmark finished, visited: ", size)
}
