package crawlers

import (
	"fmt"
	"github.com/fixwa/go-prices-tracker/models"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"strings"
	"sync"
	"time"
)

var (
	productsLinks        map[string]bool
	categoriesLinks      map[string]bool
	categoriesPagesLinks map[string]bool
	currentSource        *models.ProductSource
)

func init() {
	productsLinks = map[string]bool{}
	categoriesLinks = map[string]bool{}
	categoriesPagesLinks = map[string]bool{}
	currentSource = models.ProductsSources[1]
}

func CrawlImportadoraRonson(w *sync.WaitGroup) {
	log.Println("Crawling " + currentSource.Name)

	c := colly.NewCollector(
		colly.AllowedDomains(currentSource.AllowedDomains),
	)

	q, _ := queue.New(
		1,
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	detailCollector := c.Clone()
	productCollector := c.Clone()
	categoryCollector := c.Clone()

	c.OnHTML(".nav", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			categoryLink := el.Attr("href")
			if strings.Index(categoryLink, "/p/") == -1 {
				return
			}

			if _, found := categoriesLinks[categoryLink]; !found {
				fmt.Println("Category found: " + categoryLink)
				categoryCollector.Visit(categoryLink)
				categoriesLinks[categoryLink] = true
			} else {
				categoriesLinks[categoryLink] = false
			}

		})
	})

	categoryCollector.OnHTML(".woocommerce-pagination", func(e *colly.HTMLElement) {
		e.ForEach("a.page-number", func(_ int, el *colly.HTMLElement) {
			categoryLink := el.Attr("href")
			if strings.Index(categoryLink, "/page/") == -1 {
				return
			}

			if _, found := categoriesPagesLinks[categoryLink]; !found {
				fmt.Println("Category Page# Found: " + categoryLink)
				productCollector.Visit(categoryLink)
				categoriesPagesLinks[categoryLink] = true
			} else {
				categoriesPagesLinks[categoryLink] = false
			}
		})

	})

	productCollector.OnHTML(".product-small", func(e *colly.HTMLElement) {
		productLink := e.ChildAttr("a.woocommerce-LoopProduct-link", "href")
		if strings.Index(productLink, "/tienda/") == -1 {
			return
		}

		fmt.Println("Product found: " + productLink)
		if _, found := productsLinks[productLink]; !found {
			detailCollector.Visit(productLink)
			productsLinks[productLink] = true
		} else {
			productsLinks[productLink] = false
		}
	})

	detailCollector.OnHTML(".product-main", func(e *colly.HTMLElement) {
		title := e.ChildText("h1.product-title")
		price := "[" + e.ChildText("p.price:first-child") + "]"
		thumbnail := e.ChildAttr("img.wp-post-image", "data-src")
		publishedAt := time.Now()
		categoryName := e.ChildText("span.posted_in")

		product := &models.Product{
			Title:        title,
			Content:      title,
			Source:       currentSource.ID,
			URL:          e.Request.URL.String(),
			Price:        price,
			CategoryName: categoryName,
			Thumbnail:    thumbnail,
			PublishedAt:  publishedAt,
		}
		fmt.Printf("%v\n", product)
	})

	q.AddURL(currentSource.BaseURL)

	// Consume
	q.Run(c)
	log.Println("Finished " + currentSource.Name)
	w.Done()
}
