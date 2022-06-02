package crawlers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/fixwa/go-prices-tracker/models"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"strings"
	"sync"
	"time"
)

func CrawlDistriland(w *sync.WaitGroup) {
	productsLinks := map[string]bool{}
	categoriesLinks := map[string]bool{}
	//categoriesPagesLinks := map[string]bool{}
	var totalProductsCollected int = 0

	currentSource := models.ProductsSources[3]
	fmt.Printf("%v\n", currentSource)
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

	c.OnHTML(".nav-item", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			categoryLink := el.Attr("href")
			if strings.Index(categoryLink, currentSource.AllowedDomains) == -1 {
				return
			}

			if _, found := categoriesLinks[categoryLink]; !found {
				fmt.Println("Category found: " + categoryLink)
				productCollector.Visit(categoryLink)
				categoriesLinks[categoryLink] = true
			} else {
				categoriesLinks[categoryLink] = false
			}
		})
	})

	productCollector.OnHTML("body.template-category", func(e *colly.HTMLElement) {
		e.ForEach("div.item-product", func(_ int, el *colly.HTMLElement) {
			productLink := el.ChildAttr("a", "href")
			if strings.Index(productLink, "/productos/") == -1 {
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
	})

	detailCollector.OnHTML("body.template-product", InspectAndStore)
	totalProductsCollected++ //@fixme
	q.AddURL(currentSource.BaseURL)

	// Consume
	q.Run(c)

	log.Printf("\x1b[%dm%s %s\x1b[0m", 31, currentSource.Name, "Finished!")
	log.Println("Total Products Collected: ", totalProductsCollected)
	w.Done()
}

func InspectAndStore(e *colly.HTMLElement) {
	title := e.ChildText("h1#product-name")
	price := e.ChildText(".price-container h4")
	thumbnail := e.ChildAttr("img.product-slider-image", "data-srcset")

	if strings.Index(thumbnail, "//") == 0 {
		thumbnail = "https:" + thumbnail
	}
	publishedAt := time.Now()
	categoryName := ""
	e.DOM.Find("div.breadcrumbs .crumb").Each(func(i int, s *goquery.Selection) {
		categoryName = categoryName + "-" + s.Text()
	})
	e.ChildText("div.breadcrumbs .crumb:nth-child(3)")
	description := e.ChildText(".product-description")
	product := &models.Product{
		Title:        title,
		Description:  description,
		Source:       3, //currentSource.ID,
		URL:          e.Request.URL.String(),
		Price:        price,
		CategoryName: categoryName,
		Thumbnail:    thumbnail,
		PublishedAt:  publishedAt,
	}
	fmt.Printf("%v\n", product)
	//	storeProduct(product)
}
