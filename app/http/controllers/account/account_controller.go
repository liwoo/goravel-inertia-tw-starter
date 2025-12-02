package account

import (
	"regexp"
	nethttp "net/http"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"

	"smedi-sme-db/app/http/requests"
	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
)

// AccountController handles account-related API endpoints
type AccountController struct {
	activityService  *services.UserActivityService
	dashboardService *services.DashboardService
}

// NewAccountController creates a new account controller
func NewAccountController() *AccountController {
	return &AccountController{
		activityService:  services.NewUserActivityService(),
		dashboardService: services.NewDashboardService(),
	}
}

// UpdateProfile godoc
// @Summary      Update user profile
// @Description  Update the authenticated user's profile (name and email)
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        profile  body  requests.UpdateProfileRequest  true  "Profile data"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /account/profile [put]
func (c *AccountController) UpdateProfile(ctx http.Context) http.Response {
	// Get authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(nethttp.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
		})
	}

	// Bind and validate request
	var request requests.UpdateProfileRequest
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Json(nethttp.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request data",
			"errors":  err.Error(),
		})
	}

	// Validate using Goravel's validator
	validator, err := facades.Validation().Make(map[string]interface{}{
		"name":  request.Name,
		"email": request.Email,
	}, map[string]string{
		"name":  "required|min_len:2|max_len:255",
		"email": "required|email|max_len:255",
	})
	if err != nil {
		return ctx.Response().Json(nethttp.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Validation error",
		})
	}

	if validator.Fails() {
		return ctx.Response().Json(nethttp.StatusUnprocessableEntity, map[string]interface{}{
			"success": false,
			"message": "Validation failed",
			"errors":  validator.Errors().All(),
		})
	}

	// Check if email is being changed and if it's already taken
	if request.Email != user.Email {
		var existingUser models.User
		err := facades.Orm().Query().Where("email = ? AND id != ?", request.Email, user.ID).First(&existingUser)
		if err == nil && existingUser.ID > 0 {
			return ctx.Response().Json(nethttp.StatusUnprocessableEntity, map[string]interface{}{
				"success": false,
				"message": "Validation failed",
				"errors": map[string][]string{
					"email": {"This email is already in use"},
				},
			})
		}
	}

	// Store old values for activity log
	oldName := user.Name
	oldEmail := user.Email

	// Update user
	user.Name = request.Name
	user.Email = request.Email

	if err := facades.Orm().Query().Save(&user); err != nil {
		return ctx.Response().Json(nethttp.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to update profile",
		})
	}

	// Log the activity
	metadata := map[string]interface{}{
		"changes": map[string]interface{}{},
	}
	if oldName != request.Name {
		metadata["changes"].(map[string]interface{})["name"] = map[string]string{
			"from": oldName,
			"to":   request.Name,
		}
	}
	if oldEmail != request.Email {
		metadata["changes"].(map[string]interface{})["email"] = map[string]string{
			"from": oldEmail,
			"to":   request.Email,
		}
	}

	c.activityService.LogActivity(ctx, user.ID, models.ActivityProfileUpdate, "Profile updated", metadata)

	// Reload user with relations
	facades.Orm().Query().With("Roles").Find(&user, user.ID)

	return ctx.Response().Json(nethttp.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Profile updated successfully",
		"data": map[string]interface{}{
			"id":            user.ID,
			"name":          user.Name,
			"email":         user.Email,
			"is_active":     user.IsActive,
			"is_super_admin": user.IsSuperAdmin,
			"roles":         user.Roles,
		},
	})
}

// ChangePassword godoc
// @Summary      Change user password
// @Description  Change the authenticated user's password
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        password  body  requests.ChangePasswordRequest  true  "Password data"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /account/password [put]
func (c *AccountController) ChangePassword(ctx http.Context) http.Response {
	// Get authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(nethttp.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
		})
	}

	// Bind and validate request
	var request requests.ChangePasswordRequest
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Json(nethttp.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request data",
			"errors":  err.Error(),
		})
	}

	// Validate using Goravel's validator
	validator, err := facades.Validation().Make(map[string]interface{}{
		"current_password": request.CurrentPassword,
		"new_password":     request.NewPassword,
		"confirm_password": request.ConfirmPassword,
	}, map[string]string{
		"current_password": "required",
		"new_password":     "required|min_len:8|max_len:128",
		"confirm_password": "required",
	})
	if err != nil {
		return ctx.Response().Json(nethttp.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Validation error",
		})
	}

	if validator.Fails() {
		return ctx.Response().Json(nethttp.StatusUnprocessableEntity, map[string]interface{}{
			"success": false,
			"message": "Validation failed",
			"errors":  validator.Errors().All(),
		})
	}

	// Check if passwords match
	if !request.ValidatePasswordMatch() {
		return ctx.Response().Json(nethttp.StatusUnprocessableEntity, map[string]interface{}{
			"success": false,
			"message": "Validation failed",
			"errors": map[string][]string{
				"confirm_password": {"Passwords do not match"},
			},
		})
	}

	// Verify current password - need to load the password field
	var userWithPassword models.User
	if err := facades.Orm().Query().Select("id", "password").Find(&userWithPassword, user.ID); err != nil {
		return ctx.Response().Json(nethttp.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to verify password",
		})
	}

	if !facades.Hash().Check(request.CurrentPassword, userWithPassword.Password) {
		return ctx.Response().Json(nethttp.StatusUnprocessableEntity, map[string]interface{}{
			"success": false,
			"message": "Validation failed",
			"errors": map[string][]string{
				"current_password": {"Current password is incorrect"},
			},
		})
	}

	// Hash new password
	hashedPassword, err := facades.Hash().Make(request.NewPassword)
	if err != nil {
		return ctx.Response().Json(nethttp.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to process password",
		})
	}

	// Update password
	if _, err := facades.Orm().Query().Model(&models.User{}).Where("id = ?", user.ID).Update("password", hashedPassword); err != nil {
		return ctx.Response().Json(nethttp.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to update password",
		})
	}

	// Log the activity
	c.activityService.LogActivity(ctx, user.ID, models.ActivityPasswordChange, "Password changed", nil)

	return ctx.Response().Json(nethttp.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Password changed successfully",
	})
}

