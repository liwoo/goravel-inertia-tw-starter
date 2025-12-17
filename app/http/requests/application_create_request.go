package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// ApplicationCreateRequest handles application creation validation
type ApplicationCreateRequest struct {
	SME                        string `form:"sme" json:"sme"`
	RegistrantName             string `form:"registrant_name" json:"registrant_name"`
	Email                      string `form:"email" json:"email"`
	Phone                      string `form:"phone" json:"phone"`
	SMERegistrationNumber      string `form:"sme_registration_number" json:"sme_registration_number"`
	SMETaxIdentificationNumber string `form:"sme_tax_identification_number" json:"sme_tax_identification_number"`

	// Primary Business Owner Details
	FirstName              string `form:"first_name" json:"first_name"`
	LastName               string `form:"last_name" json:"last_name"`
	OtherNames             string `form:"other_names" json:"other_names"`
	Nationality            string `form:"nationality" json:"nationality"`
	NationalIDNumber       string `form:"national_id_number" json:"national_id_number"`
	DateOfBirth            string `form:"date_of_birth" json:"date_of_birth"`
	Gender                 string `form:"gender" json:"gender"`
	EducationLevel         string `form:"education_level" json:"education_level"`
	MalawianStatus         string `form:"malawian_status" json:"malawian_status"`
	HasSpecialNeeds        bool   `form:"has_special_needs" json:"has_special_needs"`
	LandlineNumber         string `form:"landline_number" json:"landline_number"`
	PhysicalAddress        string `form:"physical_address" json:"physical_address"`
	PostalAddress          string `form:"postal_address" json:"postal_address"`
	Region                 string `form:"region" json:"region"`
	District               string `form:"district" json:"district"`
	TraditionalAuthority   string `form:"traditional_authority" json:"traditional_authority"`
	AltContactName         string `form:"alt_contact_name" json:"alt_contact_name"`
	AltContactRelationship string `form:"alt_contact_relationship" json:"alt_contact_relationship"`
	AltContactPhone        string `form:"alt_contact_phone" json:"alt_contact_phone"`
}

// Rules defines validation rules for application creation
func (r *ApplicationCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"sme":                           "required|string|max_len:255",
		"registrant_name":               "required|string|max_len:255",
		"email":                         "required|string|max_len:255",
		"phone":                         "required|string|max_len:255",
		"sme_registration_number":       "string|max_len:255",
		"sme_tax_identification_number": "string|max_len:255",
		"first_name":                    "required|string|max_len:100",
		"last_name":                     "required|string|max_len:100",
		"nationality":                   "required|string|max_len:100",
		"national_id_number":            "required|string|max_len:50",
		"date_of_birth":                 "required|date",
		"gender":                        "required|in:MALE,FEMALE",
		"education_level":               "required|string|max_len:100",
		"malawian_status":               "required|string|max_len:50",
	}
}

// Messages defines custom validation messages
func (r *ApplicationCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"sme.required":                           "Sme is required",
		"sme.max_len":                            "Sme cannot exceed 255 characters",
		"registrant_name.required":               "Registrant Name is required",
		"registrant_name.max_len":                "Registrant Name cannot exceed 255 characters",
		"email.required":                         "Email is required",
		"email.max_len":                          "Email cannot exceed 255 characters",
		"phone.required":                         "Phone is required",
		"phone.max_len":                          "Phone cannot exceed 255 characters",
		"sme_registration_number.max_len":        "Business Registration Number cannot exceed 255 characters",
		"sme_tax_identification_number.max_len":  "Tax Identification Number cannot exceed 255 characters",
		"status.required":                        "Status is required",
	}
}

// Attributes defines custom attribute names
func (r *ApplicationCreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to make this request
func (r *ApplicationCreateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user has permission to create application
	// return facades.Gate().Allows("create.applications", ctx)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *ApplicationCreateRequest) PrepareForValidation(ctx http.Context) error {
	// TODO: Add data preparation logic
	// Example: Normalize, trim, or set default values
	return nil
}

// PassedValidation is called after validation passes
func (r *ApplicationCreateRequest) PassedValidation(ctx http.Context) error {
	// TODO: Add post-validation logic if needed
	return nil
}

// ToCreateData converts the request to create data map
func (r *ApplicationCreateRequest) ToCreateData() map[string]interface{} {
	data := map[string]interface{}{
		"sme":                           r.SME,
		"registrant_name":               r.RegistrantName,
		"email":                         r.Email,
		"phone":                         r.Phone,
		"sme_registration_number":       r.SMERegistrationNumber,
		"sme_tax_identification_number": r.SMETaxIdentificationNumber,
		"first_name":                    r.FirstName,
		"last_name":                     r.LastName,
		"other_names":                   r.OtherNames,
		"nationality":                   r.Nationality,
		"national_id_number":            r.NationalIDNumber,
		"date_of_birth":                 r.DateOfBirth,
		"gender":                        r.Gender,
		"education_level":               r.EducationLevel,
		"malawian_status":               r.MalawianStatus,
		"has_special_needs":             r.HasSpecialNeeds,
		"landline_number":               r.LandlineNumber,
		"physical_address":              r.PhysicalAddress,
		"postal_address":                r.PostalAddress,
		"region":                        r.Region,
		"district":                      r.District,
		"traditional_authority":         r.TraditionalAuthority,
		"alt_contact_name":              r.AltContactName,
		"alt_contact_relationship":      r.AltContactRelationship,
		"alt_contact_phone":             r.AltContactPhone,
	}

	return data
}
