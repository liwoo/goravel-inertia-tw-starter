---
name: inertia-page-ctrl
description: Create a Go page controller that serves an Inertia.js page. Configures the EntityPageController with service, permissions, and optional stats.
argument-hint: "[EntityName]"
allowed-tools: Bash, Read, Write, Edit, Grep, Glob
---

# Inertia Page Controller (Go Backend)

Create page controller for `$ARGUMENTS`.

## Step 1: Generate Page Controller

```bash
go run . artisan make:page-ctrl --controller=$ARGUMENTS
```

This creates `app/http/controllers/<entity>s/<entity>s_page_controller.go`.

## Step 2: Configure the Controller

```go
package entitynames

import (
    "smedi-sme-db/app/auth"
    "smedi-sme-db/app/contracts"
    "smedi-sme-db/app/services"
)

type EntityPageController struct {
    *contracts.GenericPageController
    entityService *services.EntityService
}

func NewEntityPageController() *EntityPageController {
    entityService := services.NewEntityService()

    return &EntityPageController{
        GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
            ResourceType:      "entities",               // Plural snake_case
            PageComponent:     "Entity/Index",           // React component path under resources/js/pages/
            Service:           entityService,
            ServiceIdentifier: auth.ServiceEntity,       // From permission_constants.go
            StatsEnabled:      false,                    // Set true for badge counts
            // StatsBuilder: func(controller *contracts.GenericPageController) map[string]interface{} {
            //     stats, _ := entityService.GetEntityStatistics()
            //     return stats
            // },
        }),
        entityService: entityService,
    }
}
```

## Configuration Fields

| Field | Purpose | Example |
|---|---|---|
| `ResourceType` | Plural snake_case resource name | `"configs"`, `"books"` |
| `PageComponent` | React component path | `"Config/Index"`, `"Book/Index"` |
| `Service` | CRUD service instance | `configService` |
| `ServiceIdentifier` | Permission service constant | `auth.ServiceConfig` |
| `StatsEnabled` | Enable statistics | `true` / `false` |
| `StatsBuilder` | Function returning stats map | See template above |

## Props Passed to React

The `GenericPageController.Index()` method automatically passes:

```typescript
{
    data: PaginatedResult<Entity>,     // Paginated entity list
    filters: ListRequest,               // Current page/sort/search/filter state
    permissions: {
        canCreate: boolean,
        canEdit: boolean,
        canDelete: boolean,
    },
    stats?: Record<string, any>,        // If StatsEnabled=true
    meta: {
        pagination: {
            defaultPageSize: number,
            maxPageSize: number,
            allowedSizes: number[],
        },
    },
}
```

## Step 3: Register Web Route

In `routes/web.go`:

```go
// Import
entityPageController := entitynames.NewEntityPageController()

// Route (inside authenticated middleware group)
router.Get("/admin/entity-names", entityPageController.Index)
```

## Step 4: Adding Statistics (Optional)

If you want badge counts on simple filters, create a statistics method on the service:

```go
func (s *EntityService) GetEntityStatistics() (map[string]interface{}, error) {
    var totalCount, activeCount, inactiveCount int64

    totalCount, _ = facades.Orm().Query().Model(&models.Entity{}).Count()
    activeCount, _ = facades.Orm().Query().Model(&models.Entity{}).Where("status = ?", "ACTIVE").Count()
    inactiveCount, _ = facades.Orm().Query().Model(&models.Entity{}).Where("status = ?", "INACTIVE").Count()

    return map[string]interface{}{
        "totalCount":    totalCount,
        "activeCount":   activeCount,
        "inactiveCount": inactiveCount,
    }, nil
}
```

Then enable in page controller:
```go
StatsEnabled: true,
StatsBuilder: func(controller *contracts.GenericPageController) map[string]interface{} {
    stats, _ := entityService.GetEntityStatistics()
    return stats
},
```

## Reference

See `app/http/controllers/configs/configs_page_controller.go` for a complete example.
