package contracts

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// FilterType represents the data type of a filter field
type FilterType string

const (
	FilterTypeString   FilterType = "string"
	FilterTypeNumber   FilterType = "number"
	FilterTypeDate     FilterType = "date"
	FilterTypeDateTime FilterType = "datetime"
	FilterTypeBoolean  FilterType = "boolean"
	FilterTypeEnum     FilterType = "enum"
	FilterTypeArray    FilterType = "array"
)

// String returns the string representation of the filter type
func (ft FilterType) String() string {
	return string(ft)
}

// FilterOperator represents an operation that can be performed on a filter
type FilterOperator string

const (
	// Common operators
	OperatorEquals    FilterOperator = "equals"
	OperatorNotEquals FilterOperator = "not_equals"
	OperatorIsNull    FilterOperator = "is_null"
	OperatorIsNotNull FilterOperator = "is_not_null"

	// String operators
	OperatorContains    FilterOperator = "contains"
	OperatorNotContains FilterOperator = "not_contains"
	OperatorStartsWith  FilterOperator = "starts_with"
	OperatorEndsWith    FilterOperator = "ends_with"
	OperatorIsEmpty     FilterOperator = "is_empty"
	OperatorIsNotEmpty  FilterOperator = "is_not_empty"
	OperatorRegexMatch  FilterOperator = "regex_match"

	// Number operators
	OperatorGreaterThan        FilterOperator = "greater_than"
	OperatorLessThan           FilterOperator = "less_than"
	OperatorGreaterThanOrEqual FilterOperator = "greater_than_or_equal"
	OperatorLessThanOrEqual    FilterOperator = "less_than_or_equal"
	OperatorBetween            FilterOperator = "between"
	OperatorNotBetween         FilterOperator = "not_between"

	// Date/DateTime operators
	OperatorBefore      FilterOperator = "before"
	OperatorAfter       FilterOperator = "after"
	OperatorIsToday     FilterOperator = "is_today"
	OperatorIsYesterday FilterOperator = "is_yesterday"
	OperatorIsThisWeek  FilterOperator = "is_this_week"
	OperatorIsThisMonth FilterOperator = "is_this_month"
	OperatorIsThisYear  FilterOperator = "is_this_year"
	OperatorLastNDays   FilterOperator = "last_n_days"
	OperatorNextNDays   FilterOperator = "next_n_days"

	// Boolean operators
	OperatorIsTrue  FilterOperator = "is_true"
	OperatorIsFalse FilterOperator = "is_false"

	// Array operators
	OperatorIn          FilterOperator = "in"
	OperatorNotIn       FilterOperator = "not_in"
	OperatorContainsAny FilterOperator = "contains_any"
	OperatorContainsAll FilterOperator = "contains_all"
)

// LogicOperator represents AND/OR logic for compound filters
type LogicOperator string

const (
	LogicAND LogicOperator = "AND"
	LogicOR  LogicOperator = "OR"
)

// GetOperatorsForType returns valid operators for a given filter type
func GetOperatorsForType(filterType FilterType) []FilterOperator {
	commonOps := []FilterOperator{
		OperatorEquals,
		OperatorNotEquals,
		OperatorIsNull,
		OperatorIsNotNull,
	}

	switch filterType {
	case FilterTypeString:
		return append(commonOps,
			OperatorContains,
			OperatorNotContains,
			OperatorStartsWith,
			OperatorEndsWith,
			OperatorIsEmpty,
			OperatorIsNotEmpty,
			OperatorRegexMatch,
		)

	case FilterTypeNumber:
		return append(commonOps,
			OperatorGreaterThan,
			OperatorLessThan,
			OperatorGreaterThanOrEqual,
			OperatorLessThanOrEqual,
			OperatorBetween,
			OperatorNotBetween,
		)

	case FilterTypeDate, FilterTypeDateTime:
		return append(commonOps,
			OperatorBefore,
			OperatorAfter,
			OperatorBetween,
			OperatorNotBetween,
			OperatorIsToday,
			OperatorIsYesterday,
			OperatorIsThisWeek,
			OperatorIsThisMonth,
			OperatorIsThisYear,
			OperatorLastNDays,
			OperatorNextNDays,
		)

	case FilterTypeBoolean:
		return []FilterOperator{
			OperatorIsTrue,
			OperatorIsFalse,
			OperatorIsNull,
			OperatorIsNotNull,
		}

	case FilterTypeEnum:
		return []FilterOperator{
			OperatorEquals,
			OperatorNotEquals,
			OperatorIn,
			OperatorNotIn,
			OperatorIsNull,
			OperatorIsNotNull,
		}

	case FilterTypeArray:
		return []FilterOperator{
			OperatorContains,
			OperatorNotContains,
			OperatorContainsAny,
			OperatorContainsAll,
			OperatorIsEmpty,
			OperatorIsNotEmpty,
		}

	default:
		return commonOps
	}
}

