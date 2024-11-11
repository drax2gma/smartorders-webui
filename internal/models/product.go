package models

import (
	"github.com/drax2gma/smartorders-webui/internal/utils"
)

type Product struct {
	ID          string  `json:"id"`
	Megnevezes  string  `json:"megnevezes"`
	Parameterek string  `json:"parameterek"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func GenerateProductID(megnevezes, parameterek string) string {
	input := megnevezes + "|" + parameterek
	return utils.GenerateXXH64Hash(input)
}
