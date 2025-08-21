package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250722073324AddCreatedByToResources struct {
}

// Signature The unique signature for the migration.
func (r *M20250722073324AddCreatedByToResources) Signature() string {
	return "20250722073324_add_created_by_to_resources"
}

// Up Run the migrations.
func (r *M20250722073324AddCreatedByToResources) Up() error {
	// Add created_by to books table
	err := facades.Schema().Table("books", func(table schema.Blueprint) {
		table.UnsignedBigInteger("created_by").Nullable().Comment("User who created this book")
		table.Foreign("created_by").References("id").On("users")
		table.Index("created_by")
	})
	if err != nil {
		return err
	}

	// Add created_by to users table (for tracking who created other users)
	err = facades.Schema().Table("users", func(table schema.Blueprint) {
		table.UnsignedBigInteger("created_by").Nullable().Comment("User who created this user")
		table.Foreign("created_by").References("id").On("users")
		table.Index("created_by")
	})
	if err != nil {
		return err
	}

	// Add created_by to roles table
	err = facades.Schema().Table("roles", func(table schema.Blueprint) {
		table.UnsignedBigInteger("created_by").Nullable().Comment("User who created this role")
		table.Foreign("created_by").References("id").On("users")
		table.Index("created_by")
	})
	if err != nil {
		return err
	}

	// Add created_by to permissions table
	err = facades.Schema().Table("permissions", func(table schema.Blueprint) {
		table.UnsignedBigInteger("created_by").Nullable().Comment("User who created this permission")
		table.Foreign("created_by").References("id").On("users")
		table.Index("created_by")
	})

	return err
}

// Down Reverse the migrations.
func (r *M20250722073324AddCreatedByToResources) Down() error {
	// Remove from permissions
	err := facades.Schema().Table("permissions", func(table schema.Blueprint) {
		table.DropForeign("created_by")
		table.DropIndex("created_by")
		table.DropColumn("created_by")
	})
	if err != nil {
		return err
	}

	// Remove from roles
	err = facades.Schema().Table("roles", func(table schema.Blueprint) {
		table.DropForeign("created_by")
		table.DropIndex("created_by")
		table.DropColumn("created_by")
	})
	if err != nil {
		return err
	}

	// Remove from users
	err = facades.Schema().Table("users", func(table schema.Blueprint) {
		table.DropForeign("created_by")
		table.DropIndex("created_by")
		table.DropColumn("created_by")
	})
	if err != nil {
		return err
	}

	// Remove from books
	err = facades.Schema().Table("books", func(table schema.Blueprint) {
		table.DropForeign("created_by")
		table.DropIndex("created_by")
		table.DropColumn("created_by")
	})

	return err
}
