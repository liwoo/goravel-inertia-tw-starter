package services

import (
	"fmt"
	"strconv"
	"time"

	redisfacades "github.com/goravel/redis/facades"
)

const (
	// AttemptWindow is the time window for tracking attempts (5 minutes)
	AttemptWindow = 5 * time.Minute
	// FirstLockoutDuration is the duration for the first lockout (1 hour)
	FirstLockoutDuration = 1 * time.Hour
	// SecondLockoutDuration is the duration for the second lockout in the same day (24 hours)
	SecondLockoutDuration = 24 * time.Hour
	// MaxAttempts is the maximum number of attempts before lockout
	MaxAttempts = 3
)

// PasswordAttemptsService handles password attempt tracking using Redis
type PasswordAttemptsService struct{}

// NewPasswordAttemptsService creates a new password attempts service
func NewPasswordAttemptsService() *PasswordAttemptsService {
	return &PasswordAttemptsService{}
}

// CheckLocked checks if an account is locked and returns the unlock time
func (s *PasswordAttemptsService) CheckLocked(email string) (bool, time.Time) {
	cacheDriver, err := redisfacades.Cache("default")
	if err != nil {
		return false, time.Time{}
	}

	lockoutKey := s.getLockoutKey(email)

	unlockTimeStr := cacheDriver.GetString(lockoutKey)
	if unlockTimeStr == "" {
		return false, time.Time{}
	}

	unlockTimestamp, err := strconv.ParseInt(unlockTimeStr, 10, 64)
	if err != nil {
		return false, time.Time{}
	}

	unlockTime := time.Unix(unlockTimestamp, 0).UTC()
	now := time.Now().UTC()

	if now.Before(unlockTime) {
		return true, unlockTime
	}

	// Lockout has expired, clean it up
	cacheDriver.Forget(lockoutKey)
	return false, time.Time{}
}

// RecordFailedAttempt records a failed login attempt and returns the current attempt count and whether to warn
func (s *PasswordAttemptsService) RecordFailedAttempt(email string) (int, bool) {
	cacheDriver, err := redisfacades.Cache("default")
	if err != nil {
		return 0, false
	}

	attemptsKey := s.getAttemptsKey(email)

	// Check if account is already locked
	locked, _ := s.CheckLocked(email)
	if locked {
		return MaxAttempts, false
	}

	// Get current attempt count
	count := cacheDriver.GetInt(attemptsKey, 0)

	// Increment count
	count++
	newCount := int64(count)

	// Set the count with TTL of 5 minutes (attempt window)
	cacheDriver.Put(attemptsKey, newCount, AttemptWindow)

	// If this is the 3rd attempt, lock the account
	if count >= MaxAttempts {
		s.lockAccount(email)
		return count, false
	}

	// Return count and whether to warn (2nd attempt)
	return count, count == 2
}

// ClearAttempts clears all tracking data for an email (called on successful login)
func (s *PasswordAttemptsService) ClearAttempts(email string) {
	cacheDriver, err := redisfacades.Cache("default")
	if err != nil {
		return
	}

	attemptsKey := s.getAttemptsKey(email)
	lockoutKey := s.getLockoutKey(email)

	cacheDriver.Forget(attemptsKey)
	cacheDriver.Forget(lockoutKey)
	// Note: We don't clear locksTodayKey as we want to track daily lock count
}

// GetRemainingAttempts returns the number of remaining attempts before lockout
func (s *PasswordAttemptsService) GetRemainingAttempts(email string) int {
	// Check if locked
	locked, _ := s.CheckLocked(email)
	if locked {
		return 0
	}

	cacheDriver, err := redisfacades.Cache("default")
	if err != nil {
		return MaxAttempts
	}

	attemptsKey := s.getAttemptsKey(email)
	count := cacheDriver.GetInt(attemptsKey, 0)

	remaining := MaxAttempts - count
	if remaining < 0 {
		return 0
	}

	return remaining
}

// lockAccount locks the account for the appropriate duration
func (s *PasswordAttemptsService) lockAccount(email string) {
	cacheDriver, err := redisfacades.Cache("default")
	if err != nil {
		return
	}

	lockoutKey := s.getLockoutKey(email)
	locksTodayKey := s.getLocksTodayKey(email)

	// Get current lock count for today
	locksToday := cacheDriver.GetInt(locksTodayKey, 0)

	// Increment lock count
	locksToday++
	cacheDriver.Put(locksTodayKey, int64(locksToday), s.getTimeUntilMidnight())

	// Determine lockout duration
	var lockoutDuration time.Duration
	if locksToday >= 2 {
		// Second lock in the same day = 24 hour lock
		lockoutDuration = SecondLockoutDuration
	} else {
		// First lock = 1 hour lock
		lockoutDuration = FirstLockoutDuration
	}

	// Set unlock time
	unlockTime := time.Now().UTC().Add(lockoutDuration)
	unlockTimestamp := unlockTime.Unix()
	cacheDriver.Put(lockoutKey, unlockTimestamp, lockoutDuration)

	// Clear attempt counter since we're locking
	attemptsKey := s.getAttemptsKey(email)
	cacheDriver.Forget(attemptsKey)
}

// getTimeUntilMidnight returns the duration until midnight UTC
func (s *PasswordAttemptsService) getTimeUntilMidnight() time.Duration {
	now := time.Now().UTC()
	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	return midnight.Sub(now)
}

// getAttemptsKey returns the Redis key for tracking attempts
func (s *PasswordAttemptsService) getAttemptsKey(email string) string {
	return fmt.Sprintf("password_attempts:%s", email)
}

// getLockoutKey returns the Redis key for lockout status
func (s *PasswordAttemptsService) getLockoutKey(email string) string {
	return fmt.Sprintf("password_lockout:%s", email)
}

// getLocksTodayKey returns the Redis key for tracking locks today
func (s *PasswordAttemptsService) getLocksTodayKey(email string) string {
	return fmt.Sprintf("password_locks_today:%s", email)
}
