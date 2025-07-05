package requests

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

// ProfilePasswordUpdateRequest handles password updates for user's own profile
type ProfilePasswordUpdateRequest struct {
	CurrentPassword      string `form:"current_password" json:"current_password"`
	Password             string `form:"password" json:"password"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation"`
}

// Authorize determines if the user can make this request
func (r *ProfilePasswordUpdateRequest) Authorize(ctx http.Context) error {
	// User must be authenticated to update their own profile
	// Additional authorization handled in controller
	return nil
}

// Rules returns the validation rules for the request
func (r *ProfilePasswordUpdateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"current_password":      "required|string",
		"password":              "required|string|min:8",
		"password_confirmation": "required|string",
	}
}

// Messages returns custom validation messages
func (r *ProfilePasswordUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"current_password.required":      "Current password is required",
		"password.required":              "New password is required",
		"password.min":                   "Password must be at least 8 characters",
		"password_confirmation.required": "Password confirmation is required",
		"password_confirmation.same":     "Password confirmation must match the new password",
	}
}

// Attributes returns custom attribute names for validation
func (r *ProfilePasswordUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"current_password":      "Current Password",
		"password":              "New Password",
		"password_confirmation": "Password Confirmation",
	}
}

// PrepareForValidation allows you to modify the data before validation
func (r *ProfilePasswordUpdateRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	// Could add password strength validation or other preprocessing here
	return nil
}

// ValidatePasswordConfirmation checks if password and password_confirmation match
func (r *ProfilePasswordUpdateRequest) ValidatePasswordConfirmation() error {
	if r.Password != r.PasswordConfirmation {
		return fmt.Errorf("password confirmation does not match the new password")
	}
	return nil
}

// ValidatePasswordRequirements checks if password meets minimum requirements
func (r *ProfilePasswordUpdateRequest) ValidatePasswordRequirements() error {
	if len(r.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	return nil
}
