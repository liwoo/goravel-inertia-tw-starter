package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251020061003CreateBusinessEmployeeSummaryTable struct{}

// Signature The unique signature for the migration.
func (r *M20251020061003CreateBusinessEmployeeSummaryTable) Signature() string {
	return "20251020061003_create_business_employee_summary_table"
}

// Up Run the migrations.
func (r *M20251020061003CreateBusinessEmployeeSummaryTable) Up() error {
	if !facades.Schema().HasTable("business_employee_summary") {
		return facades.Schema().Create("business_employee_summary", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("sme_id")
			table.Foreign("sme_id").References("id").On("smes")
			table.Integer("full_time_males").Default(0)
			table.Integer("full_time_females").Default(0)
			table.Integer("part_time_males").Default(0)
			table.Integer("part_time_females").Default(0)
			table.Integer("intern_males").Default(0)
			table.Integer("intern_females").Default(0)
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251020061003CreateBusinessEmployeeSummaryTable) Down() error {
	return facades.Schema().DropIfExists("business_employee_summary")
}
