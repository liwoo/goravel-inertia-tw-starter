package messages

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

// uniqueEmail generates a unique email address using timestamp to avoid conflicts
func uniqueEmail(prefix string) string {
	return fmt.Sprintf("%s_%d@example.com", prefix, time.Now().UnixNano())
}

type MessagePermissionTestSuite struct {
	suite.Suite
	tests.TestCase

	messageService *services.MessageService

	// Users with different role levels
	superAdmin   *models.User
	admin        *models.User // Level 80
	moderator    *models.User // Level 40
	member       *models.User // Level 20
	guest        *models.User // Level 10
	noRoleUser   *models.User // Level 0 (no role)

	// Roles
	adminRole     *models.Role
	moderatorRole *models.Role
	memberRole    *models.Role
	guestRole     *models.Role
}

func TestMessagePermissionTestSuite(t *testing.T) {
	suite.Run(t, new(MessagePermissionTestSuite))
}

func (s *MessagePermissionTestSuite) SetupTest() {
	s.RefreshDatabase()
	s.messageService = services.NewMessageService()

	password, _ := facades.Hash().Make("password123")

	// Create roles with different levels
	s.adminRole = &models.Role{
		Name:        "Administrator",
		Slug:        "admin",
		Description: "Administrator role",
		Level:       80,
		IsActive:    true,
	}
	facades.Orm().Query().Create(s.adminRole)

	s.moderatorRole = &models.Role{
		Name:        "Moderator",
		Slug:        "moderator",
		Description: "Moderator role",
		Level:       40,
		IsActive:    true,
	}
	facades.Orm().Query().Create(s.moderatorRole)

	s.memberRole = &models.Role{
		Name:        "Member",
		Slug:        "member",
		Description: "Member role",
		Level:       20,
		IsActive:    true,
	}
	facades.Orm().Query().Create(s.memberRole)

	s.guestRole = &models.Role{
		Name:        "Guest",
		Slug:        "guest",
		Description: "Guest role",
		Level:       10,
		IsActive:    true,
	}
	facades.Orm().Query().Create(s.guestRole)

	// Create super admin user
	s.superAdmin = &models.User{
		Name:         "Super Admin",
		Email:        uniqueEmail("superadmin"),
		Password:     password,
		IsActive:     true,
		IsSuperAdmin: true,
	}
	facades.Orm().Query().Create(s.superAdmin)

	// Create admin user (level 80)
	s.admin = &models.User{
		Name:     "Admin User",
		Email:    uniqueEmail("admin"),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.admin)
	s.assignRole(s.admin, s.adminRole)

	// Create moderator user (level 40)
	s.moderator = &models.User{
		Name:     "Moderator User",
		Email:    uniqueEmail("moderator"),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.moderator)
	s.assignRole(s.moderator, s.moderatorRole)

	// Create member user (level 20)
	s.member = &models.User{
		Name:     "Member User",
		Email:    uniqueEmail("member"),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.member)
	s.assignRole(s.member, s.memberRole)

	// Create guest user (level 10)
	s.guest = &models.User{
		Name:     "Guest User",
		Email:    uniqueEmail("guest"),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.guest)
	s.assignRole(s.guest, s.guestRole)

	// Create user with no role (level 0)
	s.noRoleUser = &models.User{
		Name:     "No Role User",
		Email:    uniqueEmail("norole"),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.noRoleUser)

	// Reload users with roles
	s.reloadUsersWithRoles()
}

func (s *MessagePermissionTestSuite) reloadUsersWithRoles() {
	// Load user and roles separately since GORM many2many has issues with our custom user_roles table
	s.loadUserWithRoles(s.superAdmin)
	s.loadUserWithRoles(s.admin)
	s.loadUserWithRoles(s.moderator)
	s.loadUserWithRoles(s.member)
	s.loadUserWithRoles(s.guest)
	s.loadUserWithRoles(s.noRoleUser)
}

func (s *MessagePermissionTestSuite) loadUserWithRoles(user *models.User) {
	// Reload user
	facades.Orm().Query().Find(user, user.ID)

	// Load roles manually via user_roles pivot table using raw SQL
	var roles []models.Role
	facades.Orm().Query().Raw(`
		SELECT r.* FROM roles r
		INNER JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = ? AND ur.is_active = true AND r.is_active = true
	`, user.ID).Scan(&roles)

	user.Roles = roles
}

// assignRole assigns a role to a user via the user_roles table
func (s *MessagePermissionTestSuite) assignRole(user *models.User, role *models.Role) {
	// Insert into user_roles pivot table (includes assigned_at and notes to satisfy NOT NULL constraints)
	facades.Orm().Query().Exec(
		"INSERT INTO user_roles (user_id, role_id, is_active, assigned_at, notes, created_at, updated_at) VALUES (?, ?, true, NOW(), '', NOW(), NOW())",
		user.ID, role.ID,
	)
}

