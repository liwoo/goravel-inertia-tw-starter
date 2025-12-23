package portal

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/services"
)

// PortalEventsController handles event attendance from SME portal
type PortalEventsController struct {
	eventService *services.EventService
	smeService   *services.SmeService
}

// NewPortalEventsController creates a new portal events controller
func NewPortalEventsController() *PortalEventsController {
	return &PortalEventsController{
		eventService: services.NewEventService(),
		smeService:   services.NewSmeService(),
	}
}

// Attend adds the authenticated user's SME to the event's attending list
// POST /api/portal/events/{id}/attend
func (c *PortalEventsController) Attend(ctx http.Context) http.Response {
	// Get the authenticated user
	user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
	if user == nil {
		return ctx.Response().Status(http.StatusUnauthorized).Json(http.Json{
			"success": false,
			"message": "Authentication required",
		})
	}

	// Get the user's linked SME
	sme, err := c.smeService.GetSmeByUserEmail(user.Email)
	if err != nil || sme == nil {
		return ctx.Response().Status(http.StatusBadRequest).Json(http.Json{
			"success": false,
			"message": "No SME linked to your account",
		})
	}

	// Get event ID from URL parameter
	eventIdStr := ctx.Request().Route("id")
	if eventIdStr == "" {
		return ctx.Response().Status(http.StatusBadRequest).Json(http.Json{
			"success": false,
			"message": "Event ID is required",
		})
	}

	eventId, err := strconv.ParseUint(eventIdStr, 10, 32)
	if err != nil {
		return ctx.Response().Status(http.StatusBadRequest).Json(http.Json{
			"success": false,
			"message": "Invalid event ID",
		})
	}

	// Add SME to event's attending list
	if err := c.eventService.AttendEvent(uint(eventId), sme.ID); err != nil {
		return ctx.Response().Status(http.StatusInternalServerError).Json(http.Json{
			"success": false,
			"message": "Failed to attend event: " + err.Error(),
		})
	}

	return ctx.Response().Status(http.StatusOK).Json(http.Json{
		"success": true,
		"message": "Successfully registered for event",
		"data": http.Json{
			"eventId": eventId,
			"smeId":   sme.ID,
		},
	})
}

// Unattend removes the authenticated user's SME from the event's attending list
// DELETE /api/portal/events/{id}/attend
func (c *PortalEventsController) Unattend(ctx http.Context) http.Response {
	// Get the authenticated user
	user := auth.GetPermissionHelper().GetAuthenticatedUser(ctx)
	if user == nil {
		return ctx.Response().Status(http.StatusUnauthorized).Json(http.Json{
			"success": false,
			"message": "Authentication required",
		})
	}

	// Get the user's linked SME
	sme, err := c.smeService.GetSmeByUserEmail(user.Email)
	if err != nil || sme == nil {
		return ctx.Response().Status(http.StatusBadRequest).Json(http.Json{
			"success": false,
			"message": "No SME linked to your account",
		})
	}

	// Get event ID from URL parameter
	eventIdStr := ctx.Request().Route("id")
	if eventIdStr == "" {
		return ctx.Response().Status(http.StatusBadRequest).Json(http.Json{
			"success": false,
			"message": "Event ID is required",
		})
	}

	eventId, err := strconv.ParseUint(eventIdStr, 10, 32)
	if err != nil {
		return ctx.Response().Status(http.StatusBadRequest).Json(http.Json{
			"success": false,
			"message": "Invalid event ID",
		})
	}

	// Remove SME from event's attending list
	if err := c.eventService.UnattendEvent(uint(eventId), sme.ID); err != nil {
		return ctx.Response().Status(http.StatusInternalServerError).Json(http.Json{
			"success": false,
			"message": "Failed to unattend event: " + err.Error(),
		})
	}

	return ctx.Response().Status(http.StatusOK).Json(http.Json{
		"success": true,
		"message": "Successfully unregistered from event",
		"data": http.Json{
			"eventId": eventId,
			"smeId":   sme.ID,
		},
	})
}
