# HTTP Scoped Permissions - Quick Reference

## 🚀 Quick Start

```bash
# Setup and run all tests
mkdir -p resources/views && touch resources/views/dummy.tmpl
APP_ENV=testing go test -v ./tests/feature -run TestHTTPScopedPermissionsTestSuite
```

## 📋 Test Commands

| Purpose | Command |
|---------|---------|
| Run all tests | `APP_ENV=testing go test -v ./tests/feature -run TestHTTPScopedPermissionsTestSuite` |
| Run specific test | `APP_ENV=testing go test -v ./tests/feature -run "TestHTTPScopedPermissionsTestSuite/TestShowBookWithScopedPermission"` |
| Check stability | `APP_ENV=testing go test -v ./tests/feature -run TestHTTPScopedPermissionsTestSuite -count=5` |
| With timeout | `APP_ENV=testing go test -v ./tests/feature -run TestHTTPScopedPermissionsTestSuite -timeout=300s` |

## 🔑 Key Patterns

### Creating Test Users
```go
// ✅ CORRECT
user, err := SetupJWTUser("test@example.com", "password", role)

// ❌ WRONG - Won't work with JWT auth
user := &models.User{Email: "test@example.com"}
facades.Orm().Query().Create(user)
```

### API Query Parameters
```go
// Pagination
"/api/books?page=1&pageSize=20"

// Sorting (use 'direction' NOT 'order')
"/api/books?sort=title&direction=asc"  // ✅
"/api/books?sort=title&order=asc"      // ❌

// Search (min 2 chars)
"/api/books/search?q=Go&page=1&pageSize=20"   // ✅
"/api/books/search?q=G&page=1&pageSize=20"    // ❌
```

### Response Structure
```go
// Paginated responses are nested
dataMap := result["data"].(map[string]interface{})
items := dataMap["data"].([]interface{})          // The actual items
pagination := dataMap["pagination"].(map[string]interface{})

// Single item responses
item := result["data"].(map[string]interface{})
```

## 🔍 Permission Scopes

| Scope | Access Level | Example |
|-------|--------------|---------|
| `by_all` | All resources | Admin can see all books |
| `by_my_role` | Same/lower role resources | Manager sees manager & user books |
| `by_me` | Own resources only | User sees only their books |

## 🐛 Common Issues

| Issue | Solution |
|-------|----------|
| Template error | `mkdir -p resources/views && touch resources/views/dummy.tmpl` |
| 403 on owned resource | Ensure using `SetupJWTUser()` not direct user creation |
| Interface conversion panic | Check response structure (nested vs flat) |
| Search returns 500 | Search query must be ≥ 2 characters |
| Tests fail on re-run | Use `GreaterOrEqual` not `Equal` for counts |

## 📊 Test Coverage

All 15 tests cover:
- ✅ Authentication & Authorization
- ✅ CRUD operations with scopes
- ✅ Pagination, Search, Sorting
- ✅ Concurrent access
- ✅ Mixed permissions
- ✅ Bulk operations

## 💡 Debug Tips

```bash
# See SQL queries
APP_ENV=testing go test -v ./tests/feature -run TestHTTPScopedPermissionsTestSuite 2>&1 | grep "SELECT"

# Check specific test output
APP_ENV=testing go test -v ./tests/feature -run "TestHTTPScopedPermissionsTestSuite/TestDeleteBookWithScopedPermission" 2>&1 | less
```

## 🧹 Cleaning Test Data

```bash
# Clean test database manually
./clean_test_db.sh

# Check what's in test database
sqlite3 database/test.sqlite "SELECT COUNT(*) as count, 'books' as table FROM books UNION ALL SELECT COUNT(*), 'users' FROM users;"
```

### Automatic Cleanup
Tests automatically clean up after themselves by:
- Deleting books with test patterns (title/ISBN)
- Removing test users (emails with @example.com or @test.)
- Cleaning test roles and permissions
- Removing related junction table records

## 📝 Checklist for New Tests

- [ ] Use `SetupJWTUser()` for user creation
- [ ] Set `CreatedBy` field for ownership tests
- [ ] Handle nested response structures
- [ ] Use flexible assertions for counts
- [ ] Test with multiple runs (`-count=3`)
- [ ] Add to this guide if introducing new patterns