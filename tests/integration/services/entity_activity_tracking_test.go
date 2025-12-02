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

// EntityActivityTrackingTestSuite tests that CRUD operations on core entities
// are properly tracked in the activity history via audit fields
type EntityActivityTrackingTestSuite struct {
	suite.Suite
	tests.TestCase

	dashboardService *services.DashboardService
	testUser         *models.User
	testUser2        *models.User
}

func TestEntityActivityTrackingTestSuite(t *testing.T) {
	suite.Run(t, new(EntityActivityTrackingTestSuite))
}

func (s *EntityActivityTrackingTestSuite) SetupTest() {
	s.RefreshDatabase()
	s.dashboardService = services.NewDashboardService()

	// Clean existing test data
	facades.Orm().Query().Exec("DELETE FROM smes")
	facades.Orm().Query().Exec("DELETE FROM events")
	facades.Orm().Query().Exec("DELETE FROM procurement_notices")
	facades.Orm().Query().Exec("DELETE FROM users WHERE email LIKE '%@test.com'")

	// Create test users
	password, _ := facades.Hash().Make("password")
	s.testUser = &models.User{
		Name:     "Test User One",
		Email:    fmt.Sprintf("entity-test-1-%d@test.com", time.Now().UnixNano()),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.testUser)

	s.testUser2 = &models.User{
		Name:     "Test User Two",
		Email:    fmt.Sprintf("entity-test-2-%d@test.com", time.Now().UnixNano()),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.testUser2)
}

func (s *EntityActivityTrackingTestSuite) TearDownTest() {
	// Clean up test data
	facades.Orm().Query().Exec("DELETE FROM smes")
	facades.Orm().Query().Exec("DELETE FROM events")
	facades.Orm().Query().Exec("DELETE FROM procurement_notices")
	if s.testUser != nil && s.testUser.ID > 0 {
		facades.Orm().Query().Exec("DELETE FROM users WHERE id = ?", s.testUser.ID)
	}
	if s.testUser2 != nil && s.testUser2.ID > 0 {
		facades.Orm().Query().Exec("DELETE FROM users WHERE id = ?", s.testUser2.ID)
	}
}

// ===========================================
// SME CRUD Activity Tracking Tests
// ===========================================

func (s *EntityActivityTrackingTestSuite) TestSME_Create_AppearsInActivityHistory() {
	// Create an SME with audit fields set
	userID := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:       fmt.Sprintf("USME-CREATE-%d", time.Now().UnixNano()),
		Name:             "New SME Company",
		ContactPhone:     "+265999111222",
		ContactEmail:     "newsme@test.com",
		BusinessCategory: "Agriculture",
		Sector:           "Agribusiness",
		CreatedBy:        &userID,
		UpdatedBy:        &userID,
	}
	err := facades.Orm().Query().Create(sme)
	s.NoError(err)
	s.NotZero(sme.ID)

	// Get recent activities
	activities := s.dashboardService.GetRecentActivities(20)

	// Find the SME in activity history
	found := false
	for _, activity := range activities {
		if activity.EntityType == "sme" && activity.EntityID == sme.ID {
			found = true
			s.Equal("New SME Company", activity.EntityName)
			s.Equal("created", activity.Action, "New SME should show as 'created'")
			s.Equal(s.testUser.ID, activity.UserID)
			s.Equal(s.testUser.Name, activity.UserName)
			break
		}
	}
	s.True(found, "Created SME should appear in activity history")
}

