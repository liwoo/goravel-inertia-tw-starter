package auth

import (
	"fmt"
	"smedi-sme-db/app/models"
	"time"

	"smedi-sme-db/app/services"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
)

type APIAuthController struct {
	passwordAttemptsService *services.PasswordAttemptsService
}

func NewAPIAuthController() *APIAuthController {
	return &APIAuthController{
		passwordAttemptsService: services.NewPasswordAttemptsService(),
	}
}

// APILoginRequest defines the structure for API login requests.
type APILoginRequest struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

// Authorize determines if the user is authorized to make this request.
func (r *APILoginRequest) Authorize(ctx http.Context) error {
	return nil
}

// Rules validation rules for APILoginRequest
func (r *APILoginRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"email":    "required|email",
		"password": "required",
	}
}

func (r *APILoginRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"email.required":    "Email is required.",
		"email.email":       "Please provide a valid email address.",
		"password.required": "Password is required.",
	}
}

// PrepareForValidation (optional data sanitization/modification before validation)
func (r *APILoginRequest) PrepareForValidation(data validation.Data) error {
	return nil
}

// Login handles API login and returns JSON response with token
func (c *APIAuthController) Login(ctx http.Context) http.Response {
	var loginRequest APILoginRequest
	errors, err := ctx.Request().ValidateRequest(&loginRequest)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"success": false,
			"message": "Validation error",
			"error":   err.Error(),
		})
	}
	if errors != nil {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"success": false,
			"message": "Validation failed",
			"errors":  errors.All(),
		})
	}

	// Check if account is locked before attempting login
	locked, unlockTime := c.passwordAttemptsService.CheckLocked(loginRequest.Email)
	if locked {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"success":          false,
			"message":          fmt.Sprintf("Account is locked due to too many failed login attempts. Please try again after %s", unlockTime.Format("2006-01-02T15:04:05Z07:00")),
			"locked":           true,
			"unlock_time":      unlockTime.Format("2006-01-02T15:04:05Z07:00"),
			"unlock_timestamp": unlockTime.Unix(),
		})
	}

	var user models.User
	// Find user by email
	if err := facades.Orm().Query().Where("email", loginRequest.Email).First(&user); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"success": false,
			"message": "Invalid credentials",
		})
	}

	// Check password
	if !facades.Hash().Check(loginRequest.Password, user.Password) {
		// Record failed attempt
		_, shouldWarn := c.passwordAttemptsService.RecordFailedAttempt(loginRequest.Email)

		// Check if account was just locked
		locked, unlockTime := c.passwordAttemptsService.CheckLocked(loginRequest.Email)
		if locked {
			return ctx.Response().Json(http.StatusForbidden, http.Json{
				"success":          false,
				"message":          fmt.Sprintf("Account has been locked due to too many failed login attempts. Please try again after %s", unlockTime.Format("2006-01-02T15:04:05Z07:00")),
				"locked":           true,
				"unlock_time":      unlockTime.Format("2006-01-02T15:04:05Z07:00"),
				"unlock_timestamp": unlockTime.Unix(),
			})
		}

		// Prepare response
		response := http.Json{
			"success":            false,
			"message":            "Invalid credentials",
			"remaining_attempts": c.passwordAttemptsService.GetRemainingAttempts(loginRequest.Email),
		}

		if shouldWarn {
			response["warning"] = "One more failed attempt will lock your account."
		}

		return ctx.Response().Json(http.StatusUnauthorized, response)
	}

	// Check if user is active
	if !user.IsActive {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"success": false,
			"message": "Account is deactivated",
		})
	}

	// Clear password attempts on successful login
	c.passwordAttemptsService.ClearAttempts(loginRequest.Email)

	// Log the user in and get the token
	token, err := facades.Auth(ctx).Login(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"success": false,
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
	}

	// Set token in HTTP-only cookie for web compatibility
	ttl := facades.Config().GetInt("jwt.ttl", 720) // Default to 12 hours (720 minutes)
	ctx.Response().Cookie(http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Duration(ttl) * time.Minute),
		Path:     "/",
		HttpOnly: true,
	})

	// Return token in JSON response for API clients
	return ctx.Response().Json(http.StatusOK, http.Json{
		"success": true,
		"message": "Login successful",
		"data": http.Json{
			"token": token,
			"user": http.Json{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"role":  user.Role,
			},
			"expires_in": ttl * 60, // Convert to seconds
		},
	})
}

// Logout handles API logout
func (c *APIAuthController) Logout(ctx http.Context) http.Response {
	if err := facades.Auth(ctx).Logout(); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"success": false,
			"message": "Failed to logout",
			"error":   err.Error(),
		})
	}

	// Clear the cookie
	ctx.Response().Cookie(http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/",
		HttpOnly: true,
	})

	return ctx.Response().Json(http.StatusOK, http.Json{
		"success": true,
		"message": "Logout successful",
	})
}

// Me returns the authenticated user's information
func (c *APIAuthController) Me(ctx http.Context) http.Response {
	// Get the authenticated user ID
	userId, err := facades.Auth(ctx).ID()
	if err != nil || userId == "" {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"success": false,
			"message": "Unauthenticated",
		})
	}

	// Load the user from database
	var user models.User
	if err := facades.Orm().Query().Where("id", userId).First(&user); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"success": false,
			"message": "User not found",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"success": true,
		"data": http.Json{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
