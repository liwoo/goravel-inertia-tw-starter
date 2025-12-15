package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
	"sync"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
)

type AuthController struct {
	activityService *services.UserActivityService
	totpService     *services.TOTPService

	// In-memory storage for pending 2FA logins (web route)
	pending2FALogins     map[string]pending2FALoginWeb
	pending2FALoginsLock sync.RWMutex
}

// pending2FALoginWeb stores user info during the 2FA verification step for web login
type pending2FALoginWeb struct {
	UserID    uint
	ExpiresAt time.Time
}

func NewAuthController() *AuthController {
	controller := &AuthController{
		activityService:  services.NewUserActivityService(),
		totpService:      services.NewTOTPService(),
		pending2FALogins: make(map[string]pending2FALoginWeb),
	}
	// Start cleanup goroutine
	go controller.cleanupExpired2FALogins()
	return controller
}

// generateTempToken generates a secure random token
func (r *AuthController) generateTempToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// cleanupExpired2FALogins removes expired pending 2FA logins periodically
func (r *AuthController) cleanupExpired2FALogins() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		r.pending2FALoginsLock.Lock()
		now := time.Now()
		for token, login := range r.pending2FALogins {
			if now.After(login.ExpiresAt) {
				delete(r.pending2FALogins, token)
			}
		}
		r.pending2FALoginsLock.Unlock()
	}
}

// LoginRequest defines the structure for login requests.
// It implements the http.FormRequest interface.
type LoginRequest struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

// Authorize determines if the user is authorized to make this request.
// For login, it's typically true as any user can attempt to login.
func (r *LoginRequest) Authorize(ctx http.Context) error {
	// You can add complex authorization logic here if needed.
	// For example, check if the IP is blacklisted, etc.
	// Returning nil means the request is authorized.
	// Returning an error will stop the request and return a 403 Forbidden response.
	return nil
}

// Validation rules for LoginRequest
func (r *LoginRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"email":    "required|email",
		"password": "required",
	}
}

func (r *LoginRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"email.required":    "Email is required.",
		"email.email":       "Please provide a valid email address.",
		"password.required": "Password is required.",
	}
}

// PrepareForValidation (optional data sanitization/modification before validation)
func (r *LoginRequest) PrepareForValidation(data validation.Data) error {
	// Example: trim spaces from email
	// if email, ok := data.Get("email").(string); ok {
	// 	data.Set("email", strings.TrimSpace(email))
	// }
	return nil
}

