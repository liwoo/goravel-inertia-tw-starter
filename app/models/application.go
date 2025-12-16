package models

// Application type constants
const (
	ApplicationTypeSignup              = "signup"
	ApplicationTypeAmendFormalisation  = "amend_formalisation"
)

// Application status constants
const (
	ApplicationStatusPending  = "Pending"
	ApplicationStatusApproved = "Approved"
	ApplicationStatusRejected = "Rejected"
)

type Application struct {
	BaseAuditableModel

	// Application type and metadata
	Type            string  `json:"type" gorm:"column:type;default:signup"`
	Data            *string `json:"data" gorm:"column:data"`                  // JSON field for amendment details
	SmeID           *uint   `json:"sme_id" gorm:"column:sme_id"`              // Reference to SME for amendments
	RejectionReason *string `json:"rejection_reason" gorm:"column:rejection_reason"`

	// Signup application fields (used when type = 'signup')
	SME                        string  `json:"sme" gorm:"column:sme"`
	RegistrantName             string  `json:"registrant_name" gorm:"column:registrant_name"`
	Email                      string  `json:"email" gorm:"column:email"`
	Phone                      string  `json:"phone" gorm:"column:phone"`
	SMERegistrationNumber      string  `json:"sme_registration_number" gorm:"column:sme_registration_number"`
	SMETaxIdentificationNumber string  `json:"sme_tax_identification_number" gorm:"column:sme_tax_identification_number"`
	Status                     string  `json:"status" gorm:"column:status"`

	// Primary Business Owner Details (used when type = 'signup')
	// Pointer types used for nullable fields to support amendment applications
	FirstName              *string `json:"first_name" gorm:"column:first_name"`
	LastName               *string `json:"last_name" gorm:"column:last_name"`
	OtherNames             *string `json:"other_names" gorm:"column:other_names"`
	Nationality            *string `json:"nationality" gorm:"column:nationality"`
	NationalIDNumber       *string `json:"national_id_number" gorm:"column:national_id_number"`
	DateOfBirth            *string `json:"date_of_birth" gorm:"column:date_of_birth"`
	Gender                 *string `json:"gender" gorm:"column:gender"`
	EducationLevel         *string `json:"education_level" gorm:"column:education_level"`
	MalawianStatus         *string `json:"malawian_status" gorm:"column:malawian_status"`
	HasSpecialNeeds        bool    `json:"has_special_needs" gorm:"column:has_special_needs"`
	LandlineNumber         *string `json:"landline_number" gorm:"column:landline_number"`
	PhysicalAddress        *string `json:"physical_address" gorm:"column:physical_address"`
	PostalAddress          *string `json:"postal_address" gorm:"column:postal_address"`
	Region                 *string `json:"region" gorm:"column:region"`
	District               *string `json:"district" gorm:"column:district"`
	TraditionalAuthority   *string `json:"traditional_authority" gorm:"column:traditional_authority"`
	AltContactName         *string `json:"alt_contact_name" gorm:"column:alt_contact_name"`
	AltContactRelationship *string `json:"alt_contact_relationship" gorm:"column:alt_contact_relationship"`
	AltContactPhone        *string `json:"alt_contact_phone" gorm:"column:alt_contact_phone"`

	// Relations
	Sme *Sme `json:"sme_record,omitempty" gorm:"foreignKey:SmeID"`
}

func (Application) TableName() string {
	return "applications"
}
