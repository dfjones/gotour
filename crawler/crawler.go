package crawler

type CrawlFunc func(url string, depth int, fetcher Fetcher) map[string]bool

type VisitMap map[string]bool
