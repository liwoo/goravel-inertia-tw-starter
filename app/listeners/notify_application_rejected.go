package listeners

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"

	"books-database/app/services"
)

// NotifyApplicationRejected listener sends notification to SME when their application is rejected
type NotifyApplicationRejected struct {
}

// Signature returns the listener's signature
func (listener *NotifyApplicationRejected) Signature() string {
	return "notify_application_rejected"
}

// Queue returns the queue details for async processing
func (listener *NotifyApplicationRejected) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     true,
		Connection: "",
		Queue:      "",
	}
}

// Handle processes the event asynchronously
// Expected args: rejectorUserID (uint), applicationID (uint), applicationType (string), smeName (string), recipientEmail (string), reason (string)
func (listener *NotifyApplicationRejected) Handle(args ...any) error {
	if len(args) < 6 {
		facades.Log().Warning("NotifyApplicationRejected: Insufficient arguments")
		return nil
	}

	// Note: After JSON serialization in queue, numbers come back as float64
	rejectorUserID, ok := toUint(args[0])
	if !ok {
		facades.Log().Warning("NotifyApplicationRejected: Invalid rejectorUserID type")
		return nil
	}

	applicationID, ok := toUint(args[1])
	if !ok {
		facades.Log().Warning("NotifyApplicationRejected: Invalid applicationID type")
		return nil
	}

	applicationType, ok := args[2].(string)
	if !ok {
		facades.Log().Warning("NotifyApplicationRejected: Invalid applicationType type")
		return nil
	}

	smeName, ok := args[3].(string)
	if !ok {
		facades.Log().Warning("NotifyApplicationRejected: Invalid smeName type")
		return nil
	}

	recipientEmail, ok := args[4].(string)
	if !ok {
		facades.Log().Warning("NotifyApplicationRejected: Invalid recipientEmail type")
		return nil
	}

	reason, ok := args[5].(string)
	if !ok {
		facades.Log().Warning("NotifyApplicationRejected: Invalid reason type")
		return nil
	}

	facades.Log().Info("NotifyApplicationRejected: Processing application rejection notification", map[string]interface{}{
		"application_id": applicationID,
		"type":           applicationType,
		"sme_name":       smeName,
		"rejector":       rejectorUserID,
		"reason":         reason,
	})

	// Use the notification service to send notification
	notificationService := services.NewNotificationService()
	notificationService.NotifyApplicationStatusSync(
		rejectorUserID,
		applicationID,
		applicationType,
		smeName,
		recipientEmail,
		"rejected",
		reason,
	)

	return nil
}
