// internal/handlers/session.go
package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/go-redis/redis/v8"
)

type contextKey string

const (
	sessionContextKey contextKey = "session_id"
	userIDContextKey  contextKey = "user_id"

	sessionIDLength = 32
	sessionDuration = 24 * time.Hour
)

var redisClient *redis.Client

func InitRedis(addr string) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Ellenőrizzük a kapcsolatot
	_, err := redisClient.Ping(context.Background()).Result()
	return err
}

func CreateSession(userID uint) (string, error) {
	if database.RedisClient == nil {
		return "", fmt.Errorf("redis client is not initialized")
	}

	sessionID := generateSessionID()
	ctx := context.Background()
	err := database.RedisClient.Set(ctx, fmt.Sprintf("session:%s", sessionID), fmt.Sprintf("%d", userID), sessionDuration).Err()
	if err != nil {
		return "", fmt.Errorf("failed to set session in Redis: %v", err)
	}

	return sessionID, nil
}

func generateSessionID() string {
	b := make([]byte, sessionIDLength)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func GetUserIDFromSession(sessionID string) (uint, error) {
	userIDStr, err := redisClient.Get(redisClient.Context(), "session:"+sessionID).Result()
	if err != nil {
		return 0, err
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(userID), nil
}

func DeleteSession(sessionID string) error {
	return redisClient.Del(redisClient.Context(), "session:"+sessionID).Err()
}

func SessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID, err := GetUserIDFromSession(sessionID.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
