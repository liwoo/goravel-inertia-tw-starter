package requests

import (
	"errors"
	"regexp"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"
)

// AdditionalBusinessMemberCreateRequest handles additional_business_member creation validation
type AdditionalBusinessMemberCreateRequest struct {
	FirstName        string  `form:"first_name" json:"first_name"`
	LastName         string  `form:"last_name" json:"last_name"`
	OtherNames       *string `form:"other_names" json:"other_names"`
	Gender           *string `form:"gender" json:"gender"`
	Nationality      string  `form:"nationality" json:"nationality"`
	NationalIdNumber string  `form:"national_id_number" json:"national_id_number"`
	DateOfBirth      *string `form:"date_of_birth" json:"date_of_birth"`
	Email            *string `form:"email" json:"email"`
	PhoneNumber      string  `form:"phone_number" json:"phone_number"`
	IsIntern         bool    `form:"is_intern" json:"is_intern"`
	IsPartTime       bool    `form:"is_part_time" json:"is_part_time"`
	SmeId            int     `form:"sme_id" json:"sme_id"`
}

// Rules defines validation rules for additional_business_member creation
func (r *AdditionalBusinessMemberCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"first_name":         "required|max_len:100",
		"last_name":          "required|max_len:100",
		"other_names":        "max_len:100",
		"gender":             "required",
		"nationality":        "required|max_len:100",
		"national_id_number": "max_len:50",
		// Note: date_of_birth validation removed because the 'date' validator doesn't work with *string
		// Date parsing/validation happens in ToCreateData() method instead
		// Note: email validation handled in PrepareForValidation since 'nullable' validator doesn't exist
		"phone_number": "required|max_len:20",
		"is_intern":    "boolean",
		"is_part_time": "boolean",
		"sme_id":       "required|numeric",
	}
}

// Messages defines custom validation messages
func (r *AdditionalBusinessMemberCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"first_name.required":        "First Name is required",
		"first_name.max_len":         "First Name cannot exceed 100 characters",
		"last_name.required":         "Last Name is required",
		"last_name.max_len":          "Last Name cannot exceed 100 characters",
		"other_names.max_len":        "Other Names cannot exceed 100 characters",
		"gender.required":            "Gender is required",
		"nationality.required":       "Nationality is required",
		"nationality.max_len":        "Nationality cannot exceed 100 characters",
		"national_id_number.max_len": "National ID Number cannot exceed 50 characters",
		"date_of_birth.date":         "Date of Birth must be a valid date",
		"email.email":                "Email must be a valid email address",
		"email.max_len":              "Email cannot exceed 100 characters",
		"phone_number.required":      "Phone Number is required",
		"phone_number.max_len":       "Phone Number cannot exceed 20 characters",
		"is_intern.boolean":          "Is Intern must be true or false",
		"is_part_time.boolean":       "Is Part Time must be true or false",
		"sme_id.required":            "MSME ID is required",
		"sme_id.numeric":             "MSME ID must be a number",
	}
}

// Attributes defines custom attribute names
func (r *AdditionalBusinessMemberCreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to make this request
func (r *AdditionalBusinessMemberCreateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user has permission to create additional_business_member
	// return facades.Gate().Allows("create.additional_business_members", ctx)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *AdditionalBusinessMemberCreateRequest) PrepareForValidation(ctx http.Context) error {
	// Validate date_of_birth format if provided
	if r.DateOfBirth != nil && *r.DateOfBirth != "" {
		parsedDate := carbon.Parse(*r.DateOfBirth)
		if parsedDate.Error != nil {
			return errors.New("date_of_birth: Date of Birth must be a valid date")
		}
	}

	// Validate email format if provided (since 'nullable' validator doesn't exist)
	if r.Email != nil && *r.Email != "" {
		// Basic email validation using regex
		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		matched, _ := regexp.MatchString(emailRegex, *r.Email)
		if !matched {
			return errors.New("email: Email must be a valid email address")
		}
		if len(*r.Email) > 100 {
			return errors.New("email: Email cannot exceed 100 characters")
		}
	}

	return nil
}

// PassedValidation is called after validation passes
func (r *AdditionalBusinessMemberCreateRequest) PassedValidation(ctx http.Context) error {
	// TODO: Add post-validation logic if needed
	return nil
}

// ToCreateData converts the request to create data map
func (r *AdditionalBusinessMemberCreateRequest) ToCreateData() map[string]interface{} {
	data := map[string]interface{}{
		"first_name":         r.FirstName,
		"last_name":          r.LastName,
		"other_names":        r.OtherNames,
		"gender":             r.Gender,
		"nationality":        r.Nationality,
		"national_id_number": r.NationalIdNumber,
		"email":              r.Email,
		"phone_number":       r.PhoneNumber,
		"is_intern":          r.IsIntern,
		"is_part_time":       r.IsPartTime,
		"sme_id":             r.SmeId,
	}

	// Convert date_of_birth string to carbon.DateTime if provided
	if r.DateOfBirth != nil && *r.DateOfBirth != "" {
		parsedDate := carbon.Parse(*r.DateOfBirth)
		if parsedDate.Error == nil {
			data["date_of_birth"] = carbon.NewDateTime(parsedDate)
		} else {
			data["date_of_birth"] = nil
		}
	} else {
		data["date_of_birth"] = nil
	}

	return data
}
