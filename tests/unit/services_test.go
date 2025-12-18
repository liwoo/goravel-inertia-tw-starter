package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"starter-project/app/services"
)

// Test BookService field mapping functionality
func TestBookServiceFieldMapping(t *testing.T) {
	bookService := services.NewBookService()

	tests := []struct {
		name          string
		frontendField string
		expectedField string
		shouldMap     bool
	}{
		{"Map publishedAt", "publishedAt", "published_at", true},
		{"Map createdAt", "createdAt", "created_at", true},
		{"Map updatedAt", "updatedAt", "updated_at", true},
		{"Map title", "title", "title", true},
		{"Invalid field", "invalidField", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, ok := bookService.MapSortField(tt.frontendField)
			assert.Equal(t, tt.shouldMap, ok, "MapSortField should return correct boolean for %s", tt.frontendField)
			if tt.shouldMap {
				assert.Equal(t, tt.expectedField, field, "MapSortField should map %s to %s", tt.frontendField, tt.expectedField)
			}
		})
	}
}

// Test BookService sort validation through public interface
func TestBookServiceSortValidation(t *testing.T) {
	bookService := services.NewBookService()

	// Test that sortable fields are exposed
	sortableFields := bookService.GetSortableFields()
	assert.Contains(t, sortableFields, "title")
	assert.Contains(t, sortableFields, "author")
	assert.Contains(t, sortableFields, "published_at")
	assert.Contains(t, sortableFields, "created_at")
	assert.Contains(t, sortableFields, "updated_at")

	// MapSortField is in the public interface
	field, ok := bookService.MapSortField("title")
	assert.True(t, ok)
	assert.Equal(t, "title", field)

	// Invalid field should not map
	_, ok = bookService.MapSortField("invalid_field")
	assert.False(t, ok)
}

// Test RoleService functionality through public interface
func TestRoleServiceMethods(t *testing.T) {
	roleService := services.NewRoleService()

	// Test that sortable fields are exposed
	sortableFields := roleService.GetSortableFields()
	assert.Contains(t, sortableFields, "name")
	assert.Contains(t, sortableFields, "created_at")

	// Test sort field mapping (RoleService doesn't have special mapping)
	field, ok := roleService.MapSortField("name")
	assert.True(t, ok)
	assert.Equal(t, "name", field)

	field, ok = roleService.MapSortField("invalid")
	assert.False(t, ok)
}

// Test UserService functionality through public interface
func TestUserServiceMethods(t *testing.T) {
	userService := services.NewUserService()

	// Test that sortable fields are exposed
	sortableFields := userService.GetSortableFields()
	assert.Contains(t, sortableFields, "name")
	assert.Contains(t, sortableFields, "email")
	assert.Contains(t, sortableFields, "created_at")

	// Test field mapping
	field, ok := userService.MapSortField("createdAt")
	assert.True(t, ok)
	assert.Equal(t, "created_at", field)

	// Invalid field should not map
	_, ok = userService.MapSortField("invalid_field")
	assert.False(t, ok)
}
