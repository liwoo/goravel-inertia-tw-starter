package database

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/database/seeder"

	"players/database/migrations"
	"players/database/seeders"
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
	}
}

func (kernel Kernel) Seeders() []seeder.Seeder {
	return []seeder.Seeder{
		&seeders.DatabaseSeeder{},
		&seeders.BookSeeder{},
	}
}
