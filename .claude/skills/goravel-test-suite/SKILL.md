---
name: goravel-test-suite
description: Write and run comprehensive Go test suites for Goravel entities using testcontainers (PostgreSQL). Covers CRUD tests, permission tests, scoped access tests, and integration tests with the full test infrastructure.
argument-hint: "[EntityName] or [test-type]"
allowed-tools: Bash, Read, Write, Edit, Grep, Glob
---

# Goravel Test Suite (Testcontainers)

Write/run tests for `$ARGUMENTS`.

## Running Tests

### Using the Shell Script (recommended)

```bash
# Run ALL tests
./scripts/run_tests.sh -v ./tests/...

# Run specific CRUD test suite
./scripts/run_tests.sh -v ./tests/feature/crud -run TestEntityCRUDTestSuite

# Run specific test within a suite
./scripts/run_tests.sh -v ./tests/feature/crud -run TestEntityCRUDTestSuite/TestCreateEntity

# Run integration tests
./scripts/run_tests.sh -v ./tests/integration/...

# Run unit tests (no DB required but script still works)
./scripts/run_tests.sh -v ./tests/unit/...

# Run with verbose + count
./scripts/run_tests.sh -v -count=1 ./tests/feature/crud -run TestEntityCRUDTestSuite
```

**What `run_tests.sh` does:**
1. Starts a PostgreSQL 16 Alpine Docker container with random port
2. Creates isolated temp storage directory (`/tmp/books-test-storage-{PID}`)
3. Sets environment variables (`APP_ENV=testing`, `DB_*`, `JWT_SECRET`, etc.)
4. Runs `go test -p=1` (sequential packages to avoid migration conflicts)
5. Cleans up container + temp dir on exit (even on CTRL+C)

### Key Flags

| Flag | Purpose |
|---|---|
| `-v` | Verbose output |
| `-count=1` | Disable test caching |
| `-run TestSuiteName` | Run specific suite |
| `-run TestSuiteName/TestMethod` | Run specific test method |
| `-p=1` | Sequential packages (already set by script) |
| `-timeout 5m` | Increase timeout for slow tests |

## Test Directory Structure

```
tests/
├── test_case.go                    # Base TestCase (embeds Goravel's TestCase)
├── bootstrap/                      # Testcontainer init (env_setup.go, test_bootstrap.go)
├── helpers/                        # Shared utilities
│   ├── auth_test_helper.go        # Mock auth contexts
│   ├── database_cleaner.go        # TRUNCATE CASCADE cleanup
│   ├── jwt_workaround.go          # SetupJWTUser helper
│   └── http_scoped_permissions_helper.go  # AssignPermissionToRole
├── feature/crud/                   # CRUD operation tests
│   ├── main_test.go               # Package bootstrap
│   ├── <entity>_crud_test.go      # Entity CRUD test suite
│   └── ...
├── integration/                    # Service/API integration tests
│   ├── api/                       # API endpoint tests
│   └── services/                  # Service layer tests
└── unit/                          # Pure unit tests (no DB)
```

## Writing a CRUD Test Suite

### Step 1: Create Test File

Create `tests/feature/crud/<entity>_crud_test.go`:

