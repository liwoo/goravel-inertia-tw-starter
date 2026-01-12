package services

import (
	"fmt"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/models"
)

type NotificationService struct {
	*contracts.BaseCrudService
	sseService   *SSEService
	cacheService *CacheService
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		BaseCrudService: contracts.NewBaseCrudService("notification", "id"),
		sseService:      NewSSEService(),
		cacheService:    GetCacheService(),
	}
}

// GetUserNotifications retrieves notifications for a user
func (s *NotificationService) GetUserNotifications(userID uint, request contracts.ListRequest) (*contracts.PaginatedResult, error) {
	query := facades.Orm().Query().Model(&models.Notification{}).
		With("TriggerUser").
		Where("user_id = ? AND is_dismissed = ?", userID, false)

	// Filter by unread if specified
	if unreadOnly, exists := request.Filters["unread_only"].(bool); exists && unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	// Filter by type if specified
	if notifType, exists := request.Filters["type"].(string); exists && notifType != "" {
		query = query.Where("type = ?", notifType)
	}

	// Filter by priority if specified
	if priority, exists := request.Filters["priority"].(string); exists && priority != "" {
		query = query.Where("priority = ?", priority)
	}

	// Apply search (using ILIKE for case-insensitive PostgreSQL search)
	if request.Search != "" {
		searchPattern := "%" + request.Search + "%"
		query = query.Where("title ILIKE ? OR message ILIKE ?", searchPattern, searchPattern)
	}

	// Apply sorting
	if request.Sort == "" {
		request.Sort = "created_at"
		request.Direction = "DESC"
	}
	orderClause := fmt.Sprintf("%s %s", request.Sort, request.Direction)
	query = query.Order(orderClause)

	// Get total count
	total, _ := query.Model(&models.Notification{}).Count()

	// Apply pagination
	offset := (request.Page - 1) * request.PageSize
	var notifications []models.Notification
	if err := query.Offset(offset).Limit(request.PageSize).Find(&notifications); err != nil {
		return nil, fmt.Errorf("failed to retrieve notifications: %v", err)
	}

	// Convert to interface slice
	data := make([]interface{}, len(notifications))
	for i, notification := range notifications {
		data[i] = notification
	}

	lastPage := int((total + int64(request.PageSize) - 1) / int64(request.PageSize))
	if lastPage < 1 {
		lastPage = 1
	}

	return &contracts.PaginatedResult{
		Data:        data,
		Total:       total,
		CurrentPage: request.Page,
		LastPage:    lastPage,
		PerPage:     request.PageSize,
		From:        offset + 1,
		To:          offset + len(notifications),
		HasNext:     request.Page < lastPage,
		HasPrev:     request.Page > 1,
	}, nil
}

// CreateNotification creates a new notification
func (s *NotificationService) CreateNotification(
	userID uint,
	title, message, notificationType string,
	triggerUserID *uint,
	relatedType *string,
	relatedID *uint,
	priority string,
	expiresAt *time.Time,
	data string,
) (*models.Notification, error) {

	// Validate required fields
	if title == "" {
		return nil, fmt.Errorf("notification title is required")
	}
	if notificationType == "" {
		return nil, fmt.Errorf("notification type is required")
	}
	if priority == "" {
		priority = "normal"
	}

	// Validate user exists
	var user models.User
	if err := facades.Orm().Query().Where("id = ? AND is_active = ?", userID, true).First(&user); err != nil {
		return nil, fmt.Errorf("user not found or inactive")
	}

	// Validate trigger user if provided
	if triggerUserID != nil {
		var triggerUser models.User
		if err := facades.Orm().Query().Where("id = ? AND is_active = ?", *triggerUserID, true).First(&triggerUser); err != nil {
			return nil, fmt.Errorf("trigger user not found or inactive")
		}
	}

	// Create notification
	notification := &models.Notification{
		Title:         title,
		Message:       message,
		Type:          notificationType,
		UserID:        userID,
		TriggerUserID: triggerUserID,
		RelatedType:   "",
		RelatedID:     relatedID,
		Priority:      priority,
		ExpiresAt:     expiresAt,
		Data:          data,
	}

	if relatedType != nil {
		notification.RelatedType = *relatedType
	}

	if err := facades.Orm().Query().Create(notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %v", err)
	}

	// Load relations for response
	facades.Orm().Query().Model(&models.Notification{}).With("TriggerUser").Where("id = ?", notification.ID).First(notification)

	// Emit SSE event for new notification
	s.sseService.NotifyNewNotification(userID, notification)

	// Update notification counts
	if counts, err := s.GetNotificationCounts(userID); err == nil {
		s.sseService.UpdateNotificationCounts(userID, counts)
	}

	return notification, nil
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(notificationID, userID uint) error {
	var notification models.Notification
	if err := facades.Orm().Query().Where("id = ? AND user_id = ?", notificationID, userID).First(&notification); err != nil {
		return fmt.Errorf("notification not found or unauthorized")
	}

	if notification.IsRead {
		return nil // Already read
	}

	notification.MarkAsRead()
	err := facades.Orm().Query().Save(&notification)

	if err == nil {
		// Invalidate notification counts cache
		s.cacheService.InvalidateNotificationCounts(userID)

		// Emit SSE event
		s.sseService.NotifyNotificationRead(userID, notificationID)

		// Update notification counts (will re-cache)
		if counts, err := s.GetNotificationCounts(userID); err == nil {
			s.sseService.UpdateNotificationCounts(userID, counts)
		}
	}

	return err
}

// MarkAllAsRead marks all notifications as read for a user
func (s *NotificationService) MarkAllAsRead(userID uint) error {
	now := time.Now()
	_, err := facades.Orm().Query().
		Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		})

	if err == nil {
		// Invalidate notification counts cache
		s.cacheService.InvalidateNotificationCounts(userID)

		// Update notification counts (will re-cache)
		if counts, err := s.GetNotificationCounts(userID); err == nil {
			s.sseService.UpdateNotificationCounts(userID, counts)
		}
	}

	return err
}

