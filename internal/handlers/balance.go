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

func BalanceHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(uint)

	if r.Method == http.MethodPost {
		amount, err := strconv.ParseInt(r.FormValue("amount"), 32, 32)
		if err != nil {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}

		// Get user from Redis
		userJSON, err := database.RedisClient.Get(context.Background(), "user:"+strconv.FormatUint(uint64(userID), 10)).Result()
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
		// Note: In a real application, you'd want to store balance separately and use atomic operations
		user.Balance += amount

		// Save updated user to Redis
		updatedUserJSON, _ := json.Marshal(user)
		err = database.RedisClient.Set(context.Background(), "user:"+strconv.FormatUint(uint64(userID), 10), updatedUserJSON, 0).Err()
		if err != nil {
			http.Error(w, "Failed to update balance", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/balance", http.StatusSeeOther)
		return
	}

	// Get user from Redis
	userJSON, err := database.RedisClient.Get(context.Background(), "user:"+strconv.FormatUint(uint64(userID), 10)).Result()
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var user models.User
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		http.Error(w, "Invalid user data", http.StatusInternalServerError)
		return
	}

	tmpl, _ := template.ParseFiles("web/templates/balance.gohtml")
	tmpl.Execute(w, user.Balance)
}
