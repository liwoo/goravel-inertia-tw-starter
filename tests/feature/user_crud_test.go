package feature

import (
	"fmt"
	"testing"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"players/app/contracts"
	"players/app/models"
	"players/app/services"
	"players/tests"
)

type UserCRUDTestSuite struct {
	suite.Suite
	tests.TestCase

	userService *services.UserService
	adminUser   *models.User
	testRole    *models.Role
}

func TestUserCRUDTestSuite(t *testing.T) {
	suite.Run(t, new(UserCRUDTestSuite))
}

func (s *UserCRUDTestSuite) SetupTest() {
	s.RefreshDatabase()
	s.userService = services.NewUserService()

	// Create admin user for tests
	password, _ := facades.Hash().Make("password123")
	s.adminUser = &models.User{
		Name:         "Admin User",
		Email:        "admin@test.com",
		Password:     password,
		IsActive:     true,
		IsSuperAdmin: true,
	}
	facades.Orm().Query().Create(s.adminUser)

	// Create test role
	s.testRole = &models.Role{
		Name:        "Test Role",
		Slug:        "test-role",
		Description: "Role for testing",
		IsActive:    true,
	}
	facades.Orm().Query().Create(s.testRole)
}

func (s *UserCRUDTestSuite) TestCreateUser_ValidData_Success() {
	// Test data that matches the UI error case
	userData := map[string]interface{}{
		"name":           "Test User",        // 9 characters
		"email":          "test@example.com", // 16 characters
		"password":       "12345678",         // 8 characters (minimum)
		"is_active":      true,
		"is_super_admin": false,
		"role_id":        float64(s.testRole.ID),
	}

	// Create user using service
	result, err := s.userService.Create(userData)
	s.NoError(err, "User creation should succeed with valid data")
	s.NotNil(result)

	// Verify user was created
	createdUser := result.(*models.User)
	s.Equal("Test User", createdUser.Name)
	s.Equal("test@example.com", createdUser.Email)
	s.True(createdUser.IsActive)
	s.False(createdUser.IsSuperAdmin)

	// Verify in database
	var dbUser models.User
	err = facades.Orm().Query().Where("email = ?", "test@example.com").First(&dbUser)
	s.NoError(err)
	s.Equal("Test User", dbUser.Name)
}

func (s *UserCRUDTestSuite) TestCreateUser_LongFieldValues_ValidationError() {
	// Create strings that exceed 255 characters
	longName := make([]byte, 300)
	for i := range longName {
		longName[i] = 'a'
	}

	longEmail := make([]byte, 250)
	for i := range longEmail {
		longEmail[i] = 'a'
	}
	longEmailStr := string(longEmail) + "@example.com" // Over 255 chars

	userData := map[string]interface{}{
		"name":           string(longName),
		"email":          longEmailStr,
		"password":       "12345678",
		"is_active":      true,
		"is_super_admin": false,
	}

	// Should fail validation
	_, err := s.userService.Create(userData)
	s.Error(err, "Should fail validation for fields exceeding max length")
	s.Contains(err.Error(), "max", "Error should mention max length validation")
}

func (s *UserCRUDTestSuite) TestCreateUser_ShortPassword_ValidationError() {
	userData := map[string]interface{}{
		"name":           "Test User",
		"email":          "test@example.com",
		"password":       "123",    // Too short
		"is_active":      true,
		"is_super_admin": false,
	}

	_, err := s.userService.Create(userData)
	s.Error(err, "Should fail validation for short password")
	s.Contains(err.Error(), "min", "Error should mention minimum length validation")
}

func (s *UserCRUDTestSuite) TestCreateUser_InvalidEmail_ValidationError() {
	userData := map[string]interface{}{
		"name":           "Test User",
		"email":          "invalid-email",
		"password":       "12345678",
		"is_active":      true,
		"is_super_admin": false,
	}

	_, err := s.userService.Create(userData)
	s.Error(err, "Should fail validation for invalid email format")
}

func (s *UserCRUDTestSuite) TestCreateUser_MissingRequiredFields_ValidationError() {
	userData := map[string]interface{}{
		"is_active": true,
	}

	_, err := s.userService.Create(userData)
	s.Error(err, "Should fail validation for missing required fields")
	s.Contains(err.Error(), "required", "Error should mention required fields")
}

