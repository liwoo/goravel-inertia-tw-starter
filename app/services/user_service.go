package services

import (
	"fmt"
	"time"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/facades"
	"players/app/contracts"
	"players/app/models"
)

// UserService - Simplified version using generic CRUD service
// From ~600 lines to ~150 lines!
type UserService struct {
	*contracts.GenericCrudService[models.User]
}

// NewUserService creates a new simplified user service
func NewUserService() *UserService {
	// Create the generic service
	genericService := contracts.NewGenericCrudService[models.User]("user", "id")

	// Configure the service
	genericService.
		SetSearchFields("name", "email").
		SetSortFields("id", "name", "email", "is_active", "is_super_admin", "created_at", "updated_at").
		SetFilterFields("is_active", "is_super_admin", "role").
		SetRelations("Roles").
		SetValidationRules(map[string]interface{}{
			"name":           "required|string|max:255",
			"email":          "required|email|max:255",
			"password":       "string|min:8",
			"is_active":      "boolean",
			"is_super_admin": "boolean",
			"role_id":        "numeric",
		}).
		SetBeforeCreate(func(data map[string]interface{}) error {
			// Set defaults
			if _, exists := data["is_active"]; !exists {
				data["is_active"] = true
			}
			if _, exists := data["is_super_admin"]; !exists {
				data["is_super_admin"] = false
			}

			// Check email uniqueness
			var count int64
			err := facades.Orm().Query().Model(&models.User{}).
				Where("email = ?", data["email"]).
				Count(&count)
			if err != nil {
				return fmt.Errorf("failed to check email uniqueness: %w", err)
			}
			if count > 0 {
				return fmt.Errorf("email already exists")
			}

			// Hash password if provided
			if password, ok := data["password"].(string); ok && password != "" {
				hashedPassword, err := facades.Hash().Make(password)
				if err != nil {
					return fmt.Errorf("failed to hash password: %w", err)
				}
				data["password"] = hashedPassword
			}

			return nil
		}).
		SetAfterCreate(func(user *models.User) error {
			// Assign role if provided in the original data
			// Note: We'd need to pass role_id through context or handle separately
			return nil
		}).
		SetBeforeUpdate(func(id uint, data map[string]interface{}) error {
			// Get existing user
			var existingUser models.User
			err := facades.Orm().Query().Model(&models.User{}).
				Where("id = ?", id).
				First(&existingUser)
			if err != nil {
				return err
			}

			// Check email uniqueness if being changed
			if email, ok := data["email"].(string); ok && email != existingUser.Email {
				var count int64
				err := facades.Orm().Query().Model(&models.User{}).
					Where("email = ? AND id != ?", email, id).
					Count(&count)
				if err != nil {
					return fmt.Errorf("failed to check email uniqueness: %w", err)
				}
				if count > 0 {
					return fmt.Errorf("email already exists")
				}
			}

			// Hash password if provided
			if password, ok := data["password"].(string); ok && password != "" {
				hashedPassword, err := facades.Hash().Make(password)
				if err != nil {
					return fmt.Errorf("failed to hash password: %w", err)
				}
				data["password"] = hashedPassword
			} else {
				// Remove password from update if empty
				delete(data, "password")
			}

			return nil
		}).
		SetCustomSearch(func(query orm.Query, search string) orm.Query {
			searchValue := "%" + search + "%"
			return query.Where("name LIKE ? OR email LIKE ?", searchValue, searchValue)
		}).
		SetCustomFilters(func(query orm.Query, filters map[string]interface{}) orm.Query {
			fmt.Printf("DEBUG UserService.CustomFilters: Received filters=%+v\n", filters)
			for field, value := range filters {
				switch field {
				case "is_active":
					// Convert string to boolean
					fmt.Printf("DEBUG UserService.CustomFilters: Processing is_active with value='%v' (type=%T)\n", value, value)
					switch v := value.(type) {
					case string:
						if v == "true" {
							query = query.Where("is_active = ?", true)
						} else if v == "false" {
							query = query.Where("is_active = ?", false)
						}
					case bool:
						query = query.Where("is_active = ?", v)
					}
				case "is_super_admin":
					// Convert string to boolean
					switch v := value.(type) {
					case string:
						if v == "true" {
							query = query.Where("is_super_admin = ?", true)
						} else if v == "false" {
							query = query.Where("is_super_admin = ?", false)
						}
					case bool:
						query = query.Where("is_super_admin = ?", v)
					}
				case "role":
					// Filter by role slug
					query = query.Where("EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = users.id AND r.slug = ?)", value)
				case "level_min":
					// For super admin filtering - check if user is super admin
					if level, ok := value.(string); ok {
						if level == "90" {
							// Filter for super admins
							query = query.Where("is_super_admin = ?", true)
						}
					}
				case "level_max":
					// For role level filtering - not implemented for now
					// Would need a different approach without joins
				}
			}
			return query
		})

	service := &UserService{
		GenericCrudService: genericService,
	}

	// Set the actual service reference so method resolution works correctly
	genericService.SetActualService(service)

	// Register service
	contracts.MustRegisterCrudService("users", service)

	return service
}

