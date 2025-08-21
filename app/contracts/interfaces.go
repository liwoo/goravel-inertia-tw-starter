package contracts

import (
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"
)

// ============================================================================
// CRUD Service Interfaces
// ============================================================================

// CrudServiceContract defines the required methods for any CRUD service
type CrudServiceContract interface {
	// Core CRUD operations
	GetList(req ListRequest) (*PaginatedResult, error)
	GetByID(id uint) (interface{}, error)
	Create(data map[string]interface{}) (interface{}, error)
	Update(id uint, data map[string]interface{}) (interface{}, error)
	Delete(id uint) error

	// Search and filtering
	Search(query string, req ListRequest) (*PaginatedResult, error)
	GetListAdvanced(req ListRequest, filters map[string]interface{}) (*PaginatedResult, error)

	// Metadata
	GetSearchableFields() []string
	GetSortableFields() []string
	GetFilterableFields() []string
	GetValidationRules() map[string]interface{}

	// Column mapping for frontend/backend field names
	GetColumnMapping() map[string]string
	MapSortField(frontendField string) (string, bool)
}

// GenericCrudServiceContract extends CrudServiceContract with type safety
type GenericCrudServiceContract[T any] interface {
	CrudServiceContract

	// Type-safe methods
	GetByIDTyped(id uint) (*T, error)
	CreateTyped(data map[string]interface{}) (*T, error)
	UpdateTyped(id uint, data map[string]interface{}) (*T, error)
}

// ============================================================================
// Controller Contracts
// ============================================================================

// SearchableControllerContract defines search capabilities
type SearchableControllerContract interface {
	GetSearchableFields() []string
	GetValidationRules() map[string]interface{}
}

// FullCrudControllerContract combines all controller contracts
type FullCrudControllerContract interface {
	CrudControllerContract
	AuthorizationControllerContract
	SearchableControllerContract
}

// ============================================================================
// Request Validation Interfaces
// ============================================================================

// RequestContract defines validation for incoming requests
type RequestContract interface {
	// Validation rules
	Rules(ctx http.Context) map[string]string
	Messages(ctx http.Context) map[string]string
	Attributes(ctx http.Context) map[string]string

	// Authorization
	Authorize(ctx http.Context) error

	// Lifecycle hooks
	PrepareForValidation(ctx http.Context) error
	PassedValidation(ctx http.Context) error
}

// CreateRequestContract for create operations
type CreateRequestContract interface {
	RequestContract
	ToCreateData() map[string]interface{}
}

// UpdateRequestContract for update operations
type UpdateRequestContract interface {
	RequestContract
	ToUpdateData() map[string]interface{}
	GetResourceID() interface{}
}

// ============================================================================
// Helper/Utility Interfaces
// ============================================================================

// AuditHelper defines audit trail methods
type AuditHelper interface {
	SetCreateAuditFields(ctx http.Context, data map[string]interface{})
	SetUpdateAuditFields(ctx http.Context, data map[string]interface{})
	SetDeleteAuditFields(ctx http.Context, data map[string]interface{})
}

// ============================================================================
// Service Builder Interfaces (for chaining)
// ============================================================================

// ServiceBuilderContract ensures all required methods are called during service setup
type ServiceBuilderContract interface {
	// Required configuration methods that MUST be called
	SetSearchFields(fields ...string) ServiceBuilderContract
	SetSortFields(fields ...string) ServiceBuilderContract
	SetFilterFields(fields ...string) ServiceBuilderContract
	SetValidationRules(rules map[string]interface{}) ServiceBuilderContract

	// Optional configuration methods
	SetRelations(relations ...string) ServiceBuilderContract
	SetDefaultSort(field string, direction string) ServiceBuilderContract
	SetDefaultPageSize(size int) ServiceBuilderContract
	EnableSoftDeletes() ServiceBuilderContract
	EnableScopeFiltering(serviceRegistry string, userField string) ServiceBuilderContract

	// Hook methods
	SetBeforeCreate(hook func(data map[string]interface{}) error) ServiceBuilderContract
	SetAfterCreate(hook func(model interface{}) error) ServiceBuilderContract
	SetBeforeUpdate(hook func(id uint, data map[string]interface{}) error) ServiceBuilderContract
	SetAfterUpdate(hook func(model interface{}) error) ServiceBuilderContract
	SetBeforeDelete(hook func(id uint) error) ServiceBuilderContract
	SetAfterDelete(hook func(id uint) error) ServiceBuilderContract

	// Custom behavior
	SetCustomSearch(fn func(query orm.Query, search string) orm.Query) ServiceBuilderContract
	SetCustomFilters(fn func(query orm.Query, filters map[string]interface{}) orm.Query) ServiceBuilderContract
	SetCustomQuery(fn func(query orm.Query) orm.Query) ServiceBuilderContract

	// Build method to finalize and return the service
	Build() (CrudServiceContract, error)
}

// ============================================================================
// Factory Functions to ensure contracts are met
// ============================================================================

// MustImplementCrudService panics if the service doesn't implement CrudServiceContract
func MustImplementCrudService(service interface{}) CrudServiceContract {
	if s, ok := service.(CrudServiceContract); ok {
		return s
	}
	panic("Service does not implement CrudServiceContract")
}

// MustImplementCrudController panics if the controller doesn't implement FullCrudControllerContract
func MustImplementCrudController(controller interface{}) FullCrudControllerContract {
	if c, ok := controller.(FullCrudControllerContract); ok {
		return c
	}
	panic("Controller does not implement FullCrudControllerContract")
}

// MustImplementPageController panics if the controller doesn't implement PageControllerContract
func MustImplementPageController(controller interface{}) PageControllerContract {
	if c, ok := controller.(PageControllerContract); ok {
		return c
	}
	panic("Page controller does not implement PageControllerContract")
}

// MustImplementRequest panics if the request doesn't implement appropriate contract
func MustImplementCreateRequest(request interface{}) CreateRequestContract {
	if r, ok := request.(CreateRequestContract); ok {
		return r
	}
	panic("Request does not implement CreateRequestContract")
}

func MustImplementUpdateRequest(request interface{}) UpdateRequestContract {
	if r, ok := request.(UpdateRequestContract); ok {
		return r
	}
	panic("Request does not implement UpdateRequestContract")
}
