package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	_ "embed"

	"github.com/xeipuuv/gojsonschema"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed schemas/payment.json
var paymentSchema []byte

//https://json-schema.org/understanding-json-schema/reference/conditionals.html
func main() {
	schema, err := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(paymentSchema))
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

	c := jsonschema.NewCompiler()
	if err := c.AddResource("payment.json", bytes.NewReader(paymentSchema)); err != nil {
		panic(err)
	}
	sch, err := c.Compile("payment.json")
	if err != nil {
		log.Fatalf("%#v", err)
	}

	docs := []string{
		`
		{
			"name": "John Doe",
			"credit_card": 5555555555555555,
			"billing_address": "555 Debtor's Lane"
		}`,
		`{
			"name": "John Doe",
			"credit_card": 5555555555555555
		}`, // This should fail, but returns true ...
		`
		{
			"name": "John Doe"
		}`,
		`
		{
			"name": "John Doe",
			"billing_address": "555 Debtor's Lane"
		}`,
	}
	for _, doc := range docs {
		validate([]byte(doc))

		var v interface{}
		if err := json.Unmarshal([]byte(doc), &v); err != nil {
			log.Fatal(err)
		}
		if err := sch.Validate(v); err != nil {
			log.Printf("%#v\n", err)

			var verr *jsonschema.ValidationError
			if errors.As(err, &verr) {
				fmt.Println("message:", verr.Message)
				fmt.Println("error:", verr.Error())
				fmt.Println("detailed output:", verr.DetailedOutput())

				b, _ := json.MarshalIndent(verr.FlagOutput(), "", "  ")
				fmt.Println("flag output:", string(b))

				b, _ = json.MarshalIndent(verr.BasicOutput(), "", "  ")
				fmt.Println("basic output:", string(b))
			}
		}
	}
}
