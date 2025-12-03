package controllers

import (
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/services"
)

// PresenceController handles user presence/online status endpoints
type PresenceController struct {
	sseService *services.SSEService
}

// NewPresenceController creates a new presence controller
func NewPresenceController() *PresenceController {
	return &PresenceController{
		sseService: services.NewSSEService(),
	}
}

// GetOnlineUsers returns a list of all currently online user IDs
// GET /api/presence/online
func (c *PresenceController) GetOnlineUsers(ctx http.Context) http.Response {
	// Verify user is authenticated
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return ctx.Response().Status(401).Json(http.Json{
			"success": false,
			"message": "Authentication required",
		})
	}

	onlineUserIDs := c.sseService.GetOnlineUserIDs()

	return ctx.Response().Json(http.StatusOK, http.Json{
		"success": true,
		"data": http.Json{
			"online_users":  onlineUserIDs,
			"online_count":  len(onlineUserIDs),
		},
	})
}

// GetUserStatus returns the presence status of a specific user
// GET /api/presence/user/:id
func (c *PresenceController) GetUserStatus(ctx http.Context) http.Response {
	// Verify user is authenticated
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return ctx.Response().Status(401).Json(http.Json{
			"success": false,
			"message": "Authentication required",
		})
	}

	userIDStr := ctx.Request().Route("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"success": false,
			"message": "Invalid user ID",
		})
	}

	status := c.sseService.GetPresenceStatus(uint(userID))

	return ctx.Response().Json(http.StatusOK, http.Json{
		"success": true,
		"data":    status,
	})
}

// GetBulkStatus returns the presence status of multiple users
// POST /api/presence/bulk
// Body: { "user_ids": [1, 2, 3] }
func (c *PresenceController) GetBulkStatus(ctx http.Context) http.Response {
	// Verify user is authenticated
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return ctx.Response().Status(401).Json(http.Json{
			"success": false,
			"message": "Authentication required",
		})
	}

	// Parse request body
	var request struct {
		UserIDs []uint `json:"user_ids"`
	}

	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"success": false,
			"message": "Invalid request body",
		})
	}

	if len(request.UserIDs) == 0 {
		return ctx.Response().Status(400).Json(http.Json{
			"success": false,
			"message": "user_ids is required",
		})
	}

	// Limit bulk requests to prevent abuse
	if len(request.UserIDs) > 100 {
		return ctx.Response().Status(400).Json(http.Json{
			"success": false,
			"message": "Maximum 100 user IDs allowed per request",
		})
	}

	statuses := c.sseService.GetBulkPresenceStatus(request.UserIDs)

	return ctx.Response().Json(http.StatusOK, http.Json{
		"success": true,
		"data":    statuses,
	})
}

// GetStats returns presence statistics
// GET /api/presence/stats
func (c *PresenceController) GetStats(ctx http.Context) http.Response {
	// Verify user is authenticated
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return ctx.Response().Status(401).Json(http.Json{
			"success": false,
			"message": "Authentication required",
		})
	}

	onlineCount := c.sseService.GetOnlineCount()

	return ctx.Response().Json(http.StatusOK, http.Json{
		"success": true,
		"data": http.Json{
			"online_count": onlineCount,
		},
	})
}

// CheckOnline checks if specific users are online (via query param)
// GET /api/presence/check?ids=1,2,3
func (c *PresenceController) CheckOnline(ctx http.Context) http.Response {
	// Verify user is authenticated
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil {
		return ctx.Response().Status(401).Json(http.Json{
			"success": false,
			"message": "Authentication required",
		})
	}

	idsParam := ctx.Request().Query("ids", "")
	if idsParam == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"success": false,
			"message": "ids query parameter is required",
		})
	}

	idStrings := strings.Split(idsParam, ",")
	onlineStatus := make(map[uint]bool)

	for _, idStr := range idStrings {
		id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32)
		if err != nil {
			continue
		}
		onlineStatus[uint(id)] = c.sseService.IsUserOnline(uint(id))
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"success": true,
		"data":    onlineStatus,
	})
}
