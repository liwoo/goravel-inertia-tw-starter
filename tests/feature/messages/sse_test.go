package messages

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"starter-project/app/models"
	"starter-project/app/services"
	"starter-project/tests"
)

type SSETestSuite struct {
	suite.Suite
	tests.TestCase

	sseService     *services.SSEService
	messageService *services.MessageService

	// Users
	user1 *models.User
	user2 *models.User
	user3 *models.User
}

func TestSSETestSuite(t *testing.T) {
	suite.Run(t, new(SSETestSuite))
}

// waitForEventType waits for a specific event type, skipping presence events
func (s *SSETestSuite) waitForEventType(client *services.SSEClient, expectedType string, timeout time.Duration) (services.SSEEvent, bool) {
	deadline := time.After(timeout)
	for {
		select {
		case event := <-client.Events:
			// Skip presence events when waiting for other event types
			if strings.HasPrefix(event.Type, "presence:") && !strings.HasPrefix(expectedType, "presence:") {
				continue
			}
			return event, true
		case <-deadline:
			return services.SSEEvent{}, false
		}
	}
}

// waitForAnyEvent waits for any non-presence event
func (s *SSETestSuite) waitForAnyEvent(client *services.SSEClient, timeout time.Duration) (services.SSEEvent, bool) {
	deadline := time.After(timeout)
	for {
		select {
		case event := <-client.Events:
			// Skip presence events
			if strings.HasPrefix(event.Type, "presence:") {
				continue
			}
			return event, true
		case <-deadline:
			return services.SSEEvent{}, false
		}
	}
}

func (s *SSETestSuite) SetupTest() {
	s.RefreshDatabase()
	s.sseService = services.NewSSEService()
	s.sseService.ResetForTesting() // Reset singleton state
	s.messageService = services.NewMessageService()

	password, _ := facades.Hash().Make("password123")

	// Create users
	s.user1 = &models.User{
		Name:     "User One",
		Email:    fmt.Sprintf("user1_%d@example.com", time.Now().UnixNano()),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.user1)

	s.user2 = &models.User{
		Name:     "User Two",
		Email:    fmt.Sprintf("user2_%d@example.com", time.Now().UnixNano()),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.user2)

	s.user3 = &models.User{
		Name:     "User Three",
		Email:    fmt.Sprintf("user3_%d@example.com", time.Now().UnixNano()),
		Password: password,
		IsActive: true,
	}
	facades.Orm().Query().Create(s.user3)
}

// Test: CreateClient creates a client with unique ID
func (s *SSETestSuite) TestCreateClient_CreatesUniqueClient() {
	client1 := s.sseService.CreateClient(s.user1.ID)
	client2 := s.sseService.CreateClient(s.user1.ID)

	s.NotNil(client1)
	s.NotNil(client2)
	s.NotEqual(client1.ID, client2.ID, "Each client should have a unique ID")
	s.Equal(s.user1.ID, client1.UserID)
	s.Equal(s.user1.ID, client2.UserID)

	// Cleanup
	s.sseService.RemoveClient(client1.ID)
	s.sseService.RemoveClient(client2.ID)
}

// Test: CreateClient creates client with events channel
func (s *SSETestSuite) TestCreateClient_HasEventsChannel() {
	client := s.sseService.CreateClient(s.user1.ID)
	s.NotNil(client.Events, "Client should have an events channel")

	// Cleanup
	s.sseService.RemoveClient(client.ID)
}

// Test: RemoveClient closes the events channel
func (s *SSETestSuite) TestRemoveClient_ClosesEventsChannel() {
	client := s.sseService.CreateClient(s.user1.ID)
	clientID := client.ID
	events := client.Events

	s.sseService.RemoveClient(clientID)

	// Channel should be closed
	_, open := <-events
	s.False(open, "Events channel should be closed after removing client")
}

