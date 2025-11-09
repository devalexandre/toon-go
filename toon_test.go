package main

import (
	"reflect"
	"testing"
)

func TestMarshalUnmarshal(t *testing.T) {
	type Person struct {
		Name   string `toon:"name"`
		Age    int    `toon:"age"`
		Active bool   `toon:"active"`
	}

	person := Person{
		Name:   "John Doe",
		Age:    30,
		Active: true,
	}

	data, err := Marshal(person)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	t.Logf("Marshaled TOON:\n%s", string(data))

	var unmarshaled Person
	err = Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(person, unmarshaled) {
		t.Errorf("Expected %+v, got %+v", person, unmarshaled)
	}
}

func TestComplexStruct(t *testing.T) {
	type Address struct {
		Street string `toon:"street"`
		City   string `toon:"city"`
	}

	type User struct {
		Name    string  `toon:"name"`
		Age     int     `toon:"age"`
		Active  bool    `toon:"active"`
		Address Address `toon:"address"`
	}

	user := User{
		Name:   "John Doe",
		Age:    30,
		Active: true,
		Address: Address{
			Street: "123 Main St",
			City:   "New York",
		},
	}

	data, err := Marshal(user)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	t.Logf("Marshaled TOON:\n%s", string(data))

	var unmarshaled User
	err = Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(user, unmarshaled) {
		t.Errorf("Expected %+v, got %+v", user, unmarshaled)
	}
}

func TestParseSimpleTOON(t *testing.T) {
	toonData := `name: "John Doe"
age: 30
active: true`

	var result map[string]interface{}
	err := Unmarshal([]byte(toonData), &result)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	expected := map[string]interface{}{
		"name":   "John Doe",
		"age":    int64(30),
		"active": true,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestParseNestedTOON(t *testing.T) {
	toonData := `name: "John Doe"
age: 30
active: true
address:
  street: "123 Main St"
  city: "New York"`

	var result map[string]interface{}
	err := Unmarshal([]byte(toonData), &result)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	expected := map[string]interface{}{
		"name":   "John Doe",
		"age":    int64(30),
		"active": true,
		"address": map[string]interface{}{
			"street": "123 Main St",
			"city":   "New York",
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}
