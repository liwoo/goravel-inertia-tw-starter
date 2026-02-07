package auth

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"books-database/app/models"
	"github.com/goravel/framework/facades"
)

// Cache TTL for permissions
const permissionCacheTTL = 15 * time.Minute

// PermissionService handles role-based access control
type PermissionService struct {
	// In-memory cache for performance (fallback when Redis unavailable)
	permissionCache map[string][]string
	roleCache       map[string]*models.Role
	cacheMutex      sync.RWMutex
	cacheExpiry     time.Duration
	lastCacheUpdate time.Time
}

// NewPermissionService creates a new permission service
func NewPermissionService() *PermissionService {
	service := &PermissionService{
		permissionCache: make(map[string][]string),
		roleCache:       make(map[string]*models.Role),
		cacheExpiry:     15 * time.Minute, // Cache for 15 minutes
	}

	// Initialize cache
	service.refreshCache()

	return service
}

// HasPermission checks if a user has a specific permission
func (s *PermissionService) HasPermission(user *models.User, permission string) bool {
	if user == nil {
		return false
	}

	// Super admin has all permissions
	if user.IsSuperAdminUser() {
		return true
	}

	// Use cached permissions when available (loaded once per request)
	// Fall back to loading fresh if no cache is available
	permissions := s.loadUserPermissions(user)

	// Check direct permission match
	for _, perm := range permissions {
		if perm == permission {
			return true
		}
	}

	// NOTE: We do NOT automatically grant base permissions when user has scoped versions
	// This ensures proper scope validation happens in CheckScopedPermission
	// For example, having "books_delete_by_me" does NOT grant "books_delete"

	// Check wildcard permissions
	return s.hasWildcardPermission(permissions, permission)
}

// HasRole checks if a user has a specific role
func (s *PermissionService) HasRole(user *models.User, roleSlug string) bool {
	if user == nil {
		return false
	}

	return user.HasRole(roleSlug)
}

// CanAccessResource checks if user can perform action on a specific resource
func (s *PermissionService) CanAccessResource(user *models.User, action string, resourceType string, resourceID uint) bool {
	if user == nil {
		return false
	}

	// Build permission strings to check
	permissions := []string{
		fmt.Sprintf("%s.%s", resourceType, action),   // books.read
		fmt.Sprintf("%s.%s.*", resourceType, action), // books.read.*
		fmt.Sprintf("%s.*", resourceType),            // books.*
		fmt.Sprintf("*.%s", action),                  // *.read
		"*.*",                                        // *.*
	}

	// Check each permission
	for _, perm := range permissions {
		if s.HasPermission(user, perm) {
			// Check ownership if required
			if s.requiresOwnership(user, perm) {
				return s.isResourceOwner(user, resourceType, resourceID)
			}
			return true
		}
	}

	return false
}

// CanManageUser checks if user can manage another user
func (s *PermissionService) CanManageUser(manager *models.User, target *models.User) bool {
	if manager == nil || target == nil {
		return false
	}

	// Super admin can manage anyone
	if manager.IsSuperAdminUser() {
		return true
	}

	// Check user management permission
	if !s.HasPermission(manager, "users.manage") {
		return false
	}

	// Check role hierarchy
	return manager.CanManageUser(target)
}

