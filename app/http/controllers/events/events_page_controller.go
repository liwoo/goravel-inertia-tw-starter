package events

import (
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/services"
)

// EventPageController handles the events page
type EventPageController struct {
	*contracts.GenericPageController
	eventService *services.EventService
}

// NewEventPageController creates a new events page controller
func NewEventPageController() *EventPageController {
	eventService := services.NewEventService()

	return &EventPageController{
		GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
			ResourceType:      "events",
			PageComponent:     "Event/Index",
			Service:           eventService,
			ServiceIdentifier: auth.ServiceEvents,
			StatsEnabled:      false,
		}),
		eventService: eventService,
	}
}
