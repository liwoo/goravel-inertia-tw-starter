package bdsps

import (
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/services"
)

// BdspPageController handles the bdsps page
type BdspPageController struct {
	*contracts.GenericPageController
	bdspService *services.BdspService
}

// NewBdspPageController creates a new bdsps page controller
func NewBdspPageController() *BdspPageController {
	bdspService := services.NewBdspService()

	return &BdspPageController{
		GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
			ResourceType:      "bdsps",
			PageComponent:     "Bdsp/Index",
			Service:           bdspService,
			ServiceIdentifier: auth.ServiceBdsps,
			StatsEnabled:      true,
			StatsBuilder: func(controller *contracts.GenericPageController) map[string]interface{} {
				stats, _ := bdspService.GetBdspStatistics()
				return stats
			},
		}),
		bdspService: bdspService,
	}
}
