package models

import (
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"
)

type PrimaryBusinessOwner struct {
	orm.Model
	orm.SoftDeletes
	FirstName              string          `json:"first_name" db:"first_name"`
	LastName               string          `json:"last_name" db:"last_name"`
	OtherNames             *string         `json:"other_names" db:"other_names"`
	Nationality            string          `json:"nationality" db:"nationality"`
	NationalIdNumber       string          `json:"national_id_number" db:"national_id_number"`
	DateOfBirth            carbon.DateTime `json:"date_of_birth" db:"date_of_birth"`
	Gender                 string          `json:"gender" db:"gender"`
	EducationLevel         string          `json:"education_level" db:"education_level"`
	MalawianStatus         string          `json:"malawian_status" db:"malawian_status"`
	HasSpecialNeeds        bool            `json:"has_special_needs" db:"has_special_needs"`
	PhoneNumber            string          `json:"phone_number" db:"phone_number"`
	LandlineNumber         *string         `json:"landline_number" db:"landline_number"`
	Email                  *string         `json:"email" db:"email"`
	PhysicalAddress        *string         `json:"physical_address" db:"physical_address"`
	PostalAddress          *string         `json:"postal_address" db:"postal_address"`
	Region                 *string         `json:"region" db:"region"`
	District               *string         `json:"district" db:"district"`
	TraditionalAuthority   *string         `json:"traditional_authority" db:"traditional_authority"`
	AltContactName         *string         `json:"alt_contact_name" db:"alt_contact_name"`
	AltContactRelationship *string         `json:"alt_contact_relationship" db:"alt_contact_relationship"`
	AltContactPhone        *string         `json:"alt_contact_phone" db:"alt_contact_phone"`
	CreatedBy              *int            `json:"created_by" db:"created_by"`
	UpdatedBy              *int            `json:"updated_by" db:"updated_by"`
	DeletedBy              *int            `json:"deleted_by" db:"deleted_by"`
	IpAddress              *string         `json:"ip_address" db:"ip_address"`
	UserAgent              *string         `json:"user_agent" db:"user_agent"`
	SmeID                  int             `json:"sme_id" db:"sme_id"`
}

func (r *PrimaryBusinessOwner) TableName() string {
	return "primary_business_owner"
}
