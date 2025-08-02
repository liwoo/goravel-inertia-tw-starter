# Test Organization Plan

## Current Structure Issues
- All tests are in a single `feature` directory
- Many duplicate/debug test files
- Session files polluting the directory
- No clear separation between test types

## Proposed Structure

```
tests/
├── unit/               # Unit tests for individual components
│   ├── models/        # Model tests
│   ├── services/      # Service layer tests
│   └── helpers/       # Helper function tests
├── integration/       # Integration tests
│   ├── api/          # API endpoint tests
│   ├── database/     # Database integration tests
│   └── services/     # Service integration tests
├── feature/          # Feature/E2E tests
│   ├── auth/         # Authentication feature tests
│   ├── permissions/  # Permission system tests
│   └── crud/         # CRUD operation tests
├── helpers/          # Test helper functions
└── fixtures/         # Test data fixtures
```

## Files to Move/Organize

### Authentication Tests → `feature/auth/`
- http_auth_debug_test.go
- http_jwt_workaround_test.go
- jwt_debug_test.go
- jwt_workaround.go (helper)

### Permission Tests → `feature/permissions/`
- http_scoped_permissions_test.go (main test suite)
- permission_scope_test.go
- scoped_permissions_test.go
- http_permission_debug_test.go
- http_permission_fix_test.go
- scoped_permissions_demo_test.go
- scoped_permissions_simple_test.go
- http_scoped_permissions_simple_test.go
- scoped_statistics_regression_test.go

### CRUD Tests → `feature/crud/`
- simple_crud_test.go
- crud_sorting_test.go

### Debug/Temporary Tests → Remove
- http_debug_test.go
- example_test.go

### Integration Tests → `integration/`
- Move API-specific tests from feature to integration/api
- Move service integration tests to integration/services

## Benefits
1. Clear separation of concerns
2. Easier navigation
3. Better test discovery
4. Cleaner imports
5. Isolated test helpers per domain