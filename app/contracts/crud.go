package contracts

// CrudConfig defines configuration for CRUD operations
type CrudConfig struct {
	// Database configuration
	TableName   string
	PrimaryKey  string
	
	// Search configuration
	SearchFields []string
	DefaultSort  string
	
	// Validation
	CreateValidation string
	UpdateValidation string
	
	// Pagination
	DefaultPageSize int
	MaxPageSize     int
	
	// Timestamps
	Timestamps bool
	
	// Soft deletes
	SoftDeletes bool
	SoftDelete  bool // Alias for backward compatibility
	
	// Relationships to load
	With []string
}

