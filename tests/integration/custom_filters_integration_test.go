package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"books-database/app/models"
	"books-database/tests"
	"books-database/tests/helpers"
)

type CustomFiltersIntegrationTestSuite struct {
	suite.Suite
	tests.TestCase

	server     *httptest.Server
	client     *http.Client
	authCookie *http.Cookie
	testUser   *models.User
}

func TestCustomFiltersIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &CustomFiltersIntegrationTestSuite{})
}

func (s *CustomFiltersIntegrationTestSuite) SetupSuite() {
	// Start test server
	s.server = httptest.NewServer(facades.Route())
	s.client = &http.Client{
		Timeout: 10 * time.Second,
	}
}

func (s *CustomFiltersIntegrationTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}

func (s *CustomFiltersIntegrationTestSuite) SetupTest() {
	// Clean database
	s.RefreshDatabase()

	// Clear any existing books
	facades.Orm().Query().Exec("DELETE FROM books")

	// Create test user with permissions
	s.setupTestUser()

	// Create test books
	s.createTestBooks()
}

func (s *CustomFiltersIntegrationTestSuite) setupTestUser() {
	// Create admin role
	adminRole := &models.Role{
		Name:  "Admin",
		Slug:  "admin",
		Level: 100,
	}
	s.Nil(facades.Orm().Query().Create(adminRole))

	// Create book permissions
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
	user, err := helpers.SetupJWTUser("admin@test.com", "password", adminRole)
	s.Nil(err)
	s.testUser = user

	// Login
	s.authCookie = s.loginUser("admin@test.com", "password")
}

func (s *CustomFiltersIntegrationTestSuite) loginUser(email, password string) *http.Cookie {
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

func (s *CustomFiltersIntegrationTestSuite) createTestBooks() {
	// Create books with various prices, statuses, and dates
	books := []struct {
		title       string
		author      string
		price       float64
		status      string
		publishedAt *time.Time
	}{
		{"Go Programming", "John Doe", 29.99, "AVAILABLE", timePtr(time.Now().AddDate(-1, 0, 0))},
		{"Python Basics", "Jane Smith", 19.99, "BORROWED", timePtr(time.Now().AddDate(-2, 0, 0))},
		{"JavaScript Guide", "Bob Johnson", 39.99, "AVAILABLE", timePtr(time.Now().AddDate(0, -6, 0))},
		{"Rust Programming", "Alice Brown", 49.99, "MAINTENANCE", timePtr(time.Now().AddDate(0, -3, 0))},
		{"TypeScript Advanced", "Charlie Wilson", 34.99, "AVAILABLE", timePtr(time.Now())},
	}

	for i, b := range books {
		book := &models.Book{
			Title:       b.title,
			Author:      b.author,
			ISBN:        fmt.Sprintf("978900%07d", i+1),
			Price:       b.price,
			Status:      b.status,
			PublishedAt: b.publishedAt,
		}
		book.CreatedBy = &s.testUser.ID
		s.Nil(facades.Orm().Query().Create(book))
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func (s *CustomFiltersIntegrationTestSuite) makeRequest(method, path string, body interface{}) (*http.Response, map[string]interface{}) {
	var bodyReader *bytes.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(jsonBody)
	} else {
		bodyReader = bytes.NewReader([]byte{})
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

	var result map[string]interface{}
	if resp.Body != nil {
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&result)
	}

	return resp, result
}

// Test filter metadata endpoint
func (s *CustomFiltersIntegrationTestSuite) TestFilterMetadataEndpoint() {
	resp, result := s.makeRequest("GET", "/api/books/filters", nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	// Check metadata structure
	data := result["data"].(map[string]interface{})
	s.NotNil(data["filters"])
	s.NotNil(data["logic_operators"])
	s.NotNil(data["searchable_fields"])
	s.NotNil(data["sortable_fields"])
	s.NotNil(data["filterable_fields"])

	// Check filter definitions
	filters := data["filters"].([]interface{})
	s.Greater(len(filters), 0)

	// Check that price filter exists
	var foundPriceFilter bool
	for _, f := range filters {
		filter := f.(map[string]interface{})
		if filter["field"] == "price" {
			foundPriceFilter = true
			s.Equal("Price", filter["label"])
			s.Equal("number", filter["type"])
			s.NotEmpty(filter["operators"])
		}
	}
	s.True(foundPriceFilter, "Price filter should be defined")
}

// Test simple custom filter
func (s *CustomFiltersIntegrationTestSuite) TestSimpleCustomFilter() {
	// Create filter for price > 30
	filter := map[string]interface{}{
		"field":    "price",
		"operator": "greater_than",
		"value":    30,
	}

	filterJSON, _ := json.Marshal(filter)
	params := url.Values{}
	params.Set("filters", string(filterJSON))

	resp, result := s.makeRequest("GET", "/api/books?"+params.Encode(), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	// Check results
	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})

	// Should have 2 books with price > 30 (JavaScript Guide: 39.99, Rust Programming: 49.99, TypeScript Advanced: 34.99)
	s.Equal(3, len(items))

	// Verify all returned books have price > 30
	for _, item := range items {
		book := item.(map[string]interface{})
		s.Greater(book["price"].(float64), 30.0)
	}
}

// Test compound custom filter
func (s *CustomFiltersIntegrationTestSuite) TestCompoundCustomFilter() {
	// Create compound filter: (price > 25 AND status = 'AVAILABLE')
	filter := map[string]interface{}{
		"logic": "AND",
		"conditions": []map[string]interface{}{
			{
				"field":    "price",
				"operator": "greater_than",
				"value":    25,
			},
			{
				"field":    "status",
				"operator": "equals",
				"value":    "AVAILABLE",
			},
		},
	}

	filterJSON, _ := json.Marshal(filter)
	params := url.Values{}
	params.Set("filters", string(filterJSON))

	resp, result := s.makeRequest("GET", "/api/books?"+params.Encode(), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	// Check results
	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})

	// Should have books with price > 25 AND status = AVAILABLE
	// Go Programming: 29.99, AVAILABLE
	// JavaScript Guide: 39.99, AVAILABLE
	// TypeScript Advanced: 34.99, AVAILABLE
	s.Equal(3, len(items))

	// Verify all returned books match criteria
	for _, item := range items {
		book := item.(map[string]interface{})
		s.Greater(book["price"].(float64), 25.0)
		s.Equal("AVAILABLE", book["status"])
	}
}

// Test date range filter
func (s *CustomFiltersIntegrationTestSuite) TestDateRangeFilter() {
	// Create filter for books published in the last year
	filter := map[string]interface{}{
		"field":    "published_at",
		"operator": "after",
		"value":    time.Now().AddDate(-1, 0, 0).Format("2006-01-02"),
	}

	filterJSON, _ := json.Marshal(filter)
	params := url.Values{}
	params.Set("filters", string(filterJSON))

	resp, result := s.makeRequest("GET", "/api/books?"+params.Encode(), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	// Check results
	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})

	// Should have books published in the last year
	s.Greater(len(items), 0)

	// Verify all returned books were published in the last year
	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	for _, item := range items {
		book := item.(map[string]interface{})
		if publishedAtStr, ok := book["published_at"].(string); ok {
			publishedAt, err := time.Parse(time.RFC3339, publishedAtStr)
			s.Nil(err)
			s.True(publishedAt.After(oneYearAgo))
		}
	}
}