```go
package crud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"books-database/app/models"
	"books-database/tests"
	"books-database/tests/helpers"
)

// ============================================================================
// Suite Definition
// ============================================================================

type EntityCRUDTestSuite struct {
	suite.Suite
	tests.TestCase

	server     *httptest.Server
	client     *http.Client
	authCookie *http.Cookie
	testUser   *models.User
}

func TestEntityCRUDTestSuite(t *testing.T) {
	suite.Run(t, &EntityCRUDTestSuite{})
}

// ============================================================================
// Suite Lifecycle (runs ONCE per suite)
// ============================================================================

func (s *EntityCRUDTestSuite) SetupSuite() {
	// Start test HTTP server
	s.server = httptest.NewServer(facades.Route())

	// Create HTTP client with cookie jar for session management
	jar, _ := cookiejar.New(nil)
	s.client = &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects
		},
	}
}

func (s *EntityCRUDTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}

// ============================================================================
// Test Lifecycle (runs BEFORE/AFTER each test)
// ============================================================================

func (s *EntityCRUDTestSuite) SetupTest() {
	// Clean database for isolation
	s.RefreshDatabase()

	// Create test user with permissions
	s.setupTestUser()
}

func (s *EntityCRUDTestSuite) TearDownTest() {
	if orm := facades.Orm(); orm != nil {
		orm.Query().Exec("DELETE FROM entities")
		orm.Query().Exec("DELETE FROM users WHERE email = 'entitytest@example.com'")
		orm.Query().Exec("DELETE FROM user_roles")
		orm.Query().Exec("DELETE FROM role_permissions")
		orm.Query().Exec("DELETE FROM roles WHERE slug = 'entity_admin'")
		orm.Query().Exec("DELETE FROM permissions WHERE slug LIKE 'entities_%'")
	}
}

// ============================================================================
// Helper: Test User Setup
// ============================================================================

func (s *EntityCRUDTestSuite) setupTestUser() {
	// Create admin role
	adminRole := &models.Role{
		Name:  "Entity Admin",
		Slug:  "entity_admin",
		Level: 100,
	}
	s.Nil(facades.Orm().Query().Create(adminRole))

	// Create CRUD permissions
	permissions := []string{"create", "read", "update", "delete"}
	for _, perm := range permissions {
		permission := &models.Permission{
			Name:  fmt.Sprintf("entities_%s", perm),
			Slug:  fmt.Sprintf("entities_%s", perm),
			Scope: "by_all",
		}
		s.Nil(facades.Orm().Query().Create(permission))
		s.Nil(helpers.AssignPermissionToRole(adminRole, permission, "by_all"))
	}

	// Create test user with JWT auth
	user, err := helpers.SetupJWTUser("entitytest@example.com", "password", adminRole)
	s.Nil(err)
	s.testUser = user

	// Login to get auth cookie
	s.authCookie = s.loginUser("entitytest@example.com", "password")
	s.NotNil(s.authCookie, "Login should return auth cookie")
}

// ============================================================================
// Helper: Login
// ============================================================================

func (s *EntityCRUDTestSuite) loginUser(email, password string) *http.Cookie {
	loginData := map[string]string{
		"email":    email,
		"password": password,
	}
	jsonData, _ := json.Marshal(loginData)

	resp, err := s.client.Post(
		s.server.URL+"/api/auth/login",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	s.Nil(err)
	defer resp.Body.Close()

	// Extract token cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "token" {
			return cookie
		}
	}

	return nil
}

// ============================================================================
// Helper: Make Authenticated Request
// ============================================================================

func (s *EntityCRUDTestSuite) makeRequest(method, path string, body interface{}) (*http.Response, map[string]interface{}) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, s.server.URL+path, bodyReader)
	s.Nil(err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	// Add auth cookie
	if s.authCookie != nil {
		req.AddCookie(s.authCookie)
	}

	resp, err := s.client.Do(req)
	s.Nil(err)

	respBody, err := io.ReadAll(resp.Body)
	s.Nil(err)
	resp.Body.Close()

	var result map[string]interface{}
	if len(respBody) > 0 {
		json.Unmarshal(respBody, &result)
	}

	return resp, result
}

// ============================================================================
// Helper: Create Test Entity
// ============================================================================

func (s *EntityCRUDTestSuite) createTestEntity(name string) map[string]interface{} {
	data := map[string]interface{}{
		"name":        name,
		"description": "Test description for " + name,
		"status":      "ACTIVE",
		// Add other required fields...
	}

	resp, result := s.makeRequest("POST", "/api/entity-names", data)
	s.Equal(http.StatusCreated, resp.StatusCode, "Create should return 201")

	return result
}

// ============================================================================
// CRUD Tests
// ============================================================================

func (s *EntityCRUDTestSuite) TestCreateEntity() {
	data := map[string]interface{}{
		"name":        "Test Entity",
		"description": "A test entity",
		"status":      "ACTIVE",
	}

	resp, result := s.makeRequest("POST", "/api/entity-names", data)

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.NotNil(result["data"])
}

func (s *EntityCRUDTestSuite) TestGetEntity() {
	// Create an entity first
	created := s.createTestEntity("Get Test")

	// Extract ID from response
	data := created["data"].(map[string]interface{})
	id := int(data["id"].(float64))

	// GET by ID
	resp, result := s.makeRequest("GET", fmt.Sprintf("/api/entity-names/%d", id), nil)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.NotNil(result["data"])
}

func (s *EntityCRUDTestSuite) TestUpdateEntity() {
	// Create first
	created := s.createTestEntity("Update Test")
	data := created["data"].(map[string]interface{})
	id := int(data["id"].(float64))

	// Update
	updateData := map[string]interface{}{
		"name": "Updated Name",
	}
	resp, _ := s.makeRequest("PUT", fmt.Sprintf("/api/entity-names/%d", id), updateData)

	s.Equal(http.StatusOK, resp.StatusCode)
}

func (s *EntityCRUDTestSuite) TestDeleteEntity() {
	// Create first
	created := s.createTestEntity("Delete Test")
	data := created["data"].(map[string]interface{})
	id := int(data["id"].(float64))

	// Delete
	resp, _ := s.makeRequest("DELETE", fmt.Sprintf("/api/entity-names/%d", id), nil)

	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify deleted (should be 404)
	resp2, _ := s.makeRequest("GET", fmt.Sprintf("/api/entity-names/%d", id), nil)
	s.Equal(http.StatusNotFound, resp2.StatusCode)
}

func (s *EntityCRUDTestSuite) TestListEntities() {
	// Create multiple entities
	for i := 1; i <= 5; i++ {
		s.createTestEntity(fmt.Sprintf("List Test %d", i))
	}

	// List
	resp, result := s.makeRequest("GET", "/api/entity-names?page=1&pageSize=10", nil)

	s.Equal(http.StatusOK, resp.StatusCode)

	// Check pagination structure
	if dataWrapper, ok := result["data"].(map[string]interface{}); ok {
		if items, ok := dataWrapper["data"].([]interface{}); ok {
			s.GreaterOrEqual(len(items), 5)
		}
	}
}

func (s *EntityCRUDTestSuite) TestSearchEntities() {
	s.createTestEntity("Searchable Alpha")
	s.createTestEntity("Searchable Beta")
	s.createTestEntity("Other Entity")

	// Search
	resp, result := s.makeRequest("GET", "/api/entity-names/search?q=Searchable", nil)

	s.Equal(http.StatusOK, resp.StatusCode)

	if dataWrapper, ok := result["data"].(map[string]interface{}); ok {
		if items, ok := dataWrapper["data"].([]interface{}); ok {
			s.Equal(2, len(items))
		}
	}
}

func (s *EntityCRUDTestSuite) TestSortEntities() {
	s.createTestEntity("Charlie")
	s.createTestEntity("Alpha")
	s.createTestEntity("Beta")

	// Sort ascending by name
	resp, result := s.makeRequest("GET", "/api/entity-names?sort=name&direction=ASC", nil)

	s.Equal(http.StatusOK, resp.StatusCode)

	if dataWrapper, ok := result["data"].(map[string]interface{}); ok {
		if items, ok := dataWrapper["data"].([]interface{}); ok {
			s.GreaterOrEqual(len(items), 3)
			// Verify order
			first := items[0].(map[string]interface{})
			s.Equal("Alpha", first["name"])
		}
	}
}

func (s *EntityCRUDTestSuite) TestPagination() {
	// Create 15 entities
	for i := 1; i <= 15; i++ {
		s.createTestEntity(fmt.Sprintf("Page Test %02d", i))
	}

	// Get page 1 with size 5
	resp, result := s.makeRequest("GET", "/api/entity-names?page=1&pageSize=5", nil)

	s.Equal(http.StatusOK, resp.StatusCode)

	if dataWrapper, ok := result["data"].(map[string]interface{}); ok {
		if items, ok := dataWrapper["data"].([]interface{}); ok {
			s.Equal(5, len(items))
		}
		// Check pagination metadata
		if pagination, ok := dataWrapper["pagination"].(map[string]interface{}); ok {
			s.Equal(float64(15), pagination["total"])
			s.Equal(float64(3), pagination["last_page"])
		}
	}
}

func (s *EntityCRUDTestSuite) TestCreateValidationError() {
	// Send empty data (missing required fields)
	data := map[string]interface{}{}

	resp, _ := s.makeRequest("POST", "/api/entity-names", data)

	s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
}

func (s *EntityCRUDTestSuite) TestGetNonExistentEntity() {
	resp, _ := s.makeRequest("GET", "/api/entity-names/99999", nil)

	s.Equal(http.StatusNotFound, resp.StatusCode)
}
```

