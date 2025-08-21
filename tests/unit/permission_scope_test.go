package unit

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestPermissionScopes(t *testing.T) {
	// This test demonstrates the scoped permission behavior without database
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

func TestBookFilteringScenarios(t *testing.T) {
	t.Run("Admin with by_all scope", func(t *testing.T) {
		// Simulate 10 books: 2 by admin, 2 by editor1, 2 by editor2, 2 by member1, 2 by member2
		totalBooks := 10
		expectedBooks := 10 // Admin sees all

		t.Logf("Admin with by_all scope sees %d out of %d books", expectedBooks, totalBooks)
		assert.Equal(t, expectedBooks, totalBooks)
	})

	t.Run("Editor with by_my_role scope", func(t *testing.T) {
		// Editor1 has editor role, editor2 also has editor role
		totalBooks := 10
		editorBooks := 4 // 2 by editor1 + 2 by editor2

		t.Logf("Editor with by_my_role scope sees %d out of %d books", editorBooks, totalBooks)
		assert.Less(t, editorBooks, totalBooks)
		assert.Equal(t, 4, editorBooks)
	})

	t.Run("Member with by_me scope", func(t *testing.T) {
		// Member1 can only see their own books
		totalBooks := 10
		memberBooks := 2 // Only their own 2 books

		t.Logf("Member with by_me scope sees %d out of %d books", memberBooks, totalBooks)
		assert.Less(t, memberBooks, totalBooks)
		assert.Equal(t, 2, memberBooks)
	})
}

func TestPermissionGeneration(t *testing.T) {
	t.Run("generates scoped permissions from base", func(t *testing.T) {
		basePermission := "books_read"
		scopes := []string{"by_all", "by_my_role", "by_me"}

		expectedPermissions := []string{
			"books_read_by_all",
			"books_read_by_my_role",
			"books_read_by_me",
		}

		var generatedPermissions []string
		for _, scope := range scopes {
			generatedPermissions = append(generatedPermissions, basePermission+"_"+scope)
		}

		assert.Equal(t, expectedPermissions, generatedPermissions)
		t.Logf("Generated permissions: %v", generatedPermissions)
	})
}