// Test GetRoleLevel returns correct levels
func (s *MessagePermissionTestSuite) TestGetRoleLevel_ReturnsCorrectLevels() {
	// Debug: Check user_roles table
	userRoleCount, _ := facades.Orm().Query().Table("user_roles").Count()
	s.T().Logf("user_roles count: %d", userRoleCount)

	// Check admin's user_roles
	var adminUserRoles []struct {
		UserID   uint `gorm:"column:user_id"`
		RoleID   uint `gorm:"column:role_id"`
		IsActive bool `gorm:"column:is_active"`
	}
	facades.Orm().Query().Table("user_roles").Where("user_id = ?", s.admin.ID).Scan(&adminUserRoles)
	s.T().Logf("Admin user ID: %d, user_roles: %+v", s.admin.ID, adminUserRoles)

	// Check roles table
	roleCount, _ := facades.Orm().Query().Model(&models.Role{}).Count()
	s.T().Logf("roles count: %d", roleCount)

	s.T().Logf("Admin roles: %+v", s.admin.Roles)
	s.T().Logf("Admin role level: %d", s.admin.GetRoleLevel())
	s.T().Logf("Member roles: %+v", s.member.Roles)
	s.T().Logf("Member role level: %d", s.member.GetRoleLevel())

	s.Equal(80, s.admin.GetRoleLevel(), "Admin should have level 80")
	s.Equal(40, s.moderator.GetRoleLevel(), "Moderator should have level 40")
	s.Equal(20, s.member.GetRoleLevel(), "Member should have level 20")
	s.Equal(10, s.guest.GetRoleLevel(), "Guest should have level 10")
	s.Equal(0, s.noRoleUser.GetRoleLevel(), "User with no role should have level 0")
}

// Test: Super admin can message anyone
func (s *MessagePermissionTestSuite) TestSuperAdmin_CanMessageAnyone() {
	s.True(s.superAdmin.CanMessageUser(s.admin), "Super admin should be able to message admin")
	s.True(s.superAdmin.CanMessageUser(s.moderator), "Super admin should be able to message moderator")
	s.True(s.superAdmin.CanMessageUser(s.member), "Super admin should be able to message member")
	s.True(s.superAdmin.CanMessageUser(s.guest), "Super admin should be able to message guest")
	s.True(s.superAdmin.CanMessageUser(s.noRoleUser), "Super admin should be able to message no-role user")
}

// Test: Admin (level 80) can message lower level users
func (s *MessagePermissionTestSuite) TestAdmin_CanMessageLowerLevelUsers() {
	s.True(s.admin.CanMessageUser(s.moderator), "Admin should be able to message moderator (lower level)")
	s.True(s.admin.CanMessageUser(s.member), "Admin should be able to message member (lower level)")
	s.True(s.admin.CanMessageUser(s.guest), "Admin should be able to message guest (lower level)")
	s.True(s.admin.CanMessageUser(s.noRoleUser), "Admin should be able to message no-role user (lower level)")
}

