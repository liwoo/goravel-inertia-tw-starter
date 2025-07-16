package services

import (
	"fmt"
	"strings"
	"time"
	"regexp"

	"github.com/goravel/framework/facades"
	"players/app/contracts"
	"players/app/models"
)

type MessageService struct {
	*contracts.BaseCrudService
}

func NewMessageService() *MessageService {
	return &MessageService{
		BaseCrudService: contracts.NewBaseCrudService("message", "id"),
	}
}

// SendMessage creates and sends a new message
func (s *MessageService) SendMessage(senderID uint, recipientID uint, content string, messageType models.MessageType) (*models.Message, error) {
	// Validate sender exists and get with roles
	var sender models.User
	if err := facades.Orm().Query().Model(&models.User{}).With("Roles").Where("id = ? AND is_active = ?", senderID, true).First(&sender); err != nil {
		return nil, fmt.Errorf("sender not found or inactive")
	}

	// Validate recipient exists and get with roles
	var recipient models.User
	if err := facades.Orm().Query().Model(&models.User{}).With("Roles").Where("id = ? AND is_active = ?", recipientID, true).First(&recipient); err != nil {
		return nil, fmt.Errorf("recipient not found or inactive")
	}

	// Check messaging permissions
	if !sender.CanMessageUser(&recipient) {
		return nil, fmt.Errorf("insufficient permissions to message this user")
	}

	// Validate content
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("message content cannot be empty")
	}
	if len(content) > 5000 {
		return nil, fmt.Errorf("message content too long (max 5000 characters)")
	}

	// Create message
	message := &models.Message{
		Content:     content,
		Type:        messageType,
		Status:      models.MessageStatusSent,
		SenderID:    senderID,
		RecipientID: &recipientID,
	}

	if err := facades.Orm().Query().Create(message); err != nil {
		return nil, fmt.Errorf("failed to create message: %v", err)
	}

	// Process mentions
	if err := s.processMentions(message, content); err != nil {
		// Log error but don't fail the message send
		facades.Log().Warning("Failed to process mentions for message %d: %v", message.ID, err)
	}

	// Mark as delivered immediately (for now)
	message.MarkAsDelivered()
	facades.Orm().Query().Save(message)

	// Create notification for recipient
	s.createMessageNotification(message, &recipient)

	// Load relations for response
	facades.Orm().Query().Model(&models.Message{}).With("Sender").With("Recipient").Where("id = ?", message.ID).First(message)

	return message, nil
}

// GetConversation retrieves messages between two users
func (s *MessageService) GetConversation(userID, otherUserID uint, request contracts.ListRequest) (*contracts.PaginatedResult, error) {
	// Validate both users exist and have proper permissions
	var user, otherUser models.User
	if err := facades.Orm().Query().Model(&models.User{}).With("Roles").Where("id = ? AND is_active = ?", userID, true).First(&user); err != nil {
		return nil, fmt.Errorf("user not found or inactive")
	}
	if err := facades.Orm().Query().Model(&models.User{}).With("Roles").Where("id = ? AND is_active = ?", otherUserID, true).First(&otherUser); err != nil {
		return nil, fmt.Errorf("other user not found or inactive")
	}

	// Check if user can view conversation with other user
	if !user.CanMessageUser(&otherUser) {
		return nil, fmt.Errorf("insufficient permissions to view conversation")
	}

	// Build query for conversation messages
	query := facades.Orm().Query().Model(&models.Message{}).
		With("Sender").
		With("Recipient").
		Where("(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)", 
			userID, otherUserID, otherUserID, userID).
		Where("status != ?", models.MessageStatusDeleted)

	// Apply search if provided
	if request.Search != "" {
		searchPattern := "%" + request.Search + "%"
		query = query.Where("content LIKE ?", searchPattern)
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
	query.Model(&models.Message{}).Count(&total)

	// Apply pagination
	offset := (request.Page - 1) * request.PageSize
	var messages []models.Message
	if err := query.Offset(offset).Limit(request.PageSize).Find(&messages); err != nil {
		return nil, fmt.Errorf("failed to retrieve messages: %v", err)
	}

	// Convert to interface slice
	data := make([]interface{}, len(messages))
	for i, message := range messages {
		data[i] = message
	}

	// Calculate pagination metadata
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
		To:          offset + len(messages),
		HasNext:     request.Page < lastPage,
		HasPrev:     request.Page > 1,
	}, nil
}

