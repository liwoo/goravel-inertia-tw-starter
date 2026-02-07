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

The module path is `smedi-sme-db`.

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

### Step 4: Add Protected Routes (Auth-Required Mutations)

Add inside the `router.Middleware(jwtAuth, require2FA).Group(...)` block:

```go
// Entity routes
protectedRouter.Post("/<entity-names>", entityController.Store)
protectedRouter.Put("/<entity-names>/{id}", entityController.Update)
protectedRouter.Delete("/<entity-names>/{id}", entityController.Delete)
```

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

## Next Step

Run `/goravel-crud-test` to generate comprehensive CRUD tests, OR run `/goravel-crud-page` if you need the UI.