### Step 2: Verify main_test.go Exists

Each test package needs a `main_test.go`. Check `tests/feature/crud/main_test.go` exists:

```go
package crud

import (
	"os"
	"testing"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/foundation"

	_ "books-database/config"
	"books-database/tests/helpers"
)

func TestMain(m *testing.M) {
	os.Setenv("AUTH_REQUIRE_2FA", "false")

	app := foundation.NewApplication()
	app.Boot()

	if err := facades.Artisan().Call("migrate"); err != nil {
		os.Stderr.WriteString("Warning: Migration failed: " + err.Error() + "\n")
	}

	if err := helpers.CleanTestDatabase(); err != nil {
		os.Stderr.WriteString("Warning: Database cleanup failed: " + err.Error() + "\n")
	}

	code := m.Run()

	helpers.CleanTestDatabase()

	os.Exit(code)
}
```

**Important**: If adding tests to a NEW package (not `tests/feature/crud/`), you must create a `main_test.go` for that package.

### Step 3: Run Tests

```bash
./scripts/run_tests.sh -v ./tests/feature/crud -run TestEntityCRUDTestSuite
```

## Test Helpers Reference

### `helpers.SetupJWTUser(email, password, role) (*models.User, error)`

Creates or updates a test user with hashed password and assigned role. Returns fully populated user.

