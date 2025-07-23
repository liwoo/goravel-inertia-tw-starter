package roles

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"players/app/auth"
	"players/app/helpers"
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

	// Log the request for debugging
	facades.Log().Debug("RolesController.Store called", map[string]interface{}{
		"method":            ctx.Request().Method(),
		"path":              ctx.Request().Path(),
		"url":               ctx.Request().Url(),
		"referer":           ctx.Request().Header("Referer"),
		"user_agent":        ctx.Request().Header("User-Agent"),
		"x_inertia":         ctx.Request().Header("X-Inertia"),
		"x_inertia_version": ctx.Request().Header("X-Inertia-Version"),
		"content_type":      ctx.Request().Header("Content-Type"),
		"data":              requestData,
		"query_params":      ctx.Request().Queries(),
	})

	// Check if request data is empty or contains only empty values
	if len(requestData) == 0 {
		facades.Log().Error("RolesController.Store: Rejecting empty request body", map[string]interface{}{
			"stack_trace": string(debug.Stack()),
		})
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Request body cannot be empty",
		})
	}

	// Validate required fields
	name, nameOk := requestData["name"].(string)
	if !nameOk || strings.TrimSpace(name) == "" {
		facades.Log().Error("RolesController.Store: Name validation failed", map[string]interface{}{
			"name_ok":    nameOk,
			"name":       name,
			"request_id": ctx.Request().Header("X-Request-ID"),
		})
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Role name is required",
		})
	}

	description, _ := requestData["description"].(string)
	level := 1
	if levelFloat, ok := requestData["level"].(float64); ok {
		level = int(levelFloat)
	}

	// Create slug from name with more robust handling
	// Remove all non-alphanumeric characters except spaces and hyphens
	cleanName := strings.TrimSpace(name)
	slug := strings.ToLower(cleanName)

	// Replace spaces with hyphens and remove consecutive hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove any character that's not alphanumeric or hyphen
	var slugBuilder strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			slugBuilder.WriteRune(r)
		}
	}
	slug = slugBuilder.String()

	// Remove leading/trailing hyphens and collapse multiple hyphens
	slug = strings.Trim(slug, "-")

	// Validate slug is not empty after cleaning
	if slug == "" || cleanName == "" {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Role name must contain at least one alphanumeric character",
		})
	}

	// Check if role with this slug already exists
	var existingRole models.Role
	err = facades.Orm().Query().Where("slug = ?", slug).First(&existingRole)
	if err == nil && existingRole.ID > 0 {
		return ctx.Response().Json(http.StatusConflict, map[string]string{
			"error": "A role with this name already exists",
		})
	}

	// Additional validation to ensure name and slug are never empty
	if name == "" || slug == "" {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Role name and slug cannot be empty",
		})
	}

	// Get authenticated user and audit info
	user := permHelper.GetAuthenticatedUser(ctx)
	auditHelper := helpers.GetAuditHelper()

	var createdBy *uint
	if user != nil {
		createdBy = &user.ID
	}

	// Get request metadata
	ipAddr, userAgent := auditHelper.GetRequestMetadata(ctx)

	// Create new role with full audit information
	role := models.Role{
		BaseAuditableModel: models.BaseAuditableModel{
			CreatedBy: createdBy,
			IPAddress: ipAddr,
			UserAgent: userAgent,
		},
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
	if name, ok := requestData["name"].(string); ok {
		cleanName := strings.TrimSpace(name)
		if cleanName != "" {
			// Create slug with same robust handling as in Store
			slug := strings.ToLower(cleanName)
			slug = strings.ReplaceAll(slug, " ", "-")

			// Remove any character that's not alphanumeric or hyphen
			var slugBuilder strings.Builder
			for _, r := range slug {
				if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
					slugBuilder.WriteRune(r)
				}
			}
			slug = slugBuilder.String()
			slug = strings.Trim(slug, "-")

			// Validate slug is not empty after cleaning
			if slug == "" {
				return ctx.Response().Json(http.StatusBadRequest, map[string]string{
					"error": "Role name must contain at least one alphanumeric character",
				})
			}

			// Check if another role already has this slug
			var existingRole models.Role
			err := facades.Orm().Query().
				Where("slug = ? AND id != ?", slug, roleID).
				First(&existingRole)
			if err == nil && existingRole.ID > 0 {
				return ctx.Response().Json(http.StatusConflict, map[string]string{
					"error": "A role with this name already exists",
				})
			}

			role.Name = cleanName
			role.Slug = slug
		}
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

	// Soft delete the role by using Delete method which handles soft deletes
	_, err = facades.Orm().Query().Model(&models.Role{}).Where("id = ?", roleID).Delete(&models.Role{})
	if err != nil {
		facades.Log().Error("Failed to soft delete role", map[string]interface{}{
			"role_id": roleID,
			"error":   err.Error(),
			"stack":   string(debug.Stack()),
		})
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete role: " + err.Error(),
		})
	}

	facades.Log().Info("Role soft deleted successfully", map[string]interface{}{
		"role_id":   roleID,
		"role_name": role.Name,
	})

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

	// Parse scoped permissions (e.g., "books_read_by_all" -> permission: "books_read", scope: "by_all")
	type scopedPermission struct {
		slug  string
		scope string
	}

	scopedPerms := make([]scopedPermission, 0)
	for _, p := range permissions {
		if permSlug, ok := p.(string); ok && strings.TrimSpace(permSlug) != "" {
			permSlug = strings.TrimSpace(permSlug)

			// Check if this is a scoped permission
			if strings.Contains(permSlug, "_by_") {
				// Find the last occurrence of "_by_"
				lastIndex := strings.LastIndex(permSlug, "_by_")
				if lastIndex > 0 {
					baseSlug := permSlug[:lastIndex]
					scopePart := permSlug[lastIndex+1:] // includes "by_"
					scopedPerms = append(scopedPerms, scopedPermission{
						slug:  baseSlug,
						scope: scopePart,
					})
					fmt.Printf("DEBUG: Parsed scoped permission: %s -> base: %s, scope: %s\n", permSlug, baseSlug, scopePart)
				} else {
					// Fallback to non-scoped
					scopedPerms = append(scopedPerms, scopedPermission{
						slug:  permSlug,
						scope: "by_all",
					})
				}
			} else {
				// Non-scoped permission
				scopedPerms = append(scopedPerms, scopedPermission{
					slug:  permSlug,
					scope: "by_all",
				})
			}
		}
	}

	fmt.Printf("DEBUG: UpdatePermissions - Parsed %d scoped permissions\n", len(scopedPerms))

	// Get current active permissions from role_permissions table with scopes
	type currentPermData struct {
		permissionID uint
		slug         string
		scope        string
		compositeKey string // slug + "_" + scope for comparison
	}

	currentPerms := make(map[string]currentPermData) // map[compositeKey]data
	for _, rp := range rolePermissions {
		if rp.IsActive && rp.Permission.ID > 0 {
			compositeKey := rp.Permission.Slug + "_" + rp.Scope
			currentPerms[compositeKey] = currentPermData{
				permissionID: rp.Permission.ID,
				slug:         rp.Permission.Slug,
				scope:        rp.Scope,
				compositeKey: compositeKey,
			}
			fmt.Printf("DEBUG: Current active permission: %s (scope: %s)\n", rp.Permission.Slug, rp.Scope)
		}
	}

	// Create map of new permissions with scopes
	newPerms := make(map[string]scopedPermission) // map[compositeKey]data
	for _, sp := range scopedPerms {
		compositeKey := sp.slug + "_" + sp.scope
		newPerms[compositeKey] = sp
	}

	// Find permissions to add and remove based on composite keys
	type permToAdd struct {
		slug  string
		scope string
	}
	type permToRemove struct {
		permissionID uint
		slug         string
		scope        string
	}

	var toAdd []permToAdd
	var toRemove []permToRemove

	// Find permissions to add
	for compositeKey, sp := range newPerms {
		if _, exists := currentPerms[compositeKey]; !exists {
			toAdd = append(toAdd, permToAdd{
				slug:  sp.slug,
				scope: sp.scope,
			})
			fmt.Printf("DEBUG: Will add permission: %s (scope: %s)\n", sp.slug, sp.scope)
		}
	}

	// Find permissions to remove
	for compositeKey, cp := range currentPerms {
		if _, exists := newPerms[compositeKey]; !exists {
			toRemove = append(toRemove, permToRemove{
				permissionID: cp.permissionID,
				slug:         cp.slug,
				scope:        cp.scope,
			})
			fmt.Printf("DEBUG: Will remove permission: %s (scope: %s)\n", cp.slug, cp.scope)
		}
	}

	// Remove old permission assignments
	if len(toRemove) > 0 {
		for _, ptr := range toRemove {
			// Update role_permission records to inactive based on permission_id AND scope
			_, updateErr := facades.Orm().Query().
				Model(&models.RolePermission{}).
				Where("role_id = ? AND permission_id = ? AND scope = ?", roleID, ptr.permissionID, ptr.scope).
				Update("is_active", false)

			if updateErr != nil {
				fmt.Printf("DEBUG: Failed to remove permission %s (scope: %s) for role %d: %v\n", ptr.slug, ptr.scope, roleID, updateErr)
			} else {
				fmt.Printf("DEBUG: Removed permission %s (scope: %s) for role %d\n", ptr.slug, ptr.scope, roleID)
			}
		}
	}

	// Add new permission assignments
	if len(toAdd) > 0 {
		fmt.Printf("DEBUG: UpdatePermissions - Permissions to add: %v\n", toAdd)

		// Extract unique slugs to fetch from database
		uniqueSlugs := make(map[string]bool)
		for _, pta := range toAdd {
			uniqueSlugs[pta.slug] = true
		}

		var slugsToFetch []string
		for slug := range uniqueSlugs {
			slugsToFetch = append(slugsToFetch, slug)
		}

		// Get permission records from database
		var permsFromDB []models.Permission
		err := facades.Orm().Query().
			Where("slug IN ? AND is_active = ?", slugsToFetch, true).
			Find(&permsFromDB)

		if err != nil {
			fmt.Printf("DEBUG: UpdatePermissions - Error finding permissions: %v\n", err)
		}

		// Create a map for quick lookup
		permMap := make(map[string]models.Permission)
		for _, perm := range permsFromDB {
			permMap[perm.Slug] = perm
		}

		fmt.Printf("DEBUG: UpdatePermissions - Found %d permissions in database for %d unique slugs\n", len(permsFromDB), len(slugsToFetch))

		// Process each permission to add
		for _, pta := range toAdd {
			perm, exists := permMap[pta.slug]
			if !exists {
				fmt.Printf("DEBUG: Permission not found in database: %s\n", pta.slug)
				continue
			}

			// Check if role_permission record already exists (maybe inactive)
			var existingRP models.RolePermission
			err := facades.Orm().Query().
				Where("role_id = ? AND permission_id = ? AND scope = ?", roleID, perm.ID, pta.scope).
				First(&existingRP)

			if err == nil && existingRP.ID > 0 {
				// Record exists, update it to active
				fmt.Printf("DEBUG: Found existing RolePermission record ID=%d for permission %s (scope: %s, IsActive=%v)\n", existingRP.ID, perm.Slug, pta.scope, existingRP.IsActive)

				// Use direct update to set is_active = true
				updateResult, updateErr := facades.Orm().Query().
					Model(&models.RolePermission{}).
					Where("id = ?", existingRP.ID).
					Update("is_active", true)

				if updateErr != nil {
					fmt.Printf("DEBUG: Failed to update permission %s (scope: %s) to active for role %d: %v\n", perm.Slug, pta.scope, roleID, updateErr)
				} else {
					fmt.Printf("DEBUG: Updated permission %s (scope: %s) to active for role %d (rows affected: %d)\n", perm.Slug, pta.scope, roleID, updateResult.RowsAffected)
				}
			} else {
				// Create new role_permission record with scope
				rolePermission := models.RolePermission{
					RoleID:       uint(roleID),
					PermissionID: perm.ID,
					Scope:        pta.scope,
					IsActive:     true,
				}
				createErr := facades.Orm().Query().Create(&rolePermission)
				if createErr != nil {
					fmt.Printf("DEBUG: Failed to create permission %s (scope: %s) for role %d: %v\n", perm.Slug, pta.scope, roleID, createErr)
				} else {
					fmt.Printf("DEBUG: Created permission %s (scope: %s) for role %d\n", perm.Slug, pta.scope, roleID)
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
