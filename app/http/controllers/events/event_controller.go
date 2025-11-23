package events

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/http/requests"
	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
)

// EventController handles API endpoints for event management
type EventController struct {
	*contracts.CrudController[models.Event, *requests.EventCreateRequest, *requests.EventUpdateRequest]
	eventService *services.EventService
}

// NewEventController creates a new event controller with compile-time enforcement
// This ensures all required methods are implemented and configured
func NewEventController() *EventController {
	eventService := services.NewEventService()

	// Build controller with compile-time enforcement
	// The builder pattern ensures all required steps are completed
	crudController := contracts.NewCrudController[models.Event, *requests.EventCreateRequest, *requests.EventUpdateRequest](
		"event",
		eventService,
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
				_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceEvents, permAction, resource)
				return err
			}

			// For general actions, check without specific resource
			_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceEvents, permAction, nil)
			return err
		}).
		Build()

	controller := &EventController{
		CrudController: crudController,
		eventService:    eventService,
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
		// Any custom logic before updating a event
		return nil
	})

	return controller
}

// Add custom domain-specific methods below this line
// Examples:
//
// // CustomAction GET /api/events/{id}/custom-action
// func (c *EventController) CustomAction(ctx http.Context) http.Response {
//     // Get ID from URL
//     id, err := c.ValidateID(ctx, "id")
//     if err != nil {
//         return c.BadRequestResponse(ctx, "Invalid event ID", nil)
//     }
//
//     // Check permissions
//     if err := c.CheckAuth(ctx, "view", nil); err != nil {
//         return c.ForbiddenResponse(ctx, "Access denied")
//     }
//
//     // Your custom logic here
//     // result, err := c.eventService.CustomMethod(id)
//
//     return c.SuccessResponse(ctx, nil, "Action completed successfully")
// }
