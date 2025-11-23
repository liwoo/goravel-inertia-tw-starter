# Goravel CRUD Engineer - End-to-End Guide

This guide documents the **complete** process for generating CRUD functionality from scratch (migration to API tests) for an entity in the Goravel meta-framework.

## Complete E2E Process (Migration → API Tests)

This is the full workflow for creating a new entity from scratch, including migration, model, service, controller, routes, and tests.

## Step-by-Step Process

### 0. Create Migration

**Command:**
```bash
go run . artisan make:migration create_<table_name>_table
```

**Example:**
```bash
go run . artisan make:migration create_events_table
```

**What to do:**
1. Edit the generated migration file in `database/migrations/`
2. Add table fields in the `Up()` method:
   ```go
   table.String("title", 255)
   table.Text("description")
   table.Date("date")  // For carbon.DateTime fields
   table.Json("partners")  // For []string or []int fields
   table.Text("notes").Nullable()  // Optional fields
   ```
3. Add soft deletes: `table.SoftDeletesTz()`
4. Run migration: `go run . artisan migrate`

**Important field type mappings:**
- String → `table.String("field", length)`
- Text → `table.Text("field")`
- Integer → `table.Integer("field")`
- Date → `table.Date("field")` (for carbon.DateTime)
- JSON → `table.Json("field")` (for arrays: []string, []int, etc.)
- Boolean → `table.Boolean("field")`
- Nullable → add `.Nullable()` after the field definition

### 0.1 Add Audit Fields Migration

**Command:**
```bash
go run . artisan make:audit --table=<table_name>
```

**Example:**
```bash
go run . artisan make:audit --table=events
```

**What it generates:**
- Migration file that adds: created_by, updated_by, deleted_by, ip_address, user_agent
- Indexes for audit fields

**Post-generation:**
1. Review the migration file
2. **Important**: Remove `table.Timestamp("deleted_at")` if present (SoftDeletesTz already creates it)
3. Run migration: `go run . artisan migrate`

### 0.2 Generate Model from Table

**Command:**
```bash
go run . artisan make:model --table=<table_name> <ModelName>
```

**Example:**
```bash
go run . artisan make:model --table=events Event
```

**What it generates:**
- `app/models/<model_name>.go`

**Post-generation fixes:**
1. **Fix array fields** - Generator creates them as `string`, change to proper types:
   ```go
   // WRONG (what generator creates):
   Partners string `json:"partners"`

   // CORRECT:
   Partners []string `json:"partners" db:"partners" gorm:"type:json;serializer:json"`
   AttendingSmes []int `json:"attending_smes" db:"attending_smes" gorm:"type:json;serializer:json"`
   ```

2. **Fix audit field types** - Generator may use wrong types:
   ```go
   // CORRECT types for audit fields:
   CreatedBy *int  `json:"created_by" db:"created_by"`
   UpdatedBy *int  `json:"updated_by" db:"updated_by"`
   DeletedBy *int  `json:"deleted_by" db:"deleted_by"`
   ```

3. **Verify carbon.DateTime fields** - Should be non-pointer:
   ```go
   // CORRECT:
   Date carbon.DateTime `json:"date" db:"date"`
   ```

### 0.3 Generate Service

**Command:**
```bash
go run . artisan make:svc --model=<ModelName> <entity_name>
```

**Example:**
```bash
go run . artisan make:svc --model=Event event
```

**What it generates:**
- `app/services/<entity_name>_service.go`

**Post-generation customization:**
1. **Add search fields** (required - at least one):
   ```go
   WithSearchFields("title", "description", "venue", "district", "notes")
   ```

2. **Add sort fields**:
   ```go
   WithSortFields("id", "created_at", "updated_at", "date", "title")
   ```

3. **Add filter fields** (optional):
   ```go
   WithFilterFields("district", "date")
   ```

4. **Add validation rules** (required - at least one):
   ```go
   WithValidationRules(map[string]interface{}{
       "title":       "required|max_len:255",
       "description": "required",
       "venue":       "required|max_len:255",
       "district":    "required|max_len:100",
   })
   ```

5. **Set default sort**:
   ```go
   WithDefaultSort("created_at", "DESC")
   ```

6. **Enable scope filtering**:
   ```go
   WithScopeFiltering("<service_name>", "created_by")
   ```

