package auth

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"players/app/auth"
	"players/app/models"
)

// RolesController handles API endpoints for role management
type RolesController struct {
}

// Index GET /api/roles - List all roles
func (c *RolesController) Index(ctx http.Context) http.Response {
	// Check permissions
	permHelper := auth.GetPermissionHelper()
	_, err := permHelper.RequireServicePermission(ctx, auth.ServiceRoles, auth.PermissionRead)
	if err != nil {
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": "Insufficient permissions",
		})
	}

	var roles []models.Role
	err = facades.Orm().Query().
		Where("is_active = ?", true).
		With("Permissions").
		Find(&roles)

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to load roles",
		})
	}

	return ctx.Response().Json(http.StatusOK, map[string]interface{}{
		"roles": roles,
	})
}

// Store POST /api/roles - Create a new role
func (c *RolesController) Store(ctx http.Context) http.Response {
	// Check permissions
	permHelper := auth.GetPermissionHelper()
	_, err := permHelper.RequireServicePermission(ctx, auth.ServiceRoles, auth.PermissionCreate)
	if err != nil {
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": "Insufficient permissions",
		})
	}

	// Parse request data
	var requestData map[string]interface{}
	if err := ctx.Request().Bind(&requestData); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
	}

	// Validate required fields
	name, nameOk := requestData["name"].(string)
	if !nameOk || strings.TrimSpace(name) == "" {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Role name is required",
		})
	}

	description, _ := requestData["description"].(string)
	level := 1
	if levelFloat, ok := requestData["level"].(float64); ok {
		level = int(levelFloat)
	}

	// Create slug from name
	slug := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "-"))

	// Validate slug is not empty
	if slug == "" {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Role name cannot be empty",
		})
	}

	// Check if role with this slug already exists (only check non-empty slugs)
	var existingRole models.Role
	err = facades.Orm().Query().Where("slug = ?", slug).First(&existingRole)
	if err == nil && existingRole.ID > 0 && existingRole.Slug != "" {
		return ctx.Response().Json(http.StatusConflict, map[string]string{
			"error": "A role with this name already exists",
		})
	}

	// Create new role
	role := models.Role{
		Name:        name,
		Slug:        slug,
		Description: description,
		Level:       level,
		IsActive:    true,
	}

	err = facades.Orm().Query().Create(&role)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create role",
		})
	}

	// Handle permission assignments if provided
	if permissions, ok := requestData["permissions"].([]interface{}); ok && len(permissions) > 0 {
		for _, p := range permissions {
			if permSlug, ok := p.(string); ok {
				// Find the permission by slug
				var permission models.Permission
				err := facades.Orm().Query().
					Where("slug = ? AND is_active = ?", permSlug, true).
					First(&permission)

				if err == nil {
					// Create role-permission assignment
					rolePermission := models.RolePermission{
						RoleID:       role.ID,
						PermissionID: permission.ID,
						IsActive:     true,
					}
					facades.Orm().Query().Create(&rolePermission)
				}
			}
		}
	}

	return ctx.Response().Json(http.StatusCreated, map[string]interface{}{
		"message": "Role created successfully",
		"role":    role,
	})
}

// Show GET /api/roles/{id} - Get a specific role
func (c *RolesController) Show(ctx http.Context) http.Response {
	// Check permissions
	permHelper := auth.GetPermissionHelper()
	_, err := permHelper.RequireServicePermission(ctx, auth.ServiceRoles, auth.PermissionRead)
	if err != nil {
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": "Insufficient permissions",
		})
	}

	// Get role ID from URL
	roleID, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 32)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Invalid role ID",
		})
	}

	var role models.Role
	err = facades.Orm().Query().
		Where("id = ? AND is_active = ?", roleID, true).
		With("Permissions").
		First(&role)

	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Role not found",
		})
	}

	return ctx.Response().Json(http.StatusOK, map[string]interface{}{
		"role": role,
	})
}

