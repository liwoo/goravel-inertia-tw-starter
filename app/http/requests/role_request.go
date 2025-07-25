package requests

import (
	"fmt"
	
	"github.com/goravel/framework/contracts/http"
	"players/app/contracts"
)

// RoleCreateRequest represents the request structure for creating a new role
type RoleCreateRequest struct {
	Name        string   `form:"name" json:"name"`
	Description string   `form:"description" json:"description"`
	Permissions []string `form:"permissions" json:"permissions"`
}

// Rules defines validation rules for role creation
func (r *RoleCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":        fmt.Sprintf("%s|string|%s|%s|unique:roles,name", contracts.Required, fmt.Sprintf(contracts.MinLength, 3), fmt.Sprintf(contracts.MaxLength, 100)),
		"description": fmt.Sprintf("string|%s", fmt.Sprintf(contracts.MaxLength, 255)),
		"permissions": contracts.Array,
	}
}

// Messages defines custom validation messages
func (r *RoleCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required": "Role name is required",
		"name.min":      "Role name must be at least 3 characters",
		"name.max":      "Role name cannot exceed 100 characters",
		"name.unique":   "A role with this name already exists",
	}
}

// Attributes defines friendly names for form fields
func (r *RoleCreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"name":        "role name",
		"description": "role description",
		"permissions": "permissions",
	}
}

// Authorize checks if the user is authorized to make this request
func (r *RoleCreateRequest) Authorize(ctx http.Context) error {
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *RoleCreateRequest) PrepareForValidation(ctx http.Context) error {
	// Could normalize role name, etc.
	return nil
}

// PassedValidation is called after validation passes
func (r *RoleCreateRequest) PassedValidation(ctx http.Context) error {
	// Could log the validation success, perform additional checks, etc.
	return nil
}

// ToCreateData converts the request to a map for database insertion
func (r *RoleCreateRequest) ToCreateData() map[string]interface{} {
	return map[string]interface{}{
		"name":        r.Name,
		"description": r.Description,
	}
}

// RoleUpdateRequest represents the request structure for updating a role
type RoleUpdateRequest struct {
	ID          uint     `route:"id"`
	Name        string   `form:"name" json:"name"`
	Description string   `form:"description" json:"description"`
	Permissions []string `form:"permissions" json:"permissions"`
}

// Rules defines validation rules for role update
func (r *RoleUpdateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":        fmt.Sprintf("%s|string|%s|%s|unique:roles,name,%d", contracts.Required, fmt.Sprintf(contracts.MinLength, 3), fmt.Sprintf(contracts.MaxLength, 100), r.ID),
		"description": fmt.Sprintf("string|%s", fmt.Sprintf(contracts.MaxLength, 255)),
		"permissions": contracts.Array,
	}
}

// Messages defines custom validation messages
func (r *RoleUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required": "Role name is required",
		"name.min":      "Role name must be at least 3 characters",
		"name.max":      "Role name cannot exceed 100 characters",
		"name.unique":   "A role with this name already exists",
	}
}

// Attributes defines friendly names for form fields
func (r *RoleUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"name":        "role name",
		"description": "role description",
		"permissions": "permissions",
	}
}

// Authorize checks if the user is authorized to make this request
func (r *RoleUpdateRequest) Authorize(ctx http.Context) error {
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *RoleUpdateRequest) PrepareForValidation(ctx http.Context) error {
	// Could normalize role name, etc.
	return nil
}

// PassedValidation is called after validation passes
func (r *RoleUpdateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// ToUpdateData converts the request to a map for database update
func (r *RoleUpdateRequest) ToUpdateData() map[string]interface{} {
	return map[string]interface{}{
		"name":        r.Name,
		"description": r.Description,
	}
}

// GetResourceID returns the resource ID for update
func (r *RoleUpdateRequest) GetResourceID() interface{} {
	return r.ID
}