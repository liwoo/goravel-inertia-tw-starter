package messages

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"players/app/auth"
	"players/app/contracts"
	"players/app/http/requests"
	"players/app/models"
	"players/app/services"
)

// MessageController handles messaging endpoints
type MessageController struct {
	*contracts.CrudController[models.Message, *requests.MessageCreateRequest, *requests.MessageUpdateRequest]
	messageService *services.MessageService
}

// NewMessageController creates a new message controller
func NewMessageController() *MessageController {
	messageService := services.NewMessageService()

	// Build controller with compile-time enforcement
	crudController := contracts.NewCrudController[models.Message, *requests.MessageCreateRequest, *requests.MessageUpdateRequest](
		"message",
		messageService,
	).
		WithAuthChecker(func(ctx http.Context, action string, resource interface{}) error {
			permHelper := auth.GetPermissionHelper()
			user := permHelper.GetAuthenticatedUser(ctx)
			if user == nil {
				return fmt.Errorf("authentication required")
			}

			// For viewing/updating/deleting, check if user is sender or recipient
			if action == "view" || action == "update" || action == "delete" {
				if msg, ok := resource.(*models.Message); ok {
					if msg.SenderID != user.ID && (msg.RecipientID == nil || *msg.RecipientID != user.ID) {
						return fmt.Errorf("unauthorized to access this message")
					}
				}
			}

			return nil
		}).
		Build()

	controller := &MessageController{
		CrudController: crudController,
		messageService: messageService,
	}

	// Override beforeStore to set sender_id from authenticated user
	controller.SetBeforeStore(func(ctx http.Context, data map[string]interface{}) error {
		permHelper := auth.GetPermissionHelper()
		user := permHelper.GetAuthenticatedUser(ctx)
		if user != nil {
			data["sender_id"] = user.ID
		}
		return nil
	})

	// Prevent direct use of Index - messages should use GetInbox
	controller.SetBeforeIndex(func(ctx http.Context) error {
		return fmt.Errorf("use /messages/conversations or /messages/inbox endpoints instead")
	})

	return controller
}

// Custom endpoints for messaging

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

	return c.ResourceCreatedResponse(ctx, message, "message")
}

// SendBroadcast handles POST /api/messages/broadcast
func (c *MessageController) SendBroadcast(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Check if user has broadcast permission
	if !user.IsSuperAdmin {
		hasPermission := permHelper.CheckPermission(ctx, "messages_broadcast")
		if !hasPermission {
			return c.ForbiddenResponse(ctx, "Broadcast permission required")
		}
	}

	// Parse request
	var request struct {
		RecipientIDs []uint `json:"recipient_ids" validate:"required"`
		Content      string `json:"content" validate:"required,max:5000"`
		Subject      string `json:"subject" validate:"max:255"`
	}

	if err := ctx.Request().Bind(&request); err != nil {
		return c.BadRequestResponse(ctx, "Invalid request data", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Send broadcast
	messages, err := c.messageService.SendBroadcast(user.ID, request.RecipientIDs, request.Content, request.Subject)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to send broadcast: "+err.Error())
	}

	return c.SuccessResponse(ctx, map[string]interface{}{
		"sent_count": len(messages),
		"messages":   messages,
	}, fmt.Sprintf("Broadcast sent to %d recipients", len(messages)))
}

// GetConversations handles GET /api/messages/conversations
func (c *MessageController) GetConversations(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Get pagination parameters
	page := ctx.Request().QueryInt("page", 1)
	pageSize := ctx.Request().QueryInt("pageSize", 20)
	if pageSize > 100 {
		pageSize = 100
	}

	// Get conversations from message service
	conversations, err := c.messageService.GetConversations(user.ID)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve conversations: "+err.Error())
	}

	// Paginate the results
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > len(conversations) {
		end = len(conversations)
	}

	paginatedConversations := conversations
	if start < len(conversations) {
		paginatedConversations = conversations[start:end]
	} else {
		paginatedConversations = []interface{}{}
	}

	// Build paginated response
	paginatedResult := map[string]interface{}{
		"data":         paginatedConversations,
		"total":        len(conversations),
		"current_page": page,
		"per_page":     pageSize,
		"last_page":    (len(conversations) + pageSize - 1) / pageSize,
		"from":         start + 1,
		"to":           end,
	}

	return c.SuccessResponse(ctx, paginatedResult, "Conversations retrieved successfully")
}

// GetInbox handles GET /api/messages/inbox
func (c *MessageController) GetInbox(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Get pagination request
	req, err := c.ValidatePaginationRequest(ctx)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid pagination parameters", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Get inbox messages
	result, err := c.messageService.GetInbox(user.ID, *req)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve inbox: "+err.Error())
	}

	// Get unread count
	unreadCount, _ := c.messageService.GetUnreadCount(user.ID)

	// Build response with unread count
	response := c.BuildPaginatedResponse(result, req)
	// The response is already a map, so we can add unread_count directly
	response["unread_count"] = unreadCount

	return c.SuccessResponse(ctx, response, "Inbox retrieved successfully")
}

