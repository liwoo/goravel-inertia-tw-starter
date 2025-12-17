package requests

import (
	"github.com/goravel/framework/support/carbon"
)

// PrimaryBusinessOwner represents the primary owner of an SME in request payload
type PrimaryBusinessOwner struct {
	FirstName              string           `json:"first_name" form:"first_name"`
	LastName               string           `json:"last_name" form:"last_name"`
	OtherNames             *string          `json:"other_names" form:"other_names"`
	Nationality            string           `json:"nationality" form:"nationality"`
	NationalIdNumber       string           `json:"national_id_number" form:"national_id_number"`
	DateOfBirth            *carbon.DateTime `json:"date_of_birth" form:"date_of_birth"`
	Gender                 string           `json:"gender" form:"gender"`
	EducationLevel         string           `json:"education_level" form:"education_level"`
	MalawianStatus         string           `json:"malawian_status" form:"malawian_status"`
	HasSpecialNeeds        bool             `json:"has_special_needs" form:"has_special_needs"`
	PhoneNumber            string           `json:"phone_number" form:"phone_number"`
	LandlineNumber         *string          `json:"landline_number" form:"landline_number"`
	Email                  *string          `json:"email" form:"email"`
	PhysicalAddress        *string          `json:"physical_address" form:"physical_address"`
	PostalAddress          *string          `json:"postal_address" form:"postal_address"`
	Region                 *string          `json:"region" form:"region"`
	District               *string          `json:"district" form:"district"`
	TraditionalAuthority   *string          `json:"traditional_authority" form:"traditional_authority"`
	AltContactName         *string          `json:"alt_contact_name" form:"alt_contact_name"`
	AltContactRelationship *string          `json:"alt_contact_relationship" form:"alt_contact_relationship"`
	AltContactPhone        *string          `json:"alt_contact_phone" form:"alt_contact_phone"`
}

// AdditionalBusinessMember represents additional members of an SME in request payload
type AdditionalBusinessMember struct {
	FirstName        string           `json:"first_name" form:"first_name"`
	LastName         string           `json:"last_name" form:"last_name"`
	OtherNames       *string          `json:"other_names" form:"other_names"`
	Nationality      string           `json:"nationality" form:"nationality"`
	NationalIdNumber string           `json:"national_id_number" form:"national_id_number"`
	DateOfBirth      *carbon.DateTime `json:"date_of_birth" form:"date_of_birth"`
	Email            *string          `json:"email" form:"email"`
	PhoneNumber      string           `json:"phone_number" form:"phone_number"`
	IsIntern         bool             `json:"is_intern" form:"is_intern"`
	IsPartTime       bool             `json:"is_part_time" form:"is_part_time"`
}

// BusinessFormalisation represents business formalization details in request payload
type BusinessFormalisation struct {
	HasBankAccount         bool    `json:"has_bank_account" form:"has_bank_account"`
	HasTaxClarification    bool    `json:"has_tax_clarification" form:"has_tax_clarification"`
	IsRegisteredForVat     bool    `json:"is_registered_for_vat" form:"is_registered_for_vat"`
	IsMemberOfAssociation  bool    `json:"is_member_of_association" form:"is_member_of_association"`
	IsAffiliated           bool    `json:"is_affiliated" form:"is_affiliated"`
	HasExportLicense       bool    `json:"has_export_license" form:"has_export_license"`
	HasAccessedBds         bool    `json:"has_accessed_bds" form:"has_accessed_bds"`
	AnnualTurnover         float64 `json:"annual_turnover" form:"annual_turnover"`
	EstimatedValueOfAssets float64 `json:"estimated_value_of_assets" form:"estimated_value_of_assets"`
	FormalisationScore     int     `json:"formalisation_score" form:"formalisation_score"`
}

// BusinessEmployeeSummary represents employee summary in request payload
type BusinessEmployeeSummary struct {
	FullTimeMales   int `json:"full_time_males" form:"full_time_males"`
	FullTimeFemales int `json:"full_time_females" form:"full_time_females"`
	PartTimeMales   int `json:"part_time_males" form:"part_time_males"`
	PartTimeFemales int `json:"part_time_females" form:"part_time_females"`
	InternMales     int `json:"intern_males" form:"intern_males"`
	InternFemales   int `json:"intern_females" form:"intern_females"`
}
