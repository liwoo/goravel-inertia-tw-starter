package auth

import (
	"net/http"
	"sync"
	"time"

	goravelhttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"smedi-sme-db/app/http/requests"
	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
)

// TOTPController handles TOTP 2FA operations
type TOTPController struct {
	totpService     *services.TOTPService
	activityService *services.UserActivityService

	// In-memory storage for pending TOTP setup (secret key before verification)
	// In production, consider using Redis or database for this
	pendingSetups     map[uint]pendingSetup
	pendingSetupsLock sync.RWMutex
}

// pendingSetup stores the TOTP key during the setup process
type pendingSetup struct {
	Secret    string
	ExpiresAt time.Time
}

// NewTOTPController creates a new TOTP controller
func NewTOTPController() *TOTPController {
	return &TOTPController{
		totpService:     services.NewTOTPService(),
		activityService: services.NewUserActivityService(),
		pendingSetups:   make(map[uint]pendingSetup),
	}
}

// Status godoc
// @Summary      Get 2FA status
// @Description  Get the current 2FA status for the authenticated user
// @Tags         2fa
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /2fa/status [get]
func (c *TOTPController) Status(ctx goravelhttp.Context) goravelhttp.Response {
	// Get authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, goravelhttp.Json{
			"success": false,
			"message": "Unauthorized",
		})
	}

	// Get backup codes remaining count
	backupCodesRemaining, err := c.totpService.GetBackupCodesRemaining(user.ID)
	if err != nil {
		backupCodesRemaining = 0
	}

	return ctx.Response().Json(http.StatusOK, goravelhttp.Json{
		"success": true,
		"data": goravelhttp.Json{
			"enabled":                user.TOTPEnabled,
			"verified_at":            user.TOTPVerifiedAt,
			"backup_codes_remaining": backupCodesRemaining,
		},
	})
}

// Setup godoc
// @Summary      Initiate 2FA setup
// @Description  Start the 2FA setup process - requires password confirmation
// @Tags         2fa
// @Accept       json
// @Produce      json
// @Param        request  body  requests.SetupTOTPRequest  true  "Password confirmation"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /2fa/setup [post]
func (c *TOTPController) Setup(ctx goravelhttp.Context) goravelhttp.Response {
	// Get authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, goravelhttp.Json{
			"success": false,
			"message": "Unauthorized",
		})
	}

	// Validate request
	var request requests.SetupTOTPRequest
	errors, err := ctx.Request().ValidateRequest(&request)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, goravelhttp.Json{
			"success": false,
			"message": "Validation error",
			"error":   err.Error(),
		})
	}
	if errors != nil {
		return ctx.Response().Json(http.StatusUnprocessableEntity, goravelhttp.Json{
			"success": false,
			"message": "Validation failed",
			"errors":  errors.All(),
		})
	}

	// Verify password - need to load the password field
	var userWithPassword models.User
	if err := facades.Orm().Query().Select("id", "password").Find(&userWithPassword, user.ID); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, goravelhttp.Json{
			"success": false,
			"message": "Failed to verify password",
		})
	}

	if !facades.Hash().Check(request.Password, userWithPassword.Password) {
		return ctx.Response().Json(http.StatusUnprocessableEntity, goravelhttp.Json{
			"success": false,
			"message": "Validation failed",
			"errors": map[string][]string{
				"password": {"Password is incorrect"},
			},
		})
	}

	// Check if already enabled
	if user.TOTPEnabled {
		return ctx.Response().Json(http.StatusBadRequest, goravelhttp.Json{
			"success": false,
			"message": "2FA is already enabled for this account",
		})
	}

	// Generate TOTP secret
	key, err := c.totpService.GenerateSecret(user.Email)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, goravelhttp.Json{
			"success": false,
			"message": "Failed to generate 2FA secret",
		})
	}

	// Generate QR code
	qrCode, err := c.totpService.GenerateQRCode(key)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, goravelhttp.Json{
			"success": false,
			"message": "Failed to generate QR code",
		})
	}

	// Store pending setup (expires in 10 minutes)
	c.pendingSetupsLock.Lock()
	c.pendingSetups[user.ID] = pendingSetup{
		Secret:    key.Secret(),
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	c.pendingSetupsLock.Unlock()

	// Clean up expired setups periodically
	go c.cleanupExpiredSetups()

	return ctx.Response().Json(http.StatusOK, goravelhttp.Json{
		"success": true,
		"message": "2FA setup initiated",
		"data": goravelhttp.Json{
			"qr_code": qrCode,
			"secret":  key.Secret(), // For manual entry
		},
	})
}

