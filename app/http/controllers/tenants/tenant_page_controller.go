package tenants

import (
	"books-database/app/auth"
	"books-database/app/contracts"
	"books-database/app/services"
)

type TenantPageController struct {
	*contracts.GenericPageController
	tenantService *services.TenantService
}

func NewTenantPageController() *TenantPageController {
	tenantService := services.NewTenantService()
	return &TenantPageController{
		GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
			ResourceType:      "tenants",
			PageComponent:     "Tenant/Index",
			Service:           tenantService,
			ServiceIdentifier: auth.ServiceTenants,
			StatsEnabled:      false,
		}),
		tenantService: tenantService,
	}
}
