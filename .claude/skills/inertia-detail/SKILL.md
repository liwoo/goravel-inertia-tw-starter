---
name: inertia-detail
description: Create the read-only detail view component for a Goravel entity. Uses useTranslation hook for i18n labels, formatted dates, status badges, and metadata section.
argument-hint: "[EntityName]"
allowed-tools: Read, Write, Edit, Grep, Glob
---

# Inertia Detail View Component (i18n-aware)

Create detail view for `$ARGUMENTS`.

## File Location

`resources/js/pages/<EntityName>/sections/<EntityName>DetailView.tsx`

## Complete Template

```tsx
import React from 'react';
import { useTranslation } from 'react-i18next';
import { Calendar, User, FileText, Tag, Hash, CheckCircle, XCircle, Clock } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { CrudDetailViewProps } from '@/types/crud';
import { Entity } from '@/types/entity';

export function EntityDetailView({
    item: entity,
    onEdit,
    onClose,
    canEdit,
}: CrudDetailViewProps<Entity>) {
    const { t } = useTranslation('entities');

    const formatDate = (date: string | Date | null | undefined) => {
        if (!date) return t('form.notSpecified');
        return new Date(date).toLocaleDateString('en-US', {
            month: 'long',
            day: 'numeric',
            year: 'numeric',
        });
    };

    const getStatusBadge = (status: string) => {
        const statusConfig: Record<string, { color: string; icon: React.ReactNode; label: string }> = {
            'ACTIVE': {
                color: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
                icon: <CheckCircle className="h-3 w-3" />,
                label: t('status.active'),
            },
            'INACTIVE': {
                color: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
                icon: <XCircle className="h-3 w-3" />,
                label: t('status.inactive'),
            },
        };
        const config = statusConfig[status] || {
            color: 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400',
            icon: null,
            label: t('status.unknown'),
        };
        return (
            <Badge className={`${config.color} flex items-center gap-1`}>
                {config.icon}
                {config.label}
            </Badge>
        );
    };

    return (
        <div className="space-y-6">
            {/* Main Information Section */}
            <div>
                <h3 className="text-lg font-semibold mb-4 text-foreground">
                    {t('form.entityInfo')}
                </h3>
                <div className="space-y-4">

                    {/* Text Field */}
                    <div className="flex items-start gap-3">
                        <div className="p-2 rounded-lg bg-muted">
                            <User className="h-4 w-4 text-muted-foreground" />
                        </div>
                        <div className="flex-1 space-y-1">
                            <p className="text-sm text-muted-foreground">
                                {t('form.name').replace(' *', '')}
                            </p>
                            <p className="font-medium text-foreground">{entity.name}</p>
                        </div>
                    </div>

                    {/* Enum/Status Field with Badge */}
                    {/*
                    <div className="flex items-start gap-3">
                        <div className="p-2 rounded-lg bg-muted">
                            <Tag className="h-4 w-4 text-muted-foreground" />
                        </div>
                        <div className="flex-1 space-y-1">
                            <p className="text-sm text-muted-foreground">{t('form.status')}</p>
                            <div className="mt-1">{getStatusBadge(entity.status)}</div>
                        </div>
                    </div>
                    */}

                    {/* Long Text Field */}
                    <div className="flex items-start gap-3">
                        <div className="p-2 rounded-lg bg-muted">
                            <FileText className="h-4 w-4 text-muted-foreground" />
                        </div>
                        <div className="flex-1 space-y-1">
                            <p className="text-sm text-muted-foreground">{t('form.description')}</p>
                            <p className="font-medium text-foreground">
                                {entity.description || (
                                    <span className="italic text-muted-foreground">
                                        {t('form.notSpecified')}
                                    </span>
                                )}
                            </p>
                        </div>
                    </div>

                    {/* Array/Tags Field */}
                    {/*
                    <div className="flex items-start gap-3">
                        <div className="p-2 rounded-lg bg-muted">
                            <Tag className="h-4 w-4 text-muted-foreground" />
                        </div>
                        <div className="flex-1 space-y-1">
                            <p className="text-sm text-muted-foreground">{t('form.tags')}</p>
                            <div className="flex flex-wrap gap-2 mt-1">
                                {entity.tags?.map((tag, i) => (
                                    <Badge key={i} variant="secondary">{tag}</Badge>
                                )) || (
                                    <span className="text-muted-foreground italic">
                                        {t('form.notSpecified')}
                                    </span>
                                )}
                            </div>
                        </div>
                    </div>
                    */}

                </div>
            </div>

            <Separator />

            {/* Metadata Section */}
            <div>
                <h3 className="text-lg font-semibold mb-4 text-foreground">
                    {t('form.metadata')}
                </h3>
                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <p className="text-sm text-muted-foreground">{t('form.entityId')}</p>
                        <p className="font-medium text-sm text-foreground">#{entity.id}</p>
                    </div>
                    <div>
                        <p className="text-sm text-muted-foreground">{t('form.created')}</p>
                        <p className="font-medium text-sm text-foreground">
                            {formatDate(entity.createdAt || entity.created_at)}
                        </p>
                    </div>
                    <div>
                        <p className="text-sm text-muted-foreground">{t('form.lastUpdated')}</p>
                        <p className="font-medium text-sm text-foreground">
                            {formatDate(entity.updatedAt || entity.updated_at)}
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
}
```

## i18n Pattern: `useTranslation` in Detail Views

Detail views are React components (not plain functions), so they use the hook directly:

```tsx
const { t } = useTranslation('entities');
```

### Label Reuse from Forms
Detail views reuse `form.*` keys (minus the `*` suffix for required markers):
```tsx
<p className="text-sm text-muted-foreground">
    {t('form.name').replace(' *', '')}
</p>
```

### Status Badges
Build status config inline using `t('status.*')` keys:
```tsx
const statusConfig = {
    'ACTIVE': { label: t('status.active'), ... },
    'INACTIVE': { label: t('status.inactive'), ... },
};
```

### Empty Values
```tsx
{entity.field || (
    <span className="italic text-muted-foreground">{t('form.notSpecified')}</span>
)}
```

## Required Translation Keys

```json
{
  "form": {
    "entityInfo": "Entity Information",
    "metadata": "Metadata",
    "name": "Name *",
    "description": "Description",
    "status": "Status",
    "entityId": "Entity ID",
    "created": "Created",
    "lastUpdated": "Last Updated",
    "notSpecified": "Not specified"
  },
  "status": {
    "active": "Active",
    "inactive": "Inactive",
    "unknown": "Unknown"
  }
}
```

## Key Differences from Forms

- **Functional component** (NOT forwardRef)
- **Read-only** — no inputs, just display
- Uses `CrudDetailViewProps<Entity>` interface
- Same icon layout pattern as forms for visual consistency
- Always include Metadata section with ID, Created, Updated

## Reference

See `resources/js/pages/Books/sections/BookDetailView.tsx` for a complete i18n-aware example.
