package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251125232825CreateApplicationsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251125232825CreateApplicationsTable) Signature() string {
	return "20251125232825_create_applications_table"
}

// Up Run the migrations.
func (r *M20251125232825CreateApplicationsTable) Up() error {
	if !facades.Schema().HasTable("applications") {
		return facades.Schema().Create("applications", func(table schema.Blueprint) {
			table.ID()
			table.String("sme")
			table.String("registrant_name")
			table.String("email")
			table.String("phone").Nullable()
			table.String("sme_registration_number").Nullable()
			table.String("sme_tax_identification_number").Nullable()
			table.Enum("status", []any{"Pending", "Approved", "Rejected"}).Default("Pending")

			// Primary Business Owner Details
			table.String("first_name", 100)
			table.String("last_name", 100)
			table.String("other_names", 100).Nullable()
			table.String("nationality", 100)
			table.String("national_id_number", 50)
			table.Date("date_of_birth")
			table.Enum("gender", []any{"MALE", "FEMALE"})
			table.String("education_level", 100)
			table.String("malawian_status", 50)
			table.Boolean("has_special_needs").Default(false)
			table.String("landline_number", 12).Nullable()
			table.String("physical_address", 255).Nullable()
			table.String("postal_address", 255).Nullable()
			table.String("region", 100).Nullable()
			table.String("district", 100).Nullable()
			table.String("traditional_authority", 100).Nullable()
			table.String("alt_contact_name", 200).Nullable()
			table.String("alt_contact_relationship", 100).Nullable()
			table.String("alt_contact_phone", 12).Nullable()
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251125232825CreateApplicationsTable) Down() error {
	return facades.Schema().DropIfExists("applications")
}
