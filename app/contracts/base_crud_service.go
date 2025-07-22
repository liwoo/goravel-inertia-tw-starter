package contracts

import (
	"errors"
	"fmt"
	"strings"
	
	"github.com/goravel/framework/contracts/database/orm"
)

// BaseCrudService provides common implementations for CRUD services
// Services MUST embed this and implement the abstract methods
type BaseCrudService struct {
	tableName         string
	primaryKey        string
	maxPageSize       int
	defaultPageSize   int
	searchBuilder     *SearchBuilder
	sortBuilder       *SortBuilder
	paginationBuilder *PaginationBuilder
}

// NewBaseCrudService creates a new base CRUD service
func NewBaseCrudService(tableName, primaryKey string) *BaseCrudService {
	return &BaseCrudService{
		tableName:         tableName,
		primaryKey:        primaryKey,
		maxPageSize:       100,
		defaultPageSize:   20,
		searchBuilder:     NewSearchBuilder(),
		sortBuilder:       NewSortBuilder(),
		paginationBuilder: NewPaginationBuilder(),
	}
}

// PAGINATION CONTRACT IMPLEMENTATION (enforced)

func (b *BaseCrudService) ValidatePaginationParams(page, pageSize int) error {
	if page <= 0 {
		return errors.New("page must be greater than 0")
	}
	if pageSize <= 0 {
		return errors.New("pageSize must be greater than 0")
	}
	if pageSize > b.maxPageSize {
		return fmt.Errorf("pageSize cannot exceed %d", b.maxPageSize)
	}
	return nil
}

func (b *BaseCrudService) GetMaxPageSize() int {
	return b.maxPageSize
}

func (b *BaseCrudService) GetDefaultPageSize() int {
	return b.defaultPageSize
}

func (b *BaseCrudService) SetMaxPageSize(size int) {
	if size > 0 {
		b.maxPageSize = size
	}
}

func (b *BaseCrudService) SetDefaultPageSize(size int) {
	if size > 0 && size <= b.maxPageSize {
		b.defaultPageSize = size
	}
}

// SORTING CONTRACT IMPLEMENTATION (enforced)

func (b *BaseCrudService) ValidateSortDirection(direction string) bool {
	upper := strings.ToUpper(direction)
	result := upper == "ASC" || upper == "DESC"
	return result
}

func (b *BaseCrudService) GetDefaultSort() (field string, direction string) {
	return b.primaryKey, "DESC"
}

// FILTERING CONTRACT IMPLEMENTATION (enforced)

func (b *BaseCrudService) ValidateFilterValue(field string, value interface{}) bool {
	// Basic validation - can be overridden by specific services
	if value == nil {
		return false
	}
	
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) != ""
	case int, int32, int64, uint, uint32, uint64:
		return true
	case float32, float64:
		return true
	case bool:
		return true
	default:
		return false
	}
}

// CONFIGURATION IMPLEMENTATION

func (b *BaseCrudService) GetTableName() string {
	return b.tableName
}

func (b *BaseCrudService) GetPrimaryKey() string {
	return b.primaryKey
}

// VALIDATION HELPERS

func (b *BaseCrudService) ValidateListRequest(req *ListRequest) error {
	// Set defaults first
	req.SetDefaults()
	
	// Validate pagination
	if err := b.ValidatePaginationParams(req.Page, req.PageSize); err != nil {
		return fmt.Errorf("pagination validation failed: %w", err)
	}
	
	// Validate sort direction if provided
	if req.Direction != "" && !b.ValidateSortDirection(req.Direction) {
		return fmt.Errorf("invalid sort direction: %s", req.Direction)
	}
	
	return nil
}

func (b *BaseCrudService) SanitizeListRequest(req *ListRequest) {
	// Ensure page is at least 1
	if req.Page <= 0 {
		req.Page = 1
	}
	
	// Ensure pageSize is within bounds
	if req.PageSize <= 0 {
		req.PageSize = b.defaultPageSize
	}
	if req.PageSize > b.maxPageSize {
		req.PageSize = b.maxPageSize
	}
	
	// Normalize sort direction
	if req.Direction != "" {
		req.Direction = strings.ToUpper(req.Direction)
		if req.Direction != "ASC" && req.Direction != "DESC" {
			req.Direction = "DESC"
		}
	}
	
	// Trim search query
	req.Search = strings.TrimSpace(req.Search)
}

// BULK OPERATIONS VALIDATION

func (b *BaseCrudService) ValidateBulkOperation(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("no IDs provided for bulk operation")
	}
	
	if len(ids) > 1000 { // Prevent massive bulk operations
		return errors.New("bulk operation cannot exceed 1000 items")
	}
	
	// Check for duplicates
	seen := make(map[uint]bool)
	for _, id := range ids {
		if id == 0 {
			return errors.New("invalid ID (0) in bulk operation")
		}
		if seen[id] {
			return fmt.Errorf("duplicate ID %d in bulk operation", id)
		}
		seen[id] = true
	}
	
	return nil
}