func (s *EntityActivityTrackingTestSuite) TestSME_Update_AppearsInActivityHistory() {
	// Create an SME first
	userID := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:       fmt.Sprintf("USME-UPDATE-%d", time.Now().UnixNano()),
		Name:             "Original SME Name",
		ContactPhone:     "+265999111222",
		ContactEmail:     "original@test.com",
		BusinessCategory: "Agriculture",
		Sector:           "Agribusiness",
		CreatedBy:        &userID,
		UpdatedBy:        &userID,
	}
	facades.Orm().Query().Create(sme)

	// Wait long enough so updated_at will differ from created_at in seconds
	// (ToDateTimeString() only has second precision)
	time.Sleep(1100 * time.Millisecond)

	// Update the SME with a different user using raw SQL for reliability
	user2ID := int(s.testUser2.ID)
	_, err := facades.Orm().Query().Exec(
		"UPDATE smes SET name = ?, updated_by = ?, updated_at = NOW() WHERE id = ?",
		"Updated SME Name", user2ID, sme.ID,
	)
	s.NoError(err)

	// Get recent activities
	activities := s.dashboardService.GetRecentActivities(20)

	// Find the SME in activity history - should show as updated by user2
	found := false
	for _, activity := range activities {
		if activity.EntityType == "sme" && activity.EntityID == sme.ID {
			found = true
			s.Equal("Updated SME Name", activity.EntityName)
			s.Equal("updated", activity.Action, "Modified SME should show as 'updated'")
			s.Equal(s.testUser2.ID, activity.UserID, "Should show the user who updated")
			s.Equal(s.testUser2.Name, activity.UserName)
			break
		}
	}
	s.True(found, "Updated SME should appear in activity history")
}

func (s *EntityActivityTrackingTestSuite) TestSME_SoftDelete_RemovedFromActivityHistory() {
	// Create an SME
	userID := int(s.testUser.ID)
	sme := &models.Sme{
		UsmeNumber:       fmt.Sprintf("USME-DELETE-%d", time.Now().UnixNano()),
		Name:             "SME To Delete",
		ContactPhone:     "+265999111222",
		ContactEmail:     "delete@test.com",
		BusinessCategory: "Agriculture",
		Sector:           "Agribusiness",
		CreatedBy:        &userID,
		UpdatedBy:        &userID,
	}
	facades.Orm().Query().Create(sme)
	smeID := sme.ID

	// Verify it appears in activity history
	activities := s.dashboardService.GetRecentActivities(20)
	found := false
	for _, activity := range activities {
		if activity.EntityType == "sme" && activity.EntityID == smeID {
			found = true
			break
		}
	}
	s.True(found, "SME should appear before deletion")

	// Soft delete the SME (update deleted_at like the service does)
	_, err := facades.Orm().Query().Exec(
		"UPDATE smes SET deleted_at = NOW() WHERE id = ?", smeID,
	)
	s.NoError(err)

	// Verify it no longer appears in activity history
	activities = s.dashboardService.GetRecentActivities(20)
	found = false
	for _, activity := range activities {
		if activity.EntityType == "sme" && activity.EntityID == smeID {
			found = true
			break
		}
	}
	s.False(found, "Soft-deleted SME should not appear in activity history")
}

// ===========================================
// Event CRUD Activity Tracking Tests
// ===========================================

func (s *EntityActivityTrackingTestSuite) TestEvent_Create_AppearsInActivityHistory() {
	userID := int(s.testUser.ID)

	// Create an Event with proper DateTime type and required fields
	futureDate := carbon.NewDateTime(carbon.Now().AddDays(30))
	event := &models.Event{
		Title:         "New Conference Event",
		Venue:         "Conference Center",
		Date:          *futureDate,
		Partners:      []string{},
		AttendingSmes: []int{},
		CreatedBy:     &userID,
		UpdatedBy:     &userID,
	}
	err := facades.Orm().Query().Create(event)
	s.NoError(err)
	s.NotZero(event.ID)

	// Get recent activities
	activities := s.dashboardService.GetRecentActivities(20)

	// Find the Event in activity history
	found := false
	for _, activity := range activities {
		if activity.EntityType == "event" && activity.EntityID == event.ID {
			found = true
			s.Equal("New Conference Event", activity.EntityName)
			s.Equal("created", activity.Action)
			s.Equal(s.testUser.ID, activity.UserID)
			s.Equal(s.testUser.Name, activity.UserName)
			break
		}
	}
	s.True(found, "Created Event should appear in activity history")
}

