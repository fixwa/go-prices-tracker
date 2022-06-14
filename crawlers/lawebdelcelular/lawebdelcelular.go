package lawebdelcelular

import (
	"fmt"
	"github.com/fixwa/go-prices-tracker/crawlers"
	"github.com/fixwa/go-prices-tracker/models"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"strings"
	"sync"
	"time"
)

var (
	totalProductsCollected int
	currentSource          *models.ProductSource
)

func Crawl(w *sync.WaitGroup) {
	currentSource = models.ProductsSources[4]
	existingProducts := crawlers.GetProductsBySource(currentSource)
	productsLinks := map[string]bool{}
	for _, product := range existingProducts {
		productsLinks[product.URL] = true
	}

	categoriesLinks := map[string]bool{}

	//fmt.Printf("%v\n", currentSource)
	log.Println("Crawling " + currentSource.Name)

	categoriesCollector := colly.NewCollector(
		colly.AllowedDomains(currentSource.AllowedDomains),
	)

	q, _ := queue.New(
		1,
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	subCategoriesCollector := categoriesCollector.Clone()
	productsCollector := categoriesCollector.Clone()
	productDetailsCollector := categoriesCollector.Clone()

	categoriesCollector.OnHTML("#new", func(e *colly.HTMLElement) {
		//fmt.Printf("%v\n", e.ChildText("li"))
		e.ForEach("a.agile-icon", func(_ int, el *colly.HTMLElement) {
			categoryLink := el.Attr("href")
			if len(categoryLink) < 1 {
				return
			}

			categoryLink = currentSource.BaseURL + "inicio/" + categoryLink

			if _, found := categoriesLinks[categoryLink]; !found {
				fmt.Println("Category found: " + categoryLink)
				subCategoriesCollector.Visit(categoryLink)
				categoriesLinks[categoryLink] = true
			} else {
				categoriesLinks[categoryLink] = false
			}

		})
	})

	subCategoriesCollector.OnHTML("ul.box-subcat", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			subCategoryLink := el.Attr("href")
			if len(subCategoryLink) < 1 {
				return
			}

			subCategoryLink = currentSource.BaseURL + "inicio/" + subCategoryLink
			subCategoryLink = strings.ReplaceAll(subCategoryLink, "/./", "/")

			if _, found := categoriesLinks[subCategoryLink]; !found {
				fmt.Println("Sub category found: " + subCategoryLink)
				productsCollector.Visit(subCategoryLink)
				categoriesLinks[subCategoryLink] = true
			} else {
				categoriesLinks[subCategoryLink] = false
			}
		})
	})

	productsCollector.OnHTML("table.table", func(e *colly.HTMLElement) {
		e.ForEach("tr td > a", func(_ int, el *colly.HTMLElement) {
			productLink := e.ChildAttr("a", "href")
			if len(productLink) < 1 {
				return
			}
			productLink = currentSource.BaseURL + "inicio/" + productLink
			productLink = strings.ReplaceAll(productLink, "/./", "/")

			if _, found := productsLinks[productLink]; !found {
				productDetailsCollector.Visit(productLink)
				productsLinks[productLink] = true
				fmt.Println("New Product found: " + productLink)
			} else {
				productsLinks[productLink] = false
			}
		})
	})

	productDetailsCollector.OnHTML("div.container:first-child", inspectAndStore)

	q.AddURL(currentSource.BaseURL + "inicio/index.php")
	q.Run(categoriesCollector)
	//
	//q.AddURL("https://www.lawebdelcelular.com.ar/inicio/productos.php?cat=18&sub_cat=222&codigo=3850020102")
	//q.Run(productDetailsCollector)

	log.Printf("\x1b[%dm%s %s\x1b[0m", 31, currentSource.Name, "Finished!")
	log.Println("Total Products Collected: ", totalProductsCollected)
	w.Done()
}

func Clear() {
	crawlers.DeleteAllBySource(currentSource)
}

func StandardizeSpaces(s string) string {
	return strings.Trim(strings.Join(strings.Fields(s), " "), " ")
}

func inspectAndStore(e *colly.HTMLElement) {
	title := e.ChildText("h5.title-w3")
	if title == "" {
		return
	}
	price := ""
	description := ""
	e.ForEach("li.item1", func(_ int, el *colly.HTMLElement) {
		text := StandardizeSpaces(el.Text)
		fmt.Println(text)

		if strings.Contains(text, "Precio:") {
			price = text
		}

		if strings.Contains(text, "Descripcion:") {
			description = text
		}
	})
	thumbnail := e.ChildAttr(".zoom-img:first-child", "src")
	publishedAt := time.Now()
	categoryName := e.ChildText("h3.hdg")

	product := &models.Product{
		Title:        title,
		Description:  description,
		Source:       currentSource.ID,
		URL:          e.Request.URL.String(),
		Price:        price,
		CategoryName: categoryName,
		Thumbnail:    currentSource.BaseURL + "images/" + thumbnail,
		PublishedAt:  publishedAt,
	}
	_, err := crawlers.StoreProduct(product)

	if err == nil {
		totalProductsCollected++
	}
	//log.Printf("\x1b[%dm Stored product: %s\x1b[0m", 32, insert.InsertedID.(primitive.ObjectID))
	//fmt.Printf("%v\n", product)
}
