package services

import (
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/models"
)

// ProcurementNoticeService implements business logic for procurement notices using the builder pattern
type ProcurementNoticeService struct {
	contracts.CrudServiceContract // Embedded interface - automatically exposes all CRUD methods!
}

// NewProcurementNoticeService creates a new ProcurementNotice service using the builder pattern
func NewProcurementNoticeService() *ProcurementNoticeService {
	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.ProcurementNotice]("procurement_notices", "id").
		WithSearchFields("procured_by", "procurement_type", "ref_no", "organization", "details", "application_details"). // Searchable fields
		WithSortFields("id", "created_at", "updated_at", "open_date", "close_date", "procured_by", "organization"). // Sortable fields
		WithFilterFields("market_approach", "invitation", "is_published", "qualifying_districts", "open_date", "close_date"). // Filterable fields
		WithValidationRules(map[string]interface{}{ // Validation rules for create/update operations
			"procured_by":              "required|max_len:255",
			"procurement_type":         "required|max_len:255",
			"market_approach":          "required|max_len:50",
			"invitation":               "required|max_len:50",
			"ref_no":                   "required|max_len:255",
			"organization":             "required|max_len:255",
			"details":                  "required",
			"application_details":      "required|max_len:500",
			"minimum_qualifying_score": "required",
			"is_published":             "boolean",
		}).
		WithDefaultSort("created_at", "DESC"). // Default sorting when none specified
		WithScopeFiltering("procurement_notices", "created_by"). // Enable permission-based filtering

		Build() // Returns a fully configured CrudServiceContract

	procurementNoticeServiceInstance := &ProcurementNoticeService{
		CrudServiceContract: service, // Set the embedded interface
	}

	// Set the actual service instance for proper method resolution
	contracts.SetActualServiceHelper(service, procurementNoticeServiceInstance, "ProcurementNoticeService")

	return procurementNoticeServiceInstance
}

// GetFilterDefinitions returns filter definitions for the procurementnotices resource
// This method is OPTIONAL - only implement if you need custom filter UI components
// Uncomment and customize the implementation below if needed:
//
// func (s *Procurement_noticeService) GetFilterDefinitions() []contracts.FilterDefinition {
// 	return []contracts.FilterDefinition{
// 		// String field example
// 		contracts.NewFilterDefinition(
// 			"field_name",           // Field name in database
// 			"Display Name",         // Human-readable label
// 			contracts.FilterTypeString,
// 			nil,                    // nil = use all string operators (equals, contains, starts_with, etc.)
// 		),
//
// 		// Enum field example (dropdown)
// 		contracts.NewFilterDefinition(
// 			"status",
// 			"Status",
// 			contracts.FilterTypeEnum,
// 			&[]string{"ACTIVE", "INACTIVE", "PENDING"}, // Available options
// 		),
//
// 		// Number field example
// 		contracts.NewFilterDefinition(
// 			"price",
// 			"Price",
// 			contracts.FilterTypeNumber,
// 			nil,                    // nil = use all number operators (equals, greater_than, less_than, etc.)
// 		),
//
// 		// Date field example
// 		contracts.NewFilterDefinition(
// 			"created_at",
// 			"Created Date",
// 			contracts.FilterTypeDate,
// 			nil,                    // nil = use all date operators (equals, before, after, between, etc.)
// 		),
//
// 		// Boolean field example
// 		contracts.NewFilterDefinition(
// 			"is_active",
// 			"Active Status",
// 			contracts.FilterTypeBoolean,
// 			nil,
// 		),
// 	}
// }

// Add domain-specific methods below this line
// Examples:
// - GetByStatus(status string) ([]*models.ProcurementNotice, error)
// - GetActive() ([]*models.ProcurementNotice, error)
// - Custom business logic methods specific to procurementnotices
