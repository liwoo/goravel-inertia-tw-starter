package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

// SetupTOTPRequest handles TOTP setup initiation validation
type SetupTOTPRequest struct {
	Password string `form:"password" json:"password"`
}

// Authorize determines if the user is authorized to make this request
func (r *SetupTOTPRequest) Authorize(ctx http.Context) error {
	return nil
}

// Rules defines validation rules for TOTP setup
func (r *SetupTOTPRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"password": "required",
	}
}

// Messages defines custom validation messages
func (r *SetupTOTPRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"password.required": "Password is required to set up 2FA",
	}
}

// Attributes defines custom attribute names
func (r *SetupTOTPRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"password": "password",
	}
}

// PrepareForValidation allows modification of input before validation
func (r *SetupTOTPRequest) PrepareForValidation(data validation.Data) error {
	return nil
}

// VerifyTOTPRequest handles TOTP code verification
type VerifyTOTPRequest struct {
	Code string `form:"code" json:"code"`
}

// Authorize determines if the user is authorized to make this request
func (r *VerifyTOTPRequest) Authorize(ctx http.Context) error {
	return nil
}

// Rules defines validation rules for TOTP verification
func (r *VerifyTOTPRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"code": "required|min_len:6|max_len:8",
	}
}

// Messages defines custom validation messages
func (r *VerifyTOTPRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"code.required": "Verification code is required",
		"code.min_len":  "Verification code must be at least 6 characters",
		"code.max_len":  "Verification code cannot exceed 8 characters",
	}
}

// Attributes defines custom attribute names
func (r *VerifyTOTPRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"code": "verification code",
	}
}

// PrepareForValidation allows modification of input before validation
func (r *VerifyTOTPRequest) PrepareForValidation(data validation.Data) error {
	return nil
}

// DisableTOTPRequest handles TOTP disabling validation
type DisableTOTPRequest struct {
	Password string `form:"password" json:"password"`
	Code     string `form:"code" json:"code"`
}

// Authorize determines if the user is authorized to make this request
func (r *DisableTOTPRequest) Authorize(ctx http.Context) error {
	return nil
}

// Rules defines validation rules for TOTP disabling
func (r *DisableTOTPRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"password": "required",
		"code":     "required|min_len:6|max_len:8",
	}
}

// Messages defines custom validation messages
func (r *DisableTOTPRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"password.required": "Password is required to disable 2FA",
		"code.required":     "Current 2FA code is required",
		"code.min_len":      "2FA code must be at least 6 characters",
		"code.max_len":      "2FA code cannot exceed 8 characters",
	}
}

// Attributes defines custom attribute names
func (r *DisableTOTPRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"password": "password",
		"code":     "2FA code",
	}
}

// PrepareForValidation allows modification of input before validation
func (r *DisableTOTPRequest) PrepareForValidation(data validation.Data) error {
	return nil
}

// RegenerateBackupCodesRequest handles backup code regeneration validation
type RegenerateBackupCodesRequest struct {
	Password string `form:"password" json:"password"`
}

// Authorize determines if the user is authorized to make this request
func (r *RegenerateBackupCodesRequest) Authorize(ctx http.Context) error {
	return nil
}

// Rules defines validation rules for backup code regeneration
func (r *RegenerateBackupCodesRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"password": "required",
	}
}

// Messages defines custom validation messages
func (r *RegenerateBackupCodesRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"password.required": "Password is required to regenerate backup codes",
	}
}

// Attributes defines custom attribute names
func (r *RegenerateBackupCodesRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"password": "password",
	}
}

// PrepareForValidation allows modification of input before validation
func (r *RegenerateBackupCodesRequest) PrepareForValidation(data validation.Data) error {
	return nil
}

// Verify2FALoginRequest handles 2FA verification during login
type Verify2FALoginRequest struct {
	TempToken string `form:"temp_token" json:"temp_token"`
	Code      string `form:"code" json:"code"`
}

// Authorize determines if the user is authorized to make this request
func (r *Verify2FALoginRequest) Authorize(ctx http.Context) error {
	return nil
}

// Rules defines validation rules for 2FA login verification
func (r *Verify2FALoginRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"temp_token": "required",
		"code":       "required|min_len:6|max_len:8",
	}
}

// Messages defines custom validation messages
func (r *Verify2FALoginRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"temp_token.required": "Temporary token is required",
		"code.required":       "2FA code is required",
		"code.min_len":        "2FA code must be at least 6 characters",
		"code.max_len":        "2FA code cannot exceed 8 characters",
	}
}

// Attributes defines custom attribute names
func (r *Verify2FALoginRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"temp_token": "temporary token",
		"code":       "2FA code",
	}
}

// PrepareForValidation allows modification of input before validation
func (r *Verify2FALoginRequest) PrepareForValidation(data validation.Data) error {
	return nil
}
