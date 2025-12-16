package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// SmeUpdateRequest handles sme update validation
type SmeUpdateRequest struct {
	// UsmeNumber is auto-generated and cannot be updated
	Name                    *string `form:"name" json:"name"`
	RegistrationNumber      *string `form:"registration_number" json:"registration_number"`
	TaxIdentificationNumber *string `form:"tax_identification_number" json:"tax_identification_number"`
	OperationalStartDate    *string `form:"operational_start_date" json:"operational_start_date"`
	BusinessCategory        *string `form:"business_category" json:"business_category"`
	Sector                  *string `form:"sector" json:"sector"`
	SubSector               *string `form:"sub_sector" json:"sub_sector"`
	BusinessDescription     *string `form:"business_description" json:"business_description"`
	ContactPhone            *string `form:"contact_phone" json:"contact_phone"`
	ContactEmail            *string `form:"contact_email" json:"contact_email"`
	PhysicalAddress         *string `form:"physical_address" json:"physical_address"`
	PostalAddress           *string `form:"postal_address" json:"postal_address"`
	Website                 *string `form:"website" json:"website"`
	// Region is inferred from District - not needed in update
	District                   *string   `form:"district" json:"district"`
	TraditionalAuthority       *string   `form:"traditional_authority" json:"traditional_authority"`
	IpAddress                  *string   `form:"ip_address" json:"ip_address"`
	UserAgent                  *string   `form:"user_agent" json:"user_agent"`
	BusinessImprovementAspects *[]string `form:"business_improvement_aspects" json:"business_improvement_aspects"`
	BusinessAccessedFinancing  *[]string `form:"business_accessed_financing" json:"business_accessed_financing"`
	Classification             *string   `form:"classification" json:"classification"`
	ID                         uint      `form:"-" json:"-"` // Set by controller
}

// Rules defines validation rules for sme updates
func (r *SmeUpdateRequest) Rules(ctx http.Context) map[string]string {
	rules := map[string]string{}

	// UsmeNumber is auto-generated and cannot be updated

	// Only validate Name if provided
	if r.Name != nil {
		rules["name"] = "required|max_len:255"
	}
	// Only validate RegistrationNumber if provided
	if r.RegistrationNumber != nil {
		rules["registration_number"] = "max_len:100"
	}
	// Only validate TaxIdentificationNumber if provided
	if r.TaxIdentificationNumber != nil {
		rules["tax_identification_number"] = "max_len:100"
	}
	// Only validate OperationalStartDate if provided
	if r.OperationalStartDate != nil {
		rules["operational_start_date"] = "date"
	}
	// Only validate BusinessCategory if provided
	if r.BusinessCategory != nil {
		rules["business_category"] = "required|max_len:100"
	}
	// Only validate Sector if provided
	if r.Sector != nil {
		rules["sector"] = "required|max_len:100"
	}
	// Only validate SubSector if provided
	if r.SubSector != nil {
		rules["sub_sector"] = "max_len:100"
	}
	// Only validate BusinessDescription if provided
	if r.BusinessDescription != nil {
		rules["business_description"] = "max_len:1000"
	}
	// Only validate ContactPhone if provided
	if r.ContactPhone != nil {
		rules["contact_phone"] = "required|max_len:20"
	}
	// Only validate ContactEmail if provided
	if r.ContactEmail != nil {
		rules["contact_email"] = "required|email|max_len:100"
	}
	// Only validate PhysicalAddress if provided
	if r.PhysicalAddress != nil {
		rules["physical_address"] = "max_len:255"
	}
	// Only validate PostalAddress if provided
	if r.PostalAddress != nil {
		rules["postal_address"] = "max_len:255"
	}
	// Only validate Website if provided
	if r.Website != nil {
		rules["website"] = "url|max_len:255"
	}
	// Region is inferred from district - not validated
	// Only validate District if provided
	if r.District != nil {
		rules["district"] = "max_len:100"
	}
	// Only validate TraditionalAuthority if provided
	if r.TraditionalAuthority != nil {
		rules["traditional_authority"] = "max_len:100"
	}
	// Only validate BusinessImprovementAspects if provided
	if r.BusinessImprovementAspects != nil {
		rules["business_improvement_aspects"] = "required|array"
		rules["business_improvement_aspects.*"] = "max_len:100"
	}
	// Only validate BusinessAccessedFinancing if provided
	if r.BusinessAccessedFinancing != nil {
		rules["business_accessed_financing"] = "required|array"
		rules["business_accessed_financing.*"] = "max_len:100"
	}
	// Only validate Classification if provided
	if r.Classification != nil {
		rules["classification"] = "in:Micro,Small,Medium,Unclassified"
	}

	// If no rules were added, add a dummy rule to prevent empty rules error
	if len(rules) == 0 {
		rules["_at_least_one_field"] = "sometimes"
	}

	return rules
}

