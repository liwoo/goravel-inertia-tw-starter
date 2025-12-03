package contracts

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/str"
	"smedi-sme-db/app/auth"
)

// GenericCrudService provides a complete CRUD service implementation with minimal code
type GenericCrudService[T any] struct {
	*BaseCrudService
	modelType       reflect.Type
	tableName       string
	searchFields    []string
	sortFields      []string
	filterFields    []string
	relations       []string
	validationRules map[string]interface{}

	// Reference to the actual service (e.g., BookService) for method resolution
	actualService interface{}

	// Permission scope configuration
	enableScopeFiltering bool
	scopeUserField       string // Field that contains the user ID (e.g., "created_by", "user_id", "owner_id")
	serviceRegistry      auth.ServiceRegistry

	// Customizable hooks
	beforeCreate func(data map[string]interface{}) error
	afterCreate  func(model *T) error
	beforeUpdate func(id uint, data map[string]interface{}) error
	afterUpdate  func(model *T) error
	beforeDelete func(id uint) error
	afterDelete  func(id uint) error

	// Custom query builders
	customQuery   func(query orm.Query) orm.Query
	customSearch  func(query orm.Query, search string) orm.Query
	customFilters func(query orm.Query, filters map[string]interface{}) orm.Query
}

// NewGenericCrudService creates a new generic CRUD service
func NewGenericCrudService[T any](resourceName string, primaryKey string) *GenericCrudService[T] {
	var model T
	modelType := reflect.TypeOf(model)

	// Use the provided resourceName as the table name
	// This allows services to specify custom table names (e.g., "sme_config" instead of "configs")
	tableName := resourceName

	return &GenericCrudService[T]{
		BaseCrudService:      NewBaseCrudService(resourceName, primaryKey),
		modelType:            modelType,
		tableName:            tableName,
		searchFields:         []string{},
		sortFields:           []string{"id", "created_at", "updated_at"},
		filterFields:         []string{},
		relations:            []string{},
		validationRules:      make(map[string]interface{}),
		enableScopeFiltering: false,
		scopeUserField:       "created_by", // Default to created_by
	}
}

