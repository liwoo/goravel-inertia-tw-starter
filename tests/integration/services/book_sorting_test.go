package integration

import (
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"
	
	"players/app/contracts"
	"players/app/models"
	"players/app/services"
	"players/tests"
	"players/tests/helpers"
)

type BookSortingTestSuite struct {
	suite.Suite
	tests.TestCase
	
	// Test users
	testUser     *models.User
	testUserRole *models.Role
	
	// Test books with different statuses
	availableBook1   *models.Book
	availableBook2   *models.Book
	borrowedBook1    *models.Book
	borrowedBook2    *models.Book
	maintenanceBook  *models.Book
	reservedBook     *models.Book
}

func TestBookSortingTestSuite(t *testing.T) {
	suite.Run(t, new(BookSortingTestSuite))
}

func (s *BookSortingTestSuite) SetupTest() {
	s.RefreshDatabase()
	// Clean existing books to ensure test isolation
	facades.Orm().Query().Exec("DELETE FROM books")
	s.setupRolesAndPermissions()
	s.setupTestUser()
	s.createTestBooks()
}

func (s *BookSortingTestSuite) setupRolesAndPermissions() {
	// Create a role
	adminRole := &models.Role{Name: "Admin", Slug: "admin", IsActive: true, Level: 100}
	facades.Orm().Query().Create(adminRole)
	
	// Create the books_read permission with by_all scope
	perm := &models.Permission{
		Slug:     "books_read",
		Name:     "Read Books",
		Resource: "books",
		Action:   "read",
		IsActive: true,
	}
	facades.Orm().Query().Create(perm)
	
	// Assign permission to role with by_all scope
	facades.Orm().Query().Create(&models.RolePermission{
		RoleID:       adminRole.ID,
		PermissionID: perm.ID,
		Scope:        "by_all",
		IsActive:     true,
	})
	
	// Store role for user assignment
	s.testUserRole = adminRole
}

func (s *BookSortingTestSuite) setupTestUser() {
	password, _ := facades.Hash().Make("password")
	
	// Create test user
	s.testUser = &models.User{
		Name:     "Test User",
		Email:    "testuser@test.com",
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.testUser)
	
	// Assign role to user
	facades.Orm().Query().Create(&models.UserRole{
		UserID:     s.testUser.ID,
		RoleID:     s.testUserRole.ID,
		IsActive:   true,
		AssignedAt: time.Now(),
	})
	
	// Reload user with roles
	facades.Orm().Query().With("Roles.Permissions").Where("id = ?", s.testUser.ID).First(s.testUser)
}

func (s *BookSortingTestSuite) createTestBooks() {
	// Create books with different statuses
	// Using different prices and titles to test other sorts work correctly
	
	// Available books
	s.availableBook1 = &models.Book{
		Title:       "Available Book A",
		Author:      "Author A",
		ISBN:        "ISBN-A1",
		Status:      "AVAILABLE",
		Price:       20.00,
		Description: "Test book",
		BaseAuditableModel: models.BaseAuditableModel{
			CreatedBy: &s.testUser.ID,
		},
	}
	facades.Orm().Query().Create(s.availableBook1)
	// Sleep to ensure different timestamps
	time.Sleep(10 * time.Millisecond)
	
	s.availableBook2 = &models.Book{
		Title:       "Available Book B", 
		Author:      "Author B",
		ISBN:        "ISBN-A2",
		Status:      "AVAILABLE",
		Price:       15.00,
		Description: "Test book",
		BaseAuditableModel: models.BaseAuditableModel{
			CreatedBy: &s.testUser.ID,
		},
	}
	facades.Orm().Query().Create(s.availableBook2)
	time.Sleep(10 * time.Millisecond)
	
	// Borrowed books
	s.borrowedBook1 = &models.Book{
		Title:       "Borrowed Book A",
		Author:      "Author C",
		ISBN:        "ISBN-B1",
		Status:      "BORROWED",
		Price:       25.00,
		Description: "Test book",
		BaseAuditableModel: models.BaseAuditableModel{
			CreatedBy: &s.testUser.ID,
		},
	}
	facades.Orm().Query().Create(s.borrowedBook1)
	time.Sleep(10 * time.Millisecond)
	
	s.borrowedBook2 = &models.Book{
		Title:       "Borrowed Book B",
		Author:      "Author D",
		ISBN:        "ISBN-B2",
		Status:      "BORROWED", 
		Price:       30.00,
		Description: "Test book",
		BaseAuditableModel: models.BaseAuditableModel{
			CreatedBy: &s.testUser.ID,
		},
	}
	facades.Orm().Query().Create(s.borrowedBook2)
	time.Sleep(10 * time.Millisecond)
	
	// Maintenance book
	s.maintenanceBook = &models.Book{
		Title:       "Maintenance Book",
		Author:      "Author E",
		ISBN:        "ISBN-M1",
		Status:      "MAINTENANCE",
		Price:       10.00,
		Description: "Test book",
		BaseAuditableModel: models.BaseAuditableModel{
			CreatedBy: &s.testUser.ID,
		},
	}
	facades.Orm().Query().Create(s.maintenanceBook)
	time.Sleep(10 * time.Millisecond)
	
	// Reserved book
	s.reservedBook = &models.Book{
		Title:       "Reserved Book",
		Author:      "Author F",
		ISBN:        "ISBN-R1",
		Status:      "RESERVED",
		Price:       35.00,
		Description: "Test book",
		BaseAuditableModel: models.BaseAuditableModel{
			CreatedBy: &s.testUser.ID,
		},
	}
	facades.Orm().Query().Create(s.reservedBook)
}

