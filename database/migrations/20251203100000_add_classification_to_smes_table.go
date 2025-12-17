package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251203100000AddClassificationToSmesTable struct{}

// Signature The unique signature for the migration.
func (r *M20251203100000AddClassificationToSmesTable) Signature() string {
	return "20251203100000_add_classification_to_smes_table"
}

// Up Run the migrations.
func (r *M20251203100000AddClassificationToSmesTable) Up() error {
	return facades.Schema().Table("smes", func(table schema.Blueprint) {
		table.String("classification", 20).Default("Unclassified")
	})
}

// Down Reverse the migrations.
func (r *M20251203100000AddClassificationToSmesTable) Down() error {
	return facades.Schema().Table("smes", func(table schema.Blueprint) {
		table.DropColumn("classification")
	})
}
