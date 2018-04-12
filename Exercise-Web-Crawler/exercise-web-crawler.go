package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

func Crawl(url string, depth int, fetcher Fetcher, ch chan string, uc *UrlCache) {
	_crawl(url, depth, fetcher, ch, uc)
	close(ch)
	return
}

func _crawl(url string, depth int, fetcher Fetcher, ch chan string, uc *UrlCache) {
	if depth <= 0 || uc.Check(url) {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		//fmt.Println(err)
		ch <- fmt.Sprint(err)
		uc.Add(url)
		return
	}

	//fmt.Println("found: %s %q", url, body)
	ch <- fmt.Sprintf("found: %s %q", url, body)
	uc.Add(url)

	for _, u := range urls {
		_crawl(u, depth-1, fetcher, ch, uc)
	}

	return
}

type UrlCache struct {
	cache map[string]bool
	m     sync.Mutex
}

func (uc *UrlCache) Add(url string) {
	uc.m.Lock()
	defer uc.m.Unlock()
	uc.cache[url] = true
	return
}

func (uc *UrlCache) Check(url string) bool {
	uc.m.Lock()
	defer uc.m.Unlock()
	return uc.cache[url]
}

func main() {
	ch := make(chan string)
	uc := &UrlCache{cache: make(map[string]bool)}
	go Crawl("https://golang.org/", 4, fetcher, ch, uc)

	for v := range ch {
		fmt.Println(v)
	}
}
