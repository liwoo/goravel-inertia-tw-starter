package books

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/http/requests"
	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
)

// BookController handles API endpoints for book management
type BookController struct {
	*contracts.CrudController[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest]
	bookService *services.BookService
}

// NewBookController creates a new book controller with compile-time enforcement
// This ensures all required methods are implemented and configured
func NewBookController() *BookController {
	bookService := services.NewBookService()

	// Build controller with compile-time enforcement
	// The builder pattern ensures all required steps are completed
	crudController := contracts.NewCrudController[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest](
		"book",
		bookService,
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
				_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceBooks, permAction, resource)
				return err
			}

			// For general actions, check without specific resource
			_, err := scopedHelper.RequireScopedPermission(ctx, auth.ServiceBooks, permAction, nil)
			return err
		}).
		Build()

	controller := &BookController{
		CrudController: crudController,
		bookService:    bookService,
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
		// Any custom logic before updating a book
		return nil
	})

	return controller
}

// ============================================================================
// CRUD Methods with Swagger annotations
// ============================================================================

// Index godoc
// @Summary      List all books
// @Description  Get paginated list of books with filtering, sorting, and search
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        page      query  int     false  "Page number" default(1)
// @Param        pageSize  query  int     false  "Items per page" default(20) Enums(5, 10, 20, 30, 50, 100)
// @Param        search    query  string  false  "Search query"
// @Param        sort      query  string  false  "Sort field"
// @Param        direction query  string  false  "Sort direction" Enums(ASC, DESC)
// @Success      200  {object}  contracts.ResponseFormat{data=contracts.PaginatedResponse{data=[]models.Book}}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Router       /books [get]
func (c *BookController) Index(ctx http.Context) http.Response {
	return c.CrudController.Index(ctx)
}

// Show godoc
// @Summary      Get book by ID
// @Description  Retrieve a specific book
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Book ID"
// @Success      200  {object}  contracts.ResponseFormat{data=models.Book}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      404  {object}  contracts.ResponseFormat
// @Router       /books/{id} [get]
func (c *BookController) Show(ctx http.Context) http.Response {
	return c.CrudController.Show(ctx)
}

// Store godoc
// @Summary      Create a new book
// @Description  Create a book with validation
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book  body  requests.BookCreateRequest  true  "book data"
// @Success      201  {object}  contracts.ResponseFormat{data=models.Book}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      422  {object}  contracts.ResponseFormat
// @Security     BearerAuth
// @Router       /books [post]
func (c *BookController) Store(ctx http.Context) http.Response {
	return c.CrudController.Store(ctx)
}

// Update godoc
// @Summary      Update a book
// @Description  Update an existing book
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id    path  int  true  "Book ID"
// @Param        book  body  requests.BookUpdateRequest  true  "book data"
// @Success      200  {object}  contracts.ResponseFormat{data=models.Book}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      404  {object}  contracts.ResponseFormat
// @Failure      422  {object}  contracts.ResponseFormat
// @Security     BearerAuth
// @Router       /books/{id} [put]
func (c *BookController) Update(ctx http.Context) http.Response {
	return c.CrudController.Update(ctx)
}

// Delete godoc
// @Summary      Delete a book
// @Description  Delete a book by ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Book ID"
// @Success      204  "Book deleted"
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Failure      404  {object}  contracts.ResponseFormat
// @Security     BearerAuth
// @Router       /books/{id} [delete]
func (c *BookController) Delete(ctx http.Context) http.Response {
	return c.CrudController.Delete(ctx)
}

// Search godoc
// @Summary      Search books
// @Description  Full-text search for books
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        q         query  string  true   "Search query (min 2 chars)"
// @Param        page      query  int     false  "Page number"
// @Param        pageSize  query  int     false  "Items per page"
// @Success      200  {object}  contracts.ResponseFormat{data=contracts.PaginatedResponse{data=[]models.Book}}
// @Failure      400  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Router       /books/search [get]
func (c *BookController) Search(ctx http.Context) http.Response {
	return c.CrudController.Search(ctx)
}

// FilterMetadata godoc
// @Summary      Get book filter metadata
// @Description  Returns available filters for books
// @Tags         books
// @Accept       json
// @Produce      json
// @Success      200  {object}  contracts.ResponseFormat
// @Failure      403  {object}  contracts.ResponseFormat
// @Router       /books/filters [get]
func (c *BookController) FilterMetadata(ctx http.Context) http.Response {
	return c.CrudController.FilterMetadata(ctx)
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

// GetByISBN GET /api/books/isbn/{isbn} - Get a book by ISBN
func (c *BookController) GetByISBN(ctx http.Context) http.Response {
	// Check permissions
	if err := c.CheckAuth(ctx, "view", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Get ISBN from URL
	isbn := ctx.Request().Route("isbn")
	if isbn == "" {
		return c.BadRequestResponse(ctx, "ISBN is required", nil)
	}

	// Get the book
	book, err := c.bookService.GetByISBN(isbn)
	if err != nil {
		return c.NotFoundResponse(ctx, "Book not found")
	}

	return c.SuccessResponse(ctx, book, "Book retrieved successfully")
}

// GetByAuthor GET /api/books/author/{author} - Get books by author
func (c *BookController) GetByAuthor(ctx http.Context) http.Response {
	// Check permissions
	if err := c.CheckAuth(ctx, "viewAny", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}

	// Get author from URL
	author := ctx.Request().Route("author")
	if author == "" {
		return c.BadRequestResponse(ctx, "Author is required", nil)
	}

	// Validate pagination request
	req, err := c.ValidatePaginationRequest(ctx)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid pagination parameters", nil)
	}

	// Get books by author
	result, err := c.bookService.GetByAuthor(author, *req)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve books")
	}

	// Build response
	response := c.BuildPaginatedResponse(result, req)
	return c.SuccessResponse(ctx, response, "Books retrieved successfully")
}
