package contracts

import (
	"github.com/goravel/framework/contracts/http"
)

// StaticEnforcedController uses the type system to enforce all requirements at compile time
// The zero value of this struct is NOT usable - you MUST use the builder
type StaticEnforcedController[T any, C CreateRequestContract, U UpdateRequestContract] struct {
	*EnforcedCrudController[T, C, U]
	
	// These fields ensure the controller was properly constructed
	_ requiredCreateRequest[C]
	_ requiredUpdateRequest[U]
	_ requiredAuthChecker
}

// GetFilters returns filter metadata for the resource
func (c *StaticEnforcedController[T, C, U]) GetFilters(ctx http.Context) http.Response {
	// Return empty metadata by default
	// Concrete controllers should override this method to provide actual filter definitions
	metadata := map[string]interface{}{
		"filters": []interface{}{},
	}
	
	return ctx.Response().Json(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    metadata,
	})
}

// Marker types that can only be created by the builder
type requiredCreateRequest[C CreateRequestContract] struct{}
type requiredUpdateRequest[U UpdateRequestContract] struct{}
type requiredAuthChecker struct{}

// StaticControllerBuilder uses a type-state pattern to enforce compile-time checks
type StaticControllerBuilder[T any, C CreateRequestContract, U UpdateRequestContract, State any] struct {
	resourceName string
	service      CrudServiceContract
	authChecker  func(ctx http.Context, action string, resource interface{}) error
}

// Initial state markers
type (
	needsCreateRequest struct{}
	needsUpdateRequest struct{}
	needsAuthChecker   struct{}
	readyToBuild       struct{}
)

// NewStaticControllerBuilder starts the builder in a state that requires all methods to be called
func NewStaticControllerBuilder[T any, C CreateRequestContract, U UpdateRequestContract](
	resourceName string,
	service CrudServiceContract,
) *StaticControllerBuilder[T, C, U, needsCreateRequest] {
	return &StaticControllerBuilder[T, C, U, needsCreateRequest]{
		resourceName: resourceName,
		service:      service,
	}
}

// ValidateCreateRequest moves from needsCreateRequest to needsUpdateRequest state
// This method is ONLY available when State = needsCreateRequest
func (b *StaticControllerBuilder[T, C, U, needsCreateRequest]) ValidateCreateRequest() *StaticControllerBuilder[T, C, U, needsUpdateRequest] {
	// The fact that C must implement CreateRequestContract is enforced by the type parameter
	return &StaticControllerBuilder[T, C, U, needsUpdateRequest]{
		resourceName: b.resourceName,
		service:      b.service,
		authChecker:  b.authChecker,
	}
}

// ValidateUpdateRequest moves from needsUpdateRequest to needsAuthChecker state
// This method is ONLY available when State = needsUpdateRequest
func (b *StaticControllerBuilder[T, C, U, needsUpdateRequest]) ValidateUpdateRequest() *StaticControllerBuilder[T, C, U, needsAuthChecker] {
	// The fact that U must implement UpdateRequestContract is enforced by the type parameter
	return &StaticControllerBuilder[T, C, U, needsAuthChecker]{
		resourceName: b.resourceName,
		service:      b.service,
		authChecker:  b.authChecker,
	}
}

// WithAuthChecker moves from needsAuthChecker to readyToBuild state
// This method is ONLY available when State = needsAuthChecker
func (b *StaticControllerBuilder[T, C, U, needsAuthChecker]) WithAuthChecker(
	checker func(ctx http.Context, action string, resource interface{}) error,
) *StaticControllerBuilder[T, C, U, readyToBuild] {
	return &StaticControllerBuilder[T, C, U, readyToBuild]{
		resourceName: b.resourceName,
		service:      b.service,
		authChecker:  checker,
	}
}

// Build is ONLY available when State = readyToBuild
// This ensures at compile time that all required methods have been called
func (b *StaticControllerBuilder[T, C, U, readyToBuild]) Build() *StaticEnforcedController[T, C, U] {
	controller := NewEnforcedCrudController[T, C, U](b.resourceName, b.service)
	controller.CheckAuth = b.authChecker
	
	return &StaticEnforcedController[T, C, U]{
		EnforcedCrudController: controller,
	}
}

// Example usage that MUST follow the exact sequence:
//
// controller := NewStaticControllerBuilder[Model, *CreateReq, *UpdateReq]("resource", service).
//     ValidateCreateRequest().     // Must be called first
//     ValidateUpdateRequest().     // Must be called second
//     WithAuthChecker(authFunc).   // Must be called third
//     Build()                      // Can only be called after all above
//
// Any other sequence will NOT compile!