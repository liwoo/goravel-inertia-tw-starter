package roles

import (
	"fmt"

	"starter-project/app/auth"
	"starter-project/app/contracts"
	"starter-project/app/http/requests"
	"starter-project/app/models"
	"starter-project/app/services"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

// RolesController handles API endpoints for role management
type RolesController struct {
	*contracts.CrudController[models.Role, *requests.RoleCreateRequest, *requests.RoleUpdateRequest]
	roleService *services.RoleService
}

// NewRolesController creates a new roles controller with compile-time enforcement
// This ensures all required methods are implemented and configured
func NewRolesController() *RolesController {
	roleService := services.NewRoleService()

	// Build controller with compile-time enforcement
	// The builder pattern ensures all required steps are completed
	crudController := contracts.NewCrudController[models.Role, *requests.RoleCreateRequest, *requests.RoleUpdateRequest](
		"role",
		roleService,
	).
		WithAuthChecker(func(ctx http.Context, action string, resource interface{}) error {
			scopedHelper := auth.GetScopedPermissionHelper()

			// Map generic actions to permission actions
			var permAction auth.CorePermissionAction
			switch action {
			case "viewAny", "view":
				permAction = auth.PermissionRead
			case "create":
				permAction = auth.PermissionCreate
			case "update":
				permAction = auth.PermissionUpdate
			case "delete":
				permAction = auth.PermissionDelete
			default:
				permAction = auth.PermissionManage
			}

			// For specific resource actions, pass the resource
			if resource != nil {
				_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceRoles, permAction, resource)
				return err
			}

			// For general actions, check without specific resource
			_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceRoles, permAction, nil)
			return err
		}).
		Build()

	controller := &RolesController{
		CrudController: crudController,
		roleService:    roleService,
	}

	// Set custom hooks for handling permissions
	controller.SetAfterStore(func(ctx http.Context, result interface{}) http.Response {
		// Handle permission assignment if provided
		if role, ok := result.(*models.Role); ok {
			var requestData requests.RoleCreateRequest
			ctx.Request().Bind(&requestData)

			if len(requestData.Permissions) > 0 {
				// Assign permissions to the newly created role
				if err := assignPermissionsFromSlugs(controller.roleService, role.ID, requestData.Permissions); err != nil {
					// Log error but don't fail the role creation
					facades.Log().Error("Failed to assign permissions to role", map[string]interface{}{
						"role_id": role.ID,
						"error":   err.Error(),
					})
				}
			}
		}

		return nil
	})

	return controller
}

// UpdatePermissions PUT /api/roles/{id}/permissions - Update role permissions
func (c *RolesController) UpdatePermissions(ctx http.Context) http.Response {
	// Check permissions - require super admin for permission management
	permHelper := auth.GetPermissionHelper()
	user, err := permHelper.RequireAuthentication(ctx)
	if err != nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	if !user.IsSuperAdminUser() {
		return c.ForbiddenResponse(ctx, "Super admin access required")
	}

	// Get role ID from URL
	roleIDStr := ctx.Request().Route("id")
	if roleIDStr == "" {
		return c.BadRequestResponse(ctx, "Role ID is required", nil)
	}

	// Find existing role
	roleInterface, err := c.roleService.GetByID(parseUint(roleIDStr))
	if err != nil {
		return c.NotFoundResponse(ctx, "Role not found")
	}

	role, ok := roleInterface.(*models.Role)
	if !ok {
		return c.InternalErrorResponse(ctx, "Invalid role type")
	}

	// Parse request data
	var requestData map[string]interface{}
	if err := ctx.Request().Bind(&requestData); err != nil {
		return c.BadRequestResponse(ctx, "Invalid request data", nil)
	}

	// Get permissions from request
	permissions, ok := requestData["permissions"].([]interface{})
	if !ok {
		return c.BadRequestResponse(ctx, "Permissions array is required", nil)
	}

	facades.Log().Info("UpdatePermissions - Processing permissions", map[string]interface{}{
		"role_id":          role.ID,
		"role_name":        role.Name,
		"permission_count": len(permissions),
	})

	// Parse permissions and scopes
	permissionIDs := make([]uint, 0)
	scopes := make(map[uint]string)

	for _, p := range permissions {
		if permSlug, ok := p.(string); ok {
			// Parse scoped permission (e.g., "books_read_by_all")
			baseSlug, scope := parsePermissionWithScope(permSlug)

			// Find the permission by slug
			var permission models.Permission
			err := facades.Orm().Query().
				Where("slug = ? AND is_active = ?", baseSlug, true).
				First(&permission)

			if err == nil {
				permissionIDs = append(permissionIDs, permission.ID)
				scopes[permission.ID] = scope
			}
		}
	}

	// Use the role service method to assign permissions
	err = c.roleService.AssignPermissions(role.ID, permissionIDs, scopes)
	if err != nil {
		return c.InternalErrorResponse(ctx, fmt.Sprintf("Failed to update permissions: %v", err))
	}

	return c.SuccessResponse(ctx, map[string]interface{}{
		"message": fmt.Sprintf("Permissions updated successfully for role '%s'", role.Name),
		"count":   len(permissionIDs),
	}, "Permissions updated")
}

// Helper functions

func parseUint(s string) uint {
	var id uint64
	fmt.Sscanf(s, "%d", &id)
	return uint(id)
}

func parsePermissionWithScope(permSlug string) (baseSlug string, scope string) {
	// Default scope
	scope = "by_all"
	baseSlug = permSlug

	// Check if this is a scoped permission
	if idx := findLastIndex(permSlug, "_by_"); idx > 0 {
		baseSlug = permSlug[:idx]
		scope = permSlug[idx+1:] // includes "by_"
	}

	return baseSlug, scope
}

func findLastIndex(s, substr string) int {
	lastIdx := -1
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			lastIdx = i
		}
	}
	return lastIdx
}

func assignPermissionsFromSlugs(roleService *services.RoleService, roleID uint, permissionSlugs []string) error {
	permissionIDs := make([]uint, 0)
	scopes := make(map[uint]string)

	for _, permSlug := range permissionSlugs {
		// Parse scoped permission
		baseSlug, scope := parsePermissionWithScope(permSlug)

		// Find the permission by slug
		var permission models.Permission
		err := facades.Orm().Query().
			Where("slug = ? AND is_active = ?", baseSlug, true).
			First(&permission)

		if err == nil {
			permissionIDs = append(permissionIDs, permission.ID)
			scopes[permission.ID] = scope
		}
	}

	if len(permissionIDs) > 0 {
		return roleService.AssignPermissions(roleID, permissionIDs, scopes)
	}

	return nil
}
