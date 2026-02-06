package migrations

import (
	"github.com/goravel/framework/facades"
)

type M20260119100000ExpandPhoneNumberFields struct{}

// Signature returns the unique signature for the migration.
func (r *M20260119100000ExpandPhoneNumberFields) Signature() string {
	return "20260119100000_expand_phone_number_fields"
}

// Up executes the migration.
// Expands phone number fields from varchar(20) to varchar(50) to accommodate
// longer phone numbers that may include country codes, extensions, or formatting.
func (r *M20260119100000ExpandPhoneNumberFields) Up() error {
	// Expand phone fields in smes table
	_, err := facades.Orm().Query().Exec(`
		ALTER TABLE smes
		ALTER COLUMN contact_phone TYPE varchar(50)
	`)
	if err != nil {
		return err
	}

	// Expand phone fields in primary_business_owner table
	_, err = facades.Orm().Query().Exec(`
		ALTER TABLE primary_business_owner
		ALTER COLUMN phone_number TYPE varchar(50),
		ALTER COLUMN landline_number TYPE varchar(50),
		ALTER COLUMN alt_contact_phone TYPE varchar(50)
	`)
	if err != nil {
		return err
	}

	// Expand phone fields in additional_business_members table
	_, err = facades.Orm().Query().Exec(`
		ALTER TABLE additional_business_members
		ALTER COLUMN phone_number TYPE varchar(50)
	`)
	if err != nil {
		return err
	}

	// Expand phone fields in applications table
	_, err = facades.Orm().Query().Exec(`
		ALTER TABLE applications
		ALTER COLUMN landline_number TYPE varchar(50),
		ALTER COLUMN alt_contact_phone TYPE varchar(50)
	`)
	if err != nil {
		return err
	}

	return nil
}

// Down reverses the migration.
func (r *M20260119100000ExpandPhoneNumberFields) Down() error {
	// We don't want to shrink fields as it could cause data loss
	return nil
}
