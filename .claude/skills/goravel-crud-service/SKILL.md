---
name: goravel-crud-service
description: Generate the service layer for a Goravel entity with builder pattern configuration. Use after model generation.
argument-hint: "[ModelName] [entity_name]"
allowed-tools: Bash, Read, Write, Edit, Grep, Glob
---

# Goravel CRUD Service Generator

Generate service for `$ARGUMENTS`.

## Step 1: Generate Service

```bash
go run . artisan make:svc --model=<ModelName> <entity_name>
```

This creates `app/services/<entity_name>_service.go`.

## Step 2: Configure the Builder Pattern

The generated service uses `contracts.NewServiceBuilder`. You MUST configure these:

### Required Configuration

```go
service := contracts.NewServiceBuilder[models.ModelName]("table_name", "id").
    // REQUIRED: At least one search field
    WithSearchFields("title", "description", "name").
    // REQUIRED: Define sortable columns
    WithSortFields("id", "created_at", "updated_at", "title").
    // REQUIRED: At least one validation rule
    WithValidationRules(map[string]interface{}{
        "title": "required|max_len:255",
        "name":  "required|max_len:100",
    }).
    Build()
```

### Optional Configuration

```go
    // Filter fields (for URL param-based filtering)
    WithFilterFields("status", "category", "district").
    // Default sort order
    WithDefaultSort("created_at", "DESC").
    // Eager load relationships
    WithRelations("Creator", "Updater").
    // Enable scoped permissions (user/role/global)
    WithScopeFiltering("service_name", "created_by").
    // Lifecycle hooks
    WithBeforeCreate(func(data map[string]interface{}) error {
        // Custom pre-create logic
        return nil
    }).
    WithAfterCreate(func(model *models.ModelName) error {
        // Custom post-create logic
        return nil
    }).
    WithBeforeUpdate(func(id uint, data map[string]interface{}) error {
        return nil
    }).
    // Custom query modifications
    WithCustomQuery(func(query orm.Query) orm.Query {
        return query.Where("is_active = ?", true)
    }).
```

## Reference Pattern

See `app/services/book_service.go` for the canonical service implementation.

Key things to match:
- `NewServiceBuilder[models.Book]("books", "id")` - generic type + table + primary key
- All required `With*` methods called
- Service constructor returns the concrete type

## Column Mapping (if needed for camelCase sorting)

If your frontend sends camelCase sort fields, add a column mapping:

```go
func (s *EntityService) GetColumnMapping() map[string]string {
    mapping := s.CrudServiceContract.GetColumnMapping()
    mapping["configType"] = "config_type"
    mapping["createdAt"] = "created_at"
    mapping["updatedAt"] = "updated_at"
    return mapping
}
```

## Statistics Method (optional, for dashboard/filters)

```go
func (s *EntityService) GetEntityStatistics() (map[string]interface{}, error) {
    var totalCount int64
    totalCount, _ = facades.Orm().Query().Model(&models.Entity{}).Count()

    return map[string]interface{}{
        "totalCount": totalCount,
    }, nil
}
```

## Verify

After configuring the service builder:

```bash
# Check for issues in the service
go vet ./app/services/...

# Confirm full project compiles
go build ./...
```

## Next Step

Run `/goravel-crud-permissions` to register permissions for this entity.
