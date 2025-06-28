package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250626020324CreateRolesTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250626020324CreateRolesTable) Signature() string {
	return "20250626020324_create_roles_table"
}

// Up Run the migrations.
func (r *M20250626020324CreateRolesTable) Up() error {
	return facades.Schema().Create("roles", func(table schema.Blueprint) {
		table.ID()
		table.String("name")
		table.String("slug")
		table.Text("description")
		table.Integer("level").Default(0)
		table.Boolean("is_active").Default(true)
		table.UnsignedBigInteger("parent_id").Nullable()
		table.Timestamps()
		table.SoftDeletes()
	})
}

// Down Reverse the migrations.
func (r *M20250626020324CreateRolesTable) Down() error {
 	return facades.Schema().DropIfExists("roles")
}
