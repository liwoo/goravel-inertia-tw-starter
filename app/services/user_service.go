package services

import (
	"fmt"
	"strings"
	"regexp"
	"time"

	"github.com/goravel/framework/facades"
	"players/app/contracts"
	"players/app/models"
)

// UserService handles user business logic with contract enforcement
type UserService struct {
	*contracts.BaseCrudService
}

// NewUserService creates a new user service that implements all contracts
func NewUserService() *UserService {
	service := &UserService{
		BaseCrudService: contracts.NewBaseCrudService("user", "id"),
	}

	// Register service with validation
	contracts.MustRegisterCrudService("users", service)

	return service
}

// GetList with built-in pagination, sorting, filtering using GORM directly
// Implements CrudServiceContract interface
func (s *UserService) GetList(req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	// Use base service validation
	if err := s.ValidateListRequest(&req); err != nil {
		return nil, err
	}
	s.SanitizeListRequest(&req)

	// Build query
	query := facades.Orm().Query().Model(&models.User{}).With("Roles")

	// Apply search if provided using searchable fields
	if req.Search != "" {
		if err := s.ValidateSearchQuery(req.Search); err != nil {
			return nil, err
		}
		searchFields := s.GetSearchableFields()
		if len(searchFields) > 0 {
			conditions := make([]string, len(searchFields))
			values := make([]interface{}, len(searchFields))
			for i, field := range searchFields {
				conditions[i] = field + " LIKE ?"
				values[i] = "%" + req.Search + "%"
			}
			query = query.Where(strings.Join(conditions, " OR "), values...)
		}
	}

	// Apply sorting with field validation and mapping
	if req.Sort != "" && req.Direction != "" {
		if s.ValidateSortField(req.Sort) && s.ValidateSortDirection(req.Direction) {
			dbColumn, valid := s.MapSortField(req.Sort)
			if valid {
				orderClause := dbColumn + " " + strings.ToUpper(req.Direction)
				query = query.Order(orderClause)
			} else {
				// Use default sort
				defaultField, defaultDir := s.GetDefaultSort()
				query = query.Order(defaultField + " " + defaultDir)
			}
		} else {
			// Use default sort
			defaultField, defaultDir := s.GetDefaultSort()
			query = query.Order(defaultField + " " + defaultDir)
		}
	} else {
		// Default sorting
		defaultField, defaultDir := s.GetDefaultSort()
		query = query.Order(defaultField + " " + defaultDir)
	}

	// Get all users with applied filters and sorting
	var allUsers []models.User
	if err := query.Find(&allUsers); err != nil {
		return nil, err
	}

	// Manual pagination
	total := int64(len(allUsers))
	offset := (req.Page - 1) * req.PageSize
	end := offset + req.PageSize

	if offset > len(allUsers) {
		offset = len(allUsers)
	}
	if end > len(allUsers) {
		end = len(allUsers)
	}

	var pageUsers []models.User
	if offset < len(allUsers) {
		pageUsers = allUsers[offset:end]
	}

	lastPage := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	// Convert to interface slice
	data := make([]interface{}, len(pageUsers))
	for i, user := range pageUsers {
		data[i] = user
	}

	return &contracts.PaginatedResult{
		Data:        data,
		Total:       total,
		PerPage:     req.PageSize,
		CurrentPage: req.Page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + len(pageUsers),
		HasNext:     req.Page < lastPage,
		HasPrev:     req.Page > 1,
	}, nil
}

