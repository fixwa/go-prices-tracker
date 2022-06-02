package crawlers

import (
	"context"
	"fmt"
	"github.com/fixwa/go-prices-tracker/database"
	"github.com/fixwa/go-prices-tracker/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func storeProduct(product *models.Product) {
	productsCollection := database.Db.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	//p := models.Product{
	//	//ID:           primitive.ObjectID{},
	//	Title:        "The Title",
	//	Description:      "The content",
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
		"description":  product.Description,
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
