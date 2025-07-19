package messages

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"players/app/auth"
	"players/app/contracts"
	"players/app/models"
	"players/app/services"
)

// MessageController - Simplified version using generic CRUD controller
type MessageController struct {
	*contracts.GenericCrudController[models.Message, interface{}, interface{}]
	messageService *services.MessageService
}

// NewMessageController creates a new simplified message controller
func NewMessageController() *MessageController {
	messageService := services.NewMessageService()

	// Create the generic controller - using interface{} for request types since messages have custom handling
	genericController := contracts.NewGenericCrudController[models.Message, interface{}, interface{}](
		"message",
		messageService,
	)

	controller := &MessageController{
		GenericCrudController: genericController,
		messageService:        messageService,
	}

	// Configure authorization - all message operations require authentication
	genericController.SetAuthCheck(func(ctx http.Context, action string, resource interface{}) error {
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
	})

	// Override Index to show user's inbox by default
	genericController.SetBeforeIndex(func(ctx http.Context) error {
		// This will be handled in custom GetInbox method
		return fmt.Errorf("use /messages/inbox endpoint instead")
	})

	// Register controller
	contracts.MustRegisterCrudController("messages", controller)

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

// Contract method implementations
func (c *MessageController) GetSearchableFields() []string {
	return c.messageService.GetSearchableFields()
}

func (c *MessageController) GetValidationRules() map[string]interface{} {
	return c.messageService.GetValidationRules()
}

func (c *MessageController) CheckPermission(ctx http.Context, permission string, resource interface{}) error {
	if c.GenericCrudController.CheckAuth != nil {
		return c.GenericCrudController.CheckAuth(ctx, permission, resource)
	}
	return nil
}

func (c *MessageController) GetCurrentUser(ctx http.Context) interface{} {
	permHelper := auth.GetPermissionHelper()
	return permHelper.GetAuthenticatedUser(ctx)
}

func (c *MessageController) RequireAuthentication(ctx http.Context) error {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		return fmt.Errorf("authentication required")
	}
	return nil
}

func (c *MessageController) BuildPermissionsMap(ctx http.Context, resourceType string) map[string]bool {
	permHelper := auth.GetPermissionHelper()
	return permHelper.BuildPermissionsMap(ctx, resourceType)
}
