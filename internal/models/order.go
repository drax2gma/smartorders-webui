package models

import (
	"time"
)

type Order struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	ProductID  uint      `json:"product_id"` // Only store the product ID
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
