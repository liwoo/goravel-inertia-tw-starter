package requests

import (
	"errors"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"
)

// EventUpdateRequest handles event update validation
type EventUpdateRequest struct {
	Title         *string   `form:"title" json:"title"`
	Description   *string   `form:"description" json:"description"`
	Date          *string   `form:"date" json:"date"` // String for validation, convert to carbon.DateTime
	EndDate       *string   `form:"end_date" json:"end_date"`
	Venue         *string   `form:"venue" json:"venue"`
	Partners      *[]string `form:"partners" json:"partners"`
	District      *string   `form:"district" json:"district"`
	AttendingSmes *[]int    `form:"attending_smes" json:"attending_smes"`
	Notes         *string   `form:"notes" json:"notes"`
	ID            uint      `form:"-" json:"-"` // Set by controller
}

// Rules defines validation rules for event updates
func (r *EventUpdateRequest) Rules(ctx http.Context) map[string]string {
	rules := map[string]string{}

	// Only validate Title if provided
	if r.Title != nil {
		rules["title"] = "required|max_len:255"
	}
	// Only validate Description if provided
	if r.Description != nil {
		rules["description"] = "required"
	}
	// Only validate Date if provided (date validation in PrepareForValidation)
	if r.Date != nil {
		rules["date"] = "required"
	}
	// Only validate Venue if provided
	if r.Venue != nil {
		rules["venue"] = "required|max_len:255"
	}
	// Only validate District if provided
	if r.District != nil {
		rules["district"] = "required|max_len:100"
	}

	// If no rules were added, add a dummy rule to prevent empty rules error
	if len(rules) == 0 {
		rules["_at_least_one_field"] = "sometimes"
	}

	return rules
}

// Messages defines custom validation messages for updates
func (r *EventUpdateRequest) Messages(ctx http.Context) map[string]string {
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

// Attributes defines custom attribute names for updates
func (r *EventUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to update this event
func (r *EventUpdateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user can update this specific event
	// return facades.Gate().Allows("update.events", event)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *EventUpdateRequest) PrepareForValidation(ctx http.Context) error {
	// Validate date format if provided
	if r.Date != nil && *r.Date != "" {
		parsedDate := carbon.Parse(*r.Date)
		if parsedDate.Error != nil {
			return errors.New("date: Event date must be a valid date")
		}
	}
	if r.EndDate != nil && *r.EndDate != "" {
		parsedDate := carbon.Parse(*r.EndDate)
		if parsedDate.Error != nil {
			return errors.New("end_date: Event end date must be a valid date")
		}
	}
	return nil
}

// PassedValidation is called after validation passes
func (r *EventUpdateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// ToUpdateData converts the request to update data map
func (r *EventUpdateRequest) ToUpdateData() map[string]interface{} {
	data := map[string]interface{}{}

	// Only include Title if provided
	if r.Title != nil {
		data["title"] = *r.Title
	}
	// Only include Description if provided
	if r.Description != nil {
		data["description"] = *r.Description
	}
	// Only include Date if provided
	if r.Date != nil && *r.Date != "" {
		parsedDate := carbon.Parse(*r.Date)
		if parsedDate.Error == nil {
			data["date"] = *carbon.NewDateTime(parsedDate)
		}
	}
	// Only include EndDate if provided
	if r.EndDate != nil && *r.EndDate != "" {
		parsedDate := carbon.Parse(*r.EndDate)
		if parsedDate.Error == nil {
			data["end_date"] = *carbon.NewDateTime(parsedDate)
		}
	}
	// Only include Venue if provided
	if r.Venue != nil {
		data["venue"] = *r.Venue
	}
	// Only include Partners if provided
	if r.Partners != nil {
		data["partners"] = *r.Partners
	}
	// Only include District if provided
	if r.District != nil {
		data["district"] = *r.District
	}
	// Only include AttendingSmes if provided
	if r.AttendingSmes != nil {
		data["attending_smes"] = *r.AttendingSmes
	}
	// Only include Notes if provided
	if r.Notes != nil {
		data["notes"] = *r.Notes
	}

	return data
}

// GetResourceID returns the resource ID for update
func (r *EventUpdateRequest) GetResourceID() interface{} {
	return r.ID
}
