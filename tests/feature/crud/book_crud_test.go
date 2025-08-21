package crud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"players/app/models"
	"players/tests"
	"players/tests/helpers"
)

type BookCRUDTestSuite struct {
	suite.Suite
	tests.TestCase

	server     *httptest.Server
	client     *http.Client
	authCookie *http.Cookie
	testUser   *models.User
}

func TestBookCRUDTestSuite(t *testing.T) {
	suite.Run(t, &BookCRUDTestSuite{})
}

func (s *BookCRUDTestSuite) SetupSuite() {
	// Start test server
	s.server = httptest.NewServer(facades.Route())

	// Create HTTP client with cookie jar
	jar, _ := cookiejar.New(nil)
	s.client = &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func (s *BookCRUDTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}

func (s *BookCRUDTestSuite) SetupTest() {
	// Clean database
	s.RefreshDatabase()

	// Create test user with permissions
	s.setupTestUser()
}

func (s *BookCRUDTestSuite) TearDownTest() {
	// Clean up test data
	if orm := facades.Orm(); orm != nil {
		orm.Query().Exec("DELETE FROM books WHERE isbn LIKE '978%'")
		orm.Query().Exec("DELETE FROM users WHERE email = 'booktest@example.com'")
		orm.Query().Exec("DELETE FROM user_roles")
		orm.Query().Exec("DELETE FROM role_permissions")
		orm.Query().Exec("DELETE FROM roles WHERE slug = 'book_admin'")
		orm.Query().Exec("DELETE FROM permissions WHERE slug LIKE 'books_%'")
	}
}

func (s *BookCRUDTestSuite) setupTestUser() {
	// Create admin role
	adminRole := &models.Role{
		Name:  "Book Admin",
		Slug:  "book_admin",
		Level: 100,
	}
	s.Nil(facades.Orm().Query().Create(adminRole))

	// Create all book permissions
	permissions := []string{"create", "read", "update", "delete"}
	for _, perm := range permissions {
		permission := &models.Permission{
			Name:  fmt.Sprintf("books_%s", perm),
			Slug:  fmt.Sprintf("books_%s", perm),
			Scope: "by_all",
		}
		s.Nil(facades.Orm().Query().Create(permission))
		s.Nil(helpers.AssignPermissionToRole(adminRole, permission, "by_all"))
	}

	// Create test user
	user, err := helpers.SetupJWTUser("booktest@example.com", "password", adminRole)
	s.Nil(err)
	s.testUser = user

	// Login
	s.authCookie = s.loginUser("booktest@example.com", "password")
	s.NotNil(s.authCookie)
}

func (s *BookCRUDTestSuite) loginUser(email, password string) *http.Cookie {
	loginData := map[string]string{
		"email":    email,
		"password": password,
	}

	jsonData, _ := json.Marshal(loginData)
	resp, err := s.client.Post(s.server.URL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
	s.Nil(err)
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "token" {
			return cookie
		}
	}

	return nil
}

func (s *BookCRUDTestSuite) makeRequest(method, path string, body interface{}) (*http.Response, map[string]interface{}) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, s.server.URL+path, bodyReader)
	s.Nil(err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	if s.authCookie != nil {
		req.AddCookie(s.authCookie)
	}

	resp, err := s.client.Do(req)
	s.Nil(err)

	respBody, err := io.ReadAll(resp.Body)
	s.Nil(err)
	resp.Body.Close()

	var result map[string]interface{}
	if len(respBody) > 0 {
		err = json.Unmarshal(respBody, &result)
		s.Nil(err)
	}

	return resp, result
}

// Test basic CRUD operations
func (s *BookCRUDTestSuite) TestCreateBook() {
	bookData := map[string]interface{}{
		"title":       "Test Book",
		"author":      "Test Author",
		"isbn":        "978-3-16-148410-0",
		"description": "A test book description",
		"price":       29.99,
		"status":      "AVAILABLE",
		"publishedAt": "2024-01-01",
		"tags":        []string{"test", "fiction"},
	}

	resp, result := s.makeRequest("POST", "/api/books", bookData)

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Create book failed: %d - %+v\n", resp.StatusCode, result)
	}
	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	// Check if data exists before asserting
	if result["data"] == nil {
		s.FailNow("Response missing 'data' field", "Response: %+v", result)
	}
	data := result["data"].(map[string]interface{})
	s.Equal("Test Book", data["title"])
	s.Equal("Test Author", data["author"])
	s.Equal("9783161484100", data["isbn"])
	s.Equal("A test book description", data["description"])
	s.Equal(29.99, data["price"])
	s.Equal("AVAILABLE", data["status"])
	s.NotNil(data["id"])
	s.Equal(float64(s.testUser.ID), data["created_by"])
}

