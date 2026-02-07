package models

import (
	"encoding/json"
	"time"
)

type Author struct {
	BaseAuditableModel

	FirstName   string     `json:"firstName" gorm:"column:first_name;not null"`
	LastName    string     `json:"lastName" gorm:"column:last_name;not null"`
	Bio         *string    `json:"bio" gorm:"type:text"`
	Email       *string    `json:"email" gorm:"unique"`
	Website     *string    `json:"website"`
	BirthDate   *time.Time `json:"birthDate" gorm:"column:birth_date"`
	Nationality *string    `json:"nationality"`
	PhotoURL    *string    `json:"photoUrl" gorm:"column:photo_url"`
	Status      string     `json:"status" gorm:"default:'ACTIVE'"`

	// Relationships
	Books []Book `json:"books,omitempty" gorm:"foreignKey:AuthorID"`
}

// TableName returns the table name for this model
func (a Author) TableName() string {
	return "authors"
}

// SearchFields returns the fields that can be searched
func (a Author) SearchFields() []string {
	return []string{"first_name", "last_name", "email", "nationality"}
}

// FullName returns the author's full name
func (a Author) FullName() string {
	return a.FirstName + " " + a.LastName
}

// MarshalJSON custom JSON marshaling to handle date formatting
func (a Author) MarshalJSON() ([]byte, error) {
	type Alias Author
	return json.Marshal(&struct {
		*Alias
		CreatedAt *time.Time `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		BirthDate *time.Time `json:"birthDate,omitempty"`
	}{
		Alias:     (*Alias)(&a),
		CreatedAt: timeToPtr(a.CreatedAt.StdTime()),
		UpdatedAt: timeToPtr(a.UpdatedAt.StdTime()),
		BirthDate: a.BirthDate,
	})
}
