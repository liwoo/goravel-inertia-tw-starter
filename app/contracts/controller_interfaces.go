package contracts

import (
	"github.com/goravel/framework/contracts/http"
)

// AuthRequiredController ensures controllers implement auth checking
type AuthRequiredController interface {
	// GetAuthChecker must return a non-nil auth checker function
	GetAuthChecker() func(ctx http.Context, action string, resource interface{}) error
}

// ValidatedCreateController ensures controllers validate create requests
type ValidatedCreateController interface {
	// GetCreateValidator returns the validation rules for create requests
	GetCreateValidator() RequestValidator
}

// ValidatedUpdateController ensures controllers validate update requests
type ValidatedUpdateController interface {
	// GetUpdateValidator returns the validation rules for update requests
	GetUpdateValidator() RequestValidator
}

// RequestValidator defines the validation contract
type RequestValidator interface {
	ValidateCreate(ctx http.Context, data interface{}) error
	ValidateUpdate(ctx http.Context, id string, data interface{}) error
}

// SecuredCrudController combines all security requirements
type SecuredCrudController interface {
	AuthRequiredController
	ValidatedCreateController
	ValidatedUpdateController
}

// CompileTimeValidator is a helper type that ensures all interfaces are implemented
type CompileTimeValidator struct {
	_ func(SecuredCrudController) // This field ensures the controller implements all required interfaces
}

// ValidateController is a compile-time check that ensures controller implements all required interfaces
// Usage: var _ = ValidateController(&YourController{})
func ValidateController(controller SecuredCrudController) SecuredCrudController {
	return controller
}
