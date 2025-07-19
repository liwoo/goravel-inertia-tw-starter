package services

import (
	"fmt"
	"strings"
	"players/app/contracts"
	"players/app/helpers"
	"players/app/models"

	"github.com/goravel/framework/facades"
)

// PermissionsService manages role-permission matrix operations
type PermissionsService struct {
	*contracts.BaseCrudService
	authHelper *helpers.AuthHelper
}

// NewPermissionsService creates a new permissions service
func NewPermissionsService() *PermissionsService {
	service := &PermissionsService{
		BaseCrudService: contracts.NewBaseCrudService("permissions_matrix", "id"),
		authHelper:      helpers.NewAuthHelper().(*helpers.AuthHelper),
	}
	return service
}

// PermissionMatrixData represents the complete permission matrix
type PermissionMatrixData struct {
	Roles       []RoleWithPermissions `json:"roles"`
	Permissions []PermissionGrouped   `json:"permissions"`
	Matrix      map[uint][]uint       `json:"matrix"` // RoleID -> []PermissionID
	Stats       MatrixStats           `json:"stats"`
}

// RoleWithPermissions includes role data with current permission assignments
type RoleWithPermissions struct {
	models.Role
	PermissionIDs []uint `json:"permission_ids"`
	PermissionCount int  `json:"permission_count"`
}

// PermissionGrouped represents permissions grouped by category
type PermissionGrouped struct {
	Category    string               `json:"category"`
	Permissions []models.Permission  `json:"permissions"`
}

// MatrixStats provides overview statistics
type MatrixStats struct {
	TotalRoles       int `json:"total_roles"`
	TotalPermissions int `json:"total_permissions"`
	TotalAssignments int `json:"total_assignments"`
	ActiveRoles      int `json:"active_roles"`
	ActivePermissions int `json:"active_permissions"`
}

// BulkAssignmentRequest represents bulk permission assignment operations
type BulkAssignmentRequest struct {
	RoleID        uint   `json:"role_id"`
	PermissionIDs []uint `json:"permission_ids"`
	Action        string `json:"action"` // "assign" or "revoke"
}

// GetPermissionMatrix retrieves the complete permission matrix
func (s *PermissionsService) GetPermissionMatrix() (*PermissionMatrixData, error) {
	// First, sync permissions from registered gates
	if err := s.SyncPermissionsFromGates(); err != nil {
		return nil, fmt.Errorf("failed to sync permissions from gates: %w", err)
	}

	// Get all roles with their permissions
	var roles []models.Role
	err := facades.Orm().Query().
		With("Permissions").
		Where("is_active = ?", true).
		Order("level DESC, name ASC").
		Find(&roles)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roles: %w", err)
	}

	// Get all permissions grouped by category
	var permissions []models.Permission
	err = facades.Orm().Query().
		Where("is_active = ?", true).
		Order("category ASC, action ASC").
		Find(&permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch permissions: %w", err)
	}

	// Build matrix data
	matrix := make(map[uint][]uint)
	rolesWithPermissions := make([]RoleWithPermissions, len(roles))
	
	for i, role := range roles {
		permissionIDs := make([]uint, len(role.Permissions))
		for j, perm := range role.Permissions {
			permissionIDs[j] = perm.ID
		}
		
		rolesWithPermissions[i] = RoleWithPermissions{
			Role:            role,
			PermissionIDs:   permissionIDs,
			PermissionCount: len(permissionIDs),
		}
		
		matrix[role.ID] = permissionIDs
	}

	// Group permissions by category
	permissionGroups := s.groupPermissionsByCategory(permissions)

	// Calculate stats
	stats := s.calculateMatrixStats(roles, permissions, matrix)

	return &PermissionMatrixData{
		Roles:       rolesWithPermissions,
		Permissions: permissionGroups,
		Matrix:      matrix,
		Stats:       stats,
	}, nil
}

