package configs

import (
	"starter-project/app/auth"
	"starter-project/app/contracts"
	"starter-project/app/services"
)

// ConfigPageController handles the configs page
type ConfigPageController struct {
	*contracts.GenericPageController
	configService *services.ConfigService
}

// NewConfigPageController creates a new configs page controller
func NewConfigPageController() *ConfigPageController {
	configService := services.NewConfigService()

	return &ConfigPageController{
		GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
			ResourceType:      "configs",
			PageComponent:     "Config/Index",
			Service:           configService,
			ServiceIdentifier: auth.ServiceConfig,
			StatsEnabled:      true,
			StatsBuilder: func(controller *contracts.GenericPageController) map[string]interface{} {
				stats, _ := configService.GetConfigStatistics()
				return stats
			},
		}),
		configService: configService,
	}
}
