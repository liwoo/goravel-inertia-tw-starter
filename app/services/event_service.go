package services

import (
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/http/requests"
	"smedi-sme-db/app/models"
	"sort"
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
		WithSortFields("id", "created_at", "updated_at", "date", "title").      // Fields that can be used for sorting results
		WithFilterFields("title", "venue", "district", "date").                 // Fields that can be filtered on
		WithValidationRules(map[string]interface{}{                             // Validation rules for create/update operations
			"title":       "required|max_len:255",
			"description": "required",
			"venue":       "required|max_len:255",
			"district":    "required|max_len:100",
		}).
		WithDefaultSort("date", "DESC").            // Default sorting when none specified
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
func (s *EventService) GetFilterDefinitions() []contracts.FilterDefinition {
	// Get all districts as strings for the enum
	districts := requests.GetAllDistricts()
	sort.Slice(districts, func(i, j int) bool {
		return districts[i] < districts[j]
	})
	districtStrings := make([]string, len(districts))
	for i, d := range districts {
		districtStrings[i] = string(d)
	}

	return []contracts.FilterDefinition{
		// Title - string search
		contracts.NewFilterDefinition(
			"title",
			"Event Title",
			contracts.FilterTypeString,
			nil,
		),
		// Venue - string search
		contracts.NewFilterDefinition(
			"venue",
			"Event Venue",
			contracts.FilterTypeString,
			nil,
		),
		// District - enum with all 28 Malawian districts
		contracts.NewFilterDefinition(
			"district",
			"District",
			contracts.FilterTypeEnum,
			&districtStrings,
		),
		// Event Date - date filter
		contracts.NewFilterDefinition(
			"date",
			"Event Date",
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
// - GetByStatus(status string) ([]*models.event, error)
// - GetActive() ([]*models.event, error)
// - Custom business logic methods specific to events
