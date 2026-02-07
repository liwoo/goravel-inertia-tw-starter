package controllers

import (
	"books-database/app/http/inertia"
	"books-database/app/services"
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
	// Get dashboard statistics
	stats := r.dashboardService.GetDashboardStats()

	return inertia.Render(ctx, "dashboard/Index", map[string]interface{}{
		"pageTitle": "Dashboard",
		"stats":     stats,
	})
}
