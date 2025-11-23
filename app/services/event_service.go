package services

import (
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/models"
)

// EventService implements business logic for events using the builder pattern
type EventService struct {
	contracts.CrudServiceContract // Embedded interface - automatically exposes all CRUD methods!
}

// NewEventService creates a new event service using the builder pattern
func NewEventService() *EventService {
	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.Event]("events", "id").
		WithSearchFields("title", "description", "venue", "district", "notes"). // Fields that will be searchable via the search query parameter
		WithSortFields("id", "created_at", "updated_at", "date", "title"). // Fields that can be used for sorting results
		WithFilterFields("district", "date"). // Fields that can be filtered on
		WithValidationRules(map[string]interface{}{ // Validation rules for create/update operations
			"title":       "required|max_len:255",
			"description": "required",
			"venue":       "required|max_len:255",
			"district":    "required|max_len:100",
		}).
		WithDefaultSort("date", "DESC"). // Default sorting when none specified
		WithScopeFiltering("events", "created_by"). // Enable permission-based filtering

		Build() // Returns a fully configured CrudServiceContract

	eventServiceInstance := &EventService{
		CrudServiceContract: service, // Set the embedded interface
	}

	// Set the actual service instance for proper method resolution
	contracts.SetActualServiceHelper(service, eventServiceInstance, "EventService")

	return eventServiceInstance
}

// GetFilterDefinitions returns filter definitions for the events resource
// This method is OPTIONAL - only implement if you need custom filter UI components
// Uncomment and customize the implementation below if needed:
//
// func (s *EventService) GetFilterDefinitions() []contracts.FilterDefinition {
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
// - GetByStatus(status string) ([]*models.event, error)
// - GetActive() ([]*models.event, error)
// - Custom business logic methods specific to events
