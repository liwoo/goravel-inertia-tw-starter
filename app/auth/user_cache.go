package auth

import (
	"fmt"
	"sync"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"starter-project/app/models"
)

const (
	// Context keys for request-scoped caching
	ContextKeyAuthUser         = "auth_user_cached"
	ContextKeyUserPermissions  = "auth_user_permissions_cached"
	ContextKeyPermissionsLoaded = "auth_permissions_loaded"
)

// RequestScopedCache provides per-request caching for user data and permissions
// to eliminate N+1 queries in middleware and permission checks
type RequestScopedCache struct {
	mutex sync.RWMutex
}

// NewRequestScopedCache creates a new request-scoped cache instance
func NewRequestScopedCache() *RequestScopedCache {
	return &RequestScopedCache{}
}

// GetCachedUser retrieves the authenticated user from context cache or loads it once
func (c *RequestScopedCache) GetCachedUser(ctx http.Context) *models.User {
	if ctx == nil {
		return nil
	}

	// Check for test user context first (for integration tests)
	if testUser, ok := ctx.Value("test_user").(*models.User); ok && testUser != nil {
		return c.ensureUserWithRoles(testUser)
	}

	// Check if user is already cached in context
	if cached := ctx.Value(ContextKeyAuthUser); cached != nil {
		if user, ok := cached.(*models.User); ok {
			return user
		}
	}

	// Load user from auth system
	var user models.User
	err := facades.Auth(ctx).User(&user)
	if err != nil || user.ID == 0 {
		return nil
	}

	// Load user with roles (single query with preloading)
	var userWithRoles models.User
	err = facades.Orm().Query().
		Where("id = ?", user.ID).
		With("Roles").
		First(&userWithRoles)

	if err != nil {
		return nil
	}

	// Cache in context for subsequent calls
	ctx.WithValue(ContextKeyAuthUser, &userWithRoles)

	return &userWithRoles
}

// GetCachedPermissions retrieves user permissions from cache or loads them once
func (c *RequestScopedCache) GetCachedPermissions(ctx http.Context, user *models.User) []string {
	if ctx == nil || user == nil {
		return []string{}
	}

	// Check if permissions are already cached
	if cached := ctx.Value(ContextKeyUserPermissions); cached != nil {
		if permissions, ok := cached.([]string); ok {
			return permissions
		}
	}

	// Load permissions once
	permissions := c.loadUserPermissionsOptimized(user)

	// Cache in context
	ctx.WithValue(ContextKeyUserPermissions, permissions)
	ctx.WithValue(ContextKeyPermissionsLoaded, true)

	return permissions
}

// loadUserPermissionsOptimized loads all user permissions in a single optimized query
func (c *RequestScopedCache) loadUserPermissionsOptimized(user *models.User) []string {
	permissionMap := make(map[string]bool)

	// Collect all active role IDs
	roleIDs := make([]uint, 0, len(user.Roles))
	for _, role := range user.Roles {
		if role.IsActive {
			roleIDs = append(roleIDs, role.ID)
		}
	}

	if len(roleIDs) == 0 {
		return []string{}
	}

	// Single query to get all permissions for all roles
	// This replaces the N+1 pattern where we query permissions for each role separately
	// Convert roleIDs to interface slice for WhereIn
	roleIDsInterface := make([]interface{}, len(roleIDs))
	for i, id := range roleIDs {
		roleIDsInterface[i] = id
	}

	var rolePermissions []models.RolePermission
	err := facades.Orm().Query().
		WhereIn("role_id", roleIDsInterface).
		Where("is_active = ?", true).
		With("Permission").
		Find(&rolePermissions)

	if err != nil {
		facades.Log().Error("Failed to load user permissions", map[string]interface{}{
			"error":   err.Error(),
			"user_id": user.ID,
		})
		return []string{}
	}

	// Build permission slugs with scopes
	for _, rp := range rolePermissions {
		if rp.Permission.IsActive {
			permSlug := rp.Permission.Slug
			if rp.Scope != "" && rp.Scope != "by_all" {
				// Include scope in the permission slug
				permSlug = fmt.Sprintf("%s_%s", permSlug, rp.Scope)
			} else if rp.Scope == "by_all" || rp.Scope == "" {
				// For by_all scope, include both the base permission and the explicit scoped version
				permissionMap[permSlug] = true
				permissionMap[fmt.Sprintf("%s_by_all", permSlug)] = true
				continue
			}
			permissionMap[permSlug] = true
		}
	}

	// Convert map to slice
	permissions := make([]string, 0, len(permissionMap))
	for permission := range permissionMap {
		permissions = append(permissions, permission)
	}

	return permissions
}

// ensureUserWithRoles ensures the user has roles loaded
func (c *RequestScopedCache) ensureUserWithRoles(user *models.User) *models.User {
	if len(user.Roles) > 0 {
		return user
	}

	// Load user with roles
	var userWithRoles models.User
	err := facades.Orm().Query().
		Where("id = ?", user.ID).
		With("Roles").
		First(&userWithRoles)

	if err != nil {
		return user
	}

	return &userWithRoles
}

// ClearCache clears the request-scoped cache (useful for testing)
func (c *RequestScopedCache) ClearCache(ctx http.Context) {
	if ctx == nil {
		return
	}
	ctx.WithValue(ContextKeyAuthUser, nil)
	ctx.WithValue(ContextKeyUserPermissions, nil)
	ctx.WithValue(ContextKeyPermissionsLoaded, false)
}

// HasCachedPermission checks if a user has a specific permission using cached data
func (c *RequestScopedCache) HasCachedPermission(ctx http.Context, user *models.User, permission string) bool {
	if user == nil {
		return false
	}

	// Super admin has all permissions
	if user.IsSuperAdminUser() {
		return true
	}

	permissions := c.GetCachedPermissions(ctx, user)

	// Direct match
	for _, perm := range permissions {
		if perm == permission {
			return true
		}
	}

	// Check wildcard permissions
	return hasWildcardPermissionMatch(permissions, permission)
}

// hasWildcardPermissionMatch checks for wildcard permission patterns
func hasWildcardPermissionMatch(permissions []string, target string) bool {
	for _, perm := range permissions {
		if matchesWildcardPattern(perm, target) {
			return true
		}
	}
	return false
}

// matchesWildcardPattern checks if a permission pattern matches a target
func matchesWildcardPattern(pattern, target string) bool {
	// Simple wildcard matching - can be extended for more complex patterns
	if pattern == "*" || pattern == "*.*" {
		return true
	}
	return false
}

// Global singleton instance
var requestScopedCacheInstance *RequestScopedCache
var requestScopedCacheOnce sync.Once

// GetRequestScopedCache returns the global request-scoped cache instance
func GetRequestScopedCache() *RequestScopedCache {
	requestScopedCacheOnce.Do(func() {
		requestScopedCacheInstance = NewRequestScopedCache()
	})
	return requestScopedCacheInstance
}