// GetList retrieves paginated list of resources
func (s *GenericCrudService[T]) GetList(req ListRequest) (*PaginatedResult, error) {

	// Validate and sanitize request
	if err := s.ValidateListRequest(&req); err != nil {
		return nil, err
	}
	s.SanitizeListRequest(&req)

	// Build base query
	var model T
	query := facades.Orm().Query().Model(&model)

	// Load relations if configured
	for _, relation := range s.relations {
		query = query.With(relation)
	}

	// Apply custom query if set
	if s.customQuery != nil {
		query = s.customQuery(query)
	}

	// Apply permission-based scope filtering if enabled and context is available
	if s.enableScopeFiltering && req.Context != nil {
		var err error
		query, err = s.applyScopeFilter(req.Context, query, auth.PermissionRead)
		if err != nil {
			facades.Log().Warning("Failed to apply scope filter", map[string]interface{}{
				"error":   err.Error(),
				"service": s.tableName,
			})
			// Continue without scope filtering rather than failing the request
		}
	} else {
	}

	// Apply search
	if req.Search != "" {
		if s.customSearch != nil {
			query = s.customSearch(query, req.Search)
		} else {
			var err error
			query, err = s.ApplySearch(query, req.Search, s)
			if err != nil {
				return nil, err
			}
		}
	}

	// Apply filters from request
	var preservedCustomFilters interface{} = nil
	if len(req.Filters) > 0 {
		// Check for custom filters first
		if customFilters, hasCustom := req.Filters["__custom_filters"]; hasCustom {
			// Preserve custom filters for count query
			preservedCustomFilters = customFilters
			// Apply custom filters
			query = s.applyCustomFilters(query, customFilters)
			// Remove custom filters from regular filters
			delete(req.Filters, "__custom_filters")
		}

		// Apply field mapping to remaining filters
		mappedFilters := s.applyFieldMapping(req.Filters)

		if s.customFilters != nil {
			query = s.customFilters(query, mappedFilters)
		} else {
			// Default filter implementation
			for field, value := range mappedFilters {
				if s.ValidateFilterField(field) {
					// Convert string boolean values to actual booleans for boolean fields
					if strVal, ok := value.(string); ok {
						if strVal == "true" || strVal == "false" {
							// Check if this is a boolean field (common boolean field names)
							if strings.HasPrefix(field, "is_") || strings.HasPrefix(field, "has_") || field == "active" || field == "verified" || strings.HasSuffix(field, "_verified") {
								value = strVal == "true"
							}
						}
					}

					query = query.Where(field+" = ?", value)
				} else {
					if logger := facades.Log(); logger != nil {
						logger.Warning("Invalid filter field", map[string]interface{}{
							"field":            field,
							"filterableFields": s.filterFields,
						})
					}
				}
			}
		}
	}

	// Apply sorting - always use self (GenericCrudService) which has all the methods
	// The actualService is used internally by MapSortField to get column mappings
	query, err := s.ApplySort(query, req.Sort, req.Direction, s)
	if err != nil {
		facades.Log().Error("GenericCrudService sort error", map[string]interface{}{
			"service":   s.tableName,
			"sort":      req.Sort,
			"direction": req.Direction,
			"error":     err.Error(),
		})
		return nil, err
	}

	// Get total count before pagination
	var total int64
	// Create a new query instance for counting to avoid modifying the original
	var countModel T
	countQuery := facades.Orm().Query().Model(&countModel)

	// Re-apply the same conditions for counting
	// Load relations if configured
	for _, relation := range s.relations {
		countQuery = countQuery.With(relation)
	}

	// Apply custom query if set
	if s.customQuery != nil {
		countQuery = s.customQuery(countQuery)
	}

	// Apply permission-based scope filtering if enabled
	if s.enableScopeFiltering && req.Context != nil {
		var err error
		countQuery, err = s.applyScopeFilter(req.Context, countQuery, auth.PermissionRead)
		if err != nil {
			// Log but continue without scope filtering
			facades.Log().Warning("Failed to apply scope filter to count query", map[string]interface{}{
				"error":   err.Error(),
				"service": s.tableName,
			})
		}
	}

	// Apply search if present
	if req.Search != "" {
		if s.customSearch != nil {
			countQuery = s.customSearch(countQuery, req.Search)
		} else {
			countQuery, _ = s.ApplySearch(countQuery, req.Search, s)
		}
	}

	// Apply filters if present
	if len(req.Filters) > 0 {
		// Apply field mapping to filters for count query
		mappedCountFilters := s.applyFieldMapping(req.Filters)
		if s.customFilters != nil {
			countQuery = s.customFilters(countQuery, mappedCountFilters)
		} else {
			for field, value := range mappedCountFilters {
				if s.ValidateFilterField(field) {
					// Convert string boolean values to actual booleans for boolean fields
					if strVal, ok := value.(string); ok {
						if strVal == "true" || strVal == "false" {
							// Check if this is a boolean field (common boolean field names)
							if strings.HasPrefix(field, "is_") || strings.HasPrefix(field, "has_") || field == "active" || field == "verified" || strings.HasSuffix(field, "_verified") {
								value = strVal == "true"
							}
						}
					}
					countQuery = countQuery.Where(field+" = ?", value)
				}
			}
		}
	}

	// Apply preserved custom filters to count query
	if preservedCustomFilters != nil {
		countQuery = s.applyCustomFilters(countQuery, preservedCustomFilters)
	}

	total, err = countQuery.Count()
	if err != nil {
		facades.Log().Error("Count query failed", map[string]interface{}{
			"service": s.tableName,
			"error":   err.Error(),
		})
		return nil, err
	}

	// Apply pagination at database level
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	// Get paginated items
	var items []T

	// Note: SQL query will be logged by GORM if logging is enabled

	if err := query.Find(&items); err != nil {
		facades.Log().Error("Query Find failed", map[string]interface{}{
			"service": s.tableName,
			"error":   err.Error(),
		})
		return nil, err
	}

	// Convert items to interface slice
	data := make([]interface{}, len(items))
	for i, item := range items {
		data[i] = item
	}

	// Build pagination metadata
	pb := s.GetPaginationBuilder()
	lastPage := pb.CalculateLastPage(total, int64(req.PageSize))

	result := &PaginatedResult{
		Data:        data,
		Total:       total,
		PerPage:     req.PageSize,
		CurrentPage: req.Page,
		LastPage:    lastPage,
		From:        pb.CalculateFrom(offset, len(items)),
		To:          pb.CalculateTo(offset, len(items)),
		HasNext:     pb.HasNextPage(req.Page, lastPage),
		HasPrev:     pb.HasPrevPage(req.Page),
	}

	return result, nil
}

