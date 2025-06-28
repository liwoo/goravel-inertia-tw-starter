package auth

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"players/app/auth"
	"players/app/contracts"
	"players/app/helpers"
	"players/app/http/inertia"
	"players/app/models"
	"players/app/services"
)

// PermissionsPageController handles Inertia.js page rendering for permission matrix
type PermissionsPageController struct {
	*contracts.BaseCrudController
	permissionsService *services.PermissionsService
	authHelper         contracts.AuthHelper
}

// GetServiceIdentifier returns the service identifier for this controller
func (c *PermissionsPageController) GetServiceIdentifier() auth.ServiceRegistry {
	return auth.ServicePermissions
}

// NewPermissionsPageController creates a new permissions page controller
func NewPermissionsPageController() *PermissionsPageController {
	return &PermissionsPageController{
		BaseCrudController: contracts.NewBaseCrudController("permissions"),
		permissionsService: services.NewPermissionsService(),
		authHelper:         helpers.NewAuthHelper(),
	}
}

// Index GET /admin/permissions - Roles list page
func (c *PermissionsPageController) Index(ctx http.Context) http.Response {
	// Super-admin only check
	if err := c.requireSuperAdmin(ctx); err != nil {
		return ctx.Response().Redirect(302, "/login")
	}

	// Get all roles (only those with valid names)
	var roles []models.Role
	err := facades.Orm().Query().
		Where("is_active = ? AND name != '' AND name IS NOT NULL", true).
		Find(&roles)

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to load roles: " + err.Error(),
		})
	}

	// Format roles for frontend with user counts
	rolesData := make([]map[string]interface{}, 0)
	for _, role := range roles {
		// Count active users for this role
		var userCount int64
		facades.Orm().Query().Model(&models.UserRole{}).
			Where("role_id = ? AND is_active = ?", role.ID, true).
			Count(&userCount)

		rolesData = append(rolesData, map[string]interface{}{
			"id":          role.ID,
			"name":        role.Name,
			"slug":        role.Slug,
			"description": role.Description,
			"level":       role.Level,
			"is_active":   role.IsActive,
			"users_count": int(userCount),
			"created_at":  role.CreatedAt,
			"updated_at":  role.UpdatedAt,
		})
	}

	// Render Inertia page
	return inertia.Render(ctx, "Permissions/RolesIndex", map[string]interface{}{
		"title":    "Roles",
		"subtitle": "Manage user roles and permissions",
		"roles":    rolesData,
	})
}

// RoleShow GET /admin/permissions/roles/:id - View role details page
func (c *PermissionsPageController) RoleShow(ctx http.Context) http.Response {
	// Super-admin only check
	if err := c.requireSuperAdmin(ctx); err != nil {
		return ctx.Response().Redirect(302, "/login")
	}

	// Get role ID from route
	roleID := ctx.Request().Route("id")

	var role models.Role
	err := facades.Orm().Query().
		Where("id = ?", roleID).
		With("Permissions").
		First(&role)

	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Role not found",
		})
	}

	// Get users with this role
	var userRoles []models.UserRole
	facades.Orm().Query().
		Where("role_id = ? AND is_active = ?", role.ID, true).
		With("User").
		Find(&userRoles)

	// Format users data
	usersData := make([]map[string]interface{}, 0)
	for _, userRole := range userRoles {
		if userRole.User.ID != 0 {
			usersData = append(usersData, map[string]interface{}{
				"id":          userRole.User.ID,
				"name":        userRole.User.Name,
				"email":       userRole.User.Email,
				"assigned_at": userRole.AssignedAt,
				"is_active":   userRole.IsActive,
			})
		}
	}

	// Format role permissions
	permissionsData := make([]map[string]interface{}, 0)
	for _, perm := range role.Permissions {
		if perm.IsActive {
			permissionsData = append(permissionsData, map[string]interface{}{
				"id":          perm.ID,
				"name":        perm.Name,
				"slug":        perm.Slug,
				"description": perm.Description,
				"category":    perm.Category,
				"action":      perm.Action,
			})
		}
	}

	// Render Inertia page
	return inertia.Render(ctx, "Permissions/RoleShow", map[string]interface{}{
		"role": map[string]interface{}{
			"id":          role.ID,
			"name":        role.Name,
			"slug":        role.Slug,
			"description": role.Description,
			"level":       role.Level,
			"is_active":   role.IsActive,
			"users_count": len(usersData),
			"created_at":  role.CreatedAt,
			"updated_at":  role.UpdatedAt,
		},
		"users":       usersData,
		"permissions": permissionsData,
	})
}

