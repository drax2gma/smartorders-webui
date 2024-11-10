// internal/database/redis.go
package database

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/drax2gma/smartorders-webui/internal/models"
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis(addr string) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Ellenőrizzük a kapcsolatot
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	return nil
}

func InitializeProducts() error {
	ctx := context.Background()

	// Példa termékek
	products := []models.Product{
		{ID: 1, Name: "Laptop", Price: 999.99, Stock: 10},
		{ID: 2, Name: "Smartphone", Price: 499.99, Stock: 20},
		{ID: 3, Name: "Headphones", Price: 99.99, Stock: 50},
	}

	// Termékek mentése egyenként
	for _, product := range products {
		productJSON, err := json.Marshal(product)
		if err != nil {
			return err
		}

		err = RedisClient.Set(ctx, "product:"+strconv.FormatUint(uint64(product.ID), 10), productJSON, 0).Err()
		if err != nil {
			return err
		}
	}

	// Termékek listájának mentése
	productsJSON, err := json.Marshal(products)
	if err != nil {
		return err
	}

	err = RedisClient.Set(ctx, "products", productsJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
