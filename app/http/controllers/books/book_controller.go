package books

import (
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"players/app/auth"
	"players/app/contracts"
	"players/app/http/requests"
	"players/app/models"
	"players/app/services"
)

// BookController - Simplified version using generic CRUD controller
// From ~400 lines to ~150 lines!
type BookController struct {
	*contracts.GenericCrudController[models.Book, requests.BookCreateRequest, requests.BookUpdateRequest]
	bookService *services.BookService
}

// NewBookController creates a new simplified book controller
func NewBookController() *BookController {
	bookService := services.NewBookService()
	
	// Create the generic controller
	genericController := contracts.NewGenericCrudController[models.Book, requests.BookCreateRequest, requests.BookUpdateRequest](
		"book",
		bookService,
	)
	
	controller := &BookController{
		GenericCrudController: genericController,
		bookService:          bookService,
	}
	
	// Configure authorization using the new permission format
	genericController.SetAuthCheck(func(ctx http.Context, action string, resource interface{}) error {
		// Public actions don't need authorization
		if action == "viewAny" || action == "view" {
			return nil
		}
		
		// Map actions to permissions
		permissionMap := map[string]string{
			"create": "books_create",
			"update": "books_update", 
			"delete": "books_delete",
		}
		
		if permission, ok := permissionMap[action]; ok {
			permHelper := auth.GetPermissionHelper()
			_, err := permHelper.RequirePermission(ctx, permission)
			return err
		}
		
		return nil
	})
	
	// Configure request bindings
	genericController.SetRequestBindings(
		// Bind create request
		func(ctx http.Context) (requests.BookCreateRequest, error) {
			var req requests.BookCreateRequest
			err := ctx.Request().Bind(&req)
			return req, err
		},
		// Transform create request
		func(req requests.BookCreateRequest) map[string]interface{} {
			return req.ToCreateData()
		},
		// Bind update request
		func(ctx http.Context, id uint) (requests.BookUpdateRequest, error) {
			var req requests.BookUpdateRequest
			req.ID = id
			err := ctx.Request().Bind(&req)
			return req, err
		},
		// Transform update request
		func(req requests.BookUpdateRequest) map[string]interface{} {
			return req.ToUpdateData()
		},
	)
	
	// Register controller
	contracts.MustRegisterCrudController("books", controller)
	
	return controller
}

// Custom endpoints beyond basic CRUD

// GetByISBN GET /books/isbn/{isbn}
func (c *BookController) GetByISBN(ctx http.Context) http.Response {
	isbn := ctx.Request().Route("isbn")
	if isbn == "" {
		return c.BadRequestResponse(ctx, "ISBN is required", nil)
	}
	
	book, err := c.bookService.GetByISBN(isbn)
	if err != nil {
		return c.ResourceNotFoundResponse(ctx, "book", 0)
	}
	
	return c.SuccessResponse(ctx, book, "Book found")
}

// GetByAuthor GET /books/author/{author}
func (c *BookController) GetByAuthor(ctx http.Context) http.Response {
	author := ctx.Request().Route("author")
	if author == "" {
		return c.BadRequestResponse(ctx, "Author is required", nil)
	}
	
	req, _ := c.ValidatePaginationRequest(ctx)
	if req == nil {
		req = &contracts.ListRequest{Page: 1, PageSize: 20}
	}
	
	result, err := c.bookService.GetByAuthor(author, *req)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve books by author")
	}
	
	return c.SuccessResponse(ctx, result, "Books retrieved")
}

// GetAvailable GET /books/available
func (c *BookController) GetAvailable(ctx http.Context) http.Response {
	req, _ := c.ValidatePaginationRequest(ctx)
	if req == nil {
		req = &contracts.ListRequest{Page: 1, PageSize: 20}
	}
	
	result, err := c.bookService.GetAvailable(*req)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve available books")
	}
	
	return c.SuccessResponse(ctx, result, "Available books retrieved")
}

// Borrow POST /books/{id}/borrow
func (c *BookController) Borrow(ctx http.Context) http.Response {
	idStr := ctx.Request().Route("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid book ID", nil)
	}
	
	// Check authorization
	permHelper := auth.GetPermissionHelper()
	if _, err := permHelper.RequirePermission(ctx, "books_borrow"); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}
	
	err = c.bookService.BorrowBook(uint(id))
	if err != nil {
		if err.Error() == "book is not available for borrowing" {
			return c.BadRequestResponse(ctx, err.Error(), nil)
		}
		return c.InternalErrorResponse(ctx, err.Error())
	}
	
	return c.SuccessResponse(ctx, nil, "Book borrowed successfully")
}

// Return POST /books/{id}/return
func (c *BookController) Return(ctx http.Context) http.Response {
	idStr := ctx.Request().Route("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid book ID", nil)
	}
	
	// Check authorization
	permHelper := auth.GetPermissionHelper()
	if _, err := permHelper.RequirePermission(ctx, "books_return"); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied")
	}
	
	err = c.bookService.ReturnBook(uint(id))
	if err != nil {
		if err.Error() == "book is not currently borrowed" {
			return c.BadRequestResponse(ctx, err.Error(), nil)
		}
		return c.InternalErrorResponse(ctx, err.Error())
	}
	
	return c.SuccessResponse(ctx, nil, "Book returned successfully")
}

// Advanced GET /books/advanced - with filters
func (c *BookController) Advanced(ctx http.Context) http.Response {
	// Public endpoint - no authorization needed for viewing
	req, _ := c.ValidatePaginationRequest(ctx)
	if req == nil {
		req = &contracts.ListRequest{Page: 1, PageSize: 20}
	}
	
	// Parse filters from query parameters
	filters := make(map[string]interface{})
	
	if status := ctx.Request().Query("status"); status != "" {
		filters["status"] = status
	}
	if author := ctx.Request().Query("author"); author != "" {
		filters["author"] = author
	}
	if minPrice := ctx.Request().Query("minPrice"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filters["minPrice"] = price
		}
	}
	if maxPrice := ctx.Request().Query("maxPrice"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filters["maxPrice"] = price
		}
	}
	
	result, err := c.bookService.GetListAdvanced(*req, filters)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve books: "+err.Error())
	}
	
	return c.SuccessResponse(ctx, result, "Books retrieved with filters")
}

// Contract method implementations (required but minimal)

func (c *BookController) GetSearchableFields() []string {
	return c.bookService.GetSearchableFields()
}

func (c *BookController) GetValidationRules() map[string]interface{} {
	return c.bookService.GetValidationRules()
}

func (c *BookController) CheckPermission(ctx http.Context, permission string, resource interface{}) error {
	if c.GenericCrudController.CheckAuth != nil {
		return c.GenericCrudController.CheckAuth(ctx, permission, resource)
	}
	return nil
}

func (c *BookController) GetCurrentUser(ctx http.Context) interface{} {
	permHelper := auth.GetPermissionHelper()
	return permHelper.GetAuthenticatedUser(ctx)
}

func (c *BookController) RequireAuthentication(ctx http.Context) error {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		return fmt.Errorf("authentication required")
	}
	return nil
}

func (c *BookController) BuildPermissionsMap(ctx http.Context, resourceType string) map[string]bool {
	permHelper := auth.GetPermissionHelper()
	return permHelper.BuildPermissionsMap(ctx, resourceType)
}