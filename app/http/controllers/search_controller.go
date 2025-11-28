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

// Note: This controller searches SMEs, BDSPs, Events, Procurements, Users, and Configs

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

	// Search SMEs if user has permission
	if permHelper.CheckServicePermission(ctx, auth.ServiceSMEs, auth.PermissionRead) {
		smeResults := c.searchSMEs(query)
		results = append(results, smeResults...)
	}

	// Search BDSPs if user has permission
	if permHelper.CheckServicePermission(ctx, auth.ServiceBdsps, auth.PermissionRead) {
		bdspResults := c.searchBDSPs(query)
		results = append(results, bdspResults...)
	}

	// Search Events if user has permission
	if permHelper.CheckServicePermission(ctx, auth.ServiceEvents, auth.PermissionRead) {
		eventResults := c.searchEvents(query)
		results = append(results, eventResults...)
	}

	// Search Procurements if user has permission
	if permHelper.CheckServicePermission(ctx, auth.ServiceProcurementNotices, auth.PermissionRead) {
		procurementResults := c.searchProcurements(query)
		results = append(results, procurementResults...)
	}

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

	return ctx.Response().Json(http.StatusOK, SearchResponse{
		Results: results,
		Total:   len(results),
	})
}

// searchSMEs performs fuzzy search on SMEs using the SmeService
func (c *SearchController) searchSMEs(query string) []SearchResult {
	results := []SearchResult{}

	smeService := services.NewSmeService()

	// Use the service's search functionality
	paginatedResult, err := smeService.Search(query, contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	})

	if err != nil || paginatedResult == nil {
		return results
	}

	// Convert service results to search results
	for _, item := range paginatedResult.Data {
		if sme, ok := item.(models.Sme); ok {
			subtitle := sme.BusinessCategory
			if sme.Sector != "" {
				subtitle = fmt.Sprintf("%s • %s", sme.BusinessCategory, sme.Sector)
			}

			results = append(results, SearchResult{
				ID:       sme.ID,
				Title:    sme.Name,
				Subtitle: subtitle,
				Type:     "sme",
				URL:      fmt.Sprintf("/admin/smes?search=%s", query),
			})
		}
	}

	return results
}

// searchBDSPs performs fuzzy search on BDSPs using the BdspService
func (c *SearchController) searchBDSPs(query string) []SearchResult {
	results := []SearchResult{}

	bdspService := services.NewBdspService()

	// Use the service's search functionality
	paginatedResult, err := bdspService.Search(query, contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	})

	if err != nil || paginatedResult == nil {
		return results
	}

	// Convert service results to search results
	for _, item := range paginatedResult.Data {
		if bdsp, ok := item.(models.Bdsp); ok {
			subtitle := bdsp.UbdspNumber
			if bdsp.RegistrationStatus != nil && *bdsp.RegistrationStatus != "" {
				subtitle = fmt.Sprintf("%s • %s", bdsp.UbdspNumber, *bdsp.RegistrationStatus)
			}

			results = append(results, SearchResult{
				ID:       bdsp.ID,
				Title:    bdsp.Name,
				Subtitle: subtitle,
				Type:     "bdsp",
				URL:      fmt.Sprintf("/admin/bdsps?search=%s", query),
			})
		}
	}

	return results
}

// searchEvents performs fuzzy search on Events using the EventService
func (c *SearchController) searchEvents(query string) []SearchResult {
	results := []SearchResult{}

	eventService := services.NewEventService()

	// Use the service's search functionality
	paginatedResult, err := eventService.Search(query, contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	})

	if err != nil || paginatedResult == nil {
		return results
	}

	// Convert service results to search results
	for _, item := range paginatedResult.Data {
		if event, ok := item.(models.Event); ok {
			subtitle := event.Venue
			if event.District != "" {
				subtitle = fmt.Sprintf("%s • %s", event.Venue, event.District)
			}

			results = append(results, SearchResult{
				ID:       event.ID,
				Title:    event.Title,
				Subtitle: subtitle,
				Type:     "event",
				URL:      fmt.Sprintf("/admin/events?search=%s", query),
			})
		}
	}

	return results
}

// searchProcurements performs fuzzy search on Procurements using the ProcurementNoticeService
func (c *SearchController) searchProcurements(query string) []SearchResult {
	results := []SearchResult{}

	procurementService := services.NewProcurementNoticeService()

	// Use the service's search functionality
	paginatedResult, err := procurementService.Search(query, contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	})

	if err != nil || paginatedResult == nil {
		return results
	}

	// Convert service results to search results
	for _, item := range paginatedResult.Data {
		if procurement, ok := item.(models.ProcurementNotice); ok {
			subtitle := procurement.ProcurementType
			if procurement.Organization != "" {
				subtitle = fmt.Sprintf("%s • %s", procurement.ProcurementType, procurement.Organization)
			}

			results = append(results, SearchResult{
				ID:       procurement.ID,
				Title:    procurement.Invitation,
				Subtitle: subtitle,
				Type:     "procurement",
				URL:      fmt.Sprintf("/admin/procurements?search=%s", query),
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
				URL:      fmt.Sprintf("/admin/configs?search=%s", query),
			})
		}
	}

	return results
}
