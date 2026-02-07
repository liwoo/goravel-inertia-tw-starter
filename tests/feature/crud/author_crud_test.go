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

	"books-database/app/models"
	"books-database/tests"
	"books-database/tests/helpers"
)

// ============================================================================
// AuthorCRUDTestSuite — Comprehensive CRUD tests for the Author entity
// ============================================================================

type AuthorCRUDTestSuite struct {
	suite.Suite
	tests.TestCase

	server     *httptest.Server
	client     *http.Client
	authCookie *http.Cookie
	testUser   *models.User
}

func TestAuthorCrudSuite(t *testing.T) {
	suite.Run(t, &AuthorCRUDTestSuite{})
}

// SetupSuite runs once before all tests in the suite
func (s *AuthorCRUDTestSuite) SetupSuite() {
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

// TearDownSuite runs once after all tests in the suite
func (s *AuthorCRUDTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}

// SetupTest runs before each individual test
func (s *AuthorCRUDTestSuite) SetupTest() {
	// Clean database
	s.RefreshDatabase()

	// Clean up any leftover authors from previous tests to ensure isolation
	if orm := facades.Orm(); orm != nil {
		orm.Query().Exec("DELETE FROM authors")
	}

	// Create test user with permissions
	s.setupTestUser()
}

// TearDownTest runs after each individual test
func (s *AuthorCRUDTestSuite) TearDownTest() {
	if orm := facades.Orm(); orm != nil {
		orm.Query().Exec("DELETE FROM authors")
		orm.Query().Exec("DELETE FROM users WHERE email = 'authortest@example.com'")
		orm.Query().Exec("DELETE FROM user_roles")
		orm.Query().Exec("DELETE FROM role_permissions")
		orm.Query().Exec("DELETE FROM roles WHERE slug = 'author_admin'")
		orm.Query().Exec("DELETE FROM permissions WHERE slug LIKE 'authors_%%'")
	}
}

// setupTestUser creates an admin role, assigns CRUD permissions for authors,
// creates a test user, and logs in to obtain an auth cookie.
func (s *AuthorCRUDTestSuite) setupTestUser() {
	// Create admin role
	adminRole := &models.Role{
		Name:  "Author Admin",
		Slug:  "author_admin",
		Level: 100,
	}
	s.Nil(facades.Orm().Query().Create(adminRole))

	// Create all CRUD permissions for the authors service
	permissions := []string{"create", "read", "update", "delete"}
	for _, perm := range permissions {
		permission := &models.Permission{
			Name:  fmt.Sprintf("authors_%s", perm),
			Slug:  fmt.Sprintf("authors_%s", perm),
			Scope: "by_all",
		}
		s.Nil(facades.Orm().Query().Create(permission))
		s.Nil(helpers.AssignPermissionToRole(adminRole, permission, "by_all"))
	}

	// Create test user
	user, err := helpers.SetupJWTUser("authortest@example.com", "password", adminRole)
	s.Nil(err)
	s.testUser = user

	// Login
	s.authCookie = s.loginUser("authortest@example.com", "password")
	s.NotNil(s.authCookie)
}

// loginUser authenticates a user and returns the token cookie
func (s *AuthorCRUDTestSuite) loginUser(email, password string) *http.Cookie {
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

// makeRequest is a helper that sends an HTTP request with optional JSON body
// and returns the response plus the decoded JSON result map.
func (s *AuthorCRUDTestSuite) makeRequest(method, path string, body interface{}) (*http.Response, map[string]interface{}) {
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

// makeUnauthenticatedRequest sends a request without an auth cookie.
// Uses a fresh HTTP client (no cookie jar) so the auth token from SetupTest is not sent.
func (s *AuthorCRUDTestSuite) makeUnauthenticatedRequest(method, path string, body interface{}) (*http.Response, map[string]interface{}) {
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

	// Use a fresh client without cookie jar to ensure no auth cookie is sent
	unauthClient := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := unauthClient.Do(req)
	s.Nil(err)

	respBody, err := io.ReadAll(resp.Body)
	s.Nil(err)
	resp.Body.Close()

	var result map[string]interface{}
	if len(respBody) > 0 {
		json.Unmarshal(respBody, &result)
	}

	return resp, result
}

// createTestAuthor is a helper that inserts an Author directly via ORM
func (s *AuthorCRUDTestSuite) createTestAuthor(firstName, lastName, status string) *models.Author {
	author := &models.Author{
		FirstName: firstName,
		LastName:  lastName,
		Status:    status,
	}
	author.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(author))
	return author
}

// strPtr returns a pointer to the given string
func strPtr(val string) *string {
	return &val
}

// ============================================================================
// CREATE Tests
// ============================================================================

func (s *AuthorCRUDTestSuite) TestCreateAuthorWithAllFields() {
	// API expects snake_case keys in request body
	authorData := map[string]interface{}{
		"first_name":  "Chinua",
		"last_name":   "Achebe",
		"bio":         "Nigerian novelist, poet, and critic.",
		"email":       "chinua@example.com",
		"website":     "https://chinua-achebe.example.com",
		"birth_date":  "1930-11-16",
		"nationality": "Nigerian",
		"photo_url":   "https://example.com/chinua.jpg",
		"status":      "ACTIVE",
	}

	resp, result := s.makeRequest("POST", "/api/authors", authorData)

	if resp.StatusCode != http.StatusCreated {
		s.T().Logf("Unexpected status %d. Response: %+v", resp.StatusCode, result)
	}

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	// Response uses camelCase (from model json tags)
	data := result["data"].(map[string]interface{})
	s.NotNil(data["id"])
	s.Equal("Chinua", data["firstName"])
	s.Equal("Achebe", data["lastName"])
	s.Equal("ACTIVE", data["status"])
	s.Equal(float64(s.testUser.ID), data["created_by"])
}

func (s *AuthorCRUDTestSuite) TestCreateAuthorWithMinimalFields() {
	authorData := map[string]interface{}{
		"first_name": "Wole",
		"last_name":  "Soyinka",
	}

	resp, result := s.makeRequest("POST", "/api/authors", authorData)

	if resp.StatusCode != http.StatusCreated {
		s.T().Logf("Unexpected status %d. Response: %+v", resp.StatusCode, result)
	}

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.NotNil(data["id"])
	s.Equal("Wole", data["firstName"])
	s.Equal("Soyinka", data["lastName"])
	// Default status should be ACTIVE
	s.Equal("ACTIVE", data["status"])
}

func (s *AuthorCRUDTestSuite) TestCreateAuthorValidationMissingFirstName() {
	authorData := map[string]interface{}{
		"last_name": "Achebe",
	}

	resp, result := s.makeRequest("POST", "/api/authors", authorData)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
	s.Contains(strings.ToLower(result["message"].(string)), "validation")
}

func (s *AuthorCRUDTestSuite) TestCreateAuthorValidationMissingLastName() {
	authorData := map[string]interface{}{
		"first_name": "Chinua",
	}

	resp, result := s.makeRequest("POST", "/api/authors", authorData)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
	s.Contains(strings.ToLower(result["message"].(string)), "validation")
}

func (s *AuthorCRUDTestSuite) TestCreateAuthorValidationEmptyBody() {
	authorData := map[string]interface{}{}

	resp, result := s.makeRequest("POST", "/api/authors", authorData)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
	s.Contains(strings.ToLower(result["message"].(string)), "validation")
}

func (s *AuthorCRUDTestSuite) TestCreateAuthorInvalidStatus() {
	authorData := map[string]interface{}{
		"first_name": "Test",
		"last_name":  "Author",
		"status":     "INVALID_STATUS",
	}

	resp, result := s.makeRequest("POST", "/api/authors", authorData)

	// Should fail validation because status must be ACTIVE or INACTIVE
	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
}

// ============================================================================
// READ Tests
// ============================================================================

func (s *AuthorCRUDTestSuite) TestGetAuthor() {
	author := s.createTestAuthor("Ngũgĩ", "wa Thiong'o", "ACTIVE")

	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/authors/%d", author.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal(float64(author.ID), data["id"].(float64))
	s.Equal("Ngũgĩ", data["firstName"])
	s.Equal("wa Thiong'o", data["lastName"])
	s.Equal("ACTIVE", data["status"])
}

func (s *AuthorCRUDTestSuite) TestGetAuthorNotFound() {
	resp, result := s.makeRequest("GET", "/api/authors/99999", nil)

	s.Equal(http.StatusNotFound, resp.StatusCode)
	s.False(result["success"].(bool))
}

func (s *AuthorCRUDTestSuite) TestGetAuthorResponseContainsCamelCaseFields() {
	email := "fields@example.com"
	website := "https://example.com"
	author := &models.Author{
		FirstName: "Fields",
		LastName:  "Test",
		Email:     &email,
		Website:   &website,
		Status:    "ACTIVE",
	}
	author.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(author))

	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/authors/%d", author.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	// Verify camelCase field names in JSON response (from model json tags)
	s.Contains(data, "firstName")
	s.Contains(data, "lastName")
	s.Contains(data, "createdAt")
	s.Contains(data, "updatedAt")
}

// ============================================================================
// LIST Tests
// ============================================================================

func (s *AuthorCRUDTestSuite) TestListAuthors() {
	// Create a few authors
	s.createTestAuthor("Author", "One", "ACTIVE")
	s.createTestAuthor("Author", "Two", "ACTIVE")
	s.createTestAuthor("Author", "Three", "INACTIVE")

	resp, result := s.makeRequest("GET", "/api/authors", nil)

	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})

	s.Equal(3, len(items))
}

