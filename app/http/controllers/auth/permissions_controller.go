package auth

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"players/app/auth"
	"players/app/models"
)

// PermissionsController handles API endpoints for permission assignment
type PermissionsController struct {
}

// Assign POST /api/permissions/assign - Assign a permission to a role
func (c *PermissionsController) Assign(ctx http.Context) http.Response {
	// Check permissions
	permHelper := auth.GetPermissionHelper()
	_, err := permHelper.RequireServicePermission(ctx, auth.ServicePermissions, auth.PermissionUpdate)
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

	// Extract role_id, service, and action
	roleIDFloat, roleOk := requestData["role_id"].(float64)
	service, serviceOk := requestData["service"].(string)
	action, actionOk := requestData["action"].(string)

	if !roleOk || !serviceOk || !actionOk {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "role_id, service, and action are required",
		})
	}

	roleID := uint(roleIDFloat)

	// Find the role
	var role models.Role
	err = facades.Orm().Query().
		Where("id = ? AND is_active = ?", roleID, true).
		First(&role)

	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Role not found",
		})
	}

	// Build permission slug
	permissionSlug := auth.BuildPermissionSlug(auth.ServiceRegistry(service), auth.CorePermissionAction(action))

	// Find the permission
	var permission models.Permission
	err = facades.Orm().Query().
		Where("slug = ? AND is_active = ?", permissionSlug, true).
		First(&permission)

	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": fmt.Sprintf("Permission '%s' not found", permissionSlug),
		})
	}

	// Check if permission is already assigned
	var count int64
	facades.Orm().Query().Model(&models.RolePermission{}).
		Where("role_id = ? AND permission_id = ?", roleID, permission.ID).
		Count(&count)

	if count > 0 {
		return ctx.Response().Json(http.StatusConflict, map[string]string{
			"error": "Permission already assigned to role",
		})
	}

	// Create role-permission assignment
	rolePermission := models.RolePermission{
		RoleID:       roleID,
		PermissionID: permission.ID,
		IsActive:     true,
	}

	err = facades.Orm().Query().Create(&rolePermission)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to assign permission",
		})
	}

	return ctx.Response().Json(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("Permission '%s' assigned to role '%s' successfully", permissionSlug, role.Name),
	})
}

// Revoke DELETE /api/permissions/revoke - Revoke a permission from a role
func (c *PermissionsController) Revoke(ctx http.Context) http.Response {
	// Check permissions
	permHelper := auth.GetPermissionHelper()
	_, err := permHelper.RequireServicePermission(ctx, auth.ServicePermissions, auth.PermissionUpdate)
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

	// Extract role_id, service, and action
	roleIDFloat, roleOk := requestData["role_id"].(float64)
	service, serviceOk := requestData["service"].(string)
	action, actionOk := requestData["action"].(string)

	if !roleOk || !serviceOk || !actionOk {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "role_id, service, and action are required",
		})
	}

	roleID := uint(roleIDFloat)

	// Find the role
	var role models.Role
	err = facades.Orm().Query().
		Where("id = ? AND is_active = ?", roleID, true).
		First(&role)

	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Role not found",
		})
	}

	// Build permission slug
	permissionSlug := auth.BuildPermissionSlug(auth.ServiceRegistry(service), auth.CorePermissionAction(action))

	// Find the permission
	var permission models.Permission
	err = facades.Orm().Query().
		Where("slug = ? AND is_active = ?", permissionSlug, true).
		First(&permission)

	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": fmt.Sprintf("Permission '%s' not found", permissionSlug),
		})
	}

	// Remove role-permission assignment
	_, err = facades.Orm().Query().
		Where("role_id = ? AND permission_id = ?", roleID, permission.ID).
		Delete(&models.RolePermission{})

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to revoke permission",
		})
	}

	return ctx.Response().Json(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("Permission '%s' revoked from role '%s' successfully", permissionSlug, role.Name),
	})
}
