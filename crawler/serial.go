package main

import (
    "fmt"
)

type visitMap map[string]bool

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.

func SerialCrawl(url string, depth int, fetcher Fetcher) (map[string]bool) {
  vm := make(visitMap)
  serialcrawl(url, depth, fetcher, vm)
  return vm
}

func serialcrawl(url string, depth int, fetcher Fetcher, visited visitMap) {
    if depth <= 0 {
        return
    }
    if _, ok := visited[url]; ok {
      return
    }
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
