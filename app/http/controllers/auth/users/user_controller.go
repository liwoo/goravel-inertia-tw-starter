package users

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

// AssignToSme assigns a user to an SME by creating a PrimaryBusinessOwner record
func (c *UserController) AssignToSme(ctx http.Context) http.Response {
	// Check authentication and authorization
	permHelper := auth.GetPermissionHelper()
	currentUser, err := permHelper.RequireAuthentication(ctx)
	if err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]interface{}{
			"message": "Unauthorized",
		})
	}

	if !currentUser.IsSuperAdminUser() {
		return ctx.Response().Json(http.StatusForbidden, map[string]interface{}{
			"message": "Super admin access required",
		})
	}

	// Get user ID from URL
	userID := ctx.Request().RouteInt("id")
	if userID == 0 {
		return ctx.Response().Json(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid user ID",
		})
	}

	// Parse request body
	var requestBody struct {
		SmeID int `json:"sme_id"`
	}
	if err := ctx.Request().Bind(&requestBody); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
		})
	}

	if requestBody.SmeID == 0 {
		return ctx.Response().Json(http.StatusBadRequest, map[string]interface{}{
			"message": "MSME ID is required",
		})
	}

	// Fetch the user
	var user models.User
	if err := facades.Orm().Query().Where("id = ?", userID).First(&user); err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]interface{}{
			"message": "User not found",
		})
	}

	// Fetch the SME
	var sme models.Sme
	if err := facades.Orm().Query().Where("id = ?", requestBody.SmeID).First(&sme); err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]interface{}{
			"message": "MSME not found",
		})
	}

	// Check if user is already assigned to this SME (has a PrimaryBusinessOwner record)
	var existingOwner models.PrimaryBusinessOwner
	checkErr := facades.Orm().Query().Where("sme_id = ? AND email = ?", requestBody.SmeID, user.Email).First(&existingOwner)
	if checkErr == nil && existingOwner.ID > 0 {
		return ctx.Response().Json(http.StatusConflict, map[string]interface{}{
			"message": "User is already assigned to this MSME",
		})
	}

	// Create PrimaryBusinessOwner record linking user to SME
	// Split user name into first and last name
	nameParts := splitName(user.Name)
	createdByID := int(currentUser.ID)

	primaryOwner := models.PrimaryBusinessOwner{
		FirstName:   nameParts[0],
		LastName:    nameParts[1],
		Email:       &user.Email,
		PhoneNumber: "", // Required field - empty as we don't have it from user
		SmeID:       requestBody.SmeID,
		CreatedBy:   &createdByID,
	}

	if err := facades.Orm().Query().Create(&primaryOwner); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]interface{}{
			"message": fmt.Sprintf("Failed to assign user to MSME: %v", err),
		})
	}

	facades.Log().Info("User assigned to SME", map[string]interface{}{
		"user_id":          userID,
		"sme_id":           requestBody.SmeID,
		"primary_owner_id": primaryOwner.ID,
		"assigned_by":      currentUser.ID,
	})

	return ctx.Response().Json(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("User %s has been assigned to MSME %s", user.Name, sme.Name),
		"data": map[string]interface{}{
			"user_id":          userID,
			"sme_id":           requestBody.SmeID,
			"primary_owner_id": primaryOwner.ID,
		},
	})
}

// splitName splits a full name into first name and last name
func splitName(fullName string) []string {
	parts := make([]string, 2)
	parts[0] = fullName
	parts[1] = ""

	// Simple split by space
	for i, c := range fullName {
		if c == ' ' {
			parts[0] = fullName[:i]
			parts[1] = fullName[i+1:]
			break
		}
	}

	return parts
}
