import React from 'react';
import { Bdsp } from '@/types/bdsp';
import { CrudColumn, CrudFilter } from '@/types/crud';
import { Badge } from '@/components/ui/badge';
import { User, MapPin, Tag, FileText } from 'lucide-react';

/**
 * Bdsp table columns configuration
 */
export const bdspColumns: CrudColumn<Bdsp>[] = [
  {
    key: 'name',
    label: 'Name',
    sortable: true,
    className: 'min-w-[150px]',
    render: (bdsp) => (
      <div className="flex items-start gap-3">
        <div className="flex items-center gap-2">
          <User className="h-3 w-3 text-muted-foreground" />
          <p className="font-medium text-foreground">{bdsp.name}</p>
        </div>
      </div>
    ),
  },
  {
    key: 'postalAddress',
    label: 'Postal Address',
    sortable: true,
    className: 'min-w-[150px]',
    render: (bdsp) => (
      <div className="flex items-start gap-3">
        <div className="flex items-center gap-2">
          <MapPin className="h-3 w-3 text-muted-foreground" />
          <p className="font-medium text-foreground">{bdsp.postal_address}</p>
        </div>
      </div>
    ),
  },
  {
    key: 'physicalAddress',
    label: 'Physical Address',
    sortable: true,
    className: 'min-w-[150px]',
    render: (bdsp) => (
      <div className="flex items-start gap-3">
        <div className="flex items-center gap-2">
          <MapPin className="h-3 w-3 text-muted-foreground" />
          <p className="font-medium text-foreground">{bdsp.physical_address}</p>
        </div>
      </div>
    ),
  },
  {
    key: 'registrationStatus',
    label: 'Registration Status',
    sortable: true,
    className: 'min-w-[150px]',
    render: (bdsp) => {
      const status = bdsp.registration_status || 'Pending';
      let variant: "default" | "secondary" | "destructive" | "outline" = "secondary";

      switch (status.toLowerCase()) {
        case 'active':
          variant = "default";
          break;
        case 'rejected':
          variant = "destructive";
          break;
        case 'inactive':
          variant = "outline";
          break;
        default:
          variant = "secondary";
      }

      return (
        <Badge variant={variant} className="capitalize">
          {status}
        </Badge>
      );
    },
  },
  {
    key: 'partners',
    label: 'Partners',
    sortable: false,
    className: 'min-w-[200px]',
    render: (bdsp) => (
      <div className="flex flex-wrap gap-1">
        {bdsp.partners && bdsp.partners.length > 0 ? (
          bdsp.partners.slice(0, 2).map((partner, i) => (
            <Badge key={i} variant="outline" className="text-xs px-1.5 py-0 h-5">
              {partner}
            </Badge>
          ))
        ) : (
          <span className="text-xs text-muted-foreground">-</span>
        )}
        {bdsp.partners && bdsp.partners.length > 2 && (
          <Badge variant="secondary" className="text-xs px-1.5 py-0 h-5">
            +{bdsp.partners.length - 2}
          </Badge>
        )}
      </div>
    ),
  },
  {
    key: 'productTypes',
    label: 'Product Types',
    sortable: false,
    className: 'min-w-[200px]',
    render: (bdsp) => (
      <div className="flex flex-wrap gap-1">
        {bdsp.product_types && bdsp.product_types.length > 0 ? (
          bdsp.product_types.slice(0, 2).map((type, i) => (
            <Badge key={i} variant="outline" className="text-xs px-1.5 py-0 h-5">
              {type}
            </Badge>
          ))
        ) : (
          <span className="text-xs text-muted-foreground">-</span>
        )}
        {bdsp.product_types && bdsp.product_types.length > 2 && (
          <Badge variant="secondary" className="text-xs px-1.5 py-0 h-5">
            +{bdsp.product_types.length - 2}
          </Badge>
        )}
      </div>
    ),
  },
  {
    key: 'createdAt',
    label: 'Created',
    sortable: true,
    className: 'w-32',
    render: (bdsp) => {
      const dateValue = bdsp.createdAt || bdsp.created_at;
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
export const bdspColumnsMobile: CrudColumn<Bdsp>[] = [
  {
    key: 'combined',
    label: 'Bdsp',
    sortable: false,
    render: (bdsp) => (
      <div className="space-y-2">
        <div className="font-medium text-foreground">{bdsp.name}</div>
        <div className="text-sm text-muted-foreground">
          {new Date(bdsp.createdAt || bdsp.created_at || '').toLocaleDateString()}
        </div>
      </div>
    ),
  },
];

/**
 * Bdsp filters configuration
 */
export const bdspFilters: CrudFilter[] = [
  {
    key: 'name',
    label: 'Name',
    type: 'text',
    placeholder: 'Enter name',
  },
  {
    key: 'postalAddress',
    label: 'Postal Address',
    type: 'text',
    placeholder: 'Enter address',
  },
  {
    key: 'physicalAddress',
    label: 'Physical Address',
    type: 'text',
    placeholder: 'Enter address',
  },
  {
    key: 'registrationStatus',
    label: 'Registration Status',
    type: 'text',
    placeholder: 'Enter registration status',
  },
  {
    key: 'partners',
    label: 'Partners',
    type: 'text',
    placeholder: 'Enter partners',
  },
  {
    key: 'productTypes',
    label: 'Product Types',
    type: 'text',
    placeholder: 'Enter product types',
  }
];

