package commands

import (
	"fmt"
	"strings"
)

// generateTypeContent generates the TypeScript type file
func (receiver *UIMaker) generateTypeContent(pageName, pluralPage string, fields []Field) string {
	resourceName := strings.ToLower(pageName)

	// Generate field interfaces
	var fieldDefs strings.Builder
	for _, field := range fields {
		optional := ""
		if !field.Required {
			optional = "?"
		}
		fieldDefs.WriteString(fmt.Sprintf("  %s%s: %s;\n", strings.ToLower(field.Name[:1])+field.Name[1:], optional, field.TSType))
	}

	return fmt.Sprintf(`// TypeScript interfaces for %s entities and operations
import { BaseModel, PaginatedResult, ListRequest } from './crud';

// Core %s interface matching the backend model
export interface %s extends BaseModel {
%s}

// %s creation data (matches %sCreateRequest)
export interface %sCreateData {
%s}

// %s update data (matches %sUpdateRequest - all optional)
export interface %sUpdateData {
%s}

// %s list response (matches service GetList response)
export interface %sListResponse extends PaginatedResult<%s> {}

// %s list request (extends base ListRequest with %s-specific filters)
export interface %sListRequest extends ListRequest {
  // Add your custom filters here
}

// Form validation types
export interface %sFormErrors {
%s  general?: string;
}

// %s statistics (if provided by backend)
export interface %sStats {
  total%s: number;
  // Add your custom stats here
}
`, pageName, pageName, pageName,
		fieldDefs.String(),
		pageName, pageName, pageName,
		fieldDefs.String(),
		pageName, pageName, pageName,
		receiver.makeFieldsOptional(fieldDefs.String()),
		pageName, pageName, pageName,
		pageName, resourceName, pageName,
		pageName,
		receiver.generateErrorFields(fields),
		pageName, pageName, pluralPage)
}

// generateIndexContent generates the Index.tsx file
func (receiver *UIMaker) generateIndexContent(pageName, pluralPage string, fields []Field) string {
	resourceName := strings.ToLower(pageName)
	camelResource := strings.ToLower(pageName[:1]) + pageName[1:]

	return fmt.Sprintf(`import React, { useState } from 'react';
import { Head, router } from '@inertiajs/react';
import {
  %s,
  %sListResponse,
  %sListRequest
} from '@/types/%s';
import { CrudPage } from '@/components/Crud/CrudPage';
import {
  %sCreateForm,
  %sEditForm,
  %sDetailView,
  %sColumns,
  %sColumnsMobile,
  %sFilters
} from './sections';
import { useIsMobile } from '@/hooks/use-mobile';
import Admin from '@/layouts/Admin';

// Props interface for the %s Index page
interface %sIndexProps {
  data: %sListResponse;
  filters: %sListRequest;
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

export default function %sIndex({
  data,
  filters,
  permissions,
  meta
}: %sIndexProps) {
  const isMobile = useIsMobile();

  const handleRefresh = () => {
    router.reload({ only: ['data'] });
  };

  return (
    <Admin title={"%s"}>
      <Head title="%s - Management" />

      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        {/* Main CRUD Component */}
        <div className="px-0">
          <CrudPage<%s>
            data={data}
            filters={filters}
            title="%s"
            resourceName="%s"
            columns={isMobile ? %sColumnsMobile : %sColumns}
            customFilters={%sFilters}
            paginationConfig={meta?.pagination}
            createForm={%sCreateForm}
            editForm={%sEditForm}
            detailView={%sDetailView}
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
`, pageName, pageName, pageName, resourceName,
		pageName, pageName, pageName, camelResource, camelResource, camelResource,
		pageName, pageName, pageName, pageName,
		pageName, pageName,
		pageName, pageName,
		pageName, pluralPage, pluralPage, camelResource, camelResource, camelResource,
		pageName, pageName, pageName)
}

