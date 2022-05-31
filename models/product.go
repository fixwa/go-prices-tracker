package models

import (
	"time"
)

type Product struct {
	ID           int `gorm:"primary_key;"`
	Title        string
	Content      string
	Source       int
	URL          string
	Price        string
	CategoryName string
	Thumbnail    string
	PublishedAt  time.Time
	CreatedAt    time.Time
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
	ProductsSources = make(map[int]*ProductSource)
	ProductsSources[1] = &s1
}
