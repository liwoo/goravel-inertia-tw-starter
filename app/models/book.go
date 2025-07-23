package models

// Book entity - now using BaseAuditableModel for consistent audit fields
type Book struct {
	BaseAuditableModel
	
	Title       string    `json:"title" gorm:"not null"`
	Author      string    `json:"author" gorm:"not null"`
	ISBN        string    `json:"isbn" gorm:"unique;not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"default:0"`
	Status      string    `json:"status" gorm:"default:'AVAILABLE'"` // AVAILABLE, BORROWED, MAINTENANCE
	PublishedAt string    `json:"publishedAt" gorm:"column:published_at"`
	Tags        string    `json:"tags" gorm:"-"` // Ignore this field in database operations for now
}

// SearchFields returns the fields that can be searched
func (b Book) SearchFields() []string {
	return []string{"title", "author", "isbn", "description"}
}

// TableName returns the table name for this model
func (b Book) TableName() string {
	return "books"
}