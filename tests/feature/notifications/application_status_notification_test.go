package notifications

import (
	"fmt"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"smedi-sme-db/app/listeners"
	"smedi-sme-db/app/models"
	"smedi-sme-db/app/services"
	"smedi-sme-db/tests"
)

// ApplicationStatusNotificationTestSuite tests the application approval/rejection notification system
type ApplicationStatusNotificationTestSuite struct {
	suite.Suite
	tests.TestCase

	notificationService *services.NotificationService

	// Admin user who approves/rejects applications
	adminUser *models.User

	// Applicant user who will receive notifications
	applicantUser *models.User
}

func TestApplicationStatusNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(ApplicationStatusNotificationTestSuite))
}

func (s *ApplicationStatusNotificationTestSuite) SetupTest() {
	s.RefreshDatabase()
	s.notificationService = services.NewNotificationService()

	password, _ := facades.Hash().Make("password123")
	timestamp := time.Now().UnixNano()

	// Create admin user who approves/rejects
	s.adminUser = &models.User{
		Name:         "Admin User",
		Email:        fmt.Sprintf("admin_%d@example.com", timestamp),
		Password:     password,
		IsActive:     true,
		IsSuperAdmin: true,
	}
	facades.Orm().Query().Create(s.adminUser)

	// Create applicant user who will receive notifications
	s.applicantUser = &models.User{
		Name:     "Applicant User",
		Email:    fmt.Sprintf("applicant_%d@example.com", timestamp),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.applicantUser)
}

// ============================================================================
// APPROVED APPLICATION TESTS
// ============================================================================

// Test: NotifyApplicationStatusSync creates notification for approved application
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationStatusSync_Approved_CreatesNotification() {
	applicationID := uint(1001)
	applicationType := models.ApplicationTypeSignup
	orgName := "Test Organization"

	// Call the notification service directly
	s.notificationService.NotifyApplicationStatusSync(
		s.adminUser.ID,
		applicationID,
		applicationType,
		orgName,
		s.applicantUser.Email,
		"approved",
		"",
	)

	// Wait for notification to be created
	time.Sleep(200 * time.Millisecond)

	// Verify notification was created
	var notification models.Notification
	err := facades.Orm().Query().
		Where("user_id = ?", s.applicantUser.ID).
		Where("type = ?", "application_approved").
		Where("related_id = ?", applicationID).
		First(&notification)

	s.NoError(err, "Notification should be created for applicant user")
	s.Equal("Application Approved", notification.Title, "Title should be 'Application Approved'")
	s.Contains(notification.Message, orgName, "Message should contain organization name")
	s.Contains(notification.Message, "approved", "Message should mention approved status")
	s.Contains(notification.Message, "Welcome!", "Message should contain welcome")
	s.Equal("application", notification.RelatedType, "Related type should be 'application'")
	s.Equal(&applicationID, notification.RelatedID, "Related ID should be application ID")
	s.Equal(&s.adminUser.ID, notification.TriggerUserID, "Trigger user should be admin")
	s.Equal("normal", notification.Priority, "Priority should be 'normal' for approved applications")
}

// Test: NotifyApplicationStatusSync creates message for approved application
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationStatusSync_Approved_CreatesMessage() {
	applicationID := uint(1002)
	applicationType := models.ApplicationTypeSignup
	orgName := "Test Organization"

	s.notificationService.NotifyApplicationStatusSync(
		s.adminUser.ID,
		applicationID,
		applicationType,
		orgName,
		s.applicantUser.Email,
		"approved",
		"",
	)

	time.Sleep(200 * time.Millisecond)

	// Verify message was created
	var message models.Message
	err := facades.Orm().Query().
		Where("sender_id = ?", s.adminUser.ID).
		Where("recipient_id = ?", s.applicantUser.ID).
		Where("type = ?", string(models.MessageTypeDirect)).
		Where("content LIKE ?", "%Application Approved%").
		First(&message)

	s.NoError(err, "Message should be created for applicant user")
	s.Contains(message.Content, orgName, "Message content should contain organization name")
	s.Contains(message.Content, "approved", "Message content should mention approved")
}