// FilterDefinition defines a filterable field with its metadata
type FilterDefinition struct {
	Field      string                 `json:"field"`
	Label      string                 `json:"label"`
	Type       FilterType             `json:"type"`
	Operators  []FilterOperator       `json:"operators"`
	Validation map[string]interface{} `json:"validation,omitempty"`
	Format     string                 `json:"format,omitempty"`      // For date/datetime formatting
	EnumValues []string               `json:"enum_values,omitempty"` // For enum type
	Default    interface{}            `json:"default,omitempty"`     // Default value
}

// NewFilterDefinition creates a new filter definition
func NewFilterDefinition(field, label string, filterType FilterType, enumValues *[]string) FilterDefinition {
	operators := GetOperatorsForType(filterType)
	var finalEnumValues []string
	if enumValues != nil {
		finalEnumValues = *enumValues
	}
	return FilterDefinition{
		Field:      field,
		Label:      label,
		Type:       filterType,
		Operators:  operators,
		EnumValues: finalEnumValues,
	}
}

// ValidateValue validates a filter value based on the field type
func (fd FilterDefinition) ValidateValue(operator FilterOperator, value interface{}) error {
	// Check if operator is valid for this type
	validOps := fd.Operators
	if len(validOps) == 0 {
		validOps = GetOperatorsForType(fd.Type)
	}

	operatorValid := false
	for _, op := range validOps {
		if op == operator {
			operatorValid = true
			break
		}
	}
	if !operatorValid {
		return fmt.Errorf("operator %s is not valid for field %s of type %s", operator, fd.Field, fd.Type)
	}

	// Some operators don't require values
	if operator == OperatorIsNull || operator == OperatorIsNotNull ||
		operator == OperatorIsEmpty || operator == OperatorIsNotEmpty ||
		operator == OperatorIsTrue || operator == OperatorIsFalse ||
		operator == OperatorIsToday || operator == OperatorIsYesterday ||
		operator == OperatorIsThisWeek || operator == OperatorIsThisMonth ||
		operator == OperatorIsThisYear {
		return nil
	}

	// Validate value based on type
	switch fd.Type {
	case FilterTypeNumber:
		return fd.validateNumberValue(operator, value)
	case FilterTypeDate, FilterTypeDateTime:
		return fd.validateDateValue(operator, value)
	case FilterTypeBoolean:
		return fd.validateBooleanValue(value)
	case FilterTypeEnum:
		return fd.validateEnumValue(operator, value)
	case FilterTypeString:
		return fd.validateStringValue(operator, value)
	case FilterTypeArray:
		return fd.validateArrayValue(operator, value)
	}

	return nil
}

