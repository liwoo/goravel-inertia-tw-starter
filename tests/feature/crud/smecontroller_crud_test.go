package crud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/stretchr/testify/suite"

	"smedi-sme-db/app/models"
	"smedi-sme-db/tests"
	"smedi-sme-db/tests/helpers"
)

type SmeControllerCRUDTestSuite struct {
	suite.Suite
	tests.TestCase

	server     *httptest.Server
	client     *http.Client
	authCookie *http.Cookie
	testUser   *models.User
}

func TestSmeControllerCRUDTestSuite(t *testing.T) {
	suite.Run(t, &SmeControllerCRUDTestSuite{})
}

func (s *SmeControllerCRUDTestSuite) SetupSuite() {
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

func (s *SmeControllerCRUDTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}

func (s *SmeControllerCRUDTestSuite) SetupTest() {
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

func (s *SmeControllerCRUDTestSuite) TearDownTest() {
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
		orm.Query().Exec("DELETE FROM permissions WHERE slug LIKE 'smes_%'")
	}
}

func (s *SmeControllerCRUDTestSuite) setupTestUser() {
	// Create admin role
	adminRole := &models.Role{
		Name:  "SME Admin",
		Slug:  "sme_admin",
		Level: 100,
	}
	s.Nil(facades.Orm().Query().Create(adminRole))

	// Create all CRUD permissions
	permissions := []string{"create", "read", "update", "delete", "view"}
	for _, perm := range permissions {
		permission := &models.Permission{
			Name:  fmt.Sprintf("smes_%s", perm),
			Slug:  fmt.Sprintf("smes_%s", perm),
			Scope: "by_all",
		}
		s.Nil(facades.Orm().Query().Create(permission))
		s.Nil(helpers.AssignPermissionToRole(adminRole, permission, "by_all"))
	}

	// Create test user
	user, err := helpers.SetupJWTUser("smetest@example.com", "password", adminRole)
	s.Nil(err)
	s.testUser = user

	// Login
	s.authCookie = s.loginUser("smetest@example.com", "password")
	s.NotNil(s.authCookie)
}

func (s *SmeControllerCRUDTestSuite) loginUser(email, password string) *http.Cookie {
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

func (s *SmeControllerCRUDTestSuite) makeRequest(method, path string, body interface{}) (*http.Response, map[string]interface{}) {
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

func (s *SmeControllerCRUDTestSuite) getValidSmeData() map[string]interface{} {
	return map[string]interface{}{
		"name":                         "Test SME Business",
		"business_category":            "Retail",
		"sector":                       "Commerce",
		"contact_phone":                "+265999123456",
		"contact_email":                "contact@testbusiness.mw",
		"business_improvement_aspects": []string{"Marketing", "Operations"},
		"business_accessed_financing":  []string{"Bank Loan", "Microfinance"},
	}
}

// ============================================================================
// CREATE Tests
// ============================================================================

func (s *SmeControllerCRUDTestSuite) TestCreateSme() {
	smeData := s.getValidSmeData()

	resp, result := s.makeRequest("POST", "/api/smes", smeData)

	// Debug: print the response if it's not successful
	if resp.StatusCode != http.StatusCreated {
		s.T().Logf("Response status: %d", resp.StatusCode)
		s.T().Logf("Response body: %+v", result)
	}

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.NotNil(data["id"])
	s.Equal(smeData["name"], data["name"])
	s.Equal(float64(s.testUser.ID), data["created_by"])
}

func (s *SmeControllerCRUDTestSuite) TestCreateSmeValidation() {
	// Missing required fields
	smeData := map[string]interface{}{
		"name": "Test Business",
	}

	resp, result := s.makeRequest("POST", "/api/smes", smeData)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
	s.Contains(strings.ToLower(result["message"].(string)), "validation")
}

func (s *SmeControllerCRUDTestSuite) TestCreateSmeInvalidEmail() {
	smeData := s.getValidSmeData()
	smeData["contact_email"] = "invalid-email"

	resp, result := s.makeRequest("POST", "/api/smes", smeData)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
}

// ============================================================================
// READ Tests
// ============================================================================

func (s *SmeControllerCRUDTestSuite) TestGetSme() {
	// Create an SME first
	createdBy := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:                 "USME-TEST-002",
		Name:                       "Test SME Item",
		BusinessCategory:           "Manufacturing",
		Sector:                     "Industrial",
		ContactPhone:               "+265999111222",
		ContactEmail:               "test@sme.mw",
		BusinessImprovementAspects: []string{"Quality", "Efficiency"},
		BusinessAccessedFinancing:  []string{"Grant"},
		CreatedBy:                  &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(sme))

	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/smes/%d", sme.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal(float64(sme.ID), data["id"].(float64))
	s.Equal(sme.UsmeNumber, data["usme_number"])
	s.Equal(sme.Name, data["name"])
}

func (s *SmeControllerCRUDTestSuite) TestGetSmeNotFound() {
	resp, result := s.makeRequest("GET", "/api/smes/99999", nil)

	s.Equal(http.StatusNotFound, resp.StatusCode)
	s.False(result["success"].(bool))
}

// ============================================================================
// UPDATE Tests
// ============================================================================

func (s *SmeControllerCRUDTestSuite) TestUpdateSme() {
	// Create an SME first
	createdBy := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:                 "USME-TEST-003",
		Name:                       "Original Name",
		BusinessCategory:           "Services",
		Sector:                     "Technology",
		ContactPhone:               "+265999333444",
		ContactEmail:               "original@sme.mw",
		BusinessImprovementAspects: []string{"Innovation"},
		BusinessAccessedFinancing:  []string{"Venture Capital"},
		CreatedBy:                  &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(sme))

	updateData := map[string]interface{}{
		"name":          "Updated SME Name",
		"contact_email": "updated@sme.mw",
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/smes/%d", sme.ID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal("Updated SME Name", data["name"])
	s.Equal("updated@sme.mw", data["contact_email"])
}

func (s *SmeControllerCRUDTestSuite) TestUpdateSmeNotFound() {
	updateData := map[string]interface{}{
		"name": "Updated Name",
	}

	resp, result := s.makeRequest("PUT", "/api/smes/99999", updateData)

	s.Equal(http.StatusNotFound, resp.StatusCode)
	s.False(result["success"].(bool))
}

func (s *SmeControllerCRUDTestSuite) TestUpdateSmeValidationError() {
	// Create an SME first
	createdBy := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:                 "USME-TEST-VAL-001",
		Name:                       "Validation Test SME",
		BusinessCategory:           "Services",
		Sector:                     "Technology",
		ContactPhone:               "+265999333555",
		ContactEmail:               "validation@sme.mw",
		BusinessImprovementAspects: []string{"Innovation"},
		BusinessAccessedFinancing:  []string{"Venture Capital"},
		CreatedBy:                  &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(sme))

	// Try to update with an invalid classification value
	updateData := map[string]interface{}{
		"classification": "INVALID_CLASS",
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/smes/%d", sme.ID), updateData)

	// Should return 422 Unprocessable Entity due to validation error
	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
}

func (s *SmeControllerCRUDTestSuite) TestUpdateSmeValidClassification() {
	// Create an SME first
	createdBy := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:                 "USME-TEST-VAL-002",
		Name:                       "Valid Classification Test SME",
		BusinessCategory:           "Manufacturing",
		Sector:                     "Industrial",
		ContactPhone:               "+265999333666",
		ContactEmail:               "validclass@sme.mw",
		BusinessImprovementAspects: []string{"Quality"},
		BusinessAccessedFinancing:  []string{"Bank Loan"},
		Classification:             "Unclassified",
		CreatedBy:                  &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(sme))

	// Update with a valid classification value
	updateData := map[string]interface{}{
		"classification": "Micro",
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/smes/%d", sme.ID), updateData)

	// Should return 200 OK
	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal("Micro", data["classification"])
}

// ============================================================================
// DELETE Tests
// ============================================================================

func (s *SmeControllerCRUDTestSuite) TestDeleteSme() {
	// Create an SME first
	createdBy := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:                 "USME-TEST-004",
		Name:                       "To Be Deleted",
		BusinessCategory:           "Retail",
		Sector:                     "Commerce",
		ContactPhone:               "+265999555666",
		ContactEmail:               "delete@sme.mw",
		BusinessImprovementAspects: []string{"Customer Service"},
		BusinessAccessedFinancing:  []string{"Personal Savings"},
		CreatedBy:                  &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(sme))

	resp, _ := s.makeRequest("DELETE", fmt.Sprintf("/api/smes/%d", sme.ID), nil)

	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify soft delete
	var deletedSme models.Sme
	err := facades.Orm().Query().WithTrashed().Where("id", sme.ID).First(&deletedSme)
	s.Nil(err)
	s.NotNil(deletedSme.DeletedAt)
}

func (s *SmeControllerCRUDTestSuite) TestDeleteSmeNotFound() {
	resp, result := s.makeRequest("DELETE", "/api/smes/99999", nil)

	s.Equal(http.StatusNotFound, resp.StatusCode)
	s.False(result["success"].(bool))
}

// ============================================================================
// PAGINATION Tests
// ============================================================================

func (s *SmeControllerCRUDTestSuite) TestPagination() {
	// Create 25 SMEs
	createdBy := int(s.testUser.ID)
	for i := 1; i <= 25; i++ {
		sme := &models.Sme{
			UsmeNumber:                 fmt.Sprintf("USME-PAGE-%03d", i),
			Name:                       fmt.Sprintf("SME Business %02d", i),
			BusinessCategory:           "Retail",
			Sector:                     "Commerce",
			ContactPhone:               fmt.Sprintf("+265999%06d", i),
			ContactEmail:               fmt.Sprintf("sme%02d@test.mw", i),
			BusinessImprovementAspects: []string{"Growth"},
			BusinessAccessedFinancing:  []string{"Self-funded"},
			CreatedBy:                  &createdBy,
		}
		s.Nil(facades.Orm().Query().Create(sme))
	}

	// Test first page (default page size 20)
	resp, result := s.makeRequest("GET", "/api/smes?page=1", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	pagination := data["pagination"].(map[string]interface{})

	s.Equal(20, len(items))
	s.Equal(float64(1), pagination["current_page"])
	s.Equal(float64(25), pagination["total"])
	s.Equal(float64(2), pagination["last_page"])

	// Test custom page size
	resp, result = s.makeRequest("GET", "/api/smes?page=1&pageSize=10", nil)
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

func (s *SmeControllerCRUDTestSuite) TestSorting() {
	// Create SMEs with different names
	createdBy := int(s.testUser.ID)
	names := []string{"Zebra Business", "Alpha Company", "Beta Enterprises"}
	for i, name := range names {
		sme := &models.Sme{
			UsmeNumber:                 fmt.Sprintf("USME-SORT-%03d", i+1),
			Name:                       name,
			BusinessCategory:           "Services",
			Sector:                     "Technology",
			ContactPhone:               fmt.Sprintf("+265991%06d", i),
			ContactEmail:               fmt.Sprintf("sort%d@test.mw", i),
			BusinessImprovementAspects: []string{"Innovation"},
			BusinessAccessedFinancing:  []string{"Angel Investment"},
			CreatedBy:                  &createdBy,
		}
		s.Nil(facades.Orm().Query().Create(sme))
	}

	// Test sort ascending by name
	resp, result := s.makeRequest("GET", "/api/smes?sort=name&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 3)

	// Verify first item is "Alpha Company"
	firstItem := items[0].(map[string]interface{})
	s.Equal("Alpha Company", firstItem["name"])

	// Test sort descending by name
	resp, result = s.makeRequest("GET", "/api/smes?sort=name&direction=desc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})

	// Verify first item is "Zebra Business"
	firstItem = items[0].(map[string]interface{})
	s.Equal("Zebra Business", firstItem["name"])
}

func (s *SmeControllerCRUDTestSuite) TestSortingByFormalisationScore() {
	createdBy := int(s.testUser.ID)

	// Create SMEs with different formalisation scores
	// Score -1 means no formalisation record (should be treated as 0)
	testData := []struct {
		name  string
		score int // -1 means no formalisation record
	}{
		{"No Formalisation SME", -1}, // Should sort as 0
		{"Low Score SME", 20},
		{"High Score SME", 90},
		{"Medium Score SME", 50},
	}

	for i, td := range testData {
		// Create SME
		sme := &models.Sme{
			UsmeNumber:                 fmt.Sprintf("USME-SCORE-%03d", i+1),
			Name:                       td.name,
			BusinessCategory:           "Services",
			Sector:                     "Technology",
			ContactPhone:               fmt.Sprintf("+265993%06d", i),
			ContactEmail:               fmt.Sprintf("score%d@test.mw", i),
			BusinessImprovementAspects: []string{"Innovation"},
			BusinessAccessedFinancing:  []string{"Grant"},
			CreatedBy:                  &createdBy,
		}
		s.Nil(facades.Orm().Query().Create(sme))

		// Create corresponding business formalisation with score (skip if score is -1)
		if td.score >= 0 {
			formalisation := &models.BusinessFormalisation{
				HasBankAccount:         true,
				HasTaxClarification:    td.score > 50,
				IsRegisteredForVat:     td.score > 70,
				IsMemberOfAssociation:  false,
				IsAffiliated:           false,
				HasExportLicense:       false,
				HasAccessedBds:         td.score > 30,
				AnnualTurnover:         float64(td.score * 1000),
				EstimatedValueOfAssets: float64(td.score * 5000),
				FormalisationScore:     td.score,
				SmeID:                  int(sme.ID),
				CreatedBy:              &createdBy,
			}
			s.Nil(facades.Orm().Query().Create(formalisation))
		}
	}

	// Test sort ascending by formalisationScore (camelCase as frontend sends it)
	resp, result := s.makeRequest("GET", "/api/smes?sort=formalisationScore&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 4, "Should have at least 4 SMEs")

	// Extract formalisation scores in order (0 for SMEs without formalisation)
	var scores []int
	var names []string
	for _, item := range items {
		itemData := item.(map[string]interface{})
		name := itemData["name"].(string)
		names = append(names, name)

		// Check if business_formalisation is present
		if bf, ok := itemData["business_formalisation"].(map[string]interface{}); ok && bf != nil {
			if score, ok := bf["formalisation_score"].(float64); ok {
				scores = append(scores, int(score))
			} else {
				scores = append(scores, 0) // No score means 0
			}
		} else {
			scores = append(scores, 0) // No formalisation record means 0
		}
	}

	s.T().Logf("Ascending order - Names: %v, Scores: %v", names, scores)

	// Verify ascending order - scores should be in increasing order
	for i := 1; i < len(scores); i++ {
		s.LessOrEqual(scores[i-1], scores[i],
			"Formalisation scores should be in ascending order: %v (names: %v)", scores, names)
	}

	// Verify first item has lowest score (0 - the one without formalisation)
	s.Equal(0, scores[0], "First item should have score 0 (no formalisation record)")
	s.Equal("No Formalisation SME", names[0], "First item should be the SME without formalisation")

	// Test sort descending by formalisationScore
	resp, result = s.makeRequest("GET", "/api/smes?sort=formalisationScore&direction=desc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})

	// Extract formalisation scores in order for descending check
	scores = nil
	names = nil
	for _, item := range items {
		itemData := item.(map[string]interface{})
		name := itemData["name"].(string)
		names = append(names, name)

		if bf, ok := itemData["business_formalisation"].(map[string]interface{}); ok && bf != nil {
			if score, ok := bf["formalisation_score"].(float64); ok {
				scores = append(scores, int(score))
			} else {
				scores = append(scores, 0)
			}
		} else {
			scores = append(scores, 0)
		}
	}

	s.T().Logf("Descending order - Names: %v, Scores: %v", names, scores)

	// Verify descending order - scores should be in decreasing order
	for i := 1; i < len(scores); i++ {
		s.GreaterOrEqual(scores[i-1], scores[i],
			"Formalisation scores should be in descending order: %v (names: %v)", scores, names)
	}

	// Verify first item has highest score (90)
	s.Equal(90, scores[0], "First item should have the highest formalisation score")
	s.Equal("High Score SME", names[0], "First item should be the High Score SME")

	// Verify last item has lowest score (0 - the one without formalisation)
	lastIdx := len(scores) - 1
	s.Equal(0, scores[lastIdx], "Last item should have score 0 (no formalisation record)")
	s.Equal("No Formalisation SME", names[lastIdx], "Last item should be the SME without formalisation")
}

func (s *SmeControllerCRUDTestSuite) TestSortingByFormalisationScoreSnakeCase() {
	createdBy := int(s.testUser.ID)

	// Create SMEs with different formalisation scores (using snake_case sort parameter)
	testData := []struct {
		name   string
		score  int
	}{
		{"Snake Low", 15},
		{"Snake High", 85},
		{"Snake Mid", 45},
	}

	for i, td := range testData {
		sme := &models.Sme{
			UsmeNumber:                 fmt.Sprintf("USME-SNAKE-%03d", i+1),
			Name:                       td.name,
			BusinessCategory:           "Manufacturing",
			Sector:                     "Industrial",
			ContactPhone:               fmt.Sprintf("+265994%06d", i),
			ContactEmail:               fmt.Sprintf("snake%d@test.mw", i),
			BusinessImprovementAspects: []string{"Quality"},
			BusinessAccessedFinancing:  []string{"Loan"},
			CreatedBy:                  &createdBy,
		}
		s.Nil(facades.Orm().Query().Create(sme))

		formalisation := &models.BusinessFormalisation{
			HasBankAccount:         true,
			HasTaxClarification:    td.score > 50,
			FormalisationScore:     td.score,
			SmeID:                  int(sme.ID),
			CreatedBy:              &createdBy,
		}
		s.Nil(facades.Orm().Query().Create(formalisation))
	}

	// Test sort using snake_case parameter (formalisation_score)
	resp, result := s.makeRequest("GET", "/api/smes?sort=formalisation_score&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 3, "Should have at least 3 SMEs")

	// Extract and verify scores are in ascending order
	var scores []int
	for _, item := range items {
		itemData := item.(map[string]interface{})
		if bf, ok := itemData["business_formalisation"].(map[string]interface{}); ok {
			if score, ok := bf["formalisation_score"].(float64); ok {
				scores = append(scores, int(score))
			}
		}
	}

	// Verify ascending order
	for i := 1; i < len(scores); i++ {
		s.LessOrEqual(scores[i-1], scores[i],
			"Formalisation scores should be in ascending order when using snake_case: %v", scores)
	}
}

// ============================================================================
// SEARCH Tests
// ============================================================================

func (s *SmeControllerCRUDTestSuite) TestSearch() {
	// Create searchable SMEs
	searchableItems := []struct {
		usmeNumber string
		name       string
	}{
		{"USME-SEARCH-001", "Searchable Tech Company"},
		{"USME-SEARCH-002", "Another Searchable Business"},
		{"USME-OTHER-003", "Non-matching Business"},
	}

	createdBy := int(s.testUser.ID)
	for _, item := range searchableItems {
		sme := &models.Sme{
			UsmeNumber:                 item.usmeNumber,
			Name:                       item.name,
			BusinessCategory:           "Technology",
			Sector:                     "IT",
			ContactPhone:               "+265999000111",
			ContactEmail:               "search@test.mw",
			BusinessImprovementAspects: []string{"Digital Transformation"},
			BusinessAccessedFinancing:  []string{"Crowdfunding"},
			CreatedBy:                  &createdBy,
		}
		s.Nil(facades.Orm().Query().Create(sme))
	}

	// Search by query
	resp, result := s.makeRequest("GET", "/api/smes/search?q=Searchable", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})

	// Should find 2 items with "Searchable" in the name
	s.Equal(2, len(items))

	// Verify both items contain "Searchable"
	for _, item := range items {
		itemData := item.(map[string]interface{})
		s.Contains(itemData["name"].(string), "Searchable")
	}
}

func (s *SmeControllerCRUDTestSuite) TestSearchByUsmeNumber() {
	// Create SME with unique USME number
	createdBy := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:                 "USME-UNIQUE-999",
		Name:                       "Unique Business",
		BusinessCategory:           "Manufacturing",
		Sector:                     "Industrial",
		ContactPhone:               "+265999888777",
		ContactEmail:               "unique@test.mw",
		BusinessImprovementAspects: []string{"Automation"},
		BusinessAccessedFinancing:  []string{"Government Grant"},
		CreatedBy:                  &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(sme))

	// Search by USME number
	resp, result := s.makeRequest("GET", "/api/smes/search?q=UNIQUE-999", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})

	s.GreaterOrEqual(len(items), 1)
	firstItem := items[0].(map[string]interface{})
	s.Contains(firstItem["usme_number"].(string), "UNIQUE-999")
}

// ============================================================================
// LIST/INDEX Tests
// ============================================================================

func (s *SmeControllerCRUDTestSuite) TestListSmes() {
	// Create a few SMEs
	createdBy := int(s.testUser.ID)
	for i := 1; i <= 3; i++ {
		sme := &models.Sme{
			UsmeNumber:                 fmt.Sprintf("USME-LIST-%03d", i),
			Name:                       fmt.Sprintf("List Business %d", i),
			BusinessCategory:           "Services",
			Sector:                     "Consulting",
			ContactPhone:               fmt.Sprintf("+265992%06d", i),
			ContactEmail:               fmt.Sprintf("list%d@test.mw", i),
			BusinessImprovementAspects: []string{"Training"},
			BusinessAccessedFinancing:  []string{"Bootstrap"},
			CreatedBy:                  &createdBy,
		}
		s.Nil(facades.Orm().Query().Create(sme))
	}

	resp, result := s.makeRequest("GET", "/api/smes", nil)
	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})

	s.GreaterOrEqual(len(items), 3)
}

// ============================================================================
// RELATIONSHIP FETCHING Tests
// ============================================================================

func (s *SmeControllerCRUDTestSuite) TestFetchPrimaryBusinessOwner() {
	// Create an SME first
	createdBy := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:                 "USME-REL-001",
		Name:                       "Test Business With Owner",
		BusinessCategory:           "Technology",
		Sector:                     "IT",
		ContactPhone:               "+265999111222",
		ContactEmail:               "owner@test.mw",
		BusinessImprovementAspects: []string{"Innovation"},
		BusinessAccessedFinancing:  []string{"Grant"},
		CreatedBy:                  &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(sme))

	// Create a primary business owner for this SME
	owner := &models.PrimaryBusinessOwner{
		FirstName:        "John",
		LastName:         "Doe",
		DateOfBirth:      *carbon.NewDateTime(carbon.Parse("1985-06-15")),
		Nationality:      "Malawian",
		NationalIdNumber: "MWI123456789",
		Gender:           "MALE",
		EducationLevel:   "Tertiary",
		MalawianStatus:   "Citizen",
		PhoneNumber:      "+265991234567",
		SmeID:            int(sme.ID),
		CreatedBy:        &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(owner))

	// Fetch the primary business owner
	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/smes/%d/primary_business_owner", sme.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.NotNil(data)
	s.Equal("John", data["first_name"])
	s.Equal("Doe", data["last_name"])
	s.Equal("MWI123456789", data["national_id_number"])
}

func (s *SmeControllerCRUDTestSuite) TestFetchPrimaryBusinessOwnerNotFound() {
	// Create an SME without a primary business owner
	createdBy := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:                 "USME-REL-002",
		Name:                       "Business Without Owner",
		BusinessCategory:           "Retail",
		Sector:                     "Commerce",
		ContactPhone:               "+265999222333",
		ContactEmail:               "noowner@test.mw",
		BusinessImprovementAspects: []string{"Growth"},
		BusinessAccessedFinancing:  []string{"Self-funded"},
		CreatedBy:                  &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(sme))

	// Try to fetch primary business owner (should return null)
	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/smes/%d/primary_business_owner", sme.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))
	// Data should be null when no owner exists
	s.Nil(result["data"])
}

func (s *SmeControllerCRUDTestSuite) TestFetchAdditionalBusinessMembers() {
	// Create an SME
	createdBy := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:                 "USME-REL-003",
		Name:                       "Business With Members",
		BusinessCategory:           "Manufacturing",
		Sector:                     "Industrial",
		ContactPhone:               "+265999333444",
		ContactEmail:               "members@test.mw",
		BusinessImprovementAspects: []string{"Quality"},
		BusinessAccessedFinancing:  []string{"Bank Loan"},
		CreatedBy:                  &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(sme))

	// Create multiple additional business members
	members := []models.AdditionalBusinessMember{
		{
			FirstName:        "Jane",
			LastName:         "Smith",
			Nationality:      "Malawian",
			NationalIdNumber: "MWI987654321",
			PhoneNumber:      "+265991111111",
			IsIntern:         false,
			IsPartTime:       false,
			SmeId:            int(sme.ID),
			CreatedBy:        &createdBy,
		},
		{
			FirstName:        "Peter",
			LastName:         "Johnson",
			Nationality:      "Malawian",
			NationalIdNumber: "MWI111222333",
			PhoneNumber:      "+265992222222",
			IsIntern:         true,
			IsPartTime:       true,
			SmeId:            int(sme.ID),
			CreatedBy:        &createdBy,
		},
	}

	for _, member := range members {
		s.Nil(facades.Orm().Query().Create(&member))
	}

	// Fetch additional business members
	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/smes/%d/additional_business_members", sme.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].([]interface{})
	s.Equal(2, len(data))

	// Verify first member
	firstMember := data[0].(map[string]interface{})
	s.Equal("Jane", firstMember["first_name"])
	s.Equal("Smith", firstMember["last_name"])
}

