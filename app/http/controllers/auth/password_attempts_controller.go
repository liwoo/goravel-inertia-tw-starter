package auth

import (
	"smedi-sme-db/app/services"

	"github.com/goravel/framework/contracts/http"
)

type PasswordAttemptsController struct {
	passwordAttemptsService *services.PasswordAttemptsService
}

func NewPasswordAttemptsController() *PasswordAttemptsController {
	return &PasswordAttemptsController{
		passwordAttemptsService: services.NewPasswordAttemptsService(),
	}
}

// CheckStatus checks if an email is locked and returns lockout information
func (c *PasswordAttemptsController) CheckStatus(ctx http.Context) http.Response {
	email := ctx.Request().Query("email", "")
	if email == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"success": false,
			"message": "Email parameter is required",
		})
	}

	locked, unlockTime := c.passwordAttemptsService.CheckLocked(email)
	remainingAttempts := c.passwordAttemptsService.GetRemainingAttempts(email)

	response := http.Json{
		"success":            true,
		"locked":             locked,
		"remaining_attempts": remainingAttempts,
	}

	if locked {
		response["unlock_time"] = unlockTime.Format("2006-01-02T15:04:05Z07:00")
		response["unlock_timestamp"] = unlockTime.Unix()
	}

	return ctx.Response().Json(http.StatusOK, response)
}

// GetRemainingAttempts returns the number of remaining attempts before lockout
func (c *PasswordAttemptsController) GetRemainingAttempts(ctx http.Context) http.Response {
	email := ctx.Request().Query("email", "")
	if email == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"success": false,
			"message": "Email parameter is required",
		})
	}

	locked, unlockTime := c.passwordAttemptsService.CheckLocked(email)
	remainingAttempts := c.passwordAttemptsService.GetRemainingAttempts(email)

	response := http.Json{
		"success":            true,
		"remaining_attempts": remainingAttempts,
		"locked":             locked,
	}

	if locked {
		response["unlock_time"] = unlockTime.Format("2006-01-02T15:04:05Z07:00")
		response["unlock_timestamp"] = unlockTime.Unix()
	}

	return ctx.Response().Json(http.StatusOK, response)
}
