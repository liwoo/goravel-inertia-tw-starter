package models

import (
	"github.com/goravel/framework/database/orm"
)

type BusinessEmployeeSummary struct {
	orm.Model
	orm.SoftDeletes
	SmeID           int     `json:"sme_id" db:"sme_id"`
	FullTimeMales   int     `json:"full_time_males" db:"full_time_males"`
	FullTimeFemales int     `json:"full_time_females" db:"full_time_females"`
	PartTimeMales   int     `json:"part_time_males" db:"part_time_males"`
	PartTimeFemales int     `json:"part_time_females" db:"part_time_females"`
	InternMales     int     `json:"intern_males" db:"intern_males"`
	InternFemales   int     `json:"intern_females" db:"intern_females"`
	CreatedBy       *int    `json:"created_by" db:"created_by"`
	UpdatedBy       *int    `json:"updated_by" db:"updated_by"`
	DeletedBy       *int    `json:"deleted_by" db:"deleted_by"`
	IpAddress       *string `json:"ip_address" db:"ip_address"`
	UserAgent       *string `json:"user_agent" db:"user_agent"`
}

func (r *BusinessEmployeeSummary) TableName() string {
	return "business_employee_summary"
}