func (s *SmeControllerCRUDTestSuite) TestFetchAdditionalBusinessMembersEmpty() {
	// Create an SME without additional members
	createdBy := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:                 "USME-REL-004",
		Name:                       "Business Without Members",
		BusinessCategory:           "Services",
		Sector:                     "Consulting",
		ContactPhone:               "+265999444555",
		ContactEmail:               "nomembers@test.mw",
		BusinessImprovementAspects: []string{"Training"},
		BusinessAccessedFinancing:  []string{"Bootstrap"},
		CreatedBy:                  &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(sme))

	// Fetch additional members (should return empty array)
	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/smes/%d/additional_business_members", sme.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].([]interface{})
	s.Equal(0, len(data))
}

func (s *SmeControllerCRUDTestSuite) TestFetchBusinessFormalisation() {
	// Create an SME
	createdBy := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:                 "USME-REL-005",
		Name:                       "Business With Formalisation",
		BusinessCategory:           "Technology",
		Sector:                     "IT",
		ContactPhone:               "+265999555666",
		ContactEmail:               "formal@test.mw",
		BusinessImprovementAspects: []string{"Compliance"},
		BusinessAccessedFinancing:  []string{"Investment"},
		CreatedBy:                  &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(sme))

	// Create business formalisation record
	formalisation := &models.BusinessFormalisation{
		HasBankAccount:         true,
		HasTaxClarification:    true,
		IsRegisteredForVat:     true,
		IsMemberOfAssociation:  false,
		IsAffiliated:           true,
		HasExportLicense:       false,
		HasAccessedBds:         true,
		AnnualTurnover:         150000.50,
		EstimatedValueOfAssets: 500000.00,
		FormalisationScore:     75,
		SmeID:                  int(sme.ID),
		CreatedBy:              &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(formalisation))

	// Fetch business formalisation
	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/smes/%d/business_formalisation", sme.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.NotNil(data)
	s.Equal(true, data["has_bank_account"])
	s.Equal(true, data["has_tax_clarification"])
	s.Equal(float64(75), data["formalisation_score"])
}

