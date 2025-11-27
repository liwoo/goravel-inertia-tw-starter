package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/http/requests"
	"smedi-sme-db/app/models"
)

// SmeService implements business logic for smes using the builder pattern
type SmeService struct {
	contracts.CrudServiceContract // Embedded interface - automatically exposes all CRUD methods!
}

// NewSmeService creates a new Sme service using the builder pattern
func NewSmeService() *SmeService {
	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.Sme]("smes", "id").
		WithSearchFields("usme_number", "name", "registration_number", "tax_identification_number", "business_category", "sector", "contact_email", "contact_phone", "region", "district"). // Fields that will be searchable via the search query parameter
		WithSortFields("id", "created_at", "updated_at", "usme_number", "name", "operational_start_date", "business_category", "sector", "region", "district").                             // Fields that can be used for sorting results
		WithFilterFields("business_category", "sector", "region", "district", "created_by").                                                                                                // Fields that can be filtered on
		WithValidationRules(map[string]interface{}{                                                                                                                                         // Validation rules for create/update operations
			// usme_number is omitted from validation - will be auto-generated in BeforeCreate hook if not provided
			"name":                           "required|string|max:255",
			"registration_number":            "string|max:100",
			"tax_identification_number":      "string|max:100",
			"operational_start_date":         "date",
			"business_category":              "required|string|max:100",
			"sector":                         "required|string|max:100",
			"sub_sector":                     "sometimes|string|max:100",
			"business_description":           "sometimes|string|max:1000",
			"contact_phone":                  "required|string|max:20",
			"contact_email":                  "required|email|max:100",
			"physical_address":               "string|max:255",
			"postal_address":                 "string|max:255",
			"website":                        "url|max:255",
			"region":                         "string|max:100",
			"district":                       "string|max:100",
			"traditional_authority":          "string|max:100",
			"business_improvement_aspects":   "required|array",
			"business_improvement_aspects.*": "string|max:100",
			"business_accessed_financing":    "required|array",
			"business_accessed_financing.*":  "string|max:100",
		}).
		WithRelations("PrimaryBusinessOwner", "AdditionalBusinessMembers", "BusinessFormalisation", "BusinessEmployeeSummary"). // Optional - load relationships
		WithDefaultSort("created_at", "DESC").                                                                                  // Default sorting when none specified
		WithScopeFiltering("smes", "created_by").                                                                               // Enable permission-based filtering
		WithSoftDeletes().                                                                                                      // Enable soft delete support
		WithBeforeCreate(func(data map[string]interface{}) error {                                                              // Handle JSON array fields and UBI generation
			// Generate UBI if usme_number is not provided
			if _, exists := data["usme_number"]; !exists || data["usme_number"] == "" {
				ubi, err := generateUBI(data)
				if err != nil {
					return fmt.Errorf("failed to generate UBI: %w", err)
				}
				data["usme_number"] = ubi
			}

			// Handle business_improvement_aspects array to JSON conversion
			if aspects, exists := data["business_improvement_aspects"]; exists && aspects != nil {
				if aspectsArray, ok := aspects.([]interface{}); ok && len(aspectsArray) > 0 {
					stringAspects := make([]string, len(aspectsArray))
					for i, aspect := range aspectsArray {
						stringAspects[i] = fmt.Sprintf("%v", aspect)
					}
					aspectsJSON, _ := json.Marshal(stringAspects)
					data["business_improvement_aspect_json"] = string(aspectsJSON)
				} else if aspectsArray, ok := aspects.([]string); ok && len(aspectsArray) > 0 {
					aspectsJSON, _ := json.Marshal(aspectsArray)
					data["business_improvement_aspect_json"] = string(aspectsJSON)
				} else {
					data["business_improvement_aspect_json"] = "[]"
				}
			} else {
				data["business_improvement_aspect_json"] = "[]"
			}

			// Handle business_accessed_financing array to JSON conversion
			if financing, exists := data["business_accessed_financing"]; exists && financing != nil {
				if financingArray, ok := financing.([]interface{}); ok && len(financingArray) > 0 {
					stringFinancing := make([]string, len(financingArray))
					for i, item := range financingArray {
						stringFinancing[i] = fmt.Sprintf("%v", item)
					}
					financingJSON, _ := json.Marshal(stringFinancing)
					data["business_accessed_financing_json"] = string(financingJSON)
				} else if financingArray, ok := financing.([]string); ok && len(financingArray) > 0 {
					financingJSON, _ := json.Marshal(financingArray)
					data["business_accessed_financing_json"] = string(financingJSON)
				} else {
					data["business_accessed_financing_json"] = "[]"
				}
			} else {
				data["business_accessed_financing_json"] = "[]"
			}

			return nil
		}).
		WithBeforeUpdate(func(id uint, data map[string]interface{}) error { // Handle JSON array fields and generate USME if missing
			// Check if SME exists and has a usme_number
			var existingSme models.Sme
			err := facades.Orm().Query().Where("id = ?", id).First(&existingSme)
			if err == nil && existingSme.UsmeNumber == "" {
				// Generate USME number if missing
				// Need to merge existing data with update data for generation
				mergedData := make(map[string]interface{})

				// Copy essential fields from existing record
				if existingSme.District != nil {
					mergedData["district"] = *existingSme.District
				}
				if existingSme.Region != nil {
					mergedData["region"] = *existingSme.Region
				}
				mergedData["business_category"] = existingSme.BusinessCategory

				// Override with any new values from update data
				if district, exists := data["district"]; exists {
					mergedData["district"] = district
				}
				if region, exists := data["region"]; exists {
					mergedData["region"] = region
				}
				if category, exists := data["business_category"]; exists {
					mergedData["business_category"] = category
				}

				// Generate UBI
				ubi, err := generateUBI(mergedData)
				if err != nil {
					facades.Log().Warning("Failed to generate USME number during update", map[string]interface{}{"error": err.Error(), "sme_id": id})
				} else {
					data["usme_number"] = ubi
				}
			}

			// Handle business_improvement_aspects array to JSON conversion
			if aspects, exists := data["business_improvement_aspects"]; exists && aspects != nil {
				if aspectsArray, ok := aspects.([]interface{}); ok && len(aspectsArray) > 0 {
					stringAspects := make([]string, len(aspectsArray))
					for i, aspect := range aspectsArray {
						stringAspects[i] = fmt.Sprintf("%v", aspect)
					}
					aspectsJSON, _ := json.Marshal(stringAspects)
					data["business_improvement_aspect_json"] = string(aspectsJSON)
				} else if aspectsArray, ok := aspects.([]string); ok && len(aspectsArray) > 0 {
					aspectsJSON, _ := json.Marshal(aspectsArray)
					data["business_improvement_aspect_json"] = string(aspectsJSON)
				} else {
					data["business_improvement_aspect_json"] = "[]"
				}
				delete(data, "business_improvement_aspects") // Remove the array field to avoid column conflict
			}

			// Handle business_accessed_financing array to JSON conversion
			if financing, exists := data["business_accessed_financing"]; exists && financing != nil {
				if financingArray, ok := financing.([]interface{}); ok && len(financingArray) > 0 {
					stringFinancing := make([]string, len(financingArray))
					for i, item := range financingArray {
						stringFinancing[i] = fmt.Sprintf("%v", item)
					}
					financingJSON, _ := json.Marshal(stringFinancing)
					data["business_accessed_financing_json"] = string(financingJSON)
				} else if financingArray, ok := financing.([]string); ok && len(financingArray) > 0 {
					financingJSON, _ := json.Marshal(financingArray)
					data["business_accessed_financing_json"] = string(financingJSON)
				} else {
					data["business_accessed_financing_json"] = "[]"
				}
				delete(data, "business_accessed_financing") // Remove the array field to avoid column conflict
			}

			return nil
		}).
		Build() // Returns a fully configured CrudServiceContract

	smeServiceInstance := &SmeService{
		CrudServiceContract: service, // Set the embedded interface
	}

	// Set the actual service instance for proper method resolution
	contracts.SetActualServiceHelper(service, smeServiceInstance, "SmeService")

	return smeServiceInstance
}

