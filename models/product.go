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
	ProductsSources = make(map[int]*ProductSource)
	ProductsSources[1] = &ProductSource{
		ID:             1,
		Name:           "Importadora Ronson",
		BaseURL:        "https://importadoraronson.com/",
		AllowedDomains: "importadoraronson.com",
	}
	ProductsSources[2] = &ProductSource{
		ID:             2,
		Name:           "Geeker",
		BaseURL:        "https://www.geeker.com.ar/",
		AllowedDomains: "www.geeker.com.ar",
	}
	ProductsSources[3] = &ProductSource{
		ID:             3,
		Name:           "Distriland",
		BaseURL:        "https://www.distriland.com.ar/",
		AllowedDomains: "www.distriland.com.ar",
	}
	ProductsSources[4] = &ProductSource{
		ID:             4,
		Name:           "La Web Del Celular",
		BaseURL:        "https://www.lawebdelcelular.com.ar/",
		AllowedDomains: "www.lawebdelcelular.com.ar",
	}
}
