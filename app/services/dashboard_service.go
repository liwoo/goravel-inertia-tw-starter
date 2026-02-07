package services

import (
	"books-database/app/models"
	"github.com/goravel/framework/facades"
)

// DashboardService provides methods for dashboard statistics and widgets
type DashboardService struct{}

// NewDashboardService creates a new dashboard service
func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

// GetDashboardStats returns all statistics for the dashboard
func (s *DashboardService) GetDashboardStats() map[string]interface{} {
	var totalBooks int64
	totalBooks, _ = facades.Orm().Query().Model(&models.Book{}).Count()

	var totalUsers int64
	totalUsers, _ = facades.Orm().Query().Model(&models.User{}).Count()

	return map[string]interface{}{
		"totalBooks": totalBooks,
		"totalUsers": totalUsers,
	}
}