// SEARCH CONTRACT IMPLEMENTATION (enforced)

func (b *BaseCrudService) ValidateSearchQuery(query string) error {
	// Trim whitespace
	query = strings.TrimSpace(query)
	
	// Check minimum length
	if len(query) < 2 {
		return errors.New("search query must be at least 2 characters long")
	}
	
	// Check maximum length
	if len(query) > 200 {
		return errors.New("search query cannot exceed 200 characters")
	}
	
	// Check for SQL injection patterns (basic)
	lowerQuery := strings.ToLower(query)
	dangerousPatterns := []string{
		"drop ", "delete ", "insert ", "update ", "alter ", "create ",
		"truncate ", "exec ", "execute ", "--", "/*", "*/", "xp_", "sp_",
	}
	
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerQuery, pattern) {
			return fmt.Errorf("search query contains invalid pattern: %s", pattern)
		}
	}
	
	return nil
}

func (b *BaseCrudService) BuildSearchQuery(query string, searchableFields []string) string {
	// Escape special characters for LIKE queries
	query = strings.ReplaceAll(query, "%", "\\%")
	query = strings.ReplaceAll(query, "_", "\\_")
	
	// Add wildcards for fuzzy search
	return "%" + query + "%"
}

// ApplySearch applies search conditions to a query using the service's searchable fields
func (b *BaseCrudService) ApplySearch(query orm.Query, search string, service interface{}) (orm.Query, error) {
	return b.searchBuilder.ApplySearchWithService(query, search, service)
}

// ApplySort applies sorting to a query using the service's sort configuration
func (b *BaseCrudService) ApplySort(query orm.Query, sort, direction string, service interface{}) (orm.Query, error) {
	return b.sortBuilder.ApplySortWithService(query, sort, direction, service)
}

// PaginateResults performs manual pagination on a slice of items
func (b *BaseCrudService) PaginateResults(items interface{}, req ListRequest) *PaginatedResult {
	// This will be overridden by specific implementations that know the concrete type
	return b.paginationBuilder.PaginateSlice(items, req.Page, req.PageSize)
}

// GetPaginationBuilder returns the pagination builder for advanced usage
func (b *BaseCrudService) GetPaginationBuilder() *PaginationBuilder {
	return b.paginationBuilder
}

func (b *BaseCrudService) ValidateSearchRequest(req *SearchRequest) error {
	// Validate search query
	if err := b.ValidateSearchQuery(req.Query); err != nil {
		return fmt.Errorf("search query validation failed: %w", err)
	}
	
	// Validate pagination
	if err := b.ValidatePaginationParams(req.Page, req.PageSize); err != nil {
		return fmt.Errorf("pagination validation failed: %w", err)
	}
	
	// Validate sort direction if provided
	if req.Direction != "" && !b.ValidateSortDirection(req.Direction) {
		return fmt.Errorf("invalid sort direction: %s", req.Direction)
	}
	
	return nil
}

// METADATA GENERATION

func (b *BaseCrudService) GenerateMetadata(name, version string, service CompleteCrudService) ServiceMetadata {
	return ServiceMetadata{
		Name:             name,
		Version:          version,
		SupportedOps:     []string{"CREATE", "READ", "UPDATE", "DELETE", "LIST", "SEARCH", "BULK"},
		SortableFields:   service.GetSortableFields(),
		FilterableFields: []string{}, // TODO: Add GetFilterableFields to CompleteCrudService interface
		SearchableFields: service.GetSearchableFields(),
		MaxPageSize:      b.maxPageSize,
		DefaultPageSize:  b.defaultPageSize,
	}
}

// SERVICE VALIDATION

func ValidateServiceImplementation(service interface{}) ServiceValidationResult {
	result := ServiceValidationResult{
		Valid:   true,
		Errors:  []string{},
		Missing: []string{},
	}
	
	// Check if service implements CompleteCrudService
	if _, ok := service.(CompleteCrudService); !ok {
		result.Valid = false
		result.Errors = append(result.Errors, "service does not implement CompleteCrudService interface")
	}
	
	// Check if service implements individual contracts
	if _, ok := service.(CrudServiceContract); !ok {
		result.Valid = false
		result.Missing = append(result.Missing, "CrudServiceContract")
	}
	
	if _, ok := service.(PaginationServiceContract); !ok {
		result.Valid = false
		result.Missing = append(result.Missing, "PaginationServiceContract")
	}
	
	if _, ok := service.(SortableServiceContract); !ok {
		result.Valid = false
		result.Missing = append(result.Missing, "SortableServiceContract")
	}
	
	if _, ok := service.(FilterableServiceContract); !ok {
		result.Valid = false
		result.Missing = append(result.Missing, "FilterableServiceContract")
	}
	
	return result
}