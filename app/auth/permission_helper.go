package auth

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"players/app/models"
	)

// PermissionHelper provides permission checking utilities
type PermissionHelper struct {
	permissionService *PermissionService
}

// NewPermissionHelper creates a new permission helper
func NewPermissionHelper() *PermissionHelper {
	return &PermissionHelper{
		permissionService: GetPermissionService(),
	}
}

// GetAuthenticatedUser gets the current authenticated user with roles
func (h *PermissionHelper) GetAuthenticatedUser(ctx http.Context) *models.User {
	var user models.User
	err := facades.Auth(ctx).User(&user)
	if err != nil || user.ID == 0 {
		return nil
	}
	
	// Load user with roles and permissions
	var userWithRoles models.User
	err = facades.Orm().Query().
		Where("id = ?", user.ID).
		With("Roles.Permissions").  // Preload roles AND their permissions
		First(&userWithRoles)
	
	if err != nil {
		return nil
	}
	
	return &userWithRoles
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
	
	// Use the new service_action format
	readSlug := BuildPermissionSlug(ServiceRegistry(resourceType), PermissionRead)
	createSlug := BuildPermissionSlug(ServiceRegistry(resourceType), PermissionCreate)
	updateSlug := BuildPermissionSlug(ServiceRegistry(resourceType), PermissionUpdate)
	deleteSlug := BuildPermissionSlug(ServiceRegistry(resourceType), PermissionDelete)
	
	fmt.Printf("DEBUG BuildPermissionsMap for %s: checking permissions %s, %s, %s, %s\n", 
		resourceType, readSlug, createSlug, updateSlug, deleteSlug)
	
	perms := map[string]bool{
		"canView":   h.permissionService.HasPermission(user, readSlug),
		"canCreate": h.permissionService.HasPermission(user, createSlug),
		"canEdit":   h.permissionService.HasPermission(user, updateSlug),
		"canDelete": h.permissionService.HasPermission(user, deleteSlug),
		"canManage": h.permissionService.HasPermission(user, BuildPermissionSlug(ServiceRegistry(resourceType), PermissionManage)),
		
		// Additional permissions
		"canExport":     h.permissionService.HasPermission(user, BuildPermissionSlug(ServiceRegistry(resourceType), PermissionExport)),
		"canBulkUpdate": h.permissionService.HasPermission(user, BuildPermissionSlug(ServiceRegistry(resourceType), PermissionBulkUpdate)),
		"canBulkDelete": h.permissionService.HasPermission(user, BuildPermissionSlug(ServiceRegistry(resourceType), PermissionBulkDelete)),
		
		// Special report permissions
		"canViewReports": h.permissionService.HasPermission(user, BuildPermissionSlug(ServiceReports, PermissionView)),
		
		// Admin permissions (legacy)
		"isAdmin":      user.IsAdmin(),
		"isSuperAdmin": user.IsSuperAdmin(),
	}
	
	fmt.Printf("DEBUG BuildPermissionsMap result for %s: %+v\n", resourceType, perms)
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
	
	permissionSlug := BuildPermissionSlug(service, action)
	if !h.permissionService.HasPermission(user, permissionSlug) {
		return nil, fmt.Errorf("insufficient permissions: %s required", permissionSlug)
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