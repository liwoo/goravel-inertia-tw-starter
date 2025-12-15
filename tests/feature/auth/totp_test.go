// Package feature provides feature tests for the TOTP 2FA functionality.
// These tests cover the TOTP service functions, API endpoint responses,
// password confirmation requirements, and backup code handling.
package feature

import (
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
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/suite"

	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
	"smedi-sme-db/tests/helpers"
)

// TOTPTestSuite provides comprehensive tests for TOTP 2FA functionality
type TOTPTestSuite struct {
	suite.Suite
	server      *httptest.Server
	client      *http.Client
	totpService *services.TOTPService
}

// SetupSuite initializes the test server
func (s *TOTPTestSuite) SetupSuite() {
	s.server = httptest.NewServer(facades.Route())

	// Create HTTP client with cookie jar for session handling
	jar, _ := cookiejar.New(nil)
	s.client = &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	s.totpService = services.NewTOTPService()
}

// TearDownSuite cleans up after all tests
func (s *TOTPTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}

// SetupTest cleans database before each test
func (s *TOTPTestSuite) SetupTest() {
	helpers.CleanTestDatabase()
}

// TearDownTest cleans database after each test
func (s *TOTPTestSuite) TearDownTest() {
	helpers.CleanupTestData()
}

// createTestUser creates a user for testing and returns the user and plain password
func (s *TOTPTestSuite) createTestUser(email string) (*models.User, string) {
	plainPassword := "TestPassword123!"
	hashedPassword, err := facades.Hash().Make(plainPassword)
	s.Require().NoError(err)

	user := &models.User{
		Name:     "Test User",
		Email:    email,
		Password: hashedPassword,
		IsActive: true,
	}
	err = facades.Orm().Query().Create(user)
	s.Require().NoError(err)

	return user, plainPassword
}

// loginUser performs login and stores auth cookie
func (s *TOTPTestSuite) loginUser(email, password string) *http.Cookie {
	loginBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)
	req, _ := http.NewRequest("POST", s.server.URL+"/api/auth/login", strings.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	// Find auth cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "token" {
			return cookie
		}
	}

	return nil
}

// makeAuthenticatedRequest makes an authenticated HTTP request
func (s *TOTPTestSuite) makeAuthenticatedRequest(method, path string, body string, cookie *http.Cookie) (*http.Response, []byte) {
	var reqBody io.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}

	req, _ := http.NewRequest(method, s.server.URL+path, reqBody)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if cookie != nil {
		req.AddCookie(cookie)
	}

	resp, err := s.client.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	return resp, respBody
}

// parseJSONResponse parses a JSON response body
func (s *TOTPTestSuite) parseJSONResponse(body []byte) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal(body, &result)
	s.Require().NoError(err)
	return result
}

// ============================================================================
// TOTP Service Unit Tests
// ============================================================================

func (s *TOTPTestSuite) TestTOTPService_GenerateSecret() {
	email := "test@example.com"

	key, err := s.totpService.GenerateSecret(email)

	s.NoError(err)
	s.NotNil(key)
	s.NotEmpty(key.Secret())
	s.Contains(key.URL(), email)
}

func (s *TOTPTestSuite) TestTOTPService_GenerateQRCode() {
	key, err := s.totpService.GenerateSecret("test@example.com")
	s.Require().NoError(err)

	qrCode, err := s.totpService.GenerateQRCode(key)

	s.NoError(err)
	s.NotEmpty(qrCode)
	s.Contains(qrCode, "data:image/png;base64,")
}

func (s *TOTPTestSuite) TestTOTPService_ValidateCode_ValidCode() {
	key, err := s.totpService.GenerateSecret("test@example.com")
	s.Require().NoError(err)

	// Generate a valid TOTP code
	validCode, err := totp.GenerateCode(key.Secret(), time.Now())
	s.Require().NoError(err)

	isValid := s.totpService.ValidateCode(key.Secret(), validCode)

	s.True(isValid)
}

func (s *TOTPTestSuite) TestTOTPService_ValidateCode_InvalidCode() {
	key, err := s.totpService.GenerateSecret("test@example.com")
	s.Require().NoError(err)

	isValid := s.totpService.ValidateCode(key.Secret(), "000000")

	s.False(isValid)
}