// GetListAdvanced with additional filters using GORM directly
// Implements CrudServiceContract interface
func (s *UserService) GetListAdvanced(req contracts.ListRequest, filters map[string]interface{}) (*contracts.PaginatedResult, error) {
	// Validate and sanitize request
	if err := s.ValidateListRequest(&req); err != nil {
		return nil, err
	}
	s.SanitizeListRequest(&req)

	// Validate filters
	validatedFilters, err := s.BuildFilterQuery(filters)
	if err != nil {
		return nil, err
	}

	// Create separate queries for count and data
	countQuery := facades.Orm().Query().Model(&models.User{})
	dataQuery := facades.Orm().Query().Model(&models.User{}).With("Roles")

	// Apply search to both queries if provided
	if req.Search != "" {
		searchCondition := "name LIKE ? OR email LIKE ?"
		searchValue := "%" + req.Search + "%"
		countQuery = countQuery.Where(searchCondition, searchValue, searchValue)
		dataQuery = dataQuery.Where(searchCondition, searchValue, searchValue)
	}

	// Apply validated filters to both queries
	for field, value := range validatedFilters {
		var condition string
		switch field {
		case "is_active":
			condition = "is_active = ?"
		case "is_super_admin":
			condition = "is_super_admin = ?"
		case "role":
			// Filter by role slug
			countQuery = countQuery.Where("EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = users.id AND r.slug = ?)", value)
			dataQuery = dataQuery.Where("EXISTS (SELECT 1 FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = users.id AND r.slug = ?)", value)
			continue
		default:
			continue
		}
		countQuery = countQuery.Where(condition, value)
		dataQuery = dataQuery.Where(condition, value)
	}

	// Count total records
	var total int64
	if err := countQuery.Count(&total); err != nil {
		return nil, err
	}

	// Add sorting to data query only
	if req.Sort != "" {
		dataQuery = dataQuery.Order("id DESC") // Should parse req.Sort properly
	} else {
		dataQuery = dataQuery.Order("id DESC")
	}

	// Calculate pagination
	offset := (req.Page - 1) * req.PageSize
	lastPage := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	// Get paginated data
	var users []models.User
	if err := dataQuery.Offset(offset).Limit(req.PageSize).Find(&users); err != nil {
		return nil, err
	}

	// Convert to interface slice
	data := make([]interface{}, len(users))
	for i, user := range users {
		data[i] = user
	}

	return &contracts.PaginatedResult{
		Data:        data,
		Total:       total,
		PerPage:     req.PageSize,
		CurrentPage: req.Page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + len(users),
		HasNext:     req.Page < lastPage,
		HasPrev:     req.Page > 1,
	}, nil
}

// GetByID - Implements CrudServiceContract interface
func (s *UserService) GetByID(id uint) (interface{}, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid ID: %d", id)
	}

	return s.getUserByID(id)
}

// getUserByID is a helper method that returns the actual model type
func (s *UserService) getUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := facades.Orm().Query().Model(&models.User{}).With("Roles").Where("id = ?", id).First(&user); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user, nil
}

// Create - Implements CrudServiceContract interface
func (s *UserService) Create(data map[string]interface{}) (interface{}, error) {
	// Validate using validation rules
	if err := s.validateWithRules(data, false); err != nil {
		return nil, err
	}

	return s.createUser(data)
}

// createUser is a helper method that returns the actual model type
func (s *UserService) createUser(data map[string]interface{}) (*models.User, error) {
	// Basic validation
	if err := s.validateUserData(data, false); err != nil {
		return nil, err
	}

	// Check if email already exists (GORM automatically excludes soft-deleted users)
	var existingCount int64
	err := facades.Orm().Query().Model(&models.User{}).Where("email = ?", data["email"].(string)).Count(&existingCount)
	if err != nil {
		return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
	}
	if existingCount > 0 {
		return nil, fmt.Errorf("email already exists")
	}

	// Set default values if not provided
	if _, exists := data["is_active"]; !exists {
		data["is_active"] = true
	}
	if _, exists := data["is_super_admin"]; !exists {
		data["is_super_admin"] = false
	}

	// Create user struct from data
	user := models.User{
		Name:         data["name"].(string),
		Email:        data["email"].(string),
		IsActive:     data["is_active"].(bool),
		IsSuperAdmin: data["is_super_admin"].(bool),
	}

	// Hash password if provided
	if password, ok := data["password"].(string); ok && password != "" {
		hashedPassword, err := facades.Hash().Make(password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
	}

	// Create using GORM
	if err := facades.Orm().Query().Create(&user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Assign role if provided
	if roleID, ok := data["role_id"].(float64); ok && roleID > 0 {
		userRole := models.UserRole{
			UserID:     user.ID,
			RoleID:     uint(roleID),
			AssignedAt: time.Now(),
			IsActive:   true,
		}
		if err := facades.Orm().Query().Create(&userRole); err != nil {
			// Log error but don't fail user creation
			facades.Log().Error("Failed to assign role to user", map[string]interface{}{
				"user_id": user.ID,
				"role_id": roleID,
				"error":   err.Error(),
			})
		}
	}

	// Reload user with roles
	if err := facades.Orm().Query().Model(&models.User{}).With("Roles").Where("id = ?", user.ID).First(&user); err != nil {
		facades.Log().Error("Failed to reload user with roles", map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		})
	}

	return &user, nil
}

// Update - Implements CrudServiceContract interface
func (s *UserService) Update(id uint, data map[string]interface{}) (interface{}, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid ID: %d", id)
	}

	// Validate using validation rules
	if err := s.validateWithRules(data, true); err != nil {
		return nil, err
	}

	return s.updateUser(id, data)
}

