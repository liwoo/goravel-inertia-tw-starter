package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251112001728CreateBdspsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251112001728CreateBdspsTable) Signature() string {
	return "20251112001728_create_bdsps_table"
}

// Up Run the migrations.
func (r *M20251112001728CreateBdspsTable) Up() error {
	if !facades.Schema().HasTable("bdsps") {
		return facades.Schema().Create("bdsps", func(table schema.Blueprint) {
			table.ID()
			table.TimestampsTz()

			// Basic Information
			table.String("name")
			table.String("postal_address").Nullable()
			table.String("physical_address").Nullable()

			// Registration Status
			table.String("registration_status")

			// Product Types - stored as JSON
			table.Text("product_types_json").Nullable()

			// Service List - stored as JSON
			table.Text("service_list_json").Nullable()

			// Associated Partners - stored as JSON
			table.Text("associated_partners_json").Nullable()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251112001728CreateBdspsTable) Down() error {
 	return facades.Schema().DropIfExists("bdsps")
}