func (s *BookCRUDTestSuite) TestCreateBookValidation() {
	// Missing required fields
	bookData := map[string]interface{}{
		"description": "A test book description",
	}

	resp, result := s.makeRequest("POST", "/api/books", bookData)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
	s.Contains(strings.ToLower(result["message"].(string)), "validation")
}

func (s *BookCRUDTestSuite) TestGetBook() {
	// Create a book first
	book := &models.Book{
		Title:       "Get Test Book",
		Author:      "Get Test Author",
		ISBN:        "9780596520687",
		Description: "Test description",
		Price:       19.99,
		Status:      "AVAILABLE",
	}
	book.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(book))

	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/books/%d", book.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal("Get Test Book", data["title"])
	s.Equal("Get Test Author", data["author"])
	s.Equal("9780596520687", data["isbn"])
}

func (s *BookCRUDTestSuite) TestUpdateBook() {
	// Create a book first
	book := &models.Book{
		Title:       "Update Test Book",
		Author:      "Update Test Author",
		ISBN:        "9781566199094",
		Description: "Test description",
		Price:       19.99,
		Status:      "AVAILABLE",
	}
	book.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(book))

	updateData := map[string]interface{}{
		"title":       "Updated Book Title",
		"price":       24.99,
		"status":      "BORROWED",
		"description": "Updated description",
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/books/%d", book.ID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal("Updated Book Title", data["title"])
	s.Equal(24.99, data["price"])
	s.Equal("BORROWED", data["status"])
	s.Equal("Updated description", data["description"])
	s.Equal("Update Test Author", data["author"]) // Unchanged
}

func (s *BookCRUDTestSuite) TestDeleteBook() {
	// Create a book first
	book := &models.Book{
		Title:       "Delete Test Book",
		Author:      "Delete Test Author",
		ISBN:        "9780131103627",
		Description: "Test description",
		Price:       19.99,
		Status:      "AVAILABLE",
	}
	book.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(book))

	resp, _ := s.makeRequest("DELETE", fmt.Sprintf("/api/books/%d", book.ID), nil)

	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify soft delete
	var deletedBook models.Book
	err := facades.Orm().Query().WithTrashed().Where("id", book.ID).First(&deletedBook)
	s.Nil(err)
	s.NotNil(deletedBook.DeletedAt)
}

