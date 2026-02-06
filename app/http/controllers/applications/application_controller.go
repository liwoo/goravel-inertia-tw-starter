package applications

import (
	"strconv"

	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/events"
	"smedi-sme-db/app/http/requests"
	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"

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

	// Fetch application details first to determine the type
	var application models.Application
	if err := facades.Orm().Query().Where("id = ?", id).First(&application); err != nil {
		return c.BadRequestResponse(ctx, "Application not found", nil)
	}

	// Determine sme_id based on application type
	var smeID uint
	if application.Type == models.ApplicationTypeAmendFormalisation {
		// For amendment applications, get sme_id from the application record
		if application.SmeID == nil || *application.SmeID == 0 {
			return c.BadRequestResponse(ctx, "Amendment application is missing SME reference", nil)
		}
		smeID = *application.SmeID
	} else {
		// For signup applications, require sme_id from the request body
		var approvalRequest requests.ApplicationApprovalRequest
		errors, err := ctx.Request().ValidateRequest(&approvalRequest)
		if err != nil {
			return c.BadRequestResponse(ctx, "Validation error", map[string]interface{}{"error": err.Error()})
		}
		if errors != nil {
			errorMap := make(map[string]interface{})
			for key, val := range errors.All() {
				errorMap[key] = val
			}
			return c.BadRequestResponse(ctx, "Validation failed", errorMap)
		}
		smeID = approvalRequest.SmeID
	}

	// Call service method to approve application
	userData, err := c.applicationService.ApproveApplication(id, smeID, uint(user.ID))
	if err != nil {
		facades.Log().Error("Failed to approve application", map[string]interface{}{
			"application_id": id,
			"sme_id":         smeID,
			"user_id":        user.ID,
			"error":          err.Error(),
		})
		return c.BadRequestResponse(ctx, "Failed to approve application: "+err.Error(), nil)
	}

	// Dispatch ApplicationApproved event for async notification (not for signup - handled differently)
	if application.Type != models.ApplicationTypeSignup {
		if err := facades.Event().Job(&events.ApplicationApproved{}, []event.Arg{
			{Type: "uint", Value: uint(user.ID)},
			{Type: "uint", Value: id},
			{Type: "string", Value: application.Type},
			{Type: "string", Value: application.SME},
			{Type: "string", Value: application.Email},
		}).Dispatch(); err != nil {
			facades.Log().Warningf("Failed to dispatch ApplicationApproved event: %v", err)
		}
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

	// Parse optional rejection reason
	var rejectRequest requests.ApplicationRejectRequest
	ctx.Request().Bind(&rejectRequest) // Ignore bind errors, reason is optional

	// Call service method to reject application
	if err := c.applicationService.RejectApplication(id, uint(user.ID), rejectRequest.Reason); err != nil {
		facades.Log().Error("Failed to reject application", map[string]interface{}{
			"application_id": id,
			"user_id":        user.ID,
			"error":          err.Error(),
		})
		return c.BadRequestResponse(ctx, "Failed to reject application: "+err.Error(), nil)
	}

	// Dispatch ApplicationRejected event for async notification (not for signup - no user exists yet)
	if application.Type != models.ApplicationTypeSignup {
		if err := facades.Event().Job(&events.ApplicationRejected{}, []event.Arg{
			{Type: "uint", Value: uint(user.ID)},
			{Type: "uint", Value: id},
			{Type: "string", Value: application.Type},
			{Type: "string", Value: application.SME},
			{Type: "string", Value: application.Email},
			{Type: "string", Value: rejectRequest.Reason},
		}).Dispatch(); err != nil {
			facades.Log().Warningf("Failed to dispatch ApplicationRejected event: %v", err)
		}
	}

	return c.SuccessResponse(ctx, nil, "Application rejected successfully")
}

// SubmitAmendment POST /api/applications/amendment
// Allows SME users to submit formalisation amendment requests
func (c *ApplicationController) SubmitAmendment(ctx http.Context) http.Response {
	// Get authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil || user.ID == 0 {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Validate request
	var amendmentRequest requests.FormalisationAmendmentRequest
	errors, err := ctx.Request().ValidateRequest(&amendmentRequest)
	if err != nil {
		return c.BadRequestResponse(ctx, "Validation error", map[string]interface{}{"error": err.Error()})
	}
	if errors != nil {
		errorMap := make(map[string]interface{})
		for key, val := range errors.All() {
			errorMap[key] = val
		}
		return c.BadRequestResponse(ctx, "Validation failed", errorMap)
	}

	// Build the proposed changes struct
	proposed := services.FormalisationAmendmentData{
		RegistrationNumber:      amendmentRequest.RegistrationNumber,
		TaxIdentificationNumber: amendmentRequest.TaxIdentificationNumber,
		HasBankAccount:          amendmentRequest.HasBankAccount,
		HasTaxClarification:     amendmentRequest.HasTaxClarification,
		IsRegisteredForVat:      amendmentRequest.IsRegisteredForVat,
		IsMemberOfAssociation:   amendmentRequest.IsMemberOfAssociation,
		IsAffiliated:            amendmentRequest.IsAffiliated,
		HasExportLicense:        amendmentRequest.HasExportLicense,
		HasAccessedBds:          amendmentRequest.HasAccessedBds,
		AnnualTurnover:          amendmentRequest.AnnualTurnover,
		EstimatedValueOfAssets:  amendmentRequest.EstimatedValueOfAssets,
		ChangeReason:            amendmentRequest.ChangeReason,
	}

	// Create the amendment application
	application, err := c.applicationService.CreateAmendmentApplication(
		amendmentRequest.SmeID,
		uint(user.ID),
		proposed,
	)
	if err != nil {
		facades.Log().Error("Failed to create amendment application", map[string]interface{}{
			"sme_id":  amendmentRequest.SmeID,
			"user_id": user.ID,
			"error":   err.Error(),
		})
		return c.BadRequestResponse(ctx, "Failed to submit amendment: "+err.Error(), nil)
	}

	return c.ResourceCreatedResponse(ctx, application, "amendment application")
}

// CheckPendingAmendment GET /api/applications/amendment/pending/:smeId
// Checks if an SME has a pending formalisation amendment application
func (c *ApplicationController) CheckPendingAmendment(ctx http.Context) http.Response {
	// Get authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil || user.ID == 0 {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Get SME ID from path parameter
	smeID, err := strconv.ParseUint(ctx.Request().Input("smeId"), 10, 32)
	if err != nil || smeID == 0 {
		return c.BadRequestResponse(ctx, "Invalid MSME ID", nil)
	}

	// Check for pending amendment
	hasPending, err := c.applicationService.HasPendingAmendment(uint(smeID))
	if err != nil {
		facades.Log().Error("Failed to check pending amendment", map[string]interface{}{
			"sme_id": smeID,
			"error":  err.Error(),
		})
		return c.InternalErrorResponse(ctx, "Failed to check pending amendment status")
	}

	return c.SuccessResponse(ctx, map[string]interface{}{
		"has_pending_amendment": hasPending,
	}, "Pending amendment status retrieved")
}
