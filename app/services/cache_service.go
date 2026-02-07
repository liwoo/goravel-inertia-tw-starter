package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/goravel/framework/facades"
)

// CacheService provides centralized caching functionality
// Uses Redis when available, falls back to memory cache
type CacheService struct{}

// Cache TTL constants
const (
	CacheTTLShort  = 5 * time.Minute  // For frequently changing data
	CacheTTLMedium = 15 * time.Minute // For moderately stable data
	CacheTTLLong   = 1 * time.Hour    // For rarely changing data
	CacheTTLDay    = 24 * time.Hour   // For static data
)

// Cache key prefixes for namespacing
const (
	CacheKeyUser              = "user:%d"
	CacheKeyUserPermissions   = "user:%d:permissions"
	CacheKeyUserRoles         = "user:%d:roles"
	CacheKeyRole              = "role:%s"
	CacheKeyRoleByID          = "role:id:%d"
	CacheKeyRolePermissions   = "role:%d:permissions"
	CacheKeyNotificationCount = "user:%d:notification_counts"
	CacheKeyMessageUnread     = "user:%d:message_unread"
	CacheKeyAllRoles          = "roles:all"
	CacheKeyAllPermissions    = "permissions:all"
)

// NewCacheService creates a new cache service
func NewCacheService() *CacheService {
	return &CacheService{}
}

// Get retrieves a value from cache
func (s *CacheService) Get(key string, defaultValue interface{}) interface{} {
	return facades.Cache().Get(key, defaultValue)
}

// GetString retrieves a string value from cache
func (s *CacheService) GetString(key string, defaultValue string) string {
	return facades.Cache().GetString(key, defaultValue)
}

// GetInt retrieves an int value from cache
func (s *CacheService) GetInt(key string, defaultValue int) int {
	return facades.Cache().GetInt(key, defaultValue)
}

// GetInt64 retrieves an int64 value from cache
func (s *CacheService) GetInt64(key string, defaultValue int64) int64 {
	return facades.Cache().GetInt64(key, defaultValue)
}

// Put stores a value in cache with TTL
func (s *CacheService) Put(key string, value interface{}, ttl time.Duration) error {
	return facades.Cache().Put(key, value, ttl)
}

// PutJSON stores a JSON-serialized value in cache
func (s *CacheService) PutJSON(key string, value interface{}, ttl time.Duration) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value to JSON: %w", err)
	}
	return facades.Cache().Put(key, string(jsonBytes), ttl)
}

// GetJSON retrieves and deserializes a JSON value from cache
func (s *CacheService) GetJSON(key string, dest interface{}) error {
	jsonStr := facades.Cache().GetString(key, "")
	if jsonStr == "" {
		return fmt.Errorf("cache miss for key: %s", key)
	}
	return json.Unmarshal([]byte(jsonStr), dest)
}

// Forever stores a value in cache permanently
func (s *CacheService) Forever(key string, value interface{}) error {
	if !facades.Cache().Forever(key, value) {
		return fmt.Errorf("failed to store value in cache for key: %s", key)
	}
	return nil
}

// Forget removes a value from cache
func (s *CacheService) Forget(key string) error {
	if !facades.Cache().Forget(key) {
		return fmt.Errorf("failed to remove value from cache for key: %s", key)
	}
	return nil
}

// Has checks if a key exists in cache
func (s *CacheService) Has(key string) bool {
	return facades.Cache().Has(key)
}

// Remember gets from cache or executes callback and caches result
func (s *CacheService) Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	return facades.Cache().Remember(key, ttl, callback)
}

// RememberForever gets from cache or executes callback and caches permanently
func (s *CacheService) RememberForever(key string, callback func() (interface{}, error)) (interface{}, error) {
	return facades.Cache().RememberForever(key, callback)
}

// Flush clears all cache
func (s *CacheService) Flush() error {
	if !facades.Cache().Flush() {
		return fmt.Errorf("failed to flush cache")
	}
	return nil
}

// ========================================
// User-related cache operations
// ========================================

// GetUserPermissions retrieves user permissions from cache
func (s *CacheService) GetUserPermissions(userID uint) ([]string, bool) {
	key := fmt.Sprintf(CacheKeyUserPermissions, userID)
	var permissions []string
	err := s.GetJSON(key, &permissions)
	if err != nil {
		return nil, false
	}
	return permissions, true
}

// SetUserPermissions caches user permissions
func (s *CacheService) SetUserPermissions(userID uint, permissions []string) error {
	key := fmt.Sprintf(CacheKeyUserPermissions, userID)
	return s.PutJSON(key, permissions, CacheTTLMedium)
}

// InvalidateUserPermissions removes user permissions from cache
func (s *CacheService) InvalidateUserPermissions(userID uint) error {
	key := fmt.Sprintf(CacheKeyUserPermissions, userID)
	return s.Forget(key)
}

// ========================================
// Notification cache operations
// ========================================

// GetNotificationCounts retrieves notification counts from cache
func (s *CacheService) GetNotificationCounts(userID uint) (map[string]int64, bool) {
	key := fmt.Sprintf(CacheKeyNotificationCount, userID)
	var counts map[string]int64
	err := s.GetJSON(key, &counts)
	if err != nil {
		return nil, false
	}
	return counts, true
}

