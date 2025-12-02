package models

import (
	"encoding/json"
	"time"

	"github.com/goravel/framework/database/orm"
)

// Activity type constants
const (
	ActivityLogin          = "login"
	ActivityLogout         = "logout"
	ActivityProfileUpdate  = "profile_update"
	ActivityPasswordChange = "password_change"
	ActivityEmailChange    = "email_change"
	ActivityRoleChange     = "role_change"
	ActivityAccountCreated = "account_created"
)

// UserActivity tracks user activities for audit purposes
type UserActivity struct {
	orm.Model
	orm.SoftDeletes

	UserID       uint   `gorm:"not null;index" json:"user_id"`
	ActivityType string `gorm:"type:varchar(50);not null;index" json:"activity_type"`
	Description  string `gorm:"type:varchar(255);not null" json:"description"`
	Metadata     string `gorm:"type:text" json:"metadata,omitempty"`
	IPAddress    string `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent    string `gorm:"type:text" json:"user_agent,omitempty"`
	RelatedType  string `gorm:"type:varchar(50)" json:"related_type,omitempty"`
	RelatedID    *uint  `json:"related_id,omitempty"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name for UserActivity model
func (UserActivity) TableName() string {
	return "user_activities"
}

// SearchFields returns the fields that can be searched
func (a UserActivity) SearchFields() []string {
	return []string{"activity_type", "description"}
}

// SetMetadata sets the metadata from a map
func (a *UserActivity) SetMetadata(data map[string]interface{}) error {
	if data == nil {
		a.Metadata = ""
		return nil
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	a.Metadata = string(jsonBytes)
	return nil
}

// GetMetadata parses metadata JSON into a map
func (a *UserActivity) GetMetadata() (map[string]interface{}, error) {
	if a.Metadata == "" {
		return nil, nil
	}

	var result map[string]interface{}
	err := json.Unmarshal([]byte(a.Metadata), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// MarshalJSON custom JSON marshaling to handle date formatting
func (a UserActivity) MarshalJSON() ([]byte, error) {
	type Alias UserActivity

	// Parse metadata JSON string to object for API responses
	var metadataObj interface{}
	if a.Metadata != "" {
		json.Unmarshal([]byte(a.Metadata), &metadataObj)
	}

	return json.Marshal(&struct {
		*Alias
		CreatedAt   *time.Time  `json:"createdAt,omitempty"`
		UpdatedAt   *time.Time  `json:"updatedAt,omitempty"`
		MetadataObj interface{} `json:"metadataObj,omitempty"`
	}{
		Alias:       (*Alias)(&a),
		CreatedAt:   activityTimeToPtr(a.CreatedAt.StdTime()),
		UpdatedAt:   activityTimeToPtr(a.UpdatedAt.StdTime()),
		MetadataObj: metadataObj,
	})
}

// activityTimeToPtr helper function to convert time.Time to *time.Time
func activityTimeToPtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}
