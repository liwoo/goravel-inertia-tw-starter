package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M2026020716315001CreateTenantsTable struct{}

func (r *M2026020716315001CreateTenantsTable) Signature() string {
	return "2026020716315001_create_tenants_table"
}

func (r *M2026020716315001CreateTenantsTable) Up() error {
	if !facades.Schema().HasTable("tenants") {
		return facades.Schema().Create("tenants", func(table schema.Blueprint) {
			table.ID()
			table.String("name", 255)
			table.String("slug", 100)
			table.Index("slug")
			table.Boolean("is_active").Default(true)
			table.Boolean("is_main").Default(false)
			table.Text("description").Nullable()
			table.Text("settings").Nullable()
			table.String("logo_url", 500).Nullable()

			// Audit fields
			table.UnsignedBigInteger("tenant_id").Nullable()
			table.UnsignedBigInteger("created_by").Nullable()
			table.UnsignedBigInteger("updated_by").Nullable()
			table.UnsignedBigInteger("deleted_by").Nullable()
			table.String("ip_address", 45).Nullable()
			table.Text("user_agent").Nullable()

			table.TimestampsTz()
			table.SoftDeletesTz()
		})
	}
	return nil
}

func (r *M2026020716315001CreateTenantsTable) Down() error {
	return facades.Schema().DropIfExists("tenants")
}
