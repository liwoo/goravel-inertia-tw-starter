package services

import (
	"fmt"
	"players/app/contracts"
	"players/app/helpers"
	"players/app/models"
	"strconv"
	"strings"

	"github.com/goravel/framework/facades"
)

// BookService - using GORM directly with contract enforcement
type BookService struct {
	*contracts.BaseCrudService
	authHelper *helpers.AuthHelper
}

// NewBookService creates a new book service that implements all contracts
func NewBookService() *BookService {
	service := &BookService{
		BaseCrudService: contracts.NewBaseCrudService("books", "id"),
		authHelper:      helpers.NewAuthHelper().(*helpers.AuthHelper),
	}

	// Register service with validation
	contracts.MustRegisterCrudService("books", service)

	return service
}

// GetList with built-in pagination, sorting, filtering using GORM directly
// Implements CrudServiceContract interface
func (s *BookService) GetList(req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	// Use base service validation
	if err := s.ValidateListRequest(&req); err != nil {
		return nil, err
	}
	s.SanitizeListRequest(&req)
	
	// Debug logging
	fmt.Printf("BookService.GetList - Sort: '%s', Direction: '%s', Search: '%s'\n", req.Sort, req.Direction, req.Search)

	// Build query with sorting
	query := facades.Orm().Query().Model(&models.Book{})
	
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
				fmt.Printf("BookService.GetList - Applied sorting: %s\n", orderClause)
			} else {
				// Use default sort
				defaultField, defaultDir := s.GetDefaultSort()
				query = query.Order(defaultField + " " + defaultDir)
				fmt.Printf("BookService.GetList - Invalid sort mapping, using default\n")
			}
		} else {
			// Use default sort
			defaultField, defaultDir := s.GetDefaultSort()
			query = query.Order(defaultField + " " + defaultDir)
			fmt.Printf("BookService.GetList - Invalid sort field '%s', using default\n", req.Sort)
		}
	} else {
		// Default sorting
		defaultField, defaultDir := s.GetDefaultSort()
		query = query.Order(defaultField + " " + defaultDir)
		fmt.Printf("BookService.GetList - Using default sorting\n")
	}
	
	// Get all books with applied filters and sorting
	var allBooks []models.Book
	if err := query.Find(&allBooks); err != nil {
		return nil, err
	}

	// Manual pagination
	total := int64(len(allBooks))
	offset := (req.Page - 1) * req.PageSize
	end := offset + req.PageSize
	
	if offset > len(allBooks) {
		offset = len(allBooks)
	}
	if end > len(allBooks) {
		end = len(allBooks)
	}
	
	var pageBooks []models.Book
	if offset < len(allBooks) {
		pageBooks = allBooks[offset:end]
	}
	
	lastPage := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	
	// Convert to interface slice
	data := make([]interface{}, len(pageBooks))
	for i, book := range pageBooks {
		data[i] = book
	}

	return &contracts.PaginatedResult{
		Data:        data,
		Total:       total,
		PerPage:     req.PageSize,
		CurrentPage: req.Page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + len(pageBooks),
		HasNext:     req.Page < lastPage,
		HasPrev:     req.Page > 1,
	}, nil
}

