package controllers

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"players/app/auth"
	"players/app/contracts"
	"players/app/models"
	"players/app/services"
)

type MessageController struct {
	*contracts.BaseCrudController
	messageService *services.MessageService
}

func NewMessageController() *MessageController {
	return &MessageController{
		BaseCrudController: contracts.NewBaseCrudController("message"),
		messageService:     services.NewMessageService(),
	}
}

// SendMessage handles POST /api/messages
func (c *MessageController) SendMessage(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Parse request
	var request struct {
		RecipientID uint               `json:"recipient_id" validate:"required"`
		Content     string             `json:"content" validate:"required,max:5000"`
		Type        models.MessageType `json:"type"`
	}

	if err := ctx.Request().Bind(&request); err != nil {
		return c.BadRequestResponse(ctx, "Invalid request data", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Set default message type
	if request.Type == "" {
		request.Type = models.MessageTypeDirect
	}

	// Send message
	message, err := c.messageService.SendMessage(user.ID, request.RecipientID, request.Content, request.Type)
	if err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	return c.CreatedResponse(ctx, message, "Message sent successfully")
}

// GetConversation handles GET /api/messages/conversation/{userId}
func (c *MessageController) GetConversation(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Parse user ID from route
	otherUserIDStr := ctx.Request().Route("userId")
	otherUserID, err := strconv.ParseUint(otherUserIDStr, 10, 32)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid user ID", nil)
	}

	// Validate pagination request
	req, err := c.ValidatePaginationRequest(ctx)
	if err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	// Get conversation
	result, err := c.messageService.GetConversation(user.ID, uint(otherUserID), *req)
	if err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	return c.SuccessResponse(ctx, c.BuildPaginatedResponse(result, req), "Conversation retrieved successfully")
}

// GetConversations handles GET /api/messages/conversations
func (c *MessageController) GetConversations(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Validate pagination request
	req, err := c.ValidatePaginationRequest(ctx)
	if err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	// Get conversations
	result, err := c.messageService.GetUserConversations(user.ID, *req)
	if err != nil {
		return c.InternalErrorResponse(ctx, err.Error())
	}

	return c.SuccessResponse(ctx, c.BuildPaginatedResponse(result, req), "Conversations retrieved successfully")
}

// GetMessagableUsers handles GET /api/messages/users
func (c *MessageController) GetMessagableUsers(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Validate pagination request
	req, err := c.ValidatePaginationRequest(ctx)
	if err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	// Get messagable users
	result, err := c.messageService.GetMessagableUsers(user.ID, *req)
	if err != nil {
		return c.InternalErrorResponse(ctx, err.Error())
	}

	return c.SuccessResponse(ctx, c.BuildPaginatedResponse(result, req), "Users retrieved successfully")
}

// MarkAsRead handles PUT /api/messages/conversation/{userId}/read
func (c *MessageController) MarkAsRead(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Parse user ID from route
	senderIDStr := ctx.Request().Route("userId")
	senderID, err := strconv.ParseUint(senderIDStr, 10, 32)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid user ID", nil)
	}

	// Mark messages as read
	if err := c.messageService.MarkMessagesAsRead(user.ID, uint(senderID)); err != nil {
		return c.InternalErrorResponse(ctx, err.Error())
	}

	return c.SuccessResponse(ctx, nil, "Messages marked as read")
}

// DeleteMessage handles DELETE /api/messages/{id}
func (c *MessageController) DeleteMessage(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Parse message ID from route
	messageID, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	// Delete message
	if err := c.messageService.DeleteMessage(messageID, user.ID); err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	return c.SuccessResponse(ctx, nil, "Message deleted successfully")
}

// EditMessage handles PUT /api/messages/{id}
func (c *MessageController) EditMessage(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Parse message ID from route
	messageID, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	// Parse request
	var request struct {
		Content string `json:"content" validate:"required,max:5000"`
	}

	if err := ctx.Request().Bind(&request); err != nil {
		return c.BadRequestResponse(ctx, "Invalid request data", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Edit message
	message, err := c.messageService.EditMessage(messageID, user.ID, request.Content)
	if err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	return c.SuccessResponse(ctx, message, "Message updated successfully")
}

// GetUnreadCount handles GET /api/messages/unread-count
func (c *MessageController) GetUnreadCount(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Get unread count
	count, err := c.messageService.GetUnreadMessageCount(user.ID)
	if err != nil {
		return c.InternalErrorResponse(ctx, err.Error())
	}

	return c.SuccessResponse(ctx, map[string]interface{}{
		"unread_count": count,
	}, "Unread count retrieved successfully")
}

// SearchUsers handles GET /api/messages/search-users
func (c *MessageController) SearchUsers(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Get search query
	query := ctx.Request().Query("q", "")
	if query == "" {
		return c.BadRequestResponse(ctx, "Search query is required", nil)
	}

	// Create search request
	req := contracts.ListRequest{
		Page:     1,
		PageSize: 10,
		Search:   query,
		Sort:     "name",
		Direction: "ASC",
	}

	// Get messagable users
	result, err := c.messageService.GetMessagableUsers(user.ID, req)
	if err != nil {
		return c.InternalErrorResponse(ctx, err.Error())
	}

	return c.SuccessResponse(ctx, result.Data, "Users found")
}

// GetMessage handles GET /api/messages/{id}
func (c *MessageController) GetMessage(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Parse message ID from route
	messageID, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	// Get message
	var message models.Message
	if err := facades.Orm().Query().Model(&models.Message{}).
		With("Sender").
		With("Recipient").
		With("ParentMessage").
		Where("id = ?", messageID).
		First(&message); err != nil {
		return c.NotFoundResponse(ctx, "Message not found")
	}

	// Check if user can read this message
	if !message.CanBeReadBy(user) {
		return c.ForbiddenResponse(ctx, "Insufficient permissions to view this message")
	}

	return c.SuccessResponse(ctx, message, "Message retrieved successfully")
}

// ReplyToMessage handles POST /api/messages/{id}/reply
func (c *MessageController) ReplyToMessage(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Parse message ID from route
	parentMessageID, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	// Get parent message to determine recipient
	var parentMessage models.Message
	if err := facades.Orm().Query().Model(&models.Message{}).
		With("Sender").
		Where("id = ?", parentMessageID).
		First(&parentMessage); err != nil {
		return c.NotFoundResponse(ctx, "Parent message not found")
	}

	// Parse request
	var request struct {
		Content string `json:"content" validate:"required,max:5000"`
	}

	if err := ctx.Request().Bind(&request); err != nil {
		return c.BadRequestResponse(ctx, "Invalid request data", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Determine recipient (reply to sender of parent message)
	var recipientID uint
	if parentMessage.SenderID == user.ID {
		// If replying to own message, need recipient ID in request
		var recipRequest struct {
			RecipientID uint `json:"recipient_id" validate:"required"`
		}
		if err := ctx.Request().Bind(&recipRequest); err != nil {
			return c.BadRequestResponse(ctx, "Recipient ID required when replying to own message", nil)
		}
		recipientID = recipRequest.RecipientID
	} else {
		recipientID = parentMessage.SenderID
	}

	// Create reply message
	message, err := c.messageService.SendMessage(user.ID, recipientID, request.Content, models.MessageTypeDirect)
	if err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	// Update reply to reference parent message
	message.ParentMessageID = &parentMessageID
	if err := facades.Orm().Query().Save(message); err != nil {
		// Log error but don't fail the response
		facades.Log().Warning("Failed to update parent message reference: %v", err)
	}

	return c.CreatedResponse(ctx, message, "Reply sent successfully")
}