// generateColumnsContent generates the Columns.tsx file
func (receiver *UIMaker) generateColumnsContent(pageName, pluralPage string, fields []Field) string {
	camelResource := strings.ToLower(pageName[:1]) + pageName[1:]
	resourceName := strings.ToLower(pageName)

	// Get unique icons needed
	iconsMap := make(map[string]bool)
	for _, field := range fields {
		if field.ShowInTable {
			iconsMap[field.Icon] = true
		}
	}
	var icons []string
	for icon := range iconsMap {
		icons = append(icons, icon)
	}

	// Generate column definitions
	var columnDefs strings.Builder
	for i, field := range fields {
		if !field.ShowInTable {
			continue
		}
		fieldLower := strings.ToLower(field.Name[:1]) + field.Name[1:]

		columnDefs.WriteString(fmt.Sprintf(`  {
    key: '%s',
    label: '%s',
    sortable: %v,
    className: 'min-w-[150px]',
    render: (%s) => (
      <div className="flex items-start gap-3">
        <div className="p-2 rounded-lg bg-muted">
          <%s className="h-4 w-4 text-muted-foreground" />
        </div>
        <div className="flex-1 space-y-1">
          <p className="font-medium text-foreground">{%s.%s}</p>
        </div>
      </div>
    ),
  }`, fieldLower, field.Label, field.Sortable, camelResource, field.Icon, camelResource, fieldLower))

		if i < len(fields)-1 {
			columnDefs.WriteString(",\n")
		}
	}

	// Generate filter definitions
	var filterDefs strings.Builder
	for i, field := range fields {
		if field.FilterType == "" {
			continue
		}
		fieldLower := strings.ToLower(field.Name[:1]) + field.Name[1:]

		filterDefs.WriteString(fmt.Sprintf(`  {
    key: '%s',
    label: '%s',
    type: '%s',
    placeholder: '%s',
  }`, fieldLower, field.Label, field.FilterType, field.Placeholder))

		if i < len(fields)-1 {
			filterDefs.WriteString(",\n")
		}
	}

	return fmt.Sprintf(`import React from 'react';
import { %s } from '@/types/%s';
import { CrudColumn, CrudFilter } from '@/types/crud';
import { Badge } from '@/components/ui/badge';
import { %s } from 'lucide-react';

/**
 * %s table columns configuration
 */
export const %sColumns: CrudColumn<%s>[] = [
%s,
  {
    key: 'createdAt',
    label: 'Created',
    sortable: true,
    className: 'w-32',
    render: (%s) => {
      const dateValue = %s.createdAt || %s.created_at;
      if (!dateValue) return <span className="text-sm text-muted-foreground">-</span>;

      const date = new Date(dateValue);
      return (
        <div className="text-sm text-muted-foreground">
          {date.toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric'
          })}
        </div>
      );
    },
  },
];

/**
 * Compact columns for mobile/smaller screens
 */
export const %sColumnsMobile: CrudColumn<%s>[] = [
  {
    key: 'combined',
    label: '%s',
    sortable: false,
    render: (%s) => (
      <div className="space-y-2">
        <div className="font-medium text-foreground">{%s.%s}</div>
        <div className="text-sm text-muted-foreground">
          {new Date(%s.createdAt || %s.created_at || '').toLocaleDateString()}
        </div>
      </div>
    ),
  },
];

/**
 * %s filters configuration
 */
export const %sFilters: CrudFilter[] = [
%s
];
`, pageName, resourceName, strings.Join(icons, ", "),
		pageName, camelResource, pageName,
		columnDefs.String(),
		camelResource, camelResource, camelResource,
		camelResource, pageName,
		pageName, camelResource,
		camelResource, receiver.getFirstTableField(fields), camelResource, camelResource,
		pageName, camelResource,
		filterDefs.String())
}

