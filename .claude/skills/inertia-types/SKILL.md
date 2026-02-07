---
name: inertia-types
description: Create TypeScript type definitions and i18n translation namespace for a Goravel entity. Maps Go model fields to TS interfaces and creates the translation JSON file.
argument-hint: "[EntityName]"
allowed-tools: Read, Write, Edit, Grep, Glob
---

# Inertia TypeScript Types & i18n Namespace

Create types and translations for `$ARGUMENTS`.

## Step 1: TypeScript Type File

Create `resources/js/types/<entity_name>.ts`:

```typescript
// =============================================
// Entity Type (matches Go model)
// =============================================
export interface Entity {
    id: number;
    name: string;
    description?: string;
    // status: EntityStatus;
    // price: number;
    // tags: string[];

    // Dual-case audit fields (Go sends snake_case, TS uses camelCase)
    createdAt?: string;
    created_at?: string;
    updatedAt?: string;
    updated_at?: string;
    deletedAt?: string;
    deleted_at?: string;
    createdBy?: number;
    created_by?: number;
    updatedBy?: number;
    updated_by?: number;
}

// =============================================
// Create/Update Data (sent to API)
// =============================================
export interface EntityCreateData {
    name: string;
    description?: string;
    // status?: EntityStatus;
}

export interface EntityUpdateData {
    name?: string;
    description?: string;
    // status?: EntityStatus;
}

// =============================================
// List Response (from paginated API)
// =============================================
export type EntityListResponse = {
    data: Entity[];
    total: number;
    per_page: number;
    current_page: number;
    last_page: number;
};

export type EntityListRequest = {
    page?: number;
    per_page?: number;
    sort?: string;
    order?: 'asc' | 'desc';
    search?: string;
    [key: string]: any;
};

// =============================================
// Stats (from page controller StatsBuilder)
// =============================================
export interface EntityStats {
    totalCount: number;
    // activeCount: number;
    // inactiveCount: number;
}
```

## Step 2: i18n Translation Namespace

Create `resources/js/locales/en/<entities>.json`:

```json
{
  "page": {
    "title": "Entities",
    "headTitle": "Entities - Management",
    "myEntities": "My Entities"
  },
  "stats": {
    "total": "Total Entities",
    "active": "Active",
    "inactive": "Inactive"
  },
  "status": {
    "active": "Active",
    "inactive": "Inactive",
    "unknown": "Unknown"
  },
  "columns": {
    "name": "Name",
    "status": "Status",
    "description": "Description",
    "created": "Created"
  },
  "filters": {
    "status": "Status",
    "allStatus": "All Status",
    "allEntities": "All Entities"
  },
  "actions": {
    "importEntities": "Import Entities",
    "exportEntities": "Export Entities"
  },
  "form": {
    "entityInfo": "Entity Information",
    "metadata": "Metadata",
    "name": "Name *",
    "description": "Description",
    "status": "Status",
    "enterName": "Enter name",
    "enterDescription": "Enter description",
    "selectStatus": "Select status",
    "entityId": "Entity ID",
    "created": "Created",
    "lastUpdated": "Last Updated",
    "notSpecified": "Not specified"
  },
  "validation": {
    "nameRequired": "Name is required"
  },
  "toast": {
    "created": "Entity created successfully",
    "updated": "Entity updated successfully"
  },
  "confirm": {
    "bulkDelete": "Are you sure you want to delete {{count}} item(s)? This action cannot be undone."
  }
}
```

### Translation Key Conventions

| Section | Purpose | Example Keys |
|---------|---------|--------------|
| `page.*` | Page titles | `page.title`, `page.headTitle` |
| `stats.*` | Stats card titles | `stats.total`, `stats.active` |
| `status.*` | Status/enum display labels | `status.active`, `status.unknown` |
| `columns.*` | Table column headers | `columns.name`, `columns.status` |
| `filters.*` | Filter labels/placeholders | `filters.status`, `filters.allStatus` |
| `actions.*` | Page action buttons | `actions.importEntities` |
| `form.*` | Form labels, placeholders, metadata | `form.name`, `form.enterName` |
| `validation.*` | Client-side validation errors | `validation.nameRequired` |
| `toast.*` | Success/error messages | `toast.created`, `toast.updated` |
| `confirm.*` | Confirmation dialogs | `confirm.bulkDelete` |

Use `{{variable}}` for interpolation: `"{{count}} items"`, `"{{percent}}% complete"`.

Use `_one`/`_other` suffixes for pluralization: `"totalItems_one": "{{count}} item"`, `"totalItems_other": "{{count}} items"`.

## Step 3: Register i18n Namespace

Edit `resources/js/locales/index.ts`:

```typescript
// Add import
import entities from './en/entities.json';

// Add to ns array
ns: ['common', 'crud', 'nav', ..., 'entities'],

// Add to resources.en
resources: {
  en: {
    ...,
    entities,
  },
},
```

## Go-to-TypeScript Field Mapping

| Go Type | TypeScript Type |
|---------|----------------|
| `string` | `string` |
| `*string` | `string?` (optional) |
| `int`, `uint` | `number` |
| `float64` | `number` |
| `bool` | `boolean` |
| `*bool` | `boolean?` (optional) |
| `[]string` | `string[]` |
| `datatypes.JSON` | `string[]` or `Record<string, any>` |
| `carbon.DateTime` | `string` (ISO date) |
| `*carbon.DateTime` | `string?` (optional) |
| `gorm.DeletedAt` | `string?` (optional) |

## Dual-Case Convention

Go sends snake_case JSON, but TypeScript prefers camelCase. Include both:

```typescript
export interface Entity {
    configType?: string;    // camelCase for TS code
    config_type?: string;   // snake_case from Go JSON
}
```

## Reference

- Types: `resources/js/types/book.ts`
- Translations: `resources/js/locales/en/books.json`
- i18n config: `resources/js/locales/index.ts`
