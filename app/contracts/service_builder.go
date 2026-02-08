package contracts

import (
	"books-database/app/auth"

	"github.com/goravel/framework/contracts/database/orm"
)

// ServiceBuilder uses a step-by-step builder pattern to ensure all required fields are set
type ServiceBuilder[T any] struct {
	service *GenericCrudService[T]
}

// Step 1: Required - Set search fields
type ServiceBuilderWithSearch[T any] struct {
	builder *ServiceBuilder[T]
}

// Step 2: Required - Set sort fields
type ServiceBuilderWithSort[T any] struct {
	builder *ServiceBuilder[T]
}

// Step 3: Required - Set filter fields
type ServiceBuilderWithFilter[T any] struct {
	builder *ServiceBuilder[T]
}

// Step 4: Required - Set validation rules
type ServiceBuilderWithValidation[T any] struct {
	builder *ServiceBuilder[T]
}

// Step 5: Optional configurations and build
type ServiceBuilderComplete[T any] struct {
	builder *ServiceBuilder[T]
}

// NewServiceBuilder creates a new service builder that enforces required fields
func NewServiceBuilder[T any](resourceName string, primaryKey string) *ServiceBuilder[T] {
	// Create the service using the existing constructor
	service := NewGenericCrudService[T](resourceName, primaryKey)

	return &ServiceBuilder[T]{
		service: service,
	}
}

// Step 1: Set search fields (REQUIRED)
func (b *ServiceBuilder[T]) WithSearchFields(fields ...string) *ServiceBuilderWithSearch[T] {
	if len(fields) == 0 {
		panic("At least one search field must be specified")
	}
	b.service.searchFields = fields
	return &ServiceBuilderWithSearch[T]{builder: b}
}

// Step 2: Set sort fields (REQUIRED)
func (b *ServiceBuilderWithSearch[T]) WithSortFields(fields ...string) *ServiceBuilderWithSort[T] {
	if len(fields) == 0 {
		panic("At least one sort field must be specified")
	}
	b.builder.service.sortFields = fields
	return &ServiceBuilderWithSort[T]{builder: b.builder}
}

// Step 3: Set filter fields (REQUIRED)
func (b *ServiceBuilderWithSort[T]) WithFilterFields(fields ...string) *ServiceBuilderWithFilter[T] {
	if len(fields) == 0 {
		panic("At least one filter field must be specified")
	}
	b.builder.service.filterFields = fields
	return &ServiceBuilderWithFilter[T]{builder: b.builder}
}

// Step 4: Set validation rules (REQUIRED)
func (b *ServiceBuilderWithFilter[T]) WithValidationRules(rules map[string]interface{}) *ServiceBuilderComplete[T] {
	if len(rules) == 0 {
		panic("At least one validation rule must be specified")
	}
	b.builder.service.validationRules = rules
	return &ServiceBuilderComplete[T]{builder: b.builder}
}

// Optional configurations - can be chained
func (b *ServiceBuilderComplete[T]) WithRelations(relations ...string) *ServiceBuilderComplete[T] {
	b.builder.service.relations = relations
	return b
}

func (b *ServiceBuilderComplete[T]) WithDefaultSort(field string, direction string) *ServiceBuilderComplete[T] {
	// Note: defaultSort and defaultDirection are not fields in GenericCrudService
	// They would need to be handled in custom query logic
	// This is kept for API compatibility but doesn't set any fields
	return b
}

func (b *ServiceBuilderComplete[T]) WithSoftDeletes() *ServiceBuilderComplete[T] {
	// Note: soft deletes are handled at the model level in GORM
	// This is kept for API compatibility
	return b
}

func (b *ServiceBuilderComplete[T]) WithScopeFiltering(serviceRegistry string, userField string) *ServiceBuilderComplete[T] {
	b.builder.service.enableScopeFiltering = true
	// Map service name to auth.ServiceRegistry
	serviceMap := map[string]auth.ServiceRegistry{
		"books":       "books",
		"users":       "users",
		"roles":       "roles",
		"permissions": "permissions",
		"lenders":     "lenders",
		"bdsps":       "bdsps",
	}

	if registry, ok := serviceMap[serviceRegistry]; ok {
		b.builder.service.serviceRegistry = registry
	}
	b.builder.service.scopeUserField = userField
	return b
}

func (b *ServiceBuilderComplete[T]) WithTenantAwareness() *ServiceBuilderComplete[T] {
	b.builder.service.tenantAware = true
	return b
}

func (b *ServiceBuilderComplete[T]) WithBeforeCreate(hook func(data map[string]interface{}) error) *ServiceBuilderComplete[T] {
	b.builder.service.beforeCreate = hook
	return b
}

func (b *ServiceBuilderComplete[T]) WithAfterCreate(hook func(model *T) error) *ServiceBuilderComplete[T] {
	b.builder.service.afterCreate = hook
	return b
}

func (b *ServiceBuilderComplete[T]) WithBeforeUpdate(hook func(id uint, data map[string]interface{}) error) *ServiceBuilderComplete[T] {
	b.builder.service.beforeUpdate = hook
	return b
}

func (b *ServiceBuilderComplete[T]) WithAfterUpdate(hook func(model *T) error) *ServiceBuilderComplete[T] {
	b.builder.service.afterUpdate = hook
	return b
}

func (b *ServiceBuilderComplete[T]) WithBeforeDelete(hook func(id uint) error) *ServiceBuilderComplete[T] {
	b.builder.service.beforeDelete = hook
	return b
}

