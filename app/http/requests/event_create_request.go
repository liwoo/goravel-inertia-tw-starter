package requests

import (
	"errors"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"
)

// EventCreateRequest handles event creation validation
type EventCreateRequest struct {
	Title         string   `form:"title" json:"title"`
	Description   string   `form:"description" json:"description"`
	Date          string   `form:"date" json:"date"` // String for validation, convert to carbon.DateTime
	Venue         string   `form:"venue" json:"venue"`
	Partners      []string `form:"partners" json:"partners"`
	District      string   `form:"district" json:"district"`
	AttendingSmes []int    `form:"attending_smes" json:"attending_smes"`
	Notes         *string  `form:"notes" json:"notes"`
}

// Rules defines validation rules for event creation
func (r *EventCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"title":       "required|max_len:255",
		"description": "required",
		"date":        "required", // Date validation in PrepareForValidation
		"venue":       "required|max_len:255",
		"district":    "required|max_len:100",
		// partners and attending_smes are arrays, no validation needed
	}
}

// Messages defines custom validation messages
func (r *EventCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"title.required":       "Event title is required",
		"title.max_len":        "Event title cannot exceed 255 characters",
		"description.required": "Event description is required",
		"date.required":        "Event date is required",
		"date.date":            "Event date must be a valid date",
		"venue.required":       "Event venue is required",
		"venue.max_len":        "Event venue cannot exceed 255 characters",
		"district.required":    "Event district is required",
		"district.max_len":     "District cannot exceed 100 characters",
	}
}

// Attributes defines custom attribute names
func (r *EventCreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to make this request
func (r *EventCreateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user has permission to create events
	// return facades.Gate().Allows("create.events", ctx)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *EventCreateRequest) PrepareForValidation(ctx http.Context) error {
	// Validate date format
	if r.Date != "" {
		parsedDate := carbon.Parse(r.Date)
		if parsedDate.Error != nil {
			return errors.New("date: Event date must be a valid date")
		}
	}
	return nil
}

// PassedValidation is called after validation passes
func (r *EventCreateRequest) PassedValidation(ctx http.Context) error {
	// TODO: Add post-validation logic if needed
	return nil
}

// ToCreateData converts the request to create data map
func (r *EventCreateRequest) ToCreateData() map[string]interface{} {
	data := map[string]interface{}{
		"title":       r.Title,
		"description": r.Description,
		"venue":       r.Venue,
		"partners":    r.Partners,
		"district":    r.District,
		"notes":       r.Notes,
	}

	// Convert date string to carbon.DateTime
	if r.Date != "" {
		parsedDate := carbon.Parse(r.Date)
		if parsedDate.Error == nil {
			data["date"] = *carbon.NewDateTime(parsedDate)
		}
	}

	// Handle attending_smes - convert to JSON
	if len(r.AttendingSmes) > 0 {
		data["attending_smes"] = r.AttendingSmes
	} else {
		data["attending_smes"] = []int{} // Empty array if not provided
	}

	return data
}
