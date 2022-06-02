package main

import (
	"github.com/fixwa/go-prices-tracker/crawlers"
	"github.com/fixwa/go-prices-tracker/database"
	"sync"
)

func main() {
	database.ConnectDatabase()

	var waiter sync.WaitGroup

	waiter.Add(1)
	go crawlers.CrawlDistriland(&waiter)

	waiter.Add(1)
	go crawlers.CrawlImportadoraRonson(&waiter)

	waiter.Add(1)
	go crawlers.CrawlGeeker(&waiter)

	waiter.Wait()
}
