package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// AuthorCreateRequest handles author creation validation
type AuthorCreateRequest struct {
	FirstName   string  `form:"first_name" json:"first_name"`
	LastName    string  `form:"last_name" json:"last_name"`
	Bio         *string `form:"bio" json:"bio"`
	Email       *string `form:"email" json:"email"`
	Website     *string `form:"website" json:"website"`
	BirthDate   *string `form:"birth_date" json:"birth_date"`
	Nationality *string `form:"nationality" json:"nationality"`
	PhotoURL    *string `form:"photo_url" json:"photo_url"`
	Status      string  `form:"status" json:"status"`
}

// Rules defines validation rules for author creation
// Keys must match ToCreateData() output keys (camelCase) since the controller
// validates the data map, not the struct fields directly.
func (r *AuthorCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"firstName": "required|max_len:100",
		"lastName":  "required|max_len:100",
		"status":    "in:ACTIVE,INACTIVE",
	}
}

// Messages defines custom validation messages
func (r *AuthorCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"firstName.required": "First name is required",
		"firstName.max_len":  "First name cannot exceed 100 characters",
		"lastName.required":  "Last name is required",
		"lastName.max_len":   "Last name cannot exceed 100 characters",
		"email.max_len":      "Email cannot exceed 255 characters",
		"website.max_len":    "Website cannot exceed 255 characters",
		"nationality.max":    "Nationality cannot exceed 100 characters",
		"photoUrl.max_len":   "Photo URL cannot exceed 500 characters",
		"status.in":          "Status must be one of: ACTIVE, INACTIVE",
	}
}

// Attributes defines custom attribute names
func (r *AuthorCreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"firstName": "first name",
		"lastName":  "last name",
		"birthDate": "birth date",
		"photoUrl":  "photo URL",
	}
}

// Authorize determines if the user is authorized to make this request
func (r *AuthorCreateRequest) Authorize(ctx http.Context) error {
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *AuthorCreateRequest) PrepareForValidation(ctx http.Context) error {
	// Set default status if not provided
	if r.Status == "" {
		r.Status = "ACTIVE"
	}
	return nil
}

// PassedValidation is called after validation passes
func (r *AuthorCreateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// ToCreateData converts the request to create data map
// Keys use camelCase to match model json tags for setFieldsRecursively
func (r *AuthorCreateRequest) ToCreateData() map[string]interface{} {
	data := map[string]interface{}{
		"firstName": r.FirstName,
		"lastName":  r.LastName,
		"status":    r.Status,
	}

	// Only include optional fields if they have values
	if r.Bio != nil && *r.Bio != "" {
		data["bio"] = *r.Bio
	}
	if r.Email != nil && *r.Email != "" {
		data["email"] = *r.Email
	}
	if r.Website != nil && *r.Website != "" {
		data["website"] = *r.Website
	}
	if r.BirthDate != nil && *r.BirthDate != "" {
		data["birthDate"] = *r.BirthDate
	}
	if r.Nationality != nil && *r.Nationality != "" {
		data["nationality"] = *r.Nationality
	}
	if r.PhotoURL != nil && *r.PhotoURL != "" {
		data["photoUrl"] = *r.PhotoURL
	}

	return data
}

// AuthorUpdateRequest handles author update validation
type AuthorUpdateRequest struct {
	FirstName   *string `form:"first_name" json:"first_name"`
	LastName    *string `form:"last_name" json:"last_name"`
	Bio         *string `form:"bio" json:"bio"`
	Email       *string `form:"email" json:"email"`
	Website     *string `form:"website" json:"website"`
	BirthDate   *string `form:"birth_date" json:"birth_date"`
	Nationality *string `form:"nationality" json:"nationality"`
	PhotoURL    *string `form:"photo_url" json:"photo_url"`
	Status      *string `form:"status" json:"status"`
	ID          uint    `form:"-" json:"-"` // Set by controller
}

// Rules defines validation rules for author updates
// Keys must match ToUpdateData() output keys (snake_case) since the controller
// validates the data map, not the struct fields directly.
func (r *AuthorUpdateRequest) Rules(ctx http.Context) map[string]string {
	rules := map[string]string{}

	// Only validate fields that are provided
	if r.FirstName != nil {
		rules["first_name"] = "required|max_len:100"
	}
	if r.LastName != nil {
		rules["last_name"] = "required|max_len:100"
	}
	if r.Status != nil {
		rules["status"] = "in:ACTIVE,INACTIVE"
	}

	// If no rules were added, add a dummy rule to prevent empty rules error
	if len(rules) == 0 {
		rules["_at_least_one_field"] = "sometimes"
	}

	return rules
}

// Messages defines custom validation messages for updates
func (r *AuthorUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"first_name.required": "First name is required",
		"first_name.max_len":  "First name cannot exceed 100 characters",
		"last_name.required":  "Last name is required",
		"last_name.max_len":   "Last name cannot exceed 100 characters",
		"email.max_len":       "Email cannot exceed 255 characters",
		"website.max_len":     "Website cannot exceed 255 characters",
		"nationality.max":     "Nationality cannot exceed 100 characters",
		"photo_url.max_len":   "Photo URL cannot exceed 500 characters",
		"status.in":           "Status must be one of: ACTIVE, INACTIVE",
	}
}

// Attributes defines custom attribute names for updates
func (r *AuthorUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"first_name": "first name",
		"last_name":  "last name",
		"birth_date": "birth date",
		"photo_url":  "photo URL",
	}
}

// Authorize determines if the user is authorized to update this author
func (r *AuthorUpdateRequest) Authorize(ctx http.Context) error {
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *AuthorUpdateRequest) PrepareForValidation(ctx http.Context) error {
	return nil
}

// PassedValidation is called after validation passes
func (r *AuthorUpdateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// ToUpdateData converts the request to update data map
// Keys use snake_case to match DB column names for GORM's Update()
func (r *AuthorUpdateRequest) ToUpdateData() map[string]interface{} {
	data := map[string]interface{}{}

	// Only include fields that are provided (not nil)
	if r.FirstName != nil {
		data["first_name"] = *r.FirstName
	}
	if r.LastName != nil {
		data["last_name"] = *r.LastName
	}
	if r.Bio != nil {
		data["bio"] = *r.Bio
	}
	if r.Email != nil {
		data["email"] = *r.Email
	}
	if r.Website != nil {
		data["website"] = *r.Website
	}
	if r.BirthDate != nil {
		data["birth_date"] = *r.BirthDate
	}
	if r.Nationality != nil {
		data["nationality"] = *r.Nationality
	}
	if r.PhotoURL != nil {
		data["photo_url"] = *r.PhotoURL
	}
	if r.Status != nil {
		data["status"] = *r.Status
	}

	return data
}

// GetResourceID returns the resource ID for update
func (r *AuthorUpdateRequest) GetResourceID() interface{} {
	return r.ID
}