// Test: SendToUser sends event to specific user's clients
func (s *SSETestSuite) TestSendToUser_SendsToCorrectUser() {
	client1 := s.sseService.CreateClient(s.user1.ID)
	client2 := s.sseService.CreateClient(s.user2.ID)

	defer s.sseService.RemoveClient(client1.ID)
	defer s.sseService.RemoveClient(client2.ID)

	testData := map[string]interface{}{"test": "data"}

	// Use goroutine and channel to handle potential blocking
	var wg sync.WaitGroup
	var user1Received, user2Received bool

	wg.Add(2)

	go func() {
		defer wg.Done()
		event, received := s.waitForEventType(client1, "test:event", 500*time.Millisecond)
		if received {
			user1Received = true
			s.Equal("test:event", event.Type)
			s.NotNil(event.Data)
		}
	}()

	go func() {
		defer wg.Done()
		_, received := s.waitForEventType(client2, "test:event", 500*time.Millisecond)
		user2Received = received
	}()

	// Send event to user1 only
	s.sseService.SendToUser(s.user1.ID, "test:event", testData)

	wg.Wait()

	s.True(user1Received, "User1 should receive the event")
	s.False(user2Received, "User2 should NOT receive user1's event")
}

// Test: SendToUsers sends to multiple users
func (s *SSETestSuite) TestSendToUsers_SendsToMultipleUsers() {
	client1 := s.sseService.CreateClient(s.user1.ID)
	client2 := s.sseService.CreateClient(s.user2.ID)
	client3 := s.sseService.CreateClient(s.user3.ID)

	defer s.sseService.RemoveClient(client1.ID)
	defer s.sseService.RemoveClient(client2.ID)
	defer s.sseService.RemoveClient(client3.ID)

	var wg sync.WaitGroup
	receivedCount := 0
	var mu sync.Mutex

	wg.Add(3)

	checkReceived := func(client *services.SSEClient) {
		defer wg.Done()
		_, received := s.waitForEventType(client, "multi:event", 500*time.Millisecond)
		if received {
			mu.Lock()
			receivedCount++
			mu.Unlock()
		}
	}

	go checkReceived(client1)
	go checkReceived(client2)
	go checkReceived(client3)

	// Send to user1 and user2 only
	s.sseService.SendToUsers([]uint{s.user1.ID, s.user2.ID}, "multi:event", nil)

	wg.Wait()

	s.Equal(2, receivedCount, "Should send to exactly 2 users")
}

// Test: Multiple clients for same user receive events
func (s *SSETestSuite) TestSendToUser_MultipleClientsReceiveEvent() {
	client1 := s.sseService.CreateClient(s.user1.ID)
	client2 := s.sseService.CreateClient(s.user1.ID)

	defer s.sseService.RemoveClient(client1.ID)
	defer s.sseService.RemoveClient(client2.ID)

	var wg sync.WaitGroup
	receivedCount := 0
	var mu sync.Mutex

	wg.Add(2)

	checkReceived := func(client *services.SSEClient) {
		defer wg.Done()
		_, received := s.waitForEventType(client, "test:event", 500*time.Millisecond)
		if received {
			mu.Lock()
			receivedCount++
			mu.Unlock()
		}
	}

	go checkReceived(client1)
	go checkReceived(client2)

	s.sseService.SendToUser(s.user1.ID, "test:event", nil)

	wg.Wait()

	s.Equal(2, receivedCount, "Both clients for same user should receive the event")
}

// Test: NotifyNewMessage sends correct event type
func (s *SSETestSuite) TestNotifyNewMessage_SendsCorrectEventType() {
	client := s.sseService.CreateClient(s.user2.ID)
	defer s.sseService.RemoveClient(client.ID)

	mockMessage := map[string]interface{}{"id": 1, "content": "Hello"}
	s.sseService.NotifyNewMessage(s.user1.ID, s.user2.ID, mockMessage)

	receivedEvent, received := s.waitForEventType(client, "message:new", 500*time.Millisecond)

	s.True(received, "Should receive new message notification")
	s.Equal("message:new", receivedEvent.Type)
}

