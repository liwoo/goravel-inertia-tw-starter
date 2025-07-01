package books

import (
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"players/app/auth"
	"players/app/contracts"
	"players/app/helpers"
	"players/app/http/requests"
	"players/app/services"
)

// BookController - enhanced with validation and authorization
// Implements ResourceControllerContract interface
type BookController struct {
	*contracts.BaseCrudController
	bookService *services.BookService
	authHelper  contracts.AuthHelper
}

// NewBookController creates a new book controller that implements all contracts
func NewBookController() *BookController {
	controller := &BookController{
		BaseCrudController: contracts.NewBaseCrudController("book"),
		bookService:        services.NewBookService(),
		authHelper:         helpers.NewAuthHelper(),
	}

	// Register controller with validation
	contracts.MustRegisterCrudController("books", controller)

	return controller
}

// Index GET /books - Implements CrudControllerContract
func (c *BookController) Index(ctx http.Context) http.Response {
	// Validate pagination request using contract
	req, err := c.ValidatePaginationRequest(ctx)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid pagination parameters", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Get books using service
	result, err := c.bookService.GetList(*req)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve books: "+err.Error())
	}

	// Build standardized paginated response
	response := c.BuildPaginatedResponse(result, req)
	return c.SuccessResponse(ctx, response, "Books retrieved successfully")
}

// Show GET /books/{id} - Implements CrudControllerContract (JSON for modals)
func (c *BookController) Show(ctx http.Context) http.Response {
	// Validate ID parameter using contract
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid book ID", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Public endpoint - no authorization needed for viewing individual books

	// Get the book
	book, err := c.bookService.GetByID(id)
	if err != nil {
		return c.ResourceNotFoundResponse(ctx, "book", id)
	}

	return c.SuccessResponse(ctx, book, "Book details retrieved successfully")
}

// Store POST /books - Implements CrudControllerContract
func (c *BookController) Store(ctx http.Context) http.Response {
	// Check authorization using new permission format
	if err := c.CheckPermission(ctx, "books_create", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
	}

	// Validate create request using contract
	data, err := c.ValidateCreateRequest(ctx)
	if err != nil {
		return c.ValidationErrorResponse(ctx, map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Create the book using validated data
	book, err := c.bookService.Create(data)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to create book: "+err.Error())
	}

	return c.ResourceCreatedResponse(ctx, book, "book")
}

// Update PUT /books/{id} - Implements CrudControllerContract
func (c *BookController) Update(ctx http.Context) http.Response {
	// Validate ID parameter using contract
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid book ID", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Check if book exists
	_, err = c.bookService.GetByID(id)
	if err != nil {
		return c.ResourceNotFoundResponse(ctx, "book", id)
	}

	// Check authorization using new permission format
	if err := c.CheckPermission(ctx, "books_update", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
	}

	// Validate update request using contract
	data, err := c.ValidateUpdateRequest(ctx, id)
	if err != nil {
		return c.ValidationErrorResponse(ctx, map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Update the book using validated data
	updatedBook, err := c.bookService.Update(id, data)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to update book: "+err.Error())
	}

	return c.ResourceUpdatedResponse(ctx, updatedBook, "book")
}

// Delete DELETE /books/{id} - Implements CrudControllerContract
func (c *BookController) Delete(ctx http.Context) http.Response {
	// Validate ID parameter using contract
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid book ID", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Check if book exists
	_, err = c.bookService.GetByID(id)
	if err != nil {
		return c.ResourceNotFoundResponse(ctx, "book", id)
	}

	// Check authorization using new permission format
	if err := c.CheckPermission(ctx, "books_delete", nil); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: "+err.Error())
	}

	// Delete the book
	err = c.bookService.Delete(id)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to delete book: "+err.Error())
	}

	return c.ResourceDeletedResponse(ctx, "book", id)
}

// GetByISBN GET /books/isbn/{isbn}
func (c *BookController) GetByISBN(ctx http.Context) http.Response {
	// Public endpoint - no authorization needed for viewing
	isbn := ctx.Request().Route("isbn")
	if isbn == "" {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "ISBN is required",
		})
	}

	book, err := c.bookService.GetByISBN(isbn)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Book not found with ISBN: " + isbn,
		})
	}

	return ctx.Response().Json(http.StatusOK, book)
}

// GetByAuthor GET /books/author/{author}
func (c *BookController) GetByAuthor(ctx http.Context) http.Response {
	// Public endpoint - no authorization needed
	author := ctx.Request().Route("author")
	if author == "" {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Author is required",
		})
	}

	var req helpers.ListRequest
	if err := ctx.Request().Bind(&req); err != nil {
		req = helpers.ListRequest{} // Use defaults
	}

	result, err := c.bookService.GetByAuthor(author, req)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve books by author",
		})
	}

	return ctx.Response().Json(http.StatusOK, result)
}

