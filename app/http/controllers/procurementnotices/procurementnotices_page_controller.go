package procurementnotices

import (
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/services"
)

// ProcurementNoticePageController handles the procurementnotices page
type ProcurementNoticePageController struct {
	*contracts.GenericPageController
	procurementnoticeService *services.ProcurementNoticeService
}

// NewProcurementNoticePageController creates a new procurementnotices page controller
func NewProcurementNoticePageController() *ProcurementNoticePageController {
	procurementnoticeService := services.NewProcurementNoticeService()

	return &ProcurementNoticePageController{
		GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
			ResourceType:      "procurementnotices",
			PageComponent:     "ProcurementNotice/Index",
			Service:           procurementnoticeService,
			ServiceIdentifier: auth.ServiceProcurementNotices,
			StatsEnabled:      false,
		}),
		procurementnoticeService: procurementnoticeService,
	}
}
