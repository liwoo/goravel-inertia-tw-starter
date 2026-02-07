import React from 'react';
import { Badge } from '@/components/ui/badge';
import { Shield, User as UserIcon, CheckCircle, XCircle } from 'lucide-react';
import { User } from '@/types/user';
import { CrudColumn, CrudAction } from '@/types/crud';
import { TFunction } from 'i18next';

// Column definitions for User table with improved theming
export function getUserColumns(t: TFunction): CrudColumn<User>[] {
  return [
    {
      key: 'id',
      label: t('columns.id'),
      sortable: true,
      width: '80px',
      render: (user) => (
        <span className="font-mono text-sm text-muted-foreground">#{user.id.toString().padStart(6, '0')}</span>
      ),
    },
    {
      key: 'name',
      label: t('columns.name'),
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
      label: t('columns.roles'),
      render: (user) => (
        <div className="flex flex-wrap gap-1">
          {user.roles && user.roles.length > 0 ? (
            user.roles.map((role, index) => (
              <Badge
                key={`${user.id}-role-${role.id}-${index}`}
                variant="secondary"
                className="text-xs bg-secondary/50 dark:bg-secondary/30"
              >
                {role.name}
              </Badge>
            ))
          ) : (
            <span className="text-sm text-muted-foreground">{t('columns.noRoles')}</span>
          )}
        </div>
      ),
    },
    {
      key: 'is_active',
      label: t('columns.status'),
      sortable: true,
      width: '120px',
      render: (user) => (
        <div className="flex items-center gap-2">
          {user.is_active ? (
            <Badge className="bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400 flex items-center gap-1">
              <CheckCircle className="h-3 w-3" />
              {t('columns.active')}
            </Badge>
          ) : (
            <Badge className="bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400 flex items-center gap-1">
              <XCircle className="h-3 w-3" />
              {t('columns.inactive')}
            </Badge>
          )}
        </div>
      ),
    },
    {
      key: 'email_verified',
      label: t('columns.verified'),
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
      label: t('columns.memberSince'),
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
}

// Mobile-friendly columns with better visual hierarchy
export function getUserColumnsMobile(t: TFunction): CrudColumn<User>[] {
  return [
    {
      key: 'user_info',
      label: t('columns.user'),
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
                {t('columns.active')}
              </Badge>
            ) : (
              <Badge className="bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400 text-xs">
                {t('columns.inactive')}
              </Badge>
            )}
            {user.email_verified && (
              <Badge variant="secondary" className="text-xs">
                {t('columns.verified')}
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
}

// Additional actions factory for User-specific actions
export function createUserAdditionalActions(t: TFunction, callbacks: {
  onActivate?: (id: number) => void;
  onDeactivate?: (id: number) => void;
  onResetPassword?: (id: number) => void;
  onImpersonate?: (id: number) => void;
  onSendWelcomeEmail?: (id: number) => void;
}): CrudAction<User>[] {
  const actions: CrudAction<User>[] = [];

  if (callbacks.onActivate) {
    actions.push({
      key: 'activate',
      label: t('actions.activate'),
      icon: <CheckCircle className="h-4 w-4 text-green-600" />,
      onClick: (user: User) => callbacks.onActivate!(user.id),
      disabled: (user: User) => user.is_active,
    });
  }

  if (callbacks.onDeactivate) {
    actions.push({
      key: 'deactivate',
      label: t('actions.deactivate'),
      icon: <XCircle className="h-4 w-4 text-orange-600" />,
      onClick: (user: User) => callbacks.onDeactivate!(user.id),
      disabled: (user: User) => !user.is_active,
    });
  }

  if (callbacks.onResetPassword) {
    actions.push({
      key: 'reset-password',
      label: t('actions.resetPassword'),
      icon: <Shield className="h-4 w-4 text-blue-600" />,
      onClick: (user: User) => callbacks.onResetPassword!(user.id),
    });
  }

  if (callbacks.onImpersonate) {
    actions.push({
      key: 'impersonate',
      label: t('actions.impersonate'),
      icon: <UserIcon className="h-4 w-4 text-purple-600" />,
      onClick: (user: User) => callbacks.onImpersonate!(user.id),
      disabled: (user: User) => user.is_super_admin,
    });
  }

  if (callbacks.onSendWelcomeEmail) {
    actions.push({
      key: 'send-welcome',
      label: t('actions.sendWelcomeEmail'),
      icon: <UserIcon className="h-4 w-4 text-blue-500" />,
      onClick: (user: User) => callbacks.onSendWelcomeEmail!(user.id),
      disabled: (user: User) => !user.is_active,
    });
  }

  return actions;
}

// Filter definitions with improved styling
export function getUserFilters(t: TFunction) {
  return [
    {
      key: 'is_active',
      label: t('filters.status'),
      type: 'select' as const,
      options: [
        { label: t('filters.allStatus'), value: '__all__' },
        { label: t('filters.active'), value: 'true' },
        { label: t('filters.inactive'), value: 'false' },
      ],
    },
    {
      key: 'is_super_admin',
      label: t('filters.adminType'),
      type: 'select' as const,
      options: [
        { label: t('filters.allTypes'), value: '__all__' },
        { label: t('filters.superAdmin'), value: 'true' },
        { label: t('filters.regularUser'), value: 'false' },
      ],
    },
    {
      key: 'email_verified',
      label: t('filters.emailStatus'),
      type: 'select' as const,
      options: [
        { label: t('filters.allEmailStatus'), value: '__all__' },
        { label: t('filters.verified'), value: 'true' },
        { label: t('filters.unverified'), value: 'false' },
      ],
    },
    {
      key: 'role',
      label: t('filters.role'),
      type: 'select' as const,
      options: [
        { label: t('filters.allRoles'), value: '__all__' },
        // Options will be populated dynamically
      ],
    },
  ];
}

// Quick filter buttons with improved icons
export function getUserQuickFilters(t: TFunction) {
  return [
    {
      key: 'all',
      label: t('filters.allUsers'),
      icon: <UserIcon className="h-4 w-4" />,
      filters: {},
    },
    {
      key: 'active',
      label: t('filters.active'),
      icon: <CheckCircle className="h-4 w-4" />,
      filters: { is_active: 'true' },
    },
    {
      key: 'inactive',
      label: t('filters.inactive'),
      icon: <XCircle className="h-4 w-4" />,
      filters: { is_active: 'false' },
    },
    {
      key: 'super_admins',
      label: t('filters.superAdmins'),
      icon: <Shield className="h-4 w-4" />,
      filters: { is_super_admin: 'true' },
    },
    {
      key: 'verified',
      label: t('filters.verified'),
      icon: <CheckCircle className="h-4 w-4" />,
      filters: { email_verified: 'true' },
    },
  ];
}