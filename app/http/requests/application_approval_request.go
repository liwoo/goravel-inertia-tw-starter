package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// ApplicationApprovalRequest handles validation for application approval
type ApplicationApprovalRequest struct {
	SmeID uint `form:"sme_id" json:"sme_id"`
}

// Rules returns the validation rules for the request
func (r *ApplicationApprovalRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"sme_id": "required|numeric|min:1",
	}
}

// Messages returns custom validation messages
func (r *ApplicationApprovalRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"sme_id.required": "MSME selection is required",
		"sme_id.numeric":  "MSME ID must be a number",
		"sme_id.min":      "Invalid MSME ID",
	}
}

// Authorize determines if the user is authorized to make this request
func (r *ApplicationApprovalRequest) Authorize(ctx http.Context) error {
	return nil
}