func (fd FilterDefinition) validateNumberValue(operator FilterOperator, value interface{}) error {
	switch operator {
	case OperatorBetween, OperatorNotBetween:
		// Expect array of two numbers
		slice, ok := value.([]float64)
		if !ok {
			// Try interface slice
			if iSlice, ok := value.([]interface{}); ok && len(iSlice) == 2 {
				for _, v := range iSlice {
					if _, err := toFloat64(v); err != nil {
						return fmt.Errorf("between operator requires two numeric values")
					}
				}
				return nil
			}
			return fmt.Errorf("between operator requires array of two numbers")
		}
		if len(slice) != 2 {
			return fmt.Errorf("between operator requires exactly two values")
		}
	default:
		// Single number expected
		if _, err := toFloat64(value); err != nil {
			return fmt.Errorf("invalid number value: %v", value)
		}
	}
	return nil
}

func (fd FilterDefinition) validateDateValue(operator FilterOperator, value interface{}) error {
	switch operator {
	case OperatorBetween, OperatorNotBetween:
		// Expect array of two dates
		if slice, ok := value.([]string); ok {
			if len(slice) != 2 {
				return fmt.Errorf("between operator requires exactly two dates")
			}
			for _, dateStr := range slice {
				if _, err := time.Parse("2006-01-02", dateStr); err != nil {
					return fmt.Errorf("invalid date format: %s", dateStr)
				}
			}
		} else {
			return fmt.Errorf("between operator requires array of two date strings")
		}
	case OperatorLastNDays, OperatorNextNDays:
		// Expect a number
		if _, err := toFloat64(value); err != nil {
			return fmt.Errorf("operator %s requires a numeric value", operator)
		}
	default:
		// Single date string expected
		dateStr, ok := value.(string)
		if !ok {
			return fmt.Errorf("date value must be a string")
		}
		if _, err := time.Parse("2006-01-02", dateStr); err != nil {
			return fmt.Errorf("invalid date format: %s", dateStr)
		}
	}
	return nil
}

func (fd FilterDefinition) validateBooleanValue(value interface{}) error {
	if value == nil {
		return nil
	}
	_, ok := value.(bool)
	if !ok {
		return fmt.Errorf("boolean value must be true or false")
	}
	return nil
}

func (fd FilterDefinition) validateEnumValue(operator FilterOperator, value interface{}) error {
	if len(fd.EnumValues) == 0 {
		return nil // No enum values defined, accept any
	}

	switch operator {
	case OperatorIn, OperatorNotIn:
		// Expect array of enum values
		values, ok := value.([]string)
		if !ok {
			// Try interface slice
			if iSlice, ok := value.([]interface{}); ok {
				for _, v := range iSlice {
					str, ok := v.(string)
					if !ok {
						return fmt.Errorf("enum values must be strings")
					}
					if !contains(fd.EnumValues, str) {
						return fmt.Errorf("invalid enum value: %s", str)
					}
				}
				return nil
			}
			return fmt.Errorf("in operator requires array of enum values")
		}
		for _, v := range values {
			if !contains(fd.EnumValues, v) {
				return fmt.Errorf("invalid enum value: %s", v)
			}
		}
	default:
		// Single enum value
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("enum value must be a string")
		}
		if !contains(fd.EnumValues, str) {
			return fmt.Errorf("invalid enum value: %s", str)
		}
	}
	return nil
}

func (fd FilterDefinition) validateStringValue(operator FilterOperator, value interface{}) error {
	_, ok := value.(string)
	if !ok {
		return fmt.Errorf("string value expected")
	}
	return nil
}

func (fd FilterDefinition) validateArrayValue(operator FilterOperator, value interface{}) error {
	// Arrays can contain various types
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return fmt.Errorf("array value expected")
	}
	return nil
}

// Helper functions
func toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// FilterCondition represents a single filter condition
type FilterCondition struct {
	Field    string         `json:"field"`
	Operator FilterOperator `json:"operator"`
	Value    interface{}    `json:"value"`
	Type     FilterType     `json:"type,omitempty"`
}

