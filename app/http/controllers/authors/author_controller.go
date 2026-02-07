package authors

import (
	"books-database/app/auth"
	"books-database/app/contracts"
	"books-database/app/http/requests"
	"books-database/app/models"
	"books-database/app/services"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

// AuthorController handles API endpoints for author management
type AuthorController struct {
	*contracts.CrudController[models.Author, *requests.AuthorCreateRequest, *requests.AuthorUpdateRequest]
	authorService *services.AuthorService
}

// NewAuthorController creates a new author controller with compile-time enforcement
// This ensures all required methods are implemented and configured
func NewAuthorController() *AuthorController {
	authorService := services.NewAuthorService()

	// Build controller with compile-time enforcement
	// The builder pattern ensures all required steps are completed
	crudController := contracts.NewCrudController[models.Author, *requests.AuthorCreateRequest, *requests.AuthorUpdateRequest](
		"author",
		authorService,
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
				_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceAuthors, permAction, resource)
				return err
			}

			// For general actions, check without specific resource
			_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceAuthors, permAction, nil)
			return err
		}).
		Build()

	controller := &AuthorController{
		CrudController: crudController,
		authorService:  authorService,
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
		var user models.User
		if err := facades.Auth(ctx).User(&user); err == nil && user.ID > 0 {
			data["updated_by"] = user.ID
		}
		return nil
	})

	return controller
}

// ============================================================================
// CRUD Methods with Swagger annotations
// ============================================================================

// Index godoc
// @Summary      List all authors
// @Description  Get paginated list of authors with filtering, sorting, and search
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        page      query  int     false  "Page number" default(1)
// @Param        pageSize  query  int     false  "Items per page" default(20) Enums(5, 10, 20, 30, 50, 100)
// @Param        search    query  string  false  "Search query"
// @Param        sort      query  string  false  "Sort field"
// @Param        direction query  string  false  "Sort direction" Enums(ASC, DESC)
// @Success      200  {object}  contracts.ResponseFormat{data=contracts.PaginatedResponse{data=[]models.Author}}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Router       /authors [get]
func (c *AuthorController) Index(ctx http.Context) http.Response {
	return c.CrudController.Index(ctx)
}

// Show godoc
// @Summary      Get author by ID
// @Description  Retrieve a specific author
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Author ID"
// @Success      200  {object}  contracts.ResponseFormat{data=models.Author}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      404  {object}  contracts.ResponseFormat
// @Router       /authors/{id} [get]
func (c *AuthorController) Show(ctx http.Context) http.Response {
	return c.CrudController.Show(ctx)
}

// Store godoc
// @Summary      Create a new author
// @Description  Create an author with validation
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        author  body  requests.AuthorCreateRequest  true  "author data"
// @Success      201  {object}  contracts.ResponseFormat{data=models.Author}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      422  {object}  contracts.ResponseFormat
// @Security     BearerAuth
// @Router       /authors [post]
func (c *AuthorController) Store(ctx http.Context) http.Response {
	return c.CrudController.Store(ctx)
}

// Update godoc
// @Summary      Update an author
// @Description  Update an existing author
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        id      path  int  true  "Author ID"
// @Param        author  body  requests.AuthorUpdateRequest  true  "author data"
// @Success      200  {object}  contracts.ResponseFormat{data=models.Author}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      404  {object}  contracts.ResponseFormat
// @Failure      422  {object}  contracts.ResponseFormat
// @Security     BearerAuth
// @Router       /authors/{id} [put]
func (c *AuthorController) Update(ctx http.Context) http.Response {
	return c.CrudController.Update(ctx)
}

// Delete godoc
// @Summary      Delete an author
// @Description  Delete an author by ID
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Author ID"
// @Success      204  "Author deleted"
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      404  {object}  contracts.ResponseFormat
// @Security     BearerAuth
// @Router       /authors/{id} [delete]
func (c *AuthorController) Delete(ctx http.Context) http.Response {
	return c.CrudController.Delete(ctx)
}

// Search godoc
// @Summary      Search authors
// @Description  Full-text search for authors
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        q         query  string  true   "Search query (min 2 chars)"
// @Param        page      query  int     false  "Page number"
// @Param        pageSize  query  int     false  "Items per page"
// @Success      200  {object}  contracts.ResponseFormat{data=contracts.PaginatedResponse{data=[]models.Author}}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Router       /authors/search [get]
func (c *AuthorController) Search(ctx http.Context) http.Response {
	return c.CrudController.Search(ctx)
}

// FilterMetadata godoc
// @Summary      Get author filter metadata
// @Description  Returns available filters for authors
// @Tags         authors
// @Accept       json
// @Produce      json
// @Success      200  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Router       /authors/filters [get]
func (c *AuthorController) FilterMetadata(ctx http.Context) http.Response {
	return c.CrudController.FilterMetadata(ctx)
}

// Statistics GET /api/authors/statistics - Get author statistics
func (c *AuthorController) Statistics(ctx http.Context) http.Response {
	// Check permissions
	if err := c.CheckAuth(ctx, "viewAny", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Get statistics
	stats, err := c.authorService.GetAuthorStatistics()
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve statistics")
	}

	return c.SuccessResponse(ctx, stats, "Statistics retrieved successfully")
}
