---
name: goravel-crud-controller
description: Generate API controller for a Goravel entity with permission checking and Swagger annotations. Use after creating request validators.
argument-hint: "[entity_name]"
allowed-tools: Bash, Read, Write, Edit, Grep, Glob
---

# Goravel CRUD Controller Generator

Generate controller for `$ARGUMENTS`.

## Step 1: Generate Controller

```bash
go run . artisan make:ctrl --model=<entity_name> <entity_name>
```

This creates `app/http/controllers/<entity_name>s/<entity_name>_controller.go`.

## Step 2: Post-Generation Fixes

### Fix 1: Naming Issues

Generator may use underscores instead of CamelCase:

```go
// WRONG (generator may create):
auth.ServiceBusiness_formalisations

// CORRECT (must match permission_constants.go):
auth.ServiceBusinessFormalisation
```

### Fix 2: Service Constant Must Match

The `auth.Service*` constant in the controller MUST match what you registered in `permission_constants.go`:

```go
// Check permission_constants.go for the exact constant name
crudController := contracts.NewCrudController[
    models.Entity,
    *requests.EntityCreateRequest,
    *requests.EntityUpdateRequest,
](
    "entity",          // resource name for logging
    entityService,     // service instance
).
    WithAuthChecker(func(ctx http.Context, action string, resource interface{}) error {
        permHelper := auth.GetPermissionHelper()
        switch action {
        case "viewAny", "view":
            if !permHelper.CheckServicePermission(ctx, auth.ServiceEntity, auth.PermissionRead) {
                return fmt.Errorf("unauthorized")
            }
        case "create":
            if !permHelper.CheckServicePermission(ctx, auth.ServiceEntity, auth.PermissionCreate) {
                return fmt.Errorf("unauthorized")
            }
        case "update":
            if !permHelper.CheckServicePermission(ctx, auth.ServiceEntity, auth.PermissionUpdate) {
                return fmt.Errorf("unauthorized")
            }
        case "delete":
            if !permHelper.CheckServicePermission(ctx, auth.ServiceEntity, auth.PermissionDelete) {
                return fmt.Errorf("unauthorized")
            }
        }
        return nil
    }).
    Build()
```

### Fix 3: Constructor Function

```go
func NewEntityController() *EntityController {
    entityService := services.NewEntityService()
    // ... build crudController
    return &EntityController{
        CrudController: crudController,
    }
}
```

## Reference Pattern

See `app/http/controllers/books/book_controller.go` for:
- Generic controller construction
- Permission checking via `WithAuthChecker`
- Custom endpoint methods (beyond CRUD)
- Swagger annotations

## Controller Methods (Auto-Provided by Framework)

The `CrudController` automatically provides:
- `Index` - List with pagination, sorting, filtering
- `Show` - Get by ID
- `Store` - Create new record
- `Update` - Update existing record
- `Delete` - Soft delete
- `Search` - Keyword search
- `FilterMetadata` - Available filter options

## Adding Custom Endpoints

```go
func (c *EntityController) Statistics(ctx http.Context) http.Response {
    // Check permission
    permHelper := auth.GetPermissionHelper()
    if !permHelper.CheckServicePermission(ctx, auth.ServiceEntity, auth.PermissionRead) {
        return ctx.Response().Status(403).Json(map[string]interface{}{
            "error": "Unauthorized",
        })
    }

    stats, err := c.entityService.GetEntityStatistics()
    if err != nil {
        return ctx.Response().Status(500).Json(map[string]interface{}{
            "error": err.Error(),
        })
    }

    return ctx.Response().Json(200, map[string]interface{}{
        "success": true,
        "data":    stats,
    })
}
```

## Verify

After fixing the controller:

```bash
# Vet the controller package
go vet ./app/http/controllers/...

# Confirm full project compiles (catches service constant mismatches)
go build ./...
```

## Next Step

Run `/goravel-crud-routes` to register the API routes.
