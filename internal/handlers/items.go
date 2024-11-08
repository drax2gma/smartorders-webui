package handlers

import (
	"smartorders-webui/internal/database"

	"github.com/labstack/echo/v4"
)

func ListItems(c echo.Context) error {
	items, err := database.GetItems()
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to fetch items"})
	}

	return c.Render(200, "items.html", map[string]interface{}{
		"items": items,
	})
}
