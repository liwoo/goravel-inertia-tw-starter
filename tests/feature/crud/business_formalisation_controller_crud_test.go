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

type BusinessFormalisationControllerCRUDTestSuite struct {
	suite.Suite
	tests.TestCase

	server     *httptest.Server
	client     *http.Client
	authCookie *http.Cookie
	testUser   *models.User
	testSme    *models.Sme
}

func TestBusinessFormalisationControllerCRUDTestSuite(t *testing.T) {
	suite.Run(t, &BusinessFormalisationControllerCRUDTestSuite{})
}

func (s *BusinessFormalisationControllerCRUDTestSuite) SetupSuite() {
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

func (s *BusinessFormalisationControllerCRUDTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}

func (s *BusinessFormalisationControllerCRUDTestSuite) SetupTest() {
	// Clean any existing test data first
	if orm := facades.Orm(); orm != nil {
		// Delete in correct order to respect foreign key constraints
		orm.Query().Exec("DELETE FROM business_formalisation")
		orm.Query().Exec("DELETE FROM smes")
		orm.Query().Exec("DELETE FROM user_roles")
		orm.Query().Exec("DELETE FROM role_permissions")
		orm.Query().Exec("DELETE FROM roles WHERE slug LIKE 'bf_%'")
		orm.Query().Exec("DELETE FROM permissions WHERE slug LIKE 'business_formalisation_%'")
		orm.Query().Exec("DELETE FROM users WHERE email LIKE 'bftest%@example.com'")
	}

	// Refresh database (runs migrations)
	s.RefreshDatabase()

	// Create test SME and user with permissions
	s.setupTestData()
}

func (s *BusinessFormalisationControllerCRUDTestSuite) TearDownTest() {
	// Clean up test data
	if orm := facades.Orm(); orm != nil {
		// Delete in correct order
		orm.Query().Exec("DELETE FROM business_formalisation")
		orm.Query().Exec("DELETE FROM smes")
		orm.Query().Exec("DELETE FROM user_roles")
		orm.Query().Exec("DELETE FROM role_permissions")
		orm.Query().Exec("DELETE FROM roles WHERE slug LIKE 'bf_%'")
		orm.Query().Exec("DELETE FROM permissions WHERE slug LIKE 'business_formalisation_%'")
		orm.Query().Exec("DELETE FROM users WHERE email LIKE 'bftest%@example.com'")
	}
}

func (s *BusinessFormalisationControllerCRUDTestSuite) setupTestData() {
	// Create admin role
	adminRole := &models.Role{
		Name:  "BF Admin",
		Slug:  "bf_admin",
		Level: 100,
	}
	s.Nil(facades.Orm().Query().Create(adminRole))

	// Create all CRUD permissions
	permissions := []string{"create", "read", "update", "delete", "view"}
	for _, perm := range permissions {
		permission := &models.Permission{
			Name:  fmt.Sprintf("business_formalisation_%s", perm),
			Slug:  fmt.Sprintf("business_formalisation_%s", perm),
			Scope: "by_all",
		}
		s.Nil(facades.Orm().Query().Create(permission))
		s.Nil(helpers.AssignPermissionToRole(adminRole, permission, "by_all"))
	}

	// Create test user
	user, err := helpers.SetupJWTUser("bftest@example.com", "password", adminRole)
	s.Nil(err)
	s.testUser = user

	// Create test SME
	createdBy := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:                 "USME-BF-TEST-001",
		Name:                       "Test SME for Business Formalisation",
		BusinessCategory:           "Technology",
		Sector:                     "IT Services",
		ContactPhone:               "+265999000111",
		ContactEmail:               "sme@test.mw",
		BusinessImprovementAspects: []string{"Innovation", "Growth"},
		BusinessAccessedFinancing:  []string{"Bank Loan", "Grant"},
		CreatedBy:                  &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(sme))
	s.testSme = sme

	// Login
	s.authCookie = s.loginUser("bftest@example.com", "password")
	s.NotNil(s.authCookie)
}

