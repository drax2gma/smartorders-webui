package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/models"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDContextKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user's order IDs from Redis
	orderIDs, err := database.RedisClient.LRange(context.Background(), fmt.Sprintf("user:%s:orders", userID), 0, -1).Result()
	if err != nil {
		log.Printf("Failed to get order IDs: %v", err)
		http.Error(w, "Failed to get orders", http.StatusInternalServerError)
		return
	}

	var orders []models.Order
	for _, orderID := range orderIDs {
		orderJSON, err := database.RedisClient.Get(context.Background(), fmt.Sprintf("order:%s", orderID)).Result()
		if err != nil {
			log.Printf("Failed to get order %s: %v", orderID, err)
			continue
		}

		var order models.Order
		if err := json.Unmarshal([]byte(orderJSON), &order); err != nil {
			log.Printf("Failed to unmarshal order %s: %v", orderID, err)
			continue
		}

		orders = append(orders, order)
	}

	tmpl, err := template.ParseFiles("web/templates/status.gohtml")
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, orders)
	if err != nil {
		log.Printf("Failed to execute template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
