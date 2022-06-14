package main

import (
	"flag"
	"fmt"
	"github.com/fixwa/go-prices-tracker/crawlers"
	"github.com/fixwa/go-prices-tracker/crawlers/distriland"
	"github.com/fixwa/go-prices-tracker/crawlers/geeker"
	"github.com/fixwa/go-prices-tracker/crawlers/importadoraronson"
	"sync"
)

func main() {
	crawlerName := flag.String("crawler", "default", "Crawler to execute.")
	clear := flag.String("clear", "default", "Clear collection first.")
	flag.Parse()

	switch *clear {
	case "all":
		crawlers.DeleteAll()
	case "distriland":
		distriland.Clear()
	case "importadora-ronson":
		importadoraronson.Clear()
	case "geeker":
		geeker.Clear()
	}

	if *crawlerName == "default" {
		fmt.Println("Crawler not specified.")
		return
	}

	var waiter sync.WaitGroup

	switch *crawlerName {
	case "distriland":
		waiter.Add(1)
		if *clear == "yes" {
			distriland.Clear()
		}
		go distriland.Crawl(&waiter)

	case "importadora-ronson":
		waiter.Add(1)
		if *clear == "yes" {
			importadoraronson.Clear()
		}
		go importadoraronson.Crawl(&waiter)

	case "geeker":
		waiter.Add(1)
		if *clear == "yes" {
			geeker.Clear()
		}
		go geeker.Crawl(&waiter)

	case "all":
		waiter.Add(3)
		go importadoraronson.Crawl(&waiter)
		go distriland.Crawl(&waiter)
		go geeker.Crawl(&waiter)
	}

	waiter.Wait()
}
