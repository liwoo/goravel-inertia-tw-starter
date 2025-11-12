package feature

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/goravel/framework/facades"
	redisfacades "github.com/goravel/redis/facades"
	"github.com/stretchr/testify/assert"

	"smedi-sme-db/app/models"
	"smedi-sme-db/tests"
)

type PasswordAttemptsTestSuite struct {
	tests.TestCase
	server *httptest.Server
	client *http.Client
}

// clearRedisKeys clears all password attempts related Redis keys
func (suite *PasswordAttemptsTestSuite) clearRedisKeys() {
	cacheDriver, err := redisfacades.Cache("default")
	if err != nil {
		return
	}

	// Get all password attempt keys (this is a simplified approach)
	// In production, you'd want to use SCAN or KEYS carefully
	patterns := []string{"password_attempts:*", "password_lockout:*", "password_locks_today:*"}

	for _, pattern := range patterns {
		// For testing purposes, we'll clear known test emails
		// In a real scenario, you might want to clear all matching keys
		testEmails := []string{"test@example.com", "locked@example.com", "remaining@example.com", "success@example.com", "fresh@example.com", "expire@example.com"}
		for _, email := range testEmails {
			key := ""
			switch pattern {
			case "password_attempts:*":
				key = "password_attempts:" + email
			case "password_lockout:*":
				key = "password_lockout:" + email
			case "password_locks_today:*":
				key = "password_locks_today:" + email
			}
			cacheDriver.Forget(key)
		}
	}
}

