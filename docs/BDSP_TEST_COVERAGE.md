# BDSP Service Test Coverage

This document describes the comprehensive test coverage for the BDSP (Business Development Service Provider) service.

## Test Organization

Tests are organized into three categories following the project's testing standards:

### 1. Unit Tests (`tests/unit/bdsp_service_test.go`)
Tests isolated service functionality without database operations.

### 2. Integration Tests (`tests/integration/bdsp_service_integration_test.go`)
Tests service operations with real database interactions.

### 3. Feature Tests (`tests/feature/crud/bdsp_crud_test.go`)
Tests complete HTTP CRUD workflows with authentication.

---

## Unit Tests

**Location:** `tests/unit/bdsp_service_test.go`

### Test Coverage

#### Field Mapping Tests
- ✅ `TestBdspServiceFieldMapping` - Tests field name mapping
  - Maps database field names correctly
  - Validates sortable fields (name, postal_address, physical_address, etc.)
  - Rejects invalid field names

#### Sort Validation Tests
- ✅ `TestBdspServiceSortValidation` - Tests sort field validation
  - Verifies sortable fields include: id, name, postal_address, physical_address, registration_status, created_at, updated_at
  - Validates MapSortField returns correct boolean for valid/invalid fields

#### Search Configuration Tests
- ✅ `TestBdspServiceSearchFields` - Tests searchable field configuration
  - Verifies searchable fields: name, postal_address, physical_address, registration_status
  - Validates JSON field searching: product_types_json, service_list_json, associated_partners_json

#### Filter Configuration Tests
- ✅ `TestBdspServiceFilterFields` - Tests filterable field configuration
  - Validates all filterable fields are properly configured
  - Includes audit fields: created_by, updated_by

- ✅ `TestBdspServiceHasFilterDefinitions` - Verifies service provides filter definitions
  - Confirms service initialization
  - Filter definitions tested in integration tests (require database)

### Running Unit Tests

```bash
APP_ENV=testing go test ./tests/unit/bdsp_service_test.go -v
```

---

## Integration Tests

**Location:** `tests/integration/bdsp_service_integration_test.go`

### Test Coverage

#### CREATE Operations
- ✅ `TestCreate` - Creates BDSP with full data
  - Validates all fields are saved correctly
  - Verifies JSON serialization (product_types, service_list, associated_partners)
  - Checks audit fields (created_by)

- ✅ `TestCreateWithMinimalData` - Creates BDSP with required fields only
  - Tests optional fields can be omitted
  - Validates minimal valid record

#### READ Operations
- ✅ `TestGetByID` - Fetches BDSP by ID
  - Verifies all fields are retrieved
  - Validates JSON deserialization

- ✅ `TestGetByIDNotFound` - Handles non-existent ID
  - Returns appropriate error

#### UPDATE Operations
- ✅ `TestUpdate` - Updates BDSP record
  - Modifies multiple fields
  - Validates updated_by audit field
  - Verifies JSON arrays can be updated

- ✅ `TestUpdateServices` - Updates service list
  - Tests complex nested JSON updates
  - Validates service details (name, cost, duration)

#### DELETE Operations
- ✅ `TestDelete` - Soft deletes BDSP
  - Verifies soft delete (deleted_at set)
  - Confirms record not found in normal queries
  - Validates record exists with WithTrashed()

#### LIST & PAGINATION
- ✅ `TestList` - Lists BDSPs with pagination
  - Returns correct number of records
  - Validates pagination metadata

- ✅ `TestListPagination` - Tests pagination edge cases
  - Page 1, 2, 3 with different page sizes
  - Validates last_page calculation

#### SEARCH Operations
- ✅ `TestSearch` - Searches across multiple fields
  - Search by name
  - Search by location (physical_address)
  - Search by registration status

#### SORTING
- ✅ `TestSorting` - Tests sorting functionality
  - Sort by name (ASC/DESC)
  - Verifies correct order

#### FILTERING
- ✅ `TestFiltering` - Tests filter conditions
  - Filters by registration_status
  - Validates filtered results

- ✅ `TestComplexFiltering` - Tests multiple filter combinations
  - Multiple AND conditions
  - Verifies all conditions are applied

#### SCOPED PERMISSIONS
- ✅ `TestScopedFilteringByAll` - Tests "by_all" scope
  - User sees all records regardless of creator

- ✅ `TestScopedFilteringByMe` - Tests "by_me" scope
  - User sees only their own records

- ✅ `TestScopedFilteringByMyRole` - Tests "by_my_role" scope
  - User sees records created by users with same role

