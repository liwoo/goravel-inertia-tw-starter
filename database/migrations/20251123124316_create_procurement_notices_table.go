package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251123124316CreateProcurementNoticesTable struct{}

// Signature The unique signature for the migration.
func (r *M20251123124316CreateProcurementNoticesTable) Signature() string {
	return "20251123124316_create_procurement_notices_table"
}

// Up Run the migrations.
func (r *M20251123124316CreateProcurementNoticesTable) Up() error {
	if !facades.Schema().HasTable("procurement_notices") {
		return facades.Schema().Create("procurement_notices", func(table schema.Blueprint) {
			table.ID()
			table.String("procured_by", 255)
			table.String("procurement_type", 255)
			table.String("market_approach", 50) // enum: National/International
			table.String("invitation", 50)      // enum: open/limited/single-source
			table.String("ref_no", 255)
			table.Date("open_date")
			table.Date("close_date")
			table.Json("partners")              // []string
			table.Json("qualifying_districts")  // []string
			table.Boolean("is_published").Default(false)
			table.String("organization", 255)
			table.Json("classification")        // []string
			table.Json("interested_smes")       // []string
			table.Text("details")
			table.String("application_details", 500)
			table.Integer("minimum_qualifying_score")
			table.TimestampsTz()
			table.SoftDeletesTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251123124316CreateProcurementNoticesTable) Down() error {
 	return facades.Schema().DropIfExists("procurement_notices")
}