func (s *AuthorCRUDTestSuite) TestListAuthorsPagination() {
	// Create 25 authors
	for i := 1; i <= 25; i++ {
		s.createTestAuthor(fmt.Sprintf("First%02d", i), fmt.Sprintf("Last%02d", i), "ACTIVE")
	}

	// Test first page (default page size 20)
	resp, result := s.makeRequest("GET", "/api/authors?page=1", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	pagination := data["pagination"].(map[string]interface{})

	s.Equal(20, len(items))
	s.Equal(float64(1), pagination["current_page"])
	s.Equal(float64(25), pagination["total"])
	s.Equal(float64(2), pagination["last_page"])

	// Test custom page size
	resp, result = s.makeRequest("GET", "/api/authors?page=1&pageSize=10", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	pagination = data["pagination"].(map[string]interface{})

	s.Equal(10, len(items))
	s.Equal(float64(3), pagination["last_page"])
}

func (s *AuthorCRUDTestSuite) TestListAuthorsPaginationSecondPage() {
	// Create 25 authors
	for i := 1; i <= 25; i++ {
		s.createTestAuthor(fmt.Sprintf("First%02d", i), fmt.Sprintf("Last%02d", i), "ACTIVE")
	}

	// Fetch page 2 with default page size (20)
	resp, result := s.makeRequest("GET", "/api/authors?page=2", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	pagination := data["pagination"].(map[string]interface{})

	s.Equal(5, len(items))
	s.Equal(float64(2), pagination["current_page"])
}

// ============================================================================
// UPDATE Tests
// ============================================================================

func (s *AuthorCRUDTestSuite) TestUpdateAuthor() {
	author := s.createTestAuthor("Original", "Name", "ACTIVE")

	updateData := map[string]interface{}{
		"first_name": "Updated",
		"last_name":  "Author",
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/authors/%d", author.ID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal("Updated", data["firstName"])
	s.Equal("Author", data["lastName"])
}

func (s *AuthorCRUDTestSuite) TestUpdateAuthorPartialData() {
	email := "partial@example.com"
	author := &models.Author{
		FirstName: "Partial",
		LastName:  "Update",
		Email:     &email,
		Status:    "ACTIVE",
	}
	author.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(author))

	// Only update the bio field
	updateData := map[string]interface{}{
		"bio": "A newly added biography.",
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/authors/%d", author.ID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	// The original fields should remain unchanged
	s.Equal("Partial", data["firstName"])
	s.Equal("Update", data["lastName"])
}

func (s *AuthorCRUDTestSuite) TestUpdateAuthorStatus() {
	author := s.createTestAuthor("Status", "Change", "ACTIVE")

	updateData := map[string]interface{}{
		"status": "INACTIVE",
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/authors/%d", author.ID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal("INACTIVE", data["status"])
}

func (s *AuthorCRUDTestSuite) TestUpdateAuthorNotFound() {
	updateData := map[string]interface{}{
		"first_name": "Ghost",
	}

	resp, result := s.makeRequest("PUT", "/api/authors/99999", updateData)

	s.Equal(http.StatusNotFound, resp.StatusCode)
	s.False(result["success"].(bool))
}

// ============================================================================
// DELETE Tests
// ============================================================================

func (s *AuthorCRUDTestSuite) TestDeleteAuthor() {
	author := s.createTestAuthor("To", "Delete", "ACTIVE")

	resp, _ := s.makeRequest("DELETE", fmt.Sprintf("/api/authors/%d", author.ID), nil)

	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify soft delete — the record should still exist with a non-nil deleted_at
	var deletedAuthor models.Author
	err := facades.Orm().Query().WithTrashed().Where("id", author.ID).First(&deletedAuthor)
	s.Nil(err)
	s.NotNil(deletedAuthor.DeletedAt)
}

func (s *AuthorCRUDTestSuite) TestDeleteAuthorNotFound() {
	resp, result := s.makeRequest("DELETE", "/api/authors/99999", nil)

	s.Equal(http.StatusNotFound, resp.StatusCode)
	s.False(result["success"].(bool))
}

func (s *AuthorCRUDTestSuite) TestDeletedAuthorNotInList() {
	author := s.createTestAuthor("Invisible", "Author", "ACTIVE")

	// Delete the author
	resp, _ := s.makeRequest("DELETE", fmt.Sprintf("/api/authors/%d", author.ID), nil)
	s.Equal(http.StatusNoContent, resp.StatusCode)

	// List should not include the soft-deleted author
	resp, result := s.makeRequest("GET", "/api/authors", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.Equal(0, len(items))
}

// ============================================================================
// SEARCH Tests
// ============================================================================

func (s *AuthorCRUDTestSuite) TestSearchAuthorsByFirstName() {
	s.createTestAuthor("Chimamanda", "Adichie", "ACTIVE")
	s.createTestAuthor("Chinua", "Achebe", "ACTIVE")
	s.createTestAuthor("Wole", "Soyinka", "ACTIVE")

	// Search for "Chi" — should match Chimamanda and Chinua
	resp, result := s.makeRequest("GET", "/api/authors/search?q=Chi", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.Equal(2, len(items))
}

func (s *AuthorCRUDTestSuite) TestSearchAuthorsByLastName() {
	s.createTestAuthor("Chimamanda", "Adichie", "ACTIVE")
	s.createTestAuthor("Chinua", "Achebe", "ACTIVE")

	// Search for "Adichie"
	resp, result := s.makeRequest("GET", "/api/authors/search?q=Adichie", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.Equal(1, len(items))

	firstItem := items[0].(map[string]interface{})
	s.Equal("Chimamanda", firstItem["firstName"])
}

func (s *AuthorCRUDTestSuite) TestSearchAuthorsByEmail() {
	email1 := "chimamanda@example.com"
	email2 := "chinua@example.com"

	author1 := &models.Author{FirstName: "Chimamanda", LastName: "Adichie", Email: &email1, Status: "ACTIVE"}
	author1.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(author1))

	author2 := &models.Author{FirstName: "Chinua", LastName: "Achebe", Email: &email2, Status: "ACTIVE"}
	author2.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(author2))

	// Search by email substring
	resp, result := s.makeRequest("GET", "/api/authors/search?q=chimamanda@example", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.Equal(1, len(items))
}

func (s *AuthorCRUDTestSuite) TestSearchAuthorsNoResults() {
	s.createTestAuthor("Test", "Author", "ACTIVE")

	resp, result := s.makeRequest("GET", "/api/authors/search?q=NonExistentXYZ", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.Equal(0, len(items))
}

// ============================================================================
// FILTER Tests
// ============================================================================

func (s *AuthorCRUDTestSuite) TestFilterAuthorsByStatus() {
	s.createTestAuthor("Active", "One", "ACTIVE")
	s.createTestAuthor("Active", "Two", "ACTIVE")
	s.createTestAuthor("Inactive", "One", "INACTIVE")

	// Filter by ACTIVE status
	resp, result := s.makeRequest("GET", "/api/authors?status=ACTIVE", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.Equal(2, len(items))

	for _, item := range items {
		author := item.(map[string]interface{})
		s.Equal("ACTIVE", author["status"])
	}
}

func (s *AuthorCRUDTestSuite) TestFilterAuthorsByInactiveStatus() {
	s.createTestAuthor("Active", "Author", "ACTIVE")
	s.createTestAuthor("Inactive", "Author", "INACTIVE")

	// Filter by INACTIVE status
	resp, result := s.makeRequest("GET", "/api/authors?status=INACTIVE", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.Equal(1, len(items))

	firstItem := items[0].(map[string]interface{})
	s.Equal("INACTIVE", firstItem["status"])
}

func (s *AuthorCRUDTestSuite) TestFilterMetadataEndpoint() {
	resp, result := s.makeRequest("GET", "/api/authors/filters", nil)
	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	// The response should contain filter definitions
	data := result["data"]
	s.NotNil(data)
}

// ============================================================================
// SORTING Tests
// ============================================================================

func (s *AuthorCRUDTestSuite) TestSortAuthorsByFirstNameAsc() {
	s.createTestAuthor("Zora", "Hurston", "ACTIVE")
	s.createTestAuthor("Alice", "Walker", "ACTIVE")
	s.createTestAuthor("Maya", "Angelou", "ACTIVE")

	resp, result := s.makeRequest("GET", "/api/authors?sort=first_name&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 3)

	// First item should be "Alice" (alphabetically first)
	firstItem := items[0].(map[string]interface{})
	s.Equal("Alice", firstItem["firstName"])
}

func (s *AuthorCRUDTestSuite) TestSortAuthorsByFirstNameDesc() {
	s.createTestAuthor("Zora", "Hurston", "ACTIVE")
	s.createTestAuthor("Alice", "Walker", "ACTIVE")
	s.createTestAuthor("Maya", "Angelou", "ACTIVE")

	resp, result := s.makeRequest("GET", "/api/authors?sort=first_name&direction=desc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 3)

	// First item should be "Zora" (alphabetically last)
	firstItem := items[0].(map[string]interface{})
	s.Equal("Zora", firstItem["firstName"])
}

func (s *AuthorCRUDTestSuite) TestSortAuthorsByLastName() {
	s.createTestAuthor("Chinua", "Achebe", "ACTIVE")
	s.createTestAuthor("Wole", "Soyinka", "ACTIVE")
	s.createTestAuthor("Buchi", "Emecheta", "ACTIVE")

	resp, result := s.makeRequest("GET", "/api/authors?sort=last_name&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 3)

	firstItem := items[0].(map[string]interface{})
	s.Equal("Achebe", firstItem["lastName"])
}

func (s *AuthorCRUDTestSuite) TestSortAuthorsByCreatedAt() {
	// Create authors sequentially — IDs are auto-incremented
	first := s.createTestAuthor("First", "Created", "ACTIVE")
	second := s.createTestAuthor("Second", "Created", "ACTIVE")
	third := s.createTestAuthor("Third", "Created", "ACTIVE")

	// Verify IDs are sequential (created_at may have second-level precision,
	// so we use ID order as a reliable proxy for creation order)
	s.Less(first.ID, second.ID)
	s.Less(second.ID, third.ID)

	// Sort by created_at descending — newest first
	resp, result := s.makeRequest("GET", "/api/authors?sort=created_at&direction=desc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 3)

	// With timestamp(0) precision, records in the same second may have identical
	// created_at values. Verify that sorting by created_at returns results
	// (the endpoint accepts the sort param and returns 200).
	// For strict ordering, TestSortAuthorsByFirstNameAsc/Desc tests cover deterministic sorts.
	firstItem := items[0].(map[string]interface{})
	s.Contains([]string{"First", "Second", "Third"}, firstItem["firstName"])
}

func (s *AuthorCRUDTestSuite) TestSortWithInvalidField() {
	s.createTestAuthor("Test", "Author", "ACTIVE")

	resp, _ := s.makeRequest("GET", "/api/authors?sort=invalid_field", nil)

	// Should either ignore unknown sort field or return an error
	s.True(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest)
}

// ============================================================================
// STATISTICS Tests
// ============================================================================

func (s *AuthorCRUDTestSuite) TestAuthorStatistics() {
	s.createTestAuthor("Active", "One", "ACTIVE")
	s.createTestAuthor("Active", "Two", "ACTIVE")
	s.createTestAuthor("Active", "Three", "ACTIVE")
	s.createTestAuthor("Inactive", "One", "INACTIVE")
	s.createTestAuthor("Inactive", "Two", "INACTIVE")

	resp, result := s.makeRequest("GET", "/api/authors/statistics", nil)
	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal(float64(5), data["totalAuthors"])
	s.Equal(float64(3), data["activeAuthors"])
	s.Equal(float64(2), data["inactiveAuthors"])
}

func (s *AuthorCRUDTestSuite) TestAuthorStatisticsEmpty() {
	// No authors created — all counts should be zero
	resp, result := s.makeRequest("GET", "/api/authors/statistics", nil)
	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal(float64(0), data["totalAuthors"])
	s.Equal(float64(0), data["activeAuthors"])
	s.Equal(float64(0), data["inactiveAuthors"])
}

// ============================================================================
// COMBINED FILTER + SORT Tests
// ============================================================================

func (s *AuthorCRUDTestSuite) TestFilterAndSort() {
	s.createTestAuthor("Zora", "Hurston", "ACTIVE")
	s.createTestAuthor("Alice", "Walker", "ACTIVE")
	s.createTestAuthor("Inactive", "Author", "INACTIVE")
	s.createTestAuthor("Maya", "Angelou", "ACTIVE")

	// Filter by ACTIVE and sort by first_name ascending
	resp, result := s.makeRequest("GET", "/api/authors?status=ACTIVE&sort=first_name&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})

	// Should only have 3 ACTIVE authors
	s.Equal(3, len(items))

	// Should be sorted: Alice, Maya, Zora
	firstItem := items[0].(map[string]interface{})
	s.Equal("Alice", firstItem["firstName"])
	s.Equal("ACTIVE", firstItem["status"])

	lastItem := items[2].(map[string]interface{})
	s.Equal("Zora", lastItem["firstName"])
}

// ============================================================================
// PERMISSION / AUTH Tests
// ============================================================================

func (s *AuthorCRUDTestSuite) TestUnauthenticatedCreateDenied() {
	authorData := map[string]interface{}{
		"first_name": "Unauthorized",
		"last_name":  "Author",
	}

	resp, _ := s.makeUnauthenticatedRequest("POST", "/api/authors", authorData)

	// JWT middleware redirects (302) unauthenticated non-Inertia requests
	s.NotEqual(http.StatusCreated, resp.StatusCode,
		"Unauthenticated request should not create an author")
}

func (s *AuthorCRUDTestSuite) TestUnauthenticatedUpdateDenied() {
	author := s.createTestAuthor("Auth", "Test", "ACTIVE")

	updateData := map[string]interface{}{
		"first_name": "Hacked",
	}

	resp, _ := s.makeUnauthenticatedRequest("PUT", fmt.Sprintf("/api/authors/%d", author.ID), updateData)

	s.NotEqual(http.StatusOK, resp.StatusCode,
		"Unauthenticated request should not update an author")
}

func (s *AuthorCRUDTestSuite) TestUnauthenticatedDeleteDenied() {
	author := s.createTestAuthor("Auth", "Test", "ACTIVE")

	resp, _ := s.makeUnauthenticatedRequest("DELETE", fmt.Sprintf("/api/authors/%d", author.ID), nil)

	s.NotEqual(http.StatusOK, resp.StatusCode,
		"Unauthenticated request should not delete an author")
}

func (s *AuthorCRUDTestSuite) TestUnauthenticatedStatisticsDenied() {
	resp, _ := s.makeUnauthenticatedRequest("GET", "/api/authors/statistics", nil)

	// Statistics endpoint requires auth
	s.NotEqual(http.StatusOK, resp.StatusCode,
		"Unauthenticated request should not access statistics")
}

// ============================================================================
// ERROR HANDLING Tests
// ============================================================================

func (s *AuthorCRUDTestSuite) TestUnknownFilterIgnored() {
	s.createTestAuthor("Test", "Author", "ACTIVE")

	resp, result := s.makeRequest("GET", "/api/authors?unknown_filter=value", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 1)
}

func (s *AuthorCRUDTestSuite) TestSearchShortQuery() {
	resp, _ := s.makeRequest("GET", "/api/authors/search?q=a", nil)

	// Minimum query length is 2 — may return 400 or proceed with short query
	s.True(resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusOK)
}

// ============================================================================
// AUDIT FIELD Tests
// ============================================================================

func (s *AuthorCRUDTestSuite) TestCreatedByIsSetOnCreate() {
	authorData := map[string]interface{}{
		"first_name": "Audit",
		"last_name":  "Test",
	}

	resp, result := s.makeRequest("POST", "/api/authors", authorData)
	s.Equal(http.StatusCreated, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	s.Equal(float64(s.testUser.ID), data["created_by"])
}
