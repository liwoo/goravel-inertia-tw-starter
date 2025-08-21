package migrations

import (
	"github.com/goravel/framework/facades"
)

type PopulateCreatedByForBooks20250723110000 struct{}

// Signature The unique signature for the migration.
func (r *PopulateCreatedByForBooks20250723110000) Signature() string {
	return "20250723110000_populate_created_by_for_books"
}

// Up Run the migrations.
func (r *PopulateCreatedByForBooks20250723110000) Up() error {
	// In test environment with SQLite in-memory, skip this migration
	if facades.Config().GetString("database.default") == "sqlite" &&
		facades.Config().GetString("database.connections.sqlite.database") == ":memory:" {
		return nil
	}

	// Get the first admin/super admin user
	var adminUserID uint
	err := facades.Orm().Query().Raw(`
		SELECT id FROM users 
		WHERE is_super_admin = true OR id IN (
			SELECT user_id FROM user_roles 
			WHERE role_id IN (
				SELECT id FROM roles WHERE slug IN ('admin', 'super-admin')
			) AND is_active = true
		)
		LIMIT 1
	`).Scan(&adminUserID)

	if err != nil || adminUserID == 0 {
		// If no admin found, use the first user
		err = facades.Orm().Query().Raw(`SELECT id FROM users LIMIT 1`).Scan(&adminUserID)
		if err != nil || adminUserID == 0 {
			// No users in database, skip
			return nil
		}
	}

	// Update all books without created_by to use this admin user
	_, err = facades.Orm().Query().Exec(`
		UPDATE books 
		SET created_by = ? 
		WHERE created_by IS NULL
	`, adminUserID)

	return err
}

// Down Reverse the migrations.
func (r *PopulateCreatedByForBooks20250723110000) Down() error {
	// We can't really reverse this migration
	// as we don't know which books originally had NULL created_by
	return nil
}