// RoleEdit GET /admin/permissions/roles/:id/edit - Edit role page
func (c *PermissionsPageController) RoleEdit(ctx http.Context) http.Response {
	// Super-admin only check
	if err := c.requireSuperAdmin(ctx); err != nil {
		return ctx.Response().Redirect(302, "/login")
	}

	// Get role ID from route
	roleID := ctx.Request().Route("id")

	var role models.Role
	err := facades.Orm().Query().
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
	facades.Orm().Query().Find(&allPermissions)

	// Format role permissions as array of slugs
	permissionSlugs := make([]string, 0)
	for _, perm := range role.Permissions {
		if perm.IsActive {
			permissionSlugs = append(permissionSlugs, perm.Slug)
		}
	}

	// Get all services and actions for the permission matrix
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

	// Render Inertia page
	return inertia.Render(ctx, "Permissions/RoleEdit", map[string]interface{}{
		"role": map[string]interface{}{
			"id":          role.ID,
			"name":        role.Name,
			"slug":        role.Slug,
			"description": role.Description,
			"level":       role.Level,
			"is_active":   role.IsActive,
			"users_count": func() int {
				var count int64
				facades.Orm().Query().Model(&models.UserRole{}).
					Where("role_id = ? AND is_active = ?", role.ID, true).
					Count(&count)
				return int(count)
			}(),
			"permissions": permissionSlugs,
		},
		"allPermissions": allPermissions,
		"services":       servicesData,
		"actions":        actionsData,
	})
}

// RoleCreate GET /admin/permissions/roles/create - Create role page
func (c *PermissionsPageController) RoleCreate(ctx http.Context) http.Response {
	// Super-admin only check
	if err := c.requireSuperAdmin(ctx); err != nil {
		return ctx.Response().Redirect(302, "/login")
	}

	// Get all permissions
	var allPermissions []models.Permission
	facades.Orm().Query().Find(&allPermissions)

	// Get all services and actions for the permission matrix
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

	// Render Inertia page
	return inertia.Render(ctx, "Permissions/RoleCreate", map[string]interface{}{
		"allPermissions": allPermissions,
		"services":       servicesData,
		"actions":        actionsData,
	})
}

// getPermissionMatrixData builds the permission matrix using the new service/action structure
func (c *PermissionsPageController) getPermissionMatrixData() (interface{}, error) {
	// Get all services and actions
	services := auth.GetAllServiceRegistries()

	// Build the matrix structure
	matrix := make(map[string]interface{})

	// Services (rows)
	servicesList := make([]map[string]interface{}, 0, len(services))
	for _, service := range services {
		servicesList = append(servicesList, map[string]interface{}{
			"id":      string(service),
			"name":    auth.GetServiceDisplayName(service),
			"slug":    string(service),
			"actions": auth.GetServiceActions(service),
		})
	}

	// Actions (columns) - get all unique actions
	actionsMap := make(map[auth.CorePermissionAction]bool)
	for _, service := range services {
		for _, action := range auth.GetServiceActions(service) {
			actionsMap[action] = true
		}
	}

	actionsList := make([]map[string]interface{}, 0, len(actionsMap))
	for action := range actionsMap {
		actionsList = append(actionsList, map[string]interface{}{
			"id":   string(action),
			"name": auth.GetActionDisplayName(action),
			"slug": string(action),
		})
	}

	// Get all roles with their permissions
	rolesList, err := c.getRolesWithPermissions()
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}

	matrix["services"] = servicesList
	matrix["actions"] = actionsList
	matrix["roles"] = rolesList
	matrix["stats"] = map[string]interface{}{
		"total_services":    len(servicesList),
		"total_actions":     len(actionsList),
		"total_roles":       len(rolesList),
		"total_permissions": len(servicesList) * len(actionsList),
	}

	return matrix, nil
}

