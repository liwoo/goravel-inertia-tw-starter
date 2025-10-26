package database

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/database/seeder"

	"smedi-sme-db/database/migrations"
	"smedi-sme-db/database/seeders"
)

type Kernel struct {
}

func (kernel Kernel) Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20240915060148CreateUsersTable{},
		&migrations.M20250626015224AddUserFieldsToUsersTable{},
		&migrations.CreateBooksTable{},
		&migrations.M20250626020324CreateRolesTable{},
		&migrations.M20250626020332CreatePermissionsTable{},
		&migrations.M20250626020339CreateUserRolesTable{},
		&migrations.M20250626020345CreateRolePermissionsTable{},
		&migrations.M20250628091858AddIsSuperAdminToUsersTable{},
		&migrations.CreateMessagesTable{},
		&migrations.CreateMessageMentionsTable{},
		&migrations.CreateNotificationsTable{},
		&migrations.M20250722065931AddConstraintsToRolesTable{},
		&migrations.M20250722073118AddScopeToPermissionsTable{},
		&migrations.M20250722073324AddCreatedByToResources{},
		&migrations.M20250722104500AddAuditFieldsToAllTables{},
		&migrations.M20250723000000AddScopeToRolePermissionsTable{},
		&migrations.UpdateBooksPublishedAtToDatetime{},
		&migrations.RestoreBookPublishedDates{},
		&migrations.M20251015211802CreateLenders{},
		&migrations.M20251015225803AddAuditFieldsToLendersTable{},
		&migrations.M20251019163714CreateSmeConfig{},
		&migrations.M20251019180624AddAuditFieldsToSmeConfigTable{},
	}
}

func (kernel Kernel) Seeders() []seeder.Seeder {
	return []seeder.Seeder{
		&seeders.DatabaseSeeder{},
		&seeders.BookSeeder{},
		&seeders.RBACSeeder{},
		&seeders.ConfigSeeder{},
	}
}
