package messages

import (
	"fmt"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"starter-project/app/contracts"
	"starter-project/app/models"
	"starter-project/app/services"
	"starter-project/tests"
)

type BroadcastTestSuite struct {
	suite.Suite
	tests.TestCase

	messageService *services.MessageService

	// Users
	superAdmin *models.User
	admin      *models.User
	member1    *models.User
	member2    *models.User
	member3    *models.User
	guest      *models.User

	// Roles
	adminRole  *models.Role
	memberRole *models.Role
	guestRole  *models.Role
}

func TestBroadcastTestSuite(t *testing.T) {
	suite.Run(t, new(BroadcastTestSuite))
}

func (s *BroadcastTestSuite) SetupTest() {
	s.RefreshDatabase()
	s.messageService = services.NewMessageService()

	password, _ := facades.Hash().Make("password123")

	// Create roles
	s.adminRole = &models.Role{
		Name:        "Administrator",
		Slug:        fmt.Sprintf("admin_%d", time.Now().UnixNano()),
		Description: "Administrator role",
		Level:       80,
		IsActive:    true,
	}
	facades.Orm().Query().Create(s.adminRole)

	s.memberRole = &models.Role{
		Name:        "Member",
		Slug:        fmt.Sprintf("member_%d", time.Now().UnixNano()),
		Description: "Member role",
		Level:       20,
		IsActive:    true,
	}
	facades.Orm().Query().Create(s.memberRole)

	s.guestRole = &models.Role{
		Name:        "Guest",
		Slug:        fmt.Sprintf("guest_%d", time.Now().UnixNano()),
		Description: "Guest role",
		Level:       10,
		IsActive:    true,
	}
	facades.Orm().Query().Create(s.guestRole)

	// Create super admin
	s.superAdmin = &models.User{
		Name:         "Super Admin",
		Email:        fmt.Sprintf("superadmin_%d@example.com", time.Now().UnixNano()),
		Password:     password,
		IsActive:     true,
		IsSuperAdmin: true,
	}
	facades.Orm().Query().Create(s.superAdmin)

	// Create admin
	s.admin = &models.User{
		Name:     "Admin User",
		Email:    fmt.Sprintf("admin_%d@example.com", time.Now().UnixNano()),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.admin)
	s.assignRole(s.admin, s.adminRole)

	// Create multiple members
	s.member1 = &models.User{
		Name:     "Member One",
		Email:    fmt.Sprintf("member1_%d@example.com", time.Now().UnixNano()),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.member1)
	s.assignRole(s.member1, s.memberRole)

	s.member2 = &models.User{
		Name:     "Member Two",
		Email:    fmt.Sprintf("member2_%d@example.com", time.Now().UnixNano()),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.member2)
	s.assignRole(s.member2, s.memberRole)

	s.member3 = &models.User{
		Name:     "Member Three",
		Email:    fmt.Sprintf("member3_%d@example.com", time.Now().UnixNano()),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.member3)
	s.assignRole(s.member3, s.memberRole)

	// Create guest
	s.guest = &models.User{
		Name:     "Guest User",
		Email:    fmt.Sprintf("guest_%d@example.com", time.Now().UnixNano()),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.guest)
	s.assignRole(s.guest, s.guestRole)
}

func (s *BroadcastTestSuite) assignRole(user *models.User, role *models.Role) {
	facades.Orm().Query().Exec(
		"INSERT INTO user_roles (user_id, role_id, is_active, assigned_at, notes, created_at, updated_at) VALUES (?, ?, true, NOW(), '', NOW(), NOW())",
		user.ID, role.ID,
	)
}

// Test: BroadcastToRole sends messages to all users with the specified role
func (s *BroadcastTestSuite) TestBroadcastToRole_SendsToAllUsersInRole() {
	content := "Test broadcast to all members"

	result, err := s.messageService.BroadcastToRole(
		s.superAdmin.ID,
		s.memberRole.ID,
		content,
		"",
	)

	s.NoError(err, "BroadcastToRole should not return an error")
	s.NotNil(result)

	// Should have sent to 3 members
	sentCount := result["sent_count"].(int)
	s.Equal(3, sentCount, "Should send to all 3 members")

	// Verify messages were created
	var messages []models.Message
	facades.Orm().Query().
		Where("sender_id = ?", s.superAdmin.ID).
		Where("type = ?", string(models.MessageTypeSystem)).
		Where("content = ?", content).
		Find(&messages)

	s.Len(messages, 3, "Should create 3 messages")

	// Verify recipients
	recipientIDs := make(map[uint]bool)
	for _, msg := range messages {
		recipientIDs[*msg.RecipientID] = true
	}
	s.True(recipientIDs[s.member1.ID], "Member 1 should receive message")
	s.True(recipientIDs[s.member2.ID], "Member 2 should receive message")
	s.True(recipientIDs[s.member3.ID], "Member 3 should receive message")
}

