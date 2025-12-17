package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251216210411AddContractFieldsToBusinessEmployeeSummary struct{}

// Signature The unique signature for the migration.
func (r *M20251216210411AddContractFieldsToBusinessEmployeeSummary) Signature() string {
	return "20251216210411_add_contract_fields_to_business_employee_summary"
}

// Up Run the migrations.
func (r *M20251216210411AddContractFieldsToBusinessEmployeeSummary) Up() error {
	return facades.Schema().Table("business_employee_summary", func(table schema.Blueprint) {
		table.Integer("full_time_with_contract_males").Default(0)
		table.Integer("full_time_with_contract_females").Default(0)
		table.Integer("temporary_males").Default(0)
		table.Integer("temporary_females").Default(0)
	})
}

// Down Reverse the migrations.
func (r *M20251216210411AddContractFieldsToBusinessEmployeeSummary) Down() error {
	return facades.Schema().Table("business_employee_summary", func(table schema.Blueprint) {
		table.DropColumn("full_time_with_contract_males")
		table.DropColumn("full_time_with_contract_females")
		table.DropColumn("temporary_males")
		table.DropColumn("temporary_females")
	})
}
