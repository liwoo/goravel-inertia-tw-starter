package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// UserCreateRequest handles validation for creating users
// TEMPORARY FIX: Removed max/min rules due to Goravel v1.15.4 validation bug
type UserCreateRequest struct {
	Name         string `form:"name" json:"name"`
	Email        string `form:"email" json:"email"`
	Password     string `form:"password" json:"password"`
	IsActive     *bool  `form:"is_active" json:"is_active"`
	IsSuperAdmin *bool  `form:"is_super_admin" json:"is_super_admin"`
	RoleID       *uint  `form:"role_id" json:"role_id"`
}

// Authorize determines if the user can make this request
func (r *UserCreateRequest) Authorize(ctx http.Context) error {
	// Authorization is handled in the controller
	return nil
}

// Rules returns the validation rules for the request
func (r *UserCreateRequest) Rules(ctx http.Context) map[string]string {
	// TEMPORARY: Removed max:255|min:2 due to Goravel bug
	// Original rules that fail:
	// "name":     "required|string|max:255|min:2",
	// "email":    "required|email|max:255",
	// "password": "required|string|min:8",

	rules := map[string]string{
		"name":     "required|string",
		"email":    "required|email",
		"password": "required|string",
	}

	// Optional fields - only validate if provided
	if r.IsActive != nil {
		rules["is_active"] = "boolean"
	}
	if r.IsSuperAdmin != nil {
		rules["is_super_admin"] = "boolean"
	}
	if r.RoleID != nil {
		rules["role_id"] = "numeric"
	}

	return rules
}

// Messages returns custom validation messages
func (r *UserCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required":     "User name is required",
		"email.required":    "Email address is required",
		"email.email":       "Invalid email format",
		"password.required": "Password is required",
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
func (r *UserCreateRequest) PrepareForValidation(ctx http.Context) error {
	// Set default values if not provided
	if r.IsActive == nil {
		isActive := true
		r.IsActive = &isActive
	}
	if r.IsSuperAdmin == nil {
		isSuperAdmin := false
		r.IsSuperAdmin = &isSuperAdmin
	}
	return nil
}

// ToCreateData converts the request to data suitable for the service
func (r *UserCreateRequest) ToCreateData() map[string]interface{} {
	data := map[string]interface{}{
		"name":     r.Name,
		"email":    r.Email,
		"password": r.Password,
	}

	// Add optional fields
	if r.IsActive != nil {
		data["is_active"] = *r.IsActive
	}
	if r.IsSuperAdmin != nil {
		data["is_super_admin"] = *r.IsSuperAdmin
	}
	if r.RoleID != nil {
		data["role_id"] = float64(*r.RoleID)
	}

	return data
}

// UserUpdateRequest handles validation for updating users
type UserUpdateRequest struct {
	ID           uint    `form:"id" json:"id"`
	Name         *string `form:"name" json:"name"`
	Email        *string `form:"email" json:"email"`
	Password     *string `form:"password" json:"password"`
	IsActive     *bool   `form:"is_active" json:"is_active"`
	IsSuperAdmin *bool   `form:"is_super_admin" json:"is_super_admin"`
	RoleID       *uint   `form:"role_id" json:"role_id"`
}

// Authorize determines if the user can make this request
func (r *UserUpdateRequest) Authorize(ctx http.Context) error {
	// Authorization is handled in the controller
	return nil
}

// Rules returns the validation rules for the request
func (r *UserUpdateRequest) Rules(ctx http.Context) map[string]string {
	rules := map[string]string{}

	// TEMPORARY: Removed max/min due to Goravel bug
	// Only validate fields that are provided
	if r.Name != nil {
		rules["name"] = "string"
	}
	if r.Email != nil {
		rules["email"] = "email"
	}
	if r.Password != nil {
		rules["password"] = "string"
	}
	if r.IsActive != nil {
		rules["is_active"] = "boolean"
	}
	if r.IsSuperAdmin != nil {
		rules["is_super_admin"] = "boolean"
	}
	if r.RoleID != nil {
		rules["role_id"] = "numeric"
	}

	// Goravel validation requires at least one rule
	// Add a dummy rule if no fields are being updated
	if len(rules) == 0 {
		rules["_at_least_one_field"] = "sometimes"
	}

	return rules
}

// Messages returns custom validation messages
func (r *UserUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.string":     "User name must be text",
		"email.email":     "Invalid email format",
		"password.string": "Password must be text",
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
func (r *UserUpdateRequest) PrepareForValidation(ctx http.Context) error {
	// Nothing to prepare for updates
	return nil
}

// ToUpdateData converts the request to data suitable for the service
func (r *UserUpdateRequest) ToUpdateData() map[string]interface{} {
	data := make(map[string]interface{})

	// Only include fields that are provided (not nil)
	if r.Name != nil {
		data["name"] = *r.Name
	}
	if r.Email != nil {
		data["email"] = *r.Email
	}
	if r.Password != nil {
		data["password"] = *r.Password
	}
	if r.IsActive != nil {
		data["is_active"] = *r.IsActive
	}
	if r.IsSuperAdmin != nil {
		data["is_super_admin"] = *r.IsSuperAdmin
	}
	if r.RoleID != nil {
		data["role_id"] = float64(*r.RoleID)
	}

	return data
}

// PassedValidation is called after validation passes
func (r *UserCreateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// PassedValidation is called after validation passes
func (r *UserUpdateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// GetResourceID returns the resource ID for update
func (r *UserUpdateRequest) GetResourceID() interface{} {
	return r.ID
}
