package models

import (
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"
)

type ProcurementNotice struct {
	orm.Model
	orm.SoftDeletes
	ProcuredBy             string          `json:"procured_by" db:"procured_by"`
	ProcurementType        string          `json:"procurement_type" db:"procurement_type"`
	MarketApproach         string          `json:"market_approach" db:"market_approach"`
	Invitation             string          `json:"invitation" db:"invitation"`
	RefNo                  string          `json:"ref_no" db:"ref_no"`
	OpenDate               carbon.DateTime `json:"open_date" db:"open_date"`
	CloseDate              carbon.DateTime `json:"close_date" db:"close_date"`
	Partners               []string        `json:"partners" db:"partners" gorm:"type:json;serializer:json"`
	QualifyingDistricts    []string        `json:"qualifying_districts" db:"qualifying_districts" gorm:"type:json;serializer:json"`
	IsPublished            bool            `json:"is_published" db:"is_published"`
	Organization           string          `json:"organization" db:"organization"`
	Classification         []string        `json:"classification" db:"classification" gorm:"type:json;serializer:json"`
	InterestedSmes         []string        `json:"interested_smes" db:"interested_smes" gorm:"type:json;serializer:json"`
	Details                string          `json:"details" db:"details"`
	ApplicationDetails     string          `json:"application_details" db:"application_details"`
	MinimumQualifyingScore int             `json:"minimum_qualifying_score" db:"minimum_qualifying_score"`
	CreatedBy              *int            `json:"created_by" db:"created_by"`
	UpdatedBy              *int            `json:"updated_by" db:"updated_by"`
	DeletedBy              *int            `json:"deleted_by" db:"deleted_by"`
	IpAddress              *string         `json:"ip_address" db:"ip_address"`
	UserAgent              *string         `json:"user_agent" db:"user_agent"`
}

func (r *ProcurementNotice) TableName() string {
	return "procurement_notices"
}
