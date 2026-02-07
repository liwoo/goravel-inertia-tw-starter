package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// ApplicationUpdateRequest handles application update validation
type ApplicationUpdateRequest struct {
	SME            *string `form:"sme" json:"sme"`
	RegistrantName *string `form:"registrant_name" json:"registrant_name"`
	Email          *string `form:"email" json:"email"`
	Phone          *string `form:"phone" json:"phone"`
	Status         *string `form:"status" json:"status"`
	ID             uint    `form:"-" json:"-"` // Set by controller
}

// Rules defines validation rules for application updates
func (r *ApplicationUpdateRequest) Rules(ctx http.Context) map[string]string {
	rules := map[string]string{}

	// Only validate SME if provided
	if r.SME != nil {
		rules["sme"] = "required|string|max:255"
	}
	// Only validate RegistrantName if provided
	if r.RegistrantName != nil {
		rules["registrant_name"] = "required|string|max:255"
	}
	// Only validate Email if provided
	if r.Email != nil {
		rules["email"] = "required|string|max:255"
	}
	// Only validate Phone if provided
	if r.Phone != nil {
		rules["phone"] = "required|string|max:255"
	}
	// Only validate Status if provided
	if r.Status != nil {
		rules["status"] = "required|string|max:255"
	}

	// If no rules were added, add a dummy rule to prevent empty rules error
	if len(rules) == 0 {
		rules["_at_least_one_field"] = "sometimes"
	}

	return rules
}

// Messages defines custom validation messages for updates
func (r *ApplicationUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"sme.required":             "Organization name is required",
		"sme.max":                  "Organization name cannot exceed 255 characters",
		"registrant_name.required": "Registrant Name is required",
		"registrant_name.max":      "Registrant Name cannot exceed 255 characters",
		"email.required":           "Email is required",
		"email.max":                "Email cannot exceed 255 characters",
		"phone.required":           "Phone is required",
		"phone.max":                "Phone cannot exceed 255 characters",
		"status.required":          "Status is required",
		"status.max":               "Status cannot exceed 255 characters",
	}
}

// Attributes defines custom attribute names for updates
func (r *ApplicationUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
	}
}

// Authorize determines if the user is authorized to update this application
func (r *ApplicationUpdateRequest) Authorize(ctx http.Context) error {
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *ApplicationUpdateRequest) PrepareForValidation(ctx http.Context) error {
	return nil
}

// PassedValidation is called after validation passes
func (r *ApplicationUpdateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// ToUpdateData converts the request to update data map
func (r *ApplicationUpdateRequest) ToUpdateData() map[string]interface{} {
	data := map[string]interface{}{}

	// Only include SME if provided
	if r.SME != nil {
		data["sme"] = *r.SME
	}
	// Only include RegistrantName if provided
	if r.RegistrantName != nil {
		data["registrant_name"] = *r.RegistrantName
	}
	// Only include Email if provided
	if r.Email != nil {
		data["email"] = *r.Email
	}
	// Only include Phone if provided
	if r.Phone != nil {
		data["phone"] = *r.Phone
	}
	// Only include Status if provided
	if r.Status != nil {
		data["status"] = *r.Status
	}

	return data
}

// GetResourceID returns the resource ID for update
func (r *ApplicationUpdateRequest) GetResourceID() interface{} {
	return r.ID
}
