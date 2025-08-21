package sse

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"players/app/models"
	"players/app/services"
)

type SSEController struct {
	sseService          *services.SSEService
	messageService      *services.MessageService
	notificationService *services.NotificationService
}

func NewSSEController() *SSEController {
	return &SSEController{
		sseService:          services.NewSSEService(),
		messageService:      services.NewMessageService(),
		notificationService: services.NewNotificationService(),
	}
}

func (r *SSEController) Stream(ctx http.Context) http.Response {
	// Get authenticated user
	userID := ctx.Value("user_id")
	if userID == nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"error": "Unauthorized",
		})
	}

	// Set SSE headers
	ctx.Response().Header("Content-Type", "text/event-stream")
	ctx.Response().Header("Cache-Control", "no-cache")
	ctx.Response().Header("Connection", "keep-alive")
	ctx.Response().Header("X-Accel-Buffering", "no")

	// Create a client for this connection
	client := r.sseService.CreateClient(userID.(uint))
	defer r.sseService.RemoveClient(client.ID)

	// Send initial connection event
	r.sendEvent(ctx, "connected", map[string]interface{}{
		"message": "Connected to SSE stream",
		"time":    time.Now().Format(time.RFC3339),
	})

	// Send initial data
	r.sendInitialData(ctx, userID.(uint))

	// Keep connection alive and send events
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case event := <-client.Events:
			// Send the event to the client
			if err := r.sendEvent(ctx, event.Type, event.Data); err != nil {
				facades.Log().Error(fmt.Sprintf("Failed to send SSE event: %v", err))
				return nil
			}

		case <-ticker.C:
			// Send heartbeat to keep connection alive
			if err := r.sendEvent(ctx, "heartbeat", map[string]interface{}{
				"time": time.Now().Format(time.RFC3339),
			}); err != nil {
				return nil
			}

		case <-ctx.Request().Context().Done():
			// Client disconnected
			return nil
		}
	}
}

func (r *SSEController) sendEvent(ctx http.Context, eventType string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Format SSE event
	event := fmt.Sprintf("event: %s\ndata: %s\n\n", eventType, string(jsonData))

	// Write to response
	if _, err := ctx.Response().Writer().Write([]byte(event)); err != nil {
		return err
	}

	// Flush the response
	if flusher, ok := ctx.Response().Writer().(http.Flusher); ok {
		flusher.Flush()
	}

	return nil
}

func (r *SSEController) sendInitialData(ctx http.Context, userID uint) {
	// Send unread message count
	unreadCount, _ := r.messageService.GetUnreadCount(userID)
	r.sendEvent(ctx, "message:unread_count", map[string]interface{}{
		"count": unreadCount,
	})

	// Send notification counts
	counts, _ := r.notificationService.GetNotificationCounts(userID)
	r.sendEvent(ctx, "notification:counts", counts)

	// Send recent notifications
	var user models.User
	if err := facades.Orm().Query().Find(&user, userID); err == nil {
		notifications, _ := r.notificationService.GetUserNotifications(&user, 1, 10, "unread")
		r.sendEvent(ctx, "notification:initial", notifications)
	}
}