// Test sorting by status field
func (s *BookSortingTestSuite) TestSortByStatus() {
	bookService := services.NewBookService()
	
	// Create context for authenticated user
	ctx := helpers.NewTestContext(s.testUser)
	
	// Test ascending sort by status
	s.T().Log("Testing ascending sort by status")
	req := contracts.ListRequest{
		Page:     1,
		PageSize: 100,
		Sort:     "status",
		Direction: "ASC",
		Context:  ctx,
	}
	
	result, err := bookService.GetList(req)
	s.NoError(err, "Should not get error when sorting by status")
	s.NotNil(result, "Result should not be nil")
	s.Equal(int64(6), result.Total, "Should have 6 books")
	
	// Check the order - should be AVAILABLE, AVAILABLE, BORROWED, BORROWED, MAINTENANCE, RESERVED
	books := result.Data
	s.Len(books, 6, "Should return all 6 books")
	
	// Convert to book models to check status
	if len(books) == 6 {
		// Handle both pointer and non-pointer types
		var book0, book1, book2, book3, book4, book5 *models.Book
		
		// Type assertion with fallback
		if b, ok := books[0].(*models.Book); ok {
			book0 = b
		} else if b, ok := books[0].(models.Book); ok {
			book0 = &b
		}
		if b, ok := books[1].(*models.Book); ok {
			book1 = b
		} else if b, ok := books[1].(models.Book); ok {
			book1 = &b
		}
		if b, ok := books[2].(*models.Book); ok {
			book2 = b
		} else if b, ok := books[2].(models.Book); ok {
			book2 = &b
		}
		if b, ok := books[3].(*models.Book); ok {
			book3 = b
		} else if b, ok := books[3].(models.Book); ok {
			book3 = &b
		}
		if b, ok := books[4].(*models.Book); ok {
			book4 = b
		} else if b, ok := books[4].(models.Book); ok {
			book4 = &b
		}
		if b, ok := books[5].(*models.Book); ok {
			book5 = b
		} else if b, ok := books[5].(models.Book); ok {
			book5 = &b
		}
		
		s.Equal("AVAILABLE", book0.Status, "First book should be AVAILABLE")
		s.Equal("AVAILABLE", book1.Status, "Second book should be AVAILABLE")
		s.Equal("BORROWED", book2.Status, "Third book should be BORROWED")
		s.Equal("BORROWED", book3.Status, "Fourth book should be BORROWED")
		s.Equal("MAINTENANCE", book4.Status, "Fifth book should be MAINTENANCE")
		s.Equal("RESERVED", book5.Status, "Sixth book should be RESERVED")
	}
	
	// Test descending sort by status
	s.T().Log("Testing descending sort by status")
	req.Direction = "DESC"
	
	result, err = bookService.GetList(req)
	s.NoError(err, "Should not get error when sorting by status DESC")
	s.NotNil(result, "Result should not be nil")
	
	// Check the order - should be RESERVED, MAINTENANCE, BORROWED, BORROWED, AVAILABLE, AVAILABLE
	books = result.Data
	s.Len(books, 6, "Should return all 6 books")
	
	if len(books) == 6 {
		// Handle both pointer and non-pointer types
		var book0, book1, book2, book3, book4, book5 *models.Book
		
		// Type assertion with fallback
		if b, ok := books[0].(*models.Book); ok {
			book0 = b
		} else if b, ok := books[0].(models.Book); ok {
			book0 = &b
		}
		if b, ok := books[1].(*models.Book); ok {
			book1 = b
		} else if b, ok := books[1].(models.Book); ok {
			book1 = &b
		}
		if b, ok := books[2].(*models.Book); ok {
			book2 = b
		} else if b, ok := books[2].(models.Book); ok {
			book2 = &b
		}
		if b, ok := books[3].(*models.Book); ok {
			book3 = b
		} else if b, ok := books[3].(models.Book); ok {
			book3 = &b
		}
		if b, ok := books[4].(*models.Book); ok {
			book4 = b
		} else if b, ok := books[4].(models.Book); ok {
			book4 = &b
		}
		if b, ok := books[5].(*models.Book); ok {
			book5 = b
		} else if b, ok := books[5].(models.Book); ok {
			book5 = &b
		}
		
		s.Equal("RESERVED", book0.Status, "First book should be RESERVED")
		s.Equal("MAINTENANCE", book1.Status, "Second book should be MAINTENANCE")
		s.Equal("BORROWED", book2.Status, "Third book should be BORROWED")
		s.Equal("BORROWED", book3.Status, "Fourth book should be BORROWED")
		s.Equal("AVAILABLE", book4.Status, "Fifth book should be AVAILABLE")
		s.Equal("AVAILABLE", book5.Status, "Sixth book should be AVAILABLE")
	}
}

