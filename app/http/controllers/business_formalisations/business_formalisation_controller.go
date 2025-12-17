package business_formalisations

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/http/requests"
	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
)

// BusinessFormalisationController handles API endpoints for business_formalisation management
type BusinessFormalisationController struct {
	*contracts.CrudController[models.BusinessFormalisation, *requests.BusinessFormalisationCreateRequest, *requests.BusinessFormalisationUpdateRequest]
	businessFormalisationService *services.BusinessFormalisationService
}

// NewBusinessFormalisationController creates a new business_formalisation controller with compile-time enforcement
// This ensures all required methods are implemented and configured
func NewBusinessFormalisationController() *BusinessFormalisationController {
	businessFormalisationService := services.NewBusinessFormalisationService()

	// Build controller with compile-time enforcement
	// The builder pattern ensures all required steps are completed
	crudController := contracts.NewCrudController[models.BusinessFormalisation, *requests.BusinessFormalisationCreateRequest, *requests.BusinessFormalisationUpdateRequest](
		"business_formalisation",
		businessFormalisationService,
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
				_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceBusinessFormalisation, permAction, resource)
				return err
			}

			// For general actions, check without specific resource
			_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceBusinessFormalisation, permAction, nil)
			return err
		}).
		Build()

	controller := &BusinessFormalisationController{
		CrudController:               crudController,
		businessFormalisationService: businessFormalisationService,
	}

	// Set custom hooks
	controller.SetBeforeStore(func(ctx http.Context, data map[string]interface{}) error {
		// Set created_by from authenticated user
		var user models.User
		if err := facades.Auth(ctx).User(&user); err == nil && user.ID > 0 {
			data["created_by"] = user.ID
		}
		return nil
	})

	controller.SetBeforeUpdate(func(ctx http.Context, id uint, data map[string]interface{}) error {
		// Any custom logic before updating a business_formalisation
		return nil
	})

	return controller
}

// StoreOrUpdate creates a new business_formalisation or updates an existing one for the given SME
// This is an upsert operation that prevents duplicate records
// POST /api/business-formalisations/upsert
func (c *BusinessFormalisationController) StoreOrUpdate(ctx http.Context) http.Response {
	// Check permissions
	if err := c.CheckAuth(ctx, "create", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Parse request body
	var data map[string]interface{}
	if err := ctx.Request().Bind(&data); err != nil {
		return c.BadRequestResponse(ctx, "Invalid request body", nil)
	}

	// Set created_by from authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err == nil && user.ID > 0 {
		data["created_by"] = user.ID
	}

	// Use the service's CreateOrUpdate method
	result, err := c.businessFormalisationService.CreateOrUpdate(data)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to create or update business formalisation: "+err.Error())
	}

	return c.SuccessResponse(ctx, result, "Business formalisation saved successfully")
}
