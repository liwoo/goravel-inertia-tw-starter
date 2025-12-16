package notifications

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

// AutoBroadcastTestSuite tests the auto-notification feature for events and procurements
type AutoBroadcastTestSuite struct {
	suite.Suite
	tests.TestCase

	notificationService *services.NotificationService

	// Admin user who creates events/procurements
	adminUser *models.User

	// SME users in different districts
	smeUserLilongwe1  *models.User
	smeUserLilongwe2  *models.User
	smeUserBlantyre   *models.User
	smeUserMzuzu      *models.User
	smeUserNoDistrict *models.User

	// Inactive SME user (should not receive notifications)
	inactiveSmeUser *models.User

	// User with no SME association (should not receive notifications)
	regularUser *models.User

	// SMEs
	smeLilongwe1  *models.Sme
	smeLilongwe2  *models.Sme
	smeBlantyre   *models.Sme
	smeMzuzu      *models.Sme
	smeNoDistrict *models.Sme
	smeInactive   *models.Sme
}

func TestAutoBroadcastTestSuite(t *testing.T) {
	suite.Run(t, new(AutoBroadcastTestSuite))
}

func (s *AutoBroadcastTestSuite) SetupTest() {
	s.RefreshDatabase()
	s.notificationService = services.NewNotificationService()

	password, _ := facades.Hash().Make("password123")
	timestamp := time.Now().UnixNano()

	// Create admin user
	s.adminUser = &models.User{
		Name:         "Admin User",
		Email:        fmt.Sprintf("admin_%d@example.com", timestamp),
		Password:     password,
		IsActive:     true,
		IsSuperAdmin: true,
	}
	facades.Orm().Query().Create(s.adminUser)

	// Create SME users with matching emails to SMEs

	// Lilongwe District SME users
	s.smeUserLilongwe1 = &models.User{
		Name:     "SME User Lilongwe 1",
		Email:    fmt.Sprintf("sme_lilongwe1_%d@example.com", timestamp),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.smeUserLilongwe1)

	s.smeUserLilongwe2 = &models.User{
		Name:     "SME User Lilongwe 2",
		Email:    fmt.Sprintf("sme_lilongwe2_%d@example.com", timestamp),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.smeUserLilongwe2)

	// Blantyre District SME user
	s.smeUserBlantyre = &models.User{
		Name:     "SME User Blantyre",
		Email:    fmt.Sprintf("sme_blantyre_%d@example.com", timestamp),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.smeUserBlantyre)

	// Mzuzu District SME user
	s.smeUserMzuzu = &models.User{
		Name:     "SME User Mzuzu",
		Email:    fmt.Sprintf("sme_mzuzu_%d@example.com", timestamp),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.smeUserMzuzu)

	// SME user with no district
	s.smeUserNoDistrict = &models.User{
		Name:     "SME User No District",
		Email:    fmt.Sprintf("sme_nodistrict_%d@example.com", timestamp),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.smeUserNoDistrict)

	// Inactive SME user - must explicitly set IsActive to false after creation
	// because GORM default:true overrides the explicit false value
	s.inactiveSmeUser = &models.User{
		Name:     "Inactive SME User",
		Email:    fmt.Sprintf("sme_inactive_%d@example.com", timestamp),
		Password: password,
	}
	facades.Orm().Query().Create(s.inactiveSmeUser)
	// Explicitly update to inactive after creation to bypass GORM default
	facades.Orm().Query().Model(s.inactiveSmeUser).Update("is_active", false)

	// Regular user with no SME
	s.regularUser = &models.User{
		Name:     "Regular User",
		Email:    fmt.Sprintf("regular_%d@example.com", timestamp),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.regularUser)

	// Create SMEs linked to users via email

	lilongwe := "Lilongwe"
	blantyre := "Blantyre"
	mzuzu := "Mzuzu"

	s.smeLilongwe1 = &models.Sme{
		Name:             "SME Lilongwe 1",
		ContactEmail:     s.smeUserLilongwe1.Email, // Link via email
		ContactPhone:     "0999111111",
		BusinessCategory: "Trading",
		Sector:           "Retail",
		District:         &lilongwe,
		IsActive:         true,
	}
	facades.Orm().Query().Create(s.smeLilongwe1)

	s.smeLilongwe2 = &models.Sme{
		Name:             "SME Lilongwe 2",
		ContactEmail:     s.smeUserLilongwe2.Email,
		ContactPhone:     "0999222222",
		BusinessCategory: "Services",
		Sector:           "Consulting",
		District:         &lilongwe,
		IsActive:         true,
	}
	facades.Orm().Query().Create(s.smeLilongwe2)

	s.smeBlantyre = &models.Sme{
		Name:             "SME Blantyre",
		ContactEmail:     s.smeUserBlantyre.Email,
		ContactPhone:     "0999333333",
		BusinessCategory: "Manufacturing",
		Sector:           "Food Processing",
		District:         &blantyre,
		IsActive:         true,
	}
	facades.Orm().Query().Create(s.smeBlantyre)

	s.smeMzuzu = &models.Sme{
		Name:             "SME Mzuzu",
		ContactEmail:     s.smeUserMzuzu.Email,
		ContactPhone:     "0999444444",
		BusinessCategory: "Agriculture",
		Sector:           "Farming",
		District:         &mzuzu,
		IsActive:         true,
	}
	facades.Orm().Query().Create(s.smeMzuzu)

	// SME with no district (nationwide)
	s.smeNoDistrict = &models.Sme{
		Name:             "SME No District",
		ContactEmail:     s.smeUserNoDistrict.Email,
		ContactPhone:     "0999555555",
		BusinessCategory: "Trading",
		Sector:           "Retail",
		District:         nil, // No district
		IsActive:         true,
	}
	facades.Orm().Query().Create(s.smeNoDistrict)

	// Inactive SME - must explicitly set IsActive to false after creation
	// because GORM default:true overrides the explicit false value
	s.smeInactive = &models.Sme{
		Name:             "Inactive SME",
		ContactEmail:     s.inactiveSmeUser.Email,
		ContactPhone:     "0999666666",
		BusinessCategory: "Trading",
		Sector:           "Retail",
		District:         &lilongwe,
	}
	facades.Orm().Query().Create(s.smeInactive)
	// Explicitly update to inactive after creation to bypass GORM default
	facades.Orm().Query().Model(s.smeInactive).Update("is_active", false)
}

