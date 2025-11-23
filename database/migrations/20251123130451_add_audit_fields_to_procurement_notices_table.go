package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251123130451AddAuditFieldsToProcurementNoticesTable struct {
}

// Signature The unique signature for the migration.
func (r *M20251123130451AddAuditFieldsToProcurementNoticesTable) Signature() string {
	return "20251123130451_add_audit_fields_to_procurement_notices_table"
}

// Up Run the migrations.
func (r *M20251123130451AddAuditFieldsToProcurementNoticesTable) Up() error {
	return facades.Schema().Table("procurement_notices", func(table schema.Blueprint) {
		table.UnsignedBigInteger("created_by").Nullable().Comment("User who created this procurement_notice")
		table.UnsignedBigInteger("updated_by").Nullable().Comment("User who last updated this procurement_notice")
		table.UnsignedBigInteger("deleted_by").Nullable().Comment("User who deleted this procurement_notice")
		// deleted_at already exists from SoftDeletesTz() in create migration
		table.String("ip_address", 45).Nullable().Comment("IP address from where the action was performed")
		table.Text("user_agent").Nullable().Comment("User agent from where the action was performed")

		table.Foreign("created_by").References("id").On("users")
		table.Foreign("updated_by").References("id").On("users")
		table.Foreign("deleted_by").References("id").On("users")

		table.Index("created_by")
		table.Index("updated_by")
		table.Index("deleted_by")
		table.Index("deleted_at")
	})
}

// Down Reverse the migrations.
func (r *M20251123130451AddAuditFieldsToProcurementNoticesTable) Down() error {
	return facades.Schema().Table("procurement_notices", func(table schema.Blueprint) {
		table.DropForeign("created_by")
		table.DropForeign("updated_by")
		table.DropForeign("deleted_by")

		table.DropIndex("created_by")
		table.DropIndex("updated_by")
		table.DropIndex("deleted_by")
		table.DropIndex("deleted_at")

		table.DropColumn("created_by")
		table.DropColumn("updated_by")
		table.DropColumn("deleted_by")
		// deleted_at will be dropped by SoftDeletesTz in create migration rollback
		table.DropColumn("ip_address")
		table.DropColumn("user_agent")
	})
}
