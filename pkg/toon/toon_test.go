package toon

import (
	"strings"
	"testing"
)

func TestMarshalSimpleTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "string",
			input:    "hello",
			expected: `"hello"`,
		},
		{
			name:     "integer",
			input:    42,
			expected: `42`,
		},
		{
			name:     "float",
			input:    3.14,
			expected: `3.14`,
		},
		{
			name:     "boolean",
			input:    true,
			expected: `true`,
		},
		{
			name:     "null",
			input:    nil,
			expected: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.input)
			if err != nil {
				t.Errorf("Marshal() error = %v", err)
				return
			}
			resultStr := string(result)
			if resultStr != tt.expected {
				t.Errorf("Marshal() = %v, want %v", resultStr, tt.expected)
			}
		})
	}
}

func TestMarshalMap(t *testing.T) {
	input := map[string]interface{}{
		"name":  "John",
		"age":   30,
		"alive": true,
	}

	result, err := Marshal(input)
	if err != nil {
		t.Errorf("Marshal() error = %v", err)
		return
	}

	resultStr := string(result)
	if !strings.Contains(resultStr, `name: "John"`) {
		t.Errorf("Marshal() missing name field, got: %s", resultStr)
	}
	if !strings.Contains(resultStr, `age: 30`) {
		t.Errorf("Marshal() missing age field, got: %s", resultStr)
	}
	if !strings.Contains(resultStr, `alive: true`) {
		t.Errorf("Marshal() missing alive field, got: %s", resultStr)
	}
}

func TestMarshalStruct(t *testing.T) {
	type Person struct {
		Name string `toon:"name"`
		Age  int    `toon:"age"`
	}

	input := Person{
		Name: "Alice",
		Age:  25,
	}

	result, err := Marshal(input)
	if err != nil {
		t.Errorf("Marshal() error = %v", err)
		return
	}

	resultStr := string(result)
	if !strings.Contains(resultStr, `name: "Alice"`) {
		t.Errorf("Marshal() missing name field, got: %s", resultStr)
	}
	if !strings.Contains(resultStr, `age: 25`) {
		t.Errorf("Marshal() missing age field, got: %s", resultStr)
	}
}

func TestUnmarshalSimple(t *testing.T) {
	toonData := []byte(`name: "John"
age: 30
active: true`)

	var result map[string]interface{}
	err := Unmarshal(toonData, &result)
	if err != nil {
		t.Errorf("Unmarshal() error = %v", err)
		return
	}

	if result["name"] != "John" {
		t.Errorf("Unmarshal() name = %v, want %v", result["name"], "John")
	}
	if result["age"] != int64(30) {
		t.Errorf("Unmarshal() age = %v, want %v", result["age"], 30)
	}
	if result["active"] != true {
		t.Errorf("Unmarshal() active = %v, want %v", result["active"], true)
	}
}
