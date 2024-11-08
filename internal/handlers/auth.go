package handlers

import (
	"github.com/labstack/echo/v4"

	db "smartorders-webui/internal/database"
)

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := db.GetUser(username)
	if err != nil || !user.CheckPassword(password) {
		return c.JSON(401, map[string]string{"error": "Invalid credentials"})
	}

	// Set session or JWT token here

	return c.Redirect(302, "/home")
}
