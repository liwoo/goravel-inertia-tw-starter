package controllers

import (
	"fmt"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
	"players/app/models" // Assuming your User model is here
	"time"
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