// Test: Admin (level 80) can message same level users
func (s *MessagePermissionTestSuite) TestAdmin_CanMessageSameLevelUsers() {
	// Create another admin
	password, _ := facades.Hash().Make("password123")
	anotherAdmin := &models.User{
		Name:     "Another Admin",
		Email:    uniqueEmail("anotheradmin"),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(anotherAdmin)
	s.assignRole(anotherAdmin, s.adminRole)
	facades.Orm().Query().With("Roles").Find(anotherAdmin, anotherAdmin.ID)

	s.True(s.admin.CanMessageUser(anotherAdmin), "Admin should be able to message another admin (same level)")
}

// Test: Member (level 20) CANNOT message higher level users
func (s *MessagePermissionTestSuite) TestMember_CannotMessageHigherLevelUsers() {
	s.False(s.member.CanMessageUser(s.admin), "Member should NOT be able to message admin (higher level)")
	s.False(s.member.CanMessageUser(s.moderator), "Member should NOT be able to message moderator (higher level)")
}

// Test: Member (level 20) can message same or lower level users
func (s *MessagePermissionTestSuite) TestMember_CanMessageSameOrLowerLevel() {
	// Create another member
	password, _ := facades.Hash().Make("password123")
	anotherMember := &models.User{
		Name:     "Another Member",
		Email:    uniqueEmail("anothermember"),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(anotherMember)
	s.assignRole(anotherMember, s.memberRole)
	facades.Orm().Query().With("Roles").Find(anotherMember, anotherMember.ID)

	s.True(s.member.CanMessageUser(anotherMember), "Member should be able to message another member (same level)")
	s.True(s.member.CanMessageUser(s.guest), "Member should be able to message guest (lower level)")
	s.True(s.member.CanMessageUser(s.noRoleUser), "Member should be able to message no-role user (lower level)")
}

// Test: Guest (level 10) can only message guests or no-role users
func (s *MessagePermissionTestSuite) TestGuest_CanOnlyMessageLowerOrSameLevel() {
	s.False(s.guest.CanMessageUser(s.admin), "Guest should NOT be able to message admin")
	s.False(s.guest.CanMessageUser(s.moderator), "Guest should NOT be able to message moderator")
	s.False(s.guest.CanMessageUser(s.member), "Guest should NOT be able to message member")

	// Create another guest
	password, _ := facades.Hash().Make("password123")
	anotherGuest := &models.User{
		Name:     "Another Guest",
		Email:    uniqueEmail("anotherguest"),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(anotherGuest)
	s.assignRole(anotherGuest, s.guestRole)
	facades.Orm().Query().With("Roles").Find(anotherGuest, anotherGuest.ID)

	s.True(s.guest.CanMessageUser(anotherGuest), "Guest should be able to message another guest (same level)")
	s.True(s.guest.CanMessageUser(s.noRoleUser), "Guest should be able to message no-role user (lower level)")
}

// Test: User with no role can only message other no-role users
func (s *MessagePermissionTestSuite) TestNoRoleUser_CanOnlyMessageOtherNoRoleUsers() {
	s.False(s.noRoleUser.CanMessageUser(s.admin), "No-role user should NOT be able to message admin")
	s.False(s.noRoleUser.CanMessageUser(s.moderator), "No-role user should NOT be able to message moderator")
	s.False(s.noRoleUser.CanMessageUser(s.member), "No-role user should NOT be able to message member")
	s.False(s.noRoleUser.CanMessageUser(s.guest), "No-role user should NOT be able to message guest")

	// Create another no-role user
	password, _ := facades.Hash().Make("password123")
	anotherNoRole := &models.User{
		Name:     "Another No Role",
		Email:    uniqueEmail("anothernorole"),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(anotherNoRole)
	facades.Orm().Query().With("Roles").Find(anotherNoRole, anotherNoRole.ID)

	s.True(s.noRoleUser.CanMessageUser(anotherNoRole), "No-role user should be able to message another no-role user")
}

// Test: User can message themselves
func (s *MessagePermissionTestSuite) TestUser_CanMessageSelf() {
	s.True(s.member.CanMessageUser(s.member), "User should be able to message themselves")
}

// Test: SendMessage respects permission - allowed case
func (s *MessagePermissionTestSuite) TestSendMessage_AllowedByPermission() {
	// Admin (level 80) sending to member (level 20) - should succeed
	message, err := s.messageService.SendMessage(
		s.admin.ID,
		s.member.ID,
		"Hello from admin to member",
		models.MessageTypeDirect,
	)

	s.NoError(err, "Admin should be able to send message to member")
	s.NotNil(message)
	s.Equal("Hello from admin to member", message.Content)
	s.Equal(s.admin.ID, message.SenderID)
	s.Equal(s.member.ID, *message.RecipientID)
}

// Test: SendMessage respects permission - denied case
func (s *MessagePermissionTestSuite) TestSendMessage_DeniedByPermission() {
	// Member (level 20) trying to send to admin (level 80) - should fail
	message, err := s.messageService.SendMessage(
		s.member.ID,
		s.admin.ID,
		"Hello from member to admin",
		models.MessageTypeDirect,
	)

	s.Error(err, "Member should NOT be able to send message to admin")
	s.Nil(message)
	s.Contains(err.Error(), "insufficient permissions", "Error should mention insufficient permissions")
}

// Test: Notification is created when message is sent
func (s *MessagePermissionTestSuite) TestSendMessage_CreatesNotification() {
	// Admin sending to member
	_, err := s.messageService.SendMessage(
		s.admin.ID,
		s.member.ID,
		"Test notification message",
		models.MessageTypeDirect,
	)
	s.NoError(err)

	// Check that a notification was created for the recipient
	var notification models.Notification
	err = facades.Orm().Query().
		Where("user_id = ?", s.member.ID).
		Where("type = ?", "message").
		Where("trigger_user_id = ?", s.admin.ID).
		First(&notification)

	s.NoError(err, "Notification should be created for the message recipient")
	s.Contains(notification.Title, s.admin.Name, "Notification title should contain sender's name")
	s.Equal("message", notification.Type)
}
