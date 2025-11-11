package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251020052629CreateSmeTable struct{}

// Signature The unique signature for the migration.
func (r *M20251020052629CreateSmeTable) Signature() string {
	return "20251020052629_create_sme_table"
}

// Up Run the migrations.
func (r *M20251020052629CreateSmeTable) Up() error {
	if !facades.Schema().HasTable("smes") {
		return facades.Schema().Create("smes", func(table schema.Blueprint) {
			table.ID()
			table.String("usme_number", 100)
			table.String("name", 255)
			table.String("registration_number", 100).Nullable()
			table.String("tax_identification_number", 100).Nullable()
			table.Date("operational_start_date").Nullable()
			table.String("business_category", 100)
			table.String("sector", 100)
			table.String("sub_sector", 100).Nullable()
			table.Text("business_description").Nullable()
			table.String("contact_phone", 12)
			table.String("contact_email", 100)
			table.String("physical_address", 255).Nullable()
			table.String("postal_address", 255).Nullable()
			table.String("website", 255).Nullable()
			table.String("region", 100).Nullable()
			table.String("district", 100).Nullable()
			table.String("traditional_authority", 100).Nullable()
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251020052629CreateSmeTable) Down() error {
	return facades.Schema().DropIfExists("sme")
}
