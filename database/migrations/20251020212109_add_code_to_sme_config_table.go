package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251020212109AddCodeToSmeConfigTable struct{}

// Signature The unique signature for the migration.
func (r *M20251020212109AddCodeToSmeConfigTable) Signature() string {
	return "20251020212109_add_code_to_sme_config_table"
}

// Up Run the migrations.
func (r *M20251020212109AddCodeToSmeConfigTable) Up() error {
	return facades.Schema().Table("sme_config", func(table schema.Blueprint) {
		// Add code field to store short codes for configs (e.g., "AGR" for Agriculture sector)
		table.String("code", 50).Nullable().After("name")
		table.Index("code")
	})
}

// Down Reverse the migrations.
func (r *M20251020212109AddCodeToSmeConfigTable) Down() error {
	return facades.Schema().Table("sme_config", func(table schema.Blueprint) {
		table.DropColumn("code")
	})
}