// GetActivities godoc
// @Summary      Get user activities
// @Description  Get paginated activity history for the authenticated user
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        page       query  int     false  "Page number" default(1)
// @Param        pageSize   query  int     false  "Items per page" default(20)
// @Param        type       query  string  false  "Filter by activity type"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /account/activities [get]
func (c *AccountController) GetActivities(ctx http.Context) http.Response {
	// Get authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(nethttp.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
		})
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.Request().Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Request().Query("pageSize", "20"))
	activityType := ctx.Request().Query("type", "")

	// Get activities
	result, err := c.activityService.GetUserActivities(user.ID, page, pageSize, activityType)
	if err != nil {
		return ctx.Response().Json(nethttp.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to retrieve activities",
		})
	}

	return ctx.Response().Json(nethttp.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Activities retrieved successfully",
		"data":    result.Data,
		"meta": map[string]interface{}{
			"total":      result.Total,
			"page":       result.Page,
			"pageSize":   result.PageSize,
			"totalPages": result.TotalPages,
		},
	})
}

// GetActivityTypes godoc
// @Summary      Get activity types
// @Description  Get list of activity types for the authenticated user
// @Tags         account
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /account/activities/types [get]
func (c *AccountController) GetActivityTypes(ctx http.Context) http.Response {
	// Get authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(nethttp.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
		})
	}

	// Get activity types
	types, err := c.activityService.GetActivityTypes(user.ID)
	if err != nil {
		return ctx.Response().Json(nethttp.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to retrieve activity types",
		})
	}

	return ctx.Response().Json(nethttp.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Activity types retrieved successfully",
		"data":    types,
	})
}

// GetActivitySummary godoc
// @Summary      Get activity summary
// @Description  Get summary of activities grouped by type for the authenticated user
// @Tags         account
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /account/activities/summary [get]
func (c *AccountController) GetActivitySummary(ctx http.Context) http.Response {
	// Get authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(nethttp.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
		})
	}

	// Get activity summary
	summary, err := c.activityService.GetActivitySummary(user.ID)
	if err != nil {
		return ctx.Response().Json(nethttp.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to retrieve activity summary",
		})
	}

	return ctx.Response().Json(nethttp.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Activity summary retrieved successfully",
		"data":    summary,
	})
}

// GetProfile godoc
// @Summary      Get user profile
// @Description  Get the authenticated user's profile
// @Tags         account
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /account/profile [get]
func (c *AccountController) GetProfile(ctx http.Context) http.Response {
	// Get authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(nethttp.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
		})
	}

	// Load user with relations
	if err := facades.Orm().Query().With("Roles").Find(&user, user.ID); err != nil {
		return ctx.Response().Json(nethttp.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to load profile",
		})
	}

	return ctx.Response().Json(nethttp.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Profile retrieved successfully",
		"data": map[string]interface{}{
			"id":             user.ID,
			"name":           user.Name,
			"email":          user.Email,
			"is_active":      user.IsActive,
			"is_super_admin": user.IsSuperAdmin,
			"email_verified": user.EmailVerified,
			"last_login_at":  user.LastLoginAt,
			"roles":          user.Roles,
			"created_at":     user.CreatedAt,
			"updated_at":     user.UpdatedAt,
		},
	})
}

// GetRecentActivities godoc
// @Summary      Get user's recent entity activities
// @Description  Get recent activities (SMEs, events, procurements) created or updated by the authenticated user
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        limit  query  int  false  "Number of activities to return" default(20)
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /account/recent-activities [get]
func (c *AccountController) GetRecentActivities(ctx http.Context) http.Response {
	// Get authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(nethttp.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
		})
	}

	// Parse limit parameter
	limit, _ := strconv.Atoi(ctx.Request().Query("limit", "20"))
	if limit < 1 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}

	// Get recent activities for this user
	activities := c.dashboardService.GetUserRecentActivities(user.ID, limit)

	return ctx.Response().Json(nethttp.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Recent activities retrieved successfully",
		"data":    activities,
	})
}

// helper variable to silence "imported and not used" error
var _ validation.Validator
var _ = regexp.Compile