// GetListAdvanced retrieves paginated list with custom filters
func (s *GenericCrudService[T]) GetListAdvanced(req ListRequest, filters map[string]interface{}) (*PaginatedResult, error) {
	// Validate and sanitize request
	if err := s.ValidateListRequest(&req); err != nil {
		return nil, err
	}
	s.SanitizeListRequest(&req)

	// Build base query
	var model T
	query := facades.Orm().Query().Model(&model)

	// Load relations if configured
	for _, relation := range s.relations {
		query = query.With(relation)
	}

	// Apply custom query if set
	if s.customQuery != nil {
		query = s.customQuery(query)
	}

	// Apply filters
	if s.customFilters != nil {
		query = s.customFilters(query, filters)
	} else {
		// Default filter implementation
		for field, value := range filters {
			if s.ValidateFilterField(field) {
				// Convert string boolean values to actual booleans for boolean fields
				if strVal, ok := value.(string); ok {
					if strVal == "true" || strVal == "false" {
						// Check if this is a boolean field (common boolean field names)
						if strings.HasPrefix(field, "is_") || strings.HasPrefix(field, "has_") || field == "active" || field == "verified" || strings.HasSuffix(field, "_verified") {
							value = strVal == "true"
						}
					}
				}
				query = query.Where(field+" = ?", value)
			}
		}
	}

	// Apply search
	if req.Search != "" {
		if s.customSearch != nil {
			query = s.customSearch(query, req.Search)
		} else {
			var err error
			query, err = s.ApplySearch(query, req.Search, s)
			if err != nil {
				return nil, err
			}
		}
	}

	// Apply sorting - always use self (GenericCrudService) which has all the methods
	// The actualService is used internally by MapSortField to get column mappings
	query, err := s.ApplySort(query, req.Sort, req.Direction, s)
	if err != nil {
		facades.Log().Error("GenericCrudService sort error", map[string]interface{}{
			"service":   s.tableName,
			"sort":      req.Sort,
			"direction": req.Direction,
			"error":     err.Error(),
		})
		return nil, err
	}

	// Get total count before pagination
	var total int64
	// Create count query with same filters
	var countModel T
	countQuery := facades.Orm().Query().Model(&countModel)

	// Re-apply all conditions for counting
	for _, relation := range s.relations {
		countQuery = countQuery.With(relation)
	}
	if s.customQuery != nil {
		countQuery = s.customQuery(countQuery)
	}
	if s.customFilters != nil {
		countQuery = s.customFilters(countQuery, filters)
	} else {
		for field, value := range filters {
			if s.ValidateFilterField(field) {
				// Convert string boolean values to actual booleans for boolean fields
				if strVal, ok := value.(string); ok {
					if strVal == "true" || strVal == "false" {
						// Check if this is a boolean field (common boolean field names)
						if strings.HasPrefix(field, "is_") || strings.HasPrefix(field, "has_") || field == "active" || field == "verified" || strings.HasSuffix(field, "_verified") {
							value = strVal == "true"
						}
					}
				}
				countQuery = countQuery.Where(field+" = ?", value)
			}
		}
	}
	if req.Search != "" {
		if s.customSearch != nil {
			countQuery = s.customSearch(countQuery, req.Search)
		} else {
			countQuery, _ = s.ApplySearch(countQuery, req.Search, s)
		}
	}

	total, err = countQuery.Count()
	if err != nil {
		facades.Log().Error("Count query failed", map[string]interface{}{
			"service": s.tableName,
			"error":   err.Error(),
		})
		return nil, err
	}

	// Apply pagination at database level
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	// Get paginated items
	var items []T

	// Note: SQL query will be logged by GORM if logging is enabled

	if err := query.Find(&items); err != nil {
		facades.Log().Error("Query Find failed", map[string]interface{}{
			"service": s.tableName,
			"error":   err.Error(),
		})
		return nil, err
	}

	// Convert items to interface slice
	data := make([]interface{}, len(items))
	for i, item := range items {
		data[i] = item
	}

	// Build result
	pb := s.GetPaginationBuilder()
	lastPage := pb.CalculateLastPage(total, int64(req.PageSize))

	result := &PaginatedResult{
		Data:        data,
		Total:       total,
		PerPage:     req.PageSize,
		CurrentPage: req.Page,
		LastPage:    lastPage,
		From:        pb.CalculateFrom(offset, len(items)),
		To:          pb.CalculateTo(offset, len(items)),
		HasNext:     pb.HasNextPage(req.Page, lastPage),
		HasPrev:     pb.HasPrevPage(req.Page),
	}

	return result, nil
}

// GetByID retrieves a single resource by ID
func (s *GenericCrudService[T]) GetByID(id uint) (interface{}, error) {
	fmt.Printf("DEBUG: GetByID called with ID: %d, table: %s\n", id, s.BaseCrudService.tableName)
	if id == 0 {
		return nil, fmt.Errorf("invalid ID: %d", id)
	}

	var model T
	// Create a new instance for the query
	query := facades.Orm().Query()

	// Load relations if configured
	for _, relation := range s.relations {
		query = query.With(relation)
	}

	// Apply custom query if set
	if s.customQuery != nil {
		query = s.customQuery(query)
	}

	// Use qualified primary key (table.column) to avoid ambiguity when JOINs are used
	primaryKey := s.BaseCrudService.GetQualifiedPrimaryKey()
	if err := query.Where(primaryKey+" = ?", id).First(&model); err != nil {
		fmt.Printf("DEBUG: GetByID - Failed to find model with ID %d: %v\n", id, err)
		return nil, fmt.Errorf("%s not found: %w", s.BaseCrudService.tableName, err)
	}

	// Check if we got a valid model by checking if it has the expected ID
	// This is needed because GORM sometimes returns zero models instead of errors for soft-deleted records
	modelInterface := interface{}(&model)
	if idField := reflect.ValueOf(modelInterface).Elem().FieldByName("ID"); idField.IsValid() {
		if idField.Uint() == 0 {
			fmt.Printf("DEBUG: GetByID - Found zero model for ID %d (likely soft deleted)\n", id)
			return nil, fmt.Errorf("%s not found", s.BaseCrudService.tableName)
		}
		if idField.Uint() != uint64(id) {
			fmt.Printf("DEBUG: GetByID - ID mismatch: expected %d, got %d\n", id, idField.Uint())
			return nil, fmt.Errorf("%s not found", s.BaseCrudService.tableName)
		}
	}

	fmt.Printf("DEBUG: GetByID - Found valid model with ID %d\n", id)

	return &model, nil
}

