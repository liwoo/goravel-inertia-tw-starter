package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/facades"
	"players/app/contracts"
	"players/app/models"
)

// BookService implements book-specific business logic using the builder pattern
type BookService struct {
	contracts.CrudServiceContract
	baseService contracts.CrudServiceContract
}

// NewBookService creates a new book service using the builder pattern
func NewBookService() *BookService {
	// Build the service with all required configurations
	service := contracts.NewServiceBuilder[models.Book]("books", "id").
		WithSearchFields("title", "author", "isbn", "description").                                   // REQUIRED
		WithSortFields("id", "title", "author", "price", "created_at", "updated_at", "published_at"). // REQUIRED
		WithFilterFields("status", "author").                                                         // REQUIRED
		WithValidationRules(map[string]interface{}{                                                   // REQUIRED
			"title":       "required|string|max:255",
			"author":      "required|string|max:100",
			"isbn":        "required|string|max:20",
			"status":      "required|string|in:AVAILABLE,BORROWED,MAINTENANCE,RESERVED",
			"price":       "numeric|min:0",
			"publishedAt": "date",
			"description": "string|max:1000",
			"tags":        "array",
			"tags.*":      "string|max:50",
		}).
		WithRelations("Creator", "Updater").                       // Optional
		WithDefaultSort("created_at", "DESC").                     // Optional
		WithSoftDeletes().                                         // Optional
		WithScopeFiltering("books", "created_by").                 // Optional
		WithBeforeCreate(func(data map[string]interface{}) error { // Optional
			// The created_by field should already be set by the controller
			// Just ensure it's properly formatted if present
			if createdBy, exists := data["created_by"]; exists && createdBy != nil {
				// Ensure it's a valid uint
				switch v := createdBy.(type) {
				case float64:
					data["created_by"] = uint(v)
				case int:
					data["created_by"] = uint(v)
				case uint:
					// Already correct type
				default:
					// Try to convert
					if fmt.Sprintf("%v", v) != "" {
						// Keep the value as is, let GORM handle the conversion
					}
				}
			}

			// Handle tags array to JSON conversion
			if tags, exists := data["tags"]; exists && tags != nil {
				if tagsArray, ok := tags.([]interface{}); ok && len(tagsArray) > 0 {
					// Convert []interface{} to []string
					stringTags := make([]string, len(tagsArray))
					for i, tag := range tagsArray {
						stringTags[i] = fmt.Sprintf("%v", tag)
					}
					tagsJSON, _ := json.Marshal(stringTags)
					data["tags"] = string(tagsJSON)
				} else if tagsArray, ok := tags.([]string); ok && len(tagsArray) > 0 {
					tagsJSON, _ := json.Marshal(tagsArray)
					data["tags"] = string(tagsJSON)
				} else {
					data["tags"] = "[]"
				}
			} else {
				data["tags"] = "[]"
			}

			return nil
		}).
		WithBeforeUpdate(func(id uint, data map[string]interface{}) error { // Optional
			// Handle tags array to JSON conversion
			if tags, exists := data["tags"]; exists && tags != nil {
				if tagsArray, ok := tags.([]interface{}); ok && len(tagsArray) > 0 {
					// Convert []interface{} to []string
					stringTags := make([]string, len(tagsArray))
					for i, tag := range tagsArray {
						stringTags[i] = fmt.Sprintf("%v", tag)
					}
					tagsJSON, _ := json.Marshal(stringTags)
					data["tags"] = string(tagsJSON)
				} else if tagsArray, ok := tags.([]string); ok && len(tagsArray) > 0 {
					tagsJSON, _ := json.Marshal(tagsArray)
					data["tags"] = string(tagsJSON)
				} else {
					data["tags"] = "[]"
				}
			}

			return nil
		}).
		Build() // Returns a fully configured CrudServiceContract

	bookServiceInstance := &BookService{
		CrudServiceContract: service,
		baseService:         service,
	}

	// Set the actual service instance for proper method resolution
	contracts.SetActualServiceHelper(service, bookServiceInstance, "BookService")

	return bookServiceInstance
}

// Override GetColumnMapping to include book-specific mappings
func (s *BookService) GetColumnMapping() map[string]string {
	mapping := s.baseService.GetColumnMapping()
	// Add book-specific mappings
	mapping["publishedAt"] = "published_at"
	mapping["createdAt"] = "created_at"
	mapping["updatedAt"] = "updated_at"
	mapping["title"] = "title"
	mapping["author"] = "author"
	mapping["price"] = "price"
	return mapping
}

// Override MapSortField to handle frontend field names
func (s *BookService) MapSortField(frontendField string) (string, bool) {
	// Check if we have a mapping for this field
	mapping := s.GetColumnMapping()

	if dbField, exists := mapping[frontendField]; exists {
		// Check if the mapped field is sortable
		sortableFields := s.baseService.GetSortableFields()
		for _, field := range sortableFields {
			if field == dbField {
				return dbField, true
			}
		}
	}

	// If no mapping exists, check if the field itself is sortable
	sortableFields := s.baseService.GetSortableFields()
	for _, field := range sortableFields {
		if field == frontendField {
			return frontendField, true
		}
	}

	// Field is not sortable
	return "", false
}

