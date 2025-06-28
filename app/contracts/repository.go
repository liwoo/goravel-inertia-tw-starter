package contracts

// Repository defines the core data access contract
type Repository interface {
	// Basic CRUD operations
	Find(id interface{}) (interface{}, error)
	FindMany(ids []interface{}) ([]interface{}, error)
	Create(data map[string]interface{}) (interface{}, error)
	Update(id interface{}, data map[string]interface{}) (interface{}, error)
	Delete(id interface{}) error

	// Query operations
	Count(conditions map[string]interface{}) (int64, error)
	Exists(id interface{}) bool

	// Bulk operations
	BulkCreate(data []map[string]interface{}) ([]interface{}, error)
	BulkUpdate(conditions map[string]interface{}, data map[string]interface{}) error
	BulkDelete(conditions map[string]interface{}) error
}

// QueryableRepository extends Repository with advanced querying
type QueryableRepository interface {
	Repository

	// Query building
	Where(field string, operator string, value interface{}) QueryableRepository
	WhereIn(field string, values []interface{}) QueryableRepository
	WhereBetween(field string, min, max interface{}) QueryableRepository
	WhereNull(field string) QueryableRepository
	WhereNotNull(field string) QueryableRepository

	// Relationships
	With(relations ...string) QueryableRepository

	// Ordering and limiting
	OrderBy(field string, direction string) QueryableRepository
	Limit(limit int) QueryableRepository
	Offset(offset int) QueryableRepository

	// Execution
	Get() ([]interface{}, error)
	First() (interface{}, error)
	Paginate(page, pageSize int) (*PaginatedResult, error)

	// Reset query builder
	Reset() QueryableRepository
	Clone() QueryableRepository
}

// SearchableRepository adds search capabilities
type SearchableRepository interface {
	QueryableRepository

	// Search operations
	Search(term string, fields []string) SearchableRepository
	SearchWithFilters(term string, fields []string, filters map[string]interface{}) SearchableRepository
}

// AuditableRepository handles soft deletes and audit trails
type AuditableRepository interface {
	QueryableRepository

	// Soft delete operations
	SoftDelete(id interface{}) error
	Restore(id interface{}) error
	ForceDelete(id interface{}) error
	OnlyTrashed() AuditableRepository
	WithTrashed() AuditableRepository
}