package auth

import (
	"fmt"
	"players/app/models" // Assuming your User model is here
	"time"

	inertiaHelper "players/app/http/inertia"

	"github.com/google/uuid"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
)

type AuthController struct {
	// Dependencies can be injected here
}

func NewAuthController() *AuthController {
	return &AuthController{}
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
		// For Inertia, it's often better to redirect back with errors
		// or return a JSON response that Inertia can handle to show errors on the form.
		// However, for a direct API-like error, this is fine.
		return ctx.Response().Status(http.StatusInternalServerError).Json(http.Json{
			"message": "Error validating request: " + err.Error(),
		})
	}
	if errors != nil {
		// Redirect back with validation errors for Inertia to display
		// This assumes your frontend is set up to handle these errors.
		// If using inertia-react, errors are typically passed as props.
		// For simplicity in this step, we'll return JSON, but a redirect back is common.
		return ctx.Response().Status(http.StatusUnprocessableEntity).Json(errors.All())
	}

	var user models.User
	// Find user by email
	if err := facades.Orm().Query().Where("email", loginRequest.Email).First(&user); err != nil {
		// Return error that can be displayed on the login form
		return ctx.Response().Status(http.StatusUnauthorized).Json(http.Json{
			"errors": map[string]string{"email": "Invalid credentials (email not found)"},
		})
	}

	// Check password
	if !facades.Hash().Check(loginRequest.Password, user.Password) {
		// Return error that can be displayed on the login form
		return ctx.Response().Status(http.StatusUnauthorized).Json(http.Json{
			"errors": map[string]string{"password": "Invalid credentials (password mismatch)"},
		})
	}

	// Log the user in and get the token
	token, err := facades.Auth(ctx).Login(&user)
	if err != nil {
		return ctx.Response().Status(http.StatusInternalServerError).Json(http.Json{
			"message": "Error during login: " + err.Error(),
		})
	}

	// Set token in HTTP-only cookie
	ttl := facades.Config().GetInt("jwt.ttl", 720) // Default to 12 hours (720 minutes) if not set
	ctx.Response().Cookie(http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Duration(ttl) * time.Minute),
		Path:     "/",
		HttpOnly: true,
	})

	// Redirect to dashboard on successful login.
	// Use 303 See Other to ensure the next request is a GET, which is best practice for Inertia.
	return ctx.Response().Redirect(http.StatusSeeOther, "/dashboard")
}

func (r *AuthController) Logout(ctx http.Context) http.Response {
	if err := facades.Auth(ctx).Logout(); err != nil {
		// It's good to log this, but for the user, redirecting is usually best.
		facades.Log().Error("Error during logout: " + err.Error())
		fmt.Println("Error during logout: " + err.Error())
		// Even if logout fails on the server, try to clear client-side session by redirecting.
		return ctx.Response().Redirect(http.StatusFound, "/")
	}

	fmt.Println("Logout successful")

	return ctx.Response().Redirect(http.StatusFound, "/")

}

// ForgotPasswordRequest defines the structure for forgot password requests.
type ForgotPasswordRequest struct {
	Email string `form:"email" json:"email"`
}

func (r *ForgotPasswordRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *ForgotPasswordRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"email": "required|email",
	}
}

func (r *ForgotPasswordRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"email.required": "Email is required.",
		"email.email":    "Please provide a valid email address.",
	}
}

// ForgotPassword handles sending a password reset link (logs to console)
func (r *AuthController) ForgotPassword(ctx http.Context) http.Response {
	var req ForgotPasswordRequest
	errors, err := ctx.Request().ValidateRequest(&req)
	if err != nil {
		return ctx.Response().Status(http.StatusInternalServerError).Json(http.Json{
			"message": "Error validating request: " + err.Error(),
		})
	}
	if errors != nil {
		return ctx.Response().Status(http.StatusUnprocessableEntity).Json(errors.All())
	}

	// Find user by email (do not reveal if not found)
	var user models.User
	userExists := facades.Orm().Query().Where("email", req.Email).First(&user) == nil

	if userExists {
		// Generate a reset token and expiration (1 hour from now)
		token := uuid.NewString()
		expires := time.Now().Add(1 * time.Hour)
		// Store token in Redis: key = "reset_token:<email>", value = token, expires in 1 hour
		cacheKey := fmt.Sprintf("reset_token:%s", req.Email)
		facades.Cache().Put(cacheKey, token, time.Hour)
		// Build base URL from env vars
		appURL := facades.Config().GetString("app.url", "http://localhost")
		appPort := facades.Config().GetString("app.port", "3000")
		baseURL := appURL
		if appPort != "80" && appPort != "443" && appPort != "" {
			baseURL = fmt.Sprintf("%s:%s", appURL, appPort)
		}
		resetLink := fmt.Sprintf("%s/reset-password?token=%s&email=%s", baseURL, token, req.Email)
		facades.Log().Info(fmt.Sprintf("Password reset link for %s (expires %s): %s", req.Email, expires.Format(time.RFC3339), resetLink))
		fmt.Printf("Password reset link for %s (expires %s): %s\n", req.Email, expires.Format(time.RFC3339), resetLink)
	}

	// Always return success (do not reveal if email exists)
	return ctx.Response().Redirect(http.StatusFound, "/forgot-password-confirmation")
}

