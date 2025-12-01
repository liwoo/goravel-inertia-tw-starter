package requests

import (
	"errors"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"
)

// ProcurementNoticeUpdateRequest handles procurement notice update validation
type ProcurementNoticeUpdateRequest struct {
	ID                     uint     `form:"-" json:"-"` // Set by controller
	ProcuredBy             *string  `form:"procured_by" json:"procured_by"`
	ProcurementType        *string  `form:"procurement_type" json:"procurement_type"`
	MarketApproach         *string  `form:"market_approach" json:"market_approach"`
	Invitation             *string  `form:"invitation" json:"invitation"`
	OpenDate               *string  `form:"open_date" json:"open_date"`   // String for validation, convert to carbon.DateTime
	CloseDate              *string  `form:"close_date" json:"close_date"` // String for validation, convert to carbon.DateTime
	Partners               []string `form:"partners" json:"partners"`
	QualifyingDistricts    []string `form:"qualifying_districts" json:"qualifying_districts"`
	IsPublished            *bool    `form:"is_published" json:"is_published"`
	Organization           *string  `form:"organization" json:"organization"`
	Classification         []string `form:"classification" json:"classification"`
	InterestedSmes         []string `form:"interested_smes" json:"interested_smes"`
	Details                *string  `form:"details" json:"details"`
	ApplicationDetails     *string  `form:"application_details" json:"application_details"`
	MinimumQualifyingScore *int     `form:"minimum_qualifying_score" json:"minimum_qualifying_score"`
}

// Rules defines validation rules for procurement notice updates
func (r *ProcurementNoticeUpdateRequest) Rules(ctx http.Context) map[string]string {
	rules := map[string]string{}

	// Only validate fields that are provided
	if r.ProcuredBy != nil {
		rules["procured_by"] = "required|max_len:255"
	}
	if r.ProcurementType != nil {
		rules["procurement_type"] = "required|max_len:255"
	}
	if r.MarketApproach != nil {
		rules["market_approach"] = "required|in:National,International"
	}
	if r.Invitation != nil {
		rules["invitation"] = "required|in:open,limited,single-source"
	}
	if r.OpenDate != nil {
		rules["open_date"] = "required"
	}
	if r.CloseDate != nil {
		rules["close_date"] = "required"
	}
	if r.Organization != nil {
		rules["organization"] = "required|max_len:255"
	}
	if r.Details != nil {
		rules["details"] = "required"
	}
	if r.ApplicationDetails != nil {
		rules["application_details"] = "required"
	}
	if r.MinimumQualifyingScore != nil {
		rules["minimum_qualifying_score"] = "required|numeric|min:0|max:100"
	}

	// Array validation
	rules["partners"] = "array"
	rules["partners.*"] = "max_len:255"
	rules["qualifying_districts"] = "array"
	rules["qualifying_districts.*"] = "max_len:255"
	rules["classification"] = "array"
	rules["classification.*"] = "max_len:255"
	rules["interested_smes"] = "array"
	rules["interested_smes.*"] = "max_len:255"

	// If no rules were added, add a dummy rule to prevent empty rules error
	if len(rules) == 0 {
		rules["_at_least_one_field"] = "sometimes"
	}

	return rules
}

// Messages defines custom validation messages for updates
func (r *ProcurementNoticeUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"procured_by.required":              "Procured By is required",
		"procured_by.max_len":               "Procured By cannot exceed 255 characters",
		"procurement_type.required":         "Procurement Type is required",
		"procurement_type.max_len":          "Procurement Type cannot exceed 255 characters",
		"market_approach.required":          "Market Approach is required",
		"market_approach.in":                "Market Approach must be National or International",
		"invitation.required":               "Invitation type is required",
		"invitation.in":                     "Invitation must be open, limited, or single-source",
		"open_date.required":                "Open Date is required",
		"close_date.required":               "Close Date is required",
		"organization.required":             "Organization is required",
		"organization.max_len":              "Organization cannot exceed 255 characters",
		"details.required":                  "Details are required",
		"application_details.required":      "Application Details are required",
		"minimum_qualifying_score.required": "Minimum Qualifying Score is required",
		"minimum_qualifying_score.numeric":  "Minimum Qualifying Score must be a number",
		"minimum_qualifying_score.min":      "Minimum Qualifying Score must be at least 0",
		"minimum_qualifying_score.max":      "Minimum Qualifying Score cannot exceed 100",
	}
}

