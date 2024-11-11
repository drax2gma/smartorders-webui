package handlers

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/models"
	"github.com/drax2gma/smartorders-webui/web/templates"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func HomeHandler(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	var user models.User
	err := database.DB.QueryRow("SELECT * FROM users WHERE id = ?", userID).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Balance, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching user data")
	}

	return templates.Home(user).Render(c.Request().Context(), c.Response().Writer)
}

func LoginPageHandler(c echo.Context) error {
	return templates.Login().Render(c.Request().Context(), c.Response().Writer)
}

func LoginHandler(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	var user models.User
	var hashedPassword string
	err := database.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&user.ID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}

	// Create session
	sessionID := generateSessionID(user.ID)
	_, err = database.DB.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)",
		sessionID, user.ID, time.Now().Add(24*time.Hour))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error creating session"})
	}

	c.SetCookie(&http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   int(24 * time.Hour.Seconds()),
	})

	return c.JSON(http.StatusOK, map[string]string{"redirect": "/"})
}

func LogoutHandler(c echo.Context) error {
	cookie, err := c.Cookie("session_id")
	if err == nil {
		_, err = database.DB.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
		if err != nil {
			// Log the error, but continue with logout process
			fmt.Printf("Error deleting session: %v\n", err)
		}
	}

	c.SetCookie(&http.Cookie{
		Name:     "session_id",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1,
	})

	return c.Redirect(http.StatusSeeOther, "/login")
}

func ValidateEmailHandler(c echo.Context) error {
	email := c.FormValue("email")
	if !isValidEmail(email) {
		return c.String(http.StatusOK, "Érvénytelen email cím")
	}
	return c.String(http.StatusOK, "")
}

func generateSessionID(userID string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(userID+time.Now().String())))[:16]
}

func isValidEmail(email string) bool {
	// Implement email validation logic here
	return true // Placeholder
}
