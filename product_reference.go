package main

import (
	_ "embed"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed schemas/product.schema.json
var productSchema []byte

//go:embed schemas/geographical-location.schema.json
var geographicalLocationSchema []byte

// https://json-schema.org/learn/getting-started-step-by-step.html
func main() {
	sl := gojsonschema.NewSchemaLoader()
	sl.AddSchemas(gojsonschema.NewBytesLoader(geographicalLocationSchema))
	mainSchema := gojsonschema.NewBytesLoader(productSchema)

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

	goodDoc := []byte(`
		{
			"productId": 1,
			"productName": "An ice sculpture",
			"price": 12.50,
			"tags": [ "cold", "ice" ],
			"dimensions": {
				"length": 7.0,
				"width": 12.0,
				"height": 9.5
			},
			"warehouseLocation": {
				"latitude": -78.75,
				"longitude": 20.4
			}
		}
	`)
	validate(goodDoc)

	badDoc := []byte(`
		{
			"productId": 1,
			"productName": "An ice sculpture",
			"price": 12.50,
			"tags": [ "cold", "ice" ],
			"dimensions": {
				"length": 7.0,
				"width": 12.0,
				"height": 9.5
			},
			"warehouseLocation": {
				"latitude": -78.75
			}
		}
	`)
	validate(badDoc)
}
