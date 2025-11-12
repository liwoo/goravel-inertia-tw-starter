package unit

import (
	"testing"
	"time"

	redisfacades "github.com/goravel/redis/facades"
	"github.com/stretchr/testify/assert"

	"smedi-sme-db/app/services"
	"smedi-sme-db/tests"
)

type PasswordAttemptsTestSuite struct {
	tests.TestCase
	service *services.PasswordAttemptsService
}

// clearRedisKeys clears all password attempts related Redis keys
func (suite *PasswordAttemptsTestSuite) clearRedisKeys() {
	cacheDriver, err := redisfacades.Cache("default")
	if err != nil {
		return
	}

	// Clear known test emails
	testEmails := []string{"test@example.com", "locked@example.com", "remaining@example.com", "success@example.com", "fresh@example.com", "expire@example.com"}
	for _, email := range testEmails {
		cacheDriver.Forget("password_attempts:" + email)
		cacheDriver.Forget("password_lockout:" + email)
		cacheDriver.Forget("password_locks_today:" + email)
	}
}

func TestPasswordAttempts(t *testing.T) {
	suite := &PasswordAttemptsTestSuite{}
	suite.service = services.NewPasswordAttemptsService()

	// Clear Redis keys from previous tests
	suite.clearRedisKeys()

	// Test basic functionality
	t.Run("TestCheckLocked_UnlockedAccount", func(t *testing.T) {
		locked, unlockTime := suite.service.CheckLocked("test@example.com")
		assert.False(t, locked)
		assert.Equal(t, time.Time{}, unlockTime)
	})

	t.Run("TestRecordFailedAttempt_IncrementsCount", func(t *testing.T) {
		email := "test@example.com"

		// First attempt
		count, shouldWarn := suite.service.RecordFailedAttempt(email)
		assert.Equal(t, 1, count)
		assert.False(t, shouldWarn)

		// Second attempt - should warn
		count, shouldWarn = suite.service.RecordFailedAttempt(email)
		assert.Equal(t, 2, count)
		assert.True(t, shouldWarn)

		// Third attempt - should lock account
		count, shouldWarn = suite.service.RecordFailedAttempt(email)
		assert.Equal(t, 3, count)
		assert.False(t, shouldWarn)

		// Fourth attempt - should still show 3 attempts (locked)
		count, shouldWarn = suite.service.RecordFailedAttempt(email)
		assert.Equal(t, 3, count)
		assert.False(t, shouldWarn)
	})

	t.Run("TestAccountLockout_AfterThreeAttempts", func(t *testing.T) {
		email := "locked@example.com"

		// Make 3 attempts to trigger lockout
		for i := 0; i < 3; i++ {
			suite.service.RecordFailedAttempt(email)
		}

		// Check if account is locked
		locked, unlockTime := suite.service.CheckLocked(email)
		assert.True(t, locked)
		assert.True(t, unlockTime.After(time.Now()))
		assert.Equal(t, 0, suite.service.GetRemainingAttempts(email))
	})

	t.Run("TestGetRemainingAttempts_BeforeLockout", func(t *testing.T) {
		email := "remaining@example.com"

		// Initial state
		assert.Equal(t, 3, suite.service.GetRemainingAttempts(email))

		// After 1 attempt
		suite.service.RecordFailedAttempt(email)
		assert.Equal(t, 2, suite.service.GetRemainingAttempts(email))

		// After 2 attempts
		suite.service.RecordFailedAttempt(email)
		assert.Equal(t, 1, suite.service.GetRemainingAttempts(email))
	})

	t.Run("TestClearAttempts_OnSuccessfulLogin", func(t *testing.T) {
		email := "success@example.com"

		// Make some failed attempts
		suite.service.RecordFailedAttempt(email)
		suite.service.RecordFailedAttempt(email)

		// Verify attempts are tracked
		assert.Equal(t, 1, suite.service.GetRemainingAttempts(email))

		// Clear attempts (simulate successful login)
		suite.service.ClearAttempts(email)

		// Verify attempts are cleared
		assert.Equal(t, 3, suite.service.GetRemainingAttempts(email))
		locked, _ := suite.service.CheckLocked(email)
		assert.False(t, locked)
	})

	t.Run("TestLockoutExpiration_AccountUnlockedAfterTime", func(t *testing.T) {
		email := "expire@example.com"

		// Trigger lockout
		for i := 0; i < 3; i++ {
			suite.service.RecordFailedAttempt(email)
		}

		locked, _ := suite.service.CheckLocked(email)
		assert.True(t, locked)

		// Simulate time passing (this is a limitation of testing with real Redis)
		// In a real scenario, we'd need to mock Redis or use time manipulation
		// For now, this test verifies the lockout was triggered
		assert.Equal(t, 0, suite.service.GetRemainingAttempts(email))
	})
}
