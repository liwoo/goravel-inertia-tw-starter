package services

import (
	"encoding/json"

	"starter-project/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

// UserActivityService handles user activity logging and retrieval
type UserActivityService struct{}

// NewUserActivityService creates a new user activity service
func NewUserActivityService() *UserActivityService {
	return &UserActivityService{}
}

// LogActivity logs a user activity
func (s *UserActivityService) LogActivity(
	ctx http.Context,
	userID uint,
	activityType string,
	description string,
	metadata map[string]interface{},
) error {
	activity := models.UserActivity{
		UserID:       userID,
		ActivityType: activityType,
		Description:  description,
	}

	// Set metadata if provided
	if metadata != nil {
		if err := activity.SetMetadata(metadata); err != nil {
			facades.Log().Warning("Failed to set activity metadata", map[string]interface{}{
				"error":   err.Error(),
				"user_id": userID,
			})
		}
	}

	// Extract IP address and user agent from context if available
	if ctx != nil {
		activity.IPAddress = ctx.Request().Ip()
		activity.UserAgent = ctx.Request().Header("User-Agent", "")
	}

	// Save the activity
	if err := facades.Orm().Query().Create(&activity); err != nil {
		facades.Log().Error("Failed to log user activity", map[string]interface{}{
			"error":         err.Error(),
			"user_id":       userID,
			"activity_type": activityType,
		})
		return err
	}

	facades.Log().Info("User activity logged", map[string]interface{}{
		"user_id":       userID,
		"activity_type": activityType,
		"description":   description,
	})

	return nil
}

// LogActivityWithRelation logs an activity with a related entity
func (s *UserActivityService) LogActivityWithRelation(
	ctx http.Context,
	userID uint,
	activityType string,
	description string,
	metadata map[string]interface{},
	relatedType string,
	relatedID uint,
) error {
	activity := models.UserActivity{
		UserID:       userID,
		ActivityType: activityType,
		Description:  description,
		RelatedType:  relatedType,
		RelatedID:    &relatedID,
	}

	// Set metadata if provided
	if metadata != nil {
		if err := activity.SetMetadata(metadata); err != nil {
			facades.Log().Warning("Failed to set activity metadata", map[string]interface{}{
				"error":   err.Error(),
				"user_id": userID,
			})
		}
	}

	// Extract IP address and user agent from context if available
	if ctx != nil {
		activity.IPAddress = ctx.Request().Ip()
		activity.UserAgent = ctx.Request().Header("User-Agent", "")
	}

	// Save the activity
	if err := facades.Orm().Query().Create(&activity); err != nil {
		facades.Log().Error("Failed to log user activity", map[string]interface{}{
			"error":         err.Error(),
			"user_id":       userID,
			"activity_type": activityType,
		})
		return err
	}

	return nil
}

// GetUserActivitiesResult represents a paginated result of user activities
type GetUserActivitiesResult struct {
	Data       []models.UserActivity `json:"data"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"pageSize"`
	TotalPages int                   `json:"totalPages"`
}

// GetUserActivities retrieves paginated activities for a user
func (s *UserActivityService) GetUserActivities(
	userID uint,
	page int,
	pageSize int,
	activityTypeFilter string,
) (*GetUserActivitiesResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var activities []models.UserActivity
	var total int64

	// Build query
	query := facades.Orm().Query().
		Model(&models.UserActivity{}).
		Where("user_id = ?", userID)

	// Apply activity type filter if provided
	if activityTypeFilter != "" {
		query = query.Where("activity_type = ?", activityTypeFilter)
	}

	// Get total count
	total, err := query.Count()
	if err != nil {
		return nil, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Fetch activities with pagination
	err = query.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&activities)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &GetUserActivitiesResult{
		Data:       activities,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetRecentActivities retrieves the most recent activities for a user
func (s *UserActivityService) GetRecentActivities(userID uint, limit int) ([]models.UserActivity, error) {
	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	var activities []models.UserActivity
	err := facades.Orm().Query().
		Model(&models.UserActivity{}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&activities)

	if err != nil {
		return nil, err
	}

	return activities, nil
}

// GetActivityTypes returns a list of distinct activity types for a user
func (s *UserActivityService) GetActivityTypes(userID uint) ([]string, error) {
	var activityTypes []string

	err := facades.Orm().Query().
		Model(&models.UserActivity{}).
		Where("user_id = ?", userID).
		Distinct("activity_type").
		Pluck("activity_type", &activityTypes)

	if err != nil {
		return nil, err
	}

	return activityTypes, nil
}

// DeleteOldActivities deletes activities older than the specified days
func (s *UserActivityService) DeleteOldActivities(daysOld int) (int64, error) {
	if daysOld < 30 {
		daysOld = 30 // Minimum retention period
	}

	result, err := facades.Orm().Query().
		Model(&models.UserActivity{}).
		Where("created_at < NOW() - INTERVAL '? days'", daysOld).
		Delete(&models.UserActivity{})

	if err != nil {
		return 0, err
	}

	return result.RowsAffected, nil
}

// ActivitySummary represents a summary of user activities
type ActivitySummary struct {
	ActivityType string `json:"activity_type"`
	Count        int64  `json:"count"`
	LastOccurred string `json:"last_occurred"`
}

// GetActivitySummary returns a summary of activities by type for a user
func (s *UserActivityService) GetActivitySummary(userID uint) ([]ActivitySummary, error) {
	var summaries []ActivitySummary

	// Use Goravel's query builder to get activity type counts
	activityTypes, err := s.GetActivityTypes(userID)
	if err != nil {
		return nil, err
	}

	for _, activityType := range activityTypes {
		var count int64
		count, err = facades.Orm().Query().
			Model(&models.UserActivity{}).
			Where("user_id = ?", userID).
			Where("activity_type = ?", activityType).
			Count()
		if err != nil {
			continue
		}

		// Get the most recent activity of this type
		var lastActivity models.UserActivity
		err = facades.Orm().Query().
			Model(&models.UserActivity{}).
			Where("user_id = ?", userID).
			Where("activity_type = ?", activityType).
			Order("created_at DESC").
			First(&lastActivity)

		lastOccurred := ""
		if err == nil && lastActivity.ID != 0 {
			lastOccurred = lastActivity.CreatedAt.Format("2006-01-02 15:04:05")
		}

		summaries = append(summaries, ActivitySummary{
			ActivityType: activityType,
			Count:        count,
			LastOccurred: lastOccurred,
		})
	}

	return summaries, nil
}

// FormatActivityForDisplay formats an activity for display purposes
func (s *UserActivityService) FormatActivityForDisplay(activity *models.UserActivity) map[string]interface{} {
	result := map[string]interface{}{
		"id":           activity.ID,
		"activityType": activity.ActivityType,
		"description":  activity.Description,
		"ipAddress":    activity.IPAddress,
		"createdAt":    activity.CreatedAt,
	}

	// Parse and include metadata
	if activity.Metadata != "" {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(activity.Metadata), &metadata); err == nil {
			result["metadata"] = metadata
		}
	}

	// Include related entity info
	if activity.RelatedType != "" {
		result["relatedType"] = activity.RelatedType
		if activity.RelatedID != nil {
			result["relatedId"] = *activity.RelatedID
		}
	}

	return result
}
