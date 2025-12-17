package models

import (
	"github.com/goravel/framework/database/orm"
)

type BusinessFormalisation struct {
	orm.Model
	orm.SoftDeletes
	SmeID                  int     `json:"sme_id" db:"sme_id"`
	HasBankAccount         bool    `json:"has_bank_account" db:"has_bank_account"`
	HasTaxClarification    bool    `json:"has_tax_clarification" db:"has_tax_clarification"`
	IsRegisteredForVat     bool    `json:"is_registered_for_vat" db:"is_registered_for_vat"`
	IsMemberOfAssociation  bool    `json:"is_member_of_association" db:"is_member_of_association"`
	IsAffiliated           bool    `json:"is_affiliated" db:"is_affiliated"`
	HasExportLicense       bool    `json:"has_export_license" db:"has_export_license"`
	HasAccessedBds         bool    `json:"has_accessed_bds" db:"has_accessed_bds"`
	AnnualTurnover         float64 `json:"annual_turnover" db:"annual_turnover"`
	EstimatedValueOfAssets float64 `json:"estimated_value_of_assets" db:"estimated_value_of_assets"`
	FormalisationScore     int     `json:"formalisation_score" db:"formalisation_score"`
	ComplianceScore        int     `json:"compliance_score" db:"compliance_score"`
	TeamStructureScore     int     `json:"team_structure_score" db:"team_structure_score"`
	FinancialScore         int     `json:"financial_score" db:"financial_score"`
	CreatedBy              *int    `json:"created_by" db:"created_by"`
	UpdatedBy              *int    `json:"updated_by" db:"updated_by"`
	DeletedBy              *int    `json:"deleted_by" db:"deleted_by"`
	IpAddress              *string `json:"ip_address" db:"ip_address"`
	UserAgent              *string `json:"user_agent" db:"user_agent"`
}

func (r *BusinessFormalisation) TableName() string {
	return "business_formalisation"
}