// Override Create to handle role assignment
func (s *UserService) Create(data map[string]interface{}) (interface{}, error) {
	// Extract role_id before creating user
	var roleID *uint
	if rid, ok := data["role_id"].(float64); ok && rid > 0 {
		roleIDVal := uint(rid)
		roleID = &roleIDVal
		// Don't save role_id to user table
		delete(data, "role_id")
	}

	// Create user using generic implementation
	result, err := s.GenericCrudService.Create(data)
	if err != nil {
		return nil, err
	}

	// Assign role if provided
	if roleID != nil {
		user := result.(models.User)
		userRole := models.UserRole{
			UserID:     user.ID,
			RoleID:     *roleID,
			AssignedAt: time.Now(),
			IsActive:   true,
		}
		if err := facades.Orm().Query().Create(&userRole); err != nil {
			// Log error but don't fail user creation
			facades.Log().Error("Failed to assign role to user", map[string]interface{}{
				"user_id": user.ID,
				"role_id": *roleID,
				"error":   err.Error(),
			})
		}
	}

	// Reload with roles
	return s.GetByID(result.(models.User).ID)
}

// Override Update to handle role changes
func (s *UserService) Update(id uint, data map[string]interface{}) (interface{}, error) {
	// Extract role_id before updating user
	var roleID *uint
	if rid, ok := data["role_id"].(float64); ok {
		roleIDVal := uint(rid)
		roleID = &roleIDVal
		// Don't save role_id to user table
		delete(data, "role_id")
	}

	// Update user using generic implementation
	_, err := s.GenericCrudService.Update(id, data)
	if err != nil {
		return nil, err
	}

	// Update role if provided
	if roleID != nil {
		// Remove existing roles
		facades.Orm().Query().Where("user_id = ?", id).Delete(&models.UserRole{})

		// Assign new role
		userRole := models.UserRole{
			UserID:     id,
			RoleID:     *roleID,
			AssignedAt: time.Now(),
			IsActive:   true,
		}
		if err := facades.Orm().Query().Create(&userRole); err != nil {
			// Log error but don't fail user update
			facades.Log().Error("Failed to update user role", map[string]interface{}{
				"user_id": id,
				"role_id": *roleID,
				"error":   err.Error(),
			})
		}
	}

	// Reload with roles
	return s.GetByID(id)
}

// GetAllRoles returns all available roles for assignment
func (s *UserService) GetAllRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := facades.Orm().Query().Find(&roles); err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	return roles, nil
}

// GetColumnMapping returns database column mappings
func (s *UserService) GetColumnMapping() map[string]string {
	return map[string]string{
		"id":             "id",
		"name":           "name",
		"email":          "email",
		"isActive":       "is_active",
		"isSuperAdmin":   "is_super_admin",
		"createdAt":      "created_at",
		"updatedAt":      "updated_at",
		"created_at":     "created_at",
		"updated_at":     "updated_at",
		"is_active":      "is_active",
		"is_super_admin": "is_super_admin",
	}
}

// MapSortField maps frontend field names to database field names
// This handles camelCase to snake_case conversion
func (s *UserService) MapSortField(frontendField string) (string, bool) {
	// Map camelCase fields to snake_case
	fieldMap := map[string]string{
		"isActive":     "is_active",
		"isSuperAdmin": "is_super_admin",
		"createdAt":    "created_at",
		"updatedAt":    "updated_at",
	}

	// Check if we have a mapping
	if dbField, exists := fieldMap[frontendField]; exists {
		// Validate the mapped field
		if s.ValidateSortField(dbField) {
			return dbField, true
		}
	}

	// Otherwise delegate to the base implementation
	return s.GenericCrudService.MapSortField(frontendField)
}