func (s *TOTPTestSuite) TestTOTPService_ValidateCode_EmptyCode() {
	key, err := s.totpService.GenerateSecret("test@example.com")
	s.Require().NoError(err)

	isValid := s.totpService.ValidateCode(key.Secret(), "")

	s.False(isValid)
}

func (s *TOTPTestSuite) TestTOTPService_GenerateBackupCodes() {
	plainCodes, hashedCodes, err := s.totpService.GenerateBackupCodes(10)

	s.NoError(err)
	s.Len(plainCodes, 10)
	s.Len(hashedCodes, 10)

	// Verify each code is unique
	codeSet := make(map[string]bool)
	for _, code := range plainCodes {
		s.NotEmpty(code)
		s.Len(code, 8)
		s.False(codeSet[code], "Duplicate backup code found")
		codeSet[code] = true
	}

	// Verify hashed codes are different from plain codes
	for i, plain := range plainCodes {
		s.NotEqual(plain, hashedCodes[i])
	}
}

func (s *TOTPTestSuite) TestTOTPService_GenerateBackupCodes_DefaultCount() {
	plainCodes, hashedCodes, err := s.totpService.GenerateBackupCodes(0)

	s.NoError(err)
	s.Len(plainCodes, 10)
	s.Len(hashedCodes, 10)
}

func (s *TOTPTestSuite) TestTOTPService_EncryptDecryptSecret() {
	originalSecret := "JBSWY3DPEHPK3PXP"

	encrypted, err := s.totpService.EncryptSecret(originalSecret)
	s.Require().NoError(err)
	s.NotEmpty(encrypted)
	s.NotEqual(originalSecret, encrypted)

	decrypted, err := s.totpService.DecryptSecret(encrypted)
	s.Require().NoError(err)
	s.Equal(originalSecret, decrypted)
}

func (s *TOTPTestSuite) TestTOTPService_DecryptSecret_InvalidData() {
	_, err := s.totpService.DecryptSecret("not-valid-base64!!!")
	s.Error(err)
}

func (s *TOTPTestSuite) TestTOTPService_SaveAndValidateBackupCode() {
	user, _ := s.createTestUser("backup-test@example.com")

	plainCodes, hashedCodes, err := s.totpService.GenerateBackupCodes(5)
	s.Require().NoError(err)

	err = s.totpService.SaveBackupCodes(user.ID, hashedCodes)
	s.Require().NoError(err)

	// Verify count
	count, err := s.totpService.GetBackupCodesRemaining(user.ID)
	s.Require().NoError(err)
	s.Equal(int64(5), count)

	// Validate a backup code
	isValid, err := s.totpService.ValidateBackupCode(user.ID, plainCodes[0])
	s.NoError(err)
	s.True(isValid)

	// Count should decrease
	count, err = s.totpService.GetBackupCodesRemaining(user.ID)
	s.Require().NoError(err)
	s.Equal(int64(4), count)

	// Same code should not work twice
	isValid, err = s.totpService.ValidateBackupCode(user.ID, plainCodes[0])
	s.NoError(err)
	s.False(isValid)
}

func (s *TOTPTestSuite) TestTOTPService_ValidateBackupCode_Invalid() {
	user, _ := s.createTestUser("backup-invalid@example.com")

	_, hashedCodes, err := s.totpService.GenerateBackupCodes(5)
	s.Require().NoError(err)

	err = s.totpService.SaveBackupCodes(user.ID, hashedCodes)
	s.Require().NoError(err)

	isValid, err := s.totpService.ValidateBackupCode(user.ID, "WRONGCODE")
	s.NoError(err)
	s.False(isValid)
}

func (s *TOTPTestSuite) TestTOTPService_EnableTOTP() {
	user, _ := s.createTestUser("enable-totp@example.com")

	key, err := s.totpService.GenerateSecret(user.Email)
	s.Require().NoError(err)

	err = s.totpService.EnableTOTP(user, key.Secret())
	s.NoError(err)

	// Reload user and verify
	var updatedUser models.User
	err = facades.Orm().Query().Find(&updatedUser, user.ID)
	s.Require().NoError(err)

	s.True(updatedUser.TOTPEnabled)
	s.NotEmpty(updatedUser.TOTPSecret)
	s.NotNil(updatedUser.TOTPVerifiedAt)
}

