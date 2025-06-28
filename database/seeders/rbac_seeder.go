package seeders

import (
	"fmt"
	
	"github.com/goravel/framework/facades"
	"players/app/auth"
	"players/app/models"
)

// RBACSeeder seeds the database with default roles and permissions
type RBACSeeder struct{}

// Signature implements the Seeder interface
func (s *RBACSeeder) Signature() string {
	return "rbac"
}

// Run seeds default roles and permissions
func (s *RBACSeeder) Run() error {
	facades.Log().Info("Starting RBAC Seeder...")
	
	// Clear existing data first
	facades.Orm().Query().Exec("DELETE FROM role_permissions")
	facades.Orm().Query().Exec("DELETE FROM user_roles")
	facades.Orm().Query().Exec("DELETE FROM permissions")
	facades.Orm().Query().Exec("DELETE FROM roles")
	
	// Create roles directly with raw SQL
	rolesSQL := []string{
		"INSERT INTO roles (name, slug, description, level, is_active, created_at, updated_at) VALUES ('Super Administrator', 'super-admin', 'Full system access with all permissions', 100, 1, datetime('now'), datetime('now'))",
		"INSERT INTO roles (name, slug, description, level, is_active, created_at, updated_at) VALUES ('Administrator', 'admin', 'Administrative access to most features', 80, 1, datetime('now'), datetime('now'))",
		"INSERT INTO roles (name, slug, description, level, is_active, created_at, updated_at) VALUES ('Librarian', 'librarian', 'Full book management access', 60, 1, datetime('now'), datetime('now'))",
		"INSERT INTO roles (name, slug, description, level, is_active, created_at, updated_at) VALUES ('Moderator', 'moderator', 'Limited administrative access', 40, 1, datetime('now'), datetime('now'))",
		"INSERT INTO roles (name, slug, description, level, is_active, created_at, updated_at) VALUES ('Member', 'member', 'Regular user with borrowing privileges', 20, 1, datetime('now'), datetime('now'))",
		"INSERT INTO roles (name, slug, description, level, is_active, created_at, updated_at) VALUES ('Guest', 'guest', 'Basic read-only access', 10, 1, datetime('now'), datetime('now'))",
	}
	
	for _, sql := range rolesSQL {
		if _, err := facades.Orm().Query().Exec(sql); err != nil {
			facades.Log().Error("Failed to create role", map[string]interface{}{
				"error": err.Error(),
				"sql": sql[:minInt(50, len(sql))],
			})
		}
	}
	
	// Create permissions dynamically from registered services
	if err := s.createPermissionsFromServices(); err != nil {
		facades.Log().Error("Failed to create permissions from services", map[string]interface{}{
			"error": err.Error(),
		})
		// Fall back to hardcoded permissions
		s.createHardcodedPermissions()
	}
	
	// Assign all permissions to super-admin role
	_, err := facades.Orm().Query().Exec(`
		INSERT INTO role_permissions (role_id, permission_id, is_active, created_at, updated_at)
		SELECT r.id, p.id, 1, datetime('now'), datetime('now')
		FROM roles r, permissions p
		WHERE r.slug = 'super-admin'
	`)
	if err != nil {
		facades.Log().Error("Failed to assign permissions to super-admin", map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	// Assign admin user (if exists) to super-admin role
	var adminUser models.User
	if err := facades.Orm().Query().Where("role = ?", "ADMIN").First(&adminUser); err == nil {
		_, err = facades.Orm().Query().Exec(`
			INSERT INTO user_roles (user_id, role_id, assigned_at, is_active, notes, created_at, updated_at)
			SELECT ?, r.id, datetime('now'), 1, 'Assigned during RBAC seeding', datetime('now'), datetime('now')
			FROM roles r
			WHERE r.slug = 'super-admin'
		`, adminUser.ID)
		if err != nil {
			facades.Log().Error("Failed to assign user to super-admin role", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			facades.Log().Info("Assigned admin user to super-admin role", map[string]interface{}{
				"user_id": adminUser.ID,
				"email": adminUser.Email,
			})
		}
	}
	
	facades.Log().Info("RBAC seeding completed")
	return nil
}

// createPermissions creates default permissions
func (s *RBACSeeder) createPermissions() error {
	permissions := []models.Permission{
		// Books permissions (using service_action format)
		{Name: "Create Books", Slug: "books_create", Category: "books", Action: "create", Description: "Create new books"},
		{Name: "Read Books", Slug: "books_read", Category: "books", Action: "read", Description: "View books"},
		{Name: "Update Books", Slug: "books_update", Category: "books", Action: "update", Description: "Update existing books"},
		{Name: "Delete Books", Slug: "books_delete", Category: "books", Action: "delete", Description: "Delete books"},
		{Name: "Export Books", Slug: "books_export", Category: "books", Action: "export", Description: "Export books data"},
		{Name: "Bulk Update Books", Slug: "books_bulk_update", Category: "books", Action: "bulk_update", Description: "Bulk update books"},
		{Name: "Bulk Delete Books", Slug: "books_bulk_delete", Category: "books", Action: "bulk_delete", Description: "Bulk delete books"},

		// Users permissions (using service_action format)
		{Name: "Create Users", Slug: "users_create", Category: "users", Action: "create", Description: "Create new users"},
		{Name: "Read Users", Slug: "users_read", Category: "users", Action: "read", Description: "View users"},
		{Name: "Update Users", Slug: "users_update", Category: "users", Action: "update", Description: "Update existing users"},
		{Name: "Delete Users", Slug: "users_delete", Category: "users", Action: "delete", Description: "Delete users"},
		{Name: "Export Users", Slug: "users_export", Category: "users", Action: "export", Description: "Export users data"},
		{Name: "Bulk Update Users", Slug: "users_bulk_update", Category: "users", Action: "bulk_update", Description: "Bulk update users"},
		{Name: "Bulk Delete Users", Slug: "users_bulk_delete", Category: "users", Action: "bulk_delete", Description: "Bulk delete users"},

		// Roles permissions (using service_action format)
		{Name: "Create Roles", Slug: "roles_create", Category: "roles", Action: "create", Description: "Create new roles"},
		{Name: "Read Roles", Slug: "roles_read", Category: "roles", Action: "read", Description: "View roles"},
		{Name: "Update Roles", Slug: "roles_update", Category: "roles", Action: "update", Description: "Update existing roles"},
		{Name: "Delete Roles", Slug: "roles_delete", Category: "roles", Action: "delete", Description: "Delete roles"},
		{Name: "Export Roles", Slug: "roles_export", Category: "roles", Action: "export", Description: "Export roles data"},
		{Name: "Bulk Update Roles", Slug: "roles_bulk_update", Category: "roles", Action: "bulk_update", Description: "Bulk update roles"},
		{Name: "Bulk Delete Roles", Slug: "roles_bulk_delete", Category: "roles", Action: "bulk_delete", Description: "Bulk delete roles"},

		// Permissions permissions (using service_action format)
		{Name: "Create Permissions", Slug: "permissions_create", Category: "permissions", Action: "create", Description: "Create new permissions"},
		{Name: "Read Permissions", Slug: "permissions_read", Category: "permissions", Action: "read", Description: "View permissions"},
		{Name: "Update Permissions", Slug: "permissions_update", Category: "permissions", Action: "update", Description: "Update existing permissions"},
		{Name: "Delete Permissions", Slug: "permissions_delete", Category: "permissions", Action: "delete", Description: "Delete permissions"},
		{Name: "Export Permissions", Slug: "permissions_export", Category: "permissions", Action: "export", Description: "Export permissions data"},
		{Name: "Bulk Update Permissions", Slug: "permissions_bulk_update", Category: "permissions", Action: "bulk_update", Description: "Bulk update permissions"},
		{Name: "Bulk Delete Permissions", Slug: "permissions_bulk_delete", Category: "permissions", Action: "bulk_delete", Description: "Bulk delete permissions"},

		// System permissions (using service_action format)
		{Name: "System Manage", Slug: "system_manage", Category: "system", Action: "manage", Description: "Full system management"},

		// Reports permissions (using service_action format)
		{Name: "Read Reports", Slug: "reports_read", Category: "reports", Action: "read", Description: "View reports and analytics"},
		{Name: "Create Reports", Slug: "reports_create", Category: "reports", Action: "create", Description: "Create custom reports"},
		{Name: "Export Reports", Slug: "reports_export", Category: "reports", Action: "export", Description: "Export reports"},
	}

	for _, permission := range permissions {
		var existing models.Permission
		err := facades.Orm().Query().Where("slug = ?", permission.Slug).First(&existing)
		if err != nil {
			// Permission doesn't exist, create it
			permission.IsActive = true // Make sure it's active
			permission.Resource = permission.Category // Set resource field
			permission.RequiresOwnership = false
			permission.CanDelegate = false
			
			// Debug logging
			facades.Log().Info("Creating permission", map[string]interface{}{
				"name": permission.Name,
				"slug": permission.Slug,
				"category": permission.Category,
				"action": permission.Action,
			})
			
			// Create using raw SQL to avoid GORM issues
			query := `INSERT INTO permissions (name, slug, description, category, resource, action, is_active, requires_ownership, can_delegate, created_at, updated_at) 
			         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`
			
			_, err = facades.Orm().Query().Exec(query, 
				permission.Name, 
				permission.Slug, 
				permission.Description, 
				permission.Category, 
				permission.Resource, 
				permission.Action, 
				permission.IsActive, 
				permission.RequiresOwnership, 
				permission.CanDelegate,
			)
			if err != nil {
				facades.Log().Error("Failed to create permission", map[string]interface{}{
					"error": err.Error(),
					"permission": permission,
				})
				return err
			}
		} else {
			// Permission exists, update it to ensure it's active and has correct data
			existing.Name = permission.Name
			existing.Category = permission.Category
			existing.Action = permission.Action
			existing.Description = permission.Description
			existing.IsActive = true
			err = facades.Orm().Query().Save(&existing)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// createRoles creates default roles with hierarchy
func (s *RBACSeeder) createRoles() error {
	type roleData struct {
		Name        string
		Slug        string
		Description string
		Level       int
	}
	
	roles := []roleData{
		{Name: "Super Administrator", Slug: "super-admin", Description: "Full system access with all permissions", Level: 100},
		{Name: "Administrator", Slug: "admin", Description: "Administrative access to most features", Level: 80},
		{Name: "Librarian", Slug: "librarian", Description: "Full book management access", Level: 60},
		{Name: "Moderator", Slug: "moderator", Description: "Limited administrative access", Level: 40},
		{Name: "Member", Slug: "member", Description: "Regular user with borrowing privileges", Level: 20},
		{Name: "Guest", Slug: "guest", Description: "Basic read-only access", Level: 10},
	}

	for _, role := range roles {
		var existing models.Role
		err := facades.Orm().Query().Where("slug = ?", role.Slug).First(&existing)
		if err != nil {
			// Role doesn't exist, create it using raw SQL to avoid GORM issues
			query := `INSERT INTO roles (name, slug, description, level, is_active, created_at, updated_at) 
			         VALUES (?, ?, ?, ?, 1, datetime('now'), datetime('now'))`
			
			_, err = facades.Orm().Query().Exec(query, role.Name, role.Slug, role.Description, role.Level)
			if err != nil {
				facades.Log().Error("Failed to create role", map[string]interface{}{
					"error": err.Error(),
					"role": role,
				})
				return fmt.Errorf("failed to create role %s: %w", role.Slug, err)
			}
			
			facades.Log().Info("Created role", map[string]interface{}{
				"name": role.Name,
				"slug": role.Slug,
			})
		} else {
			facades.Log().Info("Role already exists", map[string]interface{}{
				"name": existing.Name,
				"slug": existing.Slug,
			})
		}
	}

	// Set up role hierarchy (admin inherits from librarian, etc.)
	s.setupRoleHierarchy()

	return nil
}

// setupRoleHierarchy sets up parent-child relationships between roles
func (s *RBACSeeder) setupRoleHierarchy() error {
	hierarchyMap := map[string]string{
		"admin":     "librarian",
		"librarian": "moderator",
		"moderator": "member",
		"member":    "guest",
	}

	for childSlug, parentSlug := range hierarchyMap {
		var child, parent models.Role
		
		// Get child role
		err := facades.Orm().Query().Where("slug = ?", childSlug).First(&child)
		if err != nil {
			continue
		}
		
		// Get parent role
		err = facades.Orm().Query().Where("slug = ?", parentSlug).First(&parent)
		if err != nil {
			continue
		}
		
		// Update child with parent ID
		child.ParentID = &parent.ID
		facades.Orm().Query().Save(&child)
	}

	return nil
}

// assignPermissionsToRoles assigns permissions to roles based on their level
func (s *RBACSeeder) assignPermissionsToRoles() error {
	// Super Admin gets all permissions
	if err := s.assignAllPermissionsToRole("super-admin"); err != nil {
		return err
	}

	// Admin permissions
	adminPerms := []string{
		"books.viewAny", "books.view", "books.create", "books.update", "books.delete", "books.manage", "books.export",
		"users.viewAny", "users.view", "users.create", "users.update", "users.manage",
		"roles.viewAny", "roles.view", "roles.assign",
		"reports.view", "reports.export", "reports.create",
	}
	if err := s.assignPermissionsToRole("admin", adminPerms); err != nil {
		return err
	}

	// Librarian permissions
	librarianPerms := []string{
		"books.viewAny", "books.view", "books.create", "books.update", "books.delete", "books.manage", "books.export",
		"users.viewAny", "users.view",
		"reports.view", "reports.export",
	}
	if err := s.assignPermissionsToRole("librarian", librarianPerms); err != nil {
		return err
	}

	// Moderator permissions
	moderatorPerms := []string{
		"books.viewAny", "books.view", "books.create", "books.update", "books.borrow", "books.return",
		"users.view",
		"reports.view",
	}
	if err := s.assignPermissionsToRole("moderator", moderatorPerms); err != nil {
		return err
	}

	// Member permissions
	memberPerms := []string{
		"books.viewAny", "books.view", "books.borrow", "books.return",
	}
	if err := s.assignPermissionsToRole("member", memberPerms); err != nil {
		return err
	}

	// Guest permissions
	guestPerms := []string{
		"books.viewAny", "books.view",
	}
	if err := s.assignPermissionsToRole("guest", guestPerms); err != nil {
		return err
	}

	return nil
}

// assignAllPermissionsToRole assigns all permissions to a role
func (s *RBACSeeder) assignAllPermissionsToRole(roleSlug string) error {
	var role models.Role
	err := facades.Orm().Query().Where("slug = ?", roleSlug).First(&role)
	if err != nil {
		return err
	}

	var permissions []models.Permission
	err = facades.Orm().Query().Where("is_active = ?", true).Find(&permissions)
	if err != nil {
		return err
	}

	for _, permission := range permissions {
		s.assignPermissionToRole(role.ID, permission.ID)
	}

	return nil
}

// assignPermissionsToRole assigns specific permissions to a role
func (s *RBACSeeder) assignPermissionsToRole(roleSlug string, permissionSlugs []string) error {
	var role models.Role
	err := facades.Orm().Query().Where("slug = ?", roleSlug).First(&role)
	if err != nil {
		return err
	}

	for _, permSlug := range permissionSlugs {
		var permission models.Permission
		err := facades.Orm().Query().Where("slug = ?", permSlug).First(&permission)
		if err != nil {
			continue // Skip if permission doesn't exist
		}

		s.assignPermissionToRole(role.ID, permission.ID)
	}

	return nil
}

// assignPermissionToRole creates a role-permission relationship
func (s *RBACSeeder) assignPermissionToRole(roleID, permissionID uint) error {
	// Check if relationship already exists
	var existing models.RolePermission
	err := facades.Orm().Query().Where("role_id = ? AND permission_id = ?", roleID, permissionID).First(&existing)
	if err != nil {
		// Relationship doesn't exist, create it
		rolePermission := models.RolePermission{
			RoleID:       roleID,
			PermissionID: permissionID,
			IsActive:     true,
		}
		return facades.Orm().Query().Create(&rolePermission)
	}

	return nil
}

// createPermissionsFromServices dynamically creates permissions from registered services
func (s *RBACSeeder) createPermissionsFromServices() error {
	facades.Log().Info("Creating permissions from registered services...")
	
	// Get all registered services
	services := auth.GetAllServiceRegistries()
	facades.Log().Info("Found registered services", map[string]interface{}{
		"count": len(services),
		"services": services,
	})
	
	for _, service := range services {
		// Get actions for this service
		actions := auth.GetServiceActions(service)
		serviceName := auth.GetServiceDisplayName(service)
		
		facades.Log().Info("Processing service", map[string]interface{}{
			"service": string(service),
			"name": serviceName,
			"actions": actions,
		})
		
		for _, action := range actions {
			actionName := auth.GetActionDisplayName(action)
			slug := fmt.Sprintf("%s_%s", string(service), string(action))
			name := fmt.Sprintf("%s %s", actionName, serviceName)
			description := fmt.Sprintf("%s %s in the system", actionName, string(service))
			
			// Create permission using raw SQL
			sql := `INSERT INTO permissions (name, slug, description, category, resource, action, is_active, requires_ownership, can_delegate, created_at, updated_at) 
			       VALUES (?, ?, ?, ?, ?, ?, 1, 0, 0, datetime('now'), datetime('now'))`
			
			_, err := facades.Orm().Query().Exec(sql, name, slug, description, string(service), string(service), string(action))
			if err != nil {
				facades.Log().Error("Failed to create permission", map[string]interface{}{
					"error": err.Error(),
					"service": string(service),
					"action": string(action),
					"slug": slug,
				})
			} else {
				facades.Log().Info("Created permission", map[string]interface{}{
					"name": name,
					"slug": slug,
				})
			}
		}
	}
	
	return nil
}

// createHardcodedPermissions creates a basic set of hardcoded permissions as fallback
func (s *RBACSeeder) createHardcodedPermissions() {
	facades.Log().Info("Creating hardcoded permissions as fallback...")
	
	hardcodedPermissions := []struct {
		name, slug, description, category, action string
	}{
		{"Create Books", "books_create", "Create new books", "books", "create"},
		{"Read Books", "books_read", "View books", "books", "read"},
		{"Update Books", "books_update", "Update existing books", "books", "update"},
		{"Delete Books", "books_delete", "Delete books", "books", "delete"},
		{"Export Books", "books_export", "Export books data", "books", "export"},
		
		{"Create Users", "users_create", "Create new users", "users", "create"},
		{"Read Users", "users_read", "View users", "users", "read"},
		{"Update Users", "users_update", "Update existing users", "users", "update"},
		{"Delete Users", "users_delete", "Delete users", "users", "delete"},
		
		{"Create Roles", "roles_create", "Create new roles", "roles", "create"},
		{"Read Roles", "roles_read", "View roles", "roles", "read"},
		{"Update Roles", "roles_update", "Update existing roles", "roles", "update"},
		{"Delete Roles", "roles_delete", "Delete roles", "roles", "delete"},
		
		{"System Manage", "system_manage", "Full system management", "system", "manage"},
		{"Read Reports", "reports_read", "View reports and analytics", "reports", "read"},
	}
	
	for _, perm := range hardcodedPermissions {
		sql := `INSERT INTO permissions (name, slug, description, category, resource, action, is_active, requires_ownership, can_delegate, created_at, updated_at) 
		       VALUES (?, ?, ?, ?, ?, ?, 1, 0, 0, datetime('now'), datetime('now'))`
		
		_, err := facades.Orm().Query().Exec(sql, perm.name, perm.slug, perm.description, perm.category, perm.category, perm.action)
		if err != nil {
			facades.Log().Error("Failed to create hardcoded permission", map[string]interface{}{
				"error": err.Error(),
				"slug": perm.slug,
			})
		}
	}
}

// minInt returns the minimum of two integers
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}