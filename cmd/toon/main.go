package main

import (
	"fmt"
	"log"

	"github.com/devalexandre/toon-go/pkg/toon"
)

func main() {
	fmt.Println("=== TOON Format Example ===")
	fmt.Println()

	fmt.Println("1. Simple key-value pairs:")
	simpleData := map[string]interface{}{
		"name":    "John Doe",
		"age":     30,
		"active":  true,
		"balance": 1234.56,
		"tags":    []string{"developer", "golang", "toon"},
	}

	toonData1, err := toon.Marshal(simpleData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(toonData1))
	fmt.Println()

	fmt.Println("2. Nested objects:")
	nestedData := map[string]interface{}{
		"user": map[string]interface{}{
			"id":    123,
			"login": "johndoe",
			"profile": map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
				"settings": map[string]interface{}{
					"theme":    "dark",
					"language": "en",
				},
			},
		},
		"timestamp": "2024-01-15T10:30:00Z",
	}

	toonData2, err := toon.Marshal(nestedData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(toonData2))
	fmt.Println()

	fmt.Println("3. Parsing TOON data:")
	toonString := `name: "John Doe"
age: 30
active: true
address:
  street: "123 Main St"
  city: "New York"`

	fmt.Println("Parsing TOON string:")
	fmt.Println(toonString)
	fmt.Println()

	var parsedData map[string]interface{}
	err = toon.Unmarshal([]byte(toonString), &parsedData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Parsed data: %+v\n", parsedData)
}
