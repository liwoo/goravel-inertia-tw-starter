package opportunities

import (
	"smedi-sme-db/app/auth"
	inertiaHelper "smedi-sme-db/app/http/inertia"
	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

// OpportunitiesPageController handles the opportunities page for SME users
type OpportunitiesPageController struct {
	smeService *services.SmeService
}

// NewOpportunitiesPageController creates a new opportunities page controller
func NewOpportunitiesPageController() *OpportunitiesPageController {
	return &OpportunitiesPageController{
		smeService: services.NewSmeService(),
	}
}

// OpportunityItem represents a simplified procurement notice for display
type OpportunityItem struct {
	ID                     uint      `json:"id"`
	RefNo                  string    `json:"ref_no"`
	Organization           string    `json:"organization"`
	ProcuredBy             string    `json:"procured_by"`
	ProcurementType        string    `json:"procurement_type"`
	MarketApproach         string    `json:"market_approach"`
	Invitation             string    `json:"invitation"`
	Details                string    `json:"details"`
	ApplicationDetails     string    `json:"application_details"`
	OpenDate               string    `json:"open_date"`
	CloseDate              string    `json:"close_date"`
	MinimumQualifyingScore int       `json:"minimum_qualifying_score"`
	Classification         []string  `json:"classification"`
	QualifyingDistricts    []string  `json:"qualifying_districts"`
	IsOpen                 bool      `json:"is_open"`
	DaysRemaining          int       `json:"days_remaining"`
}

// Index renders the opportunities page
func (c *OpportunitiesPageController) Index(ctx http.Context) http.Response {
	// Get the authenticated user
	user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
	if user == nil {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	// Get user's linked SME
	sme, err := c.smeService.GetSmeByUserEmail(user.Email)
	if err != nil || sme == nil {
		return inertiaHelper.Render(ctx, "Opportunities/Index", map[string]interface{}{
			"upcomingOpportunities": []OpportunityItem{},
			"pastOpportunities":     []OpportunityItem{},
			"formalisationScore":    0,
			"userName":              user.Name,
			"smeName":               "",
			"error":                 "No SME linked to your account",
		})
	}

	// Get formalisation score
	formalisation, _ := c.smeService.GetBusinessFormalisation(sme.ID)
	formalisationScore := 0
	if formalisation != nil {
		formalisationScore = formalisation.FormalisationScore
	}

	// Get district for filtering
	district := ""
	if sme.District != nil {
		district = *sme.District
	}

	// Fetch opportunities where minimum_qualifying_score <= formalisation_score and is_published = true
	var allOpportunities []models.ProcurementNotice
	query := facades.Orm().Query().
		Where("is_published = ?", true).
		Where("minimum_qualifying_score <= ?", formalisationScore).
		Order("close_date DESC")

	if err := query.Find(&allOpportunities); err != nil {
		facades.Log().Error("Error fetching opportunities", map[string]interface{}{
			"error": err.Error(),
		})
	}

	now := time.Now()
	upcomingOpportunities := make([]OpportunityItem, 0)
	pastOpportunities := make([]OpportunityItem, 0)

	for _, opp := range allOpportunities {
		closeDate := opp.CloseDate.StdTime()
		openDate := opp.OpenDate.StdTime()
		isOpen := closeDate.After(now) && openDate.Before(now)
		daysRemaining := 0
		if closeDate.After(now) {
			daysRemaining = int(closeDate.Sub(now).Hours() / 24)
		}

		// Check if district qualifies (if qualifying_districts is empty, all districts qualify)
		qualifies := true
		if len(opp.QualifyingDistricts) > 0 && district != "" {
			qualifies = false
			for _, d := range opp.QualifyingDistricts {
				if d == district {
					qualifies = true
					break
				}
			}
		}

		if !qualifies {
			continue
		}

		item := OpportunityItem{
			ID:                     opp.ID,
			RefNo:                  opp.RefNo,
			Organization:           opp.Organization,
			ProcuredBy:             opp.ProcuredBy,
			ProcurementType:        opp.ProcurementType,
			MarketApproach:         opp.MarketApproach,
			Invitation:             opp.Invitation,
			Details:                opp.Details,
			ApplicationDetails:     opp.ApplicationDetails,
			OpenDate:               opp.OpenDate.ToDateString(),
			CloseDate:              opp.CloseDate.ToDateString(),
			MinimumQualifyingScore: opp.MinimumQualifyingScore,
			Classification:         opp.Classification,
			QualifyingDistricts:    opp.QualifyingDistricts,
			IsOpen:                 isOpen,
			DaysRemaining:          daysRemaining,
		}

		if closeDate.After(now) {
			upcomingOpportunities = append(upcomingOpportunities, item)
		} else {
			pastOpportunities = append(pastOpportunities, item)
		}
	}

	return inertiaHelper.Render(ctx, "Opportunities/Index", map[string]interface{}{
		"upcomingOpportunities": upcomingOpportunities,
		"pastOpportunities":     pastOpportunities,
		"formalisationScore":    formalisationScore,
		"userName":              user.Name,
		"smeName":               sme.Name,
		"usmeNumber":            sme.UsmeNumber,
		"classification":        sme.Classification,
		"district":              district,
	})
}
