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

	// Build permissions for the current user
	permissions := c.BuildPermissionsMap(ctx, "roles")

	// Get all services and actions for the permission matrix (using hardcoded auth constants)
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

	// Get all permissions for reference
	var allPermissions []models.Permission
	facades.Orm().Query().Where("is_active = ?", true).Find(&allPermissions)

	// Calculate stats
	var totalRoles, activeRoles, inactiveRoles int64
	var totalUsersWithRoles int64
	
	facades.Orm().Query().Model(&models.Role{}).Count(&totalRoles)
	facades.Orm().Query().Model(&models.Role{}).Where("is_active = ?", true).Count(&activeRoles)
	facades.Orm().Query().Model(&models.Role{}).Where("is_active = ?", false).Count(&inactiveRoles)
	facades.Orm().Query().Model(&models.UserRole{}).Where("is_active = ?", true).Count(&totalUsersWithRoles)

	// Format data for CrudPage component
	data := map[string]interface{}{
		"data":        rolesData,
		"total":       len(rolesData),
		"perPage":     len(rolesData), // For now, show all
		"currentPage": 1,
		"lastPage":    1,
		"from":        1,
		"to":          len(rolesData),
	}

	stats := map[string]interface{}{
		"total_roles":            int(totalRoles),
		"active_roles":           int(activeRoles),
		"inactive_roles":         int(inactiveRoles),
		"total_users_with_roles": int(totalUsersWithRoles),
	}

	// Get roles with permissions using the fixed method
	rolesWithPermissions, err := c.getRolesWithPermissions()
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to load roles with permissions: " + err.Error(),
		})
	}

	// Prepare matrix data
	matrixData := map[string]interface{}{
		"roles":    rolesWithPermissions,
		"services": servicesData,
		"actions":  actionsData,
		"stats": map[string]interface{}{
			"total_services":    len(servicesData),
			"total_actions":     len(actionsData),
			"total_roles":       len(rolesWithPermissions),
			"total_permissions": len(servicesData) * len(actionsData),
		},
	}

	// Render Inertia page
	return inertia.Render(ctx, "Permissions/Index", map[string]interface{}{
		"data":           data,
		"filters":        map[string]interface{}{},
		"stats":          stats,
		"permissions":    permissions,
		"allPermissions": allPermissions,
		"services":       servicesData,
		"actions":        actionsData,
		"matrixData":     matrixData,
		"title":          "Role & Permission Management",
		"subtitle":       "Manage roles and their permissions",
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
		First(&role)

	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Role not found",
		})
	}

	// Get users with this role
	var userRoles []models.UserRole
	facades.Orm().Query().
		Model(&models.UserRole{}).
		Where("role_id = ? AND is_active = ?", role.ID, true).
		Find(&userRoles)

	// Collect user IDs
	userIDs := make([]uint, 0)
	for _, ur := range userRoles {
		userIDs = append(userIDs, ur.UserID)
	}
	
	// Load users
	var users []models.User
	if len(userIDs) > 0 {
		facades.Orm().Query().
			Where("id IN ?", userIDs).
			Find(&users)
	}
	
	// Create user map for easy lookup
	userMap := make(map[uint]models.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// Format users data
	usersData := make([]map[string]interface{}, 0)
	for _, userRole := range userRoles {
		if user, exists := userMap[userRole.UserID]; exists {
			usersData = append(usersData, map[string]interface{}{
				"id":          user.ID,
				"name":        user.Name,
				"email":       user.Email,
				"assigned_at": userRole.AssignedAt,
				"is_active":   userRole.IsActive,
			})
		}
	}

	// Get active permissions for this role from the pivot table
	var rolePermissions []models.RolePermission
	facades.Orm().Query().
		Model(&models.RolePermission{}).
		Where("role_id = ? AND is_active = ?", role.ID, true).
		Find(&rolePermissions)
	
	// Collect permission IDs
	permissionIDs := make([]uint, 0)
	for _, rp := range rolePermissions {
		permissionIDs = append(permissionIDs, rp.PermissionID)
	}
	
	// Load permissions
	var permissions []models.Permission
	if len(permissionIDs) > 0 {
		facades.Orm().Query().
			Where("id IN ? AND is_active = ?", permissionIDs, true).
			Find(&permissions)
	}
	
	// Format role permissions
	permissionsData := make([]map[string]interface{}, 0)
	for _, perm := range permissions {
		permissionsData = append(permissionsData, map[string]interface{}{
			"id":          perm.ID,
			"name":        perm.Name,
			"slug":        perm.Slug,
			"description": perm.Description,
			"category":    perm.Category,
			"action":      perm.Action,
		})
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
		First(&role)

	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Role not found",
		})
	}

	// Get all services and actions for the permission matrix (using hardcoded auth constants)
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

	// Get all permissions for reference
	var allPermissions []models.Permission
	facades.Orm().Query().Where("is_active = ?", true).Find(&allPermissions)

	// Get active permissions for this role from the pivot table
	var rolePermissions []models.RolePermission
	err = facades.Orm().Query().
		Model(&models.RolePermission{}).
		Where("role_id = ? AND is_active = ?", role.ID, true).
		Find(&rolePermissions)
	
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to load role permissions: " + err.Error(),
		})
	}

	// Collect permission IDs
	permissionIDs := make([]uint, 0)
	for _, rp := range rolePermissions {
		permissionIDs = append(permissionIDs, rp.PermissionID)
	}
	
	// Load permissions
	var permissions []models.Permission
	if len(permissionIDs) > 0 {
		err = facades.Orm().Query().
			Where("id IN ? AND is_active = ?", permissionIDs, true).
			Find(&permissions)
		
		if err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
				"error": "Failed to load permissions: " + err.Error(),
			})
		}
	}

	// Format role permissions as array of slugs
	permissionSlugs := make([]string, 0)
	for _, perm := range permissions {
		permissionSlugs = append(permissionSlugs, perm.Slug)
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

	// Get all services and actions for the permission matrix (using hardcoded auth constants)
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

	// Get all permissions for reference
	var allPermissions []models.Permission
	facades.Orm().Query().Where("is_active = ?", true).Find(&allPermissions)

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
		Find(&roles)

	if err != nil {
		return nil, err
	}

	rolesList := make([]map[string]interface{}, 0, len(roles))
	for _, role := range roles {
		// Get active permissions for this role from the pivot table
		var rolePermissions []models.RolePermission
		err := facades.Orm().Query().
			Model(&models.RolePermission{}).
			Where("role_id = ? AND is_active = ?", role.ID, true).
			Find(&rolePermissions)
		
		if err != nil {
			continue // Skip this role on error
		}

		// Collect permission IDs
		permissionIDs := make([]uint, 0)
		for _, rp := range rolePermissions {
			permissionIDs = append(permissionIDs, rp.PermissionID)
		}
		
		// Load permissions
		var perms []models.Permission
		if len(permissionIDs) > 0 {
			facades.Orm().Query().
				Where("id IN ? AND is_active = ?", permissionIDs, true).
				Find(&perms)
		}

		// Build permission matrix for this role
		permissions := make(map[string]bool)
		for _, perm := range perms {
			permissions[perm.Slug] = true
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
	fmt.Printf("DEBUG: IsSuperAdmin() result: %t\n", user.IsSuperAdminUser())
	fmt.Printf("DEBUG: HasRole('super-admin') result: %t\n", user.HasRole("super-admin"))
	fmt.Printf("DEBUG: Legacy role check (role == 'ADMIN'): %t\n", user.Role == "ADMIN")

	if !user.IsSuperAdminUser() && user.Role != "ADMIN" {
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
	if ok && userModel.IsSuperAdminUser() {
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

// RolePermissions GET /admin/roles/:id/permissions - Manage role permissions page
func (c *PermissionsPageController) RolePermissions(ctx http.Context) http.Response {
	// Super-admin only check
	if err := c.requireSuperAdmin(ctx); err != nil {
		return ctx.Response().Redirect(302, "/login")
	}

	// Get role ID from route
	roleID := ctx.Request().Route("id")
	fmt.Printf("DEBUG: RolePermissions - requested role ID: %s\n", roleID)

	var role models.Role
	err := facades.Orm().Query().
		Where("id = ?", roleID).
		First(&role)

	if err != nil {
		fmt.Printf("DEBUG: RolePermissions - role not found: %v\n", err)
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Role not found",
		})
	}
	
	fmt.Printf("DEBUG: RolePermissions - found role: ID=%d, Name=%s, Slug=%s\n", role.ID, role.Name, role.Slug)

	// Get all services and actions for the permission matrix (using hardcoded auth constants)
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

	// Get all permissions for reference
	var allPermissions []models.Permission
	facades.Orm().Query().Where("is_active = ?", true).Find(&allPermissions)

	// Get active permissions for this role from the pivot table
	var rolePermissions []models.RolePermission
	err = facades.Orm().Query().
		Model(&models.RolePermission{}).
		Where("role_id = ? AND is_active = ?", role.ID, true).
		Find(&rolePermissions)
	
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to load role permissions: " + err.Error(),
		})
	}
	
	// Now load the permissions manually
	permissions := make([]models.Permission, 0)
	permissionIDs := make([]uint, 0)
	
	// Collect all permission IDs
	for _, rp := range rolePermissions {
		permissionIDs = append(permissionIDs, rp.PermissionID)
	}
	
	// Load all permissions at once if we have any
	if len(permissionIDs) > 0 {
		err = facades.Orm().Query().
			Where("id IN ? AND is_active = ?", permissionIDs, true).
			Find(&permissions)
		
		if err != nil {
			// Skip if error loading permissions
		}
	}
	
	// Format role permissions as array of slugs
	permissionSlugs := make([]string, 0)
	for _, perm := range permissions {
		permissionSlugs = append(permissionSlugs, perm.Slug)
	}
	
	fmt.Printf("DEBUG: RolePermissions - services count: %d\n", len(servicesData))
	fmt.Printf("DEBUG: RolePermissions - actions count: %d\n", len(actionsData))
	fmt.Printf("DEBUG: RolePermissions - permission slugs count: %d\n", len(permissionSlugs))
	fmt.Printf("DEBUG: RolePermissions - permission slugs: %v\n", permissionSlugs)
	
	// Log first service for debugging
	if len(servicesData) > 0 {
		fmt.Printf("DEBUG: RolePermissions - first service: %+v\n", servicesData[0])
	}
	
	// Log first action for debugging
	if len(actionsData) > 0 {
		fmt.Printf("DEBUG: RolePermissions - first action: %+v\n", actionsData[0])
	}

	// Render Inertia page for permission management
	return inertia.Render(ctx, "Permissions/RolePermissions", map[string]interface{}{
		"role": map[string]interface{}{
			"id":          role.ID,
			"name":        role.Name,
			"slug":        role.Slug,
			"description": role.Description,
			"level":       role.Level,
			"is_active":   role.IsActive,
			"permissions": permissionSlugs,
		},
		"allPermissions": allPermissions,
		"services":       servicesData,
		"actions":        actionsData,
	})
}