func (s *TOTPTestSuite) TestTOTPService_DisableTOTP() {
	user, _ := s.createTestUser("disable-totp@example.com")

	// First enable TOTP
	key, err := s.totpService.GenerateSecret(user.Email)
	s.Require().NoError(err)
	err = s.totpService.EnableTOTP(user, key.Secret())
	s.Require().NoError(err)

	// Save some backup codes
	_, hashedCodes, _ := s.totpService.GenerateBackupCodes(5)
	s.totpService.SaveBackupCodes(user.ID, hashedCodes)

	// Disable TOTP
	err = s.totpService.DisableTOTP(user.ID)
	s.NoError(err)

	// Verify user state
	var updatedUser models.User
	err = facades.Orm().Query().Find(&updatedUser, user.ID)
	s.Require().NoError(err)

	s.False(updatedUser.TOTPEnabled)
	s.Empty(updatedUser.TOTPSecret)

	// Verify backup codes are deleted
	count, _ := s.totpService.GetBackupCodesRemaining(user.ID)
	s.Equal(int64(0), count)
}

func (s *TOTPTestSuite) TestTOTPService_ValidateTOTPForUser_TOTPCode() {
	user, _ := s.createTestUser("validate-totp@example.com")

	key, err := s.totpService.GenerateSecret(user.Email)
	s.Require().NoError(err)
	err = s.totpService.EnableTOTP(user, key.Secret())
	s.Require().NoError(err)

	// Reload user with secret
	var updatedUser models.User
	err = facades.Orm().Query().Find(&updatedUser, user.ID)
	s.Require().NoError(err)

	// Generate valid code
	validCode, _ := totp.GenerateCode(key.Secret(), time.Now())

	valid, isBackup, err := s.totpService.ValidateTOTPForUser(&updatedUser, validCode)
	s.NoError(err)
	s.True(valid)
	s.False(isBackup)
}

func (s *TOTPTestSuite) TestTOTPService_ValidateTOTPForUser_BackupCode() {
	user, _ := s.createTestUser("validate-backup@example.com")

	key, err := s.totpService.GenerateSecret(user.Email)
	s.Require().NoError(err)
	err = s.totpService.EnableTOTP(user, key.Secret())
	s.Require().NoError(err)

	plainCodes, hashedCodes, _ := s.totpService.GenerateBackupCodes(5)
	s.totpService.SaveBackupCodes(user.ID, hashedCodes)

	// Reload user with secret
	var updatedUser models.User
	err = facades.Orm().Query().Find(&updatedUser, user.ID)
	s.Require().NoError(err)

	valid, isBackup, err := s.totpService.ValidateTOTPForUser(&updatedUser, plainCodes[0])
	s.NoError(err)
	s.True(valid)
	s.True(isBackup)
}

// ============================================================================
// API Endpoint Tests - 2FA Status
// ============================================================================

func (s *TOTPTestSuite) TestAPI_GetStatus_Unauthenticated() {
	resp, _ := s.makeAuthenticatedRequest("GET", "/api/2fa/status", "", nil)
	s.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *TOTPTestSuite) TestAPI_GetStatus_TOTPDisabled() {
	user, password := s.createTestUser("status-disabled@example.com")
	_ = user
	cookie := s.loginUser("status-disabled@example.com", password)
	s.Require().NotNil(cookie)

	resp, body := s.makeAuthenticatedRequest("GET", "/api/2fa/status", "", cookie)

	s.Equal(http.StatusOK, resp.StatusCode)

	result := s.parseJSONResponse(body)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.False(data["enabled"].(bool))
	s.Nil(data["verified_at"])
	s.Equal(float64(0), data["backup_codes_remaining"].(float64))
}

