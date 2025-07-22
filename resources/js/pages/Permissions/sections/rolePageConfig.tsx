import React from 'react';
import { Shield, CheckCircle, XCircle, Users, Plus, Copy, UserPlus, ShieldCheck } from 'lucide-react';
import { 
  StatsCardConfig, 
  PageActionConfig, 
  SimpleFilterConfig,
  handleResourceAction 
} from '@/lib/crud-page-utils';
import { router } from '@inertiajs/react';

/**
 * Stats card configurations for roles page
 */
export const roleStatsConfigs: StatsCardConfig[] = [
  {
    title: 'Total Roles',
    getValue: (stats) => stats.total_roles,
    icon: <Shield className="h-4 w-4 text-muted-foreground" />,
  },
  {
    title: 'Active',
    getValue: (stats) => stats.active_roles,
    icon: <div className="h-4 w-4 bg-green-500 rounded-full" />,
    getDescription: (stats) => 
      stats.total_roles > 0 
        ? `${Math.round((stats.active_roles / stats.total_roles) * 100)}% active`
        : undefined,
    valueClassName: 'text-green-600',
  },
  {
    title: 'Inactive',
    getValue: (stats) => stats.inactive_roles,
    icon: <div className="h-4 w-4 bg-gray-500 rounded-full" />,
    getDescription: (stats) => 
      stats.total_roles > 0 
        ? `${Math.round((stats.inactive_roles / stats.total_roles) * 100)}% inactive`
        : undefined,
    valueClassName: 'text-gray-600',
  },
  {
    title: 'Users with Roles',
    getValue: (stats) => stats.total_users_with_roles,
    icon: <div className="h-4 w-4 bg-blue-500 rounded-full" />,
    getDescription: () => 'Total assigned',
    valueClassName: 'text-blue-600',
  },
];

/**
 * Simple filter configurations for roles
 */
export const roleSimpleFilters = (stats?: any): SimpleFilterConfig[] => [
  {
    key: 'active',
    label: 'Active',
    value: 'active',
    badge: stats?.active_roles || 0,
    filterParams: { is_active: true },
  },
  {
    key: 'inactive',
    label: 'Inactive',
    value: 'inactive',
    badge: stats?.inactive_roles || 0,
    filterParams: { is_active: false },
  },
  {
    key: 'super_admin',
    label: 'Super Admin',
    value: 'super_admin',
    filterParams: { type: 'super_admin' },
  },
  {
    key: 'admin',
    label: 'Admin',
    value: 'admin',
    filterParams: { type: 'admin' },
  },
  {
    key: 'user',
    label: 'User Roles',
    value: 'user',
    filterParams: { type: 'user' },
  },
];

/**
 * Page action configurations
 */
export const rolePageActions = (permissions: any): PageActionConfig[] => [
  // No page actions currently defined for roles page
  // Add here if needed in future
];

/**
 * Additional action configurations for roles
 */
export interface RoleActionHandlers {
  onActivate?: (id: number) => void;
  onDeactivate?: (id: number) => void;
  onDuplicate?: (id: number) => void;
  onAssignUsers?: (id: number) => void;
  onManagePermissions?: (id: number) => void;
}

export const createRoleAdditionalActions = (handlers: RoleActionHandlers) => [
  {
    key: 'activate',
    label: 'Activate',
    icon: <CheckCircle className="h-4 w-4" />,
    action: (item: any) => handlers.onActivate?.(item.id),
    show: (item: any) => !item.is_active && handlers.onActivate,
  },
  {
    key: 'deactivate',
    label: 'Deactivate',
    icon: <XCircle className="h-4 w-4" />,
    action: (item: any) => handlers.onDeactivate?.(item.id),
    show: (item: any) => item.is_active && handlers.onDeactivate,
  },
  {
    key: 'duplicate',
    label: 'Duplicate',
    icon: <Copy className="h-4 w-4" />,
    action: (item: any) => handlers.onDuplicate?.(item.id),
    show: () => handlers.onDuplicate,
  },
  {
    key: 'assign-users',
    label: 'Assign Users',
    icon: <UserPlus className="h-4 w-4" />,
    action: (item: any) => handlers.onAssignUsers?.(item.id),
    show: () => handlers.onAssignUsers,
  },
  {
    key: 'manage-permissions',
    label: 'Manage Permissions',
    icon: <ShieldCheck className="h-4 w-4" />,
    action: (item: any) => handlers.onManagePermissions?.(item.id),
    show: () => handlers.onManagePermissions,
  },
];

/**
 * Role action handlers
 */
export const roleActionHandlers = {
  handleActivateRole: async (id: number) => {
    await handleResourceAction('roles', id, 'activate');
  },

  handleDeactivateRole: async (id: number) => {
    await handleResourceAction('roles', id, 'deactivate');
  },

  handleDuplicateRole: async (id: number) => {
    await handleResourceAction('roles', id, 'duplicate', {
      confirmMessage: 'Are you sure you want to duplicate this role?',
      successMessage: 'Role duplicated successfully',
    });
  },

  handleAssignUsers: (id: number) => {
    window.location.href = `/admin/roles/${id}/assign-users`;
  },

  handleManagePermissions: (id: number) => {
    router.visit(`/admin/roles/${id}/permissions`);
  },
};

/**
 * Bulk action handlers for roles
 */
export const roleBulkActionHandlers = {
  onDelete: async (ids: number[]) => {
    if (confirm(`Are you sure you want to delete ${ids.length} role(s)?`)) {
      try {
        const response = await fetch('/api/roles/bulk-delete', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'X-Requested-With': 'XMLHttpRequest',
          },
          body: JSON.stringify({ ids }),
        });
        
        if (response.ok) {
          window.location.reload();
        }
      } catch (error) {
        console.error('Bulk delete error:', error);
      }
    }
  },
  
  onStatusUpdate: async (ids: number[], status: 'active' | 'inactive') => {
    try {
      const response = await fetch('/api/roles/bulk-status', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify({ ids, status }),
      });
      
      if (response.ok) {
        window.location.reload();
      }
    } catch (error) {
      console.error('Bulk status update error:', error);
    }
  },
};