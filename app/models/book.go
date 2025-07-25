package models

import (
	"encoding/json"
	"time"
	"gorm.io/gorm"
)

// Book entity - now using BaseAuditableModel for consistent audit fields
type Book struct {
	BaseAuditableModel
	
	Title       string    `json:"title" gorm:"not null"`
	Author      string    `json:"author" gorm:"not null"`
	ISBN        string    `json:"isbn" gorm:"unique;not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"default:0"`
	Status      string    `json:"status" gorm:"default:'AVAILABLE'"` // AVAILABLE, BORROWED, MAINTENANCE, RESERVED
	PublishedAt *time.Time `json:"publishedAt" gorm:"column:published_at"`
	TagsJSON    string    `json:"-" gorm:"column:tags;type:text"` // Store as JSON string in database
	Tags        []string  `json:"tags" gorm:"-"` // Virtual field for API
}

// AfterFind hook to convert tags JSON to array
func (b *Book) AfterFind(tx *gorm.DB) error {
	if b.TagsJSON != "" {
		err := json.Unmarshal([]byte(b.TagsJSON), &b.Tags)
		if err != nil {
			// If unmarshal fails, treat as empty array
			b.Tags = []string{}
		}
	} else {
		b.Tags = []string{}
	}
	return nil
}

// SearchFields returns the fields that can be searched
func (b Book) SearchFields() []string {
	return []string{"title", "author", "isbn", "description"}
}

// TableName returns the table name for this model
func (b Book) TableName() string {
	return "books"
}

// MarshalJSON custom JSON marshaling to handle date formatting
func (b Book) MarshalJSON() ([]byte, error) {
	type Alias Book
	return json.Marshal(&struct {
		*Alias
		CreatedAt   *time.Time `json:"createdAt,omitempty"`
		UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
		PublishedAt *time.Time `json:"publishedAt,omitempty"`
	}{
		Alias:       (*Alias)(&b),
		CreatedAt:   timeToPtr(b.CreatedAt.StdTime()),
		UpdatedAt:   timeToPtr(b.UpdatedAt.StdTime()),
		PublishedAt: b.PublishedAt,
	})
}

// Helper function to convert time.Time to *time.Time
func timeToPtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}