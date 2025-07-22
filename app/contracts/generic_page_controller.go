package contracts

import (
	"fmt"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"
	"players/app/auth"
	"players/app/http/inertia"
)

// GenericPageController provides a comprehensive base for Inertia.js page controllers
// with built-in contract implementations and common patterns
type GenericPageController struct {
	*BasePageController
	
	// Service integration
	service          CrudServiceContract
	serviceIdentifier auth.ServiceRegistry
	authHelper       AuthHelper
	
	// Configuration
	requireSuperAdmin bool
	customPermissionCheck func(ctx http.Context) error
	
	// Statistics configuration
	statsEnabled     bool
	statsBuilder     func(controller *GenericPageController) map[string]interface{}
	
	// Additional data providers
	extraDataProviders map[string]func(ctx http.Context) (interface{}, error)
}

// GenericPageConfig holds configuration for creating a generic page controller
type GenericPageConfig struct {
	ResourceType      string
	PageComponent     string
	Service           CrudServiceContract
	ServiceIdentifier auth.ServiceRegistry
	RequireSuperAdmin bool
	CustomPermissionCheck func(ctx http.Context) error
	StatsEnabled      bool
	StatsBuilder      func(controller *GenericPageController) map[string]interface{}
}

// NewGenericPageController creates a new generic page controller with reduced boilerplate
func NewGenericPageController(config GenericPageConfig) *GenericPageController {
	controller := &GenericPageController{
		BasePageController:    NewBasePageController(config.ResourceType, config.PageComponent),
		service:              config.Service,
		serviceIdentifier:    config.ServiceIdentifier,
		authHelper:           nil, // Will be set by the concrete controller
		requireSuperAdmin:    config.RequireSuperAdmin,
		customPermissionCheck: config.CustomPermissionCheck,
		statsEnabled:         config.StatsEnabled,
		statsBuilder:         config.StatsBuilder,
		extraDataProviders:   make(map[string]func(ctx http.Context) (interface{}, error)),
	}
	
	return controller
}

// Index provides a complete implementation of the page index method
func (c *GenericPageController) Index(ctx http.Context) http.Response {
	// 1. Permission/Access Check
	if err := c.performPermissionCheck(ctx); err != nil {
		return c.renderForbidden(ctx, err)
	}
	
	// 2. Validate Request
	req, err := c.ValidatePageRequest(ctx)
	if err != nil {
		req = &ListRequest{Page: 1, PageSize: c.defaultPageSize}
		req.SetDefaults()
	}
	
	// 3. Build Permissions Map
	permissions := c.BuildPermissionsMap(ctx, string(c.serviceIdentifier))
	
	// 4. Fetch Data
	result, err := c.service.GetList(*req)
	if err != nil {
		result = c.emptyResult(req)
	}
	
	// 5. Get Statistics (if enabled and permitted)
	var stats map[string]interface{}
	if c.statsEnabled && (permissions["canViewReports"] || permissions["canManage"]) {
		stats = c.getStatistics()
	}
	
	// 6. Build Typed Permissions
	typedPermissions := c.BuildTypedPermissions(permissions)
	
	// 7. Build Props
	props := c.GetProps(result, req, typedPermissions, stats)
	
	// 8. Add Extra Data
	propsMap := props.ToMap()
	if err := c.addExtraData(ctx, propsMap); err != nil {
		// Log error but continue
		fmt.Printf("Error adding extra data: %v\n", err)
	}
	
	// 9. Render
	return inertia.Render(ctx, c.pageComponent, propsMap)
}

// performPermissionCheck handles the permission checking logic
func (c *GenericPageController) performPermissionCheck(ctx http.Context) error {
	// Custom permission check takes precedence
	if c.customPermissionCheck != nil {
		return c.customPermissionCheck(ctx)
	}
	
	// Super admin check
	if c.requireSuperAdmin {
		return c.checkSuperAdmin(ctx)
	}
	
	// Service-based permission check
	if c.serviceIdentifier != "" {
		permHelper := auth.GetPermissionHelper()
		_, err := permHelper.RequireServicePermission(ctx, c.serviceIdentifier, auth.PermissionView)
		return err
	}
	
	// Default: require authentication only
	return c.RequireAuthentication(ctx)
}

// checkSuperAdmin verifies if the current user is a super admin
func (c *GenericPageController) checkSuperAdmin(ctx http.Context) error {
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil || !user.IsSuperAdmin {
		return fmt.Errorf("super admin access required")
	}
	return nil
}

// renderForbidden renders a forbidden response
func (c *GenericPageController) renderForbidden(ctx http.Context, err error) http.Response {
	// Try Inertia error page first
	if c.requireSuperAdmin {
		return inertia.Render(ctx, "Errors/403", map[string]interface{}{
			"message": "Access denied: Super admin privileges required",
		})
	}
	
	// Default JSON response
	return ctx.Response().Status(403).Json(map[string]interface{}{
		"error":   "Forbidden",
		"message": "You don't have permission to access this page",
	})
}

