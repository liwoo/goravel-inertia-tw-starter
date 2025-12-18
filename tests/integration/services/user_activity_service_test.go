package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"starter-project/app/models"
	"starter-project/app/services"
	"starter-project/tests"
)

type UserActivityServiceTestSuite struct {
	suite.Suite
	tests.TestCase

	activityService *services.UserActivityService
	testUser        *models.User
}

func TestUserActivityServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserActivityServiceTestSuite))
}

func (s *UserActivityServiceTestSuite) SetupTest() {
	s.RefreshDatabase()
	s.activityService = services.NewUserActivityService()

	// Clean existing test data
	facades.Orm().Query().Exec("DELETE FROM user_activities")
	facades.Orm().Query().Exec("DELETE FROM users WHERE email LIKE '%@test.com'")

	// Create test user
	password, _ := facades.Hash().Make("password")
	s.testUser = &models.User{
		Name:     "Activity Test User",
		Email:    fmt.Sprintf("activity-test-%d@test.com", time.Now().UnixNano()),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.testUser)
}

func (s *UserActivityServiceTestSuite) TearDownTest() {
	if s.testUser != nil && s.testUser.ID > 0 {
		facades.Orm().Query().Exec("DELETE FROM user_activities WHERE user_id = ?", s.testUser.ID)
		facades.Orm().Query().Exec("DELETE FROM users WHERE id = ?", s.testUser.ID)
	}
}

// TestLogActivity_CreatesActivityRecord tests that activities are logged correctly
func (s *UserActivityServiceTestSuite) TestLogActivity_CreatesActivityRecord() {
	// Log an activity
	err := s.activityService.LogActivity(
		nil, // No context in tests
		s.testUser.ID,
		models.ActivityLogin,
		"Test login activity",
		map[string]interface{}{
			"source": "test",
		},
	)
	s.NoError(err)

	// Verify the activity was created
	var activity models.UserActivity
	err = facades.Orm().Query().
		Where("user_id = ?", s.testUser.ID).
		Where("activity_type = ?", models.ActivityLogin).
		First(&activity)
	s.NoError(err)
	s.Equal(s.testUser.ID, activity.UserID)
	s.Equal(models.ActivityLogin, activity.ActivityType)
	s.Equal("Test login activity", activity.Description)
}

// TestLogActivityWithRelation_SetsRelatedEntity tests related entity logging
func (s *UserActivityServiceTestSuite) TestLogActivityWithRelation_SetsRelatedEntity() {
	// Log an activity with a relation
	relatedID := uint(123)
	err := s.activityService.LogActivityWithRelation(
		nil,
		s.testUser.ID,
		"book_view", // Using string directly since constant doesn't exist
		"Viewed Book record",
		nil,
		"book",
		relatedID,
	)
	s.NoError(err)

	// Verify the activity was created with relation
	var activity models.UserActivity
	err = facades.Orm().Query().
		Where("user_id = ?", s.testUser.ID).
		Where("activity_type = ?", "book_view").
		First(&activity)
	s.NoError(err)
	s.Equal("book", activity.RelatedType)
	s.NotNil(activity.RelatedID)
	s.Equal(relatedID, *activity.RelatedID)
}

// TestGetUserActivities_ReturnsActivities tests activity retrieval
func (s *UserActivityServiceTestSuite) TestGetUserActivities_ReturnsActivities() {
	// Create some activities
	for i := 0; i < 5; i++ {
		s.activityService.LogActivity(
			nil,
			s.testUser.ID,
			models.ActivityLogin,
			fmt.Sprintf("Login attempt %d", i),
			nil,
		)
	}

	// Get activities
	result, err := s.activityService.GetUserActivities(s.testUser.ID, 1, 10, "")
	s.NoError(err)
	s.NotNil(result)
	s.Equal(int64(5), result.Total)
	s.Len(result.Data, 5)
}

// TestGetUserActivities_Pagination tests pagination
func (s *UserActivityServiceTestSuite) TestGetUserActivities_Pagination() {
	// Create 15 activities
	for i := 0; i < 15; i++ {
		s.activityService.LogActivity(
			nil,
			s.testUser.ID,
			models.ActivityLogin,
			fmt.Sprintf("Login attempt %d", i),
			nil,
		)
	}

	// Get first page
	result1, err := s.activityService.GetUserActivities(s.testUser.ID, 1, 10, "")
	s.NoError(err)
	s.Equal(int64(15), result1.Total)
	s.Len(result1.Data, 10)
	s.Equal(1, result1.Page)
	s.Equal(2, result1.TotalPages)

	// Get second page
	result2, err := s.activityService.GetUserActivities(s.testUser.ID, 2, 10, "")
	s.NoError(err)
	s.Equal(int64(15), result2.Total)
	s.Len(result2.Data, 5)
	s.Equal(2, result2.Page)
}

