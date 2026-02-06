package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

// FormalisationAmendmentRequest handles validation for formalisation amendment submissions
type FormalisationAmendmentRequest struct {
	SmeID uint `form:"sme_id" json:"sme_id"`

	// SME-level formalisation fields
	RegistrationNumber      *string `form:"registration_number" json:"registration_number"`
	TaxIdentificationNumber *string `form:"tax_identification_number" json:"tax_identification_number"`

	// BusinessFormalisation fields
	HasBankAccount         *bool    `form:"has_bank_account" json:"has_bank_account"`
	HasTaxClarification    *bool    `form:"has_tax_clarification" json:"has_tax_clarification"`
	IsRegisteredForVat     *bool    `form:"is_registered_for_vat" json:"is_registered_for_vat"`
	IsMemberOfAssociation  *bool    `form:"is_member_of_association" json:"is_member_of_association"`
	IsAffiliated           *bool    `form:"is_affiliated" json:"is_affiliated"`
	HasExportLicense       *bool    `form:"has_export_license" json:"has_export_license"`
	HasAccessedBds         *bool    `form:"has_accessed_bds" json:"has_accessed_bds"`
	AnnualTurnover         *float64 `form:"annual_turnover" json:"annual_turnover"`
	EstimatedValueOfAssets *float64 `form:"estimated_value_of_assets" json:"estimated_value_of_assets"`

	// Metadata
	ChangeReason string `form:"change_reason" json:"change_reason"`
}

func (r *FormalisationAmendmentRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *FormalisationAmendmentRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"sme_id": "required|numeric|min:1",
		// All other fields are optional - no validation rules needed for optional fields
		// The service layer handles processing of these optional values
	}
}

func (r *FormalisationAmendmentRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"sme_id.required": "MSME ID is required",
		"sme_id.integer":  "MSME ID must be a valid number",
	}
}

func (r *FormalisationAmendmentRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *FormalisationAmendmentRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}

// ApplicationRejectRequest handles validation for application rejection
type ApplicationRejectRequest struct {
	Reason string `form:"reason" json:"reason"`
}

func (r *ApplicationRejectRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *ApplicationRejectRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"reason": "string|max:1000",
	}
}

func (r *ApplicationRejectRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ApplicationRejectRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ApplicationRejectRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
