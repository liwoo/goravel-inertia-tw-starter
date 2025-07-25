package perimissions

import (
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"players/app/auth"
	"players/app/contracts"
	"players/app/http/inertia"
	"players/app/models"
	"players/app/services"
)

// PermissionsPageController handles Inertia.js page rendering for permission matrix
type PermissionsPageController struct {
	*contracts.GenericPageController
	permissionsService *services.PermissionsService
}

// GetServiceIdentifier returns the service identifier for this controller
func (c *PermissionsPageController) GetServiceIdentifier() auth.ServiceRegistry {
	return auth.ServicePermissions
}

// NewPermissionsPageController creates a new permissions page controller
func NewPermissionsPageController() *PermissionsPageController {
	roleService := services.NewRoleService()

	return &PermissionsPageController{
		GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
			ResourceType:      "roles",
			PageComponent:     "Permissions/RolesIndex",
			Service:           roleService,
			ServiceIdentifier: auth.ServicePermissions,
			RequireSuperAdmin: true,
			StatsEnabled:      true,
			StatsBuilder: func(controller *contracts.GenericPageController) map[string]interface{} {
				stats, _ := buildRoleStatistics()
				return stats
			},
		}),
		permissionsService: services.NewPermissionsService(),
	}
}

// buildRoleStatistics builds statistics for the roles page
func buildRoleStatistics() (map[string]interface{}, error) {
	var totalRoles, activeRoles, inactiveRoles int64
	var totalUsersWithRoles int64

	facades.Orm().Query().Model(&models.Role{}).Count(&totalRoles)
	facades.Orm().Query().Model(&models.Role{}).Where("is_active = ?", true).Count(&activeRoles)
	facades.Orm().Query().Model(&models.Role{}).Where("is_active = ?", false).Count(&inactiveRoles)
	facades.Orm().Query().Model(&models.UserRole{}).Where("is_active = ?", true).Count(&totalUsersWithRoles)

	return map[string]interface{}{
		"total_roles":            int(totalRoles),
		"active_roles":           int(activeRoles),
		"inactive_roles":         int(inactiveRoles),
		"total_users_with_roles": int(totalUsersWithRoles),
	}, nil
}

// Index GET /admin/permissions - Roles list page
// The generic page controller handles everything for us
func (c *PermissionsPageController) Index(ctx http.Context) http.Response {
	facades.Log().Debug("PermissionsPageController.Index called", map[string]interface{}{
		"url": ctx.Request().Url(),
		"method": ctx.Request().Method(),
		"query": ctx.Request().Queries(),
	})
	return c.GenericPageController.Index(ctx)
}