func (s *EntityActivityTrackingTestSuite) TestEvent_Update_AppearsInActivityHistory() {
	userID := int(s.testUser.ID)

	// Create an Event
	futureDate := carbon.NewDateTime(carbon.Now().AddDays(30))
	event := &models.Event{
		Title:         "Original Event Title",
		Venue:         "Original Venue",
		Date:          *futureDate,
		Partners:      []string{},
		AttendingSmes: []int{},
		CreatedBy:     &userID,
		UpdatedBy:     &userID,
	}
	facades.Orm().Query().Create(event)

	// Wait long enough so updated_at will differ from created_at in seconds
	time.Sleep(1100 * time.Millisecond)

	user2ID := int(s.testUser2.ID)
	_, err := facades.Orm().Query().Exec(
		"UPDATE events SET title = ?, venue = ?, updated_by = ?, updated_at = NOW() WHERE id = ?",
		"Updated Event Title", "New Venue", user2ID, event.ID,
	)
	s.NoError(err)

	// Get recent activities
	activities := s.dashboardService.GetRecentActivities(20)

	// Find the updated Event
	found := false
	for _, activity := range activities {
		if activity.EntityType == "event" && activity.EntityID == event.ID {
			found = true
			s.Equal("Updated Event Title", activity.EntityName)
			s.Equal("updated", activity.Action)
			s.Equal(s.testUser2.ID, activity.UserID)
			break
		}
	}
	s.True(found, "Updated Event should appear in activity history")
}

func (s *EntityActivityTrackingTestSuite) TestEvent_SoftDelete_RemovedFromActivityHistory() {
	userID := int(s.testUser.ID)
	futureDate := carbon.NewDateTime(carbon.Now().AddDays(30))

	event := &models.Event{
		Title:         "Event To Delete",
		Venue:         "Some Venue",
		Date:          *futureDate,
		Partners:      []string{},
		AttendingSmes: []int{},
		CreatedBy:     &userID,
		UpdatedBy:     &userID,
	}
	facades.Orm().Query().Create(event)
	eventID := event.ID

	// Verify it appears
	activities := s.dashboardService.GetRecentActivities(20)
	found := false
	for _, activity := range activities {
		if activity.EntityType == "event" && activity.EntityID == eventID {
			found = true
			break
		}
	}
	s.True(found, "Event should appear before deletion")

	// Soft delete (update deleted_at like the service does)
	_, err := facades.Orm().Query().Exec(
		"UPDATE events SET deleted_at = NOW() WHERE id = ?", eventID,
	)
	s.NoError(err)

	// Verify removed
	activities = s.dashboardService.GetRecentActivities(20)
	found = false
	for _, activity := range activities {
		if activity.EntityType == "event" && activity.EntityID == eventID {
			found = true
			break
		}
	}
	s.False(found, "Soft-deleted Event should not appear in activity history")
}

// ===========================================
// Procurement Notice CRUD Activity Tracking Tests
// ===========================================

func (s *EntityActivityTrackingTestSuite) TestProcurement_Create_AppearsInActivityHistory() {
	userID := int(s.testUser.ID)

	// Create a Procurement Notice with all required fields
	futureDate := carbon.NewDateTime(carbon.Now().AddDays(30))
	openDate := carbon.NewDateTime(carbon.Now())
	procurement := &models.ProcurementNotice{
		Organization:        "Government Agency",
		RefNo:               fmt.Sprintf("REF-%d", time.Now().UnixNano()),
		OpenDate:            *openDate,
		CloseDate:           *futureDate,
		IsPublished:         true,
		ProcurementType:     "Goods",
		Partners:            []string{},
		QualifyingDistricts: []string{},
		Classification:      []string{},
		InterestedSmes:      []string{},
		CreatedBy:           &userID,
		UpdatedBy:           &userID,
	}
	err := facades.Orm().Query().Create(procurement)
	s.NoError(err)
	s.NotZero(procurement.ID)

	// Get recent activities
	activities := s.dashboardService.GetRecentActivities(20)

	// Find the Procurement in activity history
	found := false
	for _, activity := range activities {
		if activity.EntityType == "procurement" && activity.EntityID == procurement.ID {
			found = true
			s.Contains(activity.EntityName, "Government Agency")
			s.Equal("created", activity.Action)
			s.Equal(s.testUser.ID, activity.UserID)
			s.Equal(s.testUser.Name, activity.UserName)
			break
		}
	}
	s.True(found, "Created Procurement should appear in activity history")
}

