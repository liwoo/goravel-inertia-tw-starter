package helpers

import (
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

// ValidationHelper provides validation utilities and error formatting
type ValidationHelper struct{}

// NewValidationHelper creates a new validation helper
func NewValidationHelper() *ValidationHelper {
	return &ValidationHelper{}
}

// FormatValidationErrors formats Goravel validation errors for frontend consumption
// Converts nested map structures like "map[email:Please provide a valid email address.]"
// to clean string messages: "Please provide a valid email address."
func (v *ValidationHelper) FormatValidationErrors(errors validation.Errors) map[string]string {
	formattedErrors := make(map[string]string)

	// Convert errors.All() to a simple string map
	// errors.All() returns the errors in a nested map format
	allErrors := errors.All()

	for field, errorValue := range allErrors {
		// Handle the nested map structure: map[email:Please provide a valid email address.]
		errorStr := fmt.Sprintf("%v", errorValue)

		// Extract the actual error message from the map string format
		// Pattern: "map[field:error message]"
		if strings.HasPrefix(errorStr, "map[") && strings.HasSuffix(errorStr, "]") {
			// Remove "map[" and "]"
			inner := errorStr[4 : len(errorStr)-1]

			// Find the colon that separates field from message
			colonIndex := strings.Index(inner, ":")
			if colonIndex != -1 && colonIndex < len(inner)-1 {
				// Extract the message part after the colon
				message := inner[colonIndex+1:]
				formattedErrors[field] = message
			} else {
				// Fallback to the original string if parsing fails
				formattedErrors[field] = errorStr
			}
		} else {
			// If it's not in the expected format, use as is
			formattedErrors[field] = errorStr
		}
	}

	return formattedErrors
}

// FlashErrorsAndRedirect flashes formatted errors to session and redirects
func (v *ValidationHelper) FlashErrorsAndRedirect(ctx http.Context, errors map[string]string, redirectTo string) http.Response {
	ctx.Request().Session().Flash("errors", errors)
	return ctx.Response().Redirect(http.StatusFound, redirectTo)
}

// FlashValidationErrorsAndRedirect formats validation errors and redirects
func (v *ValidationHelper) FlashValidationErrorsAndRedirect(ctx http.Context, errors validation.Errors, redirectTo string) http.Response {
	formattedErrors := v.FormatValidationErrors(errors)
	return v.FlashErrorsAndRedirect(ctx, formattedErrors, redirectTo)
}

// CreateFieldError creates a single field error for custom validation
func (v *ValidationHelper) CreateFieldError(field, message string) map[string]string {
	return map[string]string{
		field: message,
	}
}

// CreateGeneralError creates a general error message
func (v *ValidationHelper) CreateGeneralError(message string) map[string]string {
	return map[string]string{
		"general": message,
	}
}

// MergeErrors merges multiple error maps into one
func (v *ValidationHelper) MergeErrors(errorMaps ...map[string]string) map[string]string {
	merged := make(map[string]string)

	for _, errorMap := range errorMaps {
		for field, message := range errorMap {
			merged[field] = message
		}
	}

	return merged
}

// ValidateAndFormatErrors performs validation and returns formatted errors
func (v *ValidationHelper) ValidateAndFormatErrors(ctx http.Context, request http.FormRequest) (map[string]string, error) {
	errors, err := ctx.Request().ValidateRequest(request)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if errors != nil {
		return v.FormatValidationErrors(errors), fmt.Errorf("validation errors occurred")
	}

	return nil, nil
}

// Common validation error messages
const (
	ErrInvalidCredentials = "Invalid credentials"
	ErrEmailNotFound      = "Invalid credentials (email not found)"
	ErrPasswordMismatch   = "Invalid credentials (password mismatch)"
	ErrAccountInactive    = "Account is inactive"
	ErrAccountLocked      = "Account is locked"
	ErrTooManyAttempts    = "Too many login attempts. Please try again later"
)
