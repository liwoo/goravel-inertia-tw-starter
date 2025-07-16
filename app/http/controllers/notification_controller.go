package controllers

import (
	"strconv"
	"time"

	"github.com/goravel/framework/contracts/http"
	"players/app/auth"
	"players/app/contracts"
	"players/app/services"
)

type NotificationController struct {
	*contracts.BaseCrudController
	notificationService *services.NotificationService
}

func NewNotificationController() *NotificationController {
	return &NotificationController{
		BaseCrudController:  contracts.NewBaseCrudController("notification"),
		notificationService: services.NewNotificationService(),
	}
}

// GetNotifications handles GET /api/notifications
func (c *NotificationController) GetNotifications(ctx http.Context) http.Response {
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

	// Parse filters
	if ctx.Request().Query("unread_only", "false") == "true" {
		req.Filters["unread_only"] = true
	}
	if notifType := ctx.Request().Query("type", ""); notifType != "" {
		req.Filters["type"] = notifType
	}
	if priority := ctx.Request().Query("priority", ""); priority != "" {
		req.Filters["priority"] = priority
	}

	// Get notifications
	result, err := c.notificationService.GetUserNotifications(user.ID, *req)
	if err != nil {
		return c.InternalErrorResponse(ctx, err.Error())
	}

	return c.SuccessResponse(ctx, c.BuildPaginatedResponse(result, req), "Notifications retrieved successfully")
}

// MarkAsRead handles PUT /api/notifications/{id}/read
func (c *NotificationController) MarkAsRead(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Parse notification ID from route
	notificationID, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	// Mark as read
	if err := c.notificationService.MarkAsRead(notificationID, user.ID); err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	return c.SuccessResponse(ctx, nil, "Notification marked as read")
}

// MarkAllAsRead handles PUT /api/notifications/read-all
func (c *NotificationController) MarkAllAsRead(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Mark all as read
	if err := c.notificationService.MarkAllAsRead(user.ID); err != nil {
		return c.InternalErrorResponse(ctx, err.Error())
	}

	return c.SuccessResponse(ctx, nil, "All notifications marked as read")
}

// DismissNotification handles DELETE /api/notifications/{id}
func (c *NotificationController) DismissNotification(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Parse notification ID from route
	notificationID, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	// Dismiss notification
	if err := c.notificationService.DismissNotification(notificationID, user.ID); err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	return c.SuccessResponse(ctx, nil, "Notification dismissed")
}

// DismissAllNotifications handles DELETE /api/notifications
func (c *NotificationController) DismissAllNotifications(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Dismiss all notifications
	if err := c.notificationService.DismissAllNotifications(user.ID); err != nil {
		return c.InternalErrorResponse(ctx, err.Error())
	}

	return c.SuccessResponse(ctx, nil, "All notifications dismissed")
}

// GetCounts handles GET /api/notifications/counts
func (c *NotificationController) GetCounts(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Get notification counts
	counts, err := c.notificationService.GetNotificationCounts(user.ID)
	if err != nil {
		return c.InternalErrorResponse(ctx, err.Error())
	}

	return c.SuccessResponse(ctx, counts, "Notification counts retrieved successfully")
}

