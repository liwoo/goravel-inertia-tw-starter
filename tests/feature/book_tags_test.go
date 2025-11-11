package feature

import (
	"testing"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"smedi-sme-db/app/models"
	"smedi-sme-db/tests"
)

type BookTagsTestSuite struct {
	suite.Suite
	tests.TestCase

	testUser *models.User
}

func TestBookTagsTestSuite(t *testing.T) {
	suite.Run(t, new(BookTagsTestSuite))
}

func (s *BookTagsTestSuite) SetupSuite() {
	// No need to call parent SetupSuite - it doesn't exist
}

func (s *BookTagsTestSuite) SetupTest() {
	// Refresh database to ensure clean state with migrations
	s.RefreshDatabase()
}

// Test_BookCreation_WithTags tests creating a book with tags
func (s *BookTagsTestSuite) Test_BookCreation_WithTags() {
	// Create a book with tags directly in the database
	book := &models.Book{
		Title:       "Test Book With Tags",
		Author:      "Test Author",
		ISBN:        "978-0123456789",
		Price:       19.99,
		Status:      "AVAILABLE",
		Description: "A test book with tags",
		Tags:        []string{"fiction", "test", "golang"},
	}

	// This should succeed now that tags column exists
	err := facades.Orm().Query().Create(&book)
	s.Require().NoError(err, "Should create book with tags")

	// Verify the book was created with correct ID
	s.NotZero(book.ID, "Book should have an ID after creation")

	// Retrieve the book to verify tags were stored correctly
	var createdBook models.Book
	err = facades.Orm().Query().Where("id", book.ID).First(&createdBook)
	s.Require().NoError(err, "Should find the created book")

	// Verify tags were stored correctly
	expectedTags := []string{"fiction", "test", "golang"}
	s.Equal(expectedTags, createdBook.Tags, "Tags should be stored correctly")
}