// AssignRole assigns a role to a user
func (s *PermissionService) AssignRole(user *models.User, roleSlug string, assignedBy *models.User) error {
	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	// Check if assigner has permission
	if assignedBy != nil && !s.HasPermission(assignedBy, "roles.assign") {
		return fmt.Errorf("insufficient permissions to assign roles")
	}

	// Get role
	role, err := s.getRoleBySlug(roleSlug)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	// Check if user already has this role (active or inactive)
	var existingRole models.UserRole
	err = facades.Orm().Query().
		Where("user_id = ? AND role_id = ?", user.ID, role.ID).
		First(&existingRole)

	if err == nil && existingRole.IsActive {
		return fmt.Errorf("user already has role: %s", roleSlug)
	}

	// Check role hierarchy (can't assign higher role than your own)
	if assignedBy != nil && !assignedBy.IsSuperAdminUser() {
		assignerHighest := assignedBy.GetHighestRole()
		if assignerHighest == nil || !assignerHighest.IsHigherThan(role) {
			return fmt.Errorf("cannot assign role higher than your own")
		}
	}

	// Deactivate all existing active roles for this user
	_, err = facades.Orm().Query().Model(&models.UserRole{}).
		Where("user_id = ? AND is_active = ?", user.ID, true).
		Update("is_active", false)
	if err != nil {
		return fmt.Errorf("failed to deactivate existing roles: %w", err)
	}

	// Create or reactivate user-role assignment
	if existingRole.ID != 0 {
		// Reactivate existing role
		existingRole.IsActive = true
		existingRole.AssignedAt = time.Now()
		if assignedBy != nil {
			existingRole.AssignedByID = &assignedBy.ID
		}
		err = facades.Orm().Query().Save(&existingRole)
	} else {
		// Create new user-role assignment
		userRole := models.UserRole{
			UserID:     user.ID,
			RoleID:     role.ID,
			AssignedAt: time.Now(),
			IsActive:   true,
		}

		if assignedBy != nil {
			userRole.AssignedByID = &assignedBy.ID
		}

		err = facades.Orm().Query().Create(&userRole)
	}
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	// Clear cache
	s.clearUserCache(user.ID)

	return nil
}

// RemoveRole removes a role from a user
func (s *PermissionService) RemoveRole(user *models.User, roleSlug string, removedBy *models.User) error {
	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	// Check permissions
	if removedBy != nil && !s.HasPermission(removedBy, "roles.assign") {
		return fmt.Errorf("insufficient permissions to remove roles")
	}

	// Get role
	role, err := s.getRoleBySlug(roleSlug)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	// Remove user-role assignment
	_, err = facades.Orm().Query().Where("user_id = ? AND role_id = ?", user.ID, role.ID).Delete(&models.UserRole{})
	if err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}

	// Clear cache
	s.clearUserCache(user.ID)

	return nil
}

// GetUserPermissions returns all permissions for a user
func (s *PermissionService) GetUserPermissions(user *models.User) []string {
	if user == nil {
		return []string{}
	}

	return s.loadUserPermissions(user)
}

// CreateRole creates a new role
func (s *PermissionService) CreateRole(name, slug, description string, level int, parentSlug string) (*models.Role, error) {
	role := &models.Role{
		Name:        name,
		Slug:        slug,
		Description: description,
		Level:       level,
		IsActive:    true,
	}

	// Set parent if specified
	if parentSlug != "" {
		parent, err := s.getRoleBySlug(parentSlug)
		if err != nil {
			return nil, fmt.Errorf("parent role not found: %w", err)
		}
		role.ParentID = &parent.ID
	}

	err := facades.Orm().Query().Create(role)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	// Clear cache
	s.refreshCache()

	return role, nil
}

// CreatePermission creates a new permission
func (s *PermissionService) CreatePermission(name, slug, category, action, resource, description string) (*models.Permission, error) {
	permission := &models.Permission{
		Name:        name,
		Slug:        slug,
		Category:    category,
		Action:      action,
		Resource:    resource,
		Description: description,
		IsActive:    true,
	}

	err := facades.Orm().Query().Create(permission)
	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return permission, nil
}

