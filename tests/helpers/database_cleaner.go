package helpers

import (
	"github.com/goravel/framework/facades"
)

// CleanTestDatabase removes all test data from the database
// Uses PostgreSQL (testcontainers)
func CleanTestDatabase() error {
	orm := facades.Orm()
	if orm == nil {
		return nil
	}

	return cleanPostgresDatabase()
}

// cleanPostgresDatabase removes all test data from PostgreSQL database
func cleanPostgresDatabase() error {
	orm := facades.Orm()
	if orm == nil {
		return nil
	}

	// Get all tables from PostgreSQL
	var tables []struct {
		TableName string `gorm:"column:tablename"`
	}
	err := orm.Query().Raw(`
		SELECT tablename FROM pg_tables
		WHERE schemaname = 'public'
		AND tablename != 'migrations'
		ORDER BY tablename
	`).Scan(&tables)

	if err != nil {
		return err
	}

	// Use TRUNCATE CASCADE for PostgreSQL (faster and handles FKs)
	for _, table := range tables {
		if table.TableName != "migrations" {
			_, err := orm.Query().Exec("TRUNCATE TABLE " + table.TableName + " CASCADE")
			if err != nil {
				// If TRUNCATE fails, try DELETE
				orm.Query().Exec("DELETE FROM " + table.TableName)
			}
		}
	}

	return nil
}

// ResetTestDatabase completely recreates the test database
// Uses PostgreSQL (testcontainers)
func ResetTestDatabase() error {
	return resetPostgresDatabase()
}

// resetPostgresDatabase resets the PostgreSQL test database
func resetPostgresDatabase() error {
	// Clean all tables first
	err := cleanPostgresDatabase()
	if err != nil {
		return err
	}

	// Run migrations with --force flag for non-interactive execution
	return facades.Artisan().Call("migrate --force")
}

// CleanupTestData removes only test-specific data
func CleanupTestData() {
	orm := facades.Orm()
	if orm == nil {
		return
	}

	// Delete in reverse order of dependencies
	queries := []string{
		// Junction tables first
		"DELETE FROM role_permissions",
		"DELETE FROM user_roles",

		// Main tables with test data patterns
		"DELETE FROM books WHERE title LIKE 'Test %' OR author LIKE 'Test %' OR isbn LIKE 'TEST%'",
		"DELETE FROM users WHERE email LIKE '%@test.%' OR email LIKE '%@example.com' OR email LIKE 'test%'",
		"DELETE FROM permissions WHERE slug LIKE 'test_%' OR name LIKE 'Test %'",
		"DELETE FROM roles WHERE slug LIKE 'test_%' OR name LIKE 'Test %'",

		// Clean up any orphaned audit records
		"DELETE FROM books WHERE created_by NOT IN (SELECT id FROM users)",
		"DELETE FROM users WHERE created_by IS NOT NULL AND created_by NOT IN (SELECT id FROM users)",
		"DELETE FROM roles WHERE created_by IS NOT NULL AND created_by NOT IN (SELECT id FROM users)",
		"DELETE FROM permissions WHERE created_by IS NOT NULL AND created_by NOT IN (SELECT id FROM users)",
	}

	// Execute each cleanup query
	for _, query := range queries {
		orm.Query().Exec(query)
	}
}
