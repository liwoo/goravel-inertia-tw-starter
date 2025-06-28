package contracts

import "github.com/goravel/framework/contracts/http"

// ValidationRequest defines the contract for validation requests
type ValidationRequest interface {
	// Goravel's native validation methods
	Rules(ctx http.Context) map[string]string
	Messages(ctx http.Context) map[string]string
	Attributes(ctx http.Context) map[string]string

	// Optional authorization
	Authorize(ctx http.Context) error

	// Optional custom validation
	PrepareForValidation(ctx http.Context) error
	PassedValidation(ctx http.Context) error
}

// CreateRequest for resource creation
type CreateRequest interface {
	ValidationRequest
	ToCreateData() map[string]interface{}
}

// UpdateRequest for resource updates
type UpdateRequest interface {
	ValidationRequest
	ToUpdateData() map[string]interface{}
	GetResourceID() interface{}
}

// BulkRequest for bulk operations
type BulkRequest interface {
	ValidationRequest
	GetIDs() []interface{}
	ToData() []map[string]interface{}
}

// Common validation rules constants
const (
	// Required fields
	Required = "required"

	// String validations
	MinLength = "min:%d"        // min:3
	MaxLength = "max:%d"        // max:255
	Email     = "email"
	URL       = "url"
	Alpha     = "alpha"
	AlphaNum  = "alpha_num"
	Regex     = "regex:%s"      // regex:^[a-zA-Z0-9_]+$

	// Numeric validations
	Numeric  = "numeric"
	Integer  = "integer"
	MinValue = "min:%v"         // min:0
	MaxValue = "max:%v"         // max:100
	Between  = "between:%v,%v"  // between:1,100

	// Date validations
	Date       = "date"
	DateFormat = "date_format:%s" // date_format:Y-m-d
	Before     = "before:%s"      // before:2024-12-31
	After      = "after:%s"       // after:2024-01-01

	// Database validations
	Unique = "unique:%s,%s" // unique:users,email
	Exists = "exists:%s,%s" // exists:categories,id

	// File validations
	File    = "file"
	Image   = "image"
	Mimes   = "mimes:%s" // mimes:pdf,doc,docx
	MaxSize = "max:%d"   // max:2048 (KB)

	// Array validations
	Array    = "array"
	ArrayMin = "min:%d" // For arrays: min:1
	ArrayMax = "max:%d" // For arrays: max:10
)