// DismissNotification dismisses a notification
func (s *NotificationService) DismissNotification(notificationID, userID uint) error {
	var notification models.Notification
	if err := facades.Orm().Query().Where("id = ? AND user_id = ?", notificationID, userID).First(&notification); err != nil {
		return fmt.Errorf("notification not found or unauthorized")
	}

	notification.Dismiss()
	err := facades.Orm().Query().Save(&notification)

	if err == nil {
		// Invalidate notification counts cache
		s.cacheService.InvalidateNotificationCounts(userID)

		// Emit SSE event
		s.sseService.NotifyNotificationDismissed(userID, notificationID)

		// Update notification counts (will re-cache)
		if counts, err := s.GetNotificationCounts(userID); err == nil {
			s.sseService.UpdateNotificationCounts(userID, counts)
		}
	}

	return err
}

// DismissAllNotifications dismisses all notifications for a user
func (s *NotificationService) DismissAllNotifications(userID uint) error {
	now := time.Now()
	_, err := facades.Orm().Query().
		Model(&models.Notification{}).
		Where("user_id = ? AND is_dismissed = ?", userID, false).
		Update(map[string]interface{}{
			"is_dismissed": true,
			"dismissed_at": now,
		})

	if err == nil {
		// Invalidate notification counts cache
		s.cacheService.InvalidateNotificationCounts(userID)
	}

	return err
}

// GetUnreadNotificationCount returns unread notification count for a user
func (s *NotificationService) GetUnreadNotificationCount(userID uint) (int64, error) {
	return facades.Orm().Query().Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ? AND is_dismissed = ?", userID, false, false).
		Count()
}

