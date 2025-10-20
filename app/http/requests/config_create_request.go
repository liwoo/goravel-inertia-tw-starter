package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// ConfigCreateRequest handles config creation validation
type ConfigCreateRequest struct {
	Name        string      `form:"name" json:"name"`
	ConfigType  ConfigTypes `form:"config_type" json:"config_type"`
	Description *string     `form:"description" json:"description"`
}

// Rules defines validation rules for config creation
func (r *ConfigCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":        "required|max_len:255",
		"config_type": "required|max_len:255",
	}
}

// Messages defines custom validation messages
func (r *ConfigCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required":        "Name is required",
		"name.max":             "Name cannot exceed 255 characters",
		"config_type.required": "Config Type is required",
		"config_type.max":      "Config Type cannot exceed 255 characters",
	}
}

// Attributes defines custom attribute names
func (r *ConfigCreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to make this request
func (r *ConfigCreateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user has permission to create config
	// return facades.Gate().Allows("create.configs", ctx)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *ConfigCreateRequest) PrepareForValidation(ctx http.Context) error {
	// TODO: Add data preparation logic
	// Example: Normalize, trim, or set default values
	return nil
}

// PassedValidation is called after validation passes
func (r *ConfigCreateRequest) PassedValidation(ctx http.Context) error {
	// TODO: Add post-validation logic if needed
	return nil
}

// ToCreateData converts the request to create data map
func (r *ConfigCreateRequest) ToCreateData() map[string]interface{} {
	data := map[string]interface{}{
		"name":        r.Name,
		"config_type": r.ConfigType,
		"description": r.Description,
	}

	return data
}
