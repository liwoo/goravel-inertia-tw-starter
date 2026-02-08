package database

import (
	"books-database/database/migrations"
	"books-database/database/seeders"

	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/database/seeder"
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
		&migrations.CreateConfigsTable{},
		&migrations.M20251015211802CreateLenders{},
		&migrations.M20251015225803AddAuditFieldsToLendersTable{},
		&migrations.M20251125232825CreateApplicationsTable{},
		&migrations.M20251126013257AddAuditFieldsToApplicationsTable{},
		&migrations.M20251127085740AddPerformanceIndexes{},
		&migrations.M20251201170115CreateUserActivitiesTable{},
		&migrations.M20251214160000AddTotpFieldsToUsersTable{},
		&migrations.M20251214160001CreateTotpBackupCodesTable{},
		&migrations.M20251216102504AddTypeToApplicationsTable{},
		&migrations.M20251216143111AlterApplicationsNullableFields{},
		&migrations.M20260207000001CreateAuthorsTable{},
		&migrations.M20260207000002AddAuditFieldsToAuthorsTable{},
		&migrations.M20260207000003AddAuthorIdToBooksTable{},
		&migrations.M2026020716315001CreateTenantsTable{},
		&migrations.M2026020716315002CreateUserTenantsTable{},
		&migrations.M2026020716315003AddTenantIdToAllTables{},
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
