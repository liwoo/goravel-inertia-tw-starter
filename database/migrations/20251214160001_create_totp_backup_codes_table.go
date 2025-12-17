package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251214160001CreateTotpBackupCodesTable struct{}

// Signature The unique signature for the migration.
func (r *M20251214160001CreateTotpBackupCodesTable) Signature() string {
	return "20251214160001_create_totp_backup_codes_table"
}

// Up Run the migrations.
func (r *M20251214160001CreateTotpBackupCodesTable) Up() error {
	if !facades.Schema().HasTable("totp_backup_codes") {
		return facades.Schema().Create("totp_backup_codes", func(table schema.Blueprint) {
			table.ID()

			// User relationship
			table.UnsignedBigInteger("user_id")

			// Bcrypt hashed backup code
			table.String("code_hash", 255)

			// When the code was used (NULL if not used)
			table.TimestampTz("used_at").Nullable()

			// Timestamps
			table.TimestampsTz()

			// Indexes
			table.Index("user_id")
			table.Index("used_at")

			// Foreign key constraint with cascade delete
			table.Foreign("user_id").References("id").On("users").CascadeOnDelete()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251214160001CreateTotpBackupCodesTable) Down() error {
	return facades.Schema().DropIfExists("totp_backup_codes")
}
