package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

// initializing a data structure to keep the scraped data
type Product struct {
	URL         string
	Category    string
	SubCategory string
	ProductType string
	Name        string
	Attributes  map[string]string
	Variants    map[string]string
}

func main() {

	// initializing the slice of structs to store the data to scrape
	var products []Product

	c := colly.NewCollector(colly.AllowedDomains("www.hydroscand.dk"))

	// Create another collector to scrape product details
	detailCollector := c.Clone()

	// scraping logic
	c.OnHTML("li.sub-categories__item", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		e.Request.Visit(link)
	})

	c.OnHTML("div.breadcrumbs", func(e *colly.HTMLElement) {
		e.ForEach("li[class*=cat]", func(_ int, el *colly.HTMLElement) {
			fmt.Println(strings.TrimSpace(el.Text))
		})

	})

	c.OnHTML("a.product-item-link", func(e *colly.HTMLElement) {
		detailsLink := e.Attr("href")
		detailCollector.Visit(detailsLink)
	})

	detailCollector.OnHTML("div.product-info-wrapper", func(e *colly.HTMLElement) {
		log.Println("Product description found", e.Request.URL)

		product := Product{URL: e.Request.URL.String()}
		product.Name = e.ChildText("h1.page-title > span.base")
		product.Attributes = make(map[string]string)
		product.Variants = make(map[string]string)

		// fmt.Println(e.ChildText("h1.page-title > span.base"))

		// Extract the details of the product
		e.ForEach("table.data.table.additional-attributes:not(#product-attribute-specs-table) tr", func(_ int, el *colly.HTMLElement) {
			key := el.ChildText("th")
			val := el.ChildText("td")
			// fmt.Println(el.ChildText("th"), ":", el.ChildText("td"))

			product.Attributes[key] = val

		})

		e.ForEach("div.product-variants table.table td", func(_ int, el *colly.HTMLElement) {
			key := el.Attr("data-th")
			value := strings.TrimSpace(el.Text)
			// fmt.Println(el.Attr("data-th"), ":", strings.TrimSpace(el.Text))
			product.Variants[key] = value
		})

		products = append(products, product)

	})
	c.Visit("https://www.hydroscand.dk/dk_dk/produkter")

	for _, p := range products {
		fmt.Printf("Product: %+v\n", p)
	}
}
