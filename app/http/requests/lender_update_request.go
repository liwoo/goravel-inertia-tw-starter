package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// LenderUpdateRequest handles lender update validation
type LenderUpdateRequest struct {
	Name    *string `form:"name" json:"name"`
	Email   *string `form:"email" json:"email"`
	Phone   *string `form:"phone" json:"phone"`
	Address *string `form:"address" json:"address"`
	Gender  *string `form:"gender" json:"gender"`
	ID      uint    `form:"-" json:"-"` // Set by controller
}

// Rules defines validation rules for lender updates
func (r *LenderUpdateRequest) Rules(ctx http.Context) map[string]string {
	rules := map[string]string{}

	// Only validate Name if provided
	if r.Name != nil {
		rules["name"] = "required|string|max:255"
	}
	// Only validate Email if provided
	if r.Email != nil {
		rules["email"] = "required|string|max:255"
	}

	// If no rules were added, add a dummy rule to prevent empty rules error
	if len(rules) == 0 {
		rules["_at_least_one_field"] = "sometimes"
	}

	return rules
}

// Messages defines custom validation messages for updates
func (r *LenderUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required":  "Name is required",
		"name.max":       "Name cannot exceed 255 characters",
		"email.required": "Email is required",
		"email.max":      "Email cannot exceed 255 characters",
	}
}

// Attributes defines custom attribute names for updates
func (r *LenderUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to update this lender
func (r *LenderUpdateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user can update this specific lender
	// return facades.Gate().Allows("update.lenders", lender)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *LenderUpdateRequest) PrepareForValidation(ctx http.Context) error {
	// TODO: Add data preparation logic for updates
	// Example: Normalize data if provided
	return nil
}

// PassedValidation is called after validation passes
func (r *LenderUpdateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// ToUpdateData converts the request to update data map
func (r *LenderUpdateRequest) ToUpdateData() map[string]interface{} {
	data := map[string]interface{}{}

	// Only include Name if provided
	if r.Name != nil {
		data["name"] = *r.Name
	}
	// Only include Email if provided
	if r.Email != nil {
		data["email"] = *r.Email
	}
	// Only include Phone if provided
	if r.Phone != nil {
		data["phone"] = *r.Phone
	}
	// Only include Address if provided
	if r.Address != nil {
		data["address"] = *r.Address
	}
	// Only include Gender if provided
	if r.Gender != nil {
		data["gender"] = *r.Gender
	}

	return data
}

// GetResourceID returns the resource ID for update
func (r *LenderUpdateRequest) GetResourceID() interface{} {
	return r.ID
}
