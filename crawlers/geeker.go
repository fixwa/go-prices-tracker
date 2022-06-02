package crawlers

import (
	"fmt"
	"github.com/fixwa/go-prices-tracker/models"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"sync"
	"time"
)

func CrawlGeeker(w *sync.WaitGroup) {
	var totalProductsCollected int = 0
	currentSource = models.ProductsSources[2]

	log.Println("Crawling " + currentSource.Name)

	frontPageCollector := colly.NewCollector(
		colly.AllowedDomains(currentSource.AllowedDomains),
	)

	detailCollector := colly.NewCollector(
		colly.AllowedDomains(currentSource.AllowedDomains),
	)

	productCollector := colly.NewCollector(
		colly.AllowedDomains(currentSource.AllowedDomains),
	)

	categoryCollector := colly.NewCollector(
		colly.AllowedDomains(currentSource.AllowedDomains),
	)

	q, _ := queue.New(
		1,
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	//detailCollector := c.Clone()
	//productCollector := c.Clone()
	//categoryCollector := c.Clone()

	frontPageCollector.OnHTML(".nav", func(e *colly.HTMLElement) {
		//fmt.Printf("%v\n", e.ChildText("li"))
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			categoryLink := el.Attr("href")
			if len(categoryLink) < 1 {
				return
			}

			categoryLink = currentSource.BaseURL + categoryLink

			if _, found := categoriesLinks[categoryLink]; !found {
				fmt.Println("Category found: " + categoryLink)
				categoryCollector.Visit(categoryLink)
				categoriesLinks[categoryLink] = true
			} else {
				categoriesLinks[categoryLink] = false
			}

		})
	})

	categoryCollector.OnHTML(".pagination", func(e *colly.HTMLElement) {
		//fmt.Printf("%v\n", e.ChildText("li"))

		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			categoryLink := el.Attr("href")
			if len(categoryLink) < 1 {
				return
			}

			categoryLink = currentSource.BaseURL + categoryLink

			if _, found := categoriesPagesLinks[categoryLink]; !found {
				fmt.Println("Category Page# Found: " + categoryLink)
				productCollector.Visit(categoryLink)
				categoriesPagesLinks[categoryLink] = true
			} else {
				categoriesPagesLinks[categoryLink] = false
			}
		})

	})
	//body > div.container.general
	//body > div.container.general > div > div > div.recomendadosrow.row
	productCollector.OnHTML("div.recomendadosrow.row", func(e *colly.HTMLElement) {
		//fmt.Printf("%v\n", e.ChildText(".product"))
		e.ForEach(".product", func(_ int, el *colly.HTMLElement) {
			productLink := e.ChildAttr("a", "href")
			if len(productLink) < 1 {
				return
			}
			productLink = currentSource.BaseURL + productLink

			fmt.Println("Product found: " + productLink)
			if _, found := productsLinks[productLink]; !found {
				detailCollector.Visit(productLink)
				productsLinks[productLink] = true
			} else {
				productsLinks[productLink] = false
			}
		})
	})

	detailCollector.OnHTML("body.detalle", func(e *colly.HTMLElement) {
		title := e.ChildText("#detalle > div.detalle-rpincipal h1.product-title")
		description := e.ChildText("#detalle > div.detalle-rpincipal div.details-description")
		price := e.ChildText("#precio")
		thumbnail := e.ChildAttr("#img_prod img", "src")
		publishedAt := time.Now()
		categoryName := e.ChildText("div.container.general ul.breadcrumb > a")

		product := &models.Product{
			Title:        title,
			Description:  description,
			Source:       currentSource.ID,
			URL:          e.Request.URL.String(),
			Price:        price,
			CategoryName: categoryName,
			Thumbnail:    thumbnail,
			PublishedAt:  publishedAt,
		}
		storeProduct(product)
		totalProductsCollected++
	})

	q.AddURL(currentSource.BaseURL)

	// Consume
	q.Run(frontPageCollector)
	log.Println("Finished " + currentSource.Name)
	log.Println("Total Products Collected: ", totalProductsCollected)
	w.Done()
}
