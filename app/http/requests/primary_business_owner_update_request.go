package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// PrimaryBusinessOwnerUpdateRequest handles primary_business_owner update validation
type PrimaryBusinessOwnerUpdateRequest struct {
	FirstName              *string `form:"first_name" json:"first_name"`
	LastName               *string `form:"last_name" json:"last_name"`
	OtherNames             *string `form:"other_names" json:"other_names"`
	Nationality            *string `form:"nationality" json:"nationality"`
	NationalIdNumber       *string `form:"national_id_number" json:"national_id_number"`
	DateOfBirth            *string `form:"date_of_birth" json:"date_of_birth"`
	Gender                 *string `form:"gender" json:"gender"`
	EducationLevel         *string `form:"education_level" json:"education_level"`
	MalawianStatus         *string `form:"malawian_status" json:"malawian_status"`
	HasSpecialNeeds        *bool   `form:"has_special_needs" json:"has_special_needs"`
	PhoneNumber            *string `form:"phone_number" json:"phone_number"`
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
	SmeID                  *int    `form:"sme_id" json:"sme_id"`
	ID                     uint    `form:"-" json:"-"` // Set by controller
}

// Rules defines validation rules for primary_business_owner updates
func (r *PrimaryBusinessOwnerUpdateRequest) Rules(ctx http.Context) map[string]string {
	rules := map[string]string{}

	// Only validate FirstName if provided
	if r.FirstName != nil {
		rules["first_name"] = "required|max_len:100"
	}
	// Only validate LastName if provided
	if r.LastName != nil {
		rules["last_name"] = "required|max_len:100"
	}
	// Only validate OtherNames if provided
	if r.OtherNames != nil {
		rules["other_names"] = "max_len:100"
	}
	// Only validate Nationality if provided
	if r.Nationality != nil {
		rules["nationality"] = "required|max_len:100"
	}
	// Only validate NationalIdNumber if provided
	if r.NationalIdNumber != nil {
		rules["national_id_number"] = "required|max_len:50"
	}
	// Only validate DateOfBirth if provided
	if r.DateOfBirth != nil {
		rules["date_of_birth"] = "required|date"
	}
	// Only validate Gender if provided
	if r.Gender != nil {
		rules["gender"] = "required|max_len:20"
	}
	// Only validate EducationLevel if provided
	if r.EducationLevel != nil {
		rules["education_level"] = "required|max_len:50"
	}
	// Only validate MalawianStatus if provided
	if r.MalawianStatus != nil {
		rules["malawian_status"] = "required|max_len:50"
	}
	// Only validate HasSpecialNeeds if provided
	if r.HasSpecialNeeds != nil {
		rules["has_special_needs"] = "boolean"
	}
	// Only validate PhoneNumber if provided
	if r.PhoneNumber != nil {
		rules["phone_number"] = "required|max_len:20"
	}
	// Only validate LandlineNumber if provided
	if r.LandlineNumber != nil {
		rules["landline_number"] = "max_len:20"
	}
	// Only validate Email if provided
	if r.Email != nil {
		rules["email"] = "email|max_len:100"
	}
	// Only validate PhysicalAddress if provided
	if r.PhysicalAddress != nil {
		rules["physical_address"] = "max_len:255"
	}
	// Only validate PostalAddress if provided
	if r.PostalAddress != nil {
		rules["postal_address"] = "max_len:255"
	}
	// Only validate Region if provided
	if r.Region != nil {
		rules["region"] = "max_len:100"
	}
	// Only validate District if provided
	if r.District != nil {
		rules["district"] = "max_len:100"
	}
	// Only validate TraditionalAuthority if provided
	if r.TraditionalAuthority != nil {
		rules["traditional_authority"] = "max_len:100"
	}
	// Only validate AltContactName if provided
	if r.AltContactName != nil {
		rules["alt_contact_name"] = "max_len:100"
	}
	// Only validate AltContactRelationship if provided
	if r.AltContactRelationship != nil {
		rules["alt_contact_relationship"] = "max_len:50"
	}
	// Only validate AltContactPhone if provided
	if r.AltContactPhone != nil {
		rules["alt_contact_phone"] = "max_len:20"
	}
	// Only validate SmeID if provided
	if r.SmeID != nil {
		rules["sme_id"] = "required|numeric"
	}

	// If no rules were added, add a dummy rule to prevent empty rules error
	if len(rules) == 0 {
		rules["_at_least_one_field"] = "sometimes"
	}

	return rules
}

