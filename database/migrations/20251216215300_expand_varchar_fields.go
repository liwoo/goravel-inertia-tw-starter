package migrations

import (
	"github.com/goravel/framework/facades"
)

type M20251216215300ExpandVarcharFields struct{}

// Signature returns the unique signature for the migration.
func (r *M20251216215300ExpandVarcharFields) Signature() string {
	return "20251216215300_expand_varchar_fields"
}

// Up executes the migration.
func (r *M20251216215300ExpandVarcharFields) Up() error {
	// Expand SME fields
	_, err := facades.Orm().Query().Exec(`
		ALTER TABLE smes
		ALTER COLUMN sector TYPE varchar(255),
		ALTER COLUMN sub_sector TYPE varchar(255),
		ALTER COLUMN business_category TYPE varchar(255),
		ALTER COLUMN traditional_authority TYPE varchar(255),
		ALTER COLUMN contact_email TYPE varchar(255),
		ALTER COLUMN registration_number TYPE varchar(255),
		ALTER COLUMN tax_identification_number TYPE varchar(255),
		ALTER COLUMN region TYPE varchar(255),
		ALTER COLUMN district TYPE varchar(255)
	`)
	if err != nil {
		return err
	}

	// Expand primary_business_owner fields
	_, err = facades.Orm().Query().Exec(`
		ALTER TABLE primary_business_owner
		ALTER COLUMN email TYPE varchar(255),
		ALTER COLUMN education_level TYPE varchar(255),
		ALTER COLUMN traditional_authority TYPE varchar(255),
		ALTER COLUMN nationality TYPE varchar(255),
		ALTER COLUMN malawian_status TYPE varchar(255),
		ALTER COLUMN region TYPE varchar(255),
		ALTER COLUMN district TYPE varchar(255),
		ALTER COLUMN alt_contact_name TYPE varchar(255),
		ALTER COLUMN alt_contact_relationship TYPE varchar(255)
	`)
	if err != nil {
		return err
	}

	// Expand additional_business_members fields (only columns that exist)
	_, err = facades.Orm().Query().Exec(`
		ALTER TABLE additional_business_members
		ALTER COLUMN email TYPE varchar(255),
		ALTER COLUMN nationality TYPE varchar(255)
	`)
	if err != nil {
		return err
	}

	return nil
}

// Down reverses the migration.
func (r *M20251216215300ExpandVarcharFields) Down() error {
	// We don't want to shrink fields as it could cause data loss
	return nil
}