// GetListAdvanced with additional filters using GORM directly
// Implements CrudServiceContract interface
func (s *BookService) GetListAdvanced(req contracts.ListRequest, filters map[string]interface{}) (*contracts.PaginatedResult, error) {
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
	countQuery := facades.Orm().Query().Model(&models.Book{})
	dataQuery := facades.Orm().Query().Model(&models.Book{})

	// Apply search to both queries if provided
	if req.Search != "" {
		searchCondition := "title LIKE ?"
		searchValue := "%" + req.Search + "%"
		countQuery = countQuery.Where(searchCondition, searchValue)
		dataQuery = dataQuery.Where(searchCondition, searchValue)
	}
	
	// Apply validated filters to both queries
	for field, value := range validatedFilters {
		var condition string
		switch field {
		case "status", "author":
			condition = field + " = ?"
		case "minPrice":
			condition = "price >= ?"
		case "maxPrice":
			condition = "price <= ?"
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
	var books []models.Book
	if err := dataQuery.Offset(offset).Limit(req.PageSize).Find(&books); err != nil {
		return nil, err
	}
	

	// Convert to interface slice
	data := make([]interface{}, len(books))
	for i, book := range books {
		data[i] = book
	}

	return &contracts.PaginatedResult{
		Data:        data,
		Total:       total,
		PerPage:     req.PageSize,
		CurrentPage: req.Page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + len(books),
		HasNext:     req.Page < lastPage,
		HasPrev:     req.Page > 1,
	}, nil
}

// GetByID - using GORM directly
// Implements CrudServiceContract interface
func (s *BookService) GetByID(id uint) (interface{}, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid ID: %d", id)
	}

	return s.getBookByID(id)
}

// getBookByID is a helper method that returns the actual model type
func (s *BookService) getBookByID(id uint) (*models.Book, error) {
	var book models.Book
	if err := facades.Orm().Query().Model(&models.Book{}).Where("id = ?", id).First(&book); err != nil {
		return nil, fmt.Errorf("book not found: %w", err)
	}

	return &book, nil
}

// GetByISBN retrieves a book by ISBN using GORM directly
func (s *BookService) GetByISBN(isbn string) (*models.Book, error) {
	var book models.Book
	if err := facades.Orm().Query().Model(&models.Book{}).Where("isbn = ?", isbn).First(&book); err != nil {
		return nil, fmt.Errorf("book not found with ISBN %s: %w", isbn, err)
	}

	return &book, nil
}

// GetByAuthor retrieves books by author using repository
func (s *BookService) GetByAuthor(author string, req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	filters := map[string]interface{}{
		"author": author,
	}
	return s.GetListAdvanced(req, filters)
}

// GetAvailable retrieves available books using repository
func (s *BookService) GetAvailable(req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	filters := map[string]interface{}{
		"status": "AVAILABLE",
	}
	return s.GetListAdvanced(req, filters)
}

// Create - using GORM directly with validation
// Implements CrudServiceContract interface
func (s *BookService) Create(data map[string]interface{}) (interface{}, error) {
	// Validate using validation rules
	if err := s.validateWithRules(data, false); err != nil {
		return nil, err
	}

	return s.createBook(data)
}

// createBook is a helper method that returns the actual model type
func (s *BookService) createBook(data map[string]interface{}) (*models.Book, error) {

	// Set default status if not provided
	if _, exists := data["status"]; !exists {
		data["status"] = "AVAILABLE"
	}

	// Create book struct from data
	book := models.Book{
		Title:       data["title"].(string),
		Author:      data["author"].(string),
		ISBN:        data["isbn"].(string),
		Status:      data["status"].(string),
	}

	if desc, ok := data["description"].(string); ok {
		book.Description = desc
	}
	if price, ok := data["price"].(float64); ok {
		book.Price = price
	}
	if published, ok := data["publishedAt"].(string); ok {
		book.PublishedAt = published
	}

	// Create using GORM
	if err := facades.Orm().Query().Create(&book); err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
	}

	return &book, nil
}

// Update - using GORM directly with validation
// Implements CrudServiceContract interface
func (s *BookService) Update(id uint, data map[string]interface{}) (interface{}, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid ID: %d", id)
	}

	// Validate using validation rules
	if err := s.validateWithRules(data, true); err != nil {
		return nil, err
	}

	return s.updateBook(id, data)
}

// updateBook is a helper method that returns the actual model type
func (s *BookService) updateBook(id uint, data map[string]interface{}) (*models.Book, error) {
	// Check if book exists
	_, err := s.getBookByID(id)
	if err != nil {
		return nil, err
	}

	// Apply column mapping to transform frontend field names to database column names
	columnMapping := s.GetColumnMapping()
	mappedData := make(map[string]interface{})
	
	// Fields to ignore (not saved to database)
	ignoredFields := map[string]bool{
		"tags": true, // Tags are not stored in the books table yet
	}
	
	for frontendField, value := range data {
		// Skip ignored fields
		if ignoredFields[frontendField] {
			continue
		}
		
		if dbColumn, exists := columnMapping[frontendField]; exists {
			mappedData[dbColumn] = value
		} else {
			// If no mapping exists, use the field name as-is
			mappedData[frontendField] = value
		}
	}

	// Update using GORM with properly mapped column names
	var book models.Book
	if _, err := facades.Orm().Query().Model(&book).Where("id = ?", id).Update(mappedData); err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}

	// Return updated book
	return s.getBookByID(id)
}

// Delete - using GORM directly
// Implements CrudServiceContract interface
func (s *BookService) Delete(id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid ID: %d", id)
	}

	// Check if book exists
	_, err := s.getBookByID(id)
	if err != nil {
		return err
	}

	// Delete using GORM (soft delete)
	if _, err := facades.Orm().Query().Model(&models.Book{}).Where("id = ?", id).Delete(&models.Book{}); err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	return nil
}

// BorrowBook marks a book as borrowed using GORM directly
func (s *BookService) BorrowBook(id uint) error {
	bookData, err := s.getBookByID(id)
	if err != nil {
		return err
	}

	if bookData.Status != "AVAILABLE" {
		return fmt.Errorf("book is not available for borrowing")
	}

	// Update status using GORM
	if _, err := facades.Orm().Query().Model(&models.Book{}).Where("id = ?", id).Update("status", "BORROWED"); err != nil {
		return fmt.Errorf("failed to update book status: %w", err)
	}

	return nil
}

// ReturnBook marks a book as available using GORM directly
func (s *BookService) ReturnBook(id uint) error {
	bookData, err := s.getBookByID(id)
	if err != nil {
		return err
	}

	if bookData.Status != "BORROWED" {
		return fmt.Errorf("book is not currently borrowed")
	}

	// Update status using GORM
	if _, err := facades.Orm().Query().Model(&models.Book{}).Where("id = ?", id).Update("status", "AVAILABLE"); err != nil {
		return fmt.Errorf("failed to update book status: %w", err)
	}

	return nil
}


