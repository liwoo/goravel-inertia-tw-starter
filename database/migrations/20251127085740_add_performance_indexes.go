package migrations

import (
	"github.com/goravel/framework/facades"
)

type M20251127085740AddPerformanceIndexes struct{}

// Signature The unique signature for the migration.
func (r *M20251127085740AddPerformanceIndexes) Signature() string {
	return "20251127085740_add_performance_indexes"
}

// createIndexIfNotExists creates an index using raw SQL with IF NOT EXISTS for PostgreSQL compatibility
func (r *M20251127085740AddPerformanceIndexes) createIndexIfNotExists(indexName, tableName, columns string) {
	query := "CREATE INDEX IF NOT EXISTS " + indexName + " ON " + tableName + " (" + columns + ")"
	if _, err := facades.Orm().Query().Exec(query); err != nil {
		facades.Log().Error("Failed to create index "+indexName, map[string]interface{}{
			"error": err.Error(),
		})
	}
}

// Up Run the migrations.
// Performance indexes to address slow queries identified in logs:
// - SLOW SELECT * FROM smes WHERE deleted_at IS NULL ORDER BY id DESC LIMIT 20 (217.885ms)
// - SLOW SELECT * FROM users WHERE id = 1... (283.916ms and 209.803ms)
// - N+1 queries on user_roles and role_permissions
func (r *M20251127085740AddPerformanceIndexes) Up() error {
	// Add indexes to smes table
	r.createIndexIfNotExists("smes_deleted_at_id_index", "smes", "deleted_at, id")
	r.createIndexIfNotExists("smes_deleted_at_created_at_index", "smes", "deleted_at, created_at")
	r.createIndexIfNotExists("smes_created_by_index", "smes", "created_by")
	r.createIndexIfNotExists("smes_region_index", "smes", "region")
	r.createIndexIfNotExists("smes_district_index", "smes", "district")
	r.createIndexIfNotExists("smes_business_category_index", "smes", "business_category")

	// Add indexes to users table
	r.createIndexIfNotExists("users_deleted_at_index", "users", "deleted_at")

	// Add indexes to user_roles table
	r.createIndexIfNotExists("user_roles_user_id_is_active_index", "user_roles", "user_id, is_active")
	r.createIndexIfNotExists("user_roles_role_id_is_active_index", "user_roles", "role_id, is_active")

	// Add indexes to roles table
	r.createIndexIfNotExists("roles_deleted_at_is_active_index", "roles", "deleted_at, is_active")
	r.createIndexIfNotExists("roles_slug_is_active_index", "roles", "slug, is_active")

	// Add indexes to role_permissions table
	r.createIndexIfNotExists("role_permissions_role_id_is_active_index", "role_permissions", "role_id, is_active")
	r.createIndexIfNotExists("role_permissions_permission_id_is_active_index", "role_permissions", "permission_id, is_active")

	// Add indexes to permissions table
	r.createIndexIfNotExists("permissions_slug_is_active_index", "permissions", "slug, is_active")

	return nil
}

// dropIndexIfExists drops an index using raw SQL with IF EXISTS for PostgreSQL compatibility
func (r *M20251127085740AddPerformanceIndexes) dropIndexIfExists(indexName string) {
	query := "DROP INDEX IF EXISTS " + indexName
	if _, err := facades.Orm().Query().Exec(query); err != nil {
		facades.Log().Error("Failed to drop index "+indexName, map[string]interface{}{
			"error": err.Error(),
		})
	}
}

// Down Reverse the migrations.
func (r *M20251127085740AddPerformanceIndexes) Down() error {
	// Drop smes indexes
	r.dropIndexIfExists("smes_deleted_at_id_index")
	r.dropIndexIfExists("smes_deleted_at_created_at_index")
	r.dropIndexIfExists("smes_created_by_index")
	r.dropIndexIfExists("smes_region_index")
	r.dropIndexIfExists("smes_district_index")
	r.dropIndexIfExists("smes_business_category_index")

	// Drop users indexes
	r.dropIndexIfExists("users_deleted_at_index")

	// Drop user_roles indexes
	r.dropIndexIfExists("user_roles_user_id_is_active_index")
	r.dropIndexIfExists("user_roles_role_id_is_active_index")

	// Drop roles indexes
	r.dropIndexIfExists("roles_deleted_at_is_active_index")
	r.dropIndexIfExists("roles_slug_is_active_index")

	// Drop role_permissions indexes
	r.dropIndexIfExists("role_permissions_role_id_is_active_index")
	r.dropIndexIfExists("role_permissions_permission_id_is_active_index")

	// Drop permissions indexes
	r.dropIndexIfExists("permissions_slug_is_active_index")

	return nil
}
