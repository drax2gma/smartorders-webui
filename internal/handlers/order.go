package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/models"
	"github.com/labstack/echo/v4"
)

func OrderHandler(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	if c.Request().Method == http.MethodPost {
		return handleOrderSubmission(c, userID)
	}

	// Fetch products
	rows, err := database.DB.Query("SELECT * FROM products")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching products")
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Megnevezes, &p.Parameterek, &p.Price, &p.Stock)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error scanning products")
		}
		products = append(products, p)
	}

	tmpl := template.Must(template.ParseFiles(
		"web/templates/layout.gohtml",
		"web/templates/order.gohtml",
	))

	data := TemplateData{
		Title: "Rendelés",
		Data: map[string]interface{}{
			"Products": products,
		},
	}

	return tmpl.Execute(c.Response().Writer, data)
}

func handleOrderSubmission(c echo.Context, userID string) error {
	productID := c.FormValue("product_id")

	// Get product details
	var product models.Product
	err := database.DB.QueryRow("SELECT * FROM products WHERE id = ?", productID).Scan(
		&product.ID, &product.Megnevezes, &product.Parameterek, &product.Price, &product.Stock,
	)
	if err != nil {
		return c.HTML(http.StatusBadRequest, `
            <div class="alert alert-danger">
                A kiválasztott termék nem található.
            </div>
        `)
	}

	// Create order
	orderID := models.GenerateOrderID(userID, productID, time.Now())
	_, err = database.DB.Exec(`
        INSERT INTO orders (id, user_id, product_id, total_price, status, created_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `, orderID, userID, productID, product.Price, "pending", time.Now())

	if err != nil {
		return c.HTML(http.StatusInternalServerError, `
            <div class="alert alert-danger">
                Hiba történt a rendelés feldolgozása során.
            </div>
        `)
	}

	return c.HTML(http.StatusOK, `
        <div class="alert alert-success">
            A rendelés sikeresen leadva!
        </div>
    `)
}
