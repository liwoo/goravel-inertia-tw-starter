package auth

import (
	"fmt"
	"starter-project/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

// ScopedPermissionChecker extends PermissionHelper with scope-aware permission checking
type ScopedPermissionChecker struct {
	*PermissionHelper
	permissionService *PermissionService
}

// CheckScopedPermission checks if a user has permission with the appropriate scope
// Note: This returns bool, so 2FA check is not enforced here (use RequireScopedPermission for that)
func (h *ScopedPermissionChecker) CheckScopedPermission(
	ctx http.Context,
	service ServiceRegistry,
	action CorePermissionAction,
	resource interface{}, // The resource being accessed (e.g., *models.Book)
) bool {
	user := h.GetAuthenticatedUser(ctx)
	if user == nil {
		return false
	}

	// Check 2FA requirement - if not met, return false
	if h.Is2FARequired() && !h.Check2FAEnabled(user) {
		return false
	}

	// Super admins bypass all scope checks (but still need 2FA if required)
	if user.IsSuperAdminUser() {
		return true
	}

	// Check all three scopes from most to least restrictive
	scopes := []PermissionScope{ScopeByAll, ScopeByMyRole, ScopeByMe}

	// Debug logging
	facades.Log().Debug("CheckScopedPermission", map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"service": service,
		"action":  action,
	})

	for _, scope := range scopes {
		permissionSlug := GetPermissionSlug(service, action, scope)

		// Check if user has this scoped permission
		hasPermission := h.permissionService.HasPermission(user, permissionSlug)
		facades.Log().Debug("Checking permission", map[string]interface{}{
			"permission":     permissionSlug,
			"has_permission": hasPermission,
			"scope":          scope,
		})

		if hasPermission {
			// Now verify the scope requirements
			validated := h.validateScope(user, scope, resource)
			facades.Log().Debug("Scope validation", map[string]interface{}{
				"scope":     scope,
				"validated": validated,
			})
			if validated {
				return true
			}
		}

		// Also check the old-style permission format for backward compatibility
		oldStyleSlug := string(service) + "_" + string(action)
		if h.permissionService.HasPermission(user, oldStyleSlug) {
			// Old style permissions are treated as "by_all" scope
			if scope == ScopeByAll {
				return true
			}
		}
	}

	return false
}

// validateScope checks if the user meets the requirements for a specific scope
func (h *ScopedPermissionChecker) validateScope(user *models.User, scope PermissionScope, resource interface{}) bool {
	// If no specific resource is provided (e.g., listing resources),
	// the scope is valid - actual filtering will be done at the service level
	if resource == nil {
		return true
	}

	switch scope {
	case ScopeByAll:
		// No additional validation needed - user can access all resources
		return true

	case ScopeByMe:
		// User can only access resources they created
		return h.isResourceCreatedBy(resource, user.ID)

	case ScopeByMyRole:
		// User can access resources created by users with the same role or lower level
		return h.isResourceCreatedByRoleLevel(resource, user)

	default:
		return false
	}
}

// isResourceCreatedBy checks if a resource was created by a specific user
func (h *ScopedPermissionChecker) isResourceCreatedBy(resource interface{}, userID uint) bool {
	// Use reflection to check if the resource has a CreatedBy field
	switch r := resource.(type) {
	case *models.Book:
		return r.CreatedBy != nil && *r.CreatedBy == userID
	case *models.Lender:
		return r.CreatedBy != nil && *r.CreatedBy == userID
	case *models.User:
		return r.CreatedBy != nil && *r.CreatedBy == userID
	case *models.Role:
		return r.CreatedBy != nil && *r.CreatedBy == userID
	case *models.Permission:
		return r.CreatedBy != nil && *r.CreatedBy == userID
	default:
		// If we can't determine ownership, deny access for safety
		return false
	}
}

// isResourceCreatedByRoleLevel checks if a resource was created by someone at the same role level or lower
func (h *ScopedPermissionChecker) isResourceCreatedByRoleLevel(resource interface{}, user *models.User) bool {
	var creatorID *uint

	// Extract creator ID from resource
	switch r := resource.(type) {
	case *models.Book:
		creatorID = r.CreatedBy
	case *models.Lender:
		creatorID = r.CreatedBy
	case *models.User:
		creatorID = r.CreatedBy
	case *models.Role:
		creatorID = r.CreatedBy
	case *models.Permission:
		creatorID = r.CreatedBy
	default:
		return false
	}

	if creatorID == nil {
		// If no creator is set, we can't validate role level
		return false
	}

	// Get the creator's user record
	var creator models.User
	err := facades.Orm().Query().Where("id = ?", *creatorID).With("Roles").First(&creator)
	if err != nil {
		return false
	}

	// Compare role levels
	userMaxLevel := h.getUserMaxRoleLevel(user)
	creatorMaxLevel := h.getUserMaxRoleLevel(&creator)

	// User can access if their role level is >= creator's role level
	return userMaxLevel >= creatorMaxLevel
}

// getUserMaxRoleLevel gets the highest role level for a user
func (h *ScopedPermissionChecker) getUserMaxRoleLevel(user *models.User) int {
	maxLevel := 0
	for _, role := range user.Roles {
		if role.IsActive && role.Level > maxLevel {
			maxLevel = role.Level
		}
	}
	return maxLevel
}

// RequireScopedPermission checks permission and returns error if denied
// Also enforces 2FA if enabled in config
func (h *ScopedPermissionChecker) RequireScopedPermission(
	ctx http.Context,
	service ServiceRegistry,
	action CorePermissionAction,
	resource interface{},
) (*models.User, error) {
	user := h.GetAuthenticatedUser(ctx)
	if user == nil {
		return nil, fmt.Errorf("authentication required")
	}

	// Check 2FA requirement for permission-protected resources
	// Even super admins must have 2FA enabled when required
	if h.Is2FARequired() && !h.Check2FAEnabled(user) {
		return nil, NewTwoFactorRequiredError()
	}

	if !h.CheckScopedPermission(ctx, service, action, resource) {
		return user, fmt.Errorf("insufficient permissions for %s.%s", service, action)
	}

	return user, nil
}

// GetScopedPermissionHelper returns a singleton instance of ScopedPermissionChecker
var scopedPermissionHelper *ScopedPermissionChecker

func GetScopedPermissionHelper() *ScopedPermissionChecker {
	if scopedPermissionHelper == nil {
		scopedPermissionHelper = &ScopedPermissionChecker{
			PermissionHelper:  GetPermissionHelper(),
			permissionService: NewPermissionService(),
		}
	}
	return scopedPermissionHelper
}
