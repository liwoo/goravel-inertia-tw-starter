package applications

import (
	"books-database/app/auth"
	"books-database/app/contracts"
	"books-database/app/events"
	"books-database/app/http/requests"
	"books-database/app/models"
	"books-database/app/services"

	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

// ApplicationController handles API endpoints for application management
type ApplicationController struct {
	*contracts.CrudController[models.Application, *requests.ApplicationCreateRequest, *requests.ApplicationUpdateRequest]
	applicationService *services.ApplicationService
}

// NewApplicationController creates a new application controller with compile-time enforcement
// This ensures all required methods are implemented and configured
func NewApplicationController() *ApplicationController {
	applicationService := services.NewApplicationService()

	// Build controller with compile-time enforcement
	// The builder pattern ensures all required steps are completed
	crudController := contracts.NewCrudController[models.Application, *requests.ApplicationCreateRequest, *requests.ApplicationUpdateRequest](
		"application",
		applicationService,
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
				_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceApplications, permAction, resource)
				return err
			}

			// For general actions, check without specific resource
			_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceApplications, permAction, nil)
			return err
		}).
		Build()

	controller := &ApplicationController{
		CrudController:     crudController,
		applicationService: applicationService,
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
		// Any custom logic before updating a application
		return nil
	})

	return controller
}

// Add custom domain-specific methods below this line

// PublicStore handles public application submissions without authentication
// POST /api/public/applications
func (c *ApplicationController) PublicStore(ctx http.Context) http.Response {
	// Create a new instance of the create request type
	var createReq requests.ApplicationCreateRequest
	if err := ctx.Request().Bind(&createReq); err != nil {
		return c.ValidationErrorResponse(ctx, map[string]interface{}{
			"binding_error": err.Error(),
		})
	}

	// Manual validation using the request's Rules method
	rules := createReq.Rules(ctx)

	// Convert the bound request to validation data
	requestData := createReq.ToCreateData()

	// Validate using Goravel's validator
	validator, err := facades.Validation().Make(requestData, rules)
	if err != nil {
		return c.ValidationErrorResponse(ctx, map[string]interface{}{
			"validation_setup_error": err.Error(),
		})
	}

	if validator.Fails() {
		// Convert validation errors to map[string]interface{}
		errorMap := make(map[string]interface{})
		for key, val := range validator.Errors().All() {
			errorMap[key] = val
		}
		return c.ValidationErrorResponse(ctx, errorMap)
	}

	// Set default status to Pending for public submissions
	requestData["status"] = "Pending"

	// Call service Create
	result, err := c.applicationService.Create(requestData)
	if err != nil {
		facades.Log().Error("Failed to create public application", map[string]interface{}{
			"error": err.Error(),
		})
		return c.BadRequestResponse(ctx, "Failed to create application: "+err.Error(), nil)
	}

	return c.ResourceCreatedResponse(ctx, result, "application")
}

// ApproveApplication POST /api/applications/{id}/approve
func (c *ApplicationController) ApproveApplication(ctx http.Context) http.Response {
	// Get ID from URL
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid application ID", nil)
	}

	// Check permissions - user must have update or manage permission
	if err := c.CheckAuth(ctx, "update", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Get authenticated user ID
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil || user.ID == 0 {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Fetch application details first
	var application models.Application
	if err := facades.Orm().Query().Where("id = ?", id).First(&application); err != nil {
		return c.BadRequestResponse(ctx, "Application not found", nil)
	}

	// Call service method to approve application
	userData, err := c.applicationService.ApproveApplication(id, uint(user.ID))
	if err != nil {
		facades.Log().Error("Failed to approve application", map[string]interface{}{
			"application_id": id,
			"user_id":        user.ID,
			"error":          err.Error(),
		})
		return c.BadRequestResponse(ctx, "Failed to approve application: "+err.Error(), nil)
	}

	// Dispatch ApplicationApproved event for async notification
	if err := facades.Event().Job(&events.ApplicationApproved{}, []event.Arg{
		{Type: "uint", Value: uint(user.ID)},
		{Type: "uint", Value: id},
		{Type: "string", Value: application.Type},
		{Type: "string", Value: application.SME},
		{Type: "string", Value: application.Email},
	}).Dispatch(); err != nil {
		facades.Log().Warningf("Failed to dispatch ApplicationApproved event: %v", err)
	}

	return c.SuccessResponse(ctx, userData, "Application approved successfully")
}

// RejectApplication POST /api/applications/{id}/reject
func (c *ApplicationController) RejectApplication(ctx http.Context) http.Response {
	// Get ID from URL
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid application ID", nil)
	}

	// Check permissions - user must have update or manage permission
	if err := c.CheckAuth(ctx, "update", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Get authenticated user ID
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil || user.ID == 0 {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Fetch application details before rejecting (for event dispatch)
	var application models.Application
	if err := facades.Orm().Query().Where("id = ?", id).First(&application); err != nil {
		return c.BadRequestResponse(ctx, "Application not found", nil)
	}

	// Parse optional rejection reason from request body
	var rejectBody struct {
		Reason string `json:"reason"`
	}
	ctx.Request().Bind(&rejectBody) // Ignore bind errors, reason is optional

	// Call service method to reject application
	if err := c.applicationService.RejectApplication(id, uint(user.ID), rejectBody.Reason); err != nil {
		facades.Log().Error("Failed to reject application", map[string]interface{}{
			"application_id": id,
			"user_id":        user.ID,
			"error":          err.Error(),
		})
		return c.BadRequestResponse(ctx, "Failed to reject application: "+err.Error(), nil)
	}

	// Dispatch ApplicationRejected event for async notification
	if err := facades.Event().Job(&events.ApplicationRejected{}, []event.Arg{
		{Type: "uint", Value: uint(user.ID)},
		{Type: "uint", Value: id},
		{Type: "string", Value: application.Type},
		{Type: "string", Value: application.SME},
		{Type: "string", Value: application.Email},
		{Type: "string", Value: rejectBody.Reason},
	}).Dispatch(); err != nil {
		facades.Log().Warningf("Failed to dispatch ApplicationRejected event: %v", err)
	}

	return c.SuccessResponse(ctx, nil, "Application rejected successfully")
}
