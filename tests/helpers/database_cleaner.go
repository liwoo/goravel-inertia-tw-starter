package helpers

import (
	"os"
	"github.com/goravel/framework/facades"
)

// CleanTestDatabase removes all test data from the database
func CleanTestDatabase() error {
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

// ResetTestDatabase completely recreates the test database
func ResetTestDatabase() error {
	// For SQLite, we can just delete and recreate the file
	if os.Getenv("DB_CONNECTION") == "sqlite" {
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
		err = facades.Artisan().Call("migrate")
		if err != nil {
			return err
		}
	}
	
	return nil
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