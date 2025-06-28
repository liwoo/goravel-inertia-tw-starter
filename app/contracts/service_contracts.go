package contracts

// CrudServiceContract defines the mandatory interface for all CRUD services
// This contract FORCES implementation of pagination, sorting, and filtering
type CrudServiceContract interface {
	// Core CRUD operations - ALL must be implemented
	GetList(req ListRequest) (*PaginatedResult, error)
	GetListAdvanced(req ListRequest, filters map[string]interface{}) (*PaginatedResult, error)
	GetByID(id uint) (interface{}, error)
	Create(data map[string]interface{}) (interface{}, error)
	Update(id uint, data map[string]interface{}) (interface{}, error)
	Delete(id uint) error
	
	// Pagination contract - MUST be implemented
	PaginationServiceContract
	
	// Sorting contract - MUST be implemented
	SortableServiceContract
	
	// Filtering contract - MUST be implemented
	FilterableServiceContract
}

// PaginationServiceContract enforces pagination functionality
type PaginationServiceContract interface {
	// GetPaginatedList MUST implement proper pagination
	GetPaginatedList(req ListRequest) (*PaginatedResult, error)
	
	// ValidatePaginationParams ensures valid pagination parameters
	ValidatePaginationParams(page, pageSize int) error
	
	// GetMaxPageSize returns the maximum allowed page size
	GetMaxPageSize() int
	
	// GetDefaultPageSize returns the default page size
	GetDefaultPageSize() int
}

// SortableServiceContract enforces sorting functionality
type SortableServiceContract interface {
	// GetSortableFields returns list of fields that can be sorted
	GetSortableFields() []string
	
	// ValidateSortField checks if field is valid for sorting
	ValidateSortField(field string) bool
	
	// ValidateSortDirection checks if direction is valid (ASC/DESC)
	ValidateSortDirection(direction string) bool
	
	// GetDefaultSort returns the default sort configuration
	GetDefaultSort() (field string, direction string)
	
	// MapSortField maps frontend field names to database column names
	MapSortField(frontendField string) (dbColumn string, valid bool)
}

// FilterableServiceContract enforces filtering functionality
type FilterableServiceContract interface {
	// GetFilterableFields returns list of fields that can be filtered
	GetFilterableFields() []string
	
	// ValidateFilterField checks if field is valid for filtering
	ValidateFilterField(field string) bool
	
	// ValidateFilterValue checks if value is valid for the field
	ValidateFilterValue(field string, value interface{}) bool
	
	// GetSearchableFields returns fields that support text search
	GetSearchableFields() []string
	
	// BuildFilterQuery builds the appropriate filter query
	BuildFilterQuery(filters map[string]interface{}) (map[string]interface{}, error)
}

// SearchableServiceContract enforces search functionality
type SearchableServiceContract interface {
	// Search performs full-text search across searchable fields
	Search(query string, req ListRequest) (*PaginatedResult, error)
	
	// GetSearchableFields returns fields that support search
	GetSearchableFields() []string
	
	// ValidateSearchQuery validates the search query
	ValidateSearchQuery(query string) error
}

// BulkOperationsContract enforces bulk operations
type BulkOperationsContract interface {
	// BulkCreate creates multiple records
	BulkCreate(data []map[string]interface{}) ([]interface{}, error)
	
	// BulkUpdate updates multiple records
	BulkUpdate(ids []uint, data map[string]interface{}) error
	
	// BulkDelete deletes multiple records
	BulkDelete(ids []uint) error
	
	// ValidateBulkOperation validates bulk operation parameters
	ValidateBulkOperation(ids []uint) error
}

// CrudServiceConfiguration defines configuration that services must provide
type CrudServiceConfiguration interface {
	// GetTableName returns the primary table name
	GetTableName() string
	
	// GetPrimaryKey returns the primary key field name
	GetPrimaryKey() string
	
	// GetModel returns the model struct for this service
	GetModel() interface{}
	
	// GetValidationRules returns validation rules for create/update
	GetValidationRules() map[string]interface{}
	
	// GetColumnMapping returns frontend->database column mapping
	GetColumnMapping() map[string]string
}

// CompleteCrudService combines all contracts into one interface
// Any service implementing this interface MUST implement ALL CRUD features
type CompleteCrudService interface {
	CrudServiceContract
	SearchableServiceContract
	BulkOperationsContract
	CrudServiceConfiguration
}

// ServiceMetadata provides information about the service capabilities
type ServiceMetadata struct {
	Name            string   `json:"name"`
	Version         string   `json:"version"`
	SupportedOps    []string `json:"supported_operations"`
	SortableFields  []string `json:"sortable_fields"`
	FilterableFields []string `json:"filterable_fields"`
	SearchableFields []string `json:"searchable_fields"`
	MaxPageSize     int      `json:"max_page_size"`
	DefaultPageSize int      `json:"default_page_size"`
}

// ServiceValidationResult represents the result of service validation
type ServiceValidationResult struct {
	Valid   bool     `json:"valid"`
	Errors  []string `json:"errors"`
	Missing []string `json:"missing_methods"`
}