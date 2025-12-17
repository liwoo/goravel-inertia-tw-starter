package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type AddScoreComponentsToBusinessFormalisation struct {
}

// Signature The unique signature for the migration.
func (r *AddScoreComponentsToBusinessFormalisation) Signature() string {
	return "20251217100000_add_score_components_to_business_formalisation"
}

// Up Run the migrations.
func (r *AddScoreComponentsToBusinessFormalisation) Up() error {
	return facades.Schema().Table("business_formalisation", func(table schema.Blueprint) {
		// Add individual score components
		table.Integer("compliance_score").Default(0).Comment("Score from compliance checkboxes (0-70)")
		table.Integer("team_structure_score").Default(0).Comment("Score from team structure (0-20)")
		table.Integer("financial_score").Default(0).Comment("Score from financial data (0-10)")
	})
}

// Down Reverse the migrations.
func (r *AddScoreComponentsToBusinessFormalisation) Down() error {
	return facades.Schema().Table("business_formalisation", func(table schema.Blueprint) {
		table.DropColumn("compliance_score")
		table.DropColumn("team_structure_score")
		table.DropColumn("financial_score")
	})
}
