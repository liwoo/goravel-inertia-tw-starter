package users

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"players/app/http/controllers/auth/users"
	"players/app/models"
	"players/tests/helpers"
)

type UserCRUDTestSuite struct {
	suite.Suite
	testHelper *helpers.TestHelper
	adminUser  *models.User
}

func (suite *UserCRUDTestSuite) SetupSuite() {
	suite.testHelper = helpers.NewTestHelper()
	suite.testHelper.SetupTest()

	// Create admin user for authentication
	var err error
	suite.adminUser, err = suite.testHelper.CreateAdminUser("admin@example.com", "password123", true)
	if err != nil {
		suite.T().Fatalf("Failed to create admin user: %v", err)
	}
}

func (suite *UserCRUDTestSuite) TearDownSuite() {
	suite.testHelper.TearDownTest()
}

func (suite *UserCRUDTestSuite) SetupTest() {
	// Clean users table but keep the admin user
	facades.Orm().Query().Where("id != ?", suite.adminUser.ID).Delete(&models.User{})
}

func (suite *UserCRUDTestSuite) TestCreateUser_ValidData_Success() {
	// Create a librarian role
	role := &models.Role{
		Name:        "Librarian",
		Slug:        "librarian",
		Description: "Library management role",
		IsActive:    true,
	}
	err := facades.Orm().Query().Create(role)
	suite.Require().NoError(err)

	// Test data matching the error case from UI
	userData := map[string]interface{}{
		"name":           "Test User",
		"email":          "test@example.com",
		"password":       "12345678",
		"is_active":      true,
		"is_super_admin": false,
		"role_id":        role.ID,
	}

	// Make request
	token := suite.testHelper.LoginUser(suite.adminUser)
	response := suite.testHelper.PostJSON("/api/users", userData, token)

	// Assertions
	suite.Equal(http.StatusCreated, response.StatusCode)

	var responseBody map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	suite.Require().NoError(err)

	suite.True(responseBody["success"].(bool))
	suite.NotNil(responseBody["data"])

	// Verify user was created in database
	var createdUser models.User
	err = facades.Orm().Query().Where("email = ?", "test@example.com").First(&createdUser)
	suite.NoError(err)
	suite.Equal("Test User", createdUser.Name)
	suite.Equal("test@example.com", createdUser.Email)
	suite.True(createdUser.IsActive)
	suite.False(createdUser.IsSuperAdmin)
}

func (suite *UserCRUDTestSuite) TestCreateUser_InvalidEmail_ValidationError() {
	userData := map[string]interface{}{
		"name":           "Test User",
		"email":          "invalid-email",
		"password":       "12345678",
		"is_active":      true,
		"is_super_admin": false,
	}

	token := suite.testHelper.LoginUser(suite.adminUser)
	response := suite.testHelper.PostJSON("/api/users", userData, token)

	suite.Equal(http.StatusBadRequest, response.StatusCode)

	var responseBody map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&responseBody)
	suite.Require().NoError(err)

	suite.False(responseBody["success"].(bool))
	suite.Equal("Validation failed", responseBody["message"])
	suite.NotNil(responseBody["errors"])
}

func (suite *UserCRUDTestSuite) TestCreateUser_ShortPassword_ValidationError() {
	userData := map[string]interface{}{
		"name":           "Test User",
		"email":          "test@example.com",
		"password":       "123", // Too short
		"is_active":      true,
		"is_super_admin": false,
	}

	token := suite.testHelper.LoginUser(suite.adminUser)
	response := suite.testHelper.PostJSON("/api/users", userData, token)

	suite.Equal(http.StatusBadRequest, response.StatusCode)

	var responseBody map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&responseBody)
	suite.Require().NoError(err)

	suite.False(responseBody["success"].(bool))
	suite.Contains(responseBody["message"], "Validation")

	// Check specific password error
	errors, exists := responseBody["errors"].(map[string]interface{})
	suite.True(exists)
	passwordErrors, exists := errors["password"].(map[string]interface{})
	suite.True(exists)
	suite.Contains(fmt.Sprintf("%v", passwordErrors), "min")
}

