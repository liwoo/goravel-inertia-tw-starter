package migrations

import (
	"github.com/goravel/framework/facades"
)

type FixBooksCreatedBy20250124 struct {
}

// Signature The unique signature for the migration.
func (r *FixBooksCreatedBy20250124) Signature() string {
	return "20250124_fix_books_created_by"
}

// Up Run the migrations.
func (r *FixBooksCreatedBy20250124) Up() error {
	// Get the first super admin user
	var superAdminID uint
	err := facades.Orm().Query().Table("users").Where("is_super_admin = ?", true).Select("id").First(&superAdminID)
	if err != nil {
		// If no super admin, get the first user
		err = facades.Orm().Query().Table("users").Select("id").Order("id").First(&superAdminID)
		if err != nil {
			// No users exist, skip the update
			return nil
		}
	}
	
	// Update all books without created_by
	_, err = facades.Orm().Query().Table("books").Where("created_by IS NULL").Update("created_by", superAdminID)
	if err != nil {
		return err
	}
	
	// Also update updated_by
	_, err = facades.Orm().Query().Table("books").Where("updated_by IS NULL").Update("updated_by", superAdminID)
	if err != nil {
		return err
	}
	
	// Also update other tables that might have the same issue
	tables := []string{"users", "roles", "permissions"}
	for _, table := range tables {
		facades.Orm().Query().Table(table).Where("created_by IS NULL").Update("created_by", superAdminID)
		facades.Orm().Query().Table(table).Where("updated_by IS NULL").Update("updated_by", superAdminID)
	}
	
	return nil
}

// Down Reverse the migrations.
func (r *FixBooksCreatedBy20250124) Down() error {
	// This migration is not reversible as we don't know which records had NULL created_by
	return nil
}