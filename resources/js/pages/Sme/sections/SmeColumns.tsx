import React from 'react';
import { Sme } from '@/types/sme';
import { CrudColumn, CrudFilter } from '@/types/crud';
import { Badge } from '@/components/ui/badge';
import { CheckCircle2, XCircle, CircleCheck, CircleX } from 'lucide-react';
import { ScoreMeter } from '@/components/ui/score-meter';

/**
 * Sme table columns configuration
 * Focused on business-critical information for quick scanning
 */
export const smeColumns: CrudColumn<Sme>[] = [
  {
    key: 'usmeNumber',
    label: 'USME Number',
    sortable: true,
    className: 'min-w-[140px]',
    render: (sme) => (
      <div className="font-mono text-sm font-medium text-foreground">
        {sme.usmeNumber || (sme as any).usme_number || '-'}
      </div>
    ),
  },
  {
    key: 'name',
    label: 'Business Name',
    sortable: true,
    className: 'min-w-[200px]',
    render: (sme) => (
      <div className="font-medium text-foreground">{sme.name}</div>
    ),
  },
  {
    key: 'isActive',
    label: 'Status',
    sortable: true,
    className: 'w-24 text-center',
    render: (sme) => {
      // Default to true if is_active is undefined (backward compatibility)
      const isActive = sme.isActive ?? (sme as any).is_active ?? true;
      return (
        <div className="flex justify-center">
          {isActive ? (
            <Badge variant="default" className="bg-emerald-600 hover:bg-emerald-600 text-xs">
              <CircleCheck className="h-3 w-3 mr-1" />
              Active
            </Badge>
          ) : (
            <Badge variant="secondary" className="text-muted-foreground text-xs">
              <CircleX className="h-3 w-3 mr-1" />
              Inactive
            </Badge>
          )}
        </div>
      );
    },
  },
  {
    key: 'businessCategory',
    label: 'Category',
    sortable: true,
    className: 'min-w-[130px]',
    render: (sme) => (
      <Badge variant="secondary" className="font-normal">
        {sme.businessCategory || (sme as any).business_category || '-'}
      </Badge>
    ),
  },
  {
    key: 'sector',
    label: 'Sector',
    sortable: true,
    className: 'min-w-[130px]',
    render: (sme) => (
      <Badge variant="outline" className="font-normal">
        {sme.sector}
      </Badge>
    ),
  },
  {
    key: 'district',
    label: 'District',
    sortable: true,
    className: 'min-w-[120px]',
    render: (sme) => (
      <span className="text-sm">{sme.district || '-'}</span>
    ),
  },
  {
    key: 'contactPhone',
    label: 'Contact',
    sortable: true,
    className: 'min-w-[130px]',
    render: (sme) => (
      <span className="text-sm font-mono">{sme.contactPhone}</span>
    ),
  },
  {
    key: 'registrationNumber',
    label: 'Registered',
    sortable: true,
    className: 'w-28 text-center',
    render: (sme) => (
      <div className="flex justify-center">
        {sme.registrationNumber ? (
          <CheckCircle2 className="h-4 w-4 text-green-600" />
        ) : (
          <XCircle className="h-4 w-4 text-gray-400" />
        )}
      </div>
    ),
  },
  {
    key: 'formalisationScore',
    label: 'Score',
    sortable: true,
    className: 'w-20 text-center',
    render: (sme) => {
      // Get formalisation score from relationship (camelCase or snake_case)
      const score =
        sme.businessFormalisation?.formalisationScore ??
        sme.businessFormalisation?.formalisation_score ??
        sme.business_formalisation?.formalisationScore ??
        sme.business_formalisation?.formalisation_score ??
        null;

      if (score === null || score === undefined) {
        return (
          <div className="flex justify-center items-center">
            <span className="text-sm text-muted-foreground">-</span>
          </div>
        );
      }

      return (
        <div className="flex justify-center">
          <ScoreMeter score={score} size="sm" showLabel={true} />
        </div>
      );
    },
  },
  {
    key: 'createdAt',
    label: 'Added',
    sortable: true,
    className: 'w-28',
    render: (sme) => {
      const dateValue = sme.createdAt || sme.created_at;
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
export const smeColumnsMobile: CrudColumn<Sme>[] = [
  {
    key: 'combined',
    label: 'SME',
    sortable: false,
    render: (sme) => {
      const isActive = sme.isActive ?? (sme as any).is_active ?? true;

      // Get formalisation score
      const score =
        sme.businessFormalisation?.formalisationScore ??
        sme.businessFormalisation?.formalisation_score ??
        sme.business_formalisation?.formalisationScore ??
        sme.business_formalisation?.formalisation_score ??
        null;

      return (
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <div className="font-medium text-foreground">{sme.name}</div>
            <div className="flex items-center gap-2 flex-shrink-0 ml-2">
              {!isActive && (
                <Badge variant="secondary" className="text-xs">Inactive</Badge>
              )}
              {sme.registrationNumber && (
                <CheckCircle2 className="h-4 w-4 text-green-600" />
              )}
            </div>
          </div>
          <div className="text-sm text-muted-foreground font-mono">{sme.usmeNumber || (sme as any).usme_number || '-'}</div>
          <div className="flex gap-2 flex-wrap items-center">
            <Badge variant="secondary" className="text-xs">
              {sme.businessCategory || (sme as any).business_category || '-'}
            </Badge>
            <Badge variant="outline" className="text-xs">
              {sme.sector}
            </Badge>
            {sme.district && (
              <span className="text-xs text-muted-foreground">
                {sme.district}
              </span>
            )}
          </div>
          {score !== null && score !== undefined && (
            <div className="flex items-center gap-2 pt-1">
              <span className="text-xs text-muted-foreground min-w-[40px]">Score:</span>
              <ScoreMeter score={score} size="sm" showLabel={true} />
            </div>
          )}
        </div>
      );
    },
  },
];

/**
 * Sme filters configuration
 * These match the filter fields configured in sme_service.go
 */
export const smeFilters: CrudFilter[] = [
  {
    key: 'businessCategory',
    label: 'Category',
    type: 'text',
    placeholder: 'Filter by category',
  },
  {
    key: 'sector',
    label: 'Sector',
    type: 'text',
    placeholder: 'Filter by sector',
  },
  {
    key: 'district',
    label: 'District',
    type: 'text',
    placeholder: 'Filter by district',
  },
];
