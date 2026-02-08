package models

import "time"

type UserTenant struct {
	BaseAuditableModel

	UserID   uint      `gorm:"not null;index;uniqueIndex:idx_user_tenant" json:"user_id"`
	TenantID uint      `gorm:"not null;index;uniqueIndex:idx_user_tenant" json:"tenant_id"`
	RoleID   *uint     `gorm:"index" json:"role_id,omitempty"`
	IsActive bool      `gorm:"default:true" json:"is_active"`
	JoinedAt time.Time `gorm:"not null;autoCreateTime" json:"joined_at"`

	// Relations
	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Role   *Role   `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

func (UserTenant) TableName() string {
	return "user_tenants"
}