func (s *TOTPTestSuite) TestAPI_GetStatus_TOTPEnabled() {
	user, password := s.createTestUser("status-enabled@example.com")

	// Enable TOTP directly
	key, _ := s.totpService.GenerateSecret(user.Email)
	s.totpService.EnableTOTP(user, key.Secret())
	plainCodes, hashedCodes, _ := s.totpService.GenerateBackupCodes(10)
	_ = plainCodes
	s.totpService.SaveBackupCodes(user.ID, hashedCodes)

	cookie := s.loginUser2FARequired("status-enabled@example.com", password, key.Secret())
	s.Require().NotNil(cookie)

	resp, body := s.makeAuthenticatedRequest("GET", "/api/2fa/status", "", cookie)

	s.Equal(http.StatusOK, resp.StatusCode)

	result := s.parseJSONResponse(body)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.True(data["enabled"].(bool))
	s.NotNil(data["verified_at"])
	s.Equal(float64(10), data["backup_codes_remaining"].(float64))
}

// ============================================================================
// API Endpoint Tests - 2FA Setup
// ============================================================================

func (s *TOTPTestSuite) TestAPI_Setup_Unauthenticated() {
	body := `{"password":"test"}`
	resp, _ := s.makeAuthenticatedRequest("POST", "/api/2fa/setup", body, nil)
	s.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *TOTPTestSuite) TestAPI_Setup_MissingPassword() {
	_, password := s.createTestUser("setup-missing@example.com")
	cookie := s.loginUser("setup-missing@example.com", password)
	s.Require().NotNil(cookie)

	resp, body := s.makeAuthenticatedRequest("POST", "/api/2fa/setup", `{}`, cookie)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)

	result := s.parseJSONResponse(body)
	s.False(result["success"].(bool))
}

func (s *TOTPTestSuite) TestAPI_Setup_IncorrectPassword() {
	_, password := s.createTestUser("setup-wrong@example.com")
	_ = password
	cookie := s.loginUser("setup-wrong@example.com", password)
	s.Require().NotNil(cookie)

	resp, body := s.makeAuthenticatedRequest("POST", "/api/2fa/setup", `{"password":"WrongPassword123!"}`, cookie)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)

	result := s.parseJSONResponse(body)
	s.False(result["success"].(bool))

	errors := result["errors"].(map[string]interface{})
	s.Contains(errors, "password")
}

func (s *TOTPTestSuite) TestAPI_Setup_Success() {
	_, password := s.createTestUser("setup-success@example.com")
	cookie := s.loginUser("setup-success@example.com", password)
	s.Require().NotNil(cookie)

	reqBody := fmt.Sprintf(`{"password":"%s"}`, password)
	resp, body := s.makeAuthenticatedRequest("POST", "/api/2fa/setup", reqBody, cookie)

	s.Equal(http.StatusOK, resp.StatusCode)

	result := s.parseJSONResponse(body)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.NotEmpty(data["qr_code"])
	s.NotEmpty(data["secret"])
	s.Contains(data["qr_code"].(string), "data:image/png;base64,")
}

func (s *TOTPTestSuite) TestAPI_Setup_AlreadyEnabled() {
	user, password := s.createTestUser("setup-already@example.com")

	// Enable TOTP first
	key, _ := s.totpService.GenerateSecret(user.Email)
	s.totpService.EnableTOTP(user, key.Secret())

	cookie := s.loginUser2FARequired("setup-already@example.com", password, key.Secret())
	s.Require().NotNil(cookie)

	reqBody := fmt.Sprintf(`{"password":"%s"}`, password)
	resp, body := s.makeAuthenticatedRequest("POST", "/api/2fa/setup", reqBody, cookie)

	s.Equal(http.StatusBadRequest, resp.StatusCode)

	result := s.parseJSONResponse(body)
	s.False(result["success"].(bool))
	s.Contains(result["message"].(string), "already enabled")
}

// ============================================================================
// API Endpoint Tests - 2FA Verify
// ============================================================================

func (s *TOTPTestSuite) TestAPI_Verify_MissingCode() {
	_, password := s.createTestUser("verify-missing@example.com")
	cookie := s.loginUser("verify-missing@example.com", password)
	s.Require().NotNil(cookie)

	// Setup first
	reqBody := fmt.Sprintf(`{"password":"%s"}`, password)
	s.makeAuthenticatedRequest("POST", "/api/2fa/setup", reqBody, cookie)

	resp, body := s.makeAuthenticatedRequest("POST", "/api/2fa/verify", `{}`, cookie)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)

	result := s.parseJSONResponse(body)
	s.False(result["success"].(bool))
}

