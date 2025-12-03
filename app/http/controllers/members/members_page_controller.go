package members

import (
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/services"

	"github.com/goravel/framework/contracts/http"
)

// MemberPageController handles the members page
type MemberPageController struct {
	*contracts.GenericPageController
	memberService *services.AdditionalBusinessMemberService
}

// NewMemberPageController creates a new members page controller
func NewMemberPageController() *MemberPageController {
	memberService := services.NewAdditionalBusinessMemberService()
	dashboardService := services.NewDashboardService()
	smeService := services.NewSmeService()

	controller := &MemberPageController{
		GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
			ResourceType:      "members",
			PageComponent:     "Member/Index",
			Service:           memberService,
			ServiceIdentifier: auth.ServiceAdditionalBusinessMembers,
			StatsEnabled:      false,
		}),
		memberService: memberService,
	}

	controller.AddExtraDataProvider("events", func(ctx http.Context) (interface{}, error) {
		return dashboardService.GetUpcomingEvents(3), nil
	})

	controller.AddExtraDataProvider("procurements", func(ctx http.Context) (interface{}, error) {
		return dashboardService.GetUpcomingProcurements(3), nil
	})

	controller.AddExtraDataProvider("smeId", func(ctx http.Context) (interface{}, error) {
		user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
		if user == nil {
			return nil, nil
		}

		sme, err := smeService.GetSmeByUserEmail(user.Email)
		if err != nil || sme == nil {
			return nil, nil
		}

		return sme.ID, nil
	})

	controller.AddExtraDataProvider("formalisation", func(ctx http.Context) (interface{}, error) {
		user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
		if user == nil {
			return nil, nil
		}

		sme, err := smeService.GetSmeByUserEmail(user.Email)
		if err != nil || sme == nil {
			return nil, nil
		}

		formalisation, err := smeService.GetBusinessFormalisation(sme.ID)
		if err != nil {
			return nil, nil
		}

		return formalisation, nil
	})

	return controller
}
