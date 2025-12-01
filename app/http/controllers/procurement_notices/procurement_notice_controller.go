package procurement_notices

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/http/requests"
	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
)

// ProcurementNoticeController handles API endpoints for procurement notice management
type ProcurementNoticeController struct {
	*contracts.CrudController[models.ProcurementNotice, *requests.ProcurementNoticeCreateRequest, *requests.ProcurementNoticeUpdateRequest]
	procurementNoticeService *services.ProcurementNoticeService
}

// NewProcurementNoticeController creates a new procurement notice controller with compile-time enforcement
// This ensures all required methods are implemented and configured
func NewProcurementNoticeController() *ProcurementNoticeController {
	procurementNoticeService := services.NewProcurementNoticeService()

	// Build controller with compile-time enforcement
	// The builder pattern ensures all required steps are completed
	crudController := contracts.NewCrudController[models.ProcurementNotice, *requests.ProcurementNoticeCreateRequest, *requests.ProcurementNoticeUpdateRequest](
		"procurement_notice",
		procurementNoticeService,
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
				_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceProcurementNotices, permAction, resource)
				return err
			}

			// For general actions, check without specific resource
			_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceProcurementNotices, permAction, nil)
			return err
		}).
		Build()

	controller := &ProcurementNoticeController{
		CrudController: crudController,
		procurementNoticeService:    procurementNoticeService,
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
		// Any custom logic before updating a procurement_notice
		return nil
	})

	return controller
}

// TogglePublish toggles the is_published status of a procurement notice
// POST /api/procurement-notices/{id}/toggle-publish
func (c *ProcurementNoticeController) TogglePublish(ctx http.Context) http.Response {
	// Get ID from URL
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid procurement notice ID", nil)
	}

	// Check update permissions
	if err := c.CheckAuth(ctx, "update", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Get current procurement notice
	result, err := c.procurementNoticeService.GetByID(id)
	if err != nil {
		return c.NotFoundResponse(ctx, "Procurement notice not found")
	}

	// Type assert to get the model
	item, ok := result.(*models.ProcurementNotice)
	if !ok {
		return c.BadRequestResponse(ctx, "Invalid procurement notice data", nil)
	}

	// Toggle the is_published status
	newStatus := !item.IsPublished

	// Update only the is_published field
	updateData := map[string]interface{}{
		"is_published": newStatus,
	}

	_, err = c.procurementNoticeService.Update(id, updateData)
	if err != nil {
		return c.BadRequestResponse(ctx, "Failed to update publish status", nil)
	}

	statusText := "unpublished"
	if newStatus {
		statusText = "published"
	}

	return c.SuccessResponse(ctx, map[string]interface{}{
		"id":           id,
		"is_published": newStatus,
	}, "Procurement notice "+statusText+" successfully")
}
