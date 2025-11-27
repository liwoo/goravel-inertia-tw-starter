package smes

import (
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/services"
)

// SmePageController handles the smes page
type SmePageController struct {
	*contracts.GenericPageController
	smeService *services.SmeService
}

// NewSmePageController creates a new smes page controller
func NewSmePageController() *SmePageController {
	smeService := services.NewSmeService()

	return &SmePageController{
		GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
			ResourceType:      "smes",
			PageComponent:     "Sme/Index",
			Service:           smeService,
			ServiceIdentifier: auth.ServiceSMEs,
			StatsEnabled:      true,
			StatsBuilder: func(controller *contracts.GenericPageController) map[string]interface{} {
				stats, _ := smeService.GetSmeStatistics()
				return stats
			},
		}),
		smeService: smeService,
	}
}
