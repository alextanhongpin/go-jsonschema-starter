package main

import (
	_ "embed"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed schemas/product.json
var productSchema []byte

// https://json-schema.org/learn/getting-started-step-by-step.html
func main() {
	sl := gojsonschema.NewSchemaLoader()
	sl.AddSchemas(gojsonschema.NewBytesLoader(productSchema))

	mainSchema := gojsonschema.NewStringLoader(`{
		"$id" : "https://example.com/main.json",
		"allOf" : [
			{ "$ref" : "https://example.com/product.schema.json" }
		]
	}`)

	schema, err := sl.Compile(mainSchema)
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
}