// Delete overrides the base Delete method to cascade soft delete all related records
func (s *SmeService) Delete(id uint) error {
	// Call the base delete method first
	if err := s.CrudServiceContract.Delete(id); err != nil {
		return err
	}

	// Cascade soft delete all related records
	// Soft delete primary business owner
	_, err := facades.Orm().Query().Where("sme_id = ?", id).Delete(&models.PrimaryBusinessOwner{})
	if err != nil {
		facades.Log().Error("Failed to cascade delete primary business owner", map[string]interface{}{"error": err.Error(), "sme_id": id})
		return err
	}

	// Soft delete additional business members
	_, err = facades.Orm().Query().Where("sme_id = ?", id).Delete(&models.AdditionalBusinessMember{})
	if err != nil {
		facades.Log().Error("Failed to cascade delete additional business members", map[string]interface{}{"error": err.Error(), "sme_id": id})
		return err
	}

	// Soft delete business formalisation
	_, err = facades.Orm().Query().Where("sme_id = ?", id).Delete(&models.BusinessFormalisation{})
	if err != nil {
		facades.Log().Error("Failed to cascade delete business formalisation", map[string]interface{}{"error": err.Error(), "sme_id": id})
		return err
	}

	// Soft delete business employee summary
	_, err = facades.Orm().Query().Where("sme_id = ?", id).Delete(&models.BusinessEmployeeSummary{})
	if err != nil {
		facades.Log().Error("Failed to cascade delete business employee summary", map[string]interface{}{"error": err.Error(), "sme_id": id})
		return err
	}

	return nil
}

