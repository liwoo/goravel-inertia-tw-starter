package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M2026020716315003AddTenantIdToAllTables struct{}

func (r *M2026020716315003AddTenantIdToAllTables) Signature() string {
	return "2026020716315003_add_tenant_id_to_all_tables"
}

func (r *M2026020716315003AddTenantIdToAllTables) Up() error {
	tables := []string{
		"users", "books", "lenders", "applications",
		"configs", "messages", "notifications", "user_activities",
		"authors",
	}

	for _, tableName := range tables {
		if facades.Schema().HasTable(tableName) && !facades.Schema().HasColumn(tableName, "tenant_id") {
			if err := facades.Schema().Table(tableName, func(table schema.Blueprint) {
				table.UnsignedBigInteger("tenant_id").Nullable()
				table.Index("tenant_id")
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *M2026020716315003AddTenantIdToAllTables) Down() error {
	tables := []string{
		"users", "books", "lenders", "applications",
		"configs", "messages", "notifications", "user_activities",
		"authors",
	}

	for _, tableName := range tables {
		if facades.Schema().HasTable(tableName) && facades.Schema().HasColumn(tableName, "tenant_id") {
			facades.Schema().Table(tableName, func(table schema.Blueprint) {
				table.DropColumn("tenant_id")
			})
		}
	}

	return nil
}