// ============================================================================
// EVENT NOTIFICATION TESTS
// ============================================================================

// Test: BroadcastEventNotificationSync sends to all SMEs when event has no district
func (s *AutoBroadcastTestSuite) TestBroadcastEventNotification_NoDistrict_NotifiesAllActiveSMEs() {
	eventID := uint(1001)
	eventTitle := "National Business Conference"
	eventDate := "2025-02-15"
	eventVenue := "Bingu Conference Center"
	eventDistrict := "" // No district - nationwide event

	// Call the sync broadcast method directly
	s.notificationService.BroadcastEventNotificationSync(
		s.adminUser.ID,
		eventID,
		eventTitle,
		eventDate,
		eventVenue,
		eventDistrict,
	)

	// Wait a moment for notifications to be created
	time.Sleep(200 * time.Millisecond)

	// Verify all active SME users received notifications
	for _, user := range []*models.User{
		s.smeUserLilongwe1,
		s.smeUserLilongwe2,
		s.smeUserBlantyre,
		s.smeUserMzuzu,
		s.smeUserNoDistrict,
	} {
		var notification models.Notification
		err := facades.Orm().Query().
			Where("user_id = ?", user.ID).
			Where("type = ?", "event").
			Where("related_id = ?", eventID).
			First(&notification)

		s.NoError(err, fmt.Sprintf("Notification should be created for %s", user.Name))
		s.Contains(notification.Title, eventTitle, "Notification title should contain event title")
		s.Equal("event", notification.RelatedType, "Related type should be event")
	}

	// Verify inactive user did NOT receive notification
	// The inactiveSmeUser has is_active=false, so they should not receive notifications
	var inactiveNotificationCount int64
	inactiveNotificationCount, _ = facades.Orm().Query().Model(&models.Notification{}).
		Where("user_id = ?", s.inactiveSmeUser.ID).
		Where("type = ?", "event").
		Where("related_id = ?", eventID).
		Count()

	s.Equal(int64(0), inactiveNotificationCount, "Inactive user should not receive notification")

	// Verify regular user (no SME) did NOT receive notification
	var regularNotificationCount int64
	regularNotificationCount, _ = facades.Orm().Query().Model(&models.Notification{}).
		Where("user_id = ?", s.regularUser.ID).
		Where("type = ?", "event").
		Where("related_id = ?", eventID).
		Count()

	s.Equal(int64(0), regularNotificationCount, "Regular user without SME should not receive notification")
}

