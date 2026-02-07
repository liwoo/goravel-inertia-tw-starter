package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type CreateConfigsTable struct {
}

// Signature The name and signature of the console command.
func (receiver *CreateConfigsTable) Signature() string {
	return "20251019163714_create_configs_table"
}

// Description The console command description.
func (receiver *CreateConfigsTable) Description() string {
	return "Create configs table for lookup/configuration values"
}

// Up Run the migrations.
func (receiver *CreateConfigsTable) Up() error {
	return facades.Schema().Create("configs", func(table schema.Blueprint) {
		table.ID()
		table.String("name")
		table.String("code").Nullable()
		table.String("config_type")
		table.String("description").Nullable()

		// Audit fields
		table.UnsignedBigInteger("created_by").Nullable()
		table.UnsignedBigInteger("updated_by").Nullable()
		table.UnsignedBigInteger("deleted_by").Nullable()
		table.String("ip_address", 45).Nullable()
		table.Text("user_agent").Nullable()

		table.Timestamps()
		table.SoftDeletes()

		table.Index("code")
		table.Index("config_type")
		table.Index("created_by")

		table.Foreign("created_by").References("id").On("users")
		table.Foreign("updated_by").References("id").On("users")
		table.Foreign("deleted_by").References("id").On("users")
	})
}

// Down Reverse the migrations.
func (receiver *CreateConfigsTable) Down() error {
	return facades.Schema().DropIfExists("configs")
}
