package auth

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"players/app/auth"
	"players/app/contracts"
	"players/app/helpers"
	"players/app/http/requests"
	"players/app/services"
)

// UserController - API controller with contract enforcement for user management
// Only accessible by super admins
type UserController struct {
	*contracts.BaseCrudController
	userService *services.UserService
	authHelper  contracts.AuthHelper
}

// NewUserController creates a new user controller that implements all contracts
func NewUserController() *UserController {
	controller := &UserController{
		BaseCrudController: contracts.NewBaseCrudController("user"),
		userService:        services.NewUserService(),
		authHelper:         helpers.NewAuthHelper(),
	}

	// Register controller with validation
	contracts.MustRegisterCrudController("users", controller)

	return controller
}

// checkSuperAdmin verifies if the current user is a super admin
func (c *UserController) checkSuperAdmin(ctx http.Context) error {
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	if user == nil || !user.IsSuperAdmin {
		return fmt.Errorf("super admin access required")
	}
	return nil
}

// Index GET /users - Implements CrudControllerContract
func (c *UserController) Index(ctx http.Context) http.Response {
	// Check super admin access
	if err := c.checkSuperAdmin(ctx); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: Super admin privileges required")
	}

	// Validate pagination request using contract
	req, err := c.ValidatePaginationRequest(ctx)
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid pagination parameters", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Get users using service
	result, err := c.userService.GetList(*req)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve users: "+err.Error())
	}

	// Build standardized paginated response
	response := c.BuildPaginatedResponse(result, req)
	return c.SuccessResponse(ctx, response, "Users retrieved successfully")
}

// Show GET /users/{id} - Implements CrudControllerContract
func (c *UserController) Show(ctx http.Context) http.Response {
	// Check super admin access
	if err := c.checkSuperAdmin(ctx); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: Super admin privileges required")
	}

	// Validate ID parameter using contract
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid user ID", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Get the user
	user, err := c.userService.GetByID(id)
	if err != nil {
		return c.ResourceNotFoundResponse(ctx, "user", id)
	}

	return c.SuccessResponse(ctx, user, "User details retrieved successfully")
}

// Store POST /users - Implements CrudControllerContract
func (c *UserController) Store(ctx http.Context) http.Response {
	// Check super admin access
	if err := c.checkSuperAdmin(ctx); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: Super admin privileges required")
	}

	// Validate create request using contract
	data, err := c.ValidateCreateRequest(ctx)
	if err != nil {
		return c.ValidationErrorResponse(ctx, map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Create the user using validated data
	user, err := c.userService.Create(data)
	if err != nil {
		// Check for specific validation errors
		if err.Error() == "email already exists" {
			return c.ValidationErrorResponse(ctx, map[string]interface{}{
				"validation_error": "The email address is already in use",
			})
		}
		return c.InternalErrorResponse(ctx, "Failed to create user: "+err.Error())
	}

	return c.ResourceCreatedResponse(ctx, user, "user")
}

// Update PUT /users/{id} - Implements CrudControllerContract
func (c *UserController) Update(ctx http.Context) http.Response {
	// Check super admin access
	if err := c.checkSuperAdmin(ctx); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: Super admin privileges required")
	}

	// Validate ID parameter using contract
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid user ID", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Check if user exists
	_, err = c.userService.GetByID(id)
	if err != nil {
		return c.ResourceNotFoundResponse(ctx, "user", id)
	}

	// Validate update request using contract
	data, err := c.ValidateUpdateRequest(ctx, id)
	if err != nil {
		return c.ValidationErrorResponse(ctx, map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Update the user using validated data
	updatedUser, err := c.userService.Update(id, data)
	if err != nil {
		// Check for specific validation errors
		if err.Error() == "email already exists" {
			return c.ValidationErrorResponse(ctx, map[string]interface{}{
				"validation_error": "The email address is already in use",
			})
		}
		return c.InternalErrorResponse(ctx, "Failed to update user: "+err.Error())
	}

	return c.ResourceUpdatedResponse(ctx, updatedUser, "user")
}

// Delete DELETE /users/{id} - Implements CrudControllerContract
func (c *UserController) Delete(ctx http.Context) http.Response {
	// Check super admin access
	if err := c.checkSuperAdmin(ctx); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: Super admin privileges required")
	}

	// Validate ID parameter using contract
	id, err := c.ValidateID(ctx, "id")
	if err != nil {
		return c.BadRequestResponse(ctx, "Invalid user ID", map[string]interface{}{
			"validation_error": err.Error(),
		})
	}

	// Check if user exists
	_, err = c.userService.GetByID(id)
	if err != nil {
		return c.ResourceNotFoundResponse(ctx, "user", id)
	}

	// Delete the user
	err = c.userService.Delete(id)
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to delete user: "+err.Error())
	}

	return c.ResourceDeletedResponse(ctx, "user", id)
}

