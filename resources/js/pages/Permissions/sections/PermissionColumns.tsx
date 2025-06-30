import React from 'react';
import { Permission } from '@/types/permissions';
import { CrudColumn, CrudFilter } from '@/types/crud';
import { Badge } from '@/components/ui/badge';
import { Shield, Tag, Calendar, Check, X, Settings } from 'lucide-react';

/**
 * Permission table columns configuration
 */
export const permissionColumns: CrudColumn<Permission>[] = [
  {
    key: 'name',
    label: 'Permission Details',
    sortable: true,
    className: 'min-w-[300px]',
    render: (permission) => (
      <div className="flex items-start gap-3">
        <div className="p-2 rounded-lg bg-muted">
          <Shield className="h-5 w-5 text-muted-foreground" />
        </div>
        <div className="space-y-1">
          <div className="font-medium text-foreground">{permission.name}</div>
          <div className="text-sm text-muted-foreground">{permission.description}</div>
          <div className="text-xs text-muted-foreground">
            Slug: {permission.slug}
          </div>
        </div>
      </div>
    ),
  },
  {
    key: 'category',
    label: 'Category',
    sortable: true,
    className: 'w-32',
    render: (permission) => (
      <Badge variant="secondary" className="flex items-center gap-1">
        <Tag className="h-3 w-3" />
        {permission.category}
      </Badge>
    ),
  },
  {
    key: 'action',
    label: 'Action',
    sortable: true,
    className: 'w-24',
    render: (permission) => (
      <Badge variant="outline" className="font-mono text-xs">
        {permission.action}
      </Badge>
    ),
  },
  {
    key: 'resource',
    label: 'Resource',
    sortable: true,
    className: 'w-28',
    render: (permission) => (
      <span className="text-sm text-muted-foreground">
        {permission.resource || '-'}
      </span>
    ),
  },
  {
    key: 'is_active',
    label: 'Status',
    sortable: true,
    className: 'w-20',
    render: (permission) => (
      <Badge 
        className={permission.is_active 
          ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400 flex items-center gap-1'
          : 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400 flex items-center gap-1'
        }
      >
        {permission.is_active ? <Check className="h-3 w-3" /> : <X className="h-3 w-3" />}
        {permission.is_active ? 'Active' : 'Inactive'}
      </Badge>
    ),
  },
  {
    key: 'requires_ownership',
    label: 'Ownership',
    className: 'w-24 text-center',
    render: (permission) => (
      <div className="flex justify-center">
        {permission.requires_ownership ? (
          <Check className="h-4 w-4 text-green-500" />
        ) : (
          <X className="h-4 w-4 text-gray-400" />
        )}
      </div>
    ),
  },
  {
    key: 'can_delegate',
    label: 'Delegatable',
    className: 'w-24 text-center',
    render: (permission) => (
      <div className="flex justify-center">
        {permission.can_delegate ? (
          <Check className="h-4 w-4 text-green-500" />
        ) : (
          <X className="h-4 w-4 text-gray-400" />
        )}
      </div>
    ),
  },
  {
    key: 'created_at',
    label: 'Created',
    sortable: true,
    className: 'w-28',
    render: (permission) => (
      <div className="text-sm text-muted-foreground flex items-center">
        <Calendar className="w-3 h-3 mr-1" />
        {new Date(permission.created_at).toLocaleDateString('en-US', {
          month: 'short',
          day: 'numeric',
          year: 'numeric'
        })}
      </div>
    ),
  },
];

/**
 * Compact permission columns for mobile/smaller screens
 */
export const permissionColumnsMobile: CrudColumn<Permission>[] = [
  {
    key: 'name',
    label: 'Permission',
    sortable: true,
    render: (permission) => (
      <div className="space-y-3">
        <div className="flex items-start gap-3">
          <div className="p-2 rounded-lg bg-muted">
            <Shield className="h-5 w-5 text-muted-foreground" />
          </div>
          <div className="flex-1 space-y-1">
            <div className="font-medium text-foreground">{permission.name}</div>
            <div className="text-sm text-muted-foreground">{permission.description}</div>
          </div>
        </div>
        <div className="flex items-center justify-between pl-12">
          <div className="flex items-center gap-2">
            <Badge variant="secondary" className="text-xs">
              {permission.category}
            </Badge>
            <Badge variant="outline" className="text-xs font-mono">
              {permission.action}
            </Badge>
          </div>
          <div>
            <Badge 
              className={permission.is_active 
                ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400 flex items-center gap-1 text-xs'
                : 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400 flex items-center gap-1 text-xs'
              }
            >
              {permission.is_active ? <Check className="h-3 w-3" /> : <X className="h-3 w-3" />}
              {permission.is_active ? 'Active' : 'Inactive'}
            </Badge>
          </div>
        </div>
      </div>
    ),
  },
];

/**
 * Permission filters configuration
 */
export const permissionFilters: CrudFilter[] = [
  {
    key: 'category',
    label: 'Category',
    type: 'select',
    options: [
      { value: '', label: 'All Categories' },
      { value: 'books', label: 'Books' },
      { value: 'users', label: 'Users' },
      { value: 'roles', label: 'Roles' },
      { value: 'system', label: 'System' },
      { value: 'reports', label: 'Reports' },
    ],
  },
  {
    key: 'action',
    label: 'Action',
    type: 'select',
    options: [
      { value: '', label: 'All Actions' },
      { value: 'create', label: 'Create' },
      { value: 'read', label: 'Read' },
      { value: 'update', label: 'Update' },
      { value: 'delete', label: 'Delete' },
      { value: 'manage', label: 'Manage' },
      { value: 'assign', label: 'Assign' },
      { value: 'export', label: 'Export' },
      { value: 'configure', label: 'Configure' },
    ],
  },
  {
    key: 'is_active',
    label: 'Status',
    type: 'select',
    options: [
      { value: '', label: 'All Status' },
      { value: 'true', label: 'Active' },
      { value: 'false', label: 'Inactive' },
    ],
  },
  {
    key: 'requires_ownership',
    label: 'Requires Ownership',
    type: 'select',
    options: [
      { value: '', label: 'All' },
      { value: 'true', label: 'Yes' },
      { value: 'false', label: 'No' },
    ],
  },
  {
    key: 'can_delegate',
    label: 'Can Delegate',
    type: 'select',
    options: [
      { value: '', label: 'All' },
      { value: 'true', label: 'Yes' },
      { value: 'false', label: 'No' },
    ],
  },
];

/**
 * Quick filter buttons for common permission queries
 */
export const permissionQuickFilters = [
  {
    key: 'all',
    label: 'All Permissions',
    icon: <Shield className="h-4 w-4" />,
    filters: {},
  },
  {
    key: 'active',
    label: 'Active',
    icon: <Check className="h-4 w-4 text-green-500" />,
    filters: { is_active: 'true' },
  },
  {
    key: 'inactive',
    label: 'Inactive',
    icon: <X className="h-4 w-4 text-gray-500" />,
    filters: { is_active: 'false' },
  },
  {
    key: 'books',
    label: 'Books',
    icon: <Tag className="h-4 w-4 text-blue-500" />,
    filters: { category: 'books' },
  },
  {
    key: 'users',
    label: 'Users',
    icon: <Tag className="h-4 w-4 text-purple-500" />,
    filters: { category: 'users' },
  },
  {
    key: 'system',
    label: 'System',
    icon: <Settings className="h-4 w-4 text-orange-500" />,
    filters: { category: 'system' },
  },
] as const;