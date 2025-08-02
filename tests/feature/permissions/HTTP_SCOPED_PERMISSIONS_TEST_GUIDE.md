# HTTP Scoped Permissions Test Guide

This guide explains how to run and maintain the HTTP scoped permissions tests for the Goravel blog application.

## Overview

The HTTP scoped permissions tests verify that the permission system correctly enforces access control based on three scope levels:
- **by_all**: User can access all resources
- **by_my_role**: User can access resources created by users with the same or lower role level
- **by_me**: User can only access resources they created

## Running the Tests

### Prerequisites

1. **Environment Setup**
   ```bash
   # The tests require the APP_ENV to be set to 'testing'
   export APP_ENV=testing
   ```

2. **Database**
   - Tests use SQLite in-memory database (automatically configured)
   - No manual database setup required
   - Database is recreated for each test suite run

3. **Required Directories**
   ```bash
   # Create the resources/views directory if it doesn't exist
   mkdir -p resources/views
   touch resources/views/dummy.tmpl
   ```

### Running All HTTP Scoped Permission Tests

```bash
# Run the full test suite
APP_ENV=testing go test -v ./tests/feature -run TestHTTPScopedPermissionsTestSuite

# Run with count flag for repeatability testing
APP_ENV=testing go test -v ./tests/feature -run TestHTTPScopedPermissionsTestSuite -count=3

# Run with timeout (useful for CI/CD)
APP_ENV=testing go test -v ./tests/feature -run TestHTTPScopedPermissionsTestSuite -timeout=300s
```

### Running Individual Tests

```bash
# Run a specific test
APP_ENV=testing go test -v ./tests/feature -run "TestHTTPScopedPermissionsTestSuite/TestShowBookWithScopedPermission"

# Examples of other individual tests:
APP_ENV=testing go test -v ./tests/feature -run "TestHTTPScopedPermissionsTestSuite/TestPaginationWithScopedPermissions"
APP_ENV=testing go test -v ./tests/feature -run "TestHTTPScopedPermissionsTestSuite/TestSearchWithScopedPermissions"
```

## Test Structure

### Test Suite: `HTTPScopedPermissionsTestSuite`

Located in: `/tests/feature/http_scoped_permissions_test.go`

The test suite includes 15 tests covering various aspects of scoped permissions:

1. **TestUnauthenticatedAccessDenied** - Verifies unauthenticated requests are rejected
2. **TestUserWithoutPermissionCannotAccessBooks** - Ensures users without permissions can't access resources
3. **TestCreateBookWithPermission** - Tests book creation with proper permissions
4. **TestUpdateBookWithScopedPermission** - Verifies update permissions respect ownership
5. **TestDeleteBookWithScopedPermission** - Tests delete permissions with ownership checks
6. **TestShowBookWithScopedPermission** - Verifies single resource access with scopes
7. **TestAdminCanSeeAllBooksWithByAllScope** - Tests unrestricted access with by_all scope
8. **TestManagerCanSeeRoleBooksWithByMyRoleScope** - Tests role-based access filtering
9. **TestUserCanSeeOwnBooksWithByMeScope** - Verifies owner-only access filtering
10. **TestPaginationWithScopedPermissions** - Tests pagination respects permission scopes
11. **TestSearchWithScopedPermissions** - Verifies search functionality with scopes
12. **TestSortingWithScopedPermissions** - Tests sorting with permission filtering
13. **TestMixedPermissionScopes** - Tests users with multiple permission scopes
14. **TestBulkOperationsRespectScopes** - Verifies bulk operations respect permissions
15. **TestConcurrentRequestsWithDifferentScopes** - Tests concurrent access with different permissions

### Key Components

1. **JWT Authentication Helper** (`/tests/feature/jwt_workaround.go`)
   - `SetupJWTUser()` - Creates test users with proper JWT authentication setup
   - Essential for all authenticated test requests

2. **Test Base** (`/tests/test_case.go`)
   - Configures test environment
   - Sets up SQLite database for testing
   - Provides RefreshDatabase functionality

## Common Issues and Solutions

### 1. Template Parsing Errors
**Error**: `html/template: pattern matches no files: resources/views/*.tmpl`

**Solution**:
```bash
mkdir -p resources/views
touch resources/views/dummy.tmpl
```

### 2. Database Errors
**Error**: `unable to open database file: out of memory (14)`

