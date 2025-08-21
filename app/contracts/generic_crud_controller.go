package contracts

import (
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/http"
)

// GenericCrudController provides a complete CRUD implementation with minimal code
type GenericCrudController[T any, C any, U any] struct {
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

	// Request binding
	bindCreate func(ctx http.Context) (C, error)
	bindUpdate func(ctx http.Context, id uint) (U, error)

	// Data transformation
	transformCreate func(C) map[string]interface{}
	transformUpdate func(U) map[string]interface{}
}

// NewGenericCrudController creates a new generic CRUD controller
func NewGenericCrudController[T any, C any, U any](
	resourceName string,
	service CrudServiceContract,
) *GenericCrudController[T, C, U] {
	return &GenericCrudController[T, C, U]{
		BaseCrudController: NewBaseCrudController(resourceName),
		service:            service,
		resourceName:       resourceName,
	}
}

// Index GET /resources
func (c *GenericCrudController[T, C, U]) Index(ctx http.Context) http.Response {
	// Run before hook if set
	if c.beforeIndex != nil {
		if err := c.beforeIndex(ctx); err != nil {
			return c.ForbiddenResponse(ctx, err.Error())
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

	// Get list using service
	result, err := c.service.GetList(*req)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve "+c.resourceName+"s: "+err.Error())
	}

	// Build standardized paginated response
	response := c.BuildPaginatedResponse(result, req)
	return c.SuccessResponse(ctx, response, c.resourceName+"s retrieved successfully")
}

// Show GET /resources/{id}
func (c *GenericCrudController[T, C, U]) Show(ctx http.Context) http.Response {
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
			return c.ForbiddenResponse(ctx, err.Error())
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

	return c.SuccessResponse(ctx, resource, c.resourceName+" details retrieved successfully")
}

// Store POST /resources
func (c *GenericCrudController[T, C, U]) Store(ctx http.Context) http.Response {
	// Default authorization check
	if c.CheckAuth != nil {
		if err := c.CheckAuth(ctx, "create", nil); err != nil {
			return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
		}
	}

	// Validate create request
	data, err := c.ValidateCreateRequest(ctx)
	if err != nil {
		return c.ValidationErrorResponse(ctx, map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Run before hook if set
	if c.beforeStore != nil {
		if err := c.beforeStore(ctx, data); err != nil {
			return c.BadRequestResponse(ctx, err.Error(), nil)
		}
	}

	// Create the resource using validated data
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
func (c *GenericCrudController[T, C, U]) Update(ctx http.Context) http.Response {
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

	// Validate update request
	data, err := c.ValidateUpdateRequest(ctx, id)
	if err != nil {
		return c.ValidationErrorResponse(ctx, map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Run before hook if set
	if c.beforeUpdate != nil {
		if err := c.beforeUpdate(ctx, id, data); err != nil {
			return c.BadRequestResponse(ctx, err.Error(), nil)
		}
	}

	// Update the resource using validated data
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
func (c *GenericCrudController[T, C, U]) Delete(ctx http.Context) http.Response {
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
	err = c.service.Delete(id)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to delete "+c.resourceName+": "+err.Error())
	}

	return c.ResourceDeletedResponse(ctx, c.resourceName, id)
}

// Search GET /resources/search
func (c *GenericCrudController[T, C, U]) Search(ctx http.Context) http.Response {
	// Run before hook if set
	if c.beforeSearch != nil {
		query := ctx.Request().Query("q")
		if err := c.beforeSearch(ctx, query); err != nil {
			return c.ForbiddenResponse(ctx, err.Error())
		}
	}

	// Validate search request
	req, err := c.ValidateSearchRequest(ctx)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid search parameters", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Perform search using service
	searchService, ok := c.service.(SearchableServiceContract)
	if !ok {
		return c.InternalErrorResponse(ctx, "Search not supported for "+c.resourceName)
	}

	result, err := searchService.Search(req.Query, req.ToListRequest())
	if err != nil {
		return c.InternalErrorResponse(ctx, "Search failed: "+err.Error())
	}

	// Build standardized search response
	response := c.BuildSearchResponse(result, req)
	return c.SuccessResponse(ctx, response, fmt.Sprintf("Found %d results for '%s'", result.Total, req.Query))
}

// ValidateCreateRequest validates and transforms the create request
func (c *GenericCrudController[T, C, U]) ValidateCreateRequest(ctx http.Context) (map[string]interface{}, error) {
	if c.bindCreate == nil || c.transformCreate == nil {
		// Fallback to generic binding
		var data map[string]interface{}
		if err := ctx.Request().Bind(&data); err != nil {
			return nil, fmt.Errorf("data binding failed: %w", err)
		}
		return data, nil
	}

	// Use custom binding
	createReq, err := c.bindCreate(ctx)
	if err != nil {
		return nil, err
	}

	return c.transformCreate(createReq), nil
}

// ValidateUpdateRequest validates and transforms the update request
func (c *GenericCrudController[T, C, U]) ValidateUpdateRequest(ctx http.Context, id uint) (map[string]interface{}, error) {
	if c.bindUpdate == nil || c.transformUpdate == nil {
		// Fallback to generic binding
		var data map[string]interface{}
		if err := ctx.Request().Bind(&data); err != nil {
			return nil, fmt.Errorf("data binding failed: %w", err)
		}
		return data, nil
	}

	// Use custom binding
	updateReq, err := c.bindUpdate(ctx, id)
	if err != nil {
		return nil, err
	}

	return c.transformUpdate(updateReq), nil
}

// SetBeforeIndex sets a hook to run before index action
func (c *GenericCrudController[T, C, U]) SetBeforeIndex(hook func(ctx http.Context) error) *GenericCrudController[T, C, U] {
	c.beforeIndex = hook
	return c
}

// SetBeforeShow sets a hook to run before show action
func (c *GenericCrudController[T, C, U]) SetBeforeShow(hook func(ctx http.Context, id uint) error) *GenericCrudController[T, C, U] {
	c.beforeShow = hook
	return c
}

// SetBeforeStore sets a hook to run before store action
func (c *GenericCrudController[T, C, U]) SetBeforeStore(hook func(ctx http.Context, data map[string]interface{}) error) *GenericCrudController[T, C, U] {
	c.beforeStore = hook
	return c
}

// SetBeforeUpdate sets a hook to run before update action
func (c *GenericCrudController[T, C, U]) SetBeforeUpdate(hook func(ctx http.Context, id uint, data map[string]interface{}) error) *GenericCrudController[T, C, U] {
	c.beforeUpdate = hook
	return c
}

// SetBeforeDelete sets a hook to run before delete action
func (c *GenericCrudController[T, C, U]) SetBeforeDelete(hook func(ctx http.Context, id uint) error) *GenericCrudController[T, C, U] {
	c.beforeDelete = hook
	return c
}

// SetBeforeSearch sets a hook to run before search action
func (c *GenericCrudController[T, C, U]) SetBeforeSearch(hook func(ctx http.Context, query string) error) *GenericCrudController[T, C, U] {
	c.beforeSearch = hook
	return c
}

// SetAuthCheck sets a custom authorization check
func (c *GenericCrudController[T, C, U]) SetAuthCheck(check func(ctx http.Context, action string, resource interface{}) error) *GenericCrudController[T, C, U] {
	c.CheckAuth = check
	return c
}

// SetRequestBindings sets custom request binding functions
func (c *GenericCrudController[T, C, U]) SetRequestBindings(
	bindCreate func(ctx http.Context) (C, error),
	transformCreate func(C) map[string]interface{},
	bindUpdate func(ctx http.Context, id uint) (U, error),
	transformUpdate func(U) map[string]interface{},
) *GenericCrudController[T, C, U] {
	c.bindCreate = bindCreate
	c.transformCreate = transformCreate
	c.bindUpdate = bindUpdate
	c.transformUpdate = transformUpdate
	return c
}

// SetAfterStore sets a custom response for store action
func (c *GenericCrudController[T, C, U]) SetAfterStore(hook func(ctx http.Context, result interface{}) http.Response) *GenericCrudController[T, C, U] {
	c.afterStore = hook
	return c
}

// SetAfterUpdate sets a custom response for update action
func (c *GenericCrudController[T, C, U]) SetAfterUpdate(hook func(ctx http.Context, result interface{}) http.Response) *GenericCrudController[T, C, U] {
	c.afterUpdate = hook
	return c
}

// Helper to get ID from route
func (c *GenericCrudController[T, C, U]) ValidateID(ctx http.Context, param string) (uint, error) {
	idStr := ctx.Request().Route(param)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format")
	}
	return uint(id), nil
}
