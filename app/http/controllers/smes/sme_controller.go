package smes

import (
	"fmt"

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

// BulkDeactivate POST /api/smes/bulk-deactivate
func (c *SmeController) BulkDeactivate(ctx http.Context) http.Response {
	// Check permissions for bulk update
	if err := c.CheckAuth(ctx, "update", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Parse request body
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := ctx.Request().Bind(&request); err != nil {
		return c.BadRequestResponse(ctx, "Invalid request body", nil)
	}

	if len(request.IDs) == 0 {
		return c.BadRequestResponse(ctx, "No IDs provided", nil)
	}

	// Get current user ID
	var user models.User
	var updatedBy *int
	if err := facades.Auth(ctx).User(&user); err == nil && user.ID > 0 {
		userID := int(user.ID)
		updatedBy = &userID
	}

	// Call service layer
	affected, err := c.smeService.BulkDeactivate(request.IDs, updatedBy)
	if err != nil {
		return c.BadRequestResponse(ctx, "Failed to deactivate SMEs: "+err.Error(), nil)
	}

	return c.SuccessResponse(ctx, map[string]interface{}{
		"affected": affected,
	}, fmt.Sprintf("%d SME(s) deactivated successfully", affected))
}

// BulkActivate POST /api/smes/bulk-activate
func (c *SmeController) BulkActivate(ctx http.Context) http.Response {
	// Check permissions for bulk update
	if err := c.CheckAuth(ctx, "update", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Parse request body
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := ctx.Request().Bind(&request); err != nil {
		return c.BadRequestResponse(ctx, "Invalid request body", nil)
	}

	if len(request.IDs) == 0 {
		return c.BadRequestResponse(ctx, "No IDs provided", nil)
	}

	// Get current user ID
	var user models.User
	var updatedBy *int
	if err := facades.Auth(ctx).User(&user); err == nil && user.ID > 0 {
		userID := int(user.ID)
		updatedBy = &userID
	}

	// Call service layer
	affected, err := c.smeService.BulkActivate(request.IDs, updatedBy)
	if err != nil {
		return c.BadRequestResponse(ctx, "Failed to activate SMEs: "+err.Error(), nil)
	}

	return c.SuccessResponse(ctx, map[string]interface{}{
		"affected": affected,
	}, fmt.Sprintf("%d SME(s) activated successfully", affected))
}

// RecalculateClassification POST /api/smes/{id}/recalculate-classification
func (c *SmeController) RecalculateClassification(ctx http.Context) http.Response {
	// Get ID from URL
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid SME ID", nil)
	}

	// Check permissions
	if err := c.CheckAuth(ctx, "update", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Call service layer
	classification, err := c.smeService.CalculateClassification(id)
	if err != nil {
		return c.BadRequestResponse(ctx, "Failed to calculate classification: "+err.Error(), nil)
	}

	return c.SuccessResponse(ctx, map[string]interface{}{
		"id":             id,
		"classification": classification,
	}, "Classification recalculated successfully")
}

// RecalculateAllClassifications POST /api/smes/recalculate-all-classifications
func (c *SmeController) RecalculateAllClassifications(ctx http.Context) http.Response {
	// Check permissions (require manage permission for bulk operation)
	if err := c.CheckAuth(ctx, "update", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Call service layer
	count, err := c.smeService.RecalculateAllClassifications()
	if err != nil {
		return c.BadRequestResponse(ctx, "Failed to recalculate classifications: "+err.Error(), nil)
	}

	return c.SuccessResponse(ctx, map[string]interface{}{
		"processed": count,
	}, fmt.Sprintf("%d SME(s) classifications recalculated", count))
}

// FetchByUserEmail GET /api/smes/by-email - Returns the SME linked to the authenticated user's email
func (c *SmeController) FetchByUserEmail(ctx http.Context) http.Response {
	// Get the authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil || user.ID == 0 {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Call service layer to get SME by user's email
	sme, err := c.smeService.GetSmeByUserEmail(user.Email)
	if err != nil {
		return c.NotFoundResponse(ctx, "No SME linked to this user")
	}

	if sme == nil {
		return c.NotFoundResponse(ctx, "No SME linked to this user")
	}

	return c.SuccessResponse(ctx, sme, "SME retrieved successfully")
}

// FetchMySme GET /api/my-sme - Returns the full SME details for the authenticated user's linked SME
// This endpoint bypasses admin permissions and allows SME users to view their own SME
func (c *SmeController) FetchMySme(ctx http.Context) http.Response {
	// Get the authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil || user.ID == 0 {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Call service layer to get SME by user's email
	sme, err := c.smeService.GetSmeByUserEmail(user.Email)
	if err != nil {
		return c.NotFoundResponse(ctx, "No SME linked to this user")
	}

	if sme == nil {
		return c.NotFoundResponse(ctx, "No SME linked to this user")
	}

	return c.SuccessResponse(ctx, sme, "SME retrieved successfully")
}

// FetchMySmeOwner GET /api/my-sme/primary-owner - Returns the primary owner of the user's linked SME
func (c *SmeController) FetchMySmeOwner(ctx http.Context) http.Response {
	// Get the authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil || user.ID == 0 {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Get the user's linked SME
	sme, err := c.smeService.GetSmeByUserEmail(user.Email)
	if err != nil || sme == nil {
		return c.NotFoundResponse(ctx, "No SME linked to this user")
	}

	// Get the primary owner for this SME
	owner, err := c.smeService.GetPrimaryBusinessOwner(sme.ID)
	if err != nil {
		return c.NotFoundResponse(ctx, "Primary owner not found")
	}

	if owner == nil {
		return c.SuccessResponse(ctx, nil, "No primary business owner found for this SME")
	}

	return c.SuccessResponse(ctx, owner, "Primary business owner retrieved successfully")
}

// UpdateMySme PUT /api/my-sme - Updates the authenticated user's linked SME
func (c *SmeController) UpdateMySme(ctx http.Context) http.Response {
	// Get the authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil || user.ID == 0 {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Get the user's linked SME
	sme, err := c.smeService.GetSmeByUserEmail(user.Email)
	if err != nil || sme == nil {
		return c.NotFoundResponse(ctx, "No SME linked to this user")
	}

	// Parse request body
	var updateData map[string]interface{}
	if err := ctx.Request().Bind(&updateData); err != nil {
		return c.BadRequestResponse(ctx, "Invalid request body", nil)
	}

	// Update the SME
	updatedSme, err := c.smeService.Update(sme.ID, updateData)
	if err != nil {
		return c.BadRequestResponse(ctx, "Failed to update SME: "+err.Error(), nil)
	}

	return c.SuccessResponse(ctx, updatedSme, "SME updated successfully")
}

// FetchSmesByPrimaryOwnerEmail GET /api/smes/by-owner-email - Returns SMEs where primary owner email matches
// This is used during application approval to show only SMEs that belong to the applicant
func (c *SmeController) FetchSmesByPrimaryOwnerEmail(ctx http.Context) http.Response {
	// Check permissions - user must have read permission on applications
	if err := c.CheckAuth(ctx, "read", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	email := ctx.Request().Query("email", "")
	if email == "" {
		return c.BadRequestResponse(ctx, "Email parameter is required", nil)
	}

	// Find SMEs where primary owner email matches
	smes, err := c.smeService.GetSmesByPrimaryOwnerEmail(email)
	if err != nil {
		return c.BadRequestResponse(ctx, "Failed to fetch SMEs: "+err.Error(), nil)
	}

	return c.SuccessResponse(ctx, smes, "SMEs retrieved successfully")
}

// UpdateMySmeOwner PUT /api/my-sme/primary-owner - Updates the primary owner of the user's linked SME
func (c *SmeController) UpdateMySmeOwner(ctx http.Context) http.Response {
	// Get the authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil || user.ID == 0 {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Get the user's linked SME
	sme, err := c.smeService.GetSmeByUserEmail(user.Email)
	if err != nil || sme == nil {
		return c.NotFoundResponse(ctx, "No SME linked to this user")
	}

	// Get the primary owner for this SME
	owner, err := c.smeService.GetPrimaryBusinessOwner(sme.ID)
	if err != nil || owner == nil {
		return c.NotFoundResponse(ctx, "Primary owner not found")
	}

	// Parse request body
	var updateData map[string]interface{}
	if err := ctx.Request().Bind(&updateData); err != nil {
		return c.BadRequestResponse(ctx, "Invalid request body", nil)
	}

	// Update the primary owner using the service
	updatedOwner, err := c.smeService.UpdatePrimaryBusinessOwner(owner.ID, updateData)
	if err != nil {
		return c.BadRequestResponse(ctx, "Failed to update primary owner: "+err.Error(), nil)
	}

	return c.SuccessResponse(ctx, updatedOwner, "Primary owner updated successfully")
}
