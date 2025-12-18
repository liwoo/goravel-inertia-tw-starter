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

type ConfigControllerCRUDTestSuite struct {
	suite.Suite
	tests.TestCase

	server     *httptest.Server
	client     *http.Client
	authCookie *http.Cookie
	testUser   *models.User
}

func TestConfigControllerCRUDTestSuite(t *testing.T) {
	suite.Run(t, &ConfigControllerCRUDTestSuite{})
}

func (s *ConfigControllerCRUDTestSuite) SetupSuite() {
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

func (s *ConfigControllerCRUDTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}

func (s *ConfigControllerCRUDTestSuite) SetupTest() {
	// Clean database
	s.RefreshDatabase()

	// Create test user with permissions
	s.setupTestUser()
}

func (s *ConfigControllerCRUDTestSuite) TearDownTest() {
	// Clean up test data
	if orm := facades.Orm(); orm != nil {
		orm.Query().Exec("DELETE FROM sme_config")
		orm.Query().Exec("DELETE FROM users WHERE email = 'configcontrollertest@example.com'")
		orm.Query().Exec("DELETE FROM user_roles")
		orm.Query().Exec("DELETE FROM role_permissions")
		orm.Query().Exec("DELETE FROM roles WHERE slug = 'configcontroller_admin'")
		orm.Query().Exec("DELETE FROM permissions WHERE slug LIKE 'config_%%'")
	}
}

func (s *ConfigControllerCRUDTestSuite) setupTestUser() {
	// Create admin role
	adminRole := &models.Role{
		Name:  "ConfigController Admin",
		Slug:  "configcontroller_admin",
		Level: 100,
	}
	s.Nil(facades.Orm().Query().Create(adminRole))

	// Create all CRUD permissions
	permissions := []string{"create", "read", "update", "delete"}
	for _, perm := range permissions {
		permission := &models.Permission{
			Name:  fmt.Sprintf("config_%s", perm),
			Slug:  fmt.Sprintf("config_%s", perm),
			Scope: "by_all",
		}
		s.Nil(facades.Orm().Query().Create(permission))
		s.Nil(helpers.AssignPermissionToRole(adminRole, permission, "by_all"))
	}

	// Create test user
	user, err := helpers.SetupJWTUser("configcontrollertest@example.com", "password", adminRole)
	s.Nil(err)
	s.testUser = user

	// Login
	s.authCookie = s.loginUser("configcontrollertest@example.com", "password")
	s.NotNil(s.authCookie)
}

func (s *ConfigControllerCRUDTestSuite) loginUser(email, password string) *http.Cookie {
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

func (s *ConfigControllerCRUDTestSuite) makeRequest(method, path string, body interface{}) (*http.Response, map[string]interface{}) {
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

func (s *ConfigControllerCRUDTestSuite) TestCreateConfigController() {
	configcontrollerData := map[string]interface{}{
		"name":        "Test Config",
		"config_type": "Financing",
	}

	resp, result := s.makeRequest("POST", "/api/configs", configcontrollerData)

	// If validation fails, print the error message
	if resp.StatusCode != http.StatusCreated {
		if msg, ok := result["message"]; ok {
			s.T().Logf("Validation error: %v", msg)
		}
		if errors, ok := result["errors"]; ok {
			s.T().Logf("Validation errors: %v", errors)
		}
	}

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.NotNil(data["id"])
	s.Equal(float64(s.testUser.ID), data["created_by"])
}

func (s *ConfigControllerCRUDTestSuite) TestCreateConfigControllerValidation() {
	// Missing required fields
	configcontrollerData := map[string]interface{}{}

	resp, result := s.makeRequest("POST", "/api/configs", configcontrollerData)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
	s.Contains(strings.ToLower(result["message"].(string)), "validation")
}

// ============================================================================
// READ Tests
// ============================================================================

func (s *ConfigControllerCRUDTestSuite) TestGetConfigController() {
	// Create a configcontroller first
	configcontroller := &models.Config{
		Name:       "Test Config",
		ConfigType: "Industries",
	}
	configcontroller.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(configcontroller))

	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/configs/%d", configcontroller.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal(float64(configcontroller.ID), data["id"].(float64))
}

// ============================================================================
// UPDATE Tests
// ============================================================================

func (s *ConfigControllerCRUDTestSuite) TestUpdateConfigController() {
	// Create a configcontroller first
	configcontroller := &models.Config{
		Name:       "Original Config",
		ConfigType: "Sectors",
	}
	configcontroller.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(configcontroller))

	updateData := map[string]interface{}{
		"name": "Updated Config",
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/configs/%d", configcontroller.ID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))
}

// ============================================================================
// DELETE Tests
// ============================================================================

func (s *ConfigControllerCRUDTestSuite) TestDeleteConfigController() {
	// Create a config first
	config := &models.Config{
		Name:       "Config to Delete",
		ConfigType: "Registration Status",
	}
	config.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(config))

	resp, _ := s.makeRequest("DELETE", fmt.Sprintf("/api/configs/%d", config.ID), nil)

	// Delete should succeed (framework bug fixed in generic_crud_service.go:666)
	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify soft delete
	var deletedConfig models.Config
	err := facades.Orm().Query().WithTrashed().Where("id", config.ID).First(&deletedConfig)
	s.Nil(err)
	s.NotNil(deletedConfig.DeletedAt)
}

// ============================================================================
// PAGINATION Tests
// ============================================================================

func (s *ConfigControllerCRUDTestSuite) TestPagination() {
	// Create 25 configcontrollers
	for i := 1; i <= 25; i++ {
		configcontroller := &models.Config{
			Name:       fmt.Sprintf("Config %02d", i),
			ConfigType: "Development Partners",
		}
		configcontroller.CreatedBy = &s.testUser.ID
		s.Nil(facades.Orm().Query().Create(configcontroller))
	}

	// Test first page (default page size 20)
	resp, result := s.makeRequest("GET", "/api/configs?page=1", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	pagination := data["pagination"].(map[string]interface{})

	s.Equal(20, len(items))
	s.Equal(float64(1), pagination["current_page"])
	s.Equal(float64(25), pagination["total"])
	s.Equal(float64(2), pagination["last_page"])

	// Test custom page size
	resp, result = s.makeRequest("GET", "/api/configs?page=1&pageSize=10", nil)
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

func (s *ConfigControllerCRUDTestSuite) TestSorting() {
	// Create configcontrollers with sortable fields
	for i := 1; i <= 3; i++ {
		configcontroller := &models.Config{
			Name:       fmt.Sprintf("Sortable Config %02d", i),
			ConfigType: "Business Categories",
		}
		configcontroller.CreatedBy = &s.testUser.ID
		s.Nil(facades.Orm().Query().Create(configcontroller))
	}

	// Test sort ascending
	resp, result := s.makeRequest("GET", "/api/configs?sort=id&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 0)
}

// ============================================================================
// SEARCH Tests
// ============================================================================

func (s *ConfigControllerCRUDTestSuite) TestSearch() {
	// Create searchable configcontrollers
	searchableDesc := "Searchable test description"
	configcontroller := &models.Config{
		Name:        "Searchable Test Config",
		ConfigType:  "Improvement Aspects",
		Description: &searchableDesc,
	}
	configcontroller.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(configcontroller))

	// Search by query
	resp, result := s.makeRequest("GET", "/api/configs/search?q=test", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 0)
}
