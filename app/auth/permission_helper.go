package auth

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/models"
)

// TwoFactorRequiredError is returned when 2FA is required but not enabled
type TwoFactorRequiredError struct {
	Message string
}

func (e *TwoFactorRequiredError) Error() string {
	return e.Message
}

// NewTwoFactorRequiredError creates a new 2FA required error
func NewTwoFactorRequiredError() *TwoFactorRequiredError {
	return &TwoFactorRequiredError{
		Message: "Two-factor authentication is required to access this resource",
	}
}

// IsTwoFactorRequiredError checks if an error is a 2FA required error
func IsTwoFactorRequiredError(err error) bool {
	_, ok := err.(*TwoFactorRequiredError)
	return ok
}

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

// Is2FARequired checks if 2FA is required for permission-protected resources
// This is controlled by the "auth.require_2fa" config setting
func (h *PermissionHelper) Is2FARequired() bool {
	return facades.Config().GetBool("auth.require_2fa", false)
}

// Check2FAEnabled checks if the user has 2FA enabled
func (h *PermissionHelper) Check2FAEnabled(user *models.User) bool {
	return user.TOTPEnabled
}

// Require2FA ensures user has 2FA enabled (if 2FA is required by config)
// Returns TwoFactorRequiredError if 2FA is required but not enabled
func (h *PermissionHelper) Require2FA(ctx http.Context) (*models.User, error) {
	user, err := h.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	// Check if 2FA is required globally
	if h.Is2FARequired() && !h.Check2FAEnabled(user) {
		return nil, NewTwoFactorRequiredError()
	}

	return user, nil
}

// RequirePermission ensures user has specific permission
// Also enforces 2FA if enabled in config
func (h *PermissionHelper) RequirePermission(ctx http.Context, permission string) (*models.User, error) {
	user, err := h.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	// Check 2FA requirement for permission-protected resources
	if h.Is2FARequired() && !h.Check2FAEnabled(user) {
		return nil, NewTwoFactorRequiredError()
	}

	if !h.permissionService.HasPermission(user, permission) {
		return nil, fmt.Errorf("insufficient permissions: %s required", permission)
	}

	return user, nil
}

// RequireRole ensures user has specific role
// Also enforces 2FA if enabled in config
func (h *PermissionHelper) RequireRole(ctx http.Context, role string) (*models.User, error) {
	user, err := h.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	// Check 2FA requirement for role-protected resources
	if h.Is2FARequired() && !h.Check2FAEnabled(user) {
		return nil, NewTwoFactorRequiredError()
	}

	if !h.permissionService.HasRole(user, role) {
		return nil, fmt.Errorf("insufficient role: %s required", role)
	}

	return user, nil
}

// RequireResourceAccess ensures user can access specific resource
// Also enforces 2FA if enabled in config
func (h *PermissionHelper) RequireResourceAccess(ctx http.Context, action string, resourceType string, resourceID uint) (*models.User, error) {
	user, err := h.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	// Check 2FA requirement for resource-protected access
	if h.Is2FARequired() && !h.Check2FAEnabled(user) {
		return nil, NewTwoFactorRequiredError()
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

	// Super admin has all permissions - short circuit
	if user.IsSuperAdminUser() {
		return map[string]bool{
			"canView":       true,
			"canCreate":     true,
			"canEdit":       true,
			"canDelete":     true,
			"canManage":     true,
			"canExport":     true,
			"canBulkUpdate": true,
			"canBulkDelete": true,
			"isAdmin":       true,
			"isSuperAdmin":  true,
		}
	}

	// Get user's permissions once (cached)
	userPermissions := h.permissionService.GetUserPermissions(user)

	// Build a set for O(1) lookup
	permissionSet := make(map[string]bool, len(userPermissions))
	for _, perm := range userPermissions {
		permissionSet[perm] = true
	}

	// Helper function to check permission with scoped variants using the set
	hasPermissionWithScope := func(basePermission string) bool {
		// Check base permission and all scoped variants
		if permissionSet[basePermission] ||
			permissionSet[basePermission+"_by_all"] ||
			permissionSet[basePermission+"_by_my_role"] ||
			permissionSet[basePermission+"_by_me"] {
			return true
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
// Also enforces 2FA if enabled in config
func (h *PermissionHelper) RequireServicePermission(ctx http.Context, service ServiceRegistry, action CorePermissionAction) (*models.User, error) {
	user, err := h.RequireAuthentication(ctx)
	if err != nil {
		return nil, err
	}

	// Check 2FA requirement for permission-protected resources
	// Even super admins must have 2FA enabled when required
	if h.Is2FARequired() && !h.Check2FAEnabled(user) {
		return nil, NewTwoFactorRequiredError()
	}

	// Super admin has all permissions
	if user.IsSuperAdminUser() {
		return user, nil
	}

	// Get user's permissions once (cached) and build a set for O(1) lookup
	userPermissions := h.permissionService.GetUserPermissions(user)
	permissionSet := make(map[string]bool, len(userPermissions))
	for _, perm := range userPermissions {
		permissionSet[perm] = true
	}

	// Check for any scoped variant of the permission using set lookup
	basePermission := BuildPermissionSlug(service, action)
	if permissionSet[basePermission] ||
		permissionSet[basePermission+"_by_all"] ||
		permissionSet[basePermission+"_by_my_role"] ||
		permissionSet[basePermission+"_by_me"] {
		return user, nil
	}

	return nil, fmt.Errorf("insufficient permissions: %s required", basePermission)
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