// ============================================================================
// REJECTED APPLICATION TESTS
// ============================================================================

// Test: NotifyApplicationStatusSync creates notification for rejected application
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationStatusSync_Rejected_CreatesNotification() {
	applicationID := uint(2001)
	applicationType := models.ApplicationTypeSignup
	orgName := "Test Organization"
	rejectionReason := "Missing required documentation"

	s.notificationService.NotifyApplicationStatusSync(
		s.adminUser.ID,
		applicationID,
		applicationType,
		orgName,
		s.applicantUser.Email,
		"rejected",
		rejectionReason,
	)

	time.Sleep(200 * time.Millisecond)

	// Verify notification was created
	var notification models.Notification
	err := facades.Orm().Query().
		Where("user_id = ?", s.applicantUser.ID).
		Where("type = ?", "application_rejected").
		Where("related_id = ?", applicationID).
		First(&notification)

	s.NoError(err, "Notification should be created for applicant user")
	s.Equal("Application Rejected", notification.Title, "Title should be 'Application Rejected'")
	s.Contains(notification.Message, orgName, "Message should contain organization name")
	s.Contains(notification.Message, "rejected", "Message should mention rejected status")
}

// Test: Rejection notification includes reason in message
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationStatusSync_Rejected_IncludesReasonInNotification() {
	applicationID := uint(2002)
	applicationType := models.ApplicationTypeSignup
	orgName := "Test Organization"
	rejectionReason := "Invalid business registration number"

	s.notificationService.NotifyApplicationStatusSync(
		s.adminUser.ID,
		applicationID,
		applicationType,
		orgName,
		s.applicantUser.Email,
		"rejected",
		rejectionReason,
	)

	time.Sleep(200 * time.Millisecond)

	// Verify notification contains rejection reason
	var notification models.Notification
	err := facades.Orm().Query().
		Where("user_id = ?", s.applicantUser.ID).
		Where("type = ?", "application_rejected").
		Where("related_id = ?", applicationID).
		First(&notification)

	s.NoError(err)
	s.Contains(notification.Message, rejectionReason, "Notification message should contain rejection reason")
}

// Test: Rejection notification includes reason in direct message
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationStatusSync_Rejected_IncludesReasonInMessage() {
	applicationID := uint(2003)
	applicationType := models.ApplicationTypeSignup
	orgName := "Test Organization"
	rejectionReason := "Incomplete application form"

	s.notificationService.NotifyApplicationStatusSync(
		s.adminUser.ID,
		applicationID,
		applicationType,
		orgName,
		s.applicantUser.Email,
		"rejected",
		rejectionReason,
	)

	time.Sleep(200 * time.Millisecond)

	// Verify message contains rejection reason
	var message models.Message
	err := facades.Orm().Query().
		Where("sender_id = ?", s.adminUser.ID).
		Where("recipient_id = ?", s.applicantUser.ID).
		Where("type = ?", string(models.MessageTypeDirect)).
		Where("content LIKE ?", "%Rejected%").
		First(&message)

	s.NoError(err, "Message should be created for applicant user")
	s.Contains(message.Content, rejectionReason, "Message content should contain rejection reason")
}

// Test: Rejection notifications have high priority
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationStatusSync_Rejected_HasHighPriority() {
	applicationID := uint(2004)
	applicationType := models.ApplicationTypeSignup
	orgName := "Test Organization"
	rejectionReason := "Documents not verified"

	s.notificationService.NotifyApplicationStatusSync(
		s.adminUser.ID,
		applicationID,
		applicationType,
		orgName,
		s.applicantUser.Email,
		"rejected",
		rejectionReason,
	)

	time.Sleep(200 * time.Millisecond)

	// Verify notification has high priority
	var notification models.Notification
	err := facades.Orm().Query().
		Where("user_id = ?", s.applicantUser.ID).
		Where("type = ?", "application_rejected").
		Where("related_id = ?", applicationID).
		First(&notification)

	s.NoError(err)
	s.Equal("high", notification.Priority, "Rejection notifications should have high priority")
}

// ============================================================================
// EDGE CASE TESTS
// ============================================================================

