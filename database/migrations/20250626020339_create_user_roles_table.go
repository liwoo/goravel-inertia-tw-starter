package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250626020339CreateUserRolesTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250626020339CreateUserRolesTable) Signature() string {
	return "20250626020339_create_user_roles_table"
}

// Up Run the migrations.
func (r *M20250626020339CreateUserRolesTable) Up() error {
	return facades.Schema().Create("user_roles", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("user_id")
		table.UnsignedBigInteger("role_id")
		table.UnsignedBigInteger("assigned_by_id").Nullable()
		table.Timestamp("assigned_at")
		table.Timestamp("expires_at").Nullable()
		table.Boolean("is_active").Default(true)
		table.Text("notes")
		table.Timestamps()
		table.SoftDeletes()
	})
}

// Down Reverse the migrations.
func (r *M20250626020339CreateUserRolesTable) Down() error {
 	return facades.Schema().DropIfExists("user_roles")
}
