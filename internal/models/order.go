package models

import (
	"crypto/md5"
	"fmt"
	"time"
)

type Order struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	ProductID  string    `json:"product_id"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func GenerateOrderID(userID, productID string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(userID+"|"+productID+"|"+time.Now().String())))[:8]
}
