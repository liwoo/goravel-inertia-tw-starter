package lenders

import (
	"starter-project/app/auth"
	"starter-project/app/contracts"
	"starter-project/app/services"
)

// LenderPageController handles the lenders page
type LenderPageController struct {
	*contracts.GenericPageController
	lenderService *services.LenderService
}

// NewLenderPageController creates a new lenders page controller
func NewLenderPageController() *LenderPageController {
	lenderService := services.NewLenderService()

	return &LenderPageController{
		GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
			ResourceType:      "lenders",
			PageComponent:     "Lender/Index",
			Service:           lenderService,
			ServiceIdentifier: auth.ServiceLenders,
			StatsEnabled:      false,
		}),
		lenderService: lenderService,
	}
}
