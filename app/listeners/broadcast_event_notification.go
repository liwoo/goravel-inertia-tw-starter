package listeners

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"

	"smedi-sme-db/app/services"
)

// BroadcastEventNotification listener sends notifications to eligible SMEs when an event is created
type BroadcastEventNotification struct {
}

// Signature returns the listener's signature
func (listener *BroadcastEventNotification) Signature() string {
	return "broadcast_event_notification"
}

// Queue returns the queue details for async processing
func (listener *BroadcastEventNotification) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     true,
		Connection: "",
		Queue:      "",
	}
}

// Handle processes the event asynchronously
// Expected args: senderID (uint), eventID (uint), title (string), date (string), venue (string), district (string)
func (listener *BroadcastEventNotification) Handle(args ...any) error {
	if len(args) < 6 {
		facades.Log().Warning("BroadcastEventNotification: Insufficient arguments")
		return nil
	}

	senderID, ok := args[0].(uint)
	if !ok {
		facades.Log().Warning("BroadcastEventNotification: Invalid senderID type")
		return nil
	}

	eventID, ok := args[1].(uint)
	if !ok {
		facades.Log().Warning("BroadcastEventNotification: Invalid eventID type")
		return nil
	}

	title, ok := args[2].(string)
	if !ok {
		facades.Log().Warning("BroadcastEventNotification: Invalid title type")
		return nil
	}

	date, ok := args[3].(string)
	if !ok {
		facades.Log().Warning("BroadcastEventNotification: Invalid date type")
		return nil
	}

	venue, ok := args[4].(string)
	if !ok {
		facades.Log().Warning("BroadcastEventNotification: Invalid venue type")
		return nil
	}

	district, ok := args[5].(string)
	if !ok {
		facades.Log().Warning("BroadcastEventNotification: Invalid district type")
		return nil
	}

	facades.Log().Info("BroadcastEventNotification: Processing event notification", map[string]interface{}{
		"event_id": eventID,
		"title":    title,
		"sender":   senderID,
	})

	// Use the notification service to broadcast
	notificationService := services.NewNotificationService()
	notificationService.BroadcastEventNotificationSync(senderID, eventID, title, date, venue, district)

	return nil
}
