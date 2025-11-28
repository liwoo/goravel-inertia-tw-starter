// Package bootstrap provides test environment setup using testcontainers.
//
// IMPORTANT: This file must be processed FIRST before any config loading.
// Go processes files in alphabetical order within a package, and init() functions
// run in the order they are compiled. By naming this file with an early alphabetic
// prefix and using an init() function, we ensure env vars are set early.
package bootstrap

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// earlyContainer holds the container started during package init
var earlyContainer *postgres.PostgresContainer
var earlyContainerErr error

func init() {
	// Start testcontainer during package init, BEFORE config is loaded
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(TestDBName),
		postgres.WithUsername(TestDBUser),
		postgres.WithPassword(TestDBPassword),
		postgres.WithSQLDriver("pgx"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		earlyContainerErr = fmt.Errorf("failed to start postgres container: %w", err)
		return
	}

	host, err := container.Host(ctx)
	if err != nil {
		earlyContainerErr = fmt.Errorf("failed to get container host: %w", err)
		return
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		earlyContainerErr = fmt.Errorf("failed to get container port: %w", err)
		return
	}

	// Set environment variables BEFORE config package is loaded
	os.Setenv("APP_ENV", "testing")
	os.Setenv("APP_DEBUG", "true")
	os.Setenv("APP_KEY", "testkeyfortestingonlyabc12345678")
	os.Setenv("DB_CONNECTION", "postgres")
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_DATABASE", TestDBName)
	os.Setenv("DB_USERNAME", TestDBUser)
	os.Setenv("DB_PASSWORD", TestDBPassword)
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("REDIS_HOST", "")
	os.Setenv("JWT_SECRET", "test-jwt-secret-key-for-testing-only-do-not-use-in-production")

	earlyContainer = container
	log.Printf("Test container started: %s:%s (database: %s)", host, port.Port(), TestDBName)
}
