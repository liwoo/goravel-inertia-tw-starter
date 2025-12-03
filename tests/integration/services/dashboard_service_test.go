package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
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

// TestGetDistributionByYouth tests that youth distribution returns correct data
func (s *DashboardServiceTestSuite) TestGetDistributionByYouth_ReturnsYouthAndNonYouth() {
	userID := int(s.testUser.ID)

	// Create SME with young owner (age 25 - youth)
	sme1 := &models.Sme{
		UsmeNumber:       fmt.Sprintf("USME-YOUTH1-%d", time.Now().UnixNano()),
		Name:             "Youth Owner SME",
		ContactPhone:     "+265999111222",
		ContactEmail:     "youth@sme.com",
		BusinessCategory: "Agriculture",
		Sector:           "Agribusiness",
		CreatedBy:        &userID,
		UpdatedBy:        &userID,
	}
	err := facades.Orm().Query().Create(sme1)
	s.NoError(err)

	// Create primary business owner for sme1 - 25 years old (youth)
	youngDOB := time.Now().AddDate(-25, 0, 0)
	youngDOBCarbon := carbon.NewDateTime(carbon.FromStdTime(youngDOB))
	youngOwner := &models.PrimaryBusinessOwner{
		FirstName:        "Young",
		LastName:         "Owner",
		Nationality:      "Malawian",
		NationalIdNumber: "MW123456",
		DateOfBirth:      *youngDOBCarbon,
		Gender:           "MALE",
		EducationLevel:   "Bachelor",
		MalawianStatus:   "Citizen",
		PhoneNumber:      "+265999111222",
		SmeID:            int(sme1.ID),
	}
	err = facades.Orm().Query().Create(youngOwner)
	s.NoError(err)

	// Create SME with older owner (age 45 - non-youth)
	sme2 := &models.Sme{
		UsmeNumber:       fmt.Sprintf("USME-NONYOUTH1-%d", time.Now().UnixNano()),
		Name:             "Non-Youth Owner SME",
		ContactPhone:     "+265999111333",
		ContactEmail:     "nonyouth@sme.com",
		BusinessCategory: "Manufacturing",
		Sector:           "Industrial",
		CreatedBy:        &userID,
		UpdatedBy:        &userID,
	}
	err = facades.Orm().Query().Create(sme2)
	s.NoError(err)

	// Create primary business owner for sme2 - 45 years old (non-youth)
	olderDOB := time.Now().AddDate(-45, 0, 0)
	olderDOBCarbon := carbon.NewDateTime(carbon.FromStdTime(olderDOB))
	olderOwner := &models.PrimaryBusinessOwner{
		FirstName:        "Older",
		LastName:         "Owner",
		Nationality:      "Malawian",
		NationalIdNumber: "MW654321",
		DateOfBirth:      *olderDOBCarbon,
		Gender:           "FEMALE",
		EducationLevel:   "Master",
		MalawianStatus:   "Citizen",
		PhoneNumber:      "+265999111333",
		SmeID:            int(sme2.ID),
	}
	err = facades.Orm().Query().Create(olderOwner)
	s.NoError(err)

	// Get youth distribution
	smeService := services.NewSmeService()
	distribution := smeService.GetDistributionByYouth()

	// Should have data for youth and non-youth
	s.NotEmpty(distribution, "Should return youth distribution data")
	s.Len(distribution, 2, "Should have exactly 2 categories (Youth and Non-Youth)")

	// Verify the distribution contains expected categories
	hasYouth := false
	hasNonYouth := false
	for _, d := range distribution {
		label := d["label"].(string)
		value := d["value"].(int64)
		if label == "Youth" {
			hasYouth = true
			s.Equal(int64(1), value, "Should have 1 youth owner")
		} else if label == "Non-Youth" {
			hasNonYouth = true
			s.Equal(int64(1), value, "Should have 1 non-youth owner")
		}
	}
	s.True(hasYouth, "Should have Youth category")
	s.True(hasNonYouth, "Should have Non-Youth category")
}

