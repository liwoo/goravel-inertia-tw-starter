package sme

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
	"github.com/goravel/framework/support/carbon"
	"github.com/stretchr/testify/suite"

	"smedi-sme-db/app/models"
	"smedi-sme-db/tests"
	"smedi-sme-db/tests/helpers"
)

type PrimaryBusinessOwnerControllerCRUDTestSuite struct {
	suite.Suite
	tests.TestCase

	server     *httptest.Server
	client     *http.Client
	authCookie *http.Cookie
	testUser   *models.User
}

func TestPrimaryBusinessOwnerControllerCRUDTestSuite(t *testing.T) {
	suite.Run(t, &PrimaryBusinessOwnerControllerCRUDTestSuite{})
}

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) SetupSuite() {
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

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) SetupTest() {
	// Clean any existing test data first (in case previous test failed to clean up)
	if orm := facades.Orm(); orm != nil {
		// Delete child records first to respect foreign key constraints
		orm.Query().Exec("DELETE FROM primary_business_owner")
		orm.Query().Exec("DELETE FROM additional_business_members")
		orm.Query().Exec("DELETE FROM business_formalisation")
		orm.Query().Exec("DELETE FROM business_employee_summary")
		orm.Query().Exec("DELETE FROM smes")
	}

	// Refresh database (runs migrations)
	s.RefreshDatabase()

	// Create test user with permissions
	s.setupTestUser()
}

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TearDownTest() {
	// Clean up test data in the correct order (respecting foreign key constraints)
	if orm := facades.Orm(); orm != nil {
		// First delete child records that reference smes
		orm.Query().Exec("DELETE FROM primary_business_owner")
		orm.Query().Exec("DELETE FROM additional_business_members")
		orm.Query().Exec("DELETE FROM business_formalisation")
		orm.Query().Exec("DELETE FROM business_employee_summary")

		// Now we can safely delete smes
		orm.Query().Exec("DELETE FROM smes")

		// Clean up user-related records
		orm.Query().Exec("DELETE FROM user_roles")
		orm.Query().Exec("DELETE FROM users WHERE email = 'smetest@example.com'")

		// Clean up permissions and roles
		orm.Query().Exec("DELETE FROM role_permissions")
		orm.Query().Exec("DELETE FROM roles WHERE slug = 'sme_admin'")
		orm.Query().Exec("DELETE FROM permissions WHERE slug LIKE 'smes_%' OR slug LIKE 'primary_business_owners_%'")
	}
}

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) setupTestUser() {
	// Create admin role
	adminRole := &models.Role{
		Name:  "SME Admin",
		Slug:  "sme_admin",
		Level: 100,
	}
	s.Nil(facades.Orm().Query().Create(adminRole))

	// Create all CRUD permissions for both SMEs and Primary Business Owners
	permissionTypes := []string{"smes", "primary_business_owners"}

	for _, permType := range permissionTypes {
		for _, action := range []string{"create", "read", "update", "delete", "view"} {
			permission := &models.Permission{
				Name:  fmt.Sprintf("%s_%s", permType, action),
				Slug:  fmt.Sprintf("%s_%s", permType, action),
				Scope: "by_all",
			}
			s.Nil(facades.Orm().Query().Create(permission))
			s.Nil(helpers.AssignPermissionToRole(adminRole, permission, "by_all"))
		}
	}

	// Create test user using the helper function
	user, err := helpers.SetupJWTUser("smetest@example.com", "password", adminRole)
	s.Nil(err)
	s.testUser = user

	// Login to get the auth cookie
	s.authCookie = s.loginUser("smetest@example.com", "password")
	s.NotNil(s.authCookie, "Failed to get auth cookie")
}

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) loginUser(email, password string) *http.Cookie {
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

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) makeRequest(method, path string, body interface{}) (*http.Response, map[string]interface{}) {
	var jsonData []byte
	if body != nil {
		var err error
		jsonData, err = json.Marshal(body)
		if err != nil {
			s.T().Fatalf("Failed to marshal request body: %v", err)
		}
		fmt.Printf("Request %s %s\nBody: %s\n", method, path, string(jsonData))
	}

	req, _ := http.NewRequest(method, s.server.URL+path, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if s.authCookie != nil {
		req.AddCookie(s.authCookie)
	}

	resp, _ := s.client.Do(req)
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response Status: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n", string(bodyBytes))

	var result map[string]interface{}
	if len(bodyBytes) > 0 {
		json.Unmarshal(bodyBytes, &result)
	}

	return resp, result
}

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) getValidPBOData(smeId uint) map[string]interface{} {

	return map[string]interface{}{
		"first_name":               "John",
		"last_name":                "Doe",
		"other_names":              "Middle",
		"nationality":              "Malawian",
		"national_id_number":       "T6N8SARR",
		"date_of_birth":            "1990-01-15",
		"gender":                   "male",
		"education_level":          "tertiary",
		"malawian_status":          "citizen",
		"has_special_needs":        false,
		"phone_number":             "+265991234567",
		"email":                    "john.doe@example.com",
		"physical_address":         "123 Business St",
		"postal_address":           "P.O. Box 123",
		"region":                   "Southern",
		"district":                 "Blantyre",
		"traditional_authority":    "Kunthembwe",
		"alt_contact_name":         "Jane Doe",
		"alt_contact_relationship": "Spouse",
		"alt_contact_phone":        "+265992345678",
		"sme_id":                   smeId,
	}
}

