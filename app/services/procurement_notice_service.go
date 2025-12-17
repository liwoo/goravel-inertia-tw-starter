package services

import (
	"encoding/json"
	"fmt"

	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/http/requests"
	"smedi-sme-db/app/models"

	"github.com/dromara/carbon/v2"
	"github.com/goravel/framework/facades"
)

// ProcurementNoticeService implements business logic for procurement notices using the builder pattern
type ProcurementNoticeService struct {
	contracts.CrudServiceContract // Embedded interface - automatically exposes all CRUD methods!
}

// NewProcurementNoticeService creates a new ProcurementNotice service using the builder pattern
func NewProcurementNoticeService() *ProcurementNoticeService {
	// Build the service with all required configurations
	// Helper function to ensure array fields are properly typed as []string for GORM serializer
	// The gorm:"serializer:json" tag handles JSON conversion, but we need []string type, not []interface{}
	ensureStringArray := func(data map[string]interface{}, field string) {
		if value, exists := data[field]; exists {
			switch v := value.(type) {
			case []interface{}:
				// Convert []interface{} to []string
				stringArr := make([]string, len(v))
				for i, item := range v {
					stringArr[i] = fmt.Sprintf("%v", item)
				}
				data[field] = stringArr
			case []string:
				// Already []string, keep as-is
			case nil:
				// Set to empty slice instead of nil
				data[field] = []string{}
			default:
				// Unknown type, set to empty array
				data[field] = []string{}
			}
		} else {
			// Field not provided, set to empty slice
			data[field] = []string{}
		}
	}

	// Helper function to convert array fields to JSON string for raw SQL updates
	// GORM's Update() method doesn't use model serializers, so we need to serialize arrays manually
	convertArrayToJSONString := func(data map[string]interface{}, field string) {
		if value, exists := data[field]; exists {
			switch v := value.(type) {
			case []interface{}:
				// Convert []interface{} to []string then JSON
				stringArr := make([]string, len(v))
				for i, item := range v {
					stringArr[i] = fmt.Sprintf("%v", item)
				}
				jsonBytes, _ := json.Marshal(stringArr)
				data[field] = string(jsonBytes)
			case []string:
				// Convert []string to JSON
				jsonBytes, _ := json.Marshal(v)
				data[field] = string(jsonBytes)
			case nil:
				// Set to empty JSON array
				data[field] = "[]"
			case string:
				// Already a JSON string, keep as-is
			default:
				// Unknown type, set to empty JSON array
				data[field] = "[]"
			}
		}
		// Don't set default for update - only update fields that are provided
	}

	service := contracts.NewServiceBuilder[models.ProcurementNotice]("procurement_notices", "id").
		WithSearchFields("procured_by", "procurement_type", "ref_no", "organization", "details", "application_details").                      // Searchable fields
		WithSortFields("id", "created_at", "updated_at", "open_date", "close_date", "procured_by", "organization").                           // Sortable fields
		WithFilterFields("organization", "market_approach", "invitation", "is_published", "qualifying_districts", "open_date", "close_date"). // Filterable fields
		WithValidationRules(map[string]interface{}{                                                                                           // Validation rules for create/update operations
			"procured_by":              "required|max_len:255",
			"procurement_type":         "required|max_len:255",
			"market_approach":          "required|max_len:50",
			"invitation":               "required|max_len:50",
			"organization":             "required|max_len:255",
			"details":                  "required",
			"application_details":      "required|max_len:500",
			"minimum_qualifying_score": "required",
			"is_published":             "boolean",
		}).
		WithDefaultSort("created_at", "DESC").                   // Default sorting when none specified
		WithScopeFiltering("procurement_notices", "created_by"). // Enable permission-based filtering
		WithBeforeCreate(func(data map[string]interface{}) error {
			// Ensure array fields are properly typed as []string for GORM serializer
			ensureStringArray(data, "partners")
			ensureStringArray(data, "qualifying_districts")
			ensureStringArray(data, "classification")
			ensureStringArray(data, "interested_smes")

			// Generate Reference Number if not provided
			if _, exists := data["ref_no"]; !exists || data["ref_no"] == "" {
				today := carbon.Now().ToDateString()     // YYYY-MM-DD
				datePrefix := carbon.Now().Format("Ymd") // YYYYMMDD

				var count int64
				// Assuming 'db' (a *gorm.DB instance) is accessible in this scope.
				// This query counts existing notices created today to determine the next index.
				var procurementNotice models.ProcurementNotice
				count, err := facades.Orm().Query().
					Model(&procurementNotice).
					Where("DATE(created_at) = ?", today).
					Count()
				if err != nil {
					return fmt.Errorf("failed to get daily procurement notice count: %w", err)
				}

				data["ref_no"] = fmt.Sprintf("PN-%s-%d", datePrefix, count+1)
			}

			// Default IsPublished to false
			if _, exists := data["is_published"]; !exists {
				data["is_published"] = false
			}

			return nil
		}).
		WithBeforeUpdate(func(id uint, data map[string]interface{}) error {
			// Convert array fields to JSON strings for raw SQL updates
			// GORM's Update() doesn't use model serializers
			convertArrayToJSONString(data, "partners")
			convertArrayToJSONString(data, "qualifying_districts")
			convertArrayToJSONString(data, "classification")
			convertArrayToJSONString(data, "interested_smes")
			return nil
		}).
		Build() // Returns a fully configured CrudServiceContract

	procurementNoticeServiceInstance := &ProcurementNoticeService{
		CrudServiceContract: service, // Set the embedded interface
	}

	// Set the actual service instance for proper method resolution
	contracts.SetActualServiceHelper(service, procurementNoticeServiceInstance, "ProcurementNoticeService")

	return procurementNoticeServiceInstance
}

