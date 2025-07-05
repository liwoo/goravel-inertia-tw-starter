package auth

import (
	"context"
	"fmt"
	"players/app/models" // Assuming your User model is here
	"players/app/services"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
)

type AuthController struct {
	// Dependencies can be injected here
	passwordAttemptService *services.PasswordAttemptService
}

func NewAuthController() *AuthController {
	return &AuthController{
		passwordAttemptService: services.NewPasswordAttemptService(),
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

	// Check if user is locked before attempting login
	attemptResult, err := r.passwordAttemptService.CheckUserStatus(context.Background(), loginRequest.Email)
	if err != nil {
		facades.Log().Error("Error checking user status: " + err.Error())
		ctx.Request().Session().Flash("error", "Authentication service error")
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	if attemptResult.IsLocked {
		lockMessage := "Account temporarily locked due to too many failed attempts."
		if attemptResult.LockExpiresAt != nil {
			lockMessage = fmt.Sprintf("Account locked until %s due to too many failed attempts.",
				attemptResult.LockExpiresAt.Format("3:04 PM"))
		}
		ctx.Request().Session().Flash("error", lockMessage)
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	var user models.User
	// Find user by email
	if err := facades.Orm().Query().Where("email", loginRequest.Email).First(&user); err != nil {

		// Record failed attempt for non-existent user too (security measure)
		if _, attemptErr := r.passwordAttemptService.RecordFailedAttempt(context.Background(), loginRequest.Email); attemptErr != nil {
			facades.Log().Error("Error recording failed attempt: " + attemptErr.Error())
		}
		// Flash field-specific error
		ctx.Request().Session().Flash("errors", map[string]interface{}{
			"email": "Invalid credentials (Email not found)",
		})
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	// Check password
	if !facades.Hash().Check(loginRequest.Password, user.Password) {
		// Record failed attempt
		attemptResult, err := r.passwordAttemptService.RecordFailedAttempt(context.Background(), loginRequest.Email)
		if err != nil {
			facades.Log().Error("Error recording failed attempt: " + err.Error())
		}

		// Prepare error message based on attempt result
		errorMessage := "Password is incorrect"
		if attemptResult != nil {
			if attemptResult.IsLocked {
				lockMessage := "Account temporarily locked due to too many failed attempts."
				if attemptResult.LockExpiresAt != nil {
					lockMessage = fmt.Sprintf("Account locked until %s due to too many failed attempts.",
						attemptResult.LockExpiresAt.Format("3:04 PM"))
				}
				ctx.Request().Session().Flash("error", lockMessage)
				return ctx.Response().Redirect(http.StatusFound, "/login")
			} else if attemptResult.ShouldWarn {
				errorMessage = fmt.Sprintf("Invalid password. Warning: %d more failed attempts will lock your account for 1 hour.",
					attemptResult.RemainingAttempts)
			}
		}

		// Flash field-specific error
		ctx.Request().Session().Flash("errors", map[string]interface{}{
			"password": errorMessage,
		})
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	// Clear any failed attempts on successful login
	if err := r.passwordAttemptService.ClearAttempts(context.Background(), loginRequest.Email); err != nil {
		facades.Log().Error("Error clearing failed attempts: " + err.Error())
	}

	// Log the user in and get the token
	token, err := facades.Auth(ctx).Login(&user)
	if err != nil {
		ctx.Request().Session().Flash("errors", map[string]interface{}{
			"general": "Error during login: " + err.Error(),
		})
		return ctx.Response().Redirect(http.StatusFound, "/login")
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