// RolePermissions GET /admin/roles/:id/permissions - Role permissions page
func (c *PermissionsPageController) RolePermissions(ctx http.Context) http.Response {
	// Debug logging
	roleID := ctx.Request().Route("id")
	facades.Log().Info("=== ROLE PERMISSIONS PAGE LOADED ===", map[string]interface{}{
		"role_id": roleID,
		"url":     ctx.Request().Url(),
		"method":  ctx.Request().Method(),
		"time":    time.Now().Format("15:04:05"),
	})

	// Super-admin only check
	permHelper := auth.GetPermissionHelper()
	user, err := permHelper.RequireAuthentication(ctx)
	if err != nil {
		return ctx.Response().Redirect(302, "/login")
	}

	if !user.IsSuperAdminUser() {
		return ctx.Response().Redirect(302, "/")
	}

	// Get role ID from URL (already declared above)
	if roleID == "" {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Role ID is required",
		})
	}

	// Find the role
	var role models.Role
	err = facades.Orm().Query().
		Where("id = ?", roleID).
		With("Permissions").
		First(&role)

	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Role not found",
		})
	}

	// Get all permissions
	var allPermissions []models.Permission
	facades.Orm().Query().Where("is_active = ?", true).Find(&allPermissions)

	// Get all services and actions
	services := auth.GetAllServiceRegistries()
	actions := auth.GetAllCorePermissionActions()

	// Build service data
	servicesData := make([]map[string]interface{}, 0)
	for _, service := range services {
		serviceActions := auth.GetServiceActions(service)
		actionsMap := make(map[string]bool)
		for _, action := range serviceActions {
			actionsMap[string(action)] = true
		}

		servicesData = append(servicesData, map[string]interface{}{
			"id":      string(service),
			"name":    auth.GetServiceDisplayName(service),
			"slug":    string(service),
			"actions": actionsMap,
		})
	}

	// Build actions data
	actionsData := make([]map[string]interface{}, 0)
	for _, action := range actions {
		actionsData = append(actionsData, map[string]interface{}{
			"id":   string(action),
			"name": auth.GetActionDisplayName(action),
			"slug": string(action),
		})
	}

	// Build the current permissions map with scopes
	// First, we need to load role_permissions with scope information
	var rolePermissions []models.RolePermission
	facades.Orm().Query().
		Where("role_id = ? AND is_active = ?", role.ID, true).
		With("Permission").
		Find(&rolePermissions)

	facades.Log().Info("LOADING PERMISSIONS FOR ROLE", map[string]interface{}{
		"role_id":          role.ID,
		"role_name":        role.Name,
		"permission_count": len(rolePermissions),
		"time":             time.Now().Format("15:04:05"),
	})

	currentPermissions := make(map[string]bool)
	for _, rp := range rolePermissions {
		if rp.Permission.ID > 0 && rp.Permission.IsActive {
			// Create the composite key with scope
			// Default to "by_all" if scope is empty
			scope := rp.Scope
			if scope == "" {
				scope = "by_all"
			}
			permissionKey := rp.Permission.Slug + "_" + scope
			currentPermissions[permissionKey] = true
			facades.Log().Debug("Loaded permission for UI", map[string]interface{}{
				"permission": rp.Permission.Slug,
				"scope":      scope,
				"key":        permissionKey,
			})
		}
	}

	// Get user's permissions
	permissions := c.BuildPermissionsMap(ctx, "roles")

	// Debug the current permissions being sent to frontend
	var permissionKeys []string
	for key := range currentPermissions {
		permissionKeys = append(permissionKeys, key)
	}

	facades.Log().Debug("=== Sending permissions to frontend ===", map[string]interface{}{
		"role_id":        role.ID,
		"role_name":      role.Name,
		"totalCount":     len(currentPermissions),
		"permissionKeys": permissionKeys,
	})

	return inertia.Render(ctx, "Permissions/RolePermissions", map[string]interface{}{
		"role":               role,
		"services":           servicesData,
		"actions":            actionsData,
		"allPermissions":     allPermissions,
		"currentPermissions": currentPermissions,
		"permissions":        permissions,
	})
}

// BuildPermissionsMap builds permissions map for the view
func (c *PermissionsPageController) BuildPermissionsMap(ctx http.Context, resourceType string) map[string]bool {
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return map[string]bool{}
	}

	// Super admins have all permissions
	if user.IsSuperAdminUser() {
		return map[string]bool{
			"canView":       true,
			"canCreate":     true,
			"canEdit":       true,
			"canDelete":     true,
			"canManage":     true,
			"canExport":     true,
			"canBulkUpdate": true,
			"canBulkDelete": true,
			"isAdmin":       true,
			"isSuperAdmin":  true,
		}
	}

	// For other users, check specific permissions
	return map[string]bool{
		"canView":       permHelper.CheckServicePermission(ctx, auth.ServiceRoles, auth.PermissionRead),
		"canCreate":     permHelper.CheckServicePermission(ctx, auth.ServiceRoles, auth.PermissionCreate),
		"canEdit":       permHelper.CheckServicePermission(ctx, auth.ServiceRoles, auth.PermissionUpdate),
		"canDelete":     permHelper.CheckServicePermission(ctx, auth.ServiceRoles, auth.PermissionDelete),
		"canManage":     permHelper.CheckServicePermission(ctx, auth.ServiceRoles, auth.PermissionManage),
		"canExport":     permHelper.CheckServicePermission(ctx, auth.ServiceRoles, auth.PermissionExport),
		"canBulkUpdate": permHelper.CheckServicePermission(ctx, auth.ServiceRoles, auth.PermissionBulkUpdate),
		"canBulkDelete": permHelper.CheckServicePermission(ctx, auth.ServiceRoles, auth.PermissionBulkDelete),
		"isAdmin":       user.Role == "ADMIN" || user.IsSuperAdminUser(),
		"isSuperAdmin":  user.IsSuperAdminUser(),
	}
}
