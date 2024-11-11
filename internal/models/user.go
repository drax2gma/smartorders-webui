package models

import (
	"time"

	"github.com/drax2gma/smartorders-webui/internal/utils"
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
	return utils.GenerateXXH64Hash(email)
}
