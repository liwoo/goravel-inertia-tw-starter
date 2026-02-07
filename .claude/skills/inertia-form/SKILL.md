---
name: inertia-form
description: Create create and edit form components for a Goravel entity. Uses forwardRef pattern with useTranslation hook for i18n labels, placeholders, validation messages, and toast notifications.
argument-hint: "[EntityName]"
allowed-tools: Read, Write, Edit, Grep, Glob
---

# Inertia Form Components (i18n-aware)

Create form components for `$ARGUMENTS`.

## File 1: Create Form

`resources/js/pages/<EntityName>/sections/<EntityName>CreateForm.tsx`

### Complete Template

```tsx
import React, { useState, forwardRef, useImperativeHandle } from 'react';
import { useTranslation } from 'react-i18next';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Separator } from '@/components/ui/separator';
import {
    Select, SelectContent, SelectItem,
    SelectTrigger, SelectValue,
} from '@/components/ui/select';
import { CrudFormProps } from '@/types/crud';
import { EntityCreateData } from '@/types/entity';
// Import enum options if needed:
// import { StatusType, STATUS_OPTIONS } from '@/types/status_type';
import { User, FileText, Calendar, Tag } from 'lucide-react';

interface EntityCreateFormProps extends CrudFormProps {
    setIsSaving?: (saving: boolean) => void;
}

export const EntityCreateForm = forwardRef<any, EntityCreateFormProps>(({
    onSuccess,
    onError,
    onCancel,
    isLoading = false,
    setIsSaving,
}, ref) => {
    const { t } = useTranslation('entities');  // <-- entity namespace

    const [formData, setFormData] = useState<EntityCreateData>({
        name: '',
        description: '',
        // status: 'ACTIVE',
    });

    const [errors, setErrors] = useState<Record<string, string>>({});

    // ========================================
    // Validation (i18n error messages)
    // ========================================
    const validate = (): boolean => {
        const newErrors: Record<string, string> = {};

        if (!formData.name?.trim()) {
            newErrors.name = t('validation.nameRequired');
        }
        // if (!formData.email?.includes('@')) newErrors.email = t('validation.emailInvalid');
        // if (formData.price < 0) newErrors.price = t('validation.pricePositive');

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    // ========================================
    // Submit Handler
    // ========================================
    const handleSubmit = async () => {
        if (!validate()) return;

        setErrors({});
        setIsSaving?.(true);

        try {
            const requestData = {
                name: formData.name,
                description: formData.description,
                // config_type: formData.configType,  // camelCase → snake_case
            };

            const response = await fetch('/api/entity-names', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                    'X-Requested-With': 'XMLHttpRequest',
                },
                body: JSON.stringify(requestData),
            });

            if (response.ok) {
                onSuccess(t('toast.created'));  // <-- i18n toast
            } else {
                const errorData = await response.json().catch(() => ({}));
                if (errorData.errors) {
                    setErrors(errorData.errors);
                }
                onError?.(errorData);
            }
        } catch (error) {
            onError?.(error);
        } finally {
            setIsSaving?.(false);
        }
    };

    useImperativeHandle(ref, () => ({ handleSubmit }));

    // ========================================
    // Render
    // ========================================
    return (
        <form onSubmit={(e) => e.preventDefault()} className="space-y-6">
            <div className="space-y-6">
                <div>
                    <h3 className="text-lg font-semibold mb-4 text-foreground">
                        {t('form.entityInfo')}
                    </h3>
                    <div className="space-y-4">

                        {/* Text Input Field */}
                        <div className="flex items-start gap-3">
                            <div className="p-2 rounded-lg bg-muted">
                                <User className="h-4 w-4 text-muted-foreground" />
                            </div>
                            <div className="flex-1 space-y-2">
                                <Label htmlFor="name">{t('form.name')}</Label>
                                <Input
                                    id="name"
                                    value={formData.name}
                                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                    placeholder={t('form.enterName')}
                                    className={errors.name ? 'border-destructive' : ''}
                                />
                                {errors.name && (
                                    <p className="text-sm text-destructive">{errors.name}</p>
                                )}
                            </div>
                        </div>

                        {/* Textarea Field */}
                        <div className="flex items-start gap-3">
                            <div className="p-2 rounded-lg bg-muted">
                                <FileText className="h-4 w-4 text-muted-foreground" />
                            </div>
                            <div className="flex-1 space-y-2">
                                <Label htmlFor="description">{t('form.description')}</Label>
                                <Textarea
                                    id="description"
                                    value={formData.description}
                                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                                    placeholder={t('form.enterDescription')}
                                    rows={4}
                                    className="resize-none"
                                />
                            </div>
                        </div>

                        {/* SELECT/DROPDOWN TEMPLATE (for enums) */}
                        {/*
                        <div className="flex items-start gap-3">
                            <div className="p-2 rounded-lg bg-muted">
                                <Tag className="h-4 w-4 text-muted-foreground" />
                            </div>
                            <div className="flex-1 space-y-2">
                                <Label htmlFor="status">{t('form.status')}</Label>
                                <Select
                                    value={formData.status}
                                    onValueChange={(value) =>
                                        setFormData({ ...formData, status: value as StatusType })
                                    }
                                >
                                    <SelectTrigger className={errors.status ? 'border-destructive' : ''}>
                                        <SelectValue placeholder={t('form.selectStatus')} />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="ACTIVE">{t('status.active')}</SelectItem>
                                        <SelectItem value="INACTIVE">{t('status.inactive')}</SelectItem>
                                    </SelectContent>
                                </Select>
                                {errors.status && (
                                    <p className="text-sm text-destructive">{errors.status}</p>
                                )}
                            </div>
                        </div>
                        */}

                    </div>
                </div>
            </div>
        </form>
    );
});

EntityCreateForm.displayName = 'EntityCreateForm';
```

