package services

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"image/png"
	"io"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"starter-project/app/models"
)

// TOTPService handles TOTP 2FA operations
type TOTPService struct {
	activityService *UserActivityService
}

// NewTOTPService creates a new TOTP service
func NewTOTPService() *TOTPService {
	return &TOTPService{
		activityService: NewUserActivityService(),
	}
}

// GenerateSecret generates a new TOTP secret for a user
func (s *TOTPService) GenerateSecret(email string) (*otp.Key, error) {
	appName := facades.Config().GetString("app.name", "Goravel App")

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      appName,
		AccountName: email,
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})

	if err != nil {
		facades.Log().Error("Failed to generate TOTP secret", map[string]interface{}{
			"error": err.Error(),
			"email": email,
		})
		return nil, fmt.Errorf("failed to generate TOTP secret: %w", err)
	}

	return key, nil
}

// GenerateQRCode generates a QR code image as base64 PNG for the TOTP key
func (s *TOTPService) GenerateQRCode(key *otp.Key) (string, error) {
	// Generate QR code image
	img, err := key.Image(200, 200)
	if err != nil {
		facades.Log().Error("Failed to generate QR code image", map[string]interface{}{
			"error": err.Error(),
		})
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Encode image to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		facades.Log().Error("Failed to encode QR code to PNG", map[string]interface{}{
			"error": err.Error(),
		})
		return "", fmt.Errorf("failed to encode QR code: %w", err)
	}

	// Convert to base64
	base64Img := base64.StdEncoding.EncodeToString(buf.Bytes())
	return fmt.Sprintf("data:image/png;base64,%s", base64Img), nil
}

// ValidateCode validates a TOTP code against a secret
func (s *TOTPService) ValidateCode(secret, code string) bool {
	return totp.Validate(code, secret)
}

