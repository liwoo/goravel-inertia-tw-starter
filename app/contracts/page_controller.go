package contracts

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
)

// TYPED PAGE PROPS STRUCTURES

// PageProps represents the typed structure for Inertia.js page props
type PageProps struct {
	Data        PaginatedData  `json:"data"`
	Filters     FilterState    `json:"filters"`
	Permissions PermissionsMap `json:"permissions"`
	Meta        PageMetadata   `json:"meta"`
	// Additional props can be added via composition
	Stats interface{}            `json:"stats,omitempty"`
	Extra map[string]interface{} `json:"-"` // For custom fields
}

// PermissionMatrixPageProps represents specialized props for permission matrix pages
type PermissionMatrixPageProps struct {
	PageProps
	AllPermissions []interface{}          `json:"allPermissions"`
	Services       []interface{}          `json:"services"`
	Actions        []interface{}          `json:"actions"`
	MatrixData     map[string]interface{} `json:"matrixData"`
	Title          string                 `json:"title"`
	Subtitle       string                 `json:"subtitle"`
}

// ToMap converts PermissionMatrixPageProps to map for Inertia rendering
func (p PermissionMatrixPageProps) ToMap() map[string]interface{} {
	// Get base props map
	result := p.PageProps.ToMap()

	// Add permission matrix specific fields
	result["allPermissions"] = p.AllPermissions
	result["services"] = p.Services
	result["actions"] = p.Actions
	result["matrixData"] = p.MatrixData
	result["title"] = p.Title
	result["subtitle"] = p.Subtitle

	return result
}

// PaginatedData represents the paginated data structure
type PaginatedData struct {
	Data        interface{} `json:"data"`
	Total       int64       `json:"total"`
	CurrentPage int         `json:"currentPage"`
	LastPage    int         `json:"lastPage"`
	PerPage     int         `json:"perPage"`
	From        int         `json:"from"`
	To          int         `json:"to"`
	HasNext     bool        `json:"hasNext"`
	HasPrev     bool        `json:"hasPrev"`
}

// FilterState represents the current filter/search/sort state
type FilterState struct {
	Page      int                    `json:"page"`
	PageSize  int                    `json:"pageSize"`
	Search    string                 `json:"search"`
	Sort      string                 `json:"sort"`
	Direction string                 `json:"direction"`
	Filters   map[string]interface{} `json:"filters"`
}

// PermissionsMap represents user permissions for the resource
type PermissionsMap struct {
	CanView       bool `json:"canView"`
	CanCreate     bool `json:"canCreate"`
	CanEdit       bool `json:"canEdit"`
	CanDelete     bool `json:"canDelete"`
	CanManage     bool `json:"canManage"`
	CanExport     bool `json:"canExport"`
	CanBulkUpdate bool `json:"canBulkUpdate"`
	CanBulkDelete bool `json:"canBulkDelete"`
	IsAdmin       bool `json:"isAdmin"`
	IsSuperAdmin  bool `json:"isSuperAdmin"`
	// Additional custom permissions
	Custom map[string]bool `json:"-"`
}

// PageMetadata represents page metadata
type PageMetadata struct {
	Version          string                 `json:"version"`
	Component        string                 `json:"component"`
	ResourceType     string                 `json:"resourceType"`
	Timestamp        string                 `json:"timestamp"`
	PaginationConfig PaginationMetadata     `json:"pagination"`
	Custom           map[string]interface{} `json:"-"` // For custom metadata
}

// PaginationMetadata represents pagination configuration
type PaginationMetadata struct {
	DefaultPageSize int   `json:"defaultPageSize"`
	MaxPageSize     int   `json:"maxPageSize"`
	AllowedSizes    []int `json:"allowedSizes"`
}

// Helper method to convert PageProps to map for Inertia rendering
func (p PageProps) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"data":        p.Data,
		"filters":     p.Filters,
		"permissions": p.Permissions.ToMap(),
		"meta":        p.Meta.ToMap(),
	}

	// Add stats if present
	if p.Stats != nil {
		result["stats"] = p.Stats
	}

	// Add any extra fields
	for k, v := range p.Extra {
		result[k] = v
	}

	return result
}

