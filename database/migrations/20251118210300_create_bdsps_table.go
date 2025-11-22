package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251118210300CreateBdspsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251118210300CreateBdspsTable) Signature() string {
	return "20251118210300_create_bdsps_table"
}

// Up Run the migrations.
func (r *M20251118210300CreateBdspsTable) Up() error {
	if !facades.Schema().HasTable("bdsps") {
		return facades.Schema().Create("bdsps", func(table schema.Blueprint) {
			table.ID()
			table.String("name", 255)
			table.String("postal_address", 255).Nullable()
			table.String("physical_address", 255).Nullable()
			table.Enum("registration_status", []any{"Active", "Inactive", "Pending", "Rejected"}).Default("Pending")
			table.Text("partners_json").Nullable()
			table.Text("product_types_json").Nullable()
			table.Text("service_list_json").Nullable()

			// service list
			// 	name
			// 	cost
			// 	duration
			// Product types
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251118210300CreateBdspsTable) Down() error {
	return facades.Schema().DropIfExists("bdsps")
}
