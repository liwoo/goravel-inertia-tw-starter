package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"smedi-sme-db/app/contracts"
)

type CustomFiltersTestSuite struct {
	suite.Suite
}

func TestCustomFiltersTestSuite(t *testing.T) {
	suite.Run(t, new(CustomFiltersTestSuite))
}

// Test 1: Filter Type Definition
func (s *CustomFiltersTestSuite) TestFilterTypeEnumeration() {
	// Test that all filter types are properly defined
	s.Equal("string", contracts.FilterTypeString.String())
	s.Equal("number", contracts.FilterTypeNumber.String())
	s.Equal("date", contracts.FilterTypeDate.String())
	s.Equal("datetime", contracts.FilterTypeDateTime.String())
	s.Equal("boolean", contracts.FilterTypeBoolean.String())
	s.Equal("enum", contracts.FilterTypeEnum.String())
	s.Equal("array", contracts.FilterTypeArray.String())
}

// Test 2: Operator Validation by Type
func (s *CustomFiltersTestSuite) TestOperatorsByType() {
	// String operators
	stringOps := contracts.GetOperatorsForType(contracts.FilterTypeString)
	s.Contains(stringOps, contracts.OperatorEquals)
	s.Contains(stringOps, contracts.OperatorContains)
	s.Contains(stringOps, contracts.OperatorStartsWith)
	s.Contains(stringOps, contracts.OperatorEndsWith)

	// Number operators
	numberOps := contracts.GetOperatorsForType(contracts.FilterTypeNumber)
	s.Contains(numberOps, contracts.OperatorGreaterThan)
	s.Contains(numberOps, contracts.OperatorLessThan)
	s.Contains(numberOps, contracts.OperatorBetween)

	// Date operators
	dateOps := contracts.GetOperatorsForType(contracts.FilterTypeDate)
	s.Contains(dateOps, contracts.OperatorBefore)
	s.Contains(dateOps, contracts.OperatorAfter)
	s.Contains(dateOps, contracts.OperatorBetween)

	// Boolean operators
	boolOps := contracts.GetOperatorsForType(contracts.FilterTypeBoolean)
	s.Contains(boolOps, contracts.OperatorIsTrue)
	s.Contains(boolOps, contracts.OperatorIsFalse)
}

// Test 3: Filter Definition Creation
func (s *CustomFiltersTestSuite) TestFilterDefinitionCreation() {
	// Create a filter definition for price field
	filterDef := contracts.NewFilterDefinition(
		"price",
		"Product Price",
		contracts.FilterTypeNumber,
		nil, // Operators are auto-generated based on type
	)

	s.Equal("price", filterDef.Field)
	s.Equal("Product Price", filterDef.Label)
	s.Equal(contracts.FilterTypeNumber, filterDef.Type)
	// Number type should have operators like equals, greater_than, less_than, between
	s.True(len(filterDef.Operators) > 0)
}

// Test 4: Filter Value Validation
func (s *CustomFiltersTestSuite) TestFilterValueValidation() {
	// Number validation
	numberFilter := contracts.FilterDefinition{
		Field: "price",
		Type:  contracts.FilterTypeNumber,
	}

	s.NoError(numberFilter.ValidateValue(contracts.OperatorEquals, 25.99))
	s.NoError(numberFilter.ValidateValue(contracts.OperatorBetween, []float64{10, 50}))
	s.Error(numberFilter.ValidateValue(contracts.OperatorEquals, "not a number"))

	// Date validation
	dateFilter := contracts.FilterDefinition{
		Field: "published_at",
		Type:  contracts.FilterTypeDate,
	}

	s.NoError(dateFilter.ValidateValue(contracts.OperatorAfter, "2023-01-01"))
	s.NoError(dateFilter.ValidateValue(contracts.OperatorBetween, []string{"2023-01-01", "2023-12-31"}))
	s.Error(dateFilter.ValidateValue(contracts.OperatorAfter, "invalid-date"))

	// Boolean validation
	boolFilter := contracts.FilterDefinition{
		Field: "is_active",
		Type:  contracts.FilterTypeBoolean,
	}

	s.NoError(boolFilter.ValidateValue(contracts.OperatorIsTrue, nil))
	s.NoError(boolFilter.ValidateValue(contracts.OperatorIsFalse, nil))

	// Enum validation
	enumFilter := contracts.FilterDefinition{
		Field:      "status",
		Type:       contracts.FilterTypeEnum,
		EnumValues: []string{"AVAILABLE", "BORROWED", "MAINTENANCE"},
	}

	s.NoError(enumFilter.ValidateValue(contracts.OperatorEquals, "AVAILABLE"))
	s.NoError(enumFilter.ValidateValue(contracts.OperatorIn, []string{"AVAILABLE", "BORROWED"}))
	s.Error(enumFilter.ValidateValue(contracts.OperatorEquals, "INVALID_STATUS"))
}