func TestPasswordAttempts(t *testing.T) {
	suite := &PasswordAttemptsTestSuite{}

	suite.RefreshDatabase()

	// Clear Redis keys from previous tests
	suite.clearRedisKeys()

	// Create test server
	suite.server = httptest.NewServer(facades.Route())
	defer suite.server.Close()

	// Create client
	suite.client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Create test user
	hashedPassword, err := facades.Hash().Make("testpass")
	assert.NoError(t, err)

	user := &models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: hashedPassword,
		IsActive: true,
	}
	assert.NoError(t, facades.Orm().Query().Create(user))

	t.Run("TestCheckStatus_Endpoint_UnlockedAccount", func(t *testing.T) {
		req, _ := http.NewRequest("GET", suite.server.URL+"/api/auth/password-attempts/status?email=test@example.com", nil)

		resp, err := suite.client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.True(t, response["success"].(bool))
		assert.False(t, response["locked"].(bool))
		assert.Equal(t, float64(3), response["remaining_attempts"].(float64))
		_, hasUnlockTime := response["unlock_time"]
		assert.False(t, hasUnlockTime)
	})

	t.Run("TestGetRemainingAttempts_Endpoint", func(t *testing.T) {
		req, _ := http.NewRequest("GET", suite.server.URL+"/api/auth/password-attempts/remaining?email=test@example.com", nil)

		resp, err := suite.client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.True(t, response["success"].(bool))
		assert.Equal(t, float64(3), response["remaining_attempts"].(float64))
		assert.False(t, response["locked"].(bool))
	})

	t.Run("TestFailedLogin_Attempts_Tracking", func(t *testing.T) {
		// First failed login
		loginBody := `{"email":"test@example.com","password":"wrongpass"}`
		req, _ := http.NewRequest("POST", suite.server.URL+"/api/auth/login", strings.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.False(t, response["success"].(bool))
		assert.Equal(t, "Invalid credentials", response["message"].(string))
		assert.Equal(t, float64(2), response["remaining_attempts"].(float64))

		// Check status after first attempt
		req2, _ := http.NewRequest("GET", suite.server.URL+"/api/auth/password-attempts/status?email=test@example.com", nil)
		resp2, err := suite.client.Do(req2)
		assert.NoError(t, err)
		defer resp2.Body.Close()

		var statusResponse map[string]interface{}
		json.NewDecoder(resp2.Body).Decode(&statusResponse)

		assert.Equal(t, float64(2), statusResponse["remaining_attempts"].(float64))
		assert.False(t, statusResponse["locked"].(bool))
	})

	t.Run("TestFailedLogin_Warning_OnSecondAttempt", func(t *testing.T) {
		// Second failed login - should show warning
		loginBody := `{"email":"test@example.com","password":"wrongpass"}`
		req, _ := http.NewRequest("POST", suite.server.URL+"/api/auth/login", strings.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.False(t, response["success"].(bool))
		assert.Equal(t, "Invalid credentials", response["message"].(string))
		assert.Equal(t, float64(1), response["remaining_attempts"].(float64))
		assert.Equal(t, "One more failed attempt will lock your account.", response["warning"].(string))
	})

	t.Run("TestFailedLogin_AccountLockout_OnThirdAttempt", func(t *testing.T) {
		// Third failed login - should lock account
		loginBody := `{"email":"test@example.com","password":"wrongpass"}`
		req, _ := http.NewRequest("POST", suite.server.URL+"/api/auth/login", strings.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.False(t, response["success"].(bool))
		assert.True(t, response["locked"].(bool))
		assert.Contains(t, response["message"].(string), "Account has been locked")

		// Verify unlock time is provided
		_, hasUnlockTime := response["unlock_time"]
		assert.True(t, hasUnlockTime)
		_, hasUnlockTimestamp := response["unlock_timestamp"]
		assert.True(t, hasUnlockTimestamp)
	})

	t.Run("TestLockedAccount_PreventsLogin", func(t *testing.T) {
		// Try to login with correct password while account is locked
		loginBody := `{"email":"test@example.com","password":"testpass"}`
		req, _ := http.NewRequest("POST", suite.server.URL+"/api/auth/login", strings.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.False(t, response["success"].(bool))
		assert.True(t, response["locked"].(bool))
		assert.Contains(t, response["message"].(string), "locked")
	})

	t.Run("TestSuccessfulLogin_ClearsAttempts", func(t *testing.T) {
		// Create a fresh user for successful login test
		freshUser := &models.User{
			Name:     "Fresh User",
			Email:    "fresh@example.com",
			Password: hashedPassword,
			IsActive: true,
		}
		assert.NoError(t, facades.Orm().Query().Create(freshUser))

		// Make a failed attempt first
		failedLoginBody := `{"email":"fresh@example.com","password":"wrongpass"}`
		req, _ := http.NewRequest("POST", suite.server.URL+"/api/auth/login", strings.NewReader(failedLoginBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.client.Do(req)
		assert.NoError(t, err)
		resp.Body.Close()

		// Verify attempt was recorded
		statusReq, _ := http.NewRequest("GET", suite.server.URL+"/api/auth/password-attempts/status?email=fresh@example.com", nil)
		statusResp, err := suite.client.Do(statusReq)
		assert.NoError(t, err)

		var statusResponse map[string]interface{}
		json.NewDecoder(statusResp.Body).Decode(&statusResponse)
		statusResp.Body.Close()

		assert.Equal(t, float64(2), statusResponse["remaining_attempts"].(float64))

		// Now login successfully
		successLoginBody := `{"email":"fresh@example.com","password":"testpass"}`
		req2, _ := http.NewRequest("POST", suite.server.URL+"/api/auth/login", strings.NewReader(successLoginBody))
		req2.Header.Set("Content-Type", "application/json")

		resp2, err := suite.client.Do(req2)
		assert.NoError(t, err)
		defer resp2.Body.Close()

		assert.Equal(t, http.StatusOK, resp2.StatusCode)

		var successResponse map[string]interface{}
		json.NewDecoder(resp2.Body).Decode(&successResponse)

		assert.True(t, successResponse["success"].(bool))
		assert.Equal(t, "Login successful", successResponse["message"].(string))

		// Verify attempts were cleared
		statusReq2, _ := http.NewRequest("GET", suite.server.URL+"/api/auth/password-attempts/status?email=fresh@example.com", nil)
		statusResp2, err := suite.client.Do(statusReq2)
		assert.NoError(t, err)

		var statusResponse2 map[string]interface{}
		json.NewDecoder(statusResp2.Body).Decode(&statusResponse2)
		statusResp2.Body.Close()

		assert.Equal(t, float64(3), statusResponse2["remaining_attempts"].(float64))
		assert.False(t, statusResponse2["locked"].(bool))
	})

	t.Run("TestCheckStatus_MissingEmailParameter", func(t *testing.T) {
		req, _ := http.NewRequest("GET", suite.server.URL+"/api/auth/password-attempts/status", nil)

		resp, err := suite.client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.False(t, response["success"].(bool))
		assert.Equal(t, "Email parameter is required", response["message"].(string))
	})
}