// GetNotificationCounts returns various notification counts for a user
// ctx is optional - if provided, it's used to check permissions for pending applications
func (s *NotificationService) GetNotificationCounts(userID uint, ctx ...http.Context) (map[string]int64, error) {
	counts := make(map[string]int64)

	// Try to get notification counts from Redis cache first
	if cachedCounts, found := s.cacheService.GetNotificationCounts(userID); found {
		facades.Log().Debug("Notification counts cache hit", map[string]interface{}{
			"user_id": userID,
		})
		// Copy cached notification counts
		for k, v := range cachedCounts {
			counts[k] = v
		}
	} else {
		// Cache miss - fetch notification counts from database
		// Total unread
		unread, err := facades.Orm().Query().Model(&models.Notification{}).
			Where("user_id = ? AND is_read = ? AND is_dismissed = ?", userID, false, false).
			Count()
		if err != nil {
			return nil, err
		}
		counts["unread"] = unread

		// Unread high priority
		unreadHigh, err := facades.Orm().Query().Model(&models.Notification{}).
			Where("user_id = ? AND is_read = ? AND is_dismissed = ? AND priority = ?", userID, false, false, "high").
			Count()
		if err != nil {
			return nil, err
		}
		counts["unread_high"] = unreadHigh

		// Unread messages
		unreadMessages, err := facades.Orm().Query().Model(&models.Notification{}).
			Where("user_id = ? AND is_read = ? AND is_dismissed = ? AND type = ?", userID, false, false, "message").
			Count()
		if err != nil {
			return nil, err
		}
		counts["unread_messages"] = unreadMessages

		// Unread mentions
		unreadMentions, err := facades.Orm().Query().Model(&models.Notification{}).
			Where("user_id = ? AND is_read = ? AND is_dismissed = ? AND type = ?", userID, false, false, "mention").
			Count()
		if err != nil {
			return nil, err
		}
		counts["unread_mentions"] = unreadMentions

		// Cache the notification counts in Redis (excluding pending_applications)
		if err := s.cacheService.SetNotificationCounts(userID, counts); err != nil {
			facades.Log().Warning("Failed to cache notification counts", map[string]interface{}{
				"user_id": userID,
				"error":   err.Error(),
			})
		}
	}

	// Check if user has permission to manage applications before including pending count
	// Only users with update permission on applications can see pending applications count
	canManageApplications := false
	if len(ctx) > 0 && ctx[0] != nil {
		scopedHelper := auth.GetScopedPermissionHelper()
		canManageApplications = scopedHelper.CheckScopedPermission(ctx[0], auth.ServiceApplications, auth.PermissionUpdate, nil)
	}

	if canManageApplications {
		// Fetch pending applications count fresh (not cached per-user)
		// This ensures the count is always up-to-date when applications are created/processed
		pendingApplications, err := facades.Orm().Query().Model(&models.Application{}).
			Where("status = ?", "Pending").
			Count()
		if err != nil {
			// Log but don't fail - this is optional
			facades.Log().Warning("Failed to get pending applications count", map[string]interface{}{
				"error": err.Error(),
			})
			pendingApplications = 0
		}
		counts["pending_applications"] = pendingApplications
	}

	return counts, nil
}

// CleanupExpiredNotifications removes expired notifications
func (s *NotificationService) CleanupExpiredNotifications() error {
	_, err := facades.Orm().Query().
		Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).
		Delete(&models.Notification{})
	return err
}

// CreateSystemNotification creates a system-wide notification for all users
func (s *NotificationService) CreateSystemNotification(
	title, message string,
	priority string,
	expiresAt *time.Time,
	roleFilter []string, // Optional: only send to users with these roles
) error {

	if priority == "" {
		priority = "normal"
	}

	// Get users to notify
	query := facades.Orm().Query().Where("is_active = ?", true)

	if len(roleFilter) > 0 {
		query = query.
			Where("EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = users.id AND r.slug IN ?)", roleFilter)
	}

	var users []models.User
	if err := query.Find(&users); err != nil {
		return fmt.Errorf("failed to get users for system notification: %v", err)
	}

	// Create notifications for each user
	for _, user := range users {
		notification := &models.Notification{
			Title:     title,
			Message:   message,
			Type:      "system",
			UserID:    user.ID,
			Priority:  priority,
			ExpiresAt: expiresAt,
		}

		if err := facades.Orm().Query().Create(notification); err != nil {
			facades.Log().Warning("Failed to create system notification for user %d: %v", user.ID, err)
		} else {
			// Emit SSE event for system notification
			s.sseService.NotifyNewNotification(user.ID, notification)
		}
	}

	return nil
}

// GetNotificationsByType gets notifications of a specific type for a user
func (s *NotificationService) GetNotificationsByType(userID uint, notificationType string, limit int) ([]models.Notification, error) {
	var notifications []models.Notification

	query := facades.Orm().Query().Model(&models.Notification{}).
		With("TriggerUser").
		Where("user_id = ? AND type = ? AND is_dismissed = ?", userID, notificationType, false).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&notifications); err != nil {
		return nil, fmt.Errorf("failed to retrieve notifications: %v", err)
	}

	return notifications, nil
}

// BatchMarkAsRead marks multiple notifications as read
func (s *NotificationService) BatchMarkAsRead(notificationIDs []uint, userID uint) error {
	if len(notificationIDs) == 0 {
		return nil
	}

	now := time.Now()
	_, err := facades.Orm().Query().
		Model(&models.Notification{}).
		Where("id IN ? AND user_id = ? AND is_read = ?", notificationIDs, userID, false).
		Update(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		})
	return err
}

