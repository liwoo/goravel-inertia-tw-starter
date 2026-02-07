---
name: broadcast-notification
description: Create a broadcast notification system that sends notifications and messages to eligible users based on entity-specific criteria. Use when adding a new entity type that needs to notify users when records are created, updated, or published.
argument-hint: [entity-name] [trigger-description]
---

# Broadcast Notification Pattern

Implement a broadcast notification system for a new entity. This follows the established pattern in this Goravel codebase where entity lifecycle events trigger notifications + direct messages to eligible users.

Entity name: $ARGUMENTS

## Architecture Overview

The broadcast notification system has 4 layers:

```
Controller (afterStore/afterUpdate hook)
  -> facades.Event().Job(&events.EntityEvent{}, []event.Arg{...}).Dispatch()
    -> Listener (async via queue, extracts args, calls service)
      -> NotificationService method (finds eligible users, creates notifications + messages)
```

## Step-by-Step Implementation

### 1. Create the Event (`app/events/`)

Follow the pattern from `app/events/application_approved.go`:

```go
package events

import "github.com/goravel/framework/contracts/event"

type EntityPublished struct{}

func (e *EntityPublished) Handle(args []event.Arg) ([]event.Arg, error) {
    return args, nil
}
```

### 2. Create the Listener (`app/listeners/`)

Follow the pattern from `app/listeners/notify_application_approved.go`. Key points:
- Use the existing `toUint()` helper from `app/listeners/helpers.go` for safe type conversion (handles float64 from JSON queue serialization)
- Always check `len(args)` before accessing
- Return `nil` (not error) on invalid args to avoid queue retries
- Instantiate `services.NewNotificationService()` inside Handle

```go
package listeners

import (
    "github.com/goravel/framework/contracts/event"
    "github.com/goravel/framework/facades"
    "smedi-sme-db/app/services"
)

type BroadcastEntityNotification struct{}

func (l *BroadcastEntityNotification) Signature() string {
    return "broadcast_entity_notification"
}

func (l *BroadcastEntityNotification) Queue(args ...any) event.Queue {
    return event.Queue{
        Enable:     true,  // Always async for broadcasts
        Connection: "",
        Queue:      "",
    }
}

func (l *BroadcastEntityNotification) Handle(args ...any) error {
    if len(args) < 3 {
        facades.Log().Warning("BroadcastEntityNotification: Insufficient arguments")
        return nil
    }

    senderID, ok := toUint(args[0])
    if !ok {
        facades.Log().Warning("BroadcastEntityNotification: Invalid senderID type")
        return nil
    }

    entityID, ok := toUint(args[1])
    if !ok {
        facades.Log().Warning("BroadcastEntityNotification: Invalid entityID type")
        return nil
    }

    entityTitle, ok := args[2].(string)
    if !ok {
        facades.Log().Warning("BroadcastEntityNotification: Invalid entityTitle type")
        return nil
    }

    notificationService := services.NewNotificationService()
    notificationService.BroadcastEntityNotificationSync(senderID, entityID, entityTitle)

    return nil
}
```

### 3. Register Event-Listener in `app/providers/event_service_provider.go`

Add the mapping to the `listen()` method:

```go
func (receiver *EventServiceProvider) listen() map[event.Event][]event.Listener {
    return map[event.Event][]event.Listener{
        // ... existing mappings ...
        &events.EntityPublished{}: {
            &listeners.BroadcastEntityNotification{},
        },
    }
}
```

### 4. Add the Broadcast Method to `app/services/notification_service.go`

Follow the dual-delivery pattern (notification + message):

```go
func (s *NotificationService) BroadcastEntityNotificationSync(
    senderID uint, entityID uint, entityTitle string,
) {
    facades.Log().Info("Starting entity broadcast notification", map[string]interface{}{
        "entity_id":    entityID,
        "entity_title": entityTitle,
        "sender_id":    senderID,
    })

    // Step 1: Find eligible users via raw SQL query
    userIDs := s.getEligibleUsersForEntity(/* filtering criteria */)
    if len(userIDs) == 0 {
        facades.Log().Info("No eligible users found for entity notification", map[string]interface{}{
            "entity_id": entityID,
        })
        return
    }

    // Step 2: Build notification content
    title := "New Entity: " + entityTitle
    message := fmt.Sprintf("A new entity \"%s\" has been created.", entityTitle)

    // Step 3: Create notification for each eligible user
    relatedType := "entity"
    for _, userID := range userIDs {
        _, err := s.CreateNotification(
            userID,        // recipient
            title,
            message,
            "entity",      // notification type (used for filtering)
            &senderID,     // trigger user
            &relatedType,
            &entityID,
            "normal",      // priority: "normal" or "high"
            nil,           // expiresAt
            "",            // data (JSON)
        )
        if err != nil {
            facades.Log().Warning("Failed to create entity notification", map[string]interface{}{
                "user_id":   userID,
                "entity_id": entityID,
                "error":     err.Error(),
            })
        }
    }

    // Step 4: Send broadcast messages (richer markdown content)
    messageService := NewMessageService()
    messageContent := fmt.Sprintf("**New Entity**\n\n**%s**\n\nCheck it out!", entityTitle)
    _, err := messageService.SendBroadcast(senderID, userIDs, messageContent, "New Entity: "+entityTitle)
    if err != nil {
        facades.Log().Warning("Failed to send entity broadcast messages", map[string]interface{}{
            "entity_id": entityID,
            "error":     err.Error(),
        })
    }

    facades.Log().Info("Entity broadcast notification completed", map[string]interface{}{
        "entity_id":      entityID,
        "notified_users": len(userIDs),
    })
}
```