// generateCreateFormContent generates the CreateForm.tsx file
func (receiver *UIMaker) generateCreateFormContent(pageName, pluralPage string, fields []Field) string {
	resourceName := strings.ToLower(pageName)

	// Generate form fields
	var formFields strings.Builder
	for _, field := range fields {
		fieldLower := strings.ToLower(field.Name[:1]) + field.Name[1:]
		requiredLabel := ""
		if field.Required {
			requiredLabel = " *"
		}

		formFields.WriteString(fmt.Sprintf(`
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <%s className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="%s">%s%s</Label>`, field.Icon, fieldLower, field.Label, requiredLabel))

		if field.FormControl == "textarea" {
			formFields.WriteString(fmt.Sprintf(`
                <Textarea
                  id="%s"
                  value={formData.%s}
                  onChange={(e) => setFormData({ ...formData, %s: e.target.value })}
                  placeholder="%s"
                  rows={4}
                  className={errors.%s ? 'border-destructive' : ''}
                />`, fieldLower, fieldLower, fieldLower, field.Placeholder, fieldLower))
		} else {
			formFields.WriteString(fmt.Sprintf(`
                <Input
                  id="%s"
                  value={formData.%s}
                  onChange={(e) => setFormData({ ...formData, %s: e.target.value })}
                  placeholder="%s"
                  className={errors.%s ? 'border-destructive' : ''}
                />`, fieldLower, fieldLower, fieldLower, field.Placeholder, fieldLower))
		}

		formFields.WriteString(fmt.Sprintf(`
                {errors.%s && (
                  <p className="text-sm text-destructive">{errors.%s}</p>
                )}
              </div>
            </div>
`, fieldLower, fieldLower))
	}

	// Generate initial form data
	var initialData strings.Builder
	for i, field := range fields {
		fieldLower := strings.ToLower(field.Name[:1]) + field.Name[1:]
		defaultVal := "''"
		if field.TSType == "number" {
			defaultVal = "0"
		} else if field.TSType == "boolean" {
			defaultVal = "false"
		}

		initialData.WriteString(fmt.Sprintf("    %s: %s", fieldLower, defaultVal))
		if i < len(fields)-1 {
			initialData.WriteString(",\n")
		}
	}

	return fmt.Sprintf(`import React, { useState, forwardRef, useImperativeHandle } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { CrudFormProps } from '@/types/crud';
import { %sCreateData } from '@/types/%s';
import { %s } from 'lucide-react';
import { Separator } from '@/components/ui/separator';

interface %sCreateFormProps extends CrudFormProps {
  setIsSaving?: (saving: boolean) => void;
}

export const %sCreateForm = forwardRef<any, %sCreateFormProps>(({
  onSuccess,
  onError,
  onCancel,
  isLoading = false,
  setIsSaving
}, ref) => {
  const [formData, setFormData] = useState<%sCreateData>({
%s,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleSubmit = async () => {
    // Basic validation
    const newErrors: Record<string, string> = {};

    // TODO: Add validation rules
    %s

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setErrors({});
    setIsSaving?.(true);

    try {
      const response = await fetch('/api/%s', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        onSuccess('%s created successfully');
      } else {
        const errorData = await response.json().catch(() => ({}));
        onError?.(errorData);
      }
    } catch (error) {
      onError?.(error);
    } finally {
      setIsSaving?.(false);
    }
  };

  // Expose handleSubmit to parent component
  useImperativeHandle(ref, () => ({
    handleSubmit
  }));

  return (
    <form onSubmit={(e) => e.preventDefault()} className="space-y-6">
      <div className="space-y-6">
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">%s Information</h3>
          <div className="space-y-4">%s
          </div>
        </div>
      </div>
    </form>
  );
});

%sCreateForm.displayName = '%sCreateForm';
`, pageName, resourceName, receiver.getUniqueIcons(fields),
		pageName, pageName, pageName,
		pageName, initialData.String(),
		receiver.generateValidationRules(fields),
		pluralPage, pageName, pageName,
		formFields.String(),
		pageName, pageName)
}

