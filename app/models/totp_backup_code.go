package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

// TOTPBackupCode represents a single-use backup code for TOTP 2FA
type TOTPBackupCode struct {
	orm.Model

	// User relationship
	UserID uint `gorm:"not null;index" json:"user_id"`

	// Bcrypt hashed backup code
	CodeHash string `gorm:"not null" json:"-"`

	// When the code was used (NULL if not used yet)
	UsedAt *time.Time `json:"used_at,omitempty"`

	// Relationship
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name for TOTPBackupCode model
func (TOTPBackupCode) TableName() string {
	return "totp_backup_codes"
}

// IsUsed returns true if the backup code has been used
func (b *TOTPBackupCode) IsUsed() bool {
	return b.UsedAt != nil
}

// MarkAsUsed marks the backup code as used
func (b *TOTPBackupCode) MarkAsUsed() {
	now := time.Now()
	b.UsedAt = &now
}
