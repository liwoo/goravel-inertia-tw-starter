package models

import (
	"github.com/goravel/framework/database/orm"
)

// Role represents a user role with hierarchical structure
type Role struct {
	orm.Model
	Name        string `gorm:"uniqueIndex;not null" json:"name"`
	Slug        string `gorm:"uniqueIndex;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	Level       int    `gorm:"default:0" json:"level"` // Higher = more permissions
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	
	// Hierarchical structure
	ParentID *uint  `gorm:"index" json:"parent_id,omitempty"`
	Parent   *Role  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Role `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	
	// Relationships
	Users       []User       `gorm:"many2many:user_roles" json:"users,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions,omitempty"`
	
	orm.SoftDeletes
}

// TableName returns the table name for Role model
func (Role) TableName() string {
	return "roles"
}

// HasPermission checks if role has a specific permission
func (r *Role) HasPermission(permission string) bool {
	for _, perm := range r.Permissions {
		if perm.Slug == permission {
			return true
		}
	}
	return false
}

// GetAllPermissions returns all permissions including inherited ones
func (r *Role) GetAllPermissions() []string {
	permissions := make(map[string]bool)
	
	// Add direct permissions
	for _, perm := range r.Permissions {
		if perm.IsActive {
			permissions[perm.Slug] = true
		}
	}
	
	// Add inherited permissions from parent roles
	if r.Parent != nil {
		parentPerms := r.Parent.GetAllPermissions()
		for _, perm := range parentPerms {
			permissions[perm] = true
		}
	}
	
	// Convert map to slice
	result := make([]string, 0, len(permissions))
	for perm := range permissions {
		result = append(result, perm)
	}
	
	return result
}

// IsHigherThan checks if this role has higher level than another role
func (r *Role) IsHigherThan(other *Role) bool {
	return r.Level > other.Level
}

// CanManage checks if this role can manage another role (higher level can manage lower)
func (r *Role) CanManage(other *Role) bool {
	return r.Level > other.Level
}