// Verify godoc
// @Summary      Verify and enable 2FA
// @Description  Verify the TOTP code and enable 2FA, returning backup codes
// @Tags         2fa
// @Accept       json
// @Produce      json
// @Param        request  body  requests.VerifyTOTPRequest  true  "TOTP code"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /2fa/verify [post]
func (c *TOTPController) Verify(ctx goravelhttp.Context) goravelhttp.Response {
	// Get authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, goravelhttp.Json{
			"success": false,
			"message": "Unauthorized",
		})
	}

	// Validate request
	var request requests.VerifyTOTPRequest
	errors, err := ctx.Request().ValidateRequest(&request)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, goravelhttp.Json{
			"success": false,
			"message": "Validation error",
			"error":   err.Error(),
		})
	}
	if errors != nil {
		return ctx.Response().Json(http.StatusUnprocessableEntity, goravelhttp.Json{
			"success": false,
			"message": "Validation failed",
			"errors":  errors.All(),
		})
	}

	// Get pending setup
	c.pendingSetupsLock.RLock()
	setup, exists := c.pendingSetups[user.ID]
	c.pendingSetupsLock.RUnlock()

	if !exists || time.Now().After(setup.ExpiresAt) {
		return ctx.Response().Json(http.StatusBadRequest, goravelhttp.Json{
			"success": false,
			"message": "2FA setup has expired. Please start the setup process again.",
		})
	}

	// Validate the TOTP code
	if !c.totpService.ValidateCode(setup.Secret, request.Code) {
		return ctx.Response().Json(http.StatusUnprocessableEntity, goravelhttp.Json{
			"success": false,
			"message": "Validation failed",
			"errors": map[string][]string{
				"code": {"Invalid verification code"},
			},
		})
	}

	// Generate backup codes
	plainCodes, hashedCodes, err := c.totpService.GenerateBackupCodes(10)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, goravelhttp.Json{
			"success": false,
			"message": "Failed to generate backup codes",
		})
	}

	// Save backup codes
	if err := c.totpService.SaveBackupCodes(user.ID, hashedCodes); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, goravelhttp.Json{
			"success": false,
			"message": "Failed to save backup codes",
		})
	}

	// Enable TOTP
	if err := c.totpService.EnableTOTP(&user, setup.Secret); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, goravelhttp.Json{
			"success": false,
			"message": "Failed to enable 2FA",
		})
	}

	// Remove pending setup
	c.pendingSetupsLock.Lock()
	delete(c.pendingSetups, user.ID)
	c.pendingSetupsLock.Unlock()

	// Log activity
	c.activityService.LogActivity(ctx, user.ID, models.ActivityTwoFAEnabled, "Two-factor authentication enabled", nil)

	return ctx.Response().Json(http.StatusOK, goravelhttp.Json{
		"success": true,
		"message": "2FA has been enabled successfully",
		"data": goravelhttp.Json{
			"backup_codes": plainCodes,
			"enabled_at":   time.Now(),
		},
	})
}

// Disable godoc
// @Summary      Disable 2FA
// @Description  Disable 2FA - requires password and current TOTP code
// @Tags         2fa
// @Accept       json
// @Produce      json
// @Param        request  body  requests.DisableTOTPRequest  true  "Password and TOTP code"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /2fa/disable [post]
func (c *TOTPController) Disable(ctx goravelhttp.Context) goravelhttp.Response {
	// Get authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, goravelhttp.Json{
			"success": false,
			"message": "Unauthorized",
		})
	}

	// Validate request
	var request requests.DisableTOTPRequest
	errors, err := ctx.Request().ValidateRequest(&request)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, goravelhttp.Json{
			"success": false,
			"message": "Validation error",
			"error":   err.Error(),
		})
	}
	if errors != nil {
		return ctx.Response().Json(http.StatusUnprocessableEntity, goravelhttp.Json{
			"success": false,
			"message": "Validation failed",
			"errors":  errors.All(),
		})
	}

	// Check if 2FA is enabled
	if !user.TOTPEnabled {
		return ctx.Response().Json(http.StatusBadRequest, goravelhttp.Json{
			"success": false,
			"message": "2FA is not enabled for this account",
		})
	}

	// Verify password
	var userWithPassword models.User
	if err := facades.Orm().Query().Select("id", "password", "totp_secret").Find(&userWithPassword, user.ID); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, goravelhttp.Json{
			"success": false,
			"message": "Failed to verify password",
		})
	}

	if !facades.Hash().Check(request.Password, userWithPassword.Password) {
		return ctx.Response().Json(http.StatusUnprocessableEntity, goravelhttp.Json{
			"success": false,
			"message": "Validation failed",
			"errors": map[string][]string{
				"password": {"Password is incorrect"},
			},
		})
	}

	// Load the full user for TOTP validation
	user.TOTPSecret = userWithPassword.TOTPSecret

	// Validate TOTP code or backup code
	valid, isBackup, err := c.totpService.ValidateTOTPForUser(&user, request.Code)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, goravelhttp.Json{
			"success": false,
			"message": "Failed to validate 2FA code",
		})
	}

	if !valid {
		return ctx.Response().Json(http.StatusUnprocessableEntity, goravelhttp.Json{
			"success": false,
			"message": "Validation failed",
			"errors": map[string][]string{
				"code": {"Invalid 2FA code"},
			},
		})
	}

	// Disable TOTP
	if err := c.totpService.DisableTOTP(user.ID); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, goravelhttp.Json{
			"success": false,
			"message": "Failed to disable 2FA",
		})
	}

	// Log activity
	metadata := map[string]interface{}{
		"used_backup_code": isBackup,
	}
	c.activityService.LogActivity(ctx, user.ID, models.ActivityTwoFADisabled, "Two-factor authentication disabled", metadata)

	return ctx.Response().Json(http.StatusOK, goravelhttp.Json{
		"success": true,
		"message": "2FA has been disabled successfully",
	})
}

