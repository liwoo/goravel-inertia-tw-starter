package services

import (
	"fmt"
	"strings"
	"time"

	"starter-project/app/contracts"
	"starter-project/app/models"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/facades"
)

// MessageService - Simplified version using generic CRUD service
type MessageService struct {
	*contracts.GenericCrudService[models.Message]
	sseService   *SSEService
	cacheService *CacheService
}

// NewMessageService creates a new simplified message service
func NewMessageService() *MessageService {
	// Create the generic service
	genericService := contracts.NewGenericCrudService[models.Message]("messages", "id")

	// Configure the service
	genericService.
		SetSearchFields("content", "subject").
		SetSortFields("id", "created_at", "updated_at", "read_at").
		SetFilterFields("sender_id", "recipient_id", "type", "status").
		SetRelations("Sender", "Recipient", "ParentMessage").
		SetValidationRules(map[string]interface{}{
			"sender_id":    "required|numeric",
			"recipient_id": "required|numeric",
			"content":      "required|string|max:5000",
			"type":         "required|string|in:direct,broadcast,system",
			"subject":      "string|max:255",
		}).
		SetBeforeCreate(func(data map[string]interface{}) error {
			// Set defaults
			if _, exists := data["status"]; !exists {
				data["status"] = string(models.MessageStatusSent)
			}
			if _, exists := data["type"]; !exists {
				data["type"] = string(models.MessageTypeDirect)
			}

			// Validate sender and recipient
			senderID, ok := data["sender_id"].(float64)
			if !ok {
				return fmt.Errorf("invalid sender_id")
			}
			recipientID, ok := data["recipient_id"].(float64)
			if !ok {
				return fmt.Errorf("invalid recipient_id")
			}

			// Check if sender exists and is active
			var sender models.User
			if err := facades.Orm().Query().Model(&models.User{}).
				With("Roles").
				Where("id = ? AND is_active = ?", uint(senderID), true).
				First(&sender); err != nil {
				return fmt.Errorf("sender not found or inactive")
			}

			// Check if recipient exists and is active
			var recipient models.User
			if err := facades.Orm().Query().Model(&models.User{}).
				With("Roles").
				Where("id = ? AND is_active = ?", uint(recipientID), true).
				First(&recipient); err != nil {
				return fmt.Errorf("recipient not found or inactive")
			}

			// Check messaging permissions
			if !sender.CanMessageUser(&recipient) {
				return fmt.Errorf("insufficient permissions to message this user")
			}

			// Validate content
			content, ok := data["content"].(string)
			if !ok {
				return fmt.Errorf("content must be a string")
			}
			content = strings.TrimSpace(content)
			if content == "" {
				return fmt.Errorf("message content cannot be empty")
			}
			data["content"] = content

			return nil
		}).
		SetCustomQuery(func(query orm.Query) orm.Query {
			// Always order messages by created_at desc by default
			return query.Order("created_at DESC")
		}).
		SetCustomFilters(func(query orm.Query, filters map[string]interface{}) orm.Query {
			for field, value := range filters {
				switch field {
				case "sender_id", "recipient_id":
					query = query.Where(field+" = ?", value)
				case "type":
					query = query.Where("type = ?", value)
				case "status":
					query = query.Where("status = ?", value)
				case "unread":
					if unread, ok := value.(bool); ok && unread {
						query = query.Where("read_at IS NULL")
					}
				case "thread_id":
					query = query.Where("thread_id = ?", value)
				case "parent_id":
					query = query.Where("parent_id = ?", value)
				}
			}
			return query
		})

	service := &MessageService{
		GenericCrudService: genericService,
		sseService:         NewSSEService(),
		cacheService:       GetCacheService(),
	}

	// Register service
	contracts.MustRegisterCrudService("messages", service)

	return service
}

// Custom methods for messaging functionality

// truncateString truncates a string to maxLen characters with ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// SendMessage creates and sends a new message
func (s *MessageService) SendMessage(senderID uint, recipientID uint, content string, messageType models.MessageType) (*models.Message, error) {
	data := map[string]interface{}{
		"sender_id":    float64(senderID),
		"recipient_id": float64(recipientID),
		"content":      content,
		"type":         string(messageType),
	}

	result, err := s.Create(data)
	if err != nil {
		return nil, err
	}

	message := result.(*models.Message)

	// Invalidate unread count cache for recipient
	s.cacheService.InvalidateMessageUnreadCount(recipientID)

	// Emit SSE events
	s.sseService.NotifyNewMessage(senderID, recipientID, message)

	// Update unread count for recipient (will re-cache)
	if count, err := s.GetUnreadCount(recipientID); err == nil {
		s.sseService.UpdateUnreadCount(recipientID, count)
	}

	// Create notification for the recipient
	s.createMessageNotification(senderID, recipientID, message)

	return message, nil
}