func (r *AuthController) Login(ctx http.Context) http.Response {
	var loginRequest LoginRequest
	errors, err := ctx.Request().ValidateRequest(&loginRequest)
	if err != nil {
		// Flash general error message and redirect back
		ctx.Request().Session().Flash("errors", map[string]interface{}{
			"general": "Error validating request: " + err.Error(),
		})
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}
	if errors != nil {
		// Flash validation errors directly
		ctx.Request().Session().Flash("errors", errors.All())
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	var user models.User
	// Find user by email
	if err := facades.Orm().Query().Where("email", loginRequest.Email).First(&user); err != nil {
		// Flash field-specific error
		ctx.Request().Session().Flash("errors", map[string]interface{}{
			"email": "Invalid credentials (Email not found)",
		})
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	// Check password
	if !facades.Hash().Check(loginRequest.Password, user.Password) {
		// Flash field-specific error
		ctx.Request().Session().Flash("errors", map[string]interface{}{
			"password": "Password is incorrect",
		})
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	// Check if user is active
	if !user.IsActive {
		ctx.Request().Session().Flash("errors", map[string]interface{}{
			"general": "Account is deactivated",
		})
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	// Check if 2FA is enabled
	if user.TOTPEnabled {
		// Generate a temporary token for 2FA verification
		tempToken, err := r.generateTempToken()
		if err != nil {
			facades.Log().Error("Failed to generate temp token", map[string]interface{}{
				"error":   err.Error(),
				"user_id": user.ID,
			})
			ctx.Request().Session().Flash("errors", map[string]interface{}{
				"general": "Failed to process login",
			})
			return ctx.Response().Redirect(http.StatusFound, "/login")
		}

		// Store pending 2FA login (expires in 5 minutes)
		r.pending2FALoginsLock.Lock()
		r.pending2FALogins[tempToken] = pending2FALoginWeb{
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}
		r.pending2FALoginsLock.Unlock()

		// For Inertia, return with 2FA required props
		// The frontend will handle showing the 2FA form
		ctx.Request().Session().Flash("requires_2fa", true)
		ctx.Request().Session().Flash("temp_token", tempToken)

		// Return JSON for Inertia to handle
		return ctx.Response().Json(http.StatusOK, http.Json{
			"requires_2fa": true,
			"temp_token":   tempToken,
		})
	}

	// No 2FA - proceed with normal login
	return r.completeLogin(ctx, &user)
}

// completeLogin finalizes the login process
func (r *AuthController) completeLogin(ctx http.Context, user *models.User) http.Response {
	// Log the user in and get the token
	token, err := facades.Auth(ctx).Login(user)
	if err != nil {
		ctx.Request().Session().Flash("errors", map[string]interface{}{
			"general": "Error during login: " + err.Error(),
		})
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	// Set token in HTTP-only cookie
	ttl := facades.Config().GetInt("jwt.ttl", 720) // Default to 12 hours (720 minutes) if not set
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

	// Update last login timestamp
	now := time.Now()
	facades.Orm().Query().Model(user).Update("last_login_at", now)

	// Log login activity
	r.activityService.LogActivity(
		ctx,
		user.ID,
		models.ActivityLogin,
		"User logged in",
		map[string]interface{}{
			"email": user.Email,
		},
	)

	// Check if 2FA is required globally but user hasn't enabled it
	require2FA := facades.Config().GetBool("auth.require_2fa", false)
	if require2FA && !user.TOTPEnabled {
		// Redirect to 2FA setup page
		return ctx.Response().Redirect(http.StatusSeeOther, "/2fa-required")
	}

	// Redirect to dashboard on successful login.
	// Use 303 See Other to ensure the next request is a GET, which is best practice for Inertia.
	return ctx.Response().Redirect(http.StatusSeeOther, "/dashboard")
}

// Verify2FA handles 2FA verification for web login
func (r *AuthController) Verify2FA(ctx http.Context) http.Response {
	tempToken := ctx.Request().Input("temp_token")
	code := ctx.Request().Input("code")

	if tempToken == "" || code == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"success": false,
			"message": "Missing required fields",
		})
	}

	// Find pending 2FA login
	r.pending2FALoginsLock.RLock()
	pendingLogin, exists := r.pending2FALogins[tempToken]
	r.pending2FALoginsLock.RUnlock()

	if !exists || time.Now().After(pendingLogin.ExpiresAt) {
		return ctx.Response().Json(http.StatusGone, http.Json{
			"success": false,
			"message": "Login session expired. Please log in again.",
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
	valid, isBackup, err := r.totpService.ValidateTOTPForUser(&user, code)
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
	r.pending2FALoginsLock.Lock()
	delete(r.pending2FALogins, tempToken)
	r.pending2FALoginsLock.Unlock()

	// Log activity
	if isBackup {
		r.activityService.LogActivity(ctx, user.ID, models.ActivityBackupCodeUsed, "Backup code used during login", nil)
	} else {
		r.activityService.LogActivity(ctx, user.ID, models.ActivityTwoFAUsed, "2FA code used during login", nil)
	}

	// Complete login - generate token and set cookies
	token, err := facades.Auth(ctx).Login(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"success": false,
			"message": "Failed to complete login",
		})
	}

	// Set token in HTTP-only cookie
	ttl := facades.Config().GetInt("jwt.ttl", 720)
	expiry := time.Now().Add(time.Duration(ttl) * time.Minute)
	ctx.Response().Cookie(http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expiry,
		Path:     "/",
		HttpOnly: true,
	})

	// Set a non-HttpOnly cookie for SSE connections
	ctx.Response().Cookie(http.Cookie{
		Name:     "sse_token",
		Value:    token,
		Expires:  expiry,
		Path:     "/",
		HttpOnly: false,
	})

	// Update last login timestamp
	now := time.Now()
	facades.Orm().Query().Model(&user).Update("last_login_at", now)

	// Log login activity
	r.activityService.LogActivity(
		ctx,
		user.ID,
		models.ActivityLogin,
		"User logged in with 2FA",
		map[string]interface{}{
			"email":            user.Email,
			"used_backup_code": isBackup,
		},
	)

	return ctx.Response().Json(http.StatusOK, http.Json{
		"success": true,
		"message": "Login successful",
		"data": http.Json{
			"redirect": "/dashboard",
		},
	})
}

func (r *AuthController) Logout(ctx http.Context) http.Response {
	// Get the current user before logging out (for activity logging)
	var user models.User
	facades.Auth(ctx).User(&user)

	if err := facades.Auth(ctx).Logout(); err != nil {
		// It's good to log this, but for the user, redirecting is usually best.
		facades.Log().Error("Error during logout: " + err.Error())
		fmt.Println("Error during logout: " + err.Error())
		// Even if logout fails on the server, try to clear client-side session by redirecting.
		return ctx.Response().Redirect(http.StatusFound, "/")
	}

	// Log logout activity (if we had a valid user)
	if user.ID != 0 {
		r.activityService.LogActivity(
			ctx,
			user.ID,
			models.ActivityLogout,
			"User logged out",
			nil,
		)
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

	fmt.Println("Logout successful")

	return ctx.Response().Redirect(http.StatusFound, "/")
}