// validateBookData performs simple validation
func (s *BookService) validateBookData(data map[string]interface{}, isUpdate bool) error {
	// Required fields for creation
	if !isUpdate {
		requiredFields := []string{"title", "author", "isbn"}
		for _, field := range requiredFields {
			if value, exists := data[field]; !exists || value == "" {
				return fmt.Errorf("%s is required", field)
			}
		}
	}

	// Validate ISBN format if provided
	if isbn, ok := data["isbn"].(string); ok && isbn != "" {
		if len(isbn) < 10 || len(isbn) > 17 {
			return fmt.Errorf("invalid ISBN format")
		}
	}

	// Validate status if provided
	if status, exists := data["status"]; exists {
		validStatuses := []string{"AVAILABLE", "BORROWED", "MAINTENANCE"}
		statusStr, ok := status.(string)
		if !ok {
			return fmt.Errorf("status must be a string")
		}

		valid := false
		for _, validStatus := range validStatuses {
			if statusStr == validStatus {
				valid = true
				break
			}
		}

		if !valid {
			return fmt.Errorf("status must be one of: AVAILABLE, BORROWED, MAINTENANCE")
		}
	}

	// Validate price if provided
	if price, exists := data["price"]; exists {
		switch v := price.(type) {
		case float64:
			if v < 0 {
				return fmt.Errorf("price cannot be negative")
			}
		case string:
			if p, err := strconv.ParseFloat(v, 64); err != nil || p < 0 {
				return fmt.Errorf("invalid price format")
			}
		default:
			return fmt.Errorf("price must be a number")
		}
	}

	return nil
}

// CONTRACT IMPLEMENTATIONS - Required by CompleteCrudService interface

// PaginationServiceContract implementation
func (s *BookService) GetPaginatedList(req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	return s.GetList(req)
}

// SortableServiceContract implementation
func (s *BookService) GetSortableFields() []string {
	return []string{"id", "title", "author", "isbn", "price", "status", "createdAt", "updatedAt", "publishedAt"}
}

func (s *BookService) ValidateSortField(field string) bool {
	sortableFields := s.GetSortableFields()
	for _, validField := range sortableFields {
		if field == validField {
			return true
		}
	}
	return false
}

func (s *BookService) MapSortField(frontendField string) (string, bool) {
	columnMapping := s.GetColumnMapping()
	if dbColumn, exists := columnMapping[frontendField]; exists {
		return dbColumn, true
	}
	return "", false
}

// FilterableServiceContract implementation
func (s *BookService) GetFilterableFields() []string {
	return []string{"status", "author", "minPrice", "maxPrice", "isbn"}
}

func (s *BookService) ValidateFilterField(field string) bool {
	filterableFields := s.GetFilterableFields()
	for _, validField := range filterableFields {
		if field == validField {
			return true
		}
	}
	return false
}

func (s *BookService) GetSearchableFields() []string {
	return []string{"title", "author", "description", "isbn"}
}

func (s *BookService) BuildFilterQuery(filters map[string]interface{}) (map[string]interface{}, error) {
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
func (s *BookService) Search(query string, req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	if err := s.ValidateSearchQuery(query); err != nil {
		return nil, err
	}
	
	req.Search = query
	return s.GetList(req)
}

func (s *BookService) ValidateSearchQuery(query string) error {
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
func (s *BookService) BulkCreate(data []map[string]interface{}) ([]interface{}, error) {
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

func (s *BookService) BulkUpdate(ids []uint, data map[string]interface{}) error {
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

func (s *BookService) BulkDelete(ids []uint) error {
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
func (s *BookService) GetModel() interface{} {
	return &models.Book{}
}

func (s *BookService) GetValidationRules() map[string]interface{} {
	return map[string]interface{}{
		"title":       "required|string|max:255",
		"author":      "required|string|max:255",
		"isbn":        "required|string|min:10|max:17",
		"description": "string|max:1000",
		"price":       "numeric|min:0",
		"status":      "in:AVAILABLE,BORROWED,MAINTENANCE",
		"publishedAt": "string",
	}
}

func (s *BookService) GetColumnMapping() map[string]string {
	return map[string]string{
		"id":          "id",
		"title":       "title",
		"author":      "author",
		"isbn":        "isbn",
		"price":       "price",
		"status":      "status",
		"description": "description",
		"createdAt":   "created_at",
		"updatedAt":   "updated_at",
		"publishedAt": "published_at",
		"created_at":  "created_at",
		"updated_at":  "updated_at",
		"published_at": "published_at",
	}
}

// HELPER METHODS

// validateWithRules uses the validation rules from the contract
func (s *BookService) validateWithRules(data map[string]interface{}, isUpdate bool) error {
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
	return s.validateBookData(data, isUpdate)
}