// GetByIDWithContext retrieves a single resource by ID with permission scope check
func (s *GenericCrudService[T]) GetByIDWithContext(ctx http.Context, id uint) (interface{}, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid ID: %d", id)
	}

	var model T
	query := facades.Orm().Query().Model(&model)

	// Load relations if configured
	for _, relation := range s.relations {
		query = query.With(relation)
	}

	// Apply custom query if set
	if s.customQuery != nil {
		query = s.customQuery(query)
	}

	// Apply scope filtering if enabled
	if s.enableScopeFiltering && ctx != nil {
		var err error
		query, err = s.applyScopeFilter(ctx, query, auth.PermissionRead)
		if err != nil {
			return nil, fmt.Errorf("failed to apply scope filter: %w", err)
		}
	}

	// Use qualified primary key (table.column) to avoid ambiguity when JOINs are used
	primaryKey := s.BaseCrudService.GetQualifiedPrimaryKey()
	if err := query.Where(primaryKey+" = ?", id).First(&model); err != nil {
		return nil, fmt.Errorf("%s not found: %w", s.BaseCrudService.tableName, err)
	}

	return &model, nil
}

// Create creates a new resource
func (s *GenericCrudService[T]) Create(data map[string]interface{}) (interface{}, error) {
	// Apply field mapping first (frontend -> database)
	mappedData := s.applyFieldMapping(data)

	// Validate using validation rules on mapped data
	if err := s.validateWithRules(mappedData, false); err != nil {
		return nil, err
	}

	// Run before hook if set (with mapped data)
	if s.beforeCreate != nil {
		if err := s.beforeCreate(mappedData); err != nil {
			return nil, err
		}
	}

	// Create model instance
	model := new(T)

	// Use reflection to set fields (including embedded structs)
	if err := s.setFieldsRecursively(model, mappedData); err != nil {
		return nil, fmt.Errorf("failed to set fields: %w", err)
	}

	// Create using GORM
	if err := facades.Orm().Query().Create(model); err != nil {
		return nil, fmt.Errorf("failed to create %s: %w", s.BaseCrudService.tableName, err)
	}

	// Run after hook if set
	if s.afterCreate != nil {
		if err := s.afterCreate(model); err != nil {
			// Log error but don't fail
			facades.Log().Error("After create hook failed", map[string]interface{}{
				"resource": s.BaseCrudService.tableName,
				"error":    err.Error(),
			})
		}
	}

	// Reload with relations
	return s.GetByID(s.getIDFromModel(model))
}

// Update updates an existing resource
func (s *GenericCrudService[T]) Update(id uint, data map[string]interface{}) (interface{}, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid ID: %d", id)
	}

	// Check if exists
	_, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Apply field mapping first (frontend -> database)
	mappedData := s.applyFieldMapping(data)

	// Validate using validation rules on mapped data
	if err := s.validateWithRules(mappedData, true); err != nil {
		return nil, err
	}

	// Run before hook if set (with mapped data)
	if s.beforeUpdate != nil {
		if err := s.beforeUpdate(id, mappedData); err != nil {
			return nil, err
		}
	}

	// Update using GORM
	var model T
	if _, err := facades.Orm().Query().Model(&model).Where(s.BaseCrudService.GetPrimaryKey()+" = ?", id).Update(mappedData); err != nil {
		return nil, fmt.Errorf("failed to update %s: %w", s.BaseCrudService.tableName, err)
	}

	// Get updated model
	updated, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Run after hook if set
	if s.afterUpdate != nil {
		if updatedModel, ok := updated.(*T); ok {
			if err := s.afterUpdate(updatedModel); err != nil {
				// Log error but don't fail
				facades.Log().Error("After update hook failed", map[string]interface{}{
					"resource": s.BaseCrudService.tableName,
					"id":       id,
					"error":    err.Error(),
				})
			}
		} else {
			facades.Log().Warning("After update hook: type assertion failed", map[string]interface{}{
				"resource":    s.BaseCrudService.tableName,
				"id":          id,
				"actualType":  fmt.Sprintf("%T", updated),
				"expectedPtr": fmt.Sprintf("*%T", *new(T)),
			})
		}
	}

	return updated, nil
}

// Delete soft deletes a resource
func (s *GenericCrudService[T]) Delete(id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid ID: %d", id)
	}

	// Check if exists
	_, err := s.GetByID(id)
	if err != nil {
		return err
	}

	// Run before hook if set
	if s.beforeDelete != nil {
		if err := s.beforeDelete(id); err != nil {
			return err
		}
	}

	// Delete using GORM (soft delete)
	// Use Model() to ensure GORM uses the TableName() method from the model
	var model T
	result, err := facades.Orm().Query().
		Model(&model).
		Where(s.BaseCrudService.GetPrimaryKey()+" = ?", id).
		Update("deleted_at", time.Now())

	if err != nil {
		return fmt.Errorf("failed to delete %s: %w", s.tableName, err)
	}

	// Check if any row was actually deleted
	if result.RowsAffected == 0 {
		return fmt.Errorf("no %s found with id %d", s.tableName, id)
	}

	// Run after hook if set
	if s.afterDelete != nil {
		if err := s.afterDelete(id); err != nil {
			// Log error but don't fail
			facades.Log().Error("After delete hook failed", map[string]interface{}{
				"resource": s.BaseCrudService.tableName,
				"id":       id,
				"error":    err.Error(),
			})
		}
	}

	return nil
}