```go
user, err := helpers.SetupJWTUser("test@example.com", "password", adminRole)
```

### `helpers.AssignPermissionToRole(role, permission, scope) error`

Assigns a permission to a role with a specific scope.

```go
helpers.AssignPermissionToRole(adminRole, readPerm, "by_all")
```

Scopes: `"by_all"`, `"by_my_role"`, `"by_me"`

### `helpers.CleanTestDatabase() error`

TRUNCATE CASCADE all tables except migrations. PostgreSQL-specific.

### `helpers.CleanupTestData()`

Pattern-based cleanup: deletes test data matching patterns like `Test %`, `%@example.com`, `test_%`.

## Permission Setup Patterns

### Full CRUD Permissions

```go
func (s *Suite) setupTestUser() {
	role := &models.Role{Name: "Admin", Slug: "admin", Level: 100}
	facades.Orm().Query().Create(role)

	for _, action := range []string{"create", "read", "update", "delete"} {
		perm := &models.Permission{
			Name:  fmt.Sprintf("entities_%s", action),
			Slug:  fmt.Sprintf("entities_%s", action),
			Scope: "by_all",
		}
		facades.Orm().Query().Create(perm)
		helpers.AssignPermissionToRole(role, perm, "by_all")
	}

	user, _ := helpers.SetupJWTUser("test@example.com", "password", role)
	s.testUser = user
	s.authCookie = s.loginUser("test@example.com", "password")
}
```

### Scoped Permission Testing

Test that different scopes restrict data visibility:

```go
// Admin: sees ALL records (by_all)
helpers.AssignPermissionToRole(adminRole, readPerm, "by_all")

// Editor: sees records from users with same role (by_my_role)
helpers.AssignPermissionToRole(editorRole, readPerm, "by_my_role")

// Member: sees only their OWN records (by_me)
helpers.AssignPermissionToRole(memberRole, readPerm, "by_me")
```

## Response Assertion Patterns

### Nested Response Structure

The API returns `{data: {data: [...], pagination: {...}}}` for lists:

```go
if dataWrapper, ok := result["data"].(map[string]interface{}); ok {
	if items, ok := dataWrapper["data"].([]interface{}); ok {
		s.Equal(5, len(items))

		// Access first item
		first := items[0].(map[string]interface{})
		s.Equal("Expected Name", first["name"])
	}
	if pagination, ok := dataWrapper["pagination"].(map[string]interface{}); ok {
		s.Equal(float64(15), pagination["total"])
		s.Equal(float64(1), pagination["current_page"])
	}
}
```

### Single Resource Response

```go
if data, ok := result["data"].(map[string]interface{}); ok {
	s.Equal("Expected Name", data["name"])
	s.Equal(float64(1), data["id"])
}
```

### Error Response

```go
s.Equal(http.StatusUnprocessableEntity, resp.StatusCode)

if errors, ok := result["errors"].(map[string]interface{}); ok {
	s.NotNil(errors["name"], "Should have name validation error")
}
```

## Database Cleanup in TearDownTest

Always clean up test data in reverse dependency order:

```go
func (s *Suite) TearDownTest() {
	orm := facades.Orm()
	if orm == nil {
		return
	}
	// Junction tables first
	orm.Query().Exec("DELETE FROM user_roles")
	orm.Query().Exec("DELETE FROM role_permissions")

	// Entity data
	orm.Query().Exec("DELETE FROM entities")

	// Auth data
	orm.Query().Exec("DELETE FROM users WHERE email = 'entitytest@example.com'")
	orm.Query().Exec("DELETE FROM roles WHERE slug = 'entity_admin'")
	orm.Query().Exec("DELETE FROM permissions WHERE slug LIKE 'entities_%'")
}
```