// ToSQL converts the filter condition to SQL with placeholders
func (fc FilterCondition) ToSQL() (string, []interface{}) {
	var sql string
	var args []interface{}

	switch fc.Operator {
	case OperatorEquals:
		sql = fmt.Sprintf("%s = ?", fc.Field)
		args = []interface{}{fc.Value}

	case OperatorNotEquals:
		sql = fmt.Sprintf("%s != ?", fc.Field)
		args = []interface{}{fc.Value}

	case OperatorGreaterThan:
		sql = fmt.Sprintf("%s > ?", fc.Field)
		args = []interface{}{fc.Value}

	case OperatorLessThan:
		sql = fmt.Sprintf("%s < ?", fc.Field)
		args = []interface{}{fc.Value}

	case OperatorGreaterThanOrEqual:
		sql = fmt.Sprintf("%s >= ?", fc.Field)
		args = []interface{}{fc.Value}

	case OperatorLessThanOrEqual:
		sql = fmt.Sprintf("%s <= ?", fc.Field)
		args = []interface{}{fc.Value}

	case OperatorBetween:
		if slice, ok := fc.Value.([]float64); ok && len(slice) == 2 {
			sql = fmt.Sprintf("%s BETWEEN ? AND ?", fc.Field)
			args = []interface{}{slice[0], slice[1]}
		} else if slice, ok := fc.Value.([]interface{}); ok && len(slice) == 2 {
			sql = fmt.Sprintf("%s BETWEEN ? AND ?", fc.Field)
			args = slice
		}

	case OperatorNotBetween:
		if slice, ok := fc.Value.([]float64); ok && len(slice) == 2 {
			sql = fmt.Sprintf("%s NOT BETWEEN ? AND ?", fc.Field)
			args = []interface{}{slice[0], slice[1]}
		} else if slice, ok := fc.Value.([]interface{}); ok && len(slice) == 2 {
			sql = fmt.Sprintf("%s NOT BETWEEN ? AND ?", fc.Field)
			args = slice
		}

	case OperatorContains:
		sql = fmt.Sprintf("%s LIKE ?", fc.Field)
		args = []interface{}{fmt.Sprintf("%%%v%%", fc.Value)}

	case OperatorNotContains:
		sql = fmt.Sprintf("%s NOT LIKE ?", fc.Field)
		args = []interface{}{fmt.Sprintf("%%%v%%", fc.Value)}

	case OperatorStartsWith:
		sql = fmt.Sprintf("%s LIKE ?", fc.Field)
		args = []interface{}{fmt.Sprintf("%v%%", fc.Value)}

	case OperatorEndsWith:
		sql = fmt.Sprintf("%s LIKE ?", fc.Field)
		args = []interface{}{fmt.Sprintf("%%%v", fc.Value)}

	case OperatorIsNull:
		sql = fmt.Sprintf("%s IS NULL", fc.Field)
		args = []interface{}{}

	case OperatorIsNotNull:
		sql = fmt.Sprintf("%s IS NOT NULL", fc.Field)
		args = []interface{}{}

	case OperatorIsEmpty:
		sql = fmt.Sprintf("(%s IS NULL OR %s = '')", fc.Field, fc.Field)
		args = []interface{}{}

	case OperatorIsNotEmpty:
		sql = fmt.Sprintf("(%s IS NOT NULL AND %s != '')", fc.Field, fc.Field)
		args = []interface{}{}

	case OperatorIn:
		if slice, ok := fc.Value.([]string); ok {
			placeholders := make([]string, len(slice))
			for i, v := range slice {
				placeholders[i] = "?"
				args = append(args, v)
			}
			sql = fmt.Sprintf("%s IN (%s)", fc.Field, joinStrings(placeholders, ","))
		} else if slice, ok := fc.Value.([]interface{}); ok {
			placeholders := make([]string, len(slice))
			for i, v := range slice {
				placeholders[i] = "?"
				args = append(args, v)
			}
			sql = fmt.Sprintf("%s IN (%s)", fc.Field, joinStrings(placeholders, ","))
		}

	case OperatorNotIn:
		if slice, ok := fc.Value.([]string); ok {
			placeholders := make([]string, len(slice))
			for i, v := range slice {
				placeholders[i] = "?"
				args = append(args, v)
			}
			sql = fmt.Sprintf("%s NOT IN (%s)", fc.Field, joinStrings(placeholders, ","))
		} else if slice, ok := fc.Value.([]interface{}); ok {
			placeholders := make([]string, len(slice))
			for i, v := range slice {
				placeholders[i] = "?"
				args = append(args, v)
			}
			sql = fmt.Sprintf("%s NOT IN (%s)", fc.Field, joinStrings(placeholders, ","))
		}

	case OperatorBefore, OperatorAfter:
		if fc.Operator == OperatorBefore {
			sql = fmt.Sprintf("%s < ?", fc.Field)
		} else {
			sql = fmt.Sprintf("%s > ?", fc.Field)
		}
		args = []interface{}{fc.Value}

	case OperatorIsTrue:
		sql = fmt.Sprintf("%s = ?", fc.Field)
		args = []interface{}{true}

	case OperatorIsFalse:
		sql = fmt.Sprintf("%s = ?", fc.Field)
		args = []interface{}{false}

	default:
		// For date range operators and others not yet implemented
		sql = fmt.Sprintf("%s = ?", fc.Field)
		args = []interface{}{fc.Value}
	}

	return sql, args
}

