package models

import (
	"encoding/json"

	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"
	"gorm.io/gorm"
)

// SME Classification Constants based on Malawi MSME Policy
const (
	// Classification categories
	ClassificationMicro        = "Micro"
	ClassificationSmall        = "Small"
	ClassificationMedium       = "Medium"
	ClassificationUnclassified = "Unclassified"

	// Employee thresholds
	MicroEmployeeMin  = 1
	MicroEmployeeMax  = 4
	SmallEmployeeMin  = 5
	SmallEmployeeMax  = 20
	MediumEmployeeMin = 21
	MediumEmployeeMax = 99

	// Annual Turnover thresholds (MWK)
	MicroTurnoverMax  = 5000000.0    // Up to 5,000,000
	SmallTurnoverMin  = 5000000.0    // Above 5,000,000
	SmallTurnoverMax  = 50000000.0   // Up to 50,000,000
	MediumTurnoverMin = 50000000.0   // Above 50,000,000
	MediumTurnoverMax = 500000000.0  // Up to 500,000,000

	// Maximum Assets thresholds (MWK)
	MicroAssetsMax  = 1000000.0   // 1,000,000
	SmallAssetsMax  = 20000000.0  // 20,000,000
	MediumAssetsMax = 250000000.0 // 250,000,000
)

type Sme struct {
	orm.Model
	orm.SoftDeletes
	UsmeNumber                string                     `json:"usme_number" db:"usme_number"`
	Name                      string                     `json:"name" db:"name"`
	IsActive                  bool                       `json:"is_active" db:"is_active" gorm:"default:true"`
	Classification            string                     `json:"classification" db:"classification" gorm:"default:Unclassified"`
	RegistrationNumber        *string                    `json:"registration_number" db:"registration_number"`
	TaxIdentificationNumber   *string                    `json:"tax_identification_number" db:"tax_identification_number"`
	OperationalStartDate      *carbon.DateTime           `json:"operational_start_date" db:"operational_start_date"`
	BusinessCategory          string                     `json:"business_category" db:"business_category"`
	Sector                    string                     `json:"sector" db:"sector"`
	SubSector                 *string                    `json:"sub_sector" db:"sub_sector"`
	BusinessDescription       *string                    `json:"business_description" db:"business_description"`
	ContactPhone              string                     `json:"contact_phone" db:"contact_phone"`
	ContactEmail              string                     `json:"contact_email" db:"contact_email"`
	PhysicalAddress           *string                    `json:"physical_address" db:"physical_address"`
	PostalAddress             *string                    `json:"postal_address" db:"postal_address"`
	Website                   *string                    `json:"website" db:"website"`
	Region                    *string                    `json:"region" db:"region"`
	District                  *string                    `json:"district" db:"district"`
	TraditionalAuthority      *string                    `json:"traditional_authority" db:"traditional_authority"`
	CreatedBy                 *int                       `json:"created_by" db:"created_by"`
	UpdatedBy                 *int                       `json:"updated_by" db:"updated_by"`
	DeletedBy                 *int                       `json:"deleted_by" db:"deleted_by"`
	IpAddress                 *string                    `json:"ip_address" db:"ip_address"`
	UserAgent                 *string                    `json:"user_agent" db:"user_agent"`
	PrimaryBusinessOwner      *PrimaryBusinessOwner      `json:"primary_business_owner" db:"primary_business_owner"`
	AdditionalBusinessMembers []AdditionalBusinessMember `json:"additional_business_members" db:"additional_business_members"`
	BusinessFormalisation     *BusinessFormalisation     `json:"business_formalisation" db:"business_formalisation"`
	BusinessEmployeeSummary   *BusinessEmployeeSummary   `json:"business_employee_summary" db:"business_employee_summary"`
	//business improvement aspects
	BusinessImprovementAspectJSON string   `json:"-" db:"business_improvement_aspect_json" gorm:"column:business_improvement_aspect_json"`
	BusinessImprovementAspects    []string `json:"business_improvement_aspects" gorm:"-:all"`

	//business accessed financing
	BusinessAccessedFinancingJSON string   `json:"-" db:"business_accessed_financing_json" gorm:"column:business_accessed_financing_json"`
	BusinessAccessedFinancing     []string `json:"business_accessed_financing" gorm:"-:all"`
}

func (r *Sme) TableName() string {
	return "smes"
}

// BeforeSave hook to convert arrays to JSON
func (s *Sme) BeforeSave(tx *gorm.DB) error {
	// Convert BusinessImprovementAspects to JSON
	if len(s.BusinessImprovementAspects) > 0 {
		aspectsBytes, err := json.Marshal(s.BusinessImprovementAspects)
		if err != nil {
			return err
		}
		s.BusinessImprovementAspectJSON = string(aspectsBytes)
	} else {
		s.BusinessImprovementAspectJSON = ""
	}

	// Convert BusinessAccessedFinancing to JSON
	if len(s.BusinessAccessedFinancing) > 0 {
		financingBytes, err := json.Marshal(s.BusinessAccessedFinancing)
		if err != nil {
			return err
		}
		s.BusinessAccessedFinancingJSON = string(financingBytes)
	} else {
		s.BusinessAccessedFinancingJSON = ""
	}

	return nil
}

// AfterFind hook to convert JSON to arrays
func (s *Sme) AfterFind(tx *gorm.DB) error {
	// Convert BusinessImprovementAspectJSON to array
	if s.BusinessImprovementAspectJSON != "" {
		err := json.Unmarshal([]byte(s.BusinessImprovementAspectJSON), &s.BusinessImprovementAspects)
		if err != nil {
			// If unmarshal fails, treat as empty array
			s.BusinessImprovementAspects = []string{}
		}
	} else {
		s.BusinessImprovementAspects = []string{}
	}

	// Convert BusinessAccessedFinancingJSON to array
	if s.BusinessAccessedFinancingJSON != "" {
		err := json.Unmarshal([]byte(s.BusinessAccessedFinancingJSON), &s.BusinessAccessedFinancing)
		if err != nil {
			// If unmarshal fails, treat as empty array
			s.BusinessAccessedFinancing = []string{}
		}
	} else {
		s.BusinessAccessedFinancing = []string{}
	}

	return nil
}
