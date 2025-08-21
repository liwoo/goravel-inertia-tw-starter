package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250722104500AddAuditFieldsToAllTables struct {
}

// Signature The unique signature for the migration.
func (r *M20250722104500AddAuditFieldsToAllTables) Signature() string {
	return "20250722104500_add_audit_fields_to_all_tables"
}

// Up Run the migrations.
func (r *M20250722104500AddAuditFieldsToAllTables) Up() error {
	// Add audit fields to users table
	err := facades.Schema().Table("users", func(table schema.Blueprint) {
		// Check if columns don't already exist before adding
		table.UnsignedBigInteger("updated_by").Nullable().Comment("User who last updated this user")
		table.UnsignedBigInteger("deleted_by").Nullable().Comment("User who deleted this user")
		table.String("ip_address", 45).Nullable().Comment("IP address from where the action was performed")
		table.Text("user_agent").Nullable().Comment("User agent from where the action was performed")

		// Add foreign key constraints
		table.Foreign("updated_by").References("id").On("users")
		table.Foreign("deleted_by").References("id").On("users")

		// Add indexes
		table.Index("updated_by")
		table.Index("deleted_by")
	})
	if err != nil {
		return err
	}

	// Add audit fields to books table
	err = facades.Schema().Table("books", func(table schema.Blueprint) {
		table.UnsignedBigInteger("updated_by").Nullable().Comment("User who last updated this book")
		table.UnsignedBigInteger("deleted_by").Nullable().Comment("User who deleted this book")
		table.String("ip_address", 45).Nullable().Comment("IP address from where the action was performed")
		table.Text("user_agent").Nullable().Comment("User agent from where the action was performed")

		table.Foreign("updated_by").References("id").On("users")
		table.Foreign("deleted_by").References("id").On("users")

		table.Index("updated_by")
		table.Index("deleted_by")
	})
	if err != nil {
		return err
	}

	// Add audit fields to roles table
	err = facades.Schema().Table("roles", func(table schema.Blueprint) {
		table.UnsignedBigInteger("updated_by").Nullable().Comment("User who last updated this role")
		table.UnsignedBigInteger("deleted_by").Nullable().Comment("User who deleted this role")
		table.String("ip_address", 45).Nullable().Comment("IP address from where the action was performed")
		table.Text("user_agent").Nullable().Comment("User agent from where the action was performed")

		table.Foreign("updated_by").References("id").On("users")
		table.Foreign("deleted_by").References("id").On("users")

		table.Index("updated_by")
		table.Index("deleted_by")
	})
	if err != nil {
		return err
	}

	// Add audit fields to permissions table
	err = facades.Schema().Table("permissions", func(table schema.Blueprint) {
		table.UnsignedBigInteger("updated_by").Nullable().Comment("User who last updated this permission")
		table.UnsignedBigInteger("deleted_by").Nullable().Comment("User who deleted this permission")
		table.String("ip_address", 45).Nullable().Comment("IP address from where the action was performed")
		table.Text("user_agent").Nullable().Comment("User agent from where the action was performed")

		table.Foreign("updated_by").References("id").On("users")
		table.Foreign("deleted_by").References("id").On("users")

		table.Index("updated_by")
		table.Index("deleted_by")
	})
	if err != nil {
		return err
	}

	// Add audit fields to messages table if it exists
	err = facades.Schema().Table("messages", func(table schema.Blueprint) {
		table.UnsignedBigInteger("created_by").Nullable().Comment("User who created this message")
		table.UnsignedBigInteger("updated_by").Nullable().Comment("User who last updated this message")
		table.UnsignedBigInteger("deleted_by").Nullable().Comment("User who deleted this message")
		table.String("ip_address", 45).Nullable().Comment("IP address from where the action was performed")
		table.Text("user_agent").Nullable().Comment("User agent from where the action was performed")

		table.Foreign("created_by").References("id").On("users")
		table.Foreign("updated_by").References("id").On("users")
		table.Foreign("deleted_by").References("id").On("users")

		table.Index("created_by")
		table.Index("updated_by")
		table.Index("deleted_by")
	})
	if err != nil {
		return err
	}

	// Add audit fields to notifications table if it exists
	err = facades.Schema().Table("notifications", func(table schema.Blueprint) {
		table.UnsignedBigInteger("created_by").Nullable().Comment("User who created this notification")
		table.UnsignedBigInteger("updated_by").Nullable().Comment("User who last updated this notification")
		table.UnsignedBigInteger("deleted_by").Nullable().Comment("User who deleted this notification")
		table.String("ip_address", 45).Nullable().Comment("IP address from where the action was performed")
		table.Text("user_agent").Nullable().Comment("User agent from where the action was performed")

		table.Foreign("created_by").References("id").On("users")
		table.Foreign("updated_by").References("id").On("users")
		table.Foreign("deleted_by").References("id").On("users")

		table.Index("created_by")
		table.Index("updated_by")
		table.Index("deleted_by")
	})

	return err
}

