// internal/handlers/message.go
package handlers

import (
	"context"
	"html/template"
	"net/http"

	"github.com/drax2gma/smartorders-webui/internal/database"
)

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	// ... (implement session check)

	if r.Method == http.MethodPost {
		messageType := r.FormValue("message_type")

		// Save message to Redis
		err := database.RedisClient.RPush(context.Background(), "admin:messages", messageType).Err()
		if err != nil {
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/message", http.StatusSeeOther)
		return
	}

	tmpl, _ := template.ParseFiles("web/templates/message.gohtml")
	tmpl.Execute(w, nil)
}
