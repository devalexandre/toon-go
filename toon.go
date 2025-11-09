package main

import (
	"github.com/devalexandre/toon-go/pkg/toon"
)

func Marshal(v interface{}) ([]byte, error) {
	return toon.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return toon.Unmarshal(data, v)
}
