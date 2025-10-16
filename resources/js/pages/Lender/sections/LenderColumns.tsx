import React from 'react';
import { Lender } from '@/types/lender';
import { CrudColumn, CrudFilter } from '@/types/crud';
import { Badge } from '@/components/ui/badge';
import { MapPin, User, Mail, Phone } from 'lucide-react';

/**
 * Lender table columns configuration
 */
export const lenderColumns: CrudColumn<Lender>[] = [
  {
    key: 'name',
    label: 'Name',
    sortable: true,
    className: 'min-w-[150px]',
    render: (lender) => (
      <div className="flex items-start gap-3">
        <div className="p-2 rounded-lg bg-muted">
          <User className="h-4 w-4 text-muted-foreground" />
        </div>
        <div className="flex-1 space-y-1">
          <p className="font-medium text-foreground">{lender.name}</p>
        </div>
      </div>
    ),
  },
  {
    key: 'email',
    label: 'Email',
    sortable: true,
    className: 'min-w-[150px]',
    render: (lender) => (
      <div className="flex items-start gap-3">
        <div className="p-2 rounded-lg bg-muted">
          <Mail className="h-4 w-4 text-muted-foreground" />
        </div>
        <div className="flex-1 space-y-1">
          <p className="font-medium text-foreground">{lender.email}</p>
        </div>
      </div>
    ),
  },
  {
    key: 'phone',
    label: 'Phone',
    sortable: true,
    className: 'min-w-[150px]',
    render: (lender) => (
      <div className="flex items-start gap-3">
        <div className="p-2 rounded-lg bg-muted">
          <Phone className="h-4 w-4 text-muted-foreground" />
        </div>
        <div className="flex-1 space-y-1">
          <p className="font-medium text-foreground">{lender.phone}</p>
        </div>
      </div>
    ),
  },
  {
    key: 'address',
    label: 'Address',
    sortable: true,
    className: 'min-w-[150px]',
    render: (lender) => (
      <div className="flex items-start gap-3">
        <div className="p-2 rounded-lg bg-muted">
          <MapPin className="h-4 w-4 text-muted-foreground" />
        </div>
        <div className="flex-1 space-y-1">
          <p className="font-medium text-foreground">{lender.address}</p>
        </div>
      </div>
    ),
  },
  {
    key: 'gender',
    label: 'Gender',
    sortable: true,
    className: 'min-w-[150px]',
    render: (lender) => (
      <div className="flex items-start gap-3">
        <div className="p-2 rounded-lg bg-muted">
          <User className="h-4 w-4 text-muted-foreground" />
        </div>
        <div className="flex-1 space-y-1">
          <p className="font-medium text-foreground">{lender.gender}</p>
        </div>
      </div>
    ),
  },
  {
    key: 'createdAt',
    label: 'Created',
    sortable: true,
    className: 'w-32',
    render: (lender) => {
      const dateValue = lender.createdAt || lender.created_at;
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
export const lenderColumnsMobile: CrudColumn<Lender>[] = [
  {
    key: 'combined',
    label: 'Lender',
    sortable: false,
    render: (lender) => (
      <div className="space-y-2">
        <div className="font-medium text-foreground">{lender.name}</div>
        <div className="text-sm text-muted-foreground">
          {new Date(lender.createdAt || lender.created_at || '').toLocaleDateString()}
        </div>
      </div>
    ),
  },
];

/**
 * Lender filters configuration
 */
export const lenderFilters: CrudFilter[] = [
  {
    key: 'name',
    label: 'Name',
    type: 'text',
    placeholder: 'Enter name',
  },
  {
    key: 'email',
    label: 'Email',
    type: 'text',
    placeholder: 'Enter email address',
  },
  {
    key: 'phone',
    label: 'Phone',
    type: 'text',
    placeholder: 'Enter phone number',
  },
  {
    key: 'address',
    label: 'Address',
    type: 'text',
    placeholder: 'Enter address',
  },
  {
    key: 'gender',
    label: 'Gender',
    type: 'text',
    placeholder: 'Enter gender',
  }
];