### 0.4 Register Permissions

**File:** `app/auth/permission_constants.go`

**Steps:**

1. **Add service constant** (around line 30):
   ```go
   ServiceEvents ServiceRegistry = "events"
   ```

2. **Register in GetAllServiceRegistries()** (around line 100):
   ```go
   func GetAllServiceRegistries() []ServiceRegistry {
       return []ServiceRegistry{
           ServiceBooks,
           ServiceEvents,  // Add your new service
           // ... other services
       }
   }
   ```

3. **Add display name in GetServiceDisplayName()** (around line 150):
   ```go
   case ServiceEvents:
       return "Events"
   ```

4. **Add actions in GetServiceActions()** (around line 200):
   ```go
   case ServiceEvents:
       return []PermissionAction{
           PermissionCreate,
           PermissionRead,
           PermissionUpdate,
           PermissionDelete,
           PermissionManage,
       }
   ```

## Request Generation & Controller Setup

### 1. Generate Request Validation Files

**Command:**
```bash
go run . artisan make:req --model=<entity_name> <entity_name>
```

**Example:**
```bash
go run . artisan make:req --model=business_formalisation business_formalisation
```

**What it generates:**
- `app/http/requests/<entity_name>_create_request.go`
- `app/http/requests/<entity_name>_update_request.go`

**Post-generation steps:**
1. Review generated struct fields and ensure they match your model
2. Add proper validation rules in `Rules()` method
3. Add custom validation messages in `Messages()` method
4. Implement `ToCreateData()` and `ToUpdateData()` methods
5. **IMPORTANT**: Exclude any read-only/calculated fields (like scores, aggregates, etc.)
6. **CRITICAL**: Handle carbon.DateTime fields properly based on model type

**Common validation rule gotchas:**
- `date` validator doesn't work with `*string` - use custom validation in `PrepareForValidation()`
- `numeric` validator doesn't work with float64 - omit validation for numeric types
- Use pointer types (`*string`, `*bool`, etc.) for optional fields in update requests

**CRITICAL: carbon.DateTime Pointer Handling**

The model's carbon.DateTime type determines how you convert dates in requests:

```go
// If model uses NON-POINTER carbon.DateTime:
type Event struct {
    Date carbon.DateTime `json:"date"`  // NON-POINTER
}

// In request ToCreateData():
if r.Date != "" {
    parsedDate := carbon.Parse(r.Date)
    if parsedDate.Error == nil {
        data["date"] = *carbon.NewDateTime(parsedDate)  // DEREFERENCE with *
    }
}

// If model uses POINTER *carbon.DateTime:
type AdditionalBusinessMember struct {
    DateOfBirth *carbon.DateTime `json:"date_of_birth"`  // POINTER
}

// In request ToCreateData():
if r.DateOfBirth != nil && *r.DateOfBirth != "" {
    parsedDate := carbon.Parse(*r.DateOfBirth)
    if parsedDate.Error == nil {
        data["date_of_birth"] = carbon.NewDateTime(parsedDate)  // NO dereference
    }
}
```

**Rule of thumb:**
- `carbon.NewDateTime()` always returns `*carbon.DateTime` (a pointer)
- If model expects non-pointer, dereference with `*`
- If model expects pointer, use as-is

### 2. Generate Controller

**Command:**
```bash
go run . artisan make:ctrl --model=<entity_name> <entity_name>
```

**Example:**
```bash
go run . artisan make:ctrl --model=business_formalisation business_formalisation
```

**What it generates:**
- `app/http/controllers/<entity_name>s/<entity_name>_controller.go`

**Post-generation fixes:**
1. Fix naming issues (generator may use underscores instead of CamelCase)
2. Update service constant to match `permission_constants.go`:
   ```go
   // WRONG (generator may create):
   auth.ServiceBusiness_formalisations

   // CORRECT (must match permission_constants.go):
   auth.ServiceBusinessFormalisation
   ```
3. Remove `GetFilterDefinitions()` method if not needed
4. Customize `SetBeforeStore()` and `SetBeforeUpdate()` hooks as needed

### 3. Register API Routes

**File:** `routes/api.go`

**Steps:**
1. Add import:
   ```go
   "smedi-sme-db/app/http/controllers/<entity_name>s"
   ```

