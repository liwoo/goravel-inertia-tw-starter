package migrations

import (
	"github.com/goravel/framework/facades"
)

type AddTagsColumnToBooksTable struct{}

func (m *AddTagsColumnToBooksTable) Signature() string {
	return "20250723120000_add_tags_column_to_books_table"
}

func (m *AddTagsColumnToBooksTable) Up() error {
	// Use raw SQL for adding column
	_, err := facades.Orm().Query().Exec("ALTER TABLE books ADD COLUMN tags TEXT NULL")
	return err
}

func (m *AddTagsColumnToBooksTable) Down() error {
	// Use raw SQL for dropping column
	_, err := facades.Orm().Query().Exec("ALTER TABLE books DROP COLUMN tags")
	return err
}