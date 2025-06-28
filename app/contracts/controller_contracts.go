package contracts

import (
	"github.com/goravel/framework/contracts/http"
)

// CrudControllerContract defines mandatory CRUD operations for controllers
// This contract FORCES implementation of all CRUD endpoints with proper validation
type CrudControllerContract interface {
	// Core CRUD operations - ALL must be implemented
	Index(ctx http.Context) http.Response  // GET /resource - List with pagination
	Show(ctx http.Context) http.Response   // GET /resource/{id} - Get single resource (JSON for modals)
	Store(ctx http.Context) http.Response  // POST /resource - Create new resource
	Update(ctx http.Context) http.Response // PUT /resource/{id} - Update existing resource
	Delete(ctx http.Context) http.Response // DELETE /resource/{id} - Delete resource

	// Pagination contract - MUST be implemented for Index
	PaginationControllerContract

	// Validation contract - MUST be implemented for Store/Update
	ValidationControllerContract

	// Response contract - MUST be implemented for consistent responses
	ResponseControllerContract
}

// PaginationControllerContract enforces pagination in listing endpoints
type PaginationControllerContract interface {
	// ValidatePaginationRequest ensures pagination parameters are valid
	ValidatePaginationRequest(ctx http.Context) (*ListRequest, error)

	// GetPaginationDefaults returns default pagination settings
	GetPaginationDefaults() (page int, pageSize int, maxPageSize int)

	// BuildPaginatedResponse creates standardized paginated response
	BuildPaginatedResponse(result *PaginatedResult, request *ListRequest) map[string]interface{}
}

// ValidationControllerContract enforces validation for Create/Update operations
type ValidationControllerContract interface {
	// ValidateCreateRequest validates data for creation
	ValidateCreateRequest(ctx http.Context) (map[string]interface{}, error)

	// ValidateUpdateRequest validates data for updates
	ValidateUpdateRequest(ctx http.Context, id uint) (map[string]interface{}, error)

	// ValidateID ensures ID parameter is valid
	ValidateID(ctx http.Context, paramName string) (uint, error)

	// GetValidationRules returns validation rules for this controller
	GetValidationRules() map[string]interface{}
}

// ResponseControllerContract enforces consistent response formatting
type ResponseControllerContract interface {
	// Success responses
	SuccessResponse(ctx http.Context, data interface{}, message string) http.Response
	CreatedResponse(ctx http.Context, data interface{}, message string) http.Response
	NoContentResponse(ctx http.Context, message string) http.Response

	// Error responses
	BadRequestResponse(ctx http.Context, message string, errors map[string]interface{}) http.Response
	NotFoundResponse(ctx http.Context, message string) http.Response
	ForbiddenResponse(ctx http.Context, message string) http.Response
	ValidationErrorResponse(ctx http.Context, errors map[string]interface{}) http.Response
	InternalErrorResponse(ctx http.Context, message string) http.Response

	// Specialized responses for CRUD operations
	ResourceNotFoundResponse(ctx http.Context, resourceType string, id uint) http.Response
	ResourceCreatedResponse(ctx http.Context, resource interface{}, resourceType string) http.Response
	ResourceUpdatedResponse(ctx http.Context, resource interface{}, resourceType string) http.Response
	ResourceDeletedResponse(ctx http.Context, resourceType string, id uint) http.Response
}

// AuthorizationControllerContract enforces authorization checks
type AuthorizationControllerContract interface {
	// CheckPermission validates user has required permission
	CheckPermission(ctx http.Context, permission string, resource interface{}) error

	// GetCurrentUser gets authenticated user from context
	GetCurrentUser(ctx http.Context) interface{}

	// RequireAuthentication ensures user is authenticated
	RequireAuthentication(ctx http.Context) error

	// BuildPermissionsMap creates permission map for frontend
	BuildPermissionsMap(ctx http.Context, resourceType string) map[string]bool
}

// PageControllerContract defines contracts for Inertia.js page controllers
type PageControllerContract interface {
	// Index renders page with data and permissions
	Index(ctx http.Context) http.Response

	// Pagination contract - MUST be implemented
	PaginationControllerContract

	// Authorization contract - MUST be implemented
	AuthorizationControllerContract

	// Page-specific response building
	PageResponseContract
}

// PageResponseContract enforces consistent page response structure
type PageResponseContract interface {
	// BuildPageProps creates standardized props for page components
	BuildPageProps(data interface{}, filters interface{}, permissions map[string]bool, meta map[string]interface{}) map[string]interface{}

	// GetPageMetadata returns metadata for the page (version, features, etc.)
	GetPageMetadata() map[string]interface{}

	// ValidatePageRequest validates request parameters for page rendering
	ValidatePageRequest(ctx http.Context) (*ListRequest, error)
}

// ResourceControllerContract combines all controller contracts for complete resource management
type ResourceControllerContract interface {
	CrudControllerContract
	AuthorizationControllerContract
}

// ControllerMetadata provides information about controller capabilities
type ControllerMetadata struct {
	ResourceType     string   `json:"resource_type"`
	SupportedActions []string `json:"supported_actions"`
	RequiredPerms    []string `json:"required_permissions"`
	ValidationRules  map[string]interface{} `json:"validation_rules"`
	PaginationConfig PaginationConfig `json:"pagination_config"`
	ResponseFormats  []string `json:"response_formats"`
}

// PaginationConfig defines pagination configuration for controllers
type PaginationConfig struct {
	DefaultPageSize int `json:"default_page_size"`
	MaxPageSize     int `json:"max_page_size"`
	AllowedSizes    []int `json:"allowed_sizes"`
}

// ControllerValidationResult represents controller validation result
type ControllerValidationResult struct {
	Valid            bool     `json:"valid"`
	Errors           []string `json:"errors"`
	MissingMethods   []string `json:"missing_methods"`
	InvalidResponses []string `json:"invalid_responses"`
}

// ResponseFormat defines standard API response structure
type ResponseFormat struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// PaginatedResponseFormat extends ResponseFormat for paginated data
type PaginatedResponseFormat struct {
	ResponseFormat
	Pagination *PaginationMeta `json:"pagination"`
}

// PaginationMeta provides pagination metadata
type PaginationMeta struct {
	CurrentPage int  `json:"current_page"`
	LastPage    int  `json:"last_page"`
	PerPage     int  `json:"per_page"`
	Total       int64 `json:"total"`
	From        int  `json:"from"`
	To          int  `json:"to"`
	HasNext     bool `json:"has_next"`
	HasPrev     bool `json:"has_prev"`
}