func (s *SmeControllerCRUDTestSuite) TestFetchBusinessFormalisationNotFound() {
	// Create an SME without formalisation data
	createdBy := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:                 "USME-REL-006",
		Name:                       "Business Without Formalisation",
		BusinessCategory:           "Retail",
		Sector:                     "Commerce",
		ContactPhone:               "+265999666777",
		ContactEmail:               "noformal@test.mw",
		BusinessImprovementAspects: []string{"Marketing"},
		BusinessAccessedFinancing:  []string{"Personal Savings"},
		CreatedBy:                  &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(sme))

	// Try to fetch formalisation (should return null)
	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/smes/%d/business_formalisation", sme.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))
	s.Nil(result["data"])
}

func (s *SmeControllerCRUDTestSuite) TestFetchBusinessEmployeeSummary() {
	// Create an SME
	createdBy := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:                 "USME-REL-007",
		Name:                       "Business With Employees",
		BusinessCategory:           "Manufacturing",
		Sector:                     "Industrial",
		ContactPhone:               "+265999777888",
		ContactEmail:               "employees@test.mw",
		BusinessImprovementAspects: []string{"HR"},
		BusinessAccessedFinancing:  []string{"Government Grant"},
		CreatedBy:                  &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(sme))

	// Create employee summary
	summary := &models.BusinessEmployeeSummary{
		FullTimeMales:   10,
		FullTimeFemales: 8,
		PartTimeMales:   3,
		PartTimeFemales: 5,
		InternMales:     2,
		InternFemales:   4,
		SmeID:           int(sme.ID),
		CreatedBy:       &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(summary))

	// Fetch employee summary
	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/smes/%d/business_employee_summary", sme.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.NotNil(data)
	s.Equal(float64(10), data["full_time_males"])
	s.Equal(float64(8), data["full_time_females"])
	s.Equal(float64(3), data["part_time_males"])
	s.Equal(float64(5), data["part_time_females"])
}

