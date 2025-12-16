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