// ValidateSortField validates if a field can be sorted
func (s *BookService) ValidateSortField(field string) bool {
	// Check if the field is in our sortable fields list
	sortableFields := s.baseService.GetSortableFields()
	for _, sortableField := range sortableFields {
		if sortableField == field {
			return true
		}
	}
	return false
}

// ValidateSortDirection validates sort direction
func (s *BookService) ValidateSortDirection(direction string) bool {
	// Standard validation for ASC/DESC (case insensitive)
	upper := strings.ToUpper(direction)
	return upper == "ASC" || upper == "DESC"
}

// GetDefaultSort returns the default sort configuration
func (s *BookService) GetDefaultSort() (string, string) {
	// Return default sort: created_at DESC
	return "created_at", "DESC"
}

// Book-specific methods beyond basic CRUD

// GetByISBN retrieves a book by ISBN
func (s *BookService) GetByISBN(isbn string) (*models.Book, error) {
	var book models.Book
	err := facades.Orm().Query().
		Model(&models.Book{}).
		Where("isbn = ?", isbn).
		With("Creator").
		With("Updater").
		First(&book)

	if err != nil {
		return nil, err
	}

	return &book, nil
}

// GetByAuthor retrieves books by author with pagination
func (s *BookService) GetByAuthor(author string, req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	// Add author filter to the request
	if req.Filters == nil {
		req.Filters = make(map[string]interface{})
	}
	req.Filters["author"] = author

	return s.GetList(req)
}

// GetAvailable retrieves all available books
func (s *BookService) GetAvailable(req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	// Add status filter to the request
	if req.Filters == nil {
		req.Filters = make(map[string]interface{})
	}
	req.Filters["status"] = "AVAILABLE"

	return s.GetList(req)
}

// BorrowBook updates book status to borrowed
func (s *BookService) BorrowBook(id uint) error {
	// Get the book first to check if it's available
	bookInterface, err := s.GetByID(id)
	if err != nil {
		return err
	}

	book, ok := bookInterface.(*models.Book)
	if !ok {
		return errors.New("invalid book type")
	}

	if book.Status != "AVAILABLE" {
		return errors.New("book is not available for borrowing")
	}

	// Update the status
	updateData := map[string]interface{}{
		"status": "BORROWED",
	}

	_, err = s.Update(id, updateData)
	return err
}

// ReturnBook updates book status back to available
func (s *BookService) ReturnBook(id uint) error {
	// Get the book first to check if it's borrowed
	bookInterface, err := s.GetByID(id)
	if err != nil {
		return err
	}

	book, ok := bookInterface.(*models.Book)
	if !ok {
		return errors.New("invalid book type")
	}

	if book.Status != "BORROWED" {
		return errors.New("book is not currently borrowed")
	}

	// Update the status
	updateData := map[string]interface{}{
		"status": "AVAILABLE",
	}

	_, err = s.Update(id, updateData)
	return err
}

// GetBookStatistics returns statistics about books
func (s *BookService) GetBookStatistics() (map[string]interface{}, error) {
	var stats struct {
		TotalBooks       int64
		AvailableBooks   int64
		BorrowedBooks    int64
		MaintenanceBooks int64
		TotalValue       float64
		AveragePrice     float64
	}

	// Get total books
	facades.Orm().Query().Model(&models.Book{}).Count(&stats.TotalBooks)

	// Get available books
	facades.Orm().Query().Model(&models.Book{}).Where("status = ?", "AVAILABLE").Count(&stats.AvailableBooks)

	// Get borrowed books
	facades.Orm().Query().Model(&models.Book{}).Where("status = ?", "BORROWED").Count(&stats.BorrowedBooks)

	// Get maintenance books
	facades.Orm().Query().Model(&models.Book{}).Where("status = ?", "MAINTENANCE").Count(&stats.MaintenanceBooks)

	// Get total value
	facades.Orm().Query().Model(&models.Book{}).
		Select("SUM(price) as total").
		Pluck("total", &stats.TotalValue)

	// Calculate average price
	if stats.TotalBooks > 0 {
		stats.AveragePrice = stats.TotalValue / float64(stats.TotalBooks)
	}

	return map[string]interface{}{
		"totalBooks":       stats.TotalBooks,
		"availableBooks":   stats.AvailableBooks,
		"borrowedBooks":    stats.BorrowedBooks,
		"maintenanceBooks": stats.MaintenanceBooks,
		"totalValue":       fmt.Sprintf("$%.2f", stats.TotalValue),
		"averagePrice":     stats.AveragePrice,
	}, nil
}

// Custom filter for book status
func (s *BookService) ApplyCustomFilters(query orm.Query, filters map[string]interface{}) orm.Query {
	// Handle special book filters
	for key, value := range filters {
		switch key {
		case "minPrice":
			if price, ok := value.(float64); ok {
				query = query.Where("price >= ?", price)
			}
		case "maxPrice":
			if price, ok := value.(float64); ok {
				query = query.Where("price <= ?", price)
			}
		case "publishedYear":
			if year, ok := value.(string); ok {
				query = query.Where("YEAR(published_at) = ?", year)
			}
		}
	}

	return query
}
