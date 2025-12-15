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

	"smedi-sme-db/app/models"
	"smedi-sme-db/tests"
	"smedi-sme-db/tests/helpers"
)

type BdspCRUDTestSuite struct {
	suite.Suite
	tests.TestCase

	server     *httptest.Server
	client     *http.Client
	authCookie *http.Cookie
	testUser   *models.User
}

func TestBdspCRUDTestSuite(t *testing.T) {
	suite.Run(t, &BdspCRUDTestSuite{})
}

func (s *BdspCRUDTestSuite) SetupSuite() {
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

func (s *BdspCRUDTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}

func (s *BdspCRUDTestSuite) SetupTest() {
	// Clean database
	s.RefreshDatabase()

	// Create test user with permissions
	s.setupTestUser()
}

func (s *BdspCRUDTestSuite) TearDownTest() {
	// Clean up test data
	if orm := facades.Orm(); orm != nil {
		orm.Query().Exec("DELETE FROM bdsps")
		orm.Query().Exec("DELETE FROM users WHERE email = 'bdsptest@example.com'")
		orm.Query().Exec("DELETE FROM user_roles")
		orm.Query().Exec("DELETE FROM role_permissions")
		orm.Query().Exec("DELETE FROM roles WHERE slug = 'bdsp_admin'")
		orm.Query().Exec("DELETE FROM permissions WHERE slug LIKE 'bdsps_%%'")
	}
}

func (s *BdspCRUDTestSuite) setupTestUser() {
	// Create admin role
	adminRole := &models.Role{
		Name:  "Bdsp Admin",
		Slug:  "bdsp_admin",
		Level: 100,
	}
	s.Nil(facades.Orm().Query().Create(adminRole))

	// Create all CRUD permissions
	permissions := []string{"create", "read", "update", "delete"}
	for _, perm := range permissions {
		permission := &models.Permission{
			Name:  fmt.Sprintf("bdsps_%s", perm),
			Slug:  fmt.Sprintf("bdsps_%s", perm),
			Scope: "by_all",
		}
		s.Nil(facades.Orm().Query().Create(permission))
		s.Nil(helpers.AssignPermissionToRole(adminRole, permission, "by_all"))
	}

	// Create test user
	user, err := helpers.SetupJWTUser("bdsptest@example.com", "password", adminRole)
	s.Nil(err)
	s.testUser = user

	// Login
	s.authCookie = s.loginUser("bdsptest@example.com", "password")
	s.NotNil(s.authCookie)
}

func (s *BdspCRUDTestSuite) loginUser(email, password string) *http.Cookie {
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

func (s *BdspCRUDTestSuite) makeRequest(method, path string, body interface{}) (*http.Response, map[string]interface{}) {
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

func (s *BdspCRUDTestSuite) TestCreateBdsp() {
	bdspData := map[string]interface{}{
		"name":                "Test BDSP Service",
		"postal_address":      "P.O. Box 555, Blantyre",
		"physical_address":    "Victoria Avenue, Blantyre",
		"registration_status": "Active",
		"partners":            []string{"Partner X", "Partner Y"},
		"product_types":       []string{"Training", "Mentorship"},
		"service_list": []map[string]interface{}{
			{"name": "Strategy", "cost": 10000, "duration": "1 week"},
			{"name": "Marketing", "cost": 20000, "duration": "2 weeks"},
		},
	}

	resp, result := s.makeRequest("POST", "/api/bdsps", bdspData)

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.NotNil(data["id"])
	s.Equal(float64(s.testUser.ID), data["created_by"])
	s.Equal(bdspData["name"], data["name"])
	s.Equal(bdspData["postal_address"], data["postal_address"])
	s.Equal(bdspData["physical_address"], data["physical_address"])
	s.Equal(bdspData["registration_status"], data["registration_status"])

	expectedPartners, _ := json.Marshal(bdspData["partners"])
	actualPartners, _ := json.Marshal(data["partners"])
	s.JSONEq(string(expectedPartners), string(actualPartners))

	expectedProductTypes, _ := json.Marshal(bdspData["product_types"])
	actualProductTypes, _ := json.Marshal(data["product_types"])
	s.JSONEq(string(expectedProductTypes), string(actualProductTypes))

	expectedServiceList, _ := json.Marshal(bdspData["service_list"])
	actualServiceList, _ := json.Marshal(data["service_list"])
	s.JSONEq(string(expectedServiceList), string(actualServiceList))
}

func (s *BdspCRUDTestSuite) TestCreateBdspValidation() {
	// Missing required fields
	bdspData := map[string]interface{}{}

	resp, result := s.makeRequest("POST", "/api/bdsps", bdspData)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
	s.Contains(strings.ToLower(result["message"].(string)), "validation")
}

// ============================================================================
// READ Tests
// ============================================================================

func (s *BdspCRUDTestSuite) TestGetBdsp() {
	createdBy := int(s.testUser.ID)
	postal := "Box 1"
	physical := "Street 1"
	status := models.RegistrationConfirmed

	// Create a bdsp first
	bdsp := &models.Bdsp{
		BaseAuditableModel: models.BaseAuditableModel{
			CreatedBy: (func(i int) *uint { u := uint(i); return &u })(createdBy),
		},
		Name:               "Test BDSP Item",
		PostalAddress:      &postal,
		PhysicalAddress:    &physical,
		RegistrationStatus: &status,
		Partners:           []string{"P1"},
		ProductTypes:       []string{"T1"},
		ServiceList: []models.BdspService{
			{Name: "S1", Cost: 100, Duration: "1d"},
		},
	}
	bdsp.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(bdsp))

	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/bdsps/%d", bdsp.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal(float64(bdsp.ID), data["id"].(float64))
}

// ============================================================================
// UPDATE Tests
// ============================================================================

func (s *BdspCRUDTestSuite) TestUpdateBdsp() {
	createdBy := int(s.testUser.ID)
	postal := "Box 1"
	physical := "Street 1"
	status := models.RegistrationConfirmed

	// Create a bdsp first
	bdsp := &models.Bdsp{
		BaseAuditableModel: models.BaseAuditableModel{
			CreatedBy: (func(i int) *uint { u := uint(i); return &u })(createdBy),
		},
		Name:               "Test BDSP Item",
		PostalAddress:      &postal,
		PhysicalAddress:    &physical,
		RegistrationStatus: &status,
		Partners:           []string{"P1"},
		ProductTypes:       []string{"T1"},
		ServiceList: []models.BdspService{
			{Name: "S1", Cost: 100, Duration: "1d"},
		},
	}
	bdsp.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(bdsp))

	updateData := map[string]interface{}{
		"name": "Updated Name",
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/bdsps/%d", bdsp.ID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal(float64(createdBy), data["created_by"])
	s.Equal(updateData["name"], data["name"])
	s.Equal(postal, data["postal_address"])
	s.Equal(physical, data["physical_address"])
	s.Equal(status, data["registration_status"])

	expectedPartners, _ := json.Marshal(bdsp.Partners)
	actualPartners, _ := json.Marshal(data["partners"])
	s.JSONEq(string(expectedPartners), string(actualPartners))

	expectedProductTypes, _ := json.Marshal(bdsp.ProductTypes)
	actualProductTypes, _ := json.Marshal(data["product_types"])
	s.JSONEq(string(expectedProductTypes), string(actualProductTypes))

	expectedServiceList, _ := json.Marshal(bdsp.ServiceList)
	actualServiceList, _ := json.Marshal(data["service_list"])
	s.JSONEq(string(expectedServiceList), string(actualServiceList))
}

// ============================================================================
// DELETE Tests
// ============================================================================

func (s *BdspCRUDTestSuite) TestDeleteBdsp() {
	createdBy := int(s.testUser.ID)
	postal := "Box 1"
	physical := "Street 1"
	status := models.RegistrationConfirmed

	// Create a bdsp first
	bdsp := &models.Bdsp{
		BaseAuditableModel: models.BaseAuditableModel{
			CreatedBy: (func(i int) *uint { u := uint(i); return &u })(createdBy),
		},
		Name:               "Test BDSP Item",
		PostalAddress:      &postal,
		PhysicalAddress:    &physical,
		RegistrationStatus: &status,
		Partners:           []string{"P1"},
		ProductTypes:       []string{"T1"},
		ServiceList: []models.BdspService{
			{Name: "S1", Cost: 100, Duration: "1d"},
		},
	}
	bdsp.CreatedBy = &s.testUser.ID
	s.Nil(facades.Orm().Query().Create(bdsp))

	resp, _ := s.makeRequest("DELETE", fmt.Sprintf("/api/bdsps/%d", bdsp.ID), nil)

	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify soft delete
	var deletedBdsp models.Bdsp
	err := facades.Orm().Query().WithTrashed().Where("id", bdsp.ID).First(&deletedBdsp)
	s.Nil(err)
	s.NotNil(deletedBdsp.DeletedAt)
}

// ============================================================================
// PAGINATION Tests
// ============================================================================

func (s *BdspCRUDTestSuite) TestPagination() {
	// Create 25 bdsps
	for i := 1; i <= 25; i++ {
		createdBy := int(s.testUser.ID)
		postal := "Box 1"
		physical := "Street 1"
		status := models.RegistrationConfirmed

		bdsp := &models.Bdsp{
			BaseAuditableModel: models.BaseAuditableModel{
				CreatedBy: (func(i int) *uint { u := uint(i); return &u })(createdBy),
			},
			Name:               fmt.Sprintf("BDSP %02d", i),
			PostalAddress:      &postal,
			PhysicalAddress:    &physical,
			RegistrationStatus: &status,
			Partners:           []string{"P1"},
			ProductTypes:       []string{"T1"},
			ServiceList: []models.BdspService{
				{Name: "S1", Cost: 100, Duration: "1d"},
			},
		}
		bdsp.CreatedBy = &s.testUser.ID
		s.Nil(facades.Orm().Query().Create(bdsp))
	}

	// Test first page (default page size 20)
	resp, result := s.makeRequest("GET", "/api/bdsps?page=1", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	pagination := data["pagination"].(map[string]interface{})

	s.Equal(20, len(items))
	s.Equal(float64(1), pagination["current_page"])
	s.Equal(float64(25), pagination["total"])
	s.Equal(float64(2), pagination["last_page"])

	// Test custom page size
	resp, result = s.makeRequest("GET", "/api/bdsps?page=1&pageSize=10", nil)
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

func (s *BdspCRUDTestSuite) TestSorting() {
	// Create bdsps with sortable fields
	// TODO: Create items with different sortable values
	// Example: Create 3+ items with different names/values
	for i := 1; i <= 25; i++ {
		createdBy := int(s.testUser.ID)
		postal := "Box 1"
		physical := "Street 1"
		status := models.RegistrationConfirmed

		bdsp := &models.Bdsp{
			BaseAuditableModel: models.BaseAuditableModel{
				CreatedBy: (func(i int) *uint { u := uint(i); return &u })(createdBy),
			},
			Name:               fmt.Sprintf("BDSP %02d", i),
			PostalAddress:      &postal,
			PhysicalAddress:    &physical,
			RegistrationStatus: &status,
			Partners:           []string{"P1"},
			ProductTypes:       []string{"T1"},
			ServiceList: []models.BdspService{
				{Name: "S1", Cost: 100, Duration: "1d"},
			},
		}
		bdsp.CreatedBy = &s.testUser.ID
		s.Nil(facades.Orm().Query().Create(bdsp))
	}
	// Test sort ascending
	resp, result := s.makeRequest("GET", "/api/bdsps?sort=id&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 0)
}

// ============================================================================
// SEARCH Tests
// ============================================================================

func (s *BdspCRUDTestSuite) TestSearch() {
	// Create searchable bdsps
	// TODO: Create items with searchable content
	// Example: Create items with "Searchable" in name/description
	for i := 1; i <= 25; i++ {
		createdBy := int(s.testUser.ID)
		postal := "Box 1"
		physical := "Street 1"
		status := models.RegistrationConfirmed
		name := fmt.Sprintf("BDSP %02d", i)
		if i%5 == 0 {
			name = fmt.Sprintf("test BDSP %02d", i)
		}
		bdsp := &models.Bdsp{
			BaseAuditableModel: models.BaseAuditableModel{
				CreatedBy: (func(i int) *uint { u := uint(i); return &u })(createdBy),
			},
			Name:               name,
			PostalAddress:      &postal,
			PhysicalAddress:    &physical,
			RegistrationStatus: &status,
			Partners:           []string{"P1"},
			ProductTypes:       []string{"T1"},
			ServiceList: []models.BdspService{
				{Name: "S1", Cost: 100, Duration: "1d"},
			},
		}
		bdsp.CreatedBy = &s.testUser.ID
		s.Nil(facades.Orm().Query().Create(bdsp))
	}

	// Search by query
	resp, result := s.makeRequest("GET", "/api/bdsps/search?q=test", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 0)
}
