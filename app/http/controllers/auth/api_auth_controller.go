package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"starter-project/app/http/requests"
	"starter-project/app/models"
	"starter-project/app/services"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
)

type APIAuthController struct {
	totpService     *services.TOTPService
	activityService *services.UserActivityService

	// In-memory storage for pending 2FA logins
	pending2FALogins     map[string]pending2FALogin
	pending2FALoginsLock sync.RWMutex
}

// pending2FALogin stores user info during the 2FA verification step
type pending2FALogin struct {
	UserID    uint
	ExpiresAt time.Time
}

func NewAPIAuthController() *APIAuthController {
	controller := &APIAuthController{
		totpService:      services.NewTOTPService(),
		activityService:  services.NewUserActivityService(),
		pending2FALogins: make(map[string]pending2FALogin),
	}
	// Start cleanup goroutine
	go controller.cleanupExpired2FALogins()
	return controller
}

// generateTempToken generates a secure random token
func (c *APIAuthController) generateTempToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// cleanupExpired2FALogins removes expired pending 2FA logins periodically
func (c *APIAuthController) cleanupExpired2FALogins() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		c.pending2FALoginsLock.Lock()
		now := time.Now()
		for token, login := range c.pending2FALogins {
			if now.After(login.ExpiresAt) {
				delete(c.pending2FALogins, token)
			}
		}
		c.pending2FALoginsLock.Unlock()
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
// If 2FA is enabled, returns requires_2fa: true with a temp_token
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
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"success": false,
			"message": "Invalid credentials",
		})
	}

	// Check if user is active
	if !user.IsActive {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"success": false,
			"message": "Account is deactivated",
		})
	}

	// Check if 2FA is enabled
	if user.TOTPEnabled {
		// Generate a temporary token for 2FA verification
		tempToken, err := c.generateTempToken()
		if err != nil {
			facades.Log().Error("Failed to generate temp token", map[string]interface{}{
				"error":   err.Error(),
				"user_id": user.ID,
			})
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"success": false,
				"message": "Failed to process login",
			})
		}

		// Store pending 2FA login (expires in 5 minutes)
		c.pending2FALoginsLock.Lock()
		c.pending2FALogins[tempToken] = pending2FALogin{
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}
		c.pending2FALoginsLock.Unlock()

		return ctx.Response().Json(http.StatusOK, http.Json{
			"success":      true,
			"requires_2fa": true,
			"message":      "2FA verification required",
			"data": http.Json{
				"temp_token": tempToken,
			},
		})
	}

	// No 2FA - proceed with normal login
	return c.completeLogin(ctx, &user)
}

// Verify2FA godoc
// @Summary      Verify 2FA during login
// @Description  Complete the login process by providing the 2FA code
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body  requests.Verify2FALoginRequest  true  "Temp token and 2FA code"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Router       /auth/verify-2fa [post]
func (c *APIAuthController) Verify2FA(ctx http.Context) http.Response {
	// Validate request
	var request requests.Verify2FALoginRequest
	errors, err := ctx.Request().ValidateRequest(&request)
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

	// Find pending 2FA login
	c.pending2FALoginsLock.RLock()
	pendingLogin, exists := c.pending2FALogins[request.TempToken]
	c.pending2FALoginsLock.RUnlock()

	if !exists || time.Now().After(pendingLogin.ExpiresAt) {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"success": false,
			"message": "Invalid or expired login session. Please log in again.",
		})
	}

	// Load user
	var user models.User
	if err := facades.Orm().Query().Find(&user, pendingLogin.UserID); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"success": false,
			"message": "Failed to load user",
		})
	}

	// Validate TOTP code or backup code
	valid, isBackup, err := c.totpService.ValidateTOTPForUser(&user, request.Code)
	if err != nil {
		facades.Log().Error("Failed to validate 2FA code", map[string]interface{}{
			"error":   err.Error(),
			"user_id": user.ID,
		})
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"success": false,
			"message": "Failed to validate 2FA code",
		})
	}

	if !valid {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"success": false,
			"message": "Invalid 2FA code",
		})
	}

	// Remove pending login
	c.pending2FALoginsLock.Lock()
	delete(c.pending2FALogins, request.TempToken)
	c.pending2FALoginsLock.Unlock()

	// Log activity
	if isBackup {
		c.activityService.LogActivity(ctx, user.ID, models.ActivityBackupCodeUsed, "Backup code used during login", nil)
	} else {
		c.activityService.LogActivity(ctx, user.ID, models.ActivityTwoFAUsed, "2FA code used during login", nil)
	}

	// Complete login
	return c.completeLogin(ctx, &user)
}

// completeLogin finalizes the login process and returns the token
func (c *APIAuthController) completeLogin(ctx http.Context, user *models.User) http.Response {
	// Log the user in and get the token
	token, err := facades.Auth(ctx).Login(user)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"success": false,
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
	}

	// Set token in HTTP-only cookie for web compatibility
	ttl := facades.Config().GetInt("jwt.ttl", 720) // Default to 12 hours (720 minutes)
	expiry := time.Now().Add(time.Duration(ttl) * time.Minute)
	ctx.Response().Cookie(http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expiry,
		Path:     "/",
		HttpOnly: true,
	})

	// Set a non-HttpOnly cookie for SSE connections (JavaScript needs to read this)
	ctx.Response().Cookie(http.Cookie{
		Name:     "sse_token",
		Value:    token,
		Expires:  expiry,
		Path:     "/",
		HttpOnly: false,
	})

	// Log login activity
	c.activityService.LogActivity(ctx, user.ID, models.ActivityLogin, "User logged in", nil)

	// Return token in JSON response for API clients
	return ctx.Response().Json(http.StatusOK, http.Json{
		"success": true,
		"message": "Login successful",
		"data": http.Json{
			"token": token,
			"user": http.Json{
				"id":           user.ID,
				"name":         user.Name,
				"email":        user.Email,
				"role":         user.Role,
				"totp_enabled": user.TOTPEnabled,
			},
			"expires_in": ttl * 60, // Convert to seconds
		},
	})
}

// Logout handles API logout
func (c *APIAuthController) Logout(ctx http.Context) http.Response {
	// Get user for activity logging before logout
	var user models.User
	if err := facades.Auth(ctx).User(&user); err == nil {
		c.activityService.LogActivity(ctx, user.ID, models.ActivityLogout, "User logged out", nil)
	}

	if err := facades.Auth(ctx).Logout(); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"success": false,
			"message": "Failed to logout",
			"error":   err.Error(),
		})
	}

	// Clear the cookies
	ctx.Response().Cookie(http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/",
		HttpOnly: true,
	})
	ctx.Response().Cookie(http.Cookie{
		Name:     "sse_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/",
		HttpOnly: false,
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
			"id":           user.ID,
			"name":         user.Name,
			"email":        user.Email,
			"role":         user.Role,
			"totp_enabled": user.TOTPEnabled,
		},
	})
}
