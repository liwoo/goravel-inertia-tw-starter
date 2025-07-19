package contracts

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/facades"
)

// GenericCrudService provides a complete CRUD service implementation with minimal code
type GenericCrudService[T any] struct {
	*BaseCrudService
	modelType      reflect.Type
	tableName      string
	searchFields   []string
	sortFields     []string
	filterFields   []string
	relations      []string
	validationRules map[string]interface{}
	
	// Customizable hooks
	beforeCreate   func(data map[string]interface{}) error
	afterCreate    func(model *T) error
	beforeUpdate   func(id uint, data map[string]interface{}) error
	afterUpdate    func(model *T) error
	beforeDelete   func(id uint) error
	afterDelete    func(id uint) error
	
	// Custom query builders
	customQuery    func(query orm.Query) orm.Query
	customSearch   func(query orm.Query, search string) orm.Query
	customFilters  func(query orm.Query, filters map[string]interface{}) orm.Query
}

// NewGenericCrudService creates a new generic CRUD service
func NewGenericCrudService[T any](resourceName string, primaryKey string) *GenericCrudService[T] {
	var model T
	modelType := reflect.TypeOf(model)
	
	// Extract table name from model type
	tableName := strings.ToLower(modelType.Name()) + "s"
	
	return &GenericCrudService[T]{
		BaseCrudService: NewBaseCrudService(resourceName, primaryKey),
		modelType:       modelType,
		tableName:       tableName,
		searchFields:    []string{},
		sortFields:      []string{"id", "created_at", "updated_at"},
		filterFields:    []string{},
		relations:       []string{},
		validationRules: make(map[string]interface{}),
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
	
	// Apply sorting
	query, err := s.ApplySort(query, req.Sort, req.Direction, s)
	if err != nil {
		return nil, err
	}
	
	// Get all items
	var items []T
	if err := query.Find(&items); err != nil {
		return nil, err
	}
	
	// Use pagination utility
	result := PaginateSliceWithConverter(
		items,
		req.Page,
		req.PageSize,
		func(item T) interface{} { return item },
	)
	
	// Add additional pagination metadata
	pb := s.GetPaginationBuilder()
	result.HasNext = pb.HasNextPage(result.CurrentPage, result.LastPage)
	result.HasPrev = pb.HasPrevPage(result.CurrentPage)
	
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
	
	// Apply sorting
	query, err := s.ApplySort(query, req.Sort, req.Direction, s)
	if err != nil {
		return nil, err
	}
	
	// Get all items
	var items []T
	if err := query.Find(&items); err != nil {
		return nil, err
	}
	
	// Use pagination utility
	result := PaginateSliceWithConverter(
		items,
		req.Page,
		req.PageSize,
		func(item T) interface{} { return item },
	)
	
	return result, nil
}

// GetByID retrieves a single resource by ID
func (s *GenericCrudService[T]) GetByID(id uint) (interface{}, error) {
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
	
	if err := query.Where(s.BaseCrudService.GetPrimaryKey()+" = ?", id).First(&model); err != nil {
		return nil, fmt.Errorf("%s not found: %w", s.BaseCrudService.tableName, err)
	}
	
	return model, nil
}

// Create creates a new resource
func (s *GenericCrudService[T]) Create(data map[string]interface{}) (interface{}, error) {
	// Validate using validation rules
	if err := s.validateWithRules(data, false); err != nil {
		return nil, err
	}
	
	// Run before hook if set
	if s.beforeCreate != nil {
		if err := s.beforeCreate(data); err != nil {
			return nil, err
		}
	}
	
	// Create model instance
	model := new(T)
	
	// Use reflection to set fields
	modelValue := reflect.ValueOf(model).Elem()
	modelType := modelValue.Type()
	
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldName := field.Tag.Get("json")
		if fieldName == "" {
			fieldName = strings.ToLower(field.Name)
		}
		
		if value, exists := data[fieldName]; exists {
			fieldValue := modelValue.Field(i)
			if fieldValue.CanSet() {
				setValue := reflect.ValueOf(value)
				if setValue.Type().ConvertibleTo(fieldValue.Type()) {
					fieldValue.Set(setValue.Convert(fieldValue.Type()))
				}
			}
		}
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
	
	// Validate using validation rules
	if err := s.validateWithRules(data, true); err != nil {
		return nil, err
	}
	
	// Run before hook if set
	if s.beforeUpdate != nil {
		if err := s.beforeUpdate(id, data); err != nil {
			return nil, err
		}
	}
	
	// Update using GORM
	var model T
	if _, err := facades.Orm().Query().Model(&model).Where(s.BaseCrudService.GetPrimaryKey()+" = ?", id).Update(data); err != nil {
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
	var model T
	if _, err := facades.Orm().Query().Model(&model).Where(s.BaseCrudService.GetPrimaryKey()+" = ?", id).Delete(&model); err != nil {
		return fmt.Errorf("failed to delete %s: %w", s.BaseCrudService.tableName, err)
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
	// By default, we assume frontend fields map directly to database fields
	// Override this in specific services if needed
	if s.ValidateSortField(frontendField) {
		return frontendField, true
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
				q = q.Where(field+" = ?", v)
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
	return "id", "desc"
}

// ValidateSortDirection validates sort direction
func (s *GenericCrudService[T]) ValidateSortDirection(direction string) bool {
	return direction == "asc" || direction == "desc"
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

