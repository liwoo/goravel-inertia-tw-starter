package directory

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"

	inertiaHelper "smedi-sme-db/app/http/inertia"
	"smedi-sme-db/app/http/requests"
	"smedi-sme-db/app/models"
)

// DirectoryController handles the public SME directory
type DirectoryController struct{}

// NewDirectoryController creates a new directory controller
func NewDirectoryController() *DirectoryController {
	return &DirectoryController{}
}

// FilterOption represents a filter option with count
type FilterOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

// DistrictGroup represents districts grouped by region
type DistrictGroup struct {
	Region    string         `json:"region"`
	Districts []FilterOption `json:"districts"`
}

// DirectoryFilters contains all available filter options
type DirectoryFilters struct {
	Sectors   []FilterOption  `json:"sectors"`
	Districts []DistrictGroup `json:"districts"`
}

// PublicSmeDTO represents the public-safe SME data
type PublicSmeDTO struct {
	UsmeNumber           string  `json:"usme_number"`
	Name                 string  `json:"name"`
	Sector               string  `json:"sector"`
	SubSector            *string `json:"sub_sector"`
	Classification       string  `json:"classification"`
	District             *string `json:"district"`
	Region               *string `json:"region"`
	ContactPhone         string  `json:"contact_phone"`
	ContactEmail         string  `json:"contact_email"`
	Website              *string `json:"website"`
	BusinessDescription  *string `json:"business_description"`
	OperationalStartDate *string `json:"operational_start_date"`
}

// PaginatedDirectoryResponse represents the paginated response for directory
type PaginatedDirectoryResponse struct {
	Data       []PublicSmeDTO `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"pageSize"`
	TotalPages int            `json:"totalPages"`
}

// ShowDirectory renders the Inertia page with initial data
func (c *DirectoryController) ShowDirectory(ctx http.Context) http.Response {
	// Get filter options
	filters := c.getFilterOptions()

	// Start with empty results - user must search or filter
	return inertiaHelper.Render(ctx, "Directory/Index", map[string]interface{}{
		"version": support.Version,
		"smes":    []PublicSmeDTO{},
		"filters": filters,
		"pagination": map[string]interface{}{
			"total":      0,
			"page":       1,
			"pageSize":   12,
			"totalPages": 0,
		},
	})
}

// SearchDirectory handles AJAX search/filter requests
func (c *DirectoryController) SearchDirectory(ctx http.Context) http.Response {
	// Parse query parameters
	search := ctx.Request().Query("search", "")
	page, _ := strconv.Atoi(ctx.Request().Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Request().Query("pageSize", "12"))

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 12
	}

	// Get sectors filter (array)
	sectors := ctx.Request().QueryArray("sectors[]")
	if len(sectors) == 0 {
		sectors = ctx.Request().QueryArray("sectors")
	}

	// Get districts filter (array)
	districts := ctx.Request().QueryArray("districts[]")
	if len(districts) == 0 {
		districts = ctx.Request().QueryArray("districts")
	}

	// Fetch SMEs
	smes, total := c.fetchPublicSmes(search, sectors, districts, page, pageSize)
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return ctx.Response().Success().Json(http.Json{
		"smes": smes,
		"pagination": http.Json{
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": totalPages,
		},
	})
}

