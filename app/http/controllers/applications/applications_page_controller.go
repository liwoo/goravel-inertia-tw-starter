package applications

import (
	"books-database/app/auth"
	"books-database/app/contracts"
	inertiaHelper "books-database/app/http/inertia"
	"books-database/app/services"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support"
)

// ApplicationPageController handles the applications page
type ApplicationPageController struct {
	*contracts.GenericPageController
	applicationService *services.ApplicationService
}

// NewApplicationPageController creates a new applications page controller
func NewApplicationPageController() *ApplicationPageController {
	applicationService := services.NewApplicationService()

	return &ApplicationPageController{
		GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
			ResourceType:      "applications",
			PageComponent:     "Application/Index",
			Service:           applicationService,
			ServiceIdentifier: auth.ServiceApplications,
			StatsEnabled:      false,
		}),
		applicationService: applicationService,
	}
}

// ShowPublicApply renders the public application page
func (c *ApplicationPageController) ShowPublicApply(ctx http.Context) http.Response {
	return inertiaHelper.Render(ctx, "Application/ApplicationPage", map[string]interface{}{
		"version": support.Version,
	})
}