// Search implements search functionality
func (s *GenericCrudService[T]) Search(query string, req ListRequest) (*PaginatedResult, error) {
	if err := s.ValidateSearchQuery(query); err != nil {
		return nil, err
	}

	req.Search = query
	return s.GetList(req)
}

// Configuration methods

// SetActualService sets the reference to the actual service for proper method resolution
func (s *GenericCrudService[T]) SetActualService(service interface{}) *GenericCrudService[T] {
	s.actualService = service
	return s
}

// SetSearchFields sets the fields to search
func (s *GenericCrudService[T]) SetSearchFields(fields ...string) *GenericCrudService[T] {
	s.searchFields = fields
	return s
}

// SetSortFields sets the fields that can be sorted
func (s *GenericCrudService[T]) SetSortFields(fields ...string) *GenericCrudService[T] {
	s.sortFields = fields
	return s
}

// SetFilterFields sets the fields that can be filtered
func (s *GenericCrudService[T]) SetFilterFields(fields ...string) *GenericCrudService[T] {
	s.filterFields = fields
	return s
}

// SetRelations sets the relations to load
func (s *GenericCrudService[T]) SetRelations(relations ...string) *GenericCrudService[T] {
	s.relations = relations
	return s
}

// SetValidationRules sets the validation rules
func (s *GenericCrudService[T]) SetValidationRules(rules map[string]interface{}) *GenericCrudService[T] {
	s.validationRules = rules
	return s
}

// SetBeforeCreate sets a hook to run before create
func (s *GenericCrudService[T]) SetBeforeCreate(hook func(data map[string]interface{}) error) *GenericCrudService[T] {
	s.beforeCreate = hook
	return s
}

// SetAfterCreate sets a hook to run after create
func (s *GenericCrudService[T]) SetAfterCreate(hook func(model *T) error) *GenericCrudService[T] {
	s.afterCreate = hook
	return s
}

// SetBeforeUpdate sets a hook to run before update
func (s *GenericCrudService[T]) SetBeforeUpdate(hook func(id uint, data map[string]interface{}) error) *GenericCrudService[T] {
	s.beforeUpdate = hook
	return s
}

// SetAfterUpdate sets a hook to run after update
func (s *GenericCrudService[T]) SetAfterUpdate(hook func(model *T) error) *GenericCrudService[T] {
	s.afterUpdate = hook
	return s
}

// SetBeforeDelete sets a hook to run before delete
func (s *GenericCrudService[T]) SetBeforeDelete(hook func(id uint) error) *GenericCrudService[T] {
	s.beforeDelete = hook
	return s
}

// SetAfterDelete sets a hook to run after delete
func (s *GenericCrudService[T]) SetAfterDelete(hook func(id uint) error) *GenericCrudService[T] {
	s.afterDelete = hook
	return s
}

// SetCustomQuery sets a custom query builder
func (s *GenericCrudService[T]) SetCustomQuery(builder func(query orm.Query) orm.Query) *GenericCrudService[T] {
	s.customQuery = builder
	return s
}

// SetCustomSearch sets a custom search implementation
func (s *GenericCrudService[T]) SetCustomSearch(search func(query orm.Query, searchTerm string) orm.Query) *GenericCrudService[T] {
	s.customSearch = search
	return s
}

// SetCustomFilters sets a custom filter implementation
func (s *GenericCrudService[T]) SetCustomFilters(filters func(query orm.Query, filters map[string]interface{}) orm.Query) *GenericCrudService[T] {
	s.customFilters = filters
	return s
}

// EnableScopeFiltering enables permission-based scope filtering
func (s *GenericCrudService[T]) EnableScopeFiltering(serviceRegistry auth.ServiceRegistry, userField string) *GenericCrudService[T] {
	s.enableScopeFiltering = true
	s.serviceRegistry = serviceRegistry
	s.scopeUserField = userField
	return s
}

// SetScopeUserField sets the field name that contains the user ID for scope filtering
func (s *GenericCrudService[T]) SetScopeUserField(field string) *GenericCrudService[T] {
	s.scopeUserField = field
	return s
}

// PaginationServiceContract implementation
func (s *GenericCrudService[T]) GetPaginatedList(req ListRequest) (*PaginatedResult, error) {
	return s.GetList(req)
}

// Contract implementations

func (s *GenericCrudService[T]) GetSearchableFields() []string {
	return s.searchFields
}

func (s *GenericCrudService[T]) GetSortableFields() []string {
	return s.sortFields
}

func (s *GenericCrudService[T]) GetFilterableFields() []string {
	return s.filterFields
}