// CompoundFilter represents a compound filter with AND/OR logic
type CompoundFilter struct {
	Logic      LogicOperator `json:"logic"`
	Conditions []interface{} `json:"conditions"` // Can be FilterCondition or CompoundFilter
}

// ToSQL converts the compound filter to SQL with placeholders
func (cf CompoundFilter) ToSQL() (string, []interface{}) {
	var sqlParts []string
	var allArgs []interface{}

	for _, condition := range cf.Conditions {
		var sql string
		var args []interface{}

		switch c := condition.(type) {
		case FilterCondition:
			sql, args = c.ToSQL()
		case CompoundFilter:
			sql, args = c.ToSQL()
			// Nested compound filters already wrap themselves in parentheses
		case map[string]interface{}:
			// Handle map representation (from JSON parsing)
			if fc, err := parseFilterConditionFromMap(c); err == nil {
				sql, args = fc.ToSQL()
			}
		}

		if sql != "" {
			sqlParts = append(sqlParts, sql)
			allArgs = append(allArgs, args...)
		}
	}

	if len(sqlParts) == 0 {
		return "", []interface{}{}
	}

	separator := fmt.Sprintf(" %s ", cf.Logic)
	combinedSQL := fmt.Sprintf("(%s)", joinStrings(sqlParts, separator))

	return combinedSQL, allArgs
}

// ParseFilterQuery parses a filter query from a map
func ParseFilterQuery(query map[string]interface{}) (*FilterCondition, error) {
	field, ok := query["field"].(string)
	if !ok {
		return nil, fmt.Errorf("field is required")
	}

	operatorStr, ok := query["operator"].(string)
	if !ok {
		return nil, fmt.Errorf("operator is required")
	}

	operator := FilterOperator(operatorStr)

	value := query["value"]

	filter := &FilterCondition{
		Field:    field,
		Operator: operator,
		Value:    value,
	}

	// If type is specified, use it
	if typeStr, ok := query["type"].(string); ok {
		filter.Type = FilterType(typeStr)
	}

	return filter, nil
}

