package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251015211802CreateLenders struct{}

// Signature The unique signature for the migration.
func (r *M20251015211802CreateLenders) Signature() string {
	return "20251015211802_create_lenders"
}

// Up Run the migrations.
func (r *M20251015211802CreateLenders) Up() error {
	if !facades.Schema().HasTable("lenders") {
		return facades.Schema().Create("lenders", func(table schema.Blueprint) {
			table.ID()
			table.String("name", 255)
			table.String("email", 255)
			table.String("phone", 50).Nullable()
			table.String("address", 500).Nullable()
			table.Enum("gender", []any{"MALE", "FEMALE"}).Nullable()
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251015211802CreateLenders) Down() error {
	return facades.Schema().DropIfExists("lenders")
}
