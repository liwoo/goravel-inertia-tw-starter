package unit

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"players/app/models"
)

// TestReflectionWithEmbeddedStructs tests the reflection logic that handles embedded structs
func TestReflectionWithEmbeddedStructs(t *testing.T) {
	// This test verifies the core reflection logic that's used in GenericCrudService

	book := &models.Book{}
	testData := map[string]interface{}{
		"title":      "Test Book",
		"author":     "Test Author",
		"isbn":       "TEST-123",
		"created_by": uint(5), // Field in embedded BaseAuditableModel
	}

	// Apply the reflection logic
	modelValue := reflect.ValueOf(book).Elem()
	setFieldsRecursively(modelValue, testData)

	// Verify regular fields
	assert.Equal(t, "Test Book", book.Title)
	assert.Equal(t, "Test Author", book.Author)
	assert.Equal(t, "TEST-123", book.ISBN)

	// Verify embedded struct field
	assert.NotNil(t, book.CreatedBy, "created_by in embedded struct should be set")
	if book.CreatedBy != nil {
		assert.Equal(t, uint(5), *book.CreatedBy)
	}
}

// Helper function that mimics the logic in GenericCrudService.setFieldsRecursivelyHelper
func setFieldsRecursively(modelValue reflect.Value, data map[string]interface{}) {
	modelType := modelValue.Type()

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldValue := modelValue.Field(i)

		// Handle embedded structs
		if field.Anonymous && fieldValue.Kind() == reflect.Struct {
			setFieldsRecursively(fieldValue, data)
			continue
		}

		// Get field name from json tag
		fieldName := field.Tag.Get("json")
		if fieldName == "" || fieldName == "-" {
			fieldName = strings.ToLower(field.Name)
		} else {
			// Handle json tags with options
			if idx := strings.Index(fieldName, ","); idx != -1 {
				fieldName = fieldName[:idx]
			}
		}

		// Skip if json tag is "-"
		if fieldName == "-" {
			continue
		}

		// Set field value if data exists
		if value, exists := data[fieldName]; exists && fieldValue.CanSet() {
			if value == nil {
				fieldValue.Set(reflect.Zero(fieldValue.Type()))
				continue
			}

			setValue := reflect.ValueOf(value)

			// Handle pointer fields
			if fieldValue.Kind() == reflect.Ptr {
				if setValue.Kind() == reflect.Ptr {
					fieldValue.Set(setValue)
				} else {
					// Create a new pointer and set the value
					newPtr := reflect.New(fieldValue.Type().Elem())
					if setValue.Type().ConvertibleTo(fieldValue.Type().Elem()) {
						newPtr.Elem().Set(setValue.Convert(fieldValue.Type().Elem()))
						fieldValue.Set(newPtr)
					}
				}
			} else if setValue.Type().ConvertibleTo(fieldValue.Type()) {
				fieldValue.Set(setValue.Convert(fieldValue.Type()))
			}
		}
	}
}

// TestFieldDiscovery verifies that all fields including those in embedded structs are discoverable
func TestFieldDiscovery(t *testing.T) {
	bookType := reflect.TypeOf(models.Book{})
	fields := discoverAllFields(bookType)

	// Check that created_by field is discovered
	found := false
	for _, field := range fields {
		if field == "created_by" {
			found = true
			break
		}
	}

	assert.True(t, found, "created_by field should be discovered in embedded struct")
}

// Helper to discover all fields including embedded ones
func discoverAllFields(t reflect.Type) []string {
	var fields []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			// Recursively discover fields in embedded struct
			embeddedFields := discoverAllFields(field.Type)
			fields = append(fields, embeddedFields...)
		} else {
			jsonTag := field.Tag.Get("json")
			if jsonTag != "" && jsonTag != "-" {
				// Handle json tags with options
				if idx := strings.Index(jsonTag, ","); idx != -1 {
					jsonTag = jsonTag[:idx]
				}
				fields = append(fields, jsonTag)
			} else if jsonTag != "-" {
				fields = append(fields, strings.ToLower(field.Name))
			}
		}
	}

	return fields
}