// updateUser is a helper method that returns the actual model type
func (s *UserService) updateUser(id uint, data map[string]interface{}) (*models.User, error) {
	// Check if user exists
	user, err := s.getUserByID(id)
	if err != nil {
		return nil, err
	}

	// Check if email is being changed and already exists
	if email, ok := data["email"].(string); ok && email != user.Email {
		var existingCount int64
		err := facades.Orm().Query().Model(&models.User{}).Where("email = ? AND id != ?", email, id).Count(&existingCount)
		if err != nil {
			return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
		}
		if existingCount > 0 {
			return nil, fmt.Errorf("email already exists")
		}
	}

	// Hash password if provided
	if password, ok := data["password"].(string); ok && password != "" {
		hashedPassword, err := facades.Hash().Make(password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		data["password"] = hashedPassword
	} else {
		// Remove password from update if empty
		delete(data, "password")
	}

	// Update using GORM
	if _, err := facades.Orm().Query().Model(&user).Where("id = ?", id).Update(data); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Update role if provided
	if roleID, ok := data["role_id"].(float64); ok {
		// Remove existing roles
		facades.Orm().Query().Where("user_id = ?", id).Delete(&models.UserRole{})
		
		// Assign new role
		userRole := models.UserRole{
			UserID:     id,
			RoleID:     uint(roleID),
			AssignedAt: time.Now(),
			IsActive:   true,
		}
		if err := facades.Orm().Query().Create(&userRole); err != nil {
			// Log error but don't fail user update
			facades.Log().Error("Failed to update user role", map[string]interface{}{
				"user_id": id,
				"role_id": roleID,
				"error":   err.Error(),
			})
		}
	}

	// Return updated user
	return s.getUserByID(id)
}

// Delete - Implements CrudServiceContract interface
func (s *UserService) Delete(id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid ID: %d", id)
	}

	// Check if user exists
	_, err := s.getUserByID(id)
	if err != nil {
		return err
	}

	// Delete using GORM (soft delete)
	if _, err := facades.Orm().Query().Model(&models.User{}).Where("id = ?", id).Delete(&models.User{}); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// GetAllRoles returns all available roles for assignment
func (s *UserService) GetAllRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := facades.Orm().Query().Find(&roles); err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	return roles, nil
}

// CONTRACT IMPLEMENTATIONS - Required by CompleteCrudService interface

// PaginationServiceContract implementation
func (s *UserService) GetPaginatedList(req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	return s.GetList(req)
}

// SortableServiceContract implementation
func (s *UserService) GetSortableFields() []string {
	return []string{"id", "name", "email", "is_active", "is_super_admin", "createdAt", "updatedAt"}
}

func (s *UserService) ValidateSortField(field string) bool {
	sortableFields := s.GetSortableFields()
	for _, validField := range sortableFields {
		if field == validField {
			return true
		}
	}
	return false
}

func (s *UserService) MapSortField(frontendField string) (string, bool) {
	columnMapping := s.GetColumnMapping()
	if dbColumn, exists := columnMapping[frontendField]; exists {
		return dbColumn, true
	}
	return "", false
}

// FilterableServiceContract implementation
func (s *UserService) GetFilterableFields() []string {
	return []string{"name", "email", "is_active", "is_super_admin", "role"}
}

func (s *UserService) ValidateFilterField(field string) bool {
	filterableFields := s.GetFilterableFields()
	for _, validField := range filterableFields {
		if field == validField {
			return true
		}
	}
	return false
}

func (s *UserService) GetSearchableFields() []string {
	return []string{"name", "email"}
}

func (s *UserService) BuildFilterQuery(filters map[string]interface{}) (map[string]interface{}, error) {
	validatedFilters := make(map[string]interface{})

	for field, value := range filters {
		if !s.ValidateFilterField(field) {
			continue // Skip invalid fields
		}

		if !s.ValidateFilterValue(field, value) {
			continue // Skip invalid values
		}

		validatedFilters[field] = value
	}

	return validatedFilters, nil
}