#### JSON FIELDS
- ✅ `TestJSONFieldsSerialization` - Tests JSON field handling
  - Product types array
  - Service list with complex objects
  - Associated partners array
  - Verifies persistence and retrieval

- ✅ `TestBdspModelJSONHooks` - Tests BeforeSave/AfterFind hooks
  - Validates automatic JSON serialization
  - Verifies automatic JSON deserialization

- ✅ `TestEmptyJSONArrays` - Tests empty array handling
  - Creates with empty arrays
  - Verifies arrays are initialized correctly

#### ADVANCED OPERATIONS
- ✅ `TestQueryBuilderWithScopes` - Tests query builder with scopes
- ✅ `TestTransactionRollback` - Tests transaction handling
- ✅ `TestConcurrentOperations` - Tests concurrent updates
- ✅ `TestWithTrashedScope` - Tests soft delete queries
- ✅ `TestOrderByWithRelations` - Tests ordering with relations
- ✅ `TestPaginationEdgeCases` - Tests pagination boundaries
- ✅ `TestValidateBeforeCreate` - Tests validation on create
- ✅ `TestBulkOperations` - Tests multiple operations in sequence

### Running Integration Tests

```bash
APP_ENV=testing go test ./tests/integration/bdsp_service_integration_test.go -v
```

---

## Feature Tests (HTTP CRUD)

**Location:** `tests/feature/crud/bdsp_crud_test.go`

### Test Coverage

#### CREATE via HTTP
- ✅ `TestCreateBdsp` - POST /api/bdsps
  - Creates BDSP with complete data
  - Validates response structure
  - Checks HTTP 201 status

- ✅ `TestCreateBdspMinimal` - Creates with minimal required fields
  - Tests optional fields can be omitted

- ✅ `TestCreateBdspValidation` - Tests validation errors
  - Missing required fields return 422
  - Validation message is clear

- ✅ `TestCreateBdspInvalidArrayType` - Tests type validation
  - Invalid array types rejected

#### READ via HTTP
- ✅ `TestGetBdsp` - GET /api/bdsps/:id
  - Fetches single BDSP
  - Returns complete data with JSON arrays

- ✅ `TestGetBdspNotFound` - GET non-existent ID
  - Returns 404 status

- ✅ `TestListBdsps` - GET /api/bdsps
  - Lists all BDSPs with pagination
  - Returns pagination metadata

#### UPDATE via HTTP
- ✅ `TestUpdateBdsp` - PUT /api/bdsps/:id
  - Updates multiple fields
  - Returns updated data

- ✅ `TestUpdateBdspPartial` - Partial update
  - Updates single field
  - Other fields remain unchanged

- ✅ `TestUpdateBdspNotFound` - Update non-existent ID
  - Returns 404 status

#### DELETE via HTTP
- ✅ `TestDeleteBdsp` - DELETE /api/bdsps/:id
  - Soft deletes BDSP
  - Returns 204 status
  - Verifies soft delete in database

- ✅ `TestDeleteBdspNotFound` - Delete non-existent ID
  - Returns 404 status

#### PAGINATION via HTTP
- ✅ `TestPagination` - Tests pagination parameters
  - Default page size (20 items)
  - Custom page size
  - Multiple pages
  - Correct pagination metadata

#### SORTING via HTTP
- ✅ `TestSorting` - Tests sort parameters
  - Sort by name ASC
  - Sort by name DESC
  - Verifies correct order

#### SEARCH via HTTP
- ✅ `TestSearch` - GET /api/bdsps/search?q=term
  - Search by name
  - Search by location
  - Search by status

#### FILTERING via HTTP
- ✅ `TestFiltering` - Tests filter parameter
  - Filter by registration_status
  - Validates filtered results

#### AUTHENTICATION & PERMISSIONS
- ✅ `TestUnauthorizedAccess` - Tests auth requirement
  - Returns 401 without token

- ✅ `TestCreateWithoutPermission` - Tests permission enforcement
  - Verifies permission checks work

#### JSON FIELDS via HTTP
- ✅ `TestComplexServiceList` - Tests complex nested JSON
  - Multiple services with details
  - Multiple product types
  - Multiple partners

#### AUDIT FIELDS via HTTP
- ✅ `TestAuditFields` - Tests audit field population
  - created_by set on create
  - updated_at set on update

#### EDGE CASES via HTTP
- ✅ `TestEmptyArrays` - Tests empty array validation
  - Required arrays cannot be empty

- ✅ `TestLongStrings` - Tests string length validation
  - Strings exceeding max length rejected

