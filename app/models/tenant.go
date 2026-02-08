package models

type Tenant struct {
	BaseAuditableModel

	Name        string  `gorm:"not null;type:varchar(255)" json:"name"`
	Slug        string  `gorm:"uniqueIndex;not null;type:varchar(100)" json:"slug"`
	IsActive    bool    `gorm:"default:true;index" json:"is_active"`
	IsMain      bool    `gorm:"default:false;index" json:"is_main"`
	Description *string `gorm:"type:text" json:"description,omitempty"`
	Settings    *string `gorm:"type:jsonb" json:"settings,omitempty"`
	LogoURL     *string `gorm:"type:varchar(500)" json:"logo_url,omitempty"`
}

func (Tenant) TableName() string {
	return "tenants"
}
