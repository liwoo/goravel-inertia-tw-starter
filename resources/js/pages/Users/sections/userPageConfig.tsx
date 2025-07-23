import React from 'react';
import { Users, Shield, Upload, Download, RefreshCw } from 'lucide-react';
import { 
  StatsCardConfig, 
  PageActionConfig, 
  SimpleFilterConfig,
  handleResourceAction 
} from '@/lib/crud-page-utils';
import { CrudAction } from '@/types/crud';

/**
 * Stats card configurations for users
 */
export const userStatsConfigs: StatsCardConfig[] = [
  {
    title: 'Total Users',
    getValue: (stats) => stats.totalUsers,
    icon: <Users />,
  },
  {
    title: 'Active Users',
    getValue: (stats) => stats.activeUsers,
    icon: <div className="h-4 w-4 bg-green-500 rounded-full" />,
    getDescription: (stats) => 
      stats.totalUsers > 0 
        ? `${Math.round((stats.activeUsers / stats.totalUsers) * 100)}% of total` 
        : undefined,
    valueClassName: 'text-green-600',
  },
  {
    title: 'Inactive Users',
    getValue: (stats) => stats.inactiveUsers,
    icon: <div className="h-4 w-4 bg-gray-500 rounded-full" />,
    getDescription: (stats) => 
      stats.totalUsers > 0 
        ? `${Math.round((stats.inactiveUsers / stats.totalUsers) * 100)}% of total` 
        : undefined,
    valueClassName: 'text-gray-600',
  },
  {
    title: 'Super Admins',
    getValue: (stats) => stats.superAdmins,
    icon: <Shield className="h-4 w-4" />,
    iconClassName: 'text-blue-500',
    description: 'Full system access',
    valueClassName: 'text-blue-600',
  },
];

/**
 * Simple filter configurations for users
 */
export const userSimpleFilters = (stats: any): SimpleFilterConfig[] => [
  {
    key: 'active',
    label: 'Active',
    value: 'active',
    badge: stats?.activeUsers || 0,
    filterParams: { is_active: 'true' }
  },
  {
    key: 'inactive',
    label: 'Inactive',
    value: 'inactive',
    badge: stats?.inactiveUsers || 0,
    filterParams: { is_active: 'false' }
  },
  {
    key: 'super_admins',
    label: 'Super Admins',
    value: 'super_admin',
    badge: stats?.superAdmins || 0,
    filterParams: { level_min: '90' }
  },
];

/**
 * Page action configurations for users
 */
export const getUserPageActions = (
  permissions: any,
  handlers: {
    onImport: () => void;
    onExport: () => void;
    onRefresh: () => void;
  }
): PageActionConfig[] => [
  {
    key: 'import',
    label: 'Import Users',
    icon: <Upload className="h-4 w-4" />,
    handler: handlers.onImport,
    permission: permissions.canManage,
  },
  {
    key: 'export',
    label: 'Export Users',
    icon: <Download className="h-4 w-4" />,
    handler: handlers.onExport,
    permission: permissions.canManage,
  },
  {
    key: 'refresh',
    label: 'Refresh',
    icon: <RefreshCw className="h-4 w-4" />,
    handler: handlers.onRefresh,
    permission: permissions.canManage,
  },
];


/*
 * User action handlers using the generic resource handler
 */
export const userActionHandlers = {
  activate: (id: number) => 
    handleResourceAction('users', id, 'activate', {
      successMessage: 'User activated successfully',
    }),
    
  deactivate: (id: number) => 
    handleResourceAction('users', id, 'deactivate', {
      successMessage: 'User deactivated successfully',
    }),
    
  resetPassword: (id: number) => 
    handleResourceAction('users', id, 'reset-password', {
      confirmMessage: 'Are you sure you want to reset this user\'s password?',
      successMessage: 'Password reset email sent successfully',
    }),
    
  impersonate: (id: number) => 
    handleResourceAction('users', id, 'impersonate', {
      confirmMessage: 'Are you sure you want to impersonate this user?',
    }).then(() => {
      window.location.href = '/admin/dashboard';
    }),
    
  sendWelcomeEmail: (id: number) => 
    handleResourceAction('users', id, 'send-welcome', {
      successMessage: 'Welcome email sent successfully',
    }),
};