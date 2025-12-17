package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251216210410AddSpecialNeedsDescriptionToPrimaryBusinessOwner struct{}

// Signature The unique signature for the migration.
func (r *M20251216210410AddSpecialNeedsDescriptionToPrimaryBusinessOwner) Signature() string {
	return "20251216210410_add_special_needs_description_to_primary_business_owner"
}

// Up Run the migrations.
func (r *M20251216210410AddSpecialNeedsDescriptionToPrimaryBusinessOwner) Up() error {
	return facades.Schema().Table("primary_business_owner", func(table schema.Blueprint) {
		table.Text("special_needs_description").Nullable()
	})
}

// Down Reverse the migrations.
func (r *M20251216210410AddSpecialNeedsDescriptionToPrimaryBusinessOwner) Down() error {
	return facades.Schema().Table("primary_business_owner", func(table schema.Blueprint) {
		table.DropColumn("special_needs_description")
	})
}