// Override GetColumnMapping to include SME-specific mappings
func (s *SmeService) GetColumnMapping() map[string]string {
	mapping := s.CrudServiceContract.GetColumnMapping()
	// Add SME-specific camelCase to snake_case mappings
	mapping["usmeNumber"] = "usme_number"
	mapping["registrationNumber"] = "registration_number"
	mapping["taxIdentificationNumber"] = "tax_identification_number"
	mapping["operationalStartDate"] = "operational_start_date"
	mapping["businessCategory"] = "business_category"
	mapping["subSector"] = "sub_sector"
	mapping["businessDescription"] = "business_description"
	mapping["contactPhone"] = "contact_phone"
	mapping["contactEmail"] = "contact_email"
	mapping["physicalAddress"] = "physical_address"
	mapping["postalAddress"] = "postal_address"
	mapping["traditionalAuthority"] = "traditional_authority"
	mapping["createdAt"] = "created_at"
	mapping["updatedAt"] = "updated_at"
	mapping["createdBy"] = "created_by"
	mapping["updatedBy"] = "updated_by"
	mapping["deletedBy"] = "deleted_by"
	mapping["ipAddress"] = "ip_address"
	mapping["userAgent"] = "user_agent"
	// Note: We do NOT map business_improvement_aspects or business_accessed_financing here
	// because they need to go through the model's BeforeSave hook which converts the arrays to JSON
	return mapping
}

