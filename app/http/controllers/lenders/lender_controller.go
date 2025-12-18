package lenders

import (
	"starter-project/app/auth"
	"starter-project/app/contracts"
	"starter-project/app/http/requests"
	"starter-project/app/models"
	"starter-project/app/services"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

// LenderController handles API endpoints for lender management
type LenderController struct {
	*contracts.CrudController[models.Lender, *requests.LenderCreateRequest, *requests.LenderUpdateRequest]
	lenderService *services.LenderService
}

// NewLenderController creates a new lender controller with compile-time enforcement
// This ensures all required methods are implemented and configured
func NewLenderController() *LenderController {
	lenderService := services.NewLenderService()

	// Build controller with compile-time enforcement
	// The builder pattern ensures all required steps are completed
	crudController := contracts.NewCrudController[models.Lender, *requests.LenderCreateRequest, *requests.LenderUpdateRequest](
		"lender",
		lenderService,
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
				_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceLenders, permAction, resource)
				return err
			}

			// For general actions, check without specific resource
			_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceLenders, permAction, nil)
			return err
		}).
		Build()

	controller := &LenderController{
		CrudController: crudController,
		lenderService:  lenderService,
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
		// Any custom logic before updating a lender
		return nil
	})

	return controller
}

// ============================================================================
// CRUD Methods with Swagger annotations
// ============================================================================

// Index godoc
// @Summary      List all lenders
// @Description  Get paginated list of lenders with filtering, sorting, and search
// @Tags         lenders
// @Accept       json
// @Produce      json
// @Param        page      query  int     false  "Page number" default(1)
// @Param        pageSize  query  int     false  "Items per page" default(20) Enums(5, 10, 20, 30, 50, 100)
// @Param        search    query  string  false  "Search query"
// @Param        sort      query  string  false  "Sort field"
// @Param        direction query  string  false  "Sort direction" Enums(ASC, DESC)
// @Success      200  {object}  contracts.ResponseFormat{data=contracts.PaginatedResponse{data=[]models.Lender}}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Router       /lenders [get]
func (c *LenderController) Index(ctx http.Context) http.Response {
	return c.CrudController.Index(ctx)
}

// Show godoc
// @Summary      Get lender by ID
// @Description  Retrieve a specific lender
// @Tags         lenders
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Lender ID"
// @Success      200  {object}  contracts.ResponseFormat{data=models.Lender}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      404  {object}  contracts.ResponseFormat
// @Router       /lenders/{id} [get]
func (c *LenderController) Show(ctx http.Context) http.Response {
	return c.CrudController.Show(ctx)
}

// Store godoc
// @Summary      Create a new lender
// @Description  Create a lender with validation
// @Tags         lenders
// @Accept       json
// @Produce      json
// @Param        lender  body  requests.LenderCreateRequest  true  "lender data"
// @Success      201  {object}  contracts.ResponseFormat{data=models.Lender}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      422  {object}  contracts.ResponseFormat
// @Security     BearerAuth
// @Router       /lenders [post]
func (c *LenderController) Store(ctx http.Context) http.Response {
	return c.CrudController.Store(ctx)
}

// Update godoc
// @Summary      Update a lender
// @Description  Update an existing lender
// @Tags         lenders
// @Accept       json
// @Produce      json
// @Param        id    path  int  true  "Lender ID"
// @Param        lender  body  requests.LenderUpdateRequest  true  "lender data"
// @Success      200  {object}  contracts.ResponseFormat{data=models.Lender}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      404  {object}  contracts.ResponseFormat
// @Failure      422  {object}  contracts.ResponseFormat
// @Security     BearerAuth
// @Router       /lenders/{id} [put]
func (c *LenderController) Update(ctx http.Context) http.Response {
	return c.CrudController.Update(ctx)
}

// Delete godoc
// @Summary      Delete a lender
// @Description  Delete a lender by ID
// @Tags         lenders
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Lender ID"
// @Success      204  "Lender deleted"
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      404  {object}  contracts.ResponseFormat
// @Security     BearerAuth
// @Router       /lenders/{id} [delete]
func (c *LenderController) Delete(ctx http.Context) http.Response {
	return c.CrudController.Delete(ctx)
}

// Search godoc
// @Summary      Search lenders
// @Description  Full-text search for lenders
// @Tags         lenders
// @Accept       json
// @Produce      json
// @Param        q         query  string  true   "Search query (min 2 chars)"
// @Param        page      query  int     false  "Page number"
// @Param        pageSize  query  int     false  "Items per page"
// @Success      200  {object}  contracts.ResponseFormat{data=contracts.PaginatedResponse{data=[]models.Lender}}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Router       /lenders/search [get]
func (c *LenderController) Search(ctx http.Context) http.Response {
	return c.CrudController.Search(ctx)
}

// FilterMetadata godoc
// @Summary      Get lender filter metadata
// @Description  Returns available filters for lenders
// @Tags         lenders
// @Accept       json
// @Produce      json
// @Success      200  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Router       /lenders/filters [get]
func (c *LenderController) FilterMetadata(ctx http.Context) http.Response {
	return c.CrudController.FilterMetadata(ctx)
}

// Add custom domain-specific methods below this line
// Examples:
//
// // CustomAction GET /api/lenders/{id}/custom-action
// func (c *LenderController) CustomAction(ctx http.Context) http.Response {
//     // Get ID from URL
//     id, err := c.ValidateID(ctx, "id")
//     if err != nil {
//         return c.BadRequestResponse(ctx, "Invalid lender ID", nil)
//     }
//
//     // Check permissions
//     if err := c.CheckAuth(ctx, "view", nil); err != nil {
//         return c.ForbiddenResponse(ctx, "Access denied")
//     }
//
//     // Your custom logic here
//     // result, err := c.lenderService.CustomMethod(id)
//
//     return c.SuccessResponse(ctx, nil, "Action completed successfully")
// }
