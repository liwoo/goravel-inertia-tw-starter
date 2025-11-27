package migrations

import (
	"github.com/goravel/framework/facades"
)

type M20251127151630AlterPhoneColumnsSize struct{}

// Signature The unique signature for the migration.
func (r *M20251127151630AlterPhoneColumnsSize) Signature() string {
	return "20251127151630_alter_phone_columns_size"
}

// Up Run the migrations.
// Increases phone number columns from varchar(12) to varchar(20) to accommodate
// international phone number formats like +265881234567 (13 chars for Malawi)
func (r *M20251127151630AlterPhoneColumnsSize) Up() error {
	// SMEs table - contact_phone
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE smes ALTER COLUMN contact_phone TYPE varchar(20)",
	); err != nil {
		return err
	}

	// primary_business_owner table - phone_number, landline_number, alt_contact_phone
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE primary_business_owner ALTER COLUMN phone_number TYPE varchar(20)",
	); err != nil {
		return err
	}
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE primary_business_owner ALTER COLUMN landline_number TYPE varchar(20)",
	); err != nil {
		return err
	}
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE primary_business_owner ALTER COLUMN alt_contact_phone TYPE varchar(20)",
	); err != nil {
		return err
	}

	// additional_business_members table - phone_number
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE additional_business_members ALTER COLUMN phone_number TYPE varchar(20)",
	); err != nil {
		return err
	}

	// applications table - landline_number, alt_contact_phone
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE applications ALTER COLUMN landline_number TYPE varchar(20)",
	); err != nil {
		return err
	}
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE applications ALTER COLUMN alt_contact_phone TYPE varchar(20)",
	); err != nil {
		return err
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251127151630AlterPhoneColumnsSize) Down() error {
	// Revert to varchar(12) - note: this may fail if data exceeds 12 chars
	facades.Orm().Query().Exec("ALTER TABLE smes ALTER COLUMN contact_phone TYPE varchar(12)")
	facades.Orm().Query().Exec("ALTER TABLE primary_business_owner ALTER COLUMN phone_number TYPE varchar(12)")
	facades.Orm().Query().Exec("ALTER TABLE primary_business_owner ALTER COLUMN landline_number TYPE varchar(12)")
	facades.Orm().Query().Exec("ALTER TABLE primary_business_owner ALTER COLUMN alt_contact_phone TYPE varchar(12)")
	facades.Orm().Query().Exec("ALTER TABLE additional_business_members ALTER COLUMN phone_number TYPE varchar(12)")
	facades.Orm().Query().Exec("ALTER TABLE applications ALTER COLUMN landline_number TYPE varchar(12)")
	facades.Orm().Query().Exec("ALTER TABLE applications ALTER COLUMN alt_contact_phone TYPE varchar(12)")
	return nil
}
