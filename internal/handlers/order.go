package handlers

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/models"
)

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		log.Println("User ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodPost {
		handleOrderCreation(w, r, userID)
		return
	}

	// Get products from Redis
	products, err := getProducts()
	if err != nil {
		log.Printf("Failed to get products: %v", err)
		http.Error(w, "Failed to get products", http.StatusInternalServerError)
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

func handleOrderCreation(w http.ResponseWriter, r *http.Request, userID uint) {
	productID, err := strconv.ParseUint(r.FormValue("product_id"), 10, 32)
	if err != nil {
		log.Printf("Invalid product ID: %v", err)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Get product from Redis
	product, err := getProduct(productID)
	if err != nil {
		log.Printf("Failed to get product: %v", err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Create order with a single product
	order := models.Order{
		UserID:     userID,
		ProductID:  uint(productID), // Store only the product ID
		TotalPrice: product.Price,
		Status:     "pending",
	}

	// Save order to Redis
	err = saveOrder(order)
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/status", http.StatusSeeOther)
}

func getProducts() ([]models.Product, error) {
	ctx := context.Background()
	productsJSON, err := database.RedisClient.Get(ctx, "products").Result()
	if err != nil {
		return nil, err
	}

	var products []models.Product
	err = json.Unmarshal([]byte(productsJSON), &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func getProduct(productID uint64) (models.Product, error) {
	ctx := context.Background()
	productJSON, err := database.RedisClient.Get(ctx, "product:"+strconv.FormatUint(productID, 10)).Result()
	if err != nil {
		return models.Product{}, err
	}

	var product models.Product
	err = json.Unmarshal([]byte(productJSON), &product)
	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func saveOrder(order models.Order) error {
	ctx := context.Background()
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = database.RedisClient.Set(ctx, "order:"+strconv.FormatUint(uint64(order.ID), 10), orderJSON, 0).Err()
	if err != nil {
		return err
	}

	// Add order to user's order list
	err = database.RedisClient.RPush(ctx, "user:"+strconv.FormatUint(uint64(order.UserID), 10)+":orders", order.ID).Err()
	if err != nil {
		return err
	}

	return nil
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		log.Println("User ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user's orders from Redis
	orders, err := getUserOrders(userID)
	if err != nil {
		log.Printf("Failed to get orders: %v", err)
		http.Error(w, "Failed to get orders", http.StatusInternalServerError)
		return
	}

	// Fetch product details for each order
	var ordersWithProducts []struct {
		Order   models.Order
		Product models.Product
	}

	for _, order := range orders {
		product, err := getProduct(uint64(order.ProductID))
		if err != nil {
			log.Printf("Failed to get product for order %d: %v", order.ID, err)
			continue
		}
		ordersWithProducts = append(ordersWithProducts, struct {
			Order   models.Order
			Product models.Product
		}{order, product})
	}

	tmpl, err := template.ParseFiles("web/templates/status.gohtml")
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, ordersWithProducts)
	if err != nil {
		log.Printf("Failed to execute template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func getUserOrders(userID uint) ([]models.Order, error) {
	ctx := context.Background()
	orderIDs, err := database.RedisClient.LRange(ctx, "user:"+strconv.FormatUint(uint64(userID), 10)+":orders", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var orders []models.Order
	for _, orderID := range orderIDs {
		orderJSON, err := database.RedisClient.Get(ctx, "order:"+orderID).Result()
		if err != nil {
			continue
		}

		var order models.Order
		err = json.Unmarshal([]byte(orderJSON), &order)
		if err != nil {
			continue
		}

		orders = append(orders, order)
	}

	return orders, nil
}
