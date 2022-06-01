package crawlers

import (
	"context"
	"fmt"
	"github.com/fixwa/go-prices-tracker/database"
	"github.com/fixwa/go-prices-tracker/models"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		price := e.ChildText(".product-info > div.price-wrapper > p.price")
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
		storeProduct(product)
	})

	q.AddURL(currentSource.BaseURL)

	// Consume
	q.Run(c)
	log.Println("Finished " + currentSource.Name)
	w.Done()
}

func storeProduct(product *models.Product) {
	productsCollection := database.Db.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	//p := models.Product{
	//	//ID:           primitive.ObjectID{},
	//	Title:        "The Title",
	//	Content:      "The content",
	//	Source:       0,
	//	URL:          "http://example.com",
	//	Price:        "100.00",
	//	CategoryName: "Testing",
	//	Thumbnail:    "http://example.com/image.png",
	//	PublishedAt:  time.Time{},
	//	CreatedAt:    time.Time{},
	//}

	result, err := productsCollection.InsertOne(ctx, bson.M{
		"title":        product.Title,
		"content":      product.Content,
		"source":       product.Source,
		"url":          product.URL,
		"price":        product.Price,
		"categoryName": product.CategoryName,
		"thumbnail":    product.Thumbnail,
		"createdAt":    time.Now(),
		"publishedAt":  product.PublishedAt,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Stored product: ", result.InsertedID.(primitive.ObjectID))
}