// SetNotificationCounts caches notification counts
func (s *CacheService) SetNotificationCounts(userID uint, counts map[string]int64) error {
	key := fmt.Sprintf(CacheKeyNotificationCount, userID)
	return s.PutJSON(key, counts, CacheTTLShort)
}

// InvalidateNotificationCounts removes notification counts from cache
func (s *CacheService) InvalidateNotificationCounts(userID uint) error {
	key := fmt.Sprintf(CacheKeyNotificationCount, userID)
	return s.Forget(key)
}

// InvalidateAllNotificationCounts removes notification counts for all users
// This is used when global counts (like pending applications) change
func (s *CacheService) InvalidateAllNotificationCounts() error {
	// Use Flush with pattern matching - Goravel's cache may support this
	// For now, we'll rely on the short TTL to refresh counts
	// A more robust solution would use Redis SCAN to find and delete matching keys
	facades.Log().Debug("Invalidating all notification counts (relies on TTL)")
	return nil
}

// ========================================
// Message cache operations
// ========================================

// GetMessageUnreadCount retrieves unread message count from cache
func (s *CacheService) GetMessageUnreadCount(userID uint) (int64, bool) {
	key := fmt.Sprintf(CacheKeyMessageUnread, userID)
	// Check if exists first
	if !s.Has(key) {
		return 0, false
	}
	count := s.GetInt64(key, -1)
	if count == -1 {
		return 0, false
	}
	return count, true
}

// SetMessageUnreadCount caches unread message count
func (s *CacheService) SetMessageUnreadCount(userID uint, count int64) error {
	key := fmt.Sprintf(CacheKeyMessageUnread, userID)
	return s.Put(key, count, CacheTTLShort)
}

// InvalidateMessageUnreadCount removes unread message count from cache
func (s *CacheService) InvalidateMessageUnreadCount(userID uint) error {
	key := fmt.Sprintf(CacheKeyMessageUnread, userID)
	return s.Forget(key)
}

// ========================================
// Role cache operations
// ========================================

// GetRoleBySlug retrieves a role by slug from cache
func (s *CacheService) GetRoleBySlug(slug string) (map[string]interface{}, bool) {
	key := fmt.Sprintf(CacheKeyRole, slug)
	var role map[string]interface{}
	err := s.GetJSON(key, &role)
	if err != nil {
		return nil, false
	}
	return role, true
}

// SetRoleBySlug caches a role by slug
func (s *CacheService) SetRoleBySlug(slug string, role map[string]interface{}) error {
	key := fmt.Sprintf(CacheKeyRole, slug)
	return s.PutJSON(key, role, CacheTTLLong)
}

// GetRolePermissions retrieves role permissions from cache
func (s *CacheService) GetRolePermissions(roleID uint) ([]string, bool) {
	key := fmt.Sprintf(CacheKeyRolePermissions, roleID)
	var permissions []string
	err := s.GetJSON(key, &permissions)
	if err != nil {
		return nil, false
	}
	return permissions, true
}

// SetRolePermissions caches role permissions
func (s *CacheService) SetRolePermissions(roleID uint, permissions []string) error {
	key := fmt.Sprintf(CacheKeyRolePermissions, roleID)
	return s.PutJSON(key, permissions, CacheTTLLong)
}

// InvalidateRoleCache invalidates all role-related cache entries
func (s *CacheService) InvalidateRoleCache(roleID uint, slug string) error {
	// Invalidate role by ID
	if err := s.Forget(fmt.Sprintf(CacheKeyRoleByID, roleID)); err != nil {
		facades.Log().Warning("Failed to invalidate role cache by ID", map[string]interface{}{
			"role_id": roleID,
			"error":   err.Error(),
		})
	}

	// Invalidate role by slug
	if slug != "" {
		if err := s.Forget(fmt.Sprintf(CacheKeyRole, slug)); err != nil {
			facades.Log().Warning("Failed to invalidate role cache by slug", map[string]interface{}{
				"slug":  slug,
				"error": err.Error(),
			})
		}
	}

	// Invalidate role permissions
	if err := s.Forget(fmt.Sprintf(CacheKeyRolePermissions, roleID)); err != nil {
		facades.Log().Warning("Failed to invalidate role permissions cache", map[string]interface{}{
			"role_id": roleID,
			"error":   err.Error(),
		})
	}

	// Invalidate all roles cache
	if err := s.Forget(CacheKeyAllRoles); err != nil {
		facades.Log().Warning("Failed to invalidate all roles cache", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return nil
}

// ========================================
// Bulk invalidation operations
// ========================================

// InvalidateUserCache invalidates all cache entries for a user
func (s *CacheService) InvalidateUserCache(userID uint) error {
	keys := []string{
		fmt.Sprintf(CacheKeyUser, userID),
		fmt.Sprintf(CacheKeyUserPermissions, userID),
		fmt.Sprintf(CacheKeyUserRoles, userID),
		fmt.Sprintf(CacheKeyNotificationCount, userID),
		fmt.Sprintf(CacheKeyMessageUnread, userID),
	}

	for _, key := range keys {
		if err := s.Forget(key); err != nil {
			facades.Log().Warning("Failed to invalidate cache key", map[string]interface{}{
				"key":   key,
				"error": err.Error(),
			})
		}
	}

	return nil
}

// Global cache service instance
var globalCacheService *CacheService

// GetCacheService returns the global cache service instance
func GetCacheService() *CacheService {
	if globalCacheService == nil {
		globalCacheService = NewCacheService()
	}
	return globalCacheService
}