// Messages defines custom validation messages for updates
func (r *SmeUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"usme_number.required":          "USME Number is required",
		"usme_number.max":               "USME Number cannot exceed 100 characters",
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
		"classification.in":                     "Classification must be one of: Micro, Small, Medium, Unclassified",
	}
}

// Attributes defines custom attribute names for updates
func (r *SmeUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to update this sme
func (r *SmeUpdateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user can update this specific sme
	// return facades.Gate().Allows("update.smes", sme)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *SmeUpdateRequest) PrepareForValidation(ctx http.Context) error {
	// TODO: Add data preparation logic for updates
	// Example: Normalize data if provided
	return nil
}

// PassedValidation is called after validation passes
func (r *SmeUpdateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// ToUpdateData converts the request to update data map
func (r *SmeUpdateRequest) ToUpdateData() map[string]interface{} {
	data := map[string]interface{}{}

	// UsmeNumber is auto-generated and cannot be updated

	// Only include Name if provided
	if r.Name != nil {
		data["name"] = *r.Name
	}
	// Only include RegistrationNumber if provided
	if r.RegistrationNumber != nil {
		data["registration_number"] = *r.RegistrationNumber
	}
	// Only include TaxIdentificationNumber if provided
	if r.TaxIdentificationNumber != nil {
		data["tax_identification_number"] = *r.TaxIdentificationNumber
	}
	// Only include OperationalStartDate if provided
	if r.OperationalStartDate != nil {
		data["operational_start_date"] = *r.OperationalStartDate
	}
	// Only include BusinessCategory if provided
	if r.BusinessCategory != nil {
		data["business_category"] = *r.BusinessCategory
	}
	// Only include Sector if provided
	if r.Sector != nil {
		data["sector"] = *r.Sector
	}
	// Only include SubSector if provided
	if r.SubSector != nil {
		data["sub_sector"] = *r.SubSector
	}
	// Only include BusinessDescription if provided
	if r.BusinessDescription != nil {
		data["business_description"] = *r.BusinessDescription
	}
	// Only include ContactPhone if provided
	if r.ContactPhone != nil {
		data["contact_phone"] = *r.ContactPhone
	}
	// Only include ContactEmail if provided
	if r.ContactEmail != nil {
		data["contact_email"] = *r.ContactEmail
	}
	// Only include PhysicalAddress if provided
	if r.PhysicalAddress != nil {
		data["physical_address"] = *r.PhysicalAddress
	}
	// Only include PostalAddress if provided
	if r.PostalAddress != nil {
		data["postal_address"] = *r.PostalAddress
	}
	// Only include Website if provided
	if r.Website != nil {
		data["website"] = *r.Website
	}
	// Region is inferred from district - not included in update
	// Only include District if provided
	if r.District != nil {
		data["district"] = *r.District
	}
	// Only include TraditionalAuthority if provided
	if r.TraditionalAuthority != nil {
		data["traditional_authority"] = *r.TraditionalAuthority
	}
	// Only include IpAddress if provided
	if r.IpAddress != nil {
		data["ip_address"] = *r.IpAddress
	}
	// Only include UserAgent if provided
	if r.UserAgent != nil {
		data["user_agent"] = *r.UserAgent
	}
	// Only include BusinessImprovementAspects if provided
	if r.BusinessImprovementAspects != nil {
		data["business_improvement_aspects"] = *r.BusinessImprovementAspects
	}
	// Only include BusinessAccessedFinancing if provided
	if r.BusinessAccessedFinancing != nil {
		data["business_accessed_financing"] = *r.BusinessAccessedFinancing
	}
	// Only include Classification if provided
	if r.Classification != nil {
		data["classification"] = *r.Classification
	}

	return data
}

// GetResourceID returns the resource ID for update
func (r *SmeUpdateRequest) GetResourceID() interface{} {
	return r.ID
}