// GetFilterDefinitions returns filter definitions for the SMEs resource
func (s *SmeService) GetFilterDefinitions() []contracts.FilterDefinition {
	// Get all regions as strings for the enum
	regions := requests.GetAllRegions()
	regionStrings := make([]string, len(regions))
	for i, r := range regions {
		regionStrings[i] = string(r)
	}

	// Get all districts as strings for the enum
	districts := requests.GetAllDistricts()
	districtStrings := make([]string, len(districts))
	for i, d := range districts {
		districtStrings[i] = string(d)
	}

	return []contracts.FilterDefinition{
		// USME Number - string search
		contracts.NewFilterDefinition(
			"usme_number",
			"USME Number",
			contracts.FilterTypeString,
			nil, // Use all string operators (equals, contains, starts_with, etc.)
		),
		// Business Name - string search
		contracts.NewFilterDefinition(
			"name",
			"Business Name",
			contracts.FilterTypeString,
			nil,
		),
		// Business Category - string search (could be enum if categories are fixed)
		contracts.NewFilterDefinition(
			"business_category",
			"Business Category",
			contracts.FilterTypeString,
			nil,
		),
		// Sector - string search (could be enum if sectors are fixed)
		contracts.NewFilterDefinition(
			"sector",
			"Sector",
			contracts.FilterTypeString,
			nil,
		),
		// Region - enum with all 3 Malawian regions
		contracts.NewFilterDefinition(
			"region",
			"Region",
			contracts.FilterTypeEnum,
			&regionStrings,
		),
		// District - enum with all 28 Malawian districts
		contracts.NewFilterDefinition(
			"district",
			"District",
			contracts.FilterTypeEnum,
			&districtStrings,
		),
		// Operational Start Date - date filter
		contracts.NewFilterDefinition(
			"operational_start_date",
			"Operational Since",
			contracts.FilterTypeDate,
			nil, // Use all date operators (before, after, between, etc.)
		),
		// Created Date - datetime filter
		contracts.NewFilterDefinition(
			"created_at",
			"Date Registered",
			contracts.FilterTypeDateTime,
			nil,
		),
	}
}

// Add domain-specific methods below this line

// GetPrimaryBusinessOwner retrieves the primary business owner for a given SME
func (s *SmeService) GetPrimaryBusinessOwner(smeID uint) (*models.PrimaryBusinessOwner, error) {
	var sme models.Sme
	err := facades.Orm().Query().
		With("PrimaryBusinessOwner").
		Where("id = ?", smeID).
		First(&sme)

	if err != nil {
		return nil, err
	}

	// Check if SME was actually found (ID will be 0 if not found)
	if sme.ID == 0 {
		return nil, errors.New("SME not found")
	}

	return sme.PrimaryBusinessOwner, nil
}

// GetAdditionalBusinessMembers retrieves all additional business members for a given SME
func (s *SmeService) GetAdditionalBusinessMembers(smeID uint) ([]models.AdditionalBusinessMember, error) {
	var sme models.Sme
	err := facades.Orm().Query().
		With("AdditionalBusinessMembers").
		Where("id = ?", smeID).
		First(&sme)

	if err != nil {
		return nil, err
	}

	// Check if SME was actually found (ID will be 0 if not found)
	if sme.ID == 0 {
		return nil, errors.New("SME not found")
	}

	return sme.AdditionalBusinessMembers, nil
}

// GetBusinessFormalisation retrieves the business formalisation record for a given SME
func (s *SmeService) GetBusinessFormalisation(smeID uint) (*models.BusinessFormalisation, error) {
	var sme models.Sme
	err := facades.Orm().Query().
		With("BusinessFormalisation").
		Where("id = ?", smeID).
		First(&sme)

	if err != nil {
		return nil, err
	}

	// Check if SME was actually found (ID will be 0 if not found)
	if sme.ID == 0 {
		return nil, errors.New("SME not found")
	}

	return sme.BusinessFormalisation, nil
}

// GetBusinessEmployeeSummary retrieves the employee summary for a given SME
func (s *SmeService) GetBusinessEmployeeSummary(smeID uint) (*models.BusinessEmployeeSummary, error) {
	var sme models.Sme
	err := facades.Orm().Query().
		With("BusinessEmployeeSummary").
		Where("id = ?", smeID).
		First(&sme)

	if err != nil {
		return nil, err
	}

	// Check if SME was actually found (ID will be 0 if not found)
	if sme.ID == 0 {
		return nil, errors.New("SME not found")
	}

	return sme.BusinessEmployeeSummary, nil
}

