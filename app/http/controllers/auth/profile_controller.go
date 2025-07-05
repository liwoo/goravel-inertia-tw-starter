package auth

import (
	"time"

	"players/app/auth"
	"players/app/http/requests"
	"players/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type ProfileController struct{}

func NewProfileController() *ProfileController {
	return &ProfileController{}
}

// UpdatePassword handles password update requests for the authenticated user's profile
func (c *ProfileController) UpdatePassword(ctx http.Context) http.Response {
	// Validate request using manual binding and validation
	var passwordRequest requests.ProfilePasswordUpdateRequest

	// Bind the data to the struct
	if err := ctx.Request().Bind(&passwordRequest); err != nil {
		return ctx.Response().Status(http.StatusBadRequest).Json(map[string]interface{}{
			"errors": map[string]string{"server": "Data binding failed: " + err.Error()},
		})
	}

	// Manual validation - check required fields
	if passwordRequest.CurrentPassword == "" {
		return ctx.Response().Status(http.StatusUnprocessableEntity).Json(map[string]interface{}{
			"errors": map[string]string{"current_password": "Current password is required"},
		})
	}
	if passwordRequest.Password == "" {
		return ctx.Response().Status(http.StatusUnprocessableEntity).Json(map[string]interface{}{
			"errors": map[string]string{"password": "New password is required"},
		})
	}
	if passwordRequest.PasswordConfirmation == "" {
		return ctx.Response().Status(http.StatusUnprocessableEntity).Json(map[string]interface{}{
			"errors": map[string]string{"password_confirmation": "Password confirmation is required"},
		})
	}

	// Validate password requirements
	if err := passwordRequest.ValidatePasswordRequirements(); err != nil {
		return ctx.Response().Status(http.StatusUnprocessableEntity).Json(map[string]interface{}{
			"errors": map[string]string{"password": err.Error()},
		})
	}

	// Validate password confirmation
	if err := passwordRequest.ValidatePasswordConfirmation(); err != nil {
		return ctx.Response().Status(http.StatusUnprocessableEntity).Json(map[string]interface{}{
			"errors": map[string]string{"password_confirmation": err.Error()},
		})
	}

	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user, err := permHelper.RequireAuthentication(ctx)
	if err != nil {
		return ctx.Response().Status(http.StatusUnauthorized).Json(map[string]interface{}{
			"errors": map[string]string{"auth": "Authentication required"},
		})
	}

	// Validate current password
	if !facades.Hash().Check(passwordRequest.CurrentPassword, user.Password) {
		return ctx.Response().Status(http.StatusUnprocessableEntity).Json(map[string]interface{}{
			"errors": map[string]string{"current_password": "Current password is incorrect"},
		})
	}

	// Hash new password
	hashedPassword, err := facades.Hash().Make(passwordRequest.Password)
	if err != nil {
		facades.Log().Error("Failed to hash password for user", map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		})
		return ctx.Response().Status(http.StatusInternalServerError).Json(map[string]interface{}{
			"errors": map[string]string{"server": "Failed to process password"},
		})
	}

	// Update password in database
	_, err = facades.Orm().Query().Model(&models.User{}).
		Where("id = ?", user.ID).
		Update("password", hashedPassword)

	if err != nil {
		facades.Log().Error("Failed to update password in database", map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		})
		return ctx.Response().Status(http.StatusInternalServerError).Json(map[string]interface{}{
			"errors": map[string]string{"server": "Failed to update password"},
		})
	}

	// Log the password change for security auditing
	facades.Log().Info("User password updated successfully", map[string]interface{}{
		"user_id":    user.ID,
		"user_email": user.Email,
		"timestamp":  time.Now(),
		"ip_address": ctx.Request().Ip(),
		"user_agent": ctx.Request().Header("User-Agent", ""),
	})

	return ctx.Response().Status(http.StatusOK).Json(map[string]interface{}{
		"message": "Password updated successfully",
	})
}
