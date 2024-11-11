package handlers

import (
	"net/http"
	"strconv"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/web/templates"
	"github.com/labstack/echo/v4"
)

func BalanceHandler(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	if c.Request().Method == http.MethodPost {
		return handleBalanceUpdate(c)
	}

	var balance float64
	err := database.DB.QueryRow("SELECT balance FROM users WHERE id = ?", userID).Scan(&balance)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching balance")
	}

	return templates.Balance(balance).Render(c.Request().Context(), c.Response().Writer)
}

func handleBalanceUpdate(c echo.Context) error {
	userID := c.Get("user_id").(string)
	amountStr := c.FormValue("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid amount")
	}

	_, err = database.DB.Exec("UPDATE users SET balance = balance + ? WHERE id = ?", amount, userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error updating balance")
	}

	var newBalance float64
	err = database.DB.QueryRow("SELECT balance FROM users WHERE id = ?", userID).Scan(&newBalance)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching updated balance")
	}

	return c.HTML(http.StatusOK, "<p>Egyenleg sikeresen frissítve! Új egyenleg: "+strconv.FormatFloat(newBalance, 'f', 2, 64)+" Ft</p>")
}