func (s *BusinessFormalisationControllerCRUDTestSuite) loginUser(email, password string) *http.Cookie {
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

func (s *BusinessFormalisationControllerCRUDTestSuite) makeRequest(method, path string, body interface{}) (*http.Response, map[string]interface{}) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			s.T().Logf("Failed to marshal request body: %v", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
		s.T().Logf("Request body: %s", string(jsonBody))
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

	// Log raw response for debugging
	if resp.StatusCode >= 500 {
		s.T().Logf("Server error - Status: %d, Body: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if len(respBody) > 0 {
		err = json.Unmarshal(respBody, &result)
		if err != nil {
			s.T().Logf("Failed to parse JSON response: %s", string(respBody))
			s.T().Logf("Raw response body: %s", string(respBody))
		}
	}

	return resp, result
}

func (s *BusinessFormalisationControllerCRUDTestSuite) getValidFormalisationData() map[string]interface{} {
	return map[string]interface{}{
		"sme_id":                   s.testSme.ID,
		"has_bank_account":         true,
		"has_tax_clarification":    true,
		"is_registered_for_vat":    false,
		"is_member_of_association": true,
		"is_affiliated":            false,
		"has_export_license":       false,
		"has_accessed_bds":         true,
		"annual_turnover":          50000.50,
		"estimated_value_of_assets": 250000.75,
	}
}

// ============================================================================
// CREATE Tests
// ============================================================================

func (s *BusinessFormalisationControllerCRUDTestSuite) TestCreateBusinessFormalisation() {
	formalisationData := s.getValidFormalisationData()

	resp, result := s.makeRequest("POST", "/api/business-formalisations", formalisationData)

	// Debug output if not successful
	if resp.StatusCode != http.StatusCreated {
		s.T().Logf("Response status: %d", resp.StatusCode)
		s.T().Logf("Response body: %+v", result)
	}

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.NotNil(result)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.NotNil(data["id"])
	s.Equal(float64(s.testSme.ID), data["sme_id"])
	s.Equal(true, data["has_bank_account"])
	s.Equal(true, data["has_tax_clarification"])
	s.Equal(false, data["is_registered_for_vat"])
	s.Equal(50000.50, data["annual_turnover"])
	s.Equal(250000.75, data["estimated_value_of_assets"])
	s.Equal(float64(s.testUser.ID), data["created_by"])
}

// CRITICAL TEST: Verify formalisation_score CANNOT be set via API
func (s *BusinessFormalisationControllerCRUDTestSuite) TestFormalisationScoreCannotBeSetViaCreate() {
	formalisationData := s.getValidFormalisationData()
	// Try to set formalisation_score via API (should be ignored)
	formalisationData["formalisation_score"] = 99

	resp, result := s.makeRequest("POST", "/api/business-formalisations", formalisationData)

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	formalisationID := uint(data["id"].(float64))

	// Verify formalisation_score was NOT set (should be 0 or null)
	var formalisation models.BusinessFormalisation
	err := facades.Orm().Query().Find(&formalisation, formalisationID)
	s.Nil(err)
	s.Equal(0, formalisation.FormalisationScore, "formalisation_score should NOT be settable via create API")
}

func (s *BusinessFormalisationControllerCRUDTestSuite) TestCreateFormalisationValidationErrors() {
	testCases := []struct {
		name        string
		data        map[string]interface{}
		expectError bool
		errorField  string
	}{
		{
			name: "Missing required sme_id",
			data: map[string]interface{}{
				"has_bank_account": true,
			},
			expectError: true,
			errorField:  "sme_id",
		},
		{
			name: "Invalid annual_turnover type",
			data: map[string]interface{}{
				"sme_id":           s.testSme.ID,
				"annual_turnover": "not-a-number",
			},
			expectError: true,
			errorField:  "annual_turnover",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, result := s.makeRequest("POST", "/api/business-formalisations", tc.data)

			if tc.expectError {
				s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
				s.False(result["success"].(bool))
				s.Contains(strings.ToLower(result["message"].(string)), "validation")
			} else {
				s.Equal(http.StatusCreated, resp.StatusCode)
				s.True(result["success"].(bool))
			}
		})
	}
}

// ============================================================================
// READ Tests
// ============================================================================

func (s *BusinessFormalisationControllerCRUDTestSuite) TestGetBusinessFormalisation() {
	// Create a formalisation first
	createdBy := int(s.testUser.ID)
	formalisation := &models.BusinessFormalisation{
		SmeID:                  int(s.testSme.ID),
		HasBankAccount:         true,
		HasTaxClarification:    false,
		IsRegisteredForVat:     true,
		IsMemberOfAssociation:  false,
		IsAffiliated:           true,
		HasExportLicense:       false,
		HasAccessedBds:         true,
		AnnualTurnover:         75000.25,
		EstimatedValueOfAssets: 300000.00,
		FormalisationScore:     0, // Read-only, not set via API
		CreatedBy:              &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(formalisation))

	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/business-formalisations/%d", formalisation.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal(float64(formalisation.ID), data["id"].(float64))
	s.Equal(float64(s.testSme.ID), data["sme_id"])
	s.Equal(true, data["has_bank_account"])
	s.Equal(false, data["has_tax_clarification"])
	s.Equal(75000.25, data["annual_turnover"])
}

func (s *BusinessFormalisationControllerCRUDTestSuite) TestGetFormalisationNotFound() {
	resp, result := s.makeRequest("GET", "/api/business-formalisations/99999", nil)

	s.Equal(http.StatusNotFound, resp.StatusCode)
	s.False(result["success"].(bool))
}

// ============================================================================
// UPDATE Tests
// ============================================================================

func (s *BusinessFormalisationControllerCRUDTestSuite) TestUpdateBusinessFormalisation() {
	// Create a formalisation first
	createdBy := int(s.testUser.ID)
	formalisation := &models.BusinessFormalisation{
		SmeID:                  int(s.testSme.ID),
		HasBankAccount:         false,
		HasTaxClarification:    false,
		IsRegisteredForVat:     false,
		IsMemberOfAssociation:  false,
		IsAffiliated:           false,
		HasExportLicense:       false,
		HasAccessedBds:         false,
		AnnualTurnover:         10000.00,
		EstimatedValueOfAssets: 50000.00,
		CreatedBy:              &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(formalisation))

	updateData := map[string]interface{}{
		"has_bank_account":        true,
		"has_tax_clarification":   true,
		"annual_turnover":          25000.00,
		"estimated_value_of_assets": 100000.00,
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/business-formalisations/%d", formalisation.ID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal(true, data["has_bank_account"])
	s.Equal(true, data["has_tax_clarification"])
	s.Equal(25000.00, data["annual_turnover"])
	s.Equal(100000.00, data["estimated_value_of_assets"])
	// Unchanged fields should remain
	s.Equal(false, data["is_registered_for_vat"])
}

// CRITICAL TEST: Verify formalisation_score CANNOT be updated via API
func (s *BusinessFormalisationControllerCRUDTestSuite) TestFormalisationScoreCannotBeUpdatedViaAPI() {
	// Create a formalisation first
	createdBy := int(s.testUser.ID)
	formalisation := &models.BusinessFormalisation{
		SmeID:                  int(s.testSme.ID),
		HasBankAccount:         true,
		FormalisationScore:     0, // Initial value
		CreatedBy:              &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(formalisation))

	// Try to update formalisation_score via API (should be ignored)
	updateData := map[string]interface{}{
		"formalisation_score": 50,
		"has_bank_account":    false,
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/business-formalisations/%d", formalisation.ID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	// Verify formalisation_score was NOT updated (should remain 0)
	var updatedFormalisation models.BusinessFormalisation
	err := facades.Orm().Query().Find(&updatedFormalisation, formalisation.ID)
	s.Nil(err)
	s.Equal(0, updatedFormalisation.FormalisationScore, "formalisation_score should NOT be updatable via API")
	// Verify other fields were updated
	s.Equal(false, updatedFormalisation.HasBankAccount)
}

func (s *BusinessFormalisationControllerCRUDTestSuite) TestUpdateFormalisationNotFound() {
	updateData := map[string]interface{}{
		"has_bank_account": true,
	}

	resp, result := s.makeRequest("PUT", "/api/business-formalisations/99999", updateData)

	s.Equal(http.StatusNotFound, resp.StatusCode)
	s.False(result["success"].(bool))
}

// ============================================================================
// DELETE Tests
// ============================================================================

func (s *BusinessFormalisationControllerCRUDTestSuite) TestDeleteBusinessFormalisation() {
	// Create a formalisation
	createdBy := int(s.testUser.ID)
	formalisation := &models.BusinessFormalisation{
		SmeID:              int(s.testSme.ID),
		HasBankAccount:     true,
		FormalisationScore: 0,
		CreatedBy:          &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(formalisation))

	resp, _ := s.makeRequest("DELETE", fmt.Sprintf("/api/business-formalisations/%d", formalisation.ID), nil)

	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify soft delete
	var deletedFormalisation models.BusinessFormalisation
	err := facades.Orm().Query().WithTrashed().Where("id", formalisation.ID).First(&deletedFormalisation)
	s.Nil(err)
	s.NotNil(deletedFormalisation.DeletedAt)
}

func (s *BusinessFormalisationControllerCRUDTestSuite) TestDeleteFormalisationNotFound() {
	resp, result := s.makeRequest("DELETE", "/api/business-formalisations/99999", nil)

	s.Equal(http.StatusNotFound, resp.StatusCode)
	s.False(result["success"].(bool))
}

// ============================================================================
// PAGINATION Tests
// ============================================================================

func (s *BusinessFormalisationControllerCRUDTestSuite) TestPagination() {
	// Create 25 formalisations
	createdBy := int(s.testUser.ID)
	for i := 1; i <= 25; i++ {
		formalisation := &models.BusinessFormalisation{
			SmeID:              int(s.testSme.ID),
			HasBankAccount:     i%2 == 0,
			HasTaxClarification: i%3 == 0,
			AnnualTurnover:     float64(i * 1000),
			CreatedBy:          &createdBy,
		}
		s.Nil(facades.Orm().Query().Create(formalisation))
	}

	// Test first page (default page size 20)
	resp, result := s.makeRequest("GET", "/api/business-formalisations?page=1", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	pagination := data["pagination"].(map[string]interface{})

	s.Equal(20, len(items))
	s.Equal(float64(1), pagination["current_page"])
	s.Equal(float64(25), pagination["total"])
	s.Equal(float64(2), pagination["last_page"])
}

// ============================================================================
// SEARCH Tests
// ============================================================================

func (s *BusinessFormalisationControllerCRUDTestSuite) TestSearch() {
	// Create searchable formalisations
	createdBy := int(s.testUser.ID)

	// High turnover
	f1 := &models.BusinessFormalisation{
		SmeID:          int(s.testSme.ID),
		AnnualTurnover: 100000.00,
		CreatedBy:      &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(f1))

	// Low turnover
	f2 := &models.BusinessFormalisation{
		SmeID:          int(s.testSme.ID),
		AnnualTurnover: 5000.00,
		CreatedBy:      &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(f2))

	// Search by SME ID
	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/business-formalisations/search?q=%d", s.testSme.ID), nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 2)
}
