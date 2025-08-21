package books

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
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

// GetFilters returns filter metadata for books
func (c *BookController) GetFilters(ctx http.Context) http.Response {
	metadata := map[string]interface{}{
		"filters": []map[string]interface{}{
			{
				"field":     "title",
				"label":     "Title",
				"type":      "string",
				"operators": []string{"contains", "not_contains", "starts_with", "ends_with", "equals", "not_equals"},
			},
			{
				"field":     "author",
				"label":     "Author",
				"type":      "string",
				"operators": []string{"contains", "not_contains", "equals", "not_equals"},
			},
			{
				"field":       "status",
				"label":       "Status",
				"type":        "enum",
				"operators":   []string{"equals", "not_equals", "in", "not_in"},
				"enum_values": []string{"AVAILABLE", "BORROWED", "MAINTENANCE", "RESERVED"},
			},
			{
				"field":     "price",
				"label":     "Price",
				"type":      "number",
				"operators": []string{"equals", "not_equals", "greater_than", "less_than", "greater_than_or_equal", "less_than_or_equal", "between", "not_between"},
			},
			{
				"field":     "tags",
				"label":     "Tags",
				"type":      "array",
				"operators": []string{"contains", "not_contains", "is_empty", "is_not_empty"},
			},
			{
				"field":     "published_at",
				"label":     "Published Date",
				"type":      "date",
				"operators": []string{"before", "after", "between", "not_between", "is_today", "is_yesterday", "is_this_week", "is_this_month", "is_this_year", "last_n_days"},
			},
			{
				"field":     "created_at",
				"label":     "Date Added",
				"type":      "datetime",
				"operators": []string{"before", "after", "between", "not_between", "is_today", "is_this_week", "is_this_month", "last_n_days"},
			},
		},
		"logic_operators":   []string{"AND", "OR"},
		"resource":          "book",
		"filterable_fields": []string{"status", "author"},
		"searchable_fields": []string{"title", "author", "isbn", "description"},
		"sortable_fields":   []string{"id", "title", "author", "price", "created_at", "updated_at", "published_at"},
	}

	return ctx.Response().Json(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    metadata,
	})
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

// GetFilterDefinitions returns custom filter definitions for books
func (c *BookController) GetFilterDefinitions() []contracts.FilterDefinition {
	return []contracts.FilterDefinition{
		// Price filter
		contracts.NewFilterDefinition(
			"price",
			"Price",
			contracts.FilterTypeNumber,
			[]contracts.FilterOperator{
				contracts.OperatorEquals,
				contracts.OperatorGreaterThan,
				contracts.OperatorLessThan,
				contracts.OperatorGreaterThanOrEqual,
				contracts.OperatorLessThanOrEqual,
				contracts.OperatorBetween,
			},
		),
		// Status filter
		{
			Field: "status",
			Label: "Status",
			Type:  contracts.FilterTypeEnum,
			Operators: []contracts.FilterOperator{
				contracts.OperatorEquals,
				contracts.OperatorNotEquals,
				contracts.OperatorIn,
				contracts.OperatorNotIn,
			},
			EnumValues: []string{"AVAILABLE", "BORROWED", "MAINTENANCE", "RESERVED", "LOST"},
		},
		// Published date filter
		contracts.NewFilterDefinition(
			"published_at",
			"Published Date",
			contracts.FilterTypeDate,
			[]contracts.FilterOperator{
				contracts.OperatorBefore,
				contracts.OperatorAfter,
				contracts.OperatorBetween,
				contracts.OperatorIsToday,
				contracts.OperatorIsThisMonth,
				contracts.OperatorIsThisYear,
				contracts.OperatorLastNDays,
			},
		),
		// Author filter
		contracts.NewFilterDefinition(
			"author",
			"Author",
			contracts.FilterTypeString,
			[]contracts.FilterOperator{
				contracts.OperatorEquals,
				contracts.OperatorContains,
				contracts.OperatorStartsWith,
				contracts.OperatorEndsWith,
			},
		),
		// Is Verified filter
		contracts.NewFilterDefinition(
			"is_verified",
			"Verified",
			contracts.FilterTypeBoolean,
			[]contracts.FilterOperator{
				contracts.OperatorIsTrue,
				contracts.OperatorIsFalse,
			},
		),
	}
}