// Convert PermissionsMap to plain map
func (p PermissionsMap) ToMap() map[string]bool {
	result := map[string]bool{
		"canView":       p.CanView,
		"canCreate":     p.CanCreate,
		"canEdit":       p.CanEdit,
		"canDelete":     p.CanDelete,
		"canManage":     p.CanManage,
		"canExport":     p.CanExport,
		"canBulkUpdate": p.CanBulkUpdate,
		"canBulkDelete": p.CanBulkDelete,
		"isAdmin":       p.IsAdmin,
		"isSuperAdmin":  p.IsSuperAdmin,
	}

	// Add custom permissions
	for k, v := range p.Custom {
		result[k] = v
	}

	return result
}

// Convert PageMetadata to map
func (m PageMetadata) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"version":      m.Version,
		"component":    m.Component,
		"resourceType": m.ResourceType,
		"timestamp":    m.Timestamp,
		"pagination": map[string]interface{}{
			"defaultPageSize": m.PaginationConfig.DefaultPageSize,
			"maxPageSize":     m.PaginationConfig.MaxPageSize,
			"allowedSizes":    m.PaginationConfig.AllowedSizes,
		},
	}

	// Add custom metadata
	for k, v := range m.Custom {
		result[k] = v
	}

	return result
}

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
		if customFilters, err := c.ParseCustomFilters(filtersJSON); err == nil {
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

func (c *BaseCrudController) GetPaginationDefaults() (page int, pageSize int, maxPageSize int) {
	return 1, c.defaultPageSize, c.maxPageSize
}

func (c *BaseCrudController) BuildPaginatedResponse(result *PaginatedResult, request *ListRequest) map[string]interface{} {
	response := NewPaginatedResponse(result, request)
	return response.ToMap()
}

// BuildTypedPaginatedResponse returns a strongly typed paginated response
func (c *BaseCrudController) BuildTypedPaginatedResponse(result *PaginatedResult, request *ListRequest) *PaginatedResponse {
	return NewPaginatedResponse(result, request)
}

// SEARCH CONTRACT IMPLEMENTATION (enforced)

func (c *BaseCrudController) ValidateSearchRequest(ctx http.Context) (*SearchRequest, error) {
	req := &SearchRequest{}

	// Parse search parameters
	req.Query = ctx.Request().Query("q", "")
	req.Page = ctx.Request().QueryInt("page", 1)
	req.PageSize = ctx.Request().QueryInt("pageSize", c.defaultPageSize)
	req.Sort = ctx.Request().Query("sort", "relevance")
	req.Direction = ctx.Request().Query("direction", "DESC")

	// Parse searchIn fields
	searchInStr := ctx.Request().Query("searchIn", "")
	if searchInStr != "" {
		req.SearchIn = strings.Split(searchInStr, ",")
	}

	// Parse boolean flags
	req.Exact = ctx.Request().Query("exact", "false") == "true"
	req.Highlight = ctx.Request().Query("highlight", "true") == "true"

	// Parse filters from query parameters
	req.Filters = make(map[string]interface{})

	// Get all query parameters
	queries := ctx.Request().Queries()

	// List of known non-filter parameters for search
	knownParams := map[string]bool{
		"q":         true,
		"page":      true,
		"pageSize":  true,
		"sort":      true,
		"direction": true,
		"searchIn":  true,
		"exact":     true,
		"highlight": true,
	}

	// Add all other query parameters as filters
	for key, values := range queries {
		if !knownParams[key] && len(values) > 0 {
			// Use the first value if multiple are provided
			req.Filters[key] = values[0]
		}
	}

	// Validate search query
	if strings.TrimSpace(req.Query) == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	// Validate pagination
	if req.Page <= 0 {
		return nil, fmt.Errorf("page must be greater than 0")
	}

	if req.PageSize <= 0 || req.PageSize > c.maxPageSize {
		req.PageSize = c.defaultPageSize
	}

	// Set defaults
	req.SetDefaults()

	return req, nil
}

func (c *BaseCrudController) BuildSearchResponse(result *PaginatedResult, request *SearchRequest) map[string]interface{} {
	// Build typed response
	response := c.BuildTypedSearchResponse(result, request)

	// Convert to map with legacy structure for compatibility
	return map[string]interface{}{
		"data":       response.Results,
		"pagination": response.Pagination,
		"search": map[string]interface{}{
			"query":     request.Query,
			"sort":      request.Sort,
			"direction": request.Direction,
			"searchIn":  request.SearchIn,
			"exact":     request.Exact,
			"highlight": request.Highlight,
			"filters":   request.Filters,
		},
	}
}

// BuildTypedSearchResponse returns a strongly typed search response
func (c *BaseCrudController) BuildTypedSearchResponse(result *PaginatedResult, request *SearchRequest) *SearchResponse {
	pagination := PaginationInfo{
		CurrentPage: result.CurrentPage,
		LastPage:    result.LastPage,
		PerPage:     result.PerPage,
		Total:       result.Total,
		From:        result.From,
		To:          result.To,
		HasNext:     result.HasNext,
		HasPrev:     result.HasPrev,
	}

	filters := FiltersMeta{
		Page:      request.Page,
		PageSize:  request.PageSize,
		Search:    request.Query,
		Sort:      request.Sort,
		Direction: request.Direction,
		Filters:   request.Filters,
	}

	response := NewSearchResponse(result.Data, request.Query, pagination, filters)

	// Add search metadata
	if len(request.SearchIn) > 0 {
		response.Meta.SearchedIn = request.SearchIn
	}
	response.Meta.Highlighted = request.Highlight

	return response
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
	// Return 204 No Content for successful delete operations
	return ctx.Response().NoContent()
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

// DEFAULT SEARCH IMPLEMENTATION
// Controllers can override this method to provide custom search functionality

func (c *BaseCrudController) Search(ctx http.Context) http.Response {
	// This is a default implementation that returns method not implemented
	// Controllers MUST override this to provide actual search functionality
	return c.BadRequestResponse(ctx, "Search not implemented", map[string]interface{}{
		"error": fmt.Sprintf("Search functionality not implemented for %s", c.resourceType),
		"hint":  "Controller must override the Search method",
	})
}

// GetSearchableFields returns empty by default - controllers must override
func (c *BaseCrudController) GetSearchableFields() []string {
	return []string{}
}

// TYPED RESPONSE HELPERS

// PaginatedSuccessResponse returns a typed paginated success response
func (c *BaseCrudController) PaginatedSuccessResponse(ctx http.Context, result *PaginatedResult, request *ListRequest, message string) http.Response {
	paginatedResponse := NewPaginatedResponse(result, request)
	response := ResponseFormat{
		Success: true,
		Data:    paginatedResponse,
		Message: message,
	}
	return ctx.Response().Json(http.StatusOK, response)
}

// SearchSuccessResponse returns a typed search success response
func (c *BaseCrudController) SearchSuccessResponse(ctx http.Context, result *PaginatedResult, request *SearchRequest, message string) http.Response {
	searchResponse := c.BuildTypedSearchResponse(result, request)
	response := ResponseFormat{
		Success: true,
		Data:    searchResponse,
		Message: message,
	}
	return ctx.Response().Json(http.StatusOK, response)
}

// SingleResourceResponse returns a typed single resource response
func (c *BaseCrudController) SingleResourceResponse(ctx http.Context, resource interface{}, message string) http.Response {
	singleResponse := &SingleResourceResponse{
		Data: resource,
	}
	response := ResponseFormat{
		Success: true,
		Data:    singleResponse,
		Message: message,
	}
	return ctx.Response().Json(http.StatusOK, response)
}

// BulkOperationResponse returns a typed bulk operation response
func (c *BaseCrudController) BulkOperationResponse(ctx http.Context, bulkResult *BulkOperationResponse, message string) http.Response {
	response := ResponseFormat{
		Success: bulkResult.IsCompleteSuccess(),
		Data:    bulkResult,
		Message: message,
	}

	// Determine appropriate status code
	statusCode := http.StatusOK
	if bulkResult.FailedCount > 0 && bulkResult.SuccessCount == 0 {
		statusCode = http.StatusBadRequest
	} else if bulkResult.FailedCount > 0 {
		statusCode = http.StatusPartialContent
	}

	return ctx.Response().Json(statusCode, response)
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

// GetProps builds typed page props for Inertia.js
func (c *BasePageController) GetProps(result *PaginatedResult, request *ListRequest, permissions PermissionsMap, stats interface{}) PageProps {
	return PageProps{
		Data: PaginatedData{
			Data:        result.Data,
			Total:       result.Total,
			CurrentPage: result.CurrentPage,
			LastPage:    result.LastPage,
			PerPage:     result.PerPage,
			From:        result.From,
			To:          result.To,
			HasNext:     result.HasNext,
			HasPrev:     result.HasPrev,
		},
		Filters: FilterState{
			Page:      request.Page,
			PageSize:  request.PageSize,
			Search:    request.Search,
			Sort:      request.Sort,
			Direction: request.Direction,
			Filters:   request.Filters,
		},
		Permissions: permissions,
		Meta: PageMetadata{
			Version:      "1.0.0",
			Component:    c.pageComponent,
			ResourceType: c.resourceType,
			Timestamp:    "1234567890", // Could use time.Now().Unix()
			PaginationConfig: PaginationMetadata{
				DefaultPageSize: c.defaultPageSize,
				MaxPageSize:     c.maxPageSize,
				AllowedSizes:    c.allowedPageSizes,
			},
		},
		Stats: stats,
	}
}

// BuildPageProps is deprecated - use GetProps instead
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
		"version":      "1.0.0",
		"component":    c.pageComponent,
		"resourceType": c.resourceType,
		"timestamp":    "1234567890", // Could use time.Now().Unix()
		"pagination":   c.GetPaginationConfig(),
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

// BuildTypedPermissions converts a permission map to typed PermissionsMap
func (c *BasePageController) BuildTypedPermissions(perms map[string]bool) PermissionsMap {
	result := PermissionsMap{
		Custom: make(map[string]bool),
	}

	// Map known permissions
	for key, value := range perms {
		switch key {
		case "canView":
			result.CanView = value
		case "canCreate":
			result.CanCreate = value
		case "canEdit":
			result.CanEdit = value
		case "canDelete":
			result.CanDelete = value
		case "canManage":
			result.CanManage = value
		case "canExport":
			result.CanExport = value
		case "canBulkUpdate":
			result.CanBulkUpdate = value
		case "canBulkDelete":
			result.CanBulkDelete = value
		case "isAdmin":
			result.IsAdmin = value
		case "isSuperAdmin":
			result.IsSuperAdmin = value
		default:
			// Store unknown permissions in custom map
			result.Custom[key] = value
		}
	}

	return result
}

// GetPermissionMatrixProps builds typed props for permission matrix pages
func (c *BasePageController) GetPermissionMatrixProps(data interface{}, filters map[string]interface{},
	permissions PermissionsMap, stats interface{}, allPermissions []interface{},
	services []interface{}, actions []interface{}, matrixData map[string]interface{},
	title, subtitle string) PermissionMatrixPageProps {

	// Convert data to PaginatedData
	var paginatedData PaginatedData
	if dataMap, ok := data.(map[string]interface{}); ok {
		if items, hasData := dataMap["data"]; hasData {
			paginatedData.Data = items
		}
		if total, hasTotal := dataMap["total"].(int); hasTotal {
			paginatedData.Total = int64(total)
		}
		if currentPage, hasPage := dataMap["currentPage"].(int); hasPage {
			paginatedData.CurrentPage = currentPage
		}
		if lastPage, hasLastPage := dataMap["lastPage"].(int); hasLastPage {
			paginatedData.LastPage = lastPage
		}
		if perPage, hasPerPage := dataMap["perPage"].(int); hasPerPage {
			paginatedData.PerPage = perPage
		}
		if from, hasFrom := dataMap["from"].(int); hasFrom {
			paginatedData.From = from
		}
		if to, hasTo := dataMap["to"].(int); hasTo {
			paginatedData.To = to
		}
	}

	// Build filter state
	filterState := FilterState{
		Filters: filters,
	}

	return PermissionMatrixPageProps{
		PageProps: PageProps{
			Data:        paginatedData,
			Filters:     filterState,
			Permissions: permissions,
			Meta: PageMetadata{
				Version:      "1.0.0",
				Component:    c.pageComponent,
				ResourceType: c.resourceType,
				Timestamp:    "1234567890",
				PaginationConfig: PaginationMetadata{
					DefaultPageSize: c.defaultPageSize,
					MaxPageSize:     c.maxPageSize,
					AllowedSizes:    c.allowedPageSizes,
				},
			},
			Stats: stats,
		},
		AllPermissions: allPermissions,
		Services:       services,
		Actions:        actions,
		MatrixData:     matrixData,
		Title:          title,
		Subtitle:       subtitle,
	}
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

// ParseCustomFilters parses custom filter JSON string into filter conditions
func (c *BaseCrudController) ParseCustomFilters(filtersJSON string) (interface{}, error) {
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
