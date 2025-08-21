package contracts

import (
	"fmt"
	"github.com/goravel/framework/contracts/http"
)

// InterfaceEnforcedController is a base controller that implements all required security interfaces
type InterfaceEnforcedController[T any, C CreateRequestContract, U UpdateRequestContract] struct {
	*EnforcedCrudController[T, C, U]

	// These fields ensure the interfaces are satisfied
	authChecker     func(ctx http.Context, action string, resource interface{}) error
	createValidator RequestValidator
	updateValidator RequestValidator
}

// Compile-time check that our controller implements all required interfaces
var _ SecuredCrudController = (*InterfaceEnforcedController[any, CreateRequestContract, UpdateRequestContract])(nil)

// GetAuthChecker implements AuthRequiredController
func (c *InterfaceEnforcedController[T, C, U]) GetAuthChecker() func(ctx http.Context, action string, resource interface{}) error {
	if c.authChecker == nil {
		panic("InterfaceEnforcedController: authChecker not set - use builder.WithAuthChecker()")
	}
	return c.authChecker
}

// GetCreateValidator implements ValidatedCreateController
func (c *InterfaceEnforcedController[T, C, U]) GetCreateValidator() RequestValidator {
	if c.createValidator == nil {
		panic("InterfaceEnforcedController: createValidator not set - use builder.ValidateCreateRequest()")
	}
	return c.createValidator
}

// GetUpdateValidator implements ValidatedUpdateController
func (c *InterfaceEnforcedController[T, C, U]) GetUpdateValidator() RequestValidator {
	if c.updateValidator == nil {
		panic("InterfaceEnforcedController: updateValidator not set - use builder.ValidateUpdateRequest()")
	}
	return c.updateValidator
}

// defaultRequestValidator wraps the request contracts to implement RequestValidator
type defaultRequestValidator[C CreateRequestContract, U UpdateRequestContract] struct{}

func (v *defaultRequestValidator[C, U]) ValidateCreate(ctx http.Context, data interface{}) error {
	createReq, ok := data.(C)
	if !ok {
		return fmt.Errorf("invalid create request type")
	}

	// Validate using the contract methods
	if err := createReq.PrepareForValidation(ctx); err != nil {
		return err
	}

	// Here you would integrate with Goravel's validation
	// For now, we assume the Rules() method defines validation
	rules := createReq.Rules(ctx)
	if len(rules) == 0 {
		return fmt.Errorf("no validation rules defined for create request")
	}

	return nil
}

func (v *defaultRequestValidator[C, U]) ValidateUpdate(ctx http.Context, id string, data interface{}) error {
	updateReq, ok := data.(U)
	if !ok {
		return fmt.Errorf("invalid update request type")
	}

	// Validate using the contract methods
	if err := updateReq.PrepareForValidation(ctx); err != nil {
		return err
	}

	// Here you would integrate with Goravel's validation
	rules := updateReq.Rules(ctx)
	if len(rules) == 0 {
		return fmt.Errorf("no validation rules defined for update request")
	}

	return nil
}

// InterfaceEnforcedBuilder builds controllers with compile-time interface enforcement
type InterfaceEnforcedBuilder[T any, C CreateRequestContract, U UpdateRequestContract] struct {
	controller *InterfaceEnforcedController[T, C, U]

	// Track what's been set
	hasAuth   bool
	hasCreate bool
	hasUpdate bool
}

// NewInterfaceEnforcedBuilder creates a new builder
func NewInterfaceEnforcedBuilder[T any, C CreateRequestContract, U UpdateRequestContract](
	resourceName string,
	service CrudServiceContract,
) *InterfaceEnforcedBuilder[T, C, U] {
	base := NewEnforcedCrudController[T, C, U](resourceName, service)

	return &InterfaceEnforcedBuilder[T, C, U]{
		controller: &InterfaceEnforcedController[T, C, U]{
			EnforcedCrudController: base,
		},
	}
}

// WithAuthChecker sets the auth checker
func (b *InterfaceEnforcedBuilder[T, C, U]) WithAuthChecker(
	checker func(ctx http.Context, action string, resource interface{}) error,
) *InterfaceEnforcedBuilder[T, C, U] {
	if checker == nil {
		panic("auth checker cannot be nil")
	}
	b.controller.authChecker = checker
	b.controller.CheckAuth = checker // Also set on base controller
	b.hasAuth = true
	return b
}

// ValidateCreateRequest enables create validation
func (b *InterfaceEnforcedBuilder[T, C, U]) ValidateCreateRequest() *InterfaceEnforcedBuilder[T, C, U] {
	b.controller.createValidator = &defaultRequestValidator[C, U]{}
	b.hasCreate = true
	return b
}

// ValidateUpdateRequest enables update validation
func (b *InterfaceEnforcedBuilder[T, C, U]) ValidateUpdateRequest() *InterfaceEnforcedBuilder[T, C, U] {
	b.controller.updateValidator = &defaultRequestValidator[C, U]{}
	b.hasUpdate = true
	return b
}

// Build returns the controller only if all requirements are met
func (b *InterfaceEnforcedBuilder[T, C, U]) Build() *InterfaceEnforcedController[T, C, U] {
	if !b.hasAuth {
		panic("InterfaceEnforcedBuilder: WithAuthChecker() must be called before Build()")
	}
	if !b.hasCreate {
		panic("InterfaceEnforcedBuilder: ValidateCreateRequest() must be called before Build()")
	}
	if !b.hasUpdate {
		panic("InterfaceEnforcedBuilder: ValidateUpdateRequest() must be called before Build()")
	}

	// Verify the controller implements all interfaces at build time
	var _ SecuredCrudController = b.controller

	return b.controller
}