2. Initialize controller in `Api()` function:
   ```go
   <entity>Controller := <entity_name>s.New<Entity>Controller()
   ```

3. Add optional auth routes (public with permission scoping):
   ```go
   router.Middleware(optionalAuth).Group(func(optionalAuthRouter route.Router) {
       // Use hyphenated endpoints, plural
       optionalAuthRouter.Get("/<entity-name>s", <entity>Controller.Index)
       optionalAuthRouter.Get("/<entity-name>s/search", <entity>Controller.Search)
       optionalAuthRouter.Get("/<entity-name>s/{id}", <entity>Controller.Show)
   })
   ```

4. Add protected routes (require authentication):
   ```go
   router.Middleware(jwtAuth).Group(func(protectedRouter route.Router) {
       protectedRouter.Post("/<entity-name>s", <entity>Controller.Store)
       protectedRouter.Put("/<entity-name>s/{id}", <entity>Controller.Update)
       protectedRouter.Delete("/<entity-name>s/{id}", <entity>Controller.Delete)
   })
   ```

**Endpoint naming conventions:**
- Use **hyphenated** format: `/business-formalisations`, NOT `/business_formalisations`
- Use **plural** for collections
- Use singular in controller/service names

### 4. Generate CRUD Tests

**Command:**
```bash
go run . artisan make:crud-test --controller=<Entity>Controller
```

**Example:**
```bash
go run . artisan make:crud-test --controller=BusinessFormalisationController
```

**What it generates:**
- `tests/feature/crud/<entity>_controller_crud_test.go` (template)

**Post-generation customization:**

The generated test is a template. You MUST customize it:

1. **Fix table/permission names** - Generator often gets these wrong:
   ```go
   // WRONG (what generator creates):
   orm.Query().Exec("DELETE FROM businessformalisationcontrollers")
   permission.Slug = "businessformalisationcontrollers_create"

   // CORRECT:
   orm.Query().Exec("DELETE FROM business_formalisation")
   permission.Slug = "business_formalisation_create"
   ```

2. **Fix API endpoints** - Use correct hyphenated format:
   ```go
   // CORRECT:
   s.makeRequest("POST", "/api/business-formalisations", data)
   ```

3. **Create test dependencies** - If entity has foreign keys:
   ```go
   func (s *TestSuite) setupTestData() {
       // Create parent entities first (e.g., SME)
       sme := &models.Sme{...}
       s.Nil(facades.Orm().Query().Create(sme))
       s.testSme = sme
   }
   ```

4. **Add proper test data** in `getValid<Entity>Data()`:
   ```go
   func (s *TestSuite) getValidData() map[string]interface{} {
       return map[string]interface{}{
           "sme_id": s.testSme.ID,  // Use created dependencies
           "field1": value1,
           "field2": value2,
           "array_field": []string{},  // Include empty arrays for JSON fields
           // Include ALL fields EXCEPT read-only ones
       }
   }
   ```