func (b *ServiceBuilderComplete[T]) WithAfterDelete(hook func(id uint) error) *ServiceBuilderComplete[T] {
	b.builder.service.afterDelete = hook
	return b
}

func (b *ServiceBuilderComplete[T]) WithCustomSearch(fn func(query orm.Query, search string) orm.Query) *ServiceBuilderComplete[T] {
	b.builder.service.customSearch = fn
	return b
}

func (b *ServiceBuilderComplete[T]) WithCustomFilters(fn func(query orm.Query, filters map[string]interface{}) orm.Query) *ServiceBuilderComplete[T] {
	b.builder.service.customFilters = fn
	return b
}

func (b *ServiceBuilderComplete[T]) WithCustomQuery(fn func(query orm.Query) orm.Query) *ServiceBuilderComplete[T] {
	b.builder.service.customQuery = fn
	return b
}

// Build finalizes the service and ensures it implements the contract
func (b *ServiceBuilderComplete[T]) Build() CrudServiceContract {
	// Validate that all required fields are set
	if len(b.builder.service.searchFields) == 0 {
		panic("Search fields not set")
	}
	if len(b.builder.service.sortFields) == 0 {
		panic("Sort fields not set")
	}
	if len(b.builder.service.filterFields) == 0 {
		panic("Filter fields not set")
	}
	if len(b.builder.service.validationRules) == 0 {
		panic("Validation rules not set")
	}

	// Wrap the service to adapt the interface
	return &crudServiceAdapter[T]{
		service: b.builder.service,
	}
}

// crudServiceAdapter adapts GenericCrudService to implement CrudServiceContract
type crudServiceAdapter[T any] struct {
	service       *GenericCrudService[T]
	actualService interface{} // Reference to the actual service (e.g., BookService) for method resolution
}

// Implement CrudServiceContract methods by delegating to the wrapped service
func (a *crudServiceAdapter[T]) GetList(req ListRequest) (*PaginatedResult, error) {
	return a.service.GetList(req)
}

func (a *crudServiceAdapter[T]) GetByID(id uint) (interface{}, error) {
	return a.service.GetByID(id)
}

func (a *crudServiceAdapter[T]) Create(data map[string]interface{}) (interface{}, error) {
	return a.service.Create(data)
}

func (a *crudServiceAdapter[T]) Update(id uint, data map[string]interface{}) (interface{}, error) {
	return a.service.Update(id, data)
}

func (a *crudServiceAdapter[T]) Delete(id uint) error {
	return a.service.Delete(id)
}

// Search delegates to the wrapped service
func (a *crudServiceAdapter[T]) Search(query string, req ListRequest) (*PaginatedResult, error) {
	return a.service.Search(query, req)
}

func (a *crudServiceAdapter[T]) GetListAdvanced(req ListRequest, filters map[string]interface{}) (*PaginatedResult, error) {
	return a.service.GetListAdvanced(req, filters)
}

func (a *crudServiceAdapter[T]) GetSearchableFields() []string {
	return a.service.GetSearchableFields()
}

func (a *crudServiceAdapter[T]) GetSortableFields() []string {
	return a.service.GetSortableFields()
}

func (a *crudServiceAdapter[T]) GetFilterableFields() []string {
	return a.service.GetFilterableFields()
}

func (a *crudServiceAdapter[T]) GetValidationRules() map[string]interface{} {
	return a.service.GetValidationRules()
}

func (a *crudServiceAdapter[T]) GetColumnMapping() map[string]string {
	return a.service.GetColumnMapping()
}

func (a *crudServiceAdapter[T]) MapSortField(frontendField string) (string, bool) {
	// Directly delegate to the underlying service
	// Don't check actualService to avoid infinite recursion when BookService embeds the adapter
	return a.service.MapSortField(frontendField)
}

// ValidateSortField validates if a field can be sorted
func (a *crudServiceAdapter[T]) ValidateSortField(field string) bool {
	// Directly delegate to the underlying service
	return a.service.ValidateSortField(field)
}

// ValidateSortDirection validates sort direction
func (a *crudServiceAdapter[T]) ValidateSortDirection(direction string) bool {
	// Directly delegate to the underlying service
	return a.service.ValidateSortDirection(direction)
}

// GetDefaultSort returns the default sort configuration
func (a *crudServiceAdapter[T]) GetDefaultSort() (string, string) {
	// Directly delegate to the underlying service
	return a.service.GetDefaultSort()
}

// SetActualService sets the reference to the actual service for proper method resolution
func (a *crudServiceAdapter[T]) SetActualService(service interface{}) {
	a.actualService = service
	// Also set it on the underlying GenericCrudService
	a.service.SetActualService(service)
}

// Example usage that enforces compile-time checks:
//
// service := NewServiceBuilder[models.Book]("book", "id").
//     WithSearchFields("title", "author", "isbn").        // REQUIRED
//     WithSortFields("id", "title", "created_at").       // REQUIRED
//     WithFilterFields("status", "author").               // REQUIRED
//     WithValidationRules(map[string]interface{}{        // REQUIRED
//         "title": "required|string|max:255",
//         "author": "required|string|max:100",
//     }).
//     WithRelations("Tags", "Creator").                   // Optional
//     WithSoftDeletes().                                  // Optional
//     Build()
//
// If you skip any required step, you get a compile error!
