package main

import (
	"smartorders-webui/internal/common"
	"smartorders-webui/internal/db"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize Redis
	db.InitRedis()

	// Initialize i18n
	common.InitTranslator()

	// Routes
	e.Static("/static", "webui/static")
	e.File("/", "webui/templates/login.html")
	e.POST("/login", handlers.Login)
	e.GET("/home", handlers.Home)
	e.GET("/items", handlers.ListItems)
	e.POST("/cart/add", handlers.AddToCart)
	e.GET("/cart", handlers.ViewCart)
	e.POST("/order", handlers.PlaceOrder)
	e.GET("/order/:id", handlers.OrderStatus)

	e.Logger.Fatal(e.Start(":8080"))
}