func (s *TOTPTestSuite) TestAPI_Verify_InvalidCode() {
	_, password := s.createTestUser("verify-invalid@example.com")
	cookie := s.loginUser("verify-invalid@example.com", password)
	s.Require().NotNil(cookie)

	// Setup first
	reqBody := fmt.Sprintf(`{"password":"%s"}`, password)
	s.makeAuthenticatedRequest("POST", "/api/2fa/setup", reqBody, cookie)

	resp, body := s.makeAuthenticatedRequest("POST", "/api/2fa/verify", `{"code":"000000"}`, cookie)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)

	result := s.parseJSONResponse(body)
	s.False(result["success"].(bool))

	errors := result["errors"].(map[string]interface{})
	s.Contains(errors, "code")
}

func (s *TOTPTestSuite) TestAPI_Verify_ExpiredSetup() {
	_, password := s.createTestUser("verify-expired@example.com")
	cookie := s.loginUser("verify-expired@example.com", password)
	s.Require().NotNil(cookie)

	// Skip setup and go directly to verify
	resp, body := s.makeAuthenticatedRequest("POST", "/api/2fa/verify", `{"code":"123456"}`, cookie)

	s.Equal(http.StatusBadRequest, resp.StatusCode)

	result := s.parseJSONResponse(body)
	s.False(result["success"].(bool))
	s.Contains(result["message"].(string), "expired")
}

func (s *TOTPTestSuite) TestAPI_Verify_Success() {
	_, password := s.createTestUser("verify-success@example.com")
	cookie := s.loginUser("verify-success@example.com", password)
	s.Require().NotNil(cookie)

	// Setup first
	setupBody := fmt.Sprintf(`{"password":"%s"}`, password)
	_, setupResp := s.makeAuthenticatedRequest("POST", "/api/2fa/setup", setupBody, cookie)

	setupResult := s.parseJSONResponse(setupResp)
	secret := setupResult["data"].(map[string]interface{})["secret"].(string)

	// Generate valid code
	validCode, _ := totp.GenerateCode(secret, time.Now())

	verifyBody := fmt.Sprintf(`{"code":"%s"}`, validCode)
	resp, body := s.makeAuthenticatedRequest("POST", "/api/2fa/verify", verifyBody, cookie)

	s.Equal(http.StatusOK, resp.StatusCode)

	result := s.parseJSONResponse(body)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	backupCodes := data["backup_codes"].([]interface{})
	s.Len(backupCodes, 10)
	s.NotNil(data["enabled_at"])
}

// ============================================================================
// API Endpoint Tests - 2FA Disable
// ============================================================================

func (s *TOTPTestSuite) TestAPI_Disable_NotEnabled() {
	_, password := s.createTestUser("disable-not-enabled@example.com")
	cookie := s.loginUser("disable-not-enabled@example.com", password)
	s.Require().NotNil(cookie)

	reqBody := fmt.Sprintf(`{"password":"%s","code":"123456"}`, password)
	resp, body := s.makeAuthenticatedRequest("POST", "/api/2fa/disable", reqBody, cookie)

	s.Equal(http.StatusBadRequest, resp.StatusCode)

	result := s.parseJSONResponse(body)
	s.False(result["success"].(bool))
	s.Contains(result["message"].(string), "not enabled")
}

func (s *TOTPTestSuite) TestAPI_Disable_IncorrectPassword() {
	user, password := s.createTestUser("disable-wrong-pass@example.com")

	// Enable TOTP
	key, _ := s.totpService.GenerateSecret(user.Email)
	s.totpService.EnableTOTP(user, key.Secret())

	cookie := s.loginUser2FARequired("disable-wrong-pass@example.com", password, key.Secret())
	s.Require().NotNil(cookie)

	validCode, _ := totp.GenerateCode(key.Secret(), time.Now())
	reqBody := fmt.Sprintf(`{"password":"WrongPassword!","code":"%s"}`, validCode)
	resp, body := s.makeAuthenticatedRequest("POST", "/api/2fa/disable", reqBody, cookie)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)

	result := s.parseJSONResponse(body)
	s.False(result["success"].(bool))
}