func (s *GenericCrudService[T]) ValidateFilterField(field string) bool {
	for _, f := range s.filterFields {
		if f == field {
			return true
		}
	}
	return false
}

// SortableServiceContract implementation
func (s *GenericCrudService[T]) ValidateSortField(field string) bool {
	for _, f := range s.sortFields {
		if f == field {
			return true
		}
	}
	return false
}

func (s *GenericCrudService[T]) MapSortField(frontendField string) (string, bool) {

	// Get the column mapping from the actual service if available
	var mapping map[string]string
	if s.actualService != nil {
		if mappingProvider, ok := s.actualService.(interface {
			GetColumnMapping() map[string]string
		}); ok {
			mapping = mappingProvider.GetColumnMapping()
		}
	}

	// Convert camelCase to snake_case for validation
	// e.g., "formalisationScore" -> "formalisation_score"
	snakeCaseField := str.Of(frontendField).Snake().String()

	// If we have a mapping, check if this field has a mapped name
	if mapping != nil {
		if dbField, exists := mapping[frontendField]; exists {
			// The mapped field could be a complex expression like "COALESCE(table.column, 0)"
			// We validate using the original frontend field or its snake_case equivalent
			// since sortFields contains the logical field names, not the fully qualified column names
			if s.ValidateSortField(frontendField) {
				return dbField, true
			}
			// Check snake_case version (e.g., formalisationScore -> formalisation_score)
			if s.ValidateSortField(snakeCaseField) {
				return dbField, true
			}
			if logger := facades.Log(); logger != nil {
				logger.Warning("Mapped field is NOT sortable", map[string]interface{}{
					"service":        s.tableName,
					"frontendField":  frontendField,
					"snakeCaseField": snakeCaseField,
					"dbField":        dbField,
					"sortableFields": s.sortFields,
				})
			}
		}
	}

	// If no mapping exists or field wasn't in mapping, check if the field itself is sortable
	if s.ValidateSortField(frontendField) {
		return frontendField, true
	}
	// Also check snake_case version
	if s.ValidateSortField(snakeCaseField) {
		return snakeCaseField, true
	}

	// Log warning if logger is available (may be nil in unit tests)
	if logger := facades.Log(); logger != nil {
		logger.Warning("Field not sortable", map[string]interface{}{
			"service":        s.tableName,
			"frontendField":  frontendField,
			"snakeCaseField": snakeCaseField,
			"sortableFields": s.sortFields,
		})
	}
	return "", false
}

func (s *GenericCrudService[T]) GetValidationRules() map[string]interface{} {
	return s.validationRules
}

func (s *GenericCrudService[T]) GetModel() interface{} {
	var model T
	return &model
}

// BulkOperationsContract implementation