func (suite *UserCRUDTestSuite) TestCreateUser_LongFieldValues_ValidationError() {
	// Test with very long values that exceed 255 characters
	longString := string(make([]byte, 300)) // 300 characters
	for i := range longString {
		longString = longString[:i] + "a" + longString[i+1:]
	}

	userData := map[string]interface{}{
		"name":           longString,
		"email":          longString + "@example.com",
		"password":       "12345678",
		"is_active":      true,
		"is_super_admin": false,
	}

	token := suite.testHelper.LoginUser(suite.adminUser)
	response := suite.testHelper.PostJSON("/api/users", userData, token)

	suite.Equal(http.StatusBadRequest, response.StatusCode)

	var responseBody map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&responseBody)
	suite.Require().NoError(err)

	suite.False(responseBody["success"].(bool))
	suite.Contains(responseBody["message"], "Validation")

	// Check for max length errors
	errors, exists := responseBody["errors"].(map[string]interface{})
	suite.True(exists)
	
	// Should have both name and email max errors
	if nameErrors, exists := errors["name"]; exists {
		suite.Contains(fmt.Sprintf("%v", nameErrors), "max")
	}
	if emailErrors, exists := errors["email"]; exists {
		suite.Contains(fmt.Sprintf("%v", emailErrors), "max")
	}
}

func (suite *UserCRUDTestSuite) TestCreateUser_MissingRequiredFields_ValidationError() {
	userData := map[string]interface{}{
		"is_active": true,
	}

	token := suite.testHelper.LoginUser(suite.adminUser)
	response := suite.testHelper.PostJSON("/api/users", userData, token)

	suite.Equal(http.StatusBadRequest, response.StatusCode)

	var responseBody map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&responseBody)
	suite.Require().NoError(err)

	suite.False(responseBody["success"].(bool))
	suite.Contains(responseBody["message"], "Validation")

	errors, exists := responseBody["errors"].(map[string]interface{})
	suite.True(exists)
	
	// Should have required field errors
	suite.Contains(fmt.Sprintf("%v", errors), "name")
	suite.Contains(fmt.Sprintf("%v", errors), "email")
	suite.Contains(fmt.Sprintf("%v", errors), "password")
}

func (suite *UserCRUDTestSuite) TestCreateUser_DuplicateEmail_ValidationError() {
	// First create a user
	userData1 := map[string]interface{}{
		"name":           "User One",
		"email":          "duplicate@example.com",
		"password":       "12345678",
		"is_active":      true,
		"is_super_admin": false,
	}

	token := suite.testHelper.LoginUser(suite.adminUser)
	response1 := suite.testHelper.PostJSON("/api/users", userData1, token)
	suite.Equal(http.StatusCreated, response1.StatusCode)

	// Try to create another user with same email
	userData2 := map[string]interface{}{
		"name":           "User Two",
		"email":          "duplicate@example.com",
		"password":       "12345678",
		"is_active":      true,
		"is_super_admin": false,
	}

	response2 := suite.testHelper.PostJSON("/api/users", userData2, token)
	suite.Equal(http.StatusBadRequest, response2.StatusCode)

	var responseBody map[string]interface{}
	err := json.NewDecoder(response2.Body).Decode(&responseBody)
	suite.Require().NoError(err)

	suite.False(responseBody["success"].(bool))
	// Should indicate email already exists
	suite.Contains(fmt.Sprintf("%v", responseBody), "email")
}

func (suite *UserCRUDTestSuite) TestGetUser_ValidID_Success() {
	// Create a test user
	testUser := &models.User{
		Name:         "Test User",
		Email:        "get@example.com",
		Password:     "hashed_password",
		IsActive:     true,
		IsSuperAdmin: false,
	}
	err := facades.Orm().Query().Create(testUser)
	suite.Require().NoError(err)

	token := suite.testHelper.LoginUser(suite.adminUser)
	response := suite.testHelper.GetJSON(fmt.Sprintf("/api/users/%d", testUser.ID), token)

	suite.Equal(http.StatusOK, response.StatusCode)

	var responseBody map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	suite.Require().NoError(err)

	suite.True(responseBody["success"].(bool))
	data := responseBody["data"].(map[string]interface{})
	suite.Equal("Test User", data["name"])
	suite.Equal("get@example.com", data["email"])
}

