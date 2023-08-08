package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

// initializing a data structure to keep the scraped data
type Product struct {
	url         string
	category    string
	subCategory string
	productType string
	name        string
}

type Category struct {
	category    string
	subCategory string
	productType string
}

func main() {

	// initializing the slice of structs to store the data to scrape
	var Products []Product

	c := colly.NewCollector(colly.AllowedDomains("www.hydroscand.dk"))

	// Create another collector to scrape product details
	detailCollector := c.Clone()

	// scraping logic
	c.OnHTML("li.sub-categories__item", func(e *colly.HTMLElement) {

		fmt.Println(e.ChildAttr("a", "title"))
		link := e.ChildAttr("a", "href")
		e.Request.Visit(link)
	})

	c.OnHTML("a.product-item-link", func(e *colly.HTMLElement) {
		detailsLink := e.Attr("href")
		detailCollector.Visit(detailsLink)
	})

	detailCollector.OnHTML("div.product-info-wrapper", func(e *colly.HTMLElement) {
		log.Println("Product description found", e.Request.URL)

		fmt.Println(e.ChildText("h1.page-title > span.base"))

		// Extract the details of the product
		e.ForEach("table.data.table.additional-attributes:not(#product-attribute-specs-table) tr", func(_ int, el *colly.HTMLElement) {
			fmt.Println(el.ChildText("th"), ":", el.ChildText("td"))

		})

		e.ForEach("div.product-variants table.table td", func(_ int, el *colly.HTMLElement) {
			fmt.Println(el.Attr("data-th"), ":", strings.TrimSpace(el.Text))
		})

	})
	c.Visit("https://www.hydroscand.dk/dk_dk/produkter")

	fmt.Printf("+%v", Products)
}
