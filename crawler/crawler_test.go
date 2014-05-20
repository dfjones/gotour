package main

import (
	"fmt"
	"runtime"
	"testing"
)

var depth = 10
var graphSize = 100000
var graph Fetcher

func init() {
	numCpus := runtime.NumCPU()
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Running using num cpus: ", numCpus)
	graph = NewFakeFetcher(graphSize, 1, 10)
	fmt.Println("Graph built size: ", graphSize)
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
