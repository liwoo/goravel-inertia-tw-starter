import React from 'react';
import { Application } from '@/types/application';
import { CrudColumn, CrudFilter } from '@/types/crud';
import { Badge } from '@/components/ui/badge';
import { FileText, User, Mail, Phone, Tag } from 'lucide-react';

/**
 * Application table columns configuration
 */
export const applicationColumns: CrudColumn<Application>[] = [
  {
    key: 'sme',
    label: 'SME',
    sortable: true,
    className: 'min-w-[150px]',
    render: (application) => (
      <div className="flex-1 space-y-1">
        <p className="font-medium text-foreground">{application.sme}</p>
      </div>
    ),
  },
  {
    key: 'status',
    label: 'Status',
    sortable: true,
    className: 'min-w-[120px]',
    render: (application) => {
      const status = application.status;
      const variant = status === 'Approved' ? 'default' :
                      status === 'Rejected' ? 'destructive' :
                      'secondary';
      const className = status === 'Approved' ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' :
                        status === 'Rejected' ? 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200' :
                        'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200';
      return (
        <Badge variant={variant} className={className}>
          {status}
        </Badge>
      );
    },
  },
  {
    key: 'registrant_name',
    label: 'Registrant Name',
    sortable: true,
    className: 'min-w-[150px]',
    render: (application) => (
      <div className="flex-1 space-y-1">
        <p className="font-medium text-foreground">{application.registrant_name}</p>
      </div>
    ),
  },
  {
    key: 'email',
    label: 'Email',
    sortable: true,
    className: 'min-w-[150px]',
    render: (application) => (
      <div className="flex-1 space-y-1">
        <p className="font-medium text-foreground">{application.email}</p>
      </div>
    ),
  },
  {
    key: 'phone',
    label: 'Phone',
    sortable: true,
    className: 'min-w-[150px]',
    render: (application) => (
      <div className="flex-1 space-y-1">
        <p className="font-medium text-foreground">{application.phone}</p>
      </div>
    ),
  },
  {
    key: 'sme_registration_number',
    label: 'SME Registration Number',
    sortable: true,
    className: 'min-w-[150px]',
    render: (application) => (
      <div className="flex-1 space-y-1">
        <p className="font-medium text-foreground">{application.sme_registration_number}</p>
      </div>
    ),
  },
  {
    key: 'sme_tax_identification_number',
    label: 'SME Tax Identification Number',
    sortable: true,
    className: 'min-w-[150px]',
    render: (application) => (
      <div className="flex-1 space-y-1">
        <p className="font-medium text-foreground">{application.sme_tax_identification_number}</p>
      </div>
    ),
  },
  {
    key: 'createdAt',
    label: 'Created',
    sortable: true,
    className: 'w-32',
    render: (application) => {
      const dateValue = application.createdAt || application.created_at;
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
export const applicationColumnsMobile: CrudColumn<Application>[] = [
  {
    key: 'combined',
    label: 'Application',
    sortable: false,
    render: (application) => {
      const status = application.status;
      const badgeClassName = status === 'Approved' ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' :
                             status === 'Rejected' ? 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200' :
                             'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200';
      return (
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <div className="font-medium text-foreground">{application.sme}</div>
            <Badge className={badgeClassName}>{status}</Badge>
          </div>
          <div className="text-sm text-muted-foreground">
            {new Date(application.createdAt || application.created_at || '').toLocaleDateString()}
          </div>
        </div>
      );
    },
  },
];

/**
 * Application filters configuration
 */
export const applicationFilters: CrudFilter[] = [
  {
    key: 'sme',
    label: 'SME',
    type: 'text',
    placeholder: 'Enter SME',
  },
  {
    key: 'registrant_name',
    label: 'Registrant Name',
    type: 'text',
    placeholder: 'Enter registrant name',
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
    key: 'sme_registration_number',
    label: 'SME Registration Number',
    type: 'text',
    placeholder: 'Enter SME registration number',
  },
  {
    key: 'sme_tax_identification_number',
    label: 'SME Tax Identification Number',
    type: 'text',
    placeholder: 'Enter SME tax identification number',
  },
  {
    key: 'status',
    label: 'Status',
    type: 'text',
    placeholder: 'Enter status',
  }
];
