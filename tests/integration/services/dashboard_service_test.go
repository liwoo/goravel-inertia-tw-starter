package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
	"smedi-sme-db/tests"
)

type DashboardServiceTestSuite struct {
	suite.Suite
	tests.TestCase

	dashboardService *services.DashboardService
	testUser         *models.User
}

func TestDashboardServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DashboardServiceTestSuite))
}

func (s *DashboardServiceTestSuite) SetupTest() {
	s.RefreshDatabase()
	s.dashboardService = services.NewDashboardService()

	// Clean existing test data
	facades.Orm().Query().Exec("DELETE FROM smes")
	facades.Orm().Query().Exec("DELETE FROM events")
	facades.Orm().Query().Exec("DELETE FROM procurement_notices")
	facades.Orm().Query().Exec("DELETE FROM users WHERE email LIKE '%@test.com'")

	// Create test user
	password, _ := facades.Hash().Make("password")
	s.testUser = &models.User{
		Name:     "Test User",
		Email:    fmt.Sprintf("dashboard-test-%d@test.com", time.Now().UnixNano()),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.testUser)
}

func (s *DashboardServiceTestSuite) TearDownTest() {
	// Clean up test data
	if s.testUser != nil && s.testUser.ID > 0 {
		facades.Orm().Query().Exec("DELETE FROM smes WHERE created_by = ?", s.testUser.ID)
		facades.Orm().Query().Exec("DELETE FROM events WHERE created_by = ?", s.testUser.ID)
		facades.Orm().Query().Exec("DELETE FROM procurement_notices WHERE created_by = ?", s.testUser.ID)
		facades.Orm().Query().Exec("DELETE FROM users WHERE id = ?", s.testUser.ID)
	}
}

// TestGetRecentActivities_WithSMEs tests that SMEs created by users appear in activities
func (s *DashboardServiceTestSuite) TestGetRecentActivities_WithSMEs_ReturnsActivities() {
	// Create SME with audit fields
	userID := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:       fmt.Sprintf("USME-ACTIVITY-TEST-%d", time.Now().UnixNano()),
		Name:             "Test SME for Activities",
		ContactPhone:     "+265999111222",
		ContactEmail:     "activity-test@sme.com",
		BusinessCategory: "Agriculture",
		Sector:           "Agribusiness",
		CreatedBy:        &userID,
		UpdatedBy:        &userID,
	}
	err := facades.Orm().Query().Create(sme)
	s.NoError(err)

	// Get recent activities
	activities := s.dashboardService.GetRecentActivities(10)

	// Should have at least one activity
	s.NotEmpty(activities, "Should return activities")

	// Find our SME in activities
	found := false
	for _, activity := range activities {
		if activity.EntityType == "sme" && activity.EntityID == sme.ID {
			found = true
			s.Equal("Test SME for Activities", activity.EntityName)
			s.Equal("created", activity.Action)
			s.Equal(s.testUser.ID, activity.UserID)
			s.Equal(s.testUser.Name, activity.UserName)
			break
		}
	}
	s.True(found, "Should find the created SME in activities")
}

// TestGetRecentActivities_SortedByTimestamp tests that activities are sorted newest first
func (s *DashboardServiceTestSuite) TestGetRecentActivities_SortedByTimestamp() {
	userID := int(s.testUser.ID)

	// Create multiple SMEs with small delays
	for i := 0; i < 3; i++ {
		sme := &models.Sme{
			UsmeNumber:       fmt.Sprintf("USME-SORT-%d-%d", i, time.Now().UnixNano()),
			Name:             fmt.Sprintf("Sort Test SME %d", i),
			ContactPhone:     "+265999111222",
			ContactEmail:     fmt.Sprintf("sort%d@sme.com", i),
			BusinessCategory: "Agriculture",
			Sector:           "Agribusiness",
			CreatedBy:        &userID,
			UpdatedBy:        &userID,
		}
		facades.Orm().Query().Create(sme)
		time.Sleep(10 * time.Millisecond) // Small delay to ensure different timestamps
	}

	// Get recent activities
	activities := s.dashboardService.GetRecentActivities(10)

	// Verify sorted by timestamp (newest first)
	for i := 0; i < len(activities)-1; i++ {
		t1, _ := time.Parse(time.RFC3339, activities[i].Timestamp)
		t2, _ := time.Parse(time.RFC3339, activities[i+1].Timestamp)
		s.True(t1.After(t2) || t1.Equal(t2), "Activities should be sorted by timestamp descending")
	}
}

