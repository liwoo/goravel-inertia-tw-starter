package services

import (
	"fmt"
	"regexp"

	"github.com/goravel/framework/facades"
	"players/app/auth"
	"players/app/contracts"
	"players/app/models"
)

// UserService implements user-specific business logic using the builder pattern
type UserService struct {
	contracts.CrudServiceContract // Embedded - automatically exposes all methods!
}

// NewUserService creates a new user service using the builder pattern
func NewUserService() *UserService {
	// Create service instance first
	userService := &UserService{}

	// Temporary storage for role assignments between hooks
	var pendingRoleID *uint

	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.User]("user", "id").
		WithSearchFields("name", "email").                                                  // REQUIRED
		WithSortFields("id", "name", "email", "created_at", "updated_at").                  // REQUIRED
		WithFilterFields("email", "is_active", "role", "is_super_admin", "email_verified"). // REQUIRED
		WithValidationRules(map[string]interface{}{                                         // REQUIRED
			"name":     "required|string|max:255",
			"email":    "required|email|unique:users,email",
			"password": "required|string|min:8",
			"role_id":  "exists:roles,id",
		}).
		WithRelations("Roles", "Creator", "Updater"). // Optional
		WithDefaultSort("created_at", "DESC").        // Optional
		WithSoftDeletes().                            // Optional
		WithScopeFiltering("users", "id").            // Optional
		WithBeforeCreate(func(data map[string]interface{}) error {
			// Hash password if provided
			if password, exists := data["password"]; exists && password != "" {
				if passwordStr, ok := password.(string); ok {
					hashedPassword, err := facades.Hash().Make(passwordStr)
					if err != nil {
						return fmt.Errorf("failed to hash password: %v", err)
					}
					data["password"] = hashedPassword
				}
			}

			// Set default is_active status
			if _, exists := data["is_active"]; !exists {
				data["is_active"] = true
			}

			// Extract and store role_id
			if rid, exists := data["role_id"]; exists && rid != nil {
				switch v := rid.(type) {
				case float64:
					roleID := uint(v)
					pendingRoleID = &roleID
				case int:
					roleID := uint(v)
					pendingRoleID = &roleID
				case uint:
					pendingRoleID = &v
				}
				delete(data, "role_id") // Remove as it's not a user field
			}

			return nil
		}).
		WithAfterCreate(func(model *models.User) error {
			// Assign role if one was provided
			if pendingRoleID != nil {
				facades.Log().Info("Assigning role to new user", map[string]interface{}{
					"user_id": model.ID,
					"role_id": *pendingRoleID,
				})

				if err := userService.AssignRole(model.ID, *pendingRoleID); err != nil {
					facades.Log().Error("Failed to assign role to new user", map[string]interface{}{
						"user_id": model.ID,
						"role_id": *pendingRoleID,
						"error":   err.Error(),
					})
				}
				pendingRoleID = nil // Clear after use
			}

			return nil
		}).
		WithBeforeUpdate(func(id uint, data map[string]interface{}) error {
			// Hash password if being updated
			if password, exists := data["password"]; exists && password != "" {
				if passwordStr, ok := password.(string); ok {
					hashedPassword, err := facades.Hash().Make(passwordStr)
					if err != nil {
						return fmt.Errorf("failed to hash password: %v", err)
					}
					data["password"] = hashedPassword
				}
			} else {
				// Remove password field if empty to avoid updating it
				delete(data, "password")
			}

			// Extract and store role_id
			if rid, exists := data["role_id"]; exists && rid != nil {
				switch v := rid.(type) {
				case float64:
					roleID := uint(v)
					pendingRoleID = &roleID
				case int:
					roleID := uint(v)
					pendingRoleID = &roleID
				case uint:
					pendingRoleID = &v
				}
				delete(data, "role_id") // Remove as it's not a user field
			}

			return nil
		}).
		WithAfterUpdate(func(model *models.User) error {
			// Update role if one was provided
			if pendingRoleID != nil {
				facades.Log().Info("Updating user role", map[string]interface{}{
					"user_id": model.ID,
					"role_id": *pendingRoleID,
				})

				// Clear existing roles
				if err := facades.Orm().Query().Model(model).Association("Roles").Clear(); err != nil {
					facades.Log().Error("Failed to clear existing roles", map[string]interface{}{
						"user_id": model.ID,
						"error":   err.Error(),
					})
				}

				// Assign new role
				if err := userService.AssignRole(model.ID, *pendingRoleID); err != nil {
					facades.Log().Error("Failed to update user role", map[string]interface{}{
						"user_id": model.ID,
						"role_id": *pendingRoleID,
						"error":   err.Error(),
					})
				}
				pendingRoleID = nil // Clear after use
			}

			return nil
		}).
		Build() // Returns a fully configured CrudServiceContract

	userService.CrudServiceContract = service

	// Set the actual service reference for proper method resolution
	contracts.SetActualServiceHelper(service, userService, "UserService")

	return userService
}