## Writing Integration Tests (Service Layer)

For testing services directly without HTTP:

```go
package services

import (
	"testing"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"books-database/app/models"
	"books-database/app/services"
	"books-database/app/contracts"
	"books-database/tests"
)

type EntityServiceTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestEntityServiceTestSuite(t *testing.T) {
	suite.Run(t, new(EntityServiceTestSuite))
}

func (s *EntityServiceTestSuite) SetupTest() {
	s.RefreshDatabase()
}

func (s *EntityServiceTestSuite) TestCreate() {
	svc := services.NewEntityService()

	data := map[string]interface{}{
		"name":   "Test",
		"status": "ACTIVE",
	}

	result, err := svc.Create(data)
	s.NoError(err)
	s.NotNil(result)
}

func (s *EntityServiceTestSuite) TestGetList() {
	// Seed data
	for i := 0; i < 5; i++ {
		facades.Orm().Query().Create(&models.Entity{
			Name:   fmt.Sprintf("Entity %d", i),
			Status: "ACTIVE",
		})
	}

	svc := services.NewEntityService()
	result, err := svc.GetList(contracts.ListRequest{
		Page:     1,
		PageSize: 10,
	})

	s.NoError(err)
	s.NotNil(result)
	s.Equal(5, len(result.Data))
}

func (s *EntityServiceTestSuite) TestSearch() {
	facades.Orm().Query().Create(&models.Entity{Name: "Findable", Status: "ACTIVE"})
	facades.Orm().Query().Create(&models.Entity{Name: "Other", Status: "ACTIVE"})

	svc := services.NewEntityService()
	result, err := svc.Search("Find", contracts.ListRequest{Page: 1, PageSize: 10})

	s.NoError(err)
	s.Equal(1, len(result.Data))
}
```

## Writing Unit Tests (No Database)

For testing service logic without a database:

```go
package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"books-database/app/services"
)

func TestEntityServiceFieldMapping(t *testing.T) {
	svc := services.NewEntityService()

	tests := []struct {
		name     string
		input    string
		expected string
		ok       bool
	}{
		{"Map createdAt", "createdAt", "created_at", true},
		{"Map updatedAt", "updatedAt", "updated_at", true},
		{"Map unknown", "fooBar", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, ok := svc.MapSortField(tt.input)
			assert.Equal(t, tt.ok, ok)
			if tt.ok {
				assert.Equal(t, tt.expected, field)
			}
		})
	}
}
```

## Testcontainer Architecture

### Shell Script (`scripts/run_tests.sh`) — Recommended

- Starts PostgreSQL container via `docker run`
- Sets env vars BEFORE Go compiles test binary
- Cleanup trap on EXIT/INT/TERM
- Sequential packages with `-p=1`

### Go Module (`tests/bootstrap/`) — Alternative

- `env_setup.go` starts container in `init()` function
- `test_bootstrap.go` boots Goravel app and runs migrations
- Used by integration tests that import `tests/bootstrap`

### When to Use Which

| Approach | Use When |
|---|---|
| `run_tests.sh` | Running tests locally or in CI. Works with any test package. |
| Go bootstrap | Tests that need programmatic control over the container lifecycle. |

## Common Gotchas

- **JSON array fields**: Initialize as `[]string{}` not `nil` in test data
- **Permission slug format**: `entities_read` (underscore between service and action)
- **Route ordering**: Search/filter routes MUST come before `{id}` route
- **Auth cookie**: Login returns a `token` cookie — must be attached to every request
- **2FA disabled**: `AUTH_REQUIRE_2FA=false` must be set (script handles this)
- **Sequential packages**: `-p=1` prevents migration race conditions
- **Float assertion**: JSON numbers decode as `float64`, use `float64(expected)` in assertions
- **Nested response**: List endpoints return `{data: {data: [...], pagination: {...}}}`

## Reference

- Test script: `scripts/run_tests.sh`
- Test case base: `tests/test_case.go`
- Bootstrap: `tests/bootstrap/env_setup.go`, `tests/bootstrap/test_bootstrap.go`
- Helpers: `tests/helpers/jwt_workaround.go`, `tests/helpers/database_cleaner.go`
- CRUD test example: `tests/feature/crud/lender_crud_test.go`
- Advanced test example: `tests/feature/crud/book_advanced_features_test.go`
- Scoped permissions test: `tests/integration/api/api_scoped_permissions_test.go`