// createMessageNotification creates a notification for a new message
func (s *MessageService) createMessageNotification(senderID uint, recipientID uint, message *models.Message) {
	// Get sender's name for the notification title
	var sender models.User
	if err := facades.Orm().Query().Where("id = ?", senderID).First(&sender); err != nil {
		facades.Log().Warning("Failed to load sender for notification", map[string]interface{}{
			"sender_id": senderID,
			"error":     err.Error(),
		})
		return
	}

	notificationService := NewNotificationService()

	// Create the notification
	title := fmt.Sprintf("New message from %s", sender.Name)
	notificationMessage := truncateString(message.Content, 100)
	notificationType := "message"
	relatedType := "message"
	priority := "normal"

	_, err := notificationService.CreateNotification(
		recipientID,           // userID
		title,                 // title
		notificationMessage,   // message
		notificationType,      // type
		&senderID,             // triggerUserID
		&relatedType,          // relatedType
		&message.ID,           // relatedID
		priority,              // priority
		nil,                   // expiresAt
		"",                    // data
	)

	if err != nil {
		facades.Log().Warning("Failed to create message notification", map[string]interface{}{
			"sender_id":    senderID,
			"recipient_id": recipientID,
			"message_id":   message.ID,
			"error":        err.Error(),
		})
	} else {
		facades.Log().Info("Message notification created", map[string]interface{}{
			"sender_id":    senderID,
			"recipient_id": recipientID,
			"message_id":   message.ID,
		})
	}
}

// SendBroadcast sends a broadcast message to multiple recipients
func (s *MessageService) SendBroadcast(senderID uint, recipientIDs []uint, content string, subject string) ([]*models.Message, error) {
	var messages []*models.Message

	for _, recipientID := range recipientIDs {
		data := map[string]interface{}{
			"sender_id":    float64(senderID),
			"recipient_id": float64(recipientID),
			"content":      content,
			"subject":      subject,
			"type":         string(models.MessageTypeSystem),
		}

		result, err := s.Create(data)
		if err != nil {
			// Log error but continue with other recipients
			facades.Log().Error("Failed to send broadcast to recipient", map[string]interface{}{
				"recipient_id": recipientID,
				"error":        err.Error(),
			})
			continue
		}

		message := result.(*models.Message)
		messages = append(messages, message)

		// Create notification for each recipient
		s.createMessageNotification(senderID, recipientID, message)
	}

	return messages, nil
}

// BroadcastToRole sends a message to all users with a specific role (super admin only)
func (s *MessageService) BroadcastToRole(senderID uint, roleID uint, content string, subject string) (map[string]interface{}, error) {
	// Get the role to verify it exists
	var role models.Role
	if err := facades.Orm().Query().Where("id = ?", roleID).First(&role); err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}

	// Get all active users with this role
	var userIDs []uint
	err := facades.Orm().Query().Raw(`
		SELECT DISTINCT u.id
		FROM users u
		INNER JOIN user_roles ur ON ur.user_id = u.id
		WHERE ur.role_id = ?
		AND ur.is_active = true
		AND u.is_active = true
		AND u.deleted_at IS NULL
		AND u.id != ?
	`, roleID, senderID).Pluck("id", &userIDs)

	if err != nil {
		return nil, fmt.Errorf("failed to get users for role: %w", err)
	}

	if len(userIDs) == 0 {
		return map[string]interface{}{
			"sent_count":   0,
			"role_id":      roleID,
			"role_name":    role.Name,
			"failed_count": 0,
			"message":      "No active users found with this role",
		}, nil
	}

	// Send messages to all users
	var sentCount int
	var failedCount int
	var messageIDs []uint

	for _, recipientID := range userIDs {
		data := map[string]interface{}{
			"sender_id":    float64(senderID),
			"recipient_id": float64(recipientID),
			"content":      content,
			"type":         string(models.MessageTypeSystem),
		}

		result, err := s.Create(data)
		if err != nil {
			facades.Log().Error("Failed to send role broadcast to user", map[string]interface{}{
				"recipient_id": recipientID,
				"role_id":      roleID,
				"error":        err.Error(),
			})
			failedCount++
			continue
		}

		message := result.(*models.Message)
		messageIDs = append(messageIDs, message.ID)
		sentCount++

		// Create notification for the recipient
		s.createMessageNotification(senderID, recipientID, message)
	}

	facades.Log().Info("Role broadcast completed", map[string]interface{}{
		"sender_id":    senderID,
		"role_id":      roleID,
		"role_name":    role.Name,
		"sent_count":   sentCount,
		"failed_count": failedCount,
	})

	return map[string]interface{}{
		"sent_count":   sentCount,
		"failed_count": failedCount,
		"role_id":      roleID,
		"role_name":    role.Name,
		"message_ids":  messageIDs,
	}, nil
}