// GetAvailable GET /books/available
func (c *BookController) GetAvailable(ctx http.Context) http.Response {
	// Public endpoint - no authorization needed
	var req helpers.ListRequest
	if err := ctx.Request().Bind(&req); err != nil {
		req = helpers.ListRequest{} // Use defaults
	}

	result, err := c.bookService.GetAvailable(req)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve available books",
		})
	}

	return ctx.Response().Json(http.StatusOK, result)
}

// Advanced GET /books/advanced - with filters
func (c *BookController) Advanced(ctx http.Context) http.Response {
	// Public endpoint - no authorization needed for viewing
	var req helpers.ListRequest
	if err := ctx.Request().Bind(&req); err != nil {
		req = helpers.ListRequest{} // Use defaults
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

	result, err := c.bookService.GetListAdvanced(req, filters)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, result)
}

// Borrow POST /books/{id}/borrow
func (c *BookController) Borrow(ctx http.Context) http.Response {
	idStr := ctx.Request().Route("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Invalid book ID",
		})
	}

	// Check authorization for borrowing
	// TODO: Re-implement gate check
	if false && false { // Disabled: response := facades.Gate().Inspect("borrow.books", ctx); response.Denied() {
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": "Access denied",
		})
	}

	err = c.bookService.BorrowBook(uint(id))
	if err != nil {
		if err.Error() == "book is not available for borrowing" {
			return ctx.Response().Json(http.StatusConflict, map[string]string{
				"error": err.Error(),
			})
		}
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, map[string]string{
		"message": "Book borrowed successfully",
	})
}

// Return POST /books/{id}/return
func (c *BookController) Return(ctx http.Context) http.Response {
	idStr := ctx.Request().Route("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Invalid book ID",
		})
	}

	// Check authorization for returning
	// TODO: Re-implement gate check
	if false && false { // Disabled: response := facades.Gate().Inspect("return.books", ctx); response.Denied() {
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": "Access denied",
		})
	}

	err = c.bookService.ReturnBook(uint(id))
	if err != nil {
		if err.Error() == "book is not currently borrowed" {
			return ctx.Response().Json(http.StatusConflict, map[string]string{
				"error": err.Error(),
			})
		}
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, map[string]string{
		"message": "Book returned successfully",
	})
}

// CONTRACT IMPLEMENTATIONS - Required by ResourceControllerContract interface

// ValidationControllerContract implementation
func (c *BookController) ValidateCreateRequest(ctx http.Context) (map[string]interface{}, error) {
	var createRequest requests.BookCreateRequest

	// Bind the data to the struct
	if err := ctx.Request().Bind(&createRequest); err != nil {
		fmt.Printf("DEBUG: Bind error: %v\n", err)
		return nil, fmt.Errorf("data binding failed: %w", err)
	}

	fmt.Printf("DEBUG: Successfully bound data: %+v\n", createRequest)

	// Manual validation for debugging - check field lengths
	if len(createRequest.Title) > 255 {
		return nil, fmt.Errorf("validation errors: title exceeds 255 characters (%d)", len(createRequest.Title))
	}
	if len(createRequest.Author) > 100 {
		return nil, fmt.Errorf("validation errors: author exceeds 100 characters (%d)", len(createRequest.Author))
	}
	if len(createRequest.Description) > 1000 {
		return nil, fmt.Errorf("validation errors: description exceeds 1000 characters (%d)", len(createRequest.Description))
	}

	// Check required fields
	if createRequest.Title == "" {
		return nil, fmt.Errorf("validation errors: title is required")
	}
	if createRequest.Author == "" {
		return nil, fmt.Errorf("validation errors: author is required")
	}
	if createRequest.ISBN == "" {
		return nil, fmt.Errorf("validation errors: isbn is required")
	}

	fmt.Printf("DEBUG: Manual validation passed\n")
	return createRequest.ToCreateData(), nil
}

func (c *BookController) ValidateUpdateRequest(ctx http.Context, id uint) (map[string]interface{}, error) {
	var updateRequest requests.BookUpdateRequest
	updateRequest.ID = id // Set the ID for validation context

	errors, err := ctx.Request().ValidateRequest(&updateRequest)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	if errors != nil {
		return nil, fmt.Errorf("validation errors: %v", errors.All())
	}

	return updateRequest.ToUpdateData(), nil
}

func (c *BookController) GetValidationRules() map[string]interface{} {
	return c.bookService.GetValidationRules()
}

// AuthorizationControllerContract implementation
func (c *BookController) CheckPermission(ctx http.Context, permission string, resource interface{}) error {
	permHelper := auth.GetPermissionHelper()
	_, err := permHelper.RequirePermission(ctx, permission)
	return err
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