// Update PUT /api/roles/{id} - Update a role
func (c *RolesController) Update(ctx http.Context) http.Response {
	// Check permissions
	permHelper := auth.GetPermissionHelper()
	_, err := permHelper.RequireServicePermission(ctx, auth.ServiceRoles, auth.PermissionUpdate)
	if err != nil {
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": "Insufficient permissions",
		})
	}

	// Get role ID from URL
	roleID, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 32)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Invalid role ID",
		})
	}

	// Find existing role
	var role models.Role
	err = facades.Orm().Query().
		Where("id = ? AND is_active = ?", roleID, true).
		First(&role)

	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Role not found",
		})
	}

	// Parse request data
	var requestData map[string]interface{}
	if err := ctx.Request().Bind(&requestData); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
	}

	// Update fields if provided
	if name, ok := requestData["name"].(string); ok && strings.TrimSpace(name) != "" {
		role.Name = strings.TrimSpace(name)
		role.Slug = strings.ToLower(strings.ReplaceAll(role.Name, " ", "-"))
	}

	if description, ok := requestData["description"].(string); ok {
		role.Description = description
	}

	if levelFloat, ok := requestData["level"].(float64); ok {
		role.Level = int(levelFloat)
	}

	// Save changes
	err = facades.Orm().Query().Save(&role)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update role",
		})
	}

	// Handle permission updates if provided
	if permissions, ok := requestData["permissions"].([]interface{}); ok {
		// Get current role permissions
		var currentPermissions []models.Permission
		facades.Orm().Query().
			Model(&role).
			Association("Permissions").
			Find(&currentPermissions)

		// Create maps for efficient lookup
		currentPermMap := make(map[string]bool)
		for _, perm := range currentPermissions {
			currentPermMap[perm.Slug] = true
		}

		newPermMap := make(map[string]bool)
		for _, p := range permissions {
			if slug, ok := p.(string); ok {
				newPermMap[slug] = true
			}
		}

		// Find permissions to add and remove
		var toAdd []string
		var toRemove []string

		// Find permissions to add
		for slug := range newPermMap {
			if !currentPermMap[slug] {
				toAdd = append(toAdd, slug)
			}
		}

		// Find permissions to remove
		for slug := range currentPermMap {
			if !newPermMap[slug] {
				toRemove = append(toRemove, slug)
			}
		}

		// Add new permissions
		if len(toAdd) > 0 {
			var permsToAdd []models.Permission
			facades.Orm().Query().
				Where("slug IN ? AND is_active = ?", toAdd, true).
				Find(&permsToAdd)

			if len(permsToAdd) > 0 {
				facades.Orm().Query().
					Model(&role).
					Association("Permissions").
					Append(&permsToAdd)
			}
		}

		// Remove permissions
		if len(toRemove) > 0 {
			var permsToRemove []models.Permission
			facades.Orm().Query().
				Where("slug IN ?", toRemove).
				Find(&permsToRemove)

			if len(permsToRemove) > 0 {
				facades.Orm().Query().
					Model(&role).
					Association("Permissions").
					Delete(&permsToRemove)
			}
		}
	}

	return ctx.Response().Json(http.StatusOK, map[string]interface{}{
		"message": "Role updated successfully",
		"role":    role,
	})
}

