package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/models"
)

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDContextKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodPost {
		// Handle order creation
		handleOrderCreation(w, r, userID)
		return
	}

	// Get products from Redis
	productsJSON, err := database.RedisClient.Get(context.Background(), "products").Result()
	if err != nil {
		log.Printf("Failed to get products: %v", err)
		http.Error(w, "Failed to get products", http.StatusInternalServerError)
		return
	}

	var products []models.Product
	if err := json.Unmarshal([]byte(productsJSON), &products); err != nil {
		log.Printf("Failed to unmarshal products: %v", err)
		http.Error(w, "Invalid products data", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/order.gohtml")
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, products)
	if err != nil {
		log.Printf("Failed to execute template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func handleOrderCreation(w http.ResponseWriter, r *http.Request, userID string) {
	productID, err := strconv.ParseUint(r.FormValue("product_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Get product from Redis
	productJSON, err := database.RedisClient.Get(context.Background(), fmt.Sprintf("product:%d", productID)).Result()
	if err != nil {
		log.Printf("Failed to get product: %v", err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	var product models.Product
	if err := json.Unmarshal([]byte(productJSON), &product); err != nil {
		log.Printf("Failed to unmarshal product: %v", err)
		http.Error(w, "Invalid product data", http.StatusInternalServerError)
		return
	}

	// Create order
	order := models.Order{
		UserID:     userID,
		ProductID:  productID,
		TotalPrice: product.Price,
		Status:     "pending",
	}

	// Save order to Redis
	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Printf("Failed to marshal order: %v", err)
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	err = database.RedisClient.Set(context.Background(), fmt.Sprintf("order:%s", order.ID), orderJSON, 0).Err()
	if err != nil {
		log.Printf("Failed to save order: %v", err)
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	// Add order to user's order list
	err = database.RedisClient.RPush(context.Background(), fmt.Sprintf("user:%s:orders", userID), order.ID).Err()
	if err != nil {
		log.Printf("Failed to add order to user's list: %v", err)
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/status", http.StatusSeeOther)
}
