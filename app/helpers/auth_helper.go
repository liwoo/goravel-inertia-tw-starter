package helpers

import (
	"books-database/app/auth"
	"books-database/app/contracts"
	"books-database/app/models"
	"context"

	accessImpl "github.com/goravel/framework/auth/access"
	"github.com/goravel/framework/contracts/auth/access"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

// AuthHelper implements authorization helper functions with RBAC support
type AuthHelper struct {
	permissionHelper *auth.PermissionHelper
}

// NewAuthHelper creates a new auth helper
func NewAuthHelper() contracts.AuthHelper {
	return &AuthHelper{
		permissionHelper: auth.GetPermissionHelper(),
	}
}

// GetAuth returns a helper for accessing auth constants and functions
func GetAuth() *authConstantsHelper {
	return &authConstantsHelper{}
}

// authConstantsHelper provides access to auth constants
type authConstantsHelper struct{}

// GetAllServiceRegistries returns all registered services
func (a *authConstantsHelper) GetAllServiceRegistries() []auth.ServiceRegistry {
	return auth.GetAllServiceRegistries()
}

// GetServiceActions returns the valid actions for a specific service
func (a *authConstantsHelper) GetServiceActions(service auth.ServiceRegistry) []auth.CorePermissionAction {
	return auth.GetServiceActions(service)
}

// BuildPermissionSlug creates a permission slug in the format: service_action
func (a *authConstantsHelper) BuildPermissionSlug(service auth.ServiceRegistry, action auth.CorePermissionAction) string {
	return auth.BuildPermissionSlug(service, action)
}

// GetServiceDisplayName returns the human-readable name for a service
func (a *authConstantsHelper) GetServiceDisplayName(service auth.ServiceRegistry) string {
	return auth.GetServiceDisplayName(service)
}

// GetActionDisplayName returns the human-readable name for an action
func (a *authConstantsHelper) GetActionDisplayName(action auth.CorePermissionAction) string {
	return auth.GetActionDisplayName(action)
}

// Permission action constants
func (a *authConstantsHelper) PermissionView() auth.CorePermissionAction { return auth.PermissionView }
func (a *authConstantsHelper) PermissionManage() auth.CorePermissionAction {
	return auth.PermissionManage
}
func (a *authConstantsHelper) PermissionExport() auth.CorePermissionAction {
	return auth.PermissionExport
}

// Role checks
func (h *AuthHelper) HasRole(user interface{}, roles ...string) bool {
	u, ok := user.(*models.User)
	if !ok || u == nil {
		return false
	}

	for _, role := range roles {
		if u.HasRole(role) {
			return true
		}
	}
	return false
}

func (h *AuthHelper) HasAnyRole(user interface{}, roles []string) bool {
	u, ok := user.(*models.User)
	if !ok || u == nil {
		return false
	}

	// Use traditional loop since collect doesn't have Contains method for custom logic
	for _, role := range roles {
		if u.HasRole(role) {
			return true
		}
	}
	return false
}

func (h *AuthHelper) HasAllRoles(user interface{}, roles []string) bool {
	u, ok := user.(*models.User)
	if !ok || u == nil {
		return false
	}

	// Check if user has all required roles
	for _, role := range roles {
		if !u.HasRole(role) {
			return false
		}
	}
	return true
}

// Permission checks
func (h *AuthHelper) HasPermission(user interface{}, permission string) bool {
	u, ok := user.(*models.User)
	if !ok || u == nil {
		return false
	}

	return u.HasPermission(permission)
}

// Ownership checks
func (h *AuthHelper) IsOwner(user interface{}, resource interface{}) bool {
	u, ok := user.(*models.User)
	if !ok || u == nil {
		return false
	}

	// Add ownership logic based on resource type
	switch r := resource.(type) {
	case *models.Book:
		// If books had a created_by field, check it here
		// For now, return true for authenticated users
		return true
	case *models.User:
		// Users own their own profile
		return u.ID == r.ID
	default:
		return false
	}
}

// Context checks
func (h *AuthHelper) IsSuperAdmin(user interface{}) bool {
	u, ok := user.(*models.User)
	if !ok || u == nil {
		return false
	}

	return u.IsSuperAdminUser()
}

func (h *AuthHelper) IsAuthenticated(ctx http.Context) bool {
	return h.GetCurrentUser(ctx) != nil
}

func (h *AuthHelper) GetCurrentUser(ctx http.Context) interface{} {
	return h.permissionHelper.GetAuthenticatedUser(ctx)
}

// Resource-specific checks
func (h *AuthHelper) CanManageResource(user interface{}, resource string) bool {
	u, ok := user.(*models.User)
	if !ok || u == nil {
		return false
	}

	return u.HasPermission(resource + ".manage")
}

func (h *AuthHelper) CanAccessResource(user interface{}, resource interface{}) bool {
	u, ok := user.(*models.User)
	if !ok || u == nil {
		return false
	}

	// Basic access check - can be extended
	return u.IsActive
}

// GateHelper implements gate helper functions
type GateHelper struct {
	authHelper contracts.AuthHelper
}

// NewGateHelper creates a new gate helper
func NewGateHelper() contracts.GateHelper {
	return &GateHelper{
		authHelper: NewAuthHelper(),
	}
}

// RoleBasedAccess creates a role-based access gate
func (g *GateHelper) RoleBasedAccess(allowedRoles ...string) func(ctx http.Context, user interface{}) access.Response {
	return func(ctx http.Context, user interface{}) access.Response {
		if g.authHelper.HasAnyRole(user, allowedRoles) {
			return accessImpl.NewAllowResponse()
		}
		return accessImpl.NewDenyResponse("Insufficient role privileges")
	}
}

// OwnershipBasedAccess creates an ownership-based access gate
func (g *GateHelper) OwnershipBasedAccess() func(ctx http.Context, user interface{}, model interface{}) access.Response {
	return func(ctx http.Context, user interface{}, model interface{}) access.Response {
		if g.authHelper.IsSuperAdmin(user) || g.authHelper.IsOwner(user, model) {
			return accessImpl.NewAllowResponse()
		}
		return accessImpl.NewDenyResponse("You can only access your own resources")
	}
}

// ConditionalAccess creates a conditional access gate
func (g *GateHelper) ConditionalAccess(condition func(ctx http.Context, user interface{}) bool, message string) func(ctx http.Context, user interface{}) access.Response {
	return func(ctx http.Context, user interface{}) access.Response {
		if condition(ctx, user) {
			return accessImpl.NewAllowResponse()
		}
		return accessImpl.NewDenyResponse(message)
	}
}

// RegisterResourceGates registers standard CRUD gates for a resource
func (g *GateHelper) RegisterResourceGates(resource string, config contracts.GateConfig) {
	// Map resource name to ServiceRegistry
	var service auth.ServiceRegistry
	switch resource {
	case "books":
		service = auth.ServiceBooks
	case "users":
		service = auth.ServiceUsers
	case "roles":
		service = auth.ServiceRoles
	case "permissions":
		service = auth.ServicePermissions
	default:
		service = auth.ServiceRegistry(resource)
	}

	// Register viewAny gate
	if config.ViewAnyHandler != nil {
		g.RegisterGate(auth.BuildPermissionSlug(service, auth.PermissionView), config.ViewAnyHandler)
	}

	// Register view gate (for individual resources)
	if config.ViewHandler != nil {
		viewGate := auth.BuildPermissionSlug(service, auth.PermissionRead)
		facades.Gate().Define(viewGate, func(ctx context.Context, args map[string]any) access.Response {
			// TODO: Convert context.Context to http.Context when needed
			// For now, we'll skip the context conversion
			user := args["user"]
			model := args["model"]
			return config.ViewHandler(nil, user, model)
		})
	}

	// Register create gate
	if config.CreateHandler != nil {
		g.RegisterGate(auth.BuildPermissionSlug(service, auth.PermissionCreate), config.CreateHandler)
	}

	// Register update gate
	if config.UpdateHandler != nil {
		updateGate := auth.BuildPermissionSlug(service, auth.PermissionUpdate)
		facades.Gate().Define(updateGate, func(ctx context.Context, args map[string]any) access.Response {
			// TODO: Convert context.Context to http.Context when needed
			user := args["user"]
			model := args["model"]
			return config.UpdateHandler(nil, user, model)
		})
	}

	// Register delete gate
	if config.DeleteHandler != nil {
		deleteGate := auth.BuildPermissionSlug(service, auth.PermissionDelete)
		facades.Gate().Define(deleteGate, func(ctx context.Context, args map[string]any) access.Response {
			// TODO: Convert context.Context to http.Context when needed
			user := args["user"]
			model := args["model"]
			return config.DeleteHandler(nil, user, model)
		})
	}
}

// RegisterGate registers a single gate
func (g *GateHelper) RegisterGate(name string, handler func(ctx http.Context, user interface{}) access.Response) {
	facades.Gate().Define(name, func(ctx context.Context, args map[string]any) access.Response {
		// TODO: Convert context.Context to http.Context when needed
		user := args["user"]
		return handler(nil, user)
	})
}
