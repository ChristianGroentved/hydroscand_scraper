package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

// initializing a data structure to keep the scraped data
type Product struct {
	URL        string
	Name       string
	Attributes map[string]string
	Variants   map[string]string
}

type Hierachy struct {
	Category map[string]string
	URL      string
}

func main() {

	// initializing the slice of structs to store the data to scrape
	var products []Product
	var hierarchies []Hierachy

	c := colly.NewCollector(colly.AllowedDomains("www.hydroscand.dk"))

	// Create another collector to scrape product details
	detailCollector := c.Clone()

	// scraping logic
	c.OnHTML("li.sub-categories__item", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		e.Request.Visit(link)
	})

	c.OnHTML("div.page-wrapper", func(e *colly.HTMLElement) {

		hierachy := Hierachy{}
		hierachy.Category = make(map[string]string)
		levels := 0

		e.ForEach("div.breadcrumbs li[class*=cat]", func(_ int, el *colly.HTMLElement) {

			hierachy.Category[fmt.Sprintf("Category_%v", levels)] = strings.TrimSpace(el.Text)
			levels++

		})

		c.OnHTML("a.product-item-link", func(e *colly.HTMLElement) {

			detailsLink := e.Attr("href")

			hierachy.URL = detailsLink
			hierarchies = append(hierarchies, hierachy)
			detailCollector.Visit(detailsLink)
		})

	})

	detailCollector.OnHTML("div.product-info-wrapper", func(e *colly.HTMLElement) {
		log.Println("Product description found", e.Request.URL)

		product := Product{URL: e.Request.URL.String()}
		product.Name = e.ChildText("h1.page-title > span.base")
		product.Attributes = make(map[string]string)
		product.Variants = make(map[string]string)

		// Extract the details of the product
		e.ForEach("table.data.table.additional-attributes:not(#product-attribute-specs-table) tr", func(_ int, el *colly.HTMLElement) {
			key := el.ChildText("th")
			val := el.ChildText("td")

			product.Attributes[key] = val

		})

		e.ForEach("div.product-variants table.table td", func(_ int, el *colly.HTMLElement) {
			key := el.Attr("data-th")
			value := strings.TrimSpace(el.Text)
			product.Variants[key] = value
		})

		products = append(products, product)

	})
	c.Visit("https://www.hydroscand.dk/dk_dk/produkter/adaptere/bsp-koblinger-og-adaptere")

	for _, h := range hierarchies {
		fmt.Printf("Hierachy: %+v\n", h)
	}

	for _, p := range products {
		fmt.Printf("Product: %+v\n", p)
	}
}