// GetBroadcastHistory retrieves broadcast messages sent by a user (super admin only)
// Groups broadcasts by timestamp to show as threads
func (s *MessageService) GetBroadcastHistory(senderID uint, req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	// Get page parameters
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	// Get broadcast messages (type = 'system') sent by this user, grouped by content and approximate time
	// We'll use a subquery to get distinct broadcasts
	var broadcasts []struct {
		Content     string    `gorm:"column:content"`
		CreatedAt   time.Time `gorm:"column:created_at"`
		FirstMsgID  uint      `gorm:"column:first_msg_id"`
		RecipientCount int64  `gorm:"column:recipient_count"`
	}

	// Query to get unique broadcasts (grouped by content and minute)
	err := facades.Orm().Query().Raw(`
		SELECT
			content,
			MIN(created_at) as created_at,
			MIN(id) as first_msg_id,
			COUNT(*) as recipient_count
		FROM messages
		WHERE sender_id = ?
		AND type = 'system'
		AND deleted_at IS NULL
		GROUP BY content, DATE_TRUNC('minute', created_at)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, senderID, pageSize, offset).Scan(&broadcasts)

	if err != nil {
		return nil, fmt.Errorf("failed to get broadcast history: %w", err)
	}

	// Get total count
	var totalCount int64
	err = facades.Orm().Query().Raw(`
		SELECT COUNT(*) FROM (
			SELECT 1
			FROM messages
			WHERE sender_id = ?
			AND type = 'system'
			AND deleted_at IS NULL
			GROUP BY content, DATE_TRUNC('minute', created_at)
		) as broadcast_groups
	`, senderID).Scan(&totalCount)

	if err != nil {
		return nil, fmt.Errorf("failed to count broadcasts: %w", err)
	}

	// Build result with read counts for each broadcast
	var results []interface{}
	for _, broadcast := range broadcasts {
		// Calculate time window for this broadcast
		startTime := broadcast.CreatedAt.Add(-1 * time.Minute)
		endTime := broadcast.CreatedAt.Add(1 * time.Minute)

		// Get read count for this broadcast
		var readCount int64
		facades.Orm().Query().Raw(`
			SELECT COUNT(*)
			FROM messages m
			WHERE m.sender_id = ?
			AND m.type = 'system'
			AND m.content = ?
			AND m.created_at >= ?
			AND m.created_at <= ?
			AND m.deleted_at IS NULL
			AND m.read_at IS NOT NULL
		`, senderID, broadcast.Content, startTime, endTime).Scan(&readCount)

		results = append(results, map[string]interface{}{
			"id":              broadcast.FirstMsgID,
			"content":         broadcast.Content,
			"created_at":      broadcast.CreatedAt,
			"recipient_count": broadcast.RecipientCount,
			"read_count":      readCount,
		})
	}

	lastPage := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		lastPage++
	}
	if lastPage < 1 {
		lastPage = 1
	}

	return &contracts.PaginatedResult{
		Data:        results,
		Total:       totalCount,
		PerPage:     pageSize,
		CurrentPage: page,
		LastPage:    lastPage,
	}, nil
}

// GetInbox retrieves messages for a user with pagination
func (s *MessageService) GetInbox(userID uint, req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	filters := map[string]interface{}{
		"recipient_id": userID,
	}
	return s.GetListAdvanced(req, filters)
}

// GetSentMessages retrieves sent messages for a user
func (s *MessageService) GetSentMessages(userID uint, req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	filters := map[string]interface{}{
		"sender_id": userID,
	}
	return s.GetListAdvanced(req, filters)
}

// GetUnreadMessages retrieves unread messages for a user
func (s *MessageService) GetUnreadMessages(userID uint, req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	filters := map[string]interface{}{
		"recipient_id": userID,
		"unread":       true,
	}
	return s.GetListAdvanced(req, filters)
}

// GetUnreadCount gets the count of unread messages
func (s *MessageService) GetUnreadCount(userID uint) (int64, error) {
	// Try to get from cache first
	if cachedCount, found := s.cacheService.GetMessageUnreadCount(userID); found {
		facades.Log().Debug("Message unread count cache hit", map[string]interface{}{
			"user_id": userID,
			"count":   cachedCount,
		})
		return cachedCount, nil
	}

	// Query from database
	count, err := facades.Orm().Query().Model(&models.Message{}).
		Where("recipient_id = ? AND read_at IS NULL", userID).
		Count()
	if err != nil {
		return 0, err
	}

	// Cache the count
	if err := s.cacheService.SetMessageUnreadCount(userID, count); err != nil {
		facades.Log().Warning("Failed to cache message unread count", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
	}

	return count, nil
}

// GetThreadMessages retrieves all messages in a thread
func (s *MessageService) GetThreadMessages(threadID uint, req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	filters := map[string]interface{}{
		"thread_id": threadID,
	}
	return s.GetListAdvanced(req, filters)
}

// MarkAsRead marks a message as read
func (s *MessageService) MarkAsRead(messageID uint, userID uint) error {
	// Verify the user is the recipient
	result, err := s.GetByID(messageID)
	if err != nil {
		return err
	}

	message := result.(*models.Message)
	if message.RecipientID == nil || *message.RecipientID != userID {
		return fmt.Errorf("unauthorized to mark this message as read")
	}

	// Check if already read by checking ReadAt
	if message.ReadAt != nil {
		return nil // Already read
	}

	_, err = s.Update(messageID, map[string]interface{}{
		"status":  string(models.MessageStatusRead),
		"read_at": time.Now(),
	})

	if err == nil {
		// Invalidate unread count cache
		s.cacheService.InvalidateMessageUnreadCount(userID)

		// Emit SSE events
		s.sseService.NotifyMessageRead(message.SenderID, messageID)

		// Update unread count for recipient (will re-cache)
		if count, err := s.GetUnreadCount(userID); err == nil {
			s.sseService.UpdateUnreadCount(userID, count)
		}
	}

	return err
}

// MarkConversationAsRead marks all messages from a specific user as read
func (s *MessageService) MarkConversationAsRead(currentUserID uint, otherUserID uint) error {
	now := time.Now()

	// Update all unread messages from otherUserID to currentUserID
	_, err := facades.Orm().Query().
		Model(&models.Message{}).
		Where("sender_id = ?", otherUserID).
		Where("recipient_id = ?", currentUserID).
		Where("read_at IS NULL").
		Update("read_at", now)

	if err != nil {
		return fmt.Errorf("failed to mark messages as read: %w", err)
	}

	// Also update status
	_, err = facades.Orm().Query().
		Model(&models.Message{}).
		Where("sender_id = ?", otherUserID).
		Where("recipient_id = ?", currentUserID).
		Where("status != ?", string(models.MessageStatusRead)).
		Update("status", string(models.MessageStatusRead))

	if err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}

	// Invalidate unread count cache
	s.cacheService.InvalidateMessageUnreadCount(currentUserID)

	// Update unread count for recipient (will re-cache)
	if count, err := s.GetUnreadCount(currentUserID); err == nil {
		s.sseService.UpdateUnreadCount(currentUserID, count)
	}

	return nil
}

// MarkAsImportant marks a message as important/unimportant
// Note: This would need a schema change to add is_important field to messages table
func (s *MessageService) MarkAsImportant(messageID uint, userID uint, important bool) error {
	// Verify the user is either sender or recipient
	result, err := s.GetByID(messageID)
	if err != nil {
		return err
	}

	message := result.(*models.Message)
	if message.SenderID != userID && (message.RecipientID == nil || *message.RecipientID != userID) {
		return fmt.Errorf("unauthorized to modify this message")
	}

	// TODO: Add is_important field to Message model and migration
	return fmt.Errorf("marking messages as important is not yet implemented")
}

// DeleteMessage soft deletes a message (marks it as deleted for a user)
func (s *MessageService) DeleteMessage(messageID uint, userID uint) error {
	// Verify the user is either sender or recipient
	result, err := s.GetByID(messageID)
	if err != nil {
		return err
	}

	message := result.(*models.Message)
	if message.SenderID != userID && (message.RecipientID == nil || *message.RecipientID != userID) {
		return fmt.Errorf("unauthorized to delete this message")
	}

	// For now, we'll use the generic delete which soft deletes the entire message
	// In a real app, you might want to track deletion per user
	err = s.Delete(messageID)

	if err == nil {
		// Emit SSE events to both sender and recipient
		s.sseService.NotifyMessageDeleted(message.SenderID, messageID)
		if message.RecipientID != nil {
			s.sseService.NotifyMessageDeleted(*message.RecipientID, messageID)

			// Update unread count if message was unread
			if message.ReadAt == nil {
				// Invalidate unread count cache
				s.cacheService.InvalidateMessageUnreadCount(*message.RecipientID)

				if count, err := s.GetUnreadCount(*message.RecipientID); err == nil {
					s.sseService.UpdateUnreadCount(*message.RecipientID, count)
				}
			}
		}
	}

	return err
}

// GetConversations retrieves a list of conversations for a user
func (s *MessageService) GetConversations(userID uint) ([]interface{}, error) {
	// First, get all unique users this user has exchanged messages with
	var conversations []struct {
		UserID       uint      `json:"user_id"`
		LastActivity time.Time `json:"last_activity"`
	}

	// Query to get unique conversation partners
	query := `
		SELECT DISTINCT 
			CASE 
				WHEN sender_id = ? THEN recipient_id 
				ELSE sender_id 
			END as user_id,
			MAX(created_at) as last_activity
		FROM messages 
		WHERE (sender_id = ? OR recipient_id = ?) 
			AND deleted_at IS NULL
		GROUP BY user_id
		ORDER BY last_activity DESC
	`

	if err := facades.Orm().Query().Raw(query, userID, userID, userID).Scan(&conversations); err != nil {
		return nil, err
	}

	// Now build the conversation objects
	result := make([]interface{}, 0, len(conversations))

	for _, conv := range conversations {
		// Get the user details
		var user models.User
		if err := facades.Orm().Query().Model(&models.User{}).
			With("Roles").
			Where("id = ?", conv.UserID).
			First(&user); err != nil {
			continue // Skip if user not found
		}

		// Get the latest message between the two users
		var latestMessage models.Message
		if err := facades.Orm().Query().Model(&models.Message{}).
			Where("(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)",
				userID, conv.UserID, conv.UserID, userID).
			Order("created_at DESC").
			First(&latestMessage); err != nil {
			continue // Skip if no messages found
		}

		// Get unread count
		unreadCount, _ := facades.Orm().Query().Model(&models.Message{}).
			Where("sender_id = ? AND recipient_id = ? AND status = ? AND deleted_at IS NULL",
				conv.UserID, userID, models.MessageStatusSent).
			Count()

		// Build conversation object
		conversation := map[string]interface{}{
			"user": map[string]interface{}{
				"id":        user.ID,
				"name":      user.Name,
				"email":     user.Email,
				"is_active": user.IsActive,
				"roles":     user.Roles,
			},
			"latest_message": map[string]interface{}{
				"id":         latestMessage.ID,
				"content":    latestMessage.Content,
				"created_at": latestMessage.CreatedAt,
				"sender_id":  latestMessage.SenderID,
				"is_edited":  latestMessage.IsEdited,
			},
			"unread_count":  unreadCount,
			"last_activity": conv.LastActivity,
		}

		result = append(result, conversation)
	}

	return result, nil
}

// GetConversation retrieves messages between two users
func (s *MessageService) GetConversation(user1ID uint, user2ID uint, req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	// Get total count first
	total, countErr := facades.Orm().Query().Model(&models.Message{}).
		Where("((sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)) AND deleted_at IS NULL",
			user1ID, user2ID, user2ID, user1ID).
		Count()
	if countErr != nil {
		return nil, countErr
	}

	// Get paginated messages
	var messages []models.Message
	offset := (req.Page - 1) * req.PageSize
	err := facades.Orm().Query().
		With("Sender").
		With("Recipient").
		Where("((sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)) AND deleted_at IS NULL",
			user1ID, user2ID, user2ID, user1ID).
		Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&messages)

	if err != nil {
		return nil, err
	}

	// Convert to interface slice
	data := make([]interface{}, len(messages))
	for i, msg := range messages {
		data[i] = msg
	}

	// Use pagination utility
	return &contracts.PaginatedResult{
		Data:        data,
		Total:       total,
		PerPage:     req.PageSize,
		CurrentPage: req.Page,
		LastPage:    int((total + int64(req.PageSize) - 1) / int64(req.PageSize)),
		From:        offset + 1,
		To:          offset + len(messages),
		HasNext:     req.Page < int((total+int64(req.PageSize)-1)/int64(req.PageSize)),
		HasPrev:     req.Page > 1,
	}, nil
}

// GetColumnMapping returns database column mappings
func (s *MessageService) GetColumnMapping() map[string]string {
	return map[string]string{
		"id":          "id",
		"senderId":    "sender_id",
		"recipientId": "recipient_id",
		"content":     "content",
		"type":        "type",
		"subject":     "subject",
		"threadId":    "thread_id",
		"parentId":    "parent_id",
		"status":      "status",
		"readAt":      "read_at",
		"createdAt":   "created_at",
		"updatedAt":   "updated_at",
	}
}
