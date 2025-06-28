package helpers

import (
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/facades"
)

// FilterBuilder for complex queries
type FilterBuilder struct {
	query orm.Query
}

// NewFilterBuilder creates a new filter builder
func NewFilterBuilder(tableName string) *FilterBuilder {
	return &FilterBuilder{
		query: facades.Orm().Query().Table(tableName),
	}
}

// Search adds search across multiple fields
func (f *FilterBuilder) Search(fields []string, term string) *FilterBuilder {
	if term == "" || len(fields) == 0 {
		return f
	}

	searchTerm := "%" + term + "%"
	
	// Build OR conditions for search
	f.query = f.query.Where(func(query orm.Query) orm.Query {
		for i, field := range fields {
			if i == 0 {
				query = query.Where(field+" LIKE ?", searchTerm)
			} else {
				query = query.OrWhere(field+" LIKE ?", searchTerm)
			}
		}
		return query
	})
	
	return f
}

// Where adds a simple WHERE condition
func (f *FilterBuilder) Where(field string, value interface{}) *FilterBuilder {
	f.query = f.query.Where(field, value)
	return f
}

// WhereRange adds range conditions (min/max)
func (f *FilterBuilder) WhereRange(field string, min, max interface{}) *FilterBuilder {
	if min != nil {
		f.query = f.query.Where(field+" >= ?", min)
	}
	if max != nil {
		f.query = f.query.Where(field+" <= ?", max)
	}
	return f
}

// WhereIn adds a WHERE IN condition
func (f *FilterBuilder) WhereIn(field string, values []interface{}) *FilterBuilder {
	if len(values) > 0 {
		f.query = f.query.WhereIn(field, values)
	}
	return f
}

// WhereLike adds a LIKE condition
func (f *FilterBuilder) WhereLike(field string, pattern string) *FilterBuilder {
	if pattern != "" {
		f.query = f.query.Where(field+" LIKE ?", "%"+pattern+"%")
	}
	return f
}

// With preloads relationships
func (f *FilterBuilder) With(relations ...string) *FilterBuilder {
	for _, relation := range relations {
		f.query = f.query.With(relation)
	}
	return f
}

// Order adds ordering
func (f *FilterBuilder) Order(sort string) *FilterBuilder {
	if sort != "" {
		f.query = f.query.Order(sort)
	}
	return f
}

// Paginate adds pagination
func (f *FilterBuilder) Paginate(page, pageSize int) *FilterBuilder {
	offset := (page - 1) * pageSize
	f.query = f.query.Offset(offset).Limit(pageSize)
	return f
}

// Build returns the final query
func (f *FilterBuilder) Build() orm.Query {
	return f.query
}

// Count returns the count for the current query conditions
func (f *FilterBuilder) Count() (int64, error) {
	var count int64
	err := f.query.Count(&count)
	return count, err
}