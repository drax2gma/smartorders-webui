package handlers

import (
	"html/template"
	"net/http"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/models"
	"github.com/labstack/echo/v4"
)

func StatusHandler(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	// Fetch user's orders
	rows, err := database.DB.Query(`
        SELECT o.*, p.description, p.params 
        FROM orders o 
        JOIN products p ON o.product_id = p.id 
        WHERE o.user_id = ? 
        ORDER BY o.created_at DESC`, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching orders")
	}
	defer rows.Close()

	type OrderWithProduct struct {
		models.Order
		ProductName   string
		ProductParams string
	}

	var orders []OrderWithProduct
	for rows.Next() {
		var o OrderWithProduct
		err := rows.Scan(
			&o.ID, &o.UserID, &o.ProductID, &o.TotalPrice, &o.Status, &o.CreatedAt,
			&o.ProductName, &o.ProductParams,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error scanning orders")
		}
		orders = append(orders, o)
	}

	// Get user data for the template
	var user models.User
	err = database.DB.QueryRow("SELECT * FROM users WHERE id = ?", userID).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Balance, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching user data")
	}

	tmpl := template.Must(template.ParseFiles(
		"web/templates/layout.gohtml",
		"web/templates/status.gohtml",
	))

	data := TemplateData{
		Title: "Rendelések állapota",
		User:  &user,
		Data:  orders,
	}

	return tmpl.Execute(c.Response().Writer, data)
}
