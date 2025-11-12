package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251020054645CreateAdditionalBusinessMembersTable struct{}

// Signature The unique signature for the migration.
func (r *M20251020054645CreateAdditionalBusinessMembersTable) Signature() string {
	return "20251020054645_create_additional_business_members_table"
}

// Up Run the migrations.
func (r *M20251020054645CreateAdditionalBusinessMembersTable) Up() error {
	if !facades.Schema().HasTable("additional_business_members") {
		return facades.Schema().Create("additional_business_members", func(table schema.Blueprint) {
			table.ID()
			table.String("first_name", 100)
			table.String("last_name", 100)
			table.String("other_names", 100).Nullable()
			table.String("nationality", 100)
			table.String("national_id_number", 50)
			table.Date("date_of_birth").Nullable()
			table.String("email", 100).Nullable()
			table.String("phone_number", 12)
			table.Boolean("is_intern").Default(false)
			table.Boolean("is_part_time").Default(false)
			table.UnsignedBigInteger("sme_id")
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251020054645CreateAdditionalBusinessMembersTable) Down() error {
	return facades.Schema().DropIfExists("additional_business_members")
}