// AssignPermissionToRole assigns a permission to a role
func (s *PermissionsService) AssignPermissionToRole(roleID, permissionID uint) error {
	// Check if assignment already exists
	var count int64
	err := facades.Orm().Query().
		Table("role_permissions").
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Count(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing assignment: %w", err)
	}

	if count > 0 {
		return nil // Already assigned
	}

	// Create assignment using model struct
	rolePermission := struct {
		RoleID       uint `gorm:"column:role_id"`
		PermissionID uint `gorm:"column:permission_id"`
	}{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	
	err = facades.Orm().Query().Table("role_permissions").Create(&rolePermission)
	if err != nil {
		return fmt.Errorf("failed to assign permission: %w", err)
	}

	return nil
}

// RevokePermissionFromRole revokes a permission from a role
func (s *PermissionsService) RevokePermissionFromRole(roleID, permissionID uint) error {
	_, err := facades.Orm().Query().
		Table("role_permissions").
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete()
	if err != nil {
		return fmt.Errorf("failed to revoke permission: %w", err)
	}

	return nil
}

// BulkAssignPermissions handles bulk permission assignments/revocations
func (s *PermissionsService) BulkAssignPermissions(request BulkAssignmentRequest) error {
	if request.Action == "assign" {
		for _, permissionID := range request.PermissionIDs {
			if err := s.AssignPermissionToRole(request.RoleID, permissionID); err != nil {
				return fmt.Errorf("failed to assign permission %d: %w", permissionID, err)
			}
		}
	} else if request.Action == "revoke" {
		for _, permissionID := range request.PermissionIDs {
			if err := s.RevokePermissionFromRole(request.RoleID, permissionID); err != nil {
				return fmt.Errorf("failed to revoke permission %d: %w", permissionID, err)
			}
		}
	} else {
		return fmt.Errorf("invalid action: %s", request.Action)
	}

	return nil
}

// SyncRolePermissions completely replaces a role's permissions
func (s *PermissionsService) SyncRolePermissions(roleID uint, permissionIDs []uint) error {
	// Start transaction
	tx, err := facades.Orm().Query().Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Remove all existing permissions for the role
	_, err = tx.Table("role_permissions").
		Where("role_id = ?", roleID).
		Delete()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to clear existing permissions: %w", err)
	}

	// Add new permissions
	for _, permissionID := range permissionIDs {
		rolePermission := struct {
			RoleID       uint `gorm:"column:role_id"`
			PermissionID uint `gorm:"column:permission_id"`
		}{
			RoleID:       roleID,
			PermissionID: permissionID,
		}
		
		err = tx.Table("role_permissions").Create(&rolePermission)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to assign permission %d: %w", permissionID, err)
		}
	}

	tx.Commit()
	return nil
}

// GetRolePermissions gets all permissions for a specific role
func (s *PermissionsService) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	err := facades.Orm().Query().
		Model(&models.Permission{}).
		Where("id IN (SELECT permission_id FROM role_permissions WHERE role_id = ?) AND is_active = ?", roleID, true).
		Find(&permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch role permissions: %w", err)
	}

	return permissions, nil
}

// groupPermissionsByCategory groups permissions by their category
func (s *PermissionsService) groupPermissionsByCategory(permissions []models.Permission) []PermissionGrouped {
	categoryMap := make(map[string][]models.Permission)
	
	for _, perm := range permissions {
		categoryMap[perm.Category] = append(categoryMap[perm.Category], perm)
	}

	var groups []PermissionGrouped
	for category, perms := range categoryMap {
		groups = append(groups, PermissionGrouped{
			Category:    category,
			Permissions: perms,
		})
	}

	return groups
}

// calculateMatrixStats calculates statistics for the permission matrix
func (s *PermissionsService) calculateMatrixStats(roles []models.Role, permissions []models.Permission, matrix map[uint][]uint) MatrixStats {
	activeRoles := 0
	activePermissions := 0
	totalAssignments := 0

	for _, role := range roles {
		if role.IsActive {
			activeRoles++
		}
	}

	for _, perm := range permissions {
		if perm.IsActive {
			activePermissions++
		}
	}

	for _, permissionIDs := range matrix {
		totalAssignments += len(permissionIDs)
	}

	return MatrixStats{
		TotalRoles:       len(roles),
		TotalPermissions: len(permissions),
		TotalAssignments: totalAssignments,
		ActiveRoles:      activeRoles,
		ActivePermissions: activePermissions,
	}
}

// ValidatePermissionAssignment validates if a permission can be assigned to a role
func (s *PermissionsService) ValidatePermissionAssignment(roleID, permissionID uint) error {
	// Check if role exists and is active
	var role models.Role
	err := facades.Orm().Query().Where("id = ? AND is_active = ?", roleID, true).First(&role)
	if err != nil {
		return fmt.Errorf("role not found or inactive: %w", err)
	}

	// Check if permission exists and is active
	var permission models.Permission
	err = facades.Orm().Query().Where("id = ? AND is_active = ?", permissionID, true).First(&permission)
	if err != nil {
		return fmt.Errorf("permission not found or inactive: %w", err)
	}

	return nil
}

