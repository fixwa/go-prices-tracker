package main

import (
	"flag"
	"github.com/fixwa/go-prices-tracker/crawlers"
	"sync"
)

func main() {
	crawlerName := flag.String("crawler", "default", "Crawler to execute")
	flag.Parse()
	if *crawlerName == "default" {
		panic("Crawler not specified.")
	}

	var waiter sync.WaitGroup

	switch *crawlerName {
	case "distriland":
		waiter.Add(1)
		go crawlers.CrawlDistriland(&waiter)

	case "importadora-ronson":
		waiter.Add(1)
		go crawlers.CrawlImportadoraRonson(&waiter)
	case "geeker":
		waiter.Add(1)
		go crawlers.CrawlGeeker(&waiter)
	}

	waiter.Wait()
}