// Test: BroadcastToRole does not send to the sender even if they have the role
func (s *BroadcastTestSuite) TestBroadcastToRole_ExcludesSender() {
	// Assign memberRole to superAdmin as well
	s.assignRole(s.superAdmin, s.memberRole)

	content := "Test broadcast excluding sender"

	result, err := s.messageService.BroadcastToRole(
		s.superAdmin.ID,
		s.memberRole.ID,
		content,
		"",
	)

	s.NoError(err)
	sentCount := result["sent_count"].(int)
	s.Equal(3, sentCount, "Should still only send to 3 members, excluding the sender")

	// Verify super admin didn't receive the message
	var selfMessage models.Message
	facades.Orm().Query().
		Where("sender_id = ?", s.superAdmin.ID).
		Where("recipient_id = ?", s.superAdmin.ID).
		Where("content = ?", content).
		First(&selfMessage)

	s.Equal(uint(0), selfMessage.ID, "Super admin should not receive their own broadcast")
}

// Test: BroadcastToRole returns correct role name
func (s *BroadcastTestSuite) TestBroadcastToRole_ReturnsRoleName() {
	result, err := s.messageService.BroadcastToRole(
		s.superAdmin.ID,
		s.memberRole.ID,
		"Test message",
		"",
	)

	s.NoError(err)
	s.Equal("Member", result["role_name"], "Should return correct role name")
	s.Equal(s.memberRole.ID, result["role_id"], "Should return correct role ID")
}

// Test: BroadcastToRole returns 0 when no users have the role
func (s *BroadcastTestSuite) TestBroadcastToRole_NoUsersInRole() {
	// Create a new role with no users
	emptyRole := &models.Role{
		Name:        "Empty Role",
		Slug:        fmt.Sprintf("empty_%d", time.Now().UnixNano()),
		Description: "Role with no users",
		Level:       50,
		IsActive:    true,
	}
	facades.Orm().Query().Create(emptyRole)

	result, err := s.messageService.BroadcastToRole(
		s.superAdmin.ID,
		emptyRole.ID,
		"Test message to empty role",
		"",
	)

	s.NoError(err, "Should not error even with no users")
	s.Equal(0, result["sent_count"], "Should report 0 messages sent")
}

// Test: BroadcastToRole fails for non-existent role
func (s *BroadcastTestSuite) TestBroadcastToRole_NonExistentRole() {
	// Use a very large role ID that definitely doesn't exist
	nonExistentRoleID := uint(999999999)

	result, err := s.messageService.BroadcastToRole(
		s.superAdmin.ID,
		nonExistentRoleID,
		"Test message",
		"",
	)

	// The service either returns an error OR returns result with sent_count=0
	// Check that either an error occurred OR no messages were sent
	if err != nil {
		s.Contains(err.Error(), "role not found")
		s.Nil(result)
	} else {
		// If no error, the result should indicate no users/role found
		s.NotNil(result)
		sentCount := result["sent_count"].(int)
		s.Equal(0, sentCount, "Should send 0 messages for non-existent role")
	}
}

// Test: BroadcastToRole creates notifications for all recipients
func (s *BroadcastTestSuite) TestBroadcastToRole_CreatesNotifications() {
	content := "Notification test broadcast"

	_, err := s.messageService.BroadcastToRole(
		s.superAdmin.ID,
		s.memberRole.ID,
		content,
		"",
	)
	s.NoError(err)

	// Give a moment for notifications to be created
	time.Sleep(100 * time.Millisecond)

	// Check notifications for each member
	for _, member := range []*models.User{s.member1, s.member2, s.member3} {
		var notification models.Notification
		err := facades.Orm().Query().
			Where("user_id = ?", member.ID).
			Where("type = ?", "message").
			Where("trigger_user_id = ?", s.superAdmin.ID).
			First(&notification)

		s.NoError(err, fmt.Sprintf("Notification should be created for %s", member.Name))
		s.Contains(notification.Title, s.superAdmin.Name, "Notification should mention sender")
	}
}

// Test: BroadcastToRole only sends to active users
func (s *BroadcastTestSuite) TestBroadcastToRole_OnlyActiveUsers() {
	// Deactivate member2
	facades.Orm().Query().Model(s.member2).Update("is_active", false)

	content := "Test broadcast to active users only"

	result, err := s.messageService.BroadcastToRole(
		s.superAdmin.ID,
		s.memberRole.ID,
		content,
		"",
	)

	s.NoError(err)
	sentCount := result["sent_count"].(int)
	s.Equal(2, sentCount, "Should only send to 2 active members")

	// Verify member2 didn't receive message
	var member2Message models.Message
	facades.Orm().Query().
		Where("recipient_id = ?", s.member2.ID).
		Where("content = ?", content).
		First(&member2Message)

	s.Equal(uint(0), member2Message.ID, "Inactive user should not receive broadcast")
}

