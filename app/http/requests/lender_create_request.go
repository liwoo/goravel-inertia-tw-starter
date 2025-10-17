package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// LenderCreateRequest handles lender creation validation
type LenderCreateRequest struct {
	Name    string  `form:"name" json:"name"`
	Email   string  `form:"email" json:"email"`
	Phone   *string `form:"phone" json:"phone"`
	Address *string `form:"address" json:"address"`
	Gender  *string `form:"gender" json:"gender"`
}

// Rules defines validation rules for lender creation
func (r *LenderCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":  "required",
		"email": "required",
	}
}

// Messages defines custom validation messages
func (r *LenderCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required":  "Name is required",
		"name.max":       "Name cannot exceed 255 characters",
		"email.required": "Email is required",
		"email.max":      "Email cannot exceed 255 characters",
	}
}

// Attributes defines custom attribute names
func (r *LenderCreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to make this request
func (r *LenderCreateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user has permission to create lender
	// return facades.Gate().Allows("create.lenders", ctx)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *LenderCreateRequest) PrepareForValidation(ctx http.Context) error {
	// TODO: Add data preparation logic
	// Example: Normalize, trim, or set default values
	return nil
}

// PassedValidation is called after validation passes
func (r *LenderCreateRequest) PassedValidation(ctx http.Context) error {
	// TODO: Add post-validation logic if needed
	return nil
}

// ToCreateData converts the request to create data map
func (r *LenderCreateRequest) ToCreateData() map[string]interface{} {
	//TODO: make the return object strongly typed if possible
	data := map[string]interface{}{
		"name":  r.Name,
		"email": r.Email,
	}

	// Only include optional fields if they have values
	if r.Phone != nil && *r.Phone != "" {
		data["phone"] = r.Phone
	}
	if r.Address != nil && *r.Address != "" {
		data["address"] = r.Address
	}
	if r.Gender != nil && *r.Gender != "" {
		data["gender"] = r.Gender
	}

	return data
}