// Destroy DELETE /api/roles/{id} - Delete a role
func (c *RolesController) Destroy(ctx http.Context) http.Response {
	// Check permissions
	permHelper := auth.GetPermissionHelper()
	_, err := permHelper.RequireServicePermission(ctx, auth.ServiceRoles, auth.PermissionDelete)
	if err != nil {
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": "Insufficient permissions",
		})
	}

	// Get role ID from URL
	roleID, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 32)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Invalid role ID",
		})
	}

	// Find existing role
	var role models.Role
	err = facades.Orm().Query().
		Where("id = ? AND is_active = ?", roleID, true).
		First(&role)

	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Role not found",
		})
	}

	// Check if role has users assigned
	var userCount int64
	facades.Orm().Query().Model(&models.UserRole{}).
		Where("role_id = ? AND is_active = ?", roleID, true).
		Count(&userCount)

	if userCount > 0 {
		return ctx.Response().Json(http.StatusConflict, map[string]string{
			"error": fmt.Sprintf("Cannot delete role: %d users are assigned to this role", userCount),
		})
	}

	// Soft delete the role
	role.IsActive = false
	err = facades.Orm().Query().Save(&role)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete role",
		})
	}

	return ctx.Response().Json(http.StatusOK, map[string]interface{}{
		"message": "Role deleted successfully",
	})
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// UpdatePermissions PUT /api/roles/{id}/permissions - Update role permissions
func (c *RolesController) UpdatePermissions(ctx http.Context) http.Response {
	// Check permissions - require super admin for permission management
	permHelper := auth.GetPermissionHelper()
	user, err := permHelper.RequireAuthentication(ctx)
	if err != nil {
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": "Authentication required",
		})
	}

	if !user.IsSuperAdminUser() && user.Role != "ADMIN" {
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": "Super admin access required",
		})
	}

	// Get role ID from URL
	roleID, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 32)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Invalid role ID",
		})
	}

	// Find existing role
	var role models.Role
	err = facades.Orm().Query().
		Where("id = ? AND is_active = ?", roleID, true).
		First(&role)
	
	// Load ALL role_permissions for this role, not just active ones
	var rolePermissions []models.RolePermission
	facades.Orm().Query().
		Where("role_id = ?", roleID).
		With("Permission").
		Find(&rolePermissions)
	
	fmt.Printf("DEBUG: Found %d role_permissions records for role %d\n", len(rolePermissions), roleID)
	for _, rp := range rolePermissions {
		fmt.Printf("DEBUG: RolePermission ID=%d, PermissionID=%d, IsActive=%v\n", rp.ID, rp.PermissionID, rp.IsActive)
	}

	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Role not found",
		})
	}

	// Parse request data
	var requestData map[string]interface{}
	if err := ctx.Request().Bind(&requestData); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
	}

	// Get permissions from request
	permissions, ok := requestData["permissions"].([]interface{})
	if !ok {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Permissions array is required",
		})
	}

	fmt.Printf("DEBUG: UpdatePermissions - Role ID: %d, Role Name: %s\n", roleID, role.Name)
	fmt.Printf("DEBUG: UpdatePermissions - Received permissions: %v\n", permissions)

	// Convert to string array
	permissionSlugs := make([]string, 0)
	for _, p := range permissions {
		if slug, ok := p.(string); ok && strings.TrimSpace(slug) != "" {
			permissionSlugs = append(permissionSlugs, strings.TrimSpace(slug))
		}
	}
	
	fmt.Printf("DEBUG: UpdatePermissions - Permission slugs: %v\n", permissionSlugs)

	// Get current active permissions from role_permissions table
	currentPermissionSlugs := make([]string, 0)
	for _, rp := range rolePermissions {
		if rp.IsActive && rp.Permission.ID > 0 {
			currentPermissionSlugs = append(currentPermissionSlugs, rp.Permission.Slug)
			fmt.Printf("DEBUG: Active permission: %s\n", rp.Permission.Slug)
		}
	}
	fmt.Printf("DEBUG: Current active permission slugs: %v\n", currentPermissionSlugs)

	// Create maps for efficient lookup
	currentPermMap := make(map[string]bool)
	for _, slug := range currentPermissionSlugs {
		currentPermMap[slug] = true
	}

	newPermMap := make(map[string]bool)
	for _, slug := range permissionSlugs {
		newPermMap[slug] = true
	}

	// Find permissions to add and remove
	var toAdd []string
	var toRemove []string

	// Find permissions to add
	for slug := range newPermMap {
		if !currentPermMap[slug] {
			toAdd = append(toAdd, slug)
		}
	}

	// Find permissions to remove
	for slug := range currentPermMap {
		if !newPermMap[slug] {
			toRemove = append(toRemove, slug)
		}
	}

	// Remove old permission assignments
	if len(toRemove) > 0 {
		// Get permission IDs to remove
		var permsToRemove []models.Permission
		facades.Orm().Query().
			Where("slug IN ? AND is_active = ?", toRemove, true).
			Find(&permsToRemove)

		if len(permsToRemove) > 0 {
			for _, perm := range permsToRemove {
				// Update role_permission records to inactive instead of deleting
				_, updateErr := facades.Orm().Query().
					Model(&models.RolePermission{}).
					Where("role_id = ? AND permission_id = ?", roleID, perm.ID).
					Update("is_active", false)
				
				if updateErr != nil {
					fmt.Printf("DEBUG: Failed to remove permission %s for role %d: %v\n", perm.Slug, roleID, updateErr)
				} else {
					fmt.Printf("DEBUG: Removed permission %s for role %d\n", perm.Slug, roleID)
				}
			}
		}
	}

	// Add new permission assignments
	if len(toAdd) > 0 {
		fmt.Printf("DEBUG: UpdatePermissions - Permissions to add: %v\n", toAdd)
		
		// Debug what we currently have
		fmt.Printf("DEBUG: Current permission slugs from role: %v\n", currentPermissionSlugs)
		fmt.Printf("DEBUG: New permission slugs requested: %v\n", permissionSlugs)
		
		// Get permission records to add
		var permsToAdd []models.Permission
		
		// Debug: First check what permissions exist in DB
		var allPerms []models.Permission
		facades.Orm().Query().Find(&allPerms)
		fmt.Printf("DEBUG: Total permissions in DB: %d\n", len(allPerms))
		
		// List all permission slugs in DB
		dbPermSlugs := make([]string, 0)
		for _, p := range allPerms {
			dbPermSlugs = append(dbPermSlugs, p.Slug)
			if contains(toAdd, p.Slug) {
				fmt.Printf("DEBUG: Permission to add found in DB: %s (ID: %d, IsActive: %v)\n", p.Slug, p.ID, p.IsActive)
			}
		}
		fmt.Printf("DEBUG: All permission slugs in DB: %v\n", dbPermSlugs)
		
		err := facades.Orm().Query().
			Where("slug IN ? AND is_active = ?", toAdd, true).
			Find(&permsToAdd)
		
		if err != nil {
			fmt.Printf("DEBUG: UpdatePermissions - Error finding permissions to add: %v\n", err)
		}
		
		fmt.Printf("DEBUG: UpdatePermissions - Found %d permissions in database for %d slugs\n", len(permsToAdd), len(toAdd))
		fmt.Printf("DEBUG: UpdatePermissions - Looking for slugs: %v\n", toAdd)
		for _, perm := range permsToAdd {
			fmt.Printf("DEBUG: UpdatePermissions - Found permission: ID=%d, Slug=%s, IsActive=%v\n", perm.ID, perm.Slug, perm.IsActive)
		}

		if len(permsToAdd) > 0 {
			for _, perm := range permsToAdd {
				// Check if role_permission record already exists (maybe inactive)
				var existingRP models.RolePermission
				err := facades.Orm().Query().
					Where("role_id = ? AND permission_id = ?", roleID, perm.ID).
					First(&existingRP)
				
				if err == nil && existingRP.ID > 0 {
					// Record exists, update it to active
					fmt.Printf("DEBUG: Found existing RolePermission record ID=%d for permission %s (IsActive=%v)\n", existingRP.ID, perm.Slug, existingRP.IsActive)
					
					// Use direct update instead of Save
					updateResult, updateErr := facades.Orm().Query().
						Model(&models.RolePermission{}).
						Where("id = ?", existingRP.ID).
						Update("is_active", true)
					
					if updateErr != nil {
						fmt.Printf("DEBUG: Failed to update permission %s to active for role %d: %v\n", perm.Slug, roleID, updateErr)
					} else {
						// Verify the update
						var verifyRP models.RolePermission
						facades.Orm().Query().Where("id = ?", existingRP.ID).First(&verifyRP)
						fmt.Printf("DEBUG: Updated permission %s to active for role %d (ID: %d, IsActive after save: %v, rows affected: %d)\n", perm.Slug, roleID, existingRP.ID, verifyRP.IsActive, updateResult.RowsAffected)
					}
				} else {
					// Create new role_permission record
					rolePermission := models.RolePermission{
						RoleID:       uint(roleID),
						PermissionID: perm.ID,
						IsActive:     true,
					}
					createErr := facades.Orm().Query().Create(&rolePermission)
					if createErr != nil {
						fmt.Printf("DEBUG: Failed to create permission %s for role %d: %v\n", perm.Slug, roleID, createErr)
					} else {
						fmt.Printf("DEBUG: Created permission %s for role %d\n", perm.Slug, roleID)
					}
				}
			}
		}
	}

	return ctx.Response().Json(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("Permissions updated successfully. Added: %d, Removed: %d", len(toAdd), len(toRemove)),
		"added":   len(toAdd),
		"removed": len(toRemove),
	})
}
