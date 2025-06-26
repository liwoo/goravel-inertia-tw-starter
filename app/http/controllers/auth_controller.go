package controllers

import (
	"fmt"
	"players/app/http/inertia"
	"players/app/http/middleware"
	"players/app/models" // Assuming your User model is here
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
)

type AuthController struct {
	passwordAttemptTracker *middleware.PasswordAttemptTracker
}

func NewAuthController() *AuthController {
	return &AuthController{
		passwordAttemptTracker: middleware.NewPasswordAttemptTracker(),
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
		// For Inertia, return validation error
		return inertia.Render(ctx, "auth/Login", map[string]interface{}{
			"errors": map[string]string{"general": "Error validating request: " + err.Error()},
		})
	}
	if errors != nil {
		// Return the login page with validation errors
		return inertia.Render(ctx, "auth/Login", map[string]interface{}{
			"errors": errors.All(),
		})
	}

	// Check if user is locked before attempting login
	status, err := r.passwordAttemptTracker.CheckUserStatus(ctx.Context(), loginRequest.Email)
	if err != nil {
		facades.Log().Error("Error checking user status: " + err.Error())
		return inertia.Render(ctx, "auth/Login", map[string]interface{}{
			"errors": map[string]string{"general": "Internal server error. Please try again later."},
		})
	}

	if status.IsLocked {
		lockMessage := "Account temporarily locked due to too many failed login attempts."
		if status.LockExpiresAt != nil {
			lockMessage = fmt.Sprintf("Account locked until %s due to too many failed login attempts.",
				status.LockExpiresAt.Format("2006-01-02 15:04:05"))
		}

		return inertia.Render(ctx, "auth/Login", map[string]interface{}{
			"errors": map[string]string{
				"general": lockMessage,
			},
		})
	}

	var user models.User
	// Find user by email
	if err := facades.Orm().Query().Where("email", loginRequest.Email).First(&user); err != nil {
		// Record failed attempt for non-existent email too (to prevent email enumeration)
		if result, attemptErr := r.passwordAttemptTracker.RecordFailedAttempt(ctx.Context(), loginRequest.Email); attemptErr != nil {
			facades.Log().Error("Error recording failed attempt: " + attemptErr.Error())
		} else {
			// Add warning message if needed
			if result.ShouldWarn {
				return inertia.Render(ctx, "auth/Login", map[string]interface{}{
					"errors": map[string]string{
						"email":   "Invalid credentials",
						"warning": fmt.Sprintf("Warning: %d more failed attempts will result in account lockout.", result.RemainingAttempts),
					},
				})
			}
			if result.IsLocked {
				lockMessage := "Account temporarily locked due to too many failed login attempts."
				if result.LockExpiresAt != nil {
					lockMessage = fmt.Sprintf("Account locked until %s due to too many failed login attempts.",
						result.LockExpiresAt.Format("2006-01-02 15:04:05"))
				}
				return inertia.Render(ctx, "auth/Login", map[string]interface{}{
					"errors": map[string]string{
						"general": lockMessage,
					},
				})
			}
		}

		return inertia.Render(ctx, "auth/Login", map[string]interface{}{
			"errors": map[string]string{"email": "Invalid credentials"},
		})
	}

	// Check password
	if !facades.Hash().Check(loginRequest.Password, user.Password) {
		// Record failed attempt
		if result, attemptErr := r.passwordAttemptTracker.RecordFailedAttempt(ctx.Context(), loginRequest.Email); attemptErr != nil {
			facades.Log().Error("Error recording failed attempt: " + attemptErr.Error())
		} else {
			// Add warning message if needed
			if result.ShouldWarn {
				return inertia.Render(ctx, "auth/Login", map[string]interface{}{
					"errors": map[string]string{
						"password": "Invalid credentials. Password mismatch.",
						"warning":  fmt.Sprintf("Warning: %d more failed attempts will result in account lockout.", result.RemainingAttempts),
					},
				})
			}
			if result.IsLocked {
				lockMessage := "Account temporarily locked due to too many failed login attempts."
				if result.LockExpiresAt != nil {
					lockMessage = fmt.Sprintf("Account locked until %s due to too many failed login attempts.",
						result.LockExpiresAt.Format("2006-01-02 15:04:05"))
				}
				return inertia.Render(ctx, "auth/Login", map[string]interface{}{
					"errors": map[string]string{
						"general": lockMessage,
					},
				})
			}
		}

		return inertia.Render(ctx, "auth/Login", map[string]interface{}{
			"errors": map[string]string{"password": "Invalid credentials. Password mismatch."},
		})
	}

	// Clear failed attempts on successful login
	if err := r.passwordAttemptTracker.ClearAttempts(ctx.Context(), loginRequest.Email); err != nil {
		facades.Log().Error("Error clearing failed attempts: " + err.Error())
		// Don't fail the login for this, just log it
	}

	// Log the user in and get the token
	token, err := facades.Auth(ctx).Login(&user)
	if err != nil {
		return inertia.Render(ctx, "auth/Login", map[string]interface{}{
			"errors": map[string]string{"general": "Error during login: " + err.Error()},
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
