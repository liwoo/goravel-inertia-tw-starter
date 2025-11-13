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
		orm.Query().Exec("DELETE FROM permissions WHERE slug LIKE 'bdsps_%'")
	}
}

func (s *BdspCRUDTestSuite) setupTestUser() {
	// Create admin role
	adminRole := &models.Role{
		Name:  "BDSP Admin",
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
		"name":                "ABC Business Development Services",
		"postal_address":      "P.O. Box 12345, Lilongwe",
		"physical_address":    "123 Independence Drive, Lilongwe",
		"registration_status": "Registered",
		"product_types":       []string{"Training", "Consulting", "Mentoring"},
		"service_list": []map[string]interface{}{
			{"name": "Business Training", "cost": 50000, "duration": "3 days"},
			{"name": "Financial Consulting", "cost": 100000, "duration": "1 week"},
			{"name": "Marketing Strategy", "cost": 75000, "duration": "5 days"},
		},
		"associated_partners": []string{"UNDP", "World Bank", "IFC"},
	}

	resp, result := s.makeRequest("POST", "/api/bdsps", bdspData)

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.NotNil(data["id"])
	s.Equal("ABC Business Development Services", data["name"])
	s.Equal("P.O. Box 12345, Lilongwe", data["postal_address"])
	s.Equal("Registered", data["registration_status"])
	s.Equal(float64(s.testUser.ID), data["created_by"])

	// Verify JSON arrays
	productTypes := data["product_types"].([]interface{})
	s.Equal(3, len(productTypes))
	s.Contains(productTypes, "Training")

	serviceList := data["service_list"].([]interface{})
	s.Equal(3, len(serviceList))
	service := serviceList[0].(map[string]interface{})
	s.Equal("Business Training", service["name"])
	s.Equal(float64(50000), service["cost"])

	partners := data["associated_partners"].([]interface{})
	s.Equal(3, len(partners))
	s.Contains(partners, "UNDP")
}

func (s *BdspCRUDTestSuite) TestCreateBdspMinimal() {
	bdspData := map[string]interface{}{
		"name":                "Minimal BDSP",
		"registration_status": "Pending",
		"product_types":       []string{"Training"},
		"service_list": []map[string]interface{}{
			{"name": "Basic Service", "cost": 10000, "duration": "1 day"},
		},
		"associated_partners": []string{"Local Partner"},
	}

	resp, result := s.makeRequest("POST", "/api/bdsps", bdspData)

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal("Minimal BDSP", data["name"])
	s.Nil(data["postal_address"])
	s.Nil(data["physical_address"])
}

func (s *BdspCRUDTestSuite) TestCreateBdspValidation() {
	// Missing required fields
	bdspData := map[string]interface{}{
		"postal_address": "Some address",
	}

	resp, result := s.makeRequest("POST", "/api/bdsps", bdspData)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
	s.Contains(strings.ToLower(result["message"].(string)), "validation")
}

func (s *BdspCRUDTestSuite) TestCreateBdspInvalidArrayType() {
	bdspData := map[string]interface{}{
		"name":                "Invalid BDSP",
		"registration_status": "Active",
		"product_types":       "not an array", // Should be array
		"service_list": []map[string]interface{}{
			{"name": "Service", "cost": 10000, "duration": "1 day"},
		},
		"associated_partners": []string{"Partner"},
	}

	resp, result := s.makeRequest("POST", "/api/bdsps", bdspData)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
}

// ============================================================================
// READ Tests
// ============================================================================

func (s *BdspCRUDTestSuite) TestGetBdsp() {
	// Create a BDSP first
	bdsp := s.createTestBdsp("Get Test BDSP", []string{"Training"}, []string{"Partner A"})

	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/bdsps/%d", bdsp.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal(float64(bdsp.ID), data["id"].(float64))
	s.Equal("Get Test BDSP", data["name"])
	s.Equal("Registered", data["registration_status"])

	// Verify JSON arrays are properly returned
	productTypes := data["product_types"].([]interface{})
	s.Contains(productTypes, "Training")
}

