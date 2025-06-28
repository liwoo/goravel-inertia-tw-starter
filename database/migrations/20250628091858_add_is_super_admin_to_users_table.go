package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250628091858AddIsSuperAdminToUsersTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250628091858AddIsSuperAdminToUsersTable) Signature() string {
	return "20250628091858_add_is_super_admin_to_users_table"
}

// Up Run the migrations.
func (r *M20250628091858AddIsSuperAdminToUsersTable) Up() error {
	return facades.Schema().Table("users", func(table schema.Blueprint) {
		table.Boolean("is_super_admin").Default(false)
		// Add index for better query performance
		table.Index("is_super_admin")
	})
}

// Down Reverse the migrations.
func (r *M20250628091858AddIsSuperAdminToUsersTable) Down() error {
	return facades.Schema().Table("users", func(table schema.Blueprint) {
		table.DropIndex("users_is_super_admin_index")
		table.DropColumn("is_super_admin")
	})
}