package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type CreateMessageMentionsTable struct {
}

// Signature The name and signature of the console command.
func (receiver *CreateMessageMentionsTable) Signature() string {
	return "20250701120002_create_message_mentions_table"
}

// Description The console command description.
func (receiver *CreateMessageMentionsTable) Description() string {
	return "Create message_mentions table for @mentions in messages"
}

// Up Run the migrations.
func (receiver *CreateMessageMentionsTable) Up() error {
	return facades.Schema().Create("message_mentions", func(table schema.Blueprint) {
		table.ID()
		
		// Relationships
		table.UnsignedBigInteger("message_id")
		table.UnsignedBigInteger("user_id")
		
		// Mention metadata
		table.Integer("position") // Position of mention in message
		table.Integer("length")   // Length of mention text
		table.Boolean("is_read").Default(false)
		table.Timestamp("read_at").Nullable()
		
		table.Timestamps()
	})
}

// Down Reverse the migrations.
func (receiver *CreateMessageMentionsTable) Down() error {
	return facades.Schema().DropIfExists("message_mentions")
}