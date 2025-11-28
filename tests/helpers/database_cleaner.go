package helpers

import (
	"os"

	"github.com/goravel/framework/facades"
)

// CleanTestDatabase removes all test data from the database
// Works with both SQLite and PostgreSQL
func CleanTestDatabase() error {
	orm := facades.Orm()
	if orm == nil {
		return nil
	}

	dbConnection := os.Getenv("DB_CONNECTION")

	if dbConnection == "postgres" {
		return cleanPostgresDatabase()
	}

	// Default to SQLite
	return cleanSQLiteDatabase()
}

// cleanSQLiteDatabase removes all test data from SQLite database
func cleanSQLiteDatabase() error {
	orm := facades.Orm()
	if orm == nil {
		return nil
	}

	// Get all tables (SQLite specific)
	var tables []string
	err := orm.Query().Raw(`
		SELECT name FROM sqlite_master
		WHERE type='table'
		AND name NOT LIKE 'sqlite_%'
		AND name NOT LIKE 'migrations'
		ORDER BY name
	`).Scan(&tables)

	if err != nil {
		return err
	}

	// Disable foreign key constraints temporarily
	orm.Query().Exec("PRAGMA foreign_keys = OFF")

	// Clear all tables except migrations
	for _, table := range tables {
		if table != "migrations" {
			orm.Query().Exec("DELETE FROM " + table)
		}
	}

	// Re-enable foreign key constraints
	orm.Query().Exec("PRAGMA foreign_keys = ON")

	return nil
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
func ResetTestDatabase() error {
	dbConnection := os.Getenv("DB_CONNECTION")

	if dbConnection == "postgres" {
		return resetPostgresDatabase()
	}

	// Default to SQLite
	return resetSQLiteDatabase()
}

// resetSQLiteDatabase resets the SQLite test database
func resetSQLiteDatabase() error {
	dbPath := os.Getenv("DB_DATABASE")
	if dbPath == "" {
		dbPath = "database/test.sqlite"
	}

	// Delete the existing database file
	os.Remove(dbPath)

	// Create a new empty file
	file, err := os.Create(dbPath)
	if err != nil {
		return err
	}
	file.Close()

	// Run migrations
	return facades.Artisan().Call("migrate")
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
