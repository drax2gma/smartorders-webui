package main

import (
	"log"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Auth routes (no session required)
	e.GET("/login", handlers.LoginPageHandler)
	e.POST("/login", handlers.LoginHandler)
	e.POST("/validate-email", handlers.ValidateEmailHandler)

	// Protected routes (session required)
	protected := e.Group("")
	protected.Use(handlers.SessionMiddleware)

	protected.GET("/", handlers.HomeHandler)
	protected.GET("/logout", handlers.LogoutHandler)

	protected.GET("/order", handlers.OrderHandler)
	protected.POST("/order", handlers.OrderHandler)

	protected.GET("/status", handlers.StatusHandler)

	protected.GET("/balance", handlers.BalanceHandler)
	protected.POST("/balance", handlers.BalanceHandler)

	protected.GET("/message", handlers.MessageHandler)
	protected.POST("/message", handlers.MessageHandler)

	// Serve static files
	e.Static("/static", "web/static")

	// Start server
	log.Println("Server starting on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
