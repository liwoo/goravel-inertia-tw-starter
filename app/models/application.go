package models

// Application type constants
const (
	ApplicationTypeSignup = "signup"
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
	Data            *string `json:"data" gorm:"column:data"` // JSON field for extra details
	RejectionReason *string `json:"rejection_reason" gorm:"column:rejection_reason"`

	// Core application fields
	SME            string `json:"sme" gorm:"column:sme"`
	RegistrantName string `json:"registrant_name" gorm:"column:registrant_name"`
	Email          string `json:"email" gorm:"column:email"`
	Phone          string `json:"phone" gorm:"column:phone"`
	Status         string `json:"status" gorm:"column:status"`
}

func (Application) TableName() string {
	return "applications"
}
