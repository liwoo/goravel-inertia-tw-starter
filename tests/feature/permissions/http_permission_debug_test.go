package feature

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/assert"

	"players/app/auth"
	"players/app/models"
	"players/tests"
)

type HTTPPermissionDebugTestSuite struct {
	tests.TestCase
}

func TestHTTPPermissionDebug(t *testing.T) {
	suite := &HTTPPermissionDebugTestSuite{}
	suite.RefreshDatabase()

	// Create test server
	server := httptest.NewServer(facades.Route())
	defer server.Close()

	// Create client
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Create role with slug
	adminRole := &models.Role{Name: "test_admin", Slug: "test_admin", Level: 100, IsActive: true}
	assert.NoError(t, facades.Orm().Query().Create(adminRole))
	fmt.Printf("Created role: %+v\n", adminRole)

	// Create permission with slug
	permission := &models.Permission{
		Name:     "books_read_by_all",
		Slug:     "books_read", // Base slug without scope
		Scope:    "by_all",
		IsActive: true,
	}
	assert.NoError(t, facades.Orm().Query().Create(permission))
	fmt.Printf("Created permission: %+v\n", permission)

	// Assign permission to role with scope in pivot table
	rolePermission := &models.RolePermission{
		RoleID:       adminRole.ID,
		PermissionID: permission.ID,
		Scope:        "by_all",
		GrantedAt:    time.Now(),
		IsActive:     true,
	}
	assert.NoError(t, facades.Orm().Query().Create(rolePermission))

	// Verify role has permission
	var roleWithPerms models.Role
	assert.NoError(t, facades.Orm().Query().Where("id = ?", adminRole.ID).With("Permissions").First(&roleWithPerms))
	fmt.Printf("Role with permissions: %+v\n", roleWithPerms)
	fmt.Printf("Role permissions count: %d\n", len(roleWithPerms.Permissions))
	for _, perm := range roleWithPerms.Permissions {
		fmt.Printf("  Permission: %s (scope: %s, active: %v)\n", perm.Name, perm.Scope, perm.IsActive)
	}

	// Create user
	hashedPassword, err := facades.Hash().Make("testpass")
	assert.NoError(t, err)

	admin := &models.User{
		Name:     "Test Admin",
		Email:    "test@example.com",
		Password: hashedPassword,
		IsActive: true,
	}
	assert.NoError(t, facades.Orm().Query().Create(admin))
	fmt.Printf("Created user: %+v\n", admin)

	// Assign role manually
	userRole := &models.UserRole{
		UserID:     admin.ID,
		RoleID:     adminRole.ID,
		AssignedAt: time.Now(),
		IsActive:   true,
	}
	assert.NoError(t, facades.Orm().Query().Create(userRole))
	fmt.Printf("Created user_role: %+v\n", userRole)

	// Verify user has roles loaded
	var userWithRoles models.User
	assert.NoError(t, facades.Orm().Query().Where("id = ?", admin.ID).With("Roles.Permissions").First(&userWithRoles))
	fmt.Printf("User with roles loaded: %+v\n", userWithRoles)
	fmt.Printf("User roles count: %d\n", len(userWithRoles.Roles))
	for _, role := range userWithRoles.Roles {
		fmt.Printf("Role: %s (active: %v), Permissions: %d\n", role.Name, role.IsActive, len(role.Permissions))
		for _, perm := range role.Permissions {
			fmt.Printf("  Permission: %s (scope: %s, active: %v)\n", perm.Name, perm.Scope, perm.IsActive)
		}
	}

	// Test direct permission check
	permService := auth.NewPermissionService()

	// Also check what permissions are loaded
	fmt.Println("\nChecking loaded permissions:")
	var rolePermissions []models.RolePermission
	facades.Orm().Query().Where("role_id = ? AND is_active = ?", adminRole.ID, true).With("Permission").Find(&rolePermissions)
	fmt.Printf("Found %d role permissions for role %d\n", len(rolePermissions), adminRole.ID)
	for _, rp := range rolePermissions {
		fmt.Printf("  RolePermission: role_id=%d, permission_id=%d, scope=%s, permission_slug=%s\n",
			rp.RoleID, rp.PermissionID, rp.Scope, rp.Permission.Slug)
	}

	hasPermission := permService.HasPermission(&userWithRoles, "books_read_by_all")
	fmt.Printf("\nDirect permission check for 'books_read_by_all': %v\n", hasPermission)

	// Test with different formats
	testPermissions := []string{
		"books_read_by_all",
		"books_read",
		"books.read",
		"books",
	}
	for _, perm := range testPermissions {
		has := permService.HasPermission(&userWithRoles, perm)
		fmt.Printf("HasPermission('%s'): %v\n", perm, has)
	}

	// Test login
	loginBody := fmt.Sprintf(`{"email":"test@example.com","password":"testpass"}`)
	req, _ := http.NewRequest("POST", server.URL+"/api/auth/login", strings.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Login response status: %d\n", resp.StatusCode)

	// Check for cookie
	var authCookie *http.Cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "token" {
			authCookie = cookie
			break
		}
	}

	if authCookie != nil {
		fmt.Printf("Got auth cookie\n")

		// Create a book
		book := &models.Book{
			Title:  "Test Book",
			Author: "Test Author",
			ISBN:   "123456",
		}
		book.CreatedBy = &admin.ID
		assert.NoError(t, facades.Orm().Query().Create(book))
		fmt.Printf("Created book: %+v\n", book)

		// Test accessing books endpoint
		req, _ = http.NewRequest("GET", server.URL+"/api/books", nil)
		req.AddCookie(authCookie)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		body, _ = io.ReadAll(resp.Body)
		fmt.Printf("Books response status: %d\n", resp.StatusCode)
		fmt.Printf("Books response body: %s\n", body)

		// Also check what permission slug is expected
		permSlug := auth.GetPermissionSlug(auth.ServiceBooks, auth.PermissionRead, auth.ScopeByAll)
		fmt.Printf("Expected permission slug: %s\n", permSlug)

		// And verify user still has it
		hasExpectedPerm := permService.HasPermission(&userWithRoles, permSlug)
		fmt.Printf("User has permission '%s': %v\n", permSlug, hasExpectedPerm)
	} else {
		fmt.Println("No auth cookie found")
	}
}
