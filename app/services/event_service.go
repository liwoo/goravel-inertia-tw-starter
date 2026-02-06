package services

import (
	"fmt"
	"sort"

	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/http/requests"
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

// AttendEvent adds an SME to the event's attending list
func (s *EventService) AttendEvent(eventId, smeId uint) error {
	// Fetch the event by ID
	eventInterface, err := s.GetByID(eventId)
	if err != nil {
		return fmt.Errorf("event not found: %w", err)
	}
	if eventInterface == nil {
		return fmt.Errorf("event not found")
	}
	event, ok := eventInterface.(*models.Event)
	if !ok {
		return fmt.Errorf("invalid event type")
	}

	// Check if SME is already attending
	for _, id := range event.AttendingSmes {
		if id == int(smeId) {
			// Already attending, return without error
			return nil
		}
	}

	// Add SME to the attending list
	event.AttendingSmes = append(event.AttendingSmes, int(smeId))

	// Save the event using raw query to update the JSON field
	_, err = facades.Orm().Query().Model(&models.Event{}).
		Where("id = ?", eventId).
		Update(map[string]interface{}{
			"attending_smes": event.AttendingSmes,
		})
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

// UnattendEvent removes an SME from the event's attending list
func (s *EventService) UnattendEvent(eventId, smeId uint) error {
	// Fetch the event by ID
	eventInterface, err := s.GetByID(eventId)
	if err != nil {
		return fmt.Errorf("event not found: %w", err)
	}
	if eventInterface == nil {
		return fmt.Errorf("event not found")
	}
	event, ok := eventInterface.(*models.Event)
	if !ok {
		return fmt.Errorf("invalid event type")
	}

	// Find and remove SME from the attending list
	found := false
	newAttendingSmes := make([]int, 0, len(event.AttendingSmes))
	for _, id := range event.AttendingSmes {
		if id == int(smeId) {
			found = true
			continue // Skip this SME
		}
		newAttendingSmes = append(newAttendingSmes, id)
	}

	if !found {
		// SME is not in the attending list, return without error
		return nil
	}

	// Save the event using raw query to update the JSON field
	_, err = facades.Orm().Query().Model(&models.Event{}).
		Where("id = ?", eventId).
		Update(map[string]interface{}{
			"attending_smes": newAttendingSmes,
		})
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

// GetEventAttendees returns paginated list of SMEs attending an event
func (s *EventService) GetEventAttendees(eventId uint, page, perPage int, search string) ([]models.Sme, int64, error) {
	// Fetch the event by ID
	eventInterface, err := s.GetByID(eventId)
	if err != nil {
		return nil, 0, fmt.Errorf("event not found: %w", err)
	}
	if eventInterface == nil {
		return nil, 0, fmt.Errorf("event not found")
	}
	event, ok := eventInterface.(*models.Event)
	if !ok {
		return nil, 0, fmt.Errorf("invalid event type")
	}

	// If no attendees, return empty list
	if len(event.AttendingSmes) == 0 {
		return []models.Sme{}, 0, nil
	}

	// Convert []int to []interface{} for WhereIn
	smeIds := make([]interface{}, len(event.AttendingSmes))
	for i, id := range event.AttendingSmes {
		smeIds[i] = id
	}

	// Build query
	query := facades.Orm().Query().Model(&models.Sme{}).WhereIn("id", smeIds)

	// Apply search filter if provided
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	// Get total count
	var total int64
	total, err = query.Count()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count attendees: %w", err)
	}

	// Apply pagination
	offset := (page - 1) * perPage
	var smes []models.Sme
	err = query.Offset(offset).Limit(perPage).Find(&smes)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch attendees: %w", err)
	}

	return smes, total, nil
}
