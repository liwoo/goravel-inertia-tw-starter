---
name: inertia-columns
description: Create table column definitions and filter configs for a Goravel entity's CrudPage. Uses TFunction parameter for i18n-translated column headers, status labels, and filter labels.
argument-hint: "[EntityName]"
allowed-tools: Read, Write, Edit, Grep, Glob
---

# Inertia Column & Filter Definitions

Create columns and filters for `$ARGUMENTS`.

## File Location

`resources/js/pages/<EntityName>/sections/<EntityName>Columns.tsx`

## Complete Template

```tsx
import React from 'react';
import { TFunction } from 'i18next';
import { Entity } from '@/types/entity';
import { CrudColumn, CrudFilter } from '@/types/crud';
import { Badge } from '@/components/ui/badge';
import { User, Calendar, Tag, CheckCircle, Clock, XCircle } from 'lucide-react';

// =============================================
// Status Configuration (i18n-aware)
// =============================================
function getEntityStatusConfig(t: TFunction) {
    return {
        ACTIVE: {
            label: t('status.active'),
            icon: <CheckCircle className="h-3 w-3" />,
            color: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
        },
        INACTIVE: {
            label: t('status.inactive'),
            icon: <XCircle className="h-3 w-3" />,
            color: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
        },
    };
}

// =============================================
// Desktop Columns
// =============================================
export function getEntityColumns(t: TFunction): CrudColumn<Entity>[] {
    const STATUS_CONFIG = getEntityStatusConfig(t);

    return [
        {
            key: 'name',
            label: t('columns.name'),
            sortable: true,
            className: 'min-w-[250px]',
            render: (entity) => (
                <div className="flex items-start gap-3">
                    <div className="p-2 rounded-lg bg-muted">
                        <User className="h-5 w-5 text-muted-foreground" />
                    </div>
                    <div className="space-y-1">
                        <div className="font-medium text-foreground">{entity.name}</div>
                        {entity.description && (
                            <div className="text-sm text-muted-foreground line-clamp-1">
                                {entity.description}
                            </div>
                        )}
                    </div>
                </div>
            ),
        },
        {
            key: 'status',
            label: t('columns.status'),
            sortable: true,
            className: 'w-32',
            render: (entity) => {
                const config = STATUS_CONFIG[entity.status as keyof typeof STATUS_CONFIG];
                if (!config) {
                    return <Badge variant="outline">{t('status.unknown')}</Badge>;
                }
                return (
                    <Badge className={`${config.color} flex items-center gap-1`}>
                        {config.icon}
                        {config.label}
                    </Badge>
                );
            },
        },
        {
            key: 'createdAt',
            label: t('columns.created'),
            sortable: true,
            className: 'w-28',
            render: (entity) => {
                const dateValue = entity.createdAt || entity.created_at;
                if (!dateValue) return <div className="text-sm text-muted-foreground">-</div>;

                return (
                    <div className="text-sm text-muted-foreground flex items-center">
                        <Calendar className="w-3 h-3 mr-1" />
                        {new Date(dateValue).toLocaleDateString('en-US', {
                            month: 'short',
                            day: 'numeric',
                            year: 'numeric',
                        })}
                    </div>
                );
            },
        },
    ];
}

// =============================================
// Mobile Columns (Compact)
// =============================================
export function getEntityColumnsMobile(t: TFunction): CrudColumn<Entity>[] {
    const STATUS_CONFIG = getEntityStatusConfig(t);

    return [
        {
            key: 'name',
            label: t('columns.name'),
            sortable: true,
            render: (entity) => (
                <div className="space-y-2">
                    <div className="flex items-start gap-3">
                        <div className="p-2 rounded-lg bg-muted">
                            <User className="h-5 w-5 text-muted-foreground" />
                        </div>
                        <div className="flex-1 space-y-1">
                            <div className="font-medium text-foreground">{entity.name}</div>
                            {entity.description && (
                                <div className="text-sm text-muted-foreground line-clamp-1">
                                    {entity.description}
                                </div>
                            )}
                        </div>
                    </div>
                    <div className="flex items-center justify-between pl-12">
                        {(() => {
                            const config = STATUS_CONFIG[entity.status as keyof typeof STATUS_CONFIG];
                            if (!config) return <Badge variant="outline" className="text-xs">{t('status.unknown')}</Badge>;
                            return (
                                <Badge className={`${config.color} flex items-center gap-1 text-xs`}>
                                    {config.icon}
                                    {config.label}
                                </Badge>
                            );
                        })()}
                    </div>
                </div>
            ),
        },
    ];
}

// =============================================
// Filter Definitions
// =============================================
export function getEntityFilters(t: TFunction): CrudFilter[] {
    return [
        {
            key: 'status',
            label: t('filters.status'),
            type: 'select',
            options: [
                { value: '__all__', label: t('filters.allStatus') },
                { value: 'ACTIVE', label: t('status.active') },
                { value: 'INACTIVE', label: t('status.inactive') },
            ],
        },
        // Text filter example:
        // {
        //     key: 'author',
        //     label: t('filters.author'),
        //     type: 'text',
        //     placeholder: t('filters.filterByAuthor'),
        // },
        // Number filter example:
        // {
        //     key: 'minPrice',
        //     label: t('filters.minPrice'),
        //     type: 'number',
        //     placeholder: '0.00',
        // },
        // Date filter example:
        // {
        //     key: 'publishedAfter',
        //     label: t('filters.publishedAfter'),
        //     type: 'date',
        // },
    ];
}
```

