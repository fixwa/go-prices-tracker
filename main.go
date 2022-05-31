package main

import (
	"github.com/fixwa/go-prices-tracker/crawlers"
	"sync"
)

func main() {
	var waiter sync.WaitGroup

	waiter.Add(1)
	go crawlers.CrawlImportadoraRonson(&waiter)

	waiter.Wait()
}
