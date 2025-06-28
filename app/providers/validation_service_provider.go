package providers

import (
	"fmt"

	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
)

type ValidationServiceProvider struct {
}

func (receiver *ValidationServiceProvider) Register(app foundation.Application) {

}

func (receiver *ValidationServiceProvider) Boot(app foundation.Application) {
	if err := facades.Validation().AddRules(receiver.rules()); err != nil {
		facades.Log().Errorf("add rules error: %+v", err)
	}
	if err := facades.Validation().AddFilters(receiver.filters()); err != nil {
		facades.Log().Errorf("add filters error: %+v", err)
	}
}

func (receiver *ValidationServiceProvider) rules() []validation.Rule {
	return []validation.Rule{
		&UniqueRule{},
		&BeforeRule{},
	}
}

func (receiver *ValidationServiceProvider) filters() []validation.Filter {
	return []validation.Filter{}
}

// UniqueRule validates that a field value is unique in the database
type UniqueRule struct {
}

func (r *UniqueRule) Signature() string {
	return "unique"
}

func (r *UniqueRule) Passes(data validation.Data, val any, options ...any) bool {
	if len(options) < 2 {
		return false
	}

	table := options[0].(string)
	column := options[1].(string)
	
	var ignoreID string
	if len(options) > 2 {
		ignoreID = fmt.Sprint(options[2])
	}

	// Convert value to string for comparison
	value := fmt.Sprint(val)
	if value == "" {
		return true // Empty values are handled by required rule
	}

	// Build query
	query := facades.Orm().Query().Table(table).Where(column, value)
	
	// If we have an ID to ignore (for updates), exclude it
	if ignoreID != "" && ignoreID != "0" {
		query = query.Where("id", "!=", ignoreID)
	}

	// Check if record exists
	var count int64
	query.Count(&count)
	
	return count == 0
}

func (r *UniqueRule) Message() string {
	return "The :attribute has already been taken."
}

// BeforeRule validates that a date field is before a specified date
type BeforeRule struct {
}

func (r *BeforeRule) Signature() string {
	return "before"
}

func (r *BeforeRule) Passes(data validation.Data, val any, options ...any) bool {
	if len(options) < 1 {
		return false
	}

	// For "before:today", just return true for now
	// TODO: Implement proper date comparison logic
	return true
}

func (r *BeforeRule) Message() string {
	return "The :attribute must be before the specified date."
}
