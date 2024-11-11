package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type TemplateData struct {
	Title string
	User  *models.User
	Data  interface{}
}

func HomeHandler(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
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

	tmpl := template.Must(template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/home.html",
	))

	data := TemplateData{
		Title: "Főoldal",
		User:  &user,
	}

	return tmpl.Execute(c.Response().Writer, data)
}

func LoginPageHandler(c echo.Context) error {
	tmpl := template.Must(template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/login.html",
	))

	data := TemplateData{
		Title: "Bejelentkezés",
	}

	return tmpl.Execute(c.Response().Writer, data)
}

func LoginHandler(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	var user models.User
	var hashedPassword string
	err := database.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&user.ID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.HTML(http.StatusUnauthorized, `
                <div id="loginResult" class="alert alert-danger">
                    Hibás email cím vagy jelszó!
                </div>
            `)
		}
		log.Printf("Database error during login: %v", err)
		return c.HTML(http.StatusInternalServerError, `
            <div id="loginResult" class="alert alert-danger">
                Belső szerverhiba történt, kérjük próbálja újra később.
            </div>
        `)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return c.HTML(http.StatusUnauthorized, `
            <div id="loginResult" class="alert alert-danger">
                Hibás email cím vagy jelszó!
            </div>
        `)
	}

	sessionID, err := CreateSession(user.ID)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		return c.HTML(http.StatusInternalServerError, `
            <div id="loginResult" class="alert alert-danger">
                Hiba történt a bejelentkezés során, kérjük próbálja újra később.
            </div>
        `)
	}

	c.SetCookie(&http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   int(sessionDuration.Seconds()),
	})

	// HTMX átirányítás
	c.Response().Header().Set("HX-Redirect", "/")
	return c.String(http.StatusOK, "")
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
		return c.HTML(http.StatusOK, `
            <div class="alert alert-danger">
                Érvénytelen email cím formátum!
            </div>
        `)
	}

	// Ellenőrizzük, hogy az email cím már regisztrálva van-e
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, `
            <div class="alert alert-danger">
                Hiba történt az email cím ellenőrzése során.
            </div>
        `)
	}

	if exists {
		return c.HTML(http.StatusOK, `
            <div class="alert alert-warning">
                Ez az email cím már regisztrálva van!
            </div>
        `)
	}

	return c.HTML(http.StatusOK, "")
}

// isValidEmail ellenőrzi az email cím formátumát
func isValidEmail(email string) bool {
	// Egyszerű email validáció
	email = strings.TrimSpace(email)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
