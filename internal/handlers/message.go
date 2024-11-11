package handlers

import (
	"net/http"
	"time"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/web/templates"
	"github.com/labstack/echo/v4"
)

func MessageHandler(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	if c.Request().Method == http.MethodPost {
		return handleMessageSend(c)
	}

	return templates.Message().Render(c.Request().Context(), c.Response().Writer)
}

func handleMessageSend(c echo.Context) error {
	userID := c.Get("user_id").(string)
	message := c.FormValue("message")

	_, err := database.DB.Exec("INSERT INTO messages (user_id, content, created_at) VALUES (?, ?, ?)",
		userID, message, time.Now())
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error sending message")
	}

	return c.HTML(http.StatusOK, "<p>Üzenet sikeresen elküldve!</p>")
}
