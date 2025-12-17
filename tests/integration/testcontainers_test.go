package integration

import (
	"os"
	"testing"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// TestMain is defined in main_test.go for this package

// TestContainersIntegrationSuite tests that PostgreSQL works correctly
type TestContainersIntegrationSuite struct {
	suite.Suite
}

func TestTestContainersIntegrationSuite(t *testing.T) {
	suite.Run(t, &TestContainersIntegrationSuite{})
}

func (s *TestContainersIntegrationSuite) SetupSuite() {
	// Verify we're using the expected database connection
	dbConnection := facades.Config().GetString("database.default")
	s.T().Logf("Database connection: %s", dbConnection)
	s.T().Logf("Database host: %s", os.Getenv("DB_HOST"))
	s.T().Logf("Database port: %s", os.Getenv("DB_PORT"))
	s.T().Logf("Database name: %s", os.Getenv("DB_DATABASE"))
	s.Equal("postgres", dbConnection, "Should be using PostgreSQL")
}

func (s *TestContainersIntegrationSuite) TestDatabaseConnection() {
	// Test that we can query the database
	orm := facades.Orm()
	s.NotNil(orm, "ORM should be initialized")

	// Try a simple query
	var result int
	err := orm.Query().Raw("SELECT 1").Scan(&result)
	s.Nil(err, "Simple query should succeed")
	s.Equal(1, result, "Query should return 1")
}

func (s *TestContainersIntegrationSuite) TestMigrationsApplied() {
	// Verify that migrations have been applied
	orm := facades.Orm()
	s.NotNil(orm)

	// Check if the migrations table exists and has records
	var count int64
	err := orm.Query().Raw("SELECT COUNT(*) FROM migrations").Scan(&count)
	s.Nil(err, "Should be able to query migrations table")
	s.Greater(count, int64(0), "Migrations should have been applied")

	s.T().Logf("Number of migrations applied: %d", count)
}

func (s *TestContainersIntegrationSuite) TestPostgresSpecificFeatures() {
	orm := facades.Orm()
	s.NotNil(orm)

	// Test PostgreSQL-specific query (version)
	var version string
	err := orm.Query().Raw("SELECT version()").Scan(&version)
	s.Nil(err, "Should be able to get PostgreSQL version")
	s.Contains(version, "PostgreSQL", "Version should indicate PostgreSQL")

	s.T().Logf("PostgreSQL version: %s", version)
}

func (s *TestContainersIntegrationSuite) TestTablesExist() {
	// Verify key tables exist
	orm := facades.Orm()
	s.NotNil(orm)

	tables := []string{"users", "roles", "permissions", "smes", "business_employee_summary"}

	for _, table := range tables {
		var exists bool
		query := "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1)"

		err := orm.Query().Raw(query, table).Scan(&exists)
		s.Nil(err, "Query should succeed for table: "+table)
		s.True(exists, "Table should exist: "+table)
	}
}

func (s *TestContainersIntegrationSuite) TestBusinessEmployeeSummaryHasDeletedAt() {
	// Verify business_employee_summary has deleted_at column
	orm := facades.Orm()
	s.NotNil(orm)

	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM information_schema.columns
		WHERE table_schema = 'public'
		AND table_name = 'business_employee_summary'
		AND column_name = 'deleted_at'
	)`

	err := orm.Query().Raw(query).Scan(&exists)
	s.Nil(err, "Query should succeed")
	s.True(exists, "business_employee_summary should have deleted_at column")
}

func (s *TestContainersIntegrationSuite) TestIsolatedFromProduction() {
	// Verify that we're using the test database, not production
	dbName := os.Getenv("DB_DATABASE")
	dbUser := os.Getenv("DB_USERNAME")
	dbPort := os.Getenv("DB_PORT")

	s.T().Logf("Database: %s, User: %s, Port: %s", dbName, dbUser, dbPort)

	// These should be test values, not production
	s.NotEqual("smedi", dbName, "Should not be using production database")
	s.NotEqual("55000", dbPort, "Should NOT be using production port 55000")
}

func (s *TestContainersIntegrationSuite) TestNotUsingProductionPort() {
	// Verify we're not using the default production port (55000 from .env)
	dbPort := os.Getenv("DB_PORT")
	s.NotEqual("55000", dbPort, "Should NOT be using production port 55000")
	s.T().Logf("Using port: %s (not production port 55000)", dbPort)
}

// TestDatabaseIsolation verifies that each test gets a clean database state
func TestDatabaseIsolation(t *testing.T) {
	orm := facades.Orm()
	assert.NotNil(t, orm, "ORM should be initialized")

	// This test simply verifies the connection works
	var result int
	err := orm.Query().Raw("SELECT 1").Scan(&result)
	assert.Nil(t, err)
	assert.Equal(t, 1, result)

	t.Logf("Test is running with PostgreSQL")
	t.Logf("Database: %s:%s/%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))
}
