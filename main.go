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

	fmt.Println("Schema:", userSchema)
	fmt.Println("Document:", string(doc))

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
			fmt.Println("Type:", desc.Type())
			fmt.Println("Field:", desc.Field())
			fmt.Println("Description:", desc.Description())
			fmt.Println("DescriptionFormat:", desc.DescriptionFormat())
			fmt.Println("Value:", desc.Value())
			fmt.Println("Details:", desc.Details())
			fmt.Println()
		}
	}
}
