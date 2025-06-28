package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250626020345CreateRolePermissionsTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250626020345CreateRolePermissionsTable) Signature() string {
	return "20250626020345_create_role_permissions_table"
}

// Up Run the migrations.
func (r *M20250626020345CreateRolePermissionsTable) Up() error {
	return facades.Schema().Create("role_permissions", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("role_id")
		table.UnsignedBigInteger("permission_id")
		table.UnsignedBigInteger("granted_by_id").Nullable()
		table.Timestamp("granted_at").Nullable()
		table.Text("notes").Nullable()
		table.Boolean("is_active").Default(true)
		table.Timestamps()
		table.SoftDeletes()
		
		// Add indexes
		table.Index("role_id")
		table.Index("permission_id")
		table.Index("granted_by_id")
		
		// Add foreign key constraints
		table.Foreign("role_id").References("id").On("roles")
		table.Foreign("permission_id").References("id").On("permissions")
		table.Foreign("granted_by_id").References("id").On("users")
	})
}

// Down Reverse the migrations.
func (r *M20250626020345CreateRolePermissionsTable) Down() error {
 	return facades.Schema().DropIfExists("role_permissions")
}
