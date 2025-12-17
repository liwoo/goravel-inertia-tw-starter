package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/models"
)

// RoleService implements role-specific business logic using the builder pattern
type RoleService struct {
	contracts.CrudServiceContract // Embedded - automatically exposes all methods!
}

// NewRoleService creates a new role service using the builder pattern
func NewRoleService() *RoleService {
	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.Role]("roles", "id").
		WithSearchFields("name", "slug", "description").                     // REQUIRED
		WithSortFields("id", "name", "slug", "created_at", "updated_at").    // REQUIRED
		WithFilterFields("slug", "is_active", "level", "name", "parent_id"). // REQUIRED
		WithValidationRules(map[string]interface{}{                          // REQUIRED
			"name":        "required|string|max:255",
			"slug":        "required|string|max:100|unique:roles,slug",
			"description": "string|max:1000",
			"is_active":   "boolean",
		}).
		WithRelations("Permissions", "Creator", "Updater").        // Optional
		WithDefaultSort("name", "ASC").                            // Optional
		WithSoftDeletes().                                         // Optional
		WithScopeFiltering("roles", "created_by").                 // Optional
		WithBeforeCreate(func(data map[string]interface{}) error { // Optional
			// Generate slug from name if not provided
			if _, exists := data["slug"]; !exists {
				if name, ok := data["name"].(string); ok {
					data["slug"] = strings.ToLower(strings.ReplaceAll(name, " ", "_"))
				}
			}

			// Set default is_active status
			if _, exists := data["is_active"]; !exists {
				data["is_active"] = true
			}

			return nil
		}).
		WithBeforeUpdate(func(id uint, data map[string]interface{}) error { // Optional
			// Don't allow changing slug for system roles
			// Get the role directly from database to check if it's a system role
			var role models.Role
			err := facades.Orm().Query().Model(&models.Role{}).Where("id = ?", id).First(&role)
			if err == nil && isSystemRole(role.Slug) {
				delete(data, "slug")
				delete(data, "is_active") // System roles should always be active
			}

			return nil
		}).
		WithCustomFilters(func(query orm.Query, filters map[string]interface{}) orm.Query { // Optional
			facades.Log().Debug("RoleService CustomFilters called", map[string]interface{}{
				"filters": filters,
			})

			// First, apply standard filters for fields in filterFields
			standardFilters := []string{"slug", "is_active", "level", "name", "parent_id"}
			for key, value := range filters {
				// Check if it's a standard filter field
				isStandardFilter := false
				for _, field := range standardFilters {
					if field == key {
						isStandardFilter = true
						break
					}
				}

				if isStandardFilter {
					// Convert string boolean values to actual booleans for boolean fields
					if key == "is_active" {
						if strVal, ok := value.(string); ok {
							if strVal == "true" {
								value = true
							} else if strVal == "false" {
								value = false
							}
						}
					}

					facades.Log().Debug("Applying standard filter", map[string]interface{}{
						"field":     key,
						"value":     value,
						"valueType": fmt.Sprintf("%T", value),
					})
					query = query.Where(key+" = ?", value)
				}
			}

			// Then handle special role filters
			for key, value := range filters {
				switch key {
				case "level_min":
					if level, ok := value.(float64); ok {
						query = query.Where("level >= ?", int(level))
					} else if levelStr, ok := value.(string); ok && levelStr != "" {
						query = query.Where("level >= ?", levelStr)
					}
				case "level_max":
					if level, ok := value.(float64); ok {
						query = query.Where("level <= ?", int(level))
					} else if levelStr, ok := value.(string); ok && levelStr != "" {
						query = query.Where("level <= ?", levelStr)
					}
				case "has_parent":
					if hasParent, ok := value.(string); ok {
						if hasParent == "true" {
							query = query.Where("parent_id IS NOT NULL")
						} else if hasParent == "false" {
							query = query.Where("parent_id IS NULL")
						}
					}
				case "type":
					// Handle type filter as an alias for slug
					if typeStr, ok := value.(string); ok && typeStr != "" {
						facades.Log().Debug("Converting type filter to slug filter", map[string]interface{}{
							"type": typeStr,
						})
						query = query.Where("slug = ?", typeStr)
					}
				}
			}
			return query
		}).
		Build() // Returns a fully configured CrudServiceContract

	roleServiceInstance := &RoleService{
		CrudServiceContract: service, // Set the embedded interface
	}

	// Set the actual service reference for proper method resolution
	contracts.SetActualServiceHelper(service, roleServiceInstance, "RoleService")

	return roleServiceInstance
}

// Role-specific methods beyond basic CRUD

// GetBySlug retrieves a role by slug
func (s *RoleService) GetBySlug(slug string) (*models.Role, error) {
	var role models.Role
	err := facades.Orm().Query().
		Model(&models.Role{}).
		Where("slug = ?", slug).
		With("Permissions").
		With("CreatedBy").
		With("UpdatedBy").
		First(&role)

	if err != nil {
		return nil, err
	}

	return &role, nil
}

// GetActiveRoles retrieves all active roles
func (s *RoleService) GetActiveRoles(req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	// Add is_active filter to the request
	if req.Filters == nil {
		req.Filters = make(map[string]interface{})
	}
	req.Filters["is_active"] = true

	return s.GetList(req)
}