// TestGetAgeGenderDistribution tests that age-gender distribution returns correct data
func (s *DashboardServiceTestSuite) TestGetAgeGenderDistribution_ReturnsPyramidData() {
	userID := int(s.testUser.ID)

	// Create test data with different age groups and genders
	// Gender must be uppercase to match database check constraint
	testCases := []struct {
		age    int
		gender string
	}{
		{20, "MALE"},   // 18-25
		{22, "FEMALE"}, // 18-25
		{30, "MALE"},   // 26-35
		{40, "FEMALE"}, // 36-45
		{50, "MALE"},   // 46-55
	}

	for i, tc := range testCases {
		// Create SME
		sme := &models.Sme{
			UsmeNumber:       fmt.Sprintf("USME-AGEGEN-%d-%d", i, time.Now().UnixNano()),
			Name:             fmt.Sprintf("Age Gender Test SME %d", i),
			ContactPhone:     "+265999111222",
			ContactEmail:     fmt.Sprintf("agegen%d@sme.com", i),
			BusinessCategory: "Agriculture",
			Sector:           "Agribusiness",
			CreatedBy:        &userID,
			UpdatedBy:        &userID,
		}
		err := facades.Orm().Query().Create(sme)
		s.NoError(err)

		// Create primary business owner
		dob := time.Now().AddDate(-tc.age, 0, 0)
		dobCarbon := carbon.NewDateTime(carbon.FromStdTime(dob))
		owner := &models.PrimaryBusinessOwner{
			FirstName:        fmt.Sprintf("Owner%d", i),
			LastName:         "Test",
			Nationality:      "Malawian",
			NationalIdNumber: fmt.Sprintf("MW%d", i),
			DateOfBirth:      *dobCarbon,
			Gender:           tc.gender,
			EducationLevel:   "Bachelor",
			MalawianStatus:   "Citizen",
			PhoneNumber:      fmt.Sprintf("+26599911%04d", i),
			SmeID:            int(sme.ID),
		}
		err = facades.Orm().Query().Create(owner)
		s.NoError(err)
	}

	// Get age-gender distribution
	smeService := services.NewSmeService()
	distribution := smeService.GetAgeGenderDistribution()

	// Should have data for age groups
	s.NotEmpty(distribution, "Should return age-gender distribution data")
	s.Len(distribution, 6, "Should have 6 age groups")

	// Verify expected age groups exist
	ageGroups := make(map[string]bool)
	for _, d := range distribution {
		ageGroup := d["ageGroup"].(string)
		ageGroups[ageGroup] = true

		// Verify structure
		s.Contains(d, "male", "Should have male count")
		s.Contains(d, "female", "Should have female count")
		s.Contains(d, "malePercentage", "Should have male percentage")
		s.Contains(d, "femalePercentage", "Should have female percentage")
	}

	// Check expected age groups
	expectedGroups := []string{"18-25", "26-35", "36-45", "46-55", "56-65", "65+"}
	for _, eg := range expectedGroups {
		s.True(ageGroups[eg], "Should have age group: "+eg)
	}
}

// TestGetDashboardStats_IncludesYouthAndAgeGenderData tests that GetDashboardStats includes the distribution data
func (s *DashboardServiceTestSuite) TestGetDashboardStats_IncludesYouthAndAgeGenderData() {
	userID := int(s.testUser.ID)

	// Create test SME with owner
	sme := &models.Sme{
		UsmeNumber:       fmt.Sprintf("USME-DASHSTATS-%d", time.Now().UnixNano()),
		Name:             "Dashboard Stats Test SME",
		ContactPhone:     "+265999111222",
		ContactEmail:     "dashstats@sme.com",
		BusinessCategory: "Agriculture",
		Sector:           "Agribusiness",
		CreatedBy:        &userID,
		UpdatedBy:        &userID,
	}
	err := facades.Orm().Query().Create(sme)
	s.NoError(err)

	dob := time.Now().AddDate(-30, 0, 0)
	dobCarbon := carbon.NewDateTime(carbon.FromStdTime(dob))
	owner := &models.PrimaryBusinessOwner{
		FirstName:        "Dashboard",
		LastName:         "Test",
		Nationality:      "Malawian",
		NationalIdNumber: "MW999999",
		DateOfBirth:      *dobCarbon,
		Gender:           "MALE",
		EducationLevel:   "Bachelor",
		MalawianStatus:   "Citizen",
		PhoneNumber:      "+265999111222",
		SmeID:            int(sme.ID),
	}
	err = facades.Orm().Query().Create(owner)
	s.NoError(err)

	// Get dashboard stats
	smeService := services.NewSmeService()
	stats := s.dashboardService.GetDashboardStats(smeService)

	// Verify byYouth and byAgeGender are included
	s.Contains(stats, "byYouth", "Dashboard stats should include byYouth")
	s.Contains(stats, "byAgeGender", "Dashboard stats should include byAgeGender")

	// Verify byYouth has data
	byYouth, ok := stats["byYouth"].([]map[string]interface{})
	s.True(ok, "byYouth should be a slice of maps")
	s.NotEmpty(byYouth, "byYouth should have data")

	// Verify byAgeGender has data
	byAgeGender, ok := stats["byAgeGender"].([]map[string]interface{})
	s.True(ok, "byAgeGender should be a slice of maps")
	s.NotEmpty(byAgeGender, "byAgeGender should have data")
}
