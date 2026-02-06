import React from 'react';
import { Application, ApplicationType } from '@/types/application';
import { CrudColumn, CrudFilter } from '@/types/crud';
import { Badge } from '@/components/ui/badge';
import { FileText, User, Mail, Phone, Tag, UserPlus, FileEdit } from 'lucide-react';

// Helper to get type display label
const getTypeLabel = (type?: ApplicationType): string => {
  switch (type) {
    case 'amend_formalisation':
      return 'Formalisation Amendment';
    case 'signup':
    default:
      return 'Sign Up';
  }
};

// Helper to get type badge style
const getTypeBadgeStyle = (type?: ApplicationType): string => {
  switch (type) {
    case 'amend_formalisation':
      return 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200';
    case 'signup':
    default:
      return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200';
  }
};

/**
 * Application table columns configuration
 */
export const applicationColumns: CrudColumn<Application>[] = [
  {
    key: 'type',
    label: 'Type',
    sortable: true,
    className: 'min-w-[150px]',
    render: (application) => {
      const type = application.type || 'signup';
      const Icon = type === 'amend_formalisation' ? FileEdit : UserPlus;
      return (
        <div className="flex items-center gap-2">
          <Icon className="h-4 w-4 text-muted-foreground" />
          <Badge className={getTypeBadgeStyle(type as ApplicationType)}>
            {getTypeLabel(type as ApplicationType)}
          </Badge>
        </div>
      );
    },
  },
  {
    key: 'sme',
    label: 'MSME',
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
    label: 'MSME Registration Number',
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
    label: 'MSME Tax Identification Number',
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
      const type = application.type || 'signup';
      const statusBadgeClassName = status === 'Approved' ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' :
                             status === 'Rejected' ? 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200' :
                             'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200';
      return (
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <div className="font-medium text-foreground">{application.sme}</div>
            <Badge className={statusBadgeClassName}>{status}</Badge>
          </div>
          <div className="flex items-center justify-between">
            <Badge className={getTypeBadgeStyle(type as ApplicationType)} variant="outline">
              {getTypeLabel(type as ApplicationType)}
            </Badge>
            <span className="text-sm text-muted-foreground">
              {new Date(application.createdAt || application.created_at || '').toLocaleDateString()}
            </span>
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
    key: 'type',
    label: 'Type',
    type: 'select',
    options: [
      { value: 'signup', label: 'Sign Up' },
      { value: 'amend_formalisation', label: 'Formalisation Amendment' },
    ],
    placeholder: 'All types',
  },
  {
    key: 'status',
    label: 'Status',
    type: 'select',
    options: [
      { value: 'Pending', label: 'Pending' },
      { value: 'Approved', label: 'Approved' },
      { value: 'Rejected', label: 'Rejected' },
    ],
    placeholder: 'All statuses',
  },
  {
    key: 'sme',
    label: 'MSME',
    type: 'text',
    placeholder: 'Enter MSME',
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
    label: 'MSME Registration Number',
    type: 'text',
    placeholder: 'Enter MSME registration number',
  },
  {
    key: 'sme_tax_identification_number',
    label: 'MSME Tax Identification Number',
    type: 'text',
    placeholder: 'Enter MSME tax identification number',
  },
];
