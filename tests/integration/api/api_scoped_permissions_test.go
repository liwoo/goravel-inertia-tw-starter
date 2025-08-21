package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"players/app/models"
	"players/tests"
)

type APIScopedPermissionsTestSuite struct {
	suite.Suite
	tests.TestCase

	server *httptest.Server
	client *http.Client

	// Users
	admin   *models.User
	editor1 *models.User
	editor2 *models.User
	member1 *models.User
	member2 *models.User

	// Roles
	adminRole  *models.Role
	editorRole *models.Role
	memberRole *models.Role

	// Auth cookies for each user
	adminCookie   *http.Cookie
	editor1Cookie *http.Cookie
	editor2Cookie *http.Cookie
	member1Cookie *http.Cookie
	member2Cookie *http.Cookie
}

func TestAPIScopedPermissionsTestSuite(t *testing.T) {
	suite.Run(t, new(APIScopedPermissionsTestSuite))
}

func (s *APIScopedPermissionsTestSuite) SetupSuite() {
	// Application is already bootstrapped by test runner
}

func (s *APIScopedPermissionsTestSuite) SetupTest() {
	s.RefreshDatabase()
	s.setupTestData()
	s.startTestServer()
	s.loginAllUsers()
}

func (s *APIScopedPermissionsTestSuite) TearDownTest() {
	if s.server != nil {
		s.server.Close()
	}
}

func (s *APIScopedPermissionsTestSuite) startTestServer() {
	// Create a test server
	s.server = httptest.NewServer(facades.Route())

	// Create HTTP client with cookie jar
	jar, _ := cookiejar.New(nil)
	s.client = &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}
}

func (s *APIScopedPermissionsTestSuite) setupTestData() {
	// Create roles
	s.adminRole = &models.Role{Name: "Admin", Slug: "admin", IsActive: true}
	s.editorRole = &models.Role{Name: "Editor", Slug: "editor", IsActive: true}
	s.memberRole = &models.Role{Name: "Member", Slug: "member", IsActive: true}

	facades.Orm().Query().Create(s.adminRole)
	facades.Orm().Query().Create(s.editorRole)
	facades.Orm().Query().Create(s.memberRole)

	// Create base permission - the system will generate scoped versions dynamically
	perm := &models.Permission{
		Slug:     "books_read",
		Name:     "Read Books",
		Resource: "books",
		Action:   "read",
		IsActive: true,
	}
	facades.Orm().Query().Create(perm)

	// Assign permissions with different scopes
	// Admin: by_all
	facades.Orm().Query().Create(&models.RolePermission{
		RoleID:       s.adminRole.ID,
		PermissionID: perm.ID,
		Scope:        "by_all",
		IsActive:     true,
		GrantedAt:    time.Now(),
	})

	// Editor: by_my_role
	facades.Orm().Query().Create(&models.RolePermission{
		RoleID:       s.editorRole.ID,
		PermissionID: perm.ID,
		Scope:        "by_my_role",
		IsActive:     true,
		GrantedAt:    time.Now(),
	})

	// Member: by_me
	facades.Orm().Query().Create(&models.RolePermission{
		RoleID:       s.memberRole.ID,
		PermissionID: perm.ID,
		Scope:        "by_me",
		IsActive:     true,
		GrantedAt:    time.Now(),
	})

	// Create users with a common password
	password, _ := facades.Hash().Make("password123")

	s.admin = &models.User{
		Name:          "Admin User",
		Email:         "admin@test.com",
		Password:      password,
		IsActive:      true,
		EmailVerified: true,
	}
	s.editor1 = &models.User{
		Name:          "Editor One",
		Email:         "editor1@test.com",
		Password:      password,
		IsActive:      true,
		EmailVerified: true,
	}
	s.editor2 = &models.User{
		Name:          "Editor Two",
		Email:         "editor2@test.com",
		Password:      password,
		IsActive:      true,
		EmailVerified: true,
	}
	s.member1 = &models.User{
		Name:          "Member One",
		Email:         "member1@test.com",
		Password:      password,
		IsActive:      true,
		EmailVerified: true,
	}
	s.member2 = &models.User{
		Name:          "Member Two",
		Email:         "member2@test.com",
		Password:      password,
		IsActive:      true,
		EmailVerified: true,
	}

	// Save users
	facades.Orm().Query().Create(s.admin)
	facades.Orm().Query().Create(s.editor1)
	facades.Orm().Query().Create(s.editor2)
	facades.Orm().Query().Create(s.member1)
	facades.Orm().Query().Create(s.member2)

	// Assign roles
	now := time.Now()
	facades.Orm().Query().Create(&models.UserRole{UserID: s.admin.ID, RoleID: s.adminRole.ID, IsActive: true, AssignedAt: now})
	facades.Orm().Query().Create(&models.UserRole{UserID: s.editor1.ID, RoleID: s.editorRole.ID, IsActive: true, AssignedAt: now})
	facades.Orm().Query().Create(&models.UserRole{UserID: s.editor2.ID, RoleID: s.editorRole.ID, IsActive: true, AssignedAt: now})
	facades.Orm().Query().Create(&models.UserRole{UserID: s.member1.ID, RoleID: s.memberRole.ID, IsActive: true, AssignedAt: now})
	facades.Orm().Query().Create(&models.UserRole{UserID: s.member2.ID, RoleID: s.memberRole.ID, IsActive: true, AssignedAt: now})

	// Create books for each user
	s.createBooksForUser(s.admin, 2)
	s.createBooksForUser(s.editor1, 2)
	s.createBooksForUser(s.editor2, 2)
	s.createBooksForUser(s.member1, 2)
	s.createBooksForUser(s.member2, 2)
}