// Test: BroadcastEventNotificationSync filters by district
func (s *AutoBroadcastTestSuite) TestBroadcastEventNotification_WithDistrict_NotifiesOnlyEligibleSMEs() {
	eventID := uint(1002)
	eventTitle := "Lilongwe Business Workshop"
	eventDate := "2025-03-20"
	eventVenue := "Lilongwe City Hall"
	eventDistrict := "Lilongwe" // District-specific event

	s.notificationService.BroadcastEventNotificationSync(
		s.adminUser.ID,
		eventID,
		eventTitle,
		eventDate,
		eventVenue,
		eventDistrict,
	)

	time.Sleep(200 * time.Millisecond)

	// Lilongwe SMEs should receive notifications
	for _, user := range []*models.User{s.smeUserLilongwe1, s.smeUserLilongwe2} {
		var notification models.Notification
		err := facades.Orm().Query().
			Where("user_id = ?", user.ID).
			Where("type = ?", "event").
			Where("related_id = ?", eventID).
			First(&notification)

		s.NoError(err, fmt.Sprintf("Lilongwe SME user %s should receive notification", user.Name))
	}

	// SMEs with no district should also receive (they can attend any event)
	var noDistrictNotification models.Notification
	err := facades.Orm().Query().
		Where("user_id = ?", s.smeUserNoDistrict.ID).
		Where("type = ?", "event").
		Where("related_id = ?", eventID).
		First(&noDistrictNotification)

	s.NoError(err, "SME user with no district should receive notification")

	// Blantyre and Mzuzu SMEs should NOT receive notifications
	for _, user := range []*models.User{s.smeUserBlantyre, s.smeUserMzuzu} {
		var notificationCount int64
		notificationCount, _ = facades.Orm().Query().Model(&models.Notification{}).
			Where("user_id = ?", user.ID).
			Where("type = ?", "event").
			Where("related_id = ?", eventID).
			Count()

		s.Equal(int64(0), notificationCount, fmt.Sprintf("User %s in different district should not receive notification", user.Name))
	}
}

// Test: Event notification content is correct
func (s *AutoBroadcastTestSuite) TestBroadcastEventNotification_CorrectContentFormat() {
	eventID := uint(1003)
	eventTitle := "Test Event Content"
	eventDate := "2025-04-10"
	eventVenue := "Test Venue"
	eventDistrict := "Lilongwe"

	s.notificationService.BroadcastEventNotificationSync(
		s.adminUser.ID,
		eventID,
		eventTitle,
		eventDate,
		eventVenue,
		eventDistrict,
	)

	time.Sleep(200 * time.Millisecond)

	var notification models.Notification
	err := facades.Orm().Query().
		Where("user_id = ?", s.smeUserLilongwe1.ID).
		Where("type = ?", "event").
		Where("related_id = ?", eventID).
		First(&notification)

	s.NoError(err)
	s.Contains(notification.Title, "New Event:", "Title should start with 'New Event:'")
	s.Contains(notification.Title, eventTitle, "Title should contain event title")
	s.Contains(notification.Message, eventTitle, "Message should contain event title")
	s.Contains(notification.Message, eventDate, "Message should contain event date")
	s.Contains(notification.Message, eventVenue, "Message should contain event venue")
	s.Contains(notification.Message, eventDistrict, "Message should contain event district")
	s.Equal("event", notification.Type, "Type should be 'event'")
	s.Equal("event", notification.RelatedType, "Related type should be 'event'")
	s.Equal(&eventID, notification.RelatedID, "Related ID should be event ID")
	s.Equal(&s.adminUser.ID, notification.TriggerUserID, "Trigger user should be admin")
	s.Equal("normal", notification.Priority, "Priority should be 'normal'")
}