func (s *TOTPTestSuite) TestAPI_Disable_InvalidCode() {
	user, password := s.createTestUser("disable-wrong-code@example.com")

	// Enable TOTP
	key, _ := s.totpService.GenerateSecret(user.Email)
	s.totpService.EnableTOTP(user, key.Secret())

	cookie := s.loginUser2FARequired("disable-wrong-code@example.com", password, key.Secret())
	s.Require().NotNil(cookie)

	reqBody := fmt.Sprintf(`{"password":"%s","code":"000000"}`, password)
	resp, body := s.makeAuthenticatedRequest("POST", "/api/2fa/disable", reqBody, cookie)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)

	result := s.parseJSONResponse(body)
	s.False(result["success"].(bool))
}

func (s *TOTPTestSuite) TestAPI_Disable_WithTOTPCode_Success() {
	user, password := s.createTestUser("disable-totp@example.com")

	// Enable TOTP
	key, _ := s.totpService.GenerateSecret(user.Email)
	s.totpService.EnableTOTP(user, key.Secret())

	cookie := s.loginUser2FARequired("disable-totp@example.com", password, key.Secret())
	s.Require().NotNil(cookie)

	validCode, _ := totp.GenerateCode(key.Secret(), time.Now())
	reqBody := fmt.Sprintf(`{"password":"%s","code":"%s"}`, password, validCode)
	resp, body := s.makeAuthenticatedRequest("POST", "/api/2fa/disable", reqBody, cookie)

	s.Equal(http.StatusOK, resp.StatusCode)

	result := s.parseJSONResponse(body)
	s.True(result["success"].(bool))

	// Verify TOTP is disabled
	var updatedUser models.User
	facades.Orm().Query().Find(&updatedUser, user.ID)
	s.False(updatedUser.TOTPEnabled)
}

func (s *TOTPTestSuite) TestAPI_Disable_WithBackupCode_Success() {
	user, password := s.createTestUser("disable-backup@example.com")

	// Enable TOTP and save backup codes
	key, _ := s.totpService.GenerateSecret(user.Email)
	s.totpService.EnableTOTP(user, key.Secret())
	plainCodes, hashedCodes, _ := s.totpService.GenerateBackupCodes(5)
	s.totpService.SaveBackupCodes(user.ID, hashedCodes)

	cookie := s.loginUser2FARequired("disable-backup@example.com", password, key.Secret())
	s.Require().NotNil(cookie)

	reqBody := fmt.Sprintf(`{"password":"%s","code":"%s"}`, password, plainCodes[0])
	resp, body := s.makeAuthenticatedRequest("POST", "/api/2fa/disable", reqBody, cookie)

	s.Equal(http.StatusOK, resp.StatusCode)

	result := s.parseJSONResponse(body)
	s.True(result["success"].(bool))
}

// ============================================================================
// API Endpoint Tests - Regenerate Backup Codes
// ============================================================================

func (s *TOTPTestSuite) TestAPI_RegenerateBackupCodes_NotEnabled() {
	_, password := s.createTestUser("regen-not-enabled@example.com")
	cookie := s.loginUser("regen-not-enabled@example.com", password)
	s.Require().NotNil(cookie)

	reqBody := fmt.Sprintf(`{"password":"%s"}`, password)
	resp, body := s.makeAuthenticatedRequest("POST", "/api/2fa/backup-codes", reqBody, cookie)

	s.Equal(http.StatusBadRequest, resp.StatusCode)

	result := s.parseJSONResponse(body)
	s.False(result["success"].(bool))
}

func (s *TOTPTestSuite) TestAPI_RegenerateBackupCodes_IncorrectPassword() {
	user, password := s.createTestUser("regen-wrong-pass@example.com")

	// Enable TOTP
	key, _ := s.totpService.GenerateSecret(user.Email)
	s.totpService.EnableTOTP(user, key.Secret())

	cookie := s.loginUser2FARequired("regen-wrong-pass@example.com", password, key.Secret())
	s.Require().NotNil(cookie)

	resp, body := s.makeAuthenticatedRequest("POST", "/api/2fa/backup-codes", `{"password":"WrongPass!"}`, cookie)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)

	result := s.parseJSONResponse(body)
	s.False(result["success"].(bool))
}

