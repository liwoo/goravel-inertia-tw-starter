package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/database/orm"
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
		WithSearchFields("usme_number", "name", "registration_number", "tax_identification_number", "business_category", "sector", "contact_email", "contact_phone", "region", "district", "classification"). // Fields that will be searchable via the search query parameter
		WithSortFields("id", "created_at", "updated_at", "usme_number", "name", "operational_start_date", "business_category", "sector", "region", "district", "formalisation_score", "classification").       // Fields that can be used for sorting results
		WithFilterFields("business_category", "sector", "region", "district", "created_by", "is_active", "formalisation_score", "classification").                                                             // Fields that can be filtered on
		WithValidationRules(map[string]interface{}{                                                                                                                                          // Validation rules for create/update operations
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
		WithCustomQuery(func(query orm.Query) orm.Query {
			// Join business_formalisation table for sorting/filtering by formalisation_score
			return query.Join("LEFT JOIN business_formalisation ON business_formalisation.sme_id = smes.id AND business_formalisation.deleted_at IS NULL")
		}).
		WithBeforeCreate(func(data map[string]interface{}) error { // Handle JSON array fields and UBI generation
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
		WithAfterCreate(func(model *models.Sme) error {
			// Calculate classification after creating an SME
			facades.Log().Info("AfterCreate hook triggered for SME classification", map[string]interface{}{
				"sme_id": model.ID,
			})
			svc := &SmeService{}
			classification, err := svc.CalculateClassification(model.ID)
			if err != nil {
				facades.Log().Warning("Failed to calculate classification after create", map[string]interface{}{
					"sme_id": model.ID,
					"error":  err.Error(),
				})
			} else {
				facades.Log().Info("Classification calculated after create", map[string]interface{}{
					"sme_id":         model.ID,
					"classification": classification,
				})
			}
			return nil
		}).
		WithAfterUpdate(func(model *models.Sme) error {
			// Recalculate classification after updating an SME
			facades.Log().Info("AfterUpdate hook triggered for SME classification", map[string]interface{}{
				"sme_id": model.ID,
			})
			svc := &SmeService{}
			classification, err := svc.CalculateClassification(model.ID)
			if err != nil {
				facades.Log().Warning("Failed to calculate classification after update", map[string]interface{}{
					"sme_id": model.ID,
					"error":  err.Error(),
				})
			} else {
				facades.Log().Info("Classification calculated after update", map[string]interface{}{
					"sme_id":         model.ID,
					"classification": classification,
				})
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
	mapping["isActive"] = "is_active"
	// Map formalisationScore to the joined table column with COALESCE to handle NULLs
	// SMEs without formalisation records will be treated as having score 0
	mapping["formalisationScore"] = "COALESCE(business_formalisation.formalisation_score, 0)"
	mapping["formalisation_score"] = "COALESCE(business_formalisation.formalisation_score, 0)"
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
		// Formalisation Score - number filter (0-100)
		contracts.NewFilterDefinition(
			"formalisation_score",
			"Formalisation Score",
			contracts.FilterTypeNumber,
			nil, // Use all number operators (equals, greater_than, less_than, between, etc.)
		),
		// Classification - enum filter with all classification categories
		contracts.NewFilterDefinition(
			"classification",
			"Classification",
			contracts.FilterTypeEnum,
			&[]string{
				models.ClassificationMicro,
				models.ClassificationSmall,
				models.ClassificationMedium,
				models.ClassificationUnclassified,
			},
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

	// Get total SMEs (GORM automatically handles deleted_at IS NULL for soft delete models)
	totalSmes, _ = facades.Orm().Query().Model(&models.Sme{}).Count()

	// Get SMEs created this month
	newThisMonth, _ = facades.Orm().Query().Model(&models.Sme{}).
		Where("created_at >= ?", currentMonthStart).
		Count()

	// Get SMEs created last month
	newLastMonth, _ = facades.Orm().Query().Model(&models.Sme{}).
		Where("created_at >= ? AND created_at <= ?", lastMonthStart, lastMonthEnd).
		Count()

	// Get SMEs with registration number
	withRegistrationNumber, _ = facades.Orm().Query().Model(&models.Sme{}).
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
	// GORM automatically handles deleted_at IS NULL for soft delete models
	sixMonthsAgo := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -months, 0)
	var cumulativeBase int64
	cumulativeBase, _ = facades.Orm().Query().Model(&models.Sme{}).
		Where("created_at < ?", sixMonthsAgo).
		Count()

	cumulative := cumulativeBase

	for i := months - 1; i >= 0; i-- {
		// Calculate month boundaries
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -i, 0)
		monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Second)

		// Get count for this month
		// GORM automatically handles deleted_at IS NULL for soft delete models
		var count int64
		count, _ = facades.Orm().Query().Model(&models.Sme{}).
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
	// GORM automatically handles deleted_at IS NULL for soft delete models
	var total int64
	total, _ = facades.Orm().Query().Model(&models.Sme{}).Count()

	if total == 0 {
		return []map[string]interface{}{}
	}

	// Query for distribution - using raw SQL for GROUP BY
	// Note: Raw SQL requires manual deleted_at IS NULL since GORM doesn't auto-apply it
	type DistributionResult struct {
		Label string
		Value int64
	}

	var results []DistributionResult

	// Build the query based on field type
	var query string
	if field == "region" || field == "business_category" || field == "sector" || field == "district" || field == "classification" {
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

// GetSmeByUserId retrieves the SME associated with the given user ID
// It checks both PrimaryBusinessOwner and AdditionalBusinessMember tables by email
func (s *SmeService) GetSmeByUserEmail(email string) (*models.Sme, error) {
	var sme models.Sme

	err := facades.Orm().Query().
		Model(&models.Sme{}).
		Join("INNER JOIN primary_business_owner ON primary_business_owner.sme_id = smes.id").
		Where("primary_business_owner.email = ?", email).
		Where("primary_business_owner.deleted_at IS NULL").
		First(&sme)

	if err == nil && sme.ID != 0 {
		return &sme, nil
	}

	err = facades.Orm().Query().
		Model(&models.Sme{}).
		Join("INNER JOIN additional_business_members ON additional_business_members.sme_id = smes.id").
		Where("additional_business_members.email = ?", email).
		Where("additional_business_members.deleted_at IS NULL").
		First(&sme)

	if err != nil {
		return nil, err
	}

	if sme.ID == 0 {
		return nil, nil
	}

	return &sme, nil
}

// GetDistributionBySector returns SME distribution by sector field
func (s *SmeService) GetDistributionBySector() []map[string]interface{} {
	return s.getDistributionByField("sector")
}

// GetDistributionByDistrict returns SME distribution by district field
func (s *SmeService) GetDistributionByDistrict() []map[string]interface{} {
	return s.getDistributionByField("district")
}

// BulkUpdateStatus updates the is_active status for multiple SMEs
func (s *SmeService) BulkUpdateStatus(ids []uint, isActive bool, updatedBy *int) (int64, error) {
	if len(ids) == 0 {
		return 0, errors.New("no IDs provided")
	}

	// Convert []uint to []any for WhereIn
	anyIds := make([]any, len(ids))
	for i, id := range ids {
		anyIds[i] = id
	}

	// Build update data
	updateData := map[string]interface{}{
		"is_active": isActive,
	}
	if updatedBy != nil {
		updateData["updated_by"] = *updatedBy
	}

	// Perform bulk update
	result, err := facades.Orm().Query().Model(&models.Sme{}).
		WhereIn("id", anyIds).
		Update(updateData)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected, nil
}

// BulkDeactivate deactivates multiple SMEs by their IDs
func (s *SmeService) BulkDeactivate(ids []uint, updatedBy *int) (int64, error) {
	return s.BulkUpdateStatus(ids, false, updatedBy)
}

// BulkActivate activates multiple SMEs by their IDs
func (s *SmeService) BulkActivate(ids []uint, updatedBy *int) (int64, error) {
	return s.BulkUpdateStatus(ids, true, updatedBy)
}

// GetDistributionByGender returns SME distribution by primary business owner gender
func (s *SmeService) GetDistributionByGender() []map[string]interface{} {
	// Get total count for percentage calculation
	// Count SMEs that have primary business owners with gender data
	var total int64
	query := `
		SELECT COUNT(DISTINCT s.id)
		FROM smes s
		INNER JOIN primary_business_owner p ON p.sme_id = s.id
		WHERE s.deleted_at IS NULL AND p.deleted_at IS NULL AND p.gender IS NOT NULL AND p.gender != ''
	`
	err := facades.Orm().Query().Raw(query).Scan(&total)
	if err != nil || total == 0 {
		return []map[string]interface{}{}
	}

	// Query for distribution by gender
	type DistributionResult struct {
		Label string
		Value int64
	}

	var results []DistributionResult
	genderQuery := `
		SELECT COALESCE(p.gender, 'Unknown') as label, COUNT(*) as value
		FROM smes s
		INNER JOIN primary_business_owner p ON p.sme_id = s.id
		WHERE s.deleted_at IS NULL AND p.deleted_at IS NULL
		GROUP BY p.gender
		ORDER BY value DESC
	`

	err = facades.Orm().Query().Raw(genderQuery).Scan(&results)
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
			"percentage": percentage,
		}
	}

	return distribution
}

// GetDistributionByYouth returns SME distribution by youth status (18-35 years = Youth)
func (s *SmeService) GetDistributionByYouth() []map[string]interface{} {
	// Get total count for percentage calculation
	// Count SMEs that have primary business owners with valid date_of_birth
	var total int64
	query := `
		SELECT COUNT(DISTINCT s.id)
		FROM smes s
		INNER JOIN primary_business_owner p ON p.sme_id = s.id
		WHERE s.deleted_at IS NULL AND p.deleted_at IS NULL AND p.date_of_birth IS NOT NULL
	`
	err := facades.Orm().Query().Raw(query).Scan(&total)
	if err != nil || total == 0 {
		return []map[string]interface{}{}
	}

	// Query for youth vs non-youth distribution
	// Youth is defined as 18-35 years old
	// Using PostgreSQL syntax: EXTRACT(YEAR FROM AGE(CURRENT_DATE, date_of_birth))
	type YouthResult struct {
		Label string
		Value int64
	}

	var results []YouthResult
	youthQuery := `
		SELECT
			CASE
				WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.date_of_birth)) BETWEEN 18 AND 35 THEN 'Youth'
				ELSE 'Non-Youth'
			END as label,
			COUNT(*) as value
		FROM smes s
		INNER JOIN primary_business_owner p ON p.sme_id = s.id
		WHERE s.deleted_at IS NULL AND p.deleted_at IS NULL AND p.date_of_birth IS NOT NULL
		GROUP BY 1
		ORDER BY label
	`

	err = facades.Orm().Query().Raw(youthQuery).Scan(&results)
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
			"percentage": percentage,
		}
	}

	return distribution
}

// GetAgeGenderDistribution returns SME distribution by age groups and gender for population pyramid
func (s *SmeService) GetAgeGenderDistribution() []map[string]interface{} {
	// Get total count for percentage calculation
	var total int64
	query := `
		SELECT COUNT(DISTINCT s.id)
		FROM smes s
		INNER JOIN primary_business_owner p ON p.sme_id = s.id
		WHERE s.deleted_at IS NULL AND p.deleted_at IS NULL
		AND p.date_of_birth IS NOT NULL AND p.gender IS NOT NULL AND p.gender != ''
	`
	err := facades.Orm().Query().Raw(query).Scan(&total)
	if err != nil || total == 0 {
		return []map[string]interface{}{}
	}

	// Query for age-gender distribution
	// Using PostgreSQL syntax: EXTRACT(YEAR FROM AGE(CURRENT_DATE, date_of_birth))
	type AgeGenderResult struct {
		AgeGroup string
		Gender   string
		Value    int64
	}

	var results []AgeGenderResult
	ageGenderQuery := `
		SELECT
			CASE
				WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.date_of_birth)) BETWEEN 18 AND 25 THEN '18-25'
				WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.date_of_birth)) BETWEEN 26 AND 35 THEN '26-35'
				WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.date_of_birth)) BETWEEN 36 AND 45 THEN '36-45'
				WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.date_of_birth)) BETWEEN 46 AND 55 THEN '46-55'
				WHEN EXTRACT(YEAR FROM AGE(CURRENT_DATE, p.date_of_birth)) BETWEEN 56 AND 65 THEN '56-65'
				ELSE '65+'
			END as age_group,
			p.gender as gender,
			COUNT(*) as value
		FROM smes s
		INNER JOIN primary_business_owner p ON p.sme_id = s.id
		WHERE s.deleted_at IS NULL AND p.deleted_at IS NULL
		AND p.date_of_birth IS NOT NULL AND p.gender IS NOT NULL AND p.gender != ''
		GROUP BY 1, p.gender
		ORDER BY 1
	`

	err = facades.Orm().Query().Raw(ageGenderQuery).Scan(&results)
	if err != nil {
		return []map[string]interface{}{}
	}

	// Define age groups in order (reversed for pyramid - oldest at top)
	ageGroups := []string{"65+", "56-65", "46-55", "36-45", "26-35", "18-25"}

	// Build a map for easy lookup
	// Keys are uppercase to match database values (MALE, FEMALE)
	dataMap := make(map[string]map[string]int64)
	for _, ag := range ageGroups {
		dataMap[ag] = map[string]int64{"MALE": 0, "FEMALE": 0}
	}

	for _, r := range results {
		if _, exists := dataMap[r.AgeGroup]; exists {
			// Normalize gender to uppercase for consistent lookup
			gender := strings.ToUpper(r.Gender)
			dataMap[r.AgeGroup][gender] = r.Value
		}
	}

	// Convert to response format
	distribution := make([]map[string]interface{}, len(ageGroups))
	for i, ag := range ageGroups {
		maleCount := dataMap[ag]["MALE"]
		femaleCount := dataMap[ag]["FEMALE"]
		malePercentage := float64(maleCount) / float64(total) * 100
		femalePercentage := float64(femaleCount) / float64(total) * 100

		distribution[i] = map[string]interface{}{
			"ageGroup":         ag,
			"male":             maleCount,
			"female":           femaleCount,
			"malePercentage":   malePercentage,
			"femalePercentage": femalePercentage,
		}
	}

	return distribution
}

// ============================================================================
// Formalisation Score Calculation Functions
// ============================================================================

// CalculateFormalisationScore calculates and updates the formalisation score for a given SME
// The score is calculated based on:
// - Formalisation Checkboxes (70 points max): 7 boolean fields, 10 points each
// - Team Structure (20 points max): primary owner, team members, employees, team size
// - Financial Data (10 points max): annual turnover and estimated assets
// Returns the calculated score (0-100)
func (s *SmeService) CalculateFormalisationScore(smeID uint) (int, error) {
	// Load SME with all related data
	var sme models.Sme
	err := facades.Orm().Query().
		With("BusinessFormalisation").
		With("BusinessEmployeeSummary").
		With("PrimaryBusinessOwner").
		With("AdditionalBusinessMembers").
		Where("id = ?", smeID).
		First(&sme)

	if err != nil {
		return 0, fmt.Errorf("failed to load SME: %w", err)
	}

	if sme.ID == 0 {
		return 0, errors.New("SME not found")
	}

	// Calculate the score
	score := s.calculateScore(&sme)

	// Update the FormalisationScore in BusinessFormalisation table
	if sme.BusinessFormalisation != nil {
		_, err = facades.Orm().Query().
			Model(&models.BusinessFormalisation{}).
			Where("id = ?", sme.BusinessFormalisation.ID).
			Update(map[string]interface{}{
				"formalisation_score": score,
			})
		if err != nil {
			return score, fmt.Errorf("failed to update formalisation score: %w", err)
		}
	} else {
		// If no BusinessFormalisation record exists, create one with just the score
		newFormalisation := &models.BusinessFormalisation{
			SmeID:              int(smeID),
			FormalisationScore: score,
		}
		err = facades.Orm().Query().Create(newFormalisation)
		if err != nil {
			return score, fmt.Errorf("failed to create business formalisation record: %w", err)
		}
	}

	return score, nil
}

// calculateScore performs the actual score calculation based on SME data
func (s *SmeService) calculateScore(sme *models.Sme) int {
	score := 0

	// ========================================
	// Formalisation Checkboxes (70 points max)
	// ========================================
	if sme.BusinessFormalisation != nil {
		bf := sme.BusinessFormalisation

		// has_bank_account: 10 points
		if bf.HasBankAccount {
			score += 10
		}

		// has_tax_clarification: 10 points
		if bf.HasTaxClarification {
			score += 10
		}

		// is_registered_for_vat: 10 points
		if bf.IsRegisteredForVat {
			score += 10
		}

		// is_member_of_association: 10 points
		if bf.IsMemberOfAssociation {
			score += 10
		}

		// is_affiliated: 10 points
		if bf.IsAffiliated {
			score += 10
		}

		// has_export_license: 10 points
		if bf.HasExportLicense {
			score += 10
		}

		// has_accessed_bds: 10 points
		if bf.HasAccessedBds {
			score += 10
		}
	}

	// ========================================
	// Team Structure (20 points max)
	// ========================================

	// Has primary business owner: 5 points
	if sme.PrimaryBusinessOwner != nil && sme.PrimaryBusinessOwner.ID != 0 {
		score += 5
	}

	// Has additional team members (1+): 5 points
	if len(sme.AdditionalBusinessMembers) > 0 {
		score += 5
	}

	// Calculate total team size for team structure scoring
	totalTeamSize := 0

	// Count primary owner
	if sme.PrimaryBusinessOwner != nil && sme.PrimaryBusinessOwner.ID != 0 {
		totalTeamSize += 1
	}

	// Count additional business members
	totalTeamSize += len(sme.AdditionalBusinessMembers)

	// Has full-time employees (1+): 5 points
	fullTimeEmployees := 0
	if sme.BusinessEmployeeSummary != nil {
		bes := sme.BusinessEmployeeSummary
		fullTimeEmployees = bes.FullTimeMales + bes.FullTimeFemales

		// Add all employees to team size
		totalTeamSize += bes.FullTimeMales + bes.FullTimeFemales
		totalTeamSize += bes.PartTimeMales + bes.PartTimeFemales
		totalTeamSize += bes.InternMales + bes.InternFemales
	}

	if fullTimeEmployees > 0 {
		score += 5
	}

	// Team size > 5 people: 5 points
	if totalTeamSize > 5 {
		score += 5
	}

	// ========================================
	// Financial Data (10 points max)
	// ========================================
	if sme.BusinessFormalisation != nil {
		bf := sme.BusinessFormalisation

		// annual_turnover > 0: 5 points
		if bf.AnnualTurnover > 0 {
			score += 5
		}

		// estimated_value_of_assets > 0: 5 points
		if bf.EstimatedValueOfAssets > 0 {
			score += 5
		}
	}

	// Ensure score doesn't exceed 100
	if score > 100 {
		score = 100
	}

	return score
}

// RecalculateAllFormalisationScores recalculates formalisation scores for all SMEs
// This is useful for batch processing when the scoring algorithm changes
// Returns the number of SMEs processed and any error encountered
func (s *SmeService) RecalculateAllFormalisationScores() (int, error) {
	// Get all SME IDs
	var smes []models.Sme
	err := facades.Orm().Query().
		Model(&models.Sme{}).
		Select("id").
		Find(&smes)

	if err != nil {
		return 0, fmt.Errorf("failed to fetch SME IDs: %w", err)
	}

	processedCount := 0
	var lastError error

	for _, sme := range smes {
		_, err := s.CalculateFormalisationScore(sme.ID)
		if err != nil {
			// Log the error but continue processing
			facades.Log().Warning("Failed to calculate formalisation score", map[string]interface{}{
				"sme_id": sme.ID,
				"error":  err.Error(),
			})
			lastError = err
		} else {
			processedCount++
		}
	}

	if lastError != nil && processedCount == 0 {
		return 0, fmt.Errorf("failed to process any SMEs: %w", lastError)
	}

	return processedCount, nil
}

// ============================================================================
// SME Classification Functions (Malawi MSME Policy)
// ============================================================================

// CalculateClassification calculates and updates the classification for a given SME
// Classification is based on:
// - Employment criteria (must be met)
// - At least one of: Annual Turnover OR Maximum Assets criteria
//
// Classification Rules:
// | Size   | Employees | Annual Turnover (MWK)          | Max Assets (MWK) |
// |--------|-----------|--------------------------------|------------------|
// | Micro  | 1-4       | Up to 5,000,000                | 1,000,000        |
// | Small  | 5-20      | 5,000,001 - 50,000,000         | 20,000,000       |
// | Medium | 21-99     | 50,000,001 - 500,000,000       | 250,000,000      |
//
// Returns the classification string and any error encountered
func (s *SmeService) CalculateClassification(smeID uint) (string, error) {
	// Load SME with BusinessFormalisation and BusinessEmployeeSummary relationships
	var sme models.Sme
	err := facades.Orm().Query().
		With("BusinessFormalisation").
		With("BusinessEmployeeSummary").
		Where("id = ?", smeID).
		First(&sme)

	if err != nil {
		return "", fmt.Errorf("failed to load SME: %w", err)
	}

	if sme.ID == 0 {
		return "", errors.New("SME not found")
	}

	// Calculate total employees from BusinessEmployeeSummary
	totalEmployees := 0
	if sme.BusinessEmployeeSummary != nil {
		bes := sme.BusinessEmployeeSummary
		totalEmployees = bes.FullTimeMales + bes.FullTimeFemales +
			bes.PartTimeMales + bes.PartTimeFemales +
			bes.InternMales + bes.InternFemales
	}

	// Get turnover and assets from BusinessFormalisation
	var turnover, assets float64
	if sme.BusinessFormalisation != nil {
		turnover = sme.BusinessFormalisation.AnnualTurnover
		assets = sme.BusinessFormalisation.EstimatedValueOfAssets
	}

	// Determine classification
	classification := s.DetermineClassification(totalEmployees, turnover, assets)

	// Update the SME record with the classification
	_, err = facades.Orm().Query().
		Model(&models.Sme{}).
		Where("id = ?", smeID).
		Update(map[string]interface{}{
			"classification": classification,
		})
	if err != nil {
		return classification, fmt.Errorf("failed to update classification: %w", err)
	}

	return classification, nil
}

// DetermineClassification applies the classification rules based on Malawi MSME Policy
// An SME meets a classification if it satisfies:
// - Employment criteria AND
// - At least one of (Turnover OR Assets) criteria
// This method is exported to allow unit testing of the classification logic
func (s *SmeService) DetermineClassification(employees int, turnover, assets float64) string {
	// Check Medium classification first (highest)
	// Employees: 21-99
	// Turnover: Above 50,000,000 - 500,000,000
	// Assets: Up to 250,000,000
	if employees >= models.MediumEmployeeMin && employees <= models.MediumEmployeeMax {
		// Check if turnover OR assets criteria is met
		turnoverMet := turnover > models.MediumTurnoverMin && turnover <= models.MediumTurnoverMax
		assetsMet := assets <= models.MediumAssetsMax && assets > 0
		if turnoverMet || assetsMet {
			return models.ClassificationMedium
		}
	}

	// Check Small classification
	// Employees: 5-20
	// Turnover: Above 5,000,000 - 50,000,000
	// Assets: Up to 20,000,000
	if employees >= models.SmallEmployeeMin && employees <= models.SmallEmployeeMax {
		// Check if turnover OR assets criteria is met
		turnoverMet := turnover > models.SmallTurnoverMin && turnover <= models.SmallTurnoverMax
		assetsMet := assets <= models.SmallAssetsMax && assets > 0
		if turnoverMet || assetsMet {
			return models.ClassificationSmall
		}
	}

	// Check Micro classification
	// Employees: 1-4
	// Turnover: Up to 5,000,000
	// Assets: Up to 1,000,000
	if employees >= models.MicroEmployeeMin && employees <= models.MicroEmployeeMax {
		// For Micro, we require at least some financial data to classify
		// Either turnover OR assets criteria can be met (OR logic like Small/Medium)
		// But at least one must have a positive value to avoid classifying with no data
		turnoverMet := turnover > 0 && turnover <= models.MicroTurnoverMax
		assetsMet := assets > 0 && assets <= models.MicroAssetsMax
		if turnoverMet || assetsMet {
			return models.ClassificationMicro
		}
	}

	// Default to Unclassified if no criteria are met
	return models.ClassificationUnclassified
}

// RecalculateAllClassifications recalculates classifications for all SMEs
// This is useful for batch processing when the classification algorithm changes
// Returns the number of SMEs processed and any error encountered
func (s *SmeService) RecalculateAllClassifications() (int, error) {
	// Get all SME IDs
	var smes []models.Sme
	err := facades.Orm().Query().
		Model(&models.Sme{}).
		Select("id").
		Find(&smes)

	if err != nil {
		return 0, fmt.Errorf("failed to fetch SME IDs: %w", err)
	}

	processedCount := 0
	var lastError error

	for _, sme := range smes {
		_, err := s.CalculateClassification(sme.ID)
		if err != nil {
			// Log the error but continue processing
			facades.Log().Warning("Failed to calculate classification", map[string]interface{}{
				"sme_id": sme.ID,
				"error":  err.Error(),
			})
			lastError = err
		} else {
			processedCount++
		}
	}

	if lastError != nil && processedCount == 0 {
		return 0, fmt.Errorf("failed to process any SMEs: %w", lastError)
	}

	return processedCount, nil
}

// GetDistributionByClassification returns SME distribution by classification
func (s *SmeService) GetDistributionByClassification() []map[string]interface{} {
	return s.getDistributionByField("classification")
}