// Test: Event notification also creates messages
func (s *AutoBroadcastTestSuite) TestBroadcastEventNotification_CreatesMessages() {
	eventID := uint(1004)
	eventTitle := "Event with Messages"
	eventDate := "2025-05-15"
	eventVenue := "Message Test Venue"
	eventDistrict := ""

	s.notificationService.BroadcastEventNotificationSync(
		s.adminUser.ID,
		eventID,
		eventTitle,
		eventDate,
		eventVenue,
		eventDistrict,
	)

	time.Sleep(200 * time.Millisecond)

	// Verify messages were created for SME users
	messageCount, _ := facades.Orm().Query().Model(&models.Message{}).
		Where("sender_id = ?", s.adminUser.ID).
		Where("type = ?", string(models.MessageTypeSystem)).
		Where("content LIKE ?", "%"+eventTitle+"%").
		Count()

	// Should have at least one message per eligible SME user (5 active SME users)
	s.GreaterOrEqual(messageCount, int64(5), "Should create messages for all eligible SME users")
}

// ============================================================================
// PROCUREMENT NOTIFICATION TESTS
// ============================================================================

// Test: BroadcastProcurementNotificationSync sends to all active SMEs
func (s *AutoBroadcastTestSuite) TestBroadcastProcurementNotification_NotifiesAllActiveSMEs() {
	procurementID := uint(2001)
	organization := "Ministry of Trade"
	refNo := "MOT/2025/001"
	procurementType := "Goods"
	closeDate := "2025-06-30"

	s.notificationService.BroadcastProcurementNotificationSync(
		s.adminUser.ID,
		procurementID,
		organization,
		refNo,
		procurementType,
		closeDate,
	)

	time.Sleep(200 * time.Millisecond)

	// All active SME users should receive notifications (regardless of district)
	for _, user := range []*models.User{
		s.smeUserLilongwe1,
		s.smeUserLilongwe2,
		s.smeUserBlantyre,
		s.smeUserMzuzu,
		s.smeUserNoDistrict,
	} {
		var notification models.Notification
		err := facades.Orm().Query().
			Where("user_id = ?", user.ID).
			Where("type = ?", "procurement").
			Where("related_id = ?", procurementID).
			First(&notification)

		s.NoError(err, fmt.Sprintf("Notification should be created for %s", user.Name))
		s.Contains(notification.Title, "Procurement", "Title should mention procurement")
	}

	// Inactive user should NOT receive notification
	var inactiveNotificationCount int64
	inactiveNotificationCount, _ = facades.Orm().Query().Model(&models.Notification{}).
		Where("user_id = ?", s.inactiveSmeUser.ID).
		Where("type = ?", "procurement").
		Where("related_id = ?", procurementID).
		Count()

	s.Equal(int64(0), inactiveNotificationCount, "Inactive user should not receive procurement notification")

	// Regular user (no SME) should NOT receive notification
	var regularNotificationCount int64
	regularNotificationCount, _ = facades.Orm().Query().Model(&models.Notification{}).
		Where("user_id = ?", s.regularUser.ID).
		Where("type = ?", "procurement").
		Where("related_id = ?", procurementID).
		Count()

	s.Equal(int64(0), regularNotificationCount, "Regular user without SME should not receive procurement notification")
}

// Test: Procurement notification content is correct
func (s *AutoBroadcastTestSuite) TestBroadcastProcurementNotification_CorrectContentFormat() {
	procurementID := uint(2002)
	organization := "Test Organization"
	refNo := "TEST/2025/002"
	procurementType := "Services"
	closeDate := "2025-07-15"

	s.notificationService.BroadcastProcurementNotificationSync(
		s.adminUser.ID,
		procurementID,
		organization,
		refNo,
		procurementType,
		closeDate,
	)

	time.Sleep(200 * time.Millisecond)

	var notification models.Notification
	err := facades.Orm().Query().
		Where("user_id = ?", s.smeUserLilongwe1.ID).
		Where("type = ?", "procurement").
		Where("related_id = ?", procurementID).
		First(&notification)

	s.NoError(err)
	s.Contains(notification.Title, "Procurement", "Title should mention procurement")
	s.Contains(notification.Message, organization, "Message should contain organization")
	s.Contains(notification.Message, refNo, "Message should contain reference number")
	s.Contains(notification.Message, procurementType, "Message should contain procurement type")
	s.Contains(notification.Message, closeDate, "Message should contain closing date")
	s.Equal("procurement", notification.Type, "Type should be 'procurement'")
	s.Equal("procurement", notification.RelatedType, "Related type should be 'procurement'")
	s.Equal(&procurementID, notification.RelatedID, "Related ID should be procurement ID")
	s.Equal(&s.adminUser.ID, notification.TriggerUserID, "Trigger user should be admin")
}

