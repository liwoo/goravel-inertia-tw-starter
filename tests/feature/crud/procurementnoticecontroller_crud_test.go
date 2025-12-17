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
	"github.com/goravel/framework/support/carbon"
	"github.com/stretchr/testify/suite"

	"smedi-sme-db/app/models"
	"smedi-sme-db/tests"
	"smedi-sme-db/tests/helpers"
)

type ProcurementNoticeControllerCRUDTestSuite struct {
	suite.Suite
	tests.TestCase

	server     *httptest.Server
	client     *http.Client
	authCookie *http.Cookie
	testUser   *models.User
}

func TestProcurementNoticeControllerCRUDTestSuite(t *testing.T) {
	suite.Run(t, &ProcurementNoticeControllerCRUDTestSuite{})
}

func (s *ProcurementNoticeControllerCRUDTestSuite) SetupSuite() {
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

func (s *ProcurementNoticeControllerCRUDTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}

func (s *ProcurementNoticeControllerCRUDTestSuite) SetupTest() {
	// Clean database
	s.RefreshDatabase()

	// Create test user with permissions
	s.setupTestUser()
}

func (s *ProcurementNoticeControllerCRUDTestSuite) TearDownTest() {
	// Clean up test data
	if orm := facades.Orm(); orm != nil {
		orm.Query().Exec("DELETE FROM procurement_notices")
		orm.Query().Exec("DELETE FROM users WHERE email = 'procurementnoticetest@example.com'")
		orm.Query().Exec("DELETE FROM user_roles")
		orm.Query().Exec("DELETE FROM role_permissions")
		orm.Query().Exec("DELETE FROM roles WHERE slug = 'procurement_notice_admin'")
		orm.Query().Exec("DELETE FROM permissions WHERE slug LIKE 'procurement_notices_%%'")
	}
}

func (s *ProcurementNoticeControllerCRUDTestSuite) setupTestUser() {
	// Create admin role
	adminRole := &models.Role{
		Name:  "Procurement Notice Admin",
		Slug:  "procurement_notice_admin",
		Level: 100,
	}
	s.Nil(facades.Orm().Query().Create(adminRole))

	// Create all CRUD permissions
	permissions := []string{"create", "read", "update", "delete"}
	for _, perm := range permissions {
		permission := &models.Permission{
			Name:  fmt.Sprintf("procurement_notices_%s", perm),
			Slug:  fmt.Sprintf("procurement_notices_%s", perm),
			Scope: "by_all",
		}
		s.Nil(facades.Orm().Query().Create(permission))
		s.Nil(helpers.AssignPermissionToRole(adminRole, permission, "by_all"))
	}

	// Create test user
	user, err := helpers.SetupJWTUser("procurementnoticetest@example.com", "password", adminRole)
	s.Nil(err)
	s.testUser = user

	// Login
	s.authCookie = s.loginUser("procurementnoticetest@example.com", "password")
	s.NotNil(s.authCookie)
}

func (s *ProcurementNoticeControllerCRUDTestSuite) loginUser(email, password string) *http.Cookie {
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

func (s *ProcurementNoticeControllerCRUDTestSuite) makeRequest(method, path string, body interface{}) (*http.Response, map[string]interface{}) {
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

func (s *ProcurementNoticeControllerCRUDTestSuite) TestCreateProcurementNotice() {
	procurementNoticeData := map[string]interface{}{
		"procured_by":              "Ministry of Health",
		"procurement_type":         "Goods",
		"market_approach":          "National",
		"invitation":               "open",
		"ref_no":                   "MOH/2024/001",
		"open_date":                carbon.Now().Format("2006-01-02"),
		"close_date":               carbon.Now().AddDays(30).Format("2006-01-02"),
		"partners":                 []string{},
		"qualifying_districts":     []string{},
		"is_published":             true,
		"organization":             "Ministry of Health",
		"classification":           []string{},
		"interested_smes":          []string{},
		"details":                  "Supply of medical equipment",
		"application_details":      "Submit applications via email",
		"minimum_qualifying_score": 70,
	}

	resp, result := s.makeRequest("POST", "/api/procurement-notices", procurementNoticeData)

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.NotNil(data["id"])
	s.Equal(float64(s.testUser.ID), data["created_by"])
}

func (s *ProcurementNoticeControllerCRUDTestSuite) TestCreateProcurementNoticeValidation() {
	// Missing required fields
	procurementNoticeData := map[string]interface{}{}

	resp, result := s.makeRequest("POST", "/api/procurement-notices", procurementNoticeData)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
	s.Contains(strings.ToLower(result["message"].(string)), "validation")
}

func (s *ProcurementNoticeControllerCRUDTestSuite) TestCreateProcurementNoticeWithArrays() {
	procurementNoticeData := map[string]interface{}{
		"procured_by":              "Ministry of Education",
		"procurement_type":         "Services",
		"market_approach":          "International",
		"invitation":               "limited",
		"ref_no":                   "MOE/2024/002",
		"open_date":                carbon.Now().Format("2006-01-02"),
		"close_date":               carbon.Now().AddDays(45).Format("2006-01-02"),
		"partners":                 []string{"UNICEF", "World Bank"},
		"qualifying_districts":     []string{"Lilongwe", "Blantyre", "Mzuzu"},
		"is_published":             true,
		"organization":             "Ministry of Education",
		"classification":           []string{"Education", "Technology"},
		"interested_smes":          []string{},
		"details":                  "Educational technology services",
		"application_details":      "Apply through our portal",
		"minimum_qualifying_score": 80,
	}

	resp, result := s.makeRequest("POST", "/api/procurement-notices", procurementNoticeData)

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.NotNil(data["id"])

	// Check arrays are properly saved
	s.Equal(2, len(data["partners"].([]interface{})))
	s.Equal(3, len(data["qualifying_districts"].([]interface{})))
	s.Equal(2, len(data["classification"].([]interface{})))
}

// ============================================================================
// READ Tests
// ============================================================================

func (s *ProcurementNoticeControllerCRUDTestSuite) TestGetProcurementNotice() {
	// Create a procurement notice first
	userID := int(s.testUser.ID)
	procurementNotice := &models.ProcurementNotice{
		ProcuredBy:             "Test Ministry",
		ProcurementType:        "Goods",
		MarketApproach:         "National",
		Invitation:             "open",
		RefNo:                  "TEST/2024/001",
		OpenDate:               *carbon.NewDateTime(carbon.Now()),
		CloseDate:              *carbon.NewDateTime(carbon.Now().AddDays(30)),
		Partners:               []string{},
		QualifyingDistricts:    []string{},
		IsPublished:            true,
		Organization:           "Test Organization",
		Classification:         []string{},
		InterestedSmes:         []string{},
		Details:                "Test details",
		ApplicationDetails:     "Test application details",
		MinimumQualifyingScore: 60,
	}
	procurementNotice.CreatedBy = &userID
	s.Nil(facades.Orm().Query().Create(procurementNotice))

	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/procurement-notices/%d", procurementNotice.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal(float64(procurementNotice.ID), data["id"].(float64))
	s.Equal("TEST/2024/001", data["ref_no"].(string))
}

// ============================================================================
// UPDATE Tests
// ============================================================================

func (s *ProcurementNoticeControllerCRUDTestSuite) TestUpdateProcurementNotice() {
	// Create a procurement notice first
	userID := int(s.testUser.ID)
	procurementNotice := &models.ProcurementNotice{
		ProcuredBy:             "Original Ministry",
		ProcurementType:        "Goods",
		MarketApproach:         "National",
		Invitation:             "open",
		RefNo:                  "ORIG/2024/001",
		OpenDate:               *carbon.NewDateTime(carbon.Now()),
		CloseDate:              *carbon.NewDateTime(carbon.Now().AddDays(30)),
		Partners:               []string{},
		QualifyingDistricts:    []string{},
		IsPublished:            false,
		Organization:           "Original Organization",
		Classification:         []string{},
		InterestedSmes:         []string{},
		Details:                "Original details",
		ApplicationDetails:     "Original application details",
		MinimumQualifyingScore: 50,
	}
	procurementNotice.CreatedBy = &userID
	s.Nil(facades.Orm().Query().Create(procurementNotice))

	updateData := map[string]interface{}{
		"procured_by":              "Updated Ministry",
		"is_published":             true,
		"minimum_qualifying_score": 75,
		"partners":                 []string{"Partner A", "Partner B"},
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/procurement-notices/%d", procurementNotice.ID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	// Verify update
	var updated models.ProcurementNotice
	s.Nil(facades.Orm().Query().Find(&updated, procurementNotice.ID))
	s.Equal("Updated Ministry", updated.ProcuredBy)
	s.True(updated.IsPublished)
	s.Equal(75, updated.MinimumQualifyingScore)
	s.Equal(2, len(updated.Partners))
}

// ============================================================================
// DELETE Tests
// ============================================================================

func (s *ProcurementNoticeControllerCRUDTestSuite) TestDeleteProcurementNotice() {
	// Create a procurement notice first
	userID := int(s.testUser.ID)
	procurementNotice := &models.ProcurementNotice{
		ProcuredBy:             "Delete Test Ministry",
		ProcurementType:        "Services",
		MarketApproach:         "International",
		Invitation:             "single-source",
		RefNo:                  "DEL/2024/001",
		OpenDate:               *carbon.NewDateTime(carbon.Now()),
		CloseDate:              *carbon.NewDateTime(carbon.Now().AddDays(15)),
		Partners:               []string{},
		QualifyingDistricts:    []string{},
		IsPublished:            false,
		Organization:           "Delete Test Organization",
		Classification:         []string{},
		InterestedSmes:         []string{},
		Details:                "To be deleted",
		ApplicationDetails:     "To be deleted",
		MinimumQualifyingScore: 40,
	}
	procurementNotice.CreatedBy = &userID
	s.Nil(facades.Orm().Query().Create(procurementNotice))

	resp, _ := s.makeRequest("DELETE", fmt.Sprintf("/api/procurement-notices/%d", procurementNotice.ID), nil)

	// DELETE returns 204 No Content on success
	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify soft delete
	var deleted models.ProcurementNotice
	err := facades.Orm().Query().WithTrashed().Find(&deleted, procurementNotice.ID)
	s.Nil(err)
	s.NotNil(deleted.DeletedAt)
}

// ============================================================================
// PAGINATION Tests
// ============================================================================

func (s *ProcurementNoticeControllerCRUDTestSuite) TestPagination() {
	// Create multiple procurement notices
	userID := int(s.testUser.ID)
	for i := 1; i <= 15; i++ {
		procurementNotice := &models.ProcurementNotice{
			ProcuredBy:             fmt.Sprintf("Ministry %d", i),
			ProcurementType:        "Goods",
			MarketApproach:         "National",
			Invitation:             "open",
			RefNo:                  fmt.Sprintf("REF/2024/%03d", i),
			OpenDate:               *carbon.NewDateTime(carbon.Now()),
			CloseDate:              *carbon.NewDateTime(carbon.Now().AddDays(30)),
			Partners:               []string{},
			QualifyingDistricts:    []string{},
			IsPublished:            true,
			Organization:           fmt.Sprintf("Organization %d", i),
			Classification:         []string{},
			InterestedSmes:         []string{},
			Details:                fmt.Sprintf("Details for item %d", i),
			ApplicationDetails:     "Standard application",
			MinimumQualifyingScore: 50 + i,
		}
		procurementNotice.CreatedBy = &userID
		s.Nil(facades.Orm().Query().Create(procurementNotice))
	}

	resp, result := s.makeRequest("GET", "/api/procurement-notices?page=1&pageSize=10", nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	// The response format has data nested inside result["data"]
	dataWrapper := result["data"].(map[string]interface{})
	data := dataWrapper["data"].([]interface{})
	s.Equal(10, len(data))

	pagination := dataWrapper["pagination"].(map[string]interface{})
	s.Equal(float64(15), pagination["total"])
	s.Equal(float64(1), pagination["current_page"])
	s.Equal(float64(10), pagination["per_page"])
}

// ============================================================================
// SORTING Tests
// ============================================================================

func (s *ProcurementNoticeControllerCRUDTestSuite) TestSorting() {
	// Create procurement notices with different organizations
	userID := int(s.testUser.ID)
	organizations := []string{"Zebra Corp", "Alpha Ministry", "Beta Agency"}
	for i, org := range organizations {
		procurementNotice := &models.ProcurementNotice{
			ProcuredBy:             org,
			ProcurementType:        "Services",
			MarketApproach:         "National",
			Invitation:             "open",
			RefNo:                  fmt.Sprintf("SORT/2024/%03d", i+1),
			OpenDate:               *carbon.NewDateTime(carbon.Now()),
			CloseDate:              *carbon.NewDateTime(carbon.Now().AddDays(30)),
			Partners:               []string{},
			QualifyingDistricts:    []string{},
			IsPublished:            true,
			Organization:           org,
			Classification:         []string{},
			InterestedSmes:         []string{},
			Details:                "Sorting test",
			ApplicationDetails:     "Standard",
			MinimumQualifyingScore: 60,
		}
		procurementNotice.CreatedBy = &userID
		s.Nil(facades.Orm().Query().Create(procurementNotice))
	}

	// Test ascending sort
	resp, result := s.makeRequest("GET", "/api/procurement-notices?sort=organization&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)
	dataWrapper := result["data"].(map[string]interface{})
	data := dataWrapper["data"].([]interface{})
	s.Equal("Alpha Ministry", data[0].(map[string]interface{})["organization"])

	// Test descending sort
	resp, result = s.makeRequest("GET", "/api/procurement-notices?sort=organization&direction=desc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)
	dataWrapper = result["data"].(map[string]interface{})
	data = dataWrapper["data"].([]interface{})
	s.Equal("Zebra Corp", data[0].(map[string]interface{})["organization"])
}

// ============================================================================
// SEARCH Tests
// ============================================================================

func (s *ProcurementNoticeControllerCRUDTestSuite) TestSearch() {
	// Create procurement notices with different details
	userID := int(s.testUser.ID)
	testData := []struct {
		refNo   string
		org     string
		details string
	}{
		{"SEARCH/2024/001", "Health Ministry", "Medical supplies tender"},
		{"SEARCH/2024/002", "Education Department", "School furniture procurement"},
		{"UNIQUE/2024/003", "Agriculture Board", "Farming equipment purchase"},
	}

	for _, td := range testData {
		procurementNotice := &models.ProcurementNotice{
			ProcuredBy:             td.org,
			ProcurementType:        "Goods",
			MarketApproach:         "National",
			Invitation:             "open",
			RefNo:                  td.refNo,
			OpenDate:               *carbon.NewDateTime(carbon.Now()),
			CloseDate:              *carbon.NewDateTime(carbon.Now().AddDays(30)),
			Partners:               []string{},
			QualifyingDistricts:    []string{},
			IsPublished:            true,
			Organization:           td.org,
			Classification:         []string{},
			InterestedSmes:         []string{},
			Details:                td.details,
			ApplicationDetails:     "Standard application",
			MinimumQualifyingScore: 60,
		}
		procurementNotice.CreatedBy = &userID
		s.Nil(facades.Orm().Query().Create(procurementNotice))
	}

	// Search by ref_no
	resp, result := s.makeRequest("GET", "/api/procurement-notices?search=UNIQUE", nil)
	s.Equal(http.StatusOK, resp.StatusCode)
	dataWrapper := result["data"].(map[string]interface{})
	data := dataWrapper["data"].([]interface{})
	s.Equal(1, len(data))
	s.Equal("UNIQUE/2024/003", data[0].(map[string]interface{})["ref_no"])

	// Search by organization
	resp, result = s.makeRequest("GET", "/api/procurement-notices?search=Education", nil)
	s.Equal(http.StatusOK, resp.StatusCode)
	dataWrapper = result["data"].(map[string]interface{})
	data = dataWrapper["data"].([]interface{})
	s.Equal(1, len(data))
	s.Equal("Education Department", data[0].(map[string]interface{})["organization"])

	// Search by details
	resp, result = s.makeRequest("GET", "/api/procurement-notices?search=Medical", nil)
	s.Equal(http.StatusOK, resp.StatusCode)
	dataWrapper = result["data"].(map[string]interface{})
	data = dataWrapper["data"].([]interface{})
	s.Equal(1, len(data))
	s.Contains(data[0].(map[string]interface{})["details"].(string), "Medical")
}