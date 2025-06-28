```go
// Package contracts defines interfaces for CRUD operations
package contracts

import (
"github.com/goravel/framework/contracts/http"
)

// Common Request/Response Types
type ListRequest struct {
Page     int    `form:"page" json:"page"`
PageSize int    `form:"pageSize" json:"pageSize"`
Sort     string `form:"sort" json:"sort"`
Search   string `form:"search" json:"search"`
}

type ListResponse struct {
Data       interface{} `json:"data"`
Total      int64       `json:"total"`
Page       int         `json:"page"`
PageSize   int         `json:"pageSize"`
TotalPages int         `json:"totalPages"`
HasNext    bool        `json:"hasNext"`
HasPrev    bool        `json:"hasPrev"`
}

// SetDefaults applies sensible defaults to ListRequest
func (r *ListRequest) SetDefaults() {
if r.Page <= 0 {
r.Page = 1
}
if r.PageSize <= 0 {
r.PageSize = 20
}
if r.PageSize > 100 {
r.PageSize = 100
}
if r.Sort == "" {
r.Sort = "id DESC"
}
}

// Service Interfaces

// CrudService defines the contract for CRUD service operations
type CrudService interface {
GetList(req ListRequest) (*ListResponse, error)
GetByID(id interface{}) (interface{}, error)
Create(data map[string]interface{}) (interface{}, error)
Update(id interface{}, data map[string]interface{}) (interface{}, error)
Delete(id interface{}) error
}

// FilterableService extends CrudService with advanced filtering
type FilterableService interface {
CrudService
GetListWithFilters(req ListRequest, filters map[string]interface{}) (*ListResponse, error)
}

// SoftDeleteService adds soft delete capabilities
type SoftDeleteService interface {
CrudService
SoftDelete(id interface{}) error
Restore(id interface{}) error
GetTrashed(req ListRequest) (*ListResponse, error)
ForceDelete(id interface{}) error
}

// BulkOperationService adds bulk operation capabilities
type BulkOperationService interface {
CrudService
BulkCreate(data []map[string]interface{}) ([]interface{}, error)
BulkUpdate(ids []interface{}, data map[string]interface{}) error
BulkDelete(ids []interface{}) error
}

// FullFeaturedService combines all service capabilities
type FullFeaturedService interface {
FilterableService
SoftDeleteService
BulkOperationService
}

// Controller Interfaces

// CrudController defines the contract for CRUD HTTP controllers
type CrudController interface {
Index(ctx http.Context) http.Response   // GET /resource
Show(ctx http.Context) http.Response    // GET /resource/{id}
Store(ctx http.Context) http.Response   // POST /resource
Update(ctx http.Context) http.Response  // PUT /resource/{id}
Destroy(ctx http.Context) http.Response // DELETE /resource/{id}
}

// RestfulController extends CrudController with additional REST methods
type RestfulController interface {
CrudController
Create(ctx http.Context) http.Response // GET /resource/create (form)
Edit(ctx http.Context) http.Response   // GET /resource/{id}/edit (form)
}

// BulkController adds bulk operation endpoints
type BulkController interface {
CrudController
BulkStore(ctx http.Context) http.Response   // POST /resource/bulk
BulkUpdate(ctx http.Context) http.Response  // PUT /resource/bulk
BulkDestroy(ctx http.Context) http.Response // DELETE /resource/bulk
}

// SoftDeleteController adds soft delete endpoints
type SoftDeleteController interface {
CrudController
Trash(ctx http.Context) http.Response    // GET /resource/trash
Restore(ctx http.Context) http.Response  // POST /resource/{id}/restore
ForceDestroy(ctx http.Context) http.Response // DELETE /resource/{id}/force
}

// AdvancedController combines all controller capabilities
type AdvancedController interface {
RestfulController
BulkController
SoftDeleteController
Export(ctx http.Context) http.Response // GET /resource/export
Import(ctx http.Context) http.Response // POST /resource/import
}

// Repository Interfaces (for data access layer)

// Repository defines basic data access operations
type Repository interface {
Find(id interface{}) (interface{}, error)
FindAll() ([]interface{}, error)
FindBy(field string, value interface{}) ([]interface{}, error)
Create(data map[string]interface{}) (interface{}, error)
Update(id interface{}, data map[string]interface{}) error
Delete(id interface{}) error
Count() (int64, error)
}

// QueryableRepository adds advanced querying capabilities
type QueryableRepository interface {
Repository
Where(field string, value interface{}) QueryableRepository
WhereIn(field string, values []interface{}) QueryableRepository
WhereLike(field string, pattern string) QueryableRepository
OrderBy(field string, direction string) QueryableRepository
Limit(limit int) QueryableRepository
Offset(offset int) QueryableRepository
With(relations ...string) QueryableRepository
Get() ([]interface{}, error)
First() (interface{}, error)
Paginate(page, pageSize int) (*ListResponse, error)
}

// Validation Interfaces

// Validator defines validation contract
type Validator interface {
Validate(data map[string]interface{}) error
ValidateForCreate(data map[string]interface{}) error
ValidateForUpdate(data map[string]interface{}) error
}

// Permission Interfaces

// Authorizer defines authorization contract
type Authorizer interface {
CanView(userID interface{}, resourceID interface{}) bool
CanCreate(userID interface{}) bool
CanUpdate(userID interface{}, resourceID interface{}) bool
CanDelete(userID interface{}, resourceID interface{}) bool
}

// Entity Interfaces

// Entity defines basic entity contract
type Entity interface {
GetID() interface{}
GetTableName() string
}

// Searchable defines searchable entity contract
type Searchable interface {
Entity
GetSearchableFields() []string
}

// Filterable defines filterable entity contract
type Filterable interface {
Entity
GetFilterableFields() map[string]string // field -> database column mapping
}

// Sortable defines sortable entity contract
type Sortable interface {
Entity
GetSortableFields() []string
GetDefaultSort() string
}

// Timestamped defines entity with timestamps
type Timestamped interface {
Entity
GetCreatedAt() interface{}
GetUpdatedAt() interface{}
}

// SoftDeletable defines soft deletable entity
type SoftDeletable interface {
Timestamped
GetDeletedAt() interface{}
IsDeleted() bool
}

// Event Interfaces (for hooks/events)

// EventHandler defines event handling contract
type EventHandler interface {
BeforeCreate(data map[string]interface{}) error
AfterCreate(entity interface{}) error
BeforeUpdate(id interface{}, data map[string]interface{}) error
AfterUpdate(entity interface{}) error
BeforeDelete(id interface{}) error
AfterDelete(id interface{}) error
}

// Response Helper Interfaces

// ResponseBuilder helps build consistent API responses
type ResponseBuilder interface {
Success(data interface{}) http.Response
Error(message string, code int) http.Response
Created(data interface{}) http.Response
Updated(data interface{}) http.Response
Deleted() http.Response
NotFound(message string) http.Response
ValidationError(errors map[string]string) http.Response
Paginated(response *ListResponse) http.Response
}

// Cache Interfaces

// Cacheable defines cacheable service contract
type Cacheable interface {
GetCacheKey(id interface{}) string
GetCacheTTL() int
ShouldCache() bool
InvalidateCache(id interface{}) error
}

// Configuration Interfaces

// CrudConfig defines CRUD service configuration
type CrudConfig struct {
TableName        string
PrimaryKey       string
DefaultPageSize  int
MaxPageSize      int
DefaultSort      string
SearchFields     []string
FilterableFields map[string]string
SoftDelete       bool
Timestamps       bool
EnableCache      bool
CacheTTL         int
}

// ConfigurableService allows runtime configuration
type ConfigurableService interface {
GetConfig() *CrudConfig
SetConfig(config *CrudConfig)
Configure(fn func(*CrudConfig))
}

// Factory Interfaces

// ServiceFactory creates services with proper dependencies
type ServiceFactory interface {
CreateCrudService(config *CrudConfig) CrudService
CreateFilterableService(config *CrudConfig) FilterableService
CreateFullFeaturedService(config *CrudConfig) FullFeaturedService
}

// ControllerFactory creates controllers with proper dependencies
type ControllerFactory interface {
CreateCrudController(service CrudService) CrudController
CreateAdvancedController(service FullFeaturedService) AdvancedController
}

// Middleware Interfaces

// CrudMiddleware defines CRUD-specific middleware
type CrudMiddleware interface {
ValidateID(ctx http.Context) http.Response
ValidateInput(ctx http.Context) http.Response
CheckPermissions(ctx http.Context) http.Response
RateLimit(ctx http.Context) http.Response
Cache(ctx http.Context) http.Response
}

// Usage Example Interfaces

// BookService shows how to implement the contracts
type BookService interface {
FilterableService
SoftDeleteService
// Add domain-specific methods
GetByISBN(isbn string) (interface{}, error)
GetByAuthor(author string) (*ListResponse, error)
GetPopular(limit int) ([]interface{}, error)
}

// BookController shows how to implement controller contracts
type BookController interface {
AdvancedController
// Add domain-specific endpoints
GetByISBN(ctx http.Context) http.Response     // GET /books/isbn/{isbn}
GetByAuthor(ctx http.Context) http.Response   // GET /books/author/{author}
GetPopular(ctx http.Context) http.Response    // GET /books/popular
}

// Implementation hint: Services and controllers implement these interfaces
// This ensures consistent behavior across all CRUD operations while
// allowing for domain-specific extensions
```