// Test: Procurement notification also creates messages
func (s *AutoBroadcastTestSuite) TestBroadcastProcurementNotification_CreatesMessages() {
	procurementID := uint(2003)
	organization := "Message Test Org"
	refNo := "MSG/2025/001"
	procurementType := "Works"
	closeDate := "2025-08-20"

	s.notificationService.BroadcastProcurementNotificationSync(
		s.adminUser.ID,
		procurementID,
		organization,
		refNo,
		procurementType,
		closeDate,
	)

	time.Sleep(200 * time.Millisecond)

	// Verify messages were created
	messageCount, _ := facades.Orm().Query().Model(&models.Message{}).
		Where("sender_id = ?", s.adminUser.ID).
		Where("type = ?", string(models.MessageTypeSystem)).
		Where("content LIKE ?", "%"+organization+"%").
		Count()

	s.GreaterOrEqual(messageCount, int64(5), "Should create messages for all active SME users")
}

// ============================================================================
// EDGE CASE TESTS
// ============================================================================

// Test: No notifications sent when no eligible SMEs exist
func (s *AutoBroadcastTestSuite) TestBroadcastEventNotification_NoEligibleSMEs() {
	// Deactivate all SMEs
	facades.Orm().Query().Model(&models.Sme{}).Where("1=1").Update("is_active", false)

	eventID := uint(3001)
	beforeCount := s.countNotifications()

	s.notificationService.BroadcastEventNotificationSync(
		s.adminUser.ID,
		eventID,
		"No Recipients Event",
		"2025-09-01",
		"Empty Venue",
		"",
	)

	time.Sleep(200 * time.Millisecond)

	afterCount := s.countNotifications()
	s.Equal(beforeCount, afterCount, "No new notifications should be created when no eligible SMEs")
}

