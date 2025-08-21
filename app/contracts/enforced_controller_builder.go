package contracts

import (
	"github.com/goravel/framework/contracts/http"
)

// EnforcedControllerBuilder ensures all required methods are implemented at compile time
type EnforcedControllerBuilder[T any, C CreateRequestContract, U UpdateRequestContract] struct {
	resourceName string
	service      CrudServiceContract

	// These fields MUST be set before Build() can be called
	createRequestFactory func() C
	updateRequestFactory func() U
	authChecker          func(ctx http.Context, action string, resource interface{}) error

	// Flags to track what has been set
	hasCreateFactory bool
	hasUpdateFactory bool
	hasAuthChecker   bool
}

// NewEnforcedControllerBuilder creates a new builder that enforces compile-time checks
func NewEnforcedControllerBuilder[T any, C CreateRequestContract, U UpdateRequestContract](
	resourceName string,
	service CrudServiceContract,
) *EnforcedControllerBuilder[T, C, U] {
	return &EnforcedControllerBuilder[T, C, U]{
		resourceName: resourceName,
		service:      service,
	}
}

// WithCreateRequestFactory sets the factory for creating request instances
// This MUST be called or Build() will panic
func (b *EnforcedControllerBuilder[T, C, U]) WithCreateRequestFactory(factory func() C) *EnforcedControllerBuilder[T, C, U] {
	b.createRequestFactory = factory
	b.hasCreateFactory = true
	return b
}

// WithUpdateRequestFactory sets the factory for creating update request instances
// This MUST be called or Build() will panic
func (b *EnforcedControllerBuilder[T, C, U]) WithUpdateRequestFactory(factory func() U) *EnforcedControllerBuilder[T, C, U] {
	b.updateRequestFactory = factory
	b.hasUpdateFactory = true
	return b
}

// WithAuthChecker sets the authorization checker
// This MUST be called or Build() will panic
func (b *EnforcedControllerBuilder[T, C, U]) WithAuthChecker(checker func(ctx http.Context, action string, resource interface{}) error) *EnforcedControllerBuilder[T, C, U] {
	b.authChecker = checker
	b.hasAuthChecker = true
	return b
}

// Build creates the controller, panicking if any required methods haven't been set
func (b *EnforcedControllerBuilder[T, C, U]) Build() *EnforcedCrudController[T, C, U] {
	// Enforce that all required methods have been set
	if !b.hasCreateFactory {
		panic("CreateRequestFactory must be set via WithCreateRequestFactory()")
	}
	if !b.hasUpdateFactory {
		panic("UpdateRequestFactory must be set via WithUpdateRequestFactory()")
	}
	if !b.hasAuthChecker {
		panic("AuthChecker must be set via WithAuthChecker()")
	}

	// Verify that the factories produce types that implement the interfaces
	// This will cause a compile error if they don't
	testCreate := b.createRequestFactory()
	var _ CreateRequestContract = testCreate

	testUpdate := b.updateRequestFactory()
	var _ UpdateRequestContract = testUpdate

	// Create the controller
	controller := NewEnforcedCrudController[T, C, U](b.resourceName, b.service)
	controller.CheckAuth = b.authChecker

	return controller
}

// MustBuild is like Build but returns the concrete controller
func (b *EnforcedControllerBuilder[T, C, U]) MustBuild() *EnforcedCrudController[T, C, U] {
	return b.Build()
}
