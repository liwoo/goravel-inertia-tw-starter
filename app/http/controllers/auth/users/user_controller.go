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
	*contracts.CrudController[models.User, *requests.UserCreateRequest, *requests.UserUpdateRequest]
	userService *services.UserService
}

// NewUserController creates a new user controller with compile-time enforcement
// This ensures all required methods are implemented and configured
func NewUserController() *UserController {
	userService := services.NewUserService()

	// Build controller with compile-time enforcement
	// The builder pattern ensures all required steps are completed
	crudController := contracts.NewCrudController[models.User, *requests.UserCreateRequest, *requests.UserUpdateRequest](
		"user",
		userService,
	).
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
		CrudController: crudController,
		userService:    userService,
	}

	// No custom hooks needed - the UserService handles role assignment in its Create/Update methods

	return controller
}
