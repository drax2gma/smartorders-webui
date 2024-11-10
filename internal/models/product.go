package models

import (
	"fmt"
	"hash/fnv"
)

type Product struct {
	ID          string  `json:"id"`
	Megnevezes  string  `json:"megnevezes"`
	Parameterek string  `json:"parameterek"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func GenerateProductID(megnevezes, parameterek string) string {
	return xxh32(megnevezes + "|" + parameterek)
}

// xxh32 implementálja az xxh32 hash algoritmust
func xxh32(input string) string {
	h := fnv.New32a()
	h.Write([]byte(input))
	return fmt.Sprintf("%08x", h.Sum32())
}
