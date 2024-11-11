package models

import (
	"github.com/drax2gma/smartorders-webui/internal/utils"
)

type Product struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Params      string  `json:"params"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func GenerateProductID(description, params string) string {
	input := description + "|" + params
	return utils.GenerateXXH64Hash(input)
}
