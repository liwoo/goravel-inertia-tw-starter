package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250626020332CreatePermissionsTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250626020332CreatePermissionsTable) Signature() string {
	return "20250626020332_create_permissions_table"
}

// Up Run the migrations.
func (r *M20250626020332CreatePermissionsTable) Up() error {
	return facades.Schema().Create("permissions", func(table schema.Blueprint) {
		table.ID()
		table.String("name")
		table.String("slug")
		table.Text("description")
		table.String("category")
		table.String("resource")
		table.String("action")
		table.Boolean("is_active").Default(true)
		table.Boolean("requires_ownership").Default(false)
		table.Boolean("can_delegate").Default(false)
		table.Timestamps()
		table.SoftDeletes()
	})
}

// Down Reverse the migrations.
func (r *M20250626020332CreatePermissionsTable) Down() error {
 	return facades.Schema().DropIfExists("permissions")
}