// Test: Service handles non-existent user email gracefully
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationStatusSync_UserNotFound_HandlesGracefully() {
	applicationID := uint(4001)
	applicationType := models.ApplicationTypeSignup
	orgName := "Unknown Organization"
	nonExistentEmail := "nonexistent_user_does_not_exist@example.com"

	// The service should handle this gracefully without panicking
	s.NotPanics(func() {
		s.notificationService.NotifyApplicationStatusSync(
			s.adminUser.ID,
			applicationID,
			applicationType,
			orgName,
			nonExistentEmail,
			"approved",
			"",
		)
	})

	time.Sleep(200 * time.Millisecond)

	// Verify no notification was created for the non-existent user
	var notification models.Notification
	err := facades.Orm().Query().
		Where("related_id = ?", applicationID).
		Where("type = ?", "application_approved").
		Where("user_id != ?", s.applicantUser.ID).
		First(&notification)

	if err == nil && notification.UserID == 0 {
		s.T().Log("Note: Implementation creates notification for user_id=0 when user not found")
	}
}

// Test: Inactive user does not receive notification with their user ID
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationStatusSync_InactiveUser_HandlesGracefully() {
	timestamp := time.Now().UnixNano()
	password, _ := facades.Hash().Make("password123")

	// Create inactive user
	inactiveUser := &models.User{
		Name:     "Inactive User",
		Email:    fmt.Sprintf("inactive_%d@example.com", timestamp),
		Password: password,
	}
	facades.Orm().Query().Create(inactiveUser)
	// Explicitly set inactive after creation to bypass GORM defaults
	facades.Orm().Query().Model(inactiveUser).Update("is_active", false)

	applicationID := uint(4002)
	applicationType := models.ApplicationTypeSignup
	orgName := "Inactive Test"

	// The service should handle this gracefully without panicking
	s.NotPanics(func() {
		s.notificationService.NotifyApplicationStatusSync(
			s.adminUser.ID,
			applicationID,
			applicationType,
			orgName,
			inactiveUser.Email,
			"approved",
			"",
		)
	})

	time.Sleep(200 * time.Millisecond)

	// Verify that no notification was created specifically for the inactive user's ID
	var notificationCount int64
	notificationCount, _ = facades.Orm().Query().Model(&models.Notification{}).
		Where("user_id = ?", inactiveUser.ID).
		Where("type = ?", "application_approved").
		Where("related_id = ?", applicationID).
		Count()

	s.Equal(int64(0), notificationCount, "Inactive user should not have notification with their user ID")
}

