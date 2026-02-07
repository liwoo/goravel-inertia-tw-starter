---
name: goravel-crud-review
description: Review and audit a Goravel CRUD implementation for completeness, correctness, and best practices. Checks all API operations, permissions, error handling, query efficiency, route registration, and response consistency.
argument-hint: "[EntityName]"
allowed-tools: Read, Grep, Glob
---

# Goravel CRUD Backend Review

Review CRUD implementation for `$ARGUMENTS`.

## Files to Review

1. `app/http/controllers/<entities>/<entity>_controller.go` — API controller
2. `app/services/<entity>_service.go` — Service layer
3. `app/models/<entity>.go` — Model definition
4. `app/http/requests/<entity>_request.go` — Request validators (create + update)
5. `app/auth/permission_constants.go` — Permission registration
6. `routes/api.go` — Route registration
7. `routes/web.go` — Page route (if Inertia page exists)
8. `tests/feature/crud/<entity>*_test.go` — Test coverage

## Checklist

### 1. CRUD Operations Completeness

- [ ] **Index** (GET `/api/entity-names`) — List with pagination, sorting, filtering, search
- [ ] **Show** (GET `/api/entity-names/{id}`) — Get by ID
- [ ] **Store** (POST `/api/entity-names`) — Create with validation
- [ ] **Update** (PUT `/api/entity-names/{id}`) — Update with validation
- [ ] **Delete** (DELETE `/api/entity-names/{id}`) — Soft delete (if applicable)
- [ ] **Search** (GET `/api/entity-names/search`) — Quick search endpoint
- [ ] **FilterMetadata** (GET `/api/entity-names/filters`) — Available filters

### 2. Route Registration (`routes/api.go`)

- [ ] GET routes in optional auth group (or appropriate auth group)
- [ ] POST/PUT/DELETE routes in protected group (jwtAuth + require2FA)
- [ ] Endpoint naming follows hyphenated format: `/entity-names` not `/entity_names`
- [ ] Search endpoint registered **before** `{id}` route to avoid conflicts
- [ ] No duplicate routes

```go
// Correct route ordering:
optionalAuthRouter.Get("/entity-names", controller.Index)
optionalAuthRouter.Get("/entity-names/search", controller.Search)      // Before {id}!
optionalAuthRouter.Get("/entity-names/filters", controller.FilterMetadata) // Before {id}!
optionalAuthRouter.Get("/entity-names/{id}", controller.Show)          // Last

protectedRouter.Post("/entity-names", controller.Store)
protectedRouter.Put("/entity-names/{id}", controller.Update)
protectedRouter.Delete("/entity-names/{id}", controller.Delete)
```

### 3. Controller Review

- [ ] Extends `contracts.CrudController[Model, *CreateRequest, *UpdateRequest]`
- [ ] `WithAuthChecker` configured with correct `auth.Service*` constant
- [ ] Permission mapping covers: viewAny, view, create, update, delete
- [ ] Uses `auth.GetScopedPermissionHelper()` for scoped permissions
- [ ] `SetBeforeStore` sets `created_by` from authenticated user
- [ ] All custom endpoints have permission checks (`CheckAuth`)
- [ ] Custom endpoints validate input (ID, pagination, etc.)
- [ ] Proper HTTP status codes (200 OK, 201 Created, 204 No Content, 400 Bad Request, 403 Forbidden, 404 Not Found, 422 Unprocessable)

### 4. Service Review

- [ ] Uses `contracts.NewServiceBuilder[Model]` builder pattern
- [ ] `WithSearchFields()` configured (required for search to work)
- [ ] `WithSortFields()` configured with all sortable fields
- [ ] `WithFilterFields()` configured with filterable fields
- [ ] `WithValidationRules()` configured
- [ ] `WithDefaultSort()` set (typically `"created_at", "DESC"`)
- [ ] `WithSoftDeletes()` enabled if model uses soft deletes
- [ ] `WithScopeFiltering()` enabled if permission-scoped access needed
- [ ] `WithRelations()` includes audit relations (`"Creator", "Updater"`) if applicable
- [ ] `GetColumnMapping()` overridden for camelCase→snake_case sort field mapping
- [ ] `GetFilterDefinitions()` returns proper filter types (String, Enum, Number, Date, DateTime, Array)
- [ ] Custom business logic methods exist for specialized endpoints
- [ ] No N+1 query patterns (relations loaded eagerly via `With()`)

### 5. Model Review

- [ ] Extends `BaseAuditableModel` (provides ID, CreatedAt, UpdatedAt, DeletedAt, CreatedBy, UpdatedBy, Creator, Updater)
- [ ] JSON tags use camelCase: `json:"fieldName"`
- [ ] GORM tags specify column names for non-standard fields: `gorm:"column:field_name"`
- [ ] Nullable fields use pointer types: `*string`, `*float64`, `*time.Time`
- [ ] `SearchFields()` method returns searchable field names
- [ ] `TableName()` method returns correct table name
- [ ] `MarshalJSON()` custom marshaler handles date formatting (RFC3339)
- [ ] Array/JSON fields: virtual field with `gorm:"-"` + storage field with `json:"-"`
- [ ] GORM hooks (`BeforeSave`, `AfterFind`) for JSON array serialization

