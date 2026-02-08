package tenants

import (
	"fmt"
	"strconv"

	"books-database/app/auth"
	"books-database/app/contracts"
	"books-database/app/http/requests"
	"books-database/app/models"
	"books-database/app/services"

	"github.com/goravel/framework/contracts/http"
)

type TenantController struct {
	*contracts.CrudController[models.Tenant, *requests.TenantCreateRequest, *requests.TenantUpdateRequest]
	tenantService *services.TenantService
}

func NewTenantController() *TenantController {
	tenantService := services.NewTenantService()

	crudController := contracts.NewCrudController[models.Tenant, *requests.TenantCreateRequest, *requests.TenantUpdateRequest]("tenant", tenantService).
		WithAuthChecker(func(ctx http.Context, action string, resource interface{}) error {
			permHelper := auth.GetPermissionHelper()
			user, err := permHelper.RequireAuthentication(ctx)
			if err != nil {
				return err
			}
			if !user.IsSuperAdminUser() {
				return fmt.Errorf("super admin access required")
			}
			return nil
		}).
		Build()

	return &TenantController{
		CrudController: crudController,
		tenantService:  tenantService,
	}
}

func (c *TenantController) Index(ctx http.Context) http.Response {
	return c.CrudController.Index(ctx)
}

func (c *TenantController) Show(ctx http.Context) http.Response {
	return c.CrudController.Show(ctx)
}

func (c *TenantController) Store(ctx http.Context) http.Response {
	return c.CrudController.Store(ctx)
}

func (c *TenantController) Update(ctx http.Context) http.Response {
	return c.CrudController.Update(ctx)
}

func (c *TenantController) Delete(ctx http.Context) http.Response {
	return c.CrudController.Delete(ctx)
}

// AssignUser assigns a user to the tenant
func (c *TenantController) AssignUser(ctx http.Context) http.Response {
	if err := c.CheckAuth(ctx, "update", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	tenantID, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid tenant ID", nil)
	}

	userIDStr := ctx.Request().Input("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid user ID", nil)
	}

	var roleID *uint
	if roleIDStr := ctx.Request().Input("role_id"); roleIDStr != "" {
		if id, err := strconv.ParseUint(roleIDStr, 10, 64); err == nil {
			uid := uint(id)
			roleID = &uid
		}
	}

	if err := c.tenantService.AssignUserToTenant(uint(userID), tenantID, roleID); err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	return c.SuccessResponse(ctx, nil, "User assigned to tenant")
}

// RemoveUser removes a user from the tenant
func (c *TenantController) RemoveUser(ctx http.Context) http.Response {
	if err := c.CheckAuth(ctx, "update", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	tenantID, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid tenant ID", nil)
	}

	userIDStr := ctx.Request().Route("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid user ID", nil)
	}

	if err := c.tenantService.RemoveUserFromTenant(uint(userID), tenantID); err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	return c.SuccessResponse(ctx, nil, "User removed from tenant")
}

// GetUsers returns users assigned to the tenant
func (c *TenantController) GetUsers(ctx http.Context) http.Response {
	if err := c.CheckAuth(ctx, "view", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	tenantID, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid tenant ID", nil)
	}

	users, err := c.tenantService.GetTenantUsers(tenantID)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve tenant users: "+err.Error())
	}

	return c.SuccessResponse(ctx, users, "Tenant users retrieved")
}
