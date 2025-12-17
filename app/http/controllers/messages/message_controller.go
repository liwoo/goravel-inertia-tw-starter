package messages

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/http/requests"
	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
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

// BroadcastToRole handles POST /api/messages/broadcast-to-role
// Super admin only - sends a message to all users with a specific role
func (c *MessageController) BroadcastToRole(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Super admin only
	if !user.IsSuperAdmin {
		return c.ForbiddenResponse(ctx, "Super admin access required")
	}

	// Parse request
	var request struct {
		RoleID  uint   `json:"role_id" validate:"required"`
		Content string `json:"content" validate:"required,max:5000"`
		Subject string `json:"subject" validate:"max:255"`
	}

	if err := ctx.Request().Bind(&request); err != nil {
		return c.BadRequestResponse(ctx, "Invalid request data", map[string]interface{}{
			"error": err.Error(),
		})
	}

	if request.RoleID == 0 {
		return c.BadRequestResponse(ctx, "Role ID is required", nil)
	}

	if request.Content == "" {
		return c.BadRequestResponse(ctx, "Message content is required", nil)
	}

	// Send broadcast to role
	result, err := c.messageService.BroadcastToRole(user.ID, request.RoleID, request.Content, request.Subject)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to send broadcast: "+err.Error())
	}

	return c.SuccessResponse(ctx, result, fmt.Sprintf("Broadcast sent to %d users in role", result["sent_count"]))
}

// GetBroadcastHistory handles GET /api/messages/broadcast-history
// Super admin only - retrieves history of broadcasts sent by the user
func (c *MessageController) GetBroadcastHistory(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Super admin only
	if !user.IsSuperAdmin {
		return c.ForbiddenResponse(ctx, "Super admin access required")
	}

	// Get pagination request
	req, err := c.ValidatePaginationRequest(ctx)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid pagination parameters", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Get broadcast history
	result, err := c.messageService.GetBroadcastHistory(user.ID, *req)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve broadcast history: "+err.Error())
	}

	response := c.BuildPaginatedResponse(result, req)
	return c.SuccessResponse(ctx, response, "Broadcast history retrieved successfully")
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

// MarkConversationAsRead handles PUT /api/messages/conversation/{userId}/read
func (c *MessageController) MarkConversationAsRead(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Get other user ID from the conversation
	otherUserID, err := c.ValidateID(ctx, "userId")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid user ID", nil)
	}

	// Mark all messages from this user as read
	if err := c.messageService.MarkConversationAsRead(user.ID, otherUserID); err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	return c.SuccessResponse(ctx, nil, "Conversation marked as read")
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
// Returns users that the current user can message (same role level or lower)
func (c *MessageController) GetMessagableUsers(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Ensure user has roles loaded
	if err := facades.Orm().Query().With("Roles").Find(user, user.ID); err != nil {
		return c.InternalErrorResponse(ctx, "Failed to load user roles: "+err.Error())
	}

	// Get user's role level
	userLevel := user.GetRoleLevel()

	var users []models.User

	if user.IsSuperAdminUser() {
		// Super admins can see all active users except themselves
		facades.Log().Info("Super admin getting all users", map[string]interface{}{
			"current_user_id": user.ID,
		})
		if err := facades.Orm().Query().Model(&models.User{}).
			Where("is_active = ?", true).
			Where("id != ?", user.ID).
			With("Roles").
			Order("name ASC").
			Find(&users); err != nil {
			return c.InternalErrorResponse(ctx, "Failed to retrieve users: "+err.Error())
		}
	} else if userLevel == 0 {
		// Users with no role can only message other users with no role
		facades.Log().Info("User with no role getting users with no role", map[string]interface{}{
			"current_user_id": user.ID,
		})
		// Find users who have no active roles
		if err := facades.Orm().Query().
			Model(&models.User{}).
			Where("is_active = ?", true).
			Where("id != ?", user.ID).
			Where("is_super_admin = ?", false).
			With("Roles").
			Order("name ASC").
			Find(&users); err != nil {
			return c.InternalErrorResponse(ctx, "Failed to retrieve users: "+err.Error())
		}
		// Filter to only users with no roles
		filteredUsers := []models.User{}
		for _, u := range users {
			if u.GetRoleLevel() == 0 {
				filteredUsers = append(filteredUsers, u)
			}
		}
		users = filteredUsers
	} else {
		// Regular users can message users at their role level or lower
		facades.Log().Info("User getting messagable users by role level", map[string]interface{}{
			"current_user_id": user.ID,
			"user_level":      userLevel,
		})

		// Get all active users with their roles
		if err := facades.Orm().Query().
			Model(&models.User{}).
			Where("is_active = ?", true).
			Where("id != ?", user.ID).
			Where("is_super_admin = ?", false). // Can't message super admins
			With("Roles").
			Order("name ASC").
			Find(&users); err != nil {
			return c.InternalErrorResponse(ctx, "Failed to retrieve users: "+err.Error())
		}

		// Filter to users at same level or lower
		filteredUsers := []models.User{}
		for _, u := range users {
			recipientLevel := u.GetRoleLevel()
			// Can message if recipient level <= sender level
			// Users with no role (level 0) can also be messaged
			if recipientLevel <= userLevel {
				filteredUsers = append(filteredUsers, u)
			}
		}
		users = filteredUsers
	}

	facades.Log().Info("Found messagable users", map[string]interface{}{
		"current_user_id": user.ID,
		"users_count":     len(users),
	})

	// Build messagable users list with role info
	messagableUsers := []interface{}{}
	for _, u := range users {
		userRoles := []map[string]interface{}{}
		for _, role := range u.Roles {
			if role.IsActive {
				userRoles = append(userRoles, map[string]interface{}{
					"id":   role.ID,
					"name": role.Name,
					"slug": role.Slug,
				})
			}
		}

		messagableUsers = append(messagableUsers, map[string]interface{}{
			"id":             u.ID,
			"name":           u.Name,
			"email":          u.Email,
			"is_super_admin": u.IsSuperAdmin,
			"is_active":      u.IsActive,
			"roles":          userRoles,
		})
	}

	// Return in the expected format
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

// GetUnreadCount handles GET /api/messages/unread-count
func (c *MessageController) GetUnreadCount(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Get unread count
	unreadCount, err := c.messageService.GetUnreadCount(user.ID)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve unread count: "+err.Error())
	}

	return c.SuccessResponse(ctx, map[string]interface{}{
		"unread_count": unreadCount,
	}, "Unread count retrieved successfully")
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