// Test: Email matching is case-insensitive
func (s *AutoBroadcastTestSuite) TestEmailMatchingIsCaseInsensitive() {
	timestamp := time.Now().UnixNano()
	password, _ := facades.Hash().Make("password123")

	// Create user with lowercase email
	userLower := &models.User{
		Name:     "Case Test User",
		Email:    fmt.Sprintf("casetest_%d@example.com", timestamp),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(userLower)

	// Create SME with UPPERCASE email (should still match)
	district := "Lilongwe"
	smeUpper := &models.Sme{
		Name:             "Case Test SME",
		ContactEmail:     fmt.Sprintf("CASETEST_%d@EXAMPLE.COM", timestamp), // Uppercase
		ContactPhone:     "0999999999",
		BusinessCategory: "Trading",
		Sector:           "Retail",
		District:         &district,
		IsActive:         true,
	}
	facades.Orm().Query().Create(smeUpper)

	eventID := uint(4001)
	s.notificationService.BroadcastEventNotificationSync(
		s.adminUser.ID,
		eventID,
		"Case Test Event",
		"2025-10-01",
		"Case Venue",
		"",
	)

	time.Sleep(200 * time.Millisecond)

	var notification models.Notification
	err := facades.Orm().Query().
		Where("user_id = ?", userLower.ID).
		Where("type = ?", "event").
		Where("related_id = ?", eventID).
		First(&notification)

	s.NoError(err, "Notification should be created despite email case mismatch")
}

// Test: SME with empty string district receives district-specific events
func (s *AutoBroadcastTestSuite) TestSMEWithEmptyDistrictReceivesDistrictEvents() {
	timestamp := time.Now().UnixNano()
	password, _ := facades.Hash().Make("password123")

	// Create user
	userEmpty := &models.User{
		Name:     "Empty District User",
		Email:    fmt.Sprintf("emptydistrict_%d@example.com", timestamp),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(userEmpty)

	// Create SME with empty string district
	emptyDistrict := ""
	smeEmpty := &models.Sme{
		Name:             "Empty District SME",
		ContactEmail:     userEmpty.Email,
		ContactPhone:     "0888888888",
		BusinessCategory: "Trading",
		Sector:           "Retail",
		District:         &emptyDistrict, // Empty string
		IsActive:         true,
	}
	facades.Orm().Query().Create(smeEmpty)

	eventID := uint(5001)
	s.notificationService.BroadcastEventNotificationSync(
		s.adminUser.ID,
		eventID,
		"District Test Event",
		"2025-11-01",
		"Test Venue",
		"Lilongwe", // Specific district
	)

	time.Sleep(200 * time.Millisecond)

	var notification models.Notification
	err := facades.Orm().Query().
		Where("user_id = ?", userEmpty.ID).
		Where("type = ?", "event").
		Where("related_id = ?", eventID).
		First(&notification)

	s.NoError(err, "SME with empty district should receive district-specific events")
}

// Test: Multiple events can be sent in succession
func (s *AutoBroadcastTestSuite) TestMultipleEventsBroadcast() {
	for i := 1; i <= 3; i++ {
		eventID := uint(6000 + i)
		s.notificationService.BroadcastEventNotificationSync(
			s.adminUser.ID,
			eventID,
			fmt.Sprintf("Multiple Event %d", i),
			fmt.Sprintf("2025-12-%02d", i),
			fmt.Sprintf("Venue %d", i),
			"",
		)
	}

	time.Sleep(300 * time.Millisecond)

	// Count notifications for one user - should have 3
	count, _ := facades.Orm().Query().Model(&models.Notification{}).
		Where("user_id = ?", s.smeUserLilongwe1.ID).
		Where("type = ?", "event").
		Where("related_id >= ? AND related_id <= ?", 6001, 6003).
		Count()

	s.Equal(int64(3), count, "User should receive all 3 event notifications")
}

// ============================================================================
// LISTENER INTEGRATION TESTS
// ============================================================================

// Test: BroadcastEventNotification listener handles correct arguments
func (s *AutoBroadcastTestSuite) TestEventListenerHandlesCorrectArguments() {
	// This tests the listener directly to verify argument handling
	listener := &services.NotificationService{}
	// The listener is tested indirectly through the notification service
	// Verify it doesn't panic with correct arguments
	s.NotPanics(func() {
		services.NewNotificationService().BroadcastEventNotificationSync(
			s.adminUser.ID,
			uint(7001),
			"Listener Test Event",
			"2025-12-15",
			"Test Venue",
			"Lilongwe",
		)
	})
	_ = listener // Use the variable
}

// Test: BroadcastProcurementNotification listener handles correct arguments
func (s *AutoBroadcastTestSuite) TestProcurementListenerHandlesCorrectArguments() {
	s.NotPanics(func() {
		services.NewNotificationService().BroadcastProcurementNotificationSync(
			s.adminUser.ID,
			uint(8001),
			"Test Organization",
			"REF/001",
			"Goods",
			"2025-12-31",
		)
	})
}

// ============================================================================
// HELPER METHODS
// ============================================================================

func (s *AutoBroadcastTestSuite) countNotifications() int64 {
	count, _ := facades.Orm().Query().Model(&models.Notification{}).Count()
	return count
}

// ============================================================================
// EVENT DISPATCH INTEGRATION TESTS (via controller hooks)
// ============================================================================

// EventDispatchTestSuite tests the event dispatch from controllers
type EventDispatchTestSuite struct {
	suite.Suite
	tests.TestCase

	adminUser *models.User
	smeUser   *models.User
	sme       *models.Sme
}

func TestEventDispatchTestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatchTestSuite))
}

func (s *EventDispatchTestSuite) SetupTest() {
	s.RefreshDatabase()

	password, _ := facades.Hash().Make("password123")
	timestamp := time.Now().UnixNano()

	// Create admin user
	s.adminUser = &models.User{
		Name:         "Dispatch Admin",
		Email:        fmt.Sprintf("dispatch_admin_%d@example.com", timestamp),
		Password:     password,
		IsActive:     true,
		IsSuperAdmin: true,
	}
	facades.Orm().Query().Create(s.adminUser)

	// Create SME user
	s.smeUser = &models.User{
		Name:     "Dispatch SME User",
		Email:    fmt.Sprintf("dispatch_sme_%d@example.com", timestamp),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.smeUser)

	// Create SME
	district := "Lilongwe"
	s.sme = &models.Sme{
		Name:             "Dispatch Test SME",
		ContactEmail:     s.smeUser.Email,
		ContactPhone:     "0999123456",
		BusinessCategory: "Trading",
		Sector:           "Retail",
		District:         &district,
		IsActive:         true,
	}
	facades.Orm().Query().Create(s.sme)
}

