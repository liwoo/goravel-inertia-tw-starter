---
name: goravel-crud-search
description: Add a new entity to the global search (CMD+K) system. Updates backend search controller with permission-gated search and frontend search config with entity type, icon, and colors.
argument-hint: "[EntityName] [icon-name]"
allowed-tools: Read, Write, Edit, Grep, Glob
---

# Goravel Global Search Integration (CMD+K)

Add `$ARGUMENTS` to global search.

## Architecture Overview

The global search (CMD+K) has three layers:
1. **Backend**: `search_controller.go` — permission-gated search methods per entity
2. **Frontend config**: `search_config.tsx` — entity type, icon, colors, permissions
3. **Frontend UI**: `GlobalSearch.tsx` — command palette (rarely needs changes)

Adding a new entity only requires changes to layers 1 and 2.

## File 1: Backend Search — `app/http/controllers/search_controller.go`

### Step 1: Add Permission Check in `GlobalSearch()`

Find the `results := []SearchResult{}` line and add a permission-gated search block:

```go
// Search Entities if user has permission
if permHelper.CheckServicePermission(ctx, auth.ServiceEntity, auth.PermissionRead) {
    entityResults := c.searchEntities(query)
    results = append(results, entityResults...)
}
```

**Important**: `auth.ServiceEntity` must match the constant registered in `app/auth/permission_constants.go`.

### Step 2: Add Search Method

Add a new method to `SearchController`. Follow the existing pattern:

```go
// searchEntities performs fuzzy search on entities using the EntityService
func (c *SearchController) searchEntities(query string) []SearchResult {
    results := []SearchResult{}

    entityService := services.NewEntityService()

    // Use the service's search functionality
    paginatedResult, err := entityService.Search(query, contracts.ListRequest{
        Page:     1,
        PageSize: 10,
    })

    if err != nil || paginatedResult == nil {
        return results
    }

    // Convert service results to search results
    for _, item := range paginatedResult.Data {
        if entity, ok := item.(models.Entity); ok {
            results = append(results, SearchResult{
                ID:       entity.ID,
                Title:    entity.Name,           // Primary display field
                Subtitle: entity.Description,     // Secondary display field (optional)
                Type:     "entity",               // Must match frontend SearchEntityType
                URL:      fmt.Sprintf("/admin/entity-names?search=%s", url.QueryEscape(entity.Name)),
            })
        }
    }

    return results
}
```

### Field Selection Guide

| SearchResult field | Purpose | Example |
|---|---|---|
| `Title` | Primary text shown in results | `entity.Name`, `user.Email` |
| `Subtitle` | Secondary text below title | `entity.Status`, combined fields with `fmt.Sprintf` |
| `Type` | Entity type identifier (must match frontend) | `"entity"`, `"book"`, `"user"` |
| `URL` | Navigation URL when result is clicked | `/admin/entity-names?search=<encoded>` |

### Subtitle Composition Pattern

For entities with multiple display fields, compose subtitles:

```go
subtitle := entity.Status
if entity.Category != "" {
    subtitle = fmt.Sprintf("%s • %s", entity.Category, entity.Status)
}
```

For entities with nullable fields:

```go
subtitle := entity.Type
if entity.Description != nil && *entity.Description != "" {
    subtitle = fmt.Sprintf("%s • %s", entity.Type, *entity.Description)
}
```

### Step 3: Verify Imports

Ensure these are imported in `search_controller.go`:

```go
import (
    "books-database/app/models"
    "books-database/app/services"
)
```

The `fmt`, `net/url`, `strings` imports and `auth`, `contracts` imports should already be present.

## File 2: Frontend Search Config — `resources/js/config/search_config.tsx`

### Step 1: Add Entity to `SearchEntityType` Union

```typescript
export type SearchEntityType = 'user' | 'config' | 'application' | 'entity';
```

### Step 2: Add Icon Import

```typescript
import {
    Users,
    FileText,
    Settings,
    YourIcon,  // Add from lucide-react
} from 'lucide-react';
```

Browse icons at https://lucide.dev/icons

### Step 3: Add Entity Config to `SEARCH_ENTITIES` Array

