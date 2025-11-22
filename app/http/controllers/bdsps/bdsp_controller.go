package bdsps

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/http/requests"
	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
)

// BdspController handles API endpoints for bdsp management
type BdspController struct {
	*contracts.CrudController[models.Bdsp, *requests.BdspCreateRequest, *requests.BdspUpdateRequest]
	bdspService *services.BdspService
}

// NewBdspController creates a new bdsp controller with compile-time enforcement
// This ensures all required methods are implemented and configured
func NewBdspController() *BdspController {
	bdspService := services.NewBdspService()

	// Build controller with compile-time enforcement
	// The builder pattern ensures all required steps are completed
	crudController := contracts.NewCrudController[models.Bdsp, *requests.BdspCreateRequest, *requests.BdspUpdateRequest](
		"bdsp",
		bdspService,
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
				_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceBdsps, permAction, resource)
				return err
			}

			// For general actions, check without specific resource
			_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceBdsps, permAction, nil)
			return err
		}).
		Build()

	controller := &BdspController{
		CrudController: crudController,
		bdspService:    bdspService,
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
		// Any custom logic before updating a bdsp
		return nil
	})

	return controller
}

// ============================================================================
// CRUD Methods with Swagger annotations
// ============================================================================

// Index godoc
// @Summary      List all bdsps
// @Description  Get paginated list of bdsps with filtering, sorting, and search
// @Tags         bdsps
// @Accept       json
// @Produce      json
// @Param        page      query  int     false  "Page number" default(1)
// @Param        pageSize  query  int     false  "Items per page" default(20) Enums(5, 10, 20, 30, 50, 100)
// @Param        search    query  string  false  "Search query"
// @Param        sort      query  string  false  "Sort field"
// @Param        direction query  string  false  "Sort direction" Enums(ASC, DESC)
// @Success      200  {object}  contracts.ResponseFormat{data=contracts.PaginatedResponse{data=[]models.Bdsp}}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Router       /bdsps [get]
func (c *BdspController) Index(ctx http.Context) http.Response {
	return c.CrudController.Index(ctx)
}

// Show godoc
// @Summary      Get bdsp by ID
// @Description  Retrieve a specific bdsp
// @Tags         bdsps
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Bdsp ID"
// @Success      200  {object}  contracts.ResponseFormat{data=models.Bdsp}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      404  {object}  contracts.ResponseFormat
// @Router       /bdsps/{id} [get]
func (c *BdspController) Show(ctx http.Context) http.Response {
	return c.CrudController.Show(ctx)
}

// Store godoc
// @Summary      Create a new bdsp
// @Description  Create a bdsp with validation
// @Tags         bdsps
// @Accept       json
// @Produce      json
// @Param        bdsp  body  requests.BdspCreateRequest  true  "bdsp data"
// @Success      201  {object}  contracts.ResponseFormat{data=models.Bdsp}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      422  {object}  contracts.ResponseFormat
// @Security     BearerAuth
// @Router       /bdsps [post]
func (c *BdspController) Store(ctx http.Context) http.Response {
	return c.CrudController.Store(ctx)
}

// Update godoc
// @Summary      Update a bdsp
// @Description  Update an existing bdsp
// @Tags         bdsps
// @Accept       json
// @Produce      json
// @Param        id    path  int  true  "Bdsp ID"
// @Param        bdsp  body  requests.BdspUpdateRequest  true  "bdsp data"
// @Success      200  {object}  contracts.ResponseFormat{data=models.Bdsp}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      404  {object}  contracts.ResponseFormat
// @Failure      422  {object}  contracts.ResponseFormat
// @Security     BearerAuth
// @Router       /bdsps/{id} [put]
func (c *BdspController) Update(ctx http.Context) http.Response {
	return c.CrudController.Update(ctx)
}

// Delete godoc
// @Summary      Delete a bdsp
// @Description  Delete a bdsp by ID
// @Tags         bdsps
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Bdsp ID"
// @Success      204  "Bdsp deleted"
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      404  {object}  contracts.ResponseFormat
// @Security     BearerAuth
// @Router       /bdsps/{id} [delete]
func (c *BdspController) Delete(ctx http.Context) http.Response {
	return c.CrudController.Delete(ctx)
}

// Search godoc
// @Summary      Search bdsps
// @Description  Full-text search for bdsps
// @Tags         bdsps
// @Accept       json
// @Produce      json
// @Param        q         query  string  true   "Search query (min 2 chars)"
// @Param        page      query  int     false  "Page number"
// @Param        pageSize  query  int     false  "Items per page"
// @Success      200  {object}  contracts.ResponseFormat{data=contracts.PaginatedResponse{data=[]models.Bdsp}}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Router       /bdsps/search [get]
func (c *BdspController) Search(ctx http.Context) http.Response {
	return c.CrudController.Search(ctx)
}

// FilterMetadata godoc
// @Summary      Get bdsp filter metadata
// @Description  Returns available filters for bdsps
// @Tags         bdsps
// @Accept       json
// @Produce      json
// @Success      200  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Router       /bdsps/filters [get]
func (c *BdspController) FilterMetadata(ctx http.Context) http.Response {
	return c.CrudController.FilterMetadata(ctx)
}

// Statistics GET /api/bdsps/statistics - Get bdsp statistics
func (c *BdspController) Statistics(ctx http.Context) http.Response {
	// Check permissions
	if err := c.CheckAuth(ctx, "viewAny", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Get statistics
	stats, err := c.bdspService.GetBdspStatistics()
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve statistics")
	}

	return c.SuccessResponse(ctx, stats, "Statistics retrieved successfully")
}

// Add custom domain-specific methods below this line
// Examples:
//
// // CustomAction GET /api/bdsps/{id}/custom-action
// func (c *BdspController) CustomAction(ctx http.Context) http.Response {
//     // Get ID from URL
//     id, err := c.ValidateID(ctx, "id")
//     if err != nil {
//         return c.BadRequestResponse(ctx, "Invalid bdsp ID", nil)
//     }
//
//     // Check permissions
//     if err := c.CheckAuth(ctx, "view", nil); err != nil {
//         return c.ForbiddenResponse(ctx, "Access denied")
//     }
//
//     // Your custom logic here
//     // result, err := c.bdspService.CustomMethod(id)
//
//     return c.SuccessResponse(ctx, nil, "Action completed successfully")
// }
