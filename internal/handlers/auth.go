package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/models"
	"github.com/drax2gma/smartorders-webui/web/templates"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func HomeHandler(c echo.Context) error {
	// Ellenőrizzük, hogy létezik-e a user_id
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	var user models.User
	err := database.DB.QueryRow("SELECT * FROM users WHERE id = ?", userID).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Balance, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Redirect(http.StatusSeeOther, "/login")
		}
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
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Hibás email cím vagy jelszó!",
			})
		}
		log.Printf("Database error during login: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Belső szerverhiba történt, kérjük próbálja újra később.",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Hibás email cím vagy jelszó!",
		})
	}

	sessionID, err := CreateSession(user.ID)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Hiba történt a bejelentkezés során, kérjük próbálja újra később.",
		})
	}

	c.SetCookie(&http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   int(sessionDuration.Seconds()),
	})

	return c.JSON(http.StatusOK, map[string]string{
		"redirect": "/",
	})
}

func LogoutHandler(c echo.Context) error {
	cookie, err := c.Cookie("session_id")
	if err == nil {
		DeleteSession(cookie.Value)
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

// Validate email to check if it is valid
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	if !strings.Contains(email, "@") {
		return false
	}
	return true
}
