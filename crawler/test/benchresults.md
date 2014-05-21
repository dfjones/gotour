```
± go test -test.bench=. -test.benchmem
Running using num cpus:  8
Graph built size:  100000
PASS
BenchmarkSerialCrawl-8	       1	292775140168 ns/op	  230800 B/op	     192 allocs/op
BenchmarkRecursiveChanCrawl-8	       2	 616463221 ns/op	  800832 B/op	    6331 allocs/op
BenchmarkRecursiveMutexCrawl-8	       2	 610879158 ns/op	  676328 B/op	    3456 allocs/op
BenchmarkChanGraphCrawl-8	       2	 611310891 ns/op	  484080 B/op	    3265 allocs/op
BenchmarkMutexGraphCrawl-8	       2	 608324830 ns/op	  374208 B/op	     367 allocs/op
ok  	github.com/dfjones/gotour/crawler/test	302.148s
```
