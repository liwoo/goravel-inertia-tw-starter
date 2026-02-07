package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// ApplicationCreateRequest handles application creation validation
type ApplicationCreateRequest struct {
	SME            string `form:"sme" json:"sme"`
	RegistrantName string `form:"registrant_name" json:"registrant_name"`
	Email          string `form:"email" json:"email"`
	Phone          string `form:"phone" json:"phone"`
}

// Rules defines validation rules for application creation
func (r *ApplicationCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"sme":             "required|string|max_len:255",
		"registrant_name": "required|string|max_len:255",
		"email":           "required|string|max_len:255",
		"phone":           "required|string|max_len:255",
	}
}

// Messages defines custom validation messages
func (r *ApplicationCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"sme.required":             "Organization name is required",
		"sme.max_len":              "Organization name cannot exceed 255 characters",
		"registrant_name.required": "Registrant Name is required",
		"registrant_name.max_len":  "Registrant Name cannot exceed 255 characters",
		"email.required":           "Email is required",
		"email.max_len":            "Email cannot exceed 255 characters",
		"phone.required":           "Phone is required",
		"phone.max_len":            "Phone cannot exceed 255 characters",
	}
}

// Attributes defines custom attribute names
func (r *ApplicationCreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
	}
}

// Authorize determines if the user is authorized to make this request
func (r *ApplicationCreateRequest) Authorize(ctx http.Context) error {
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *ApplicationCreateRequest) PrepareForValidation(ctx http.Context) error {
	return nil
}

// PassedValidation is called after validation passes
func (r *ApplicationCreateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// ToCreateData converts the request to create data map
func (r *ApplicationCreateRequest) ToCreateData() map[string]interface{} {
	data := map[string]interface{}{
		"sme":             r.SME,
		"registrant_name": r.RegistrantName,
		"email":           r.Email,
		"phone":           r.Phone,
	}

	return data
}
