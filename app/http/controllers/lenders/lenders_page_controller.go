package lenders

import (
	"books-database/app/auth"
	"books-database/app/contracts"
	"books-database/app/services"
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
