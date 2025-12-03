package messages

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"smedi-sme-db/app/services"
	"smedi-sme-db/tests"
)

type PresenceTestSuite struct {
	suite.Suite
	tests.TestCase

	sseService *services.SSEService

	// Use simple user IDs for testing (no need for actual users in DB)
	userID1 uint
	userID2 uint
	userID3 uint
}

func TestPresenceTestSuite(t *testing.T) {
	suite.Run(t, new(PresenceTestSuite))
}

func (s *PresenceTestSuite) SetupTest() {
	s.RefreshDatabase()
	s.sseService = services.NewSSEService()

	// Reset SSE service state for clean tests
	s.sseService.ResetForTesting()

	// Use unique IDs to avoid interference between parallel test runs
	baseID := uint(time.Now().UnixNano() % 1000000)
	s.userID1 = baseID + 1
	s.userID2 = baseID + 2
	s.userID3 = baseID + 3
}

// Test: User comes online when client is created
func (s *PresenceTestSuite) TestIsUserOnline_TrueWhenClientExists() {
	// User is offline initially
	s.False(s.sseService.IsUserOnline(s.userID1), "User should be offline initially")

	// Create client - user comes online
	client := s.sseService.CreateClient(s.userID1)
	defer s.sseService.RemoveClient(client.ID)

	s.True(s.sseService.IsUserOnline(s.userID1), "User should be online after creating client")
}

// Test: User goes offline when last client is removed
func (s *PresenceTestSuite) TestIsUserOnline_FalseWhenClientRemoved() {
	client := s.sseService.CreateClient(s.userID1)

	s.True(s.sseService.IsUserOnline(s.userID1), "User should be online with active client")

	s.sseService.RemoveClient(client.ID)

	s.False(s.sseService.IsUserOnline(s.userID1), "User should be offline after client removal")
}

// Test: User stays online with multiple clients
func (s *PresenceTestSuite) TestIsUserOnline_StaysOnlineWithMultipleClients() {
	client1 := s.sseService.CreateClient(s.userID1)
	client2 := s.sseService.CreateClient(s.userID1)

	s.True(s.sseService.IsUserOnline(s.userID1), "User should be online with multiple clients")

	// Remove first client
	s.sseService.RemoveClient(client1.ID)

	s.True(s.sseService.IsUserOnline(s.userID1), "User should stay online with remaining client")

	// Remove second client
	s.sseService.RemoveClient(client2.ID)

	s.False(s.sseService.IsUserOnline(s.userID1), "User should be offline after all clients removed")
}

// Test: GetOnlineUserIDs returns correct users
func (s *PresenceTestSuite) TestGetOnlineUserIDs_ReturnsCorrectUsers() {
	// Initially no users online
	onlineUsers := s.sseService.GetOnlineUserIDs()
	initialCount := len(onlineUsers)

	// Bring users online
	client1 := s.sseService.CreateClient(s.userID1)
	defer s.sseService.RemoveClient(client1.ID)

	client2 := s.sseService.CreateClient(s.userID2)
	defer s.sseService.RemoveClient(client2.ID)

	onlineUsers = s.sseService.GetOnlineUserIDs()
	s.Equal(initialCount+2, len(onlineUsers), "Should have 2 more online users")

	// Check userID1 and userID2 are in the list
	user1Online := false
	user2Online := false
	for _, id := range onlineUsers {
		if id == s.userID1 {
			user1Online = true
		}
		if id == s.userID2 {
			user2Online = true
		}
	}
	s.True(user1Online, "User1 should be in online users list")
	s.True(user2Online, "User2 should be in online users list")
}

