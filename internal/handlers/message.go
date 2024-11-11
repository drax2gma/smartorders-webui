package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/models"
	"github.com/labstack/echo/v4"
)

func MessageHandler(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	if c.Request().Method == http.MethodPost {
		return handleMessageSubmission(c, userID)
	}

	// Get user data for the template
	var user models.User
	err := database.DB.QueryRow("SELECT id, name, email FROM users WHERE id = ?", userID).Scan(
		&user.ID, &user.Name, &user.Email,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching user data")
	}

	tmpl := template.Must(template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/message.html",
	))

	data := TemplateData{
		Title: "Üzenet küldése",
		User:  &user,
	}

	return tmpl.Execute(c.Response().Writer, data)
}

func handleMessageSubmission(c echo.Context, userID string) error {
	message := c.FormValue("message")
	if message == "" {
		return c.HTML(http.StatusBadRequest, `
            <div class="alert alert-danger">
                Kérjük, írjon üzenetet!
            </div>
        `)
	}

	_, err := database.DB.Exec(`
        INSERT INTO messages (user_id, content, created_at) 
        VALUES (?, ?, ?)
    `, userID, message, time.Now())

	if err != nil {
		return c.HTML(http.StatusInternalServerError, `
            <div class="alert alert-danger">
                Hiba történt az üzenet mentése során. Kérjük, próbálja újra később.
            </div>
        `)
	}

	return c.HTML(http.StatusOK, `
        <div class="alert alert-success">
            Az üzenet sikeresen elküldve!
        </div>
    `)
}