// ============================================================================
// CREATE Tests
// ============================================================================

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TestCreatePrimaryBusinessOwner() {
	// First create an SME
	smeData := map[string]interface{}{
		"name":                         "Test Business",
		"business_category":            "retail",
		"sector":                       "commerce",
		"contact_phone":                "+265991234567",
		"contact_email":                "test@example.com",
		"registration_number":          "BRNR-PBO001",
		"operational_start_date":       carbon.Now().Format("2006-01-02"),
		"business_improvement_aspects": []string{"financial_management"},
		"business_accessed_financing":  []string{"bank_loan"},
	}

	resp, smeResult := s.makeRequest("POST", "/api/smes", smeData)

	bodyBytes, _ := io.ReadAll(resp.Body)

	// Check if the request was successful
	if resp.StatusCode != http.StatusCreated {
		s.T().Fatalf("Failed to create SME. Status: %d, Body: %s", resp.StatusCode, string(bodyBytes))
	}

	s.True(smeResult["success"].(bool), "SME creation was not successful")
	s.NotNil(smeResult["data"], "SME creation response missing data")

	smeDataMap, ok := smeResult["data"].(map[string]interface{})
	if !ok {
		s.T().Fatalf("Failed to parse SME data: %v", smeResult["data"])
	}

	smeID, ok := smeDataMap["id"].(float64)
	if !ok {
		s.T().Fatalf("Failed to get SME ID from response: %v", smeDataMap)
	}

	// Now create PBO
	data := s.getValidPBOData(uint(smeID))
	data["sme_id"] = uint(smeID)

	pboResp, pboResult := s.makeRequest("POST", "/api/primary-business-owners", data)

	s.Equal(http.StatusCreated, pboResp.StatusCode, "Failed to create PBO")
	s.True(pboResult["success"].(bool), "PBO creation was not successful")
	s.NotEmpty(pboResult["data"].(map[string]interface{})["id"])

	// Verify data was saved correctly
	var pbo models.PrimaryBusinessOwner
	s.Nil(facades.Orm().Query().Where("id = ?", pboResult["data"].(map[string]interface{})["id"]).First(&pbo))
	s.Equal(data["first_name"], pbo.FirstName)
	s.Equal(data["last_name"], pbo.LastName)
	s.Equal(data["national_id_number"], pbo.NationalIdNumber)
	s.Equal(uint(smeID), uint(pbo.SmeID))
}

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TestCreatePrimaryBusinessOwnerValidation() {
	// First create an SME
	smeData := map[string]interface{}{
		"name":                         "Test Business for Validation",
		"business_category":            "retail",
		"sector":                       "commerce",
		"contact_phone":                "+265991234567",
		"contact_email":                "validation@example.com",
		"registration_number":          "BRNR-PBO002",
		"operational_start_date":       carbon.Now().Format("2006-01-02"),
		"business_improvement_aspects": []string{"financial_management"},
		"business_accessed_financing":  []string{"bank_loan"},
	}

	resp, smeResult := s.makeRequest("POST", "/api/smes", smeData)
	s.Equal(http.StatusCreated, resp.StatusCode, "Failed to create SME for validation test")
	smeID := uint(smeResult["data"].(map[string]interface{})["id"].(float64))

	// Test required fields
	data := s.getValidPBOData(smeID)
	data["first_name"] = "" // Make first name empty

	resp, result := s.makeRequest("POST", "/api/primary-business-owners", data)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
	s.Contains(result["message"], "Validation failed")
	s.Contains(result["errors"].(map[string]interface{}), "first_name")
}
// ============================================================================
// READ Tests
// ============================================================================

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TestGetPrimaryBusinessOwner() {
	// First create an SME
	smeData := map[string]interface{}{
		"name":                         "Test Business for Read",
		"business_category":            "retail",
		"sector":                       "commerce",
		"contact_phone":                "+265991234567",
		"contact_email":                "read@example.com",
		"registration_number":          "BRNR-PBO003",
		"operational_start_date":       carbon.Now().Format("2006-01-02"),
		"business_improvement_aspects": []string{"financial_management"},
		"business_accessed_financing":  []string{"bank_loan"},
	}

	resp, smeResult := s.makeRequest("POST", "/api/smes", smeData)
	s.Equal(http.StatusCreated, resp.StatusCode, "Failed to create SME for read test")
	smeID := uint(smeResult["data"].(map[string]interface{})["id"].(float64))

	// Create a PBO
	data := s.getValidPBOData(smeID)
	_, createResult := s.makeRequest("POST", "/api/primary-business-owners", data)
	pboID := createResult["data"].(map[string]interface{})["id"]

	// Now try to retrieve it
	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/primary-business-owners/%.0f", pboID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))
	s.Equal(pboID, result["data"].(map[string]interface{})["id"])
	s.Equal(data["first_name"], result["data"].(map[string]interface{})["first_name"])
	s.Equal(data["last_name"], result["data"].(map[string]interface{})["last_name"])
	s.Equal(data["national_id_number"], result["data"].(map[string]interface{})["national_id_number"])
}
// ============================================================================
// UPDATE Tests
// ============================================================================

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TestUpdatePrimaryBusinessOwner() {
	// First create an SME
	smeData := map[string]interface{}{
		"name":                         "Test Business for Update",
		"business_category":            "retail",
		"sector":                       "commerce",
		"contact_phone":                "+265991234567",
		"contact_email":                "update@example.com",
		"registration_number":          "BRNR-PBO004",
		"operational_start_date":       carbon.Now().Format("2006-01-02"),
		"business_improvement_aspects": []string{"financial_management"},
		"business_accessed_financing":  []string{"bank_loan"},
	}

	resp, smeResult := s.makeRequest("POST", "/api/smes", smeData)
	s.Equal(http.StatusCreated, resp.StatusCode, "Failed to create SME for update test")
	smeID := uint(smeResult["data"].(map[string]interface{})["id"].(float64))

	// Create a PBO
	data := s.getValidPBOData(smeID)
	_, createResult := s.makeRequest("POST", "/api/primary-business-owners", data)
	pboID := createResult["data"].(map[string]interface{})["id"]

	// Update the PBO
	updateData := map[string]interface{}{
		"first_name": "UpdatedFirstName",
		"last_name":  "UpdatedLastName",
		"email":      "updated.email@example.com",
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/primary-business-owners/%.0f", pboID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	// Verify the update
	var pbo models.PrimaryBusinessOwner
	s.Nil(facades.Orm().Query().Where("id = ?", pboID).First(&pbo))
	s.Equal(updateData["first_name"], pbo.FirstName)
	s.Equal(updateData["last_name"], pbo.LastName)
	s.Equal(updateData["email"], *pbo.Email)
}
// ============================================================================
// DELETE Tests
// ============================================================================

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TestDeletePrimaryBusinessOwner() {
	// First create an SME
	smeData := map[string]interface{}{
		"name":                         "Test Business for Delete",
		"business_category":            "retail",
		"sector":                       "commerce",
		"contact_phone":                "+265991234567",
		"contact_email":                "delete@example.com",
		"registration_number":          "BRNR-PBO005",
		"operational_start_date":       carbon.Now().Format("2006-01-02"),
		"business_improvement_aspects": []string{"financial_management"},
		"business_accessed_financing":  []string{"bank_loan"},
	}

	resp, smeResult := s.makeRequest("POST", "/api/smes", smeData)
	s.Equal(http.StatusCreated, resp.StatusCode, "Failed to create SME for delete test")
	smeID := uint(smeResult["data"].(map[string]interface{})["id"].(float64))

	// Create a PBO
	data := s.getValidPBOData(smeID)
	_, createResult := s.makeRequest("POST", "/api/primary-business-owners", data)
	pboID := createResult["data"].(map[string]interface{})["id"]

	// Now delete it
	delResp, _ := s.makeRequest("DELETE", fmt.Sprintf("/api/primary-business-owners/%.0f", pboID), nil)

	s.Equal(http.StatusNoContent, delResp.StatusCode)

	// Verify it's deleted (soft delete)
	var pbo models.PrimaryBusinessOwner
	err := facades.Orm().Query().WithTrashed().Where("id = ?", pboID).First(&pbo)
	s.Nil(err, "PBO should still exist in database with soft delete")
	s.NotNil(pbo.DeletedAt, "PBO should have DeletedAt timestamp set")
}
// ============================================================================
// LIST/INDEX Tests
// ============================================================================

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TestListPrimaryBusinessOwners() {
	// First create an SME
	smeData := map[string]interface{}{
		"name":                         "Test Business for List",
		"business_category":            "retail",
		"sector":                       "commerce",
		"contact_phone":                "+265991234567",
		"contact_email":                "list@example.com",
		"registration_number":          "BRNR-PBO006",
		"operational_start_date":       carbon.Now().Format("2006-01-02"),
		"business_improvement_aspects": []string{"financial_management"},
		"business_accessed_financing":  []string{"bank_loan"},
	}

	resp, smeResult := s.makeRequest("POST", "/api/smes", smeData)
	s.Equal(http.StatusCreated, resp.StatusCode, "Failed to create SME for list test")
	smeID := uint(smeResult["data"].(map[string]interface{})["id"].(float64))

	// Create multiple PBOs
	for i := 0; i < 3; i++ {
		data := s.getValidPBOData(smeID)
		data["national_id_number"] = fmt.Sprintf("ABCD00%02d", i)
		data["email"] = fmt.Sprintf("pbo%d@example.com", i)
		data["first_name"] = fmt.Sprintf("FirstName%d", i)
		s.makeRequest("POST", "/api/primary-business-owners", data)
	}

	// Get the list
	resp, result := s.makeRequest("GET", "/api/primary-business-owners", nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	// The list endpoint returns paginated data with nested structure
	dataMap := result["data"].(map[string]interface{})
	items := dataMap["data"].([]interface{})
	s.GreaterOrEqual(len(items), 3, "Should have at least 3 PBOs")
}

// ============================================================================
// SEARCH Tests
// ============================================================================

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TestSearchPrimaryBusinessOwners() {
	// First create an SME
	smeData := map[string]interface{}{
		"name":                         "Test Business for Search",
		"business_category":            "retail",
		"sector":                       "commerce",
		"contact_phone":                "+265991234567",
		"contact_email":                "search@example.com",
		"registration_number":          "BRNR-PBO007",
		"operational_start_date":       carbon.Now().Format("2006-01-02"),
		"business_improvement_aspects": []string{"financial_management"},
		"business_accessed_financing":  []string{"bank_loan"},
	}

	resp, smeResult := s.makeRequest("POST", "/api/smes", smeData)
	s.Equal(http.StatusCreated, resp.StatusCode, "Failed to create SME for search test")
	smeID := uint(smeResult["data"].(map[string]interface{})["id"].(float64))

	// Create test data
	testData := []map[string]interface{}{
		{"first_name": "John", "last_name": "Doe", "national_id_number": "SRCH0001", "email": "john.doe.search@example.com"},
		{"first_name": "Jane", "last_name": "Smith", "national_id_number": "SRCH0002", "email": "jane.smith.search@example.com"},
		{"first_name": "John", "last_name": "Smith", "national_id_number": "SRCH0003", "email": "john.smith.search@example.com"},
	}

	for _, data := range testData {
		pboData := s.getValidPBOData(smeID)
		for k, v := range data {
			pboData[k] = v
		}
		s.makeRequest("POST", "/api/primary-business-owners", pboData)
	}

	// Search by first name
	resp, result := s.makeRequest("GET", "/api/primary-business-owners?search=John", nil)
	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	// Count matching records (paginated response)
	dataMap := result["data"].(map[string]interface{})
	items := dataMap["data"].([]interface{})
	matchingJohns := 0
	for _, item := range items {
		pbo := item.(map[string]interface{})
		if firstName, ok := pbo["first_name"].(string); ok && firstName == "John" {
			matchingJohns++
		}
	}
	s.Equal(2, matchingJohns, "Should find 2 PBOs named John")

	// Search by last name
	resp, result = s.makeRequest("GET", "/api/primary-business-owners?search=Smith", nil)
	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	// Count matching Smiths (paginated response)
	dataMap = result["data"].(map[string]interface{})
	items = dataMap["data"].([]interface{})
	matchingSmiths := 0
	for _, item := range items {
		pbo := item.(map[string]interface{})
		if lastName, ok := pbo["last_name"].(string); ok && lastName == "Smith" {
			matchingSmiths++
		}
	}
	s.Equal(2, matchingSmiths, "Should find 2 PBOs with last name Smith")

	// Search by national ID
	resp, result = s.makeRequest("GET", "/api/primary-business-owners?search=SRCH0002", nil)
	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	// Check for Jane (paginated response)
	dataMap = result["data"].(map[string]interface{})
	items = dataMap["data"].([]interface{})
	foundJane := false
	for _, item := range items {
		pbo := item.(map[string]interface{})
		if nationalId, ok := pbo["national_id_number"].(string); ok && nationalId == "SRCH0002" {
			s.Equal("Jane", pbo["first_name"])
			foundJane = true
			break
		}
	}
	s.True(foundJane, "Should find Jane with national ID SRCH0002")
}

// ============================================================================
// PAGINATION Tests
// ============================================================================

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TestPagination() {
	// First create an SME
	smeData := map[string]interface{}{
		"name":                         "Test Business for Pagination",
		"business_category":            "retail",
		"sector":                       "commerce",
		"contact_phone":                "+265991234567",
		"contact_email":                "pagination@example.com",
		"registration_number":          "BRNR-PAGE01",
		"operational_start_date":       carbon.Now().Format("2006-01-02"),
		"business_improvement_aspects": []string{"financial_management"},
		"business_accessed_financing":  []string{"bank_loan"},
	}

	resp, smeResult := s.makeRequest("POST", "/api/smes", smeData)
	s.Equal(http.StatusCreated, resp.StatusCode, "Failed to create SME for pagination test")
	smeID := uint(smeResult["data"].(map[string]interface{})["id"].(float64))

	// Create 25 PBOs
	for i := 1; i <= 25; i++ {
		data := s.getValidPBOData(smeID)
		data["first_name"] = fmt.Sprintf("FirstName%02d", i)
		data["last_name"] = fmt.Sprintf("LastName%02d", i)
		data["national_id_number"] = fmt.Sprintf("PAGE0%03d", i)
		data["email"] = fmt.Sprintf("pbo%02d@pagination.com", i)
		s.makeRequest("POST", "/api/primary-business-owners", data)
	}

	// Test first page (default page size 20)
	resp, result := s.makeRequest("GET", "/api/primary-business-owners?page=1", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	dataMap := result["data"].(map[string]interface{})
	items := dataMap["data"].([]interface{})
	pagination := dataMap["pagination"].(map[string]interface{})

	s.Equal(20, len(items), "First page should have 20 items")
	s.Equal(float64(1), pagination["current_page"], "Should be on page 1")
	s.Equal(float64(25), pagination["total"], "Should have 25 total items")
	s.Equal(float64(2), pagination["last_page"], "Should have 2 pages")

	// Test second page
	resp, result = s.makeRequest("GET", "/api/primary-business-owners?page=2", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	dataMap = result["data"].(map[string]interface{})
	items = dataMap["data"].([]interface{})
	pagination = dataMap["pagination"].(map[string]interface{})

	s.Equal(5, len(items), "Second page should have 5 remaining items")
	s.Equal(float64(2), pagination["current_page"], "Should be on page 2")

	// Test custom page size
	resp, result = s.makeRequest("GET", "/api/primary-business-owners?page=1&pageSize=10", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	dataMap = result["data"].(map[string]interface{})
	items = dataMap["data"].([]interface{})
	pagination = dataMap["pagination"].(map[string]interface{})

	s.Equal(10, len(items), "Should have 10 items with custom page size")
	s.Equal(float64(3), pagination["last_page"], "Should have 3 pages with page size 10")

	// Test third page with custom page size
	resp, result = s.makeRequest("GET", "/api/primary-business-owners?page=3&pageSize=10", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	dataMap = result["data"].(map[string]interface{})
	items = dataMap["data"].([]interface{})
	pagination = dataMap["pagination"].(map[string]interface{})

	s.Equal(5, len(items), "Third page should have 5 remaining items")
	s.Equal(float64(3), pagination["current_page"], "Should be on page 3")
}

// ============================================================================
// SORTING Tests
// ============================================================================

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TestSorting() {
	// First create an SME
	smeData := map[string]interface{}{
		"name":                         "Test Business for Sorting",
		"business_category":            "retail",
		"sector":                       "commerce",
		"contact_phone":                "+265991234567",
		"contact_email":                "sorting@example.com",
		"registration_number":          "BRNR-SORT01",
		"operational_start_date":       carbon.Now().Format("2006-01-02"),
		"business_improvement_aspects": []string{"financial_management"},
		"business_accessed_financing":  []string{"bank_loan"},
	}

	resp, smeResult := s.makeRequest("POST", "/api/smes", smeData)
	s.Equal(http.StatusCreated, resp.StatusCode, "Failed to create SME for sorting test")
	smeID := uint(smeResult["data"].(map[string]interface{})["id"].(float64))

	// Create PBOs with different names
	names := []struct {
		first string
		last  string
	}{
		{"Zebra", "Young"},
		{"Alpha", "Wilson"},
		{"Beta", "Xavier"},
	}

	for i, name := range names {
		data := s.getValidPBOData(smeID)
		data["first_name"] = name.first
		data["last_name"] = name.last
		data["national_id_number"] = fmt.Sprintf("SORT0%03d", i+1)
		data["email"] = fmt.Sprintf("sort%d@example.com", i)
		s.makeRequest("POST", "/api/primary-business-owners", data)
	}

	// Test sort ascending by first name
	resp, result := s.makeRequest("GET", "/api/primary-business-owners?sort=first_name&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	dataMap := result["data"].(map[string]interface{})
	items := dataMap["data"].([]interface{})
	s.GreaterOrEqual(len(items), 3)

	// Verify first item is "Alpha"
	firstItem := items[0].(map[string]interface{})
	s.Equal("Alpha", firstItem["first_name"])

	// Test sort descending by first name
	resp, result = s.makeRequest("GET", "/api/primary-business-owners?sort=first_name&direction=desc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	dataMap = result["data"].(map[string]interface{})
	items = dataMap["data"].([]interface{})

	// Verify first item is "Zebra"
	firstItem = items[0].(map[string]interface{})
	s.Equal("Zebra", firstItem["first_name"])

	// Test sort ascending by last name
	resp, result = s.makeRequest("GET", "/api/primary-business-owners?sort=last_name&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	dataMap = result["data"].(map[string]interface{})
	items = dataMap["data"].([]interface{})
	s.GreaterOrEqual(len(items), 3)

	// Verify first item is "Wilson"
	firstItem = items[0].(map[string]interface{})
	s.Equal("Wilson", firstItem["last_name"])

	// Test sort descending by last name
	resp, result = s.makeRequest("GET", "/api/primary-business-owners?sort=last_name&direction=desc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	dataMap = result["data"].(map[string]interface{})
	items = dataMap["data"].([]interface{})

	// Verify first item is "Young"
	firstItem = items[0].(map[string]interface{})
	s.Equal("Young", firstItem["last_name"])
}

// ============================================================================
// COMBINED PAGINATION AND SORTING Tests
// ============================================================================

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TestPaginationWithSorting() {
	// First create an SME
	smeData := map[string]interface{}{
		"name":                         "Test Business for Combined Test",
		"business_category":            "retail",
		"sector":                       "commerce",
		"contact_phone":                "+265991234567",
		"contact_email":                "combined@example.com",
		"registration_number":          "BRNR-COMB01",
		"operational_start_date":       carbon.Now().Format("2006-01-02"),
		"business_improvement_aspects": []string{"financial_management"},
		"business_accessed_financing":  []string{"bank_loan"},
	}

	resp, smeResult := s.makeRequest("POST", "/api/smes", smeData)
	s.Equal(http.StatusCreated, resp.StatusCode, "Failed to create SME for combined test")
	smeID := uint(smeResult["data"].(map[string]interface{})["id"].(float64))

	// Create 15 PBOs with varied names
	for i := 1; i <= 15; i++ {
		data := s.getValidPBOData(smeID)
		data["first_name"] = fmt.Sprintf("Person%02d", 16-i) // Reverse order
		data["last_name"] = fmt.Sprintf("Family%02d", i)
		data["national_id_number"] = fmt.Sprintf("COMB0%03d", i)
		data["email"] = fmt.Sprintf("combined%02d@example.com", i)
		s.makeRequest("POST", "/api/primary-business-owners", data)
	}

	// Test pagination with sorting by first_name ascending
	resp, result := s.makeRequest("GET", "/api/primary-business-owners?page=1&pageSize=5&sort=first_name&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	dataMap := result["data"].(map[string]interface{})
	items := dataMap["data"].([]interface{})
	pagination := dataMap["pagination"].(map[string]interface{})

	s.Equal(5, len(items), "Should have 5 items on first page")
	s.Equal(float64(1), pagination["current_page"])

	// Verify items are sorted ascending
	firstItem := items[0].(map[string]interface{})
	s.Equal("Person01", firstItem["first_name"], "First item should be Person01")

	// Test second page with same sorting
	resp, result = s.makeRequest("GET", "/api/primary-business-owners?page=2&pageSize=5&sort=first_name&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	dataMap = result["data"].(map[string]interface{})
	items = dataMap["data"].([]interface{})
	pagination = dataMap["pagination"].(map[string]interface{})

	s.Equal(5, len(items), "Should have 5 items on second page")
	s.Equal(float64(2), pagination["current_page"])

	// Verify sorting continues across pages
	firstItemPage2 := items[0].(map[string]interface{})
	s.Equal("Person06", firstItemPage2["first_name"], "First item on page 2 should be Person06")

	// Test pagination with sorting by last_name descending
	resp, result = s.makeRequest("GET", "/api/primary-business-owners?page=1&pageSize=5&sort=last_name&direction=desc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	dataMap = result["data"].(map[string]interface{})
	items = dataMap["data"].([]interface{})

	s.Equal(5, len(items), "Should have 5 items")

	// Verify items are sorted descending by last name
	firstItem = items[0].(map[string]interface{})
	// Since we have Person01-15 with Family01-15, descending should start with Family15
	lastName := firstItem["last_name"].(string)
	s.Contains([]string{"Family15", "Family14", "Family13", "Family12", "Family11"}, lastName,
		"First item should have one of the highest family names")
}