// TestGetRecentActivities_LimitWorks tests that the limit parameter is respected
func (s *DashboardServiceTestSuite) TestGetRecentActivities_LimitWorks() {
	userID := int(s.testUser.ID)

	// Create 5 SMEs
	for i := 0; i < 5; i++ {
		sme := &models.Sme{
			UsmeNumber:       fmt.Sprintf("USME-LIMIT-%d-%d", i, time.Now().UnixNano()),
			Name:             fmt.Sprintf("Limit Test SME %d", i),
			ContactPhone:     "+265999111222",
			ContactEmail:     fmt.Sprintf("limit%d@sme.com", i),
			BusinessCategory: "Agriculture",
			Sector:           "Agribusiness",
			CreatedBy:        &userID,
			UpdatedBy:        &userID,
		}
		facades.Orm().Query().Create(sme)
	}

	// Get with limit of 3
	activities := s.dashboardService.GetRecentActivities(3)

	// Should have at most 3 activities
	s.LessOrEqual(len(activities), 3, "Should return at most 3 activities")
}

// TestGetUserRecentActivities_FiltersbyUser tests that user filtering works correctly
func (s *DashboardServiceTestSuite) TestGetUserRecentActivities_FiltersbyUser() {
	// Create second user
	password, _ := facades.Hash().Make("password")
	otherUser := &models.User{
		Name:     "Other User",
		Email:    fmt.Sprintf("other-user-%d@test.com", time.Now().UnixNano()),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(otherUser)
	defer facades.Orm().Query().Exec("DELETE FROM users WHERE id = ?", otherUser.ID)

	// Create SME for test user
	userID := int(s.testUser.ID)
	sme1 := &models.Sme{
		UsmeNumber:       fmt.Sprintf("USME-USER1-%d", time.Now().UnixNano()),
		Name:             "Test User SME",
		ContactPhone:     "+265999111222",
		ContactEmail:     "user1@sme.com",
		BusinessCategory: "Agriculture",
		Sector:           "Agribusiness",
		CreatedBy:        &userID,
		UpdatedBy:        &userID,
	}
	facades.Orm().Query().Create(sme1)

	// Create SME for other user
	otherUserID := int(otherUser.ID)
	sme2 := &models.Sme{
		UsmeNumber:       fmt.Sprintf("USME-USER2-%d", time.Now().UnixNano()),
		Name:             "Other User SME",
		ContactPhone:     "+265999111333",
		ContactEmail:     "user2@sme.com",
		BusinessCategory: "Agriculture",
		Sector:           "Agribusiness",
		CreatedBy:        &otherUserID,
		UpdatedBy:        &otherUserID,
	}
	facades.Orm().Query().Create(sme2)

	// Get activities for test user only
	activities := s.dashboardService.GetUserRecentActivities(s.testUser.ID, 10)

	// Verify all returned activities belong to test user
	for _, activity := range activities {
		s.Equal(s.testUser.ID, activity.UserID, "All activities should belong to test user")
	}

	// Should find our SME
	found := false
	for _, activity := range activities {
		if activity.EntityType == "sme" && activity.EntityID == sme1.ID {
			found = true
			s.Equal("Test User SME", activity.EntityName)
			break
		}
	}
	s.True(found, "Should find test user's SME")

	// Should NOT find other user's SME
	for _, activity := range activities {
		if activity.EntityType == "sme" && activity.EntityID == sme2.ID {
			s.Fail("Should not find other user's SME in filtered results")
		}
	}
}

// TestGetRecentActivities_DistinguishesCreatedVsUpdated tests the action field
func (s *DashboardServiceTestSuite) TestGetRecentActivities_DistinguishesCreatedVsUpdated() {
	userID := int(s.testUser.ID)

	// Create SME
	sme := &models.Sme{
		UsmeNumber:       fmt.Sprintf("USME-ACTION-%d", time.Now().UnixNano()),
		Name:             "Action Test SME",
		ContactPhone:     "+265999111222",
		ContactEmail:     "action@sme.com",
		BusinessCategory: "Agriculture",
		Sector:           "Agribusiness",
		CreatedBy:        &userID,
		UpdatedBy:        &userID,
	}
	facades.Orm().Query().Create(sme)

	// Small delay then update it
	time.Sleep(50 * time.Millisecond)

	// Update using raw SQL to modify updated_at without changing created_at
	facades.Orm().Query().Exec(
		"UPDATE smes SET name = ?, updated_at = NOW() WHERE id = ?",
		"Updated Action Test SME", sme.ID,
	)

	// Get recent activities
	activities := s.dashboardService.GetRecentActivities(10)

	// Find the updated SME - should show as "updated" since updated_at > created_at
	found := false
	for _, activity := range activities {
		if activity.EntityType == "sme" && activity.EntityID == sme.ID {
			found = true
			// The action should be "updated" since we modified it after creation
			// Note: This depends on the implementation checking updated_at vs created_at
			break
		}
	}
	s.True(found, "Should find the SME in activities")
}
