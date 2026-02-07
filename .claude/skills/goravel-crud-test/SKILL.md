---
name: goravel-crud-test
description: Generate and fix comprehensive CRUD tests for a Goravel entity. Covers all CRUD operations, sorting, pagination, search, and permissions.
argument-hint: "[EntityName]"
allowed-tools: Bash, Read, Write, Edit, Grep, Glob
---

# Goravel CRUD Test Generator

Generate tests for `$ARGUMENTS`.

## Step 1: Generate Test Template

```bash
go run . artisan make:crud-test --controller=$ARGUMENTS
```

This creates `tests/feature/crud/<entity>_crud_test.go`.

## Step 2: Fix Generated Test (CRITICAL)

The generator produces a template with known issues. Fix in this order:

### Fix 1: Permission Names

```go
// WRONG (generator creates):
Slug: fmt.Sprintf("entitycontrollers_%s", perm)

// CORRECT (use service name from permission_constants.go):
Slug: fmt.Sprintf("entity_%s", perm)
```

Check `app/auth/permission_constants.go` for the exact `ServiceRegistry` value.

### Fix 2: API Endpoints

```go
// WRONG (generator creates):
s.makeRequest("POST", "/api/entitycontrollers", data)

// CORRECT (use hyphenated plural from routes/api.go):
s.makeRequest("POST", "/api/entity-names", data)
```

### Fix 3: Table Names in Cleanup

```go
// WRONG:
orm.Query().Exec("DELETE FROM entitycontrollers")

// CORRECT:
orm.Query().Exec("DELETE FROM entity_names")
```

### Fix 4: Add Valid Test Data

```go
func (s *TestSuite) TestCreateEntity() {
    data := map[string]interface{}{
        "title":       "Test Entity",
        "description": "A test description",
        "status":      "ACTIVE",
        "tags":        []string{"tag1", "tag2"},  // Include arrays
    }
    resp, result := s.makeRequest("POST", "/api/entities", data)
    s.Equal(http.StatusCreated, resp.StatusCode)
}
```

### Fix 5: Initialize ALL Required Fields for ORM-Created Models

```go
entity := &models.Entity{
    Title:       "Test",
    Description: "Test",
    Tags:        []string{},  // MUST initialize arrays - NOT NULL constraint
    CreatedBy:   &createdByInt,
}
```

**Common error**: `NOT NULL constraint failed: table.array_field`
**Fix**: Always initialize JSON array fields to `[]string{}` or `[]int{}`.

### Fix 6: carbon.DateTime in Tests

```go
// Create date for test data:
date := *carbon.NewDateTime(carbon.Parse("2025-12-01"))

entity := &models.Entity{
    Date: date,  // Non-pointer carbon.DateTime
}
```

### Fix 7: Foreign Key Dependencies

If entity has foreign keys, create parents in setup:

```go
func (s *TestSuite) SetupTest() {
    s.RefreshDatabase()
    s.setupTestUser()

    // Create parent entity
    parent := &models.Parent{Name: "Test Parent"}
    s.Nil(facades.Orm().Query().Create(parent))
    s.testParent = parent
}
```

## Step 3: Add Read-Only Field Tests (if applicable)

```go
func (s *TestSuite) TestReadOnlyFieldCannotBeSet() {
    data := s.getValidData()
    data["score"] = 99  // Try to set read-only field

    resp, result := s.makeRequest("POST", "/api/entities", data)
    s.Equal(http.StatusCreated, resp.StatusCode)

    var entity models.Entity
    facades.Orm().Query().Find(&entity, result["data"].(map[string]interface{})["id"])
    s.Equal(0, entity.Score, "read-only field should NOT be settable")
}
```

## Step 4: Run Tests (TDD Loop)

```bash
APP_ENV=testing go test -v ./tests/feature/crud -run Test<Entity>CRUDTestSuite
```

Fix each failure, re-run, iterate until all pass.

## Test Coverage Checklist

- [ ] Create with valid data
- [ ] Create validation errors (missing required fields)
- [ ] Get by ID (found)
- [ ] Get by ID (not found - 404)
- [ ] Update with valid data
- [ ] Delete (soft delete)
- [ ] Pagination (default page size, custom page size)
- [ ] Sorting (ascending, descending)
- [ ] Search by keyword
- [ ] Read-only fields cannot be set/updated (if applicable)
- [ ] Foreign key constraints (if applicable)

## Next Step

After all tests pass, run `/goravel-crud-page` to generate the UI, or `/goravel-crud-nav` if backend is complete.