func (s *APIScopedPermissionsTestSuite) createBooksForUser(user *models.User, count int) {
	for i := 0; i < count; i++ {
		book := &models.Book{
			Title:       fmt.Sprintf("%s Book %d", user.Name, i+1),
			Author:      user.Name,
			ISBN:        fmt.Sprintf("ISBN-%d-%d", user.ID, i+1),
			Description: fmt.Sprintf("Book by %s", user.Name),
			Price:       float64(10 + i*5),
			Status:      "AVAILABLE",
		}
		book.CreatedBy = &user.ID
		facades.Orm().Query().Create(book)
	}
}

func (s *APIScopedPermissionsTestSuite) loginAllUsers() {
	s.adminCookie = s.login(s.admin.Email, "password123")
	s.editor1Cookie = s.login(s.editor1.Email, "password123")
	s.editor2Cookie = s.login(s.editor2.Email, "password123")
	s.member1Cookie = s.login(s.member1.Email, "password123")
	s.member2Cookie = s.login(s.member2.Email, "password123")
	
	// Debug: check permissions for admin
	var adminWithRoles models.User
	facades.Orm().Query().Where("id = ?", s.admin.ID).With("Roles").First(&adminWithRoles)
	s.T().Logf("Admin roles: %+v", adminWithRoles.Roles)
	
	// Check role permissions
	var rolePerms []models.RolePermission
	facades.Orm().Query().Where("role_id = ?", s.adminRole.ID).With("Permission").Find(&rolePerms)
	for _, rp := range rolePerms {
		s.T().Logf("Admin role permission: %s (scope: %s)", rp.Permission.Slug, rp.Scope)
	}
}

func (s *APIScopedPermissionsTestSuite) login(email, password string) *http.Cookie {
	loginData := map[string]string{
		"email":    email,
		"password": password,
	}
	body, _ := json.Marshal(loginData)

	req, err := http.NewRequest("POST", s.server.URL+"/api/auth/login", bytes.NewReader(body))
	s.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	s.NoError(err)
	defer resp.Body.Close()
	
	// Read response body
	respBody, _ := io.ReadAll(resp.Body)
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		s.T().Logf("Login failed for %s. Status: %d, Body: %s", email, resp.StatusCode, string(respBody))
		s.Fail("Login failed")
		return nil
	}
	
	// Parse JSON response
	var loginResp map[string]interface{}
	err = json.Unmarshal(respBody, &loginResp)
	s.NoError(err, "Failed to parse login response")
	
	// Check success
	success, ok := loginResp["success"].(bool)
	s.True(ok && success, "Login was not successful")
	
	// Get token from response
	data, ok := loginResp["data"].(map[string]interface{})
	s.True(ok, "No data in login response")
	
	token, ok := data["token"].(string)
	s.True(ok, "No token in login response")
	s.NotEmpty(token, "Token is empty")

	// Look for the token cookie (should also be set)
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "token" {
			s.Equal(token, cookie.Value, "Cookie token should match response token")
			return cookie
		}
	}

	// If no cookie found, create one from the token
	return &http.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
	}
}

