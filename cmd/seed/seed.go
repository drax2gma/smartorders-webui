package main

import (
	"fmt"
	"log"
	"time"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	seedUsers()
	seedProducts()

	log.Println("Seeding completed successfully!")
}

func seedUsers() {
	for i := 1; i <= 20; i++ {
		email := fmt.Sprintf("user%d@example.com", i)
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("password%d", i)), bcrypt.DefaultCost)

		user := models.User{
			ID:        models.GenerateUserID(email),
			Name:      fmt.Sprintf("Felhasználó %d", i),
			Email:     email,
			Password:  string(hashedPassword),
			Balance:   0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err := database.DB.Exec(`
            INSERT INTO users (id, name, email, password, balance, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?, ?)
        `, user.ID, user.Name, user.Email, user.Password, user.Balance, user.CreatedAt, user.UpdatedAt)

		if err != nil {
			log.Printf("Error seeding user %s: %v\n", user.Email, err)
		}
	}
}

func seedProducts() {
	products := []models.Product{
		{Megnevezes: "Laptop", Parameterek: "16GB RAM, 512GB SSD", Price: 999.99, Stock: 10},
		{Megnevezes: "Smartphone", Parameterek: "64GB Storage, Black", Price: 499.99, Stock: 20},
		{Megnevezes: "Headphones", Parameterek: "Wireless, Noise Cancelling", Price: 99.99, Stock: 50},
		{Megnevezes: "Tablet", Parameterek: "10 inch, Wi-Fi", Price: 299.99, Stock: 15},
		{Megnevezes: "Smartwatch", Parameterek: "Heart Rate Monitor, GPS", Price: 199.99, Stock: 30},
	}

	for _, product := range products {
		product.ID = models.GenerateProductID(product.Megnevezes, product.Parameterek)
		_, err := database.DB.Exec(`
            INSERT INTO products (id, megnevezes, parameterek, price, stock)
            VALUES (?, ?, ?, ?, ?)
        `, product.ID, product.Megnevezes, product.Parameterek, product.Price, product.Stock)

		if err != nil {
			log.Printf("Error seeding product %s: %v\n", product.Megnevezes, err)
		}
	}
}