func (s *TOTPTestSuite) TestAPI_RegenerateBackupCodes_Success() {
	user, password := s.createTestUser("regen-success@example.com")

	// Enable TOTP and save initial backup codes
	key, _ := s.totpService.GenerateSecret(user.Email)
	s.totpService.EnableTOTP(user, key.Secret())
	_, hashedCodes, _ := s.totpService.GenerateBackupCodes(5)
	s.totpService.SaveBackupCodes(user.ID, hashedCodes)

	cookie := s.loginUser2FARequired("regen-success@example.com", password, key.Secret())
	s.Require().NotNil(cookie)

	reqBody := fmt.Sprintf(`{"password":"%s"}`, password)
	resp, body := s.makeAuthenticatedRequest("POST", "/api/2fa/backup-codes", reqBody, cookie)

	s.Equal(http.StatusOK, resp.StatusCode)

	result := s.parseJSONResponse(body)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	backupCodes := data["backup_codes"].([]interface{})
	s.Len(backupCodes, 10)

	// Verify old codes are invalidated
	count, _ := s.totpService.GetBackupCodesRemaining(user.ID)
	s.Equal(int64(10), count)
}

// ============================================================================
// API Endpoint Tests - Login with 2FA
// ============================================================================

func (s *TOTPTestSuite) TestAPI_Login_NoTOTP() {
	_, password := s.createTestUser("login-no-totp@example.com")

	loginBody := fmt.Sprintf(`{"email":"login-no-totp@example.com","password":"%s"}`, password)
	req, _ := http.NewRequest("POST", s.server.URL+"/api/auth/login", strings.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	result := s.parseJSONResponse(body)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	// Should have token directly
	data := result["data"].(map[string]interface{})
	s.NotEmpty(data["token"])

	// Should NOT have requires_2fa
	_, has2FA := result["requires_2fa"]
	s.False(has2FA)
}

func (s *TOTPTestSuite) TestAPI_Login_RequiresTOTP() {
	user, password := s.createTestUser("login-totp@example.com")

	// Enable TOTP
	key, _ := s.totpService.GenerateSecret(user.Email)
	s.totpService.EnableTOTP(user, key.Secret())

	loginBody := fmt.Sprintf(`{"email":"login-totp@example.com","password":"%s"}`, password)
	req, _ := http.NewRequest("POST", s.server.URL+"/api/auth/login", strings.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	result := s.parseJSONResponse(body)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))
	s.True(result["requires_2fa"].(bool))

	data := result["data"].(map[string]interface{})
	s.NotEmpty(data["temp_token"])
}