func (s *SmeControllerCRUDTestSuite) TestFetchBusinessEmployeeSummaryNotFound() {
	// Create an SME without employee summary
	createdBy := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:                 "USME-REL-008",
		Name:                       "Business Without Employee Data",
		BusinessCategory:           "Services",
		Sector:                     "Consulting",
		ContactPhone:               "+265999888999",
		ContactEmail:               "noemployees@test.mw",
		BusinessImprovementAspects: []string{"Operations"},
		BusinessAccessedFinancing:  []string{"Crowdfunding"},
		CreatedBy:                  &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(sme))

	// Try to fetch employee summary (should return null)
	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/smes/%d/business_employee_summary", sme.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))
	s.Nil(result["data"])
}

func (s *SmeControllerCRUDTestSuite) TestFetchRelationshipInvalidSmeID() {
	// Test with invalid SME ID
	resp, result := s.makeRequest("GET", "/api/smes/99999/primary_business_owner", nil)

	s.Equal(http.StatusNotFound, resp.StatusCode)
	s.False(result["success"].(bool))
}

func (s *SmeControllerCRUDTestSuite) TestFetchRelationshipInvalidIDFormat() {
	// Test with invalid ID format
	resp, result := s.makeRequest("GET", "/api/smes/invalid/primary_business_owner", nil)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
	s.False(result["success"].(bool))
}

