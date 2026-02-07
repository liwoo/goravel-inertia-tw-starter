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
	// Presence tracking
	userLastSeen  map[uint]time.Time
	presenceMutex sync.RWMutex
}

type BroadcastEvent struct {
	UserIDs []uint
	Event   SSEEvent
}

var sseServiceInstance *SSEService
var sseServiceOnce sync.Once

// safeLog is a helper that safely logs messages without panicking if facades aren't initialized
func safeLog(level string, message string, context map[string]interface{}) {
	defer func() {
		recover() // Silently ignore any panics from logging
	}()

	log := facades.Log()
	if log == nil {
		return
	}

	switch level {
	case "info":
		log.Info(message, context)
	case "warning":
		log.Warning(message, context)
	case "error":
		log.Error(message, context)
	}
}

func NewSSEService() *SSEService {
	sseServiceOnce.Do(func() {
		sseServiceInstance = &SSEService{
			clients:      make(map[string]*SSEClient),
			userClients:  make(map[uint][]*SSEClient),
			eventBus:     make(chan BroadcastEvent, 1000),
			userLastSeen: make(map[uint]time.Time),
		}
		go sseServiceInstance.eventBroadcaster()
	})
	return sseServiceInstance
}

// ResetForTesting clears all clients and presence data - only use in tests
func (s *SSEService) ResetForTesting() {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()
	s.presenceMutex.Lock()
	defer s.presenceMutex.Unlock()

	// Close all client channels
	for _, client := range s.clients {
		close(client.Events)
	}

	// Reset all maps
	s.clients = make(map[string]*SSEClient)
	s.userClients = make(map[uint][]*SSEClient)
	s.userLastSeen = make(map[uint]time.Time)
}

func (s *SSEService) CreateClient(userID uint) *SSEClient {
	s.clientsMutex.Lock()

	// Check if this is the first client for this user (user coming online)
	wasOffline := len(s.userClients[userID]) == 0

	client := &SSEClient{
		ID:     uuid.New().String(),
		UserID: userID,
		Events: make(chan SSEEvent, 100),
	}

	s.clients[client.ID] = client
	s.userClients[userID] = append(s.userClients[userID], client)

	// Update last seen
	s.presenceMutex.Lock()
	s.userLastSeen[userID] = time.Now()
	s.presenceMutex.Unlock()

	s.clientsMutex.Unlock()

	safeLog("info", "SSE client connected", map[string]interface{}{
		"client_id": client.ID,
		"user_id":   userID,
	})

	// Broadcast presence change if user just came online
	if wasOffline {
		s.broadcastPresenceChange(userID, true)
	}

	return client
}

func (s *SSEService) RemoveClient(clientID string) {
	s.clientsMutex.Lock()

	client, exists := s.clients[clientID]
	if !exists {
		s.clientsMutex.Unlock()
		return
	}

	userID := client.UserID

	// Remove from userClients
	if clients, ok := s.userClients[userID]; ok {
		for i, c := range clients {
			if c.ID == clientID {
				s.userClients[userID] = append(clients[:i], clients[i+1:]...)
				break
			}
		}
	}

	// Check if this was the last client for this user (user going offline)
	isNowOffline := len(s.userClients[userID]) == 0
	if isNowOffline {
		delete(s.userClients, userID)
	}

	// Update last seen
	s.presenceMutex.Lock()
	s.userLastSeen[userID] = time.Now()
	s.presenceMutex.Unlock()

	// Close the events channel and remove the client
	close(client.Events)
	delete(s.clients, clientID)

	s.clientsMutex.Unlock()

	safeLog("info", "SSE client disconnected", map[string]interface{}{
		"client_id": clientID,
		"user_id":   userID,
	})

	// Broadcast presence change if user just went offline
	if isNowOffline {
		s.broadcastPresenceChange(userID, false)
	}
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
						safeLog("warning", "SSE client event channel full", map[string]interface{}{
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

// Presence-related methods

// broadcastPresenceChange notifies all connected users about a user's status change
func (s *SSEService) broadcastPresenceChange(userID uint, isOnline bool) {
	s.BroadcastToAll("presence:change", map[string]interface{}{
		"user_id":   userID,
		"is_online": isOnline,
		"timestamp": time.Now(),
	})
}

// IsUserOnline checks if a user has any active SSE connections
func (s *SSEService) IsUserOnline(userID uint) bool {
	s.clientsMutex.RLock()
	defer s.clientsMutex.RUnlock()
	return len(s.userClients[userID]) > 0
}

// GetOnlineUserIDs returns a list of all currently connected user IDs
func (s *SSEService) GetOnlineUserIDs() []uint {
	s.clientsMutex.RLock()
	defer s.clientsMutex.RUnlock()

	userIDs := make([]uint, 0, len(s.userClients))
	for userID := range s.userClients {
		userIDs = append(userIDs, userID)
	}
	return userIDs
}

// GetUserLastSeen returns the last time a user was seen (connected or disconnected)
func (s *SSEService) GetUserLastSeen(userID uint) (time.Time, bool) {
	s.presenceMutex.RLock()
	defer s.presenceMutex.RUnlock()
	lastSeen, exists := s.userLastSeen[userID]
	return lastSeen, exists
}

// GetPresenceStatus returns online status and last seen for a user
func (s *SSEService) GetPresenceStatus(userID uint) map[string]interface{} {
	isOnline := s.IsUserOnline(userID)
	lastSeen, hasLastSeen := s.GetUserLastSeen(userID)

	status := map[string]interface{}{
		"user_id":   userID,
		"is_online": isOnline,
	}

	if hasLastSeen {
		status["last_seen"] = lastSeen
	}

	return status
}

// GetBulkPresenceStatus returns presence status for multiple users
func (s *SSEService) GetBulkPresenceStatus(userIDs []uint) []map[string]interface{} {
	statuses := make([]map[string]interface{}, len(userIDs))
	for i, userID := range userIDs {
		statuses[i] = s.GetPresenceStatus(userID)
	}
	return statuses
}

// GetOnlineCount returns the number of currently online users
func (s *SSEService) GetOnlineCount() int {
	s.clientsMutex.RLock()
	defer s.clientsMutex.RUnlock()
	return len(s.userClients)
}

// SendPresenceToUser sends the current online users list to a specific user
func (s *SSEService) SendPresenceToUser(userID uint) {
	onlineUserIDs := s.GetOnlineUserIDs()
	s.SendToUser(userID, "presence:initial", map[string]interface{}{
		"online_users": onlineUserIDs,
		"timestamp":    time.Now(),
	})
}
