package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251201170115CreateUserActivitiesTable struct{}

// Signature The unique signature for the migration.
func (r *M20251201170115CreateUserActivitiesTable) Signature() string {
	return "20251201170115_create_user_activities_table"
}

// Up Run the migrations.
func (r *M20251201170115CreateUserActivitiesTable) Up() error {
	if !facades.Schema().HasTable("user_activities") {
		err := facades.Schema().Create("user_activities", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("user_id")
			table.String("activity_type", 50)
			table.String("description", 255)
			table.Text("metadata").Nullable()
			table.String("ip_address", 45).Nullable()
			table.Text("user_agent").Nullable()
			table.String("related_type", 50).Nullable()
			table.UnsignedBigInteger("related_id").Nullable()
			table.TimestampsTz()
			table.SoftDeletesTz()

			// Indexes
			table.Index("user_id")
			table.Index("activity_type")
			table.Index("related_type", "related_id")

			// Foreign key constraint
			table.Foreign("user_id").References("id").On("users").CascadeOnDelete()
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251201170115CreateUserActivitiesTable) Down() error {
 	return facades.Schema().DropIfExists("user_activities")
}