// BatchDismiss dismisses multiple notifications
func (s *NotificationService) BatchDismiss(notificationIDs []uint, userID uint) error {
	if len(notificationIDs) == 0 {
		return nil
	}

	now := time.Now()
	_, err := facades.Orm().Query().
		Model(&models.Notification{}).
		Where("id IN ? AND user_id = ?", notificationIDs, userID).
		Update(map[string]interface{}{
			"is_dismissed": true,
			"dismissed_at": now,
		})
	return err
}

// BroadcastEventNotificationSync sends notifications and messages to eligible SMEs when an event is created
// This is the synchronous version called by the event listener (async is handled by the queue)
// Eligible SMEs are those in the event's district OR if event has no district (nationwide)
func (s *NotificationService) BroadcastEventNotificationSync(senderID uint, eventID uint, eventTitle string, eventDate string, eventVenue string, eventDistrict string) {
	facades.Log().Info("Starting event broadcast notification", map[string]interface{}{
		"event_id":       eventID,
		"event_title":    eventTitle,
		"event_district": eventDistrict,
		"sender_id":      senderID,
	})

	// Get eligible SME users
	userIDs := s.getEligibleSMEUsersForEvent(eventDistrict)
	if len(userIDs) == 0 {
		facades.Log().Info("No eligible SME users found for event notification", map[string]interface{}{
			"event_id": eventID,
		})
		return
	}

	facades.Log().Info("Found eligible SME users for event notification", map[string]interface{}{
		"event_id":   eventID,
		"user_count": len(userIDs),
	})

	// Prepare notification content
	title := "New Event: " + eventTitle
	message := fmt.Sprintf("A new event \"%s\" has been scheduled", eventTitle)
	if eventDate != "" {
		message += fmt.Sprintf(" on %s", eventDate)
	}
	if eventVenue != "" {
		message += fmt.Sprintf(" at %s", eventVenue)
	}
	if eventDistrict != "" {
		message += fmt.Sprintf(" in %s", eventDistrict)
	}
	message += ". Check it out!"

	// Send notifications to all eligible users
	relatedType := "event"
	for _, userID := range userIDs {
		_, err := s.CreateNotification(
			userID,
			title,
			message,
			"event",
			&senderID,
			&relatedType,
			&eventID,
			"normal",
			nil,
			"",
		)
		if err != nil {
			facades.Log().Warning("Failed to create event notification for user", map[string]interface{}{
				"user_id":  userID,
				"event_id": eventID,
				"error":    err.Error(),
			})
		}
	}

	// Also send messages to users
	messageService := NewMessageService()
	messageContent := fmt.Sprintf("📅 **New Event Announcement**\n\n**%s**\n\n", eventTitle)
	if eventDate != "" {
		messageContent += fmt.Sprintf("🗓️ Date: %s\n", eventDate)
	}
	if eventVenue != "" {
		messageContent += fmt.Sprintf("📍 Venue: %s\n", eventVenue)
	}
	if eventDistrict != "" {
		messageContent += fmt.Sprintf("🌍 District: %s\n", eventDistrict)
	}
	messageContent += "\nDon't miss out on this opportunity!"

	_, err := messageService.SendBroadcast(senderID, userIDs, messageContent, "New Event: "+eventTitle)
	if err != nil {
		facades.Log().Warning("Failed to send event broadcast messages", map[string]interface{}{
			"event_id": eventID,
			"error":    err.Error(),
		})
	}

	facades.Log().Info("Event broadcast notification completed", map[string]interface{}{
		"event_id":       eventID,
		"notified_users": len(userIDs),
	})
}

