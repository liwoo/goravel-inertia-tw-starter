package feature

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/assert"

	"players/app/models"
	"players/tests"
)

type HTTPAuthDebugTestSuite struct {
	tests.TestCase
}

func TestHTTPAuthDebug(t *testing.T) {
	suite := &HTTPAuthDebugTestSuite{}
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
	adminRole := &models.Role{Name: "test_admin", Slug: "test_admin", Level: 100}
	assert.NoError(t, facades.Orm().Query().Create(adminRole))
	fmt.Printf("Created role: %+v\n", adminRole)

	// Create permission
	permission := &models.Permission{Name: "books_read_by_all", Scope: "by_all"}
	assert.NoError(t, facades.Orm().Query().Create(permission))
	fmt.Printf("Created permission: %+v\n", permission)

	// Assign permission to role
	assert.NoError(t, facades.Orm().Query().Model(adminRole).Association("Permissions").Append(permission))

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
		fmt.Printf("Role: %s, Permissions: %d\n", role.Name, len(role.Permissions))
		for _, perm := range role.Permissions {
			fmt.Printf("  Permission: %s (scope: %s)\n", perm.Name, perm.Scope)
		}
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
	fmt.Printf("Login response body: %s\n", body)

	// Check for cookie
	var authCookie *http.Cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "token" {
			authCookie = cookie
			break
		}
	}

	if authCookie != nil {
		fmt.Printf("Got auth cookie: %s\n", authCookie.Value)

		// Test /api/auth/me endpoint first
		req, _ = http.NewRequest("GET", server.URL+"/api/auth/me", nil)
		req.AddCookie(authCookie)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		body, _ = io.ReadAll(resp.Body)
		fmt.Printf("Me response status: %d\n", resp.StatusCode)
		fmt.Printf("Me response body: %s\n", body)

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
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authCookie.Value))

		resp, err = client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		body, _ = io.ReadAll(resp.Body)
		fmt.Printf("Books response status: %d\n", resp.StatusCode)
		fmt.Printf("Books response body: %s\n", body)

		// Also test with just the Authorization header
		req2, _ := http.NewRequest("GET", server.URL+"/api/books", nil)
		req2.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authCookie.Value))

		resp2, err := client.Do(req2)
		assert.NoError(t, err)
		defer resp2.Body.Close()

		body2, _ := io.ReadAll(resp2.Body)
		fmt.Printf("Books with Bearer header status: %d\n", resp2.StatusCode)
		fmt.Printf("Books with Bearer header body: %s\n", body2)
	} else {
		fmt.Println("No auth cookie found")

		// Try parsing JSON response for token
		var loginResp map[string]interface{}
		json.Unmarshal(body, &loginResp)
		if data, ok := loginResp["data"].(map[string]interface{}); ok {
			if token, ok := data["token"]; ok {
				fmt.Printf("Found token in response: %v\n", token)
			}
		}
	}
}