// ParseCompoundFilter parses a compound filter from a map
func ParseCompoundFilter(query map[string]interface{}) (*CompoundFilter, error) {
	logicStr, ok := query["logic"].(string)
	if !ok {
		return nil, fmt.Errorf("logic operator is required for compound filter")
	}

	logic := LogicOperator(logicStr)
	if logic != LogicAND && logic != LogicOR {
		return nil, fmt.Errorf("invalid logic operator: %s", logicStr)
	}

	conditionsRaw, ok := query["conditions"].([]interface{})
	if !ok {
		conditionsSlice, ok := query["conditions"].([]map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("conditions array is required for compound filter")
		}
		// Convert to []interface{}
		conditionsRaw = make([]interface{}, len(conditionsSlice))
		for i, c := range conditionsSlice {
			conditionsRaw[i] = c
		}
	}

	filter := &CompoundFilter{
		Logic:      logic,
		Conditions: []interface{}{},
	}

	for _, condRaw := range conditionsRaw {
		condMap, ok := condRaw.(map[string]interface{})
		if !ok {
			continue
		}

		// Check if it's a nested compound filter
		if _, hasLogic := condMap["logic"]; hasLogic {
			nestedFilter, err := ParseCompoundFilter(condMap)
			if err != nil {
				return nil, err
			}
			filter.Conditions = append(filter.Conditions, *nestedFilter)
		} else {
			// It's a simple filter condition
			condition, err := ParseFilterQuery(condMap)
			if err != nil {
				return nil, err
			}
			filter.Conditions = append(filter.Conditions, *condition)
		}
	}

	return filter, nil
}

// parseFilterConditionFromMap is a helper to parse FilterCondition from a map
func parseFilterConditionFromMap(m map[string]interface{}) (*FilterCondition, error) {
	return ParseFilterQuery(m)
}

// GenerateFilterMetadata generates metadata for frontend consumption
func GenerateFilterMetadata(definitions []FilterDefinition) map[string]interface{} {
	filters := make([]map[string]interface{}, len(definitions))

	for i, def := range definitions {
		filterMeta := map[string]interface{}{
			"field":     def.Field,
			"label":     def.Label,
			"type":      def.Type.String(),
			"operators": make([]string, len(def.Operators)),
		}

		// Convert operators to strings
		for j, op := range def.Operators {
			filterMeta["operators"].([]string)[j] = string(op)
		}

		// Add optional fields if present
		if def.Validation != nil {
			filterMeta["validation"] = def.Validation
		}
		if def.Format != "" {
			filterMeta["format"] = def.Format
		}
		if len(def.EnumValues) > 0 {
			filterMeta["enum_values"] = def.EnumValues
		}
		if def.Default != nil {
			filterMeta["default"] = def.Default
		}

		filters[i] = filterMeta
	}

	return map[string]interface{}{
		"filters": filters,
		"logic_operators": []string{
			string(LogicAND),
			string(LogicOR),
		},
	}
}

// GetDateRangeForOperator returns start and end dates for date range operators
func GetDateRangeForOperator(operator FilterOperator, baseDate time.Time) (time.Time, time.Time) {
	var start, end time.Time

	switch operator {
	case OperatorIsToday:
		start = baseDate.Truncate(24 * time.Hour)
		end = start.Add(24 * time.Hour).Add(-time.Nanosecond)

	case OperatorIsYesterday:
		start = baseDate.AddDate(0, 0, -1).Truncate(24 * time.Hour)
		end = start.Add(24 * time.Hour).Add(-time.Nanosecond)

	case OperatorIsThisWeek:
		// Assuming week starts on Monday
		weekday := int(baseDate.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday = 7
		}
		start = baseDate.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)
		end = start.AddDate(0, 0, 7).Add(-time.Nanosecond)

	case OperatorIsThisMonth:
		start = time.Date(baseDate.Year(), baseDate.Month(), 1, 0, 0, 0, 0, baseDate.Location())
		end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	case OperatorIsThisYear:
		start = time.Date(baseDate.Year(), 1, 1, 0, 0, 0, 0, baseDate.Location())
		end = start.AddDate(1, 0, 0).Add(-time.Nanosecond)

	default:
		// For other operators, return the base date as both start and end
		start = baseDate
		end = baseDate
	}

	return start, end
}

// Helper function to join strings
func joinStrings(strs []string, separator string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += separator + strs[i]
	}
	return result
}
