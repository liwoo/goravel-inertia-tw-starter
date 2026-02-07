package configs

import (
	"books-database/app/auth"
	"books-database/app/contracts"
	"books-database/app/http/requests"
	"books-database/app/models"
	"books-database/app/services"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

// ConfigController handles API endpoints for config management
type ConfigController struct {
	*contracts.CrudController[models.Config, *requests.ConfigCreateRequest, *requests.ConfigUpdateRequest]
	configService *services.ConfigService
}

// NewConfigController creates a new config controller with compile-time enforcement
// This ensures all required methods are implemented and configured
func NewConfigController() *ConfigController {
	configService := services.NewConfigService()

	// Build controller with compile-time enforcement
	// The builder pattern ensures all required steps are completed
	crudController := contracts.NewCrudController[models.Config, *requests.ConfigCreateRequest, *requests.ConfigUpdateRequest](
		"config",
		configService,
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
				_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceConfig, permAction, resource)
				return err
			}

			// For general actions, check without specific resource
			_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceConfig, permAction, nil)
			return err
		}).
		Build()

	controller := &ConfigController{
		CrudController: crudController,
		configService:  configService,
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
		// Any custom logic before updating a config
		return nil
	})

	return controller
}

// Add custom domain-specific methods below this line
// Examples:
//
// // CustomAction GET /api/configs/{id}/custom-action
// func (c *ConfigController) CustomAction(ctx http.Context) http.Response {
//     // Get ID from URL
//     id, err := c.ValidateID(ctx, "id")
//     if err != nil {
//         return c.BadRequestResponse(ctx, "Invalid config ID", nil)
//     }
//
//     // Check permissions
//     if err := c.CheckAuth(ctx, "view", nil); err != nil {
//         return c.ForbiddenResponse(ctx, "Access denied")
//     }
//
//     // Your custom logic here
//     // result, err := c.configService.CustomMethod(id)
//
//     return c.SuccessResponse(ctx, nil, "Action completed successfully")
// }
