package main

import (
	"encoding/json"
	"os"
)

func main() {

	categories, attributes, products := scrape()

	jsonBytes, err := json.Marshal(categories)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("categories.json", jsonBytes, 0644)
	if err != nil {
		panic(err)
	}

	jsonBytes, err = json.Marshal(attributes)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("attributes.json", jsonBytes, 0644)
	if err != nil {
		panic(err)
	}

	jsonBytes, err = json.Marshal(products)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("products.json", jsonBytes, 0644)
	if err != nil {
		panic(err)
	}

}