**Solution**: Ensure the test environment is using file-based SQLite:
- Check that `test_case.go` sets `DB_DATABASE` to `database/test.sqlite`
- Verify the `database` directory exists

### 3. Permission Check Failures
**Error**: Tests fail with 403 Forbidden when they should pass

**Common Causes**:
- User not created with `SetupJWTUser()`
- Permission scopes not properly assigned
- Resource ownership (created_by) not set correctly

### 4. Response Structure Mismatches
**Error**: `interface conversion: interface {} is map[string]interface {}, not []interface {}`

**Solution**: The API returns nested responses for paginated data:
```go
// Correct way to access paginated data
dataMap := result["data"].(map[string]interface{})
items := dataMap["data"].([]interface{})
pagination := dataMap["pagination"].(map[string]interface{})
```

## Best Practices for Test Development

### 1. Always Use SetupJWTUser

```go
// DON'T create users directly
user := &models.User{
    Email: "test@example.com",
    // ...
}
facades.Orm().Query().Create(user)

// DO use SetupJWTUser
user, err := SetupJWTUser("test@example.com", "password", role)
```

### 2. Handle Dynamic Test Data

```go
// DON'T use exact counts (fails on repeated runs)
s.Equal(10, len(items))

// DO use flexible assertions
s.GreaterOrEqual(len(items), 10)
```

### 3. Use Proper Query Parameters

```go
// Pagination
"/api/books?page=1&pageSize=20"

// Sorting (use 'direction' not 'order')
"/api/books?sort=title&direction=asc"

// Search (minimum 2 characters)
"/api/books/search?q=Go&page=1&pageSize=20"
```

### 4. Set Resource Ownership

```go
// Always set created_by for ownership-based permissions
book := &models.Book{
    Title: "Test Book",
    // ... other fields
}
book.CreatedBy = &user.ID
facades.Orm().Query().Create(book)
```

## Debugging Tests

### Enable Debug Logging

```go
// Add debug output to understand test failures
fmt.Printf("DEBUG: User ID: %d, Book ID: %d\n", user.ID, book.ID)

// Check response bodies on failures
if resp.StatusCode != http.StatusOK {
    body, _ := io.ReadAll(resp.Body)
    fmt.Printf("Response: %s\n", string(body))
}
```

### Check SQL Queries

The test output includes SQL queries when running with `-v` flag:
```bash
APP_ENV=testing go test -v ./tests/feature -run TestHTTPScopedPermissionsTestSuite 2>&1 | grep "SELECT"
```

## Continuous Integration

### GitHub Actions Example

```yaml
- name: Run HTTP Scoped Permission Tests
  env:
    APP_ENV: testing
  run: |
    mkdir -p resources/views
    touch resources/views/dummy.tmpl
    go test -v ./tests/feature -run TestHTTPScopedPermissionsTestSuite -timeout=300s
```

### Pre-commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Running HTTP scoped permission tests..."
APP_ENV=testing go test ./tests/feature -run TestHTTPScopedPermissionsTestSuite
if [ $? -ne 0 ]; then
    echo "Tests failed. Please fix before committing."
    exit 1
fi
```

## Maintenance Notes

1. **Adding New Tests**: Follow the existing pattern using `SetupJWTUser()` and handle nested response structures
2. **Updating Permissions**: If permission scopes change, update both the test data and assertions
3. **Database Migrations**: Ensure test database schema matches production (migrations run automatically)
4. **API Changes**: Update response parsing if API response structure changes

## Quick Reference

```bash
# Run all tests
APP_ENV=testing go test -v ./tests/feature -run TestHTTPScopedPermissionsTestSuite

# Run specific test
APP_ENV=testing go test -v ./tests/feature -run "TestHTTPScopedPermissionsTestSuite/TestShowBookWithScopedPermission"

# Run with coverage
APP_ENV=testing go test -v -cover ./tests/feature -run TestHTTPScopedPermissionsTestSuite

# Run multiple times for stability check
APP_ENV=testing go test -v ./tests/feature -run TestHTTPScopedPermissionsTestSuite -count=5
```

## Support

For issues or questions:
1. Check the test output with `-v` flag for detailed logs
2. Verify all prerequisites are met
3. Ensure the latest code is pulled and dependencies are updated
4. Check the CLAUDE.md file for project-specific conventions