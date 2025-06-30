import React from 'react';
import { Role } from '@/types/permissions';
import { CrudColumn, CrudFilter, CrudAction } from '@/types/crud';
import { Badge } from '@/components/ui/badge';
import { Shield, Users, Calendar, Check, X, Crown, UserCheck, Copy, UserPlus, Settings } from 'lucide-react';

/**
 * Role table columns configuration
 */
export const roleColumns: CrudColumn<Role>[] = [
  {
    key: 'name',
    label: 'Role Details',
    sortable: true,
    className: 'min-w-[300px]',
    render: (role) => (
      <div className="flex items-start gap-3">
        <div className="p-2 rounded-lg bg-muted">
          {role.level >= 90 ? (
            <Crown className="h-5 w-5 text-yellow-500" />
          ) : (
            <Shield className="h-5 w-5 text-muted-foreground" />
          )}
        </div>
        <div className="space-y-1">
          <div className="font-medium text-foreground flex items-center gap-2">
            {role.name}
            {role.level >= 90 && (
              <Badge variant="outline" className="text-xs bg-yellow-50 text-yellow-700 border-yellow-200">
                Super Admin
              </Badge>
            )}
          </div>
          <div className="text-sm text-muted-foreground">{role.description}</div>
          <div className="text-xs text-muted-foreground">
            Slug: {role.slug}
          </div>
        </div>
      </div>
    ),
  },
  {
    key: 'level',
    label: 'Level',
    sortable: true,
    className: 'w-20 text-center',
    render: (role) => (
      <div className="text-center">
        <Badge 
          variant={role.level >= 90 ? 'default' : 'secondary'}
          className={role.level >= 90 
            ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400'
            : ''
          }
        >
          {role.level}
        </Badge>
      </div>
    ),
  },
  {
    key: 'is_active',
    label: 'Status',
    sortable: true,
    className: 'w-24',
    render: (role) => (
      <Badge 
        className={role.is_active 
          ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400 flex items-center gap-1'
          : 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400 flex items-center gap-1'
        }
      >
        {role.is_active ? <Check className="h-3 w-3" /> : <X className="h-3 w-3" />}
        {role.is_active ? 'Active' : 'Inactive'}
      </Badge>
    ),
  },
  {
    key: 'parent_id',
    label: 'Parent Role',
    className: 'w-32',
    render: (role) => (
      <span className="text-sm text-muted-foreground">
        {role.parent ? role.parent.name : '-'}
      </span>
    ),
  },
  {
    key: 'created_at',
    label: 'Created',
    sortable: true,
    className: 'w-28',
    render: (role) => (
      <div className="text-sm text-muted-foreground flex items-center">
        <Calendar className="w-3 h-3 mr-1" />
        {new Date(role.created_at).toLocaleDateString('en-US', {
          month: 'short',
          day: 'numeric',
          year: 'numeric'
        })}
      </div>
    ),
  },
];

/**
 * Compact role columns for mobile/smaller screens
 */
export const roleColumnsMobile: CrudColumn<Role>[] = [
  {
    key: 'name',
    label: 'Role',
    sortable: true,
    render: (role) => (
      <div className="space-y-3">
        <div className="flex items-start gap-3">
          <div className="p-2 rounded-lg bg-muted">
            {role.level >= 90 ? (
              <Crown className="h-5 w-5 text-yellow-500" />
            ) : (
              <Shield className="h-5 w-5 text-muted-foreground" />
            )}
          </div>
          <div className="flex-1 space-y-1">
            <div className="font-medium text-foreground flex items-center gap-2">
              {role.name}
              {role.level >= 90 && (
                <Badge variant="outline" className="text-xs bg-yellow-50 text-yellow-700 border-yellow-200">
                  Super Admin
                </Badge>
              )}
            </div>
            <div className="text-sm text-muted-foreground">{role.description}</div>
          </div>
        </div>
        <div className="flex items-center justify-between pl-12">
          <div className="flex items-center gap-2">
            <Badge 
              variant={role.level >= 90 ? 'default' : 'secondary'}
              className={`text-xs ${role.level >= 90 
                ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400'
                : ''
              }`}
            >
              Level {role.level}
            </Badge>
            {role.parent && (
              <Badge variant="outline" className="text-xs">
                Child of {role.parent.name}
              </Badge>
            )}
          </div>
          <div>
            <Badge 
              className={role.is_active 
                ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400 flex items-center gap-1 text-xs'
                : 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400 flex items-center gap-1 text-xs'
              }
            >
              {role.is_active ? <Check className="h-3 w-3" /> : <X className="h-3 w-3" />}
              {role.is_active ? 'Active' : 'Inactive'}
            </Badge>
          </div>
        </div>
      </div>
    ),
  },
];

