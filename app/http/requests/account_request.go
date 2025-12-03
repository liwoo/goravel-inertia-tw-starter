package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// UpdateProfileRequest handles profile update validation
// Note: Email is intentionally excluded - users cannot change their email address
type UpdateProfileRequest struct {
	Name string `form:"name" json:"name"`
}

// Rules defines validation rules for profile update
func (r *UpdateProfileRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name": "required|min_len:2|max_len:255",
	}
}

// Messages defines custom validation messages
func (r *UpdateProfileRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required": "Name is required",
		"name.min_len":  "Name must be at least 2 characters",
		"name.max_len":  "Name cannot exceed 255 characters",
	}
}

// Attributes defines custom attribute names
func (r *UpdateProfileRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"name": "name",
	}
}

// Authorize determines if the user is authorized to make this request
func (r *UpdateProfileRequest) Authorize(ctx http.Context) error {
	// All authenticated users can update their own profile
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *UpdateProfileRequest) PrepareForValidation(ctx http.Context) error {
	return nil
}

// PassedValidation is called after validation passes
func (r *UpdateProfileRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// ToUpdateData converts the request to update data map
// Note: Only name is included - email cannot be changed
func (r *UpdateProfileRequest) ToUpdateData() map[string]interface{} {
	return map[string]interface{}{
		"name": r.Name,
	}
}

// ChangePasswordRequest handles password change validation
type ChangePasswordRequest struct {
	CurrentPassword string `form:"current_password" json:"current_password"`
	NewPassword     string `form:"new_password" json:"new_password"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password"`
}

// Rules defines validation rules for password change
func (r *ChangePasswordRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"current_password": "required",
		"new_password":     "required|min_len:8|max_len:128",
		"confirm_password": "required",
	}
}

// Messages defines custom validation messages
func (r *ChangePasswordRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"current_password.required": "Current password is required",
		"new_password.required":     "New password is required",
		"new_password.min_len":      "New password must be at least 8 characters",
		"new_password.max_len":      "New password cannot exceed 128 characters",
		"confirm_password.required": "Please confirm your new password",
	}
}

// Attributes defines custom attribute names
func (r *ChangePasswordRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"current_password": "current password",
		"new_password":     "new password",
		"confirm_password": "password confirmation",
	}
}

// Authorize determines if the user is authorized to make this request
func (r *ChangePasswordRequest) Authorize(ctx http.Context) error {
	// All authenticated users can change their own password
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *ChangePasswordRequest) PrepareForValidation(ctx http.Context) error {
	return nil
}

// PassedValidation is called after validation passes
func (r *ChangePasswordRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// ValidatePasswordMatch checks if new password and confirm password match
func (r *ChangePasswordRequest) ValidatePasswordMatch() bool {
	return r.NewPassword == r.ConfirmPassword
}