// Test that other sorts still work correctly
func (s *BookSortingTestSuite) TestSortByPrice() {
	bookService := services.NewBookService()
	
	// Create context for authenticated user
	ctx := helpers.NewTestContext(s.testUser)
	
	// Test ascending sort by price
	req := contracts.ListRequest{
		Page:     1,
		PageSize: 100,
		Sort:     "price",
		Direction: "ASC",
		Context:  ctx,
	}
	
	result, err := bookService.GetList(req)
	s.NoError(err, "Should not get error when sorting by price")
	s.NotNil(result, "Result should not be nil")
	
	// Check the order - should be sorted by price: 10, 15, 20, 25, 30, 35
	books := result.Data
	s.Len(books, 6, "Should return all 6 books")
	
	if len(books) == 6 {
		// Helper function to get book price
		getPrice := func(item interface{}) float64 {
			if b, ok := item.(*models.Book); ok {
				return b.Price
			} else if b, ok := item.(models.Book); ok {
				return b.Price
			}
			return 0
		}
		
		s.Equal(10.00, getPrice(books[0]), "First book should have price 10")
		s.Equal(15.00, getPrice(books[1]), "Second book should have price 15")
		s.Equal(20.00, getPrice(books[2]), "Third book should have price 20")
		s.Equal(25.00, getPrice(books[3]), "Fourth book should have price 25")
		s.Equal(30.00, getPrice(books[4]), "Fifth book should have price 30")
		s.Equal(35.00, getPrice(books[5]), "Sixth book should have price 35")
	}
}

// Test sorting validation
func (s *BookSortingTestSuite) TestInvalidSortField() {
	bookService := services.NewBookService()
	
	// Test that invalid sort field validation works
	isValid := bookService.ValidateSortField("invalid_field")
	s.False(isValid, "Invalid field should not be sortable")
	
	// Test that status field validation
	isValid = bookService.ValidateSortField("status")
	s.True(isValid, "Status field should be sortable")
}

// Test MapSortField functionality
func (s *BookSortingTestSuite) TestMapSortField() {
	bookService := services.NewBookService()
	
	// Test mapping for status field
	dbField, isSortable := bookService.MapSortField("status")
	s.True(isSortable, "Status should be sortable")
	s.Equal("status", dbField, "Status should map to itself")
}

func (s *BookSortingTestSuite) TearDownTest() {
	// Cleanup handled by RefreshDatabase
}