// getRolesWithPermissions gets all roles with their current permission assignments
func (c *PermissionsPageController) getRolesWithPermissions() ([]map[string]interface{}, error) {
	var roles []models.Role
	err := facades.Orm().Query().
		Where("is_active = ?", true).
		With("Permissions").
		Find(&roles)

	if err != nil {
		return nil, err
	}

	rolesList := make([]map[string]interface{}, 0, len(roles))
	for _, role := range roles {
		// Build permission matrix for this role
		permissions := make(map[string]bool)
		for _, perm := range role.Permissions {
			if perm.IsActive {
				permissions[perm.Slug] = true
			}
		}

		rolesList = append(rolesList, map[string]interface{}{
			"id":          role.ID,
			"name":        role.Name,
			"slug":        role.Slug,
			"level":       role.Level,
			"permissions": permissions,
		})
	}

	return rolesList, nil
}

// getServicesOfType groups specific services for UI organization
func (c *PermissionsPageController) getServicesOfType(services []auth.ServiceRegistry) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	for _, service := range services {
		// Get available actions for this service
		actions := auth.GetServiceActions(service)

		// Build permissions map for each action
		permissions := make(map[string]bool)
		for _, action := range actions {
			permissions[string(action)] = true
		}

		result = append(result, map[string]interface{}{
			"resource":    string(service),
			"permissions": permissions,
		})
	}

	return result
}

// requireSuperAdmin ensures the user is a super-admin
func (c *PermissionsPageController) requireSuperAdmin(ctx http.Context) error {
	permHelper := auth.GetPermissionHelper()
	user, err := permHelper.RequireAuthentication(ctx)
	if err != nil {
		fmt.Printf("DEBUG: Authentication failed: %v\n", err)
		return fmt.Errorf("authentication required: %w", err)
	}

	fmt.Printf("DEBUG: User loaded for permissions check: ID=%d, Email=%s, Role=%s\n", user.ID, user.Email, user.Role)
	fmt.Printf("DEBUG: User roles count: %d\n", len(user.Roles))
	for _, role := range user.Roles {
		fmt.Printf("DEBUG: User has role: %s (active: %t)\n", role.Slug, role.IsActive)
	}
	fmt.Printf("DEBUG: IsSuperAdmin() result: %t\n", user.IsSuperAdmin())
	fmt.Printf("DEBUG: HasRole('super-admin') result: %t\n", user.HasRole("super-admin"))
	fmt.Printf("DEBUG: Legacy role check (role == 'ADMIN'): %t\n", user.Role == "ADMIN")

	if !user.IsSuperAdmin() && user.Role != "ADMIN" {
		return fmt.Errorf("super-admin access required")
	}

	return nil
}

// CONTRACT IMPLEMENTATIONS

// ValidationControllerContract - Not needed for page controller
func (c *PermissionsPageController) ValidateCreateRequest(ctx http.Context) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented for page controller")
}

func (c *PermissionsPageController) ValidateUpdateRequest(ctx http.Context, id uint) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented for page controller")
}

func (c *PermissionsPageController) GetValidationRules() map[string]interface{} {
	return map[string]interface{}{}
}

// AuthorizationControllerContract
func (c *PermissionsPageController) CheckPermission(ctx http.Context, permission string, resource interface{}) error {
	return c.requireSuperAdmin(ctx)
}

func (c *PermissionsPageController) GetCurrentUser(ctx http.Context) interface{} {
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	return user
}

func (c *PermissionsPageController) RequireAuthentication(ctx http.Context) error {
	permHelper := auth.GetPermissionHelper()
	_, err := permHelper.RequireAuthentication(ctx)
	return err
}

func (c *PermissionsPageController) BuildPermissionsMap(ctx http.Context, resourceType string) map[string]bool {
	// For super-admin permission matrix, always return full access
	user := c.GetCurrentUser(ctx)
	if user == nil {
		return map[string]bool{
			"canCreate": false,
			"canView":   false,
			"canEdit":   false,
			"canDelete": false,
			"canManage": false,
		}
	}

	// Check if user is super-admin
	userModel, ok := user.(*models.User)
	if ok && userModel.IsSuperAdmin() {
		return map[string]bool{
			"canCreate": true,
			"canView":   true,
			"canEdit":   true,
			"canDelete": true,
			"canManage": true,
		}
	}

	return map[string]bool{
		"canCreate": false,
		"canView":   false,
		"canEdit":   false,
		"canDelete": false,
		"canManage": false,
	}
}
