---
name: goravel-crud-page
description: Generate page controller and UI files for a Goravel entity's admin page. Creates Inertia page controller and React components.
argument-hint: "[EntityName]"
allowed-tools: Bash, Read, Write, Edit, Grep, Glob
---

# Goravel CRUD Page Generator

Generate page controller and UI for `$ARGUMENTS`.

## Step 1: Generate Page Controller

```bash
go run . artisan make:page-ctrl --controller=$ARGUMENTS
```

This creates `app/http/controllers/<entity>s/<entity>s_page_controller.go`.

### Post-Generation Configuration

Edit the page controller to configure:

```go
func NewEntityPageController() *EntityPageController {
    entityService := services.NewEntityService()

    return &EntityPageController{
        GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
            ResourceType:      "entities",           // Plural snake_case
            PageComponent:     "Entity/Index",       // React component path
            Service:           entityService,
            ServiceIdentifier: auth.ServiceEntity,   // From permission_constants.go
            StatsEnabled:      false,                // Set true if you have statistics
            // Optional: Stats builder
            // StatsBuilder: func(controller *contracts.GenericPageController) map[string]interface{} {
            //     stats, _ := entityService.GetEntityStatistics()
            //     return stats
            // },
        }),
        entityService: entityService,
    }
}
```

## Step 2: Generate UI Files

```bash
go run . artisan make:ui --page=$ARGUMENTS --request=$ARGUMENTS
```

This creates the React component hierarchy:

```
resources/js/pages/<Entity>/
├── Index.tsx                    # Main page component
└── sections/
    ├── <Entity>Columns.tsx      # Table column definitions
    ├── <Entity>CreateForm.tsx   # Create form
    ├── <Entity>DetailView.tsx   # Detail/view modal
    ├── <Entity>EditForm.tsx     # Edit form
    ├── <Entity>PageConfig.tsx   # Page configuration
    └── index.ts                 # Barrel exports
```

## Step 3: Register Web Route

In `routes/web.go`, add the page route inside the authenticated group:

```go
// Import
entityPageController := entitynames.NewEntityPageController()

// Route (inside authenticated group)
router.Get("/admin/<entity-names>", entityPageController.Index)
```

## Step 4: Post-Generation UI Fixes

### Fix Column Definitions

Edit `<Entity>Columns.tsx` to match your model fields:

```tsx
export const entityColumns: ColumnDef<Entity>[] = [
    { accessorKey: "title", header: "Title" },
    { accessorKey: "status", header: "Status" },
    { accessorKey: "created_at", header: "Created", cell: ({ row }) => formatDate(row.original.created_at) },
];
```

### Fix Form Fields

Edit create/edit forms to include correct field types:
- Text inputs for strings
- Select dropdowns for enums (use `/goravel-enum` to generate options)
- Date pickers for dates
- Checkboxes for booleans

### Fix TypeScript Types

Create/update `resources/js/types/<entity>.ts`:

```typescript
export interface Entity extends BaseModel {
    title: string;
    description: string;
    status: string;
    tags: string[];
}
```

## Step 5: Add Simple Filters (Optional)

If entity has enum fields, add filter buttons in PageConfig:

```tsx
export const entitySimpleFilters = (stats: any): SimpleFilterConfig[] => [
    {
        key: 'status-active',
        label: 'Active',
        value: 'ACTIVE',
        badge: stats?.activeCount || 0,
        filterParams: { status: 'ACTIVE' },
    },
];
```

## Verify

After all page generation steps:

```bash
# Backend compiles (page controller + web route)
go build ./...

# Frontend compiles (generated UI components)
npx tsc --noEmit
```

## Next Step

Run `/goravel-crud-nav` to add navigation and search integration.
