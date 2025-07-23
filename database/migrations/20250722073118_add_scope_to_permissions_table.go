package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250722073118AddScopeToPermissionsTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250722073118AddScopeToPermissionsTable) Signature() string {
	return "20250722073118_add_scope_to_permissions_table"
}

// Up Run the migrations.
func (r *M20250722073118AddScopeToPermissionsTable) Up() error {
	return facades.Schema().Table("permissions", func(table schema.Blueprint) {
		table.String("scope").Default("by_all").Comment("Permission scope: by_me, by_my_role, by_all")
		table.Index("scope")
	})
}

// Down Reverse the migrations.
func (r *M20250722073118AddScopeToPermissionsTable) Down() error {
	return facades.Schema().Table("permissions", func(table schema.Blueprint) {
		table.DropIndex("scope")
		table.DropColumn("scope")
	})
}
