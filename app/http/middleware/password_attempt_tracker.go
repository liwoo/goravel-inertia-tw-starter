package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/goravel/framework/contracts/cache"
	redisfacades "github.com/goravel/redis/facades"
)

type PasswordAttemptTracker struct {
	redis cache.Driver
}

func NewPasswordAttemptTracker() *PasswordAttemptTracker {
	return &PasswordAttemptTracker{
		redis: redisfacades.Redis("redis"),
	}
}

// AttemptResult represents the result of a password attempt check
type AttemptResult struct {
	IsLocked          bool
	RemainingAttempts int
	LockExpiresAt     *time.Time
	ShouldWarn        bool
	AttemptCount      int
}

// RecordFailedAttempt records a failed login attempt for a user
func (s *PasswordAttemptTracker) RecordFailedAttempt(ctx context.Context, email string) (*AttemptResult, error) {
	// Keys for Redis
	attemptKey := s.getAttemptKey(email)
	lockKey := s.getLockKey(email)
	dailyLockKey := s.getDailyLockKey(email)

	// Check if user is currently locked
	if locked, err := s.isUserLocked(ctx, email); err != nil {
		return nil, err
	} else if locked {
		lockExpiry, _ := s.getLockExpiry(ctx, email)
		return &AttemptResult{
			IsLocked:      true,
			LockExpiresAt: lockExpiry,
		}, nil
	}

	// Get current attempt count and increment
	currentCount := s.redis.GetInt(attemptKey, 0)
	newCount := currentCount + 1

	// Store new count with 5 minute expiry
	err := s.redis.Put(attemptKey, newCount, 5*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to store attempt count: %w", err)
	}

	result := &AttemptResult{
		AttemptCount:      newCount,
		RemainingAttempts: 3 - newCount,
		ShouldWarn:        newCount == 2, // Warn on second attempt
	}

	// Lock user after 3 attempts
	if newCount >= 3 {
		lockDuration, err := s.determineLockDuration(ctx, email)
		if err != nil {
			return nil, err
		}

		// Set lock
		lockExpiry := time.Now().Add(lockDuration)
		err = s.redis.Put(lockKey, lockExpiry.Unix(), lockDuration)
		if err != nil {
			return nil, fmt.Errorf("failed to set lock: %w", err)
		}

		// Track daily locks
		if lockDuration > time.Hour {
			// This is a 24-hour lock, record it
			s.redis.Put(dailyLockKey, time.Now().Unix(), 24*time.Hour)
		}

		// Clear attempt counter since user is now locked
		s.redis.Forget(attemptKey)

		result.IsLocked = true
		result.LockExpiresAt = &lockExpiry
		result.RemainingAttempts = 0
	}

	return result, nil
}

// CheckUserStatus checks if a user is locked and returns their status
func (s *PasswordAttemptTracker) CheckUserStatus(ctx context.Context, email string) (*AttemptResult, error) {
	// Check if user is locked
	if locked, err := s.isUserLocked(ctx, email); err != nil {
		return nil, err
	} else if locked {
		lockExpiry, _ := s.getLockExpiry(ctx, email)
		return &AttemptResult{
			IsLocked:      true,
			LockExpiresAt: lockExpiry,
		}, nil
	}

	// Get current attempt count
	attemptKey := s.getAttemptKey(email)
	count := s.redis.GetInt(attemptKey, 0)

	return &AttemptResult{
		IsLocked:          false,
		RemainingAttempts: 3 - count,
		AttemptCount:      count,
		ShouldWarn:        count == 2,
	}, nil
}

// ClearAttempts clears all failed attempts for a user (called on successful login)
func (s *PasswordAttemptTracker) ClearAttempts(ctx context.Context, email string) error {
	attemptKey := s.getAttemptKey(email)
	s.redis.Forget(attemptKey)
	return nil
}

// isUserLocked checks if a user is currently locked
func (s *PasswordAttemptTracker) isUserLocked(ctx context.Context, email string) (bool, error) {
	lockKey := s.getLockKey(email)
	return s.redis.Has(lockKey), nil
}

// getLockExpiry gets the lock expiry time for a user
func (s *PasswordAttemptTracker) getLockExpiry(ctx context.Context, email string) (*time.Time, error) {
	lockKey := s.getLockKey(email)
	timestamp := s.redis.GetInt64(lockKey, 0)

	if timestamp == 0 {
		return nil, fmt.Errorf("lock not found")
	}

	expiry := time.Unix(timestamp, 0)
	return &expiry, nil
}

// determineLockDuration determines how long to lock a user based on their lock history
func (s *PasswordAttemptTracker) determineLockDuration(ctx context.Context, email string) (time.Duration, error) {
	dailyLockKey := s.getDailyLockKey(email)

	// Check if user has been locked today already
	if s.redis.Has(dailyLockKey) {
		// Second lock in the same day = 24 hours
		return 24 * time.Hour, nil
	}

	// First lock = 1 hour
	return time.Hour, nil
}

// Redis key helpers
func (s *PasswordAttemptTracker) getAttemptKey(email string) string {
	return fmt.Sprintf("password_attempts:%s", email)
}

func (s *PasswordAttemptTracker) getLockKey(email string) string {
	return fmt.Sprintf("user_locked:%s", email)
}

func (s *PasswordAttemptTracker) getDailyLockKey(email string) string {
	today := time.Now().Format("2006-01-02")
	return fmt.Sprintf("daily_lock:%s:%s", email, today)
}