// RegenerateBackupCodes godoc
// @Summary      Regenerate backup codes
// @Description  Generate new backup codes - requires password confirmation
// @Tags         2fa
// @Accept       json
// @Produce      json
// @Param        request  body  requests.RegenerateBackupCodesRequest  true  "Password confirmation"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /2fa/backup-codes [post]
func (c *TOTPController) RegenerateBackupCodes(ctx goravelhttp.Context) goravelhttp.Response {
	// Get authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, goravelhttp.Json{
			"success": false,
			"message": "Unauthorized",
		})
	}

	// Validate request
	var request requests.RegenerateBackupCodesRequest
	errors, err := ctx.Request().ValidateRequest(&request)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, goravelhttp.Json{
			"success": false,
			"message": "Validation error",
			"error":   err.Error(),
		})
	}
	if errors != nil {
		return ctx.Response().Json(http.StatusUnprocessableEntity, goravelhttp.Json{
			"success": false,
			"message": "Validation failed",
			"errors":  errors.All(),
		})
	}

	// Check if 2FA is enabled
	if !user.TOTPEnabled {
		return ctx.Response().Json(http.StatusBadRequest, goravelhttp.Json{
			"success": false,
			"message": "2FA is not enabled for this account",
		})
	}

	// Verify password
	var userWithPassword models.User
	if err := facades.Orm().Query().Select("id", "password").Find(&userWithPassword, user.ID); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, goravelhttp.Json{
			"success": false,
			"message": "Failed to verify password",
		})
	}

	if !facades.Hash().Check(request.Password, userWithPassword.Password) {
		return ctx.Response().Json(http.StatusUnprocessableEntity, goravelhttp.Json{
			"success": false,
			"message": "Validation failed",
			"errors": map[string][]string{
				"password": {"Password is incorrect"},
			},
		})
	}

	// Generate new backup codes
	plainCodes, hashedCodes, err := c.totpService.GenerateBackupCodes(10)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, goravelhttp.Json{
			"success": false,
			"message": "Failed to generate backup codes",
		})
	}

	// Save backup codes (replaces existing)
	if err := c.totpService.SaveBackupCodes(user.ID, hashedCodes); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, goravelhttp.Json{
			"success": false,
			"message": "Failed to save backup codes",
		})
	}

	// Log activity
	c.activityService.LogActivity(ctx, user.ID, "backup_codes_regenerated", "Backup codes regenerated", nil)

	return ctx.Response().Json(http.StatusOK, goravelhttp.Json{
		"success": true,
		"message": "Backup codes have been regenerated",
		"data": goravelhttp.Json{
			"backup_codes": plainCodes,
			"generated_at": time.Now(),
		},
	})
}

// cleanupExpiredSetups removes expired pending TOTP setups
func (c *TOTPController) cleanupExpiredSetups() {
	c.pendingSetupsLock.Lock()
	defer c.pendingSetupsLock.Unlock()

	now := time.Now()
	for userID, setup := range c.pendingSetups {
		if now.After(setup.ExpiresAt) {
			delete(c.pendingSetups, userID)
		}
	}
}
