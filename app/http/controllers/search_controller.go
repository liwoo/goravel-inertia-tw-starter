package controllers

import (
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
)

type SearchController struct{}

func NewSearchController() *SearchController {
	return &SearchController{}
}

type SearchResult struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle,omitempty"`
	Type     string `json:"type"`
	URL      string `json:"url"`
}

type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
}

// GlobalSearch performs a fuzzy search across all accessible resources
func (c *SearchController) GlobalSearch(ctx http.Context) http.Response {
	// Get search query
	query := ctx.Request().Query("q", "")
	if query == "" {
		return ctx.Response().Json(http.StatusOK, SearchResponse{
			Results: []SearchResult{},
			Total:   0,
		})
	}

	// Normalize query for case-insensitive search
	query = strings.ToLower(strings.TrimSpace(query))

	// Get current user and permissions
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return ctx.Response().Status(403).Json(map[string]interface{}{
			"error": "Unauthorized",
		})
	}

	results := []SearchResult{}

	// Search Books if user has permission
	if permHelper.CheckServicePermission(ctx, auth.ServiceBooks, auth.PermissionRead) {
		bookResults := c.searchBooks(query)
		results = append(results, bookResults...)
	}

	// Search Users if user has permission
	if permHelper.CheckServicePermission(ctx, auth.ServiceUsers, auth.PermissionRead) {
		userResults := c.searchUsers(query)
		results = append(results, userResults...)
	}

	// Search Lenders if user has permission
	if permHelper.CheckServicePermission(ctx, auth.ServiceLenders, auth.PermissionRead) {
		lenderResults := c.searchLenders(query)
		results = append(results, lenderResults...)
	}

	return ctx.Response().Json(http.StatusOK, SearchResponse{
		Results: results,
		Total:   len(results),
	})
}

// searchBooks performs fuzzy search on books using the BookService
func (c *SearchController) searchBooks(query string) []SearchResult {
	results := []SearchResult{}

	bookService := services.NewBookService()

	// Use the service's search functionality
	paginatedResult, err := bookService.Search(query, contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	})

	if err != nil || paginatedResult == nil {
		return results
	}

	// Convert service results to search results
	// Note: Data contains values, not pointers
	for _, item := range paginatedResult.Data {
		if book, ok := item.(models.Book); ok {
			subtitle := book.Author
			if book.Status != "" {
				subtitle = fmt.Sprintf("%s • %s", book.Author, book.Status)
			}

			results = append(results, SearchResult{
				ID:       book.ID,
				Title:    book.Title,
				Subtitle: subtitle,
				Type:     "book",
				URL:      fmt.Sprintf("/admin/books?search=%s", query),
			})
		}
	}

	return results
}

// searchUsers performs fuzzy search on users using the UserService
func (c *SearchController) searchUsers(query string) []SearchResult {
	results := []SearchResult{}

	userService := services.NewUserService()

	// Use the service's search functionality
	paginatedResult, err := userService.Search(query, contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	})

	if err != nil || paginatedResult == nil {
		return results
	}

	// Convert service results to search results
	// Note: Data contains values, not pointers
	for _, item := range paginatedResult.Data {
		if user, ok := item.(models.User); ok {
			results = append(results, SearchResult{
				ID:       user.ID,
				Title:    user.Name,
				Subtitle: user.Email,
				Type:     "user",
				URL:      fmt.Sprintf("/admin/users?search=%s", query),
			})
		}
	}

	return results
}

// searchLenders performs fuzzy search on lenders using the LenderService
func (c *SearchController) searchLenders(query string) []SearchResult {
	results := []SearchResult{}

	lenderService := services.NewLenderService()

	// Use the service's search functionality
	paginatedResult, err := lenderService.Search(query, contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	})

	if err != nil || paginatedResult == nil {
		return results
	}

	// Convert service results to search results
	// Note: Data contains values, not pointers
	for _, item := range paginatedResult.Data {
		if lender, ok := item.(models.Lender); ok {
			subtitle := lender.Email
			if lender.Phone != nil && *lender.Phone != "" {
				subtitle = fmt.Sprintf("%s • %s", lender.Email, *lender.Phone)
			}

			results = append(results, SearchResult{
				ID:       lender.ID,
				Title:    lender.Name,
				Subtitle: subtitle,
				Type:     "lender",
				URL:      fmt.Sprintf("/admin/lenders?search=%s", query),
			})
		}
	}

	return results
}
