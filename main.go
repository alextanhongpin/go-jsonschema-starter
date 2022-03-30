package main

import (
	_ "embed"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed schemas/user.json
var userSchema string

func main() {
	doc := []byte(`{"foo": "bar", "age": "10"}`)
	schema := gojsonschema.NewStringLoader(userSchema)

	fmt.Println("Schema:")
	fmt.Println(userSchema)

	fmt.Println("Document:")
	fmt.Println(string(doc))

	result, err := gojsonschema.Validate(schema, gojsonschema.NewBytesLoader(doc))
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