func (suite *UserCRUDTestSuite) TestGetUser_InvalidID_NotFound() {
	token := suite.testHelper.LoginUser(suite.adminUser)
	response := suite.testHelper.GetJSON("/api/users/99999", token)

	suite.Equal(http.StatusNotFound, response.StatusCode)

	var responseBody map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&responseBody)
	suite.Require().NoError(err)

	suite.False(responseBody["success"].(bool))
}

func (suite *UserCRUDTestSuite) TestUpdateUser_ValidData_Success() {
	// Create a test user
	testUser := &models.User{
		Name:         "Original Name",
		Email:        "update@example.com",
		Password:     "hashed_password",
		IsActive:     true,
		IsSuperAdmin: false,
	}
	err := facades.Orm().Query().Create(testUser)
	suite.Require().NoError(err)

	updateData := map[string]interface{}{
		"name":      "Updated Name",
		"is_active": false,
	}

	token := suite.testHelper.LoginUser(suite.adminUser)
	response := suite.testHelper.PutJSON(fmt.Sprintf("/api/users/%d", testUser.ID), updateData, token)

	suite.Equal(http.StatusOK, response.StatusCode)

	var responseBody map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	suite.Require().NoError(err)

	suite.True(responseBody["success"].(bool))

	// Verify database was updated
	var updatedUser models.User
	err = facades.Orm().Query().Find(&updatedUser, testUser.ID)
	suite.NoError(err)
	suite.Equal("Updated Name", updatedUser.Name)
	suite.False(updatedUser.IsActive)
}

func (suite *UserCRUDTestSuite) TestDeleteUser_ValidID_Success() {
	// Create a test user
	testUser := &models.User{
		Name:         "Delete Me",
		Email:        "delete@example.com",
		Password:     "hashed_password",
		IsActive:     true,
		IsSuperAdmin: false,
	}
	err := facades.Orm().Query().Create(testUser)
	suite.Require().NoError(err)

	token := suite.testHelper.LoginUser(suite.adminUser)
	response := suite.testHelper.DeleteJSON(fmt.Sprintf("/api/users/%d", testUser.ID), token)

	suite.Equal(http.StatusOK, response.StatusCode)

	var responseBody map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	suite.Require().NoError(err)

	suite.True(responseBody["success"].(bool))

	// Verify user was soft deleted
	var deletedUser models.User
	err = facades.Orm().Query().Find(&deletedUser, testUser.ID)
	suite.Error(err) // Should not find the user (soft deleted)

	// But should find with Unscoped
	err = facades.Orm().Query().Unscoped().Find(&deletedUser, testUser.ID)
	suite.NoError(err)
	suite.NotNil(deletedUser.DeletedAt)
}

func (suite *UserCRUDTestSuite) TestListUsers_WithPagination_Success() {
	// Create multiple test users
	for i := 0; i < 15; i++ {
		testUser := &models.User{
			Name:         fmt.Sprintf("User %d", i+1),
			Email:        fmt.Sprintf("user%d@example.com", i+1),
			Password:     "hashed_password",
			IsActive:     true,
			IsSuperAdmin: false,
		}
		err := facades.Orm().Query().Create(testUser)
		suite.Require().NoError(err)
	}

	token := suite.testHelper.LoginUser(suite.adminUser)
	response := suite.testHelper.GetJSON("/api/users?page=1&pageSize=10", token)

	suite.Equal(http.StatusOK, response.StatusCode)

	var responseBody map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&responseBody)
	suite.Require().NoError(err)

	suite.True(responseBody["success"].(bool))
	data := responseBody["data"].(map[string]interface{})
	
	// Should have pagination info
	suite.Equal(float64(1), data["currentPage"])
	suite.Equal(float64(10), data["perPage"])
	suite.True(data["total"].(float64) > 15) // Including admin user

	users := data["data"].([]interface{})
	suite.Equal(10, len(users)) // Should return 10 items per page
}

