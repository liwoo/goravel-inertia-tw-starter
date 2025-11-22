package services

import (
	"encoding/json"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/models"
)

// BdspService implements business logic for bdsps using the builder pattern
type BdspService struct {
	contracts.CrudServiceContract // Embedded interface - automatically exposes all CRUD methods!
}

// NewBdspService creates a new Bdsp service using the builder pattern
func NewBdspService() *BdspService {
	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.Bdsp]("bdsps", "id").
		WithSearchFields("name", "postal_address", "physical_address", "registration_status", "partners", "product_types", "service_list"). // Fields that will be searchable via the search query parameter
		WithSortFields("id", "created_at", "updated_at", "name", "postal_address", "physical_address", "registration_status").              // Fields that can be used for sorting results
		WithFilterFields("name", "postal_address", "physical_address", "registration_status", "partners", "product_types", "service_list"). // Fields that can be filtered on
		WithValidationRules(map[string]interface{}{                                                                                         // Validation rules for create/update operations
			"registration_status": "string|max:255",
			"name":                "required|string|max:255",
			"postal_address":      "required|string|max:255",
			"physical_address":    "string|max:255",
			"partners":            "array",
			"product_types":       "array",
			"service_list":        "array",
		}).
		WithDefaultSort("created_at", "DESC").     // Default sorting when none specified
		WithScopeFiltering("bdsps", "created_by"). // Enable permission-based filtering
		WithBeforeUpdate(func(id uint, data map[string]interface{}) error {
			if partners, ok := data["partners"]; ok {
				if partnersSlice, ok := partners.([]interface{}); ok {
					bytes, err := json.Marshal(partnersSlice)
					if err == nil {
						data["partners_json"] = string(bytes)
					}
				}
				delete(data, "partners")
			}

			if productTypes, ok := data["product_types"]; ok {
				if productTypesSlice, ok := productTypes.([]interface{}); ok {
					bytes, err := json.Marshal(productTypesSlice)
					if err == nil {
						data["product_types_json"] = string(bytes)
					}
				}
				delete(data, "product_types")
			}

			if serviceList, ok := data["service_list"]; ok {
				if serviceListSlice, ok := serviceList.([]interface{}); ok {
					bytes, err := json.Marshal(serviceListSlice)
					if err == nil {
						data["service_list_json"] = string(bytes)
					}
				}
				delete(data, "service_list")
			}

			return nil
		}).
		Build() // Returns a fully configured CrudServiceContract

	bdspServiceInstance := &BdspService{
		CrudServiceContract: service, // Set the embedded interface
	}

	// Set the actual service instance for proper method resolution
	contracts.SetActualServiceHelper(service, bdspServiceInstance, "BdspService")

	return bdspServiceInstance
}

// GetFilterDefinitions returns filter definitions for the bdsps resource
// This method is OPTIONAL - only implement if you need custom filter UI components
// Uncomment and customize the implementation below if needed:
//
// func (s *BdspService) GetFilterDefinitions() []contracts.FilterDefinition {
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
// - GetByStatus(status string) ([]*models.Bdsp, error)
// - GetActive() ([]*models.Bdsp, error)
// - Custom business logic methods specific to bdsps
