package handlers

import (
	db "smartorders-webui/internal/db"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := db.GetCustomer(username)
	if err != nil || !user.CheckPassword(password) {
		return c.JSON(401, map[string]string{"error": "Invalid credentials"})
	}

	// Set session or JWT token here

	return c.Redirect(302, "/home")
}

func ListItems(c echo.Context) error {
	items, err := db.GetOrders()
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to fetch items"})
	}

	return c.Render(200, "items.html", map[string]interface{}{
		"items": items,
	})
}
