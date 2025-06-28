package contracts

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
)

// BaseCrudController provides common implementations for CRUD controllers
// Controllers MUST embed this and implement the abstract methods
type BaseCrudController struct {
	resourceType     string
	maxPageSize      int
	defaultPageSize  int
	allowedPageSizes []int
}

// NewBaseCrudController creates a new base CRUD controller
func NewBaseCrudController(resourceType string) *BaseCrudController {
	return &BaseCrudController{
		resourceType:     resourceType,
		maxPageSize:      100,
		defaultPageSize:  20,
		allowedPageSizes: []int{5, 10, 20, 30, 50, 100}, // More flexible options
	}
}

// PAGINATION CONTRACT IMPLEMENTATION (enforced)

func (c *BaseCrudController) ValidatePaginationRequest(ctx http.Context) (*ListRequest, error) {
	req := &ListRequest{}
	
	// Parse pagination parameters
	req.Page = ctx.Request().QueryInt("page", 1)
	req.PageSize = ctx.Request().QueryInt("pageSize", c.defaultPageSize)
	req.Search = ctx.Request().Query("search", "")
	req.Sort = ctx.Request().Query("sort", "")
	req.Direction = ctx.Request().Query("direction", "")
	
	// Parse filters
	req.Filters = make(map[string]interface{})
	
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

func (c *BaseCrudController) GetPaginationDefaults() (page int, pageSize int, maxPageSize int) {
	return 1, c.defaultPageSize, c.maxPageSize
}

func (c *BaseCrudController) BuildPaginatedResponse(result *PaginatedResult, request *ListRequest) map[string]interface{} {
	return map[string]interface{}{
		"data": result.Data,
		"pagination": map[string]interface{}{
			"current_page": result.CurrentPage,
			"last_page":    result.LastPage,
			"per_page":     result.PerPage,
			"total":        result.Total,
			"from":         result.From,
			"to":           result.To,
			"has_next":     result.HasNext,
			"has_prev":     result.HasPrev,
		},
		"filters": map[string]interface{}{
			"page":      request.Page,
			"pageSize":  request.PageSize,
			"search":    request.Search,
			"sort":      request.Sort,
			"direction": request.Direction,
			"filters":   request.Filters,
		},
	}
}

// VALIDATION CONTRACT IMPLEMENTATION (enforced)

func (c *BaseCrudController) ValidateID(ctx http.Context, paramName string) (uint, error) {
	idStr := ctx.Request().Route(paramName)
	if idStr == "" {
		return 0, fmt.Errorf("%s parameter is required", paramName)
	}
	
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: must be a positive integer", paramName)
	}
	
	if id == 0 {
		return 0, fmt.Errorf("invalid %s: must be greater than 0", paramName)
	}
	
	return uint(id), nil
}

// RESPONSE CONTRACT IMPLEMENTATION (enforced)

func (c *BaseCrudController) SuccessResponse(ctx http.Context, data interface{}, message string) http.Response {
	response := ResponseFormat{
		Success: true,
		Data:    data,
		Message: message,
	}
	return ctx.Response().Json(http.StatusOK, response)
}

func (c *BaseCrudController) CreatedResponse(ctx http.Context, data interface{}, message string) http.Response {
	response := ResponseFormat{
		Success: true,
		Data:    data,
		Message: message,
	}
	return ctx.Response().Json(http.StatusCreated, response)
}

func (c *BaseCrudController) NoContentResponse(ctx http.Context, message string) http.Response {
	response := ResponseFormat{
		Success: true,
		Message: message,
	}
	return ctx.Response().Json(http.StatusNoContent, response)
}

func (c *BaseCrudController) BadRequestResponse(ctx http.Context, message string, errors map[string]interface{}) http.Response {
	response := ResponseFormat{
		Success: false,
		Message: message,
		Errors:  errors,
	}
	return ctx.Response().Json(http.StatusBadRequest, response)
}

func (c *BaseCrudController) NotFoundResponse(ctx http.Context, message string) http.Response {
	response := ResponseFormat{
		Success: false,
		Message: message,
	}
	return ctx.Response().Json(http.StatusNotFound, response)
}

func (c *BaseCrudController) ForbiddenResponse(ctx http.Context, message string) http.Response {
	response := ResponseFormat{
		Success: false,
		Message: message,
	}
	return ctx.Response().Json(http.StatusForbidden, response)
}

func (c *BaseCrudController) ValidationErrorResponse(ctx http.Context, errors map[string]interface{}) http.Response {
	response := ResponseFormat{
		Success: false,
		Message: "Validation failed",
		Errors:  errors,
	}
	return ctx.Response().Json(http.StatusUnprocessableEntity, response)
}

func (c *BaseCrudController) InternalErrorResponse(ctx http.Context, message string) http.Response {
	response := ResponseFormat{
		Success: false,
		Message: message,
	}
	return ctx.Response().Json(http.StatusInternalServerError, response)
}

