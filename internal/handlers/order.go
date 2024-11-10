package handlers

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/models"
)

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(uint)

	if r.Method == http.MethodPost {
		productID, _ := strconv.ParseUint(r.FormValue("product_id"), 10, 32)

		// Get product from Redis
		productJSON, err := database.RedisClient.Get(context.Background(), "product:"+strconv.FormatUint(productID, 10)).Result()
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		var product models.Product
		if err := json.Unmarshal([]byte(productJSON), &product); err != nil {
			http.Error(w, "Invalid product data", http.StatusInternalServerError)
			return
		}

		// Create order
		order := models.Order{
			UserID:     userID,
			Products:   []models.Product{product},
			TotalPrice: product.Price,
			Status:     "pending",
		}

		// Save order to Redis
		orderJSON, _ := json.Marshal(order)
		err = database.RedisClient.Set(context.Background(), "order:"+strconv.FormatUint(uint64(order.ID), 10), orderJSON, 0).Err()
		if err != nil {
			http.Error(w, "Failed to create order", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/status", http.StatusSeeOther)
		return
	}

	// Get products from Redis
	productsJSON, err := database.RedisClient.Get(context.Background(), "products").Result()
	if err != nil {
		http.Error(w, "Failed to get products", http.StatusInternalServerError)
		return
	}

	var products []models.Product
	if err := json.Unmarshal([]byte(productsJSON), &products); err != nil {
		http.Error(w, "Invalid products data", http.StatusInternalServerError)
		return
	}

	tmpl, _ := template.ParseFiles("web/templates/order.gohtml")
	tmpl.Execute(w, products)
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(uint)

	// Get user's orders from Redis
	ordersJSON, err := database.RedisClient.Get(context.Background(), "user:"+strconv.FormatUint(uint64(userID), 10)+":orders").Result()
	if err != nil {
		http.Error(w, "Failed to get orders", http.StatusInternalServerError)
		return
	}

	var orders []models.Order
	if err := json.Unmarshal([]byte(ordersJSON), &orders); err != nil {
		http.Error(w, "Invalid orders data", http.StatusInternalServerError)
		return
	}

	tmpl, _ := template.ParseFiles("web/templates/status.gohtml")
	tmpl.Execute(w, orders)
}
