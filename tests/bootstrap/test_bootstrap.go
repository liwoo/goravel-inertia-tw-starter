// Package bootstrap provides an isolated test environment using testcontainers.
// This package uses testcontainers to create a completely isolated PostgreSQL
// database for testing, separate from any production database.
//
// IMPORTANT: The env_setup.go file in this package starts the container and sets
// environment variables during its init() function, which runs before any other
// init() functions. This ensures the testcontainer settings are available when
// the config package is loaded.
package bootstrap

import (
	"context"
	"fmt"
	"sync"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/foundation"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	// Import pgx driver for faster snapshot/restore operations
	_ "github.com/jackc/pgx/v5/stdlib"
	// NOTE: Do NOT import config here - it must be imported AFTER env vars are set
)

const (
	TestDBName     = "starter_project_test"
	TestDBUser     = "testuser"
	TestDBPassword = "testpassword123"
)

// TestPostgresContainer holds the testcontainer instance
type TestPostgresContainer struct {
	Container *postgres.PostgresContainer
	Host      string
	Port      string
	DSN       string
	ctx       context.Context
}

var (
	testContainer  *TestPostgresContainer
	isBootstrapped bool
	hasSnapshot    bool
	bootMutex      sync.Mutex
)

// StartTestContainer returns the early-started container
func StartTestContainer(ctx context.Context) (*TestPostgresContainer, error) {
	if earlyContainerErr != nil {
		return nil, earlyContainerErr
	}

	if earlyContainer == nil {
		return nil, fmt.Errorf("test container not started")
	}

	if testContainer != nil {
		return testContainer, nil
	}

	host, _ := earlyContainer.Host(ctx)
	port, _ := earlyContainer.MappedPort(ctx, "5432")
	dsn, _ := earlyContainer.ConnectionString(ctx, "sslmode=disable")

	testContainer = &TestPostgresContainer{
		Container: earlyContainer,
		Host:      host,
		Port:      port.Port(),
		DSN:       dsn,
		ctx:       ctx,
	}

	return testContainer, nil
}

// TerminateTestContainer stops and removes the test container
func TerminateTestContainer() error {
	bootMutex.Lock()
	defer bootMutex.Unlock()

	if earlyContainer != nil {
		if err := testcontainers.TerminateContainer(earlyContainer); err != nil {
			return fmt.Errorf("failed to terminate container: %w", err)
		}
		earlyContainer = nil
		testContainer = nil
	}
	return nil
}

// BootTestApplication initializes the Goravel application with test-specific configuration
func BootTestApplication(ctx context.Context) error {
	bootMutex.Lock()
	defer bootMutex.Unlock()

	if isBootstrapped {
		return nil
	}

	// Check for early container errors
	if earlyContainerErr != nil {
		return fmt.Errorf("container init failed: %w", earlyContainerErr)
	}

	// Get the container reference
	container, err := StartTestContainer(ctx)
	if err != nil {
		return fmt.Errorf("failed to get test container: %w", err)
	}

	// Initialize Goravel application
	app := foundation.NewApplication()
	app.Boot()

	// CRITICAL: Override database configuration AFTER app boot
	// The config package may have read from .env file during its init()
	// We need to override with testcontainer settings
	overrideDatabaseConfig(container)

	isBootstrapped = true
	return nil
}

// overrideDatabaseConfig overrides the database configuration with testcontainer settings
func overrideDatabaseConfig(container *TestPostgresContainer) {
	config := facades.Config()

	// Override the postgres connection settings
	config.Add("database.connections.postgres.host", container.Host)
	config.Add("database.connections.postgres.port", container.Port)
	config.Add("database.connections.postgres.database", TestDBName)
	config.Add("database.connections.postgres.username", TestDBUser)
	config.Add("database.connections.postgres.password", TestDBPassword)
	config.Add("database.connections.postgres.sslmode", "disable")

	// Ensure postgres is the default connection
	config.Add("database.default", "postgres")
}

// RunMigrations runs all database migrations on the test database
func RunMigrations() error {
	return facades.Artisan().Call("migrate")
}

// CreateSnapshot creates a snapshot of the database after migrations
func CreateSnapshot(ctx context.Context) error {
	bootMutex.Lock()
	defer bootMutex.Unlock()

	if testContainer == nil || testContainer.Container == nil {
		return fmt.Errorf("test container not initialized")
	}

	if hasSnapshot {
		return nil
	}

	err := testContainer.Container.Snapshot(ctx, postgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}

	hasSnapshot = true
	return nil
}

// RestoreSnapshot restores the database to its snapshot state
func RestoreSnapshot(ctx context.Context) error {
	bootMutex.Lock()
	defer bootMutex.Unlock()

	if testContainer == nil || testContainer.Container == nil {
		return fmt.Errorf("test container not initialized")
	}

	if !hasSnapshot {
		return fmt.Errorf("no snapshot available")
	}

	err := testContainer.Container.Restore(ctx)
	if err != nil {
		return fmt.Errorf("failed to restore snapshot: %w", err)
	}

	return nil
}

// CleanDatabase truncates all tables except migrations
func CleanDatabase() error {
	orm := facades.Orm()
	if orm == nil {
		return fmt.Errorf("ORM not initialized")
	}

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

	for _, table := range tables {
		_, err := orm.Query().Exec("TRUNCATE TABLE " + table.TableName + " CASCADE")
		if err != nil {
			orm.Query().Exec("DELETE FROM " + table.TableName)
		}
	}

	return nil
}

// ResetDatabase cleans all data and re-runs migrations
func ResetDatabase() error {
	if err := CleanDatabase(); err != nil {
		return err
	}
	return RunMigrations()
}

// GetContainer returns the current test container
func GetContainer() *TestPostgresContainer {
	return testContainer
}

// IsBootstrapped returns true if the test application has been bootstrapped
func IsBootstrapped() bool {
	bootMutex.Lock()
	defer bootMutex.Unlock()
	return isBootstrapped
}

// HasSnapshot returns true if a database snapshot has been created
func HasSnapshot() bool {
	bootMutex.Lock()
	defer bootMutex.Unlock()
	return hasSnapshot
}

// GetContainerError returns any error that occurred during container startup
func GetContainerError() error {
	return earlyContainerErr
}

// GetTestContext returns the test context
func GetTestContext() context.Context {
	return context.Background()
}
