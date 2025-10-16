package contracts

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
	httpvalidate "github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/validation"
)

// CrudController provides a complete CRUD implementation with compile-time enforcement
// This is a unified controller that combines utilities, CRUD operations, and type safety
type CrudController[T any, C CreateRequestContract, U UpdateRequestContract] struct {
	service      CrudServiceContract
	resourceName string

	// Configuration
	maxPageSize      int
	defaultPageSize  int
	allowedPageSizes []int

	// Customizable hooks
	beforeIndex  func(ctx http.Context) error
	beforeShow   func(ctx http.Context, id uint) error
	beforeStore  func(ctx http.Context, data map[string]interface{}) error
	beforeUpdate func(ctx http.Context, id uint, data map[string]interface{}) error
	beforeDelete func(ctx http.Context, id uint) error
	beforeSearch func(ctx http.Context, query string) error

	// Custom authorization (REQUIRED - enforced by builder)
	CheckAuth func(ctx http.Context, action string, resource interface{}) error

	// Response customization
	afterStore  func(ctx http.Context, result interface{}) http.Response
	afterUpdate func(ctx http.Context, result interface{}) http.Response
}

// ============================================================================
// BUILDER PATTERN WITH TYPE-STATE FOR COMPILE-TIME ENFORCEMENT
// ============================================================================

// Builder states
type (
	needsAuthChecker struct{}
	readyToBuild     struct{}
)

// CrudControllerBuilder uses type-state pattern to enforce compile-time checks
type CrudControllerBuilder[T any, C CreateRequestContract, U UpdateRequestContract, State any] struct {
	resourceName string
	service      CrudServiceContract
	authChecker  func(ctx http.Context, action string, resource interface{}) error
}

// NewCrudController starts the builder - requires auth checker to be set before building
func NewCrudController[T any, C CreateRequestContract, U UpdateRequestContract](
	resourceName string,
	service CrudServiceContract,
) *CrudControllerBuilder[T, C, U, needsAuthChecker] {
	return &CrudControllerBuilder[T, C, U, needsAuthChecker]{
		resourceName: resourceName,
		service:      service,
	}
}

// WithAuthChecker sets the authorization checker (REQUIRED)
// This method is ONLY available when State = needsAuthChecker
func (b *CrudControllerBuilder[T, C, U, needsAuthChecker]) WithAuthChecker(
	checker func(ctx http.Context, action string, resource interface{}) error,
) *CrudControllerBuilder[T, C, U, readyToBuild] {
	return &CrudControllerBuilder[T, C, U, readyToBuild]{
		resourceName: b.resourceName,
		service:      b.service,
		authChecker:  checker,
	}
}

// Build creates the controller - ONLY available when State = readyToBuild
// This ensures at compile time that all required methods have been called
func (b *CrudControllerBuilder[T, C, U, readyToBuild]) Build() *CrudController[T, C, U] {
	return &CrudController[T, C, U]{
		service:          b.service,
		resourceName:     b.resourceName,
		CheckAuth:        b.authChecker,
		maxPageSize:      100,
		defaultPageSize:  20,
		allowedPageSizes: []int{5, 10, 20, 30, 50, 100},
	}
}

// ============================================================================
// CRUD OPERATIONS
// ============================================================================

// Index handles listing resources with pagination, filtering, sorting, and search
// Swagger annotations should be added in concrete controller wrapper methods
func (c *CrudController[T, C, U]) Index(ctx http.Context) http.Response {
	// Run before hook if set
	if c.beforeIndex != nil {
		if err := c.beforeIndex(ctx); err != nil {
			return c.BadRequestResponse(ctx, err.Error(), nil)
		}
	}

	// Authorization check
	if err := c.CheckAuth(ctx, "viewAny", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
	}

	// Validate pagination request
	req, err := c.ValidatePaginationRequest(ctx)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid pagination parameters", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Get paginated list
	result, err := c.service.GetList(*req)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve "+c.resourceName+" list: "+err.Error())
	}

	// Build response
	response := c.BuildPaginatedResponse(result, req)
	return c.SuccessResponse(ctx, response, c.resourceName+" list retrieved successfully")
}

