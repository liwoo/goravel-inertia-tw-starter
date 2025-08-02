package unit

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"players/app/services"
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

// Test BookService sort validation
func TestBookServiceSortValidation(t *testing.T) {
	bookService := services.NewBookService()
	
	// Test field validation
	validFields := []string{"title", "author", "published_at", "created_at", "updated_at"}
	for _, field := range validFields {
		assert.True(t, bookService.ValidateSortField(field), "Field %s should be valid", field)
	}
	
	// Test invalid fields
	assert.False(t, bookService.ValidateSortField("invalid_field"))
	
	// Test direction validation
	assert.True(t, bookService.ValidateSortDirection("ASC"))
	assert.True(t, bookService.ValidateSortDirection("asc"))
	assert.True(t, bookService.ValidateSortDirection("DESC"))
	assert.True(t, bookService.ValidateSortDirection("desc"))
	assert.False(t, bookService.ValidateSortDirection("invalid"))
	
	// Test default sort
	field, dir := bookService.GetDefaultSort()
	assert.Equal(t, "created_at", field)
	assert.Equal(t, "DESC", dir)
}

// Test RoleService functionality
func TestRoleServiceMethods(t *testing.T) {
	roleService := services.NewRoleService()
	
	// Test sort field mapping (RoleService doesn't have special mapping)
	field, ok := roleService.MapSortField("name")
	assert.True(t, ok)
	assert.Equal(t, "name", field)
	
	field, ok = roleService.MapSortField("invalid")
	assert.False(t, ok)
	
	// Test default sort
	defaultField, defaultDir := roleService.GetDefaultSort()
	assert.Equal(t, "name", defaultField)
	assert.Equal(t, "ASC", defaultDir)
}

// Test UserService functionality
func TestUserServiceMethods(t *testing.T) {
	userService := services.NewUserService()
	
	// Test field mapping
	field, ok := userService.MapSortField("createdAt")
	assert.True(t, ok)
	assert.Equal(t, "created_at", field)
	
	// Test sort validation
	assert.True(t, userService.ValidateSortField("name"))
	assert.True(t, userService.ValidateSortField("email"))
	assert.True(t, userService.ValidateSortField("created_at"))
	
	// Test default sort
	defaultField, defaultDir := userService.GetDefaultSort()
	assert.Equal(t, "created_at", defaultField)
	assert.Equal(t, "DESC", defaultDir)
}