package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type CreateBooksTable struct {
}

// Signature The name and signature of the console command.
func (receiver *CreateBooksTable) Signature() string {
	return "20250625000001_create_books_table"
}

// Description The console command description.
func (receiver *CreateBooksTable) Description() string {
	return "Create books table for CRUD example"
}

// Up Run the migrations.
func (receiver *CreateBooksTable) Up() error {
	return facades.Schema().Create("books", func(table schema.Blueprint) {
		table.ID()
		table.String("title")
		table.String("author")
		table.String("isbn")
		table.Text("description")
		table.Float("price", 10, 2)
		table.String("published_at")
		table.String("status").Default("AVAILABLE") // AVAILABLE, BORROWED, MAINTENANCE
		table.Timestamps()
		table.SoftDeletes()
	})
}

// Down Reverse the migrations.
func (receiver *CreateBooksTable) Down() error {
	return facades.Schema().DropIfExists("books")
}