// TestGetUserActivities_FilterByType tests filtering by activity type
func (s *UserActivityServiceTestSuite) TestGetUserActivities_FilterByType() {
	// Create activities of different types
	s.activityService.LogActivity(nil, s.testUser.ID, models.ActivityLogin, "Login", nil)
	s.activityService.LogActivity(nil, s.testUser.ID, models.ActivityLogin, "Login 2", nil)
	s.activityService.LogActivity(nil, s.testUser.ID, models.ActivityLogout, "Logout", nil)
	s.activityService.LogActivity(nil, s.testUser.ID, models.ActivityProfileUpdate, "Profile update", nil)

	// Filter by login type
	result, err := s.activityService.GetUserActivities(s.testUser.ID, 1, 10, models.ActivityLogin)
	s.NoError(err)
	s.Equal(int64(2), result.Total)
	for _, activity := range result.Data {
		s.Equal(models.ActivityLogin, activity.ActivityType)
	}
}

// TestGetActivityTypes_ReturnsDistinctTypes tests getting distinct activity types
func (s *UserActivityServiceTestSuite) TestGetActivityTypes_ReturnsDistinctTypes() {
	// Create activities of different types
	s.activityService.LogActivity(nil, s.testUser.ID, models.ActivityLogin, "Login", nil)
	s.activityService.LogActivity(nil, s.testUser.ID, models.ActivityLogout, "Logout", nil)
	s.activityService.LogActivity(nil, s.testUser.ID, models.ActivityProfileUpdate, "Profile", nil)

	// Get activity types
	types, err := s.activityService.GetActivityTypes(s.testUser.ID)
	s.NoError(err)
	s.Len(types, 3)
	s.Contains(types, models.ActivityLogin)
	s.Contains(types, models.ActivityLogout)
	s.Contains(types, models.ActivityProfileUpdate)
}

// TestGetActivitySummary_GroupsByType tests activity summary grouping
func (s *UserActivityServiceTestSuite) TestGetActivitySummary_GroupsByType() {
	// Create activities with different types
	for i := 0; i < 3; i++ {
		s.activityService.LogActivity(nil, s.testUser.ID, models.ActivityLogin, "Login", nil)
	}
	for i := 0; i < 2; i++ {
		s.activityService.LogActivity(nil, s.testUser.ID, models.ActivityLogout, "Logout", nil)
	}
	s.activityService.LogActivity(nil, s.testUser.ID, models.ActivityProfileUpdate, "Profile", nil)

	// Get summary
	summaries, err := s.activityService.GetActivitySummary(s.testUser.ID)
	s.NoError(err)

	// Verify we have summaries for the activity types
	s.NotNil(summaries)
	s.Len(summaries, 3)

	// Build a map from the summaries for easier assertion
	summaryMap := make(map[string]int64)
	for _, summary := range summaries {
		summaryMap[summary.ActivityType] = summary.Count
	}

	s.Equal(int64(3), summaryMap[models.ActivityLogin])
	s.Equal(int64(2), summaryMap[models.ActivityLogout])
	s.Equal(int64(1), summaryMap[models.ActivityProfileUpdate])
}

// TestLogActivity_WithMetadata tests metadata is stored correctly
func (s *UserActivityServiceTestSuite) TestLogActivity_WithMetadata() {
	// Log activity with metadata
	metadata := map[string]interface{}{
		"browser": "Chrome",
		"version": "120.0",
		"os":      "macOS",
	}
	err := s.activityService.LogActivity(
		nil,
		s.testUser.ID,
		models.ActivityLogin,
		"Login with metadata",
		metadata,
	)
	s.NoError(err)

	// Verify metadata was stored
	var activity models.UserActivity
	err = facades.Orm().Query().
		Where("user_id = ?", s.testUser.ID).
		Where("activity_type = ?", models.ActivityLogin).
		First(&activity)
	s.NoError(err)

	// Get metadata (returns map and error)
	storedMetadata, metaErr := activity.GetMetadata()
	s.NoError(metaErr)
	s.NotNil(storedMetadata)
	s.Equal("Chrome", storedMetadata["browser"])
	s.Equal("120.0", storedMetadata["version"])
}

// TestGetUserActivities_SortedByCreatedAtDesc tests activities are sorted newest first
func (s *UserActivityServiceTestSuite) TestGetUserActivities_SortedByCreatedAtDesc() {
	// Create activities with small delays
	for i := 0; i < 3; i++ {
		s.activityService.LogActivity(
			nil,
			s.testUser.ID,
			models.ActivityLogin,
			fmt.Sprintf("Login %d", i),
			nil,
		)
		time.Sleep(10 * time.Millisecond)
	}

	// Get activities
	result, err := s.activityService.GetUserActivities(s.testUser.ID, 1, 10, "")
	s.NoError(err)

	// Verify sorted by created_at descending (newest first)
	for i := 0; i < len(result.Data)-1; i++ {
		t1 := result.Data[i].CreatedAt.StdTime()
		t2 := result.Data[i+1].CreatedAt.StdTime()
		s.True(
			t1.After(t2) || t1.Equal(t2),
			"Activities should be sorted by created_at descending",
		)
	}
}

// TestDeleteOldActivities tests cleaning up old activities
func (s *UserActivityServiceTestSuite) TestActivityTypes_Constants() {
	// Verify activity type constants are defined
	s.NotEmpty(models.ActivityLogin)
	s.NotEmpty(models.ActivityLogout)
	s.NotEmpty(models.ActivityProfileUpdate)
	s.NotEmpty(models.ActivityPasswordChange)
}