- ✅ `TestSpecialCharacters` - Tests special character handling
  - Special characters in names, addresses
  - Proper encoding/decoding

### Running Feature Tests

```bash
APP_ENV=testing go test ./tests/feature/crud/bdsp_crud_test.go -v
```

---

## Test Data Structure

### Sample BDSP Record

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

---

## Test Statistics

### Coverage Summary

| Category | Test File | Test Count | Status |
|----------|-----------|------------|--------|
| Unit | `bdsp_service_test.go` | 5 tests | ✅ Passing |
| Integration | `bdsp_service_integration_test.go` | 30+ tests | 📝 Ready |
| Feature | `bdsp_crud_test.go` | 25+ tests | 📝 Ready |
| **Total** | | **60+ tests** | |

### Test Areas Covered

- ✅ CRUD Operations (Create, Read, Update, Delete)
- ✅ Pagination (multiple pages, custom sizes, edge cases)
- ✅ Sorting (ASC/DESC, multiple fields)
- ✅ Searching (across multiple fields, JSON fields)
- ✅ Filtering (single & multiple conditions)
- ✅ Scoped Permissions (by_all, by_me, by_my_role)
- ✅ JSON Serialization/Deserialization
- ✅ Soft Deletes
- ✅ Audit Fields (created_by, updated_by, deleted_by)
- ✅ Validation (required fields, types, lengths)
- ✅ Authentication & Authorization
- ✅ Edge Cases (empty arrays, special chars, concurrent ops)
- ✅ HTTP Status Codes
- ✅ Error Handling

---

## Running All Tests

```bash
# Run all BDSP tests
APP_ENV=testing go test ./tests/unit/bdsp_service_test.go \
                        ./tests/integration/bdsp_service_integration_test.go \
                        ./tests/feature/crud/bdsp_crud_test.go -v

# Or run by category
APP_ENV=testing go test ./tests/unit -v
APP_ENV=testing go test ./tests/integration -v
APP_ENV=testing go test ./tests/feature/crud -v

# Run with coverage
APP_ENV=testing go test ./tests/... -v -cover
```

---

## Continuous Integration

These tests are designed to run in CI/CD pipelines:

1. **Fast Feedback**: Unit tests run first (no database required)
2. **Integration**: Tests with real database operations
3. **End-to-End**: Full HTTP workflow tests

### CI Configuration

```yaml
test:
  script:
    - APP_ENV=testing go test ./tests/unit/bdsp_service_test.go -v
    - APP_ENV=testing go test ./tests/integration/bdsp_service_integration_test.go -v
    - APP_ENV=testing go test ./tests/feature/crud/bdsp_crud_test.go -v
```

---

## Test Maintenance

### Adding New Tests

When adding new BDSP functionality:

1. **Unit Test**: Add to `tests/unit/bdsp_service_test.go` for business logic
2. **Integration Test**: Add to `tests/integration/bdsp_service_integration_test.go` for database operations
3. **Feature Test**: Add to `tests/feature/crud/bdsp_crud_test.go` for HTTP endpoints

### Test Naming Convention

```go
// Unit: Test<ServiceName><Method>
func TestBdspServiceFieldMapping(t *testing.T) {}

// Integration: Test<Operation>
func (s *BdspServiceIntegrationTestSuite) TestCreate() {}

// Feature: Test<HTTPAction><Resource>
func (s *BdspCRUDTestSuite) TestCreateBdsp() {}
```

---

## Known Limitations

1. **Filter Definitions Test**: Unit test for `GetFilterDefinitions()` requires database access, so it's tested in integration tests instead.
2. **Concurrent Operations**: Basic test provided, but production load testing recommended.
3. **Performance**: Tests focus on correctness, not performance benchmarks.

---

## Related Documentation

- [README.md](../README.md) - Main project documentation
- [CRUD_IMPLEMENTATION.md](./CRUD_IMPLEMENTATION.md) - CRUD implementation guide
- [PERMISSION_SYSTEM_GUIDE.md](./PERMISSION_SYSTEM_GUIDE.md) - Permission system details
- [TEST_ORGANIZATION.md](../tests/TEST_ORGANIZATION.md) - General testing guidelines

---

## Summary

The BDSP service has **comprehensive test coverage** across all three testing categories:

- ✅ **Unit Tests**: Service configuration and field mapping
- ✅ **Integration Tests**: Database operations, scoped permissions, JSON handling
- ✅ **Feature Tests**: Complete HTTP CRUD workflows with authentication

All tests follow the project's established patterns and can be run independently or as part of the full test suite.


