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

	// Routes
	e.GET("/", handlers.SessionMiddleware(handlers.HomeHandler))
	e.GET("/login", handlers.LoginPageHandler)
	e.POST("/login", handlers.LoginHandler)
	e.GET("/logout", handlers.LogoutHandler)
	e.GET("/order", handlers.OrderHandler)
	e.POST("/order", handlers.OrderHandler)
	e.GET("/status", handlers.StatusHandler)
	e.GET("/balance", handlers.BalanceHandler)
	e.POST("/balance", handlers.BalanceHandler)
	e.GET("/message", handlers.MessageHandler)
	e.POST("/message", handlers.MessageHandler)
	e.POST("/validate-email", handlers.ValidateEmailHandler)

	// Serve static files
	e.Static("/static", "web/static")

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
