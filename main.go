package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

// initializing a data structure to keep the scraped data
type ProductDetails struct {
	URL      string
	Variants map[string]string
}

type ProductAttribute struct {
	URL       string
	Name      string
	Attribute map[string]string
}

type Hierachy struct {
	Category map[string]string
	URL      string
}

func main() {

	// initializing the slice of structs to store the data to scrape
	var products []ProductDetails
	var attributes []ProductAttribute
	var hierarchies []Hierachy

	c := colly.NewCollector(colly.AllowedDomains("www.hydroscand.dk"))

	// Create another collector to scrape product details
	detailCollector := c.Clone()

	// Callback to find the links to the products on the page
	c.OnHTML("li.sub-categories__item", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		e.Request.Visit(link)
	})

	// Callback to find the links to the details page of each product
	// and fetch the category hierarchy of the product
	c.OnHTML("div.page-wrapper", func(e *colly.HTMLElement) {

		hierachy := Hierachy{}
		hierachy.Category = make(map[string]string)
		levels := 0

		// Extract the category hierarchy of the product
		e.ForEach("div.breadcrumbs li[class*=cat]", func(_ int, el *colly.HTMLElement) {

			hierachy.Category[fmt.Sprintf("Category_%v", levels)] = strings.TrimSpace(el.Text)
			levels++

		})
		// Extract the links to the details page of each product
		c.OnHTML("a.product-item-link", func(e *colly.HTMLElement) {

			detailsLink := e.Attr("href")

			hierachy.URL = detailsLink
			hierarchies = append(hierarchies, hierachy)
			detailCollector.Visit(detailsLink)
		})

	})

	// Callback to extract the details of the product
	detailCollector.OnHTML("div.product-info-wrapper", func(e *colly.HTMLElement) {
		log.Println("Product description found", e.Request.URL)

		productAttributes := ProductAttribute{URL: e.Request.URL.String()}
		productAttributes.Name = e.ChildText("h1.page-title > span.base")
		productAttributes.Attribute = make(map[string]string)

		product := ProductDetails{URL: e.Request.URL.String()}
		product.Variants = make(map[string]string)

		// Extract the common attributes of the product
		e.ForEach("table.data.table.additional-attributes:not(#product-attribute-specs-table) tr", func(_ int, el *colly.HTMLElement) {
			key := el.ChildText("th")
			val := el.ChildText("td")

			productAttributes.Attribute[key] = val

		})

		// Extract the variants of the product
		e.ForEach("div.product-variants table.table td", func(_ int, el *colly.HTMLElement) {
			key := el.Attr("data-th")
			value := strings.TrimSpace(el.Text)
			product.Variants[key] = value
		})

		attributes = append(attributes, productAttributes)
		products = append(products, product)

	})
	c.Visit("https://www.hydroscand.dk/dk_dk/produkter")

	for _, h := range hierarchies {
		fmt.Printf("Hierachy: %+v\n", h)
	}

	for _, p := range products {
		fmt.Printf("Product: %+v\n", p)
	}
}
