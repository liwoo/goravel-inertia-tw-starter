package contracts

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
	httpvalidate "github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/validation"
)

// EnforcedCrudController provides a complete CRUD implementation with compile-time enforcement
// The generic constraints ensure that C and U implement the required interfaces
type EnforcedCrudController[T any, C CreateRequestContract, U UpdateRequestContract] struct {
	*BaseCrudController
	service      CrudServiceContract
	resourceName string

	// Customizable hooks
	beforeIndex  func(ctx http.Context) error
	beforeShow   func(ctx http.Context, id uint) error
	beforeStore  func(ctx http.Context, data map[string]interface{}) error
	beforeUpdate func(ctx http.Context, id uint, data map[string]interface{}) error
	beforeDelete func(ctx http.Context, id uint) error
	beforeSearch func(ctx http.Context, query string) error

	// Custom authorization
	CheckAuth func(ctx http.Context, action string, resource interface{}) error

	// Response customization
	afterStore  func(ctx http.Context, result interface{}) http.Response
	afterUpdate func(ctx http.Context, result interface{}) http.Response

	// These are no longer needed as the request types themselves have the methods
	// The controller will call methods directly on the request types
}

// NewEnforcedCrudController creates a new enforced CRUD controller
// This ensures at compile time that C implements CreateRequestContract and U implements UpdateRequestContract
func NewEnforcedCrudController[T any, C CreateRequestContract, U UpdateRequestContract](
	resourceName string,
	service CrudServiceContract,
) *EnforcedCrudController[T, C, U] {
	return &EnforcedCrudController[T, C, U]{
		BaseCrudController: NewBaseCrudController(resourceName),
		service:            service,
		resourceName:       resourceName,
	}
}