// AssignPermissions assigns permissions to a role
func (s *RoleService) AssignPermissions(roleID uint, permissionIDs []uint, scopes map[uint]string) error {
	// Note: We allow modifying permissions for all roles including system roles
	// Super admins should be able to customize all roles
	// The UI/controller layer handles super admin authorization

	// Validate that all permission IDs exist in the permissions table
	if len(permissionIDs) > 0 {
		var validPermissions []models.Permission
		err := facades.Orm().Query().
			Model(&models.Permission{}).
			WhereIn("id", interfaceSlice(permissionIDs)).
			Where("is_active = ?", true).
			Find(&validPermissions)
		if err != nil {
			return fmt.Errorf("failed to validate permissions: %v", err)
		}

		// Build a set of valid permission IDs
		validIDSet := make(map[uint]bool)
		for _, p := range validPermissions {
			validIDSet[p.ID] = true
		}

		// Filter to only valid permission IDs
		validPermissionIDs := make([]uint, 0, len(permissionIDs))
		for _, id := range permissionIDs {
			if validIDSet[id] {
				validPermissionIDs = append(validPermissionIDs, id)
			} else {
				facades.Log().Warningf("Skipping invalid permission ID %d (not found or inactive)", id)
			}
		}
		permissionIDs = validPermissionIDs
	}

	// Begin transaction
	tx, err := facades.Orm().Query().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Remove existing permissions
	_, err = tx.Table("role_permissions").
		Where("role_id = ?", roleID).
		Delete(&models.RolePermission{})
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove existing permissions: %v", err)
	}

	// Add new permissions
	for _, permID := range permissionIDs {
		scope := "by_all" // Default scope
		if s, exists := scopes[permID]; exists {
			scope = s
		}

		rolePermission := models.RolePermission{
			RoleID:       roleID,
			PermissionID: permID,
			Scope:        scope,
		}

		if err := tx.Create(&rolePermission); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to assign permission: %v", err)
		}
	}

	// Commit transaction
	tx.Commit()

	return nil
}

// interfaceSlice converts a []uint to []interface{} for WhereIn
func interfaceSlice(ids []uint) []interface{} {
	result := make([]interface{}, len(ids))
	for i, id := range ids {
		result[i] = id
	}
	return result
}

// GetPermissionsMatrix returns all permissions organized by service and action
func (s *RoleService) GetPermissionsMatrix() (map[string]interface{}, error) {
	var permissions []models.Permission
	err := facades.Orm().Query().
		Model(&models.Permission{}).
		With("CreatedBy").
		Find(&permissions)

	if err != nil {
		return nil, err
	}

	// Organize permissions by service and action
	matrix := make(map[string]map[string]*models.Permission)
	services := make(map[string]bool)
	actions := make(map[string]bool)

	for i := range permissions {
		perm := &permissions[i]

		// Parse the slug to get service and action
		parts := strings.Split(perm.Slug, "_")
		if len(parts) >= 2 {
			service := parts[0]
			action := strings.Join(parts[1:], "_")

			if matrix[service] == nil {
				matrix[service] = make(map[string]*models.Permission)
			}

			matrix[service][action] = perm
			services[service] = true
			actions[action] = true
		}
	}

	// Convert to arrays for frontend
	serviceList := make([]string, 0, len(services))
	for service := range services {
		serviceList = append(serviceList, service)
	}

	actionList := make([]string, 0, len(actions))
	for action := range actions {
		actionList = append(actionList, action)
	}

	return map[string]interface{}{
		"matrix":      matrix,
		"services":    serviceList,
		"actions":     actionList,
		"permissions": permissions,
	}, nil
}

// GetRolePermissions gets permissions for a specific role
func (s *RoleService) GetRolePermissions(roleID uint) ([]models.RolePermission, error) {
	var rolePermissions []models.RolePermission
	err := facades.Orm().Query().
		Model(&models.RolePermission{}).
		Where("role_id = ?", roleID).
		With("Permission").
		Find(&rolePermissions)

	if err != nil {
		return nil, err
	}

	return rolePermissions, nil
}

// CloneRole creates a copy of an existing role with a new name
func (s *RoleService) CloneRole(sourceRoleID uint, newName string, newSlug string) (*models.Role, error) {
	// Get source role
	sourceInterface, err := s.GetByID(sourceRoleID)
	if err != nil {
		return nil, fmt.Errorf("source role not found: %v", err)
	}

	sourceRole, ok := sourceInterface.(*models.Role)
	if !ok {
		return nil, errors.New("invalid role type")
	}

	// Create new role
	newRoleData := map[string]interface{}{
		"name":        newName,
		"slug":        newSlug,
		"description": fmt.Sprintf("Cloned from %s", sourceRole.Name),
		"is_active":   true,
	}

	newRoleInterface, err := s.Create(newRoleData)
	if err != nil {
		return nil, fmt.Errorf("failed to create new role: %v", err)
	}

	newRole, ok := newRoleInterface.(*models.Role)
	if !ok {
		return nil, errors.New("invalid new role type")
	}

	// Copy permissions
	sourcePermissions, err := s.GetRolePermissions(sourceRoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source permissions: %v", err)
	}

	permissionIDs := make([]uint, len(sourcePermissions))
	scopes := make(map[uint]string)

	for i, perm := range sourcePermissions {
		permissionIDs[i] = perm.PermissionID
		scopes[perm.PermissionID] = perm.Scope
	}

	if err := s.AssignPermissions(newRole.ID, permissionIDs, scopes); err != nil {
		// Rollback by deleting the new role
		s.Delete(newRole.ID)
		return nil, fmt.Errorf("failed to copy permissions: %v", err)
	}

	return newRole, nil
}

// isSystemRole checks if a role is a system role that shouldn't be modified
func isSystemRole(slug string) bool {
	systemRoles := []string{"super_admin", "admin", "member"}
	for _, sysRole := range systemRoles {
		if slug == sysRole {
			return true
		}
	}
	return false
}
