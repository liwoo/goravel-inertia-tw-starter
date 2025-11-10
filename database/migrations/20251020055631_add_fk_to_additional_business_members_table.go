package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251020055631AddFkToAdditionalBusinessMembersTable struct{}

// Signature The unique signature for the migration.
func (r *M20251020055631AddFkToAdditionalBusinessMembersTable) Signature() string {
	return "20251020055631_add_fk_to_additional_business_members_table"
}

// Up Run the migrations.
func (r *M20251020055631AddFkToAdditionalBusinessMembersTable) Up() error {
	return facades.Schema().Table("additional_business_members", func(table schema.Blueprint) {
		table.Foreign("sme_id").References("id").On("smes")
	})
}

// Down Reverse the migrations.
func (r *M20251020055631AddFkToAdditionalBusinessMembersTable) Down() error {
	return nil
}
