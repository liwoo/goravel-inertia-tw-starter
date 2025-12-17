package myapplications

import (
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/http/inertia"
	"smedi-sme-db/app/services"

	"github.com/goravel/framework/contracts/http"
)

// MyApplicationsPageController handles the my applications page for SME users
type MyApplicationsPageController struct {
	*contracts.GenericPageController
	smeService         *services.SmeService
	applicationService *services.ApplicationService
}

// NewMyApplicationsPageController creates a new my applications page controller
func NewMyApplicationsPageController() *MyApplicationsPageController {
	smeService := services.NewSmeService()
	applicationService := services.NewApplicationService()

	controller := &MyApplicationsPageController{
		GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
			ResourceType:  "my-applications",
			PageComponent: "MyApplications/Index",
			Service:       applicationService,
			// No service identifier - we use custom permission check for authenticated users only
			StatsEnabled: false,
		}),
		smeService:         smeService,
		applicationService: applicationService,
	}

	// Set mandatory filter provider to filter applications by the current user's SME
	controller.SetMandatoryFilterProvider(func(ctx http.Context) map[string]interface{} {
		user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
		if user == nil {
			return map[string]interface{}{"sme_id": 0} // No results
		}

		// Get user's linked SME
		sme, err := smeService.GetSmeByUserEmail(user.Email)
		if err != nil || sme == nil {
			return map[string]interface{}{"sme_id": 0} // No results
		}

		return map[string]interface{}{"sme_id": sme.ID}
	})

	// Add extra data providers for SME info
	controller.AddExtraDataProvider("smeName", func(ctx http.Context) (interface{}, error) {
		user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
		if user == nil {
			return "", nil
		}
		sme, err := smeService.GetSmeByUserEmail(user.Email)
		if err != nil || sme == nil {
			return "", nil
		}
		return sme.Name, nil
	})

	controller.AddExtraDataProvider("usmeNumber", func(ctx http.Context) (interface{}, error) {
		user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
		if user == nil {
			return "", nil
		}
		sme, err := smeService.GetSmeByUserEmail(user.Email)
		if err != nil || sme == nil {
			return "", nil
		}
		return sme.UsmeNumber, nil
	})

	controller.AddExtraDataProvider("classification", func(ctx http.Context) (interface{}, error) {
		user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
		if user == nil {
			return "", nil
		}
		sme, err := smeService.GetSmeByUserEmail(user.Email)
		if err != nil || sme == nil {
			return "", nil
		}
		return sme.Classification, nil
	})

	controller.AddExtraDataProvider("userName", func(ctx http.Context) (interface{}, error) {
		user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
		if user == nil {
			return "", nil
		}
		return user.Name, nil
	})

	controller.AddExtraDataProvider("error", func(ctx http.Context) (interface{}, error) {
		user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
		if user == nil {
			return "Authentication required", nil
		}
		sme, err := smeService.GetSmeByUserEmail(user.Email)
		if err != nil || sme == nil {
			return "No SME linked to your account", nil
		}
		return "", nil
	})

	return controller
}

// Index overrides the default Index to handle the case where user has no linked SME
func (c *MyApplicationsPageController) Index(ctx http.Context) http.Response {
	// Check authentication first
	user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
	if user == nil {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	// Check if user has linked SME - if not, show empty state with error
	sme, err := c.smeService.GetSmeByUserEmail(user.Email)
	if err != nil || sme == nil {
		// Return page with error message and empty data
		return inertia.Render(ctx, "MyApplications/Index", map[string]interface{}{
			"data": map[string]interface{}{
				"data":        []interface{}{},
				"total":       0,
				"currentPage": 1,
				"lastPage":    1,
				"perPage":     10,
				"from":        0,
				"to":          0,
				"hasNext":     false,
				"hasPrev":     false,
			},
			"filters": map[string]interface{}{
				"page":     1,
				"pageSize": 10,
			},
			"permissions": map[string]interface{}{
				"canCreate": false,
				"canEdit":   false,
				"canDelete": false,
			},
			"userName": user.Name,
			"error":    "No SME linked to your account",
		})
	}

	// Call the parent Index method which handles everything else
	return c.GenericPageController.Index(ctx)
}
