package listeners

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"

	"smedi-sme-db/app/services"
)

// BroadcastProcurementNotification listener sends notifications to all SMEs when a procurement is published
type BroadcastProcurementNotification struct {
}

// Signature returns the listener's signature
func (listener *BroadcastProcurementNotification) Signature() string {
	return "broadcast_procurement_notification"
}

// Queue returns the queue details for async processing
func (listener *BroadcastProcurementNotification) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     true,
		Connection: "",
		Queue:      "",
	}
}

// Handle processes the event asynchronously
// Expected args: senderID (uint), procurementID (uint), organization (string), refNo (string), procurementType (string), closeDate (string)
func (listener *BroadcastProcurementNotification) Handle(args ...any) error {
	if len(args) < 6 {
		facades.Log().Warning("BroadcastProcurementNotification: Insufficient arguments")
		return nil
	}

	senderID, ok := args[0].(uint)
	if !ok {
		facades.Log().Warning("BroadcastProcurementNotification: Invalid senderID type")
		return nil
	}

	procurementID, ok := args[1].(uint)
	if !ok {
		facades.Log().Warning("BroadcastProcurementNotification: Invalid procurementID type")
		return nil
	}

	organization, ok := args[2].(string)
	if !ok {
		facades.Log().Warning("BroadcastProcurementNotification: Invalid organization type")
		return nil
	}

	refNo, ok := args[3].(string)
	if !ok {
		facades.Log().Warning("BroadcastProcurementNotification: Invalid refNo type")
		return nil
	}

	procurementType, ok := args[4].(string)
	if !ok {
		facades.Log().Warning("BroadcastProcurementNotification: Invalid procurementType type")
		return nil
	}

	closeDate, ok := args[5].(string)
	if !ok {
		facades.Log().Warning("BroadcastProcurementNotification: Invalid closeDate type")
		return nil
	}

	facades.Log().Info("BroadcastProcurementNotification: Processing procurement notification", map[string]interface{}{
		"procurement_id": procurementID,
		"organization":   organization,
		"sender":         senderID,
	})

	// Use the notification service to broadcast
	notificationService := services.NewNotificationService()
	notificationService.BroadcastProcurementNotificationSync(senderID, procurementID, organization, refNo, procurementType, closeDate)

	return nil
}
