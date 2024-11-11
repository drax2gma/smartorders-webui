package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/drax2gma/smartorders-webui/internal/database"
	"github.com/drax2gma/smartorders-webui/internal/models"
	"github.com/labstack/echo/v4"
)

func BalanceHandler(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	if c.Request().Method == http.MethodPost {
		return handleBalanceUpdate(c, userID)
	}

	var user models.User
	err := database.DB.QueryRow("SELECT * FROM users WHERE id = ?", userID).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Balance, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching user data")
	}

	tmpl := template.Must(template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/balance.html",
	))

	data := TemplateData{
		Title: "Egyenleg",
		User:  &user,
		Data:  user.Balance,
	}

	return tmpl.Execute(c.Response().Writer, data)
}

func handleBalanceUpdate(c echo.Context, userID string) error {
	amountStr := c.FormValue("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		return c.HTML(http.StatusBadRequest, `
            <div class="alert alert-danger">
                Érvénytelen összeg. Kérjük, adjon meg egy pozitív számot.
            </div>
        `)
	}

	// Start transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return c.HTML(http.StatusInternalServerError, `
            <div class="alert alert-danger">
                Hiba történt a tranzakció során. Kérjük, próbálja újra.
            </div>
        `)
	}
	defer tx.Rollback()

	// Update balance
	result, err := tx.Exec("UPDATE users SET balance = balance + ? WHERE id = ?", amount, userID)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, `
            <div class="alert alert-danger">
                Hiba történt az egyenleg frissítése során. Kérjük, próbálja újra.
            </div>
        `)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return c.HTML(http.StatusInternalServerError, `
            <div class="alert alert-danger">
                Nem sikerült frissíteni az egyenleget. Kérjük, próbálja újra.
            </div>
        `)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return c.HTML(http.StatusInternalServerError, `
            <div class="alert alert-danger">
                Hiba történt a tranzakció véglegesítése során. Kérjük, próbálja újra.
            </div>
        `)
	}

	// Get updated balance
	var newBalance float64
	err = database.DB.QueryRow("SELECT balance FROM users WHERE id = ?", userID).Scan(&newBalance)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, `
            <div class="alert alert-danger">
                Az egyenleg frissítve, de nem sikerült lekérdezni az új egyenleget.
            </div>
        `)
	}

	return c.HTML(http.StatusOK, `
        <div class="alert alert-success">
            Az egyenleg sikeresen feltöltve! Új egyenleg: `+strconv.FormatFloat(newBalance, 'f', 2, 64)+` Ft
        </div>
    `)
}
