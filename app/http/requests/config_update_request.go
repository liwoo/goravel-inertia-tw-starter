package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// ConfigUpdateRequest handles config update validation
type ConfigUpdateRequest struct {
	Name        *string      `form:"name" json:"name"`
	ConfigType  *ConfigTypes `form:"config_type" json:"config_type"`
	Description *string      `form:"description" json:"description"`
	ID          uint         `form:"-" json:"-"` // Set by controller
}

// Rules defines validation rules for config updates
func (r *ConfigUpdateRequest) Rules(ctx http.Context) map[string]string {
	rules := map[string]string{}

	// Only validate Name if provided
	if r.Name != nil {
		rules["name"] = "required|max_len:255"
	}
	// Only validate ConfigType if provided
	if r.ConfigType != nil {
		rules["config_type"] = "required|max_len:255"
	}

	return rules
}

// Messages defines custom validation messages for updates
func (r *ConfigUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required":        "Name is required",
		"name.max":             "Name cannot exceed 255 characters",
		"config_type.required": "Config Type is required",
		"config_type.max":      "Config Type cannot exceed 255 characters",
	}
}

// Attributes defines custom attribute names for updates
func (r *ConfigUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to update this config
func (r *ConfigUpdateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user can update this specific config
	// return facades.Gate().Allows("update.configs", config)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *ConfigUpdateRequest) PrepareForValidation(ctx http.Context) error {
	// TODO: Add data preparation logic for updates
	// Example: Normalize data if provided
	return nil
}

// PassedValidation is called after validation passes
func (r *ConfigUpdateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// ToUpdateData converts the request to update data map
func (r *ConfigUpdateRequest) ToUpdateData() map[string]interface{} {
	data := map[string]interface{}{}

	// Only include Name if provided
	if r.Name != nil {
		data["name"] = *r.Name
	}
	// Only include ConfigType if provided
	if r.ConfigType != nil {
		data["config_type"] = *r.ConfigType
	}
	// Only include Description if provided
	if r.Description != nil {
		data["description"] = *r.Description
	}

	return data
}

// GetResourceID returns the resource ID for update
func (r *ConfigUpdateRequest) GetResourceID() interface{} {
	return r.ID
}
