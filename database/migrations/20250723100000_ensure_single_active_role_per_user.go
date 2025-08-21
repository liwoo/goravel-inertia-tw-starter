package migrations

import (
	"github.com/goravel/framework/facades"
)

type EnsureSingleActiveRolePerUser20250723100000 struct{}

// Signature The unique signature for the migration.
func (r *EnsureSingleActiveRolePerUser20250723100000) Signature() string {
	return "20250723100000_ensure_single_active_role_per_user"
}

// Up Run the migrations.
func (r *EnsureSingleActiveRolePerUser20250723100000) Up() error {
	// First, ensure only one active role per user
	// Deactivate all but the most recent active role for each user
	driver := facades.Config().GetString("database.default")

	if driver == "sqlite" {
		// SQLite-compatible query
		_, err := facades.Orm().Query().Exec(`
			UPDATE user_roles 
			SET is_active = 0
			WHERE id IN (
				SELECT ur1.id
				FROM user_roles ur1
				INNER JOIN (
					SELECT user_id, MAX(assigned_at) as latest_assigned
					FROM user_roles
					WHERE is_active = 1
					GROUP BY user_id
					HAVING COUNT(*) > 1
				) ur2 ON ur1.user_id = ur2.user_id
				WHERE ur1.is_active = 1 
				AND ur1.assigned_at < ur2.latest_assigned
			)
		`)
		if err != nil {
			return err
		}
	} else {
		// MySQL/PostgreSQL compatible query
		_, err := facades.Orm().Query().Exec(`
			UPDATE user_roles ur1
			SET is_active = false
			WHERE is_active = true
			AND EXISTS (
				SELECT 1
				FROM (
					SELECT user_id, MAX(assigned_at) as latest_assigned
					FROM user_roles
					WHERE is_active = true
					GROUP BY user_id
					HAVING COUNT(*) > 1
				) ur2
				WHERE ur1.user_id = ur2.user_id
				AND ur1.assigned_at < ur2.latest_assigned
			)
		`)
		if err != nil {
			return err
		}
	}

	// Add a partial unique index to ensure only one active role per user
	// This allows multiple inactive roles but only one active role
	// Note: SQLite doesn't support partial indexes in the same way, so we'll skip for SQLite
	if driver != "sqlite" {
		_, err := facades.Orm().Query().Exec(`
			CREATE UNIQUE INDEX idx_user_roles_one_active_per_user 
			ON user_roles(user_id) 
			WHERE is_active = true AND deleted_at IS NULL
		`)
		if err != nil {
			return err
		}
	}

	return nil
}

// Down Reverse the migrations.
func (r *EnsureSingleActiveRolePerUser20250723100000) Down() error {
	// Drop the unique index
	_, err := facades.Orm().Query().Exec(`DROP INDEX IF EXISTS idx_user_roles_one_active_per_user`)
	return err
}
