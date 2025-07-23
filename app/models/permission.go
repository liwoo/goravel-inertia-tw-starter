package models

// Permission represents a specific permission that can be granted to roles
type Permission struct {
	BaseAuditableModel
	
	Name        string `gorm:"uniqueIndex;not null" json:"name"`
	Slug        string `gorm:"uniqueIndex;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	Category    string `gorm:"index;not null" json:"category"` // e.g., "books", "users", "system"
	Resource    string `gorm:"index" json:"resource"`          // Specific resource type
	Action      string `gorm:"index;not null" json:"action"`   // e.g., "create", "read", "update", "delete"
	Scope       string `gorm:"index;default:'by_all'" json:"scope"` // by_me, by_my_role, by_all
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	
	// Permission metadata
	RequiresOwnership bool `gorm:"default:false" json:"requires_ownership"` // User must own the resource
	CanDelegate       bool `gorm:"default:false" json:"can_delegate"`       // Can grant this permission to others
	
	// Relationships
	Roles []Role `gorm:"many2many:role_permissions" json:"roles,omitempty"`
}

// TableName returns the table name for Permission model
func (Permission) TableName() string {
	return "permissions"
}

// GetFullName returns the full permission name (category.action.resource)
func (p *Permission) GetFullName() string {
	if p.Resource != "" {
		return p.Category + "." + p.Action + "." + p.Resource
	}
	return p.Category + "." + p.Action
}

// Matches checks if this permission matches a given permission string
func (p *Permission) Matches(permissionString string) bool {
	return p.Slug == permissionString || p.GetFullName() == permissionString
}

// PermissionCategory represents permission categories
type PermissionCategory struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Actions     []string `json:"actions"`
}

// GetStandardCategories returns standard permission categories
func GetStandardCategories() []PermissionCategory {
	return []PermissionCategory{
		{
			Name:        "Books Management",
			Slug:        "books",
			Description: "Permissions related to book management",
			Actions:     []string{"create", "read", "update", "delete", "borrow", "return", "manage"},
		},
		{
			Name:        "User Management",
			Slug:        "users",
			Description: "Permissions related to user management",
			Actions:     []string{"create", "read", "update", "delete", "impersonate", "manage"},
		},
		{
			Name:        "Role Management",
			Slug:        "roles",
			Description: "Permissions related to role and permission management",
			Actions:     []string{"create", "read", "update", "delete", "assign", "manage"},
		},
		{
			Name:        "System Administration",
			Slug:        "system",
			Description: "System-level permissions",
			Actions:     []string{"backup", "restore", "configure", "monitor", "logs", "manage"},
		},
		{
			Name:        "Reports & Analytics",
			Slug:        "reports",
			Description: "Access to reports and analytics",
			Actions:     []string{"view", "export", "create", "manage"},
		},
	}
}