package handlers

import (
	"net/http"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/models"
	"github.com/drax2gma/smartorders-webui/web/templates"
	"github.com/labstack/echo/v4"
)

func OrderHandler(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	if c.Request().Method == http.MethodPost {
		return handleOrderCreation(c)
	}

	// Fetch products from database
	rows, err := database.DB.Query("SELECT * FROM products")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching products")
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Megnevezes, &p.Parameterek, &p.Price, &p.Stock)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error scanning products")
		}
		products = append(products, p)
	}

	return templates.Order(products).Render(c.Request().Context(), c.Response().Writer)
}

func handleOrderCreation(c echo.Context) error {
	userID := c.Get("user_id").(string)
	productID := c.FormValue("product_id")

	// Fetch product details
	var product models.Product
	err := database.DB.QueryRow("SELECT * FROM products WHERE id = ?", productID).Scan(
		&product.ID, &product.Megnevezes, &product.Parameterek, &product.Price, &product.Stock,
	)
	if err != nil {
		return c.String(http.StatusNotFound, "Product not found")
	}

	// Create order
	order := models.Order{
		ID:         models.GenerateOrderID(userID, productID),
		UserID:     userID,
		ProductID:  productID,
		TotalPrice: product.Price,
		Status:     "pending",
	}

	_, err = database.DB.Exec(`
        INSERT INTO orders (id, user_id, product_id, total_price, status, created_at)
        VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
    `, order.ID, order.UserID, order.ProductID, order.TotalPrice, order.Status)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error creating order")
	}

	return c.HTML(http.StatusOK, "<p>Rendelés sikeresen leadva!</p>")
}