// GetSentMessages handles GET /api/messages/sent
func (c *MessageController) GetSentMessages(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Get pagination request
	req, err := c.ValidatePaginationRequest(ctx)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid pagination parameters", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Get sent messages
	result, err := c.messageService.GetSentMessages(user.ID, *req)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve sent messages: "+err.Error())
	}

	response := c.BuildPaginatedResponse(result, req)
	return c.SuccessResponse(ctx, response, "Sent messages retrieved successfully")
}

// GetUnreadMessages handles GET /api/messages/unread
func (c *MessageController) GetUnreadMessages(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Get pagination request
	req, err := c.ValidatePaginationRequest(ctx)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid pagination parameters", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Get unread messages
	result, err := c.messageService.GetUnreadMessages(user.ID, *req)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve unread messages: "+err.Error())
	}

	response := c.BuildPaginatedResponse(result, req)
	return c.SuccessResponse(ctx, response, "Unread messages retrieved successfully")
}

// MarkAsRead handles PUT /api/messages/{id}/read
func (c *MessageController) MarkAsRead(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Get message ID
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid message ID", nil)
	}

	// Mark as read
	if err := c.messageService.MarkAsRead(id, user.ID); err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	return c.SuccessResponse(ctx, nil, "Message marked as read")
}

// MarkAsImportant handles PUT /api/messages/{id}/important
func (c *MessageController) MarkAsImportant(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Get message ID
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid message ID", nil)
	}

	// Get important flag
	var request struct {
		Important bool `json:"important"`
	}
	if err := ctx.Request().Bind(&request); err != nil {
		return c.BadRequestResponse(ctx, "Invalid request data", nil)
	}

	// Mark as important/unimportant
	if err := c.messageService.MarkAsImportant(id, user.ID, request.Important); err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	status := "unimportant"
	if request.Important {
		status = "important"
	}

	return c.SuccessResponse(ctx, nil, fmt.Sprintf("Message marked as %s", status))
}

// GetConversation handles GET /api/messages/conversation/{userId}
func (c *MessageController) GetConversation(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Get other user ID
	otherUserID, err := c.ValidateID(ctx, "userId")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid user ID", nil)
	}

	// Get pagination request
	req, err := c.ValidatePaginationRequest(ctx)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid pagination parameters", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Get conversation
	result, err := c.messageService.GetConversation(user.ID, otherUserID, *req)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve conversation: "+err.Error())
	}

	response := c.BuildPaginatedResponse(result, req)
	return c.SuccessResponse(ctx, response, "Conversation retrieved successfully")
}

// GetMessagableUsers handles GET /api/messages/users
func (c *MessageController) GetMessagableUsers(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Get users based on role-based discovery rules:
	// - Super admins can message anyone
	// - Regular users can discover people in their roles or below
	var users []models.User

	if user.IsSuperAdmin {
		// Super admins can see all active users except themselves
		facades.Log().Info("Super admin getting all users", map[string]interface{}{
			"current_user_id": user.ID,
		})
		if err := facades.Orm().Query().Model(&models.User{}).
			Where("is_active = ?", true).
			Where("id != ?", user.ID).
			Order("name ASC").
			Find(&users); err != nil {
			return c.InternalErrorResponse(ctx, "Failed to retrieve users: "+err.Error())
		}
		facades.Log().Info("Super admin found users", map[string]interface{}{
			"users_count": len(users),
		})
	} else {
		// Regular users can only discover users in their roles or below
		// For now, implement a simplified version - they can see other regular users
		// TODO: Implement proper role hierarchy checking
		if err := facades.Orm().Query().Model(&models.User{}).
			Where("is_active = ?", true).
			Where("id != ?", user.ID).
			Where("is_super_admin = ?", false). // Regular users can see other regular users
			Order("name ASC").
			Find(&users); err != nil {
			return c.InternalErrorResponse(ctx, "Failed to retrieve users: "+err.Error())
		}
	}

	// Build messagable users list
	messagableUsers := []interface{}{}
	for _, u := range users {
		messagableUsers = append(messagableUsers, map[string]interface{}{
			"id":             u.ID,
			"name":           u.Name,
			"email":          u.Email,
			"is_super_admin": u.IsSuperAdmin,
			"is_active":      u.IsActive,
		})
	}

	// Return in the expected format
	// The frontend expects response.data.data to be the paginated structure
	paginatedResult := map[string]interface{}{
		"data":         messagableUsers,
		"total":        len(messagableUsers),
		"current_page": 1,
		"per_page":     100,
		"last_page":    1,
		"from":         1,
		"to":           len(messagableUsers),
	}

	return c.SuccessResponse(ctx, paginatedResult, "Messagable users retrieved successfully")
}

// Override Delete to use custom delete logic
func (c *MessageController) Delete(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Get message ID
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid message ID", nil)
	}

	// Delete message
	if err := c.messageService.DeleteMessage(id, user.ID); err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	return c.ResourceDeletedResponse(ctx, "message", id)
}