// ============================================================================
// UBI (Universal Business Identifier) Generation Functions
// ============================================================================

// generateUBI generates a UBI in the format: MW-YYYY-DD-CT-NNNNNN-C
// Components:
// - MW: Country Code (Malawi)
// - YYYY: Registration/Assignment Year (2025-2099)
// - DD: District Code (01-28)
// - CT: Category + Type (M1-M4)
// - NNNNNN: Sequential Number (000001-999999)
// - C: Check Digit (Luhn algorithm)
func generateUBI(data map[string]interface{}) (string, error) {
	// Get current year
	year := time.Now().Year()

	// Get district code
	districtCode, err := getDistrictCode(data)
	if err != nil {
		return "", err
	}

	// Get category code
	categoryCode, err := getCategoryCode(data)
	if err != nil {
		return "", err
	}

	// Get next sequential number
	sequentialNumber, err := getNextSequentialNumber(year, districtCode, categoryCode)
	if err != nil {
		return "", err
	}

	// Build UBI without check digit
	ubiWithoutCheck := fmt.Sprintf("MW-%04d-%s-%s-%06d", year, districtCode, categoryCode, sequentialNumber)

	// Calculate check digit
	checkDigit := calculateLuhnCheckDigit(ubiWithoutCheck)

	// Return complete UBI
	return fmt.Sprintf("%s-%d", ubiWithoutCheck, checkDigit), nil
}

// getDistrictCode returns the 2-letter district code for a given district name based on ISO 3166-2:MW
func getDistrictCode(data map[string]interface{}) (string, error) {
	// District code mapping for all 28 Malawian districts using ISO 3166-2:MW two-letter codes
	districtCodes := map[string]string{
		// Northern Region
		"Chitipa":    "CT",
		"Karonga":    "KR",
		"Mzuzu":      "MZ",
		"Nkhata Bay": "NB",
		"Rumphi":     "RU",
		"Likoma":     "LK",
		"Mzimba":     "MH", // Added missing district

		// Central Region
		"Dedza":      "DE",
		"Dowa":       "DO",
		"Kasungu":    "KS",
		"Lilongwe":   "LI",
		"Mchinji":    "MC",
		"Nkhotakota": "NK",
		"Ntcheu":     "NU",
		"Ntchisi":    "NI",
		"Salima":     "SA",

		// Southern Region
		"Balaka":     "BA", // Confirmed from ISO 3166-2:MW
		"Blantyre":   "BT", // Confirmed from user specification
		"Chikwawa":   "CK",
		"Chiradzulu": "CR",
		"Machinga":   "MG",
		"Mangochi":   "MN",
		"Mulanje":    "MJ",
		"Mwanza":     "MW",
		"Nsanje":     "NS",
		"Thyolo":     "TH",
		"Phalombe":   "PH",
		"Zomba":      "ZO",
		"Neno":       "NE",
	}

	district, ok := data["district"]
	if !ok || district == nil {
		return "00", nil // Default district code if not provided
	}

	var districtStr string
	switch v := district.(type) {
	case string:
		districtStr = v
	case *string:
		if v != nil {
			districtStr = *v
		} else {
			return "00", nil
		}
	default:
		return "", fmt.Errorf("district must be a string, got %T", district)
	}

	code, exists := districtCodes[districtStr]
	if !exists {
		return "", fmt.Errorf("unknown district: %s", districtStr)
	}

	return code, nil
}