// GrantPermissionToRole grants a permission to a role
func (s *PermissionService) GrantPermissionToRole(roleSlug, permissionSlug string, grantedBy *models.User) error {
	// Get role and permission
	role, err := s.getRoleBySlug(roleSlug)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	permission, err := s.getPermissionBySlug(permissionSlug)
	if err != nil {
		return fmt.Errorf("permission not found: %w", err)
	}

	// Check if already granted
	count, err := facades.Orm().Query().Model(&models.RolePermission{}).Where("role_id = ? AND permission_id = ?", role.ID, permission.ID).Count()
	if err != nil {
		return fmt.Errorf("failed to check existing permission: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("permission already granted to role")
	}

	// Create role-permission assignment
	rolePermission := models.RolePermission{
		RoleID:       role.ID,
		PermissionID: permission.ID,
		GrantedAt:    time.Now(),
		IsActive:     true,
	}

	if grantedBy != nil {
		rolePermission.GrantedByID = &grantedBy.ID
	}

	err = facades.Orm().Query().Create(&rolePermission)
	if err != nil {
		return fmt.Errorf("failed to grant permission: %w", err)
	}

	// Clear cache
	s.refreshCache()

	return nil
}

// Private helper methods

func (s *PermissionService) loadUserPermissions(user *models.User) []string {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("user:%d:permissions", user.ID)
	if cachedJSON := facades.Cache().GetString(cacheKey, ""); cachedJSON != "" {
		var cachedPermissions []string
		if err := json.Unmarshal([]byte(cachedJSON), &cachedPermissions); err == nil {
			return cachedPermissions
		}
	}

	var permissions []string

	// If user already has roles loaded, use them directly
	// This avoids an extra query if GetAuthenticatedUser already loaded roles
	userRoles := user.Roles
	if len(userRoles) == 0 {
		// Load user with roles if not already loaded
		var userWithRoles models.User
		err := facades.Orm().Query().
			Where("id = ?", user.ID).
			With("Roles").
			First(&userWithRoles)

		if err != nil {
			return permissions
		}
		userRoles = userWithRoles.Roles
	}

	// Collect all active role IDs
	roleIDs := make([]uint, 0, len(userRoles))
	for _, role := range userRoles {
		if role.IsActive {
			roleIDs = append(roleIDs, role.ID)
		}
	}

	if len(roleIDs) == 0 {
		return permissions
	}

	// OPTIMIZED: Use query builder to get permission slugs directly
	// This avoids loading 8000+ full permission objects into memory
	// Previous approach used With("Permission") which caused OOM with large permission sets
	type PermissionSlugResult struct {
		Slug  string `gorm:"column:slug"`
		Scope string `gorm:"column:scope"`
	}

	// Convert roleIDs to interface slice for WhereIn
	roleIDsInterface := make([]interface{}, len(roleIDs))
	for i, id := range roleIDs {
		roleIDsInterface[i] = id
	}

	var slugResults []PermissionSlugResult
	err := facades.Orm().Query().
		Table("role_permissions as rp").
		Select("p.slug, rp.scope").
		Join("INNER JOIN permissions p ON p.id = rp.permission_id").
		WhereIn("rp.role_id", roleIDsInterface).
		Where("rp.is_active = ?", true).
		Where("p.is_active = ?", true).
		Where("rp.deleted_at IS NULL").
		Where("p.deleted_at IS NULL").
		Get(&slugResults)

	if err != nil {
		facades.Log().Error("Failed to load permission slugs", map[string]interface{}{
			"error":    err.Error(),
			"role_ids": roleIDs,
		})
		return permissions
	}

	// Collect all permissions from all roles
	permissionMap := make(map[string]bool)

	// Build permission slugs with scopes
	for _, result := range slugResults {
		permSlug := result.Slug
		if result.Scope != "" && result.Scope != "by_all" {
			// Include scope in the permission slug
			permSlug = fmt.Sprintf("%s_%s", permSlug, result.Scope)
		} else if result.Scope == "by_all" || result.Scope == "" {
			// For by_all scope, include both the base permission and the explicit scoped version
			permissionMap[permSlug] = true
			permissionMap[fmt.Sprintf("%s_by_all", permSlug)] = true
			continue
		}
		permissionMap[permSlug] = true
	}

	// Convert map to slice
	for permission := range permissionMap {
		permissions = append(permissions, permission)
	}

	// Cache the permissions
	if len(permissions) > 0 {
		if jsonBytes, err := json.Marshal(permissions); err == nil {
			facades.Cache().Put(cacheKey, string(jsonBytes), permissionCacheTTL)
		}
	}

	return permissions
}

func (s *PermissionService) hasWildcardPermission(permissions []string, targetPermission string) bool {
	for _, perm := range permissions {
		if strings.Contains(perm, "*") {
			if s.matchesWildcard(perm, targetPermission) {
				return true
			}
		}
	}

	return false
}

func (s *PermissionService) matchesWildcard(pattern, target string) bool {
	patternParts := strings.Split(pattern, ".")
	targetParts := strings.Split(target, ".")

	if len(patternParts) != len(targetParts) {
		return false
	}

	for i, part := range patternParts {
		if part != "*" && part != targetParts[i] {
			return false
		}
	}

	return true
}

func (s *PermissionService) getRoleBySlug(slug string) (*models.Role, error) {
	s.cacheMutex.RLock()
	role, exists := s.roleCache[slug]
	s.cacheMutex.RUnlock()

	if exists && !s.isCacheExpired() {
		return role, nil
	}

	var dbRole models.Role
	err := facades.Orm().Query().Where("slug = ? AND is_active = ?", slug, true).First(&dbRole)
	if err != nil {
		return nil, err
	}

	s.cacheMutex.Lock()
	s.roleCache[slug] = &dbRole
	s.cacheMutex.Unlock()

	return &dbRole, nil
}

func (s *PermissionService) getPermissionBySlug(slug string) (*models.Permission, error) {
	var permission models.Permission
	err := facades.Orm().Query().Where("slug = ? AND is_active = ?", slug, true).First(&permission)
	return &permission, err
}

func (s *PermissionService) requiresOwnership(user *models.User, permission string) bool {
	// Check if permission requires ownership
	var perm models.Permission
	err := facades.Orm().Query().Where("slug = ?", permission).First(&perm)
	if err != nil {
		return false
	}

	return perm.RequiresOwnership
}

func (s *PermissionService) isResourceOwner(user *models.User, resourceType string, resourceID uint) bool {
	// Implementation depends on your resource ownership logic
	// For books, check if user is the one who added it, or if it's assigned to them
	switch resourceType {
	case "books":
		var book models.Book
		err := facades.Orm().Query().Where("id = ?", resourceID).First(&book)
		if err != nil {
			return false
		}
		// Add ownership logic here if you have created_by fields
		return true // For now, allow access
	default:
		return false
	}
}

func (s *PermissionService) clearUserCache(userID uint) {
	// Clear in-memory cache
	s.cacheMutex.Lock()
	userKey := fmt.Sprintf("user_%d", userID)
	delete(s.permissionCache, userKey)
	s.cacheMutex.Unlock()

	// Clear cache using facades directly
	cacheKey := fmt.Sprintf("user:%d:permissions", userID)
	if !facades.Cache().Forget(cacheKey) {
		facades.Log().Warning("Failed to invalidate permission cache", map[string]interface{}{
			"user_id": userID,
		})
	}
}

func (s *PermissionService) refreshCache() {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	// Clear existing cache
	s.permissionCache = make(map[string][]string)
	s.roleCache = make(map[string]*models.Role)
	s.lastCacheUpdate = time.Now()
}

func (s *PermissionService) isCacheExpired() bool {
	return time.Since(s.lastCacheUpdate) > s.cacheExpiry
}

// Global instance
var globalPermissionService *PermissionService

// GetPermissionService returns the global permission service instance
func GetPermissionService() *PermissionService {
	if globalPermissionService == nil {
		globalPermissionService = NewPermissionService()
	}
	return globalPermissionService
}