// GetUserConversations retrieves all conversations for a user
func (s *MessageService) GetUserConversations(userID uint, request contracts.ListRequest) (*contracts.PaginatedResult, error) {
	// Get all messages for the user to build conversations
	var allMessages []models.Message
	if err := facades.Orm().Query().Model(&models.Message{}).
		With("Sender").
		With("Recipient").
		Where("(sender_id = ? OR recipient_id = ?) AND status != ?", userID, userID, models.MessageStatusDeleted).
		Order("created_at DESC").
		Find(&allMessages); err != nil {
		return nil, fmt.Errorf("failed to retrieve messages: %v", err)
	}

	// Group messages by conversation partner
	conversationMap := make(map[uint]*models.Message)
	unreadCounts := make(map[uint]int64)

	for _, message := range allMessages {
		var otherUserID uint
		if message.SenderID == userID {
			if message.RecipientID != nil {
				otherUserID = *message.RecipientID
			} else {
				continue
			}
		} else {
			otherUserID = message.SenderID
		}

		// Keep track of latest message per conversation
		if _, exists := conversationMap[otherUserID]; !exists {
			conversationMap[otherUserID] = &message
		}

		// Count unread messages (messages sent TO the current user that are unread)
		if message.SenderID == otherUserID && message.RecipientID != nil && *message.RecipientID == userID {
			if message.Status != models.MessageStatusRead && message.ReadAt == nil {
				unreadCounts[otherUserID]++
			}
		}
	}

	// Convert map to slice and get user details
	conversations := make([]interface{}, 0, len(conversationMap))
	for otherUserID, latestMessage := range conversationMap {
		// Get other user details
		var otherUser models.User
		if err := facades.Orm().Query().Model(&models.User{}).With("Roles").Where("id = ?", otherUserID).First(&otherUser); err != nil {
			continue
		}

		conversations = append(conversations, map[string]interface{}{
			"user":           otherUser,
			"latest_message": *latestMessage,
			"unread_count":   unreadCounts[otherUserID],
			"last_activity":  latestMessage.CreatedAt.StdTime().Format(time.RFC3339),
		})
	}

	// Sort conversations by last activity (most recent first)
	// Since we're already getting messages ordered by created_at DESC, 
	// the first message we encounter for each user is the latest

	// Apply pagination
	total := int64(len(conversations))
	offset := (request.Page - 1) * request.PageSize
	end := offset + request.PageSize

	if offset > len(conversations) {
		offset = len(conversations)
	}
	if end > len(conversations) {
		end = len(conversations)
	}

	var pageConversations []interface{}
	if offset < len(conversations) {
		pageConversations = conversations[offset:end]
	}

	lastPage := int((total + int64(request.PageSize) - 1) / int64(request.PageSize))
	if lastPage < 1 {
		lastPage = 1
	}

	return &contracts.PaginatedResult{
		Data:        pageConversations,
		Total:       total,
		CurrentPage: request.Page,
		LastPage:    lastPage,
		PerPage:     request.PageSize,
		From:        offset + 1,
		To:          offset + len(pageConversations),
		HasNext:     request.Page < lastPage,
		HasPrev:     request.Page > 1,
	}, nil
}

// GetMessagableUsers returns users that the current user can message
func (s *MessageService) GetMessagableUsers(currentUserID uint, request contracts.ListRequest) (*contracts.PaginatedResult, error) {
	// Get current user with roles
	var currentUser models.User
	if err := facades.Orm().Query().Model(&models.User{}).With("Roles").Where("id = ? AND is_active = ?", currentUserID, true).First(&currentUser); err != nil {
		return nil, fmt.Errorf("current user not found or inactive")
	}

	query := facades.Orm().Query().Model(&models.User{}).
		With("Roles").
		Where("id != ? AND is_active = ?", currentUserID, true)

	// If not super admin, filter by shared roles
	if !currentUser.IsSuperAdminUser() {
		// Get current user's role IDs
		roleIDs := make([]uint, 0, len(currentUser.GetActiveRoles()))
		for _, role := range currentUser.GetActiveRoles() {
			roleIDs = append(roleIDs, role.ID)
		}

		if len(roleIDs) > 0 {
			query = query.
				Where("EXISTS (SELECT 1 FROM user_roles WHERE user_roles.user_id = users.id AND user_roles.role_id IN ?)", roleIDs)
		} else {
			// User has no roles, can't message anyone
			return &contracts.PaginatedResult{
				Data:        []interface{}{},
				Total:       0,
				CurrentPage: request.Page,
				LastPage:    1,
				PerPage:     request.PageSize,
			}, nil
		}
	}

	// Apply search
	if request.Search != "" {
		searchPattern := "%" + request.Search + "%"
		query = query.Where("name LIKE ? OR email LIKE ?", searchPattern, searchPattern)
	}

	// Apply sorting
	if request.Sort == "" {
		request.Sort = "name"
		request.Direction = "ASC"
	}
	orderClause := fmt.Sprintf("%s %s", request.Sort, request.Direction)
	query = query.Order(orderClause)

	// Get total count
	var total int64
	query.Model(&models.User{}).Count(&total)

	// Apply pagination
	offset := (request.Page - 1) * request.PageSize
	var users []models.User
	if err := query.Offset(offset).Limit(request.PageSize).Find(&users); err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %v", err)
	}

	// Convert to interface slice
	data := make([]interface{}, len(users))
	for i, user := range users {
		data[i] = user
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
		To:          offset + len(users),
		HasNext:     request.Page < lastPage,
		HasPrev:     request.Page > 1,
	}, nil
}

