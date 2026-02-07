package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251216102504AddTypeToApplicationsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251216102504AddTypeToApplicationsTable) Signature() string {
	return "20251216102504_add_type_to_applications_table"
}

// Up Run the migrations.
func (r *M20251216102504AddTypeToApplicationsTable) Up() error {
	return facades.Schema().Table("applications", func(table schema.Blueprint) {
		// Application type: 'signup' for new registration
		table.Enum("type", []any{"signup"}).Default("signup")

		// JSON field to store extra details
		table.Text("data").Nullable()

		// Rejection reason when application is rejected
		table.Text("rejection_reason").Nullable()
	})
}

// Down Reverse the migrations.
func (r *M20251216102504AddTypeToApplicationsTable) Down() error {
	return facades.Schema().Table("applications", func(table schema.Blueprint) {
		table.DropColumn("type", "data", "rejection_reason")
	})
}
