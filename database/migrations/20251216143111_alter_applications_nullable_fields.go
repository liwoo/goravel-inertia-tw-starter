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
// Make signup-specific fields nullable to support amendment applications
// which only use type, data, and sme_id fields
func (r *M20251216143111AlterApplicationsNullableFields) Up() error {
	// Use raw SQL to alter columns to be nullable
	_, err := facades.Orm().Query().Exec(`
		ALTER TABLE applications
		ALTER COLUMN sme DROP NOT NULL,
		ALTER COLUMN registrant_name DROP NOT NULL,
		ALTER COLUMN email DROP NOT NULL,
		ALTER COLUMN first_name DROP NOT NULL,
		ALTER COLUMN last_name DROP NOT NULL,
		ALTER COLUMN nationality DROP NOT NULL,
		ALTER COLUMN national_id_number DROP NOT NULL,
		ALTER COLUMN date_of_birth DROP NOT NULL,
		ALTER COLUMN gender DROP NOT NULL,
		ALTER COLUMN education_level DROP NOT NULL,
		ALTER COLUMN malawian_status DROP NOT NULL
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
		ALTER COLUMN email SET NOT NULL,
		ALTER COLUMN first_name SET NOT NULL,
		ALTER COLUMN last_name SET NOT NULL,
		ALTER COLUMN nationality SET NOT NULL,
		ALTER COLUMN national_id_number SET NOT NULL,
		ALTER COLUMN date_of_birth SET NOT NULL,
		ALTER COLUMN gender SET NOT NULL,
		ALTER COLUMN education_level SET NOT NULL,
		ALTER COLUMN malawian_status SET NOT NULL
	`)
	return err
}