// GetRoles GET /users/roles - Get all available roles for assignment
func (c *UserController) GetRoles(ctx http.Context) http.Response {
	// Check super admin access
	if err := c.checkSuperAdmin(ctx); err != nil {
		return c.ForbiddenResponse(ctx, "Access denied: Super admin privileges required")
	}

	roles, err := c.userService.GetAllRoles()
	if err != nil {
		return c.InternalErrorResponse(ctx, "Failed to retrieve roles: "+err.Error())
	}

	return c.SuccessResponse(ctx, roles, "Roles retrieved successfully")
}

// CONTRACT IMPLEMENTATIONS - Required by ResourceControllerContract interface

// ValidationControllerContract implementation
func (c *UserController) ValidateCreateRequest(ctx http.Context) (map[string]interface{}, error) {
	var createRequest requests.UserCreateRequest
	
	// Bind the data to the struct
	if err := ctx.Request().Bind(&createRequest); err != nil {
		return nil, fmt.Errorf("data binding failed: %w", err)
	}
	
	// Manual validation
	if len(createRequest.Name) < 2 || len(createRequest.Name) > 255 {
		return nil, fmt.Errorf("validation errors: name must be between 2 and 255 characters")
	}
	if createRequest.Email == "" {
		return nil, fmt.Errorf("validation errors: email is required")
	}
	if createRequest.Password == "" || len(createRequest.Password) < 8 {
		return nil, fmt.Errorf("validation errors: password must be at least 8 characters")
	}

	return createRequest.ToCreateData(), nil
}

func (c *UserController) ValidateUpdateRequest(ctx http.Context, id uint) (map[string]interface{}, error) {
	var updateRequest requests.UserUpdateRequest
	updateRequest.ID = id // Set the ID for validation context

	// Bind the data to the struct
	if err := ctx.Request().Bind(&updateRequest); err != nil {
		return nil, fmt.Errorf("data binding failed: %w", err)
	}
	
	// Manual validation for update (all fields optional)
	if updateRequest.Name != "" && (len(updateRequest.Name) < 2 || len(updateRequest.Name) > 255) {
		return nil, fmt.Errorf("validation errors: name must be between 2 and 255 characters")
	}
	if updateRequest.Password != "" && len(updateRequest.Password) < 8 {
		return nil, fmt.Errorf("validation errors: password must be at least 8 characters")
	}

	return updateRequest.ToUpdateData(), nil
}

func (c *UserController) GetValidationRules() map[string]interface{} {
	return c.userService.GetValidationRules()
}

// AuthorizationControllerContract implementation
func (c *UserController) CheckPermission(ctx http.Context, permission string, resource interface{}) error {
	// For user management, we only check super admin status
	return c.checkSuperAdmin(ctx)
}

func (c *UserController) GetCurrentUser(ctx http.Context) interface{} {
	permHelper := auth.GetPermissionHelper()
	return permHelper.GetAuthenticatedUser(ctx)
}

func (c *UserController) RequireAuthentication(ctx http.Context) error {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		return fmt.Errorf("authentication required")
	}
	return nil
}

func (c *UserController) BuildPermissionsMap(ctx http.Context, resourceType string) map[string]bool {
	// For super admin only access, all permissions are based on super admin status
	permHelper := auth.GetPermissionHelper()
	user := permHelper.GetAuthenticatedUser(ctx)
	isSuperAdmin := user != nil && user.IsSuperAdmin
	
	return map[string]bool{
		"canCreate": isSuperAdmin,
		"canEdit":   isSuperAdmin,
		"canDelete": isSuperAdmin,
		"canManage": isSuperAdmin,
		"canExport": isSuperAdmin,
		"canViewReports": isSuperAdmin,
	}
}