package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250626015224AddUserFieldsToUsersTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250626015224AddUserFieldsToUsersTable) Signature() string {
	return "20250626015224_add_user_fields_to_users_table"
}

// Up Run the migrations.
func (r *M20250626015224AddUserFieldsToUsersTable) Up() error {
	return facades.Schema().Table("users", func(table schema.Blueprint) {
		table.Boolean("is_active").Default(true)
		table.Boolean("email_verified").Default(false)
		table.Timestamp("last_login_at").Nullable()
	})
}

// Down Reverse the migrations.
func (r *M20250626015224AddUserFieldsToUsersTable) Down() error {
	return facades.Schema().Table("users", func(table schema.Blueprint) {
		table.DropColumn("is_active")
		table.DropColumn("email_verified")
		table.DropColumn("last_login_at")
	})
}