// BroadcastProcurementNotificationSync sends notifications and messages to all active SME users when a procurement is published
// This is the synchronous version called by the event listener (async is handled by the queue)
func (s *NotificationService) BroadcastProcurementNotificationSync(senderID uint, procurementID uint, organization string, refNo string, procurementType string, closeDate string) {
	facades.Log().Info("Starting procurement broadcast notification", map[string]interface{}{
		"procurement_id": procurementID,
		"organization":   organization,
		"ref_no":         refNo,
		"sender_id":      senderID,
	})

	// Get all active SME users
	userIDs := s.getAllActiveSMEUsers()
	if len(userIDs) == 0 {
		facades.Log().Info("No active SME users found for procurement notification", map[string]interface{}{
			"procurement_id": procurementID,
		})
		return
	}

	facades.Log().Info("Found active SME users for procurement notification", map[string]interface{}{
		"procurement_id": procurementID,
		"user_count":     len(userIDs),
	})

	// Prepare notification content
	title := "New Procurement Opportunity"
	message := fmt.Sprintf("New procurement notice from %s (Ref: %s)", organization, refNo)
	if procurementType != "" {
		message = fmt.Sprintf("New %s opportunity from %s (Ref: %s)", procurementType, organization, refNo)
	}
	if closeDate != "" {
		message += fmt.Sprintf(". Closing date: %s", closeDate)
	}

	// Send notifications to all active SME users
	relatedType := "procurement"
	for _, userID := range userIDs {
		_, err := s.CreateNotification(
			userID,
			title,
			message,
			"procurement",
			&senderID,
			&relatedType,
			&procurementID,
			"normal",
			nil,
			"",
		)
		if err != nil {
			facades.Log().Warning("Failed to create procurement notification for user", map[string]interface{}{
				"user_id":        userID,
				"procurement_id": procurementID,
				"error":          err.Error(),
			})
		}
	}

	// Also send messages to users
	messageService := NewMessageService()
	messageContent := fmt.Sprintf("📋 **New Procurement Opportunity**\n\n**%s**\nRef: %s\n", organization, refNo)
	if procurementType != "" {
		messageContent += fmt.Sprintf("Type: %s\n", procurementType)
	}
	if closeDate != "" {
		messageContent += fmt.Sprintf("⏰ Closing Date: %s\n", closeDate)
	}
	messageContent += "\nSubmit your bid before the deadline!"

	_, err := messageService.SendBroadcast(senderID, userIDs, messageContent, "New Procurement: "+organization)
	if err != nil {
		facades.Log().Warning("Failed to send procurement broadcast messages", map[string]interface{}{
			"procurement_id": procurementID,
			"error":          err.Error(),
		})
	}

	facades.Log().Info("Procurement broadcast notification completed", map[string]interface{}{
		"procurement_id": procurementID,
		"notified_users": len(userIDs),
	})
}

// getEligibleSMEUsersForEvent returns user IDs linked to SMEs eligible for an event
// Eligible: SMEs in the event's district OR if event has no district (available to all)
func (s *NotificationService) getEligibleSMEUsersForEvent(eventDistrict string) []uint {
	var userIDs []uint

	// Build query based on district
	var query string
	var args []interface{}

	if eventDistrict == "" {
		// Event has no district - notify all active SME users
		query = `
			SELECT DISTINCT u.id
			FROM users u
			INNER JOIN smes s ON LOWER(s.contact_email) = LOWER(u.email)
			WHERE u.is_active = true
			AND u.deleted_at IS NULL
			AND s.is_active = true
			AND s.deleted_at IS NULL
		`
	} else {
		// Event has district - notify SMEs in that district OR SMEs with no district
		query = `
			SELECT DISTINCT u.id
			FROM users u
			INNER JOIN smes s ON LOWER(s.contact_email) = LOWER(u.email)
			WHERE u.is_active = true
			AND u.deleted_at IS NULL
			AND s.is_active = true
			AND s.deleted_at IS NULL
			AND (s.district = ? OR s.district IS NULL OR s.district = '')
		`
		args = append(args, eventDistrict)
	}

	err := facades.Orm().Query().Raw(query, args...).Pluck("id", &userIDs)
	if err != nil {
		facades.Log().Error("Failed to get eligible SME users for event", map[string]interface{}{
			"district": eventDistrict,
			"error":    err.Error(),
		})
		return []uint{}
	}

	return userIDs
}

// getAllActiveSMEUsers returns user IDs for all active SMEs
func (s *NotificationService) getAllActiveSMEUsers() []uint {
	var userIDs []uint

	query := `
		SELECT DISTINCT u.id
		FROM users u
		INNER JOIN smes s ON LOWER(s.contact_email) = LOWER(u.email)
		WHERE u.is_active = true
		AND u.deleted_at IS NULL
		AND s.is_active = true
		AND s.deleted_at IS NULL
	`

	err := facades.Orm().Query().Raw(query).Pluck("id", &userIDs)
	if err != nil {
		facades.Log().Error("Failed to get all active SME users", map[string]interface{}{
			"error": err.Error(),
		})
		return []uint{}
	}

	return userIDs
}

