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

type AdditionalBusinessMemberControllerCRUDTestSuite struct {
	suite.Suite
	tests.TestCase

	server     *httptest.Server
	client     *http.Client
	authCookie *http.Cookie
	testUser   *models.User
	testSme    *models.Sme
}

func TestAdditionalBusinessMemberControllerCRUDTestSuite(t *testing.T) {
	suite.Run(t, &AdditionalBusinessMemberControllerCRUDTestSuite{})
}

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) SetupSuite() {
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

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) SetupTest() {
	// Clean any existing test data first
	if orm := facades.Orm(); orm != nil {
		// Delete in correct order to respect foreign key constraints
		orm.Query().Exec("DELETE FROM additional_business_members")
		orm.Query().Exec("DELETE FROM business_employee_summary")
		orm.Query().Exec("DELETE FROM smes")
		orm.Query().Exec("DELETE FROM user_roles")
		orm.Query().Exec("DELETE FROM role_permissions")
		orm.Query().Exec("DELETE FROM roles WHERE slug LIKE 'abm_%'")
		orm.Query().Exec("DELETE FROM permissions WHERE slug LIKE 'additional_business_members_%'")
		orm.Query().Exec("DELETE FROM users WHERE email LIKE 'abmtest%@example.com'")
	}

	// Refresh database (runs migrations)
	s.RefreshDatabase()

	// Create test SME and user with permissions
	s.setupTestData()
}

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TearDownTest() {
	// Clean up test data
	if orm := facades.Orm(); orm != nil {
		// Delete in correct order
		orm.Query().Exec("DELETE FROM additional_business_members")
		orm.Query().Exec("DELETE FROM business_employee_summary")
		orm.Query().Exec("DELETE FROM smes")
		orm.Query().Exec("DELETE FROM user_roles")
		orm.Query().Exec("DELETE FROM role_permissions")
		orm.Query().Exec("DELETE FROM roles WHERE slug LIKE 'abm_%'")
		orm.Query().Exec("DELETE FROM permissions WHERE slug LIKE 'additional_business_members_%'")
		orm.Query().Exec("DELETE FROM users WHERE email LIKE 'abmtest%@example.com'")
	}
}

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) setupTestData() {
	// Create admin role
	adminRole := &models.Role{
		Name:  "ABM Admin",
		Slug:  "abm_admin",
		Level: 100,
	}
	s.Nil(facades.Orm().Query().Create(adminRole))

	// Create all CRUD permissions
	permissions := []string{"create", "read", "update", "delete", "view"}
	for _, perm := range permissions {
		permission := &models.Permission{
			Name:  fmt.Sprintf("additional_business_members_%s", perm),
			Slug:  fmt.Sprintf("additional_business_members_%s", perm),
			Scope: "by_all",
		}
		s.Nil(facades.Orm().Query().Create(permission))
		s.Nil(helpers.AssignPermissionToRole(adminRole, permission, "by_all"))
	}

	// Create test user
	user, err := helpers.SetupJWTUser("abmtest@example.com", "password", adminRole)
	s.Nil(err)
	s.testUser = user

	// Create test SME
	createdBy := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:                 "USME-ABM-TEST-001",
		Name:                       "Test SME for Additional Members",
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
	s.authCookie = s.loginUser("abmtest@example.com", "password")
	s.NotNil(s.authCookie)
}

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) loginUser(email, password string) *http.Cookie {
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

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) makeRequest(method, path string, body interface{}) (*http.Response, map[string]interface{}) {
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

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) getValidMemberData() map[string]interface{} {
	return map[string]interface{}{
		"first_name":         "John",
		"last_name":          "Doe",
		"other_names":        "William",
		"nationality":        "Malawian",
		"national_id_number": "MWI123456789",
		"date_of_birth":      "1990-05-15",
		"email":              "john.doe@example.com",
		"phone_number":       "+265991234567",
		"is_intern":          false,
		"is_part_time":       false,
		"sme_id":             s.testSme.ID,
	}
}

