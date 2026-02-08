package services

import (
	"books-database/app/contracts"
	"books-database/app/models"
)

// LenderService implements business logic for lenders using the builder pattern
type LenderService struct {
	contracts.CrudServiceContract // Embedded interface - automatically exposes all CRUD methods!
}

// NewLenderService creates a new Lender service using the builder pattern
func NewLenderService() *LenderService {
	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.Lender]("lenders", "id").
		WithSearchFields("name", "email", "phone", "address", "gender").                                 // Fields that will be searchable via the search query parameter
		WithSortFields("id", "created_at", "updated_at", "name", "email", "phone", "address", "gender"). // Fields that can be used for sorting results
		WithFilterFields("name", "email", "phone", "address", "gender").                                 // Fields that can be filtered on
		WithValidationRules(map[string]interface{}{                                                      // Validation rules for create/update operations
			"phone":   "string|max:255",
			"address": "string|max:255",
			"gender":  "string|max:255",
			"name":    "required|string|max:255",
			"email":   "required|string|max:255",
		}).
		WithDefaultSort("created_at", "DESC").       // Default sorting when none specified
		WithScopeFiltering("lenders", "created_by"). // Enable permission-based filtering
		WithTenantAwareness().
		Build() // Returns a fully configured CrudServiceContract

	lenderServiceInstance := &LenderService{
		CrudServiceContract: service, // Set the embedded interface
	}

	// Set the actual service instance for proper method resolution
	contracts.SetActualServiceHelper(service, lenderServiceInstance, "LenderService")

	return lenderServiceInstance
}

// GetFilterDefinitions returns filter definitions for the lenders resource
// This method is OPTIONAL - only implement if you need custom filter UI components
// Uncomment and customize the implementation below if needed:
//
// func (s *LenderService) GetFilterDefinitions() []contracts.FilterDefinition {
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
// - GetByStatus(status string) ([]*models.Lender, error)
// - GetActive() ([]*models.Lender, error)
// - Custom business logic methods specific to lenders
