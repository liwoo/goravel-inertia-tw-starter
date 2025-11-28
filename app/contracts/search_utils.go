package contracts

import (
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/database/orm"
)

// SearchConfig defines the configuration for search operations
type SearchConfig struct {
	Query            string                   // The search query
	SearchableFields []string                 // Fields to search in
	ValidateFunc     func(query string) error // Optional validation function
}

// SearchBuilder provides utilities for building search queries
type SearchBuilder struct{}

// NewSearchBuilder creates a new search builder instance
func NewSearchBuilder() *SearchBuilder {
	return &SearchBuilder{}
}

// ApplySearch applies search conditions to a query
func (sb *SearchBuilder) ApplySearch(query orm.Query, config SearchConfig) (orm.Query, error) {
	// Skip if no search query
	if config.Query == "" {
		return query, nil
	}

	// Validate search query if validation function is provided
	if config.ValidateFunc != nil {
		if err := config.ValidateFunc(config.Query); err != nil {
			return query, fmt.Errorf("invalid search query: %w", err)
		}
	}

	// Skip if no searchable fields
	if len(config.SearchableFields) == 0 {
		return query, nil
	}

	// Build search conditions
	conditions := make([]string, len(config.SearchableFields))
	values := make([]interface{}, len(config.SearchableFields))
	searchPattern := "%" + config.Query + "%"

	for i, field := range config.SearchableFields {
		// Use ILIKE for case-insensitive search in PostgreSQL
		conditions[i] = field + " ILIKE ?"
		values[i] = searchPattern
	}

	// Apply search conditions with OR logic
	searchCondition := strings.Join(conditions, " OR ")
	return query.Where(searchCondition, values...), nil
}

// ApplySearchWithService applies search using service methods
func (sb *SearchBuilder) ApplySearchWithService(query orm.Query, search string, service interface{}) (orm.Query, error) {
	// Check if service implements searchable interface
	searchable, ok := service.(interface {
		ValidateSearchQuery(string) error
		GetSearchableFields() []string
	})

	if !ok {
		return query, fmt.Errorf("service does not implement searchable interface")
	}

	config := SearchConfig{
		Query:            search,
		SearchableFields: searchable.GetSearchableFields(),
		ValidateFunc:     searchable.ValidateSearchQuery,
	}

	return sb.ApplySearch(query, config)
}

// BuildSearchCondition builds a search condition string with placeholders
func (sb *SearchBuilder) BuildSearchCondition(fields []string) (condition string, valueCount int) {
	if len(fields) == 0 {
		return "", 0
	}

	conditions := make([]string, len(fields))
	for i, field := range fields {
		// Use ILIKE for case-insensitive search in PostgreSQL
		conditions[i] = field + " ILIKE ?"
	}

	return strings.Join(conditions, " OR "), len(fields)
}

// GetSearchValues generates the values array for a search query
func (sb *SearchBuilder) GetSearchValues(query string, fieldCount int) []interface{} {
	searchPattern := "%" + query + "%"
	values := make([]interface{}, fieldCount)
	for i := 0; i < fieldCount; i++ {
		values[i] = searchPattern
	}
	return values
}

// CombineSearchFields combines multiple field sets for searching
func (sb *SearchBuilder) CombineSearchFields(fieldSets ...[]string) []string {
	uniqueFields := make(map[string]bool)
	for _, fields := range fieldSets {
		for _, field := range fields {
			uniqueFields[field] = true
		}
	}

	result := make([]string, 0, len(uniqueFields))
	for field := range uniqueFields {
		result = append(result, field)
	}
	return result
}

// ValidateSearchFields validates that all search fields are in the allowed list
func (sb *SearchBuilder) ValidateSearchFields(searchFields, allowedFields []string) error {
	allowedMap := make(map[string]bool)
	for _, field := range allowedFields {
		allowedMap[field] = true
	}

	for _, field := range searchFields {
		if !allowedMap[field] {
			return fmt.Errorf("invalid search field: %s", field)
		}
	}

	return nil
}
