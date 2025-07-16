package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type CreateNotificationsTable struct {
}

// Signature The name and signature of the console command.
func (receiver *CreateNotificationsTable) Signature() string {
	return "20250701120003_create_notifications_table"
}

// Description The console command description.
func (receiver *CreateNotificationsTable) Description() string {
	return "Create notifications table for notification system"
}

// Up Run the migrations.
func (receiver *CreateNotificationsTable) Up() error {
	return facades.Schema().Create("notifications", func(table schema.Blueprint) {
		table.ID()
		
		// Core notification fields
		table.String("title")
		table.Text("message").Nullable()
		table.String("type", 50) // message, mention, system, warning, success, etc.
		
		// User relationships
		table.UnsignedBigInteger("user_id")
		table.UnsignedBigInteger("trigger_user_id").Nullable()
		
		// Related entities
		table.String("related_type", 50).Nullable() // message, user, role, etc.
		table.UnsignedBigInteger("related_id").Nullable()
		
		// Notification state
		table.Boolean("is_read").Default(false)
		table.Timestamp("read_at").Nullable()
		table.Boolean("is_dismissed").Default(false)
		table.Timestamp("dismissed_at").Nullable()
		
		// Metadata
		table.Text("data").Nullable() // Additional data as JSON string
		table.String("priority", 20).Default("normal") // high, medium, low, normal
		table.Timestamp("expires_at").Nullable()
		
		table.Timestamps()
		table.SoftDeletes()
	})
}

// Down Reverse the migrations.
func (receiver *CreateNotificationsTable) Down() error {
	return facades.Schema().DropIfExists("notifications")
}