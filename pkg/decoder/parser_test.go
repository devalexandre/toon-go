package decoder

import (
	"strings"
	"testing"
)

func TestParseRootTabularArray(t *testing.T) {
	// Test case for root-level tabular array (slice without name)
	input := `[3]{Name,Age,Email,Active}:
  John Doe,30,john.doe@example.com,true
  Alice Smith,28,alice.smith@example.com,true
  Bob Johnson,35,bob.johnson@example.com,false`

	reader := strings.NewReader(input)
	parser := NewParser(reader)

	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check that result is a slice
	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("Expected slice, got %T", result)
	}

	if len(slice) != 3 {
		t.Fatalf("Expected 3 elements, got %d", len(slice))
	}

	// Check first element
	first := slice[0].(map[string]interface{})
	if first["Name"] != "John Doe" {
		t.Errorf("Expected Name=John Doe, got %v", first["Name"])
	}
	if first["Age"] != int64(30) {
		t.Errorf("Expected Age=30, got %v", first["Age"])
	}
	if first["Email"] != "john.doe@example.com" {
		t.Errorf("Expected Email=john.doe@example.com, got %v", first["Email"])
	}
	if first["Active"] != true {
		t.Errorf("Expected Active=true, got %v", first["Active"])
	}

	// Check second element
	second := slice[1].(map[string]interface{})
	if second["Name"] != "Alice Smith" {
		t.Errorf("Expected Name=Alice Smith, got %v", second["Name"])
	}
	if second["Age"] != int64(28) {
		t.Errorf("Expected Age=28, got %v", second["Age"])
	}
	if second["Email"] != "alice.smith@example.com" {
		t.Errorf("Expected Email=alice.smith@example.com, got %v", second["Email"])
	}
	if second["Active"] != true {
		t.Errorf("Expected Active=true, got %v", second["Active"])
	}

	// Check third element
	third := slice[2].(map[string]interface{})
	if third["Name"] != "Bob Johnson" {
		t.Errorf("Expected Name=Bob Johnson, got %v", third["Name"])
	}
	if third["Age"] != int64(35) {
		t.Errorf("Expected Age=35, got %v", third["Age"])
	}
	if third["Email"] != "bob.johnson@example.com" {
		t.Errorf("Expected Email=bob.johnson@example.com, got %v", third["Email"])
	}
	if third["Active"] != false {
		t.Errorf("Expected Active=false, got %v", third["Active"])
	}
}

func TestParseNamedTabularArray(t *testing.T) {
	// Test case for named tabular array (regular object with array field)
	input := `users[2]{Name,Age}:
  John,30
  Alice,25`

	reader := strings.NewReader(input)
	parser := NewParser(reader)

	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check that result is a map
	obj, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map, got %T", result)
	}

	users := obj["users"]
	if users == nil {
		t.Fatal("Expected 'users' field")
	}

	// Check that users is a slice
	slice, ok := users.([]interface{})
	if !ok {
		t.Fatalf("Expected slice, got %T", users)
	}

	if len(slice) != 2 {
		t.Fatalf("Expected 2 elements, got %d", len(slice))
	}

	// Check elements
	first := slice[0].(map[string]interface{})
	if first["Name"] != "John" {
		t.Errorf("Expected Name=John, got %v", first["Name"])
	}
	if first["Age"] != int64(30) {
		t.Errorf("Expected Age=30, got %v", first["Age"])
	}

	second := slice[1].(map[string]interface{})
	if second["Name"] != "Alice" {
		t.Errorf("Expected Name=Alice, got %v", second["Name"])
	}
	if second["Age"] != int64(25) {
		t.Errorf("Expected Age=25, got %v", second["Age"])
	}
}
