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
