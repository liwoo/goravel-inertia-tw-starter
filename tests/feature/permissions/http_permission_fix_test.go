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

type HTTPPermissionFixTestSuite struct {
	tests.TestCase
}

func TestHTTPPermissionFix(t *testing.T) {
	suite := &HTTPPermissionFixTestSuite{}
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

	// Create permission with proper slug format
	// The slug should NOT include the scope - it's just the base permission
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

	// Check what permissions the user has
	permService := auth.NewPermissionService()
	userPerms := permService.GetUserPermissions(admin)
	fmt.Printf("User permissions: %v\n", userPerms)

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

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	} else {
		fmt.Println("No auth cookie found")
	}
}
