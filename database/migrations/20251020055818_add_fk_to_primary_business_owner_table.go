package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251020055818AddFkToPrimaryBusinessOwnerTable struct{}

// Signature The unique signature for the migration.
func (r *M20251020055818AddFkToPrimaryBusinessOwnerTable) Signature() string {
	return "20251020055818_add_fk_to_primary_business_owner_table"
}

// Up Run the migrations.
func (r *M20251020055818AddFkToPrimaryBusinessOwnerTable) Up() error {
	return facades.Schema().Table("primary_business_owner", func(table schema.Blueprint) {
		table.UnsignedBigInteger("sme_id")
		table.Foreign("sme_id").References("id").On("smes")
	})
}

// Down Reverse the migrations.
func (r *M20251020055818AddFkToPrimaryBusinessOwnerTable) Down() error {
	return nil
}