// GenerateBackupCodes generates backup codes for a user
// Returns (plaintext codes for display, hashed codes for storage)
func (s *TOTPService) GenerateBackupCodes(count int) ([]string, []string, error) {
	if count <= 0 {
		count = 10
	}

	plainCodes := make([]string, count)
	hashedCodes := make([]string, count)

	for i := 0; i < count; i++ {
		// Generate an 8-character alphanumeric code using crypto/rand
		code, err := s.generateSecureRandomCode(8)
		if err != nil {
			facades.Log().Error("Failed to generate backup code", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, nil, fmt.Errorf("failed to generate backup code: %w", err)
		}
		plainCodes[i] = code

		// Hash the code with bcrypt
		hashedCode, err := facades.Hash().Make(code)
		if err != nil {
			facades.Log().Error("Failed to hash backup code", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, nil, fmt.Errorf("failed to hash backup code: %w", err)
		}
		hashedCodes[i] = hashedCode
	}

	return plainCodes, hashedCodes, nil
}

// generateSecureRandomCode generates a cryptographically secure random alphanumeric code
func (s *TOTPService) generateSecureRandomCode(length int) (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // Excluding confusing chars like 0, O, 1, I
	code := make([]byte, length)

	// Use crypto/rand for secure random number generation
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	for i := range code {
		code[i] = charset[int(randomBytes[i])%len(charset)]
	}
	return string(code), nil
}

// ValidateBackupCode checks if a backup code is valid and marks it as used
func (s *TOTPService) ValidateBackupCode(userID uint, code string) (bool, error) {
	var backupCodes []models.TOTPBackupCode

	// Find all unused backup codes for the user
	err := facades.Orm().Query().
		Where("user_id = ?", userID).
		Where("used_at IS NULL").
		Find(&backupCodes)

	if err != nil {
		facades.Log().Error("Failed to fetch backup codes", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return false, fmt.Errorf("failed to fetch backup codes: %w", err)
	}

	// Check each unused code
	for _, bc := range backupCodes {
		if facades.Hash().Check(code, bc.CodeHash) {
			// Mark as used
			now := time.Now()
			bc.UsedAt = &now

			if _, err := facades.Orm().Query().Model(&bc).Where("id = ?", bc.ID).Update("used_at", now); err != nil {
				facades.Log().Error("Failed to mark backup code as used", map[string]interface{}{
					"error":   err.Error(),
					"code_id": bc.ID,
				})
				return false, fmt.Errorf("failed to mark backup code as used: %w", err)
			}

			facades.Log().Info("Backup code used", map[string]interface{}{
				"user_id": userID,
				"code_id": bc.ID,
			})

			return true, nil
		}
	}

	return false, nil
}

// SaveBackupCodes saves hashed backup codes for a user (deletes existing ones first)
func (s *TOTPService) SaveBackupCodes(userID uint, hashedCodes []string) error {
	// Delete existing backup codes for the user
	if _, err := facades.Orm().Query().Where("user_id = ?", userID).Delete(&models.TOTPBackupCode{}); err != nil {
		facades.Log().Error("Failed to delete existing backup codes", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return fmt.Errorf("failed to delete existing backup codes: %w", err)
	}

	// Create new backup codes
	for _, hash := range hashedCodes {
		backupCode := models.TOTPBackupCode{
			UserID:   userID,
			CodeHash: hash,
		}

		if err := facades.Orm().Query().Create(&backupCode); err != nil {
			facades.Log().Error("Failed to save backup code", map[string]interface{}{
				"error":   err.Error(),
				"user_id": userID,
			})
			return fmt.Errorf("failed to save backup code: %w", err)
		}
	}

	facades.Log().Info("Backup codes saved", map[string]interface{}{
		"user_id": userID,
		"count":   len(hashedCodes),
	})

	return nil
}

// GetBackupCodesRemaining returns the count of unused backup codes for a user
func (s *TOTPService) GetBackupCodesRemaining(userID uint) (int64, error) {
	count, err := facades.Orm().Query().
		Model(&models.TOTPBackupCode{}).
		Where("user_id = ?", userID).
		Where("used_at IS NULL").
		Count()

	if err != nil {
		facades.Log().Error("Failed to count backup codes", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return 0, fmt.Errorf("failed to count backup codes: %w", err)
	}

	return count, nil
}

// EncryptSecret encrypts a TOTP secret using AES-256-GCM with the APP_KEY
func (s *TOTPService) EncryptSecret(secret string) (string, error) {
	appKey := facades.Config().GetString("app.key", "")
	if appKey == "" {
		return "", fmt.Errorf("APP_KEY is not set")
	}

	// Derive a 32-byte key from APP_KEY using SHA-256
	hash := sha256.Sum256([]byte(appKey))
	key := hash[:]

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, []byte(secret), nil)

	// Encode to base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptSecret decrypts a TOTP secret encrypted with EncryptSecret
func (s *TOTPService) DecryptSecret(encrypted string) (string, error) {
	appKey := facades.Config().GetString("app.key", "")
	if appKey == "" {
		return "", fmt.Errorf("APP_KEY is not set")
	}

	// Decode from base64
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Derive a 32-byte key from APP_KEY using SHA-256
	hash := sha256.Sum256([]byte(appKey))
	key := hash[:]

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce and ciphertext
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// EnableTOTP enables TOTP for a user
func (s *TOTPService) EnableTOTP(user *models.User, secret string) error {
	// Encrypt the secret
	encryptedSecret, err := s.EncryptSecret(secret)
	if err != nil {
		facades.Log().Error("Failed to encrypt TOTP secret", map[string]interface{}{
			"error":   err.Error(),
			"user_id": user.ID,
		})
		return fmt.Errorf("failed to encrypt TOTP secret: %w", err)
	}

	now := time.Now()

	// Update user
	_, err = facades.Orm().Query().Model(&models.User{}).Where("id = ?", user.ID).Update(map[string]interface{}{
		"totp_enabled":     true,
		"totp_secret":      encryptedSecret,
		"totp_verified_at": now,
	})

	if err != nil {
		facades.Log().Error("Failed to enable TOTP", map[string]interface{}{
			"error":   err.Error(),
			"user_id": user.ID,
		})
		return fmt.Errorf("failed to enable TOTP: %w", err)
	}

	facades.Log().Info("TOTP enabled for user", map[string]interface{}{
		"user_id": user.ID,
	})

	return nil
}

// DisableTOTP disables TOTP for a user
func (s *TOTPService) DisableTOTP(userID uint) error {
	// Update user to disable TOTP
	_, err := facades.Orm().Query().Model(&models.User{}).Where("id = ?", userID).Update(map[string]interface{}{
		"totp_enabled":     false,
		"totp_secret":      "",
		"totp_verified_at": nil,
	})

	if err != nil {
		facades.Log().Error("Failed to disable TOTP", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return fmt.Errorf("failed to disable TOTP: %w", err)
	}

	// Delete backup codes
	if _, err := facades.Orm().Query().Where("user_id = ?", userID).Delete(&models.TOTPBackupCode{}); err != nil {
		facades.Log().Warning("Failed to delete backup codes on TOTP disable", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		// Don't fail the operation for this
	}

	facades.Log().Info("TOTP disabled for user", map[string]interface{}{
		"user_id": userID,
	})

	return nil
}

// GetDecryptedSecret gets and decrypts the TOTP secret for a user
func (s *TOTPService) GetDecryptedSecret(user *models.User) (string, error) {
	if user.TOTPSecret == "" {
		return "", fmt.Errorf("no TOTP secret found")
	}

	secret, err := s.DecryptSecret(user.TOTPSecret)
	if err != nil {
		facades.Log().Error("Failed to decrypt TOTP secret", map[string]interface{}{
			"error":   err.Error(),
			"user_id": user.ID,
		})
		return "", fmt.Errorf("failed to decrypt TOTP secret: %w", err)
	}

	return secret, nil
}

// ValidateTOTPForUser validates a TOTP code or backup code for a user
func (s *TOTPService) ValidateTOTPForUser(user *models.User, code string) (bool, bool, error) {
	// First try TOTP code
	if user.TOTPEnabled && user.TOTPSecret != "" {
		secret, err := s.GetDecryptedSecret(user)
		if err != nil {
			return false, false, err
		}

		if s.ValidateCode(secret, code) {
			return true, false, nil // Valid TOTP code
		}
	}

	// Try backup code
	isBackupValid, err := s.ValidateBackupCode(user.ID, code)
	if err != nil {
		return false, false, err
	}

	if isBackupValid {
		return true, true, nil // Valid backup code
	}

	return false, false, nil // Invalid code
}
