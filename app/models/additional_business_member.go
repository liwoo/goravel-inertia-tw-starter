package models

import (
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"
)

type AdditionalBusinessMember struct {
	orm.Model
	orm.SoftDeletes
	FirstName        string           `json:"first_name" db:"first_name"`
	LastName         string           `json:"last_name" db:"last_name"`
	OtherNames       *string          `json:"other_names" db:"other_names"`
	Gender           *string          `json:"gender" db:"gender"`
	Nationality      string           `json:"nationality" db:"nationality"`
	NationalIdNumber string           `json:"national_id_number" db:"national_id_number"`
	DateOfBirth      *carbon.DateTime `json:"date_of_birth" db:"date_of_birth"`
	Email            *string          `json:"email" db:"email"`
	PhoneNumber      string           `json:"phone_number" db:"phone_number"`
	IsIntern         bool             `json:"is_intern" db:"is_intern"`
	IsPartTime       bool             `json:"is_part_time" db:"is_part_time"`
	SmeId            int              `json:"sme_id" db:"sme_id"`
	CreatedBy        *int             `json:"created_by" db:"created_by"`
	UpdatedBy        *int             `json:"updated_by" db:"updated_by"`
	DeletedBy        *int             `json:"deleted_by" db:"deleted_by"`
	IpAddress        *string          `json:"ip_address" db:"ip_address"`
	UserAgent        *string          `json:"user_agent" db:"user_agent"`
}

func (r *AdditionalBusinessMember) TableName() string {
	return "additional_business_members"
}
