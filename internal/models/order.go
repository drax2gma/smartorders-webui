package models

import (
	"time"

	"github.com/drax2gma/smartorders-webui/internal/utils"
)

type Order struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	ProductID  string    `json:"product_id"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func GenerateOrderID(userID, productID string, timestamp time.Time) string {
	input := userID + "|" + productID + "|" + timestamp.String()
	return utils.GenerateXXH64Hash(input)
}
