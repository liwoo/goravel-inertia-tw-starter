package smes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/http/requests"
	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
)

// SmeController handles API endpoints for sme management
type SmeController struct {
	*contracts.CrudController[models.Sme, *requests.SmeCreateRequest, *requests.SmeUpdateRequest]
	smeService *services.SmeService
}

// NewSmeController creates a new sme controller with compile-time enforcement
// This ensures all required methods are implemented and configured
func NewSmeController() *SmeController {
	smeService := services.NewSmeService()

	// Build controller with compile-time enforcement
	// The builder pattern ensures all required steps are completed
	crudController := contracts.NewCrudController[models.Sme, *requests.SmeCreateRequest, *requests.SmeUpdateRequest](
		"sme",
		smeService,
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
				_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceSMEs, permAction, resource)
				return err
			}

			// For general actions, check without specific resource
			_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceSMEs, permAction, nil)
			return err
		}).
		Build()

	controller := &SmeController{
		CrudController: crudController,
		smeService:    smeService,
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
		// Any custom logic before updating a sme
		return nil
	})

	return controller
}

// Add custom domain-specific methods below this line

// FetchPrimaryBusinessOwner GET /api/smes/{id}/primary_business_owner
func (c *SmeController) FetchPrimaryBusinessOwner(ctx http.Context) http.Response {
	// Get ID from URL
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid SME ID", nil)
	}

	// Check permissions
	if err := c.CheckAuth(ctx, "view", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Call service layer
	owner, err := c.smeService.GetPrimaryBusinessOwner(id)
	if err != nil {
		return c.NotFoundResponse(ctx, "SME not found")
	}

	// Return the primary business owner
	if owner == nil {
		return c.SuccessResponse(ctx, nil, "No primary business owner found for this SME")
	}

	return c.SuccessResponse(ctx, owner, "Primary business owner retrieved successfully")
}

// FetchAdditionalBusinessMembers GET /api/smes/{id}/additional_business_members
func (c *SmeController) FetchAdditionalBusinessMembers(ctx http.Context) http.Response {
	// Get ID from URL
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid SME ID", nil)
	}

	// Check permissions
	if err := c.CheckAuth(ctx, "view", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Call service layer
	members, err := c.smeService.GetAdditionalBusinessMembers(id)
	if err != nil {
		return c.NotFoundResponse(ctx, "SME not found")
	}

	// Return the additional business members
	if members == nil || len(members) == 0 {
		return c.SuccessResponse(ctx, []interface{}{}, "No additional business members found for this SME")
	}

	return c.SuccessResponse(ctx, members, "Additional business members retrieved successfully")
}

// FetchBusinessFormalisation GET /api/smes/{id}/business_formalisation
func (c *SmeController) FetchBusinessFormalisation(ctx http.Context) http.Response {
	// Get ID from URL
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid SME ID", nil)
	}

	// Check permissions
	if err := c.CheckAuth(ctx, "view", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Call service layer
	formalisation, err := c.smeService.GetBusinessFormalisation(id)
	if err != nil {
		return c.NotFoundResponse(ctx, "SME not found")
	}

	// Return the business formalisation
	if formalisation == nil {
		return c.SuccessResponse(ctx, nil, "No business formalisation found for this SME")
	}

	return c.SuccessResponse(ctx, formalisation, "Business formalisation retrieved successfully")
}

// FetchBusinessEmployeeSummary GET /api/smes/{id}/business_employee_summary
func (c *SmeController) FetchBusinessEmployeeSummary(ctx http.Context) http.Response {
	// Get ID from URL
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid SME ID", nil)
	}

	// Check permissions
	if err := c.CheckAuth(ctx, "view", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Call service layer
	summary, err := c.smeService.GetBusinessEmployeeSummary(id)
	if err != nil {
		return c.NotFoundResponse(ctx, "SME not found")
	}

	// Return the employee summary
	if summary == nil {
		return c.SuccessResponse(ctx, nil, "No employee summary found for this SME")
	}

	return c.SuccessResponse(ctx, summary, "Employee summary retrieved successfully")
}
