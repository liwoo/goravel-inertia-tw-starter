package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	orm.Model
	UserID uint `gorm:"not null;index" json:"user_id"`
	RoleID uint `gorm:"not null;index" json:"role_id"`
	
	// Additional metadata for the relationship
	AssignedByID *uint      `gorm:"index" json:"assigned_by_id,omitempty"`
	AssignedBy   *User      `gorm:"foreignKey:AssignedByID" json:"assigned_by,omitempty"`
	AssignedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"assigned_at"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	Note         string     `gorm:"type:text;column:notes" json:"note,omitempty"`
	
	// Foreign key relationships
	User User `gorm:"foreignKey:UserID" json:"user"`
	Role Role `gorm:"foreignKey:RoleID" json:"role"`
	
	orm.SoftDeletes
}

// TableName returns the table name for UserRole model
func (UserRole) TableName() string {
	return "user_roles"
}

// IsExpired checks if the role assignment has expired
func (ur *UserRole) IsExpired() bool {
	if ur.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*ur.ExpiresAt)
}

// IsValid checks if the role assignment is currently valid
func (ur *UserRole) IsValid() bool {
	return ur.IsActive && !ur.IsExpired()
}

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	orm.Model
	RoleID       uint `gorm:"not null;index" json:"role_id"`
	PermissionID uint `gorm:"not null;index" json:"permission_id"`
	
	// Additional metadata
	GrantedByID *uint     `gorm:"index" json:"granted_by_id,omitempty"`
	GrantedBy   *User     `gorm:"foreignKey:GrantedByID" json:"granted_by,omitempty"`
	GrantedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"granted_at"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	Note        string    `gorm:"type:text;column:notes" json:"note,omitempty"`
	
	// Foreign key relationships
	Role       Role       `gorm:"foreignKey:RoleID" json:"role"`
	Permission Permission `gorm:"foreignKey:PermissionID" json:"permission"`
	
	orm.SoftDeletes
}

// TableName returns the table name for RolePermission model
func (RolePermission) TableName() string {
	return "role_permissions"
}