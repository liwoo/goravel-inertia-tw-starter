package services

import (
	"fmt"
	"time"

	"github.com/goravel/framework/facades"
	"players/app/contracts"
	"players/app/models"
)

type NotificationService struct {
	*contracts.BaseCrudService
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		BaseCrudService: contracts.NewBaseCrudService("notification", "id"),
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

	// Apply search
	if request.Search != "" {
		searchPattern := "%" + request.Search + "%"
		query = query.Where("title LIKE ? OR message LIKE ?", searchPattern, searchPattern)
	}

	// Apply sorting
	if request.Sort == "" {
		request.Sort = "created_at"
		request.Direction = "DESC"
	}
	orderClause := fmt.Sprintf("%s %s", request.Sort, request.Direction)
	query = query.Order(orderClause)

	// Get total count
	var total int64
	query.Model(&models.Notification{}).Count(&total)

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
	return facades.Orm().Query().Save(&notification)
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
	return err
}

// DismissNotification dismisses a notification
func (s *NotificationService) DismissNotification(notificationID, userID uint) error {
	var notification models.Notification
	if err := facades.Orm().Query().Where("id = ? AND user_id = ?", notificationID, userID).First(&notification); err != nil {
		return fmt.Errorf("notification not found or unauthorized")
	}

	notification.Dismiss()
	return facades.Orm().Query().Save(&notification)
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
	return err
}

// GetUnreadNotificationCount returns unread notification count for a user
func (s *NotificationService) GetUnreadNotificationCount(userID uint) (int64, error) {
	var count int64
	err := facades.Orm().Query().Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ? AND is_dismissed = ?", userID, false, false).
		Count(&count)
	return count, err
}

// GetNotificationCounts returns various notification counts for a user
func (s *NotificationService) GetNotificationCounts(userID uint) (map[string]int64, error) {
	counts := make(map[string]int64)

	// Total unread
	var unread int64
	if err := facades.Orm().Query().Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ? AND is_dismissed = ?", userID, false, false).
		Count(&unread); err != nil {
		return nil, err
	}
	counts["unread"] = unread

	// Unread high priority
	var unreadHigh int64
	if err := facades.Orm().Query().Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ? AND is_dismissed = ? AND priority = ?", userID, false, false, "high").
		Count(&unreadHigh); err != nil {
		return nil, err
	}
	counts["unread_high"] = unreadHigh

	// Unread messages
	var unreadMessages int64
	if err := facades.Orm().Query().Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ? AND is_dismissed = ? AND type = ?", userID, false, false, "message").
		Count(&unreadMessages); err != nil {
		return nil, err
	}
	counts["unread_messages"] = unreadMessages

	// Unread mentions
	var unreadMentions int64
	if err := facades.Orm().Query().Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ? AND is_dismissed = ? AND type = ?", userID, false, false, "mention").
		Count(&unreadMentions); err != nil {
		return nil, err
	}
	counts["unread_mentions"] = unreadMentions

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
			Title:       title,
			Message:     message,
			Type:        "system",
			UserID:      user.ID,
			Priority:    priority,
			ExpiresAt:   expiresAt,
		}

		if err := facades.Orm().Query().Create(notification); err != nil {
			facades.Log().Warning("Failed to create system notification for user %d: %v", user.ID, err)
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