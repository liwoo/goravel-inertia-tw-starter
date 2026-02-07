package controllers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
)

// Note: This controller searches Users, Configs, and Applications

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

	// Normalize query - trim whitespace (ILIKE handles case-insensitivity)
	query = strings.TrimSpace(query)

	// Get current user and permissions
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return ctx.Response().Status(403).Json(map[string]interface{}{
			"error": "Unauthorized",
		})
	}

	results := []SearchResult{}

	// Search Users if user has permission
	if permHelper.CheckServicePermission(ctx, auth.ServiceUsers, auth.PermissionRead) {
		userResults := c.searchUsers(query)
		results = append(results, userResults...)
	}

	// Search Configs if user has permission
	if permHelper.CheckServicePermission(ctx, auth.ServiceConfig, auth.PermissionRead) {
		configResults := c.searchConfigs(query)
		results = append(results, configResults...)
	}

	// Search Applications if user has permission
	if permHelper.CheckServicePermission(ctx, auth.ServiceApplications, auth.PermissionRead) {
		applicationResults := c.searchApplications(query)
		results = append(results, applicationResults...)
	}

	return ctx.Response().Json(http.StatusOK, SearchResponse{
		Results: results,
		Total:   len(results),
	})
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
				URL:      fmt.Sprintf("/admin/users?search=%s", url.QueryEscape(user.Name)),
			})
		}
	}

	return results
}

// searchConfigs performs fuzzy search on configs using the ConfigService
func (c *SearchController) searchConfigs(query string) []SearchResult {
	results := []SearchResult{}

	configService := services.NewConfigService()

	// Use the service's search functionality
	paginatedResult, err := configService.Search(query, contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	})

	if err != nil || paginatedResult == nil {
		return results
	}

	// Convert service results to search results
	// Note: Data contains values, not pointers
	for _, item := range paginatedResult.Data {
		if config, ok := item.(models.Config); ok {
			subtitle := config.ConfigType
			if config.Description != nil && *config.Description != "" {
				subtitle = fmt.Sprintf("%s • %s", config.ConfigType, *config.Description)
			}

			results = append(results, SearchResult{
				ID:       config.ID,
				Title:    config.Name,
				Subtitle: subtitle,
				Type:     "config",
				URL:      fmt.Sprintf("/admin/configs?search=%s", url.QueryEscape(config.Name)),
			})
		}
	}

	return results
}

// searchApplications performs fuzzy search on applications using the ApplicationService
func (c *SearchController) searchApplications(query string) []SearchResult {
	results := []SearchResult{}

	applicationService := services.NewApplicationService()

	// Use the service's search functionality
	paginatedResult, err := applicationService.Search(query, contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	})

	if err != nil || paginatedResult == nil {
		return results
	}

	// Convert service results to search results
	for _, item := range paginatedResult.Data {
		if application, ok := item.(models.Application); ok {
			subtitle := application.Status
			if application.SME != "" {
				subtitle = fmt.Sprintf("%s • %s", application.SME, application.Status)
			}

			results = append(results, SearchResult{
				ID:       application.ID,
				Title:    application.RegistrantName,
				Subtitle: subtitle,
				Type:     "application",
				URL:      fmt.Sprintf("/admin/applications?search=%s", url.QueryEscape(application.RegistrantName)),
			})
		}
	}

	return results
}
