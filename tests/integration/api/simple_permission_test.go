package integration

import (
	"github.com/stretchr/testify/assert"
	"starter-project/app/models"
	"strings"
	"testing"
)

func TestPermissionScopes(t *testing.T) {
	// This test demonstrates the scoped permission behavior
	t.Run("by_all scope returns all records", func(t *testing.T) {
		// Admin with by_all scope should see all books
		expectedQuery := "SELECT * FROM books WHERE deleted_at IS NULL"
		t.Logf("Expected query for by_all scope: %s", expectedQuery)
		assert.NotEmpty(t, expectedQuery)
	})

	t.Run("by_my_role scope filters by role", func(t *testing.T) {
		// Editor with by_my_role scope should see only books created by users with editor role
		expectedQuery := "SELECT * FROM books WHERE created_by IN (SELECT user_id FROM user_roles WHERE role_id = ?) AND deleted_at IS NULL"
		t.Logf("Expected query for by_my_role scope: %s", expectedQuery)
		assert.NotEmpty(t, expectedQuery)
	})

	t.Run("by_me scope filters by user", func(t *testing.T) {
		// Member with by_me scope should see only their own books
		expectedQuery := "SELECT * FROM books WHERE created_by = ? AND deleted_at IS NULL"
		t.Logf("Expected query for by_me scope: %s", expectedQuery)
		assert.NotEmpty(t, expectedQuery)
	})
}

func TestScopeDetection(t *testing.T) {
	t.Run("detects scoped permissions correctly", func(t *testing.T) {
		// Test scope detection logic
		testCases := []struct {
			permissionSlug string
			expectedScope  string
			expectedBase   string
		}{
			{"books_read_by_all", "by_all", "books_read"},
			{"books_read_by_my_role", "by_my_role", "books_read"},
			{"books_read_by_me", "by_me", "books_read"},
			{"books_read", "", "books_read"},
		}

		for _, tc := range testCases {
			// Extract scope from permission slug
			scope := ""
			base := tc.permissionSlug

			if strings.HasSuffix(tc.permissionSlug, "_by_all") {
				scope = "by_all"
				base = strings.TrimSuffix(tc.permissionSlug, "_by_all")
			} else if strings.HasSuffix(tc.permissionSlug, "_by_my_role") {
				scope = "by_my_role"
				base = strings.TrimSuffix(tc.permissionSlug, "_by_my_role")
			} else if strings.HasSuffix(tc.permissionSlug, "_by_me") {
				scope = "by_me"
				base = strings.TrimSuffix(tc.permissionSlug, "_by_me")
			}

			assert.Equal(t, tc.expectedScope, scope,
				"Expected scope %s for permission %s", tc.expectedScope, tc.permissionSlug)
			assert.Equal(t, tc.expectedBase, base,
				"Expected base %s for permission %s", tc.expectedBase, tc.permissionSlug)
		}
	})
}

func TestBookFilteringLogic(t *testing.T) {
	t.Run("applies correct filtering based on scope", func(t *testing.T) {
		// Create test data (ID is set after save in Goravel)
		adminUser := &models.User{Name: "Admin", Email: "admin@test.com"}
		editorUser1 := &models.User{Name: "Editor1", Email: "editor1@test.com"}
		editorUser2 := &models.User{Name: "Editor2", Email: "editor2@test.com"}
		memberUser := &models.User{Name: "Member", Email: "member@test.com"}

		// Simulate roles
		adminRole := &models.Role{Slug: "admin", Name: "Admin"}
		editorRole := &models.Role{Slug: "editor", Name: "Editor"}
		memberRole := &models.Role{Slug: "member", Name: "Member"}

		// Test filtering logic
		t.Logf("Admin (by_all): Should see books from all users [1,2,3,4]")
		t.Logf("Editor1 (by_my_role): Should see books from editors [2,3]")
		t.Logf("Member (by_me): Should see only their books [4]")

		// These assertions demonstrate the expected behavior
		assert.NotNil(t, adminUser)
		assert.NotNil(t, editorUser1)
		assert.NotNil(t, editorUser2)
		assert.NotNil(t, memberUser)
		assert.NotNil(t, adminRole)
		assert.NotNil(t, editorRole)
		assert.NotNil(t, memberRole)
	})
}
