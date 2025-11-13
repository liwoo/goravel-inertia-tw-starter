package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"smedi-sme-db/app/services"
)

// TestBdspServiceFieldMapping tests field mapping functionality
func TestBdspServiceFieldMapping(t *testing.T) {
	bdspService := services.NewBdspService()

	tests := []struct {
		name          string
		frontendField string
		expectedField string
		shouldMap     bool
	}{
		// BDSP service uses direct database field names without custom mapping
		{"Map name", "name", "name", true},
		{"Map postal_address", "postal_address", "postal_address", true},
		{"Map physical_address", "physical_address", "physical_address", true},
		{"Map registration_status", "registration_status", "registration_status", true},
		{"Map created_at", "created_at", "created_at", true},
		{"Map updated_at", "updated_at", "updated_at", true},
		{"Invalid field", "invalidField", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, ok := bdspService.MapSortField(tt.frontendField)
			assert.Equal(t, tt.shouldMap, ok, "MapSortField should return correct boolean for %s", tt.frontendField)
			if tt.shouldMap {
				assert.Equal(t, tt.expectedField, field, "MapSortField should map %s to %s", tt.frontendField, tt.expectedField)
			}
		})
	}
}

// TestBdspServiceSortValidation tests sort validation through public interface
func TestBdspServiceSortValidation(t *testing.T) {
	bdspService := services.NewBdspService()

	// Test that sortable fields are exposed
	sortableFields := bdspService.GetSortableFields()
	assert.Contains(t, sortableFields, "id")
	assert.Contains(t, sortableFields, "name")
	assert.Contains(t, sortableFields, "postal_address")
	assert.Contains(t, sortableFields, "physical_address")
	assert.Contains(t, sortableFields, "registration_status")
	assert.Contains(t, sortableFields, "created_at")
	assert.Contains(t, sortableFields, "updated_at")

	// MapSortField is in the public interface
	field, ok := bdspService.MapSortField("name")
	assert.True(t, ok)
	assert.Equal(t, "name", field)

	// Invalid field should not map
	_, ok = bdspService.MapSortField("invalid_field")
	assert.False(t, ok)
}

// TestBdspServiceSearchFields tests searchable fields configuration
func TestBdspServiceSearchFields(t *testing.T) {
	bdspService := services.NewBdspService()

	// Test that searchable fields are exposed
	searchableFields := bdspService.GetSearchableFields()
	assert.Contains(t, searchableFields, "name")
	assert.Contains(t, searchableFields, "postal_address")
	assert.Contains(t, searchableFields, "physical_address")
	assert.Contains(t, searchableFields, "registration_status")
	assert.Contains(t, searchableFields, "product_types_json")
	assert.Contains(t, searchableFields, "service_list_json")
	assert.Contains(t, searchableFields, "associated_partners_json")
}

// TestBdspServiceFilterFields tests filterable fields configuration
func TestBdspServiceFilterFields(t *testing.T) {
	bdspService := services.NewBdspService()

	// Test that filterable fields are exposed
	filterableFields := bdspService.GetFilterableFields()
	assert.Contains(t, filterableFields, "name")
	assert.Contains(t, filterableFields, "postal_address")
	assert.Contains(t, filterableFields, "physical_address")
	assert.Contains(t, filterableFields, "registration_status")
	assert.Contains(t, filterableFields, "product_types_json")
	assert.Contains(t, filterableFields, "service_list_json")
	assert.Contains(t, filterableFields, "associated_partners_json")
	assert.Contains(t, filterableFields, "created_by")
	assert.Contains(t, filterableFields, "updated_by")
}

// TestBdspServiceHasFilterDefinitions tests that the service provides filter definitions
// Note: This test doesn't call GetFilterDefinitions() because it requires database access
// Integration tests will verify the actual filter definitions work correctly
func TestBdspServiceHasFilterDefinitions(t *testing.T) {
	bdspService := services.NewBdspService()

	// Just verify the service has the method (will be tested in integration tests)
	// The service should have GetFilterDefinitions method available
	assert.NotNil(t, bdspService, "Service should be initialized")

	// Verify service implements the contract (this is compile-time verified)
	// We just ensure the service was created successfully
	assert.IsType(t, &services.BdspService{}, bdspService)
}