func (s *EntityActivityTrackingTestSuite) TestProcurement_Update_AppearsInActivityHistory() {
	userID := int(s.testUser.ID)
	futureDate := carbon.NewDateTime(carbon.Now().AddDays(30))
	openDate := carbon.NewDateTime(carbon.Now())

	procurement := &models.ProcurementNotice{
		Organization:        "Original Org",
		RefNo:               fmt.Sprintf("REF-%d", time.Now().UnixNano()),
		OpenDate:            *openDate,
		CloseDate:           *futureDate,
		IsPublished:         true,
		ProcurementType:     "Goods",
		Partners:            []string{},
		QualifyingDistricts: []string{},
		Classification:      []string{},
		InterestedSmes:      []string{},
		CreatedBy:           &userID,
		UpdatedBy:           &userID,
	}
	facades.Orm().Query().Create(procurement)

	// Wait long enough so updated_at will differ from created_at in seconds
	time.Sleep(1100 * time.Millisecond)

	user2ID := int(s.testUser2.ID)
	_, err := facades.Orm().Query().Exec(
		"UPDATE procurement_notices SET organization = ?, updated_by = ?, updated_at = NOW() WHERE id = ?",
		"Updated Organization", user2ID, procurement.ID,
	)
	s.NoError(err)

	// Get recent activities
	activities := s.dashboardService.GetRecentActivities(20)

	// Find the updated Procurement
	found := false
	for _, activity := range activities {
		if activity.EntityType == "procurement" && activity.EntityID == procurement.ID {
			found = true
			s.Contains(activity.EntityName, "Updated Organization")
			s.Equal("updated", activity.Action)
			s.Equal(s.testUser2.ID, activity.UserID)
			break
		}
	}
	s.True(found, "Updated Procurement should appear in activity history")
}

func (s *EntityActivityTrackingTestSuite) TestProcurement_SoftDelete_RemovedFromActivityHistory() {
	userID := int(s.testUser.ID)
	futureDate := carbon.NewDateTime(carbon.Now().AddDays(30))
	openDate := carbon.NewDateTime(carbon.Now())

	procurement := &models.ProcurementNotice{
		Organization:        "Org To Delete",
		RefNo:               fmt.Sprintf("REF-%d", time.Now().UnixNano()),
		OpenDate:            *openDate,
		CloseDate:           *futureDate,
		IsPublished:         true,
		ProcurementType:     "Services",
		Partners:            []string{},
		QualifyingDistricts: []string{},
		Classification:      []string{},
		InterestedSmes:      []string{},
		CreatedBy:           &userID,
		UpdatedBy:           &userID,
	}
	facades.Orm().Query().Create(procurement)
	procID := procurement.ID

	// Verify it appears
	activities := s.dashboardService.GetRecentActivities(20)
	found := false
	for _, activity := range activities {
		if activity.EntityType == "procurement" && activity.EntityID == procID {
			found = true
			break
		}
	}
	s.True(found, "Procurement should appear before deletion")

	// Soft delete (update deleted_at like the service does)
	_, err := facades.Orm().Query().Exec(
		"UPDATE procurement_notices SET deleted_at = NOW() WHERE id = ?", procID,
	)
	s.NoError(err)

	// Verify removed
	activities = s.dashboardService.GetRecentActivities(20)
	found = false
	for _, activity := range activities {
		if activity.EntityType == "procurement" && activity.EntityID == procID {
			found = true
			break
		}
	}
	s.False(found, "Soft-deleted Procurement should not appear in activity history")
}

// ===========================================
// Cross-Entity Activity Tracking Tests
// ===========================================

