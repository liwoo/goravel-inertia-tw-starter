package messages

import (
	"os"
	"testing"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/foundation"

	_ "starter-project/config"
	"starter-project/tests/helpers"
)

func TestMain(m *testing.M) {
	// Boot the Goravel application
	app := foundation.NewApplication()
	app.Boot()

	// Run migrations
	if err := facades.Artisan().Call("migrate"); err != nil {
		os.Stderr.WriteString("Warning: Migration failed: " + err.Error() + "\n")
	}

	// Clean database before running tests to ensure isolation
	if err := helpers.CleanTestDatabase(); err != nil {
		os.Stderr.WriteString("Warning: Database cleanup failed: " + err.Error() + "\n")
	}

	// Run all tests
	code := m.Run()

	// Clean up after tests complete
	helpers.CleanTestDatabase()

	os.Exit(code)
}
