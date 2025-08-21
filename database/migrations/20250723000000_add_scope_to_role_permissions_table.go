package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250723000000AddScopeToRolePermissionsTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250723000000AddScopeToRolePermissionsTable) Signature() string {
	return "20250723000000_add_scope_to_role_permissions_table"
}

// Up Run the migrations.
func (r *M20250723000000AddScopeToRolePermissionsTable) Up() error {
	return facades.Schema().Table("role_permissions", func(table schema.Blueprint) {
		table.String("scope").Default("by_all").Comment("Permission scope for this role: by_me, by_my_role, by_all")
		table.Index("scope")
	})
}

// Down Reverse the migrations.
func (r *M20250723000000AddScopeToRolePermissionsTable) Down() error {
	return facades.Schema().Table("role_permissions", func(table schema.Blueprint) {
		table.DropIndex("scope")
		table.DropColumn("scope")
	})
}
