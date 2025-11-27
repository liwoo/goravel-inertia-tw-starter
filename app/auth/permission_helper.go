package auth

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/models"
)

// PermissionHelper provides permission checking utilities
type PermissionHelper struct {
	permissionService *PermissionService
	cache             *RequestScopedCache
}

// NewPermissionHelper creates a new permission helper
func NewPermissionHelper() *PermissionHelper {
	return &PermissionHelper{
		permissionService: GetPermissionService(),
		cache:             GetRequestScopedCache(),
	}
}

// GetAuthenticatedUser gets the current authenticated user with roles
// Uses request-scoped caching to prevent N+1 queries
func (h *PermissionHelper) GetAuthenticatedUser(ctx http.Context) *models.User {
	// Use the request-scoped cache to avoid repeated DB queries
	return h.cache.GetCachedUser(ctx)
}

// RequireAuthentication ensures user is authenticated
func (h *PermissionHelper) RequireAuthentication(ctx http.Context) (*models.User, error) {
	user := h.GetAuthenticatedUser(ctx)
	if user == nil {
		return nil, fmt.Errorf("authentication required")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	return user, nil
}

// RequirePermission ensures user has specific permission
func (h *PermissionHelper) RequirePermission(ctx http.Context, permission string) (*models.User, error) {
	user, err := h.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	if !h.permissionService.HasPermission(user, permission) {
		return nil, fmt.Errorf("insufficient permissions: %s required", permission)
	}

	return user, nil
}

// RequireRole ensures user has specific role
func (h *PermissionHelper) RequireRole(ctx http.Context, role string) (*models.User, error) {
	user, err := h.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	if !h.permissionService.HasRole(user, role) {
		return nil, fmt.Errorf("insufficient role: %s required", role)
	}

	return user, nil
}

// RequireResourceAccess ensures user can access specific resource
func (h *PermissionHelper) RequireResourceAccess(ctx http.Context, action string, resourceType string, resourceID uint) (*models.User, error) {
	user, err := h.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	if !h.permissionService.CanAccessResource(user, action, resourceType, resourceID) {
		return nil, fmt.Errorf("insufficient permissions for %s.%s on resource %d", resourceType, action, resourceID)
	}

	return user, nil
}

// CheckPermission checks if user has permission (returns bool, no error)
func (h *PermissionHelper) CheckPermission(ctx http.Context, permission string) bool {
	user := h.GetAuthenticatedUser(ctx)
	if user == nil {
		return false
	}

	return h.permissionService.HasPermission(user, permission)
}

// CheckRole checks if user has role (returns bool, no error)
func (h *PermissionHelper) CheckRole(ctx http.Context, role string) bool {
	user := h.GetAuthenticatedUser(ctx)
	if user == nil {
		return false
	}

	return h.permissionService.HasRole(user, role)
}

// CheckResourceAccess checks if user can access resource (returns bool, no error)
func (h *PermissionHelper) CheckResourceAccess(ctx http.Context, action string, resourceType string, resourceID uint) bool {
	user := h.GetAuthenticatedUser(ctx)
	if user == nil {
		return false
	}

	return h.permissionService.CanAccessResource(user, action, resourceType, resourceID)
}

// BuildPermissionsMap builds a permission map for frontend using the new service_action format
func (h *PermissionHelper) BuildPermissionsMap(ctx http.Context, resourceType string) map[string]bool {
	user := h.GetAuthenticatedUser(ctx)
	if user == nil {
		return map[string]bool{
			"canView":   false,
			"canCreate": false,
			"canEdit":   false,
			"canDelete": false,
			"canManage": false,
		}
	}

	// Helper function to check permission with scoped variants
	hasPermissionWithScope := func(basePermission string) bool {
		// Check base permission and all scoped variants
		scopedPermissions := []string{
			basePermission,
			basePermission + "_by_all",
			basePermission + "_by_my_role",
			basePermission + "_by_me",
		}

		for _, perm := range scopedPermissions {
			if h.permissionService.HasPermission(user, perm) {
				return true
			}
		}
		return false
	}

	// Use the new service_action format
	viewSlug := BuildPermissionSlug(ServiceRegistry(resourceType), PermissionView)
	readSlug := BuildPermissionSlug(ServiceRegistry(resourceType), PermissionRead)
	createSlug := BuildPermissionSlug(ServiceRegistry(resourceType), PermissionCreate)
	updateSlug := BuildPermissionSlug(ServiceRegistry(resourceType), PermissionUpdate)
	deleteSlug := BuildPermissionSlug(ServiceRegistry(resourceType), PermissionDelete)

	perms := map[string]bool{
		// Use 'view' permission for listing/viewing, 'read' for accessing individual items
		"canView":   hasPermissionWithScope(viewSlug) || hasPermissionWithScope(readSlug),
		"canCreate": hasPermissionWithScope(createSlug),
		"canEdit":   hasPermissionWithScope(updateSlug),
		"canDelete": hasPermissionWithScope(deleteSlug),
		"canManage": hasPermissionWithScope(BuildPermissionSlug(ServiceRegistry(resourceType), PermissionManage)),

		// Additional permissions
		"canExport":     hasPermissionWithScope(BuildPermissionSlug(ServiceRegistry(resourceType), PermissionExport)),
		"canBulkUpdate": hasPermissionWithScope(BuildPermissionSlug(ServiceRegistry(resourceType), PermissionBulkUpdate)),
		"canBulkDelete": hasPermissionWithScope(BuildPermissionSlug(ServiceRegistry(resourceType), PermissionBulkDelete)),

		// Admin permissions (legacy)
		"isAdmin":      user.IsAdmin(),
		"isSuperAdmin": user.IsSuperAdminUser(),
	}

	return perms
}

// CheckServicePermission checks if user has permission for a specific service and action
func (h *PermissionHelper) CheckServicePermission(ctx http.Context, service ServiceRegistry, action CorePermissionAction) bool {
	user := h.GetAuthenticatedUser(ctx)
	if user == nil {
		return false
	}

	permissionSlug := BuildPermissionSlug(service, action)
	return h.permissionService.HasPermission(user, permissionSlug)
}

// RequireServicePermission ensures user has permission for a specific service and action
func (h *PermissionHelper) RequireServicePermission(ctx http.Context, service ServiceRegistry, action CorePermissionAction) (*models.User, error) {
	user, err := h.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	// Check for any scoped variant of the permission
	basePermission := BuildPermissionSlug(service, action)
	scopedPermissions := []string{
		basePermission,
		basePermission + "_by_all",
		basePermission + "_by_my_role",
		basePermission + "_by_me",
	}

	hasPermission := false
	for _, perm := range scopedPermissions {
		if h.permissionService.HasPermission(user, perm) {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return nil, fmt.Errorf("insufficient permissions: %s required", basePermission)
	}

	return user, nil
}

// GetUserRoles returns user roles as simple string slice
func (h *PermissionHelper) GetUserRoles(ctx http.Context) []string {
	user := h.GetAuthenticatedUser(ctx)
	if user == nil {
		return []string{}
	}

	roles := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		if role.IsActive {
			roles = append(roles, role.Slug)
		}
	}

	return roles
}

// GetUserPermissions returns user permissions as simple string slice
func (h *PermissionHelper) GetUserPermissions(ctx http.Context) []string {
	user := h.GetAuthenticatedUser(ctx)
	if user == nil {
		return []string{}
	}

	return h.permissionService.GetUserPermissions(user)
}

// CanManageUser checks if current user can manage another user
func (h *PermissionHelper) CanManageUser(ctx http.Context, targetUserID uint) bool {
	user := h.GetAuthenticatedUser(ctx)
	if user == nil {
		return false
	}

	// Load target user
	var targetUser models.User
	err := facades.Orm().Query().
		Where("id = ?", targetUserID).
		First(&targetUser)

	if err != nil {
		return false
	}

	return h.permissionService.CanManageUser(user, &targetUser)
}

// Global helper instance
var globalPermissionHelper *PermissionHelper

// GetPermissionHelper returns the global permission helper instance
func GetPermissionHelper() *PermissionHelper {
	if globalPermissionHelper == nil {
		globalPermissionHelper = NewPermissionHelper()
	}
	return globalPermissionHelper
}