func (s *UserCRUDTestSuite) TestCreateUser_DuplicateEmail_ValidationError() {
	// First create a user
	userData1 := map[string]interface{}{
		"name":           "User One",
		"email":          "duplicate@example.com",
		"password":       "12345678",
		"is_active":      true,
		"is_super_admin": false,
	}

	_, err := s.userService.Create(userData1)
	s.NoError(err, "First user creation should succeed")

	// Try to create another with same email
	userData2 := map[string]interface{}{
		"name":           "User Two",
		"email":          "duplicate@example.com",
		"password":       "12345678",
		"is_active":      true,
		"is_super_admin": false,
	}

	_, err = s.userService.Create(userData2)
	s.Error(err, "Should fail validation for duplicate email")
	s.Contains(err.Error(), "email", "Error should mention email conflict")
}

func (s *UserCRUDTestSuite) TestGetUser_ValidID_Success() {
	// Create test user first
	userData := map[string]interface{}{
		"name":           "Get Test User",
		"email":          "get@example.com",
		"password":       "12345678",
		"is_active":      true,
		"is_super_admin": false,
	}

	result, err := s.userService.Create(userData)
	s.NoError(err)
	createdUser := result.(*models.User)

	// Get the user by ID
	retrievedResult, err := s.userService.GetByID(createdUser.ID)
	s.NoError(err)
	s.NotNil(retrievedResult)

	retrievedUser := retrievedResult.(*models.User)
	s.Equal("Get Test User", retrievedUser.Name)
	s.Equal("get@example.com", retrievedUser.Email)
}

func (s *UserCRUDTestSuite) TestGetUser_InvalidID_NotFound() {
	_, err := s.userService.GetByID(99999)
	s.Error(err, "Should return error for non-existent user")
}

func (s *UserCRUDTestSuite) TestUpdateUser_ValidData_Success() {
	// Create test user first
	userData := map[string]interface{}{
		"name":           "Original Name",
		"email":          "update@example.com",
		"password":       "12345678",
		"is_active":      true,
		"is_super_admin": false,
	}

	result, err := s.userService.Create(userData)
	s.NoError(err)
	createdUser := result.(*models.User)

	// Update the user
	updateData := map[string]interface{}{
		"name":      "Updated Name",
		"is_active": false,
	}

	updatedResult, err := s.userService.Update(createdUser.ID, updateData)
	s.NoError(err)
	s.NotNil(updatedResult)

	updatedUser := updatedResult.(*models.User)
	s.Equal("Updated Name", updatedUser.Name)
	s.False(updatedUser.IsActive)
	s.Equal("update@example.com", updatedUser.Email) // Should remain unchanged
}

func (s *UserCRUDTestSuite) TestDeleteUser_ValidID_Success() {
	// Create test user first
	userData := map[string]interface{}{
		"name":           "Delete Me",
		"email":          "delete@example.com",
		"password":       "12345678",
		"is_active":      true,
		"is_super_admin": false,
	}

	result, err := s.userService.Create(userData)
	s.NoError(err)
	createdUser := result.(*models.User)

	// Delete the user
	err = s.userService.Delete(createdUser.ID)
	s.NoError(err)

	// Verify user is soft deleted (not found in normal query)
	_, err = s.userService.GetByID(createdUser.ID)
	s.Error(err, "Deleted user should not be found")
}

func (s *UserCRUDTestSuite) TestListUsers_WithPagination_Success() {
	// Create multiple test users
	for i := 0; i < 5; i++ {
		userData := map[string]interface{}{
			"name":           fmt.Sprintf("User %d", i+1),
			"email":          fmt.Sprintf("user%d@example.com", i+1),
			"password":       "12345678",
			"is_active":      true,
			"is_super_admin": false,
		}
		_, err := s.userService.Create(userData)
		s.NoError(err)
	}

	// Test pagination
	request := contracts.ListRequest{
		Page:     1,
		PageSize: 3,
		Sort:     "created_at",
		Direction: "DESC",
	}

	result, err := s.userService.GetList(request)
	s.NoError(err)
	s.NotNil(result)

	// Should have pagination info
	s.Equal(1, result.CurrentPage)
	s.Equal(3, result.PerPage)
	s.True(result.Total > 5) // Including admin user

	s.Equal(3, len(result.Data)) // Should return 3 items per page
	s.True(result.HasNext)       // Should have more pages
}