// Attributes defines custom attribute names for updates
func (r *ProcurementNoticeUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"procured_by":              "Procured By",
		"procurement_type":         "Procurement Type",
		"market_approach":          "Market Approach",
		"invitation":               "Invitation Type",
		"open_date":                "Open Date",
		"close_date":               "Close Date",
		"partners":                 "Partners",
		"qualifying_districts":     "Qualifying Districts",
		"is_published":             "Published Status",
		"organization":             "Organization",
		"classification":           "Classification",
		"interested_smes":          "Interested SMEs",
		"details":                  "Details",
		"application_details":      "Application Details",
		"minimum_qualifying_score": "Minimum Qualifying Score",
	}
}

// Authorize determines if the user is authorized to update this procurement notice
func (r *ProcurementNoticeUpdateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user can update this specific procurement notice
	// return facades.Gate().Allows("update.procurement_notices", procurement_notice)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *ProcurementNoticeUpdateRequest) PrepareForValidation(ctx http.Context) error {
	// Validate date formats if provided
	if r.OpenDate != nil && *r.OpenDate != "" {
		parsedDate := carbon.Parse(*r.OpenDate)
		if parsedDate.Error != nil {
			return errors.New("open_date: Open Date must be a valid date")
		}
	}

	if r.CloseDate != nil && *r.CloseDate != "" {
		parsedDate := carbon.Parse(*r.CloseDate)
		if parsedDate.Error != nil {
			return errors.New("close_date: Close Date must be a valid date")
		}
	}

	return nil
}

// PassedValidation is called after validation passes
func (r *ProcurementNoticeUpdateRequest) PassedValidation(ctx http.Context) error {
	// Validate that close_date is after open_date if both provided
	if r.OpenDate != nil && *r.OpenDate != "" && r.CloseDate != nil && *r.CloseDate != "" {
		openDate := carbon.Parse(*r.OpenDate)
		closeDate := carbon.Parse(*r.CloseDate)

		if openDate.Error == nil && closeDate.Error == nil {
			if closeDate.Lt(openDate) {
				return errors.New("Close Date must be after Open Date")
			}
		}
	}

	return nil
}

// ToUpdateData converts the request to update data map
func (r *ProcurementNoticeUpdateRequest) ToUpdateData() map[string]interface{} {
	data := map[string]interface{}{}

	// Only update fields that are provided
	if r.ProcuredBy != nil {
		data["procured_by"] = *r.ProcuredBy
	}
	if r.ProcurementType != nil {
		data["procurement_type"] = *r.ProcurementType
	}
	if r.MarketApproach != nil {
		data["market_approach"] = *r.MarketApproach
	}
	if r.Invitation != nil {
		data["invitation"] = *r.Invitation
	}
	if r.IsPublished != nil {
		data["is_published"] = *r.IsPublished
	}
	if r.Organization != nil {
		data["organization"] = *r.Organization
	}
	if r.Details != nil {
		data["details"] = *r.Details
	}
	if r.ApplicationDetails != nil {
		data["application_details"] = *r.ApplicationDetails
	}
	if r.MinimumQualifyingScore != nil {
		data["minimum_qualifying_score"] = *r.MinimumQualifyingScore
	}

	// Convert date strings to carbon.DateTime if provided
	if r.OpenDate != nil && *r.OpenDate != "" {
		parsedDate := carbon.Parse(*r.OpenDate)
		if parsedDate.Error == nil {
			data["open_date"] = *carbon.NewDateTime(parsedDate)
		}
	}

	if r.CloseDate != nil && *r.CloseDate != "" {
		parsedDate := carbon.Parse(*r.CloseDate)
		if parsedDate.Error == nil {
			data["close_date"] = *carbon.NewDateTime(parsedDate)
		}
	}

	// Handle array fields - always update if provided (even if empty)
	// This allows clearing arrays by sending empty arrays
	if r.Partners != nil {
		data["partners"] = r.Partners
	}

	if r.QualifyingDistricts != nil {
		data["qualifying_districts"] = r.QualifyingDistricts
	}

	if r.Classification != nil {
		data["classification"] = r.Classification
	}

	if r.InterestedSmes != nil {
		data["interested_smes"] = r.InterestedSmes
	}

	return data
}

// GetResourceID returns the resource ID for update
func (r *ProcurementNoticeUpdateRequest) GetResourceID() interface{} {
	return r.ID
}