// Messages defines custom validation messages for updates
func (r *PrimaryBusinessOwnerUpdateRequest) Messages(ctx http.Context) map[string]string {
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

// Attributes defines custom attribute names for updates
func (r *PrimaryBusinessOwnerUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to update this primary_business_owner
func (r *PrimaryBusinessOwnerUpdateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user can update this specific primary_business_owner
	// return facades.Gate().Allows("update.primary_business_owners", primary_business_owner)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *PrimaryBusinessOwnerUpdateRequest) PrepareForValidation(ctx http.Context) error {
	// TODO: Add data preparation logic for updates
	// Example: Normalize data if provided
	return nil
}

// PassedValidation is called after validation passes
func (r *PrimaryBusinessOwnerUpdateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// ToUpdateData converts the request to update data map
func (r *PrimaryBusinessOwnerUpdateRequest) ToUpdateData() map[string]interface{} {
	data := map[string]interface{}{}

	// Only include FirstName if provided
	if r.FirstName != nil {
		data["first_name"] = *r.FirstName
	}
	// Only include LastName if provided
	if r.LastName != nil {
		data["last_name"] = *r.LastName
	}
	// Only include OtherNames if provided
	if r.OtherNames != nil {
		data["other_names"] = *r.OtherNames
	}
	// Only include Nationality if provided
	if r.Nationality != nil {
		data["nationality"] = *r.Nationality
	}
	// Only include NationalIdNumber if provided
	if r.NationalIdNumber != nil {
		data["national_id_number"] = *r.NationalIdNumber
	}
	// Only include DateOfBirth if provided
	if r.DateOfBirth != nil {
		data["date_of_birth"] = *r.DateOfBirth
	}
	// Only include Gender if provided
	if r.Gender != nil {
		data["gender"] = *r.Gender
	}
	// Only include EducationLevel if provided
	if r.EducationLevel != nil {
		data["education_level"] = *r.EducationLevel
	}
	// Only include MalawianStatus if provided
	if r.MalawianStatus != nil {
		data["malawian_status"] = *r.MalawianStatus
	}
	// Only include HasSpecialNeeds if provided
	if r.HasSpecialNeeds != nil {
		data["has_special_needs"] = *r.HasSpecialNeeds
	}
	// Only include PhoneNumber if provided
	if r.PhoneNumber != nil {
		data["phone_number"] = *r.PhoneNumber
	}
	// Only include LandlineNumber if provided
	if r.LandlineNumber != nil {
		data["landline_number"] = *r.LandlineNumber
	}
	// Only include Email if provided
	if r.Email != nil {
		data["email"] = *r.Email
	}
	// Only include PhysicalAddress if provided
	if r.PhysicalAddress != nil {
		data["physical_address"] = *r.PhysicalAddress
	}
	// Only include PostalAddress if provided
	if r.PostalAddress != nil {
		data["postal_address"] = *r.PostalAddress
	}
	// Only include Region if provided
	if r.Region != nil {
		data["region"] = *r.Region
	}
	// Only include District if provided
	if r.District != nil {
		data["district"] = *r.District
	}
	// Only include TraditionalAuthority if provided
	if r.TraditionalAuthority != nil {
		data["traditional_authority"] = *r.TraditionalAuthority
	}
	// Only include AltContactName if provided
	if r.AltContactName != nil {
		data["alt_contact_name"] = *r.AltContactName
	}
	// Only include AltContactRelationship if provided
	if r.AltContactRelationship != nil {
		data["alt_contact_relationship"] = *r.AltContactRelationship
	}
	// Only include AltContactPhone if provided
	if r.AltContactPhone != nil {
		data["alt_contact_phone"] = *r.AltContactPhone
	}
	// Only include SmeID if provided
	if r.SmeID != nil {
		data["sme_id"] = *r.SmeID
	}

	return data
}

// GetResourceID returns the resource ID for update
func (r *PrimaryBusinessOwnerUpdateRequest) GetResourceID() interface{} {
	return r.ID
}
