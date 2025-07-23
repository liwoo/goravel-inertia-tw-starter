package users

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"players/app/auth"
	"players/app/contracts"
	"players/app/helpers"
	"players/app/http/requests"
	"players/app/models"
	"players/app/services"
)

// UserController - Simplified version using generic CRUD controller
// From ~300 lines to ~100 lines!
type UserController struct {
	*contracts.GenericCrudController[models.User, requests.UserCreateRequest, requests.UserUpdateRequest]
	userService *services.UserService
}

// NewUserController creates a new simplified user controller
func NewUserController() *UserController {
	userService := services.NewUserService()

	// Create the generic controller
	genericController := contracts.NewGenericCrudController[models.User, requests.UserCreateRequest, requests.UserUpdateRequest](
		"user",
		userService,
	)

	controller := &UserController{
		GenericCrudController: genericController,
		userService:           userService,
	}

	// Configure authorization - User management requires super admin
	genericController.SetAuthCheck(func(ctx http.Context, action string, resource interface{}) error {
		permHelper := auth.GetPermissionHelper()
		user := permHelper.GetAuthenticatedUser(ctx)
		if user == nil || !user.IsSuperAdmin {
			return fmt.Errorf("super admin access required")
		}
		return nil
	})

	// Configure BeforeStore hook to set audit fields for creation
	genericController.SetBeforeStore(func(ctx http.Context, data map[string]interface{}) error {
		auditHelper := helpers.GetAuditHelper()
		auditHelper.SetCreateAuditFields(ctx, data)
		return nil
	})

	// Configure BeforeUpdate hook to set audit fields for updates
	genericController.SetBeforeUpdate(func(ctx http.Context, id uint, data map[string]interface{}) error {
		auditHelper := helpers.GetAuditHelper()
		auditHelper.SetUpdateAuditFields(ctx, data)
		return nil
	})

	// Configure request bindings
	genericController.SetRequestBindings(
		// Bind create request
		func(ctx http.Context) (requests.UserCreateRequest, error) {
			var req requests.UserCreateRequest
			err := ctx.Request().Bind(&req)
			return req, err
		},
		// Transform create request
		func(req requests.UserCreateRequest) map[string]interface{} {
			return req.ToCreateData()
		},
		// Bind update request
		func(ctx http.Context, id uint) (requests.UserUpdateRequest, error) {
			var req requests.UserUpdateRequest
			req.ID = id
			err := ctx.Request().Bind(&req)
			return req, err
		},
		// Transform update request
		func(req requests.UserUpdateRequest) map[string]interface{} {
			return req.ToUpdateData()
		},
	)

	// Override error handling for specific cases
	genericController.SetAfterStore(func(ctx http.Context, result interface{}) http.Response {
		return controller.ResourceCreatedResponse(ctx, result, "user")
	})

	genericController.SetAfterUpdate(func(ctx http.Context, result interface{}) http.Response {
		return controller.ResourceUpdatedResponse(ctx, result, "user")
	})

	// Register controller
	contracts.MustRegisterCrudController("users", controller)

	return controller
}

// GetRoles GET /users/roles - Get all available roles for assignment
func (c *UserController) GetRoles(ctx http.Context) http.Response {
	// Check super admin access using the configured auth check
	if c.GenericCrudController.CheckAuth != nil {
		if err := c.GenericCrudController.CheckAuth(ctx, "viewAny", nil); err != nil {
			return c.ForbiddenResponse(ctx, "Access denied: Super admin privileges required")
		}
	}

	roles, err := c.userService.GetAllRoles()
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve roles: "+err.Error())
	}

	return c.SuccessResponse(ctx, roles, "Roles retrieved successfully")
}

// Contract method implementations
func (c *UserController) GetSearchableFields() []string {
	return c.userService.GetSearchableFields()
}

func (c *UserController) GetValidationRules() map[string]interface{} {
	return c.userService.GetValidationRules()
}

func (c *UserController) CheckPermission(ctx http.Context, permission string, resource interface{}) error {
	// For user management, we only check super admin status
	if c.GenericCrudController.CheckAuth != nil {
		return c.GenericCrudController.CheckAuth(ctx, permission, resource)
	}
	return nil
}

func (c *UserController) GetCurrentUser(ctx http.Context) interface{} {
	permHelper := auth.GetPermissionHelper()
	return permHelper.GetAuthenticatedUser(ctx)
}

func (c *UserController) RequireAuthentication(ctx http.Context) error {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		return fmt.Errorf("authentication required")
	}
	return nil
}

func (c *UserController) BuildPermissionsMap(ctx http.Context, resourceType string) map[string]bool {
	// For super admin only access, all permissions are based on super admin status
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	isSuperAdmin := user != nil && user.IsSuperAdmin

	return map[string]bool{
		"canCreate":      isSuperAdmin,
		"canEdit":        isSuperAdmin,
		"canDelete":      isSuperAdmin,
		"canManage":      isSuperAdmin,
		"canExport":      isSuperAdmin,
		"canViewReports": isSuperAdmin,
	}
}
