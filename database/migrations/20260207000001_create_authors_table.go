package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20260207000001CreateAuthorsTable struct{}

// Signature The unique signature for the migration.
func (r *M20260207000001CreateAuthorsTable) Signature() string {
	return "20260207000001_create_authors_table"
}

// Up Run the migrations.
func (r *M20260207000001CreateAuthorsTable) Up() error {
	if !facades.Schema().HasTable("authors") {
		return facades.Schema().Create("authors", func(table schema.Blueprint) {
			table.ID()
			table.String("first_name", 100)
			table.String("last_name", 100)
			table.Text("bio").Nullable()
			table.String("email", 255).Nullable()
			table.String("website", 255).Nullable()
			table.Date("birth_date").Nullable()
			table.String("nationality", 100).Nullable()
			table.String("photo_url", 500).Nullable()
			table.String("status", 20).Default("ACTIVE")
			table.TimestampsTz()
			table.SoftDeletesTz()
			table.Index("email")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260207000001CreateAuthorsTable) Down() error {
	return facades.Schema().DropIfExists("authors")
}