// generateEditFormContent generates the EditForm.tsx file
func (receiver *UIMaker) generateEditFormContent(pageName, pluralPage string, fields []Field) string {
	resourceName := strings.ToLower(pageName)
	camelResource := strings.ToLower(pageName[:1]) + pageName[1:]

	// Generate form fields
	var formFields strings.Builder
	for _, field := range fields {
		fieldLower := strings.ToLower(field.Name[:1]) + field.Name[1:]
		requiredLabel := ""
		if field.Required {
			requiredLabel = " *"
		}

		formFields.WriteString(fmt.Sprintf(`
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <%s className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="%s">%s%s</Label>`, field.Icon, fieldLower, field.Label, requiredLabel))

		if field.FormControl == "textarea" {
			formFields.WriteString(fmt.Sprintf(`
                <Textarea
                  id="%s"
                  value={formData.%s}
                  onChange={(e) => setFormData({ ...formData, %s: e.target.value })}
                  placeholder="%s"
                  rows={4}
                  className={errors.%s ? 'border-destructive' : ''}
                />`, fieldLower, fieldLower, fieldLower, field.Placeholder, fieldLower))
		} else {
			formFields.WriteString(fmt.Sprintf(`
                <Input
                  id="%s"
                  value={formData.%s}
                  onChange={(e) => setFormData({ ...formData, %s: e.target.value })}
                  placeholder="%s"
                  className={errors.%s ? 'border-destructive' : ''}
                />`, fieldLower, fieldLower, fieldLower, field.Placeholder, fieldLower))
		}

		formFields.WriteString(fmt.Sprintf(`
                {errors.%s && (
                  <p className="text-sm text-destructive">{errors.%s}</p>
                )}
              </div>
            </div>
`, fieldLower, fieldLower))
	}

	// Generate initial form data from item
	var initialData strings.Builder
	for i, field := range fields {
		fieldLower := strings.ToLower(field.Name[:1]) + field.Name[1:]

		// Handle optional fields with fallback
		if !field.Required {
			if field.TSType == "number" {
				initialData.WriteString(fmt.Sprintf("    %s: %s.%s || 0", fieldLower, camelResource, fieldLower))
			} else {
				initialData.WriteString(fmt.Sprintf("    %s: %s.%s || ''", fieldLower, camelResource, fieldLower))
			}
		} else {
			initialData.WriteString(fmt.Sprintf("    %s: %s.%s", fieldLower, camelResource, fieldLower))
		}

		if i < len(fields)-1 {
			initialData.WriteString(",\n")
		}
	}

	// Generate validation with optional field checks
	var validationRules strings.Builder
	for _, field := range fields {
		if !field.Required {
			continue
		}
		fieldLower := strings.ToLower(field.Name[:1]) + field.Name[1:]
		validationRules.WriteString(fmt.Sprintf(`if (!formData.%s?.trim()) {
      newErrors.%s = '%s is required';
    }
    `, fieldLower, fieldLower, field.Label))
	}

	return fmt.Sprintf(`import React, { useState, forwardRef, useImperativeHandle } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { CrudEditFormProps } from '@/types/crud';
import { %s, %sUpdateData } from '@/types/%s';
import { %s } from 'lucide-react';
import { Separator } from '@/components/ui/separator';

interface %sEditFormProps extends CrudEditFormProps<%s> {
  setIsSaving?: (saving: boolean) => void;
}

export const %sEditForm = forwardRef<any, %sEditFormProps>(({
  item: %s,
  onSuccess,
  onError,
  onCancel,
  isLoading = false,
  setIsSaving
}, ref) => {
  const [formData, setFormData] = useState<%sUpdateData>({
%s,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleSubmit = async () => {
    // Basic validation
    const newErrors: Record<string, string> = {};

    // TODO: Add validation rules
    %s

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setErrors({});
    setIsSaving?.(true);

    try {
      const response = await fetch(`+"`"+`/api/%s/${%s.id}`+"`"+`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        onSuccess('%s updated successfully');
      } else {
        const errorData = await response.json().catch(() => ({}));
        onError?.(errorData);
      }
    } catch (error) {
      onError?.(error);
    } finally {
      setIsSaving?.(false);
    }
  };

  // Expose handleSubmit to parent component
  useImperativeHandle(ref, () => ({
    handleSubmit
  }));

  return (
    <form onSubmit={(e) => e.preventDefault()} className="space-y-6">
      <div className="space-y-6">
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">%s Information</h3>
          <div className="space-y-4">%s
          </div>
        </div>

        <Separator />

        {/* Metadata Section */}
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">Metadata</h3>
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <p className="text-muted-foreground">ID</p>
              <p className="font-medium text-foreground">#{%s.id}</p>
            </div>
            <div>
              <p className="text-muted-foreground">Created</p>
              <p className="font-medium text-foreground">
                {(%s.createdAt || %s.created_at) ? new Date(%s.createdAt || %s.created_at || '').toLocaleDateString() : '-'}
              </p>
            </div>
            <div>
              <p className="text-muted-foreground">Last Updated</p>
              <p className="font-medium text-foreground">
                {(%s.updatedAt || %s.updated_at) ? new Date(%s.updatedAt || %s.updated_at || '').toLocaleDateString() : '-'}
              </p>
            </div>
          </div>
        </div>
      </div>
    </form>
  );
});

%sEditForm.displayName = '%sEditForm';
`, pageName, pageName, resourceName, receiver.getUniqueIcons(fields),
		pageName, pageName,
		pageName, pageName, camelResource,
		pageName, initialData.String(),
		validationRules.String(),
		pluralPage, camelResource,
		pageName, pageName,
		formFields.String(),
		camelResource,
		camelResource, camelResource, camelResource, camelResource,
		camelResource, camelResource, camelResource, camelResource,
		pageName, pageName)
}

