import React from 'react';
import { Badge } from '@/components/ui/badge';
import { Shield, User as UserIcon, CheckCircle, XCircle } from 'lucide-react';
import { User } from '@/types/user';
import { CrudColumn, CrudAction } from '@/types/crud';

// Column definitions for User table with improved theming
export const userColumns: CrudColumn<User>[] = [
  {
    key: 'id',
    label: 'ID',
    sortable: true,
    width: '80px',
    render: (user) => (
      <span className="font-mono text-sm text-muted-foreground">#{user.id.toString().padStart(6, '0')}</span>
    ),
  },
  {
    key: 'name',
    label: 'Name',
    sortable: true,
    render: (user) => (
      <div className="flex items-center gap-3">
        <div className="p-1.5 rounded-lg bg-muted">
          {user.is_super_admin ? (
            <Shield className="h-4 w-4 text-blue-500 dark:text-blue-400" />
          ) : (
            <UserIcon className="h-4 w-4 text-muted-foreground" />
          )}
        </div>
        <div>
          <div className="font-medium text-foreground">{user.name}</div>
          <div className="text-sm text-muted-foreground">{user.email}</div>
        </div>
      </div>
    ),
  },
  {
    key: 'roles',
    label: 'Roles',
    render: (user) => (
      <div className="flex flex-wrap gap-1">
        {user.roles && user.roles.length > 0 ? (
          user.roles.map((role) => (
            <Badge 
              key={role.id} 
              variant="secondary" 
              className="text-xs bg-secondary/50 dark:bg-secondary/30"
            >
              {role.name}
            </Badge>
          ))
        ) : (
          <span className="text-sm text-muted-foreground">No roles</span>
        )}
      </div>
    ),
  },
  {
    key: 'is_active',
    label: 'Status',
    sortable: true,
    width: '120px',
    render: (user) => (
      <div className="flex items-center gap-2">
        {user.is_active ? (
          <Badge className="bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400 flex items-center gap-1">
            <CheckCircle className="h-3 w-3" />
            Active
          </Badge>
        ) : (
          <Badge className="bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400 flex items-center gap-1">
            <XCircle className="h-3 w-3" />
            Inactive
          </Badge>
        )}
      </div>
    ),
  },
  {
    key: 'email_verified',
    label: 'Verified',
    sortable: true,
    width: '100px',
    render: (user) => (
      <div className="flex justify-center">
        {user.email_verified ? (
          <CheckCircle className="h-4 w-4 text-green-500 dark:text-green-400" />
        ) : (
          <XCircle className="h-4 w-4 text-muted-foreground" />
        )}
      </div>
    ),
  },
  {
    key: 'created_at',
    label: 'Member Since',
    sortable: true,
    width: '140px',
    render: (user) => (
      <span className="text-sm text-muted-foreground">
        {new Date(user.created_at).toLocaleDateString('en-US', {
          month: 'short',
          day: 'numeric',
          year: 'numeric'
        })}
      </span>
    ),
  },
];

// Mobile-friendly columns with better visual hierarchy
export const userColumnsMobile: CrudColumn<User>[] = [
  {
    key: 'user_info',
    label: 'User',
    render: (user) => (
      <div className="space-y-2">
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-lg bg-muted">
            {user.is_super_admin ? (
              <Shield className="h-5 w-5 text-blue-500 dark:text-blue-400" />
            ) : (
              <UserIcon className="h-5 w-5 text-muted-foreground" />
            )}
          </div>
          <div className="flex-1">
            <div className="font-medium text-foreground">{user.name}</div>
            <div className="text-sm text-muted-foreground">{user.email}</div>
          </div>
        </div>
        <div className="flex items-center gap-2 pl-12">
          {user.is_active ? (
            <Badge className="bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400 text-xs">
              Active
            </Badge>
          ) : (
            <Badge className="bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400 text-xs">
              Inactive
            </Badge>
          )}
          {user.email_verified && (
            <Badge variant="secondary" className="text-xs">
              Verified
            </Badge>
          )}
          {user.roles && user.roles.length > 0 && (
            <Badge variant="secondary" className="text-xs bg-secondary/50 dark:bg-secondary/30">
              {user.roles[0].name}
            </Badge>
          )}
        </div>
      </div>
    ),
  },
];

