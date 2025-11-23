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
	s.RefreshDatabase()
	s.setupTestUser()
}

func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TearDownTest() {
	if orm := facades.Orm(); orm != nil {
		orm.Query().Exec("DELETE FROM primary_business_owner")
		orm.Query().Exec("DELETE FROM smes")
		orm.Query().Exec("DELETE FROM users WHERE email = 'pbotest@example.com'")
		orm.Query().Exec("DELETE FROM user_roles")
		orm.Query().Exec("DELETE FROM role_permissions")
		orm.Query().Exec("DELETE FROM roles WHERE slug = 'pbo_admin'")
		orm.Query().Exec("DELETE FROM permissions WHERE slug LIKE 'pbo_%'")
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
		"national_id_number":       "MWK12345678",
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
		"registration_number":          "COMP12345",
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

//
//func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TestCreatePrimaryBusinessOwnerValidation() {
//	// Test required fields
//	data := s.getValidPBOData()
//	data["first_name"] = "" // Make first name empty
//
//	resp, result := s.makeRequest("POST", "/api/primary-business-owners", data)
//
//	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
//	s.False(result["success"].(bool))
//	s.Contains(result["message"], "validation failed")
//	s.Contains(result["errors"].(map[string]interface{}), "first_name")
//}
//
//// ============================================================================
//// READ Tests
//// ============================================================================
//
//func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TestGetPrimaryBusinessOwner() {
//	// First create a PBO
//	data := s.getValidPBOData()
//	_, createResult := s.makeRequest("POST", "/api/primary-business-owners", data)
//	pboID := createResult["data"].(map[string]interface{})["id"]
//
//	// Now try to retrieve it
//	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/primary-business-owners/%s", pboID), nil)
//
//	s.Equal(http.StatusOK, resp.StatusCode)
//	s.True(result["success"].(bool))
//	s.Equal(pboID, result["data"].(map[string]interface{})["id"])
//	s.Equal(data["first_name"], result["data"].(map[string]interface{})["first_name"])
//}
//
//// ============================================================================
//// UPDATE Tests
//// ============================================================================
//
//func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TestUpdatePrimaryBusinessOwner() {
//	// First create a PBO
//	data := s.getValidPBOData()
//	_, createResult := s.makeRequest("POST", "/api/primary-business-owners", data)
//	pboID := createResult["data"].(map[string]interface{})["id"]
//
//	// Update the PBO
//	updateData := map[string]interface{}{
//		"first_name": "UpdatedFirstName",
//		"last_name":  "UpdatedLastName",
//		"email":      "updated.email@example.com",
//	}
//
//	resp, result := s.makeRequest("PUT", fmt.Sprintf("/api/primary-business-owners/%s", pboID), updateData)
//
//	s.Equal(http.StatusOK, resp.StatusCode)
//	s.True(result["success"].(bool))
//
//	// Verify the update
//	var pbo models.PrimaryBusinessOwner
//	s.Nil(facades.Orm().Query().Where("id = ?", pboID).First(&pbo))
//	s.Equal(updateData["first_name"], pbo.FirstName)
//	s.Equal(updateData["last_name"], pbo.LastName)
//	s.Equal(updateData["email"], *pbo.Email)
//}
//
//// ============================================================================
//// DELETE Tests
//// ============================================================================
//
//func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TestDeletePrimaryBusinessOwner() {
//	// First create a PBO
//	data := s.getValidPBOData()
//	_, createResult := s.makeRequest("POST", "/api/primary-business-owners", data)
//	pboID := createResult["data"].(map[string]interface{})["id"]
//
//	// Now delete it
//	resp, result := s.makeRequest("DELETE", fmt.Sprintf("/api/primary-business-owners/%s", pboID), nil)
//
//	s.Equal(http.StatusOK, resp.StatusCode)
//	s.True(result["success"].(bool))
//
//	// Verify it's deleted
//	var pbo models.PrimaryBusinessOwner
//	s.NotNil(facades.Orm().Query().WithTrashed().Where("id = ?", pboID).First(&pbo))
//	s.NotNil(pbo.DeletedAt)
//}
//
//// ============================================================================
//// LIST/INDEX Tests
//// ============================================================================
//
//func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TestListPrimaryBusinessOwners() {
//	// Create multiple PBOs
//	for i := 0; i < 3; i++ {
//		data := s.getValidPBOData()
//		data["national_id_number"] = fmt.Sprintf("MWK1234567%d", i)
//		data["email"] = fmt.Sprintf("pbo%d@example.com", i)
//		s.makeRequest("POST", "/api/primary-business-owners", data)
//	}
//
//	// Get the list
//	resp, result := s.makeRequest("GET", "/api/primary-business-owners", nil)
//
//	s.Equal(http.StatusOK, resp.StatusCode)
//	s.True(result["success"].(bool))
//	s.Len(result["data"].([]interface{}), 3)
//}
//
//// ============================================================================
//// SEARCH Tests
//// ============================================================================
//
//func (s *PrimaryBusinessOwnerControllerCRUDTestSuite) TestSearchPrimaryBusinessOwners() {
//	// Create test data
//	testData := []map[string]interface{}{
//		{"first_name": "John", "last_name": "Doe", "national_id_number": "MWK11111111"},
//		{"first_name": "Jane", "last_name": "Smith", "national_id_number": "MWK22222222"},
//		{"first_name": "John", "last_name": "Smith", "national_id_number": "MWK33333333"},
//	}
//
//	for _, data := range testData {
//		pboData := s.getValidPBOData()
//		for k, v := range data {
//			pboData[k] = v
//		}
//		s.makeRequest("POST", "/api/primary-business-owners", pboData)
//	}
//
//	// Search by first name
//	resp, result := s.makeRequest("GET", "/api/primary-business-owners?search=John", nil)
//	s.Equal(http.StatusOK, resp.StatusCode)
//	s.True(result["success"].(bool))
//	s.Len(result["data"].([]interface{}), 2)
//
//	// Search by last name
//	resp, result = s.makeRequest("GET", "/api/primary-business-owners?search=Smith", nil)
//	s.Equal(http.StatusOK, resp.StatusCode)
//	s.True(result["success"].(bool))
//	s.Len(result["data"].([]interface{}), 2)
//
//	// Search by national ID
//	resp, result = s.makeRequest("GET", "/api/primary-business-owners?search=MWK22222222", nil)
//	s.Equal(http.StatusOK, resp.StatusCode)
//	s.True(result["success"].(bool))
//	s.Len(result["data"].([]interface{}), 1)
//	s.Equal("Jane", result["data"].([]interface{})[0].(map[string]interface{})["first_name"])
//}
