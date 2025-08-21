package models

import (
	"github.com/goravel/framework/database/orm"
	"time"
)

// MessageType represents the type of message
type MessageType string

const (
	MessageTypeDirect MessageType = "direct"
	MessageTypeGroup  MessageType = "group"
	MessageTypeSystem MessageType = "system"
)

// MessageStatus represents the status of a message
type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusDeleted   MessageStatus = "deleted"
)

// Message represents a message in the system
type Message struct {
	orm.Model

	// Core message fields
	Content string        `gorm:"type:text;not null" json:"content"`
	Type    MessageType   `gorm:"type:varchar(20);default:'direct';index" json:"type"`
	Status  MessageStatus `gorm:"type:varchar(20);default:'sent';index" json:"status"`

	// User relationships
	SenderID    uint  `gorm:"not null;index" json:"sender_id"`
	Sender      User  `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	RecipientID *uint `gorm:"index" json:"recipient_id,omitempty"`
	Recipient   *User `gorm:"foreignKey:RecipientID" json:"recipient,omitempty"`

	// Group messaging (for future use)
	GroupID *uint `gorm:"index" json:"group_id,omitempty"`

	// Message metadata
	IsEdited    bool       `gorm:"default:false" json:"is_edited"`
	EditedAt    *time.Time `json:"edited_at,omitempty"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`

	// Threading support
	ParentMessageID *uint     `gorm:"index" json:"parent_message_id,omitempty"`
	ParentMessage   *Message  `gorm:"foreignKey:ParentMessageID" json:"parent_message,omitempty"`
	Replies         []Message `gorm:"foreignKey:ParentMessageID" json:"replies,omitempty"`

	// Mentions and attachments
	MentionedUserIDs []uint `gorm:"-" json:"mentioned_user_ids,omitempty"`
	HasAttachments   bool   `gorm:"default:false" json:"has_attachments"`

	orm.SoftDeletes
}

// TableName returns the table name for Message model
func (Message) TableName() string {
	return "messages"
}

// MarkAsRead marks the message as read
func (m *Message) MarkAsRead() {
	now := time.Now()
	m.Status = MessageStatusRead
	m.ReadAt = &now
}

// MarkAsDelivered marks the message as delivered
func (m *Message) MarkAsDelivered() {
	now := time.Now()
	m.Status = MessageStatusDelivered
	m.DeliveredAt = &now
}

// IsReadBy checks if the message has been read by a specific user
func (m *Message) IsReadBy(userID uint) bool {
	return m.Status == MessageStatusRead && m.RecipientID != nil && *m.RecipientID == userID
}

// CanBeReadBy checks if a user can read this message based on role permissions
func (m *Message) CanBeReadBy(user *User) bool {
	// Super admins can read any message
	if user.IsSuperAdminUser() {
		return true
	}

	// Users can read messages they sent
	if m.SenderID == user.ID {
		return true
	}

	// Users can read messages sent to them
	if m.RecipientID != nil && *m.RecipientID == user.ID {
		return true
	}

	// For role-based messaging: users can read messages from users with the same roles
	if m.Sender.ID != 0 {
		return user.SharesRoleWith(&m.Sender)
	}

	return false
}

// CanEditMessage checks if a user can edit this message
func (m *Message) CanEditMessage(user *User) bool {
	// Only sender can edit their own messages within 15 minutes
	if m.SenderID != user.ID {
		return false
	}

	// Check if message is within edit window (15 minutes)
	editWindow := 15 * time.Minute
	return time.Since(m.CreatedAt.StdTime()) <= editWindow
}

// MessageMention represents user mentions in messages
type MessageMention struct {
	orm.Model
	MessageID uint    `gorm:"not null;index" json:"message_id"`
	Message   Message `gorm:"foreignKey:MessageID" json:"message,omitempty"`
	UserID    uint    `gorm:"not null;index" json:"user_id"`
	User      User    `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// Mention metadata
	Position int        `gorm:"not null" json:"position"`
	Length   int        `gorm:"not null" json:"length"`
	IsRead   bool       `gorm:"default:false" json:"is_read"`
	ReadAt   *time.Time `json:"read_at,omitempty"`
}

// TableName returns the table name for MessageMention model
func (MessageMention) TableName() string {
	return "message_mentions"
}

// Notification represents a notification in the system
type Notification struct {
	orm.Model

	// Core notification fields
	Title   string `gorm:"not null" json:"title"`
	Message string `gorm:"type:text" json:"message"`
	Type    string `gorm:"type:varchar(50);not null;index" json:"type"`

	// User relationships
	UserID        uint  `gorm:"not null;index" json:"user_id"`
	User          User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TriggerUserID *uint `gorm:"index" json:"trigger_user_id,omitempty"`
	TriggerUser   *User `gorm:"foreignKey:TriggerUserID" json:"trigger_user,omitempty"`

	// Related entities
	RelatedType string `gorm:"type:varchar(50)" json:"related_type,omitempty"`
	RelatedID   *uint  `gorm:"index" json:"related_id,omitempty"`

	// Notification state
	IsRead      bool       `gorm:"default:false;index" json:"is_read"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	IsDismissed bool       `gorm:"default:false" json:"is_dismissed"`
	DismissedAt *time.Time `json:"dismissed_at,omitempty"`

	// Metadata
	Data      string     `gorm:"type:json" json:"data,omitempty"`
	Priority  string     `gorm:"type:varchar(20);default:'normal'" json:"priority"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`

	orm.SoftDeletes
}

// TableName returns the table name for Notification model
func (Notification) TableName() string {
	return "notifications"
}

// MarkAsRead marks the notification as read
func (n *Notification) MarkAsRead() {
	now := time.Now()
	n.IsRead = true
	n.ReadAt = &now
}

// Dismiss marks the notification as dismissed
func (n *Notification) Dismiss() {
	now := time.Now()
	n.IsDismissed = true
	n.DismissedAt = &now
}

// IsExpired checks if the notification has expired
func (n *Notification) IsExpired() bool {
	return n.ExpiresAt != nil && time.Now().After(*n.ExpiresAt)
}
