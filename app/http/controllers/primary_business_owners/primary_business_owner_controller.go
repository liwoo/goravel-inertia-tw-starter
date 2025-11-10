package primary_business_owners

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/http/requests"
	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
)

// PrimaryBusinessOwnerController handles API endpoints for primary business owner management
type PrimaryBusinessOwnerController struct {
	*contracts.CrudController[models.PrimaryBusinessOwner, *requests.PrimaryBusinessOwnerCreateRequest, *requests.PrimaryBusinessOwnerUpdateRequest]
	primaryBusinessOwnerService *services.PrimaryBusinessOwnerService
}

// NewPrimaryBusinessOwnerController creates a new primary business owner controller with compile-time enforcement
// This ensures all required methods are implemented and configured
func NewPrimaryBusinessOwnerController() *PrimaryBusinessOwnerController {
	primaryBusinessOwnerService := services.NewPrimaryBusinessOwnerService()

	// Build controller with compile-time enforcement
	// The builder pattern ensures all required steps are completed
	crudController := contracts.NewCrudController[models.PrimaryBusinessOwner, *requests.PrimaryBusinessOwnerCreateRequest, *requests.PrimaryBusinessOwnerUpdateRequest](
		"primary_business_owner",
		primaryBusinessOwnerService,
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
				_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServicePrimaryBusinessOwners, permAction, resource)
				return err
			}

			// For general actions, check without specific resource
			_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServicePrimaryBusinessOwners, permAction, nil)
			return err
		}).
		Build()

	controller := &PrimaryBusinessOwnerController{
		CrudController:              crudController,
		primaryBusinessOwnerService: primaryBusinessOwnerService,
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
		// Any custom logic before updating a primary business owner
		return nil
	})

	return controller
}

// Add custom domain-specific methods below this line
// Examples:
//
// // CustomAction GET /api/primary_business_owners/{id}/custom-action
// func (c *PrimaryBusinessOwnerController) CustomAction(ctx http.Context) http.Response {
//     // Get ID from URL
//     id, err := c.ValidateID(ctx, "id")
//     if err != nil {
//         return c.BadRequestResponse(ctx, "Invalid primary business owner ID", nil)
//     }
//
//     // Check permissions
//     if err := c.CheckAuth(ctx, "view", nil); err != nil {
//         return c.ForbiddenResponse(ctx, "Access denied")
//     }
//
//     // Your custom logic here
//     // result, err := c.primaryBusinessOwnerService.CustomMethod(id)
//
//     return c.SuccessResponse(ctx, nil, "Action completed successfully")
// }