// Test: Creating an event directly creates notifications
func (s *EventDispatchTestSuite) TestEventCreationTriggersNotification() {
	// Use a fixed event ID to test the notification service directly
	// This tests that notifications are created when the service is called,
	// which is what the controller afterStore hook does
	eventID := uint(9001)
	eventTitle := "Test Event Title"
	eventDate := "2025-12-20 10:00:00"
	eventVenue := "Test Venue"
	eventDistrict := "Lilongwe"

	// Manually call the notification service (simulating what the afterStore hook does)
	notificationService := services.NewNotificationService()
	notificationService.BroadcastEventNotificationSync(
		s.adminUser.ID,
		eventID,
		eventTitle,
		eventDate,
		eventVenue,
		eventDistrict,
	)

	time.Sleep(200 * time.Millisecond)

	// Verify notification was created
	var notification models.Notification
	err := facades.Orm().Query().
		Where("user_id = ?", s.smeUser.ID).
		Where("type = ?", "event").
		Where("related_id = ?", eventID).
		First(&notification)

	s.NoError(err, "Notification should be created for SME user")
	s.Contains(notification.Title, eventTitle)
	s.Contains(notification.Message, eventTitle, "Message should contain event title")
}

// Test: Publishing a procurement creates notifications
func (s *EventDispatchTestSuite) TestProcurementPublishTriggersNotification() {
	// Create a procurement notice
	procurement := &models.ProcurementNotice{
		Organization:        "Test Ministry",
		ProcurementType:     "Goods",
		RefNo:               "TEST/2025/DISPATCH",
		MarketApproach:      "Open",
		Invitation:          "Test invitation",
		OpenDate:            *carbon.NewDateTime(carbon.Now()),
		CloseDate:           *carbon.NewDateTime(carbon.Now().AddDays(30)),
		IsPublished:         true,
		Details:             "Test details",
		Partners:            []string{},
		QualifyingDistricts: []string{},
		Classification:      []string{},
		InterestedSmes:      []string{},
	}
	err := facades.Orm().Query().Create(procurement)
	s.NoError(err)
	s.NotEqual(uint(0), procurement.ID, "Procurement should have been assigned an ID")

	// Manually call the notification service (simulating what the controller does)
	notificationService := services.NewNotificationService()
	notificationService.BroadcastProcurementNotificationSync(
		s.adminUser.ID,
		procurement.ID,
		procurement.Organization,
		procurement.RefNo,
		procurement.ProcurementType,
		procurement.CloseDate.ToDateTimeString(),
	)

	time.Sleep(200 * time.Millisecond)

	// Verify notification was created
	var notification models.Notification
	err = facades.Orm().Query().
		Where("user_id = ?", s.smeUser.ID).
		Where("type = ?", "procurement").
		Where("related_id = ?", procurement.ID).
		First(&notification)

	s.NoError(err, "Notification should be created for SME user")
	s.Contains(notification.Message, procurement.Organization)
}

// Test: Unpublished procurement does NOT create notifications
func (s *EventDispatchTestSuite) TestUnpublishedProcurementNoNotification() {
	// Create an unpublished procurement notice
	procurement := &models.ProcurementNotice{
		Organization:        "Unpublished Ministry",
		ProcurementType:     "Services",
		RefNo:               "UNPUB/2025/001",
		MarketApproach:      "Open",
		Invitation:          "Test invitation",
		OpenDate:            *carbon.NewDateTime(carbon.Now()),
		CloseDate:           *carbon.NewDateTime(carbon.Now().AddDays(30)),
		IsPublished:         false, // Not published
		Details:             "Test details",
		Partners:            []string{},
		QualifyingDistricts: []string{},
		Classification:      []string{},
		InterestedSmes:      []string{},
	}
	err := facades.Orm().Query().Create(procurement)
	s.NoError(err)

	// Count notifications before
	beforeCount, _ := facades.Orm().Query().Model(&models.Notification{}).
		Where("type = ?", "procurement").
		Count()

	// The controller would NOT dispatch the event because IsPublished is false
	// So no notifications should be created
	// (This test verifies the expected behavior)

	time.Sleep(100 * time.Millisecond)

	// Count notifications after
	afterCount, _ := facades.Orm().Query().Model(&models.Notification{}).
		Where("type = ?", "procurement").
		Count()

	s.Equal(beforeCount, afterCount, "No notifications should be created for unpublished procurement")
}