// fetchPublicSmes retrieves active SMEs with filters and pagination
func (c *DirectoryController) fetchPublicSmes(search string, sectors, districts []string, page, pageSize int) ([]PublicSmeDTO, int64) {
	var smes []models.Sme
	var total int64

	// Build base query for active SMEs only
	query := facades.Orm().Query().Model(&models.Sme{}).Where("is_active = ?", true)

	// Apply search filter
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(func(q interface{}) {
			// Type assert to get the underlying query builder
		}).Where("(name ILIKE ? OR usme_number ILIKE ? OR sector ILIKE ? OR district ILIKE ? OR business_description ILIKE ?)",
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Apply sectors filter
	if len(sectors) > 0 {
		// Convert []string to []interface{} for WhereIn
		sectorInterfaces := make([]interface{}, len(sectors))
		for i, s := range sectors {
			sectorInterfaces[i] = s
		}
		query = query.WhereIn("sector", sectorInterfaces)
	}

	// Apply districts filter
	if len(districts) > 0 {
		// Convert []string to []interface{} for WhereIn
		districtInterfaces := make([]interface{}, len(districts))
		for i, d := range districts {
			districtInterfaces[i] = d
		}
		query = query.WhereIn("district", districtInterfaces)
	}

	// Get total count
	total, _ = query.Count()

	// Apply pagination and fetch data
	offset := (page - 1) * pageSize
	query.Order("name ASC").Offset(offset).Limit(pageSize).Find(&smes)

	// Convert to DTOs
	dtos := make([]PublicSmeDTO, len(smes))
	for i, sme := range smes {
		dtos[i] = c.toPublicDTO(sme)
	}

	return dtos, total
}

// toPublicDTO converts an SME model to a public-safe DTO
func (c *DirectoryController) toPublicDTO(sme models.Sme) PublicSmeDTO {
	dto := PublicSmeDTO{
		UsmeNumber:     sme.UsmeNumber,
		Name:           sme.Name,
		Sector:         sme.Sector,
		SubSector:      sme.SubSector,
		Classification: sme.Classification,
		District:       sme.District,
		Region:         sme.Region,
		ContactPhone:   sme.ContactPhone,
		ContactEmail:   sme.ContactEmail,
		Website:        sme.Website,
	}

	// Truncate business description to 200 characters
	if sme.BusinessDescription != nil && len(*sme.BusinessDescription) > 200 {
		truncated := (*sme.BusinessDescription)[:200] + "..."
		dto.BusinessDescription = &truncated
	} else {
		dto.BusinessDescription = sme.BusinessDescription
	}

	// Format operational start date
	if sme.OperationalStartDate != nil {
		dateStr := sme.OperationalStartDate.ToDateString()
		dto.OperationalStartDate = &dateStr
	}

	return dto
}

// getFilterOptions retrieves available filter options with counts
func (c *DirectoryController) getFilterOptions() DirectoryFilters {
	return DirectoryFilters{
		Sectors:   c.getSectorOptions(),
		Districts: c.getDistrictOptions(),
	}
}

// getSectorOptions retrieves distinct sectors with counts from active SMEs
func (c *DirectoryController) getSectorOptions() []FilterOption {
	type SectorCount struct {
		Sector string
		Count  int
	}

	var results []SectorCount
	query := `
		SELECT sector, COUNT(*) as count
		FROM smes
		WHERE is_active = true AND deleted_at IS NULL AND sector IS NOT NULL AND sector != ''
		GROUP BY sector
		ORDER BY count DESC, sector ASC
	`
	facades.Orm().Query().Raw(query).Scan(&results)

	options := make([]FilterOption, len(results))
	for i, r := range results {
		options[i] = FilterOption{
			Value: r.Sector,
			Label: r.Sector,
			Count: r.Count,
		}
	}

	return options
}

// getDistrictOptions retrieves districts grouped by region with counts
func (c *DirectoryController) getDistrictOptions() []DistrictGroup {
	type DistrictCount struct {
		District string
		Count    int
	}

	var results []DistrictCount
	query := `
		SELECT district, COUNT(*) as count
		FROM smes
		WHERE is_active = true AND deleted_at IS NULL AND district IS NOT NULL AND district != ''
		GROUP BY district
		ORDER BY district ASC
	`
	facades.Orm().Query().Raw(query).Scan(&results)

	// Create a map of district to count
	districtCounts := make(map[string]int)
	for _, r := range results {
		districtCounts[r.District] = r.Count
	}

	// Group districts by region using the existing region/district mapping
	regions := requests.GetAllRegions()
	districtGroups := make([]DistrictGroup, len(regions))

	for i, region := range regions {
		regionDistricts := requests.GetDistrictsByRegion(region)
		districtOptions := make([]FilterOption, 0)

		for _, district := range regionDistricts {
			districtStr := string(district)
			count := districtCounts[districtStr]
			// Only include districts that have at least one active SME
			if count > 0 {
				districtOptions = append(districtOptions, FilterOption{
					Value: districtStr,
					Label: districtStr,
					Count: count,
				})
			}
		}

		districtGroups[i] = DistrictGroup{
			Region:    string(region),
			Districts: districtOptions,
		}
	}

	return districtGroups
}