// ============================================================================
// CREATE Tests
// ============================================================================

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestCreateAdditionalBusinessMember() {
	memberData := s.getValidMemberData()

	resp, result := s.makeRequest("POST", "/api/additional_business_members", memberData)

	// Debug output if not successful
	if resp.StatusCode != http.StatusCreated {
		s.T().Logf("Response status: %d", resp.StatusCode)
		s.T().Logf("Response body: %+v", result)
		if result == nil {
			s.T().Logf("Response was empty or could not be parsed")
		}
	}

	s.Equal(http.StatusCreated, resp.StatusCode)
	if result == nil {
		s.Fail("Response body was empty")
		return
	}
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.NotNil(data["id"])
	s.Equal(memberData["first_name"], data["first_name"])
	s.Equal(memberData["last_name"], data["last_name"])
	s.Equal(memberData["nationality"], data["nationality"])
	s.Equal(memberData["national_id_number"], data["national_id_number"])
	s.Equal(memberData["phone_number"], data["phone_number"])
	s.Equal(float64(s.testSme.ID), data["sme_id"])
	s.Equal(float64(s.testUser.ID), data["created_by"])
}

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestCreateMemberValidationErrors() {
	testCases := []struct {
		name        string
		memberData  map[string]interface{}
		expectError bool
		errorField  string
	}{
		{
			name: "Missing required first_name",
			memberData: map[string]interface{}{
				"last_name":          "Doe",
				"nationality":        "Malawian",
				"national_id_number": "MWI987654321",
				"phone_number":       "+265991234567",
				"sme_id":             s.testSme.ID,
			},
			expectError: true,
			errorField:  "first_name",
		},
		{
			name: "Missing required last_name",
			memberData: map[string]interface{}{
				"first_name":         "John",
				"nationality":        "Malawian",
				"national_id_number": "MWI987654321",
				"phone_number":       "+265991234567",
				"sme_id":             s.testSme.ID,
			},
			expectError: true,
			errorField:  "last_name",
		},
		{
			name: "Invalid email format",
			memberData: map[string]interface{}{
				"first_name":         "John",
				"last_name":          "Doe",
				"nationality":        "Malawian",
				"national_id_number": "MWI987654321",
				"email":              "invalid-email",
				"phone_number":       "+265991234567",
				"sme_id":             s.testSme.ID,
			},
			expectError: true,
			errorField:  "email",
		},
		{
			name: "Invalid date format",
			memberData: map[string]interface{}{
				"first_name":         "John",
				"last_name":          "Doe",
				"nationality":        "Malawian",
				"national_id_number": "MWI987654321",
				"date_of_birth":      "not-a-date",
				"phone_number":       "+265991234567",
				"sme_id":             s.testSme.ID,
			},
			expectError: true,
			errorField:  "date_of_birth",
		},
		{
			name: "Missing SME ID",
			memberData: map[string]interface{}{
				"first_name":         "John",
				"last_name":          "Doe",
				"nationality":        "Malawian",
				"national_id_number": "MWI987654321",
				"phone_number":       "+265991234567",
			},
			expectError: true,
			errorField:  "sme_id",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, result := s.makeRequest("POST", "/api/additional_business_members", tc.memberData)

			if tc.expectError {
				s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
				s.False(result["success"].(bool))
				// Validation errors should be in the response
				s.Contains(strings.ToLower(result["message"].(string)), "validation")
			} else {
				s.Equal(http.StatusCreated, resp.StatusCode)
				s.True(result["success"].(bool))
			}
		})
	}
}

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestCreateMemberWithOptionalFields() {
	memberData := map[string]interface{}{
		"first_name":         "Jane",
		"last_name":          "Smith",
		"nationality":        "Malawian",
		"national_id_number": "MWI111222333",
		"phone_number":       "+265992223333",
		"is_intern":          true,
		"is_part_time":       true,
		"sme_id":             s.testSme.ID,
		// Optional fields not provided
	}

	resp, result := s.makeRequest("POST", "/api/additional_business_members", memberData)

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal(true, data["is_intern"])
	s.Equal(true, data["is_part_time"])
	s.Nil(data["other_names"])
	s.Nil(data["email"])
	s.Nil(data["date_of_birth"])
}

// ============================================================================
// READ Tests
// ============================================================================

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestGetAdditionalBusinessMember() {
	// Create a member first
	createdBy := int(s.testUser.ID)
	dob := carbon.Parse("1985-03-20")
	email := "gettest@example.com"
	otherNames := "Middle"
	member := &models.AdditionalBusinessMember{
		FirstName:        "Get",
		LastName:         "Test",
		OtherNames:       &otherNames,
		Nationality:      "Malawian",
		NationalIdNumber: "GET123456",
		DateOfBirth:      carbon.NewDateTime(dob),
		Email:            &email,
		PhoneNumber:      "+265993334444",
		IsIntern:         false,
		IsPartTime:       true,
		SmeId:            int(s.testSme.ID),
		CreatedBy:        &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(member))

	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/additional_business_members/%d", member.ID), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal(float64(member.ID), data["id"].(float64))
	s.Equal(member.FirstName, data["first_name"])
	s.Equal(member.LastName, data["last_name"])
	s.Equal(*member.OtherNames, data["other_names"])
	s.Equal(member.Nationality, data["nationality"])
	s.Equal(member.NationalIdNumber, data["national_id_number"])
	s.Equal(*member.Email, data["email"])
	s.Equal(member.PhoneNumber, data["phone_number"])
	s.Equal(member.IsIntern, data["is_intern"])
	s.Equal(member.IsPartTime, data["is_part_time"])
}

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestGetMemberNotFound() {
	resp, result := s.makeRequest("GET", "/api/additional_business_members/99999", nil)

	s.Equal(http.StatusNotFound, resp.StatusCode)
	s.False(result["success"].(bool))
}

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestGetMemberInvalidID() {
	resp, result := s.makeRequest("GET", "/api/additional_business_members/invalid", nil)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
	s.False(result["success"].(bool))
}

