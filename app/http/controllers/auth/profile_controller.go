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
		ctx.Request().Session().Flash("errors", map[string]interface{}{
			"server": "Data binding failed: " + err.Error(),
		})
		return ctx.Response().Redirect(http.StatusSeeOther, "/account")
	}

	// Manual validation - check required fields
	if passwordRequest.CurrentPassword == "" {
		ctx.Request().Session().Flash("errors", map[string]interface{}{
			"current_password": "Current password is required",
		})
		return ctx.Response().Redirect(http.StatusSeeOther, "/account")
	}
	if passwordRequest.Password == "" {
		ctx.Request().Session().Flash("errors", map[string]interface{}{
			"password": "New password is required",
		})
		return ctx.Response().Redirect(http.StatusSeeOther, "/account")
	}
	if passwordRequest.PasswordConfirmation == "" {
		ctx.Request().Session().Flash("errors", map[string]interface{}{
			"password_confirmation": "Password confirmation is required",
		})
		return ctx.Response().Redirect(http.StatusSeeOther, "/account")
	}

	// Validate password requirements
	if err := passwordRequest.ValidatePasswordRequirements(); err != nil {
		ctx.Request().Session().Flash("errors", map[string]interface{}{
			"password": err.Error(),
		})
		return ctx.Response().Redirect(http.StatusSeeOther, "/account")
	}

	// Validate password confirmation
	if err := passwordRequest.ValidatePasswordConfirmation(); err != nil {
		ctx.Request().Session().Flash("errors", map[string]interface{}{
			"password_confirmation": err.Error(),
		})
		return ctx.Response().Redirect(http.StatusSeeOther, "/account")
	}

	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user, err := permHelper.RequireAuthentication(ctx)
	if err != nil {
		ctx.Request().Session().Flash("errors", map[string]interface{}{
			"auth": "Authentication required",
		})
		return ctx.Response().Redirect(http.StatusSeeOther, "/account")
	}

	// Validate current password
	if !facades.Hash().Check(passwordRequest.CurrentPassword, user.Password) {
		ctx.Request().Session().Flash("errors", map[string]interface{}{
			"current_password": "Current password is incorrect",
		})
		return ctx.Response().Redirect(http.StatusSeeOther, "/account")
	}

	// Hash new password
	hashedPassword, err := facades.Hash().Make(passwordRequest.Password)
	if err != nil {
		facades.Log().Error("Failed to hash password for user", map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		})
		ctx.Request().Session().Flash("errors", map[string]interface{}{
			"server": "Failed to process password",
		})
		return ctx.Response().Redirect(http.StatusSeeOther, "/account")
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
		ctx.Request().Session().Flash("errors", map[string]interface{}{
			"server": "Failed to update password",
		})
		return ctx.Response().Redirect(http.StatusSeeOther, "/account")
	}

	// Log the password change for security auditing
	facades.Log().Info("User password updated successfully", map[string]interface{}{
		"user_id":    user.ID,
		"user_email": user.Email,
		"timestamp":  time.Now(),
		"ip_address": ctx.Request().Ip(),
		"user_agent": ctx.Request().Header("User-Agent", ""),
	})

	// For Inertia.js, we need to redirect back with success message
	ctx.Request().Session().Flash("success", "Password updated successfully!")
	return ctx.Response().Redirect(http.StatusSeeOther, "/account")
}
