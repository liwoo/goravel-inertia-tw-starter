package listeners

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"

	"books-database/app/services"
)

// NotifyApplicationApproved listener sends notification to SME when their application is approved
type NotifyApplicationApproved struct {
}

// Signature returns the listener's signature
func (listener *NotifyApplicationApproved) Signature() string {
	return "notify_application_approved"
}

// Queue returns the queue details for async processing
func (listener *NotifyApplicationApproved) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     true,
		Connection: "",
		Queue:      "",
	}
}

// Handle processes the event asynchronously
// Expected args: approverUserID (uint), applicationID (uint), applicationType (string), smeName (string), recipientEmail (string)
func (listener *NotifyApplicationApproved) Handle(args ...any) error {
	if len(args) < 5 {
		facades.Log().Warning("NotifyApplicationApproved: Insufficient arguments")
		return nil
	}

	// Note: After JSON serialization in queue, numbers come back as float64
	approverUserID, ok := toUint(args[0])
	if !ok {
		facades.Log().Warning("NotifyApplicationApproved: Invalid approverUserID type")
		return nil
	}

	applicationID, ok := toUint(args[1])
	if !ok {
		facades.Log().Warning("NotifyApplicationApproved: Invalid applicationID type")
		return nil
	}

	applicationType, ok := args[2].(string)
	if !ok {
		facades.Log().Warning("NotifyApplicationApproved: Invalid applicationType type")
		return nil
	}

	smeName, ok := args[3].(string)
	if !ok {
		facades.Log().Warning("NotifyApplicationApproved: Invalid smeName type")
		return nil
	}

	recipientEmail, ok := args[4].(string)
	if !ok {
		facades.Log().Warning("NotifyApplicationApproved: Invalid recipientEmail type")
		return nil
	}

	facades.Log().Info("NotifyApplicationApproved: Processing application approval notification", map[string]interface{}{
		"application_id": applicationID,
		"type":           applicationType,
		"sme_name":       smeName,
		"approver":       approverUserID,
	})

	// Use the notification service to send notification
	notificationService := services.NewNotificationService()
	notificationService.NotifyApplicationStatusSync(
		approverUserID,
		applicationID,
		applicationType,
		smeName,
		recipientEmail,
		"approved",
		"",
	)

	return nil
}