func (s *EntityActivityTrackingTestSuite) TestMixedEntities_AllAppearInActivityHistory() {
	userID := int(s.testUser.ID)
	futureDate := carbon.NewDateTime(carbon.Now().AddDays(30))
	openDate := carbon.NewDateTime(carbon.Now())

	// Create one of each entity type
	sme := &models.Sme{
		UsmeNumber:       fmt.Sprintf("USME-MIX-%d", time.Now().UnixNano()),
		Name:             "Mixed Test SME",
		ContactPhone:     "+265999111222",
		ContactEmail:     "mixed@test.com",
		BusinessCategory: "Agriculture",
		Sector:           "Agribusiness",
		CreatedBy:        &userID,
		UpdatedBy:        &userID,
	}
	facades.Orm().Query().Create(sme)

	event := &models.Event{
		Title:         "Mixed Test Event",
		Venue:         "Test Venue",
		Date:          *futureDate,
		Partners:      []string{},
		AttendingSmes: []int{},
		CreatedBy:     &userID,
		UpdatedBy:     &userID,
	}
	facades.Orm().Query().Create(event)

	procurement := &models.ProcurementNotice{
		Organization:        "Mixed Test Org",
		RefNo:               fmt.Sprintf("MIX-REF-%d", time.Now().UnixNano()),
		OpenDate:            *openDate,
		CloseDate:           *futureDate,
		IsPublished:         true,
		ProcurementType:     "Works",
		Partners:            []string{},
		QualifyingDistricts: []string{},
		Classification:      []string{},
		InterestedSmes:      []string{},
		CreatedBy:           &userID,
		UpdatedBy:           &userID,
	}
	facades.Orm().Query().Create(procurement)

	// Get recent activities
	activities := s.dashboardService.GetRecentActivities(20)

	// Verify all three entity types appear
	foundSME := false
	foundEvent := false
	foundProcurement := false

	for _, activity := range activities {
		if activity.EntityType == "sme" && activity.EntityID == sme.ID {
			foundSME = true
			s.Equal("Mixed Test SME", activity.EntityName)
		}
		if activity.EntityType == "event" && activity.EntityID == event.ID {
			foundEvent = true
			s.Equal("Mixed Test Event", activity.EntityName)
		}
		if activity.EntityType == "procurement" && activity.EntityID == procurement.ID {
			foundProcurement = true
			s.Contains(activity.EntityName, "Mixed Test Org")
		}
	}

	s.True(foundSME, "SME should appear in mixed activity list")
	s.True(foundEvent, "Event should appear in mixed activity list")
	s.True(foundProcurement, "Procurement should appear in mixed activity list")
}

func (s *EntityActivityTrackingTestSuite) TestUserSpecificActivities_FilteredCorrectly() {
	user1ID := int(s.testUser.ID)
	user2ID := int(s.testUser2.ID)

	// Create SME by user 1
	sme1 := &models.Sme{
		UsmeNumber:       fmt.Sprintf("USME-U1-%d", time.Now().UnixNano()),
		Name:             "User1 SME",
		ContactPhone:     "+265999111222",
		ContactEmail:     "user1@test.com",
		BusinessCategory: "Agriculture",
		Sector:           "Agribusiness",
		CreatedBy:        &user1ID,
		UpdatedBy:        &user1ID,
	}
	facades.Orm().Query().Create(sme1)

	// Create SME by user 2
	sme2 := &models.Sme{
		UsmeNumber:       fmt.Sprintf("USME-U2-%d", time.Now().UnixNano()),
		Name:             "User2 SME",
		ContactPhone:     "+265999333444",
		ContactEmail:     "user2@test.com",
		BusinessCategory: "Manufacturing",
		Sector:           "Industrial",
		CreatedBy:        &user2ID,
		UpdatedBy:        &user2ID,
	}
	facades.Orm().Query().Create(sme2)

	// Get user-specific activities for user 1
	user1Activities := s.dashboardService.GetUserRecentActivities(s.testUser.ID, 20)

	// Verify user 1 only sees their own activities
	for _, activity := range user1Activities {
		s.Equal(s.testUser.ID, activity.UserID, "User 1 should only see their own activities")
	}

	// Should find user1's SME
	foundUser1SME := false
	for _, activity := range user1Activities {
		if activity.EntityType == "sme" && activity.EntityID == sme1.ID {
			foundUser1SME = true
			s.Equal("User1 SME", activity.EntityName)
			break
		}
	}
	s.True(foundUser1SME, "User 1 should see their own SME")

	// Should NOT find user2's SME
	for _, activity := range user1Activities {
		if activity.EntityType == "sme" && activity.EntityID == sme2.ID {
			s.Fail("User 1 should not see User 2's SME in filtered results")
		}
	}
}

