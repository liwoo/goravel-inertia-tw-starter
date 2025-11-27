package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251127085740AddPerformanceIndexes struct{}

// Signature The unique signature for the migration.
func (r *M20251127085740AddPerformanceIndexes) Signature() string {
	return "20251127085740_add_performance_indexes"
}

// Up Run the migrations.
// Performance indexes to address slow queries identified in logs:
// - SLOW SELECT * FROM smes WHERE deleted_at IS NULL ORDER BY id DESC LIMIT 20 (217.885ms)
// - SLOW SELECT * FROM users WHERE id = 1... (283.916ms and 209.803ms)
// - N+1 queries on user_roles and role_permissions
func (r *M20251127085740AddPerformanceIndexes) Up() error {
	// Add indexes to smes table
	err := facades.Schema().Table("smes", func(table schema.Blueprint) {
		// Composite index for soft delete + ordering (most common query pattern)
		// Handles: SELECT * FROM smes WHERE deleted_at IS NULL ORDER BY id DESC
		table.Index("deleted_at", "id")

		// Index for created_at ordering (common for dashboards)
		table.Index("deleted_at", "created_at")

		// Index for created_by field (scope filtering)
		table.Index("created_by")

		// Indexes for common filter fields
		table.Index("region")
		table.Index("district")
		table.Index("business_category")
	})
	if err != nil {
		facades.Log().Error("Failed to add indexes to smes table", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Add indexes to users table
	err = facades.Schema().Table("users", func(table schema.Blueprint) {
		// Index for user lookups with soft delete
		table.Index("deleted_at")
	})
	if err != nil {
		facades.Log().Error("Failed to add indexes to users table", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Add indexes to user_roles table
	err = facades.Schema().Table("user_roles", func(table schema.Blueprint) {
		// Composite index for user role lookups
		table.Index("user_id", "is_active")
		table.Index("role_id", "is_active")
	})
	if err != nil {
		facades.Log().Error("Failed to add indexes to user_roles table", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Add indexes to roles table
	err = facades.Schema().Table("roles", func(table schema.Blueprint) {
		// Index for role lookups with soft delete
		table.Index("deleted_at", "is_active")
		table.Index("slug", "is_active")
	})
	if err != nil {
		facades.Log().Error("Failed to add indexes to roles table", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Add indexes to role_permissions table
	err = facades.Schema().Table("role_permissions", func(table schema.Blueprint) {
		// Composite index for role permission lookups
		table.Index("role_id", "is_active")
		table.Index("permission_id", "is_active")
	})
	if err != nil {
		facades.Log().Error("Failed to add indexes to role_permissions table", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Add indexes to permissions table
	err = facades.Schema().Table("permissions", func(table schema.Blueprint) {
		// Index for permission slug lookups
		table.Index("slug", "is_active")
	})
	if err != nil {
		facades.Log().Error("Failed to add indexes to permissions table", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251127085740AddPerformanceIndexes) Down() error {
	// Drop smes indexes
	facades.Schema().Table("smes", func(table schema.Blueprint) {
		table.DropIndex("smes_deleted_at_id_index")
		table.DropIndex("smes_deleted_at_created_at_index")
		table.DropIndex("smes_created_by_index")
		table.DropIndex("smes_region_index")
		table.DropIndex("smes_district_index")
		table.DropIndex("smes_business_category_index")
	})

	// Drop users indexes
	facades.Schema().Table("users", func(table schema.Blueprint) {
		table.DropIndex("users_deleted_at_index")
	})

	// Drop user_roles indexes
	facades.Schema().Table("user_roles", func(table schema.Blueprint) {
		table.DropIndex("user_roles_user_id_is_active_index")
		table.DropIndex("user_roles_role_id_is_active_index")
	})

	// Drop roles indexes
	facades.Schema().Table("roles", func(table schema.Blueprint) {
		table.DropIndex("roles_deleted_at_is_active_index")
		table.DropIndex("roles_slug_is_active_index")
	})

	// Drop role_permissions indexes
	facades.Schema().Table("role_permissions", func(table schema.Blueprint) {
		table.DropIndex("role_permissions_role_id_is_active_index")
		table.DropIndex("role_permissions_permission_id_is_active_index")
	})

	// Drop permissions indexes
	facades.Schema().Table("permissions", func(table schema.Blueprint) {
		table.DropIndex("permissions_slug_is_active_index")
	})

	return nil
}