### 5. Eligibility Query Pattern

Use raw SQL with `facades.Orm().Query().Raw()` for complex eligibility filtering. Use case-insensitive email matching with `LOWER()`:

```go
func (s *NotificationService) getEligibleUsersForEntity(filterValue string) []uint {
    var userIDs []uint

    // Option A: All active users
    query := `
        SELECT DISTINCT u.id FROM users u
        WHERE u.is_active = true AND u.deleted_at IS NULL
    `

    // Option B: Users matching criteria via JOIN
    query = `
        SELECT DISTINCT u.id FROM users u
        INNER JOIN entity_memberships em ON em.user_id = u.id
        WHERE u.is_active = true AND u.deleted_at IS NULL
        AND em.is_active = true
    `

    // Option C: Users by role
    query = `
        SELECT DISTINCT u.id FROM users u
        INNER JOIN user_roles ur ON ur.user_id = u.id
        INNER JOIN roles r ON ur.role_id = r.id
        WHERE u.is_active = true AND u.deleted_at IS NULL
        AND r.slug = ?
    `

    err := facades.Orm().Query().Raw(query, filterValue).Pluck("id", &userIDs)
    if err != nil {
        facades.Log().Error("Failed to get eligible users", map[string]interface{}{
            "error": err.Error(),
        })
        return []uint{}
    }
    return userIDs
}
```

### 6. Dispatch from Controller

In the controller's `afterStore` or `afterUpdate` hook:

```go
import "smedi-sme-db/app/events"

// In the afterStore callback:
go func() {
    facades.Event().Job(&events.EntityPublished{}, []event.Arg{
        {Type: "uint", Value: currentUserID},
        {Type: "uint", Value: entity.ID},
        {Type: "string", Value: entity.Title},
    }).Dispatch()
}()
```

### 7. Write Tests

Create `tests/feature/notifications/<entity>_broadcast_test.go` following the pattern from `tests/feature/notifications/application_status_notification_test.go`:

- Use `suite.Suite` with `tests.TestCase` embedded
- Create test users in `SetupTest()` with unique timestamps: `fmt.Sprintf("user_%d@example.com", time.Now().UnixNano())`
- Call the notification service method directly (not through the event system)
- Use `time.Sleep(200 * time.Millisecond)` after calling the service before asserting
- Assert notification fields: `Title`, `Message`, `Type`, `RelatedType`, `RelatedID`, `TriggerUserID`, `Priority`
- Assert messages were created via `SendBroadcast`
- Test edge cases: no eligible users, inactive users excluded, case-insensitive email

## Key Services Reference

| Service | Constructor | Key Methods |
|---------|------------|-------------|
| `NotificationService` | `services.NewNotificationService()` | `CreateNotification()`, `CreateSystemNotification()` |
| `MessageService` | `services.NewMessageService()` | `SendMessage()`, `SendBroadcast()`, `BroadcastToRole()` |
| `SSEService` | `services.NewSSEService()` | Auto-called by `CreateNotification()` for real-time push |

## Notification Model Fields

| Field | Type | Description |
|-------|------|-------------|
| `Title` | `string` | Short notification title |
| `Message` | `string` | Notification body text |
| `Type` | `string` | Category for filtering (e.g., "event", "application_approved") |
| `UserID` | `uint` | Recipient user ID |
| `TriggerUserID` | `*uint` | User who caused the notification |
| `RelatedType` | `string` | Entity type (e.g., "application", "event") |
| `RelatedID` | `*uint` | Entity ID for deep linking |
| `Priority` | `string` | "normal" or "high" |
| `ExpiresAt` | `*time.Time` | Optional auto-expiry |
| `Data` | `string` | Optional JSON metadata |

## Checklist

- [ ] Event struct created in `app/events/`
- [ ] Listener created in `app/listeners/` with safe arg extraction
- [ ] Event-listener registered in `app/providers/event_service_provider.go`
- [ ] Broadcast method added to `app/services/notification_service.go`
- [ ] Eligibility query written (raw SQL with `Pluck`)
- [ ] Dual delivery: notification + message
- [ ] Controller hook dispatches event with `go func()` wrapper
- [ ] Tests written covering: happy path, no eligible users, inactive users excluded
- [ ] `go build ./...` passes
