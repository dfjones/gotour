package test

import (
  "fmt"
  "github.com/dfjones/gotour/crawler"
  "github.com/dfjones/gotour/crawler/concurrent/graph"
  "github.com/dfjones/gotour/crawler/concurrent/recursive"
  "github.com/dfjones/gotour/crawler/concurrent/pool"
  "github.com/dfjones/gotour/crawler/serial"
  "reflect"
  "runtime"
  "sort"
  "testing"
)

var depth = 6
var graphSize = 100000
var largeGraph crawler.Fetcher
var smallGraph crawler.Fetcher

func init() {
  numCpus := runtime.NumCPU()
  runtime.GOMAXPROCS(runtime.NumCPU())
  fmt.Println("Running using num cpus: ", numCpus)
  smallGraph = crawler.NewFakeFetcher(100, 1, 5)
  largeGraph = crawler.NewFakeFetcher(graphSize, 1, 10)
  fmt.Println("Graph built size: ", graphSize)
}

func TestCrawlEquiv(t *testing.T) {
  depth := 3
  vms := []crawler.VisitMap{}
  funcs := []crawler.CrawlFunc{
    serial.Crawl, 
    graph.ChanCrawl, 
    graph.MutexCrawl, 
    recursive.ChanCrawl, 
    recursive.MutexCrawl,
    pool.ChanCrawl,
    pool.MutexCrawl,
  }
  for _, f := range funcs {
    res := f("0", depth, smallGraph)
    vms = append(vms, res)
  }

  assertMapsEqual(t, vms)
}

func assertMapsEqual(t *testing.T, vms []crawler.VisitMap) {
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

func keys(m crawler.VisitMap) []string {
  keys := []string{}
  for k := range m {
    keys = append(keys, k)
  }
  return keys
}

func BenchmarkSerialCrawl(b *testing.B) {
  bench(b, serial.Crawl)
}

func BenchmarkRecursiveChanCrawl(b *testing.B) {
  bench(b, recursive.ChanCrawl)
}

func BenchmarkRecursiveMutexCrawl(b *testing.B) {
  bench(b, recursive.MutexCrawl)
}

func BenchmarkChanGraphCrawl(b *testing.B) {
  bench(b, graph.ChanCrawl)
}

func BenchmarkMutexGraphCrawl(b *testing.B) {
  bench(b, graph.MutexCrawl)
}

func BenchmarkMutexPoolCrawl(b *testing.B) {
  bench(b, pool.MutexCrawl)
}

func BenchmarkChanPoolCrawl(b *testing.B) {
  bench(b, pool.ChanCrawl)
}

func bench(b *testing.B, f crawler.CrawlFunc) {
  //var size = 0
  for i := 0; i < b.N; i++ {
    //size = len(f("0", depth, largeGraph))
    f("0", depth, largeGraph)
  }
  //b.Log("Benchmark finished, visited: ", size)
}