### 6. Request Validation Review

- [ ] **Create request**: All required fields validated, `ToCreateData()` maps all fields
- [ ] **Update request**: Uses pointer types for optional fields, `ToUpdateData()` only includes non-nil fields
- [ ] `Rules()` — Validation rules match service validation
- [ ] `Messages()` — Custom error messages for all rules
- [ ] `Attributes()` — Human-readable field names for error messages
- [ ] `PrepareForValidation()` — Input normalization (trim, defaults)
- [ ] No `|numeric` on float64 fields (known Goravel bug)
- [ ] No `|date` on `*string` date fields (known Goravel bug)
- [ ] Enum fields use `in:VALUE1,VALUE2,...` validation
- [ ] Read-only fields excluded from create/update data

### 7. Permission Review

Check `app/auth/permission_constants.go`:
- [ ] `ServiceRegistry` constant defined (e.g., `ServiceEntityNames`)
- [ ] Registered in `GetAllServiceRegistries()`
- [ ] Display name in `GetServiceDisplayName()`
- [ ] Actions registered in `GetServiceActions()` — typically: create, read, update, delete
- [ ] Service constant in controller matches permission constants exactly
- [ ] Scope filtering configured if needed (by_own, by_team, by_all)

### 8. Error Handling

- [ ] Invalid ID returns 400 (not 500)
- [ ] Not found returns 404
- [ ] Unauthorized returns 403
- [ ] Validation errors return 422 with field-level errors
- [ ] Internal errors return 500 with generic message (no sensitive info)
- [ ] All error paths have proper error responses (no silent failures)

### 9. Query Efficiency

- [ ] No N+1 queries — relations loaded with `With()` in service builder
- [ ] Search uses ILIKE with `%query%` (via service `WithSearchFields`)
- [ ] Pagination configured with sensible defaults
- [ ] Sort field mapping handles camelCase→snake_case
- [ ] Index/unique constraints exist for frequently queried fields

### 10. Audit Fields

- [ ] `created_by` set on create (via controller `SetBeforeStore` hook)
- [ ] `updated_by` set on update (if applicable)
- [ ] `CreatedAt` / `UpdatedAt` auto-managed by GORM
- [ ] `DeletedAt` used for soft deletes (via `WithSoftDeletes`)
- [ ] Audit relations (`Creator`, `Updater`) loaded in service `WithRelations`

### 11. Test Coverage

Check `tests/feature/crud/`:
- [ ] Test suite exists with `SetupSuite`, `TearDownSuite`, `SetupTest`, `TearDownTest`
- [ ] Tests cover all CRUD operations (create, read, update, delete, list)
- [ ] Tests cover search with min query length
- [ ] Tests cover pagination (empty, single page, multi-page)
- [ ] Tests cover sorting (ASC, DESC, multiple fields)
- [ ] Tests cover filtering
- [ ] Tests cover permission-based access (authorized vs unauthorized)
- [ ] Tests cover validation errors (missing required fields, invalid values)
- [ ] Tests cover not-found scenarios (invalid ID)
- [ ] Test data properly initialized (arrays as `[]string{}`, not nil)
- [ ] Database cleaned between tests

## Common Issues to Flag

### Critical
- Missing permission checks on any endpoint
- Service constant mismatch between controller and `permission_constants.go`
- Search endpoint registered after `{id}` route (will never match)
- Validation bypass (missing rules on required fields)

### Important
- Missing `WithSearchFields` (search returns no results)
- Missing sort field mapping (sorting on camelCase fields fails)
- Create request missing `created_by` in `ToCreateData`
- Nullable fields not using pointer types in update request
- Array fields not initialized to empty array (causes nil JSON)

### Nice-to-have
- Missing Swagger annotations
- Missing custom error messages in request validation
- Statistics endpoint not implemented
- Filter metadata not returning enum options

## Output

After review, produce a report with:
1. **Summary**: Overall implementation quality (Complete / Partial / Missing)
2. **Critical Issues**: Must fix before production
3. **Important Issues**: Should fix soon
4. **Minor Issues**: Nice-to-have improvements
5. **Missing Features**: Expected endpoints/functionality not implemented
6. **Test Coverage Gaps**: Untested scenarios

## Reference

- Canonical controller: `app/http/controllers/books/book_controller.go`
- Canonical service: `app/services/book_service.go`
- Canonical model: `app/models/book.go`
- Canonical request: `app/http/requests/book_request.go`
- Canonical test: `tests/feature/crud/lender_crud_test.go`
- Route patterns: `routes/api.go`
