package services

import (
	"context"
	"fmt"
	"time"

	"github.com/goravel/framework/contracts/cache"
	"github.com/goravel/framework/facades"
)

type PasswordAttemptService struct {
	cache cache.Driver
}

func NewPasswordAttemptService() *PasswordAttemptService {
	return &PasswordAttemptService{
		cache: facades.Cache(),
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
func (s *PasswordAttemptService) RecordFailedAttempt(ctx context.Context, email string) (*AttemptResult, error) {
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
	currentCount := s.cache.GetInt(attemptKey, 0)
	newCount := currentCount + 1

	// Store new count with 5 minute expiry
	err := s.cache.Put(attemptKey, newCount, 5*time.Minute)
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
		err = s.cache.Put(lockKey, lockExpiry.Unix(), lockDuration)
		if err != nil {
			return nil, fmt.Errorf("failed to set lock: %w", err)
		}

		// Track daily locks
		if lockDuration > time.Hour {
			// This is a 24-hour lock, record it
			s.cache.Put(dailyLockKey, time.Now().Unix(), 24*time.Hour)
		}

		// Clear attempt counter since user is now locked
		s.cache.Forget(attemptKey)

		result.IsLocked = true
		result.LockExpiresAt = &lockExpiry
		result.RemainingAttempts = 0
	}

	return result, nil
}

// CheckUserStatus checks if a user is locked and returns their status
func (s *PasswordAttemptService) CheckUserStatus(ctx context.Context, email string) (*AttemptResult, error) {
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
	count := s.cache.GetInt(attemptKey, 0)

	return &AttemptResult{
		IsLocked:          false,
		AttemptCount:      count,
		RemainingAttempts: 3 - count,
		ShouldWarn:        count == 2,
	}, nil
}

// ClearAttempts clears all failed attempts for a user (call on successful login)
func (s *PasswordAttemptService) ClearAttempts(ctx context.Context, email string) error {
	attemptKey := s.getAttemptKey(email)
	s.cache.Forget(attemptKey)
	return nil
}

// Helper methods
func (s *PasswordAttemptService) getAttemptKey(email string) string {
	return fmt.Sprintf("password_attempts:%s", email)
}

func (s *PasswordAttemptService) getLockKey(email string) string {
	return fmt.Sprintf("password_lock:%s", email)
}

func (s *PasswordAttemptService) getDailyLockKey(email string) string {
	return fmt.Sprintf("daily_lock:%s", email)
}

func (s *PasswordAttemptService) isUserLocked(ctx context.Context, email string) (bool, error) {
	lockKey := s.getLockKey(email)
	return s.cache.Has(lockKey), nil
}

func (s *PasswordAttemptService) getLockExpiry(ctx context.Context, email string) (*time.Time, error) {
	lockKey := s.getLockKey(email)
	timestamp := s.cache.GetInt64(lockKey, 0)

	if timestamp == 0 {
		return nil, fmt.Errorf("lock not found")
	}

	expiry := time.Unix(timestamp, 0)
	return &expiry, nil
}

func (s *PasswordAttemptService) determineLockDuration(ctx context.Context, email string) (time.Duration, error) {
	dailyLockKey := s.getDailyLockKey(email)

	// Check if user has been locked today already
	if s.cache.Has(dailyLockKey) {
		// Second lock in the same day = 24 hours
		return 24 * time.Hour, nil
	}

	// First lock = 1 hour
	return time.Hour, nil
}
