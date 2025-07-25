package perimissions

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

// NewPermissionsController creates a new permissions controller
func NewPermissionsController() *PermissionsController {
	return &PermissionsController{}
}

// Assign POST /api/permissions/assign - Assign a permission to a role
func (c *PermissionsController) Assign(ctx http.Context) http.Response {
	// Parse request data first to get the role
	var requestData map[string]interface{}
	if err := ctx.Request().Bind(&requestData); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
	}

	// Extract role_id
	roleIDFloat, roleOk := requestData["role_id"].(float64)
	if !roleOk {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "role_id is required",
		})
	}
	roleID := uint(roleIDFloat)

	// Find the role to check scoped permissions
	var role models.Role
	err := facades.Orm().Query().
		Where("id = ? AND is_active = ?", roleID, true).
		First(&role)
	
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Role not found",
		})
	}

	// Check scoped permissions - user needs permission to update this specific role
	scopedHelper := auth.GetScopedPermissionHelper()
	_, err = scopedHelper.RequireScopedPermission(ctx, auth.ServiceRoles, auth.PermissionUpdate, &role)
	if err != nil {
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": "Insufficient permissions to manage this role",
		})
	}

	// Extract service and action
	service, serviceOk := requestData["service"].(string)
	action, actionOk := requestData["action"].(string)

	if !serviceOk || !actionOk {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "service and action are required",
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
	// Parse request data
	var requestData map[string]interface{}
	if err := ctx.Request().Bind(&requestData); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
	}

	// Extract role_id
	roleIDFloat, roleOk := requestData["role_id"].(float64)
	if !roleOk {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "role_id is required",
		})
	}
	roleID := uint(roleIDFloat)

	// Find the role
	var role models.Role
	err := facades.Orm().Query().
		Where("id = ? AND is_active = ?", roleID, true).
		First(&role)

	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Role not found",
		})
	}

	// Check scoped permissions - user needs permission to update this specific role
	scopedHelper := auth.GetScopedPermissionHelper()
	_, err = scopedHelper.RequireScopedPermission(ctx, auth.ServiceRoles, auth.PermissionUpdate, &role)
	if err != nil {
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": "Insufficient permissions to manage this role",
		})
	}

	// Extract service and action
	service, serviceOk := requestData["service"].(string)
	action, actionOk := requestData["action"].(string)

	if !serviceOk || !actionOk {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "service and action are required",
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