// Additional actions factory for User-specific actions
export const createUserAdditionalActions = (callbacks: {
  onActivate?: (id: number) => void;
  onDeactivate?: (id: number) => void;
  onResetPassword?: (id: number) => void;
  onImpersonate?: (id: number) => void;
  onSendWelcomeEmail?: (id: number) => void;
}): CrudAction<User>[] => {
  const actions: CrudAction<User>[] = [];

  if (callbacks.onActivate) {
    actions.push({
      key: 'activate',
      label: 'Activate',
      icon: <CheckCircle className="h-4 w-4 text-green-600" />,
      onClick: (user: User) => callbacks.onActivate!(user.id),
      disabled: (user: User) => user.is_active,
    });
  }

  if (callbacks.onDeactivate) {
    actions.push({
      key: 'deactivate',
      label: 'Deactivate',
      icon: <XCircle className="h-4 w-4 text-orange-600" />,
      onClick: (user: User) => callbacks.onDeactivate!(user.id),
      disabled: (user: User) => !user.is_active,
    });
  }

  if (callbacks.onResetPassword) {
    actions.push({
      key: 'reset-password',
      label: 'Reset Password',
      icon: <Shield className="h-4 w-4 text-blue-600" />,
      onClick: (user: User) => callbacks.onResetPassword!(user.id),
    });
  }

  if (callbacks.onImpersonate) {
    actions.push({
      key: 'impersonate',
      label: 'Impersonate',
      icon: <UserIcon className="h-4 w-4 text-purple-600" />,
      onClick: (user: User) => callbacks.onImpersonate!(user.id),
      disabled: (user: User) => user.is_super_admin,
    });
  }

  if (callbacks.onSendWelcomeEmail) {
    actions.push({
      key: 'send-welcome',
      label: 'Send Welcome Email',
      icon: <UserIcon className="h-4 w-4 text-blue-500" />,
      onClick: (user: User) => callbacks.onSendWelcomeEmail!(user.id),
      disabled: (user: User) => !user.is_active,
    });
  }

  return actions;
};

// Filter definitions with improved styling
export const userFilters = [
  {
    key: 'is_active',
    label: 'Status',
    type: 'select' as const,
    options: [
      { label: 'All Status', value: '' },
      { label: 'Active', value: 'true' },
      { label: 'Inactive', value: 'false' },
    ],
  },
  {
    key: 'is_super_admin',
    label: 'Admin Type',
    type: 'select' as const,
    options: [
      { label: 'All Types', value: '' },
      { label: 'Super Admin', value: 'true' },
      { label: 'Regular User', value: 'false' },
    ],
  },
  {
    key: 'email_verified',
    label: 'Email Status',
    type: 'select' as const,
    options: [
      { label: 'All', value: '' },
      { label: 'Verified', value: 'true' },
      { label: 'Unverified', value: 'false' },
    ],
  },
];

// Quick filter buttons with improved icons
export const userQuickFilters = [
  {
    key: 'all',
    label: 'All Users',
    icon: <UserIcon className="h-4 w-4" />,
    filters: {},
  },
  {
    key: 'active',
    label: 'Active',
    icon: <CheckCircle className="h-4 w-4" />,
    filters: { is_active: true },
  },
  {
    key: 'inactive',
    label: 'Inactive',
    icon: <XCircle className="h-4 w-4" />,
    filters: { is_active: false },
  },
  {
    key: 'super_admins',
    label: 'Super Admins',
    icon: <Shield className="h-4 w-4" />,
    filters: { is_super_admin: true },
  },
  {
    key: 'verified',
    label: 'Verified',
    icon: <CheckCircle className="h-4 w-4" />,
    filters: { email_verified: true },
  },
];