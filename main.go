package main

import (
	"encoding/json"
	"os"
)

func main() {

	categories, attributes, products := scrape()

	saveJson("categories.json", categories)
	saveJson("attributes.json", attributes)
	saveJson("products.json", products)

}

func saveJson(filename string, v any) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filename, jsonBytes, 0644)
	if err != nil {
		panic(err)
	}
}
