package controllers

import (
	"starter-project/app/http/inertia"
	"starter-project/app/services"

	"github.com/goravel/framework/contracts/http"
)

type DashboardController struct {
	dashboardService *services.DashboardService
}

func NewDashboardController() *DashboardController {
	return &DashboardController{
		dashboardService: services.NewDashboardService(),
	}
}

// Show displays the dashboard page with statistics and widgets.
func (r *DashboardController) Show(ctx http.Context) http.Response {

	// Get recent user activities (last 10)
	recentActivities := r.dashboardService.GetRecentActivities(10)

	return inertia.Render(ctx, "dashboard/Index", map[string]interface{}{
		"pageTitle":        "Dashboard",
		"recentActivities": recentActivities,
	})
}
