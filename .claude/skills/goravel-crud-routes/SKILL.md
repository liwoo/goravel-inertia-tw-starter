---
name: goravel-crud-routes
description: Register API and web routes for a Goravel entity. Adds endpoints to routes/api.go and routes/web.go.
argument-hint: "[entity_name]"
allowed-tools: Read, Write, Edit, Grep, Glob
---

# Goravel CRUD Routes

Register routes for `$ARGUMENTS`.

## File 1: `routes/api.go`

### Step 1: Add Import

```go
import (
    // ... existing imports
    "<module>/app/http/controllers/<entity_name>s"
)
```

The module path is `books-database`.

### Step 2: Initialize Controller

Add inside the `Api()` function, with other controller initializations:

```go
entityController := entitynames.NewEntityController()
```

### Step 3: Add Optional Auth Routes (Public GET Endpoints)

Add inside the `router.Middleware(optionalAuth).Group(...)` block:

```go
// Entity routes
optionalAuthRouter.Get("/<entity-names>", entityController.Index)
optionalAuthRouter.Get("/<entity-names>/search", entityController.Search)
optionalAuthRouter.Get("/<entity-names>/filters", entityController.FilterMetadata)
optionalAuthRouter.Get("/<entity-names>/{id}", entityController.Show)
```

### Step 4: Add Protected Routes (Auth-Required Mutations + Statistics)

Add inside the `router.Middleware(jwtAuth, require2FA).Group(...)` block:

```go
// Entity routes
protectedRouter.Get("/<entity-names>/statistics", entityController.Statistics)  // Must be before {id}
protectedRouter.Post("/<entity-names>", entityController.Store)
protectedRouter.Put("/<entity-names>/{id}", entityController.Update)
protectedRouter.Delete("/<entity-names>/{id}", entityController.Delete)
```

**Statistics endpoints go in `protectedRouter`**, not `optionalAuthRouter` — they expose aggregate data that should require authentication.

### Endpoint Naming Conventions

- Use **hyphenated** format: `/business-formalisations` NOT `/business_formalisations`
- Use **plural** for collections: `/books`, `/lenders`, `/configs`
- Search endpoint MUST come before `{id}` to avoid route conflicts

## File 2: `routes/web.go` (for Inertia pages)

### Step 1: Add Import

```go
import (
    // ... existing imports
    "<module>/app/http/controllers/<entity_name>s"
)
```

### Step 2: Initialize Page Controller

```go
entityPageController := entitynames.NewEntityPageController()
```

### Step 3: Add Page Route

Inside the authenticated routes group:

```go
router.Get("/admin/<entity-names>", entityPageController.Index)
```

## Current Route Structure Reference

See `routes/api.go` and `routes/web.go` for existing patterns:

**API routes** (`routes/api.go`):
- Optional auth group: GET endpoints (Index, Search, FilterMetadata, Show)
- Protected group: POST/PUT/DELETE endpoints (Store, Update, Delete)

**Web routes** (`routes/web.go`):
- All admin pages under `router.Get("/admin/...")`

## Verify

After registering all routes:

```bash
# This is critical — catches import errors, missing controllers, wrong method signatures
go build ./...
```

If the build fails, check:
- Import path matches the controller package name
- Controller method signatures match route expectations (GET/POST/PUT/DELETE)
- Route ordering is correct (search before `{id}`)

## Next Step

Run `/goravel-crud-test` to generate and run comprehensive CRUD tests. **All API tests MUST pass before starting any UI/frontend work.** This catches data-flow bugs (Bind, validation keys, GORM mapping) that are much harder to debug through the UI.
