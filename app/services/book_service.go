package services

import (
	"fmt"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/facades"
	"players/app/contracts"
	"players/app/models"
)

// BookService - Simplified version using generic CRUD service
// From ~600 lines to ~150 lines!
type BookService struct {
	*contracts.GenericCrudService[models.Book]
}

// NewBookService creates a new simplified book service
func NewBookService() *BookService {
	// Create the generic service
	genericService := contracts.NewGenericCrudService[models.Book]("book", "id")
	
	// Configure the service
	genericService.
		SetSearchFields("title", "author", "isbn", "description").
		SetSortFields("id", "title", "author", "price", "published_at", "created_at", "updated_at").
		SetFilterFields("status", "author", "isbn", "is_available").
		SetValidationRules(map[string]interface{}{
			"title":        "required|string|max:255",
			"author":       "required|string|max:255",
			"isbn":         "required|string|max:20",
			"description":  "string|max:1000",
			"price":        "numeric|min:0",
			"status":       "string|in:available,borrowed,reserved,lost",
			"published_at": "date",
		}).
		SetBeforeCreate(func(data map[string]interface{}) error {
			// Set default status if not provided
			if _, exists := data["status"]; !exists {
				data["status"] = "available"
			}
			if _, exists := data["is_available"]; !exists {
				data["is_available"] = true
			}
			
			// Check ISBN uniqueness
			var count int64
			err := facades.Orm().Query().Model(&models.Book{}).
				Where("isbn = ?", data["isbn"]).
				Count(&count)
			if err != nil {
				return fmt.Errorf("failed to check ISBN uniqueness: %w", err)
			}
			if count > 0 {
				return fmt.Errorf("ISBN already exists")
			}
			
			return nil
		}).
		SetBeforeUpdate(func(id uint, data map[string]interface{}) error {
			// Check ISBN uniqueness if being changed
			if isbn, ok := data["isbn"].(string); ok {
				var count int64
				err := facades.Orm().Query().Model(&models.Book{}).
					Where("isbn = ? AND id != ?", isbn, id).
					Count(&count)
				if err != nil {
					return fmt.Errorf("failed to check ISBN uniqueness: %w", err)
				}
				if count > 0 {
					return fmt.Errorf("ISBN already exists")
				}
			}
			return nil
		}).
		SetCustomSearch(func(query orm.Query, search string) orm.Query {
			searchValue := "%" + search + "%"
			return query.Where("title LIKE ? OR author LIKE ? OR isbn LIKE ? OR description LIKE ?",
				searchValue, searchValue, searchValue, searchValue)
		}).
		SetCustomFilters(func(query orm.Query, filters map[string]interface{}) orm.Query {
			for field, value := range filters {
				switch field {
				case "status":
					query = query.Where("status = ?", value)
				case "author":
					query = query.Where("author = ?", value)
				case "minPrice":
					if price, ok := value.(float64); ok {
						query = query.Where("price >= ?", price)
					}
				case "maxPrice":
					if price, ok := value.(float64); ok {
						query = query.Where("price <= ?", price)
					}
				case "is_available":
					query = query.Where("is_available = ?", value)
				}
			}
			return query
		})
	
	service := &BookService{
		GenericCrudService: genericService,
	}
	
	// Register service
	contracts.MustRegisterCrudService("books", service)
	
	return service
}

// Custom methods beyond basic CRUD

// GetByISBN retrieves a book by ISBN
func (s *BookService) GetByISBN(isbn string) (*models.Book, error) {
	var book models.Book
	if err := facades.Orm().Query().Where("isbn = ?", isbn).First(&book); err != nil {
		return nil, fmt.Errorf("book not found with ISBN %s: %w", isbn, err)
	}
	return &book, nil
}

// GetByAuthor retrieves books by author with pagination
func (s *BookService) GetByAuthor(author string, req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	// Add author filter
	filters := map[string]interface{}{
		"author": author,
	}
	return s.GetListAdvanced(req, filters)
}

// GetAvailable retrieves available books with pagination
func (s *BookService) GetAvailable(req contracts.ListRequest) (*contracts.PaginatedResult, error) {
	// Add availability filter
	filters := map[string]interface{}{
		"is_available": true,
		"status":       "available",
	}
	return s.GetListAdvanced(req, filters)
}

// BorrowBook marks a book as borrowed
func (s *BookService) BorrowBook(id uint) error {
	book, err := s.GetByID(id)
	if err != nil {
		return err
	}
	
	bookModel := book.(*models.Book)
	if bookModel.Status != "available" {
		return fmt.Errorf("book is not available for borrowing")
	}
	
	// Update book status
	updateData := map[string]interface{}{
		"status":       "borrowed",
		"is_available": false,
	}
	
	_, err = s.Update(id, updateData)
	return err
}

// ReturnBook marks a book as returned
func (s *BookService) ReturnBook(id uint) error {
	book, err := s.GetByID(id)
	if err != nil {
		return err
	}
	
	bookModel := book.(*models.Book)
	if bookModel.Status != "borrowed" {
		return fmt.Errorf("book is not currently borrowed")
	}
	
	// Update book status
	updateData := map[string]interface{}{
		"status":       "available",
		"is_available": true,
	}
	
	_, err = s.Update(id, updateData)
	return err
}

// Contract method implementations
func (s *BookService) GetColumnMapping() map[string]string {
	return map[string]string{
		"id":           "id",
		"title":        "title",
		"author":       "author",
		"isbn":         "isbn",
		"price":        "price",
		"publishedAt":  "published_at",
		"createdAt":    "created_at",
		"updatedAt":    "updated_at",
		"isAvailable":  "is_available",
		"status":       "status",
	}
}