// ShowResetPassword displays the reset password page with token validation
func (r *AuthController) ShowResetPassword(ctx http.Context) http.Response {
	// Get token and email from query parameters
	token := ctx.Request().Query("token", "")
	email := ctx.Request().Query("email", "")

	// Basic validation
	if token == "" || email == "" {
		return ctx.Response().Redirect(http.StatusFound, "/invalid-token")
	}

	// Validate token format (should be a valid UUID)
	if _, err := uuid.Parse(token); err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/invalid-token")
	}

	// Check if token exists in Redis
	cacheKey := fmt.Sprintf("reset_token:%s", email)
	cachedValue := facades.Cache().Get(cacheKey, nil)

	if cachedValue == nil {
		return ctx.Response().Redirect(http.StatusFound, "/invalid-token")
	}

	// Validate cached token
	var redisToken string
	if tokenStr, ok := cachedValue.(string); ok {
		redisToken = tokenStr
	} else {
		return ctx.Response().Redirect(http.StatusFound, "/invalid-token")
	}

	// Compare tokens
	if redisToken != token {
		return ctx.Response().Redirect(http.StatusFound, "/invalid-token")
	}

	// Token is valid, show the reset password page
	return inertiaHelper.Render(ctx, "auth/ResetPassword", map[string]interface{}{
		"version": support.Version,
		"token":   token,
		"email":   email,
	})
}

// ResetPasswordRequest defines the structure for reset password requests.
type ResetPasswordRequest struct {
	Email                string `form:"email" json:"email"`
	Token                string `form:"token" json:"token"`
	Password             string `form:"password" json:"password"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation"`
}

func (r *ResetPasswordRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *ResetPasswordRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"email":                 "required|email",
		"token":                 "required",
		"password":              "required|min:8",
		"password_confirmation": "required",
	}
}

func (r *ResetPasswordRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"email.required":                 "Email is required.",
		"email.email":                    "Please provide a valid email address.",
		"token.required":                 "Reset token is required.",
		"password.required":              "Password is required.",
		"password.min":                   "Password must be at least 8 characters.",
		"password_confirmation.required": "Please confirm your password.",
		"password_confirmation.same":     "Passwords do not match.",
	}
}

// ResetPassword handles the password reset logic
func (r *AuthController) ResetPassword(ctx http.Context) http.Response {

	var req ResetPasswordRequest
	errors, err := ctx.Request().ValidateRequest(&req)
	if err != nil {
		return ctx.Response().Status(http.StatusInternalServerError).Json(http.Json{
			"message": "Error validating request: " + err.Error(),
		})
	}
	if errors != nil {
		return ctx.Response().Status(http.StatusUnprocessableEntity).Json(errors.All())
	}

	// Manual password confirmation check
	if req.Password != req.PasswordConfirmation {
		return ctx.Response().Status(http.StatusUnprocessableEntity).Json(http.Json{
			"password_confirmation": "Passwords do not match.",
		})
	}

	// Check token in Redis
	cacheKey := fmt.Sprintf("reset_token:%s", req.Email)

	// Getting the token from cache
	cachedValue := facades.Cache().Get(cacheKey, nil)

	// Enhanced token validation
	if cachedValue == nil {
		return ctx.Response().Status(http.StatusUnprocessableEntity).Json(http.Json{
			"token": "Invalid or expired reset token.",
		})
	}

	// Check if token exists and convert to string safely
	var redisToken string
	if tokenStr, ok := cachedValue.(string); ok {
		redisToken = tokenStr
	} else {
		return ctx.Response().Status(http.StatusUnprocessableEntity).Json(http.Json{
			"token": "Invalid or expired reset token.",
		})
	}

	// Validate token format (should be a valid UUID)
	if _, err := uuid.Parse(redisToken); err != nil {
		return ctx.Response().Status(http.StatusUnprocessableEntity).Json(http.Json{
			"token": "Invalid token format.",
		})
	}

	if _, err := uuid.Parse(req.Token); err != nil {
		return ctx.Response().Status(http.StatusUnprocessableEntity).Json(http.Json{
			"token": "Invalid token format.",
		})
	}

	// Compare tokens
	if redisToken != req.Token {
		return ctx.Response().Status(http.StatusUnprocessableEntity).Json(http.Json{
			"token": "Invalid or expired reset token.",
		})
	}

	// Token is valid, delete it from Redis
	if ok := facades.Cache().Forget(cacheKey); !ok {
		facades.Log().Warning("Failed to delete reset token from cache")
	}

	// Find user by email
	var user models.User
	if err := facades.Orm().Query().Where("email", req.Email).First(&user); err != nil {
		return ctx.Response().Status(http.StatusUnprocessableEntity).Json(http.Json{
			"email": "Invalid email address.",
		})
	}

	// Update password
	hashedPassword, err := facades.Hash().Make(req.Password)
	if err != nil {
		return ctx.Response().Status(http.StatusInternalServerError).Json(http.Json{
			"message": "Failed to hash password.",
		})
	}

	// Use direct update instead of Save to avoid potential model validation issues
	_, err = facades.Orm().Query().Model(&models.User{}).Where("email = ?", req.Email).Update(map[string]interface{}{
		"password": hashedPassword,
	})
	if err != nil {
		return ctx.Response().Status(http.StatusInternalServerError).Json(http.Json{
			"message": "Failed to update password.",
		})
	}

	// Redirect to success page instead of returning JSON
	return ctx.Response().Redirect(http.StatusSeeOther, "/reset-password-success")
}
