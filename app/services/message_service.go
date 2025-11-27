package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/models"
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
	genericService := contracts.NewGenericCrudService[models.Message]("message", "id")

	// Configure the service
	genericService.
		SetSearchFields("content", "subject").
		SetSortFields("id", "created_at", "updated_at", "read_at").
		SetFilterFields("sender_id", "recipient_id", "type", "status").
		SetRelations("Sender", "Recipient", "Thread", "Parent").
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

	return message, nil
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

		messages = append(messages, result.(*models.Message))
	}

	return messages, nil
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
	// Create a custom query for conversation
	var messages []models.Message
	query := facades.Orm().Query().Model(&models.Message{}).
		With("Sender", "Recipient").
		Where("(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)",
			user1ID, user2ID, user2ID, user1ID).
		Order("created_at DESC")

	// Apply pagination manually
	total, err := query.Count()
	if err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&messages); err != nil {
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
