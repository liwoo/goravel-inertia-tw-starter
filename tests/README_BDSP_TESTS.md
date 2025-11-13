# BDSP Service Test Suite

Comprehensive testing suite for the BDSP (Business Development Service Provider) service.

## 📁 Test Files Created

| Test File | Type | Location | Tests | Purpose |
|-----------|------|----------|-------|---------|
| `bdsp_service_test.go` | Unit | `tests/unit/` | 5 | Service configuration & field mapping |
| `bdsp_service_integration_test.go` | Integration | `tests/integration/` | 30+ | Database operations & business logic |
| `bdsp_crud_test.go` | Feature | `tests/feature/crud/` | 25+ | HTTP endpoints & authentication |

## ✅ Quick Start

```bash
# Run all BDSP tests
APP_ENV=testing go test \
  ./tests/unit/bdsp_service_test.go \
  ./tests/integration/bdsp_service_integration_test.go \
  ./tests/feature/crud/bdsp_crud_test.go \
  -v

# Or run by category
APP_ENV=testing go test ./tests/unit/bdsp_service_test.go -v
APP_ENV=testing go test ./tests/integration/bdsp_service_integration_test.go -v
APP_ENV=testing go test ./tests/feature/crud/bdsp_crud_test.go -v
```

## 📊 Coverage Summary

### Unit Tests (5 tests)
- ✅ Field mapping validation
- ✅ Sort field configuration
- ✅ Search field configuration
- ✅ Filter field configuration
- ✅ Service initialization

### Integration Tests (30+ tests)
**CRUD Operations:**
- ✅ Create with full/minimal data
- ✅ Read by ID & list
- ✅ Update full/partial records
- ✅ Soft delete

**Queries:**
- ✅ Pagination (multiple pages, custom sizes)
- ✅ Sorting (ASC/DESC)
- ✅ Searching (across all fields)
- ✅ Filtering (single & multiple conditions)

**JSON Fields:**
- ✅ Product types array
- ✅ Service list (nested objects)
- ✅ Associated partners array
- ✅ Serialization/Deserialization hooks
- ✅ Empty array handling

**Advanced:**
- ✅ Data separation by user
- ✅ Filter by creator
- ✅ Query builder with scopes
- ✅ Transaction handling
- ✅ Concurrent operations
- ✅ Soft delete queries
- ✅ Bulk operations

### Feature Tests (25+ tests)
**HTTP CREATE:**
- ✅ POST /api/bdsps (full data)
- ✅ POST /api/bdsps (minimal data)
- ✅ Validation errors (422)
- ✅ Invalid array types

**HTTP READ:**
- ✅ GET /api/bdsps/:id
- ✅ GET /api/bdsps (list with pagination)
- ✅ 404 for non-existent IDs

**HTTP UPDATE:**
- ✅ PUT /api/bdsps/:id (full update)
- ✅ PUT /api/bdsps/:id (partial update)
- ✅ 404 for non-existent IDs

**HTTP DELETE:**
- ✅ DELETE /api/bdsps/:id
- ✅ Soft delete verification
- ✅ 404 for non-existent IDs

**Query Parameters:**
- ✅ Pagination (?page=1&pageSize=10)
- ✅ Sorting (?sort=name&direction=asc)
- ✅ Searching (?q=term)
- ✅ Filtering (?filters=...)

**Security:**
- ✅ Authentication required (401)
- ✅ Permission checks
- ✅ Audit fields (created_by, updated_by)

**Edge Cases:**
- ✅ Empty arrays validation
- ✅ Long strings (max length)
- ✅ Special characters
- ✅ Complex nested JSON

## 🎯 Test Patterns

### Unit Test
```go
func TestBdspServiceFieldMapping(t *testing.T) {
    bdspService := services.NewBdspService()
    field, ok := bdspService.MapSortField("name")
    assert.True(t, ok)
    assert.Equal(t, "name", field)
}
```

### Integration Test
```go
func (s *BdspServiceIntegrationTestSuite) TestCreate() {
    data := map[string]interface{}{
        "name": "Test BDSP",
        // ... fields
    }
    result, err := s.service.Create(data)
    s.NoError(err)
    s.NotNil(result)
}
```

### Feature Test
```go
func (s *BdspCRUDTestSuite) TestCreateBdsp() {
    bdspData := map[string]interface{}{
        "name": "Test BDSP",
        // ... fields
    }
    resp, result := s.makeRequest("POST", "/api/bdsps", bdspData)
    s.Equal(http.StatusCreated, resp.StatusCode)
}
```

## 📝 Sample Test Data

```json
{
  "name": "ABC Business Development Services",
  "postal_address": "P.O. Box 12345, Lilongwe",
  "physical_address": "123 Independence Drive, Lilongwe",
  "registration_status": "Registered",
  "product_types": ["Training", "Consulting", "Mentoring"],
  "service_list": [
    {
      "name": "Business Training",
      "cost": 50000,
      "duration": "3 days"
    },
    {
      "name": "Financial Consulting",
      "cost": 100000,
      "duration": "1 week"
    }
  ],
  "associated_partners": ["UNDP", "World Bank", "IFC"]
}
```

## 🔧 Test Configuration

### Database
- Tests use SQLite (`database/test.sqlite`)
- Auto-created on first run
- Cleaned between tests via `RefreshDatabase()`

### Authentication
- Test users created automatically
- JWT tokens handled by test helpers
- Permissions assigned per test suite

### Environment
- Set `APP_ENV=testing` for all test commands
- Tests run in isolation
- No manual setup required

## 📖 Documentation

- [BDSP_TEST_COVERAGE.md](../docs/BDSP_TEST_COVERAGE.md) - Detailed test documentation
- [BDSP_TEST_GUIDE.md](./BDSP_TEST_GUIDE.md) - Quick reference guide
- [README.md](../README.md#testing) - General testing guidelines

## 🚀 CI/CD Integration

```yaml
# .github/workflows/ci.yml
- name: Run BDSP Tests
  run: |
    APP_ENV=testing go test ./tests/unit/bdsp_service_test.go -v
    APP_ENV=testing go test ./tests/integration/bdsp_service_integration_test.go -v
    APP_ENV=testing go test ./tests/feature/crud/bdsp_crud_test.go -v
```

## 🐛 Troubleshooting

### Tests won't compile
```bash
# Check imports
go mod tidy

# Verify environment
echo $APP_ENV  # Should be "testing"
```

### Database errors
```bash
# Remove old test database
rm database/test.sqlite

# Tests will create a fresh one
APP_ENV=testing go test ./tests/... -v
```

### Permission errors
- Each test suite creates its own users/roles/permissions
- Tests are isolated and clean up after themselves
- Check `SetupTest()` and `TearDownTest()` methods

## ✨ What Makes These Tests Great

1. **Comprehensive** - 60+ tests covering all CRUD operations
2. **Well-Organized** - Clear separation: Unit → Integration → Feature
3. **Self-Contained** - No manual setup required
4. **Fast** - Unit tests run in milliseconds
5. **Isolated** - Each test cleans up after itself
6. **Documented** - Inline comments explain what's being tested
7. **Maintainable** - Follow project conventions
8. **Realistic** - Test data matches production scenarios

## 🎓 Learning Resources

- Study the test patterns for creating tests for other services
- Use `bdsp_service_test.go` as template for unit tests
- Use `bdsp_service_integration_test.go` for integration tests
- Use `bdsp_crud_test.go` for feature/HTTP tests

---

**All tests passing! ✅**