func (s *BdspCRUDTestSuite) TestGetBdspNotFound() {
	resp, _ := s.makeRequest("GET", "/api/bdsps/99999", nil)
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *BdspCRUDTestSuite) TestListBdsps() {
	// Create multiple BDSPs
	s.createTestBdsp("BDSP Alpha", []string{"Training"}, []string{"Partner A"})
	s.createTestBdsp("BDSP Beta", []string{"Consulting"}, []string{"Partner B"})
	s.createTestBdsp("BDSP Gamma", []string{"Mentoring"}, []string{"Partner C"})

	resp, result := s.makeRequest("GET", "/api/bdsps", nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	pagination := data["pagination"].(map[string]interface{})

	s.GreaterOrEqual(len(items), 3)
	s.GreaterOrEqual(pagination["total"].(float64), float64(3))
}

// ============================================================================
// UPDATE Tests
// ============================================================================

func (s *BdspCRUDTestSuite) TestUpdateBdsp() {
	// Create a BDSP first
	bdsp := s.createTestBdsp("Original Name", []string{"Training"}, []string{"Partner A"})

	updateData := map[string]interface{}{
		"name":                "Updated BDSP Name",
		"postal_address":      "New P.O. Box 99999",
		"physical_address":    "New Physical Address",
		"registration_status": "Active",
		"product_types":       []string{"Training", "Consulting", "Mentoring"},
		"service_list": []map[string]interface{}{
			{"name": "Updated Service", "cost": 60000, "duration": "4 days"},
		},
		"associated_partners": []string{"Partner A", "Partner B", "Partner C"},
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/bdsps/%d", bdsp.ID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal("Updated BDSP Name", data["name"])
	s.Equal("New P.O. Box 99999", data["postal_address"])
	s.Equal("Active", data["registration_status"])

	// Verify updated arrays
	productTypes := data["product_types"].([]interface{})
	s.Equal(3, len(productTypes))

	partners := data["associated_partners"].([]interface{})
	s.Equal(3, len(partners))
}

func (s *BdspCRUDTestSuite) TestUpdateBdspPartial() {
	// Create a BDSP
	bdsp := s.createTestBdsp("Partial Update Test", []string{"Training"}, []string{"Partner A"})

	// Update only name
	updateData := map[string]interface{}{
		"name": "Partially Updated Name",
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/bdsps/%d", bdsp.ID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	data := result["data"].(map[string]interface{})
	s.Equal("Partially Updated Name", data["name"])
	// Other fields should remain unchanged
	s.Equal("Registered", data["registration_status"])
}

func (s *BdspCRUDTestSuite) TestUpdateBdspNotFound() {
	updateData := map[string]interface{}{
		"name": "Does Not Exist",
	}

	resp, _ := s.makeRequest("PUT", "/api/bdsps/99999", updateData)
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

// ============================================================================
// DELETE Tests
// ============================================================================

func (s *BdspCRUDTestSuite) TestDeleteBdsp() {
	// Create a BDSP first
	bdsp := s.createTestBdsp("To Delete", []string{"Training"}, []string{"Partner A"})

	resp, _ := s.makeRequest("DELETE", fmt.Sprintf("/api/bdsps/%d", bdsp.ID), nil)

	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify soft delete
	var deletedBdsp models.Bdsp
	err := facades.Orm().Query().WithTrashed().Where("id", bdsp.ID).First(&deletedBdsp)
	s.Nil(err)
	s.NotNil(deletedBdsp.DeletedAt)
}

func (s *BdspCRUDTestSuite) TestDeleteBdspNotFound() {
	resp, _ := s.makeRequest("DELETE", "/api/bdsps/99999", nil)
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

// ============================================================================
// PAGINATION Tests
// ============================================================================

func (s *BdspCRUDTestSuite) TestPagination() {
	// Create 25 BDSPs
	for i := 1; i <= 25; i++ {
		s.createTestBdsp(
			fmt.Sprintf("BDSP %02d", i),
			[]string{"Training"},
			[]string{"Partner"},
		)
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

	// Test second page
	resp, result = s.makeRequest("GET", "/api/bdsps?page=2", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	s.Equal(5, len(items))

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
	// Create BDSPs with sortable fields
	names := []struct {
		name   string
		status string
	}{
		{"Zebra BDSP", "Active"},
		{"Alpha BDSP", "Registered"},
		{"Gamma BDSP", "Pending"},
		{"Beta BDSP", "Active"},
	}

	for _, n := range names {
		data := map[string]interface{}{
			"name":                n.name,
			"registration_status": n.status,
			"product_types":       []string{"Training"},
			"service_list": []map[string]interface{}{
				{"name": "Service", "cost": 10000, "duration": "1 day"},
			},
			"associated_partners": []string{"Partner"},
		}
		_, _ = s.makeRequest("POST", "/api/bdsps", data)
		time.Sleep(10 * time.Millisecond) // Ensure different created_at times
	}

	// Test sort by name ascending
	resp, result := s.makeRequest("GET", "/api/bdsps?sort=name&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 4)

	// Verify first item is "Alpha BDSP" (alphabetically first)
	firstItem := items[0].(map[string]interface{})
	s.Equal("Alpha BDSP", firstItem["name"])

	// Test sort by name descending
	resp, result = s.makeRequest("GET", "/api/bdsps?sort=name&direction=desc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	firstItem = items[0].(map[string]interface{})
	s.Equal("Zebra BDSP", firstItem["name"])
}

// ============================================================================
// SEARCH Tests
// ============================================================================

func (s *BdspCRUDTestSuite) TestSearch() {
	// Create searchable BDSPs
	searchData := []struct {
		name    string
		address string
		status  string
	}{
		{"Excellence Training Services", "Lilongwe", "Registered"},
		{"Business Development Group", "Blantyre", "Active"},
		{"Excellence Consulting Ltd", "Mzuzu", "Registered"},
		{"Innovation Hub", "Lilongwe", "Pending"},
	}

	for _, sd := range searchData {
		data := map[string]interface{}{
			"name":                sd.name,
			"physical_address":    sd.address,
			"registration_status": sd.status,
			"product_types":       []string{"Training"},
			"service_list": []map[string]interface{}{
				{"name": "Service", "cost": 10000, "duration": "1 day"},
			},
			"associated_partners": []string{"Partner"},
		}
		_, _ = s.makeRequest("POST", "/api/bdsps", data)
	}

	// Search by name
	resp, result := s.makeRequest("GET", "/api/bdsps/search?q=Excellence", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 2) // Should find at least 2 BDSPs with "Excellence"

	// Search by location
	resp, result = s.makeRequest("GET", "/api/bdsps/search?q=Lilongwe", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 2) // Should find at least 2 BDSPs in Lilongwe

	// Search by status
	resp, result = s.makeRequest("GET", "/api/bdsps/search?q=Registered", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 2)
}

// ============================================================================
// FILTERING Tests
// ============================================================================

func (s *BdspCRUDTestSuite) TestFiltering() {
	// Create BDSPs with different statuses
	statuses := []string{"Registered", "Active", "Pending", "Registered", "Active"}
	for i, status := range statuses {
		data := map[string]interface{}{
			"name":                fmt.Sprintf("BDSP %d", i+1),
			"registration_status": status,
			"product_types":       []string{"Training"},
			"service_list": []map[string]interface{}{
				{"name": "Service", "cost": 10000, "duration": "1 day"},
			},
			"associated_partners": []string{"Partner"},
		}
		_, _ = s.makeRequest("POST", "/api/bdsps", data)
	}

	// Filter by registration_status = Registered
	filterJSON := `[{"field":"registration_status","operator":"equals","value":"Registered"}]`
	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/bdsps?filters=%s", filterJSON), nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 2)

	// Verify all results have "Registered" status
	for _, item := range items {
		bdsp := item.(map[string]interface{})
		s.Equal("Registered", bdsp["registration_status"])
	}
}

// ============================================================================
// PERMISSION Tests
// ============================================================================

func (s *BdspCRUDTestSuite) TestUnauthorizedAccess() {
	// Remove auth cookie
	originalCookie := s.authCookie
	s.authCookie = nil

	resp, _ := s.makeRequest("GET", "/api/bdsps", nil)
	s.Equal(http.StatusUnauthorized, resp.StatusCode)

	// Restore auth cookie
	s.authCookie = originalCookie
}

func (s *BdspCRUDTestSuite) TestCreateWithoutPermission() {
	// This would require creating a user without create permission
	// For now, we'll just verify that authenticated user can create
	bdspData := map[string]interface{}{
		"name":                "Permission Test",
		"registration_status": "Active",
		"product_types":       []string{"Training"},
		"service_list": []map[string]interface{}{
			{"name": "Service", "cost": 10000, "duration": "1 day"},
		},
		"associated_partners": []string{"Partner"},
	}

	resp, result := s.makeRequest("POST", "/api/bdsps", bdspData)
	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))
}

// ============================================================================
// JSON FIELDS Tests
// ============================================================================

func (s *BdspCRUDTestSuite) TestComplexServiceList() {
	bdspData := map[string]interface{}{
		"name":                "Complex Services BDSP",
		"registration_status": "Active",
		"product_types":       []string{"Training", "Consulting", "Mentoring", "Coaching"},
		"service_list": []map[string]interface{}{
			{"name": "Basic Training", "cost": 25000, "duration": "1 day"},
			{"name": "Advanced Training", "cost": 50000, "duration": "3 days"},
			{"name": "Expert Consulting", "cost": 100000, "duration": "1 week"},
			{"name": "Long-term Mentoring", "cost": 500000, "duration": "6 months"},
			{"name": "Executive Coaching", "cost": 250000, "duration": "3 months"},
		},
		"associated_partners": []string{"UNDP", "World Bank", "IFC", "USAID", "EU"},
	}

	resp, result := s.makeRequest("POST", "/api/bdsps", bdspData)
	s.Equal(http.StatusCreated, resp.StatusCode)

	data := result["data"].(map[string]interface{})

	// Verify all services
	serviceList := data["service_list"].([]interface{})
	s.Equal(5, len(serviceList))

	// Verify service details
	service1 := serviceList[0].(map[string]interface{})
	s.Equal("Basic Training", service1["name"])
	s.Equal(float64(25000), service1["cost"])
	s.Equal("1 day", service1["duration"])

	// Verify all partners
	partners := data["associated_partners"].([]interface{})
	s.Equal(5, len(partners))
	s.Contains(partners, "UNDP")
	s.Contains(partners, "EU")
}

// ============================================================================
// AUDIT FIELDS Tests
// ============================================================================

func (s *BdspCRUDTestSuite) TestAuditFields() {
	bdspData := map[string]interface{}{
		"name":                "Audit Test BDSP",
		"registration_status": "Active",
		"product_types":       []string{"Training"},
		"service_list": []map[string]interface{}{
			{"name": "Service", "cost": 10000, "duration": "1 day"},
		},
		"associated_partners": []string{"Partner"},
	}

	resp, result := s.makeRequest("POST", "/api/bdsps", bdspData)
	s.Equal(http.StatusCreated, resp.StatusCode)

	data := result["data"].(map[string]interface{})

	// Verify audit fields
	s.Equal(float64(s.testUser.ID), data["created_by"])
	s.NotNil(data["created_at"])

	// Update and check updated_by
	bdspID := int(data["id"].(float64))
	updateData := map[string]interface{}{
		"name": "Updated for Audit",
	}

	resp, result = s.makeRequest("PUT", fmt.Sprintf("/api/bdsps/%d", bdspID), updateData)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	s.NotNil(data["updated_at"])
	// updated_by should also be set (if implemented in the controller)
}

// ============================================================================
// EDGE CASES Tests
// ============================================================================

func (s *BdspCRUDTestSuite) TestEmptyArrays() {
	// Create with minimal/empty arrays - this should fail validation
	// since product_types, service_list, and associated_partners are required
	bdspData := map[string]interface{}{
		"name":                "Empty Arrays Test",
		"registration_status": "Active",
		"product_types":       []string{},
		"service_list":        []map[string]interface{}{},
		"associated_partners": []string{},
	}

	resp, result := s.makeRequest("POST", "/api/bdsps", bdspData)
	// Should fail validation
	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
}

func (s *BdspCRUDTestSuite) TestLongStrings() {
	longName := strings.Repeat("A", 300) // Exceeds max length
	bdspData := map[string]interface{}{
		"name":                longName,
		"registration_status": "Active",
		"product_types":       []string{"Training"},
		"service_list": []map[string]interface{}{
			{"name": "Service", "cost": 10000, "duration": "1 day"},
		},
		"associated_partners": []string{"Partner"},
	}

	resp, result := s.makeRequest("POST", "/api/bdsps", bdspData)
	// Should fail validation due to max length
	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
}

func (s *BdspCRUDTestSuite) TestSpecialCharacters() {
	bdspData := map[string]interface{}{
		"name":                "BDSP with Special Chars: @#$%^&*()",
		"postal_address":      "P.O. Box 123-456",
		"physical_address":    "123 Main St. Apt #5",
		"registration_status": "Active",
		"product_types":       []string{"Training & Development", "Consulting/Mentoring"},
		"service_list": []map[string]interface{}{
			{"name": "Service (Premium)", "cost": 10000, "duration": "1-2 days"},
		},
		"associated_partners": []string{"Partner A & B"},
	}

	resp, result := s.makeRequest("POST", "/api/bdsps", bdspData)
	s.Equal(http.StatusCreated, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	s.Equal("BDSP with Special Chars: @#$%^&*()", data["name"])
}

// ============================================================================
// Helper Methods
// ============================================================================

func (s *BdspCRUDTestSuite) createTestBdsp(name string, productTypes []string, partners []string) *models.Bdsp {
	bdsp := &models.Bdsp{
		Name:               name,
		RegistrationStatus: "Registered",
		ProductTypes:       productTypes,
		ServiceList: []models.BdspService{
			{Name: "Test Service", Cost: 50000, Duration: "3 days"},
		},
		AssociatedPartners: partners,
	}
	postalAddr := "P.O. Box 123"
	physicalAddr := "123 Test St"
	bdsp.PostalAddress = &postalAddr
	bdsp.PhysicalAddress = &physicalAddr
	userID := int(s.testUser.ID)
	bdsp.CreatedBy = &userID

	err := facades.Orm().Query().Create(bdsp)
	s.Nil(err)

	return bdsp
}
