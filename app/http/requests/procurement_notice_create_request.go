package requests

import (
	"errors"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"
)

// ProcurementNoticeCreateRequest handles procurement notice creation validation
type ProcurementNoticeCreateRequest struct {
	ProcuredBy             string   `form:"procured_by" json:"procured_by"`
	ProcurementType        string   `form:"procurement_type" json:"procurement_type"`
	MarketApproach         string   `form:"market_approach" json:"market_approach"`
	Invitation             string   `form:"invitation" json:"invitation"`
	OpenDate               string   `form:"open_date" json:"open_date"`   // String for validation, convert to carbon.DateTime
	CloseDate              string   `form:"close_date" json:"close_date"` // String for validation, convert to carbon.DateTime
	Partners               []string `form:"partners" json:"partners"`
	QualifyingDistricts    []string `form:"qualifying_districts" json:"qualifying_districts"`
	IsPublished            bool     `form:"is_published" json:"is_published"`
	Organization           string   `form:"organization" json:"organization"`
	Classification         []string `form:"classification" json:"classification"`
	InterestedSmes         []string `form:"interested_smes" json:"interested_smes"`
	Details                string   `form:"details" json:"details"`
	ApplicationDetails     string   `form:"application_details" json:"application_details"`
	MinimumQualifyingScore int      `form:"minimum_qualifying_score" json:"minimum_qualifying_score"`
}

// Rules defines validation rules for procurement notice creation
func (r *ProcurementNoticeCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"procured_by":              "required|max_len:255",
		"procurement_type":         "required|max_len:255",
		"market_approach":          "required|in:National,International",
		"invitation":               "required|in:open,limited,single-source",
		"open_date":                "required",
		"close_date":               "required",
		"partners":                 "array",
		"partners.*":               "max_len:255",
		"qualifying_districts":     "array",
		"qualifying_districts.*":   "max_len:255",
		"is_published":             "bool",
		"organization":             "required|max_len:255",
		"classification":           "array",
		"classification.*":         "max_len:255",
		"interested_smes":          "array",
		"interested_smes.*":        "max_len:255",
		"details":                  "required",
		"application_details":      "required",
		"minimum_qualifying_score": "required|numeric|min:0|max:100",
	}
}

// Messages defines custom validation messages
func (r *ProcurementNoticeCreateRequest) Messages(ctx http.Context) map[string]string {
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

// Attributes defines custom attribute names
func (r *ProcurementNoticeCreateRequest) Attributes(ctx http.Context) map[string]string {
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

// Authorize determines if the user is authorized to make this request
func (r *ProcurementNoticeCreateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user has permission to create procurement notices
	// return facades.Gate().Allows("create.procurement_notices", ctx)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *ProcurementNoticeCreateRequest) PrepareForValidation(ctx http.Context) error {
	// Validate date formats
	if r.OpenDate != "" {
		parsedDate := carbon.Parse(r.OpenDate)
		if parsedDate.Error != nil {
			return errors.New("open_date: Open Date must be a valid date")
		}
	}

	if r.CloseDate != "" {
		parsedDate := carbon.Parse(r.CloseDate)
		if parsedDate.Error != nil {
			return errors.New("close_date: Close Date must be a valid date")
		}
	}

	return nil
}

// PassedValidation is called after validation passes
func (r *ProcurementNoticeCreateRequest) PassedValidation(ctx http.Context) error {
	// Validate that close_date is after open_date
	if r.OpenDate != "" && r.CloseDate != "" {
		openDate := carbon.Parse(r.OpenDate)
		closeDate := carbon.Parse(r.CloseDate)

		if openDate.Error == nil && closeDate.Error == nil {
			if closeDate.Lt(openDate) {
				return errors.New("Close Date must be after Open Date")
			}
		}
	}

	return nil
}

// ToCreateData converts the request to create data map
func (r *ProcurementNoticeCreateRequest) ToCreateData() map[string]interface{} {
	data := map[string]interface{}{
		"procured_by":              r.ProcuredBy,
		"procurement_type":         r.ProcurementType,
		"market_approach":          r.MarketApproach,
		"invitation":               r.Invitation,
		"is_published":             r.IsPublished,
		"organization":             r.Organization,
		"details":                  r.Details,
		"application_details":      r.ApplicationDetails,
		"minimum_qualifying_score": r.MinimumQualifyingScore,
	}

	// Convert date strings to carbon.DateTime
	if r.OpenDate != "" {
		parsedDate := carbon.Parse(r.OpenDate)
		if parsedDate.Error == nil {
			data["open_date"] = *carbon.NewDateTime(parsedDate)
		}
	}

	if r.CloseDate != "" {
		parsedDate := carbon.Parse(r.CloseDate)
		if parsedDate.Error == nil {
			data["close_date"] = *carbon.NewDateTime(parsedDate)
		}
	}

	// Initialize array fields - always provide empty arrays if not provided
	if len(r.Partners) > 0 {
		data["partners"] = r.Partners
	} else {
		data["partners"] = []string{}
	}

	if len(r.QualifyingDistricts) > 0 {
		data["qualifying_districts"] = r.QualifyingDistricts
	} else {
		data["qualifying_districts"] = []string{}
	}

	if len(r.Classification) > 0 {
		data["classification"] = r.Classification
	} else {
		data["classification"] = []string{}
	}

	if len(r.InterestedSmes) > 0 {
		data["interested_smes"] = r.InterestedSmes
	} else {
		data["interested_smes"] = []string{}
	}

	return data
}
