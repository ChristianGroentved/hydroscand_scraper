package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

// initializing a data structure to keep the scraped data
type Product struct {
	url, name string
}

func main() {

	// initializing the slice of structs to store the data to scrape
	var Products []Product

	c := colly.NewCollector(colly.AllowedDomains("www.hydroscand.dk"))

	// Create another collector to scrape product details
	detailCollector := c.Clone()

	// scraping logic
	c.OnHTML("li.sub-categories__item", func(e *colly.HTMLElement) {
		Product := Product{}
		fmt.Println(e.ChildAttr("a", "title"))
		link := e.ChildAttr("a", "href")
		e.Request.Visit(link)
		Product.name = e.ChildAttr("a", "title")

		Products = append(Products, Product)
	})

	c.Visit("https://www.hydroscand.dk/dk_dk/produkter")

	fmt.Printf("+%v", Products)
}
