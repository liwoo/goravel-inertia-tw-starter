package models

import (
	"time"
	"github.com/goravel/framework/database/orm"
)

// Auditable provides audit fields for tracking resource lifecycle
type Auditable struct {
	// Standard GORM fields
	ID        uint                  `gorm:"primarykey" json:"id"`
	CreatedAt orm.Carbon            `json:"created_at"`
	UpdatedAt orm.Carbon            `json:"updated_at"`
	DeletedAt orm.DeletedAt         `gorm:"index" json:"deleted_at,omitempty"`
	
	// Creation tracking
	CreatedBy   *uint     `gorm:"index" json:"created_by,omitempty"`
	Creator     *User     `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	
	// Update tracking  
	UpdatedBy   *uint     `gorm:"index" json:"updated_by,omitempty"`
	Updater     *User     `gorm:"foreignKey:UpdatedBy" json:"updater,omitempty"`
	
	// Soft delete tracking
	DeletedBy   *uint     `gorm:"index" json:"deleted_by,omitempty"`
	Deleter     *User     `gorm:"foreignKey:DeletedBy" json:"deleter,omitempty"`
	
	// Additional audit fields
	IPAddress   string    `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent   string    `gorm:"type:text" json:"user_agent,omitempty"`
}

// AuditableInterface defines methods that auditable models should implement
type AuditableInterface interface {
	GetCreatedBy() *uint
	SetCreatedBy(userID *uint)
	GetUpdatedBy() *uint  
	SetUpdatedBy(userID *uint)
	GetDeletedBy() *uint
	SetDeletedBy(userID *uint)
	GetAuditInfo() *AuditInfo
}

// AuditInfo holds consolidated audit information
type AuditInfo struct {
	CreatedBy *uint `json:"created_by"`
	UpdatedBy *uint `json:"updated_by"`
	DeletedBy *uint `json:"deleted_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// Implement AuditableInterface methods
func (a *Auditable) GetCreatedBy() *uint {
	return a.CreatedBy
}

func (a *Auditable) SetCreatedBy(userID *uint) {
	a.CreatedBy = userID
}

func (a *Auditable) GetUpdatedBy() *uint {
	return a.UpdatedBy
}

func (a *Auditable) SetUpdatedBy(userID *uint) {
	a.UpdatedBy = userID
}

func (a *Auditable) GetDeletedBy() *uint {
	return a.DeletedBy
}

func (a *Auditable) SetDeletedBy(userID *uint) {
	a.DeletedBy = userID
}

func (a *Auditable) GetAuditInfo() *AuditInfo {
	var deletedAt *time.Time
	if a.DeletedAt.Valid {
		deletedAt = &a.DeletedAt.Time
	}
	return &AuditInfo{
		CreatedBy: a.CreatedBy,
		UpdatedBy: a.UpdatedBy,
		DeletedBy: a.DeletedBy,
		CreatedAt: a.CreatedAt.StdTime(),
		UpdatedAt: a.UpdatedAt.StdTime(),
		DeletedAt: deletedAt,
	}
}

// BaseAuditableModel is a convenience type that other models can embed
// This combines standard GORM fields with Auditable fields
type BaseAuditableModel struct {
	// Standard GORM fields
	ID        uint                  `gorm:"primarykey" json:"id"`
	CreatedAt orm.Carbon            `json:"created_at"`
	UpdatedAt orm.Carbon            `json:"updated_at"`
	DeletedAt orm.DeletedAt         `gorm:"index" json:"deleted_at,omitempty"`
	
	// Audit fields
	CreatedBy   *uint     `gorm:"index" json:"created_by,omitempty"`
	Creator     *User     `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	UpdatedBy   *uint     `gorm:"index" json:"updated_by,omitempty"`
	Updater     *User     `gorm:"foreignKey:UpdatedBy" json:"updater,omitempty"`
	DeletedBy   *uint     `gorm:"index" json:"deleted_by,omitempty"`
	Deleter     *User     `gorm:"foreignKey:DeletedBy" json:"deleter,omitempty"`
	IPAddress   string    `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent   string    `gorm:"type:text" json:"user_agent,omitempty"`
}

// Implement AuditableInterface for BaseAuditableModel
func (b *BaseAuditableModel) GetCreatedBy() *uint {
	return b.CreatedBy
}

func (b *BaseAuditableModel) SetCreatedBy(userID *uint) {
	b.CreatedBy = userID
}

func (b *BaseAuditableModel) GetUpdatedBy() *uint {
	return b.UpdatedBy
}

func (b *BaseAuditableModel) SetUpdatedBy(userID *uint) {
	b.UpdatedBy = userID
}

func (b *BaseAuditableModel) GetDeletedBy() *uint {
	return b.DeletedBy
}

func (b *BaseAuditableModel) SetDeletedBy(userID *uint) {
	b.DeletedBy = userID
}

func (b *BaseAuditableModel) GetAuditInfo() *AuditInfo {
	var deletedAt *time.Time
	if b.DeletedAt.Valid {
		deletedAt = &b.DeletedAt.Time
	}
	return &AuditInfo{
		CreatedBy: b.CreatedBy,
		UpdatedBy: b.UpdatedBy,
		DeletedBy: b.DeletedBy,
		CreatedAt: b.CreatedAt.StdTime(),
		UpdatedAt: b.UpdatedAt.StdTime(),
		DeletedAt: deletedAt,
	}
}