package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type UpdateBooksPublishedAtToDatetime struct {
}

// Signature The name and signature of the console command.
func (receiver *UpdateBooksPublishedAtToDatetime) Signature() string {
	return "20250724000000_update_books_published_at_to_datetime"
}

// Description The console command description.
func (receiver *UpdateBooksPublishedAtToDatetime) Description() string {
	return "Change published_at column from string to datetime in books table"
}

// Up Run the migrations.
func (receiver *UpdateBooksPublishedAtToDatetime) Up() error {
	// In test environment with SQLite in-memory, skip this migration
	if facades.Config().GetString("database.default") == "sqlite" && 
	   facades.Config().GetString("database.connections.sqlite.database") == ":memory:" {
		return nil
	}
	
	return facades.Schema().Table("books", func(table schema.Blueprint) {
		// Drop the old string column
		table.DropColumn("published_at")
		// Add new datetime column
		table.DateTime("published_at").Nullable()
	})
}

// Down Reverse the migrations.
func (receiver *UpdateBooksPublishedAtToDatetime) Down() error {
	// In test environment with SQLite in-memory, skip this migration
	if facades.Config().GetString("database.default") == "sqlite" && 
	   facades.Config().GetString("database.connections.sqlite.database") == ":memory:" {
		return nil
	}
	
	return facades.Schema().Table("books", func(table schema.Blueprint) {
		// Drop the datetime column
		table.DropColumn("published_at")
		// Add back the string column
		table.String("published_at").Nullable()
	})
}