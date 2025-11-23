package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// SmeCreateRequest handles sme creation validation
type SmeCreateRequest struct {
	Name                    string  `form:"name" json:"name"`
	RegistrationNumber      *string `form:"registration_number" json:"registration_number"`
	TaxIdentificationNumber *string `form:"tax_identification_number" json:"tax_identification_number"`
	OperationalStartDate    string  `form:"operational_start_date" json:"operational_start_date"`
	BusinessCategory        string  `form:"business_category" json:"business_category"`
	Sector                  string  `form:"sector" json:"sector"`
	SubSector               *string `form:"sub_sector" json:"sub_sector"`
	BusinessDescription     *string `form:"business_description" json:"business_description"`
	ContactPhone            string  `form:"contact_phone" json:"contact_phone"`
	ContactEmail            string  `form:"contact_email" json:"contact_email"`
	PhysicalAddress         *string `form:"physical_address" json:"physical_address"`
	PostalAddress           *string `form:"postal_address" json:"postal_address"`
	Website                 *string `form:"website" json:"website"`
	// Region is inferred from District - not needed in request
	District                   *string  `form:"district" json:"district"`
	TraditionalAuthority       *string  `form:"traditional_authority" json:"traditional_authority"`
	BusinessImprovementAspects []string `form:"business_improvement_aspects" json:"business_improvement_aspects"`
	BusinessAccessedFinancing  []string `form:"business_accessed_financing" json:"business_accessed_financing"`
}

// Rules defines validation rules for sme creation
func (r *SmeCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"usme_number":               "max_len:100", // Optional - will be auto-generated if not provided
		"name":                      "required|max_len:255",
		"registration_number":       "max_len:100",
		"tax_identification_number": "max_len:100",
		"operational_start_date":    "date",
		"business_category":         "required|max_len:100",
		"sector":                    "required|max_len:100",
		"sub_sector":                "max_len:100",
		"business_description":      "max_len:1000",
		"contact_phone":             "required|max_len:20",
		"contact_email":             "required|email|max_len:100",
		"physical_address":          "max_len:255",
		"postal_address":            "max_len:255",
		"website":                   "full_url|max_len:255",
		// region is inferred from district
		"district":                       "max_len:100",
		"traditional_authority":          "max_len:100",
		"business_improvement_aspects":   "required|array",
		"business_improvement_aspects.*": "max_len:100",
		"business_accessed_financing":    "required|array",
		"business_accessed_financing.*":  "max_len:100",
	}
}

// Messages defines custom validation messages
func (r *SmeCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required":                 "Business Name is required",
		"name.max":                      "Business Name cannot exceed 255 characters",
		"registration_number.max":       "Registration Number cannot exceed 100 characters",
		"tax_identification_number.max": "Tax Identification Number cannot exceed 100 characters",
		"operational_start_date.date":   "Operational Start Date must be a valid date",
		"business_category.required":    "Business Category is required",
		"business_category.max":         "Business Category cannot exceed 100 characters",
		"sector.required":               "Sector is required",
		"sector.max":                    "Sector cannot exceed 100 characters",
		"sub_sector.max":                "Sub Sector cannot exceed 100 characters",
		"business_description.max":      "Business Description cannot exceed 1000 characters",
		"contact_phone.required":        "Contact Phone is required",
		"contact_phone.max":             "Contact Phone cannot exceed 20 characters",
		"contact_email.required":        "Contact Email is required",
		"contact_email.email":           "Contact Email must be a valid email address",
		"contact_email.max":             "Contact Email cannot exceed 100 characters",
		"physical_address.max":          "Physical Address cannot exceed 255 characters",
		"postal_address.max":            "Postal Address cannot exceed 255 characters",
		"website.url":                   "Website must be a valid URL",
		"website.max":                   "Website cannot exceed 255 characters",
		// region is inferred from district
		"district.max":                          "District cannot exceed 100 characters",
		"traditional_authority.max":             "Traditional Authority cannot exceed 100 characters",
		"business_improvement_aspects.required": "Business Improvement Aspects is required",
		"business_improvement_aspects.array":    "Business Improvement Aspects must be an array",
		"business_improvement_aspects.*.max":    "Each Business Improvement Aspect cannot exceed 100 characters",
		"business_accessed_financing.required":  "Business Accessed Financing is required",
		"business_accessed_financing.array":     "Business Accessed Financing must be an array",
		"business_accessed_financing.*.max":     "Each Business Accessed Financing entry cannot exceed 100 characters",
	}
}

// Attributes defines custom attribute names
func (r *SmeCreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to make this request
func (r *SmeCreateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user has permission to create sme
	// return facades.Gate().Allows("create.smes", ctx)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *SmeCreateRequest) PrepareForValidation(ctx http.Context) error {
	// TODO: Add data preparation logic
	// Example: Normalize, trim, or set default values
	return nil
}

// PassedValidation is called after validation passes
func (r *SmeCreateRequest) PassedValidation(ctx http.Context) error {
	// TODO: Add post-validation logic if needed
	return nil
}

// ToCreateData converts the request to create data map
func (r *SmeCreateRequest) ToCreateData() map[string]interface{} {
	data := map[string]interface{}{
		"name":                      r.Name,
		"registration_number":       r.RegistrationNumber,
		"tax_identification_number": r.TaxIdentificationNumber,
		"operational_start_date":    r.OperationalStartDate,
		"business_category":         r.BusinessCategory,
		"sector":                    r.Sector,
		"sub_sector":                r.SubSector,
		"business_description":      r.BusinessDescription,
		"contact_phone":             r.ContactPhone,
		"contact_email":             r.ContactEmail,
		"physical_address":          r.PhysicalAddress,
		"postal_address":            r.PostalAddress,
		"website":                   r.Website,
		// region is inferred from district
		"district":                     r.District,
		"traditional_authority":        r.TraditionalAuthority,
		"business_improvement_aspects": r.BusinessImprovementAspects,
		"business_accessed_financing":  r.BusinessAccessedFinancing,
	}

	return data
}