// ============================================================================
// UPDATE Tests
// ============================================================================

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestUpdateAdditionalBusinessMember() {
	// Create a member first
	createdBy := int(s.testUser.ID)
	member := &models.AdditionalBusinessMember{
		FirstName:        "Original",
		LastName:         "Name",
		Nationality:      "Malawian",
		NationalIdNumber: "UPDATE123",
		PhoneNumber:      "+265994445555",
		IsIntern:         false,
		IsPartTime:       false,
		SmeId:            int(s.testSme.ID),
		CreatedBy:        &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(member))

	updateData := map[string]interface{}{
		"first_name":   "Updated",
		"last_name":    "Member",
		"email":        "updated@example.com",
		"is_intern":    true,
		"is_part_time": true,
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/additional_business_members/%d", member.ID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal("Updated", data["first_name"])
	s.Equal("Member", data["last_name"])
	s.Equal("updated@example.com", data["email"])
	s.Equal(true, data["is_intern"])
	s.Equal(true, data["is_part_time"])
	// Unchanged fields should remain
	s.Equal(member.Nationality, data["nationality"])
	s.Equal(member.NationalIdNumber, data["national_id_number"])
}

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestUpdateMemberPartialUpdate() {
	// Create a member
	createdBy := int(s.testUser.ID)
	email := "partial@example.com"
	member := &models.AdditionalBusinessMember{
		FirstName:        "Partial",
		LastName:         "Update",
		Nationality:      "Malawian",
		NationalIdNumber: "PARTIAL123",
		Email:            &email,
		PhoneNumber:      "+265995556666",
		IsIntern:         false,
		IsPartTime:       false,
		SmeId:            int(s.testSme.ID),
		CreatedBy:        &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(member))

	// Update only email
	updateData := map[string]interface{}{
		"email": "newemail@example.com",
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/additional_business_members/%d", member.ID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.Equal("newemail@example.com", data["email"])
	// Other fields should remain unchanged
	s.Equal(member.FirstName, data["first_name"])
	s.Equal(member.LastName, data["last_name"])
	s.Equal(member.Nationality, data["nationality"])
}

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestUpdateMemberNotFound() {
	updateData := map[string]interface{}{
		"first_name": "Updated",
	}

	resp, result := s.makeRequest("PUT", "/api/additional_business_members/99999", updateData)

	s.Equal(http.StatusNotFound, resp.StatusCode)
	s.False(result["success"].(bool))
}

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestUpdateMemberValidation() {
	// Create a member
	createdBy := int(s.testUser.ID)
	member := &models.AdditionalBusinessMember{
		FirstName:        "Validation",
		LastName:         "Test",
		Nationality:      "Malawian",
		NationalIdNumber: "VAL123",
		PhoneNumber:      "+265996667777",
		IsIntern:         false,
		IsPartTime:       false,
		SmeId:            int(s.testSme.ID),
		CreatedBy:        &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(member))

	// Try to update with invalid email
	updateData := map[string]interface{}{
		"email": "invalid-email",
	}

	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/additional_business_members/%d", member.ID), updateData)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
}

// ============================================================================
// DELETE Tests
// ============================================================================

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestDeleteAdditionalBusinessMember() {
	// Create a member
	createdBy := int(s.testUser.ID)
	member := &models.AdditionalBusinessMember{
		FirstName:        "Delete",
		LastName:         "Me",
		Nationality:      "Malawian",
		NationalIdNumber: "DEL123",
		PhoneNumber:      "+265997778888",
		IsIntern:         false,
		IsPartTime:       false,
		SmeId:            int(s.testSme.ID),
		CreatedBy:        &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(member))

	resp, _ := s.makeRequest("DELETE", fmt.Sprintf("/api/additional_business_members/%d", member.ID), nil)

	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify soft delete
	var deletedMember models.AdditionalBusinessMember
	err := facades.Orm().Query().WithTrashed().Where("id", member.ID).First(&deletedMember)
	s.Nil(err)
	s.NotNil(deletedMember.DeletedAt)
}

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestDeleteMemberNotFound() {
	resp, result := s.makeRequest("DELETE", "/api/additional_business_members/99999", nil)

	s.Equal(http.StatusNotFound, resp.StatusCode)
	s.False(result["success"].(bool))
}

// ============================================================================
// PAGINATION Tests
// ============================================================================

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestPagination() {
	// Create 25 members
	createdBy := int(s.testUser.ID)
	for i := 1; i <= 25; i++ {
		member := &models.AdditionalBusinessMember{
			FirstName:        fmt.Sprintf("Member%02d", i),
			LastName:         fmt.Sprintf("Test%02d", i),
			Nationality:      "Malawian",
			NationalIdNumber: fmt.Sprintf("PAGE%03d", i),
			PhoneNumber:      fmt.Sprintf("+26599%07d", i),
			IsIntern:         i%2 == 0,
			IsPartTime:       i%3 == 0,
			SmeId:            int(s.testSme.ID),
			CreatedBy:        &createdBy,
		}
		s.Nil(facades.Orm().Query().Create(member))
	}

	// Test first page (default page size 20)
	resp, result := s.makeRequest("GET", "/api/additional_business_members?page=1", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	pagination := data["pagination"].(map[string]interface{})

	s.Equal(20, len(items))
	s.Equal(float64(1), pagination["current_page"])
	s.Equal(float64(25), pagination["total"])
	s.Equal(float64(2), pagination["last_page"])

	// Test custom page size
	resp, result = s.makeRequest("GET", "/api/additional_business_members?page=1&pageSize=10", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	pagination = data["pagination"].(map[string]interface{})

	s.Equal(10, len(items))
	s.Equal(float64(3), pagination["last_page"])

	// Test second page
	resp, result = s.makeRequest("GET", "/api/additional_business_members?page=2&pageSize=10", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	pagination = data["pagination"].(map[string]interface{})

	s.Equal(10, len(items))
	s.Equal(float64(2), pagination["current_page"])
}

// ============================================================================
// SORTING Tests
// ============================================================================

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestSorting() {
	// Create members with different names
	createdBy := int(s.testUser.ID)
	names := []struct {
		firstName string
		lastName  string
	}{
		{"Zebra", "Anderson"},
		{"Alpha", "Brown"},
		{"Beta", "Clark"},
	}

	for i, name := range names {
		dob := carbon.Parse(fmt.Sprintf("199%d-01-01", i))
		member := &models.AdditionalBusinessMember{
			FirstName:        name.firstName,
			LastName:         name.lastName,
			Nationality:      "Malawian",
			NationalIdNumber: fmt.Sprintf("SORT%03d", i+1),
			DateOfBirth:      carbon.NewDateTime(dob),
			PhoneNumber:      fmt.Sprintf("+26598%07d", i),
			IsIntern:         false,
			IsPartTime:       false,
			SmeId:            int(s.testSme.ID),
			CreatedBy:        &createdBy,
		}
		s.Nil(facades.Orm().Query().Create(member))
		time.Sleep(10 * time.Millisecond) // Small delay to ensure different created_at times
	}

	// Test sort ascending by first_name
	resp, result := s.makeRequest("GET", "/api/additional_business_members?sort=first_name&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 3)

	// Verify first item is "Alpha"
	firstItem := items[0].(map[string]interface{})
	s.Equal("Alpha", firstItem["first_name"])

	// Test sort descending by last_name
	resp, result = s.makeRequest("GET", "/api/additional_business_members?sort=last_name&direction=desc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})

	// Verify first item is "Clark"
	firstItem = items[0].(map[string]interface{})
	s.Equal("Clark", firstItem["last_name"])

	// Test sort by date_of_birth
	resp, result = s.makeRequest("GET", "/api/additional_business_members?sort=date_of_birth&direction=asc", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 3)

	// Test default sort (created_at desc)
	resp, result = s.makeRequest("GET", "/api/additional_business_members", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	// Most recently created should be first
	firstItem = items[0].(map[string]interface{})
	lastItem := items[len(items)-1].(map[string]interface{})
	s.NotEqual(firstItem["id"], lastItem["id"])
}

// ============================================================================
// SEARCH Tests
// ============================================================================

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestSearch() {
	// Create searchable members
	createdBy := int(s.testUser.ID)
	searchableMembers := []struct {
		firstName        string
		lastName         string
		otherNames       string
		nationalIdNumber string
		email            string
		phoneNumber      string
	}{
		{
			firstName:        "Searchable",
			lastName:         "Person",
			otherNames:       "Middle",
			nationalIdNumber: "SEARCH001",
			email:            "searchable@example.com",
			phoneNumber:      "+265991111111",
		},
		{
			firstName:        "Another",
			lastName:         "Searchable",
			otherNames:       "Test",
			nationalIdNumber: "SEARCH002",
			email:            "another@searchable.com",
			phoneNumber:      "+265992222222",
		},
		{
			firstName:        "NonMatching",
			lastName:         "Member",
			otherNames:       "Other",
			nationalIdNumber: "NOMATCH001",
			email:            "nomatch@example.com",
			phoneNumber:      "+265993333333",
		},
	}

	for _, sm := range searchableMembers {
		member := &models.AdditionalBusinessMember{
			FirstName:        sm.firstName,
			LastName:         sm.lastName,
			OtherNames:       &sm.otherNames,
			Nationality:      "Malawian",
			NationalIdNumber: sm.nationalIdNumber,
			Email:            &sm.email,
			PhoneNumber:      sm.phoneNumber,
			IsIntern:         false,
			IsPartTime:       false,
			SmeId:            int(s.testSme.ID),
			CreatedBy:        &createdBy,
		}
		s.Nil(facades.Orm().Query().Create(member))
	}

	// Search by first name
	resp, result := s.makeRequest("GET", "/api/additional_business_members/search?q=Searchable", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})

	// Should find 2 items with "Searchable"
	s.Equal(2, len(items))

	// Search by email
	resp, result = s.makeRequest("GET", "/api/additional_business_members/search?q=searchable@example.com", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 1)

	// Search by national ID
	resp, result = s.makeRequest("GET", "/api/additional_business_members/search?q=SEARCH001", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 1)

	firstItem := items[0].(map[string]interface{})
	s.Equal("SEARCH001", firstItem["national_id_number"])

	// Search by phone number
	resp, result = s.makeRequest("GET", "/api/additional_business_members/search?q=265991111111", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	s.GreaterOrEqual(len(items), 1)

	// Test minimum query length (should be 2+ characters)
	resp, result = s.makeRequest("GET", "/api/additional_business_members/search?q=a", nil)
	// Should either return empty or error for too short query
	if resp.StatusCode == http.StatusOK {
		data = result["data"].(map[string]interface{})
		items = data["data"].([]interface{})
		s.Equal(0, len(items), "Single character search should return no results")
	}
}

// ============================================================================
// FILTERING Tests
// ============================================================================

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestFiltering() {
	// Create members with different attributes
	createdBy := int(s.testUser.ID)
	members := []struct {
		firstName   string
		nationality string
		isIntern    bool
		isPartTime  bool
	}{
		{"Intern1", "Malawian", true, false},
		{"Intern2", "Malawian", true, true},
		{"PartTime1", "Zimbabwean", false, true},
		{"FullTime1", "Malawian", false, false},
		{"FullTime2", "South African", false, false},
	}

	for _, m := range members {
		member := &models.AdditionalBusinessMember{
			FirstName:        m.firstName,
			LastName:         "FilterTest",
			Nationality:      m.nationality,
			NationalIdNumber: fmt.Sprintf("FLT%s", m.firstName),
			PhoneNumber:      "+265994444444",
			IsIntern:         m.isIntern,
			IsPartTime:       m.isPartTime,
			SmeId:            int(s.testSme.ID),
			CreatedBy:        &createdBy,
		}
		s.Nil(facades.Orm().Query().Create(member))
	}

	// Test filter by nationality
	resp, result := s.makeRequest("GET", "/api/additional_business_members?nationality=Malawian", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})
	// Should find 3 Malawian members
	malawianCount := 0
	for _, item := range items {
		if item.(map[string]interface{})["nationality"] == "Malawian" {
			malawianCount++
		}
	}
	s.GreaterOrEqual(malawianCount, 3)

	// Test filter by is_intern
	resp, result = s.makeRequest("GET", "/api/additional_business_members?is_intern=true", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	// Should find 2 interns
	internCount := 0
	for _, item := range items {
		if item.(map[string]interface{})["is_intern"] == true {
			internCount++
		}
	}
	s.GreaterOrEqual(internCount, 2)

	// Test filter by is_part_time
	resp, result = s.makeRequest("GET", "/api/additional_business_members?is_part_time=true", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	// Should find 2 part-time members
	partTimeCount := 0
	for _, item := range items {
		if item.(map[string]interface{})["is_part_time"] == true {
			partTimeCount++
		}
	}
	s.GreaterOrEqual(partTimeCount, 2)

	// Test combined filters
	resp, result = s.makeRequest("GET", "/api/additional_business_members?is_intern=true&is_part_time=true", nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	// Should find 1 member who is both intern and part-time
	bothCount := 0
	for _, item := range items {
		itemMap := item.(map[string]interface{})
		if itemMap["is_intern"] == true && itemMap["is_part_time"] == true {
			bothCount++
		}
	}
	s.GreaterOrEqual(bothCount, 1)

	// Test filter by sme_id
	resp, result = s.makeRequest("GET", fmt.Sprintf("/api/additional_business_members?sme_id=%d", s.testSme.ID), nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	items = data["data"].([]interface{})
	// All created members should belong to our test SME
	for _, item := range items {
		s.Equal(float64(s.testSme.ID), item.(map[string]interface{})["sme_id"])
	}
}

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestFilterMetadata() {
	// Get filter metadata
	resp, result := s.makeRequest("GET", "/api/additional_business_members/filters", nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].([]interface{})
	s.Greater(len(data), 0, "Should have filter definitions")

	// Check for expected filters
	filterNames := make(map[string]bool)
	for _, filter := range data {
		filterMap := filter.(map[string]interface{})
		filterNames[filterMap["field"].(string)] = true
	}

	// Verify expected filters are present
	expectedFilters := []string{"nationality", "is_intern", "is_part_time", "date_of_birth", "created_at"}
	for _, expected := range expectedFilters {
		s.Contains(filterNames, expected, fmt.Sprintf("Filter %s should be available", expected))
	}
}

// ============================================================================
// LIST/INDEX Tests
// ============================================================================

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestListAdditionalBusinessMembers() {
	// Create a few members
	createdBy := int(s.testUser.ID)
	for i := 1; i <= 3; i++ {
		member := &models.AdditionalBusinessMember{
			FirstName:        fmt.Sprintf("List%d", i),
			LastName:         "Member",
			Nationality:      "Malawian",
			NationalIdNumber: fmt.Sprintf("LIST%03d", i),
			PhoneNumber:      fmt.Sprintf("+26599500%04d", i),
			IsIntern:         false,
			IsPartTime:       false,
			SmeId:            int(s.testSme.ID),
			CreatedBy:        &createdBy,
		}
		s.Nil(facades.Orm().Query().Create(member))
	}

	resp, result := s.makeRequest("GET", "/api/additional_business_members", nil)
	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	items := data["data"].([]interface{})

	s.GreaterOrEqual(len(items), 3)
}

// ============================================================================
// PERMISSION Tests
// ============================================================================

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestUnauthorizedAccess() {
	// Logout first
	s.authCookie = nil

	// Try to create without authentication
	memberData := s.getValidMemberData()
	resp, result := s.makeRequest("POST", "/api/additional_business_members", memberData)

	// Should get unauthorized or forbidden
	s.True(resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden)
	if result != nil {
		s.False(result["success"].(bool))
	}

	// Try to update without authentication
	resp, _ = s.makeRequest("PUT", "/api/additional_business_members/1", memberData)
	s.True(resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden)

	// Try to delete without authentication
	resp, _ = s.makeRequest("DELETE", "/api/additional_business_members/1", nil)
	s.True(resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden)
}

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestReadOnlyPermissions() {
	// Create a user with only read permissions
	readOnlyRole := &models.Role{
		Name:  "ABM Read Only",
		Slug:  "abm_read_only",
		Level: 50,
	}
	s.Nil(facades.Orm().Query().Create(readOnlyRole))

	// Only assign read permission
	readPermission := &models.Permission{
		Name:  "additional_business_members_read",
		Slug:  "additional_business_members_read",
		Scope: "by_all",
	}
	s.Nil(facades.Orm().Query().Create(readPermission))
	s.Nil(helpers.AssignPermissionToRole(readOnlyRole, readPermission, "by_all"))

	// Create read-only user
	readOnlyUser, err := helpers.SetupJWTUser("abmreadonly@example.com", "password", readOnlyRole)
	s.Nil(err)

	// Login as read-only user
	s.authCookie = s.loginUser("abmreadonly@example.com", "password")

	// Should be able to read
	resp, result := s.makeRequest("GET", "/api/additional_business_members", nil)
	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	// Should NOT be able to create
	memberData := s.getValidMemberData()
	resp, result = s.makeRequest("POST", "/api/additional_business_members", memberData)
	s.Equal(http.StatusForbidden, resp.StatusCode)
	s.False(result["success"].(bool))

	// Should NOT be able to update
	resp, result = s.makeRequest("PUT", "/api/additional_business_members/1", memberData)
	s.Equal(http.StatusForbidden, resp.StatusCode)
	s.False(result["success"].(bool))

	// Should NOT be able to delete
	resp, result = s.makeRequest("DELETE", "/api/additional_business_members/1", nil)
	s.Equal(http.StatusForbidden, resp.StatusCode)
	s.False(result["success"].(bool))

	// Cleanup
	facades.Orm().Query().Where("id = ?", readOnlyUser.ID).Delete(&models.User{})
	facades.Orm().Query().Where("id = ?", readOnlyRole.ID).Delete(&models.Role{})
}

// ============================================================================
// AUDIT FIELD Tests
// ============================================================================

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestAuditFields() {
	// Create a member
	memberData := s.getValidMemberData()
	resp, result := s.makeRequest("POST", "/api/additional_business_members", memberData)

	s.Equal(http.StatusCreated, resp.StatusCode)
	data := result["data"].(map[string]interface{})
	memberID := uint(data["id"].(float64))

	// Verify created_by is set
	s.Equal(float64(s.testUser.ID), data["created_by"])
	s.NotNil(data["created_at"])

	// Update the member
	updateData := map[string]interface{}{
		"first_name": "UpdatedForAudit",
	}
	resp, result = s.makeRequest("PUT", fmt.Sprintf("/api/additional_business_members/%d", memberID), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
	data = result["data"].(map[string]interface{})

	// Verify updated_at is set
	s.NotNil(data["updated_at"])
	// Note: updated_by might not be set depending on implementation

	// Delete the member
	resp, _ = s.makeRequest("DELETE", fmt.Sprintf("/api/additional_business_members/%d", memberID), nil)
	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify soft delete and deleted_at is set
	var deletedMember models.AdditionalBusinessMember
	err := facades.Orm().Query().WithTrashed().Where("id", memberID).First(&deletedMember)
	s.Nil(err)
	s.NotNil(deletedMember.DeletedAt)
}

// ============================================================================
// EDGE CASES and ERROR HANDLING Tests
// ============================================================================

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestMaxLengthValidation() {
	// Test exceeding max length for various fields
	longString := make([]byte, 101)
	for i := range longString {
		longString[i] = 'a'
	}

	memberData := s.getValidMemberData()
	memberData["first_name"] = string(longString) // Exceeds max 100

	resp, result := s.makeRequest("POST", "/api/additional_business_members", memberData)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	s.False(result["success"].(bool))
	s.Contains(result["message"].(string), "validation")
}

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestDuplicateNationalID() {
	// Create first member
	memberData := s.getValidMemberData()
	memberData["national_id_number"] = "UNIQUE123"

	resp, _ := s.makeRequest("POST", "/api/additional_business_members", memberData)
	s.Equal(http.StatusCreated, resp.StatusCode)

	// Try to create another member with same national ID (if unique constraint exists)
	memberData2 := s.getValidMemberData()
	memberData2["first_name"] = "Different"
	memberData2["last_name"] = "Person"
	memberData2["national_id_number"] = "UNIQUE123" // Same ID

	resp, result := s.makeRequest("POST", "/api/additional_business_members", memberData2)

	// Depending on implementation, this might succeed or fail
	// If there's a unique constraint, it should fail
	if resp.StatusCode != http.StatusCreated {
		s.False(result["success"].(bool))
	}
}

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestInvalidSMEID() {
	memberData := s.getValidMemberData()
	memberData["sme_id"] = 99999 // Non-existent SME

	resp, result := s.makeRequest("POST", "/api/additional_business_members", memberData)

	// Should fail if foreign key constraint is enforced
	if resp.StatusCode != http.StatusCreated {
		s.False(result["success"].(bool))
	}
}

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestConcurrentUpdates() {
	// Create a member
	createdBy := int(s.testUser.ID)
	member := &models.AdditionalBusinessMember{
		FirstName:        "Concurrent",
		LastName:         "Test",
		Nationality:      "Malawian",
		NationalIdNumber: "CONC123",
		PhoneNumber:      "+265998889999",
		IsIntern:         false,
		IsPartTime:       false,
		SmeId:            int(s.testSme.ID),
		CreatedBy:        &createdBy,
	}
	s.Nil(facades.Orm().Query().Create(member))

	// Simulate concurrent updates
	update1 := map[string]interface{}{"first_name": "Update1"}
	update2 := map[string]interface{}{"last_name": "Update2"}

	// Both should succeed (last write wins)
	resp1, _ := s.makeRequest("PUT", fmt.Sprintf("/api/additional_business_members/%d", member.ID), update1)
	resp2, _ := s.makeRequest("PUT", fmt.Sprintf("/api/additional_business_members/%d", member.ID), update2)

	s.Equal(http.StatusOK, resp1.StatusCode)
	s.Equal(http.StatusOK, resp2.StatusCode)

	// Verify final state
	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/additional_business_members/%d", member.ID), nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data := result["data"].(map[string]interface{})
	// Should have both updates applied
	s.Equal("Update1", data["first_name"])
	s.Equal("Update2", data["last_name"])
}

// ============================================================================
// SPECIAL CHARACTER HANDLING Tests
// ============================================================================

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestSpecialCharactersInFields() {
	memberData := map[string]interface{}{
		"first_name":         "Jean-François",
		"last_name":          "O'Brien",
		"other_names":        "María José",
		"nationality":        "Côte d'Ivoire",
		"national_id_number": "SP#123-456/789",
		"email":              "test+special@example.com",
		"phone_number":       "+265-99-123-4567",
		"is_intern":          false,
		"is_part_time":       false,
		"sme_id":             s.testSme.ID,
	}

	resp, result := s.makeRequest("POST", "/api/additional_business_members", memberData)

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	memberID := uint(data["id"].(float64))

	// Verify special characters are preserved
	s.Equal("Jean-François", data["first_name"])
	s.Equal("O'Brien", data["last_name"])
	s.Equal("María José", data["other_names"])

	// Fetch and verify again
	resp, result = s.makeRequest("GET", fmt.Sprintf("/api/additional_business_members/%d", memberID), nil)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	s.Equal("Jean-François", data["first_name"])
	s.Equal("O'Brien", data["last_name"])
}

// ============================================================================
// NULL VALUE HANDLING Tests
// ============================================================================

func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestNullableFields() {
	// Create member with minimal required fields
	memberData := map[string]interface{}{
		"first_name":         "Minimal",
		"last_name":          "Member",
		"nationality":        "Malawian",
		"national_id_number": "MIN123",
		"phone_number":       "+265991234567",
		"is_intern":          false,
		"is_part_time":       false,
		"sme_id":             s.testSme.ID,
		// Nullable fields not provided
	}

	resp, result := s.makeRequest("POST", "/api/additional_business_members", memberData)

	s.Equal(http.StatusCreated, resp.StatusCode)
	data := result["data"].(map[string]interface{})

	// Verify nullable fields are null
	s.Nil(data["other_names"])
	s.Nil(data["date_of_birth"])
	s.Nil(data["email"])

	memberID := uint(data["id"].(float64))

	// Update to add nullable fields
	updateData := map[string]interface{}{
		"other_names":   "Added",
		"date_of_birth": "1990-01-01",
		"email":         "added@example.com",
	}

	resp, result = s.makeRequest("PUT", fmt.Sprintf("/api/additional_business_members/%d", memberID), updateData)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	s.Equal("Added", data["other_names"])
	s.NotNil(data["date_of_birth"])
	s.Equal("added@example.com", data["email"])

	// Update to remove nullable fields (set to null)
	updateData = map[string]interface{}{
		"other_names": nil,
		"email":       nil,
	}

	resp, result = s.makeRequest("PUT", fmt.Sprintf("/api/additional_business_members/%d", memberID), updateData)
	s.Equal(http.StatusOK, resp.StatusCode)

	data = result["data"].(map[string]interface{})
	s.Nil(data["other_names"])
	s.Nil(data["email"])
}

// TestBusinessEmployeeSummaryAggregates verifies that BusinessEmployeeSummary aggregates
// are correctly updated when additional business members are created, updated, or deleted
func (s *AdditionalBusinessMemberControllerCRUDTestSuite) TestBusinessEmployeeSummaryAggregates() {
	var err error

	// Create first member: full-time male
	member1Data := map[string]interface{}{
		"first_name":         "John",
		"last_name":          "Doe",
		"gender":             "MALE",
		"nationality":        "Malawian",
		"national_id_number": "FTM001",
		"phone_number":       "+265991111111",
		"is_intern":          false,
		"is_part_time":       false,
		"sme_id":             s.testSme.ID,
	}

	resp, result := s.makeRequest("POST", "/api/additional_business_members", member1Data)
	s.Equal(http.StatusCreated, resp.StatusCode)
	member1ID := uint(result["data"].(map[string]interface{})["id"].(float64))

	// Wait for async event processing
	time.Sleep(500 * time.Millisecond)

	// Verify summary was created with 1 full-time male
	var summary models.BusinessEmployeeSummary
	err = facades.Orm().Query().Where("sme_id = ?", s.testSme.ID).First(&summary)
	s.NoError(err)
	s.Equal(1, summary.FullTimeMales)
	s.Equal(0, summary.FullTimeFemales)
	s.Equal(0, summary.PartTimeMales)
	s.Equal(0, summary.PartTimeFemales)
	s.Equal(0, summary.InternMales)
	s.Equal(0, summary.InternFemales)

	// Create second member: full-time female
	member2Data := map[string]interface{}{
		"first_name":         "Jane",
		"last_name":          "Smith",
		"gender":             "FEMALE",
		"nationality":        "Malawian",
		"national_id_number": "FTF001",
		"phone_number":       "+265991111112",
		"is_intern":          false,
		"is_part_time":       false,
		"sme_id":             s.testSme.ID,
	}

	resp, result = s.makeRequest("POST", "/api/additional_business_members", member2Data)
	s.Equal(http.StatusCreated, resp.StatusCode)

	// Wait for async event processing
	time.Sleep(500 * time.Millisecond)

	// Verify summary updated with 1 full-time male and 1 full-time female
	err = facades.Orm().Query().Where("sme_id = ?", s.testSme.ID).First(&summary)
	s.NoError(err)
	s.Equal(1, summary.FullTimeMales)
	s.Equal(1, summary.FullTimeFemales)

	// Create third member: part-time male
	member3Data := map[string]interface{}{
		"first_name":         "Bob",
		"last_name":          "Johnson",
		"gender":             "MALE",
		"nationality":        "Malawian",
		"national_id_number": "PTM001",
		"phone_number":       "+265991111113",
		"is_intern":          false,
		"is_part_time":       true,
		"sme_id":             s.testSme.ID,
	}

	resp, result = s.makeRequest("POST", "/api/additional_business_members", member3Data)
	s.Equal(http.StatusCreated, resp.StatusCode)

	// Wait for async event processing
	time.Sleep(500 * time.Millisecond)

	// Verify summary updated with part-time male
	err = facades.Orm().Query().Where("sme_id = ?", s.testSme.ID).First(&summary)
	s.NoError(err)
	s.Equal(1, summary.FullTimeMales)
	s.Equal(1, summary.FullTimeFemales)
	s.Equal(1, summary.PartTimeMales)
	s.Equal(0, summary.PartTimeFemales)

	// Create fourth member: intern female
	member4Data := map[string]interface{}{
		"first_name":         "Alice",
		"last_name":          "Williams",
		"gender":             "FEMALE",
		"nationality":        "Malawian",
		"national_id_number": "INF001",
		"phone_number":       "+265991111114",
		"is_intern":          true,
		"is_part_time":       false,
		"sme_id":             s.testSme.ID,
	}

	resp, result = s.makeRequest("POST", "/api/additional_business_members", member4Data)
	s.Equal(http.StatusCreated, resp.StatusCode)
	member4ID := uint(result["data"].(map[string]interface{})["id"].(float64))

	// Wait for async event processing
	time.Sleep(500 * time.Millisecond)

	// Verify summary updated with intern female
	err = facades.Orm().Query().Where("sme_id = ?", s.testSme.ID).First(&summary)
	s.NoError(err)
	s.Equal(1, summary.FullTimeMales)
	s.Equal(1, summary.FullTimeFemales)
	s.Equal(1, summary.PartTimeMales)
	s.Equal(0, summary.PartTimeFemales)
	s.Equal(0, summary.InternMales)
	s.Equal(1, summary.InternFemales)

	// Update first member from full-time to part-time
	updateData := map[string]interface{}{
		"is_part_time": true,
	}

	resp, result = s.makeRequest("PUT", fmt.Sprintf("/api/additional_business_members/%d", member1ID), updateData)
	s.Equal(http.StatusOK, resp.StatusCode)

	// Wait for async event processing
	time.Sleep(500 * time.Millisecond)

	// Verify summary updated: full-time male decreased, part-time male increased
	err = facades.Orm().Query().Where("sme_id = ?", s.testSme.ID).First(&summary)
	s.NoError(err)
	s.Equal(0, summary.FullTimeMales) // Decreased from 1 to 0
	s.Equal(1, summary.FullTimeFemales)
	s.Equal(2, summary.PartTimeMales) // Increased from 1 to 2
	s.Equal(0, summary.PartTimeFemales)
	s.Equal(0, summary.InternMales)
	s.Equal(1, summary.InternFemales)

	// Delete the intern female
	resp, _ = s.makeRequest("DELETE", fmt.Sprintf("/api/additional_business_members/%d", member4ID), nil)
	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Wait longer for async event processing after delete
	time.Sleep(2 * time.Second)

	// Verify summary updated: intern female decreased
	err = facades.Orm().Query().Where("sme_id = ?", s.testSme.ID).First(&summary)
	s.NoError(err)
	s.Equal(0, summary.FullTimeMales)
	s.Equal(1, summary.FullTimeFemales)
	s.Equal(2, summary.PartTimeMales)
	s.Equal(0, summary.PartTimeFemales)
	s.Equal(0, summary.InternMales)
	s.Equal(0, summary.InternFemales) // Decreased from 1 to 0
}