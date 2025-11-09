package main

import (
	"fmt"
	"log"

	"github.com/devalexandre/toon-go/pkg/toon"
)

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Email   string `json:"email"`
	Active  bool   `json:"active"`
	Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	} `json:"address"`
}

type Company struct {
	Name      string   `json:"name"`
	Employees []Person `json:"employees"`
	Founded   int      `json:"founded"`
}

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

	fmt.Println("3. Arrays:")
	arrayData := map[string]interface{}{
		"numbers": []int{1, 2, 3, 4, 5},
		"strings": []string{"apple", "banana", "cherry"},
		"mixed":   []interface{}{1, "hello", true, 3.14},
	}

	toonData3, err := toon.Marshal(arrayData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(toonData3))
	fmt.Println()

	fmt.Println("4. Tabular arrays (TOON's main feature):")

	employees := []map[string]interface{}{
		{
			"id":     1,
			"name":   "Alice Smith",
			"role":   "Developer",
			"salary": 75000,
			"active": true,
		},
		{
			"id":     2,
			"name":   "Bob Johnson",
			"role":   "Designer",
			"salary": 65000,
			"active": true,
		},
		{
			"id":     3,
			"name":   "Carol Brown",
			"role":   "Manager",
			"salary": 85000,
			"active": false,
		},
	}

	tabularData := map[string]interface{}{
		"company":   "Tech Corp",
		"employees": employees,
		"metadata": map[string]interface{}{
			"total":       len(employees),
			"departments": []string{"Engineering", "Design", "Management"},
		},
	}

	toonData4, err := toon.Marshal(tabularData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(toonData4))
	fmt.Println()

	fmt.Println("5. Struct encoding:")
	person := Person{
		Name:   "John Doe",
		Age:    30,
		Email:  "john@example.com",
		Active: true,
	}
	person.Address.Street = "123 Main St"
	person.Address.City = "New York"

	toonData5, err := toon.Marshal(person)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(toonData5))
	fmt.Println()

	fmt.Println("6. Comparison with JSON:")
	fmt.Println("TOON format (more readable and token-efficient):")
	fmt.Println(string(toonData5))

	fmt.Println("\nJSON equivalent (more verbose):")
	fmt.Println(`{
  "name": "John Doe",
  "age": 30,
  "email": "john@example.com",
  "active": true,
  "address": {
    "street": "123 Main St",
    "city": "New York"
  }
}`)

	fmt.Println("\n7. Parsing TOON data:")
	toonString := `name: "John Doe"
age: 30
active: true
address:
  street: "123 Main St"
  city: "New York"`

	fmt.Println("Parsing TOON string:")
	fmt.Println(toonString)

	fmt.Println("(Parsing implementation in progress...)")
}