// SearchableServiceContract implementation
func (s *UserService) Search(query string, req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	if err := s.ValidateSearchQuery(query); err != nil {
		return nil, err
	}

	req.Search = query
	return s.GetList(req)
}

func (s *UserService) ValidateSearchQuery(query string) error {
	query = strings.TrimSpace(query)
	if len(query) < 2 {
		return fmt.Errorf("search query must be at least 2 characters")
	}
	if len(query) > 100 {
		return fmt.Errorf("search query cannot exceed 100 characters")
	}
	return nil
}

// BulkOperationsContract implementation
func (s *UserService) BulkCreate(data []map[string]interface{}) ([]interface{}, error) {
	if err := s.ValidateBulkOperation([]uint{uint(len(data))}); err != nil {
		return nil, err
	}

	results := make([]interface{}, 0, len(data))
	for _, item := range data {
		result, err := s.Create(item)
		if err != nil {
			return nil, fmt.Errorf("bulk create failed: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}

func (s *UserService) BulkUpdate(ids []uint, data map[string]interface{}) error {
	if err := s.ValidateBulkOperation(ids); err != nil {
		return err
	}

	for _, id := range ids {
		_, err := s.Update(id, data)
		if err != nil {
			return fmt.Errorf("bulk update failed for ID %d: %w", id, err)
		}
	}

	return nil
}

func (s *UserService) BulkDelete(ids []uint) error {
	if err := s.ValidateBulkOperation(ids); err != nil {
		return err
	}

	for _, id := range ids {
		if err := s.Delete(id); err != nil {
			return fmt.Errorf("bulk delete failed for ID %d: %w", id, err)
		}
	}

	return nil
}

// CrudServiceConfiguration implementation
func (s *UserService) GetModel() interface{} {
	return &models.User{}
}

func (s *UserService) GetValidationRules() map[string]interface{} {
	return map[string]interface{}{
		"name":           "required|string|max:255",
		"email":          "required|email|max:255",
		"password":       "string|min:8",
		"is_active":      "boolean",
		"is_super_admin": "boolean",
		"role_id":        "numeric",
	}
}

func (s *UserService) GetColumnMapping() map[string]string {
	return map[string]string{
		"id":            "id",
		"name":          "name",
		"email":         "email",
		"isActive":      "is_active",
		"isSuperAdmin":  "is_super_admin",
		"createdAt":     "created_at",
		"updatedAt":     "updated_at",
		"created_at":    "created_at",
		"updated_at":    "updated_at",
		"is_active":     "is_active",
		"is_super_admin": "is_super_admin",
	}
}

// HELPER METHODS

// validateWithRules uses the validation rules from the contract
func (s *UserService) validateWithRules(data map[string]interface{}, isUpdate bool) error {
	rules := s.GetValidationRules()

	// For updates, make required fields optional
	if isUpdate {
		for field, rule := range rules {
			if ruleStr, ok := rule.(string); ok {
				// Remove 'required|' from validation rules for updates
				if strings.HasPrefix(ruleStr, "required|") {
					rules[field] = strings.TrimPrefix(ruleStr, "required|")
				}
			}
		}
	}

	// Basic validation implementation (can be enhanced with proper validator)
	return s.validateUserData(data, isUpdate)
}

// validateUserData performs simple validation
func (s *UserService) validateUserData(data map[string]interface{}, isUpdate bool) error {
	// Required fields for creation
	if !isUpdate {
		requiredFields := []string{"name", "email"}
		for _, field := range requiredFields {
			if value, exists := data[field]; !exists || value == "" {
				return fmt.Errorf("%s is required", field)
			}
		}
		
		// Password is required for new users
		if password, exists := data["password"]; !exists || password == "" {
			return fmt.Errorf("password is required for new users")
		}
	}

	// Validate name if provided
	if name, ok := data["name"].(string); ok && name != "" {
		if len(name) < 2 || len(name) > 255 {
			return fmt.Errorf("name must be between 2 and 255 characters")
		}
	}

	// Validate email if provided
	if email, ok := data["email"].(string); ok && email != "" {
		// Simple email regex validation
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(email) {
			return fmt.Errorf("invalid email format")
		}
		if len(email) > 255 {
			return fmt.Errorf("email cannot exceed 255 characters")
		}
	}

	// Validate password if provided
	if password, ok := data["password"].(string); ok && password != "" {
		if len(password) < 8 {
			return fmt.Errorf("password must be at least 8 characters")
		}
	}

	return nil
}