// ============================================================================
// JSON SERIALIZATION Tests
// ============================================================================

func (s *SmeControllerCRUDTestSuite) TestJsonFieldsSerialization() {
	// Test creating an SME with multiple improvement aspects and financing options
	improvementAspects := []string{
		"Digital Transformation",
		"Quality Management",
		"Employee Training",
		"Market Expansion",
	}

	financingOptions := []string{
		"Bank Loan",
		"Government Grant",
		"Angel Investment",
		"Crowdfunding",
		"Personal Savings",
	}

	smeData := map[string]interface{}{
		"usme_number":                  "USME-JSON-001",
		"name":                         "JSON Test Business",
		"business_category":            "Technology",
		"sector":                       "IT",
		"contact_phone":                "+265999123001",
		"contact_email":                "json@test.mw",
		"business_improvement_aspects": improvementAspects,
		"business_accessed_financing":  financingOptions,
	}

	// Create the SME
	resp, result := s.makeRequest("POST", "/api/smes", smeData)
	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	smeID := uint(data["id"].(float64))

	// Verify the arrays were returned correctly
	returnedAspects := data["business_improvement_aspects"].([]interface{})
	s.Equal(len(improvementAspects), len(returnedAspects))
	for i, aspect := range improvementAspects {
		s.Equal(aspect, returnedAspects[i].(string))
	}

	returnedFinancing := data["business_accessed_financing"].([]interface{})
	s.Equal(len(financingOptions), len(returnedFinancing))
	for i, option := range financingOptions {
		s.Equal(option, returnedFinancing[i].(string))
	}

	// Fetch the SME to verify persistence
	resp, result = s.makeRequest("GET", fmt.Sprintf("/api/smes/%d", smeID), nil)
	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	fetchedData := result["data"].(map[string]interface{})

	// Verify the fetched arrays match what we saved
	fetchedAspects := fetchedData["business_improvement_aspects"].([]interface{})
	s.Equal(len(improvementAspects), len(fetchedAspects))
	for i, aspect := range improvementAspects {
		s.Equal(aspect, fetchedAspects[i].(string), "Improvement aspect at index %d should match", i)
	}

	fetchedFinancing := fetchedData["business_accessed_financing"].([]interface{})
	s.Equal(len(financingOptions), len(fetchedFinancing))
	for i, option := range financingOptions {
		s.Equal(option, fetchedFinancing[i].(string), "Financing option at index %d should match", i)
	}

	// Verify data was actually stored as JSON in the database
	var sme models.Sme
	s.Nil(facades.Orm().Query().Where("id = ?", smeID).First(&sme))

	// Check the JSON fields are properly populated
	s.NotEmpty(sme.BusinessImprovementAspectJSON)
	s.NotEmpty(sme.BusinessAccessedFinancingJSON)

	// Verify the deserialized arrays match
	s.Equal(len(improvementAspects), len(sme.BusinessImprovementAspects))
	for i, aspect := range improvementAspects {
		s.Equal(aspect, sme.BusinessImprovementAspects[i])
	}

	s.Equal(len(financingOptions), len(sme.BusinessAccessedFinancing))
	for i, option := range financingOptions {
		s.Equal(option, sme.BusinessAccessedFinancing[i])
	}
}

