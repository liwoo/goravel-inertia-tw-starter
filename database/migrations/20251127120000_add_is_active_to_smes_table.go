package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251127120000AddIsActiveToSmesTable struct{}

// Signature The unique signature for the migration.
func (r *M20251127120000AddIsActiveToSmesTable) Signature() string {
	return "20251127120000_add_is_active_to_smes_table"
}

// Up Run the migrations.
func (r *M20251127120000AddIsActiveToSmesTable) Up() error {
	return facades.Schema().Table("smes", func(table schema.Blueprint) {
		table.Boolean("is_active").Default(true)
	})
}

// Down Reverse the migrations.
func (r *M20251127120000AddIsActiveToSmesTable) Down() error {
	return facades.Schema().Table("smes", func(table schema.Blueprint) {
		table.DropColumn("is_active")
	})
}
