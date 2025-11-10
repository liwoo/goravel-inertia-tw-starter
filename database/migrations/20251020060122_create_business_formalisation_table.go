package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251020060122CreateBusinessFormalisationTable struct{}

// Signature The unique signature for the migration.
func (r *M20251020060122CreateBusinessFormalisationTable) Signature() string {
	return "20251020060122_create_business_formalisation_table"
}

// Up Run the migrations.
func (r *M20251020060122CreateBusinessFormalisationTable) Up() error {
	if !facades.Schema().HasTable("business_formalisation") {
		return facades.Schema().Create("business_formalisation", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("sme_id")
			table.Foreign("sme_id").References("id").On("smes")
			table.Boolean("has_bank_account").Default(false)
			table.Boolean("has_tax_clarification").Default(false)
			table.Boolean("is_registered_for_vat").Default(false)
			table.Boolean("is_member_of_association").Default(false)
			table.Boolean("is_affiliated").Default(false)
			table.Boolean("has_export_license").Default(false)
			table.Boolean("has_accessed_bds").Default(false)
			table.Decimal("annual_turnover").Default(0.00)
			table.Decimal("estimated_value_of_assets").Default(0.00)
			table.Integer("formalisation_score").Default(0)
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251020060122CreateBusinessFormalisationTable) Down() error {
	return facades.Schema().DropIfExists("business_formalisation")
}
