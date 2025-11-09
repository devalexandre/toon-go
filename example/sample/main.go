package main

import (
	"fmt"
	"log"

	"github.com/devalexandre/toon-go/pkg/toon"
)

type Person struct {
	Name   string
	Age    int
	Email  string
	Active bool
}

func main() {
	// Struct simples
	fmt.Println("=== TOON Format Example ===")
	fmt.Println()

	// Exemplo: slice de structs
	people := []Person{
		{Name: "John Doe", Age: 30, Email: "john.doe@example.com", Active: true},
		{Name: "Alice Smith", Age: 28, Email: "alice.smith@example.com", Active: true},
		{Name: "Bob Johnson", Age: 35, Email: "bob.johnson@example.com", Active: false},
	}

	fmt.Println("# Tabular array (slice of struct)")
	str, err := toon.Encode(people)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(str)

	// Exemplo: struct simples
	fmt.Println("\n# Simple object (struct)")
	person := Person{Name: "Alice", Age: 30, Email: "alice@example.com", Active: true}
	str, err = toon.Encode(person)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(str)
}
