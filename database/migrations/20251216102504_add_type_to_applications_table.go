package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251216102504AddTypeToApplicationsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251216102504AddTypeToApplicationsTable) Signature() string {
	return "20251216102504_add_type_to_applications_table"
}

// Up Run the migrations.
func (r *M20251216102504AddTypeToApplicationsTable) Up() error {
	return facades.Schema().Table("applications", func(table schema.Blueprint) {
		// Application type: 'signup' for new SME registration, 'amend_formalisation' for changes to existing SME
		table.Enum("type", []any{"signup", "amend_formalisation"}).Default("signup")

		// JSON field to store amendment details (current vs proposed values)
		table.Text("data").Nullable()

		// Reference to the SME being amended (only for amend_formalisation type)
		table.UnsignedBigInteger("sme_id").Nullable()

		// Rejection reason when application is rejected
		table.Text("rejection_reason").Nullable()

		// Foreign key to smes table
		table.Foreign("sme_id").References("id").On("smes")
	})
}

// Down Reverse the migrations.
func (r *M20251216102504AddTypeToApplicationsTable) Down() error {
	return facades.Schema().Table("applications", func(table schema.Blueprint) {
		table.DropForeign("applications_sme_id_foreign")
		table.DropColumn("type", "data", "sme_id", "rejection_reason")
	})
}
