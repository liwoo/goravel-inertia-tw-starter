package services

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/goravel/framework/facades"
)

type SSEEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	Time time.Time   `json:"time"`
}

type SSEClient struct {
	ID     string
	UserID uint
	Events chan SSEEvent
}

type SSEService struct {
	clients      map[string]*SSEClient
	userClients  map[uint][]*SSEClient
	clientsMutex sync.RWMutex
	eventBus     chan BroadcastEvent
}

type BroadcastEvent struct {
	UserIDs []uint
	Event   SSEEvent
}

var sseServiceInstance *SSEService
var sseServiceOnce sync.Once

func NewSSEService() *SSEService {
	sseServiceOnce.Do(func() {
		sseServiceInstance = &SSEService{
			clients:     make(map[string]*SSEClient),
			userClients: make(map[uint][]*SSEClient),
			eventBus:    make(chan BroadcastEvent, 1000),
		}
		go sseServiceInstance.eventBroadcaster()
	})
	return sseServiceInstance
}

func (s *SSEService) CreateClient(userID uint) *SSEClient {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	client := &SSEClient{
		ID:     uuid.New().String(),
		UserID: userID,
		Events: make(chan SSEEvent, 100),
	}

	s.clients[client.ID] = client
	s.userClients[userID] = append(s.userClients[userID], client)

	facades.Log().Info("SSE client connected", map[string]interface{}{
		"client_id": client.ID,
		"user_id":   userID,
	})

	return client
}

func (s *SSEService) RemoveClient(clientID string) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	client, exists := s.clients[clientID]
	if !exists {
		return
	}

	// Remove from userClients
	if clients, ok := s.userClients[client.UserID]; ok {
		for i, c := range clients {
			if c.ID == clientID {
				s.userClients[client.UserID] = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		if len(s.userClients[client.UserID]) == 0 {
			delete(s.userClients, client.UserID)
		}
	}

	// Close the events channel and remove the client
	close(client.Events)
	delete(s.clients, clientID)

	facades.Log().Info("SSE client disconnected", map[string]interface{}{
		"client_id": clientID,
		"user_id":   client.UserID,
	})
}

func (s *SSEService) SendToUser(userID uint, eventType string, data interface{}) {
	event := SSEEvent{
		Type: eventType,
		Data: data,
		Time: time.Now(),
	}

	s.eventBus <- BroadcastEvent{
		UserIDs: []uint{userID},
		Event:   event,
	}
}

func (s *SSEService) SendToUsers(userIDs []uint, eventType string, data interface{}) {
	event := SSEEvent{
		Type: eventType,
		Data: data,
		Time: time.Now(),
	}

	s.eventBus <- BroadcastEvent{
		UserIDs: userIDs,
		Event:   event,
	}
}

func (s *SSEService) BroadcastToAll(eventType string, data interface{}) {
	s.clientsMutex.RLock()
	userIDs := make([]uint, 0, len(s.userClients))
	for userID := range s.userClients {
		userIDs = append(userIDs, userID)
	}
	s.clientsMutex.RUnlock()

	if len(userIDs) > 0 {
		s.SendToUsers(userIDs, eventType, data)
	}
}

func (s *SSEService) eventBroadcaster() {
	for broadcast := range s.eventBus {
		s.clientsMutex.RLock()
		for _, userID := range broadcast.UserIDs {
			if clients, ok := s.userClients[userID]; ok {
				for _, client := range clients {
					select {
					case client.Events <- broadcast.Event:
						// Event sent successfully
					default:
						// Channel is full, skip this event
						facades.Log().Warning("SSE client event channel full", map[string]interface{}{
							"client_id":  client.ID,
							"user_id":    userID,
							"event_type": broadcast.Event.Type,
						})
					}
				}
			}
		}
		s.clientsMutex.RUnlock()
	}
}

// Message-related events
func (s *SSEService) NotifyNewMessage(senderID, recipientID uint, message interface{}) {
	s.SendToUser(recipientID, "message:new", map[string]interface{}{
		"sender_id": senderID,
		"message":   message,
	})
}

func (s *SSEService) NotifyMessageRead(userID uint, messageID uint) {
	s.SendToUser(userID, "message:read", map[string]interface{}{
		"message_id": messageID,
	})
}

func (s *SSEService) NotifyMessageUpdated(userID uint, message interface{}) {
	s.SendToUser(userID, "message:updated", message)
}

func (s *SSEService) NotifyMessageDeleted(userID uint, messageID uint) {
	s.SendToUser(userID, "message:deleted", map[string]interface{}{
		"message_id": messageID,
	})
}

func (s *SSEService) UpdateUnreadCount(userID uint, count int64) {
	s.SendToUser(userID, "message:unread_count", map[string]interface{}{
		"count": count,
	})
}

// Notification-related events
func (s *SSEService) NotifyNewNotification(userID uint, notification interface{}) {
	s.SendToUser(userID, "notification:new", notification)
}

func (s *SSEService) NotifyNotificationRead(userID uint, notificationID uint) {
	s.SendToUser(userID, "notification:read", map[string]interface{}{
		"notification_id": notificationID,
	})
}

func (s *SSEService) NotifyNotificationDismissed(userID uint, notificationID uint) {
	s.SendToUser(userID, "notification:dismissed", map[string]interface{}{
		"notification_id": notificationID,
	})
}

func (s *SSEService) UpdateNotificationCounts(userID uint, counts interface{}) {
	s.SendToUser(userID, "notification:counts", counts)
}

// System-wide events
func (s *SSEService) BroadcastSystemNotification(notification interface{}) {
	s.BroadcastToAll("notification:system", notification)
}
