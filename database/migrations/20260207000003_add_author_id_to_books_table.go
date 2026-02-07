package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20260207000003AddAuthorIdToBooksTable struct{}

// Signature The unique signature for the migration.
func (r *M20260207000003AddAuthorIdToBooksTable) Signature() string {
	return "20260207000003_add_author_id_to_books_table"
}

// Up Run the migrations.
func (r *M20260207000003AddAuthorIdToBooksTable) Up() error {
	return facades.Schema().Table("books", func(table schema.Blueprint) {
		table.UnsignedBigInteger("author_id").Nullable().Comment("Foreign key to authors table")
		table.Foreign("author_id").References("id").On("authors")
		table.Index("author_id")
	})
}

// Down Reverse the migrations.
func (r *M20260207000003AddAuthorIdToBooksTable) Down() error {
	return facades.Schema().Table("books", func(table schema.Blueprint) {
		table.DropForeign("author_id")
		table.DropIndex("author_id")
		table.DropColumn("author_id")
	})
}
