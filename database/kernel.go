package database

import (
	"smedi-sme-db/database/migrations"
	"smedi-sme-db/database/seeders"

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
		&migrations.M20251015211802CreateLenders{},
		&migrations.M20251015225803AddAuditFieldsToLendersTable{},
		&migrations.M20251019163714CreateSmeConfig{},
		&migrations.M20251019180624AddAuditFieldsToSmeConfigTable{},
		&migrations.M20251020052629CreateSmeTable{},
		&migrations.M20251020063409AddAuditFieldsToSmesTable{},
		&migrations.M20251020053512CreatePrimaryBusinessOwnerTable{},
		&migrations.M20251020064423AddAuditFieldsToPrimaryBusinessOwnerTable{},
		&migrations.M20251020054645CreateAdditionalBusinessMembersTable{},
		&migrations.M20251020065231AddAuditFieldsToAdditionalBusinessMembersTable{},
		&migrations.M20251020055631AddFkToAdditionalBusinessMembersTable{},
		&migrations.M20251020055818AddFkToPrimaryBusinessOwnerTable{},
		&migrations.M20251020060122CreateBusinessFormalisationTable{},
		&migrations.M20251020070907AddAuditFieldsToBusinessFormalisationTable{},
		&migrations.M20251020061003CreateBusinessEmployeeSummaryTable{},
		&migrations.M20251020071158AddAuditFieldsToBusinessEmployeeSummaryTable{},
		&migrations.M20251020064006AddJsonFieldsToSmesTable{},
		&migrations.M20251020113904AddCascadeDeleteToSmeRelationships{},
		&migrations.M20251020212109AddCodeToSmeConfigTable{},
		&migrations.M20251111212220FixSmeConfigIndexDropOrder{},
		&migrations.M20251123031832AddGenderToAdditionalBusinessMembersTable{},
		&migrations.M20251123050933CreateEventsTable{},
		&migrations.M20251123051707AddAuditFieldsToEventsTable{},
		&migrations.M20251123124316CreateProcurementNoticesTable{},
		&migrations.M20251123130451AddAuditFieldsToProcurementNoticesTable{},
		&migrations.M20251118210300CreateBdspsTable{},
		&migrations.M20251119090154AddAuditFieldsToBdspsTable{},
		&migrations.M20251125232825CreateApplicationsTable{},
		&migrations.M20251126013257AddAuditFieldsToApplicationsTable{},
		&migrations.M20251127085740AddPerformanceIndexes{},
		&migrations.M20251127120000AddIsActiveToSmesTable{},
		&migrations.M20251127151630AlterPhoneColumnsSize{},
		&migrations.M20251128054945AddEventEndDate{},
		&migrations.M20251201170115CreateUserActivitiesTable{},
		&migrations.M20251203100000AddClassificationToSmesTable{},
		&migrations.M20251214150511AlterCurrencyFieldsPrecision{},
		&migrations.M20251214160000AddTotpFieldsToUsersTable{},
		&migrations.M20251214160001CreateTotpBackupCodesTable{},
	}
}
func (kernel Kernel) Seeders() []seeder.Seeder {
	return []seeder.Seeder{
		&seeders.DatabaseSeeder{},
		&seeders.BookSeeder{},
		&seeders.RBACSeeder{},
		&seeders.ConfigSeeder{},
		&seeders.SmeSeeder{},
		&seeders.BdspSeeder{},
	}
}
