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

	controller.AddExtraDataProvider("calendarEvents", func(ctx http.Context) (interface{}, error) {
		user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
		if user == nil {
			return dashboardService.GetEventsForDistrict(""), nil
		}

		sme, err := smeService.GetSmeByUserEmail(user.Email)
		if err != nil || sme == nil {
			return dashboardService.GetEventsForDistrict(""), nil
		}

		district := ""
		if sme.District != nil {
			district = *sme.District
		}
		// Use the new method that includes isAttending flag
		return dashboardService.GetEventsForDistrictWithAttendance(district, sme.ID), nil
	})

	// Helper function to get user's linked SME
	getSmeForUser := func(ctx http.Context) *uint {
		user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
		if user == nil {
			return nil
		}
		sme, err := smeService.GetSmeByUserEmail(user.Email)
		if err != nil || sme == nil {
			return nil
		}
		return &sme.ID
	}

	controller.AddExtraDataProvider("smeId", func(ctx http.Context) (interface{}, error) {
		smeID := getSmeForUser(ctx)
		if smeID == nil {
			return nil, nil
		}
		return *smeID, nil
	})

	// Set mandatory filter to only show members for the user's linked SME
	controller.SetMandatoryFilterProvider(func(ctx http.Context) map[string]interface{} {
		smeID := getSmeForUser(ctx)
		if smeID == nil {
			// Return impossible filter to show no data if user has no linked SME
			return map[string]interface{}{"sme_id": 0}
		}
		return map[string]interface{}{"sme_id": *smeID}
	})

	controller.AddExtraDataProvider("userName", func(ctx http.Context) (interface{}, error) {
		user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
		if user == nil {
			return nil, nil
		}
		return user.Name, nil
	})

	controller.AddExtraDataProvider("smeName", func(ctx http.Context) (interface{}, error) {
		user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
		if user == nil {
			return nil, nil
		}

		sme, err := smeService.GetSmeByUserEmail(user.Email)
		if err != nil || sme == nil {
			return nil, nil
		}

		return sme.Name, nil
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

	controller.AddExtraDataProvider("registrationNumber", func(ctx http.Context) (interface{}, error) {
		user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
		if user == nil {
			return nil, nil
		}

		sme, err := smeService.GetSmeByUserEmail(user.Email)
		if err != nil || sme == nil {
			return nil, nil
		}

		return sme.RegistrationNumber, nil
	})

	controller.AddExtraDataProvider("taxIdentificationNumber", func(ctx http.Context) (interface{}, error) {
		user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
		if user == nil {
			return nil, nil
		}

		sme, err := smeService.GetSmeByUserEmail(user.Email)
		if err != nil || sme == nil {
			return nil, nil
		}

		return sme.TaxIdentificationNumber, nil
	})

	controller.AddExtraDataProvider("classification", func(ctx http.Context) (interface{}, error) {
		user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
		if user == nil {
			return nil, nil
		}

		sme, err := smeService.GetSmeByUserEmail(user.Email)
		if err != nil || sme == nil {
			return nil, nil
		}

		return sme.Classification, nil
	})

	controller.AddExtraDataProvider("usmeNumber", func(ctx http.Context) (interface{}, error) {
		user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
		if user == nil {
			return nil, nil
		}

		sme, err := smeService.GetSmeByUserEmail(user.Email)
		if err != nil || sme == nil {
			return nil, nil
		}

		return sme.UsmeNumber, nil
	})

	controller.AddExtraDataProvider("district", func(ctx http.Context) (interface{}, error) {
		user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
		if user == nil {
			return nil, nil
		}

		sme, err := smeService.GetSmeByUserEmail(user.Email)
		if err != nil || sme == nil {
			return nil, nil
		}

		return sme.District, nil
	})

	return controller
}