func (s *SmeControllerCRUDTestSuite) TestJsonFieldsUpdate() {
	// Create an SME with initial values
	createdBy := int(s.testUser.ID)

	initialAspects := []string{"Innovation", "Growth"}
	initialFinancing := []string{"Bootstrap", "Friends and Family"}

	sme := &models.Sme{
		UsmeNumber:                 "USME-JSON-UPDATE-001",
		Name:                       "Update Test Business",
		BusinessCategory:           "Services",
		Sector:                     "Consulting",
		ContactPhone:               "+265999123002",
		ContactEmail:               "update@test.mw",
		BusinessImprovementAspects: initialAspects,
		BusinessAccessedFinancing:  initialFinancing,
		CreatedBy:                  &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(sme))

	// Update with new values
	updatedAspects := []string{
		"Process Optimization",
		"Technology Adoption",
		"Customer Service Excellence",
	}

	updatedFinancing := []string{
		"Venture Capital",
		"Bank Loan",
		"Government Subsidy",
	}

	updateData := map[string]interface{}{
		"usme_number":                  sme.UsmeNumber,
		"name":                         sme.Name,
		"business_category":            sme.BusinessCategory,
		"sector":                       sme.Sector,
		"contact_phone":                sme.ContactPhone,
		"contact_email":                sme.ContactEmail,
		"business_improvement_aspects": updatedAspects,
		"business_accessed_financing":  updatedFinancing,
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/smes/%d", sme.ID), updateData)
	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})

	// Verify updated values in response
	returnedAspects := data["business_improvement_aspects"].([]interface{})
	s.Equal(len(updatedAspects), len(returnedAspects))
	for i, aspect := range updatedAspects {
		s.Equal(aspect, returnedAspects[i].(string))
	}

	returnedFinancing := data["business_accessed_financing"].([]interface{})
	s.Equal(len(updatedFinancing), len(returnedFinancing))
	for i, option := range updatedFinancing {
		s.Equal(option, returnedFinancing[i].(string))
	}

	// Fetch again to verify persistence
	var updatedSme models.Sme
	s.Nil(facades.Orm().Query().Where("id = ?", sme.ID).First(&updatedSme))

	s.Equal(len(updatedAspects), len(updatedSme.BusinessImprovementAspects))
	for i, aspect := range updatedAspects {
		s.Equal(aspect, updatedSme.BusinessImprovementAspects[i])
	}

	s.Equal(len(updatedFinancing), len(updatedSme.BusinessAccessedFinancing))
	for i, option := range updatedFinancing {
		s.Equal(option, updatedSme.BusinessAccessedFinancing[i])
	}
}

func (s *SmeControllerCRUDTestSuite) TestJsonFieldsEmptyArrays() {
	// Test that empty arrays stored in DB are returned correctly
	// Create directly in DB to bypass API validation
	createdBy := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:       "USME-JSON-EMPTY-001",
		Name:             "Empty Arrays Test",
		BusinessCategory: "Manufacturing",
		Sector:           "Industrial",
		ContactPhone:     "+265999123003",
		ContactEmail:     "empty@test.mw",
		// Leave JSON arrays empty - the BeforeSave hook will handle empty slices
		BusinessImprovementAspects: []string{},
		BusinessAccessedFinancing:  []string{},
		CreatedBy:                  &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(sme))

	// Fetch via API and verify empty arrays are returned correctly
	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/smes/%d", sme.ID), nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	fetchedData := result["data"].(map[string]interface{})
	fetchedAspects := fetchedData["business_improvement_aspects"].([]interface{})
	s.Equal(0, len(fetchedAspects), "Fetched improvement aspects should be empty array")

	fetchedFinancing := fetchedData["business_accessed_financing"].([]interface{})
	s.Equal(0, len(fetchedFinancing), "Fetched financing options should be empty array")
}

func (s *SmeControllerCRUDTestSuite) TestJsonFieldsSingleItem() {
	// Test that single-item arrays work correctly
	smeData := map[string]interface{}{
		"usme_number":                  "USME-JSON-SINGLE-001",
		"name":                         "Single Item Test",
		"business_category":            "Retail",
		"sector":                       "Commerce",
		"contact_phone":                "+265999123004",
		"contact_email":                "single@test.mw",
		"business_improvement_aspects": []string{"Sustainability"},
		"business_accessed_financing":  []string{"Microfinance"},
	}

	resp, result := s.makeRequest("POST", "/api/smes", smeData)
	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	smeID := uint(data["id"].(float64))

	// Verify single-item arrays
	aspects := data["business_improvement_aspects"].([]interface{})
	s.Equal(1, len(aspects))
	s.Equal("Sustainability", aspects[0].(string))

	financing := data["business_accessed_financing"].([]interface{})
	s.Equal(1, len(financing))
	s.Equal("Microfinance", financing[0].(string))

	// Fetch and verify
	var sme models.Sme
	s.Nil(facades.Orm().Query().Where("id = ?", smeID).First(&sme))

	s.Equal(1, len(sme.BusinessImprovementAspects))
	s.Equal("Sustainability", sme.BusinessImprovementAspects[0])

	s.Equal(1, len(sme.BusinessAccessedFinancing))
	s.Equal("Microfinance", sme.BusinessAccessedFinancing[0])
}

