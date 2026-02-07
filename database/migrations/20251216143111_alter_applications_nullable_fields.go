package migrations

import (
	"github.com/goravel/framework/facades"
)

type M20251216143111AlterApplicationsNullableFields struct{}

// Signature The unique signature for the migration.
func (r *M20251216143111AlterApplicationsNullableFields) Signature() string {
	return "20251216143111_alter_applications_nullable_fields"
}

// Up Run the migrations.
// Make some fields nullable to support flexibility
func (r *M20251216143111AlterApplicationsNullableFields) Up() error {
	// Use raw SQL to alter columns to be nullable
	_, err := facades.Orm().Query().Exec(`
		ALTER TABLE applications
		ALTER COLUMN sme DROP NOT NULL,
		ALTER COLUMN registrant_name DROP NOT NULL,
		ALTER COLUMN email DROP NOT NULL
	`)
	return err
}

// Down Reverse the migrations.
func (r *M20251216143111AlterApplicationsNullableFields) Down() error {
	// Note: This may fail if there are NULL values in these columns
	_, err := facades.Orm().Query().Exec(`
		ALTER TABLE applications
		ALTER COLUMN sme SET NOT NULL,
		ALTER COLUMN registrant_name SET NOT NULL,
		ALTER COLUMN email SET NOT NULL
	`)
	return err
}