/**
 * Role filters configuration
 */
export const roleFilters: CrudFilter[] = [
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
    key: 'level_min',
    label: 'Min Level',
    type: 'number',
    placeholder: '0',
  },
  {
    key: 'level_max',
    label: 'Max Level',
    type: 'number',
    placeholder: '100',
  },
  {
    key: 'has_parent',
    label: 'Has Parent Role',
    type: 'select',
    options: [
      { value: '', label: 'All' },
      { value: 'true', label: 'Yes' },
      { value: 'false', label: 'No' },
    ],
  },
];

/**
 * Quick filter buttons for common role queries
 */
export const roleQuickFilters = [
  {
    key: 'all',
    label: 'All Roles',
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
    key: 'super_admin',
    label: 'Super Admin',
    icon: <Crown className="h-4 w-4 text-yellow-500" />,
    filters: { level_min: '90' },
  },
  {
    key: 'admin',
    label: 'Admin',
    icon: <UserCheck className="h-4 w-4 text-blue-500" />,
    filters: { level_min: '50', level_max: '89' },
  },
  {
    key: 'user',
    label: 'User Roles',
    icon: <Users className="h-4 w-4 text-purple-500" />,
    filters: { level_max: '49' },
  },
] as const;

// Additional actions factory for Role-specific actions (beyond the default View/Edit/Delete)
export const createRoleAdditionalActions = (callbacks: {
  onActivate?: (id: number) => void;
  onDeactivate?: (id: number) => void;
  onDuplicate?: (id: number) => void;
  onAssignUsers?: (id: number) => void;
  onManagePermissions?: (id: number) => void;
}): CrudAction<Role>[] => {
  const actions: CrudAction<Role>[] = [];

  // Status toggle actions
  if (callbacks.onActivate) {
    actions.push({
      key: 'activate',
      label: 'Activate',
      icon: <Check className="h-4 w-4 text-green-600" />,
      onClick: (role: Role) => callbacks.onActivate!(role.id),
      disabled: (role: Role) => role.is_active, // Disable if already active
    });
  }

  if (callbacks.onDeactivate) {
    actions.push({
      key: 'deactivate',
      label: 'Deactivate',
      icon: <X className="h-4 w-4 text-orange-600" />,
      onClick: (role: Role) => callbacks.onDeactivate!(role.id),
      disabled: (role: Role) => !role.is_active, // Disable if already inactive
    });
  }

  // Duplicate role action
  if (callbacks.onDuplicate) {
    actions.push({
      key: 'duplicate',
      label: 'Duplicate Role',
      icon: <Copy className="h-4 w-4 text-blue-600" />,
      onClick: (role: Role) => callbacks.onDuplicate!(role.id),
    });
  }

  // Assign users action
  if (callbacks.onAssignUsers) {
    actions.push({
      key: 'assign-users',
      label: 'Assign Users',
      icon: <UserPlus className="h-4 w-4 text-purple-600" />,
      onClick: (role: Role) => callbacks.onAssignUsers!(role.id),
    });
  }

  // Manage permissions action
  if (callbacks.onManagePermissions) {
    actions.push({
      key: 'manage-permissions',
      label: 'Manage Permissions',
      icon: <Settings className="h-4 w-4 text-indigo-600" />,
      onClick: (role: Role) => callbacks.onManagePermissions!(role.id),
    });
  }

  return actions;
};