// NotifyApplicationStatusSync sends notification and message to an SME when their application is approved/rejected
// This is the synchronous version called by the event listener (async is handled by the queue)
func (s *NotificationService) NotifyApplicationStatusSync(
	senderUserID uint,
	applicationID uint,
	applicationType string,
	smeName string,
	recipientEmail string,
	status string, // "approved" or "rejected"
	reason string, // rejection reason (empty for approvals)
) {
	facades.Log().Info("Starting application status notification", map[string]interface{}{
		"application_id":   applicationID,
		"application_type": applicationType,
		"sme_name":         smeName,
		"recipient_email":  recipientEmail,
		"status":           status,
		"sender_id":        senderUserID,
	})

	// Find the user by email
	var user models.User
	if err := facades.Orm().Query().Where("LOWER(email) = LOWER(?) AND is_active = ?", recipientEmail, true).First(&user); err != nil {
		facades.Log().Warning("User not found for application notification - may be a signup rejection", map[string]interface{}{
			"email":          recipientEmail,
			"application_id": applicationID,
			"status":         status,
		})
		return
	}

	// Build notification content based on status and type
	var title, message, messageContent string
	var priority string = "normal"

	if status == "approved" {
		if applicationType == "amend_formalisation" {
			title = "Amendment Approved"
			message = fmt.Sprintf("Your formalisation amendment for %s has been approved. The changes have been applied to your MSME profile.", smeName)
			messageContent = fmt.Sprintf("✅ **Amendment Approved**\n\nYour formalisation amendment request for **%s** has been approved.\n\nThe requested changes have been applied to your MSME profile. You can view your updated profile in the Member Portal.", smeName)
		} else {
			title = "Application Approved"
			message = fmt.Sprintf("Your application for %s has been approved. Welcome to SMEDI!", smeName)
			messageContent = fmt.Sprintf("✅ **Application Approved**\n\nCongratulations! Your application for **%s** has been approved.\n\nWelcome to the MSME Database! You can now access all the features available in the Member Portal.", smeName)
		}
	} else { // rejected
		priority = "high"
		if applicationType == "amend_formalisation" {
			title = "Amendment Rejected"
			message = fmt.Sprintf("Your formalisation amendment for %s has been rejected.", smeName)
			if reason != "" {
				message += fmt.Sprintf(" Reason: %s", reason)
			}
			messageContent = fmt.Sprintf("❌ **Amendment Rejected**\n\nYour formalisation amendment request for **%s** has been rejected.", smeName)
			if reason != "" {
				messageContent += fmt.Sprintf("\n\n**Reason:** %s", reason)
			}
			messageContent += "\n\nIf you have questions about this decision, please contact our support team."
		} else {
			title = "Application Rejected"
			message = fmt.Sprintf("Your application for %s has been rejected.", smeName)
			if reason != "" {
				message += fmt.Sprintf(" Reason: %s", reason)
			}
			messageContent = fmt.Sprintf("❌ **Application Rejected**\n\nWe regret to inform you that your application for **%s** has been rejected.", smeName)
			if reason != "" {
				messageContent += fmt.Sprintf("\n\n**Reason:** %s", reason)
			}
			messageContent += "\n\nIf you believe this decision was made in error or have questions, please contact our support team."
		}
	}

	// Create notification
	relatedType := "application"
	_, err := s.CreateNotification(
		user.ID,
		title,
		message,
		"application_"+status,
		&senderUserID,
		&relatedType,
		&applicationID,
		priority,
		nil,
		"",
	)
	if err != nil {
		facades.Log().Warning("Failed to create application status notification", map[string]interface{}{
			"user_id":        user.ID,
			"application_id": applicationID,
			"status":         status,
			"error":          err.Error(),
		})
	}

	// Send message to the user
	messageService := NewMessageService()
	_, err = messageService.SendMessage(senderUserID, user.ID, messageContent, models.MessageTypeDirect)
	if err != nil {
		facades.Log().Warning("Failed to send application status message", map[string]interface{}{
			"user_id":        user.ID,
			"application_id": applicationID,
			"status":         status,
			"error":          err.Error(),
		})
	}

	facades.Log().Info("Application status notification completed", map[string]interface{}{
		"application_id": applicationID,
		"user_id":        user.ID,
		"status":         status,
	})
}
