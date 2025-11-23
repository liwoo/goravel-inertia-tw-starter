package requests

import (
	"errors"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"
)

// AdditionalBusinessMemberUpdateRequest handles additional_business_member update validation
type AdditionalBusinessMemberUpdateRequest struct {
	FirstName        *string `form:"first_name" json:"first_name"`
	LastName         *string `form:"last_name" json:"last_name"`
	OtherNames       *string `form:"other_names" json:"other_names"`
	Gender           *string `form:"gender" json:"gender"`
	Nationality      *string `form:"nationality" json:"nationality"`
	NationalIdNumber *string `form:"national_id_number" json:"national_id_number"`
	DateOfBirth      *string `form:"date_of_birth" json:"date_of_birth"`
	Email            *string `form:"email" json:"email"`
	PhoneNumber      *string `form:"phone_number" json:"phone_number"`
	IsIntern         *bool   `form:"is_intern" json:"is_intern"`
	IsPartTime       *bool   `form:"is_part_time" json:"is_part_time"`
	SmeId            *int    `form:"sme_id" json:"sme_id"`
	ID               uint    `form:"-" json:"-"` // Set by controller
}

// Rules defines validation rules for additional_business_member updates
func (r *AdditionalBusinessMemberUpdateRequest) Rules(ctx http.Context) map[string]string {
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
	// Note: date_of_birth validation removed because the 'date' validator doesn't work with *string
	// Date parsing/validation happens in ToUpdateData() method instead
	if r.DateOfBirth != nil {
		// Skip validation - will be validated during parsing
	}
	// Only validate Email if provided
	if r.Email != nil {
		rules["email"] = "email|max_len:100"
	}
	// Only validate PhoneNumber if provided
	if r.PhoneNumber != nil {
		rules["phone_number"] = "required|max_len:20"
	}
	// Only validate IsIntern if provided
	if r.IsIntern != nil {
		rules["is_intern"] = "boolean"
	}
	// Only validate IsPartTime if provided
	if r.IsPartTime != nil {
		rules["is_part_time"] = "boolean"
	}
	// Only validate SmeId if provided
	if r.SmeId != nil {
		rules["sme_id"] = "required|numeric"
	}

	// If no rules were added, add a dummy rule to prevent empty rules error
	if len(rules) == 0 {
		rules["_at_least_one_field"] = "sometimes"
	}

	return rules
}

// Messages defines custom validation messages for updates
func (r *AdditionalBusinessMemberUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"first_name.required":         "First Name is required",
		"first_name.max_len":          "First Name cannot exceed 100 characters",
		"last_name.required":          "Last Name is required",
		"last_name.max_len":           "Last Name cannot exceed 100 characters",
		"other_names.max_len":         "Other Names cannot exceed 100 characters",
		"nationality.required":        "Nationality is required",
		"nationality.max_len":         "Nationality cannot exceed 100 characters",
		"national_id_number.required": "National ID Number is required",
		"national_id_number.max_len":  "National ID Number cannot exceed 50 characters",
		"date_of_birth.date":          "Date of Birth must be a valid date",
		"email.email":                 "Email must be a valid email address",
		"email.max_len":               "Email cannot exceed 100 characters",
		"phone_number.required":       "Phone Number is required",
		"phone_number.max_len":        "Phone Number cannot exceed 20 characters",
		"is_intern.boolean":           "Is Intern must be true or false",
		"is_part_time.boolean":        "Is Part Time must be true or false",
		"sme_id.required":             "SME ID is required",
		"sme_id.numeric":              "SME ID must be a number",
	}
}

// Attributes defines custom attribute names for updates
func (r *AdditionalBusinessMemberUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to update this additional_business_member
func (r *AdditionalBusinessMemberUpdateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user can update this specific additional_business_member
	// return facades.Gate().Allows("update.additional_business_members", additional_business_member)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *AdditionalBusinessMemberUpdateRequest) PrepareForValidation(ctx http.Context) error {
	// Validate date_of_birth format if provided
	if r.DateOfBirth != nil && *r.DateOfBirth != "" {
		parsedDate := carbon.Parse(*r.DateOfBirth)
		if parsedDate.Error != nil {
			// Return a validation error that will be caught by the framework
			return errors.New("date_of_birth: Date of Birth must be a valid date")
		}
	}
	return nil
}

// PassedValidation is called after validation passes
func (r *AdditionalBusinessMemberUpdateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// ToUpdateData converts the request to update data map
func (r *AdditionalBusinessMemberUpdateRequest) ToUpdateData() map[string]interface{} {
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
	// Only include Gender if provided
	if r.Gender != nil {
		data["gender"] = *r.Gender
	}
	// Only include Nationality if provided
	if r.Nationality != nil {
		data["nationality"] = *r.Nationality
	}
	// Only include NationalIdNumber if provided
	if r.NationalIdNumber != nil {
		data["national_id_number"] = *r.NationalIdNumber
	}
	// Only include DateOfBirth if provided and convert to carbon.DateTime
	if r.DateOfBirth != nil {
		if *r.DateOfBirth != "" {
			parsedDate := carbon.Parse(*r.DateOfBirth)
			if parsedDate.Error == nil {
				data["date_of_birth"] = carbon.NewDateTime(parsedDate)
			} else {
				// If parsing fails, set to nil
				data["date_of_birth"] = nil
			}
		} else {
			// Empty string means clear the date
			data["date_of_birth"] = nil
		}
	}
	// Only include Email if provided
	if r.Email != nil {
		data["email"] = *r.Email
	}
	// Only include PhoneNumber if provided
	if r.PhoneNumber != nil {
		data["phone_number"] = *r.PhoneNumber
	}
	// Only include IsIntern if provided
	if r.IsIntern != nil {
		data["is_intern"] = *r.IsIntern
	}
	// Only include IsPartTime if provided
	if r.IsPartTime != nil {
		data["is_part_time"] = *r.IsPartTime
	}
	// Only include SmeId if provided
	if r.SmeId != nil {
		data["sme_id"] = *r.SmeId
	}

	return data
}

// GetResourceID returns the resource ID for update
func (r *AdditionalBusinessMemberUpdateRequest) GetResourceID() interface{} {
	return r.ID
}
