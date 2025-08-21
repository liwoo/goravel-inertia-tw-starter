package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type AddTagsColumnToBooksTable20250802235000 struct {
}

// Signature The name and signature of the migration.
func (r *AddTagsColumnToBooksTable20250802235000) Signature() string {
	return "20250802_235000_add_tags_column_to_books_table"
}

// Up Run the migrations.
func (r *AddTagsColumnToBooksTable20250802235000) Up() error {
	return facades.Schema().Table("books", func(table schema.Blueprint) {
		table.Text("tags").Nullable().Comment("JSON array of tags for the book")
	})
}

// Down Reverse the migrations.
func (r *AddTagsColumnToBooksTable20250802235000) Down() error {
	return facades.Schema().Table("books", func(table schema.Blueprint) {
		table.DropColumn("tags")
	})
}