// generateDetailViewContent generates the DetailView.tsx file
func (receiver *UIMaker) generateDetailViewContent(pageName, pluralPage string, fields []Field) string {
	camelResource := strings.ToLower(pageName[:1]) + pageName[1:]
	resourceName := strings.ToLower(pageName)

	// Generate detail fields
	var detailFields strings.Builder
	for _, field := range fields {
		if !field.ShowInDetail {
			continue
		}
		fieldLower := strings.ToLower(field.Name[:1]) + field.Name[1:]

		detailFields.WriteString(fmt.Sprintf(`
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <%s className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">%s</p>
                <p className="font-medium text-foreground">{%s.%s}</p>
              </div>
            </div>
`, field.Icon, field.Label, camelResource, fieldLower))
	}

	return fmt.Sprintf(`import React from 'react';
import { Calendar, %s } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { CrudDetailViewProps } from '@/types/crud';
import { %s } from '@/types/%s';

export function %sDetailView({
  item: %s,
  onEdit,
  onClose,
  canEdit
}: CrudDetailViewProps<%s>) {
  const formatDate = (date: string | Date | null) => {
    if (!date) return 'Not specified';
    return new Date(date).toLocaleDateString('en-US', {
      month: 'long',
      day: 'numeric',
      year: 'numeric'
    });
  };

  return (
    <div className="space-y-6">
      <div className="space-y-6">
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">%s Information</h3>
          <div className="space-y-4">%s
          </div>
        </div>

        <Separator />

        {/* Metadata Section */}
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">Metadata</h3>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <p className="text-sm text-muted-foreground">Created</p>
              <p className="font-medium text-sm text-foreground">{formatDate(%s.createdAt || %s.created_at || null)}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Last Updated</p>
              <p className="font-medium text-sm text-foreground">{formatDate(%s.updatedAt || %s.updated_at || null)}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
`, receiver.getUniqueIcons(fields), pageName, resourceName,
		pageName, camelResource, pageName,
		pageName, detailFields.String(),
		camelResource, camelResource, camelResource, camelResource)
}

