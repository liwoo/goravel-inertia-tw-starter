package services

import (
	"fmt"
	"strings"
	
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/facades"
	"players/app/contracts"
	"players/app/models"
)

// RoleService provides business logic for roles
type RoleService struct {
	*contracts.GenericCrudService[models.Role]
}

// NewRoleService creates a new instance of RoleService
func NewRoleService() *RoleService {
	// Create the generic service
	genericService := contracts.NewGenericCrudService[models.Role]("roles", "id")

	// Configure the service
	genericService.
		SetSearchFields("name", "description", "slug").
		SetSortFields("id", "name", "slug", "level", "is_active", "created_at", "updated_at").
		SetFilterFields("is_active", "type", "level").
		SetCustomFilters(func(query orm.Query, filters map[string]interface{}) orm.Query {
			for key, value := range filters {
				switch key {
				case "is_active":
					// Handle boolean conversion for is_active
					switch v := value.(type) {
					case bool:
						query = query.Where("is_active = ?", v)
					case string:
						if v == "true" {
							query = query.Where("is_active = ?", true)
						} else if v == "false" {
							query = query.Where("is_active = ?", false)
						}
					}
				case "type":
					// Handle role type filtering
					switch value {
					case "super_admin":
						query = query.Where("slug = ?", "super-admin")
					case "admin":
						query = query.Where("slug = ?", "admin")
					case "user":
						// Filter for user-type roles (not admin or super-admin)
						query = query.Where("slug NOT IN ?", []string{"super-admin", "admin"})
					}
				case "level":
					query = query.Where("level = ?", value)
				}
			}
			return query
		}).
		SetBeforeCreate(func(data map[string]interface{}) error {
			// Validate required fields before creation
			name, nameOk := data["name"].(string)
			if !nameOk || strings.TrimSpace(name) == "" {
				return fmt.Errorf("role name is required and cannot be empty")
			}
			
			slug, slugOk := data["slug"].(string) 
			if !slugOk || strings.TrimSpace(slug) == "" {
				// Auto-generate slug from name if not provided
				slug = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "-"))
				if slug == "" {
					return fmt.Errorf("role slug is required and cannot be empty")
				}
				data["slug"] = slug
			}
			
			// Set default status if not provided
			if _, exists := data["is_active"]; !exists {
				data["is_active"] = true
			}
			
			// Set created_by if provided (from context)
			if createdBy, exists := data["created_by"]; exists && createdBy != nil {
				data["created_by"] = createdBy
			}
			
			return nil
		})

	service := &RoleService{
		GenericCrudService: genericService,
	}

	// Set the actual service reference for method resolution
	genericService.SetActualService(service)
	
	// Set a custom query to show all roles (including inactive ones) by default
	// This prevents the issue where deleted roles "disappear" from the list
	genericService.SetCustomQuery(func(query orm.Query) orm.Query {
		// Don't filter by is_active by default - let the frontend handle this
		return query
	})

	return service
}

// List returns paginated roles
func (s *RoleService) List(page, pageSize int, sort, direction, search string, filters map[string]interface{}) (*contracts.PaginatedResult, error) {
	// Build list request
	req := contracts.ListRequest{
		Page:      page,
		PageSize:  pageSize,
		Sort:      sort,
		Direction: direction,
		Search:    search,
		Filters:   filters,
	}

	// Use the generic service's GetList method
	result, err := s.GenericCrudService.GetList(req)
	if err != nil {
		return nil, err
	}
	
	// Process each item to add computed fields
	processedData := make([]interface{}, len(result.Data))
	for i, item := range result.Data {
		processedData[i] = s.ProcessListItem(item)
	}
	result.Data = processedData
	
	return result, nil
}

