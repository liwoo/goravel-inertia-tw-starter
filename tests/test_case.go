package tests

import (
	"os"
	"github.com/goravel/framework/testing"

	"players/bootstrap"
)

func init() {
	// Set environment for testing if not already set
	if os.Getenv("APP_ENV") == "" {
		os.Setenv("APP_ENV", "testing")
	}
	
	// For testing, explicitly set database configuration
	if os.Getenv("APP_ENV") == "testing" {
		os.Setenv("DB_CONNECTION", "sqlite")
		os.Setenv("DB_DATABASE", "database/test.sqlite")
		
		// Create the database directory if it doesn't exist
		os.MkdirAll("database", 0755)
	}
	
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