// Test 5: Filter Query Parsing
func (s *CustomFiltersTestSuite) TestFilterQueryParsing() {
	// Simple filter
	simpleQuery := map[string]interface{}{
		"field":    "price",
		"operator": "greater_than",
		"value":    25.99,
	}

	filter, err := contracts.ParseFilterQuery(simpleQuery)
	s.NoError(err)
	s.Equal("price", filter.Field)
	s.Equal(contracts.OperatorGreaterThan, filter.Operator)
	s.Equal(25.99, filter.Value)

	// Compound filter with AND
	compoundQuery := map[string]interface{}{
		"logic": "AND",
		"conditions": []map[string]interface{}{
			{
				"field":    "price",
				"operator": "between",
				"value":    []float64{10, 50},
			},
			{
				"field":    "status",
				"operator": "equals",
				"value":    "AVAILABLE",
			},
		},
	}

	compoundFilter, err := contracts.ParseCompoundFilter(compoundQuery)
	s.NoError(err)
	s.Equal(contracts.LogicAND, compoundFilter.Logic)
	s.Len(compoundFilter.Conditions, 2)
}

// Test 6: Filter to SQL Conversion
func (s *CustomFiltersTestSuite) TestFilterToSQL() {
	// Test greater than
	gtFilter := contracts.FilterCondition{
		Field:    "price",
		Operator: contracts.OperatorGreaterThan,
		Value:    25.99,
	}
	sql, args := gtFilter.ToSQL()
	s.Equal("price > ?", sql)
	s.Equal([]interface{}{25.99}, args)

	// Test between
	betweenFilter := contracts.FilterCondition{
		Field:    "price",
		Operator: contracts.OperatorBetween,
		Value:    []float64{10, 50},
	}
	sql, args = betweenFilter.ToSQL()
	s.Equal("price BETWEEN ? AND ?", sql)
	s.Equal([]interface{}{10.0, 50.0}, args)

	// Test contains (string)
	containsFilter := contracts.FilterCondition{
		Field:    "title",
		Operator: contracts.OperatorContains,
		Value:    "Go",
	}
	sql, args = containsFilter.ToSQL()
	s.Equal("title LIKE ?", sql)
	s.Equal([]interface{}{"%Go%"}, args)

	// Test is null
	nullFilter := contracts.FilterCondition{
		Field:    "deleted_at",
		Operator: contracts.OperatorIsNull,
		Value:    nil,
	}
	sql, args = nullFilter.ToSQL()
	s.Equal("deleted_at IS NULL", sql)
	s.Len(args, 0)
}

