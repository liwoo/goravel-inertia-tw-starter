package contracts

import (
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/database/orm"
)

// SortConfig defines the configuration for sort operations
type SortConfig struct {
	Field           string // The field to sort by
	Direction       string // The sort direction (ASC/DESC)
	ValidateFunc    func(field string) bool // Field validation function
	MapFieldFunc    func(field string) (string, bool) // Field mapping function
	GetDefaultFunc  func() (string, string) // Default sort function
}

// SortBuilder provides utilities for building sort queries
type SortBuilder struct{}

// NewSortBuilder creates a new sort builder instance
func NewSortBuilder() *SortBuilder {
	return &SortBuilder{}
}

// ApplySort applies sorting to a query with validation and mapping
func (sb *SortBuilder) ApplySort(query orm.Query, config SortConfig) orm.Query {
	// Determine the field and direction to use
	field, direction := sb.determineSortParams(config)
	
	// Build and apply the order clause
	orderClause := field + " " + strings.ToUpper(direction)
	return query.Order(orderClause)
}

// ApplySortWithService applies sorting using service methods
func (sb *SortBuilder) ApplySortWithService(query orm.Query, sort, direction string, service interface{}) (orm.Query, error) {
	// Check if service implements sortable interface
	sortable, ok := service.(interface {
		ValidateSortField(string) bool
		ValidateSortDirection(string) bool
		MapSortField(string) (string, bool)
		GetDefaultSort() (string, string)
	})
	
	if !ok {
		return query, fmt.Errorf("service does not implement sortable interface")
	}

	config := SortConfig{
		Field:          sort,
		Direction:      direction,
		ValidateFunc:   sortable.ValidateSortField,
		MapFieldFunc:   sortable.MapSortField,
		GetDefaultFunc: sortable.GetDefaultSort,
	}
	
	// Validate direction separately
	if direction != "" && !sortable.ValidateSortDirection(direction) {
		// Use default if direction is invalid
		defaultField, defaultDir := sortable.GetDefaultSort()
		config.Field = defaultField
		config.Direction = defaultDir
	}

	return sb.ApplySort(query, config), nil
}

// determineSortParams determines the final sort field and direction
func (sb *SortBuilder) determineSortParams(config SortConfig) (field, direction string) {
	// If no sort field specified, use default
	if config.Field == "" || config.Direction == "" {
		if config.GetDefaultFunc != nil {
			return config.GetDefaultFunc()
		}
		return "id", "DESC" // Ultimate fallback
	}

	// Validate the sort field if validation function provided
	if config.ValidateFunc != nil && !config.ValidateFunc(config.Field) {
		// Invalid field, use default
		if config.GetDefaultFunc != nil {
			return config.GetDefaultFunc()
		}
		return "id", "DESC"
	}

	// Map the field if mapping function provided
	if config.MapFieldFunc != nil {
		mappedField, valid := config.MapFieldFunc(config.Field)
		if valid {
			return mappedField, config.Direction
		}
		// Invalid mapping, use default
		if config.GetDefaultFunc != nil {
			return config.GetDefaultFunc()
		}
		return "id", "DESC"
	}

	// Use the field as-is
	return config.Field, config.Direction
}

// BuildSortClause builds a sort clause string
func (sb *SortBuilder) BuildSortClause(field, direction string) string {
	if field == "" {
		return ""
	}
	
	direction = strings.ToUpper(direction)
	if direction != "ASC" && direction != "DESC" {
		direction = "DESC"
	}
	
	return field + " " + direction
}

// ValidateSortDirection validates sort direction
func (sb *SortBuilder) ValidateSortDirection(direction string) bool {
	normalized := strings.ToUpper(strings.TrimSpace(direction))
	return normalized == "ASC" || normalized == "DESC"
}

// NormalizeSortDirection normalizes sort direction
func (sb *SortBuilder) NormalizeSortDirection(direction string) string {
	normalized := strings.ToUpper(strings.TrimSpace(direction))
	if normalized != "ASC" && normalized != "DESC" {
		return "DESC"
	}
	return normalized
}

// CombineSortClauses combines multiple sort clauses
func (sb *SortBuilder) CombineSortClauses(clauses ...string) string {
	validClauses := make([]string, 0, len(clauses))
	for _, clause := range clauses {
		if clause != "" {
			validClauses = append(validClauses, clause)
		}
	}
	return strings.Join(validClauses, ", ")
}

// ParseSortString parses a sort string like "name:asc,created_at:desc"
func (sb *SortBuilder) ParseSortString(sortString string) []struct{ Field, Direction string } {
	var result []struct{ Field, Direction string }
	
	if sortString == "" {
		return result
	}
	
	parts := strings.Split(sortString, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		fieldDir := strings.Split(part, ":")
		if len(fieldDir) == 2 {
			result = append(result, struct{ Field, Direction string }{
				Field:     strings.TrimSpace(fieldDir[0]),
				Direction: sb.NormalizeSortDirection(fieldDir[1]),
			})
		} else if len(fieldDir) == 1 {
			// Default to DESC if no direction specified
			result = append(result, struct{ Field, Direction string }{
				Field:     strings.TrimSpace(fieldDir[0]),
				Direction: "DESC",
			})
		}
	}
	
	return result
}