func (suite *UserCRUDTestSuite) TestListUsers_WithFilters_Success() {
	// Create test users with different statuses
	activeUser := &models.User{
		Name:         "Active User",
		Email:        "active@example.com",
		Password:     "hashed_password",
		IsActive:     true,
		IsSuperAdmin: false,
	}
	err := facades.Orm().Query().Create(activeUser)
	suite.Require().NoError(err)

	inactiveUser := &models.User{
		Name:         "Inactive User",
		Email:        "inactive@example.com",
		Password:     "hashed_password",
		IsActive:     false,
		IsSuperAdmin: false,
	}
	err = facades.Orm().Query().Create(inactiveUser)
	suite.Require().NoError(err)

	token := suite.testHelper.LoginUser(suite.adminUser)
	response := suite.testHelper.GetJSON("/api/users?is_active=false", token)

	suite.Equal(http.StatusOK, response.StatusCode)

	var responseBody map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	suite.Require().NoError(err)

	suite.True(responseBody["success"].(bool))
	data := responseBody["data"].(map[string]interface{})
	users := data["data"].([]interface{})
	
	// Should only return inactive users
	for _, userInterface := range users {
		user := userInterface.(map[string]interface{})
		suite.False(user["is_active"].(bool))
	}
}

func (suite *UserCRUDTestSuite) TestUserCRUD_UnauthorizedUser_AccessDenied() {
	// Create a regular user (non-admin)
	regularUser := &models.User{
		Name:         "Regular User",
		Email:        "regular@example.com",
		Password:     "hashed_password",
		IsActive:     true,
		IsSuperAdmin: false,
	}
	err := facades.Orm().Query().Create(regularUser)
	suite.Require().NoError(err)

	token := suite.testHelper.LoginUser(regularUser)
	
	// Try to create a user
	userData := map[string]interface{}{
		"name":           "Test User",
		"email":          "test@example.com",
		"password":       "12345678",
		"is_active":      true,
		"is_super_admin": false,
	}

	response := suite.testHelper.PostJSON("/api/users", userData, token)
	suite.Equal(http.StatusForbidden, response.StatusCode)
}

func (suite *UserCRUDTestSuite) TestUserCRUD_NoAuthentication_AccessDenied() {
	// Try to access without authentication
	userData := map[string]interface{}{
		"name":           "Test User",
		"email":          "test@example.com",
		"password":       "12345678",
		"is_active":      true,
		"is_super_admin": false,
	}

	response := suite.testHelper.PostJSON("/api/users", userData, "")
	suite.Equal(http.StatusUnauthorized, response.StatusCode)
}

// Debug test to reproduce the exact UI error
func (suite *UserCRUDTestSuite) TestCreateUser_ReproduceUIError_Debug() {
	// Create a librarian role first
	role := &models.Role{
		Name:        "Librarian",
		Slug:        "librarian", 
		Description: "Library management role",
		IsActive:    true,
	}
	err := facades.Orm().Query().Create(role)
	suite.Require().NoError(err)

	// Exact data from UI error
	userData := map[string]interface{}{
		"name":           "Test User",        // 9 characters - well under 255
		"email":          "test@example.com", // 16 characters - well under 255  
		"password":       "12345678",         // 8 characters - meets minimum
		"is_active":      true,
		"is_super_admin": false,
		"role_id":        role.ID,
	}

	// Log the data being sent
	jsonData, _ := json.Marshal(userData)
	suite.T().Logf("Sending data: %s", string(jsonData))

	token := suite.testHelper.LoginUser(suite.adminUser)
	
	// Make the request and capture full response
	requestBody, _ := json.Marshal(userData)
	req := suite.testHelper.NewRequest("POST", "/api/users", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	
	response := suite.testHelper.MakeRequest(req)
	
	// Read and log the full response
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	responseBody := buf.String()
	suite.T().Logf("Response status: %d", response.StatusCode)
	suite.T().Logf("Response body: %s", responseBody)
	
	// Parse response
	var respData map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &respData)
	if err != nil {
		suite.T().Logf("Failed to parse response as JSON: %v", err)
	} else {
		if errors, exists := respData["errors"]; exists {
			suite.T().Logf("Validation errors: %+v", errors)
		}
	}

	// The test should pass with valid data
	suite.Equal(http.StatusCreated, response.StatusCode, "User creation should succeed with valid data")
}

func TestUserCRUDTestSuite(t *testing.T) {
	suite.Run(t, new(UserCRUDTestSuite))
}