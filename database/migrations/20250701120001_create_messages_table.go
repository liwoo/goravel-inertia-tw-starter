package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type CreateMessagesTable struct {
}

// Signature The name and signature of the console command.
func (receiver *CreateMessagesTable) Signature() string {
	return "20250701120001_create_messages_table"
}

// Description The console command description.
func (receiver *CreateMessagesTable) Description() string {
	return "Create messages table for messaging system"
}

// Up Run the migrations.
func (receiver *CreateMessagesTable) Up() error {
	return facades.Schema().Create("messages", func(table schema.Blueprint) {
		table.ID()
		
		// Core message fields
		table.Text("content")
		table.String("type", 20).Default("direct") // direct, group, system
		table.String("status", 20).Default("sent") // sent, delivered, read, deleted
		
		// User relationships
		table.UnsignedBigInteger("sender_id")
		table.UnsignedBigInteger("recipient_id").Nullable()
		
		// Group messaging (for future use)
		table.UnsignedBigInteger("group_id").Nullable()
		
		// Message metadata
		table.Boolean("is_edited").Default(false)
		table.Timestamp("edited_at").Nullable()
		table.Timestamp("read_at").Nullable()
		table.Timestamp("delivered_at").Nullable()
		
		// Threading support
		table.UnsignedBigInteger("parent_message_id").Nullable()
		
		// Attachments support
		table.Boolean("has_attachments").Default(false)
		
		table.Timestamps()
		table.SoftDeletes()
	})
}

// Down Reverse the migrations.
func (receiver *CreateMessagesTable) Down() error {
	return facades.Schema().DropIfExists("messages")
}