---
name: inertia-page-config
description: Create page configuration for a Goravel entity - i18n-aware simple filters with badges, stats cards, page actions, and bulk actions. All config functions receive TFunction.
argument-hint: "[EntityName]"
allowed-tools: Read, Write, Edit, Grep, Glob
---

# Inertia Page Configuration (i18n-aware)

Create page config for `$ARGUMENTS`.

## File Location

`resources/js/pages/<EntityName>/sections/<EntityName>PageConfig.tsx`

## Complete Template

```tsx
import React from 'react';
import { TFunction } from 'i18next';
import { router } from '@inertiajs/react';
import { Upload, Download, BarChart3 } from 'lucide-react';
import {
    StatsCardConfig,
    PageActionConfig,
    SimpleFilterConfig,
} from '@/lib/crud-page-utils';

// =============================================
// Stats Card Configurations
// =============================================
export function getEntityStatsConfigs(t: TFunction): StatsCardConfig[] {
    return [
        {
            title: t('stats.total'),
            getValue: (stats) => stats?.totalCount || 0,
        },
        // {
        //     title: t('stats.active'),
        //     getValue: (stats) => stats?.activeCount || 0,
        //     valueClassName: 'text-green-600',
        // },
        // {
        //     title: t('stats.inactive'),
        //     getValue: (stats) => stats?.inactiveCount || 0,
        //     valueClassName: 'text-red-600',
        // },
    ];
}

// =============================================
// Simple Filters (Tab/Dropdown with Badges)
// =============================================
export function getEntitySimpleFilters(t: TFunction, stats: any): SimpleFilterConfig[] {
    return [
        {
            key: 'status-active',
            label: t('status.active'),
            value: 'ACTIVE',
            badge: stats?.activeCount || 0,
            filterParams: { status: 'ACTIVE' },
        },
        {
            key: 'status-inactive',
            label: t('status.inactive'),
            value: 'INACTIVE',
            badge: stats?.inactiveCount || 0,
            filterParams: { status: 'INACTIVE' },
        },
    ];
}

// =============================================
// Page Actions (Export, Import, etc.)
// =============================================
export function getEntityPageActions(
    t: TFunction,
    permissions: any,
    handlers: {
        onImport?: () => void;
        onExport?: () => void;
    }
): PageActionConfig[] {
    return [
        // {
        //     key: 'import',
        //     label: t('actions.importEntities'),
        //     icon: <Upload className="h-4 w-4" />,
        //     handler: handlers.onImport || (() => {}),
        //     permission: permissions.canCreate,
        // },
        // {
        //     key: 'export',
        //     label: t('actions.exportEntities'),
        //     icon: <Download className="h-4 w-4" />,
        //     handler: handlers.onExport || (() => {}),
        // },
    ];
}

// =============================================
// Bulk Actions
// =============================================
export function getEntityBulkActions(t: TFunction) {
    return {
        handleBulkDelete: (ids: number[]) => {
            const confirmMessage = t('confirm.bulkDelete', { count: ids.length });
            if (confirm(confirmMessage)) {
                router.delete('/api/entity-names/bulk', {
                    data: { ids },
                });
            }
        },
    };
}
```

## Key i18n Pattern: `TFunction` Parameter

**All config functions receive `t: TFunction` as first parameter** because they are plain functions, not React components.

```typescript
import { TFunction } from 'i18next';

// CORRECT: TFunction as parameter
export function getEntitySimpleFilters(t: TFunction, stats: any): SimpleFilterConfig[] { ... }
export function getEntityStatsConfigs(t: TFunction): StatsCardConfig[] { ... }
export function getEntityPageActions(t: TFunction, permissions: any, handlers: {...}): PageActionConfig[] { ... }
export function getEntityBulkActions(t: TFunction) { ... }
```

Called from Index.tsx:
```tsx
const { t } = useTranslation('entities');
const simpleFilters = createSimpleFilters(getEntitySimpleFilters(t, stats));
const pageActions = createPageActions(getEntityPageActions(t, permissions, handlers));
```

### Interpolation in Confirmations
```tsx
t('confirm.bulkDelete', { count: ids.length })
// JSON: "confirm.bulkDelete": "Are you sure you want to delete {{count}} item(s)?"
```

## Using Simple Filters in Index.tsx

```tsx
import { createSimpleFilters } from '@/lib/crud-page-utils';
import { getEntitySimpleFilters } from './sections';

const simpleFilters = createSimpleFilters(getEntitySimpleFilters(t, stats));

<CrudPage
    simpleFilters={simpleFilters}
    simpleFiltersVariant="dropdown"         // or "tabs"
    simpleFiltersDropdownLabel={t('filters.status')}
/>
```

## Required Translation Keys

```json
{
  "stats": {
    "total": "Total Entities",
    "active": "Active",
    "inactive": "Inactive"
  },
  "status": {
    "active": "Active",
    "inactive": "Inactive"
  },
  "actions": {
    "importEntities": "Import Entities",
    "exportEntities": "Export Entities"
  },
  "filters": {
    "status": "Status"
  },
  "confirm": {
    "bulkDelete": "Are you sure you want to delete {{count}} item(s)? This action cannot be undone."
  }
}
```

## Barrel Export

Add to `resources/js/pages/<EntityName>/sections/index.ts`:

```typescript
export { EntityDetailView } from './EntityDetailView';
export { EntityCreateForm } from './EntityCreateForm';
export { EntityEditForm } from './EntityEditForm';
export { getEntityColumns, getEntityColumnsMobile, getEntityFilters } from './EntityColumns';
export { getEntityStatsConfigs, getEntitySimpleFilters, getEntityPageActions, getEntityBulkActions } from './EntityPageConfig';
```

## Reference

See `resources/js/pages/Books/sections/bookPageConfig.tsx` for a complete i18n-aware example.
