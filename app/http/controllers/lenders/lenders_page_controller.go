package lenders

import (
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/services"
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