// getCategoryCode returns the category code based on business category
func getCategoryCode(data map[string]interface{}) (string, error) {

	category, ok := data["business_category"]
	if !ok || category == nil {
		return "", fmt.Errorf("business_category is required")
	}

	var categoryStr string
	switch v := category.(type) {
	case string:
		categoryStr = v
	case *string:
		if v != nil {
			categoryStr = *v
		} else {
			return "", fmt.Errorf("business_category is required")
		}
	default:
		return "", fmt.Errorf("business_category must be a string, got %T", category)
	}

	// if category has spaces we get the first letter of the first two words, e.g., "Micro Enterprise" -> "ME"
	// else we get the first two letters, e.g., "Micro" -> "MI"
	words := strings.Fields(categoryStr)

	var code string
	if len(words) >= 2 {
		code = strings.ToUpper(string(words[0][0]) + string(words[1][0]))
	} else if len(words) == 1 {
		if len(words[0]) >= 2 {
			code = strings.ToUpper(words[0][:2])
		} else {
			code = strings.ToUpper(string(words[0][0]) + "X") // Pad with X if only one letter
		}
	} else {
		return "", fmt.Errorf("invalid business_category: %s", categoryStr)
	}

	return code, nil
}

// getNextSequentialNumber gets the next sequential number for the given year, district, and category
func getNextSequentialNumber(year int, districtCode, categoryCode string) (int, error) {
	// Query the database to find the highest sequential number for this year/district/category combination
	var maxSequence int

	// Pattern to match: MW-YYYY-DD-CT-
	pattern := fmt.Sprintf("MW-%04d-%s-%s-%%", year, districtCode, categoryCode)

	var sme models.Sme
	err := facades.Orm().Query().
		Where("usme_number LIKE ?", pattern).
		Order("usme_number DESC").
		First(&sme)

	if err != nil || sme.ID == 0 {
		// No existing records, start from 1
		return 1, nil
	}

	// Extract the sequential number from the usme_number
	// Format: MW-YYYY-DD-CT-NNNNNN-C
	parts := splitUBI(sme.UsmeNumber)
	if len(parts) >= 5 {
		var parseErr error
		maxSequence, parseErr = strconv.Atoi(parts[4])
		if parseErr != nil {
			// If parsing fails, start from 1
			return 1, nil
		}
	}

	// Return next number
	return maxSequence + 1, nil
}

// splitUBI splits a UBI string into its components
func splitUBI(ubi string) []string {
	return strings.Split(ubi, "-")
}

// calculateLuhnCheckDigit calculates the Luhn check digit for the given UBI string
func calculateLuhnCheckDigit(ubi string) int {
	// Remove all non-digit characters
	digits := ""
	for _, char := range ubi {
		if char >= '0' && char <= '9' {
			digits += string(char)
		}
	}

	// Luhn algorithm
	sum := 0
	alternate := false

	// Process digits from right to left
	for i := len(digits) - 1; i >= 0; i-- {
		digit := int(digits[i] - '0')

		if alternate {
			digit *= 2
			if digit > 9 {
				digit = digit - 9
			}
		}

		sum += digit
		alternate = !alternate
	}

	// The check digit is the amount needed to make the sum a multiple of 10
	checkDigit := (10 - (sum % 10)) % 10
	return checkDigit
}

// ============================================================================
// SME Statistics Functions
// ============================================================================

