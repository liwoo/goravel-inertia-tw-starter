# Test Suite Organization

## Overview
This directory contains all tests for the Goravel Blog application, organized by test type and domain.

## Directory Structure

```
tests/
├── unit/               # Unit tests for individual components
│   ├── models/        # Model unit tests
│   └── services/      # Service layer unit tests
├── integration/       # Integration tests
│   ├── api/          # API endpoint integration tests
│   ├── services/     # Service integration tests
│   └── books/        # Book-specific integration tests
├── feature/          # Feature/E2E tests
│   ├── auth/         # Authentication feature tests
│   ├── permissions/  # Permission system tests (main test suite)
│   ├── crud/         # CRUD operation tests
│   └── storage/      # Test storage (sessions, logs, etc)
└── helpers/          # Shared test helper functions
```

## Running Tests

Use the provided test runner script:

```bash
# Run all tests
./run_tests.sh

# Run specific test type
./run_tests.sh -t unit
./run_tests.sh -t integration
./run_tests.sh -t feature

# Run specific domain tests
./run_tests.sh -t permissions
./run_tests.sh -t auth
./run_tests.sh -t crud

# Run with verbose output
./run_tests.sh -v

# Run specific test
./run_tests.sh -t permissions -s TestHTTPScopedPermissionsTestSuite
```

## Manual Test Commands

```bash
# Ensure test environment
export APP_ENV=testing
mkdir -p resources/views && touch resources/views/dummy.tmpl

# Run specific test suite
APP_ENV=testing go test -v ./tests/feature/permissions -run TestHTTPScopedPermissionsTestSuite

# Run with multiple iterations (test stability)
APP_ENV=testing go test -v ./tests/feature/permissions -run TestHTTPScopedPermissionsTestSuite -count=5

# Clean test database
./clean_test_db.sh
```

## Test Organization

### Unit Tests (`/unit`)
- Model validation and business logic
- Service method testing in isolation
- Helper function tests

### Integration Tests (`/integration`)
- API endpoint testing with real database
- Service integration with dependencies
- Cross-service interactions

### Feature Tests (`/feature`)
- End-to-end user scenarios
- Authentication flows
- Permission system testing
- CRUD operations with full stack

## Key Test Files

### Permissions (`/feature/permissions`)
- `http_scoped_permissions_test.go` - Main scoped permissions test suite (15 tests)
- `HTTP_SCOPED_PERMISSIONS_TEST_GUIDE.md` - Comprehensive testing guide
- `HTTP_SCOPED_PERMISSIONS_QUICK_REFERENCE.md` - Quick reference for common tasks

### Authentication (`/feature/auth`)
- JWT authentication tests
- Login/logout flows
- Session management

### CRUD (`/feature/crud`)
- Basic CRUD operation tests
- Sorting and filtering tests

## Test Helpers

Located in `/helpers`:
- `SetupJWTUser()` - Creates JWT-compatible test users
- `AssignPermissionToRole()` - Assigns permissions with scopes
- `auth_test_helper.go` - Authentication test utilities
- `database_cleaner.go` - Test data cleanup utilities

## Test Data Cleanup

Tests automatically clean up after themselves:
- Session files are cleared after each test run
- Test data is removed based on patterns (e.g., @example.com emails)
- Database can be manually cleaned with `./clean_test_db.sh`

## Common Issues & Solutions

1. **Template Error**
   ```bash
   mkdir -p resources/views && touch resources/views/dummy.tmpl
   ```

2. **Database Out of Memory**
   - Tests use SQLite file-based database instead of in-memory
   - Clean database with `./clean_test_db.sh`

3. **JWT Authentication Issues**
   - Always use `SetupJWTUser()` for creating test users
   - Don't create users directly with ORM

4. **Import Errors After Reorganization**
   - Helpers are in `smedi-sme-db/tests/helpers` package
   - Test files remain in `package feature`

## Best Practices

1. Use the test runner script for consistency
2. Keep tests isolated - each test should set up its own data
3. Use helper functions for common operations
4. Clean up test data in TearDown methods
5. Run tests multiple times to ensure stability