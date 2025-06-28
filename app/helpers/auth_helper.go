package helpers

import (
	"players/app/auth"
	"players/app/contracts"
	"players/app/models"

	accessImpl "github.com/goravel/framework/auth/access"
	"github.com/goravel/framework/contracts/auth/access"
	"github.com/goravel/framework/contracts/http"
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
	// TODO: Implement gate registration with proper signatures
	// For now, this is disabled to allow the build to succeed
	// Gates will be handled through middleware and direct authorization checks
}

// RegisterGate registers a single gate
func (g *GateHelper) RegisterGate(name string, handler func(ctx http.Context, user interface{}) access.Response) {
	// TODO: Implement gate registration
	// For now, this is disabled to allow the build to succeed
}