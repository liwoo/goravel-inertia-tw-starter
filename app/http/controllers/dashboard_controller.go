package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"smedi-sme-db/app/http/inertia"
	"smedi-sme-db/app/services"
)

type DashboardController struct {
	dashboardService *services.DashboardService
	smeService       *services.SmeService
}

func NewDashboardController() *DashboardController {
	return &DashboardController{
		dashboardService: services.NewDashboardService(),
		smeService:       services.NewSmeService(),
	}
}

// Show displays the dashboard page with statistics and widgets.
func (r *DashboardController) Show(ctx http.Context) http.Response {
	// Get dashboard statistics
	stats := r.dashboardService.GetDashboardStats(r.smeService)

	// Get upcoming events (next 5)
	upcomingEvents := r.dashboardService.GetUpcomingEvents(5)

	// Get upcoming procurements (next 5)
	upcomingProcurements := r.dashboardService.GetUpcomingProcurements(5)

	// Get recent user activities (last 10)
	recentActivities := r.dashboardService.GetRecentActivities(10)

	return inertia.Render(ctx, "dashboard/Index", map[string]interface{}{
		"pageTitle":            "Dashboard",
		"stats":                stats,
		"upcomingEvents":       upcomingEvents,
		"upcomingProcurements": upcomingProcurements,
		"recentActivities":     recentActivities,
	})
}
