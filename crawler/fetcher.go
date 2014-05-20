package main

import (
  "fmt"
  "math/rand"
  "time"
)

type Fetcher interface {
  // Fetch returns the body of URL and
  // a slice of URLs found on that page.
  Fetch(url string) (body string, urls []string, err error)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
  body string
  urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
  time.Sleep(100 * time.Millisecond)
  if res, ok := f[url]; ok {
    return res.body, res.urls, nil
  }
  return "", nil, fmt.Errorf("not found: %s", url)
}

func NewFakeFetcher(results, minLinks, maxLinks int) Fetcher {
  ff := make(fakeFetcher)
  for i := 0; i < results; i++ {
    body := fmt.Sprint(i)
    ff[body] = &fakeResult{body, []string{}}
  }
  urls := make([]string, results)
  i := 0
  for key, _ := range ff {
    urls[i] = key
    i++
  }
  for key, _ := range ff {
    numLinks := rand.Intn(maxLinks - minLinks)
    numLinks += minLinks
    for l := 0; l < numLinks; l++ {
      link := key
      for link == key {
        i := rand.Intn(len(urls))
        link = urls[i]
      }
      fr := ff[key]
      fr.urls = append(fr.urls, link)
    }
  }
  return ff
}

// fetcher is a populated fakeFetcher.
var FakeFetcher = fakeFetcher{
  "http://golang.org/": &fakeResult{
    "The Go Programming Language",
    []string{
      "http://golang.org/pkg/",
      "http://golang.org/cmd/",
    },
  },
  "http://golang.org/pkg/": &fakeResult{
    "Packages",
    []string{
      "http://golang.org/",
      "http://golang.org/cmd/",
      "http://golang.org/pkg/fmt/",
      "http://golang.org/pkg/os/",
    },
  },
  "http://golang.org/pkg/fmt/": &fakeResult{
    "Package fmt",
    []string{
      "http://golang.org/",
      "http://golang.org/pkg/",
    },
  },
  "http://golang.org/pkg/os/": &fakeResult{
    "Package os",
    []string{
      "http://golang.org/",
      "http://golang.org/pkg/",
    },
  },
}
