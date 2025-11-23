package requests

import (
	"fmt"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/str"
)

// PrimaryBusinessOwnerCreateRequest handles primary_business_owner creation validation
type PrimaryBusinessOwnerCreateRequest struct {
	FirstName              string  `form:"first_name" json:"first_name"`
	LastName               string  `form:"last_name" json:"last_name"`
	OtherNames             *string `form:"other_names" json:"other_names"`
	Nationality            string  `form:"nationality" json:"nationality"`
	NationalIdNumber       string  `form:"national_id_number" json:"national_id_number"`
	DateOfBirth            string  `form:"date_of_birth" json:"date_of_birth"`
	Gender                 string  `form:"gender" json:"gender"`
	EducationLevel         string  `form:"education_level" json:"education_level"`
	MalawianStatus         string  `form:"malawian_status" json:"malawian_status"`
	HasSpecialNeeds        bool    `form:"has_special_needs" json:"has_special_needs"`
	PhoneNumber            string  `form:"phone_number" json:"phone_number"`
	LandlineNumber         *string `form:"landline_number" json:"landline_number"`
	Email                  *string `form:"email" json:"email"`
	PhysicalAddress        *string `form:"physical_address" json:"physical_address"`
	PostalAddress          *string `form:"postal_address" json:"postal_address"`
	Region                 *string `form:"region" json:"region"`
	District               *string `form:"district" json:"district"`
	TraditionalAuthority   *string `form:"traditional_authority" json:"traditional_authority"`
	AltContactName         *string `form:"alt_contact_name" json:"alt_contact_name"`
	AltContactRelationship *string `form:"alt_contact_relationship" json:"alt_contact_relationship"`
	AltContactPhone        *string `form:"alt_contact_phone" json:"alt_contact_phone"`
	SmeID                  int     `form:"sme_id" json:"sme_id"`
}

// Rules defines validation rules for primary_business_owner creation
func (r *PrimaryBusinessOwnerCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"first_name":               "required|max_len:100",
		"last_name":                "required|max_len:100",
		"other_names":              "max_len:100",
		"nationality":              "required|max_len:100",
		"national_id_number":       "required|max_len:50",
		"date_of_birth":            "required",
		"gender":                   "required|max_len:20",
		"education_level":          "required|max_len:50",
		"malawian_status":          "required|max_len:50",
		"has_special_needs":        "boolean",
		"phone_number":             "required|max_len:20",
		"landline_number":          "max_len:20",
		"email":                    "email|max_len:100",
		"physical_address":         "max_len:255",
		"postal_address":           "max_len:255",
		"region":                   "max_len:100",
		"district":                 "max_len:100",
		"traditional_authority":    "max_len:100",
		"alt_contact_name":         "max_len:100",
		"alt_contact_relationship": "max_len:50",
		"alt_contact_phone":        "max_len:20",
		"sme_id":                   "required|numeric",
	}
}

// Messages defines custom validation messages
func (r *PrimaryBusinessOwnerCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"first_name.required":              "First Name is required",
		"first_name.max_len":               "First Name cannot exceed 100 characters",
		"last_name.required":               "Last Name is required",
		"last_name.max_len":                "Last Name cannot exceed 100 characters",
		"other_names.max_len":              "Other Names cannot exceed 100 characters",
		"nationality.required":             "Nationality is required",
		"nationality.max_len":              "Nationality cannot exceed 100 characters",
		"national_id_number.required":      "National ID Number is required",
		"national_id_number.max_len":       "National ID Number cannot exceed 50 characters",
		"date_of_birth.required":           "Date of Birth is required",
		"date_of_birth.date":               "Date of Birth must be a valid date",
		"gender.required":                  "Gender is required",
		"gender.max_len":                   "Gender cannot exceed 20 characters",
		"education_level.required":         "Education Level is required",
		"education_level.max_len":          "Education Level cannot exceed 50 characters",
		"malawian_status.required":         "Malawian Status is required",
		"malawian_status.max_len":          "Malawian Status cannot exceed 50 characters",
		"has_special_needs.boolean":        "Has Special Needs must be true or false",
		"phone_number.required":            "Phone Number is required",
		"phone_number.max_len":             "Phone Number cannot exceed 20 characters",
		"landline_number.max_len":          "Landline Number cannot exceed 20 characters",
		"email.email":                      "Email must be a valid email address",
		"email.max_len":                    "Email cannot exceed 100 characters",
		"physical_address.max_len":         "Physical Address cannot exceed 255 characters",
		"postal_address.max_len":           "Postal Address cannot exceed 255 characters",
		"region.max_len":                   "Region cannot exceed 100 characters",
		"district.max_len":                 "District cannot exceed 100 characters",
		"traditional_authority.max_len":    "Traditional Authority cannot exceed 100 characters",
		"alt_contact_name.max_len":         "Alternate Contact Name cannot exceed 100 characters",
		"alt_contact_relationship.max_len": "Alternate Contact Relationship cannot exceed 50 characters",
		"alt_contact_phone.max_len":        "Alternate Contact Phone cannot exceed 20 characters",
		"sme_id.required":                  "SME ID is required",
		"sme_id.numeric":                   "SME ID must be a number",
	}
}

// Attributes defines custom attribute names
func (r *PrimaryBusinessOwnerCreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to make this request
func (r *PrimaryBusinessOwnerCreateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user has permission to create primary_business_owner
	// return facades.Gate().Allows("create.primary_business_owners", ctx)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *PrimaryBusinessOwnerCreateRequest) PrepareForValidation(ctx http.Context) error {
	fmt.Printf("PrepareForValidation - DateOfBirth: %T %+v\n", r.DateOfBirth, r.DateOfBirth)
	return nil
}

// PassedValidation is called after validation passes
func (r *PrimaryBusinessOwnerCreateRequest) PassedValidation(ctx http.Context) error {
	// TODO: Add post-validation logic if needed
	return nil
}

// ToCreateData converts the request to create data map
func (r *PrimaryBusinessOwnerCreateRequest) ToCreateData() map[string]interface{} {

	data := map[string]interface{}{
		"first_name":               r.FirstName,
		"last_name":                r.LastName,
		"other_names":              r.OtherNames,
		"nationality":              r.Nationality,
		"date_of_birth":            carbon.DateTime{Carbon: carbon.Now()}, // Convert to YYYY-MM-DD string
		"national_id_number":       r.NationalIdNumber,
		"gender":                   str.Of(r.Gender).Upper().String(),
		"education_level":          r.EducationLevel,
		"malawian_status":          r.MalawianStatus,
		"has_special_needs":        r.HasSpecialNeeds,
		"phone_number":             r.PhoneNumber,
		"landline_number":          r.LandlineNumber,
		"email":                    r.Email,
		"physical_address":         r.PhysicalAddress,
		"postal_address":           r.PostalAddress,
		"region":                   r.Region,
		"district":                 r.District,
		"traditional_authority":    r.TraditionalAuthority,
		"alt_contact_name":         r.AltContactName,
		"alt_contact_relationship": r.AltContactRelationship,
		"alt_contact_phone":        r.AltContactPhone,
		"sme_id":                   uint(r.SmeID),
	}

	return data
}