// Test: GetOnlineCount returns correct count
func (s *PresenceTestSuite) TestGetOnlineCount_ReturnsCorrectCount() {
	initialCount := s.sseService.GetOnlineCount()

	client1 := s.sseService.CreateClient(s.userID1)
	s.Equal(initialCount+1, s.sseService.GetOnlineCount(), "Count should increase by 1")

	client2 := s.sseService.CreateClient(s.userID2)
	s.Equal(initialCount+2, s.sseService.GetOnlineCount(), "Count should increase by 2")

	// Same user, multiple clients - count shouldn't change for unique users
	client3 := s.sseService.CreateClient(s.userID1)
	s.Equal(initialCount+2, s.sseService.GetOnlineCount(), "Count should stay same for duplicate user")

	s.sseService.RemoveClient(client1.ID)
	s.Equal(initialCount+2, s.sseService.GetOnlineCount(), "Count should stay same (user1 still has client3)")

	s.sseService.RemoveClient(client3.ID)
	s.Equal(initialCount+1, s.sseService.GetOnlineCount(), "Count should decrease when last user1 client removed")

	s.sseService.RemoveClient(client2.ID)
	s.Equal(initialCount, s.sseService.GetOnlineCount(), "Count should return to initial")
}

// Test: GetPresenceStatus returns correct status
func (s *PresenceTestSuite) TestGetPresenceStatus_ReturnsCorrectStatus() {
	// Offline status
	status := s.sseService.GetPresenceStatus(s.userID1)
	s.Equal(s.userID1, status["user_id"], "Should have correct user_id")
	s.False(status["is_online"].(bool), "Should be offline")

	// Online status
	client := s.sseService.CreateClient(s.userID1)
	defer s.sseService.RemoveClient(client.ID)

	status = s.sseService.GetPresenceStatus(s.userID1)
	s.True(status["is_online"].(bool), "Should be online")
}

// Test: GetBulkPresenceStatus returns status for multiple users
func (s *PresenceTestSuite) TestGetBulkPresenceStatus_ReturnsMultipleStatuses() {
	client1 := s.sseService.CreateClient(s.userID1)
	defer s.sseService.RemoveClient(client1.ID)

	// userID2 is offline, userID1 is online
	statuses := s.sseService.GetBulkPresenceStatus([]uint{s.userID1, s.userID2, s.userID3})

	s.Len(statuses, 3, "Should return status for 3 users")

	// Find statuses
	var user1Status, user2Status, user3Status map[string]interface{}
	for _, status := range statuses {
		switch status["user_id"].(uint) {
		case s.userID1:
			user1Status = status
		case s.userID2:
			user2Status = status
		case s.userID3:
			user3Status = status
		}
	}

	s.True(user1Status["is_online"].(bool), "User1 should be online")
	s.False(user2Status["is_online"].(bool), "User2 should be offline")
	s.False(user3Status["is_online"].(bool), "User3 should be offline")
}

// Test: Last seen is updated on connect
func (s *PresenceTestSuite) TestGetUserLastSeen_UpdatedOnConnect() {
	beforeConnect := time.Now()
	client := s.sseService.CreateClient(s.userID1)
	defer s.sseService.RemoveClient(client.ID)

	lastSeen, exists := s.sseService.GetUserLastSeen(s.userID1)
	s.True(exists, "Last seen should exist after connect")
	s.True(lastSeen.After(beforeConnect) || lastSeen.Equal(beforeConnect), "Last seen should be after connect time")
}

// Test: Last seen is updated on disconnect
func (s *PresenceTestSuite) TestGetUserLastSeen_UpdatedOnDisconnect() {
	client := s.sseService.CreateClient(s.userID1)

	beforeDisconnect := time.Now()
	time.Sleep(10 * time.Millisecond) // Small delay

	s.sseService.RemoveClient(client.ID)

	lastSeen, exists := s.sseService.GetUserLastSeen(s.userID1)
	s.True(exists, "Last seen should exist after disconnect")
	s.True(lastSeen.After(beforeDisconnect), "Last seen should be after disconnect time")
}

// Test: Non-existent user returns offline
func (s *PresenceTestSuite) TestIsUserOnline_FalseForNonExistentUser() {
	s.False(s.sseService.IsUserOnline(999999999), "Non-existent user should be offline")
}

// Test: GetPresenceStatus for non-existent user
func (s *PresenceTestSuite) TestGetPresenceStatus_NonExistentUser() {
	status := s.sseService.GetPresenceStatus(999999999)
	s.Equal(uint(999999999), status["user_id"], "Should have correct user_id")
	s.False(status["is_online"].(bool), "Should be offline")
}
