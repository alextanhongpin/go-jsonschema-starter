package main

import (
	_ "embed"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed schemas/product-nested.json
var productSchema []byte

// https://json-schema.org/learn/getting-started-step-by-step.html
func main() {
	schema, err := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(productSchema))
	if err != nil {
		panic(err)
	}

	validate := func(doc []byte) {
		result, err := schema.Validate(gojsonschema.NewBytesLoader(doc))
		if err != nil {
			panic(err)
		}

		if result.Valid() {
			println("The document is valid")
		} else {
			println("The document is not valid. see errors:")
			for _, desc := range result.Errors() {
				fmt.Println("-", desc)
			}
		}
	}

	emptyDoc := []byte(`{}`)
	validate(emptyDoc)

	badDoc := []byte(`
		{
			"id": "1",
			"price": 0
		}
	`)
	validate(badDoc)

	badPriceDoc := []byte(`
		{
			"productId": 1,
			"productName": "apple",
			"price": 0
		}
	`)
	validate(badPriceDoc)

	goodDoc := []byte(`
		{
			"productId": 1,
			"productName": "apple",
			"price": 1
		}
	`)
	validate(goodDoc)

	goodDocBadTags := []byte(`
		{
			"productId": 1,
			"productName": "apple",
			"price": 1,
			"tags": ["fruit", "red", "red"]
		}
	`)
	validate(goodDocBadTags)

	goodDocWithDimension := []byte(`
		{
			"productId": 1,
			"productName": "apple",
			"price": 1,
			"dimensions": {
				"length": 1,
				"width": 1,
				"height": "0"
			}
		}
	`)
	validate(goodDocWithDimension)
}