// User-specific methods beyond basic CRUD

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := facades.Orm().Query().
		Model(&models.User{}).
		Where("email = ?", email).
		With("Role").
		With("CreatedBy").
		With("UpdatedBy").
		First(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetActiveUsers retrieves all active users
func (s *UserService) GetActiveUsers(req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	// Add is_active filter to the request
	if req.Filters == nil {
		req.Filters = make(map[string]interface{})
	}
	req.Filters["is_active"] = true

	return s.GetList(req)
}

// DeactivateUser deactivates a user account
func (s *UserService) DeactivateUser(id uint) error {
	updateData := map[string]interface{}{
		"is_active": false,
	}

	_, err := s.Update(id, updateData)
	return err
}

// ActivateUser activates a user account
func (s *UserService) ActivateUser(id uint) error {
	updateData := map[string]interface{}{
		"is_active": true,
	}

	_, err := s.Update(id, updateData)
	return err
}

// AssignRole assigns a role to a user
func (s *UserService) AssignRole(userID uint, roleID uint) error {
	// Get the permission service to handle role assignment
	permService := auth.GetPermissionService()

	// Get user
	userInterface, err := s.CrudServiceContract.GetByID(userID)
	if err != nil {
		facades.Log().Error("AssignRole: user not found", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return fmt.Errorf("user not found: %v", err)
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		facades.Log().Error("AssignRole: invalid user type", map[string]interface{}{
			"user_id": userID,
			"type":    fmt.Sprintf("%T", userInterface),
		})
		return fmt.Errorf("invalid user type")
	}

	// Get the role by ID to get its slug
	var role models.Role
	err = facades.Orm().Query().Where("id = ?", roleID).First(&role)
	if err != nil {
		facades.Log().Error("AssignRole: failed to find role", map[string]interface{}{
			"role_id": roleID,
			"error":   err.Error(),
		})
		return fmt.Errorf("failed to find role: %v", err)
	}

	err = permService.AssignRole(user, role.Slug, nil)
	if err != nil {
		facades.Log().Error("AssignRole: permService.AssignRole failed", map[string]interface{}{
			"user_id":   userID,
			"role_id":   roleID,
			"role_slug": role.Slug,
			"error":     err.Error(),
		})
		return fmt.Errorf("failed to assign role: %v", err)
	}

	return nil
}

// GetAllRoles retrieves all available roles
func (s *UserService) GetAllRoles() ([]models.Role, error) {
	facades.Log().Info("GetAllRoles called")

	var roles []models.Role
	err := facades.Orm().Query().
		Model(&models.Role{}).
		With("Permissions").
		Find(&roles)

	if err != nil {
		facades.Log().Error("GetAllRoles failed", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	facades.Log().Info("GetAllRoles success", map[string]interface{}{
		"count": len(roles),
	})

	return roles, nil
}

// GetUserWithPermissions retrieves a user with their permissions loaded
func (s *UserService) GetUserWithPermissions(id uint) (*models.User, error) {
	var user models.User
	err := facades.Orm().Query().
		Model(&models.User{}).
		Where("id = ?", id).
		With("Role.Permissions").
		First(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdatePassword updates a user's password
func (s *UserService) UpdatePassword(userID uint, newPassword string) error {
	hashedPassword, err := facades.Hash().Make(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	updateData := map[string]interface{}{
		"password": hashedPassword,
	}

	_, err = s.Update(userID, updateData)
	return err
}

// GetUserStatistics returns statistics about users
func (s *UserService) GetUserStatistics() (map[string]interface{}, error) {
	var stats struct {
		TotalUsers    int64
		ActiveUsers   int64
		InactiveUsers int64
		AdminUsers    int64
	}

	// Get total users
	facades.Orm().Query().Model(&models.User{}).Count(&stats.TotalUsers)

	// Get active users
	facades.Orm().Query().Model(&models.User{}).Where("is_active = ?", true).Count(&stats.ActiveUsers)

	// Get inactive users
	stats.InactiveUsers = stats.TotalUsers - stats.ActiveUsers

	// Get admin users (assuming there's an admin role)
	// Note: Count the super admins instead since user_roles doesn't have is_active field
	facades.Orm().Query().
		Model(&models.User{}).
		Where("is_super_admin = ?", true).
		Count(&stats.AdminUsers)

	return map[string]interface{}{
		"totalUsers":    stats.TotalUsers,
		"activeUsers":   stats.ActiveUsers,
		"inactiveUsers": stats.InactiveUsers,
		"superAdmins":   stats.AdminUsers, // Actually counting super admins
	}, nil
}

// GetColumnMapping returns the column mapping for UserService
func (s *UserService) GetColumnMapping() map[string]string {
	mapping := make(map[string]string)
	mapping["createdAt"] = "created_at"
	mapping["updatedAt"] = "updated_at"
	return mapping
}

// Create creates a new user with manual validation to work around Goravel's max/min bug
func (s *UserService) Create(data map[string]interface{}) (interface{}, error) {
	// Manual validation for string length (workaround for Goravel bug)
	if err := s.validateUserData(data, true); err != nil {
		return nil, err
	}

	// Call the embedded service's Create method
	return s.CrudServiceContract.Create(data)
}

// Update updates a user with manual validation to work around Goravel's max/min bug
func (s *UserService) Update(id uint, data map[string]interface{}) (interface{}, error) {
	// Manual validation for string length (workaround for Goravel bug)
	if err := s.validateUserData(data, false); err != nil {
		return nil, err
	}

	// Call the embedded service's Update method
	return s.CrudServiceContract.Update(id, data)
}

// GetFilterDefinitions returns filter definitions for the users resource
func (s *UserService) GetFilterDefinitions() []contracts.FilterDefinition {
	return []contracts.FilterDefinition{
		// Name filter - uses all string operators automatically
		contracts.NewFilterDefinition(
			"name",
			"Name",
			contracts.FilterTypeString,
			nil, // Will use GetOperatorsForType(FilterTypeString)
		),
		// Email filter - uses all string operators automatically
		contracts.NewFilterDefinition(
			"email",
			"Email",
			contracts.FilterTypeString,
			nil, // Will use GetOperatorsForType(FilterTypeString)
		),
		// Status filter - uses all enum operators automatically
		contracts.NewFilterDefinition(
			"status",
			"Status",
			contracts.FilterTypeEnum,
			&[]string{
				"ACTIVE",
				"INACTIVE",
				"PENDING",
			},
		),
		// Created date filter - uses all datetime operators automatically
		contracts.NewFilterDefinition(
			"created_at",
			"Created Date",
			contracts.FilterTypeDateTime,
			nil, // Will use GetOperatorsForType(FilterTypeDateTime)
		),
	}
}

// validateUserData manually validates user data to work around Goravel's max/min validation bug
func (s *UserService) validateUserData(data map[string]interface{}, isCreate bool) error {
	// Validate name
	if name, exists := data["name"]; exists {
		nameStr, ok := name.(string)
		if !ok {
			return fmt.Errorf("name must be a string")
		}
		if isCreate && nameStr == "" {
			return fmt.Errorf("name is required")
		}
		if len(nameStr) < 2 {
			return fmt.Errorf("name min length is 2 characters")
		}
		if len(nameStr) > 255 {
			return fmt.Errorf("name max length is 255 characters")
		}
	} else if isCreate {
		return fmt.Errorf("name is required")
	}

	// Validate email
	if email, exists := data["email"]; exists {
		emailStr, ok := email.(string)
		if !ok {
			return fmt.Errorf("email must be a string")
		}
		if isCreate && emailStr == "" {
			return fmt.Errorf("email is required")
		}
		if emailStr != "" {
			// Basic email validation
			emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
			if !emailRegex.MatchString(emailStr) {
				return fmt.Errorf("email format is invalid")
			}
			// Check for duplicate email on create
			if isCreate {
				var existingUser models.User
				err := facades.Orm().Query().Where("email = ?", emailStr).First(&existingUser)
				if err == nil && existingUser.ID > 0 {
					return fmt.Errorf("email already exists")
				}
			}
		}
		if len(emailStr) > 255 {
			return fmt.Errorf("email max length is 255 characters")
		}
	} else if isCreate {
		return fmt.Errorf("email is required")
	}

	// Validate password
	if password, exists := data["password"]; exists {
		passwordStr, ok := password.(string)
		if !ok {
			return fmt.Errorf("password must be a string")
		}
		if isCreate && passwordStr == "" {
			return fmt.Errorf("password is required")
		}
		if passwordStr != "" && len(passwordStr) < 8 {
			return fmt.Errorf("password min length is 8 characters")
		}
	} else if isCreate {
		return fmt.Errorf("password is required")
	}

	return nil
}