5. **Initialize ALL required fields when creating models directly in tests**:
   ```go
   // When creating models via ORM (not API) in tests:
   event := &models.Event{
       Title:         "Test Event",
       Description:   "Description",
       Date:          *carbon.NewDateTime(carbon.Parse("2025-12-01")),
       Venue:         "Venue",
       Partners:      []string{"Partner1"},  // MUST initialize arrays
       District:      "District",
       AttendingSmes: []int{},  // MUST initialize - NOT NULL in database
       CreatedBy:     &createdByInt,
   }
   s.Nil(facades.Orm().Query().Create(event))
   ```

   **Common error**: `NOT NULL constraint failed: table.array_field`
   - Cause: Forgetting to initialize JSON array fields
   - Solution: Always set array fields to `[]string{}` or `[]int{}` even if empty
   ```

5. **Add critical read-only field tests**:
   ```go
   // Test that read-only fields CANNOT be set via create
   func (s *TestSuite) TestReadOnlyFieldCannotBeSetViaCreate() {
       data := s.getValidData()
       data["read_only_field"] = 99  // Try to set it

       resp, result := s.makeRequest("POST", "/api/entities", data)
       s.Equal(http.StatusCreated, resp.StatusCode)

       // Verify it was NOT set
       var entity models.Entity
       facades.Orm().Query().Find(&entity, result["data"].(map[string]interface{})["id"])
       s.Equal(0, entity.ReadOnlyField, "should NOT be settable via API")
   }

   // Test that read-only fields CANNOT be updated via API
   func (s *TestSuite) TestReadOnlyFieldCannotBeUpdatedViaAPI() {
       // Create entity first
       entity := &models.Entity{ReadOnlyField: 0}
       facades.Orm().Query().Create(entity)

       // Try to update read-only field
       updateData := map[string]interface{}{
           "read_only_field": 50,
           "other_field": "value",
       }

       resp, _ := s.makeRequest("PUT", fmt.Sprintf("/api/entities/%d", entity.ID), updateData)
       s.Equal(http.StatusOK, resp.StatusCode)

       // Verify read-only field was NOT updated
       var updated models.Entity
       facades.Orm().Query().Find(&updated, entity.ID)
       s.Equal(0, updated.ReadOnlyField, "should NOT be updatable via API")
       s.Equal("value", updated.OtherField)  // Other fields should update
   }
   ```

6. **Follow existing test patterns** from similar entities (e.g., `additional_business_member_controller_crud_test.go`)

### 5. Run Tests

**Command:**
```bash
APP_ENV=testing go test -v ./tests/feature/crud -run Test<Entity>ControllerCRUDTestSuite
```

**Example:**
```bash
APP_ENV=testing go test -v ./tests/feature/crud -run TestBusinessFormalisationControllerCRUDTestSuite
```

**Common test failures and fixes:**

1. **Compilation errors about undefined types:**
   - Fix naming in controller (CamelCase vs snake_case)
   - Ensure service constant matches `permission_constants.go`

2. **Validation errors for float64/numeric fields:**
   - Remove `numeric` validator from float64 fields
   - Fields with Go type `float64` don't need string validation

3. **Validation errors for date fields:**
   - Remove `date` validator from `*string` date fields
   - Add custom validation in `PrepareForValidation()` using carbon.Parse()

4. **Permission denied errors:**
   - Verify permission slugs match exactly in test setup
   - Check permission constant is registered in `permission_constants.go`
   - Ensure `permissions:setup` was run after adding new service

5. **Foreign key constraint errors:**
   - Create parent entities in test setup before main entity
   - Delete in correct order in teardown (children before parents)

## Complete Examples

### Example 1: Events (Simple Entity with Arrays)

**Full E2E implementation from scratch:**

- **Migration:** `database/migrations/20251123050933_create_events_table.go`
- **Audit Migration:** `database/migrations/20251123051707_add_audit_fields_to_events_table.go`
- **Model:** `app/models/event.go`
- **Service:** `app/services/event_service.go`
- **Requests:** `app/http/requests/event_*.go`
- **Controller:** `app/http/controllers/events/event_controller.go`
- **Routes:** `routes/api.go` (search for "/events")
- **Tests:** `tests/feature/crud/eventcontroller_crud_test.go`

**Key features demonstrated:**
- JSON array fields (`partners` as []string, `attending_smes` as []int)
- carbon.DateTime field (non-pointer) with proper dereferencing in requests
- Search fields for text fields (title, description, venue, district, notes)
- Filter by district and date
- Default sort by date DESC
- Scoped permissions by created_by
- All 8 CRUD tests passing

**Important lessons from Events:**
1. JSON array fields MUST be initialized in tests: `AttendingSmes: []int{}`
2. carbon.DateTime non-pointer requires dereferencing: `*carbon.NewDateTime(parsedDate)`
3. Service needs at least one search field and one validation rule

### Example 2: Business Formalisation (Complex with Read-Only Fields)

See the implementation for reference:

- **Model:** `app/models/business_formalisation.go`
- **Service:** `app/services/business_formalisation_service.go`
- **Requests:** `app/http/requests/business_formalisation_*.go`
- **Controller:** `app/http/controllers/business_formalisations/business_formalisation_controller.go`
- **Routes:** `routes/api.go` (search for "business-formalisations")
- **Tests:** `tests/feature/crud/business_formalisation_controller_crud_test.go`

**Key features demonstrated:**
- Read-only field (`formalisation_score`) properly excluded from API
- Tests verify read-only field cannot be set/updated
- Foreign key relationship (requires SME)
- Proper validation for all field types
- Comprehensive test coverage (14 test cases)

## Checklist for New CRUD Entity (Complete E2E)

### Database & Model Setup
- [ ] Migration created (`make:migration`)
  - [ ] Table fields defined with correct types
  - [ ] SoftDeletesTz added
  - [ ] Migration applied successfully
- [ ] Audit fields migration created (`make:audit`)
  - [ ] Removed duplicate deleted_at if present
  - [ ] Migration applied successfully
- [ ] Model generated from table (`make:model --table=`)
  - [ ] Array fields fixed (string → []string with gorm tags)
  - [ ] Audit field types verified (*int)
  - [ ] carbon.DateTime fields verified (pointer vs non-pointer)

### Service & Permissions
- [ ] Service generated (`make:svc --model=`)
  - [ ] Search fields added (at least one)
  - [ ] Sort fields configured
  - [ ] Filter fields configured (optional)
  - [ ] Validation rules added (at least one)
  - [ ] Default sort set
  - [ ] Scope filtering enabled
- [ ] Permissions registered in permission_constants.go
  - [ ] Service constant added
  - [ ] Registered in GetAllServiceRegistries()
  - [ ] Display name added in GetServiceDisplayName()
  - [ ] Actions added in GetServiceActions()

### Controller & Routes
- [ ] Request files generated (`make:req --model=`)
  - [ ] Validation rules added
  - [ ] Read-only fields excluded
  - [ ] ToCreateData/ToUpdateData implemented
  - [ ] carbon.DateTime conversion correct (pointer vs non-pointer)
  - [ ] Array fields handled properly
- [ ] Controller generated (`make:ctrl --model=`)
  - [ ] Naming issues resolved
  - [ ] Service constant matches permissions
  - [ ] Hooks customized (SetBeforeStore/SetBeforeUpdate)
- [ ] Routes registered in routes/api.go
  - [ ] Import added
  - [ ] Controller initialized
  - [ ] Optional auth routes added (GET endpoints)
  - [ ] Protected routes added (POST/PUT/DELETE)
  - [ ] Hyphenated endpoint names used

### Testing
- [ ] Tests generated (`make:crud-test --controller=`)
  - [ ] Table/permission names fixed
  - [ ] API endpoint paths corrected (hyphenated)
  - [ ] Test data created properly
  - [ ] Array fields initialized in test models
  - [ ] carbon.DateTime fields created correctly
  - [ ] Dependencies created in setup (if foreign keys)
  - [ ] Read-only field tests added (if applicable)
  - [ ] All tests passing (`APP_ENV=testing go test -v ...`)

## Tips and Best Practices

1. **Always use artisan commands** - Don't manually create files
2. **Follow existing patterns** - Look at working examples like `additional_business_member`
3. **Test early and often** - Run tests after each major step
4. **Use hyphenated endpoints** - `/entity-names` not `/entity_names`
5. **Validate read-only fields** - Critical for data integrity
6. **Clean up test data** - Proper setup/teardown prevents flaky tests
7. **Document special cases** - Add comments for non-obvious logic

## Common Patterns

### Read-Only Fields
Fields that should never be set via API (calculated, aggregates, etc.):
- Exclude from both create and update request structs
- Add tests to verify they cannot be set/updated
- Can be set internally by services/listeners

### Foreign Key Relationships
When entity references another:
- Create parent entity in test setup
- Use parent ID in test data
- Delete children before parents in teardown

### Optional vs Required Fields
- Create request: Use non-pointer for required, pointer for optional
- Update request: Use pointers for all fields (all optional in partial updates)
- Update request: Only include provided fields in `ToUpdateData()`

### Boolean Fields
- Default to false in database
- Include in request structs as `bool` (create) or `*bool` (update)
- Use `boolean` validator (works correctly)

### Date Fields
- Use `*string` in request structs
- Use `carbon.DateTime` in model
- NO `date` validator - add custom validation in `PrepareForValidation()`
- Parse with carbon.Parse() in `ToCreateData()`/`ToUpdateData()`

### Numeric Fields (float64, int)
- Use native types in request structs
- NO `numeric` validator needed for Go numeric types
- Only use `numeric` for string representations of numbers
