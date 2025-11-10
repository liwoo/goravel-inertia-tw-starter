package requests

import (
	"github.com/goravel/framework/contracts/http"
	"strings"
)

// ConfigCreateRequest handles config creation validation
type ConfigCreateRequest struct {
	Name        string      `form:"name" json:"name"`
	Code        *string     `form:"code" json:"code"`
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
	// Auto-generate code if not provided
	if r.Code == nil || *r.Code == "" {
		generatedCode := generateCodeFromName(r.Name)
		r.Code = &generatedCode
	}
	return nil
}

// generateCodeFromName generates a code from the name
// If name has spaces, take first letter of each word
// Otherwise, take first 2 letters
func generateCodeFromName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}

	// Split by spaces
	words := strings.Fields(name)

	if len(words) > 1 {
		// Multiple words: take first letter of each word
		var code strings.Builder
		for _, word := range words {
			if len(word) > 0 {
				code.WriteRune([]rune(strings.ToUpper(word))[0])
			}
		}
		return code.String()
	}

	// Single word: take first 2 letters
	upperName := strings.ToUpper(name)
	if len(upperName) >= 2 {
		return upperName[:2]
	}
	return upperName
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
		"code":        r.Code,
		"config_type": r.ConfigType,
		"description": r.Description,
	}

	return data
}