// Test: NotifyMessageRead sends correct event
func (s *SSETestSuite) TestNotifyMessageRead_SendsCorrectEvent() {
	client := s.sseService.CreateClient(s.user1.ID)
	defer s.sseService.RemoveClient(client.ID)

	s.sseService.NotifyMessageRead(s.user1.ID, 123)

	receivedEvent, received := s.waitForEventType(client, "message:read", 500*time.Millisecond)

	s.True(received, "Should receive message read notification")
	s.Equal("message:read", receivedEvent.Type)
}

// Test: UpdateUnreadCount sends correct event
func (s *SSETestSuite) TestUpdateUnreadCount_SendsCorrectEvent() {
	client := s.sseService.CreateClient(s.user1.ID)
	defer s.sseService.RemoveClient(client.ID)

	s.sseService.UpdateUnreadCount(s.user1.ID, 5)

	receivedEvent, received := s.waitForEventType(client, "message:unread_count", 500*time.Millisecond)

	s.True(received, "Should receive unread count update")
	s.Equal("message:unread_count", receivedEvent.Type)
}

// Test: NotifyNewNotification sends correct event
func (s *SSETestSuite) TestNotifyNewNotification_SendsCorrectEvent() {
	client := s.sseService.CreateClient(s.user1.ID)
	defer s.sseService.RemoveClient(client.ID)

	mockNotification := map[string]interface{}{
		"id":    1,
		"title": "Test Notification",
	}
	s.sseService.NotifyNewNotification(s.user1.ID, mockNotification)

	receivedEvent, received := s.waitForEventType(client, "notification:new", 500*time.Millisecond)

	s.True(received, "Should receive new notification")
	s.Equal("notification:new", receivedEvent.Type)
}

// Test: UpdateNotificationCounts sends correct event
func (s *SSETestSuite) TestUpdateNotificationCounts_SendsCorrectEvent() {
	client := s.sseService.CreateClient(s.user1.ID)
	defer s.sseService.RemoveClient(client.ID)

	counts := map[string]int64{"unread": 3, "total": 10}
	s.sseService.UpdateNotificationCounts(s.user1.ID, counts)

	receivedEvent, received := s.waitForEventType(client, "notification:counts", 500*time.Millisecond)

	s.True(received, "Should receive notification counts update")
	s.Equal("notification:counts", receivedEvent.Type)
}

// Test: Event includes timestamp
func (s *SSETestSuite) TestEvent_IncludesTimestamp() {
	client := s.sseService.CreateClient(s.user1.ID)
	defer s.sseService.RemoveClient(client.ID)

	beforeSend := time.Now()

	s.sseService.SendToUser(s.user1.ID, "test:event", nil)

	receivedEvent, received := s.waitForEventType(client, "test:event", 500*time.Millisecond)

	s.True(received, "Should receive event")
	s.False(receivedEvent.Time.IsZero(), "Event should have a timestamp")
	s.True(receivedEvent.Time.After(beforeSend) || receivedEvent.Time.Equal(beforeSend),
		"Event timestamp should be after or equal to send time")
}

// Test: SSEService is singleton
func (s *SSETestSuite) TestNewSSEService_ReturnsSingleton() {
	service1 := services.NewSSEService()
	service2 := services.NewSSEService()

	s.Same(service1, service2, "NewSSEService should return the same singleton instance")
}

// Test: Send to non-existent user doesn't error
func (s *SSETestSuite) TestSendToUser_NonExistentUser_NoError() {
	// This shouldn't panic or error
	s.NotPanics(func() {
		s.sseService.SendToUser(99999, "test:event", nil)
	}, "Sending to non-existent user should not panic")
}

// Test: RemoveClient with invalid ID doesn't panic
func (s *SSETestSuite) TestRemoveClient_InvalidID_NoError() {
	s.NotPanics(func() {
		s.sseService.RemoveClient("non-existent-id")
	}, "Removing non-existent client should not panic")
}