// Test Pagination
func (s *BookCRUDTestSuite) TestPagination() {
	// Create 25 books
	for i := 1; i <= 25; i++ {
		book := &models.Book{
			Title:       fmt.Sprintf("Book %02d", i),
			Author:      "Test Author",
			ISBN:        fmt.Sprintf("978000%07d", 1000000+i),
			Description: "Test description",
			Price:       float64(10 + i),
			Status:      "AVAILABLE",
		}
		book.CreatedBy = &s.testUser.ID
		s.Nil(facades.Orm().Query().Create(book))
	}

	// Test first page (default page size 20)
	resp, result := s.makeRequest("GET", "/api/books?page=1", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	// For paginated responses, data is nested
	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	pagination := data["pagination"].(map[string]interface{})

	s.Equal(20, len(items))
	s.Equal(float64(1), pagination["current_page"])
	s.Equal(float64(25), pagination["total"])
	s.Equal(float64(2), pagination["last_page"])

	// Test second page
	resp, result = s.makeRequest("GET", "/api/books?page=2", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	pagination = data["pagination"].(map[string]interface{})

	s.Equal(5, len(items))
	s.Equal(float64(2), pagination["current_page"])

	// Test custom page size
	resp, result = s.makeRequest("GET", "/api/books?page=1&pageSize=10", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	pagination = data["pagination"].(map[string]interface{})

	s.Equal(10, len(items))
	s.Equal(float64(3), pagination["last_page"]) // 25 items / 10 per page = 3 pages
}

// Test Sorting
func (s *BookCRUDTestSuite) TestSorting() {
	// Create books with different attributes
	books := []struct {
		title  string
		author string
		price  float64
		date   string
	}{
		{"Alpha Book", "Charlie Author", 30.00, "2024-01-01"},
		{"Beta Book", "Alice Author", 20.00, "2024-02-01"},
		{"Gamma Book", "Bob Author", 10.00, "2024-03-01"},
	}

	// First, clean up any existing books to ensure clean test
	facades.Orm().Query().Exec("DELETE FROM books")

	for i, b := range books {
		pubDate, _ := time.Parse("2006-01-02", b.date)
		book := &models.Book{
			Title:       b.title,
			Author:      b.author,
			ISBN:        fmt.Sprintf("978100%07d", 2000000+i),
			Price:       b.price,
			Status:      "AVAILABLE",
			PublishedAt: &pubDate,
		}
		book.CreatedBy = &s.testUser.ID
		s.Nil(facades.Orm().Query().Create(book))
		time.Sleep(10 * time.Millisecond) // Ensure different created_at times
	}

	// Test sort by title ascending
	resp, result := s.makeRequest("GET", "/api/books?sort=title&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	// For paginated responses, data is nested
	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})

	s.Equal("Alpha Book", items[0].(map[string]interface{})["title"])
	s.Equal("Beta Book", items[1].(map[string]interface{})["title"])
	s.Equal("Gamma Book", items[2].(map[string]interface{})["title"])

	// Test sort by price descending
	resp, result = s.makeRequest("GET", "/api/books?sort=price&direction=desc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})

	s.Equal(float64(30.00), items[0].(map[string]interface{})["price"].(float64))
	s.Equal(float64(20.00), items[1].(map[string]interface{})["price"].(float64))
	s.Equal(float64(10.00), items[2].(map[string]interface{})["price"].(float64))

	// Test sort by author
	resp, result = s.makeRequest("GET", "/api/books?sort=author&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})

	s.Equal("Alice Author", items[0].(map[string]interface{})["author"])
	s.Equal("Bob Author", items[1].(map[string]interface{})["author"])
	s.Equal("Charlie Author", items[2].(map[string]interface{})["author"])

	// Test default sort (created_at DESC)
	resp, result = s.makeRequest("GET", "/api/books", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})

	// Most recently created should be first
	s.Equal("Gamma Book", items[0].(map[string]interface{})["title"])
}

// Test Search
func (s *BookCRUDTestSuite) TestSearch() {
	// Create books with searchable content - more books to test pagination
	books := []struct {
		title       string
		author      string
		isbn        string
		description string
	}{
		{"Go Programming", "John Doe", "9780134190440", "Learn Go programming language"},
		{"Python Basics", "Jane Smith", "9781491946008", "Introduction to Python"},
		{"Go Web Development", "John Doe", "9781787125643", "Building web apps with Go"},
		{"JavaScript Guide", "Bob Johnson", "9781593279509", "Modern JavaScript programming"},
		{"Advanced Go", "Alice Brown", "9780977123456", "Advanced Go techniques"},
		{"Go Patterns", "Bob Wilson", "9780987654321", "Design patterns in Go"},
		{"Go Testing", "Carol Green", "9780123456789", "Testing strategies for Go"},
		{"Go Microservices", "Dave Black", "9780987123456", "Microservices with Go"},
	}

	for _, b := range books {
		book := &models.Book{
			Title:       b.title,
			Author:      b.author,
			ISBN:        b.isbn,
			Description: b.description,
			Price:       29.99,
			Status:      "AVAILABLE",
		}
		book.CreatedBy = &s.testUser.ID
		s.Nil(facades.Orm().Query().Create(book))
	}

	// Search by title
	resp, result := s.makeRequest("GET", "/api/books/search?q=Go", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	// For paginated responses, data is nested
	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})

	s.Equal(6, len(items)) // Should find 6 Go books now with the expanded test data

	// Search by author
	resp, result = s.makeRequest("GET", "/api/books/search?q=John", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})

	s.GreaterOrEqual(len(items), 3) // John Doe (2 books) and Bob Johnson (1 book)

	// Search by description
	resp, result = s.makeRequest("GET", "/api/books/search?q=programming", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})

	s.GreaterOrEqual(len(items), 2) // At least 2 books have "programming" in description

	// Search with no results
	resp, result = s.makeRequest("GET", "/api/books/search?q=NoSuchBook", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})

	s.Equal(0, len(items))

	// Search with pagination - verify pagination works correctly
	resp, result = s.makeRequest("GET", "/api/books/search?q=Go&page=1&pageSize=5", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	pagination := data["pagination"].(map[string]interface{})

	// Verify pagination structure and that we get some results
	s.True(len(items) > 0 && len(items) <= 5) // Should get some items, max 5
	s.True(pagination["total"].(float64) > 0) // Should have some total count
	s.Equal(float64(1), pagination["current_page"])
	s.True(pagination["last_page"].(float64) >= 1) // Should have at least 1 page

	// Test that pageSize parameter is being respected by checking second page
	if pagination["total"].(float64) > 5 {
		resp, result = s.makeRequest("GET", "/api/books/search?q=Go&page=2&pageSize=5", nil)
		s.Equal(http.StatusOK, resp.StatusCode)

		data = result["data"].(map[string]interface{})
		items = data["data"].([]interface{})
		pagination = data["pagination"].(map[string]interface{})

		// Verify we get remaining items and pagination is consistent
		s.True(len(items) >= 0 && len(items) <= 5)
		s.Equal(float64(2), pagination["current_page"])
	}
}

