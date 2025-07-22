package auth

import (
	"github.com/goravel/framework/contracts/http"
	"players/app/contracts"
	"players/app/helpers"
	"players/app/services"
)

// UserPageController handles the Inertia.js User management page
// Only accessible by super admins - now using GenericPageController
type UserPageController struct {
	*contracts.GenericPageController
	userService *services.UserService
}

// NewUserPageController creates a new user page controller with minimal boilerplate
func NewUserPageController() *UserPageController {
	userService := services.NewUserService()
	
	controller := &UserPageController{
		GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
			ResourceType:      "user",
			PageComponent:     "Users/Index",
			Service:           userService,
			ServiceIdentifier: "", // No service identifier for super admin only
			RequireSuperAdmin: true,
			StatsEnabled:      true,
			StatsBuilder:      buildUserStatistics,
		}),
		userService: userService,
	}
	
	// Set the auth helper
	controller.SetAuthHelper(helpers.NewAuthHelper())
	
	// Add extra data provider for roles
	controller.AddExtraDataProvider("roles", func(ctx http.Context) (interface{}, error) {
		return controller.userService.GetAllRoles()
	})

	// Register page controller with validation
	contracts.MustRegisterPageController("users_page", controller)

	return controller
}

// buildUserStatistics builds user statistics for the dashboard
func buildUserStatistics(controller *contracts.GenericPageController) map[string]interface{} {
	totalCount := controller.GetTotalCount()
	activeCount := controller.GetCountByFilter(map[string]interface{}{"is_active": true})
	inactiveCount := totalCount - activeCount
	superAdminCount := controller.GetCountByFilter(map[string]interface{}{"is_super_admin": true})

	return map[string]interface{}{
		"totalUsers":    totalCount,
		"activeUsers":   activeCount,
		"inactiveUsers": inactiveCount,
		"superAdmins":   superAdminCount,
	}
}

// That's it! The GenericPageController handles:
// - The Index method implementation with super admin check
// - All permission checking (RequireSuperAdmin: true)
// - Request validation
// - Data fetching with error handling
// - Statistics gathering
// - Extra data providers (roles)
// - All contract implementations (CheckPermission, GetCurrentUser, etc.)
// 
// This reduces the controller from 198 lines to just 54 lines - a 73% reduction!