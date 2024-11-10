// internal/database/redis.go
package database

import (
	"context"

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
