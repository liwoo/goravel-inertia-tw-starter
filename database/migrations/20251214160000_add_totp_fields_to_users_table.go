package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251214160000AddTotpFieldsToUsersTable struct{}

// Signature The unique signature for the migration.
func (r *M20251214160000AddTotpFieldsToUsersTable) Signature() string {
	return "20251214160000_add_totp_fields_to_users_table"
}

// Up Run the migrations.
func (r *M20251214160000AddTotpFieldsToUsersTable) Up() error {
	return facades.Schema().Table("users", func(table schema.Blueprint) {
		// TOTP enabled flag
		table.Boolean("totp_enabled").Default(false)

		// Encrypted TOTP secret
		table.Text("totp_secret").Nullable()

		// When TOTP was verified/enabled
		table.TimestampTz("totp_verified_at").Nullable()

		// Index for querying users with TOTP enabled
		table.Index("totp_enabled")
	})
}

// Down Reverse the migrations.
func (r *M20251214160000AddTotpFieldsToUsersTable) Down() error {
	return facades.Schema().Table("users", func(table schema.Blueprint) {
		table.DropIndex("totp_enabled")
		table.DropColumn("totp_verified_at")
		table.DropColumn("totp_secret")
		table.DropColumn("totp_enabled")
	})
}
