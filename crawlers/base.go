package crawlers

import (
	"context"
	"fmt"
	"github.com/fixwa/go-prices-tracker/database"
	"github.com/fixwa/go-prices-tracker/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func init() {
	database.ConnectDatabase()
}

func StoreProduct(product *models.Product) (*mongo.InsertOneResult, error) {
	productsCollection := database.Db.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

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

	//log.Printf("\x1b[%dm Stored product: %s\x1b[0m", 32, result.InsertedID.(primitive.ObjectID))
	return result, err
}

func GetProductsBySource(source *models.ProductSource) []models.Product {
	productsCollection := database.Db.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	cursor, err := productsCollection.Find(ctx, bson.M{"source": source.ID})
	if err != nil {
		log.Fatal(err)
	}

	var result []models.Product
	if err = cursor.All(ctx, &result); err != nil {
		log.Fatal(err)
	}

	return result
}

func DeleteAllBySource(source *models.ProductSource) {
	productsCollection := database.Db.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	deleteCollection, err := productsCollection.DeleteMany(ctx, bson.M{"source": source.ID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v products.\n", deleteCollection.DeletedCount)
}

func DeleteAll() {
	productsCollection := database.Db.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	deleteCollection, err := productsCollection.DeleteMany(ctx, bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v products.\n", deleteCollection.DeletedCount)
}
