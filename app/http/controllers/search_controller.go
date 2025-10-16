package controllers

import (
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"players/app/auth"
	"players/app/models"
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

// searchBooks performs fuzzy search on books
func (c *SearchController) searchBooks(query string) []SearchResult {
	var books []models.Book
	results := []SearchResult{}

	// Build query for case-insensitive search
	searchPattern := "%" + query + "%"

	// Search in title, author, isbn, and description
	// Using COLLATE NOCASE for SQLite compatibility
	facades.Orm().Query().
		Where("title COLLATE NOCASE LIKE ?", searchPattern).
		OrWhere("author COLLATE NOCASE LIKE ?", searchPattern).
		OrWhere("isbn COLLATE NOCASE LIKE ?", searchPattern).
		OrWhere("description COLLATE NOCASE LIKE ?", searchPattern).
		Order("title ASC").
		Limit(10).
		Find(&books)

	for _, book := range books {
		// Highlight matching parts in the subtitle
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

	return results
}

// searchUsers performs fuzzy search on users
func (c *SearchController) searchUsers(query string) []SearchResult {
	var users []models.User
	results := []SearchResult{}

	searchPattern := "%" + query + "%"

	// Search in name and email
	facades.Orm().Query().
		Where("name COLLATE NOCASE LIKE ?", searchPattern).
		OrWhere("email COLLATE NOCASE LIKE ?", searchPattern).
		Order("name ASC").
		Limit(10).
		Find(&users)

	for _, user := range users {
		results = append(results, SearchResult{
			ID:       user.ID,
			Title:    user.Name,
			Subtitle: user.Email,
			Type:     "user",
			URL:      fmt.Sprintf("/admin/users?search=%s", query),
		})
	}

	return results
}

// searchLenders performs fuzzy search on lenders
func (c *SearchController) searchLenders(query string) []SearchResult {
	var lenders []models.Lender
	results := []SearchResult{}

	searchPattern := "%" + query + "%"

	// Search in name, email, and phone
	facades.Orm().Query().
		Where("name COLLATE NOCASE LIKE ?", searchPattern).
		OrWhere("email COLLATE NOCASE LIKE ?", searchPattern).
		OrWhere("phone COLLATE NOCASE LIKE ?", searchPattern).
		Order("name ASC").
		Limit(10).
		Find(&lenders)

	for _, lender := range lenders {
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

	return results
}
