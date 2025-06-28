package contracts

import (
	"errors"
	"fmt"
	"strings"
)

// BaseCrudService provides common implementations for CRUD services
// Services MUST embed this and implement the abstract methods
type BaseCrudService struct {
	tableName       string
	primaryKey      string
	maxPageSize     int
	defaultPageSize int
}

// NewBaseCrudService creates a new base CRUD service
func NewBaseCrudService(tableName, primaryKey string) *BaseCrudService {
	return &BaseCrudService{
		tableName:       tableName,
		primaryKey:      primaryKey,
		maxPageSize:     100,
		defaultPageSize: 20,
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
	return upper == "ASC" || upper == "DESC"
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

// METADATA GENERATION

func (b *BaseCrudService) GenerateMetadata(name, version string, service CompleteCrudService) ServiceMetadata {
	return ServiceMetadata{
		Name:             name,
		Version:          version,
		SupportedOps:     []string{"CREATE", "READ", "UPDATE", "DELETE", "LIST", "SEARCH", "BULK"},
		SortableFields:   service.GetSortableFields(),
		FilterableFields: service.GetFilterableFields(),
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