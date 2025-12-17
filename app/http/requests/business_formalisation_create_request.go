package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// BusinessFormalisationCreateRequest handles business_formalisation creation validation
type BusinessFormalisationCreateRequest struct {
	SmeID                 int     `form:"sme_id" json:"sme_id"`
	HasBankAccount        bool    `form:"has_bank_account" json:"has_bank_account"`
	HasTaxClarification   bool    `form:"has_tax_clarification" json:"has_tax_clarification"`
	IsRegisteredForVat    bool    `form:"is_registered_for_vat" json:"is_registered_for_vat"`
	IsMemberOfAssociation bool    `form:"is_member_of_association" json:"is_member_of_association"`
	IsAffiliated          bool    `form:"is_affiliated" json:"is_affiliated"`
	HasExportLicense      bool    `form:"has_export_license" json:"has_export_license"`
	HasAccessedBds        bool    `form:"has_accessed_bds" json:"has_accessed_bds"`
	AnnualTurnover        float64 `form:"annual_turnover" json:"annual_turnover"`
	EstimatedValueOfAssets float64 `form:"estimated_value_of_assets" json:"estimated_value_of_assets"`
	// NOTE: formalisation_score is intentionally excluded - it's read-only and calculated internally
}

// Rules defines validation rules for business_formalisation creation
func (r *BusinessFormalisationCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"sme_id":                   "required|numeric",
		"has_bank_account":         "boolean",
		"has_tax_clarification":    "boolean",
		"is_registered_for_vat":    "boolean",
		"is_member_of_association": "boolean",
		"is_affiliated":            "boolean",
		"has_export_license":       "boolean",
		"has_accessed_bds":         "boolean",
		// Note: annual_turnover and estimated_value_of_assets are float64 and don't need string validation
	}
}

// Messages defines custom validation messages
func (r *BusinessFormalisationCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"sme_id.required":                   "SME ID is required",
		"sme_id.numeric":                    "SME ID must be a number",
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

// Attributes defines custom attribute names
func (r *BusinessFormalisationCreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		// Add custom attribute name mappings here
		// e.g., "fieldName": "Field Display Name",
	}
}

// Authorize determines if the user is authorized to make this request
func (r *BusinessFormalisationCreateRequest) Authorize(ctx http.Context) error {
	// TODO: Implement authorization logic
	// Example: Check if user has permission to create business_formalisation
	// return facades.Gate().Allows("create.business_formalisations", ctx)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *BusinessFormalisationCreateRequest) PrepareForValidation(ctx http.Context) error {
	// TODO: Add data preparation logic
	// Example: Normalize, trim, or set default values
	return nil
}

// PassedValidation is called after validation passes
func (r *BusinessFormalisationCreateRequest) PassedValidation(ctx http.Context) error {
	// TODO: Add post-validation logic if needed
	return nil
}

// ToCreateData converts the request to create data map
func (r *BusinessFormalisationCreateRequest) ToCreateData() map[string]interface{} {
	data := map[string]interface{}{
		"sme_id":                   r.SmeID,
		"has_bank_account":         r.HasBankAccount,
		"has_tax_clarification":    r.HasTaxClarification,
		"is_registered_for_vat":    r.IsRegisteredForVat,
		"is_member_of_association": r.IsMemberOfAssociation,
		"is_affiliated":            r.IsAffiliated,
		"has_export_license":       r.HasExportLicense,
		"has_accessed_bds":         r.HasAccessedBds,
		"annual_turnover":          r.AnnualTurnover,
		"estimated_value_of_assets": r.EstimatedValueOfAssets,
		// formalisation_score is intentionally excluded - it's read-only
	}

	return data
}
