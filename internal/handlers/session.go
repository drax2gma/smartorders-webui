// internal/handlers/session.go
package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/drax2gma/smartorders-webui/internal/database"
)

type contextKey string

const (
	sessionContextKey contextKey = "session_id"
	userIDContextKey  contextKey = "user_id"

	sessionIDLength = 32
	sessionDuration = 24 * time.Hour
)

func CreateSession(userID uint) (string, error) {
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

func getUserIDFromSession(sessionID string) (uint, error) {
	ctx := context.Background()
	userIDStr, err := database.RedisClient.Get(ctx, fmt.Sprintf("session:%s", sessionID)).Result()
	if err != nil {
		return 0, err
	}
	var userID uint
	_, err = fmt.Sscanf(userIDStr, "%d", &userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func DeleteSession(sessionID string) error {
	ctx := context.Background()
	return database.RedisClient.Del(ctx, "session:"+sessionID).Err()
}

func SessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID, err := getUserIDFromSession(sessionID.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
