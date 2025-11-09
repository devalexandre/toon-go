package types

import (
	"strconv"
	"strings"
)

type TOONValue interface{}

type TOONObject map[string]TOONValue

type TOONArray []TOONValue

type TabularArray struct {
	Name   string
	Count  int
	Fields []string
	Rows   [][]string
}

type PrimitiveType int

const (
	TypeString PrimitiveType = iota
	TypeNumber
	TypeBoolean
	TypeNull
)

func ParsePrimitive(value string) TOONValue {
	trimmed := strings.TrimSpace(value)
	
	if trimmed == "true" {
		return true
	}
	if trimmed == "false" {
		return false
	}
	
	if trimmed == "null" {
		return nil
	}
	
	if num, err := strconv.ParseFloat(trimmed, 64); err == nil {
		if float64(int64(num)) == num {
			return int64(num)
		}
		return num
	}
	
	if len(trimmed) >= 2 && trimmed[0] == '"' && trimmed[len(trimmed)-1] == '"' {
		return trimmed[1 : len(trimmed)-1]
	}
	
	return trimmed
}

func IsPrimitiveType(value TOONValue) bool {
	switch value.(type) {
	case string, int64, float64, bool, nil:
		return true
	default:
		return false
	}
}

func ShouldUseTabularFormat(slice []TOONValue) bool {
	if len(slice) < 2 {
		return false
	}
	
	var firstFields []string
	for i, item := range slice {
		if obj, ok := item.(TOONObject); ok {
			fields := make([]string, 0, len(obj))
			for key := range obj {
				fields = append(fields, key)
			}
			
			
			if i == 0 {
				firstFields = fields
			} else {
				if len(fields) != len(firstFields) {
					return false
				}
				for _, v := range obj {
					if !IsPrimitiveType(v) {
						return false
					}
				}
			}
		} else {
			return false
		}
	}
	
	return len(firstFields) > 0
}
