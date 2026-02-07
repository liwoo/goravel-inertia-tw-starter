---
name: goravel-crud-permissions
description: Register permissions for a new Goravel entity in the permission system. Adds service constant, display name, and actions.
argument-hint: "[ServiceName] [DisplayName]"
allowed-tools: Bash, Read, Write, Edit, Grep, Glob
---

# Goravel CRUD Permissions

Register permissions for `$ARGUMENTS`.

## File to Edit

`app/auth/permission_constants.go`

## Step 1: Add Service Registry Constant

Add to the `ServiceRegistry` const block (around line 26):

```go
const (
    ServiceBooks        ServiceRegistry = "books"
    ServiceUsers        ServiceRegistry = "users"
    // ... existing services
    ServiceYourEntity   ServiceRegistry = "your_entities"  // ADD THIS
)
```

**Naming convention**: `Service` + PascalCase entity name. Value is lowercase plural snake_case.

## Step 2: Register in GetAllServiceRegistries()

Add your service to the return slice (around line 52):

```go
func GetAllServiceRegistries() []ServiceRegistry {
    return []ServiceRegistry{
        ServiceBooks,
        ServiceUsers,
        // ... existing services
        ServiceYourEntity,  // ADD THIS
    }
}
```

## Step 3: Add Display Name

Add a case in `GetServiceDisplayName()` (around line 70):

```go
case ServiceYourEntity:
    return "Your Entity Management"
```

## Step 4: Define Actions

Add a case in `GetServiceActions()` (around line 125):

```go
case ServiceYourEntity:
    return []CorePermissionAction{
        PermissionCreate,
        PermissionRead,
        PermissionUpdate,
        PermissionDelete,
        PermissionView,
    }
```

Common action sets:
- **Basic CRUD**: Create, Read, Update, Delete, View
- **Full management**: Add PermissionManage, PermissionExport
- **Bulk operations**: Add PermissionBulkUpdate, PermissionBulkDelete

## Step 5: Sync Permissions to Database

```bash
go run . artisan permissions:setup
```

This creates the permission records in the database. Then:
1. Navigate to `/admin/permissions` as super admin
2. Assign permissions to appropriate roles
3. Test with different user accounts

## Verify

Check that the permissions were created:

```bash
go run . artisan db:table permissions
```

Look for entries like `your_entities_create`, `your_entities_read`, etc.

## Current Services in the System

Reference `app/auth/permission_constants.go` for existing patterns:
- `ServiceBooks` = "books"
- `ServiceUsers` = "users"
- `ServiceConfig` = "config"
- `ServiceLenders` = "lenders"
- `ServiceApplications` = "applications"

## Verify

After editing `permission_constants.go`:

```bash
# Vet the auth package for issues
go vet ./app/auth/...

# Confirm full project compiles (catches typos in constants)
go build ./...
```

## Next Step

Run `/goravel-crud-request` to generate request validators.