// BulkCreate creates multiple resources
func (s *GenericCrudService[T]) BulkCreate(data []map[string]interface{}) ([]interface{}, error) {
	if err := s.ValidateBulkOperation([]uint{uint(len(data))}); err != nil {
		return nil, err
	}

	results := make([]interface{}, 0, len(data))
	for _, item := range data {
		result, err := s.Create(item)
		if err != nil {
			return nil, fmt.Errorf("bulk create failed: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}

// BulkUpdate updates multiple resources
func (s *GenericCrudService[T]) BulkUpdate(ids []uint, data map[string]interface{}) error {
	if err := s.ValidateBulkOperation(ids); err != nil {
		return err
	}

	for _, id := range ids {
		_, err := s.Update(id, data)
		if err != nil {
			return fmt.Errorf("bulk update failed for ID %d: %w", id, err)
		}
	}

	return nil
}

// BulkDelete deletes multiple resources
func (s *GenericCrudService[T]) BulkDelete(ids []uint) error {
	if err := s.ValidateBulkOperation(ids); err != nil {
		return err
	}

	for _, id := range ids {
		if err := s.Delete(id); err != nil {
			return fmt.Errorf("bulk delete failed for ID %d: %w", id, err)
		}
	}

	return nil
}

// ValidateBulkOperation validates bulk operation
func (s *GenericCrudService[T]) ValidateBulkOperation(ids []uint) error {
	if len(ids) == 0 {
		return fmt.Errorf("no items provided for bulk operation")
	}
	if len(ids) > 1000 {
		return fmt.Errorf("bulk operation limited to 1000 items")
	}
	return nil
}

// Helper methods

func (s *GenericCrudService[T]) validateWithRules(data map[string]interface{}, isUpdate bool) error {
	// Basic validation - can be enhanced with proper validator
	if !isUpdate {
		// Check required fields
		for field, rule := range s.validationRules {
			if ruleStr, ok := rule.(string); ok && strings.Contains(ruleStr, "required") {
				if value, exists := data[field]; !exists || value == "" {
					return fmt.Errorf("%s is required", field)
				}
			}
		}
	}
	return nil
}

func (s *GenericCrudService[T]) getIDFromModel(model *T) uint {
	modelValue := reflect.ValueOf(model).Elem()
	idField := modelValue.FieldByName("ID")
	if idField.IsValid() && idField.CanInterface() {
		if id, ok := idField.Interface().(uint); ok {
			return id
		}
	}
	return 0
}

// applyScopeFilter applies permission-based scope filtering to queries
func (s *GenericCrudService[T]) applyScopeFilter(ctx http.Context, query orm.Query, action auth.CorePermissionAction) (orm.Query, error) {
	if !s.enableScopeFiltering || s.serviceRegistry == "" {
		return query, nil
	}

	scopeHelper := auth.GetScopeHelper()
	return scopeHelper.ApplyScopeToQuery(ctx, query, s.serviceRegistry, action, s.scopeUserField)
}

// BuildFilterQuery applies filters to the query
func (s *GenericCrudService[T]) BuildFilterQuery(query interface{}, filters map[string]interface{}) interface{} {
	q, ok := query.(orm.Query)
	if !ok {
		return query
	}

	// Apply custom filters if set
	if s.customFilters != nil {
		return s.customFilters(q, filters)
	}

	// Apply default filters
	for field, value := range filters {
		// Skip empty values
		if value == nil || value == "" {
			continue
		}

		// Validate filter field
		if !s.ValidateFilterField(field) {
			continue
		}

		// Apply filter based on type
		switch v := value.(type) {
		case string:
			if v != "" {
				// Convert string boolean values to actual booleans for boolean fields
				if (v == "true" || v == "false") && (strings.HasPrefix(field, "is_") || strings.HasPrefix(field, "has_") || field == "active" || field == "verified" || strings.HasSuffix(field, "_verified")) {
					boolValue := v == "true"
					q = q.Where(field+" = ?", boolValue)
				} else {
					q = q.Where(field+" = ?", v)
				}
			}
		case int, int64, uint, uint64:
			q = q.Where(field+" = ?", v)
		case bool:
			q = q.Where(field+" = ?", v)
		case []interface{}:
			if len(v) > 0 {
				q = q.WhereIn(field, v)
			}
		default:
			q = q.Where(field+" = ?", v)
		}
	}

	return q
}

// CrudServiceConfiguration interface methods

// GetTableName returns the primary table name
func (s *GenericCrudService[T]) GetTableName() string {
	return s.tableName
}

// GetPrimaryKey returns the primary key field name
func (s *GenericCrudService[T]) GetPrimaryKey() string {
	return s.BaseCrudService.primaryKey
}

// GetColumnMapping returns frontend->database column mapping
func (s *GenericCrudService[T]) GetColumnMapping() map[string]string {
	// Default mapping - can be overridden by services
	return map[string]string{
		"id":         "id",
		"created_at": "created_at",
		"updated_at": "updated_at",
		"deleted_at": "deleted_at",
	}
}

// PaginationServiceContract interface methods

// ValidatePaginationParams ensures valid pagination parameters
func (s *GenericCrudService[T]) ValidatePaginationParams(page, pageSize int) error {
	if page < 1 {
		return fmt.Errorf("page must be greater than 0")
	}
	if pageSize < 1 || pageSize > s.GetMaxPageSize() {
		return fmt.Errorf("page size must be between 1 and %d", s.GetMaxPageSize())
	}
	return nil
}

// GetMaxPageSize returns the maximum allowed page size
func (s *GenericCrudService[T]) GetMaxPageSize() int {
	return 100
}

// GetDefaultPageSize returns the default page size
func (s *GenericCrudService[T]) GetDefaultPageSize() int {
	return 20
}

// SortableServiceContract interface methods

// GetDefaultSort returns the default sort configuration
func (s *GenericCrudService[T]) GetDefaultSort() (field string, direction string) {
	// Delegate to base service for consistency
	return s.BaseCrudService.GetDefaultSort()
}

// ValidateSortDirection validates sort direction
func (s *GenericCrudService[T]) ValidateSortDirection(direction string) bool {
	// Delegate to base service which handles case-insensitive validation
	return s.BaseCrudService.ValidateSortDirection(direction)
}

// SearchableServiceContract interface methods

// ValidateSearchQuery validates the search query
func (s *GenericCrudService[T]) ValidateSearchQuery(query string) error {
	if len(query) < 2 {
		return fmt.Errorf("search query must be at least 2 characters")
	}
	if len(query) > 100 {
		return fmt.Errorf("search query must not exceed 100 characters")
	}
	return nil
}

// setFieldsRecursively sets fields on a struct using reflection, including embedded structs
func (s *GenericCrudService[T]) setFieldsRecursively(model interface{}, data map[string]interface{}) error {
	modelValue := reflect.ValueOf(model).Elem()
	return s.setFieldsRecursivelyHelper(modelValue, data)
}

// setFieldsRecursivelyHelper is the recursive helper for setting fields
func (s *GenericCrudService[T]) setFieldsRecursivelyHelper(modelValue reflect.Value, data map[string]interface{}) error {
	modelType := modelValue.Type()

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldValue := modelValue.Field(i)

		// Handle embedded structs
		if field.Anonymous && fieldValue.Kind() == reflect.Struct {
			// Log embedded struct processing

			// Recursively set fields in embedded struct
			if err := s.setFieldsRecursivelyHelper(fieldValue, data); err != nil {
				return err
			}
			continue
		}

		// Get the field name from json tag or use lowercase field name
		fieldName := field.Tag.Get("json")
		if fieldName == "" || fieldName == "-" {
			fieldName = strings.ToLower(field.Name)
		} else {
			// Handle json tags with options (e.g., "field,omitempty")
			if idx := strings.Index(fieldName, ","); idx != -1 {
				fieldName = fieldName[:idx]
			}
		}

		// Skip if json tag is "-"
		if fieldName == "-" {
			continue
		}

		// Check if we have data for this field
		if value, exists := data[fieldName]; exists && fieldValue.CanSet() {
			// Log field setting for debugging

			// Handle different types of values
			if value == nil {
				// Set zero value for nil
				fieldValue.Set(reflect.Zero(fieldValue.Type()))
				continue
			}

			setValue := reflect.ValueOf(value)

			// Handle pointer fields
			if fieldValue.Kind() == reflect.Ptr {
				if setValue.Kind() == reflect.Ptr {
					fieldValue.Set(setValue)
				} else {
					// Create a new pointer and set the value
					newPtr := reflect.New(fieldValue.Type().Elem())
					if setValue.Type().ConvertibleTo(fieldValue.Type().Elem()) {
						newPtr.Elem().Set(setValue.Convert(fieldValue.Type().Elem()))
						fieldValue.Set(newPtr)
					}
				}
			} else if setValue.Type().ConvertibleTo(fieldValue.Type()) {
				fieldValue.Set(setValue.Convert(fieldValue.Type()))
			} else {
				// Try to handle special cases like converting float64 to uint
				switch fieldValue.Kind() {
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					if floatVal, ok := value.(float64); ok {
						fieldValue.SetUint(uint64(floatVal))
					}
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if floatVal, ok := value.(float64); ok {
						fieldValue.SetInt(int64(floatVal))
					}
				}
			}
		}
	}

	return nil
}

// applyCustomFilters applies custom filter conditions to a query
func (s *GenericCrudService[T]) applyCustomFilters(query orm.Query, customFilters interface{}) orm.Query {

	switch filter := customFilters.(type) {
	case *FilterCondition:
		sql, args := filter.ToSQL()
		if sql != "" {
			// Use Where with raw SQL instead of WhereRaw
			query = query.Where(sql, args...)
		}

	case *CompoundFilter:
		sql, args := filter.ToSQL()
		if sql != "" {
			// Use Where with raw SQL instead of WhereRaw
			query = query.Where(sql, args...)
		}

	default:
		if logger := facades.Log(); logger != nil {
			logger.Warning("Unknown custom filter type", map[string]interface{}{
				"type": fmt.Sprintf("%T", customFilters),
			})
		}
	}

	return query
}

// applyFieldMapping applies frontend->database field mapping to data
func (s *GenericCrudService[T]) applyFieldMapping(data map[string]interface{}) map[string]interface{} {
	// Get the mapping from the actual service if available
	var mapping map[string]string
	if s.actualService != nil {
		if mappingProvider, ok := s.actualService.(interface {
			GetColumnMapping() map[string]string
		}); ok {
			mapping = mappingProvider.GetColumnMapping()
		}
	}

	// If no mapping available, return data as-is
	if mapping == nil || len(mapping) == 0 {
		return data
	}

	// Apply mapping
	result := make(map[string]interface{})
	for key, value := range data {
		// Check if this field has a mapping
		if mappedKey, exists := mapping[key]; exists {
			// Try to parse date-like strings for any field
			if dateStr, ok := value.(string); ok && value != nil && value != "" {
				// Check if the string looks like a date using common patterns
				// This will match formats like: 2024-01-15, 2024-01-15 10:30:00, 2024-01-15T10:30:00Z, etc.
				isDateLike := strings.Contains(dateStr, "-") && len(dateStr) >= 10 &&
					(dateStr[0] >= '0' && dateStr[0] <= '9') // Starts with a digit

				if isDateLike {
					// Try multiple common date/datetime formats
					formats := []string{
						"2006-01-02 15:04:05",           // MySQL datetime format
						"2006-01-02",                    // Date only (YYYY-MM-DD)
						time.RFC3339,                    // ISO 8601 (2006-01-02T15:04:05Z07:00)
						time.RFC3339Nano,                // ISO 8601 with nanoseconds
						"2006-01-02T15:04:05",           // ISO 8601 without timezone
						"2006-01-02 15:04:05.999999999", // MySQL datetime with microseconds
					}

					parsed := false
					for _, format := range formats {
						if parsedTime, err := time.Parse(format, dateStr); err == nil {
							result[mappedKey] = parsedTime
							parsed = true
							break
						}
					}

					// If all parsing attempts fail, keep the original string value
					// The database driver might be able to handle it
					if !parsed {
						result[mappedKey] = value
					}
				} else {
					// Not a date-like string, keep as-is
					result[mappedKey] = value
				}
			} else {
				result[mappedKey] = value
			}
		} else {
			// Keep unmapped fields as-is
			result[key] = value
		}
	}

	return result
}
