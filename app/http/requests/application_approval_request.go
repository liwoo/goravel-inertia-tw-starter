package requests

import (
	"github.com/goravel/framework/contracts/http"
)

// ApplicationApprovalRequest handles validation for application approval
type ApplicationApprovalRequest struct {
}

// Rules returns the validation rules for the request
func (r *ApplicationApprovalRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{}
}

// Messages returns custom validation messages
func (r *ApplicationApprovalRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

// Authorize determines if the user is authorized to make this request
func (r *ApplicationApprovalRequest) Authorize(ctx http.Context) error {
	return nil
}

// ApplicationRejectRequest handles validation for application rejection
type ApplicationRejectRequest struct {
	Reason string `form:"reason" json:"reason"`
}

// Rules returns the validation rules for the request
func (r *ApplicationRejectRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"reason": "string|max:1000",
	}
}

// Messages returns custom validation messages
func (r *ApplicationRejectRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

// Authorize determines if the user is authorized to make this request
func (r *ApplicationRejectRequest) Authorize(ctx http.Context) error {
	return nil
}
