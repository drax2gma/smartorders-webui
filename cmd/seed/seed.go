package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/drax2gma/smartorders-webui/internal/models"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

var redisClient *redis.Client

func main() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Töröljük a meglévő adatokat
	clearExistingData()

	// Új adatok beszúrása
	seedProducts()
	seedUsers()

	fmt.Println("Seeding completed successfully!")
}

func clearExistingData() {
	ctx := context.Background()

	// Töröljük az összes terméket
	redisClient.Del(ctx, "products")

	// Töröljük az összes felhasználót
	keys, _ := redisClient.Keys(ctx, "user:*").Result()
	if len(keys) > 0 {
		redisClient.Del(ctx, keys...)
	}

	fmt.Println("Existing data cleared.")
}

func seedProducts() {
	candies := []string{"Csokoládé", "Gumicukor", "Nyalóka", "Karamella", "Zselés cukor"}
	clothes := []string{"Póló", "Nadrág", "Kabát", "Sapka", "Zokni"}

	for i := 0; i < 80; i++ {
		var product models.Product
		if i < 40 {
			product = models.Product{
				ID:          uint(i + 1),
				Name:        candies[rand.Intn(len(candies))],
				Description: "Finom édesség",
				Price:       float64(rand.Intn(1000) + 100),
				Stock:       rand.Intn(100) + 50,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
		} else {
			product = models.Product{
				ID:          uint(i + 1),
				Name:        clothes[rand.Intn(len(clothes))],
				Description: "Divatos ruhadarab",
				Price:       float64(rand.Intn(10000) + 1000),
				Stock:       rand.Intn(50) + 10,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
		}

		jsonProduct, _ := json.Marshal(product)
		redisClient.HSet(context.Background(), "products", fmt.Sprintf("%d", product.ID), jsonProduct)
	}
}

func seedUsers() {
	ctx := context.Background()

	for i := 0; i < 20; i++ {
		email := fmt.Sprintf("user%d@example.com", i+1)
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("password%d", i+1)), bcrypt.DefaultCost)

		user := models.User{
			ID:        uint(i + 1),
			Name:      fmt.Sprintf("Felhasználó %d", i+1),
			Email:     email,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		jsonUser, _ := json.Marshal(user)
		err := redisClient.Set(ctx, "user:"+email, jsonUser, 0).Err()
		if err != nil {
			fmt.Printf("Error seeding user %s: %v\n", email, err)
		}

		// Külön tároljuk a jelszó hash-t
		err = redisClient.Set(ctx, "user:"+email+":password", string(hashedPassword), 0).Err()
		if err != nil {
			fmt.Printf("Error seeding password for user %s: %v\n", email, err)
		}
	}
}
