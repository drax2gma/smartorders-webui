package models

import (
	"time"
)

type Order struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	Products   []Product `json:"products" gorm:"many2many:order_products;"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
