# BDSP Service Testing Quick Guide

Quick reference for testing the BDSP (Business Development Service Provider) service.

## Quick Commands

```bash
# Run all BDSP tests
./run_tests.sh

# Or run individually:

# Unit tests (fast, no database)
APP_ENV=testing go test ./tests/unit/bdsp_service_test.go -v

# Integration tests (with database)
APP_ENV=testing go test ./tests/integration/bdsp_service_integration_test.go -v

# Feature tests (full HTTP workflow)
APP_ENV=testing go test ./tests/feature/crud/bdsp_crud_test.go -v

# Run all with coverage
APP_ENV=testing go test ./tests/... -v -cover -coverprofile=coverage.out
```

## Test Structure

```
tests/
├── unit/
│   └── bdsp_service_test.go          # 5 tests - Service configuration
├── integration/
│   └── bdsp_service_integration_test.go  # 30+ tests - Database operations
└── feature/
    └── crud/
        └── bdsp_crud_test.go         # 25+ tests - HTTP endpoints
```

## What's Tested

### ✅ Core CRUD
- Create BDSP with full data
- Read by ID and list with pagination
- Update full and partial records
- Soft delete with trashed scope

### ✅ JSON Fields
- Product types array
- Service list (complex nested objects)
- Associated partners array
- Serialization/deserialization

### ✅ Queries
- Pagination (page, pageSize, total, lastPage)
- Sorting (name, created_at, registration_status)
- Searching (name, address, status, JSON fields)
- Filtering (single and multiple conditions)

### ✅ Permissions
- Scoped permissions (by_all, by_me, by_my_role)
- Authentication checks
- Authorization enforcement

### ✅ Edge Cases
- Empty/minimal data
- Invalid data types
- Long strings
- Special characters
- Concurrent operations

## Sample Test Data

```go
// Minimal BDSP
{
    "name": "Test BDSP",
    "registration_status": "Registered",
    "product_types": ["Training"],
    "service_list": [
        {"name": "Service", "cost": 10000, "duration": "1 day"}
    ],
    "associated_partners": ["Partner A"]
}

// Full BDSP
{
    "name": "ABC Business Development Services",
    "postal_address": "P.O. Box 12345",
    "physical_address": "123 Main St, Lilongwe",
    "registration_status": "Registered",
    "product_types": ["Training", "Consulting", "Mentoring"],
    "service_list": [
        {"name": "Business Training", "cost": 50000, "duration": "3 days"},
        {"name": "Financial Consulting", "cost": 100000, "duration": "1 week"}
    ],
    "associated_partners": ["UNDP", "World Bank", "IFC"]
}
```

## Test Patterns

### Unit Test Pattern
```go
func TestBdspServiceFieldMapping(t *testing.T) {
    bdspService := services.NewBdspService()
    field, ok := bdspService.MapSortField("name")
    assert.True(t, ok)
    assert.Equal(t, "name", field)
}
```

### Integration Test Pattern
```go
func (s *BdspServiceIntegrationTestSuite) TestCreate() {
    data := map[string]interface{}{
        "name": "Test BDSP",
        // ... other fields
    }
    result, err := s.service.Create(data)
    s.NoError(err)
    s.NotNil(result)
}
```

### Feature Test Pattern
```go
func (s *BdspCRUDTestSuite) TestCreateBdsp() {
    bdspData := map[string]interface{}{
        "name": "Test BDSP",
        // ... other fields
    }
    resp, result := s.makeRequest("POST", "/api/bdsps", bdspData)
    s.Equal(http.StatusCreated, resp.StatusCode)
    s.True(result["success"].(bool))
}
```

## Common Issues

### Issue: Tests fail with "binding not found: goravel.orm"
**Solution**: Ensure `APP_ENV=testing` is set

```bash
APP_ENV=testing go test ./tests/... -v
```

### Issue: Database connection errors
**Solution**: Tests use SQLite by default (no setup needed)
- Database file: `database/test.sqlite`
- Auto-created on first run
- Cleaned between tests

### Issue: Permission-related test failures
**Solution**: Tests create their own users/roles/permissions
- Each test suite has `SetupTest()` 
- Automatically creates test user with proper permissions
- Cleans up in `TearDownTest()`

### Issue: JSON field errors
**Solution**: Ensure arrays are properly formatted
```go
// ✅ Correct
"product_types": []string{"Training"}

// ❌ Wrong
"product_types": "Training"
```

## Debugging Tests

### Run single test
```bash
APP_ENV=testing go test ./tests/unit/bdsp_service_test.go -v -run TestBdspServiceFieldMapping
```

### Run with verbose output
```bash
APP_ENV=testing go test ./tests/integration/bdsp_service_integration_test.go -v
```

### Check test coverage
```bash
APP_ENV=testing go test ./tests/... -cover
```

### Generate coverage report
```bash
APP_ENV=testing go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Helpers

### Available Helpers (`tests/helpers/`)
- `SetupJWTUser()` - Create authenticated test user
- `AssignPermissionToRole()` - Assign permissions to roles
- `CreateAuthenticatedContext()` - Mock authenticated context

### Using Helpers
```go
// Create test user with role
user, err := helpers.SetupJWTUser("test@example.com", "password", role)

// Assign permission
helpers.AssignPermissionToRole(role, permission, "by_all")
```

## CI/CD Integration

### GitHub Actions Example
```yaml
- name: Run BDSP Tests
  run: |
    APP_ENV=testing go test ./tests/unit/bdsp_service_test.go -v
    APP_ENV=testing go test ./tests/integration/bdsp_service_integration_test.go -v
    APP_ENV=testing go test ./tests/feature/crud/bdsp_crud_test.go -v
```

### Pre-commit Hook
```bash
#!/bin/bash
# Run tests before commit
APP_ENV=testing go test ./tests/... -v
```

## Next Steps

1. ✅ Run unit tests first (fastest)
2. ✅ Run integration tests (database operations)
3. ✅ Run feature tests (full workflow)
4. 📊 Check coverage report
5. 🐛 Fix any failures
6. ✨ Add new tests for new features

## Documentation

- [BDSP_TEST_COVERAGE.md](../docs/BDSP_TEST_COVERAGE.md) - Detailed test documentation
- [README.md](../README.md#testing) - General testing guidelines
- [TEST_ORGANIZATION.md](./TEST_ORGANIZATION.md) - Test structure

## Support

If tests fail:
1. Check `APP_ENV=testing` is set
2. Verify database file exists: `database/test.sqlite`
3. Look at test output for specific error
4. Check documentation for that test category
5. Review test helpers in `tests/helpers/`

---

**Happy Testing! 🧪**