// MapSortField maps frontend field names to database column names
func (s *RoleService) MapSortField(frontendField string) (string, bool) {
	// Map camelCase to snake_case
	fieldMap := map[string]string{
		"createdAt": "created_at",
		"updatedAt": "updated_at",
		"isActive":  "is_active",
	}

	if dbField, exists := fieldMap[frontendField]; exists {
		return dbField, true
	}

	// For other fields, check if they're valid
	validFields := []string{"id", "name", "slug", "description", "level", "is_active", "created_at", "updated_at"}
	for _, field := range validFields {
		if field == frontendField {
			return frontendField, true
		}
	}

	return "", false
}

// ProcessListItem processes each role item for the list view
func (s *RoleService) ProcessListItem(item interface{}) interface{} {
	role := item.(*models.Role)

	// Count users with this role
	var userCount int64
	facades.Orm().Query().Model(&models.UserRole{}).
		Where("role_id = ? AND is_active = ?", role.ID, true).
		Count(&userCount)

	// Return role with additional computed fields
	return map[string]interface{}{
		"id":          role.ID,
		"name":        role.Name,
		"slug":        role.Slug,
		"description": role.Description,
		"level":       role.Level,
		"is_active":   role.IsActive,
		"users_count": userCount,
		"created_at":  role.CreatedAt,
		"updated_at":  role.UpdatedAt,
	}
}

// ValidateCreate validates role creation
func (s *RoleService) ValidateCreate(data map[string]interface{}) error {
	// Validate name
	name, ok := data["name"].(string)
	if !ok || strings.TrimSpace(name) == "" {
		return fmt.Errorf("role name is required")
	}
	
	// Validate slug
	slug, ok := data["slug"].(string)
	if !ok || strings.TrimSpace(slug) == "" {
		return fmt.Errorf("role slug is required")
	}
	
	// Check if role with this slug already exists
	var existingRole models.Role
	err := facades.Orm().Query().Where("slug = ?", slug).First(&existingRole)
	if err == nil && existingRole.ID > 0 {
		return fmt.Errorf("a role with this slug already exists")
	}
	
	// Validate level if provided
	if level, ok := data["level"]; ok {
		switch v := level.(type) {
		case float64:
			if v < 0 || v > 100 {
				return fmt.Errorf("level must be between 0 and 100")
			}
		case int:
			if v < 0 || v > 100 {
				return fmt.Errorf("level must be between 0 and 100")
			}
		}
	}
	
	return nil
}

// ValidateUpdate validates role updates  
func (s *RoleService) ValidateUpdate(id uint, data map[string]interface{}) error {
	// If name is being updated, validate it
	if name, ok := data["name"]; ok {
		nameStr, ok := name.(string)
		if !ok || strings.TrimSpace(nameStr) == "" {
			return fmt.Errorf("role name cannot be empty")
		}
	}
	
	// If slug is being updated, validate it
	if slug, ok := data["slug"]; ok {
		slugStr, ok := slug.(string)
		if !ok || strings.TrimSpace(slugStr) == "" {
			return fmt.Errorf("role slug cannot be empty")
		}
		
		// Check if another role already has this slug
		var existingRole models.Role
		err := facades.Orm().Query().
			Where("slug = ? AND id != ?", slugStr, id).
			First(&existingRole)
		if err == nil && existingRole.ID > 0 {
			return fmt.Errorf("a role with this slug already exists")
		}
	}
	
	// Validate level if provided
	if level, ok := data["level"]; ok {
		switch v := level.(type) {
		case float64:
			if v < 0 || v > 100 {
				return fmt.Errorf("level must be between 0 and 100")
			}
		case int:
			if v < 0 || v > 100 {
				return fmt.Errorf("level must be between 0 and 100")
			}
		}
	}
	
	return nil
}

// GetSearchableFields returns the fields that can be searched
func (s *RoleService) GetSearchableFields() []string {
	return []string{"name", "description", "slug"}
}

// GetSortableFields returns the fields that can be sorted
func (s *RoleService) GetSortableFields() []string {
	return []string{"id", "name", "slug", "level", "is_active", "created_at", "updated_at"}
}