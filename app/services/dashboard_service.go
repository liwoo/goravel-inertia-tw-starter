package services

import (
	"starter-project/app/models"

	"github.com/goravel/framework/facades"
)

// DashboardService provides methods for dashboard statistics and widgets
type DashboardService struct{}

// NewDashboardService creates a new dashboard service
func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

// RecentActivityDTO represents the data structure for recent user activities
type RecentActivityDTO struct {
	ID         uint   `json:"id"`
	EntityType string `json:"entityType"`
	EntityID   uint   `json:"entityId"`
	EntityName string `json:"entityName"`
	Action     string `json:"action"` // "created", "updated"
	UserID     uint   `json:"userId"`
	UserName   string `json:"userName"`
	Timestamp  string `json:"timestamp"`
}

// GetRecentActivities returns the most recent user activities across all entities
func (s *DashboardService) GetRecentActivities(limit int) []RecentActivityDTO {
	var activities []RecentActivityDTO

	// Cache for user names to avoid repeated queries
	userNameCache := make(map[uint]string)

	// Helper to get user name by ID with caching
	getUserName := func(userID *uint) string {
		if userID == nil || *userID == 0 {
			return "System"
		}

		// Check cache first
		if name, ok := userNameCache[*userID]; ok {
			return name
		}

		var user models.User
		err := facades.Orm().Query().Model(&models.User{}).Where("id = ?", *userID).First(&user)
		if err != nil || user.ID == 0 {
			userNameCache[*userID] = "Unknown"
			return "Unknown"
		}

		userNameCache[*userID] = user.Name
		return user.Name
	}

	var books []models.Book
	facades.Orm().Query().Model(&models.Book{}).
		Order("updated_at DESC").
		Limit(limit).
		Find(&books)

	for _, book := range books {
		action := "updated"
		userID := book.UpdatedBy
		timestamp := book.UpdatedAt.ToDateTimeString()

		// Check if this is a newly created record (created_at == updated_at approximately)
		if book.CreatedAt.ToDateTimeString() == book.UpdatedAt.ToDateTimeString() {
			action = "created"
			userID = book.CreatedBy
			timestamp = book.CreatedAt.ToDateTimeString()
		}

		activities = append(activities, RecentActivityDTO{
			ID:         book.ID,
			EntityType: "book",
			EntityID:   book.ID,
			EntityName: book.Title,
			Action:     action,
			UserID:     *userID,
			UserName:   getUserName(userID),
			Timestamp:  timestamp,
		})
	}

	// Sort all activities by timestamp (newest first) and limit
	sortActivitiesByTimestamp(activities)
	if len(activities) > limit {
		activities = activities[:limit]
	}

	return activities
}

// safeInt converts *int to int safely
func safeInt(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}

// GetUserRecentActivities returns the most recent activities performed BY a specific user
func (s *DashboardService) GetUserRecentActivities(userID uint, limit int) []RecentActivityDTO {
	var activities []RecentActivityDTO

	var books []models.Book
	facades.Orm().Query().Model(&models.Book{}).
		Where("created_by = ? OR updated_by = ?", userID, userID).
		Order("updated_at DESC").
		Limit(limit).
		Find(&books)

	for _, book := range books {
		action := "updated"
		timestamp := book.UpdatedAt.ToDateTimeString()

		// Check if this user created or updated
		if book.CreatedBy != nil && *book.CreatedBy == userID &&
			book.CreatedAt.ToDateTimeString() == book.UpdatedAt.ToDateTimeString() {
			action = "created"
			timestamp = book.CreatedAt.ToDateTimeString()
		}

		activities = append(activities, RecentActivityDTO{
			ID:         book.ID,
			EntityType: "book",
			EntityID:   book.ID,
			EntityName: book.Title,
			Action:     action,
			UserID:     userID,
			UserName:   "", // Not needed for user's own activities
			Timestamp:  timestamp,
		})
	}

	// Sort all activities by timestamp (newest first) and limit
	sortActivitiesByTimestamp(activities)
	if len(activities) > limit {
		activities = activities[:limit]
	}

	return activities
}

// sortActivitiesByTimestamp sorts activities by timestamp descending
func sortActivitiesByTimestamp(activities []RecentActivityDTO) {
	for i := 0; i < len(activities)-1; i++ {
		for j := i + 1; j < len(activities); j++ {
			if activities[i].Timestamp < activities[j].Timestamp {
				activities[i], activities[j] = activities[j], activities[i]
			}
		}
	}
}

// GetDashboardStats returns all statistics for the dashboard
func (s *DashboardService) GetDashboardStats() map[string]interface{} {

	// Combine all statistics
	return map[string]interface{}{}
}
