package contracts

import (
	"fmt"
	"github.com/goravel/framework/contracts/http"
)

// ControllerBuilder uses a step-by-step builder pattern to ensure all required methods are implemented
type ControllerBuilder[TModel any, TCreate any, TUpdate any] struct {
	controller *GenericCrudController[TModel, TCreate, TUpdate]
}

// Step 1: Required - Set service
type ControllerBuilderWithService[TModel any, TCreate any, TUpdate any] struct {
	builder *ControllerBuilder[TModel, TCreate, TUpdate]
}

// Step 2: Required - Set auth check
type ControllerBuilderWithAuth[TModel any, TCreate any, TUpdate any] struct {
	builder *ControllerBuilder[TModel, TCreate, TUpdate]
}

// Step 3: Required - Set request bindings
type ControllerBuilderWithBindings[TModel any, TCreate any, TUpdate any] struct {
	builder *ControllerBuilder[TModel, TCreate, TUpdate]
}

// Step 4: Optional configurations and build
type ControllerBuilderComplete[TModel any, TCreate any, TUpdate any] struct {
	builder *ControllerBuilder[TModel, TCreate, TUpdate]
}

// NewControllerBuilder creates a new controller builder that enforces required setup
func NewControllerBuilder[TModel any, TCreate any, TUpdate any](resourceName string) *ControllerBuilder[TModel, TCreate, TUpdate] {
	return &ControllerBuilder[TModel, TCreate, TUpdate]{
		controller: &GenericCrudController[TModel, TCreate, TUpdate]{
			BaseCrudController: NewBaseCrudController(resourceName),
			resourceName:       resourceName,
		},
	}
}

// Step 1: Set service (REQUIRED)
func (b *ControllerBuilder[TModel, TCreate, TUpdate]) WithService(service CrudServiceContract) *ControllerBuilderWithService[TModel, TCreate, TUpdate] {
	if service == nil {
		panic("Service cannot be nil")
	}
	b.controller.service = service
	return &ControllerBuilderWithService[TModel, TCreate, TUpdate]{builder: b}
}

// Step 2: Set auth check (REQUIRED)
func (b *ControllerBuilderWithService[TModel, TCreate, TUpdate]) WithAuthCheck(
	authCheck func(ctx http.Context, action string, resource interface{}) error,
) *ControllerBuilderWithAuth[TModel, TCreate, TUpdate] {
	if authCheck == nil {
		panic("Auth check function cannot be nil")
	}
	b.builder.controller.CheckAuth = authCheck
	return &ControllerBuilderWithAuth[TModel, TCreate, TUpdate]{builder: b.builder}
}

// Step 3: Set request bindings (REQUIRED)
func (b *ControllerBuilderWithAuth[TModel, TCreate, TUpdate]) WithRequestBindings(
	bindCreate func(ctx http.Context) (TCreate, error),
	transformCreate func(TCreate) map[string]interface{},
	bindUpdate func(ctx http.Context, id uint) (TUpdate, error),
	transformUpdate func(TUpdate) map[string]interface{},
) *ControllerBuilderComplete[TModel, TCreate, TUpdate] {
	if bindCreate == nil || transformCreate == nil || bindUpdate == nil || transformUpdate == nil {
		panic("All request binding functions must be provided")
	}
	b.builder.controller.bindCreate = bindCreate
	b.builder.controller.transformCreate = transformCreate
	b.builder.controller.bindUpdate = bindUpdate
	b.builder.controller.transformUpdate = transformUpdate
	return &ControllerBuilderComplete[TModel, TCreate, TUpdate]{builder: b.builder}
}

// Optional configurations
func (b *ControllerBuilderComplete[TModel, TCreate, TUpdate]) WithBeforeStore(
	hook func(ctx http.Context, data map[string]interface{}) error,
) *ControllerBuilderComplete[TModel, TCreate, TUpdate] {
	b.builder.controller.beforeStore = hook
	return b
}

func (b *ControllerBuilderComplete[TModel, TCreate, TUpdate]) WithAfterStore(
	hook func(ctx http.Context, result interface{}) http.Response,
) *ControllerBuilderComplete[TModel, TCreate, TUpdate] {
	b.builder.controller.afterStore = hook
	return b
}

func (b *ControllerBuilderComplete[TModel, TCreate, TUpdate]) WithBeforeUpdate(
	hook func(ctx http.Context, id uint, data map[string]interface{}) error,
) *ControllerBuilderComplete[TModel, TCreate, TUpdate] {
	b.builder.controller.beforeUpdate = hook
	return b
}

