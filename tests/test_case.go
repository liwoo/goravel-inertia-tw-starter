// Package tests provides test utilities and helpers for the Starter Project test suite.
// Tests should be run using the scripts/run_tests.sh script which sets up the test
// database container and environment variables.
package tests

import (
	"os"
	"testing"

	goraveltesting "github.com/goravel/framework/testing"
)

// RunTestMain is a helper for test packages.
// The actual database setup should be done via the shell script (scripts/run_tests.sh)
func RunTestMain(m *testing.M) {
	os.Exit(m.Run())
}

// IsUsingPostgres returns true if using PostgreSQL
func IsUsingPostgres() bool {
	return os.Getenv("DB_CONNECTION") == "postgres"
}

// GetTestDBPort returns the test database port from environment
func GetTestDBPort() string {
	return os.Getenv("DB_PORT")
}

// GetTestDBHost returns the test database host from environment
func GetTestDBHost() string {
	return os.Getenv("DB_HOST")
}

// GetTestDBName returns the test database name from environment
func GetTestDBName() string {
	return os.Getenv("DB_DATABASE")
}

// TestCase embeds Goravel's TestCase
type TestCase struct {
	goraveltesting.TestCase
}

// RefreshDatabase is a placeholder - database refresh should be handled by test infrastructure
func (t *TestCase) RefreshDatabase() {
	// Database refresh handled externally
}

// Seed is a placeholder for seeder support
func (t *TestCase) Seed(seeders ...string) {
	// Seeders can be run here if needed
}
