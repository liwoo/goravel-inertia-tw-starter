package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250722065931AddConstraintsToRolesTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250722065931AddConstraintsToRolesTable) Signature() string {
	return "20250722065931_add_constraints_to_roles_table"
}

// Up Run the migrations.
func (r *M20250722065931AddConstraintsToRolesTable) Up() error {
	return facades.Schema().Table("roles", func(table schema.Blueprint) {

	})
}

// Down Reverse the migrations.
func (r *M20250722065931AddConstraintsToRolesTable) Down() error {
	return nil
}