// Index GET /resources
func (c *EnforcedCrudController[T, C, U]) Index(ctx http.Context) http.Response {
	// Run before hook if set
	if c.beforeIndex != nil {
		if err := c.beforeIndex(ctx); err != nil {
			return c.BadRequestResponse(ctx, err.Error(), nil)
		}
	}

	// Default authorization check
	if c.CheckAuth != nil {
		if err := c.CheckAuth(ctx, "viewAny", nil); err != nil {
			return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
		}
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

// Show GET /resources/{id}
func (c *EnforcedCrudController[T, C, U]) Show(ctx http.Context) http.Response {
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

	// Default authorization check
	if c.CheckAuth != nil {
		if err := c.CheckAuth(ctx, "view", resource); err != nil {
			return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
		}
	}

	return c.SuccessResponse(ctx, resource, c.resourceName+" retrieved successfully")
}

// Store POST /resources
func (c *EnforcedCrudController[T, C, U]) Store(ctx http.Context) http.Response {
	// Default authorization check
	if c.CheckAuth != nil {
		if err := c.CheckAuth(ctx, "create", nil); err != nil {
			return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
		}
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
	fmt.Printf("DEBUG: EnforcedCrudController calling service.Create - resource: %s, service_type: %T\n", c.resourceName, c.service)
	facades.Log().Info("EnforcedCrudController calling service.Create", map[string]interface{}{
		"resource":     c.resourceName,
		"service_type": fmt.Sprintf("%T", c.service),
	})
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

// Update PUT /resources/{id}
func (c *EnforcedCrudController[T, C, U]) Update(ctx http.Context) http.Response {
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

	// Default authorization check
	if c.CheckAuth != nil {
		if err := c.CheckAuth(ctx, "update", resource); err != nil {
			return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
		}
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
	fmt.Printf("DEBUG: EnforcedCrudController calling service.Update - resource: %s, id: %d, service_type: %T\n", c.resourceName, id, c.service)
	facades.Log().Info("EnforcedCrudController calling service.Update", map[string]interface{}{
		"resource":     c.resourceName,
		"id":           id,
		"service_type": fmt.Sprintf("%T", c.service),
	})
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

// Delete DELETE /resources/{id}
func (c *EnforcedCrudController[T, C, U]) Delete(ctx http.Context) http.Response {
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

	// Default authorization check
	if c.CheckAuth != nil {
		if err := c.CheckAuth(ctx, "delete", resource); err != nil {
			return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
		}
	}

	// Run before hook if set
	if c.beforeDelete != nil {
		if err := c.beforeDelete(ctx, id); err != nil {
			return c.BadRequestResponse(ctx, err.Error(), nil)
		}
	}

	// Delete the resource
	fmt.Printf("DEBUG: EnforcedCrudController calling service.Delete - resource: %s, id: %d\n", c.resourceName, id)
	err = c.service.Delete(id)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to delete "+c.resourceName+": "+err.Error())
	}
	fmt.Printf("DEBUG: EnforcedCrudController delete completed - resource: %s, id: %d\n", c.resourceName, id)

	return c.ResourceDeletedResponse(ctx, c.resourceName, id)
}

// Search GET /resources/search
func (c *EnforcedCrudController[T, C, U]) Search(ctx http.Context) http.Response {
	// Run before hook if set
	query := ctx.Request().Query("q")
	if c.beforeSearch != nil {
		if err := c.beforeSearch(ctx, query); err != nil {
			return c.BadRequestResponse(ctx, err.Error(), nil)
		}
	}

	// Default authorization check
	if c.CheckAuth != nil {
		if err := c.CheckAuth(ctx, "viewAny", nil); err != nil {
			return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
		}
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

// SetBeforeIndex sets a hook to run before index action
func (c *EnforcedCrudController[T, C, U]) SetBeforeIndex(hook func(ctx http.Context) error) {
	c.beforeIndex = hook
}

// SetBeforeShow sets a hook to run before show action
func (c *EnforcedCrudController[T, C, U]) SetBeforeShow(hook func(ctx http.Context, id uint) error) {
	c.beforeShow = hook
}

// SetBeforeStore sets a hook to run before store action
func (c *EnforcedCrudController[T, C, U]) SetBeforeStore(hook func(ctx http.Context, data map[string]interface{}) error) {
	c.beforeStore = hook
}

// SetBeforeUpdate sets a hook to run before update action
func (c *EnforcedCrudController[T, C, U]) SetBeforeUpdate(hook func(ctx http.Context, id uint, data map[string]interface{}) error) {
	c.beforeUpdate = hook
}

// SetBeforeDelete sets a hook to run before delete action
func (c *EnforcedCrudController[T, C, U]) SetBeforeDelete(hook func(ctx http.Context, id uint) error) {
	c.beforeDelete = hook
}

// SetBeforeSearch sets a hook to run before search action
func (c *EnforcedCrudController[T, C, U]) SetBeforeSearch(hook func(ctx http.Context, query string) error) {
	c.beforeSearch = hook
}

// SetAfterStore sets a custom response handler for after store
func (c *EnforcedCrudController[T, C, U]) SetAfterStore(hook func(ctx http.Context, result interface{}) http.Response) {
	c.afterStore = hook
}

// SetAfterUpdate sets a custom response handler for after update
func (c *EnforcedCrudController[T, C, U]) SetAfterUpdate(hook func(ctx http.Context, result interface{}) http.Response) {
	c.afterUpdate = hook
}

// Helper methods

// ValidateID validates and parses an ID from the request
func (c *EnforcedCrudController[T, C, U]) ValidateID(ctx http.Context, param string) (uint, error) {
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
func (c *EnforcedCrudController[T, C, U]) GetSearchableFields() []string {
	return c.service.GetSearchableFields()
}

// GetValidationRules returns the validation rules from the service
func (c *EnforcedCrudController[T, C, U]) GetValidationRules() map[string]interface{} {
	return c.service.GetValidationRules()
}

// FilterMetadata GET /resources/filters
func (c *EnforcedCrudController[T, C, U]) FilterMetadata(ctx http.Context) http.Response {
	// Default authorization check
	if c.CheckAuth != nil {
		if err := c.CheckAuth(ctx, "viewAny", nil); err != nil {
			return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
		}
	}

	// Get filter definitions from the service or controller
	var filterDefs []FilterDefinition

	// Check if service provides filter definitions
	if filterProvider, ok := c.service.(interface {
		GetFilterDefinitions() []FilterDefinition
	}); ok {
		filterDefs = filterProvider.GetFilterDefinitions()
	}

	// If no definitions from service, use controller's definitions
	if len(filterDefs) == 0 {
		filterDefs = c.GetFilterDefinitions()
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
