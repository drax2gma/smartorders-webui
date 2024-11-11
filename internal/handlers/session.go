package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/labstack/echo/v4"
)

const (
	sessionIDLength = 32
	sessionDuration = 24 * time.Hour
)

func CreateSession(userID string) (string, error) {
	sessionID := generateSessionID()
	_, err := database.DB.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)",
		sessionID, userID, time.Now().Add(sessionDuration))
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	return sessionID, nil
}

func generateSessionID() string {
	b := make([]byte, sessionIDLength)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func getUserIDFromSession(sessionID string) (string, error) {
	var userID string
	err := database.DB.QueryRow("SELECT user_id FROM sessions WHERE id = ? AND expires_at > ?", sessionID, time.Now()).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func DeleteSession(sessionID string) error {
	_, err := database.DB.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	return err
}

func SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_id")
		if err != nil || cookie.Value == "" {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		userID, err := getUserIDFromSession(cookie.Value)
		if err != nil || userID == "" {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		// Set the user ID in the context for later use
		c.Set("user_id", userID)
		return next(c)
	}
}