```typescript
{
    type: 'entity',
    label: 'Entities',
    icon: <YourIcon className="h-4 w-4" />,
    permissionService: 'entities',      // Must match ServiceRegistry in permission_constants.go
    permissionAction: 'read',
    colors: {
        light: 'bg-blue-100 text-blue-800',
        dark: 'dark:bg-blue-900/30 dark:text-blue-400',
    },
    urlPrefix: '/admin/entity-names',
},
```

### SearchEntityConfig Fields

| Field | Purpose | Example |
|---|---|---|
| `type` | Matches `SearchResult.Type` from backend | `'entity'` |
| `label` | Display name in search UI (quick access, no-results badges) | `'Entities'` |
| `icon` | Lucide icon JSX element | `<BookOpen className="h-4 w-4" />` |
| `permissionService` | Service name for permission gating | `'entities'` |
| `permissionAction` | Required permission action | `'read'` |
| `colors.light` | Light mode badge colors | `'bg-blue-100 text-blue-800'` |
| `colors.dark` | Dark mode badge colors | `'dark:bg-blue-900/30 dark:text-blue-400'` |
| `urlPrefix` | Base URL for entity pages | `'/admin/entity-names'` |

### Available Color Palettes

Already used:
- **cyan**: Users — `bg-cyan-100 text-cyan-800`
- **gray**: Configs — `bg-gray-100 text-gray-800`
- **amber**: Applications — `bg-amber-100 text-amber-800`

Pick from unused colors:
- **blue**: `bg-blue-100 text-blue-800` / `dark:bg-blue-900/30 dark:text-blue-400`
- **green**: `bg-green-100 text-green-800` / `dark:bg-green-900/30 dark:text-green-400`
- **purple**: `bg-purple-100 text-purple-800` / `dark:bg-purple-900/30 dark:text-purple-400`
- **rose**: `bg-rose-100 text-rose-800` / `dark:bg-rose-900/30 dark:text-rose-400`
- **indigo**: `bg-indigo-100 text-indigo-800` / `dark:bg-indigo-900/30 dark:text-indigo-400`
- **emerald**: `bg-emerald-100 text-emerald-800` / `dark:bg-emerald-900/30 dark:text-emerald-400`
- **orange**: `bg-orange-100 text-orange-800` / `dark:bg-orange-900/30 dark:text-orange-400`
- **teal**: `bg-teal-100 text-teal-800` / `dark:bg-teal-900/30 dark:text-teal-400`
- **violet**: `bg-violet-100 text-violet-800` / `dark:bg-violet-900/30 dark:text-violet-400`
- **pink**: `bg-pink-100 text-pink-800` / `dark:bg-pink-900/30 dark:text-pink-400`

## How Search Works (Reference)

1. User presses **CMD+K** (or **CTRL+K**) → `GlobalSearch` dialog opens
2. User types query → debounced (300ms) → `GET /api/search?q=<query>`
3. Backend `GlobalSearch()` checks permissions per entity, calls `searchEntities()` for each
4. Each `searchEntities()` uses the entity's service `Search()` method (ILIKE with `%query%` wildcards)
5. Results returned as `SearchResult[]` → frontend renders with icons, badges, keyboard nav
6. Frontend filters `SEARCH_ENTITIES` by user permissions for quick access section
7. Clicking a result or pressing Enter navigates to `result.url`

### Service Search Requirements

For the backend search to work, the entity's service must have search fields configured:

```go
// In the service builder (e.g., NewEntityService)
func NewEntityService() *EntityService {
    return &EntityService{
        GenericCrudService: contracts.NewServiceBuilder[models.Entity]("entities", "id").
            WithSearchFields([]string{"name", "description", "status"}).  // Required for search!
            // ... other config
            Build(),
    }
}
```

Without `WithSearchFields`, the service `Search()` method won't match any results.

## Verification

1. **Backend compiles**: `go build ./...`
2. **Type in `SearchEntityType` matches `SearchResult.Type`** from backend
3. **Permission service name matches** `ServiceRegistry` in `permission_constants.go`
4. **Test CMD+K search**: type a query, verify new entity results appear
5. **Test quick access**: open CMD+K without typing, verify entity shows in quick access
6. **Test permissions**: log in as user without permission, verify entity is hidden

## Reference

See existing search implementations in `search_controller.go`:
- `searchUsers()` — simple title + subtitle
- `searchConfigs()` — nullable field handling in subtitle
- `searchApplications()` — composed subtitle with multiple fields