func (b *ControllerBuilderComplete[TModel, TCreate, TUpdate]) WithAfterUpdate(
	hook func(ctx http.Context, model interface{}) http.Response,
) *ControllerBuilderComplete[TModel, TCreate, TUpdate] {
	b.builder.controller.afterUpdate = hook
	return b
}

func (b *ControllerBuilderComplete[TModel, TCreate, TUpdate]) WithBeforeDestroy(
	hook func(ctx http.Context, id uint) error,
) *ControllerBuilderComplete[TModel, TCreate, TUpdate] {
	b.builder.controller.beforeDelete = hook
	return b
}

func (b *ControllerBuilderComplete[TModel, TCreate, TUpdate]) WithAfterDestroy(
	hook func(ctx http.Context, id uint) error,
) *ControllerBuilderComplete[TModel, TCreate, TUpdate] {
	// Note: GenericCrudController doesn't have afterDelete hook
	// This could be added if needed
	_ = hook
	return b
}

// Build ensures the controller implements the required interface
func (b *ControllerBuilderComplete[TModel, TCreate, TUpdate]) Build() FullCrudControllerContract {
	// Create a wrapper that ensures all interface methods are implemented
	return &controllerContractWrapper[TModel, TCreate, TUpdate]{
		GenericCrudController: b.builder.controller,
	}
}

// controllerContractWrapper ensures the controller implements FullCrudControllerContract
type controllerContractWrapper[TModel any, TCreate any, TUpdate any] struct {
	*GenericCrudController[TModel, TCreate, TUpdate]
}

// Override Destroy to map to Delete
func (w *controllerContractWrapper[TModel, TCreate, TUpdate]) Destroy(ctx http.Context) http.Response {
	return w.Delete(ctx)
}

func (w *controllerContractWrapper[TModel, TCreate, TUpdate]) GetMeta(ctx http.Context) http.Response {
	// Return metadata about the resource
	meta := map[string]interface{}{
		"searchableFields": w.GetSearchableFields(),
		"validationRules":  w.GetValidationRules(),
	}
	return w.SuccessResponse(ctx, meta, "Metadata retrieved")
}

// Implement AuthorizationControllerContract
func (w *controllerContractWrapper[TModel, TCreate, TUpdate]) CheckPermission(ctx http.Context, permission string, resource interface{}) error {
	if w.CheckAuth != nil {
		return w.CheckAuth(ctx, permission, resource)
	}
	return fmt.Errorf("auth check not configured")
}

func (w *controllerContractWrapper[TModel, TCreate, TUpdate]) GetCurrentUser(ctx http.Context) interface{} {
	// This should be implemented by the specific controller
	return nil
}

func (w *controllerContractWrapper[TModel, TCreate, TUpdate]) RequireAuthentication(ctx http.Context) error {
	user := w.GetCurrentUser(ctx)
	if user == nil {
		return fmt.Errorf("authentication required")
	}
	return nil
}

func (w *controllerContractWrapper[TModel, TCreate, TUpdate]) BuildPermissionsMap(ctx http.Context, resourceType string) map[string]bool {
	// This should be implemented by the specific controller
	return map[string]bool{}
}

// Implement SearchableControllerContract
func (w *controllerContractWrapper[TModel, TCreate, TUpdate]) GetSearchableFields() []string {
	if w.service != nil {
		return w.service.GetSearchableFields()
	}
	return []string{}
}

func (w *controllerContractWrapper[TModel, TCreate, TUpdate]) GetValidationRules() map[string]interface{} {
	if w.service != nil {
		return w.service.GetValidationRules()
	}
	return map[string]interface{}{}
}

// Example usage that enforces compile-time checks:
//
// controller := NewControllerBuilder[models.Book, requests.BookCreateRequest, requests.BookUpdateRequest]("book").
//     WithService(bookService).                           // REQUIRED
//     WithAuthCheck(authCheckFunc).                       // REQUIRED
//     WithRequestBindings(                                // REQUIRED
//         bindCreateFunc,
//         transformCreateFunc,
//         bindUpdateFunc,
//         transformUpdateFunc,
//     ).
//     WithBeforeStore(auditFunc).                         // Optional
//     WithAfterStore(notifyFunc).                         // Optional
//     Build()
//
// If you skip any required step, you get a compile error!