package crud

import (
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"players/app/contracts"
	"players/app/models"
	"players/app/services"
	"players/tests"
)

type CrudSortingTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestCrudSortingTestSuite(t *testing.T) {
	suite.Run(t, new(CrudSortingTestSuite))
}

func (s *CrudSortingTestSuite) SetupTest() {
	s.RefreshDatabase()
}

// Test that BookService correctly sorts by date fields
func (s *CrudSortingTestSuite) TestBookServiceSortsByPublishedAt() {
	// Skip this test in SQLite environment as published_at is a string column
	if facades.Config().GetString("database.default") == "sqlite" {
		s.T().Skip("Skipping date sorting test in SQLite environment")
	}

	// Create test books with different published dates
	now := time.Now()
	books := []models.Book{
		{
			Title:       "Book C",
			Author:      "Author C",
			ISBN:        "ISBN-C",
			PublishedAt: &now,
			Status:      "AVAILABLE",
		},
		{
			Title:       "Book A",
			Author:      "Author A",
			ISBN:        "ISBN-A",
			PublishedAt: func() *time.Time { t := now.Add(-48 * time.Hour); return &t }(),
			Status:      "AVAILABLE",
		},
		{
			Title:       "Book B",
			Author:      "Author B",
			ISBN:        "ISBN-B",
			PublishedAt: func() *time.Time { t := now.Add(-24 * time.Hour); return &t }(),
			Status:      "AVAILABLE",
		},
	}

	// Create books
	for i := range books {
		err := facades.Orm().Query().Create(&books[i])
		s.NoError(err)
	}

	// Test sorting by publishedAt ascending
	bookService := services.NewBookService()
	req := contracts.ListRequest{
		Page:      1,
		PageSize:  10,
		Sort:      "publishedAt",
		Direction: "asc",
	}

	result, err := bookService.GetList(req)
	s.NoError(err)
	s.Equal(int64(3), result.Total)

	// Check order - oldest first
	s.Len(result.Data, 3)
	bookData := make([]models.Book, 3)
	for i, item := range result.Data {
		bookData[i] = item.(models.Book)
	}
	s.Equal("Book A", bookData[0].Title)
	s.Equal("Book B", bookData[1].Title)
	s.Equal("Book C", bookData[2].Title)

	// Test sorting by publishedAt descending
	req.Direction = "desc"
	result, err = bookService.GetList(req)
	s.NoError(err)

	// Check order - newest first
	s.Len(result.Data, 3)
	for i, item := range result.Data {
		bookData[i] = item.(models.Book)
	}
	s.Equal("Book C", bookData[0].Title)
	s.Equal("Book B", bookData[1].Title)
	s.Equal("Book A", bookData[2].Title)
}

// Test that RoleService correctly sorts by name
func (s *CrudSortingTestSuite) TestRoleServiceSortsByName() {
	// Create test roles
	roles := []models.Role{
		{Name: "Zebra Role", Slug: "zebra", IsActive: true},
		{Name: "Alpha Role", Slug: "alpha", IsActive: true},
		{Name: "Beta Role", Slug: "beta", IsActive: true},
	}

	for i := range roles {
		err := facades.Orm().Query().Create(&roles[i])
		s.NoError(err)
	}

	// Test sorting by name ascending
	roleService := services.NewRoleService()
	req := contracts.ListRequest{
		Page:      1,
		PageSize:  10,
		Sort:      "name",
		Direction: "asc",
	}

	result, err := roleService.GetList(req)
	s.NoError(err)
	s.GreaterOrEqual(result.Total, int64(3))

	// Check order
	foundAlpha := false
	foundBeta := false
	foundZebra := false

	for _, item := range result.Data {
		role := item.(models.Role)
		if role.Name == "Alpha Role" {
			foundAlpha = true
		} else if role.Name == "Beta Role" {
			foundBeta = true
			s.True(foundAlpha, "Alpha should come before Beta")
		} else if role.Name == "Zebra Role" {
			foundZebra = true
			s.True(foundAlpha && foundBeta, "Zebra should come after Alpha and Beta")
		}
	}

	s.True(foundAlpha && foundBeta && foundZebra, "All test roles should be found")
}

// Test that field mapping works correctly for frontend field names
func (s *CrudSortingTestSuite) TestBookServiceMapsFieldNames() {
	bookService := services.NewBookService()

	// Test mapping of camelCase to snake_case
	testCases := []struct {
		frontendField string
		expectedField string
		shouldMap     bool
	}{
		{"publishedAt", "published_at", true},
		{"createdAt", "created_at", true},
		{"updatedAt", "updated_at", true},
		{"title", "title", true},
		{"invalidField", "", false},
	}

	for _, tc := range testCases {
		mappedField, ok := bookService.MapSortField(tc.frontendField)
		s.Equal(tc.shouldMap, ok, "Field %s mapping result", tc.frontendField)
		if tc.shouldMap {
			s.Equal(tc.expectedField, mappedField, "Field %s should map to %s", tc.frontendField, tc.expectedField)
		}
	}
}

// Test filter functionality with boolean values
func (s *CrudSortingTestSuite) TestRoleServiceFiltersBoolean() {
	// Create active and inactive roles
	activeRole := models.Role{Name: "Active Role", Slug: "active_role", IsActive: true}
	inactiveRole := models.Role{Name: "Inactive Role", Slug: "inactive_role"}

	err := facades.Orm().Query().Create(&activeRole)
	s.NoError(err)

	// Create inactive role and explicitly set IsActive to false
	err = facades.Orm().Query().Create(&inactiveRole)
	s.NoError(err)
	// Update to set IsActive to false (to override the default)
	_, err = facades.Orm().Query().Model(&inactiveRole).Update("is_active", false)
	s.NoError(err)

	roleService := services.NewRoleService()

	// Test filtering by is_active = true
	req := contracts.ListRequest{
		Page:     1,
		PageSize: 100,
		Filters: map[string]interface{}{
			"is_active": true,
		},
	}

	result, err := roleService.GetList(req)
	s.NoError(err)

	// Should have at least the active role
	s.GreaterOrEqual(result.Total, int64(1))

	// Verify all results are active
	for _, item := range result.Data {
		role := item.(models.Role)
		s.True(role.IsActive, "All filtered roles should be active")
	}

	// Test filtering by is_active = false
	req.Filters["is_active"] = false
	result, err = roleService.GetList(req)
	s.NoError(err)

	// Should have at least the inactive role
	s.GreaterOrEqual(result.Total, int64(1))

	// Verify all results are inactive
	for _, item := range result.Data {
		role := item.(models.Role)
		s.False(role.IsActive, "All filtered roles should be inactive")
	}
}

func (s *CrudSortingTestSuite) TearDownTest() {
	// Cleanup handled by RefreshDatabase
}
