package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"players/app/http/inertia"
)

type DashboardController struct {
	// Dependencies can be injected here
}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

// Show displays the dashboard page.
func (r *DashboardController) Show(ctx http.Context) http.Response {
	// You can pass any specific props needed for the dashboard here
	return inertia.Render(ctx, "dashboard/Index", map[string]interface{}{
		"pageTitle": "Dashboard",
	})
}
