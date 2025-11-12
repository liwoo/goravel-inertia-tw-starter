package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251019163714CreateSmeConfig struct{}

// Signature The unique signature for the migration.
func (r *M20251019163714CreateSmeConfig) Signature() string {
	return "20251019163714_create_sme_config"
}

// Up Run the migrations.
func (r *M20251019163714CreateSmeConfig) Up() error {
	if !facades.Schema().HasTable("sme_config") {
		return facades.Schema().Create("sme_config", func(table schema.Blueprint) {
			table.ID()
			table.String("name", 255)
			table.String("config_type")
			table.Text("description").Nullable()
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251019163714CreateSmeConfig) Down() error {
	return facades.Schema().DropIfExists("sme_config")
}
