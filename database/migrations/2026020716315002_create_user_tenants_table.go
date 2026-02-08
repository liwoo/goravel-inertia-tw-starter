package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M2026020716315002CreateUserTenantsTable struct{}

func (r *M2026020716315002CreateUserTenantsTable) Signature() string {
	return "2026020716315002_create_user_tenants_table"
}

func (r *M2026020716315002CreateUserTenantsTable) Up() error {
	if !facades.Schema().HasTable("user_tenants") {
		return facades.Schema().Create("user_tenants", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("user_id")
			table.UnsignedBigInteger("tenant_id")
			table.UnsignedBigInteger("role_id").Nullable()
			table.Boolean("is_active").Default(true)
			table.TimestampTz("joined_at").UseCurrent()

			// Audit fields
			table.UnsignedBigInteger("tenant_id_ref").Nullable()
			table.UnsignedBigInteger("created_by").Nullable()
			table.UnsignedBigInteger("updated_by").Nullable()
			table.UnsignedBigInteger("deleted_by").Nullable()
			table.String("ip_address", 45).Nullable()
			table.Text("user_agent").Nullable()

			table.TimestampsTz()
			table.SoftDeletesTz()

			// Indexes and constraints
			table.Index("user_id")
			table.Index("tenant_id")
		})
	}
	return nil
}

func (r *M2026020716315002CreateUserTenantsTable) Down() error {
	return facades.Schema().DropIfExists("user_tenants")
}