// Down Reverse the migrations.
func (r *M20250722104500AddAuditFieldsToAllTables) Down() error {
	// Remove audit fields from notifications table
	err := facades.Schema().Table("notifications", func(table schema.Blueprint) {
		table.DropForeign("created_by")
		table.DropForeign("updated_by")
		table.DropForeign("deleted_by")

		table.DropIndex("created_by")
		table.DropIndex("updated_by")
		table.DropIndex("deleted_by")

		table.DropColumn("created_by")
		table.DropColumn("updated_by")
		table.DropColumn("deleted_by")
		table.DropColumn("ip_address")
		table.DropColumn("user_agent")
	})
	if err != nil {
		return err
	}

	// Remove audit fields from messages table
	err = facades.Schema().Table("messages", func(table schema.Blueprint) {
		table.DropForeign("created_by")
		table.DropForeign("updated_by")
		table.DropForeign("deleted_by")

		table.DropIndex("created_by")
		table.DropIndex("updated_by")
		table.DropIndex("deleted_by")

		table.DropColumn("created_by")
		table.DropColumn("updated_by")
		table.DropColumn("deleted_by")
		table.DropColumn("ip_address")
		table.DropColumn("user_agent")
	})
	if err != nil {
		return err
	}

	// Remove audit fields from permissions table
	err = facades.Schema().Table("permissions", func(table schema.Blueprint) {
		table.DropForeign("updated_by")
		table.DropForeign("deleted_by")

		table.DropIndex("updated_by")
		table.DropIndex("deleted_by")

		table.DropColumn("updated_by")
		table.DropColumn("deleted_by")
		table.DropColumn("ip_address")
		table.DropColumn("user_agent")
	})
	if err != nil {
		return err
	}

	// Remove audit fields from roles table
	err = facades.Schema().Table("roles", func(table schema.Blueprint) {
		table.DropForeign("updated_by")
		table.DropForeign("deleted_by")

		table.DropIndex("updated_by")
		table.DropIndex("deleted_by")

		table.DropColumn("updated_by")
		table.DropColumn("deleted_by")
		table.DropColumn("ip_address")
		table.DropColumn("user_agent")
	})
	if err != nil {
		return err
	}

	// Remove audit fields from books table
	err = facades.Schema().Table("books", func(table schema.Blueprint) {
		table.DropForeign("updated_by")
		table.DropForeign("deleted_by")

		table.DropIndex("updated_by")
		table.DropIndex("deleted_by")

		table.DropColumn("updated_by")
		table.DropColumn("deleted_by")
		table.DropColumn("ip_address")
		table.DropColumn("user_agent")
	})
	if err != nil {
		return err
	}

	// Remove audit fields from users table
	err = facades.Schema().Table("users", func(table schema.Blueprint) {
		table.DropForeign("updated_by")
		table.DropForeign("deleted_by")

		table.DropIndex("updated_by")
		table.DropIndex("deleted_by")

		table.DropColumn("updated_by")
		table.DropColumn("deleted_by")
		table.DropColumn("ip_address")
		table.DropColumn("user_agent")
	})

	return err
}
