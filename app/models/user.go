package models

import (
	"time"
)

type User struct {
	BaseAuditableModel

	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`

	// Legacy role field (keep for backward compatibility)
	Role string `gorm:"default:'USER'" json:"legacy_role"`

	// User status and metadata
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	IsSuperAdmin  bool       `gorm:"default:false;index" json:"is_super_admin"`
	EmailVerified bool       `gorm:"default:false" json:"email_verified"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`

	// Many-to-many relationships
	Roles []Role `gorm:"many2many:user_roles" json:"roles,omitempty"`
}

// TableName returns the table name for User model
func (User) TableName() string {
	return "users"
}

// HasRole checks if user has a specific role
func (u *User) HasRole(roleSlug string) bool {
	for _, role := range u.Roles {
		if role.Slug == roleSlug && role.IsActive {
			return true
		}
	}
	return false
}

// HasPermission checks if user has a specific permission
func (u *User) HasPermission(permission string) bool {
	for _, role := range u.Roles {
		if !role.IsActive {
			continue
		}

		allPerms := role.GetAllPermissions()
		for _, perm := range allPerms {
			if perm == permission {
				return true
			}
		}
	}
	return false
}

// GetAllPermissions returns all permissions from all user's roles
func (u *User) GetAllPermissions() []string {
	permissions := make(map[string]bool)

	for _, role := range u.Roles {
		if !role.IsActive {
			continue
		}

		rolePerms := role.GetAllPermissions()
		for _, perm := range rolePerms {
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

// GetActiveRoles returns only active roles
func (u *User) GetActiveRoles() []Role {
	activeRoles := make([]Role, 0)
	for _, role := range u.Roles {
		if role.IsActive {
			activeRoles = append(activeRoles, role)
		}
	}
	return activeRoles
}

// GetHighestRole returns the role with the highest level
func (u *User) GetHighestRole() *Role {
	var highest *Role
	for _, role := range u.Roles {
		if !role.IsActive {
			continue
		}
		if highest == nil || role.Level > highest.Level {
			highest = &role
		}
	}
	return highest
}

// CanManageUser checks if this user can manage another user (based on role levels)
func (u *User) CanManageUser(other *User) bool {
	userHighest := u.GetHighestRole()
	otherHighest := other.GetHighestRole()

	if userHighest == nil {
		return false
	}
	if otherHighest == nil {
		return true // Can manage users with no roles
	}

	return userHighest.IsHigherThan(otherHighest)
}

// IsSuperAdminUser checks if user has super admin privileges
func (u *User) IsSuperAdminUser() bool {
	// First check the direct super admin flag
	if u.IsSuperAdmin {
		return true
	}

	// Check legacy ADMIN role for backward compatibility
	legacyAdmin := u.Role == "ADMIN" || u.Role == "SUPER_ADMIN"

	return u.HasRole("super-admin") || u.HasPermission("system.manage") || legacyAdmin
}

// IsAdmin checks if user has admin privileges
func (u *User) IsAdmin() bool {
	return u.IsSuperAdminUser() || u.HasRole("admin")
}

// SharesRoleWith checks if this user shares any active role with another user
func (u *User) SharesRoleWith(other *User) bool {
	if u.ID == other.ID {
		return true
	}

	userRoles := make(map[uint]bool)
	for _, role := range u.GetActiveRoles() {
		userRoles[role.ID] = true
	}

	for _, role := range other.GetActiveRoles() {
		if userRoles[role.ID] {
			return true
		}
	}

	return false
}

// GetRoleLevel returns the user's highest role level, or 0 if no roles
func (u *User) GetRoleLevel() int {
	role := u.GetHighestRole()
	if role == nil {
		return 0
	}
	return role.Level
}

// CanMessageUser checks if this user can send messages to another user
// Users can message others at their role level or lower
func (u *User) CanMessageUser(other *User) bool {
	// Super admins can message anyone
	if u.IsSuperAdminUser() {
		return true
	}

	// Can't message yourself (handled elsewhere, but safety check)
	if u.ID == other.ID {
		return true
	}

	// Get sender's highest role level
	senderLevel := u.GetRoleLevel()

	// Get recipient's highest role level
	recipientLevel := other.GetRoleLevel()

	// Users with no role (level 0) can only message other users with no role
	if senderLevel == 0 {
		return recipientLevel == 0
	}

	// Sender can message if their level >= recipient's level
	return senderLevel >= recipientLevel
}
