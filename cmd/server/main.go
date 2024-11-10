// cmd/server/main.go
package main

import (
	"log"
	"net/http"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/handlers"
)

func main() {
	// Redis inicializálása
	if err := database.InitRedis("localhost:6379"); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	if err := database.InitializeProducts(); err != nil {
		log.Fatalf("Failed to initialize products: %v", err)
	}

	// Set up routes
	http.HandleFunc("/", handlers.SessionMiddleware(handlers.HomeHandler))
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/order", handlers.SessionMiddleware(handlers.OrderHandler))
	http.HandleFunc("/status", handlers.SessionMiddleware(handlers.StatusHandler))
	http.HandleFunc("/balance", handlers.SessionMiddleware(handlers.BalanceHandler))
	http.HandleFunc("/message", handlers.SessionMiddleware(handlers.MessageHandler))
	http.HandleFunc("/validate-email", handlers.ValidateEmailHandler)

	// Serve static files
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start the server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