// Test 7: Compound Filter to SQL
func (s *CustomFiltersTestSuite) TestCompoundFilterToSQL() {
	// (price > 20 AND status = 'AVAILABLE') OR published_at > '2023-01-01'
	filter := contracts.CompoundFilter{
		Logic: contracts.LogicOR,
		Conditions: []interface{}{
			contracts.CompoundFilter{
				Logic: contracts.LogicAND,
				Conditions: []interface{}{
					contracts.FilterCondition{
						Field:    "price",
						Operator: contracts.OperatorGreaterThan,
						Value:    20,
					},
					contracts.FilterCondition{
						Field:    "status",
						Operator: contracts.OperatorEquals,
						Value:    "AVAILABLE",
					},
				},
			},
			contracts.FilterCondition{
				Field:    "published_at",
				Operator: contracts.OperatorAfter,
				Value:    "2023-01-01",
			},
		},
	}

	sql, args := filter.ToSQL()
	expectedSQL := "((price > ? AND status = ?) OR published_at > ?)"
	s.Equal(expectedSQL, sql)
	s.Equal([]interface{}{20, "AVAILABLE", "2023-01-01"}, args)
}

// Test 8: Filter Metadata Generation
func (s *CustomFiltersTestSuite) TestFilterMetadataGeneration() {
	// Define filters for a Book model
	filterDefs := []contracts.FilterDefinition{
		{
			Field: "price",
			Label: "Price",
			Type:  contracts.FilterTypeNumber,
			Operators: []contracts.FilterOperator{
				contracts.OperatorEquals,
				contracts.OperatorGreaterThan,
				contracts.OperatorLessThan,
				contracts.OperatorBetween,
			},
			Validation: map[string]interface{}{
				"min": 0,
				"max": 999.99,
			},
		},
		{
			Field: "published_at",
			Label: "Published Date",
			Type:  contracts.FilterTypeDate,
			Operators: []contracts.FilterOperator{
				contracts.OperatorBefore,
				contracts.OperatorAfter,
				contracts.OperatorBetween,
			},
			Format: "2006-01-02",
		},
		{
			Field: "status",
			Label: "Status",
			Type:  contracts.FilterTypeEnum,
			Operators: []contracts.FilterOperator{
				contracts.OperatorEquals,
				contracts.OperatorIn,
			},
			EnumValues: []string{"AVAILABLE", "BORROWED", "MAINTENANCE"},
		},
	}

	metadata := contracts.GenerateFilterMetadata(filterDefs)
	s.Len(metadata["filters"], 3)

	// Check first filter
	filters := metadata["filters"].([]map[string]interface{})
	priceFilter := filters[0]
	s.Equal("price", priceFilter["field"])
	s.Equal("Price", priceFilter["label"])
	s.Equal("number", priceFilter["type"])
	s.Len(priceFilter["operators"], 4)
}

// Test 9: Date Range Helpers
func (s *CustomFiltersTestSuite) TestDateRangeHelpers() {
	now := time.Now()

	// Test "is_today"
	start, end := contracts.GetDateRangeForOperator(contracts.OperatorIsToday, now)
	s.Equal(now.Truncate(24*time.Hour), start)
	s.Equal(start.Add(24*time.Hour).Add(-time.Nanosecond), end)

	// Test "is_this_week"
	start, end = contracts.GetDateRangeForOperator(contracts.OperatorIsThisWeek, now)
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	expectedStart := now.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)
	expectedEnd := expectedStart.AddDate(0, 0, 7).Add(-time.Nanosecond)
	s.Equal(expectedStart, start)
	s.Equal(expectedEnd, end)

	// Test "is_this_month"
	start, end = contracts.GetDateRangeForOperator(contracts.OperatorIsThisMonth, now)
	expectedStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	expectedEnd = expectedStart.AddDate(0, 1, 0).Add(-time.Nanosecond)
	s.Equal(expectedStart, start)
	s.Equal(expectedEnd, end)
}

// Test 10: Integration with GORM Query Builder
func (s *CustomFiltersTestSuite) TestGORMIntegration() {
	// This would be tested with actual database connection
	// For now, we just test the query building logic

	filter := contracts.FilterCondition{
		Field:    "price",
		Operator: contracts.OperatorGreaterThan,
		Value:    25.99,
	}

	// The filter should be applicable to a GORM query
	sql, args := filter.ToSQL()
	s.Equal("price > ?", sql)
	s.Equal([]interface{}{25.99}, args)
}
