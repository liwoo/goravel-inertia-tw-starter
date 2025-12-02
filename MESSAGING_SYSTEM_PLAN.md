# Messaging System Completion Plan

## Current State Summary

### What's Working
- **Backend Infrastructure**: Message and Notification models, services, controllers exist
- **Database Schema**: Messages, message_mentions, notifications tables with proper fields
- **Frontend Components**: MessageSidebar, MessageChat, UserMentionInput, NotificationDrawer
- **Real-time Updates**: SSE integration for live message/notification updates
- **Basic CRUD**: Send, view conversations, delete messages, mark as read
- **Notification System**: Full CRUD, batch operations, system notifications with role filtering

### Key Issues to Fix

#### 1. Role-Based Messaging Permission (CRITICAL)
**Current**: Users can only message others with the *exact same role* (`SharesRoleWith()`)
**Required**: Users should message those at their role level OR LOWER

**Role Hierarchy (Level field)**:
- Super Administrator: 100
- Administrator: 80
- Librarian: 60
- Moderator: 40
- Member: 20
- Guest: 10

#### 2. Missing API Routes
These endpoints exist in controllers but are NOT registered in `routes/api.go`:
- `GET /api/messages/inbox`
- `GET /api/messages/sent`
- `GET /api/messages/unread`
- `POST /api/messages/broadcast`

#### 3. Automatic Notifications
Messages are sent but no notification is created for the recipient.

---

## Implementation Plan

### Phase 1: Fix Role-Based Messaging Permission

#### 1.1 Update User Model - `CanMessageUser()` method
**File**: `app/models/user.go`

Change from:
```go
func (u *User) CanMessageUser(other *User) bool {
    if u.IsSuperAdminUser() {
        return true
    }
    return u.SharesRoleWith(other)
}
```

To:
```go
func (u *User) CanMessageUser(other *User) bool {
    // Super admins can message anyone
    if u.IsSuperAdminUser() {
        return true
    }

    // Get sender's highest role level
    senderRole := u.GetHighestRole()
    if senderRole == nil {
        return false // No role = can't message
    }

    // Get recipient's highest role level
    recipientRole := other.GetHighestRole()
    if recipientRole == nil {
        return true // No role = anyone can message them
    }

    // Sender can message if their level >= recipient's level
    return senderRole.Level >= recipientRole.Level
}
```

#### 1.2 Update `GetMessagableUsers()` in MessageController
**File**: `app/http/controllers/messages/message_controller.go`

Update the query to filter users by role level:
```go
// Get current user's role level
userRole := user.GetHighestRole()
userLevel := 0
if userRole != nil {
    userLevel = userRole.Level
}

// Query users at same level or lower
var users []models.User
query := facades.Orm().Query().
    Model(&models.User{}).
    Where("is_active = ?", true).
    Where("id != ?", user.ID).
    With("Roles")

if !user.IsSuperAdminUser() {
    // Join to filter by role level
    query = query.
        Joins("LEFT JOIN user_roles ur ON ur.user_id = users.id AND ur.is_active = true").
        Joins("LEFT JOIN roles r ON r.id = ur.role_id AND r.is_active = true").
        Where("r.level IS NULL OR r.level <= ?", userLevel).
        Distinct()
}

query.Find(&users)
```

#### 1.3 Add Helper Method for Role Level
**File**: `app/models/user.go`

```go
// GetRoleLevel returns the user's highest role level, or 0 if no roles
func (u *User) GetRoleLevel() int {
    role := u.GetHighestRole()
    if role == nil {
        return 0
    }
    return role.Level
}
```

---

### Phase 2: Register Missing Routes

**File**: `routes/api.go`

Add the following routes in the messages group:

```go
// Message routes (add these missing ones)
messagesGroup.Get("/inbox", messageController.GetInbox)
messagesGroup.Get("/sent", messageController.GetSentMessages)
messagesGroup.Get("/unread", messageController.GetUnreadMessages)
messagesGroup.Post("/broadcast", messageController.SendBroadcast)
```

---

### Phase 3: Automatic Notification on Message Receipt

#### 3.1 Update MessageService.SendMessage()
**File**: `app/services/message_service.go`

After successfully creating the message, create a notification:

```go
// After message creation succeeds...

// Create notification for recipient
notificationService := NewNotificationService()
notificationData := map[string]interface{}{
    "title":           fmt.Sprintf("New message from %s", sender.Name),
    "message":         truncateString(content, 100), // First 100 chars
    "type":            "message",
    "user_id":         recipientID,
    "trigger_user_id": senderID,
    "related_type":    "message",
    "related_id":      message.ID,
    "priority":        "normal",
}
notificationService.CreateNotification(notificationData)
```

#### 3.2 Add truncateString helper
**File**: `app/services/message_service.go`

```go
func truncateString(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen-3] + "..."
}
```

---

### Phase 4: Update Frontend for Role-Based User Discovery

#### 4.1 Update MessageContext to show role info
**File**: `resources/js/contexts/MessageContext.tsx`

The `loadMessagableUsers()` function already works correctly - it calls the backend which will now filter by role level.

#### 4.2 Display role level indication in UI (Optional Enhancement)
**File**: `resources/js/components/Messages/MessageSidebar.tsx`

Consider showing a visual indicator of role level in the user list to help users understand who they can message.

---

### Phase 5: Testing

#### 5.1 Backend Tests
Create test file: `tests/feature/messages/message_permission_test.go`

Test cases:
1. Super admin can message any user
2. Admin (level 80) can message Moderator (level 40)
3. Admin (level 80) can message another Admin (level 80)
4. Member (level 20) CANNOT message Admin (level 80)
5. Member (level 20) can message Guest (level 10)
6. User with no role can message user with no role
7. GetMessagableUsers returns only users at same level or lower

#### 5.2 Frontend Testing
- Verify user list only shows messageable users
- Verify sending message to higher-level user fails gracefully
- Verify notifications appear when message received

---

## File Changes Summary

| File | Change Type | Description |
|------|-------------|-------------|
| `app/models/user.go` | Modify | Update `CanMessageUser()`, add `GetRoleLevel()` |
| `app/http/controllers/messages/message_controller.go` | Modify | Update `GetMessagableUsers()` query |
| `routes/api.go` | Modify | Register missing message routes |
| `app/services/message_service.go` | Modify | Add notification creation in `SendMessage()` |
| `tests/feature/messages/message_permission_test.go` | Create | New test file for permission tests |

---

## Execution Order

1. **Phase 1**: Role-based permission fix (most critical)
2. **Phase 2**: Register missing routes
3. **Phase 3**: Automatic notifications
4. **Phase 4**: Frontend updates (if needed)
5. **Phase 5**: Testing

---

## Risks and Considerations

1. **Performance**: Role level query with JOINs may need indexing
2. **Edge Cases**: Users with no roles - decide if they can message anyone or no one
3. **Backwards Compatibility**: Existing messages remain accessible regardless of role changes
4. **Notification Volume**: Consider rate limiting for broadcast messages to prevent notification spam