// GetPermissionsByCategory returns permissions grouped by category for easier UI rendering
func (s *PermissionsService) GetPermissionsByCategory() (map[string][]models.Permission, error) {
	var permissions []models.Permission
	err := facades.Orm().Query().
		Where("is_active = ?", true).
		Order("category ASC, action ASC").
		Find(&permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch permissions: %w", err)
	}

	result := make(map[string][]models.Permission)
	for _, perm := range permissions {
		result[perm.Category] = append(result[perm.Category], perm)
	}

	return result, nil
}

// SyncPermissionsFromGates syncs the registered gates to the permissions table
func (s *PermissionsService) SyncPermissionsFromGates() error {
	// Import the auth package to access permission constants
	auth := helpers.GetAuth()
	
	// Generate permissions dynamically from the service registry
	var permissions []models.Permission
	
	// Iterate through all registered services
	for _, service := range auth.GetAllServiceRegistries() {
		// Get valid actions for this service
		actions := auth.GetServiceActions(service)
		
		for _, action := range actions {
			// Build the permission slug using the new format
			slug := auth.BuildPermissionSlug(service, action)
			
			// Create permission name and description
			name := fmt.Sprintf("%s %s", auth.GetActionDisplayName(action), auth.GetServiceDisplayName(service))
			description := fmt.Sprintf("Permission to %s for %s", string(action), auth.GetServiceDisplayName(service))
			
			// Special handling for certain actions
			switch action {
			case auth.PermissionView():
				description = fmt.Sprintf("View and list %s", strings.ToLower(auth.GetServiceDisplayName(service)))
			case auth.PermissionManage():
				description = fmt.Sprintf("Full management access for %s", strings.ToLower(auth.GetServiceDisplayName(service)))
			case auth.PermissionExport():
				description = fmt.Sprintf("Export data from %s", strings.ToLower(auth.GetServiceDisplayName(service)))
			}
			
			permission := models.Permission{
				Name:        name,
				Slug:        slug,
				Category:    string(service),
				Resource:    string(service),
				Action:      string(action),
				Description: description,
				IsActive:    true,
			}
			
			permissions = append(permissions, permission)
		}
	}
	
	// Add any custom permissions that don't fit the standard pattern
	customPermissions := []models.Permission{
		{
			Name:        "Impersonate Users",
			Slug:        "users_impersonate",
			Category:    "users",
			Resource:    "users",
			Action:      "impersonate",
			Description: "Ability to impersonate other users",
			IsActive:    true,
		},
		{
			Name:        "Backup System",
			Slug:        "system_backup",
			Category:    "system",
			Resource:    "system",
			Action:      "backup",
			Description: "Create system backups",
			IsActive:    true,
		},
		{
			Name:        "Configure System",
			Slug:        "system_configure",
			Category:    "system",
			Resource:    "system",
			Action:      "configure",
			Description: "Configure system settings",
			IsActive:    true,
		},
	}
	
	permissions = append(permissions, customPermissions...)
	
	// Sync all permissions to the database
	for _, perm := range permissions {
		var existing models.Permission
		err := facades.Orm().Query().Where("slug = ?", perm.Slug).First(&existing)
		
		if err != nil {
			// Permission doesn't exist, create it
			if err := facades.Orm().Query().Create(&perm); err != nil {
				return fmt.Errorf("failed to create permission %s: %w", perm.Slug, err)
			}
		} else {
			// Permission exists, update it
			updateData := map[string]interface{}{
				"name":        perm.Name,
				"category":    perm.Category,
				"resource":    perm.Resource,
				"action":      perm.Action,
				"description": perm.Description,
				"is_active":   true,
			}
			
			if _, err := facades.Orm().Query().Model(&existing).Where("id = ?", existing.ID).Update(updateData); err != nil {
				return fmt.Errorf("failed to update permission %s: %w", perm.Slug, err)
			}
		}
	}
	
	// Mark any permissions not in our list as inactive
	var allDbPermissions []models.Permission
	if err := facades.Orm().Query().Find(&allDbPermissions); err != nil {
		return fmt.Errorf("failed to fetch all permissions: %w", err)
	}
	
	// Create a map of active permission slugs
	activeSlugMap := make(map[string]bool)
	for _, perm := range permissions {
		activeSlugMap[perm.Slug] = true
	}
	
	// Deactivate permissions not in our active list
	for _, dbPerm := range allDbPermissions {
		if !activeSlugMap[dbPerm.Slug] && dbPerm.IsActive {
			updateData := map[string]interface{}{
				"is_active": false,
			}
			if _, err := facades.Orm().Query().Model(&dbPerm).Where("id = ?", dbPerm.ID).Update(updateData); err != nil {
				return fmt.Errorf("failed to deactivate permission %s: %w", dbPerm.Slug, err)
			}
		}
	}
	
	return nil
}