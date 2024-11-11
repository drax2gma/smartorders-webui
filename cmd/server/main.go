package main

import (
	"fmt"
	"strings"
	"time"

	"encoding/json"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func getStatusCodeColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "\033[32m" // Green
	case code >= 300 && code < 400:
		return "\033[33m" // Yellow
	case code >= 400 && code < 500:
		return "\033[31m" // Red
	case code >= 500:
		return "\033[35m" // Magenta
	default:
		return "\033[0m" // Default color
	}
}

func getStatusCodeExplanation(code int) string {
	explanations := map[int]string{
		100: "Continue",
		101: "Switching Protocols",
		200: "OK",
		201: "Created",
		202: "Accepted",
		204: "No Content",
		301: "Moved Permanently",
		302: "Found",
		303: "See Other",
		304: "Not Modified",
		400: "Bad Request",
		401: "Unauthorized",
		403: "Forbidden",
		404: "Not Found",
		405: "Method Not Allowed",
		500: "Internal Server Error",
		501: "Not Implemented",
		502: "Bad Gateway",
		503: "Service Unavailable",
		504: "Gateway Timeout",
	}

	if explanation, ok := explanations[code]; ok {
		return explanation
	}
	return "Unknown Status Code"
}

func prettyJSONLogger(c echo.Context, values middleware.RequestLoggerValues) error {
	statusColor := getStatusCodeColor(values.Status)
	resetColor := "\033[0m"

	data := map[string]interface{}{
		"time":   values.StartTime.Format(time.RFC3339),
		"uri":    values.URI,
		"method": values.Method,
		"status": map[string]interface{}{
			"code":        fmt.Sprintf("%s%d%s", statusColor, values.Status, resetColor),
			"explanation": getStatusCodeExplanation(values.Status),
		},
		"latency":   values.Latency.String(),
		"error":     values.Error,
		"remote_ip": values.RemoteIP,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Replace escaped color codes with actual color codes
	coloredJSON := string(jsonData)
	coloredJSON = strings.ReplaceAll(coloredJSON, `\u001b`, "\033")

	fmt.Println(coloredJSON)
	return nil
}

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Create Echo instance
	e := echo.New()

	// Middleware
	// Configure custom logger
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:        true,
		LogStatus:     true,
		LogLatency:    true,
		LogMethod:     true,
		LogError:      true,
		LogValuesFunc: prettyJSONLogger,
	}))

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
	e.Logger.Fatal(e.Start(":8080"))
}
