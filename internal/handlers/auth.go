package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Get user ID from email
		userID, err := database.RedisClient.Get(context.Background(), fmt.Sprintf("email:%s", email)).Result()
		if err != nil {
			handleLoginError(w, "Invalid email or password")
			return
		}

		// Get user from Redis
		userJSON, err := database.RedisClient.Get(context.Background(), fmt.Sprintf("user:%s", userID)).Result()
		if err != nil {
			handleLoginError(w, "Invalid email or password")
			return
		}

		var user models.User
		if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
			handleLoginError(w, "Error processing user data")
			return
		}

		// Get password hash from Redis
		storedHash, err := database.RedisClient.Get(context.Background(), fmt.Sprintf("user:%s:password", userID)).Result()
		if err != nil {
			handleLoginError(w, "Invalid email or password")
			return
		}

		// Check password
		if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
			handleLoginError(w, "Invalid email or password")
			return
		}

		// Create session
		sessionID, err := CreateSession(userID)
		if err != nil {
			log.Printf("Failed to create session: %v", err)
			handleLoginError(w, "Error creating session")
			return
		}

		// Set session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			HttpOnly: true,
			Path:     "/",
			MaxAge:   int(sessionDuration.Seconds()),
		})

		// Átirányítás
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
		return
	}

	renderLoginPage(w, "")
}

func handleLoginError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func renderLoginPage(w http.ResponseWriter, errorMessage string) {
	tmpl, _ := template.ParseFiles("web/templates/login.gohtml")
	tmpl.Execute(w, map[string]string{"ErrorMessage": errorMessage})
}

func ValidateEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if !isValidEmail(email) {
		w.Write([]byte("Érvénytelen email cím"))
		return
	}
	w.Write([]byte(""))
}

func isValidEmail(email string) bool {
	// Egyszerű email validáció, a gyakorlatban használj robusztusabb megoldást
	return regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(email)
}
