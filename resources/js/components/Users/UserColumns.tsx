import React from 'react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Edit, Eye, Trash2, Shield, User as UserIcon } from 'lucide-react';
import { User } from '@/types/user';
import { CrudColumn, CrudAction } from '@/types/crud';

// Column definitions for User table
export const userColumns: CrudColumn<User>[] = [
  {
    key: 'id',
    label: 'ID',
    sortable: true,
    width: '80px',
    render: (user) => (
      <span className="font-mono text-sm">#{user.id}</span>
    ),
  },
  {
    key: 'name',
    label: 'Name',
    sortable: true,
    render: (user) => (
      <div className="flex items-center space-x-2">
        {user.is_super_admin ? (
          <Shield className="h-4 w-4 text-blue-500" />
        ) : (
          <UserIcon className="h-4 w-4 text-gray-400" />
        )}
        <div className="font-medium">{user.name}</div>
      </div>
    ),
  },
  {
    key: 'email',
    label: 'Email',
    sortable: true,
    render: (user) => (
      <div className="text-sm text-muted-foreground">{user.email}</div>
    ),
  },
  {
    key: 'roles',
    label: 'Role',
    render: (user) => (
      <div className="flex flex-wrap gap-1">
        {user.roles && user.roles.length > 0 ? (
          user.roles.map((role) => (
            <Badge key={role.id} variant="secondary" className="text-xs">
              {role.name}
            </Badge>
          ))
        ) : (
          <span className="text-sm text-muted-foreground">No role</span>
        )}
      </div>
    ),
  },
  {
    key: 'is_active',
    label: 'Status',
    sortable: true,
    width: '100px',
    render: (user) => (
      <Badge variant={user.is_active ? 'default' : 'secondary'}>
        {user.is_active ? 'Active' : 'Inactive'}
      </Badge>
    ),
  },
  {
    key: 'created_at',
    label: 'Created',
    sortable: true,
    width: '120px',
    render: (user) => (
      <span className="text-sm text-muted-foreground">
        {new Date(user.created_at).toLocaleDateString()}
      </span>
    ),
  },
];

// Mobile-friendly columns
export const userColumnsMobile: CrudColumn<User>[] = [
  {
    key: 'user_info',
    label: 'User',
    render: (user) => (
      <div>
        <div className="flex items-center space-x-2">
          {user.is_super_admin ? (
            <Shield className="h-4 w-4 text-blue-500" />
          ) : (
            <UserIcon className="h-4 w-4 text-gray-400" />
          )}
          <div className="font-medium">{user.name}</div>
        </div>
        <div className="text-sm text-muted-foreground">{user.email}</div>
        <div className="flex items-center space-x-2 mt-1">
          <Badge variant={user.is_active ? 'default' : 'secondary'} className="text-xs">
            {user.is_active ? 'Active' : 'Inactive'}
          </Badge>
          {user.roles && user.roles.length > 0 && (
            <Badge variant="secondary" className="text-xs">
              {user.roles[0].name}
            </Badge>
          )}
        </div>
      </div>
    ),
  },
];

// Additional actions factory for User-specific actions (beyond the default View/Edit/Delete)
export const createUserAdditionalActions = (callbacks: {
  onActivate?: (id: number) => void;
  onDeactivate?: (id: number) => void;
  onResetPassword?: (id: number) => void;
  onImpersonate?: (id: number) => void;
  onSendWelcomeEmail?: (id: number) => void;
}): CrudAction<User>[] => {
  const actions: CrudAction<User>[] = [];

  // Status toggle actions
  if (callbacks.onActivate) {
    actions.push({
      key: 'activate',
      label: 'Activate',
      icon: <Shield className="h-4 w-4 text-green-600" />,
      onClick: (user: User) => callbacks.onActivate!(user.id),
      disabled: (user: User) => user.is_active, // Disable if already active
    });
  }

  if (callbacks.onDeactivate) {
    actions.push({
      key: 'deactivate',
      label: 'Deactivate',
      icon: <Shield className="h-4 w-4 text-orange-600" />,
      onClick: (user: User) => callbacks.onDeactivate!(user.id),
      disabled: (user: User) => !user.is_active, // Disable if already inactive
    });
  }

  // Reset password action
  if (callbacks.onResetPassword) {
    actions.push({
      key: 'reset-password',
      label: 'Reset Password',
      icon: <UserIcon className="h-4 w-4 text-blue-600" />,
      onClick: (user: User) => callbacks.onResetPassword!(user.id),
    });
  }

  // Impersonate user action (for super admins)
  if (callbacks.onImpersonate) {
    actions.push({
      key: 'impersonate',
      label: 'Impersonate',
      icon: <Eye className="h-4 w-4 text-purple-600" />,
      onClick: (user: User) => callbacks.onImpersonate!(user.id),
      disabled: (user: User) => user.is_super_admin, // Can't impersonate super admins
    });
  }

  // Send welcome email action
  if (callbacks.onSendWelcomeEmail) {
    actions.push({
      key: 'send-welcome',
      label: 'Send Welcome Email',
      icon: <UserIcon className="h-4 w-4 text-blue-500" />,
      onClick: (user: User) => callbacks.onSendWelcomeEmail!(user.id),
      disabled: (user: User) => !user.is_active, // Only for active users
    });
  }

  return actions;
};

// Filter definitions
export const userFilters = [
  {
    key: 'is_active',
    label: 'Status',
    type: 'select' as const,
    options: [
      { label: 'All', value: '' },
      { label: 'Active', value: 'true' },
      { label: 'Inactive', value: 'false' },
    ],
  },
  {
    key: 'is_super_admin',
    label: 'Admin Type',
    type: 'select' as const,
    options: [
      { label: 'All', value: '' },
      { label: 'Super Admin', value: 'true' },
      { label: 'Regular User', value: 'false' },
    ],
  },
  {
    key: 'role',
    label: 'Role',
    type: 'select' as const,
    options: [
      { label: 'All', value: '' },
      // Role options will be populated dynamically
    ],
  },
];

// Quick filter buttons
export const userQuickFilters = [
  {
    key: 'all',
    label: 'All Users',
    icon: React.createElement('span', { className: 'text-xs' }, '👥'),
    filters: {},
  },
  {
    key: 'active',
    label: 'Active',
    icon: React.createElement('span', { className: 'text-xs' }, '✅'),
    filters: { is_active: true },
  },
  {
    key: 'inactive',
    label: 'Inactive',
    icon: React.createElement('span', { className: 'text-xs' }, '❌'),
    filters: { is_active: false },
  },
  {
    key: 'super_admins',
    label: 'Super Admins',
    icon: React.createElement('span', { className: 'text-xs' }, '🛡️'),
    filters: { is_super_admin: true },
  },
];