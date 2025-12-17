package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251123050933CreateEventsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251123050933CreateEventsTable) Signature() string {
	return "20251123050933_create_events_table"
}

// Up Run the migrations.
func (r *M20251123050933CreateEventsTable) Up() error {
	if !facades.Schema().HasTable("events") {
		return facades.Schema().Create("events", func(table schema.Blueprint) {
			table.ID()
			table.String("title", 255)
			table.Text("description")
			table.Date("date")
			table.String("venue", 255)
			table.Json("partners") // []string stored as JSON
			table.String("district", 100)
			table.Json("attending_smes") // []int stored as JSON for SME IDs
			table.Text("notes").Nullable()
			table.TimestampsTz()
			table.SoftDeletesTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251123050933CreateEventsTable) Down() error {
	return facades.Schema().DropIfExists("events")
}