// Test Filtering
func (s *BookCRUDTestSuite) TestFiltering() {
	// Create books with different statuses and authors
	books := []struct {
		title  string
		author string
		status string
	}{
		{"Book 1", "Author A", "AVAILABLE"},
		{"Book 2", "Author B", "BORROWED"},
		{"Book 3", "Author A", "MAINTENANCE"},
		{"Book 4", "Author B", "AVAILABLE"},
		{"Book 5", "Author C", "RESERVED"},
	}

	for i, b := range books {
		book := &models.Book{
			Title:  b.title,
			Author: b.author,
			ISBN:   fmt.Sprintf("978200%07d", 3000000+i),
			Status: b.status,
			Price:  19.99,
		}
		book.CreatedBy = &s.testUser.ID
		s.Nil(facades.Orm().Query().Create(book))
	}

	// Filter by status
	resp, result := s.makeRequest("GET", "/api/books?status=AVAILABLE", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	// For paginated responses, data is nested
	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})

	s.Equal(2, len(items))
	for _, item := range items {
		s.Equal("AVAILABLE", item.(map[string]interface{})["status"])
	}

	// Filter by author
	resp, result = s.makeRequest("GET", "/api/books?author=Author%20A", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})

	s.Equal(2, len(items))
	for _, item := range items {
		s.Equal("Author A", item.(map[string]interface{})["author"])
	}

	// Multiple filters
	resp, result = s.makeRequest("GET", "/api/books?status=AVAILABLE&author=Author%20B", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})

	s.Equal(1, len(items))
	s.Equal("Book 4", items[0].(map[string]interface{})["title"])
}

