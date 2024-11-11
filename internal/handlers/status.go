package handlers

import (
	"net/http"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/models"
	"github.com/drax2gma/smartorders-webui/web/templates"
	"github.com/labstack/echo/v4"
)

func StatusHandler(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	rows, err := database.DB.Query("SELECT * FROM orders WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching orders")
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		err := rows.Scan(&o.ID, &o.UserID, &o.ProductID, &o.TotalPrice, &o.Status, &o.CreatedAt)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error scanning orders")
		}
		orders = append(orders, o)
	}

	return templates.Status(orders).Render(c.Request().Context(), c.Response().Writer)
}