// GetFilterDefinitions returns filter definitions for the procurement notices resource
func (s *ProcurementNoticeService) GetFilterDefinitions() []contracts.FilterDefinition {
	// Get all districts as strings for the enum
	districts := requests.GetAllDistricts()
	districtStrings := make([]string, len(districts))
	for i, d := range districts {
		districtStrings[i] = string(d)
	}

	// Market approach options
	marketApproaches := []string{
		"Open",
		"Restricted",
		"Direct",
	}

	// Invitation options
	invitations := []string{
		"Open",
		"Selective",
		"Limited",
	}

	return []contracts.FilterDefinition{
		// Organization - string search
		contracts.NewFilterDefinition(
			"organization",
			"Organization",
			contracts.FilterTypeString,
			nil,
		),
		// Procured By - string search
		contracts.NewFilterDefinition(
			"procured_by",
			"Procured By",
			contracts.FilterTypeString,
			nil,
		),
		// Market Approach - enum
		contracts.NewFilterDefinition(
			"market_approach",
			"Market Approach",
			contracts.FilterTypeEnum,
			&marketApproaches,
		),
		// Invitation - enum
		contracts.NewFilterDefinition(
			"invitation",
			"Invitation Type",
			contracts.FilterTypeEnum,
			&invitations,
		),
		// Published Status - boolean
		contracts.NewFilterDefinition(
			"is_published",
			"Published",
			contracts.FilterTypeBoolean,
			nil,
		),
		// Qualifying Districts - enum with all 28 Malawian districts
		contracts.NewFilterDefinition(
			"qualifying_districts",
			"Qualifying Districts",
			contracts.FilterTypeEnum,
			&districtStrings,
		),
		// Open Date - date filter
		contracts.NewFilterDefinition(
			"open_date",
			"Opening Date",
			contracts.FilterTypeDate,
			nil,
		),
		// Close Date - date filter
		contracts.NewFilterDefinition(
			"close_date",
			"Closing Date",
			contracts.FilterTypeDate,
			nil,
		),
		// Created Date - datetime filter
		contracts.NewFilterDefinition(
			"created_at",
			"Date Created",
			contracts.FilterTypeDateTime,
			nil,
		),
	}
}

// Add domain-specific methods below this line
// Examples:
// - GetByStatus(status string) ([]*models.ProcurementNotice, error)
// - GetActive() ([]*models.ProcurementNotice, error)
// - Custom business logic methods specific to procurementnotices
