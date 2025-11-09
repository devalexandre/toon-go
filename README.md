# TOON - Token Optimized Object Notation (Go Implementation)

TOON is a data serialization format designed as a more token-efficient and human-readable alternative to JSON. It's particularly well-suited for AI applications, APIs, and scenarios where reducing token usage is important.

For the complete TOON specification, see: https://github.com/toon-format/spec/blob/main/SPEC.md

## Features

- **Token Efficient**: Uses fewer tokens than JSON for equivalent data
- **Human Readable**: Clean, intuitive syntax
- **Tabular Arrays**: Optimized format for array-of-objects data
- **Backward Compatible**: Can represent the same data as JSON
- **Go Integration**: Native Go support with encoding/decoding

## Installation

```bash
go get github.com/devalexandre/toon-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/devalexandre/toon-go"
)

func main() {
    data := map[string]interface{}{
        "name": "John Doe",
        "age":  30,
        "active": true,
        "tags": []string{"developer", "golang"},
    }
    
    // Marshal to TOON format
    toonData, err := toon.Marshal(data)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(string(toonData))
    // Output:
    // name: "John Doe"
    // age: 30
    // active: true
    // tags[2]: "developer","golang"
}
```

## TOON Syntax

### Basic Key-Value Pairs

```
name: "John Doe"
age: 30
active: true
balance: 1234.56
```

### Nested Objects

```
user:
  id: 123
  name: "John Doe"
  profile:
    email: "john@example.com"
    settings:
      theme: "dark"
      notifications: true
```

### Arrays

Regular arrays:
```
numbers[5]: 1,2,3,4,5
strings[3]: "apple","banana","cherry"
```

### Tabular Arrays (TOON's Key Feature)

TOON excels at representing arrays of objects with consistent fields:

```
employees[3]{id,name,role,salary,active}:
  1,"Alice Smith","Developer",75000,true
  2,"Bob Johnson","Designer",65000,true
  3,"Carol Brown","Manager",85000,false
```

This tabular format is much more token-efficient than JSON's verbose array-of-objects representation.

## API

### Marshal

```go
func Marshal(v interface{}) ([]byte, error)
```

Marshals a Go value to TOON format.

### MarshalIndent

```go
func MarshalIndent(v interface{}, indent string) ([]byte, error)
```

Marshals a Go value to TOON format with custom indentation.

### Unmarshal

```go
func Unmarshal(data []byte, v interface{}) error
```

Unmarshals TOON data into a Go value.

## Examples

See `example/toon_example.go` for comprehensive usage examples.

## Benefits Over JSON

1. **Token Efficiency**: Up to 40% fewer tokens for typical data
2. **Better Readability**: Cleaner syntax without excessive brackets
3. **Optimized Arrays**: Tabular format for array-of-objects data
4. **AI-Friendly**: Designed with LLM tokenization in mind
5. **Backward Compatible**: Can represent all JSON data

## Use Cases

- **AI/LLM Applications**: Reduce token usage in prompts and responses
- **APIs**: More efficient data transmission
- **Configuration Files**: Cleaner, more readable configs
- **Data Logging**: Compact log data representation
- **Microservices**: Efficient inter-service communication

## License

MIT

## Contributing

Contributions welcome! Please see CONTRIBUTING.md for details.