// SPECIALIZED CRUD RESPONSES

func (c *BaseCrudController) ResourceNotFoundResponse(ctx http.Context, resourceType string, id uint) http.Response {
	message := fmt.Sprintf("%s with ID %d not found", strings.Title(resourceType), id)
	return c.NotFoundResponse(ctx, message)
}

func (c *BaseCrudController) ResourceCreatedResponse(ctx http.Context, resource interface{}, resourceType string) http.Response {
	message := fmt.Sprintf("%s created successfully", strings.Title(resourceType))
	return c.CreatedResponse(ctx, resource, message)
}

func (c *BaseCrudController) ResourceUpdatedResponse(ctx http.Context, resource interface{}, resourceType string) http.Response {
	message := fmt.Sprintf("%s updated successfully", strings.Title(resourceType))
	return c.SuccessResponse(ctx, resource, message)
}

func (c *BaseCrudController) ResourceDeletedResponse(ctx http.Context, resourceType string, id uint) http.Response {
	message := fmt.Sprintf("%s with ID %d deleted successfully", strings.Title(resourceType), id)
	return c.NoContentResponse(ctx, message)
}

// CONFIGURATION

func (c *BaseCrudController) SetPaginationConfig(defaultPageSize, maxPageSize int, allowedSizes []int) {
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

func (c *BaseCrudController) GetResourceType() string {
	return c.resourceType
}

// METADATA GENERATION

func (c *BaseCrudController) GenerateMetadata(supportedActions []string, requiredPerms []string, validationRules map[string]interface{}) ControllerMetadata {
	return ControllerMetadata{
		ResourceType:     c.resourceType,
		SupportedActions: supportedActions,
		RequiredPerms:    requiredPerms,
		ValidationRules:  validationRules,
		PaginationConfig: PaginationConfig{
			DefaultPageSize: c.defaultPageSize,
			MaxPageSize:     c.maxPageSize,
			AllowedSizes:    c.allowedPageSizes,
		},
		ResponseFormats: []string{"json"},
	}
}

// BASE PAGE CONTROLLER for Inertia.js pages

type BasePageController struct {
	*BaseCrudController
	pageComponent string
}

func NewBasePageController(resourceType, pageComponent string) *BasePageController {
	return &BasePageController{
		BaseCrudController: NewBaseCrudController(resourceType),
		pageComponent:      pageComponent,
	}
}

// PAGE RESPONSE CONTRACT IMPLEMENTATION

func (c *BasePageController) BuildPageProps(data interface{}, filters interface{}, permissions map[string]bool, meta map[string]interface{}) map[string]interface{} {
	props := map[string]interface{}{
		"data":        data,
		"filters":     filters,
		"permissions": permissions,
	}
	
	if meta != nil {
		for key, value := range meta {
			props[key] = value
		}
	}
	
	// Add page metadata
	props["meta"] = c.GetPageMetadata()
	
	return props
}

func (c *BasePageController) GetPageMetadata() map[string]interface{} {
	return map[string]interface{}{
		"version":       "1.0.0",
		"component":     c.pageComponent,
		"resourceType":  c.resourceType,
		"timestamp":     fmt.Sprintf("%d", 1234567890), // Could use time.Now().Unix()
		"pagination":    c.GetPaginationConfig(),
	}
}

func (c *BasePageController) GetPaginationConfig() map[string]interface{} {
	return map[string]interface{}{
		"defaultPageSize": c.defaultPageSize,
		"maxPageSize":     c.maxPageSize,
		"allowedSizes":    c.allowedPageSizes,
	}
}

func (c *BasePageController) ValidatePageRequest(ctx http.Context) (*ListRequest, error) {
	return c.ValidatePaginationRequest(ctx)
}

// VALIDATION HELPERS

func ValidateControllerImplementation(controller interface{}) ControllerValidationResult {
	result := ControllerValidationResult{
		Valid:            true,
		Errors:           []string{},
		MissingMethods:   []string{},
		InvalidResponses: []string{},
	}
	
	// Check if controller implements CrudControllerContract
	if _, ok := controller.(CrudControllerContract); !ok {
		result.Valid = false
		result.Errors = append(result.Errors, "controller does not implement CrudControllerContract interface")
	}
	
	// Check if controller implements individual contracts
	if _, ok := controller.(PaginationControllerContract); !ok {
		result.Valid = false
		result.MissingMethods = append(result.MissingMethods, "PaginationControllerContract")
	}
	
	if _, ok := controller.(ValidationControllerContract); !ok {
		result.Valid = false
		result.MissingMethods = append(result.MissingMethods, "ValidationControllerContract")
	}
	
	if _, ok := controller.(ResponseControllerContract); !ok {
		result.Valid = false
		result.MissingMethods = append(result.MissingMethods, "ResponseControllerContract")
	}
	
	return result
}