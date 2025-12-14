package migrations

import (
	"github.com/goravel/framework/facades"
)

type M20251214150511AlterCurrencyFieldsPrecision struct{}

// Signature The unique signature for the migration.
func (r *M20251214150511AlterCurrencyFieldsPrecision) Signature() string {
	return "20251214150511_alter_currency_fields_precision"
}

// Up Run the migrations.
// Increases currency fields precision from default DECIMAL(8,2) to NUMERIC(18,2)
// to accommodate large Malawian Kwacha values (e.g., 100,000,000 MWK)
func (r *M20251214150511AlterCurrencyFieldsPrecision) Up() error {
	// business_formalisation table - annual_turnover and estimated_value_of_assets
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE business_formalisation ALTER COLUMN annual_turnover TYPE NUMERIC(18,2)",
	); err != nil {
		return err
	}

	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE business_formalisation ALTER COLUMN estimated_value_of_assets TYPE NUMERIC(18,2)",
	); err != nil {
		return err
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251214150511AlterCurrencyFieldsPrecision) Down() error {
	// Revert to default DECIMAL - note: this may fail if data exceeds the smaller precision
	facades.Orm().Query().Exec("ALTER TABLE business_formalisation ALTER COLUMN annual_turnover TYPE DECIMAL(8,2)")
	facades.Orm().Query().Exec("ALTER TABLE business_formalisation ALTER COLUMN estimated_value_of_assets TYPE DECIMAL(8,2)")
	return nil
}