func (s *EntityActivityTrackingTestSuite) TestActivityOrder_MostRecentFirst() {
	userID := int(s.testUser.ID)

	// Create entities with delays to ensure different timestamps
	// Need 1100ms delays since ToDateTimeString() only has second precision
	sme1 := &models.Sme{
		UsmeNumber:       fmt.Sprintf("USME-ORD1-%d", time.Now().UnixNano()),
		Name:             "First SME",
		ContactPhone:     "+265999111222",
		ContactEmail:     "first@test.com",
		BusinessCategory: "Agriculture",
		Sector:           "Agribusiness",
		CreatedBy:        &userID,
		UpdatedBy:        &userID,
	}
	facades.Orm().Query().Create(sme1)

	time.Sleep(1100 * time.Millisecond)

	sme2 := &models.Sme{
		UsmeNumber:       fmt.Sprintf("USME-ORD2-%d", time.Now().UnixNano()),
		Name:             "Second SME",
		ContactPhone:     "+265999333444",
		ContactEmail:     "second@test.com",
		BusinessCategory: "Manufacturing",
		Sector:           "Industrial",
		CreatedBy:        &userID,
		UpdatedBy:        &userID,
	}
	facades.Orm().Query().Create(sme2)

	time.Sleep(1100 * time.Millisecond)

	sme3 := &models.Sme{
		UsmeNumber:       fmt.Sprintf("USME-ORD3-%d", time.Now().UnixNano()),
		Name:             "Third SME",
		ContactPhone:     "+265999555666",
		ContactEmail:     "third@test.com",
		BusinessCategory: "Services",
		Sector:           "Technology",
		CreatedBy:        &userID,
		UpdatedBy:        &userID,
	}
	facades.Orm().Query().Create(sme3)

	// Get activities
	activities := s.dashboardService.GetRecentActivities(20)

	// Find positions of our SMEs
	positions := make(map[string]int)
	for i, activity := range activities {
		if activity.EntityType == "sme" {
			if activity.EntityID == sme1.ID {
				positions["first"] = i
			} else if activity.EntityID == sme2.ID {
				positions["second"] = i
			} else if activity.EntityID == sme3.ID {
				positions["third"] = i
			}
		}
	}

	// Third (most recent) should come before Second, which should come before First
	if pos1, ok1 := positions["first"]; ok1 {
		if pos2, ok2 := positions["second"]; ok2 {
			s.Greater(pos1, pos2, "First SME should appear after Second SME (most recent first)")
		}
		if pos3, ok3 := positions["third"]; ok3 {
			s.Greater(pos1, pos3, "First SME should appear after Third SME (most recent first)")
		}
	}
	if pos2, ok2 := positions["second"]; ok2 {
		if pos3, ok3 := positions["third"]; ok3 {
			s.Greater(pos2, pos3, "Second SME should appear after Third SME (most recent first)")
		}
	}
}

func (s *EntityActivityTrackingTestSuite) TestMultipleUpdates_ShowsLatestState() {
	userID := int(s.testUser.ID)
	user2ID := int(s.testUser2.ID)

	// Create an SME
	sme := &models.Sme{
		UsmeNumber:       fmt.Sprintf("USME-MULTI-%d", time.Now().UnixNano()),
		Name:             "Version 1",
		ContactPhone:     "+265999111222",
		ContactEmail:     "multi@test.com",
		BusinessCategory: "Agriculture",
		Sector:           "Agribusiness",
		CreatedBy:        &userID,
		UpdatedBy:        &userID,
	}
	facades.Orm().Query().Create(sme)

	// Multiple updates - need 1100ms delay so updated_at differs from created_at in seconds
	time.Sleep(1100 * time.Millisecond)
	facades.Orm().Query().Exec(
		"UPDATE smes SET name = ?, updated_by = ?, updated_at = NOW() WHERE id = ?",
		"Version 2", userID, sme.ID,
	)

	time.Sleep(100 * time.Millisecond)
	facades.Orm().Query().Exec(
		"UPDATE smes SET name = ?, updated_by = ?, updated_at = NOW() WHERE id = ?",
		"Version 3 - Final", user2ID, sme.ID,
	)

	// Get activities
	activities := s.dashboardService.GetRecentActivities(20)

	// Should show the latest state
	found := false
	for _, activity := range activities {
		if activity.EntityType == "sme" && activity.EntityID == sme.ID {
			found = true
			s.Equal("Version 3 - Final", activity.EntityName, "Should show latest name")
			s.Equal("updated", activity.Action)
			s.Equal(s.testUser2.ID, activity.UserID, "Should show last updater")
			break
		}
	}
	s.True(found, "SME should appear with latest state")
}