// generatePageConfigContent generates the PageConfig.tsx file
func (receiver *UIMaker) generatePageConfigContent(pageName, pluralPage string, fields []Field) string {
	camelResource := strings.ToLower(pageName[:1]) + pageName[1:]

	return fmt.Sprintf(`import React from 'react';
import { Download, Upload, BarChart3 } from 'lucide-react';
import {
  StatsCardConfig,
  PageActionConfig,
  SimpleFilterConfig
} from '@/lib/crud-page-utils';

/**
 * Stats card configurations for %s
 */
export const %sStatsConfigs: StatsCardConfig[] = [
  // TODO: Configure your stats cards here
];

/**
 * Simple filter configurations for %s
 */
export const %sSimpleFilters = (stats: any): SimpleFilterConfig[] => [
  // TODO: Configure your simple filters here
];

/**
 * Page action configurations for %s
 */
export const get%sPageActions = (
  permissions: any,
  handlers: {
    onImport?: () => void;
    onExport?: () => void;
  }
): PageActionConfig[] => [
  // TODO: Configure your page actions here
];

/**
 * Bulk action configurations for %s
 */
export const %sBulkActions = {
  handleBulkDelete: (%sIds: number[]) => {
    const confirmMessage = `+"`"+`Are you sure you want to delete ${%sIds.length} item(s)?`+"`"+`;
    if (confirm(confirmMessage)) {
      // TODO: Implement bulk delete
    }
  },
};

import { router } from '@inertiajs/react';
`, pluralPage, camelResource,
		pluralPage, camelResource,
		pluralPage, pageName,
		pluralPage, camelResource, camelResource, camelResource)
}

// generateSectionsIndexContent generates the index.ts file for sections
func (receiver *UIMaker) generateSectionsIndexContent(pageName, pluralPage string, fields []Field) string {
	camelResource := strings.ToLower(pageName[:1]) + pageName[1:]

	return fmt.Sprintf(`export { %sDetailView } from './%sDetailView';
export { %sCreateForm } from './%sCreateForm';
export { %sEditForm } from './%sEditForm';
export {
  %sColumns,
  %sColumnsMobile,
  %sFilters
} from './%sColumns';
export {
  %sStatsConfigs,
  %sSimpleFilters,
  get%sPageActions,
  %sBulkActions
} from './%sPageConfig';
`, pageName, pageName,
		pageName, pageName,
		pageName, pageName,
		camelResource, camelResource, camelResource, pageName,
		camelResource, camelResource, pageName, camelResource, pageName)
}

// Helper functions

func (receiver *UIMaker) makeFieldsOptional(fields string) string {
	// Convert required fields to optional
	// First handle already optional fields (with ?) to prevent double question marks
	result := strings.ReplaceAll(fields, "?: ", ": ")
	// Then make all fields optional
	return strings.ReplaceAll(result, ": ", "?: ")
}

func (receiver *UIMaker) generateErrorFields(fields []Field) string {
	var result strings.Builder
	for _, field := range fields {
		fieldLower := strings.ToLower(field.Name[:1]) + field.Name[1:]
		result.WriteString(fmt.Sprintf("  %s?: string;\n", fieldLower))
	}
	return result.String()
}

func (receiver *UIMaker) generateValidationRules(fields []Field) string {
	var result strings.Builder
	for _, field := range fields {
		if !field.Required {
			continue
		}
		fieldLower := strings.ToLower(field.Name[:1]) + field.Name[1:]
		result.WriteString(fmt.Sprintf(`if (!formData.%s) {
      newErrors.%s = '%s is required';
    }
    `, fieldLower, fieldLower, field.Label))
	}
	return result.String()
}

func (receiver *UIMaker) getUniqueIcons(fields []Field) string {
	iconsMap := make(map[string]bool)
	for _, field := range fields {
		iconsMap[field.Icon] = true
	}
	var icons []string
	for icon := range iconsMap {
		icons = append(icons, icon)
	}
	return strings.Join(icons, ", ")
}

func (receiver *UIMaker) getFirstTableField(fields []Field) string {
	for _, field := range fields {
		if field.ShowInTable {
			return strings.ToLower(field.Name[:1]) + field.Name[1:]
		}
	}
	return "name"
}
