package books

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"players/app/auth"
	"players/app/contracts"
	"players/app/helpers"
	"players/app/http/inertia"
	"players/app/services"
)

// BooksPageController handles the Inertia.js Books management page
// Implements PageControllerContract interface
type BooksPageController struct {
	*contracts.BasePageController
	bookService *services.BookService
	authHelper  contracts.AuthHelper
}

// GetServiceIdentifier returns the service identifier for this controller
func (c *BooksPageController) GetServiceIdentifier() auth.ServiceRegistry {
	return auth.ServiceBooks
}

// NewBooksPageController creates a new books page controller that implements all contracts
func NewBooksPageController() *BooksPageController {
	controller := &BooksPageController{
		BasePageController: contracts.NewBasePageController("book", "Books/Index"),
		bookService:        services.NewBookService(),
		authHelper:         helpers.NewAuthHelper(),
	}

	// Register page controller with validation
	contracts.MustRegisterPageController("books_page", controller)

	return controller
}

// Index renders the Books management page with data and permissions - Implements PageControllerContract
func (c *BooksPageController) Index(ctx http.Context) http.Response {
	// Check if user has at least view permission for books (for listing)
	permHelper := auth.GetPermissionHelper()
	_, err := permHelper.RequireServicePermission(ctx, auth.ServiceBooks, auth.PermissionView)
	if err != nil {
		// Return 403 Forbidden
		return ctx.Response().Status(403).Json(map[string]interface{}{
			"error":   "Forbidden",
			"message": "You don't have permission to access this page",
		})
	}

	// Validate page request using contract
	req, err := c.ValidatePageRequest(ctx)
	if err != nil {
		// Return error page or redirect with error
		req = &contracts.ListRequest{Page: 1, PageSize: 20}
		req.SetDefaults()
	}

	// Build permissions map using contract - using service identifier
	permissions := c.BuildPermissionsMap(ctx, string(c.GetServiceIdentifier()))

	// Debug logging for permissions
	fmt.Printf("DEBUG: User authentication status: %+v\n", c.GetCurrentUser(ctx) != nil)
	fmt.Printf("DEBUG: Permissions for 'books' resource: %+v\n", permissions)

	// Get books data
	booksResult, err := c.bookService.GetList(*req)
	if err != nil {
		// Handle error gracefully, provide empty result
		booksResult = &contracts.PaginatedResult{
			Data:        []interface{}{},
			Total:       0,
			CurrentPage: 1,
			LastPage:    1,
			PerPage:     req.PageSize,
		}
	}

	// Get book statistics if user can view reports
	var stats map[string]interface{}
	if permissions["canViewReports"] {
		stats = c.getBookStatistics()
	}

	// Build standardized page props using contract
	data := map[string]interface{}{
		"data":        booksResult.Data,
		"total":       booksResult.Total,
		"currentPage": booksResult.CurrentPage,
		"lastPage":    booksResult.LastPage,
		"perPage":     booksResult.PerPage,
		"from":        booksResult.From,
		"to":          booksResult.To,
		"hasNext":     booksResult.HasNext,
		"hasPrev":     booksResult.HasPrev,
	}

	filters := map[string]interface{}{
		"page":      req.Page,
		"pageSize":  req.PageSize,
		"search":    req.Search,
		"sort":      req.Sort,
		"direction": req.Direction,
		"filters":   req.Filters,
	}

	meta := map[string]interface{}{
		"stats": stats,
	}

	props := c.BuildPageProps(data, filters, permissions, meta)

	// Ensure all required props are present and not nil
	if props["data"] == nil {
		props["data"] = map[string]interface{}{
			"data":        []interface{}{},
			"total":       0,
			"currentPage": 1,
			"lastPage":    1,
			"perPage":     20,
		}
	}
	if props["filters"] == nil {
		props["filters"] = map[string]interface{}{}
	}
	if props["permissions"] == nil {
		props["permissions"] = map[string]interface{}{}
	}

	return inertia.Render(ctx, "Books/Index", props)
}

// getBookStatistics returns book statistics for the dashboard
func (c *BooksPageController) getBookStatistics() map[string]interface{} {
	// Get status counts
	availableCount := c.getBookCountByStatus("AVAILABLE")
	borrowedCount := c.getBookCountByStatus("BORROWED")
	maintenanceCount := c.getBookCountByStatus("MAINTENANCE")
	totalCount := availableCount + borrowedCount + maintenanceCount

	// Calculate total value (this would need to be implemented in the service)
	totalValue := 0.0 // Placeholder

	// Get top authors (this would need to be implemented in the service)
	topAuthors := []map[string]interface{}{
		{"name": "Sample Author 1", "count": 5},
		{"name": "Sample Author 2", "count": 3},
	}

	return map[string]interface{}{
		"totalBooks":       totalCount,
		"availableBooks":   availableCount,
		"borrowedBooks":    borrowedCount,
		"maintenanceBooks": maintenanceCount,
		"totalValue":       totalValue,
		"averagePrice":     totalValue / float64(max(totalCount, 1)),
		"topAuthors":       topAuthors,
	}
}

// getBookCountByStatus is a helper function to get book count by status
func (c *BooksPageController) getBookCountByStatus(status string) int {
	// This is a simplified version. In a real implementation,
	// you'd want to add a method to the service to get counts efficiently
	req := contracts.ListRequest{
		PageSize: 1,
		Filters: map[string]interface{}{
			"status": status,
		},
	}

	result, err := c.bookService.GetListAdvanced(req, map[string]interface{}{
		"status": status,
	})
	if err != nil {
		return 0
	}

	return int(result.Total)
}

// max helper function
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// CONTRACT IMPLEMENTATIONS - Required by PageControllerContract interface

// AuthorizationControllerContract implementation
func (c *BooksPageController) CheckPermission(ctx http.Context, permission string, resource interface{}) error {
	permHelper := auth.GetPermissionHelper()
	_, err := permHelper.RequirePermission(ctx, permission)
	return err
}

func (c *BooksPageController) GetCurrentUser(ctx http.Context) interface{} {
	permHelper := auth.GetPermissionHelper()
	return permHelper.GetAuthenticatedUser(ctx)
}

func (c *BooksPageController) RequireAuthentication(ctx http.Context) error {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		return fmt.Errorf("authentication required")
	}
	return nil
}

func (c *BooksPageController) BuildPermissionsMap(ctx http.Context, resourceType string) map[string]bool {
	permHelper := auth.GetPermissionHelper()
	return permHelper.BuildPermissionsMap(ctx, resourceType)
}