func (s *SmeControllerCRUDTestSuite) TestJsonFieldsSpecialCharacters() {
	// Test that special characters and various string formats are preserved
	specialAspects := []string{
		"Quality & Compliance",
		"R&D Innovation",
		"Customer Service (24/7)",
		"ISO 9001:2015 Certification",
		"Employee Well-being & Safety",
	}

	specialFinancing := []string{
		"SME Development Fund (2023)",
		"COVID-19 Relief Grant",
		"Women's Empowerment Loan @ 5%",
		"Trade & Export Finance",
	}

	smeData := map[string]interface{}{
		"usme_number":                  "USME-JSON-SPECIAL-001",
		"name":                         "Special Characters Test",
		"business_category":            "Technology",
		"sector":                       "IT",
		"contact_phone":                "+265999123005",
		"contact_email":                "special@test.mw",
		"business_improvement_aspects": specialAspects,
		"business_accessed_financing":  specialFinancing,
	}

	resp, result := s.makeRequest("POST", "/api/smes", smeData)
	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	smeID := uint(data["id"].(float64))

	// Fetch and verify special characters are preserved
	resp, result = s.makeRequest("GET", fmt.Sprintf("/api/smes/%d", smeID), nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	fetchedData := result["data"].(map[string]interface{})

	fetchedAspects := fetchedData["business_improvement_aspects"].([]interface{})
	s.Equal(len(specialAspects), len(fetchedAspects))
	for i, aspect := range specialAspects {
		s.Equal(aspect, fetchedAspects[i].(string), "Special character in aspect should be preserved")
	}

	fetchedFinancing := fetchedData["business_accessed_financing"].([]interface{})
	s.Equal(len(specialFinancing), len(fetchedFinancing))
	for i, option := range specialFinancing {
		s.Equal(option, fetchedFinancing[i].(string), "Special character in financing should be preserved")
	}
}

// ============================================================================
// CASCADE DELETE Tests
// ============================================================================

func (s *SmeControllerCRUDTestSuite) TestDeleteSmeCascadesToRelationships() {
	// Create an SME with all relationships
	createdBy := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:                 "USME-CASCADE-001",
		Name:                       "Business With Full Relationships",
		BusinessCategory:           "Technology",
		Sector:                     "IT",
		ContactPhone:               "+265999123006",
		ContactEmail:               "cascade@test.mw",
		BusinessImprovementAspects: []string{"Innovation", "Growth"},
		BusinessAccessedFinancing:  []string{"Bank Loan", "Grant"},
		CreatedBy:                  &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(sme))

	// Create primary business owner
	owner := &models.PrimaryBusinessOwner{
		FirstName:        "Cascade",
		LastName:         "Owner",
		DateOfBirth:      *carbon.NewDateTime(carbon.Parse("1980-01-01")),
		Nationality:      "Malawian",
		NationalIdNumber: "CASCADE001",
		Gender:           "FEMALE",
		EducationLevel:   "Secondary",
		MalawianStatus:   "Citizen",
		PhoneNumber:      "+265991111111",
		SmeID:            int(sme.ID),
		CreatedBy:        &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(owner))

	// Create additional business members
	member1 := &models.AdditionalBusinessMember{
		FirstName:        "Member",
		LastName:         "One",
		Nationality:      "Malawian",
		NationalIdNumber: "MEM001",
		PhoneNumber:      "+265992222222",
		IsIntern:         false,
		IsPartTime:       false,
		SmeId:            int(sme.ID),
		CreatedBy:        &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(member1))

	member2 := &models.AdditionalBusinessMember{
		FirstName:        "Member",
		LastName:         "Two",
		Nationality:      "Malawian",
		NationalIdNumber: "MEM002",
		PhoneNumber:      "+265993333333",
		IsIntern:         true,
		IsPartTime:       true,
		SmeId:            int(sme.ID),
		CreatedBy:        &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(member2))

	// Create business formalisation
	formalisation := &models.BusinessFormalisation{
		HasBankAccount:         true,
		HasTaxClarification:    true,
		IsRegisteredForVat:     false,
		IsMemberOfAssociation:  true,
		IsAffiliated:           false,
		HasExportLicense:       false,
		HasAccessedBds:         true,
		AnnualTurnover:         250000.00,
		EstimatedValueOfAssets: 750000.00,
		FormalisationScore:     85,
		SmeID:                  int(sme.ID),
		CreatedBy:              &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(formalisation))

	// Create business employee summary
	summary := &models.BusinessEmployeeSummary{
		FullTimeMales:   15,
		FullTimeFemales: 12,
		PartTimeMales:   4,
		PartTimeFemales: 6,
		InternMales:     3,
		InternFemales:   5,
		SmeID:           int(sme.ID),
		CreatedBy:       &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(summary))

	// Verify all relationships exist before deletion
	var checkOwner models.PrimaryBusinessOwner
	s.Nil(facades.Orm().Query().Where("sme_id = ?", sme.ID).First(&checkOwner))
	s.Equal(int(sme.ID), checkOwner.SmeID)

	var checkMembers []models.AdditionalBusinessMember
	s.Nil(facades.Orm().Query().Where("sme_id = ?", sme.ID).Find(&checkMembers))
	s.Equal(2, len(checkMembers))

	var checkFormalisation models.BusinessFormalisation
	s.Nil(facades.Orm().Query().Where("sme_id = ?", sme.ID).First(&checkFormalisation))
	s.Equal(int(sme.ID), checkFormalisation.SmeID)

	var checkSummary models.BusinessEmployeeSummary
	s.Nil(facades.Orm().Query().Where("sme_id = ?", sme.ID).First(&checkSummary))
	s.Equal(int(sme.ID), checkSummary.SmeID)

	// Delete the SME
	resp, _ := s.makeRequest("DELETE", fmt.Sprintf("/api/smes/%d", sme.ID), nil)
	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify SME is deleted (soft deleted)
	var deletedSme models.Sme
	err := facades.Orm().Query().Where("id = ?", sme.ID).First(&deletedSme)
	// Should not find the record (soft deleted) - no error means it found nothing
	s.Nil(err, "Query should succeed but find no records")
	s.Equal(uint(0), deletedSme.ID, "SME should be soft deleted (ID should be 0)")

	// Verify all relationships are also deleted (cascade soft delete)
	var deletedOwner models.PrimaryBusinessOwner
	err = facades.Orm().Query().Where("sme_id = ?", sme.ID).First(&deletedOwner)
	s.Nil(err, "Query should succeed but find no records")
	s.Equal(uint(0), deletedOwner.ID, "Primary business owner should be soft deleted when SME is deleted")

	var deletedMembers []models.AdditionalBusinessMember
	err = facades.Orm().Query().Where("sme_id = ?", sme.ID).Find(&deletedMembers)
	s.Nil(err, "Query should succeed")
	s.Equal(0, len(deletedMembers), "Additional business members should be soft deleted when SME is deleted")

	var deletedFormalisation models.BusinessFormalisation
	err = facades.Orm().Query().Where("sme_id = ?", sme.ID).First(&deletedFormalisation)
	s.Nil(err, "Query should succeed but find no records")
	s.Equal(uint(0), deletedFormalisation.ID, "Business formalisation should be soft deleted when SME is deleted")

	var deletedSummary models.BusinessEmployeeSummary
	err = facades.Orm().Query().Where("sme_id = ?", sme.ID).First(&deletedSummary)
	s.Nil(err, "Query should succeed but find no records")
	s.Equal(uint(0), deletedSummary.ID, "Business employee summary should be soft deleted when SME is deleted")
}

// ============================================================================
// UBI (Universal Business Identifier) Generation Tests
// ============================================================================

