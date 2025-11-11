package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251020064006AddJsonFieldsToSmesTable struct{}

// Signature The unique signature for the migration.
func (r *M20251020064006AddJsonFieldsToSmesTable) Signature() string {
	return "20251020064006_add_json_fields_to_smes_table"
}

// Up Run the migrations.
func (r *M20251020064006AddJsonFieldsToSmesTable) Up() error {
	return facades.Schema().Table("smes", func(table schema.Blueprint) {
		table.Text("business_improvement_aspect_json").Nullable()
		table.Text("business_accessed_financing_json").Nullable()
	})
}

// Down Reverse the migrations.
func (r *M20251020064006AddJsonFieldsToSmesTable) Down() error {
	return facades.Schema().Table("smes", func(table schema.Blueprint) {
		table.DropColumn("business_improvement_aspect_json")
		table.DropColumn("business_accessed_financing_json")
	})
}