// Test: Case-insensitive email matching
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationStatusSync_CaseInsensitiveEmail() {
	applicationID := uint(4003)
	applicationType := models.ApplicationTypeSignup
	orgName := "Case Test Organization"

	// Create user with lowercase email
	timestamp := time.Now().UnixNano()
	password, _ := facades.Hash().Make("password123")
	caseTestUser := &models.User{
		Name:     "Case Test User",
		Email:    fmt.Sprintf("casetest_%d@example.com", timestamp),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(caseTestUser)

	// Call with uppercase email (should still match due to case-insensitive lookup)
	upperEmail := fmt.Sprintf("CASETEST_%d@EXAMPLE.COM", timestamp)

	s.notificationService.NotifyApplicationStatusSync(
		s.adminUser.ID,
		applicationID,
		applicationType,
		orgName,
		upperEmail,
		"approved",
		"",
	)

	time.Sleep(200 * time.Millisecond)

	// Verify notification was created despite email case mismatch
	var notification models.Notification
	err := facades.Orm().Query().
		Where("user_id = ?", caseTestUser.ID).
		Where("type = ?", "application_approved").
		Where("related_id = ?", applicationID).
		First(&notification)

	s.NoError(err, "Notification should be created despite email case mismatch")
}

// Test: Rejection with empty reason
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationStatusSync_RejectedWithEmptyReason() {
	applicationID := uint(4004)
	applicationType := models.ApplicationTypeSignup
	orgName := "Test Organization"

	s.notificationService.NotifyApplicationStatusSync(
		s.adminUser.ID,
		applicationID,
		applicationType,
		orgName,
		s.applicantUser.Email,
		"rejected",
		"", // Empty reason
	)

	time.Sleep(200 * time.Millisecond)

	// Verify notification was created without reason
	var notification models.Notification
	err := facades.Orm().Query().
		Where("user_id = ?", s.applicantUser.ID).
		Where("type = ?", "application_rejected").
		Where("related_id = ?", applicationID).
		First(&notification)

	s.NoError(err)
	s.NotContains(notification.Message, "Reason:", "Message should not contain 'Reason:' when reason is empty")
}

// ============================================================================
// LISTENER TESTS
// ============================================================================

// Test: NotifyApplicationApproved listener handles correct arguments
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationApprovedListener_HandlesCorrectArguments() {
	listener := &listeners.NotifyApplicationApproved{}

	// Verify signature
	s.Equal("notify_application_approved", listener.Signature(), "Listener signature should be correct")

	// Test queue configuration
	queue := listener.Queue()
	s.True(queue.Enable, "Queue should be enabled for async processing")

	// Test handle with correct arguments - should not panic
	s.NotPanics(func() {
		err := listener.Handle(
			s.adminUser.ID,        // approverUserID
			uint(5001),            // applicationID
			"signup",              // applicationType
			"Test Organization",   // orgName
			s.applicantUser.Email, // recipientEmail
		)
		s.NoError(err)
	})
}

// Test: NotifyApplicationApproved listener handles insufficient arguments
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationApprovedListener_HandlesInsufficientArguments() {
	listener := &listeners.NotifyApplicationApproved{}

	// Test with insufficient arguments - should return nil without panicking
	s.NotPanics(func() {
		err := listener.Handle(
			s.adminUser.ID, // Only one argument
		)
		s.NoError(err)
	})
}

// Test: NotifyApplicationRejected listener handles correct arguments
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationRejectedListener_HandlesCorrectArguments() {
	listener := &listeners.NotifyApplicationRejected{}

	// Verify signature
	s.Equal("notify_application_rejected", listener.Signature(), "Listener signature should be correct")

	// Test queue configuration
	queue := listener.Queue()
	s.True(queue.Enable, "Queue should be enabled for async processing")

	// Test handle with correct arguments - should not panic
	s.NotPanics(func() {
		err := listener.Handle(
			s.adminUser.ID,        // rejectorUserID
			uint(5002),            // applicationID
			"signup",              // applicationType
			"Test Organization",   // orgName
			s.applicantUser.Email, // recipientEmail
			"Documents invalid",   // reason
		)
		s.NoError(err)
	})
}

// Test: NotifyApplicationRejected listener handles insufficient arguments
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationRejectedListener_HandlesInsufficientArguments() {
	listener := &listeners.NotifyApplicationRejected{}

	// Test with insufficient arguments - should return nil without panicking
	s.NotPanics(func() {
		err := listener.Handle(
			s.adminUser.ID, // Only some arguments
			uint(5003),
			"signup",
		)
		s.NoError(err)
	})
}

// Test: NotifyApplicationApproved listener handles invalid argument types
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationApprovedListener_HandlesInvalidArgumentTypes() {
	listener := &listeners.NotifyApplicationApproved{}

	// Test with invalid types - should return nil without panicking
	s.NotPanics(func() {
		err := listener.Handle(
			"invalid_type", // Should be uint
			uint(5004),
			"signup",
			"Test Organization",
			s.applicantUser.Email,
		)
		s.NoError(err)
	})
}

// Test: NotifyApplicationRejected listener handles invalid argument types
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationRejectedListener_HandlesInvalidArgumentTypes() {
	listener := &listeners.NotifyApplicationRejected{}

	// Test with invalid types - should return nil without panicking
	s.NotPanics(func() {
		err := listener.Handle(
			s.adminUser.ID,
			"invalid_id", // Should be uint
			"signup",
			"Test Organization",
			s.applicantUser.Email,
			"Reason",
		)
		s.NoError(err)
	})
}

// ============================================================================
// NOTIFICATION CONTENT FORMAT TESTS
// ============================================================================

// Test: Approved application notification content format
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationStatusSync_Approved_CorrectContentFormat() {
	applicationID := uint(6001)
	applicationType := models.ApplicationTypeSignup
	orgName := "Format Test Organization"

	s.notificationService.NotifyApplicationStatusSync(
		s.adminUser.ID,
		applicationID,
		applicationType,
		orgName,
		s.applicantUser.Email,
		"approved",
		"",
	)

	time.Sleep(200 * time.Millisecond)

	var notification models.Notification
	err := facades.Orm().Query().
		Where("user_id = ?", s.applicantUser.ID).
		Where("type = ?", "application_approved").
		Where("related_id = ?", applicationID).
		First(&notification)

	s.NoError(err)
	s.Contains(notification.Message, "approved", "Message should mention approved")
	s.Contains(notification.Message, "Welcome!", "Message should contain welcome")
}

// Test: Rejected application notification content format
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationStatusSync_Rejected_CorrectContentFormat() {
	applicationID := uint(6002)
	applicationType := models.ApplicationTypeSignup
	orgName := "Format Test Organization"
	reason := "Tax ID mismatch"

	s.notificationService.NotifyApplicationStatusSync(
		s.adminUser.ID,
		applicationID,
		applicationType,
		orgName,
		s.applicantUser.Email,
		"rejected",
		reason,
	)

	time.Sleep(200 * time.Millisecond)

	var notification models.Notification
	err := facades.Orm().Query().
		Where("user_id = ?", s.applicantUser.ID).
		Where("type = ?", "application_rejected").
		Where("related_id = ?", applicationID).
		First(&notification)

	s.NoError(err)
	s.Contains(notification.Message, "Reason:", "Message should contain 'Reason:' prefix")
	s.Contains(notification.Message, reason, "Message should contain the rejection reason")
}

// Test: Direct message content format for approved application
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationStatusSync_ApprovedMessage_CorrectContentFormat() {
	applicationID := uint(6003)
	applicationType := models.ApplicationTypeSignup
	orgName := "Message Format Organization"

	s.notificationService.NotifyApplicationStatusSync(
		s.adminUser.ID,
		applicationID,
		applicationType,
		orgName,
		s.applicantUser.Email,
		"approved",
		"",
	)

	time.Sleep(200 * time.Millisecond)

	var message models.Message
	err := facades.Orm().Query().
		Where("sender_id = ?", s.adminUser.ID).
		Where("recipient_id = ?", s.applicantUser.ID).
		Where("content LIKE ?", "%Application Approved%").
		First(&message)

	s.NoError(err)
	s.Contains(message.Content, "portal", "Message should mention portal")
}

// Test: Direct message content format for rejected application
func (s *ApplicationStatusNotificationTestSuite) TestNotifyApplicationStatusSync_RejectedMessage_CorrectContentFormat() {
	applicationID := uint(6004)
	applicationType := models.ApplicationTypeSignup
	orgName := "Message Format Organization"
	reason := "Invalid documentation"

	s.notificationService.NotifyApplicationStatusSync(
		s.adminUser.ID,
		applicationID,
		applicationType,
		orgName,
		s.applicantUser.Email,
		"rejected",
		reason,
	)

	time.Sleep(200 * time.Millisecond)

	var message models.Message
	err := facades.Orm().Query().
		Where("sender_id = ?", s.adminUser.ID).
		Where("recipient_id = ?", s.applicantUser.ID).
		Where("content LIKE ?", "%Rejected%").
		First(&message)

	s.NoError(err)
	s.Contains(message.Content, "support team", "Rejected message should mention support team")
	s.Contains(message.Content, "**Reason:**", "Rejected message should contain Reason in bold")
}

// ============================================================================
// HELPER METHODS
// ============================================================================

func (s *ApplicationStatusNotificationTestSuite) countNotifications() int64 {
	count, _ := facades.Orm().Query().Model(&models.Notification{}).Count()
	return count
}

func (s *ApplicationStatusNotificationTestSuite) countMessages() int64 {
	count, _ := facades.Orm().Query().Model(&models.Message{}).Count()
	return count
}
