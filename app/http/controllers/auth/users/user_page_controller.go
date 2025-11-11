package users

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/services"
)

// UserPageController handles the users page
type UserPageController struct {
	*contracts.GenericPageController
	userService *services.UserService
}

// NewUserPageController creates a new users page controller
func NewUserPageController() *UserPageController {
	userService := services.NewUserService()

	controller := &UserPageController{
		GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
			ResourceType:      "users",
			PageComponent:     "Users/Index",
			Service:           userService,
			ServiceIdentifier: auth.ServiceUsers,
			RequireSuperAdmin: true,
			StatsEnabled:      true,
			StatsBuilder: func(controller *contracts.GenericPageController) map[string]interface{} {
				stats, _ := userService.GetUserStatistics()
				return stats
			},
		}),
		userService: userService,
	}

	// Add roles data provider
	controller.GenericPageController.AddExtraDataProvider("roles", func(ctx http.Context) (interface{}, error) {
		roles, err := userService.GetAllRoles()
		if err != nil {
			// Log the error
			facades.Log().Error("Failed to get roles", map[string]interface{}{
				"error": err.Error(),
			})
			return []interface{}{}, nil // Return empty array on error
		}
		// Log success
		facades.Log().Info("Successfully fetched roles", map[string]interface{}{
			"count": len(roles),
		})
		return roles, nil
	})

	return controller
}
