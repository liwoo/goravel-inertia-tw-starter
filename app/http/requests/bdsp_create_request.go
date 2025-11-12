package requests

import (
	"smedi-sme-db/app/models"

	"github.com/goravel/framework/contracts/http"
)

// BdspCreateRequest handles bdsp creation validation
type BdspCreateRequest struct {
	Name               string               `form:"name" json:"name"`
	PostalAddress      *string              `form:"postal_address" json:"postal_address"`
	PhysicalAddress    *string              `form:"physical_address" json:"physical_address"`
	RegistrationStatus string               `form:"registration_status" json:"registration_status"`
	ProductTypes       []string             `form:"product_types" json:"product_types"`
	ServiceList        []models.BdspService `form:"service_list" json:"service_list"`
	AssociatedPartners []string             `form:"associated_partners" json:"associated_partners"`
	IpAddress          *string              `form:"ip_address" json:"ip_address"`
	UserAgent          *string              `form:"user_agent" json:"user_agent"`
}

// Rules defines validation rules for bdsp creation
func (r *BdspCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":                "required|string|max:255",
		"registration_status": "required|string|max:255",
		"product_types":       "required|array",
		"service_list":        "required",
		"associated_partners": "required|array",
	}
}

// Messages defines custom validation messages
func (r *BdspCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required":                "Name is required",
		"name.max":                     "Name cannot exceed 255 characters",
		"registration_status.required": "Registration Status is required",
		"registration_status.max":      "Registration Status cannot exceed 255 characters",
		"product_types.required":       "Product Types is required",
		"service_list.required":        "Service List is required",
		"associated_partners.required": "Associated Partners is required",
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
		"product_types":       r.ProductTypes,
		"service_list":        r.ServiceList,
		"associated_partners": r.AssociatedPartners,
		"ip_address":          r.IpAddress,
		"user_agent":          r.UserAgent,
	}

	return data
}