// Show handles retrieving a single resource by ID
// Swagger annotations should be added in concrete controller wrapper methods
func (c *CrudController[T, C, U]) Show(ctx http.Context) http.Response {
	// Validate ID parameter
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid "+c.resourceName+" ID", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Run before hook if set
	if c.beforeShow != nil {
		if err := c.beforeShow(ctx, id); err != nil {
			return c.BadRequestResponse(ctx, err.Error(), nil)
		}
	}

	// Get the resource
	resource, err := c.service.GetByID(id)
	if err != nil {
		return c.ResourceNotFoundResponse(ctx, c.resourceName, id)
	}

	// Authorization check
	if err := c.CheckAuth(ctx, "view", resource); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
	}

	return c.SuccessResponse(ctx, resource, c.resourceName+" retrieved successfully")
}

// Store handles creating a new resource with validation
// Swagger annotations should be added in concrete controller wrapper methods
func (c *CrudController[T, C, U]) Store(ctx http.Context) http.Response {
	// Authorization check
	if err := c.CheckAuth(ctx, "create", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
	}

	// Create a new instance of the create request type
	var createReq C
	if err := ctx.Request().Bind(&createReq); err != nil {
		return c.ValidationErrorResponse(ctx, map[string]interface{}{
			"binding_error": err.Error(),
		})
	}

	// Manual validation using the request's Rules method
	rules := createReq.Rules(ctx)

	// Convert the bound request to validation data
	requestData := createReq.ToCreateData()

	// Get custom messages and attributes from the request
	messages := createReq.Messages(ctx)
	attributes := createReq.Attributes(ctx)

	// Build validation options
	options := []httpvalidate.Option{}
	if len(messages) > 0 {
		options = append(options, validation.Messages(messages))
	}
	if len(attributes) > 0 {
		options = append(options, validation.Attributes(attributes))
	}

	// Validate using facades with options
	validator, err := facades.Validation().Make(
		requestData,
		rules,
		options...,
	)
	if err != nil {
		return c.ValidationErrorResponse(ctx, map[string]interface{}{
			"validation_error": err.Error(),
		})
	}
	if validator.Fails() {
		errors := make(map[string]interface{})
		for field, messages := range validator.Errors().All() {
			errors[field] = messages
		}
		return c.ValidationErrorResponse(ctx, errors)
	}

	// Authorization check on the request
	if err := createReq.Authorize(ctx); err != nil {
		return c.ForbiddenResponse(ctx, "Request authorization failed: "+err.Error())
	}

	// Prepare for validation
	if err := createReq.PrepareForValidation(ctx); err != nil {
		return c.ValidationErrorResponse(ctx, map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Transform to data
	data := createReq.ToCreateData()

	// Run before hook if set
	if c.beforeStore != nil {
		if err := c.beforeStore(ctx, data); err != nil {
			return c.BadRequestResponse(ctx, err.Error(), nil)
		}
	}

	// Post validation hook
	if err := createReq.PassedValidation(ctx); err != nil {
		return c.ValidationErrorResponse(ctx, map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Create the resource
	resource, err := c.service.Create(data)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to create "+c.resourceName+": "+err.Error())
	}

	// Use custom response if set
	if c.afterStore != nil {
		return c.afterStore(ctx, resource)
	}

	return c.ResourceCreatedResponse(ctx, resource, c.resourceName)
}

// Update handles updating an existing resource
// Swagger annotations should be added in concrete controller wrapper methods
func (c *CrudController[T, C, U]) Update(ctx http.Context) http.Response {
	// Validate ID parameter
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid "+c.resourceName+" ID", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Check if resource exists
	resource, err := c.service.GetByID(id)
	if err != nil {
		return c.ResourceNotFoundResponse(ctx, c.resourceName, id)
	}

	// Authorization check
	if err := c.CheckAuth(ctx, "update", resource); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
	}

	// Create a new instance of the update request type
	var updateReq U
	if err := ctx.Request().Bind(&updateReq); err != nil {
		return c.ValidationErrorResponse(ctx, map[string]interface{}{
			"binding_error": err.Error(),
		})
	}

	// Manual validation using the request's Rules method
	rules := updateReq.Rules(ctx)

	// Convert the bound request to validation data
	requestData := updateReq.ToUpdateData()

	// Get custom messages and attributes from the request
	messages := updateReq.Messages(ctx)
	attributes := updateReq.Attributes(ctx)

	// Build validation options
	options := []httpvalidate.Option{}
	if len(messages) > 0 {
		options = append(options, validation.Messages(messages))
	}
	if len(attributes) > 0 {
		options = append(options, validation.Attributes(attributes))
	}

	// Validate using facades with options
	validator, err := facades.Validation().Make(
		requestData,
		rules,
		options...,
	)
	if err != nil {
		return c.ValidationErrorResponse(ctx, map[string]interface{}{
			"validation_error": err.Error(),
		})
	}
	if validator.Fails() {
		errors := make(map[string]interface{})
		for field, messages := range validator.Errors().All() {
			errors[field] = messages
		}
		return c.ValidationErrorResponse(ctx, errors)
	}

	// Authorization check on the request
	if err := updateReq.Authorize(ctx); err != nil {
		return c.ForbiddenResponse(ctx, "Request authorization failed: "+err.Error())
	}

	// Prepare for validation
	if err := updateReq.PrepareForValidation(ctx); err != nil {
		return c.ValidationErrorResponse(ctx, map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Transform to data
	data := updateReq.ToUpdateData()

	// Run before hook if set
	if c.beforeUpdate != nil {
		if err := c.beforeUpdate(ctx, id, data); err != nil {
			return c.BadRequestResponse(ctx, err.Error(), nil)
		}
	}

	// Post validation hook
	if err := updateReq.PassedValidation(ctx); err != nil {
		return c.ValidationErrorResponse(ctx, map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Update the resource
	updatedResource, err := c.service.Update(id, data)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to update "+c.resourceName+": "+err.Error())
	}

	// Use custom response if set
	if c.afterUpdate != nil {
		return c.afterUpdate(ctx, updatedResource)
	}

	return c.ResourceUpdatedResponse(ctx, updatedResource, c.resourceName)
}

// Delete handles deleting a resource by ID
// Swagger annotations should be added in concrete controller wrapper methods
func (c *CrudController[T, C, U]) Delete(ctx http.Context) http.Response {
	// Validate ID parameter
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid "+c.resourceName+" ID", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Check if resource exists
	resource, err := c.service.GetByID(id)
	if err != nil {
		return c.ResourceNotFoundResponse(ctx, c.resourceName, id)
	}

	// Authorization check
	if err := c.CheckAuth(ctx, "delete", resource); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
	}

	// Run before hook if set
	if c.beforeDelete != nil {
		if err := c.beforeDelete(ctx, id); err != nil {
			return c.BadRequestResponse(ctx, err.Error(), nil)
		}
	}

	// Delete the resource
	err = c.service.Delete(id)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to delete "+c.resourceName+": "+err.Error())
	}

	return c.ResourceDeletedResponse(ctx, c.resourceName, id)
}

// Search handles full-text search with pagination
// Swagger annotations should be added in concrete controller wrapper methods
func (c *CrudController[T, C, U]) Search(ctx http.Context) http.Response {
	// Run before hook if set
	query := ctx.Request().Query("q")
	if c.beforeSearch != nil {
		if err := c.beforeSearch(ctx, query); err != nil {
			return c.BadRequestResponse(ctx, err.Error(), nil)
		}
	}

	// Authorization check
	if err := c.CheckAuth(ctx, "viewAny", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
	}

	// Validate search query
	if query == "" {
		return c.BadRequestResponse(ctx, "Search query is required", nil)
	}

	// Validate query length (minimum 2 characters)
	if len(strings.TrimSpace(query)) < 2 {
		return c.BadRequestResponse(ctx, "Search query must be at least 2 characters long", nil)
	}

	// Validate pagination request
	req, err := c.ValidatePaginationRequest(ctx)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid pagination parameters", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Perform search
	result, err := c.service.Search(query, *req)
	if err != nil {
		// Check if it's a validation error
		if strings.Contains(err.Error(), "must be at least") || strings.Contains(err.Error(), "validation") {
			return c.BadRequestResponse(ctx, err.Error(), nil)
		}
		return c.InternalErrorResponse(ctx, "Search failed: "+err.Error())
	}

	// Build response with search context
	response := c.BuildPaginatedResponse(result, req)
	response["query"] = query
	response["searchable_fields"] = c.service.GetSearchableFields()

	return c.SuccessResponse(ctx, response, "Search completed successfully")
}

// FilterMetadata returns available filters, searchable and sortable fields
// Swagger annotations should be added in concrete controller wrapper methods
func (c *CrudController[T, C, U]) FilterMetadata(ctx http.Context) http.Response {
	// Authorization check
	if err := c.CheckAuth(ctx, "viewAny", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
	}

	// Get filter definitions from the service
	var filterDefs []FilterDefinition

	// Check if service provides filter definitions
	if filterProvider, ok := c.service.(interface {
		GetFilterDefinitions() []FilterDefinition
	}); ok {
		filterDefs = filterProvider.GetFilterDefinitions()
	}

	// Generate metadata
	metadata := GenerateFilterMetadata(filterDefs)

	// Add additional metadata
	metadata["resource"] = c.resourceName
	metadata["searchable_fields"] = c.service.GetSearchableFields()
	metadata["sortable_fields"] = c.service.GetSortableFields()
	metadata["filterable_fields"] = c.service.GetFilterableFields()

	return c.SuccessResponse(ctx, metadata, "Filter metadata retrieved successfully")
}

// ============================================================================
// HOOK SETTERS
// ============================================================================

// SetBeforeIndex sets a hook to run before index action
func (c *CrudController[T, C, U]) SetBeforeIndex(hook func(ctx http.Context) error) {
	c.beforeIndex = hook
}

// SetBeforeShow sets a hook to run before show action
func (c *CrudController[T, C, U]) SetBeforeShow(hook func(ctx http.Context, id uint) error) {
	c.beforeShow = hook
}

// SetBeforeStore sets a hook to run before store action
func (c *CrudController[T, C, U]) SetBeforeStore(hook func(ctx http.Context, data map[string]interface{}) error) {
	c.beforeStore = hook
}

// SetBeforeUpdate sets a hook to run before update action
func (c *CrudController[T, C, U]) SetBeforeUpdate(hook func(ctx http.Context, id uint, data map[string]interface{}) error) {
	c.beforeUpdate = hook
}

// SetBeforeDelete sets a hook to run before delete action
func (c *CrudController[T, C, U]) SetBeforeDelete(hook func(ctx http.Context, id uint) error) {
	c.beforeDelete = hook
}

// SetBeforeSearch sets a hook to run before search action
func (c *CrudController[T, C, U]) SetBeforeSearch(hook func(ctx http.Context, query string) error) {
	c.beforeSearch = hook
}

// SetAfterStore sets a custom response handler for after store
func (c *CrudController[T, C, U]) SetAfterStore(hook func(ctx http.Context, result interface{}) http.Response) {
	c.afterStore = hook
}

// SetAfterUpdate sets a custom response handler for after update
func (c *CrudController[T, C, U]) SetAfterUpdate(hook func(ctx http.Context, result interface{}) http.Response) {
	c.afterUpdate = hook
}

// ============================================================================
// UTILITY METHODS (from BaseCrudController)
// ============================================================================

// ValidateID validates and parses an ID from the request
func (c *CrudController[T, C, U]) ValidateID(ctx http.Context, param string) (uint, error) {
	idStr := ctx.Request().Route(param)
	if idStr == "" {
		return 0, fmt.Errorf("missing %s parameter", param)
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format", param)
	}

	return uint(id), nil
}

// GetSearchableFields returns the searchable fields from the service
func (c *CrudController[T, C, U]) GetSearchableFields() []string {
	return c.service.GetSearchableFields()
}

// GetValidationRules returns the validation rules from the service
func (c *CrudController[T, C, U]) GetValidationRules() map[string]interface{} {
	return c.service.GetValidationRules()
}

// ValidatePaginationRequest validates and parses pagination parameters
func (c *CrudController[T, C, U]) ValidatePaginationRequest(ctx http.Context) (*ListRequest, error) {
	req := &ListRequest{}

	// Set the HTTP context for permission checks
	req.Context = ctx

	// Parse pagination parameters
	req.Page = ctx.Request().QueryInt("page", 1)
	req.PageSize = ctx.Request().QueryInt("pageSize", c.defaultPageSize)
	req.Search = ctx.Request().Query("search", "")
	req.Sort = ctx.Request().Query("sort", "")
	req.Direction = ctx.Request().Query("direction", "")

	// Normalize direction to uppercase for consistency
	if req.Direction != "" {
		req.Direction = strings.ToUpper(req.Direction)
	}

	// Parse filters from query parameters
	req.Filters = make(map[string]interface{})

	// Get all query parameters as strings
	queries := ctx.Request().Queries()

	// List of known non-filter parameters
	knownParams := map[string]bool{
		"page":      true,
		"pageSize":  true,
		"search":    true,
		"sort":      true,
		"direction": true,
		"filters":   true, // Special parameter for custom filters JSON
	}

	// Check for custom filters JSON parameter
	if filtersJSON := ctx.Request().Query("filters", ""); filtersJSON != "" {
		// Parse custom filters from JSON
		if customFilters, err := parseCustomFilters(filtersJSON); err == nil {
			req.Filters["__custom_filters"] = customFilters
		}
	}

	// Add all other query parameters as simple filters
	for key := range queries {
		if !knownParams[key] {
			// Get the value as a string
			value := ctx.Request().Query(key, "")
			if value != "" {
				req.Filters[key] = value
			}
		}
	}

	// Validate pagination parameters
	if req.Page <= 0 {
		return nil, fmt.Errorf("page must be greater than 0")
	}

	if req.PageSize <= 0 {
		req.PageSize = c.defaultPageSize
	}

	if req.PageSize > c.maxPageSize {
		return nil, fmt.Errorf("pageSize cannot exceed %d", c.maxPageSize)
	}

	// Validate page size is in allowed sizes
	validPageSize := false
	for _, size := range c.allowedPageSizes {
		if req.PageSize == size {
			validPageSize = true
			break
		}
	}

	if !validPageSize {
		req.PageSize = c.defaultPageSize
	}

	// Validate sort direction
	if req.Direction != "" {
		upper := strings.ToUpper(req.Direction)
		if upper != "ASC" && upper != "DESC" {
			req.Direction = "DESC"
		} else {
			req.Direction = upper
		}
	}

	// Set defaults
	req.SetDefaults()

	return req, nil
}

// BuildPaginatedResponse builds a paginated response
func (c *CrudController[T, C, U]) BuildPaginatedResponse(result *PaginatedResult, request *ListRequest) map[string]interface{} {
	response := NewPaginatedResponse(result, request)
	return response.ToMap()
}

// ============================================================================
// RESPONSE HELPERS
// ============================================================================

// SuccessResponse returns a success response
func (c *CrudController[T, C, U]) SuccessResponse(ctx http.Context, data interface{}, message string) http.Response {
	response := ResponseFormat{
		Success: true,
		Data:    data,
		Message: message,
	}
	return ctx.Response().Json(http.StatusOK, response)
}

// BadRequestResponse returns a bad request response
func (c *CrudController[T, C, U]) BadRequestResponse(ctx http.Context, message string, errors map[string]interface{}) http.Response {
	response := ResponseFormat{
		Success: false,
		Message: message,
		Errors:  errors,
	}
	return ctx.Response().Json(http.StatusBadRequest, response)
}

// NotFoundResponse returns a not found response
func (c *CrudController[T, C, U]) NotFoundResponse(ctx http.Context, message string) http.Response {
	response := ResponseFormat{
		Success: false,
		Message: message,
	}
	return ctx.Response().Json(http.StatusNotFound, response)
}

// ForbiddenResponse returns a forbidden response
func (c *CrudController[T, C, U]) ForbiddenResponse(ctx http.Context, message string) http.Response {
	response := ResponseFormat{
		Success: false,
		Message: message,
	}
	return ctx.Response().Json(http.StatusForbidden, response)
}

// ValidationErrorResponse returns a validation error response
func (c *CrudController[T, C, U]) ValidationErrorResponse(ctx http.Context, errors map[string]interface{}) http.Response {
	response := ResponseFormat{
		Success: false,
		Message: "Validation failed",
		Errors:  errors,
	}
	return ctx.Response().Json(http.StatusUnprocessableEntity, response)
}

// InternalErrorResponse returns an internal server error response
func (c *CrudController[T, C, U]) InternalErrorResponse(ctx http.Context, message string) http.Response {
	response := ResponseFormat{
		Success: false,
		Message: message,
	}
	return ctx.Response().Json(http.StatusInternalServerError, response)
}

// ResourceNotFoundResponse returns a resource not found response
func (c *CrudController[T, C, U]) ResourceNotFoundResponse(ctx http.Context, resourceType string, id uint) http.Response {
	message := fmt.Sprintf("%s with ID %d not found", strings.Title(resourceType), id)
	return c.NotFoundResponse(ctx, message)
}

// ResourceCreatedResponse returns a resource created response
func (c *CrudController[T, C, U]) ResourceCreatedResponse(ctx http.Context, resource interface{}, resourceType string) http.Response {
	message := fmt.Sprintf("%s created successfully", strings.Title(resourceType))
	response := ResponseFormat{
		Success: true,
		Data:    resource,
		Message: message,
	}
	return ctx.Response().Json(http.StatusCreated, response)
}

// ResourceUpdatedResponse returns a resource updated response
func (c *CrudController[T, C, U]) ResourceUpdatedResponse(ctx http.Context, resource interface{}, resourceType string) http.Response {
	message := fmt.Sprintf("%s updated successfully", strings.Title(resourceType))
	return c.SuccessResponse(ctx, resource, message)
}

// ResourceDeletedResponse returns a resource deleted response
func (c *CrudController[T, C, U]) ResourceDeletedResponse(ctx http.Context, resourceType string, id uint) http.Response {
	// Return 204 No Content for successful delete operations
	return ctx.Response().NoContent()
}

// ============================================================================
// CONFIGURATION
// ============================================================================

// SetPaginationConfig sets pagination configuration
func (c *CrudController[T, C, U]) SetPaginationConfig(defaultPageSize, maxPageSize int, allowedSizes []int) {
	if defaultPageSize > 0 {
		c.defaultPageSize = defaultPageSize
	}
	if maxPageSize > 0 {
		c.maxPageSize = maxPageSize
	}
	if len(allowedSizes) > 0 {
		c.allowedPageSizes = allowedSizes
	}
}

// GetResourceType returns the resource type
func (c *CrudController[T, C, U]) GetResourceType() string {
	return c.resourceName
}

// ============================================================================
// INTERNAL HELPERS
// ============================================================================

// parseCustomFilters parses custom filter JSON string into filter conditions
func parseCustomFilters(filtersJSON string) (interface{}, error) {
	var filterData map[string]interface{}
	if err := json.Unmarshal([]byte(filtersJSON), &filterData); err != nil {
		return nil, fmt.Errorf("invalid filter JSON: %v", err)
	}

	// Check if it's a compound filter
	if _, hasLogic := filterData["logic"]; hasLogic {
		return ParseCompoundFilter(filterData)
	}

	// Otherwise, it's a simple filter
	return ParseFilterQuery(filterData)
}
