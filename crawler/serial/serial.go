package serial

import (
  "fmt"
  "github.com/dfjones/gotour/crawler"
)

type visitMap map[string]bool

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.

// SerialCrawl performs the crawl operation with no concurrency
func Crawl(url string, depth int, fetcher crawler.Fetcher) map[string]bool {
  vm := make(visitMap)
  serialcrawl(url, depth, fetcher, vm)
  return vm
}

func serialcrawl(url string, depth int, fetcher crawler.Fetcher, visited visitMap) {
  if depth <= 0 {
    return
  }
  if _, ok := visited[url]; ok {
    return
  }
  visited[url] = true
  _, urls, err := fetcher.Fetch(url)
  if err != nil {
    fmt.Println(err)
    return
  }
  //fmt.Printf("found: %s %q\n", url, body)
  for _, u := range urls {
    serialcrawl(u, depth-1, fetcher, visited)
  }
  return
}