## Key i18n Pattern: `TFunction` Parameter

**All column/filter functions receive `t: TFunction` as first parameter** instead of calling `useTranslation` internally. This is because these are plain functions, not React components.

```typescript
import { TFunction } from 'i18next';

// CORRECT: TFunction as parameter
export function getEntityColumns(t: TFunction): CrudColumn<Entity>[] {
    return [{ label: t('columns.name'), ... }];
}

// Called from Index.tsx where useTranslation is available:
const { t } = useTranslation('entities');
<CrudPage columns={getEntityColumns(t)} />
```

## Status Config Pattern

Extract status config into a function that receives `t`:

```tsx
function getEntityStatusConfig(t: TFunction) {
    return {
        ACTIVE: {
            label: t('status.active'),
            icon: <CheckCircle className="h-3 w-3" />,
            color: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
        },
        // ... more statuses
    };
}
```

## Required Translation Keys

Ensure these keys exist in the entity's i18n namespace (`locales/en/<entities>.json`):

```json
{
  "status": { "active": "Active", "inactive": "Inactive", "unknown": "Unknown" },
  "columns": { "name": "Name", "status": "Status", "created": "Created" },
  "filters": { "status": "Status", "allStatus": "All Status" }
}
```

## Price/Currency Column Pattern

For decimal/currency fields, always format to avoid floating-point display issues like `23.989999771118164`:

```tsx
{
    key: 'price',
    label: t('columns.price'),
    sortable: true,
    className: 'w-28',
    render: (entity) => (
        <div className="text-sm text-foreground font-medium">
            ${entity.price?.toFixed(2) ?? '0.00'}
        </div>
    ),
},
```

For locale-aware formatting:
```tsx
render: (entity) => (
    <div className="text-sm text-foreground font-medium">
        {entity.price?.toLocaleString('en-US', {
            style: 'currency',
            currency: 'USD',
        }) ?? '$0.00'}
    </div>
),
```

**Never display raw `float64` values** — they may have precision artifacts from JSON deserialization.

## Verify

After creating column definitions:

```bash
# TypeScript compiles (catches wrong Entity property access)
npx tsc --noEmit

# Lint the columns file
npx eslint "resources/js/pages/<EntityName>/sections/<EntityName>Columns.tsx" --max-warnings=0
```

## Reference

See `resources/js/pages/Books/sections/BookColumns.tsx` for a complete i18n-aware example.
