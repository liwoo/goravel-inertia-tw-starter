package requests

import (
	"github.com/goravel/framework/contracts/http"

	"smedi-sme-db/app/models"
)

// BdspUpdateRequest handles bdsp update validation
type BdspUpdateRequest struct {
	Name               *string               `form:"name" json:"name"`
	PostalAddress      *string               `form:"postal_address" json:"postal_address"`
	PhysicalAddress    *string               `form:"physical_address" json:"physical_address"`
	RegistrationStatus *string               `form:"registration_status" json:"registration_status"`
	Partners           *[]string             `form:"partners" json:"partners"`
	ProductTypes       *[]string             `form:"product_types" json:"product_types"`
	ServiceList        *[]models.BdspService `form:"service_list" json:"service_list"`
	ID                 uint                  `form:"-" json:"-"` // Set by controller
}

// Rules defines validation rules for bdsp updates
func (r *BdspUpdateRequest) Rules(ctx http.Context) map[string]string {
	rules := map[string]string{}

	// Only validate Name if provided
	if r.Name != nil {
		rules["name"] = "required|string|max_len:255"
	}
	if r.PostalAddress != nil {
		rules["postal_address"] = "required|string|max_len:255"
	}
	if r.PhysicalAddress != nil {
		rules["physical_address"] = "string|max_len:255"
	}
	if r.RegistrationStatus != nil {
		rules["registration_status"] = "string|max_len:255"
	}
	if r.Partners != nil {
		rules["partners"] = "array"
		rules["partners.*"] = "string|max_len:255"
	}
	if r.ProductTypes != nil {
		rules["product_types"] = "array"
		rules["product_types.*"] = "string|max_len:255"
	}
	if r.ServiceList != nil {
		rules["service_list"] = "array"
	}

	// If no rules were added, add a dummy rule to prevent empty rules error
	if len(rules) == 0 {
		rules["_at_least_one_field"] = "sometimes"
	}

	return rules
}

// Messages defines custom validation messages for updates
func (r *BdspUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required":       "Name is required",
		"name.max":            "Name cannot exceed 255 characters",
		"partners.array":      "Partners must be an array",
		"product_types.array": "Product Types must be an array",
		"service_list.array":  "Service List must be an array",
	}
}

// Attributes defines custom attribute names for updates
func (r *BdspUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to update this bdsp
func (r *BdspUpdateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user can update this specific bdsp
	// return facades.Gate().Allows("update.bdsps", bdsp)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *BdspUpdateRequest) PrepareForValidation(ctx http.Context) error {
	// TODO: Add data preparation logic for updates
	// Example: Normalize data if provided
	return nil
}

// PassedValidation is called after validation passes
func (r *BdspUpdateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// ToUpdateData converts the request to update data map
func (r *BdspUpdateRequest) ToUpdateData() map[string]interface{} {
	data := map[string]interface{}{}

	// Only include Name if provided
	if r.Name != nil {
		data["name"] = *r.Name
	}
	// Only include PostalAddress if provided
	if r.PostalAddress != nil {
		data["postal_address"] = *r.PostalAddress
	}
	// Only include PhysicalAddress if provided
	if r.PhysicalAddress != nil {
		data["physical_address"] = *r.PhysicalAddress
	}
	// Only include RegistrationStatus if provided
	if r.RegistrationStatus != nil {
		data["registration_status"] = *r.RegistrationStatus
	}
	if r.Partners != nil {
		data["partners"] = *r.Partners
	}
	if r.ProductTypes != nil {
		data["product_types"] = *r.ProductTypes
	}
	if r.ServiceList != nil {
		data["service_list"] = *r.ServiceList
	}

	return data
}

// GetResourceID returns the resource ID for update
func (r *BdspUpdateRequest) GetResourceID() interface{} {
	return r.ID
}