// Test combined features
func (s *BookCRUDTestSuite) TestCombinedFeatures() {
	// Create a variety of books
	for i := 1; i <= 30; i++ {
		status := "AVAILABLE"
		if i%3 == 0 {
			status = "BORROWED"
		} else if i%5 == 0 {
			status = "MAINTENANCE"
		}

		author := fmt.Sprintf("Author %c", 'A'+(i%3))

		book := &models.Book{
			Title:       fmt.Sprintf("Book %02d", i),
			Author:      author,
			ISBN:        fmt.Sprintf("978300%07d", 4000000+i),
			Description: fmt.Sprintf("Description for book %d with some searchable content", i),
			Price:       float64(10 + (i % 20)),
			Status:      status,
		}
		book.CreatedBy = &s.testUser.ID
		s.Nil(facades.Orm().Query().Create(book))
	}

	// Search with filter, sort and pagination
	resp, result := s.makeRequest("GET", "/api/books/search?q=searchable&status=AVAILABLE&sort=price&direction=desc&page=1&pageSize=5", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	// For paginated responses, data is nested
	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	pagination := data["pagination"].(map[string]interface{})

	s.Equal(5, len(items))
	s.True(pagination["total"].(float64) > 0)

	// Verify all items match criteria
	prevPrice := 100.0
	for _, item := range items {
		book := item.(map[string]interface{})
		s.Equal("AVAILABLE", book["status"])
		s.Contains(book["description"], "searchable")
		// Verify descending price order
		price := book["price"].(float64)
		s.LessOrEqual(price, prevPrice)
		prevPrice = price
	}
}

// Test soft delete behavior
func (s *BookCRUDTestSuite) TestSoftDelete() {
	fmt.Printf("DEBUG: Starting TestSoftDelete\n")

	// Create a book
	book := &models.Book{
		Title:  "Soft Delete Test",
		Author: "Test Author",
		ISBN:   "9784000000000",
		Status: "AVAILABLE",
		Price:  19.99,
	}
	book.CreatedBy = &s.testUser.ID
	fmt.Printf("DEBUG: About to create book\n")
	s.Nil(facades.Orm().Query().Create(book))
	bookID := book.ID
	fmt.Printf("DEBUG: Created book with ID: %d\n", bookID)

	// Check list before delete
	resp, result := s.makeRequest("GET", "/api/books", nil)
	s.Equal(http.StatusOK, resp.StatusCode)
	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	fmt.Printf("DEBUG: Books before delete: %d items\n", len(items))
	for _, item := range items {
		fmt.Printf("DEBUG: Book before delete: ID=%v\n", item.(map[string]interface{})["id"])
	}

	// Delete the book
	fmt.Printf("DEBUG: About to delete book ID: %d\n", bookID)
	deleteResp, _ := s.makeRequest("DELETE", fmt.Sprintf("/api/books/%d", bookID), nil)
	fmt.Printf("DEBUG: DELETE response status: %d\n", deleteResp.StatusCode)
	s.Equal(http.StatusNoContent, deleteResp.StatusCode)

	// Try to get the deleted book - should fail
	resp, result = s.makeRequest("GET", fmt.Sprintf("/api/books/%d", bookID), nil)
	fmt.Printf("DEBUG: GET after delete - status: %d\n", resp.StatusCode)
	if result != nil {
		fmt.Printf("DEBUG: GET after delete - response: %+v\n", result)
	}
	s.Equal(http.StatusNotFound, resp.StatusCode)

	// Verify it's not in the list
	resp, result = s.makeRequest("GET", "/api/books", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	// For paginated responses, data is nested
	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	fmt.Printf("DEBUG: Books after delete: %d items\n", len(items))
	for _, item := range items {
		fmt.Printf("DEBUG: Book after delete: ID=%v\n", item.(map[string]interface{})["id"])
		s.NotEqual(float64(bookID), item.(map[string]interface{})["id"])
	}

	// Verify it still exists in database with deleted_at
	var deletedBook models.Book
	err := facades.Orm().Query().WithTrashed().Where("id", bookID).First(&deletedBook)
	s.Nil(err)
	s.NotNil(deletedBook.DeletedAt)
	s.Equal("Soft Delete Test", deletedBook.Title)
}

// Test validation rules
func (s *BookCRUDTestSuite) TestValidationRules() {
	testCases := []struct {
		name        string
		data        map[string]interface{}
		expectError bool
		errorField  string
	}{
		{
			name: "Valid book",
			data: map[string]interface{}{
				"title":  "Valid Book",
				"author": "Valid Author",
				"isbn":   "9785000000001",
				"status": "AVAILABLE",
				"price":  29.99,
			},
			expectError: false,
		},
		{
			name: "Missing title",
			data: map[string]interface{}{
				"author": "Valid Author",
				"isbn":   "9785000000002",
				"status": "AVAILABLE",
				"price":  19.99,
			},
			expectError: true,
			errorField:  "title",
		},
		{
			name: "Invalid status",
			data: map[string]interface{}{
				"title":  "Valid Book",
				"author": "Valid Author",
				"isbn":   "9785000000003",
				"status": "INVALID_STATUS",
				"price":  19.99,
			},
			expectError: true,
			errorField:  "status",
		},
		{
			name: "Negative price",
			data: map[string]interface{}{
				"title":  "Valid Book",
				"author": "Valid Author",
				"isbn":   "9785000000004",
				"status": "AVAILABLE",
				"price":  -10.0,
			},
			expectError: true,
			errorField:  "price",
		},
		{
			name: "Title too long",
			data: map[string]interface{}{
				"title":  string(make([]byte, 300)), // 300 chars, max is 255
				"author": "Valid Author",
				"isbn":   "9785000000005",
				"status": "AVAILABLE",
			},
			expectError: true,
			errorField:  "title",
		},
	}

	for _, tc := range testCases {
		resp, result := s.makeRequest("POST", "/api/books", tc.data)

		if tc.expectError {
			s.Equal(http.StatusUnprocessableEntity, resp.StatusCode, tc.name)
			s.False(result["success"].(bool), tc.name)
			s.Contains(strings.ToLower(result["message"].(string)), "validation", tc.name)
		} else {
			if resp.StatusCode != http.StatusCreated {
				// Print the error details for debugging
				fmt.Printf("Test %s failed with status %d: %+v\n", tc.name, resp.StatusCode, result)
			}
			s.Equal(http.StatusCreated, resp.StatusCode, tc.name)
			s.True(result["success"].(bool), tc.name)
		}
	}
}
