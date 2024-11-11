package models

import (
	"crypto/md5"
	"fmt"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GenerateUserID(email string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(email)))[:8]
}
