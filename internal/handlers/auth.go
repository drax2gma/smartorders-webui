package handlers

import (
	"context"
	"encoding/json"
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

		if database.RedisClient == nil {
			log.Println("Redis client is not initialized")
			handleLoginError(w, "Internal server error")
			return
		}

		// Get user from Redis
		userJSON, err := database.RedisClient.Get(context.Background(), "user:"+email).Result()
		if err != nil {
			log.Printf("Error getting user from Redis: %v", err)
			handleLoginError(w, "Invalid email or password")
			return
		}

		var user models.User
		if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
			log.Printf("Error unmarshaling user data: %v", err)
			handleLoginError(w, "Error processing user data")
			return
		}

		// Get password hash from Redis
		storedHash, err := database.RedisClient.Get(context.Background(), "user:"+email+":password").Result()
		if err != nil {
			log.Printf("Error getting password hash from Redis: %v", err)
			handleLoginError(w, "Invalid email or password")
			return
		}

		// Check password
		if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
			log.Printf("Password mismatch for user %s", email)
			handleLoginError(w, "Invalid email or password")
			return
		}

		// Create session
		sessionID, err := CreateSession(user.ID)
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
		})

		// Átirányítás
		w.Header().Set("HX-Redirect", "/order")
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
