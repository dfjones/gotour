package main

import (
    "fmt"
)

type Fetcher interface {
    // Fetch returns the body of URL and
    // a slice of URLs found on that page.
    Fetch(url string) (body string, urls []string, err error)
}

type visitMap map[string]bool

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.

func SerialCrawl(url string, depth int, fetcher Fetcher) (map[string]bool) {
  vm := make(visitMap)
  crawl(url, depth, fetcher, vm)
  return vm
}

func crawl(url string, depth int, fetcher Fetcher, visited visitMap) {
    if depth <= 0 {
        return
    }
    if found, ok := visitMap[url]; ok {
      return
    }
    body, urls, err := fetcher.Fetch(url)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("found: %s %q\n", url, body)
    for _, u := range urls {
        Crawl(u, depth-1, fetcher)
    }
    return
}

func main() {
    SerialCrawl("http://golang.org/", 4, fetcher)
}