func (s *SmeControllerCRUDTestSuite) TestUBIGeneration() {
	// Test UBI generation with all required components
	// Don't provide usme_number - it should be auto-generated
	smeData := map[string]interface{}{
		"name":                         "UBI Test Business",
		"business_category":            "Micro",
		"sector":                       "Technology",
		"contact_phone":                "+265999123007",
		"contact_email":                "ubi@test.mw",
		"district":                     "Lilongwe",
		"region":                       "Central",
		"business_improvement_aspects": []string{"Operations"},
		"business_accessed_financing":  []string{"Savings"},
	}

	resp, result := s.makeRequest("POST", "/api/smes", smeData)
	if resp.StatusCode != http.StatusCreated {
		s.T().Logf("Failed to create SME. Status: %d, Response: %+v", resp.StatusCode, result)
	}
	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})

	// Verify UBI was generated as usme_number
	s.NotNil(data["usme_number"], "UBI should be generated as usme_number")
	ubi := data["usme_number"].(string)

	// Verify UBI format: MW-YYYY-DD-CT-NNNNNN-C (DD is 2-3 letter district code, CT is 2-letter category code)
	s.Regexp(`^MW-\d{4}-[A-Z]{2,3}-[A-Z]{2}-\d{6}-\d$`, ubi, "UBI should match expected format")

	// Verify components
	parts := strings.Split(ubi, "-")
	s.Equal(6, len(parts), "UBI should have 6 parts")
	s.Equal("MW", parts[0], "Country code should be MW")

	// Year should be current year
	currentYear := time.Now().Year()
	year, _ := strconv.Atoi(parts[1])
	s.Equal(currentYear, year, "Year should be current year")

	// District code should be LL for Lilongwe
	s.Equal("LL", parts[2], "District code for Lilongwe should be LL")

	// Category should be MI for Micro (first 2 letters of "Micro")
	s.Equal("MI", parts[3], "Category should be MI for Micro business")

	// Sequential number should be 6 digits
	s.Len(parts[4], 6, "Sequential number should be 6 digits")

	// Check digit should be 1 digit
	s.Len(parts[5], 1, "Check digit should be 1 digit")
}

func (s *SmeControllerCRUDTestSuite) TestUBIDistrictCodes() {
	// Test that different districts get correct codes
	testCases := []struct {
		district     string
		expectedCode string
	}{
		{"Lilongwe", "LL"},
		{"Blantyre", "BT"},
		{"Mzuzu", "MZ"},
		{"Zomba", "ZA"},
	}

	for _, tc := range testCases {
		smeData := map[string]interface{}{
			"name":                         fmt.Sprintf("Business in %s", tc.district),
			"business_category":            "Small",
			"sector":                       "Retail",
			"contact_phone":                "+265999123008",
			"contact_email":                fmt.Sprintf("district-%s@test.mw", strings.ToLower(tc.district)),
			"district":                     tc.district,
			"business_improvement_aspects": []string{"Marketing"},
			"business_accessed_financing":  []string{"Bank Loan"},
		}

		resp, result := s.makeRequest("POST", "/api/smes", smeData)
		s.Equal(http.StatusCreated, resp.StatusCode)

		data := result["data"].(map[string]interface{})
		ubi := data["usme_number"].(string)
		parts := strings.Split(ubi, "-")

		s.Equal(tc.expectedCode, parts[2], fmt.Sprintf("District code for %s should be %s", tc.district, tc.expectedCode))
	}
}

func (s *SmeControllerCRUDTestSuite) TestUBIBusinessCategories() {
	// Test that different business categories get correct codes
	// Category codes are now dynamic based on the first 2 letters of the category name
	testCases := []struct {
		category     string
		expectedCode string
	}{
		{"Micro", "MI"},                // First 2 letters: MI
		{"Small", "SM"},                // First 2 letters: SM
		{"Medium", "ME"},               // First 2 letters: ME
		{"Large", "LA"},                // First 2 letters: LA
		{"Micro Enterprise", "ME"},     // First letters of first 2 words: M + E = ME
		{"Small Business", "SB"},       // First letters of first 2 words: S + B = SB
		{"Medium Enterprise", "ME"},    // First letters of first 2 words: M + E = ME
		{"Social Enterprise", "SE"},    // First letters of first 2 words: S + E = SE
	}

	for _, tc := range testCases {
		smeData := map[string]interface{}{
			"name":                         fmt.Sprintf("%s Business", tc.category),
			"business_category":            tc.category,
			"sector":                       "Services",
			"contact_phone":                "+265999123009",
			"contact_email":                fmt.Sprintf("cat-%s@test.mw", strings.ToLower(strings.ReplaceAll(tc.category, " ", ""))),
			"district":                     "Lilongwe",
			"business_improvement_aspects": []string{"Operations"},
			"business_accessed_financing":  []string{"Grants"},
		}

		resp, result := s.makeRequest("POST", "/api/smes", smeData)
		s.Equal(http.StatusCreated, resp.StatusCode)

		data := result["data"].(map[string]interface{})
		ubi := data["usme_number"].(string)
		parts := strings.Split(ubi, "-")

		s.Equal(tc.expectedCode, parts[3], fmt.Sprintf("Category code for %s should be %s", tc.category, tc.expectedCode))
	}
}

func (s *SmeControllerCRUDTestSuite) TestUBISequentialNumbers() {
	// Test that sequential numbers increment correctly
	ubis := make([]string, 3)

	for i := 0; i < 3; i++ {
		smeData := map[string]interface{}{
			"name":                         fmt.Sprintf("Sequential Business %d", i),
			"business_category":            "Micro",
			"sector":                       "Services",
			"contact_phone":                fmt.Sprintf("+265999%06d", i+100),
			"contact_email":                fmt.Sprintf("seq%d@test.mw", i),
			"district":                     "Lilongwe",
			"business_improvement_aspects": []string{"Networking"},
			"business_accessed_financing":  []string{"Loan"},
		}

		resp, result := s.makeRequest("POST", "/api/smes", smeData)
		s.Equal(http.StatusCreated, resp.StatusCode)

		data := result["data"].(map[string]interface{})
		ubis[i] = data["usme_number"].(string)
	}

	// Extract sequential numbers and verify they increment
	for i := 1; i < len(ubis); i++ {
		parts1 := strings.Split(ubis[i-1], "-")
		parts2 := strings.Split(ubis[i], "-")

		seq1, _ := strconv.Atoi(parts1[4])
		seq2, _ := strconv.Atoi(parts2[4])

		s.Greater(seq2, seq1, "Sequential numbers should increment")
	}
}

func (s *SmeControllerCRUDTestSuite) TestUBIUniqueness() {
	// Test that each SME gets a unique UBI
	ubis := make(map[string]bool)

	for i := 0; i < 5; i++ {
		smeData := map[string]interface{}{
			"name":                         fmt.Sprintf("Unique Business %d", i),
			"business_category":            "Small",
			"sector":                       "Manufacturing",
			"contact_phone":                fmt.Sprintf("+265999%06d", i+200),
			"contact_email":                fmt.Sprintf("unique%d@test.mw", i),
			"district":                     "Blantyre",
			"business_improvement_aspects": []string{"Quality management"},
			"business_accessed_financing":  []string{"None"},
		}

		resp, result := s.makeRequest("POST", "/api/smes", smeData)
		s.Equal(http.StatusCreated, resp.StatusCode)

		data := result["data"].(map[string]interface{})
		ubi := data["usme_number"].(string)

		// Verify this UBI hasn't been seen before
		s.False(ubis[ubi], fmt.Sprintf("UBI %s should be unique", ubi))
		ubis[ubi] = true
	}
}