// Helper to make authenticated API request
func (s *APIScopedPermissionsTestSuite) apiGet(path string, authCookie *http.Cookie) (map[string]interface{}, int) {
	req, err := http.NewRequest("GET", s.server.URL+path, nil)
	s.NoError(err)

	if authCookie != nil {
		// Add both cookie and Authorization header for maximum compatibility
		req.AddCookie(authCookie)
		req.Header.Set("Authorization", "Bearer "+authCookie.Value)
	}

	resp, err := s.client.Do(req)
	s.NoError(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	s.NoError(err)

	var result map[string]interface{}
	if len(body) > 0 {
		err = json.Unmarshal(body, &result)
		if err != nil {
			// If not JSON, return raw response
			result = map[string]interface{}{
				"body": string(body),
			}
		}
	}

	return result, resp.StatusCode
}

// Test that protected API endpoints require authentication
func (s *APIScopedPermissionsTestSuite) TestUnauthenticatedAccessDenied() {
	// Try to create a book without authentication (POST is protected)
	bookData := map[string]interface{}{
		"title":  "Test Book",
		"author": "Test Author",
	}
	body, _ := json.Marshal(bookData)
	
	req, err := http.NewRequest("POST", s.server.URL+"/api/books", bytes.NewReader(body))
	s.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := s.client.Do(req)
	s.NoError(err)
	defer resp.Body.Close()
	
	// Should return 302 redirect to login, 401 unauthorized, or 403 forbidden
	s.True(resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden, 
		"Unauthenticated POST request should return 302, 401, or 403, got %d", resp.StatusCode)
}

// Test admin sees all books
func (s *APIScopedPermissionsTestSuite) TestAdminSeesAllBooks() {
	// First, test if the auth endpoint works
	meResult, meStatus := s.apiGet("/api/auth/me", s.adminCookie)
	s.T().Logf("Me endpoint status: %d, response: %+v", meStatus, meResult)
	
	// Check that auth really works
	s.Equal(http.StatusOK, meStatus, "Auth should work")
	
	result, status := s.apiGet("/api/books", s.adminCookie)
	s.T().Logf("Books endpoint status: %d, response: %+v", status, result)
	
	// If we get a 403, it means the auth is required but failing
	if status == http.StatusForbidden {
		s.T().Logf("Got 403 Forbidden. Response: %+v", result)
		s.T().Logf("Admin cookie: %+v", s.adminCookie)
		s.T().Logf("Token value: %s", s.adminCookie.Value)
		s.Fail("Admin should have access to books")
		return
	}
	
	s.Equal(http.StatusOK, status, "Expected 200 OK, got %d", status)

	// Check the response structure - handle nested data.data structure
	if dataWrapper, ok := result["data"].(map[string]interface{}); ok {
		// Check for nested data array
		if data, ok := dataWrapper["data"].([]interface{}); ok {
			s.Equal(10, len(data), "Admin should see all 10 books")
		} else {
			s.T().Logf("No data array in wrapper: %+v", dataWrapper)
			s.Fail("Expected data array in response")
		}
		
		// Also check pagination info if available
		if pagination, ok := dataWrapper["pagination"].(map[string]interface{}); ok {
			if total, ok := pagination["total"].(float64); ok {
				s.Equal(float64(10), total, "Pagination total should be 10")
			}
		}
	} else if data, ok := result["data"].([]interface{}); ok {
		// Fallback to direct data array
		s.Equal(10, len(data), "Admin should see all 10 books")
	} else if items, ok := result["items"].([]interface{}); ok {
		s.Equal(10, len(items), "Admin should see all 10 books")
	} else if meta, ok := result["meta"].(map[string]interface{}); ok {
		if total, ok := meta["total"].(float64); ok {
			s.Equal(float64(10), total, "Admin should see total of 10 books")
		}
	} else {
		// Log the actual response structure
		s.T().Logf("Response structure: %+v", result)
		s.Fail("Unexpected response structure")
	}
}

// Test editor sees only editor role books
func (s *APIScopedPermissionsTestSuite) TestEditorSeesOnlyEditorBooks() {
	result, status := s.apiGet("/api/books", s.editor1Cookie)
	s.Equal(http.StatusOK, status)

	// Editor with by_my_role should see books from both editors (4 books)
	if dataWrapper, ok := result["data"].(map[string]interface{}); ok {
		// Check for nested data array
		if data, ok := dataWrapper["data"].([]interface{}); ok {
			s.Equal(4, len(data), "Editor should see 4 books (all editor books)")
			
			// Verify all books are from editors
			for _, item := range data {
				book := item.(map[string]interface{})
				author := book["author"].(string)
				s.True(author == "Editor One" || author == "Editor Two", 
					"Book should be from an editor, got: %s", author)
			}
		}
		
		// Check pagination total
		if pagination, ok := dataWrapper["pagination"].(map[string]interface{}); ok {
			if total, ok := pagination["total"].(float64); ok {
				s.Equal(float64(4), total, "Pagination total should be 4")
			}
		}
	} else if data, ok := result["data"].([]interface{}); ok {
		// Fallback to direct data array
		s.Equal(4, len(data), "Editor should see 4 books (all editor books)")
		
		// Verify all books are from editors
		for _, item := range data {
			book := item.(map[string]interface{})
			author := book["author"].(string)
			s.True(author == "Editor One" || author == "Editor Two", 
				"Book should be from an editor, got: %s", author)
		}
	} else if total, ok := result["meta"].(map[string]interface{})["total"].(float64); ok {
		s.Equal(float64(4), total, "Editor should see total of 4 books")
	}
}

// Test member sees only their own books
func (s *APIScopedPermissionsTestSuite) TestMemberSeesOnlyOwnBooks() {
	result, status := s.apiGet("/api/books", s.member1Cookie)
	s.Equal(http.StatusOK, status)

	// Member with by_me should see only their 2 books
	if dataWrapper, ok := result["data"].(map[string]interface{}); ok {
		// Check for nested data array
		if data, ok := dataWrapper["data"].([]interface{}); ok {
			s.Equal(2, len(data), "Member should see only their 2 books")
			
			// Verify all books belong to member1
			for _, item := range data {
				book := item.(map[string]interface{})
				author := book["author"].(string)
				s.Equal("Member One", author, "Book should belong to Member One")
			}
		}
		
		// Check pagination total
		if pagination, ok := dataWrapper["pagination"].(map[string]interface{}); ok {
			if total, ok := pagination["total"].(float64); ok {
				s.Equal(float64(2), total, "Pagination total should be 2")
			}
		}
	} else if data, ok := result["data"].([]interface{}); ok {
		// Fallback to direct data array
		s.Equal(2, len(data), "Member should see only their 2 books")
		
		// Verify all books belong to member1
		for _, item := range data {
			book := item.(map[string]interface{})
			author := book["author"].(string)
			s.Equal("Member One", author, "Book should belong to Member One")
		}
	} else if total, ok := result["meta"].(map[string]interface{})["total"].(float64); ok {
		s.Equal(float64(2), total, "Member should see total of 2 books")
	}
}

// Test that different members see different books
func (s *APIScopedPermissionsTestSuite) TestDifferentMembersSeeOwnBooks() {
	// Member1's books
	result1, status1 := s.apiGet("/api/books", s.member1Cookie)
	s.Equal(http.StatusOK, status1)

	// Member2's books
	result2, status2 := s.apiGet("/api/books", s.member2Cookie)
	s.Equal(http.StatusOK, status2)

	// Get the data arrays - handle nested structure
	var data1, data2 []interface{}
	
	if wrapper1, ok := result1["data"].(map[string]interface{}); ok {
		data1, _ = wrapper1["data"].([]interface{})
	} else {
		data1, _ = result1["data"].([]interface{})
	}
	
	if wrapper2, ok := result2["data"].(map[string]interface{}); ok {
		data2, _ = wrapper2["data"].([]interface{})
	} else {
		data2, _ = result2["data"].([]interface{})
	}

	s.Equal(2, len(data1), "Member1 should see 2 books")
	s.Equal(2, len(data2), "Member2 should see 2 books")

	// Verify they see different books
	book1 := data1[0].(map[string]interface{})
	book2 := data2[0].(map[string]interface{})
	
	s.NotEqual(book1["isbn"], book2["isbn"], "Members should see different books")
	s.Equal("Member One", book1["author"].(string))
	s.Equal("Member Two", book2["author"].(string))
}

// Test pagination respects scopes
func (s *APIScopedPermissionsTestSuite) TestPaginationRespectsScopes() {
	// Test with page size = 5 (smallest allowed page size)
	result, status := s.apiGet("/api/books?page=1&pageSize=5", s.editor1Cookie)
	s.Equal(http.StatusOK, status)
	
	// Debug: Log the full response
	s.T().Logf("Pagination test response: %+v", result)

	// Should still show total of 4 for editor
	if dataWrapper, ok := result["data"].(map[string]interface{}); ok {
		// Check pagination info
		if pagination, ok := dataWrapper["pagination"].(map[string]interface{}); ok {
			s.Equal(float64(4), pagination["total"].(float64), "Total should be 4 for editor")
			s.Equal(float64(5), pagination["per_page"].(float64), "Per page should be 5")
			s.Equal(float64(1), pagination["last_page"].(float64), "Should have 1 page (4 items fit in page size 5)")
		}
		
		// Also check filters to see what was passed
		if filters, ok := dataWrapper["filters"].(map[string]interface{}); ok {
			s.T().Logf("Filters in response: %+v", filters)
			// The filters should show the pageSize that was requested
			if pageSize, ok := filters["pageSize"].(float64); ok {
				s.Equal(float64(5), pageSize, "Filters should show pageSize of 5")
			}
		}
		
		// Data should have all 4 items (since 4 items fit in page size 5)
		if data, ok := dataWrapper["data"].([]interface{}); ok {
			s.Equal(4, len(data), "Should return all 4 books (fits in page size 5)")
		}
	} else if meta, ok := result["meta"].(map[string]interface{}); ok {
		// Fallback to meta structure
		s.Equal(float64(4), meta["total"].(float64), "Total should be 4 for editor")
		s.Equal(float64(5), meta["per_page"].(float64), "Per page should be 5")
		s.Equal(float64(1), meta["last_page"].(float64), "Should have 1 page")
		
		// Data should have all 4 items
		if data, ok := result["data"].([]interface{}); ok {
			s.Equal(4, len(data), "Should return all 4 books")
		}
	}
}

// Test search within scoped results
func (s *APIScopedPermissionsTestSuite) TestSearchWithinScope() {
	// Editor searches for "Editor One" - should only find Editor One's books
	result, status := s.apiGet("/api/books/search?q=Editor+One", s.editor1Cookie)
	s.Equal(http.StatusOK, status)

	if dataWrapper, ok := result["data"].(map[string]interface{}); ok {
		if data, ok := dataWrapper["data"].([]interface{}); ok {
			s.Equal(2, len(data), "Should find 2 books by Editor One")
			for _, item := range data {
				book := item.(map[string]interface{})
				s.Equal("Editor One", book["author"].(string))
			}
		}
	} else if data, ok := result["data"].([]interface{}); ok {
		s.Equal(2, len(data), "Should find 2 books by Editor One")
		for _, item := range data {
			book := item.(map[string]interface{})
			s.Equal("Editor One", book["author"].(string))
		}
	}

	// Member1 searches for "Editor" - should find nothing (can't see editor books)
	result, status = s.apiGet("/api/books/search?q=Editor", s.member1Cookie)
	s.Equal(http.StatusOK, status)

	if dataWrapper, ok := result["data"].(map[string]interface{}); ok {
		if data, ok := dataWrapper["data"].([]interface{}); ok {
			s.Equal(0, len(data), "Member should not find any editor books")
		}
	} else if data, ok := result["data"].([]interface{}); ok {
		s.Equal(0, len(data), "Member should not find any editor books")
	}
}