// emptyResult returns an empty paginated result
func (c *GenericPageController) emptyResult(req *ListRequest) *PaginatedResult {
	return &PaginatedResult{
		Data:        []interface{}{},
		Total:       0,
		CurrentPage: 1,
		LastPage:    1,
		PerPage:     req.PageSize,
	}
}

// getStatistics builds statistics using the configured builder
func (c *GenericPageController) getStatistics() map[string]interface{} {
	if c.statsBuilder != nil {
		return c.statsBuilder(c)
	}
	return make(map[string]interface{})
}

// addExtraData adds any registered extra data providers to the props
func (c *GenericPageController) addExtraData(ctx http.Context, props map[string]interface{}) error {
	for key, provider := range c.extraDataProviders {
		data, err := provider(ctx)
		if err != nil {
			return fmt.Errorf("error getting %s data: %w", key, err)
		}
		props[key] = data
	}
	return nil
}

// AddExtraDataProvider registers a function to provide additional data for the page
func (c *GenericPageController) AddExtraDataProvider(key string, provider func(ctx http.Context) (interface{}, error)) {
	c.extraDataProviders[key] = provider
}

// SetAuthHelper sets the auth helper for the controller
func (c *GenericPageController) SetAuthHelper(authHelper AuthHelper) {
	c.authHelper = authHelper
}

// Statistics Helper Methods

// GetCountByFilter gets count using a filter - useful for statistics
func (c *GenericPageController) GetCountByFilter(filters map[string]interface{}) int {
	req := ListRequest{
		PageSize: 1,
		Filters:  filters,
	}
	
	result, err := c.service.GetListAdvanced(req, filters)
	if err != nil {
		return 0
	}
	
	return int(result.Total)
}

// GetCountByStatus is a convenience method for status-based counts
func (c *GenericPageController) GetCountByStatus(status string) int {
	return c.GetCountByFilter(map[string]interface{}{"status": status})
}

// GetTotalCount returns the total count of all records
func (c *GenericPageController) GetTotalCount() int {
	req := ListRequest{PageSize: 1}
	result, err := c.service.GetList(req)
	if err != nil {
		return 0
	}
	return int(result.Total)
}

// CONTRACT IMPLEMENTATIONS - Built into the base controller

// GetServiceIdentifier returns the service identifier for this controller
func (c *GenericPageController) GetServiceIdentifier() auth.ServiceRegistry {
	return c.serviceIdentifier
}

// CheckPermission implements AuthorizationControllerContract
func (c *GenericPageController) CheckPermission(ctx http.Context, permission string, resource interface{}) error {
	if c.requireSuperAdmin {
		return c.checkSuperAdmin(ctx)
	}
	
	permHelper := auth.GetPermissionHelper()
	_, err := permHelper.RequirePermission(ctx, permission)
	return err
}

// GetCurrentUser implements AuthorizationControllerContract
func (c *GenericPageController) GetCurrentUser(ctx http.Context) interface{} {
	permHelper := auth.GetPermissionHelper()
	return permHelper.GetAuthenticatedUser(ctx)
}

// RequireAuthentication implements AuthorizationControllerContract
func (c *GenericPageController) RequireAuthentication(ctx http.Context) error {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		return fmt.Errorf("authentication required")
	}
	return nil
}

// BuildPermissionsMap implements AuthorizationControllerContract
func (c *GenericPageController) BuildPermissionsMap(ctx http.Context, resourceType string) map[string]bool {
	// Super admin override
	if c.requireSuperAdmin {
		permHelper := auth.GetPermissionHelper()
		user := permHelper.GetAuthenticatedUser(ctx)
		isSuperAdmin := user != nil && user.IsSuperAdmin
		
		return map[string]bool{
			"canView":       isSuperAdmin,
			"canCreate":     isSuperAdmin,
			"canEdit":       isSuperAdmin,
			"canDelete":     isSuperAdmin,
			"canManage":     isSuperAdmin,
			"canExport":     isSuperAdmin,
			"canViewReports": isSuperAdmin,
		}
	}
	
	// Default permission building
	permHelper := auth.GetPermissionHelper()
	return permHelper.BuildPermissionsMap(ctx, resourceType)
}

// Helper method to create standard statistics structure
func (c *GenericPageController) BuildStandardStatistics(statusField string, statuses []string) map[string]interface{} {
	stats := make(map[string]interface{})
	stats["total"] = c.GetTotalCount()
	
	// Count by each status
	for _, status := range statuses {
		key := fmt.Sprintf("%sCount", strings.ToLower(status))
		stats[key] = c.GetCountByFilter(map[string]interface{}{statusField: status})
	}
	
	// Add timestamp
	stats["lastUpdated"] = time.Now().Format(time.RFC3339)
	
	return stats
}