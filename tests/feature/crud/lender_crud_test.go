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

	"starter-project/app/models"
	"starter-project/tests"
	"starter-project/tests/helpers"
)

type LenderCRUDTestSuite struct {
	suite.Suite
	tests.TestCase

	server     *httptest.Server
	client     *http.Client
	authCookie *http.Cookie
	testUser   *models.User
}

func TestLenderCRUDTestSuite(t *testing.T) {
	suite.Run(t, &LenderCRUDTestSuite{})
}

func (s *LenderCRUDTestSuite) SetupSuite() {
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

func (s *LenderCRUDTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}

func (s *LenderCRUDTestSuite) SetupTest() {
	// Clean database
	s.RefreshDatabase()

	// Create test user with permissions
	s.setupTestUser()
}

func (s *LenderCRUDTestSuite) TearDownTest() {
	// Clean up test data
	if orm := facades.Orm(); orm != nil {
		orm.Query().Exec("DELETE FROM lenders")
		orm.Query().Exec("DELETE FROM users WHERE email = 'lendertest@example.com'")
		orm.Query().Exec("DELETE FROM user_roles")
		orm.Query().Exec("DELETE FROM role_permissions")
		orm.Query().Exec("DELETE FROM roles WHERE slug = 'lender_admin'")
		orm.Query().Exec("DELETE FROM permissions WHERE slug LIKE 'lenders_%'")
	}
}

func (s *LenderCRUDTestSuite) setupTestUser() {
	// Create admin role
	adminRole := &models.Role{
		Name:  "Lender Admin",
		Slug:  "lender_admin",
		Level: 100,
	}
	s.Nil(facades.Orm().Query().Create(adminRole))

	// Create all CRUD permissions
	permissions := []string{"create", "read", "update", "delete"}
	for _, perm := range permissions {
		permission := &models.Permission{
			Name:  fmt.Sprintf("lenders_%s", perm),
			Slug:  fmt.Sprintf("lenders_%s", perm),
			Scope: "by_all",
		}
		s.Nil(facades.Orm().Query().Create(permission))
		s.Nil(helpers.AssignPermissionToRole(adminRole, permission, "by_all"))
	}

	// Create test user
	user, err := helpers.SetupJWTUser("lendertest@example.com", "password", adminRole)
	s.Nil(err)
	s.testUser = user

	// Login
	s.authCookie = s.loginUser("lendertest@example.com", "password")
	s.NotNil(s.authCookie)
}

func (s *LenderCRUDTestSuite) loginUser(email, password string) *http.Cookie {
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

func (s *LenderCRUDTestSuite) makeRequest(method, path string, body interface{}) (*http.Response, map[string]interface{}) {
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

// ============================================================================
// CREATE Tests
// ============================================================================

func (s *LenderCRUDTestSuite) TestCreateLender() {
	lenderData := map[string]interface{}{
		// TODO: Add required fields for Lender
		// Example (based on Lender model):
		"name":  "Test Lender",
		"email": "testlender@example.com",
		// "phone": "555-0100",
		// "address": "123 Test St",
		// "gender": "other",
	}

	resp, result := s.makeRequest("POST", "/api/lenders", lenderData)

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.NotNil(data["id"])
	s.Equal(float64(s.testUser.ID), data["created_by"])
}

func (s *LenderCRUDTestSuite) TestCreateLenderValidation() {
	// Missing required fields
	lenderData := map[string]interface{}{}

	resp, result := s.makeRequest("POST", "/api/lenders", lenderData)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
	s.Contains(strings.ToLower(result["message"].(string)), "validation")
}

// ============================================================================
// READ Tests
// ============================================================================

func (s *LenderCRUDTestSuite) TestGetLender() {
	// Create a lender first
	lender := &models.Lender{
		Name:  "Get Test Lender",
		Email: "gettest@example.com",
		// TODO: Add other required fields if any
	}
	lender.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(lender))

	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/lenders/%d", lender.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal(float64(lender.ID), data["id"].(float64))
	s.Equal("Get Test Lender", data["name"])
}

// ============================================================================
// UPDATE Tests
// ============================================================================

func (s *LenderCRUDTestSuite) TestUpdateLender() {
	// Create a lender first
	lender := &models.Lender{
		Name:  "Original Name",
		Email: "original@example.com",
	}
	lender.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(lender))

	updateData := map[string]interface{}{
		"name":  "Updated Name",
		"email": "updated@example.com",
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/lenders/%d", lender.ID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal("Updated Name", data["name"])
	s.Equal("updated@example.com", data["email"])
}

// ============================================================================
// DELETE Tests
// ============================================================================

func (s *LenderCRUDTestSuite) TestDeleteLender() {
	// Create a lender first
	lender := &models.Lender{
		Name:  "To Delete",
		Email: "delete@example.com",
	}
	lender.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(lender))

	resp, _ := s.makeRequest("DELETE", fmt.Sprintf("/api/lenders/%d", lender.ID), nil)

	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify soft delete
	var deletedLender models.Lender
	err := facades.Orm().Query().WithTrashed().Where("id", lender.ID).First(&deletedLender)
	s.Nil(err)
	s.NotNil(deletedLender.DeletedAt)
}

// ============================================================================
// PAGINATION Tests
// ============================================================================

func (s *LenderCRUDTestSuite) TestPagination() {
	// Create 25 lenders
	for i := 1; i <= 25; i++ {
		lender := &models.Lender{
			Name:  fmt.Sprintf("Lender %02d", i),
			Email: fmt.Sprintf("lender%d@example.com", i),
		}
		lender.CreatedBy = &s.testUser.ID
		s.Nil(facades.Orm().Query().Create(lender))
	}

	// Test first page (default page size 20)
	resp, result := s.makeRequest("GET", "/api/lenders?page=1", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	pagination := data["pagination"].(map[string]interface{})

	s.Equal(20, len(items))
	s.Equal(float64(1), pagination["current_page"])
	s.Equal(float64(25), pagination["total"])
	s.Equal(float64(2), pagination["last_page"])

	// Test custom page size
	resp, result = s.makeRequest("GET", "/api/lenders?page=1&pageSize=10", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	pagination = data["pagination"].(map[string]interface{})

	s.Equal(10, len(items))
	s.Equal(float64(3), pagination["last_page"])
}

// ============================================================================
// SORTING Tests
// ============================================================================

func (s *LenderCRUDTestSuite) TestSorting() {
	// Create lenders with sortable fields
	// Create items with different names for sorting
	lenders := []struct {
		name  string
		email string
	}{
		{"Zebra Lender", "zebra@example.com"},
		{"Alpha Lender", "alpha@example.com"},
		{"Beta Lender", "beta@example.com"},
	}

	for _, l := range lenders {
		lender := &models.Lender{
			Name:  l.name,
			Email: l.email,
		}
		lender.CreatedBy = &s.testUser.ID
		s.Nil(facades.Orm().Query().Create(lender))
		time.Sleep(10 * time.Millisecond) // Ensure different created_at times
	}

	// Test sort by name ascending
	resp, result := s.makeRequest("GET", "/api/lenders?sort=name&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 3)

	// Verify first item is "Alpha Lender" (alphabetically first)
	firstItem := items[0].(map[string]interface{})
	s.Equal("Alpha Lender", firstItem["name"])
}

// ============================================================================
// SEARCH Tests
// ============================================================================

func (s *LenderCRUDTestSuite) TestSearch() {
	// Create searchable lenders
	// Create items with searchable content in name and email
	searchLenders := []struct {
		name  string
		email string
	}{
		{"Searchable Lender One", "searchable1@example.com"},
		{"Another Lender", "another@example.com"},
		{"Searchable Lender Two", "searchable2@example.com"},
	}

	for _, l := range searchLenders {
		lender := &models.Lender{
			Name:  l.name,
			Email: l.email,
		}
		lender.CreatedBy = &s.testUser.ID
		s.Nil(facades.Orm().Query().Create(lender))
	}

	// Search by query - should find lenders with "Searchable" in name
	resp, result := s.makeRequest("GET", "/api/lenders/search?q=Searchable", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 2) // Should find at least 2 lenders with "Searchable"
}
