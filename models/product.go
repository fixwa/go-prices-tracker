package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Product struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title        string             `json:"title" bson:"title,omitempty"`
	Description  string             `json:"description" bson:"description,omitempty"`
	Source       int                `json:"source" bson:"source,omitempty"`
	URL          string             `json:"url" bson:"url,omitempty"`
	Price        string             `json:"price" bson:"price,omitempty"`
	CategoryName string             `json:"categoryName" bson:"categoryName,omitempty"`
	Thumbnail    string             `json:"thumbnail" bson:"thumbnail,omitempty"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt,omitempty"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
}

type ProductSource struct {
	ID             int
	Name           string
	BaseURL        string
	AllowedDomains string
}

var ProductsSources map[int]*ProductSource

func init() {
	s1 := ProductSource{
		ID:             1,
		Name:           "Importadora Ronson",
		BaseURL:        "https://importadoraronson.com/",
		AllowedDomains: "importadoraronson.com",
	}

	s2 := ProductSource{
		ID:             2,
		Name:           "Geeker",
		BaseURL:        "https://www.geeker.com.ar/",
		AllowedDomains: "www.geeker.com.ar",
	}
	ProductsSources = make(map[int]*ProductSource)
	ProductsSources[1] = &s1
	ProductsSources[2] = &s2
}
