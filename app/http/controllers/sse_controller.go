package controllers

import (
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"time"

	"books-database/app/auth"
	"books-database/app/services"
	"github.com/goravel/framework/contracts/http"
)

// SSEController handles Server-Sent Events for real-time updates
type SSEController struct {
	sseService *services.SSEService
}

// NewSSEController creates a new SSE controller
func NewSSEController() *SSEController {
	return &SSEController{
		sseService: services.NewSSEService(),
	}
}

// Stream handles GET /api/sse/stream
// This establishes an SSE connection for real-time updates
func (c *SSEController) Stream(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return ctx.Response().Status(401).Json(http.Json{
			"success": false,
			"message": "Authentication required",
		})
	}

	// Set SSE headers
	ctx.Response().Header("Content-Type", "text/event-stream")
	ctx.Response().Header("Cache-Control", "no-cache")
	ctx.Response().Header("Connection", "keep-alive")
	ctx.Response().Header("X-Accel-Buffering", "no")

	// Create SSE client
	client := c.sseService.CreateClient(user.ID)
	defer c.sseService.RemoveClient(client.ID)

	// Get the underlying response writer
	writer := ctx.Response().Writer()

	// Send initial connection event
	initialEvent := map[string]interface{}{
		"type":      "connected",
		"client_id": client.ID,
		"user_id":   user.ID,
		"time":      time.Now(),
	}
	initialData, _ := json.Marshal(initialEvent)
	fmt.Fprintf(writer, "event: connected\ndata: %s\n\n", initialData)
	if flusher, ok := writer.(nethttp.Flusher); ok {
		flusher.Flush()
	}

	// Send initial presence state (list of online users)
	onlineUserIDs := c.sseService.GetOnlineUserIDs()
	presenceEvent := map[string]interface{}{
		"type":         "presence:initial",
		"online_users": onlineUserIDs,
		"time":         time.Now(),
	}
	presenceData, _ := json.Marshal(presenceEvent)
	fmt.Fprintf(writer, "event: presence:initial\ndata: %s\n\n", presenceData)
	if flusher, ok := writer.(nethttp.Flusher); ok {
		flusher.Flush()
	}

	// Keep-alive ticker
	keepAliveTicker := time.NewTicker(30 * time.Second)
	defer keepAliveTicker.Stop()

	// Create a done channel to detect client disconnect
	done := make(chan bool)
	go func() {
		// Monitor for client disconnect using the origin request context
		<-ctx.Request().Origin().Context().Done()
		done <- true
	}()

	// Helper function to safely write to the response writer
	writeToClient := func(data string) bool {
		_, err := fmt.Fprint(writer, data)
		if err != nil {
			// Write failed - connection is broken
			return false
		}
		if flusher, ok := writer.(nethttp.Flusher); ok {
			flusher.Flush()
		}
		return true
	}

	// Listen for events
	for {
		select {
		case event, ok := <-client.Events:
			if !ok {
				// Channel closed, client disconnected
				return nil
			}

			// Marshal event data
			eventData, err := json.Marshal(event)
			if err != nil {
				continue
			}

			// Write SSE event
			if !writeToClient(fmt.Sprintf("event: %s\ndata: %s\n\n", event.Type, eventData)) {
				return nil
			}

		case <-keepAliveTicker.C:
			// Send keep-alive comment - this also detects broken connections
			if !writeToClient(": keep-alive\n\n") {
				return nil
			}

		case <-done:
			// Client disconnected
			return nil
		}
	}
}
