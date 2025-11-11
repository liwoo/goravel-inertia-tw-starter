import React from 'react';
import { Sme } from '@/types/sme';
import { CrudColumn, CrudFilter } from '@/types/crud';
import { Badge } from '@/components/ui/badge';
import { FileText, User, FolderOpen, Phone, Mail, MapPin, Building2, Calendar } from 'lucide-react';

/**
 * Sme table columns configuration
 * Only showing the most relevant columns for listing
 */
export const smeColumns: CrudColumn<Sme>[] = [
  {
    key: 'usmeNumber',
    label: 'USME Number',
    sortable: true,
    className: 'min-w-[150px]',
    render: (sme) => (
      <div className="flex items-start gap-3">
        <div className="p-2 rounded-lg bg-muted">
          <FileText className="h-4 w-4 text-muted-foreground" />
        </div>
        <div className="flex-1 space-y-1">
          <p className="font-medium text-foreground">{sme.usmeNumber}</p>
        </div>
      </div>
    ),
  },
  {
    key: 'name',
    label: 'Business Name',
    sortable: true,
    className: 'min-w-[200px]',
    render: (sme) => (
      <div className="flex items-start gap-3">
        <div className="p-2 rounded-lg bg-muted">
          <Building2 className="h-4 w-4 text-muted-foreground" />
        </div>
        <div className="flex-1 space-y-1">
          <p className="font-medium text-foreground">{sme.name}</p>
        </div>
      </div>
    ),
  },
  {
    key: 'businessCategory',
    label: 'Category',
    sortable: true,
    className: 'min-w-[150px]',
    render: (sme) => (
      <Badge variant="secondary" className="font-normal">
        {sme.businessCategory}
      </Badge>
    ),
  },
  {
    key: 'sector',
    label: 'Sector',
    sortable: true,
    className: 'min-w-[150px]',
    render: (sme) => (
      <Badge variant="outline" className="font-normal">
        {sme.sector}
      </Badge>
    ),
  },
  {
    key: 'region',
    label: 'Region',
    sortable: true,
    className: 'min-w-[120px]',
    render: (sme) => (
      <div className="flex items-center gap-2">
        <MapPin className="h-3 w-3 text-muted-foreground" />
        <span className="text-sm">{sme.region || '-'}</span>
      </div>
    ),
  },
  {
    key: 'district',
    label: 'District',
    sortable: true,
    className: 'min-w-[120px]',
    render: (sme) => (
      <span className="text-sm text-muted-foreground">{sme.district || '-'}</span>
    ),
  },
  {
    key: 'contactPhone',
    label: 'Phone',
    sortable: true,
    className: 'min-w-[120px]',
    render: (sme) => (
      <div className="flex items-center gap-2">
        <Phone className="h-3 w-3 text-muted-foreground" />
        <span className="text-sm">{sme.contactPhone}</span>
      </div>
    ),
  },
  {
    key: 'contactEmail',
    label: 'Email',
    sortable: true,
    className: 'min-w-[180px]',
    render: (sme) => (
      <div className="flex items-center gap-2">
        <Mail className="h-3 w-3 text-muted-foreground" />
        <span className="text-sm truncate">{sme.contactEmail}</span>
      </div>
    ),
  },
  {
    key: 'createdAt',
    label: 'Created',
    sortable: true,
    className: 'w-32',
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
    render: (sme) => (
      <div className="space-y-2">
        <div className="font-medium text-foreground">{sme.name}</div>
        <div className="text-sm text-muted-foreground">{sme.usmeNumber}</div>
        <div className="flex gap-2 flex-wrap">
          <Badge variant="secondary" className="text-xs">
            {sme.businessCategory}
          </Badge>
          {sme.region && (
            <span className="text-xs text-muted-foreground flex items-center gap-1">
              <MapPin className="h-3 w-3" />
              {sme.region}
            </span>
          )}
        </div>
      </div>
    ),
  },
];

/**
 * Sme filters configuration
 * These match the filter fields configured in sme_service.go
 */
export const smeFilters: CrudFilter[] = [
  {
    key: 'businessCategory',
    label: 'Business Category',
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
    key: 'region',
    label: 'Region',
    type: 'text',
    placeholder: 'Filter by region',
  },
  {
    key: 'district',
    label: 'District',
    type: 'text',
    placeholder: 'Filter by district',
  },
];
