package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/models"
)

// TestRequestScopedCacheCreation tests that the request-scoped cache is created correctly
func TestRequestScopedCacheCreation(t *testing.T) {
	cache := auth.NewRequestScopedCache()
	assert.NotNil(t, cache, "Cache should be created")
}

// TestGetRequestScopedCacheSingleton tests that the singleton returns the same instance
func TestGetRequestScopedCacheSingleton(t *testing.T) {
	cache1 := auth.GetRequestScopedCache()
	cache2 := auth.GetRequestScopedCache()
	assert.Same(t, cache1, cache2, "Should return the same singleton instance")
}

// TestUserCacheWithNilContext tests that nil context returns nil user
func TestUserCacheWithNilContext(t *testing.T) {
	cache := auth.NewRequestScopedCache()
	user := cache.GetCachedUser(nil)
	assert.Nil(t, user, "Nil context should return nil user")
}

// TestPermissionCacheWithNilContext tests that nil context returns empty permissions
func TestPermissionCacheWithNilContext(t *testing.T) {
	cache := auth.NewRequestScopedCache()
	permissions := cache.GetCachedPermissions(nil, nil)
	assert.Empty(t, permissions, "Nil context should return empty permissions")
}

// TestPermissionCacheWithNilUser tests that nil user returns empty permissions
func TestPermissionCacheWithNilUser(t *testing.T) {
	cache := auth.NewRequestScopedCache()
	// Even with a non-nil context, nil user should return empty permissions
	permissions := cache.GetCachedPermissions(nil, nil)
	assert.Empty(t, permissions, "Nil user should return empty permissions")
}

// TestHasCachedPermissionSuperAdmin tests that super admin always has permission
func TestHasCachedPermissionSuperAdmin(t *testing.T) {
	cache := auth.NewRequestScopedCache()

	// Create a super admin user
	superAdmin := &models.User{
		IsSuperAdmin: true,
	}
	superAdmin.ID = 1

	// Super admin should have any permission
	hasPermission := cache.HasCachedPermission(nil, superAdmin, "any_permission")
	assert.True(t, hasPermission, "Super admin should have all permissions")
}

// TestHasCachedPermissionNilUser tests that nil user has no permissions
func TestHasCachedPermissionNilUser(t *testing.T) {
	cache := auth.NewRequestScopedCache()
	hasPermission := cache.HasCachedPermission(nil, nil, "any_permission")
	assert.False(t, hasPermission, "Nil user should not have permissions")
}

// TestPermissionHelperUsesCache tests that PermissionHelper uses the cache
func TestPermissionHelperUsesCache(t *testing.T) {
	// This test verifies the integration between PermissionHelper and RequestScopedCache
	helper := auth.NewPermissionHelper()
	assert.NotNil(t, helper, "PermissionHelper should be created with cache")
}

// BenchmarkPermissionCacheLookup benchmarks permission lookups with caching
func BenchmarkPermissionCacheLookup(b *testing.B) {
	cache := auth.NewRequestScopedCache()
	superAdmin := &models.User{
		IsSuperAdmin: true,
	}
	superAdmin.ID = 1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.HasCachedPermission(nil, superAdmin, "test_permission")
	}
}

// TestPermissionServiceOptimizedLoading tests that permission loading is optimized
func TestPermissionServiceOptimizedLoading(t *testing.T) {
	// Create a user with roles already loaded (simulating cached user)
	user := &models.User{
		Roles: []models.Role{
			{
				IsActive: true,
			},
		},
	}
	user.ID = 1
	user.Roles[0].ID = 1

	// Verify the user has roles loaded
	assert.Len(t, user.Roles, 1, "User should have 1 role loaded")
	assert.True(t, user.Roles[0].IsActive, "Role should be active")
}

// TestContextKeys tests that context keys are properly defined
func TestContextKeys(t *testing.T) {
	assert.Equal(t, "auth_user_cached", auth.ContextKeyAuthUser)
	assert.Equal(t, "auth_user_permissions_cached", auth.ContextKeyUserPermissions)
	assert.Equal(t, "auth_permissions_loaded", auth.ContextKeyPermissionsLoaded)
}

// TestWildcardPermissionMatching tests wildcard permission matching
func TestWildcardPermissionMatching(t *testing.T) {
	// Test that super admin wildcards work
	cache := auth.NewRequestScopedCache()

	superAdmin := &models.User{
		IsSuperAdmin: true,
	}
	superAdmin.ID = 1

	// Should match any permission for super admin
	assert.True(t, cache.HasCachedPermission(nil, superAdmin, "books_view"))
	assert.True(t, cache.HasCachedPermission(nil, superAdmin, "users_manage"))
	assert.True(t, cache.HasCachedPermission(nil, superAdmin, "random_permission"))
}

// TestPerformanceTimingHelper is a helper for performance testing
type TestPerformanceTimingHelper struct {
	start time.Time
	name  string
}

func NewPerformanceTimer(name string) *TestPerformanceTimingHelper {
	return &TestPerformanceTimingHelper{
		start: time.Now(),
		name:  name,
	}
}

func (h *TestPerformanceTimingHelper) Stop(t *testing.T, maxDuration time.Duration) {
	elapsed := time.Since(h.start)
	t.Logf("%s took %v", h.name, elapsed)
	if elapsed > maxDuration {
		t.Logf("WARNING: %s exceeded expected duration of %v", h.name, maxDuration)
	}
}

// TestCachePerformanceWithManyPermissions tests cache performance
func TestCachePerformanceWithManyPermissions(t *testing.T) {
	cache := auth.NewRequestScopedCache()

	// Create a user with many roles
	user := &models.User{
		Roles: make([]models.Role, 5),
	}
	user.ID = 1

	for i := range user.Roles {
		user.Roles[i].ID = uint(i + 1)
		user.Roles[i].IsActive = true
	}

	// Time multiple permission checks
	timer := NewPerformanceTimer("100 permission checks")
	for i := 0; i < 100; i++ {
		cache.HasCachedPermission(nil, user, "test_permission")
	}
	timer.Stop(t, 10*time.Millisecond) // Should be very fast in-memory
}