func (s *TOTPTestSuite) TestAPI_Verify2FA_InvalidToken() {
	reqBody := `{"temp_token":"invalid-token","code":"123456"}`
	req, _ := http.NewRequest("POST", s.server.URL+"/api/auth/verify-2fa", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	result := s.parseJSONResponse(body)

	s.Equal(http.StatusUnauthorized, resp.StatusCode)
	s.False(result["success"].(bool))
}

func (s *TOTPTestSuite) TestAPI_Verify2FA_InvalidCode() {
	user, password := s.createTestUser("verify2fa-invalid@example.com")

	// Enable TOTP
	key, _ := s.totpService.GenerateSecret(user.Email)
	s.totpService.EnableTOTP(user, key.Secret())

	// Login to get temp token
	loginBody := fmt.Sprintf(`{"email":"verify2fa-invalid@example.com","password":"%s"}`, password)
	req, _ := http.NewRequest("POST", s.server.URL+"/api/auth/login", strings.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := s.client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	loginResult := s.parseJSONResponse(body)
	tempToken := loginResult["data"].(map[string]interface{})["temp_token"].(string)

	// Verify with wrong code
	verifyBody := fmt.Sprintf(`{"temp_token":"%s","code":"000000"}`, tempToken)
	req, _ = http.NewRequest("POST", s.server.URL+"/api/auth/verify-2fa", strings.NewReader(verifyBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ = s.client.Do(req)
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()

	result := s.parseJSONResponse(body)

	s.Equal(http.StatusUnauthorized, resp.StatusCode)
	s.False(result["success"].(bool))
}

func (s *TOTPTestSuite) TestAPI_Verify2FA_WithTOTPCode_Success() {
	user, password := s.createTestUser("verify2fa-totp@example.com")

	// Enable TOTP
	key, _ := s.totpService.GenerateSecret(user.Email)
	s.totpService.EnableTOTP(user, key.Secret())

	// Login to get temp token
	loginBody := fmt.Sprintf(`{"email":"verify2fa-totp@example.com","password":"%s"}`, password)
	req, _ := http.NewRequest("POST", s.server.URL+"/api/auth/login", strings.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := s.client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	loginResult := s.parseJSONResponse(body)
	tempToken := loginResult["data"].(map[string]interface{})["temp_token"].(string)

	// Generate valid code
	validCode, _ := totp.GenerateCode(key.Secret(), time.Now())

	verifyBody := fmt.Sprintf(`{"temp_token":"%s","code":"%s"}`, tempToken, validCode)
	req, _ = http.NewRequest("POST", s.server.URL+"/api/auth/verify-2fa", strings.NewReader(verifyBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ = s.client.Do(req)
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()

	result := s.parseJSONResponse(body)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	data := result["data"].(map[string]interface{})
	s.NotEmpty(data["token"])
}

func (s *TOTPTestSuite) TestAPI_Verify2FA_WithBackupCode_Success() {
	user, password := s.createTestUser("verify2fa-backup@example.com")

	// Enable TOTP and save backup codes
	key, _ := s.totpService.GenerateSecret(user.Email)
	s.totpService.EnableTOTP(user, key.Secret())
	plainCodes, hashedCodes, _ := s.totpService.GenerateBackupCodes(5)
	s.totpService.SaveBackupCodes(user.ID, hashedCodes)

	// Login to get temp token
	loginBody := fmt.Sprintf(`{"email":"verify2fa-backup@example.com","password":"%s"}`, password)
	req, _ := http.NewRequest("POST", s.server.URL+"/api/auth/login", strings.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := s.client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	loginResult := s.parseJSONResponse(body)
	tempToken := loginResult["data"].(map[string]interface{})["temp_token"].(string)

	// Verify with backup code
	verifyBody := fmt.Sprintf(`{"temp_token":"%s","code":"%s"}`, tempToken, plainCodes[0])
	req, _ = http.NewRequest("POST", s.server.URL+"/api/auth/verify-2fa", strings.NewReader(verifyBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ = s.client.Do(req)
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()

	result := s.parseJSONResponse(body)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.True(result["success"].(bool))

	// Verify backup code was consumed
	count, _ := s.totpService.GetBackupCodesRemaining(user.ID)
	s.Equal(int64(4), count)
}

// loginUser2FARequired performs login with 2FA verification
func (s *TOTPTestSuite) loginUser2FARequired(email, password, secret string) *http.Cookie {
	// First login
	loginBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)
	req, _ := http.NewRequest("POST", s.server.URL+"/api/auth/login", strings.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := s.client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	loginResult := s.parseJSONResponse(body)

	// Check if 2FA is required
	if requires2fa, ok := loginResult["requires_2fa"].(bool); ok && requires2fa {
		tempToken := loginResult["data"].(map[string]interface{})["temp_token"].(string)

		// Generate valid code
		validCode, _ := totp.GenerateCode(secret, time.Now())

		verifyBody := fmt.Sprintf(`{"temp_token":"%s","code":"%s"}`, tempToken, validCode)
		req, _ = http.NewRequest("POST", s.server.URL+"/api/auth/verify-2fa", strings.NewReader(verifyBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ = s.client.Do(req)
		resp.Body.Close()

		for _, cookie := range resp.Cookies() {
			if cookie.Name == "token" {
				return cookie
			}
		}
	} else {
		// No 2FA required, get cookie from first response
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "token" {
				return cookie
			}
		}
	}

	return nil
}

// ============================================================================
// Test Runner
// ============================================================================

func TestTOTPSuite(t *testing.T) {
	suite.Run(t, new(TOTPTestSuite))
}
