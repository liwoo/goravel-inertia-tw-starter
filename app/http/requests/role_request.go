package requests

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
)

// RoleCreateRequest represents the request structure for creating a new role
type RoleCreateRequest struct {
	Name        string   `form:"name" json:"name"`
	Description string   `form:"description" json:"description"`
	Slug        string   `form:"slug" json:"slug"`
	Level       int      `form:"level" json:"level"`
	Permissions []string `form:"permissions" json:"permissions"`
}

// Rules defines validation rules for role creation
func (r *RoleCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":        "required|string|min_len:3|max_len:255|unique:roles,name",
		"description": "string|max_len:1000",
		"slug":        "required|string|min_len:3|max_len:255|unique:roles,slug",
		"level":       "required|numeric|min:1|max:100",
		"permissions": "array",
	}
}

// Messages defines custom validation messages
func (r *RoleCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required":       "Role name is required",
		"name.min_len":        "Role name must be at least 3 characters",
		"name.max_len":        "Role name cannot exceed 255 characters",
		"name.unique":         "A role with this name already exists",
		"description.max_len": "Description cannot exceed 1000 characters",
		"slug.required":       "Role slug is required",
		"slug.min_len":        "Role slug must be at least 3 characters",
		"slug.max_len":        "Role slug cannot exceed 255 characters",
		"slug.unique":         "A role with this slug already exists",
		"level.required":      "Role level is required",
		"level.integer":       "Role level must be an integer",
		"level.min":           "Role level must be at least 1",
		"level.max":           "Role level cannot exceed 100",
	}
}

// Attributes defines friendly names for form fields
func (r *RoleCreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"name":        "role name",
		"description": "role description",
		"permissions": "permissions",
		"slug":        "role slug",
		"level":       "role level",
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
		"slug":        r.Slug,
		"level":       r.Level,
	}
}

// RoleUpdateRequest represents the request structure for updating a role
type RoleUpdateRequest struct {
	ID          uint     `route:"id"`
	Name        string   `form:"name" json:"name"`
	Description string   `form:"description" json:"description"`
	Slug        string   `form:"slug" json:"slug"`
	Level       int      `form:"level" json:"level"`
	Permissions []string `form:"permissions" json:"permissions"`
}

// Rules defines validation rules for role update
func (r *RoleUpdateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":        fmt.Sprintf("required|string|min_len:3|max_len:255|unique:roles,name,%d", r.ID),
		"description": "string|max_len:1000",
		"slug":        "required|string|min_len:3|max_len:255",
		"permissions": "array",
	}
}

// Messages defines custom validation messages
func (r *RoleUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required":       "Role name is required",
		"name.min_len":        "Role name must be at least 3 characters",
		"name.max_len":        "Role name cannot exceed 255 characters",
		"name.unique":         "A role with this name already exists",
		"description.max_len": "Description cannot exceed 1000 characters",
		"slug.required":       "Role slug is required",
		"slug.min_len":        "Role slug must be at least 3 characters",
		"slug.max_len":        "Role slug cannot exceed 255 characters",
		"slug.unique":         "A role with this slug already exists",
		"level.required":      "Role level is required",
		"level.integer":       "Role level must be an integer",
		"level.min":           "Role level must be at least 1",
		"level.max":           "Role level cannot exceed 100",
	}
}

// Attributes defines friendly names for form fields
func (r *RoleUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"name":        "role name",
		"description": "role description",
		"permissions": "permissions",
		"slug":        "role slug",
		"level":       "role level",
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
		"slug":        r.Slug,
		"level":       r.Level,
		"permissions": r.Permissions,
	}
}

// GetResourceID returns the resource ID for update
func (r *RoleUpdateRequest) GetResourceID() interface{} {
	return r.ID
}
