package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// BusinessFormalisationUpdateRequest handles business_formalisation update validation
type BusinessFormalisationUpdateRequest struct {
	SmeID                  *int     `form:"sme_id" json:"sme_id"`
	HasBankAccount         *bool    `form:"has_bank_account" json:"has_bank_account"`
	HasTaxClarification    *bool    `form:"has_tax_clarification" json:"has_tax_clarification"`
	IsRegisteredForVat     *bool    `form:"is_registered_for_vat" json:"is_registered_for_vat"`
	IsMemberOfAssociation  *bool    `form:"is_member_of_association" json:"is_member_of_association"`
	IsAffiliated           *bool    `form:"is_affiliated" json:"is_affiliated"`
	HasExportLicense       *bool    `form:"has_export_license" json:"has_export_license"`
	HasAccessedBds         *bool    `form:"has_accessed_bds" json:"has_accessed_bds"`
	AnnualTurnover         *float64 `form:"annual_turnover" json:"annual_turnover"`
	EstimatedValueOfAssets *float64 `form:"estimated_value_of_assets" json:"estimated_value_of_assets"`
	ID                     uint     `form:"-" json:"-"` // Set by controller
	// NOTE: formalisation_score is intentionally excluded - it's read-only and calculated internally
}

// Rules defines validation rules for business_formalisation updates
func (r *BusinessFormalisationUpdateRequest) Rules(ctx http.Context) map[string]string {
	rules := map[string]string{}

	// Only validate SmeID if provided
	if r.SmeID != nil {
		rules["sme_id"] = "required|numeric"
	}
	// Only validate HasBankAccount if provided
	if r.HasBankAccount != nil {
		rules["has_bank_account"] = "boolean"
	}
	// Only validate HasTaxClarification if provided
	if r.HasTaxClarification != nil {
		rules["has_tax_clarification"] = "boolean"
	}
	// Only validate IsRegisteredForVat if provided
	if r.IsRegisteredForVat != nil {
		rules["is_registered_for_vat"] = "boolean"
	}
	// Only validate IsMemberOfAssociation if provided
	if r.IsMemberOfAssociation != nil {
		rules["is_member_of_association"] = "boolean"
	}
	// Only validate IsAffiliated if provided
	if r.IsAffiliated != nil {
		rules["is_affiliated"] = "boolean"
	}
	// Only validate HasExportLicense if provided
	if r.HasExportLicense != nil {
		rules["has_export_license"] = "boolean"
	}
	// Only validate HasAccessedBds if provided
	if r.HasAccessedBds != nil {
		rules["has_accessed_bds"] = "boolean"
	}
	// AnnualTurnover and EstimatedValueOfAssets are float64 and don't need string validation

	// If no rules were added, add a dummy rule to prevent empty rules error
	if len(rules) == 0 {
		rules["_at_least_one_field"] = "sometimes"
	}

	return rules
}

// Messages defines custom validation messages for updates
func (r *BusinessFormalisationUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"sme_id.required":                   "MSME ID is required",
		"sme_id.numeric":                    "MSME ID must be a number",
		"has_bank_account.boolean":          "Has Bank Account must be true or false",
		"has_tax_clarification.boolean":     "Has Tax Clarification must be true or false",
		"is_registered_for_vat.boolean":     "Is Registered for VAT must be true or false",
		"is_member_of_association.boolean":  "Is Member of Association must be true or false",
		"is_affiliated.boolean":             "Is Affiliated must be true or false",
		"has_export_license.boolean":        "Has Export License must be true or false",
		"has_accessed_bds.boolean":          "Has Accessed BDS must be true or false",
		"annual_turnover.numeric":           "Annual Turnover must be a number",
		"estimated_value_of_assets.numeric": "Estimated Value of Assets must be a number",
	}
}

// Attributes defines custom attribute names for updates
func (r *BusinessFormalisationUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to update this business_formalisation
func (r *BusinessFormalisationUpdateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user can update this specific business_formalisation
	// return facades.Gate().Allows("update.business_formalisations", business_formalisation)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *BusinessFormalisationUpdateRequest) PrepareForValidation(ctx http.Context) error {
	// TODO: Add data preparation logic for updates
	// Example: Normalize data if provided
	return nil
}

// PassedValidation is called after validation passes
func (r *BusinessFormalisationUpdateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// ToUpdateData converts the request to update data map
func (r *BusinessFormalisationUpdateRequest) ToUpdateData() map[string]interface{} {
	data := map[string]interface{}{}

	// Only include SmeID if provided
	if r.SmeID != nil {
		data["sme_id"] = *r.SmeID
	}
	// Only include HasBankAccount if provided
	if r.HasBankAccount != nil {
		data["has_bank_account"] = *r.HasBankAccount
	}
	// Only include HasTaxClarification if provided
	if r.HasTaxClarification != nil {
		data["has_tax_clarification"] = *r.HasTaxClarification
	}
	// Only include IsRegisteredForVat if provided
	if r.IsRegisteredForVat != nil {
		data["is_registered_for_vat"] = *r.IsRegisteredForVat
	}
	// Only include IsMemberOfAssociation if provided
	if r.IsMemberOfAssociation != nil {
		data["is_member_of_association"] = *r.IsMemberOfAssociation
	}
	// Only include IsAffiliated if provided
	if r.IsAffiliated != nil {
		data["is_affiliated"] = *r.IsAffiliated
	}
	// Only include HasExportLicense if provided
	if r.HasExportLicense != nil {
		data["has_export_license"] = *r.HasExportLicense
	}
	// Only include HasAccessedBds if provided
	if r.HasAccessedBds != nil {
		data["has_accessed_bds"] = *r.HasAccessedBds
	}
	// Only include AnnualTurnover if provided
	if r.AnnualTurnover != nil {
		data["annual_turnover"] = *r.AnnualTurnover
	}
	// Only include EstimatedValueOfAssets if provided
	if r.EstimatedValueOfAssets != nil {
		data["estimated_value_of_assets"] = *r.EstimatedValueOfAssets
	}
	// formalisation_score is intentionally excluded - it's read-only

	return data
}

// GetResourceID returns the resource ID for update
func (r *BusinessFormalisationUpdateRequest) GetResourceID() interface{} {
	return r.ID
}
