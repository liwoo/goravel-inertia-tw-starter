package requests

import (
	"github.com/goravel/framework/contracts/http"

	"smedi-sme-db/app/models"
)

// BdspCreateRequest handles bdsp creation validation
type BdspCreateRequest struct {
	Name               string               `form:"name" json:"name"`
	PostalAddress      string               `form:"postal_address" json:"postal_address"`
	PhysicalAddress    *string              `form:"physical_address" json:"physical_address"`
	RegistrationStatus *string              `form:"registration_status" json:"registration_status"`
	Partners           []string             `form:"partners" json:"partners"`
	ProductTypes       []string             `form:"product_types" json:"product_types"`
	ServiceList        []models.BdspService `form:"service_list" json:"service_list"`
}

// Rules defines validation rules for bdsp creation
func (r *BdspCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":                "required|string|max_len:255",
		"postal_address":      "required|string|max_len:255",
		"physical_address":    "string|max_len:255",
		"registration_status": "string|max_len:255",
		"partners":            "array",
		"partners.*":          "string|max_len:255",
		"product_types":       "array",
		"product_types.*":     "string|max_len:255",
		"service_list":        "array",
		// Deep validation for service_list objects would be nice but 'array' is basic check
		// If we want to validate fields inside objects, we might need custom rules or manual validation
	}
}

// Messages defines custom validation messages
func (r *BdspCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required":           "Name is required",
		"name.max_len":            "Name cannot exceed 255 characters",
		"postal_address.required": "Postal Address is required",
		"partners.array":          "Partners must be an array",
		"product_types.array":     "Product Types must be an array",
		"service_list.array":      "Service List must be an array",
	}
}

// Attributes defines custom attribute names
func (r *BdspCreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to make this request
func (r *BdspCreateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user has permission to create bdsp
	// return facades.Gate().Allows("create.bdsps", ctx)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *BdspCreateRequest) PrepareForValidation(ctx http.Context) error {
	// TODO: Add data preparation logic
	// Example: Normalize, trim, or set default values
	return nil
}

// PassedValidation is called after validation passes
func (r *BdspCreateRequest) PassedValidation(ctx http.Context) error {
	// TODO: Add post-validation logic if needed
	return nil
}

// ToCreateData converts the request to create data map
func (r *BdspCreateRequest) ToCreateData() map[string]interface{} {
	data := map[string]interface{}{
		"name":                r.Name,
		"postal_address":      r.PostalAddress,
		"physical_address":    r.PhysicalAddress,
		"registration_status": r.RegistrationStatus,
		"partners":            r.Partners,
		"product_types":       r.ProductTypes,
		"service_list":        r.ServiceList,
	}

	return data
}
