---
name: inertia-page
description: Create the Index.tsx page component for a Goravel entity. Uses CrudPage wrapper with i18n useTranslation hook, responsive columns, and section imports.
argument-hint: "[EntityName]"
allowed-tools: Bash, Read, Write, Edit, Grep, Glob
---

# Inertia Index Page Component

Create Index page for `$ARGUMENTS`.

## File Location

`resources/js/pages/<EntityName>/Index.tsx`

**Note**: Use `make:ui` to generate the initial structure, then customize:
```bash
go run . artisan make:ui --page=$ARGUMENTS --request=$ARGUMENTS
```

## Complete Template

```tsx
import React from 'react';
// @ts-ignore
import { Head, router } from '@inertiajs/react';
import { useTranslation } from 'react-i18next';
import { Entity, EntityListResponse, EntityListRequest, EntityStats } from '@/types/entity';
import { CrudPage } from '@/components/Crud/CrudPage';
import {
    EntityCreateForm,
    EntityEditForm,
    EntityDetailView,
    getEntityColumns,
    getEntityColumnsMobile,
    getEntityFilters,
    getEntityStatsConfigs,
    getEntitySimpleFilters,
    getEntityPageActions,
    getEntityBulkActions,
} from './sections';
import {
    renderStatsCards,
    createPageActions,
    createSimpleFilters,
} from '@/lib/crud-page-utils';
import { useIsMobile } from '@/hooks/use-mobile';
import Admin from '@/layouts/Admin';

interface EntityIndexProps {
    data: EntityListResponse;
    filters: EntityListRequest;
    stats?: EntityStats;
    permissions: {
        canCreate: boolean;
        canEdit: boolean;
        canDelete: boolean;
    };
    meta?: {
        pagination: {
            defaultPageSize: number;
            maxPageSize: number;
            allowedSizes: number[];
        };
    };
}

export default function EntityIndex({
    data,
    filters,
    stats,
    permissions,
    meta,
}: EntityIndexProps) {
    const isMobile = useIsMobile();
    const { t } = useTranslation('entities');  // <-- entity namespace

    const handleRefresh = () => {
        router.reload({ only: ['data', 'stats'] });
    };

    // Use i18n-aware configurations — pass t to all config functions
    const simpleFilters = createSimpleFilters(getEntitySimpleFilters(t, stats));

    const pageActions = createPageActions(
        getEntityPageActions(t, permissions, {
            onExport: () => {/* TODO */},
        })
    );

    return (
        <Admin title={t('page.title')}>
            <Head title={t('page.headTitle')} />

            <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
                {/* Statistics Cards */}
                {renderStatsCards(stats, getEntityStatsConfigs(t))}

                {/* Main CRUD Component */}
                <div className="px-0">
                    <CrudPage<Entity>
                        data={data}
                        filters={filters}
                        title={t('page.myEntities')}
                        resourceName="entities"
                        columns={isMobile ? getEntityColumnsMobile(t) : getEntityColumns(t)}
                        customFilters={getEntityFilters(t)}
                        simpleFilters={simpleFilters}
                        pageActions={pageActions}
                        paginationConfig={meta?.pagination}
                        createForm={EntityCreateForm}
                        editForm={EntityEditForm}
                        detailView={EntityDetailView}
                        onRefresh={handleRefresh}
                        canCreate={permissions.canCreate}
                        canEdit={permissions.canEdit}
                        canDelete={permissions.canDelete}
                        canView={true}
                    />
                </div>
            </div>
        </Admin>
    );
}
```

## Key i18n Patterns

### Namespace Selection
```tsx
const { t } = useTranslation('entities');  // Use the entity's own namespace
```

### Page Titles
```tsx
<Admin title={t('page.title')}>
<Head title={t('page.headTitle')} />
```

### Passing `t` to Config Functions
All section config functions receive `t: TFunction` as first parameter:
```tsx
getEntityColumns(t)                              // Columns get translated headers
getEntityFilters(t)                              // Filters get translated labels
getEntitySimpleFilters(t, stats)                 // Simple filters get translated labels
getEntityPageActions(t, permissions, handlers)    // Actions get translated labels
getEntityStatsConfigs(t)                         // Stats get translated titles
```

### CrudPage Title
```tsx
title={t('page.myEntities')}  // Translated resource title
```

## Props from Backend

The `GenericPageController.Index()` passes:
- `data` - Paginated entity list
- `filters` - Current page/sort/search state
- `permissions` - `canCreate`, `canEdit`, `canDelete`
- `stats` - Optional statistics (if `StatsEnabled=true`)
- `meta` - Pagination config

## Reference

See `resources/js/pages/Books/Index.tsx` for a complete i18n-aware example.