// Test: GetBroadcastHistory returns broadcasts sent by user
func (s *BroadcastTestSuite) TestGetBroadcastHistory_ReturnsBroadcasts() {
	// Send a broadcast
	content1 := "First broadcast message for history test"
	_, err := s.messageService.BroadcastToRole(
		s.superAdmin.ID,
		s.memberRole.ID,
		content1,
		"",
	)
	s.NoError(err)

	// Wait a moment and send another
	time.Sleep(100 * time.Millisecond)

	content2 := "Second broadcast message for history test"
	_, err = s.messageService.BroadcastToRole(
		s.superAdmin.ID,
		s.guestRole.ID,
		content2,
		"",
	)
	s.NoError(err)

	// Get broadcast history
	req := contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	}

	result, err := s.messageService.GetBroadcastHistory(s.superAdmin.ID, req)
	s.NoError(err)
	s.NotNil(result)

	// Should have at least 2 distinct broadcasts
	s.GreaterOrEqual(len(result.Data), 2, "Should return at least 2 broadcasts")

	// Verify both broadcasts are present in the history
	broadcasts := result.Data
	foundContent1 := false
	foundContent2 := false
	for _, item := range broadcasts {
		broadcast := item.(map[string]interface{})
		if broadcast["content"] == content1 {
			foundContent1 = true
		}
		if broadcast["content"] == content2 {
			foundContent2 = true
		}
	}
	s.True(foundContent1, "First broadcast should be in history")
	s.True(foundContent2, "Second broadcast should be in history")
}

// Test: GetBroadcastHistory returns correct recipient count
func (s *BroadcastTestSuite) TestGetBroadcastHistory_CorrectRecipientCount() {
	content := "Broadcast with recipient count"

	_, err := s.messageService.BroadcastToRole(
		s.superAdmin.ID,
		s.memberRole.ID,
		content,
		"",
	)
	s.NoError(err)

	req := contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	}

	result, err := s.messageService.GetBroadcastHistory(s.superAdmin.ID, req)
	s.NoError(err)

	// Find the broadcast we just sent
	found := false
	for _, item := range result.Data {
		broadcast := item.(map[string]interface{})
		if broadcast["content"] == content {
			found = true
			recipientCount := broadcast["recipient_count"].(int64)
			s.Equal(int64(3), recipientCount, "Should show 3 recipients")
			break
		}
	}
	s.True(found, "Should find the broadcast in history")
}

// Test: GetBroadcastHistory tracks read count
func (s *BroadcastTestSuite) TestGetBroadcastHistory_TracksReadCount() {
	content := "Broadcast for read tracking"

	_, err := s.messageService.BroadcastToRole(
		s.superAdmin.ID,
		s.memberRole.ID,
		content,
		"",
	)
	s.NoError(err)

	// Mark one message as read
	var message models.Message
	facades.Orm().Query().
		Where("recipient_id = ?", s.member1.ID).
		Where("content = ?", content).
		First(&message)

	err = s.messageService.MarkAsRead(message.ID, s.member1.ID)
	s.NoError(err)

	// Get broadcast history
	req := contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	}

	result, err := s.messageService.GetBroadcastHistory(s.superAdmin.ID, req)
	s.NoError(err)

	// Find the broadcast and check read count
	for _, item := range result.Data {
		broadcast := item.(map[string]interface{})
		if broadcast["content"] == content {
			readCount := broadcast["read_count"].(int64)
			s.Equal(int64(1), readCount, "Should show 1 read")
			break
		}
	}
}

// Test: GetBroadcastHistory respects pagination
func (s *BroadcastTestSuite) TestGetBroadcastHistory_Pagination() {
	// Send multiple broadcasts
	for i := 0; i < 5; i++ {
		content := fmt.Sprintf("Pagination test broadcast %d", i)
		_, err := s.messageService.BroadcastToRole(
			s.superAdmin.ID,
			s.memberRole.ID,
			content,
			"",
		)
		s.NoError(err)
		time.Sleep(50 * time.Millisecond) // Small delay to ensure distinct timestamps
	}

	// Get first page with page size of 2
	req := contracts.ListRequest{
		Page:     1,
		PageSize: 2,
	}

	result, err := s.messageService.GetBroadcastHistory(s.superAdmin.ID, req)
	s.NoError(err)

	s.Len(result.Data, 2, "First page should have 2 items")
	s.GreaterOrEqual(result.Total, int64(5), "Total should be at least 5")

	// Get second page
	req.Page = 2
	result2, err := s.messageService.GetBroadcastHistory(s.superAdmin.ID, req)
	s.NoError(err)

	s.Len(result2.Data, 2, "Second page should have 2 items")

	// Verify different content on different pages
	page1Content := result.Data[0].(map[string]interface{})["content"]
	page2Content := result2.Data[0].(map[string]interface{})["content"]
	s.NotEqual(page1Content, page2Content, "Pages should have different content")
}

// Test: GetBroadcastHistory returns empty for non-super-admin (no broadcasts)
func (s *BroadcastTestSuite) TestGetBroadcastHistory_EmptyForRegularUser() {
	req := contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	}

	// Regular user has no broadcasts
	result, err := s.messageService.GetBroadcastHistory(s.member1.ID, req)
	s.NoError(err)
	s.Len(result.Data, 0, "Regular user should have no broadcast history")
}