// CreateNotification handles POST /api/notifications (admin only)
func (c *NotificationController) CreateNotification(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Check if user is super admin
	if !user.IsSuperAdminUser() {
		return c.ForbiddenResponse(ctx, "Super admin access required")
	}

	// Parse request
	var request struct {
		UserID        uint    `json:"user_id" validate:"required"`
		Title         string  `json:"title" validate:"required,max:255"`
		Message       string  `json:"message" validate:"max:1000"`
		Type          string  `json:"type" validate:"required,max:50"`
		Priority      string  `json:"priority"`
		ExpiresInDays *int    `json:"expires_in_days"`
		Data          string  `json:"data"`
	}

	if err := ctx.Request().Bind(&request); err != nil {
		return c.BadRequestResponse(ctx, "Invalid request data", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Calculate expiration
	var expiresAt *time.Time
	if request.ExpiresInDays != nil && *request.ExpiresInDays > 0 {
		expiry := time.Now().AddDate(0, 0, *request.ExpiresInDays)
		expiresAt = &expiry
	}

	// Create notification
	notification, err := c.notificationService.CreateNotification(
		request.UserID,
		request.Title,
		request.Message,
		request.Type,
		&user.ID,
		nil,
		nil,
		request.Priority,
		expiresAt,
		request.Data,
	)
	if err != nil {
		return c.BadRequestResponse(ctx, err.Error(), nil)
	}

	return c.CreatedResponse(ctx, notification, "Notification created successfully")
}

// CreateSystemNotification handles POST /api/notifications/system (super admin only)
func (c *NotificationController) CreateSystemNotification(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Check if user is super admin
	if !user.IsSuperAdminUser() {
		return c.ForbiddenResponse(ctx, "Super admin access required")
	}

	// Parse request
	var request struct {
		Title         string   `json:"title" validate:"required,max:255"`
		Message       string   `json:"message" validate:"required,max:1000"`
		Priority      string   `json:"priority"`
		ExpiresInDays *int     `json:"expires_in_days"`
		RoleFilter    []string `json:"role_filter"`
	}

	if err := ctx.Request().Bind(&request); err != nil {
		return c.BadRequestResponse(ctx, "Invalid request data", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Calculate expiration
	var expiresAt *time.Time
	if request.ExpiresInDays != nil && *request.ExpiresInDays > 0 {
		expiry := time.Now().AddDate(0, 0, *request.ExpiresInDays)
		expiresAt = &expiry
	}

	// Create system notification
	if err := c.notificationService.CreateSystemNotification(
		request.Title,
		request.Message,
		request.Priority,
		expiresAt,
		request.RoleFilter,
	); err != nil {
		return c.InternalErrorResponse(ctx, err.Error())
	}

	return c.SuccessResponse(ctx, nil, "System notification created successfully")
}

// BatchMarkAsRead handles PUT /api/notifications/batch/read
func (c *NotificationController) BatchMarkAsRead(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Parse request
	var request struct {
		NotificationIDs []uint `json:"notification_ids" validate:"required"`
	}

	if err := ctx.Request().Bind(&request); err != nil {
		return c.BadRequestResponse(ctx, "Invalid request data", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Batch mark as read
	if err := c.notificationService.BatchMarkAsRead(request.NotificationIDs, user.ID); err != nil {
		return c.InternalErrorResponse(ctx, err.Error())
	}

	return c.SuccessResponse(ctx, nil, "Notifications marked as read")
}

// BatchDismiss handles DELETE /api/notifications/batch
func (c *NotificationController) BatchDismiss(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Parse request
	var request struct {
		NotificationIDs []uint `json:"notification_ids" validate:"required"`
	}

	if err := ctx.Request().Bind(&request); err != nil {
		return c.BadRequestResponse(ctx, "Invalid request data", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Batch dismiss
	if err := c.notificationService.BatchDismiss(request.NotificationIDs, user.ID); err != nil {
		return c.InternalErrorResponse(ctx, err.Error())
	}

	return c.SuccessResponse(ctx, nil, "Notifications dismissed")
}

// GetByType handles GET /api/notifications/type/{type}
func (c *NotificationController) GetByType(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Parse type from route
	notificationType := ctx.Request().Route("type")
	if notificationType == "" {
		return c.BadRequestResponse(ctx, "Notification type is required", nil)
	}

	// Parse limit
	limitStr := ctx.Request().Query("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	// Get notifications by type
	notifications, err := c.notificationService.GetNotificationsByType(user.ID, notificationType, limit)
	if err != nil {
		return c.InternalErrorResponse(ctx, err.Error())
	}

	return c.SuccessResponse(ctx, notifications, "Notifications retrieved successfully")
}

// CleanupExpired handles POST /api/notifications/cleanup (admin only)
func (c *NotificationController) CleanupExpired(ctx http.Context) http.Response {
	// Get authenticated user
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return c.ForbiddenResponse(ctx, "Authentication required")
	}

	// Check if user is admin
	if !user.IsAdmin() {
		return c.ForbiddenResponse(ctx, "Admin access required")
	}

	// Cleanup expired notifications
	if err := c.notificationService.CleanupExpiredNotifications(); err != nil {
		return c.InternalErrorResponse(ctx, err.Error())
	}

	return c.SuccessResponse(ctx, nil, "Expired notifications cleaned up successfully")
}