// MarkMessagesAsRead marks messages in a conversation as read
func (s *MessageService) MarkMessagesAsRead(userID, senderID uint) error {
	_, err := facades.Orm().Query().
		Model(&models.Message{}).
		Where("sender_id = ? AND recipient_id = ? AND status != ?", senderID, userID, models.MessageStatusRead).
		Update(map[string]interface{}{
			"status":  models.MessageStatusRead,
			"read_at": time.Now(),
		})
	return err
}

// DeleteMessage soft deletes a message (only sender can delete)
func (s *MessageService) DeleteMessage(messageID, userID uint) error {
	var message models.Message
	if err := facades.Orm().Query().Where("id = ?", messageID).First(&message); err != nil {
		return fmt.Errorf("message not found")
	}

	if message.SenderID != userID {
		return fmt.Errorf("unauthorized: only sender can delete message")
	}

	message.Status = models.MessageStatusDeleted
	return facades.Orm().Query().Save(&message)
}

// EditMessage edits a message (only sender within time limit)
func (s *MessageService) EditMessage(messageID, userID uint, newContent string) (*models.Message, error) {
	var message models.Message
	if err := facades.Orm().Query().Where("id = ?", messageID).First(&message); err != nil {
		return nil, fmt.Errorf("message not found")
	}

	// Get user for permission check
	var user models.User
	if err := facades.Orm().Query().Where("id = ?", userID).First(&user); err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if !message.CanEditMessage(&user) {
		return nil, fmt.Errorf("unauthorized: cannot edit this message")
	}

	// Validate new content
	newContent = strings.TrimSpace(newContent)
	if newContent == "" {
		return nil, fmt.Errorf("message content cannot be empty")
	}
	if len(newContent) > 5000 {
		return nil, fmt.Errorf("message content too long (max 5000 characters)")
	}

	// Update message
	now := time.Now()
	message.Content = newContent
	message.IsEdited = true
	message.EditedAt = &now

	if err := facades.Orm().Query().Save(&message); err != nil {
		return nil, fmt.Errorf("failed to update message: %v", err)
	}

	// Load relations for response
	facades.Orm().Query().Model(&models.Message{}).With("Sender").With("Recipient").Where("id = ?", message.ID).First(&message)

	return &message, nil
}

// processMentions extracts and creates mentions from message content
func (s *MessageService) processMentions(message *models.Message, content string) error {
	// Extract @mentions using regex
	mentionRegex := regexp.MustCompile(`@(\w+)`)
	matches := mentionRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		username := match[1]
		
		// Find user by name (you might want to use username field instead)
		var mentionedUser models.User
		if err := facades.Orm().Query().Where("name = ? AND is_active = ?", username, true).First(&mentionedUser); err != nil {
			continue // User not found, skip mention
		}

		// Create mention record
		mention := &models.MessageMention{
			MessageID: message.ID,
			UserID:    mentionedUser.ID,
			Position:  strings.Index(content, "@"+username),
			Length:    len("@" + username),
		}

		if err := facades.Orm().Query().Create(mention); err != nil {
			facades.Log().Warning("Failed to create mention for user %d in message %d: %v", mentionedUser.ID, message.ID, err)
		}

		// Create notification for mentioned user
		s.createMentionNotification(message, &mentionedUser)
	}

	return nil
}

// createMessageNotification creates a notification for new messages
func (s *MessageService) createMessageNotification(message *models.Message, recipient *models.User) {
	notification := &models.Notification{
		Title:         fmt.Sprintf("New message from %s", message.Sender.Name),
		Message:       s.truncateContent(message.Content, 100),
		Type:          "message",
		UserID:        recipient.ID,
		TriggerUserID: &message.SenderID,
		RelatedType:   "message",
		RelatedID:     &message.ID,
		Priority:      "normal",
	}

	if err := facades.Orm().Query().Create(notification); err != nil {
		facades.Log().Warning("Failed to create message notification: %v", err)
	}
}

// createMentionNotification creates a notification for mentions
func (s *MessageService) createMentionNotification(message *models.Message, mentionedUser *models.User) {
	notification := &models.Notification{
		Title:         fmt.Sprintf("You were mentioned by %s", message.Sender.Name),
		Message:       s.truncateContent(message.Content, 100),
		Type:          "mention",
		UserID:        mentionedUser.ID,
		TriggerUserID: &message.SenderID,
		RelatedType:   "message",
		RelatedID:     &message.ID,
		Priority:      "high",
	}

	if err := facades.Orm().Query().Create(notification); err != nil {
		facades.Log().Warning("Failed to create mention notification: %v", err)
	}
}

// truncateContent truncates content to specified length with ellipsis
func (s *MessageService) truncateContent(content string, maxLength int) string {
	if len(content) <= maxLength {
		return content
	}
	return content[:maxLength-3] + "..."
}

// GetUnreadMessageCount returns unread message count for a user
func (s *MessageService) GetUnreadMessageCount(userID uint) (int64, error) {
	var count int64
	err := facades.Orm().Query().Model(&models.Message{}).
		Where("recipient_id = ? AND status != ? AND (read_at IS NULL OR status != ?)", 
			userID, models.MessageStatusDeleted, models.MessageStatusRead).
		Count(&count)
	return count, err
}