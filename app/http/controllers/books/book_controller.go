package books

import (
	"github.com/goravel/framework/contracts/http"
	"players/app/auth"
	"players/app/contracts"
	"players/app/http/requests"
	"players/app/models"
	"players/app/services"
)

// BookController handles API endpoints for book management
type BookController struct {
	*contracts.StaticEnforcedController[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest]
	bookService *services.BookService
}

// NewBookController creates a new book controller with compile-time enforcement
// This ensures all required methods are implemented and configured
func NewBookController() *BookController {
	bookService := services.NewBookService()

	// Build controller with compile-time enforcement
	// The builder pattern ensures all required steps are completed
	staticController := contracts.NewStaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest](
		"book",
		bookService,
	).
		ValidateCreateRequest().
		ValidateUpdateRequest().
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
				_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceBooks, permAction, resource)
				return err
			}

			// For general actions, check without specific resource
			_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceBooks, permAction, nil)
			return err
		}).
		Build()

	controller := &BookController{
		StaticEnforcedController: staticController,
		bookService:              bookService,
	}

	// Set custom hooks
	controller.SetBeforeStore(func(ctx http.Context, data map[string]interface{}) error {
		// Any custom logic before creating a book
		return nil
	})

	controller.SetBeforeUpdate(func(ctx http.Context, id uint, data map[string]interface{}) error {
		// Any custom logic before updating a book
		return nil
	})

	return controller
}

// Borrow POST /api/books/{id}/borrow - Borrow a book
func (c *BookController) Borrow(ctx http.Context) http.Response {
	// Get book ID from URL
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid book ID", nil)
	}

	// Check permissions
	if err := c.CheckAuth(ctx, "update", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Borrow the book
	err = c.bookService.BorrowBook(id)
	if err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	return c.SuccessResponse(ctx, map[string]interface{}{
		"message": "Book borrowed successfully",
	}, "Book borrowed")
}

// Return POST /api/books/{id}/return - Return a borrowed book
func (c *BookController) Return(ctx http.Context) http.Response {
	// Get book ID from URL
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid book ID", nil)
	}

	// Check permissions
	if err := c.CheckAuth(ctx, "update", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Return the book
	err = c.bookService.ReturnBook(id)
	if err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	return c.SuccessResponse(ctx, map[string]interface{}{
		"message": "Book returned successfully",
	}, "Book returned")
}

// Available GET /api/books/available - Get available books
func (c *BookController) Available(ctx http.Context) http.Response {
	// Check permissions
	if err := c.CheckAuth(ctx, "viewAny", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Validate pagination request
	req, err := c.ValidatePaginationRequest(ctx)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid pagination parameters", nil)
	}

	// Get available books
	result, err := c.bookService.GetAvailable(*req)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve available books")
	}

	// Build response
	response := c.BuildPaginatedResponse(result, req)
	return c.SuccessResponse(ctx, response, "Available books retrieved successfully")
}

// Statistics GET /api/books/statistics - Get book statistics
func (c *BookController) Statistics(ctx http.Context) http.Response {
	// Check permissions
	if err := c.CheckAuth(ctx, "viewAny", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Get statistics
	stats, err := c.bookService.GetBookStatistics()
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve statistics")
	}

	return c.SuccessResponse(ctx, stats, "Statistics retrieved successfully")
}
