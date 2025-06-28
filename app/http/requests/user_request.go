package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

// UserCreateRequest handles validation for creating users
type UserCreateRequest struct {
	Name         string `form:"name" json:"name"`
	Email        string `form:"email" json:"email"`
	Password     string `form:"password" json:"password"`
	IsActive     bool   `form:"is_active" json:"is_active"`
	IsSuperAdmin bool   `form:"is_super_admin" json:"is_super_admin"`
	RoleID       uint   `form:"role_id" json:"role_id"`
}

// Authorize determines if the user can make this request
func (r *UserCreateRequest) Authorize(ctx http.Context) error {
	// Authorization is handled in the controller
	return nil
}

// Rules returns the validation rules for the request
func (r *UserCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":           "required|string|max:255|min:2",
		"email":          "required|email|max:255",
		"password":       "required|string|min:8",
		"is_active":      "boolean",
		"is_super_admin": "boolean",
		"role_id":        "numeric",
	}
}

// Messages returns custom validation messages
func (r *UserCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required":     "User name is required",
		"name.min":          "User name must be at least 2 characters",
		"name.max":          "User name cannot exceed 255 characters",
		"email.required":    "Email address is required",
		"email.email":       "Invalid email format",
		"email.max":         "Email cannot exceed 255 characters",
		"password.required": "Password is required",
		"password.min":      "Password must be at least 8 characters",
		"role_id.numeric":   "Invalid role ID",
	}
}

// Attributes returns custom attribute names for validation
func (r *UserCreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"name":           "Full Name",
		"email":          "Email Address",
		"password":       "Password",
		"is_active":      "Active Status",
		"is_super_admin": "Super Admin Status",
		"role_id":        "Role",
	}
}

// PrepareForValidation allows you to modify the data before validation
func (r *UserCreateRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	// Set default values if not provided
	if _, exists := data.Get("is_active"); !exists {
		data.Set("is_active", true)
	}
	if _, exists := data.Get("is_super_admin"); !exists {
		data.Set("is_super_admin", false)
	}
	return nil
}

// ToCreateData converts the request to data suitable for the service
func (r *UserCreateRequest) ToCreateData() map[string]interface{} {
	data := map[string]interface{}{
		"name":           r.Name,
		"email":          r.Email,
		"password":       r.Password,
		"is_active":      r.IsActive,
		"is_super_admin": r.IsSuperAdmin,
	}
	
	if r.RoleID > 0 {
		data["role_id"] = float64(r.RoleID)
	}
	
	return data
}

// UserUpdateRequest handles validation for updating users
type UserUpdateRequest struct {
	ID           uint   `form:"id" json:"id"`
	Name         string `form:"name" json:"name"`
	Email        string `form:"email" json:"email"`
	Password     string `form:"password" json:"password"`
	IsActive     bool   `form:"is_active" json:"is_active"`
	IsSuperAdmin bool   `form:"is_super_admin" json:"is_super_admin"`
	RoleID       uint   `form:"role_id" json:"role_id"`
}

// Authorize determines if the user can make this request
func (r *UserUpdateRequest) Authorize(ctx http.Context) error {
	// Authorization is handled in the controller
	return nil
}

// Rules returns the validation rules for the request
func (r *UserUpdateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":           "string|max:255|min:2",
		"email":          "email|max:255",
		"password":       "string|min:8",
		"is_active":      "boolean",
		"is_super_admin": "boolean",
		"role_id":        "numeric",
	}
}

// Messages returns custom validation messages
func (r *UserUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.min":        "User name must be at least 2 characters",
		"name.max":        "User name cannot exceed 255 characters",
		"email.email":     "Invalid email format",
		"email.max":       "Email cannot exceed 255 characters",
		"password.min":    "Password must be at least 8 characters",
		"role_id.numeric": "Invalid role ID",
	}
}

// Attributes returns custom attribute names for validation
func (r *UserUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"name":           "Full Name",
		"email":          "Email Address",
		"password":       "Password",
		"is_active":      "Active Status",
		"is_super_admin": "Super Admin Status",
		"role_id":        "Role",
	}
}

// PrepareForValidation allows you to modify the data before validation
func (r *UserUpdateRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	// Nothing to prepare for updates
	return nil
}

// ToUpdateData converts the request to data suitable for the service
func (r *UserUpdateRequest) ToUpdateData() map[string]interface{} {
	data := make(map[string]interface{})
	
	if r.Name != "" {
		data["name"] = r.Name
	}
	if r.Email != "" {
		data["email"] = r.Email
	}
	if r.Password != "" {
		data["password"] = r.Password
	}
	// Always include boolean fields for updates
	data["is_active"] = r.IsActive
	data["is_super_admin"] = r.IsSuperAdmin
	
	if r.RoleID > 0 {
		data["role_id"] = float64(r.RoleID)
	}
	
	return data
}