// GetSmeStatistics returns comprehensive statistics about SMEs for dashboard KPIs and charts
func (s *SmeService) GetSmeStatistics() (map[string]interface{}, error) {
	now := time.Now()
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastMonthStart := currentMonthStart.AddDate(0, -1, 0)
	lastMonthEnd := currentMonthStart.Add(-time.Second)

	// Basic counts
	var totalSmes int64
	var newThisMonth int64
	var newLastMonth int64
	var withRegistrationNumber int64

	// Get total SMEs (excluding soft deleted)
	totalSmes, _ = facades.Orm().Query().Model(&models.Sme{}).Where("deleted_at IS NULL").Count()

	// Get SMEs created this month
	newThisMonth, _ = facades.Orm().Query().Model(&models.Sme{}).
		Where("deleted_at IS NULL").
		Where("created_at >= ?", currentMonthStart).
		Count()

	// Get SMEs created last month
	newLastMonth, _ = facades.Orm().Query().Model(&models.Sme{}).
		Where("deleted_at IS NULL").
		Where("created_at >= ? AND created_at <= ?", lastMonthStart, lastMonthEnd).
		Count()

	// Get SMEs with registration number
	withRegistrationNumber, _ = facades.Orm().Query().Model(&models.Sme{}).
		Where("deleted_at IS NULL").
		Where("registration_number IS NOT NULL AND registration_number != ''").
		Count()

	// Calculate percentage with registration number
	var registrationPercentage float64
	if totalSmes > 0 {
		registrationPercentage = float64(withRegistrationNumber) / float64(totalSmes) * 100
	}

	// Get registration trend data (last 6 months)
	registrationTrend := s.getRegistrationTrend(6)

	// Get distribution by region
	regionDistribution := s.getDistributionByField("region")

	// Get distribution by business category
	categoryDistribution := s.getDistributionByField("business_category")

	return map[string]interface{}{
		"totalSmes":                 totalSmes,
		"newThisMonth":              newThisMonth,
		"newLastMonth":              newLastMonth,
		"hasRegistration":           withRegistrationNumber,
		"hasRegistrationPercentage": registrationPercentage,
		"registrationTrend":         registrationTrend,
		"byRegion":                  regionDistribution,
		"byCategory":                categoryDistribution,
	}, nil
}

// getRegistrationTrend returns SME registration counts for the last N months
func (s *SmeService) getRegistrationTrend(months int) []map[string]interface{} {
	now := time.Now()
	trend := make([]map[string]interface{}, months)

	// Calculate cumulative total up to 6 months ago
	sixMonthsAgo := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -months, 0)
	var cumulativeBase int64
	cumulativeBase, _ = facades.Orm().Query().Model(&models.Sme{}).
		Where("deleted_at IS NULL").
		Where("created_at < ?", sixMonthsAgo).
		Count()

	cumulative := cumulativeBase

	for i := months - 1; i >= 0; i-- {
		// Calculate month boundaries
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -i, 0)
		monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Second)

		// Get count for this month
		var count int64
		count, _ = facades.Orm().Query().Model(&models.Sme{}).
			Where("deleted_at IS NULL").
			Where("created_at >= ? AND created_at <= ?", monthStart, monthEnd).
			Count()

		cumulative += count

		// Format period as "Jan 2025"
		period := monthStart.Format("Jan 2006")

		trend[months-1-i] = map[string]interface{}{
			"period":     period,
			"count":      count,
			"cumulative": cumulative,
		}
	}

	return trend
}

// getDistributionByField returns distribution data for a given field (region or business_category)
func (s *SmeService) getDistributionByField(field string) []map[string]interface{} {
	// Get total count for percentage calculation
	var total int64
	total, _ = facades.Orm().Query().Model(&models.Sme{}).Where("deleted_at IS NULL").Count()

	if total == 0 {
		return []map[string]interface{}{}
	}

	// Query for distribution - using raw SQL for GROUP BY
	type DistributionResult struct {
		Label string
		Value int64
	}

	var results []DistributionResult

	// Build the query based on field type
	var query string
	if field == "region" || field == "business_category" {
		query = fmt.Sprintf(`
			SELECT COALESCE(%s, 'Unknown') as label, COUNT(*) as value
			FROM smes
			WHERE deleted_at IS NULL
			GROUP BY %s
			ORDER BY value DESC
		`, field, field)
	} else {
		return []map[string]interface{}{}
	}

	// Execute raw query
	err := facades.Orm().Query().Raw(query).Scan(&results)
	if err != nil {
		return []map[string]interface{}{}
	}

	// Convert to response format with percentages
	distribution := make([]map[string]interface{}, len(results))
	for i, r := range results {
		percentage := float64(r.Value) / float64(total) * 100
		distribution[i] = map[string]interface{}{
			"label":      r.Label,
			"value":      r.Value,
			"percentage": percentage, // Return as number, not string
		}
	}

	return distribution
}
