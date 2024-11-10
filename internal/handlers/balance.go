package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/models"
)

func BalanceHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDContextKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodPost {
		amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
		if err != nil {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}

		// Get user from Redis
		userJSON, err := database.RedisClient.Get(context.Background(), fmt.Sprintf("user:%s", userID)).Result()
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		var user models.User
		if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
			http.Error(w, "Invalid user data", http.StatusInternalServerError)
			return
		}

		// Update balance
		user.Balance += amount
		user.UpdatedAt = time.Now()

		// Save updated user to Redis
		updatedUserJSON, _ := json.Marshal(user)
		err = database.RedisClient.Set(context.Background(), fmt.Sprintf("user:%s", userID), updatedUserJSON, 0).Err()
		if err != nil {
			http.Error(w, "Failed to update balance", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/balance", http.StatusSeeOther)
		return
	}

	// Get user from Redis
	userJSON, err := database.RedisClient.Get(context.Background(), fmt.Sprintf("user:%s", userID)).Result()
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var user models.User
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		http.Error(w, "Invalid user data", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/balance.gohtml")
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, user.Balance)
	if err != nil {
		log.Printf("Failed to execute template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
