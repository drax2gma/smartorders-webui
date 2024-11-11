package models

import (
	"crypto/md5"
	"fmt"
)

type Product struct {
	ID          string  `json:"id"`
	Megnevezes  string  `json:"megnevezes"`
	Parameterek string  `json:"parameterek"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func GenerateProductID(megnevezes, parameterek string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(megnevezes+"|"+parameterek)))[:8]
}
