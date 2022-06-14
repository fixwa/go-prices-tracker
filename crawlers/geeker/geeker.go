package geeker

import (
	"fmt"
	"github.com/fixwa/go-prices-tracker/crawlers"
	"github.com/fixwa/go-prices-tracker/models"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"sync"
	"time"
)

var (
	totalProductsCollected int
	currentSource          *models.ProductSource
)

func Crawl(w *sync.WaitGroup) {
	currentSource = models.ProductsSources[2]
	existingProducts := crawlers.GetProductsBySource(currentSource)
	productsLinks := map[string]bool{}
	for _, product := range existingProducts {
		productsLinks[product.URL] = true
	}

	categoriesLinks := map[string]bool{}
	categoriesPagesLinks := map[string]bool{}

	//fmt.Printf("%v\n", currentSource)
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
		//fmt.Printf("%v\n", e.ChildText("li"))
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			categoryLink := el.Attr("href")
			if len(categoryLink) < 1 {
				return
			}

			categoryLink = currentSource.BaseURL + categoryLink

			if _, found := categoriesLinks[categoryLink]; !found {
				//fmt.Println("Category found: " + categoryLink)
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
				//fmt.Println("Category Page# Found: " + categoryLink)
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

			if _, found := productsLinks[productLink]; !found {
				detailCollector.Visit(productLink)
				productsLinks[productLink] = true
				fmt.Println("New Product found: " + productLink)
			} else {
				productsLinks[productLink] = false
			}
		})
	})

	detailCollector.OnHTML("body.detalle", inspectAndStore)

	q.AddURL(currentSource.BaseURL)

	// Consume
	q.Run(c)
	log.Printf("\x1b[%dm%s %s\x1b[0m", 31, currentSource.Name, "Finished!")
	log.Println("Total Products Collected: ", totalProductsCollected)
	w.Done()
}

func Clear() {
	crawlers.DeleteAllBySource(currentSource)
}

func inspectAndStore(e *colly.HTMLElement) {
	title := e.ChildText("#detalle > div.detalle-rpincipal h1.product-title")
	description := e.ChildText("#detalle > div.detalle-rpincipal div.details-description")
	price := e.ChildText("#precio")
	thumbnail := e.ChildAttr("#img_prod a", "href")
	publishedAt := time.Now()
	categoryName := e.ChildText("div.container.general ul.breadcrumb > a")

	product := &models.Product{
		Title:        title,
		Description:  description,
		Source:       currentSource.ID,
		URL:          e.Request.URL.String(),
		Price:        price,
		CategoryName: categoryName,
		Thumbnail:    currentSource.BaseURL + thumbnail,
		PublishedAt:  publishedAt,
	}
	crawlers.StoreProduct(product)
	totalProductsCollected++
}
