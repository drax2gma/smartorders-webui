package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := database.InitRedis("localhost:6379"); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	seedUsers()
	seedProducts()

	fmt.Println("Seeding completed successfully!")
}

func seedUsers() {
	ctx := context.Background()

	for i := 0; i < 20; i++ {
		email := fmt.Sprintf("user%d@example.com", i+1)
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("password%d", i+1)), bcrypt.DefaultCost)

		user := models.User{
			ID:        models.GenerateUserID(email),
			Name:      fmt.Sprintf("Felhasználó %d", i+1),
			Email:     email,
			Password:  string(hashedPassword),
			Balance:   0, // Kezdeti egyenleg
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		jsonUser, _ := json.Marshal(user)
		err := database.RedisClient.Set(ctx, fmt.Sprintf("user:%s", user.ID), jsonUser, 0).Err()
		if err != nil {
			log.Printf("Error seeding user %s: %v\n", user.Email, err)
		}

		// Store email to ID mapping
		err = database.RedisClient.Set(ctx, fmt.Sprintf("email:%s", user.Email), user.ID, 0).Err()
		if err != nil {
			log.Printf("Error storing email mapping for %s: %v\n", user.Email, err)
		}

		// Store password separately
		err = database.RedisClient.Set(ctx, fmt.Sprintf("user:%s:password", user.ID), user.Password, 0).Err()
		if err != nil {
			log.Printf("Error storing password for user %s: %v\n", user.Email, err)
		}
	}
}

func seedProducts() {
	ctx := context.Background()
	products := []models.Product{
		{ID: 1, Name: "Laptop", Price: 999.99, Stock: 10},
		{ID: 2, Name: "Smartphone", Price: 499.99, Stock: 20},
		{ID: 3, Name: "Headphones", Price: 99.99, Stock: 50},
		{ID: 4, Name: "Tablet", Price: 299.99, Stock: 15},
		{ID: 5, Name: "Smartwatch", Price: 199.99, Stock: 30},
	}

	for _, product := range products {
		jsonProduct, err := json.Marshal(product)
		if err != nil {
			log.Printf("Error marshaling product %s: %v\n", product.Name, err)
			continue
		}

		err = database.RedisClient.Set(ctx, fmt.Sprintf("product:%d", product.ID), jsonProduct, 0).Err()
		if err != nil {
			log.Printf("Error seeding product %s: %v\n", product.Name, err)
		}
	}

	// Store all products in a single key
	allProductsJSON, err := json.Marshal(products)
	if err != nil {
		log.Printf("Error marshaling all products: %v\n", err)
	} else {
		err = database.RedisClient.Set(ctx, "products", allProductsJSON, 0).Err()
		if err != nil {
			log.Printf("Error storing all products: %v\n", err)
		}
	}
}