func (s *UserCRUDTestSuite) TestListUsers_WithFilters_Success() {
	// Create test users with different statuses
	activeUserData := map[string]interface{}{
		"name":           "Active User",
		"email":          "active@example.com",
		"password":       "12345678",
		"is_active":      true,
		"is_super_admin": false,
	}
	_, err := s.userService.Create(activeUserData)
	s.NoError(err)

	inactiveUserData := map[string]interface{}{
		"name":           "Inactive User",
		"email":          "inactive@example.com",
		"password":       "12345678",
		"is_active":      false,
		"is_super_admin": false,
	}
	_, err = s.userService.Create(inactiveUserData)
	s.NoError(err)

	// Filter for inactive users
	request := contracts.ListRequest{
		Page:     1,
		PageSize: 10,
		Filters: map[string]interface{}{
			"is_active": false,
		},
	}

	result, err := s.userService.GetList(request)
	s.NoError(err)
	s.NotNil(result)

	// Should only return inactive users
	for _, item := range result.Data {
		user := item.(*models.User)
		s.False(user.IsActive, "Filtered results should only include inactive users")
	}
}

// Debug test that reproduces the exact UI validation issue
func (s *UserCRUDTestSuite) TestCreateUser_ReproduceUIValidationError_Debug() {
	s.T().Log("=== Reproducing UI validation error ===")

	// Exact data from the UI error report
	userData := map[string]interface{}{
		"name":           "Test User",        // 9 characters - should pass max:255
		"email":          "test@example.com", // 16 characters - should pass max:255
		"password":       "12345678",         // 8 characters - should pass min:8
		"is_active":      true,
		"is_super_admin": false,
		"role_id":        float64(s.testRole.ID),
	}

	s.T().Logf("Test data: %+v", userData)
	s.T().Logf("Name length: %d", len(userData["name"].(string)))
	s.T().Logf("Email length: %d", len(userData["email"].(string)))
	s.T().Logf("Password length: %d", len(userData["password"].(string)))

	// This should succeed with valid data
	result, err := s.userService.Create(userData)
	
	if err != nil {
		s.T().Logf("ERROR: User creation failed: %v", err)
		s.T().Logf("This indicates a validation bug - valid data is being rejected")
		
		// Log the full error details for debugging
		s.T().Logf("Full error string: %s", err.Error())
		
		// Check if it's a validation error with specific messages
		if errStr := err.Error(); errStr != "" {
			s.T().Logf("Error analysis:")
			if contains := func(s, substr string) bool { 
				return len(s) >= len(substr) && (len(substr) == 0 || s[len(s)-len(substr):] == substr || 
				       (len(s) > len(substr) && s[:len(substr)] == substr) || 
				       (len(s) > len(substr) && s[len(s)-len(substr):] == substr))
			}; contains(errStr, "max") {
				s.T().Log("- Contains 'max' validation error (unexpected for short strings)")
			}
			if contains := func(s, substr string) bool { 
				return len(s) >= len(substr) && (len(substr) == 0 || s[len(s)-len(substr):] == substr || 
				       (len(s) > len(substr) && s[:len(substr)] == substr) || 
				       (len(s) > len(substr) && s[len(s)-len(substr):] == substr))
			}; contains(errStr, "min") {
				s.T().Log("- Contains 'min' validation error")
			}
		}
		
		// This assertion will fail and show the validation issue
		s.NoError(err, "User creation should succeed with valid data - this indicates a validation bug")
	} else {
		s.T().Log("SUCCESS: User creation succeeded as expected")
		s.NotNil(result)
		
		createdUser := result.(*models.User)
		s.Equal("Test User", createdUser.Name)
		s.Equal("test@example.com", createdUser.Email)
		
		s.T().Log("User created successfully - no validation bug detected")
	}
}