## File 2: Edit Form

`resources/js/pages/<EntityName>/sections/<EntityName>EditForm.tsx`

Same pattern as create form with these differences:

```tsx
import { useTranslation } from 'react-i18next';
import { CrudEditFormProps } from '@/types/crud';
import { Entity, EntityUpdateData } from '@/types/entity';

interface EntityEditFormProps extends CrudEditFormProps<Entity> {
    setIsSaving?: (saving: boolean) => void;
}

export const EntityEditForm = forwardRef<any, EntityEditFormProps>(({
    item,
    onSuccess,
    onError,
    onCancel,
    isLoading = false,
    setIsSaving,
}, ref) => {
    const { t } = useTranslation('entities');  // <-- same namespace

    // Initialize from existing item (handle dual-case)
    const [formData, setFormData] = useState<EntityUpdateData>({
        name: item.name,
        description: item.description || '',
        // status: item.status || item.status,
    });

    const handleSubmit = async () => {
        // ... validation with t('validation.*') ...

        const requestData = { /* snake_case fields */ };

        const response = await fetch(`/api/entity-names/${item.id}`, {
            method: 'PUT',  // PUT for updates
            headers: { /* same headers */ },
            body: JSON.stringify(requestData),
        });

        if (response.ok) {
            onSuccess(t('toast.updated'));  // <-- i18n toast
        }
    };

    return (
        <form onSubmit={(e) => e.preventDefault()} className="space-y-6">
            {/* ... same field layout as create, using t() for all labels ... */}

            <Separator />

            {/* Metadata Section */}
            <div>
                <h3 className="text-lg font-semibold mb-4 text-foreground">
                    {t('form.metadata')}
                </h3>
                <div className="grid grid-cols-2 gap-4 text-sm">
                    <div>
                        <p className="text-muted-foreground">{t('form.entityId')}</p>
                        <p className="font-medium text-foreground">#{item.id}</p>
                    </div>
                    <div>
                        <p className="text-muted-foreground">{t('form.created')}</p>
                        <p className="font-medium text-foreground">
                            {(item.createdAt || item.created_at)
                                ? new Date(item.createdAt || item.created_at || '').toLocaleDateString()
                                : '-'}
                        </p>
                    </div>
                    <div>
                        <p className="text-muted-foreground">{t('form.lastUpdated')}</p>
                        <p className="font-medium text-foreground">
                            {(item.updatedAt || item.updated_at)
                                ? new Date(item.updatedAt || item.updated_at || '').toLocaleDateString()
                                : '-'}
                        </p>
                    </div>
                </div>
            </div>
        </form>
    );
});

EntityEditForm.displayName = 'EntityEditForm';
```

## i18n Patterns in Forms

### Hook (inside component, NOT TFunction parameter)
Forms are React components, so they use `useTranslation` directly:
```tsx
const { t } = useTranslation('entities');
```

### Labels & Placeholders
```tsx
<Label htmlFor="name">{t('form.name')}</Label>
<Input placeholder={t('form.enterName')} />
```

### Validation Errors
```tsx
newErrors.name = t('validation.nameRequired');
newErrors.price = t('validation.pricePositive');
```

### Toast Messages
```tsx
onSuccess(t('toast.created'));  // Create form
onSuccess(t('toast.updated'));  // Edit form
```

### Select Dropdown Options (enum values)
```tsx
<SelectItem value="ACTIVE">{t('status.active')}</SelectItem>
<SelectItem value="INACTIVE">{t('status.inactive')}</SelectItem>
```

### Section Headers
```tsx
<h3>{t('form.entityInfo')}</h3>
<h3>{t('form.metadata')}</h3>
```

## Required Translation Keys

```json
{
  "form": {
    "entityInfo": "Entity Information",
    "metadata": "Metadata",
    "name": "Name *",
    "description": "Description",
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
  "status": {
    "active": "Active",
    "inactive": "Inactive"
  }
}
```

## Critical Requirements

1. **forwardRef + useImperativeHandle**: MUST expose `handleSubmit` for CrudPage drawer
2. **displayName**: MUST set `ComponentName.displayName = 'ComponentName'`
3. **useTranslation**: MUST use entity namespace for all user-visible strings
4. **camelCase -> snake_case**: Convert field names in requestData before API call
5. **X-Requested-With header**: Required for CSRF/auth
6. **Icon layout**: Every field wrapped in `flex items-start gap-3` with icon box

## Reference

See `resources/js/pages/Books/sections/BookCreateForm.tsx` and `BookEditForm.tsx`.
