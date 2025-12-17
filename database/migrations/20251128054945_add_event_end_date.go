package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251128054945AddEventEndDate struct{}

// Signature The unique signature for the migration.
func (r *M20251128054945AddEventEndDate) Signature() string {
	return "20251128054945_add_event_end_date"
}

// Up Run the migrations.
func (r *M20251128054945AddEventEndDate) Up() error {
	return facades.Schema().Table("events", func(table schema.Blueprint) {
		table.Date("end_date").Nullable()
	})
}

// Down Reverse the migrations.
func (r *M20251128054945AddEventEndDate) Down() error {
	return facades.Schema().Table("events", func(table schema.Blueprint) {
		table.DropColumn("end_date")
	})
}
