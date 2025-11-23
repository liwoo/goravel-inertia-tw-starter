package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251123031832AddGenderToAdditionalBusinessMembersTable struct{}

// Signature The unique signature for the migration.
func (r *M20251123031832AddGenderToAdditionalBusinessMembersTable) Signature() string {
	return "20251123031832_add_gender_to_additional_business_members_table"
}

// Up Run the migrations.
func (r *M20251123031832AddGenderToAdditionalBusinessMembersTable) Up() error {
	return facades.Schema().Table("additional_business_members", func(table schema.Blueprint) {
		table.String("gender", 10).Nullable()
	})
}

// Down Reverse the migrations.
func (r *M20251123031832AddGenderToAdditionalBusinessMembersTable) Down() error {
	return facades.Schema().Table("additional_business_members", func(table schema.Blueprint) {
		table.DropColumn("gender")
	})
}
