package 73

var graph

func init() {
  numCpus := runtime.NumCPU()
  runtime.GOMAXPROCS(runtime.NumCPU())
  fmt.Println("Running using num cpus: ", numCpus)
  graph = NewFakeFetcher(100000, 1, 10)
  fmt.Println("Graph built size: ", len(graph))
}

func BenchmarkChanCrawl(b *testing.B) {
  visited := ChanCrawl("0", 4, graph)
}
