package crud

import (
	"testing"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"books-database/app/models"
	"books-database/app/services"
	"books-database/tests"
)

type SimpleCrudTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestSimpleCrudTestSuite(t *testing.T) {
	suite.Run(t, new(SimpleCrudTestSuite))
}

// SetupTest will run before each test in the suite.
func (s *SimpleCrudTestSuite) SetupTest() {
	// Refresh database - this will run migrations
	s.RefreshDatabase()
}

// Test that BookService sorts correctly
func (s *SimpleCrudTestSuite) TestBookServiceSorting() {
	// Create test books
	books := []models.Book{
		{Title: "Book C", Author: "Author C", ISBN: "ISBN-C", Status: "AVAILABLE"},
		{Title: "Book A", Author: "Author A", ISBN: "ISBN-A", Status: "AVAILABLE"},
		{Title: "Book B", Author: "Author B", ISBN: "ISBN-B", Status: "AVAILABLE"},
	}

	for i := range books {
		err := facades.Orm().Query().Create(&books[i])
		s.NoError(err)
	}

	// Test service
	bookService := services.NewBookService()

	// Test MapSortField
	field, ok := bookService.MapSortField("publishedAt")
	s.True(ok)
	s.Equal("published_at", field)

	field, ok = bookService.MapSortField("createdAt")
	s.True(ok)
	s.Equal("created_at", field)

	// Test that sortable fields are exposed
	sortableFields := bookService.GetSortableFields()
	s.Contains(sortableFields, "title")
	s.Contains(sortableFields, "author")
	s.Contains(sortableFields, "created_at")
}

// Test RoleService filters
func (s *SimpleCrudTestSuite) TestRoleServiceFilters() {
	// Create test roles
	activeRole := models.Role{Name: "Active Role", Slug: "active_test", IsActive: true}
	inactiveRole := models.Role{Name: "Inactive Role", Slug: "inactive_test", IsActive: false}

	err := facades.Orm().Query().Create(&activeRole)
	s.NoError(err)
	err = facades.Orm().Query().Create(&inactiveRole)
	s.NoError(err)

	roleService := services.NewRoleService()

	// Test that sortable fields are exposed
	sortableFields := roleService.GetSortableFields()
	s.Contains(sortableFields, "name")
	s.Contains(sortableFields, "created_at")

	// Test field mapping
	field, ok := roleService.MapSortField("name")
	s.True(ok)
	s.Equal("name", field)
}
