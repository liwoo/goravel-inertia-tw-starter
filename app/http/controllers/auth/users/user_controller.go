package users

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"players/app/auth"
	"players/app/contracts"
	"players/app/http/requests"
	"players/app/models"
	"players/app/services"
)

// UserController handles API endpoints for user management
type UserController struct {
	*contracts.StaticEnforcedController[models.User, *requests.UserCreateRequest, *requests.UserUpdateRequest]
	userService *services.UserService
}

// NewUserController creates a new user controller with compile-time enforcement
// This ensures all required methods are implemented and configured
func NewUserController() *UserController {
	userService := services.NewUserService()

	// Build controller with compile-time enforcement
	// The builder pattern ensures all required steps are completed
	staticController := contracts.NewStaticControllerBuilder[models.User, *requests.UserCreateRequest, *requests.UserUpdateRequest](
		"user",
		userService,
	).
		ValidateCreateRequest().
		ValidateUpdateRequest().
		WithAuthChecker(func(ctx http.Context, action string, resource interface{}) error {
			// Users controller is super admin only
			permHelper := auth.GetPermissionHelper()
			user, err := permHelper.RequireAuthentication(ctx)
			if err != nil {
				return err
			}

			if !user.IsSuperAdminUser() {
				return fmt.Errorf("super admin access required")
			}

			return nil
		}).
		Build()

	controller := &UserController{
		StaticEnforcedController: staticController,
		userService:              userService,
	}

	// No custom hooks needed - the UserService handles role assignment in its Create/Update methods

	return controller
}

// GetFilters returns filter metadata for users
func (c *UserController) GetFilters(ctx http.Context) http.Response {
	metadata := map[string]interface{}{
		"filters": []map[string]interface{}{
			{
				"field":     "name",
				"label":     "Name",
				"type":      "string",
				"operators": []string{"contains", "not_contains", "starts_with", "ends_with", "equals", "not_equals"},
			},
			{
				"field":     "email",
				"label":     "Email",
				"type":      "string",
				"operators": []string{"contains", "not_contains", "equals", "not_equals"},
			},
			{
				"field":       "status",
				"label":       "Status",
				"type":        "enum",
				"operators":   []string{"equals", "not_equals", "in", "not_in"},
				"enum_values": []string{"ACTIVE", "INACTIVE", "PENDING"},
			},
			{
				"field":     "created_at",
				"label":     "Created Date",
				"type":      "date",
				"operators": []string{"before", "after", "between", "not_between", "is_today", "is_yesterday", "is_this_week", "is_this_month", "is_this_year"},
			},
		},
	}

	return ctx.Response().Json(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    metadata,
	})
}
