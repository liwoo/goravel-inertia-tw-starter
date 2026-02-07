import React from 'react';
import { Users, Shield, Upload, Download, RefreshCw } from 'lucide-react';
import {
  StatsCardConfig,
  PageActionConfig,
  SimpleFilterConfig,
  handleResourceAction
} from '@/lib/crud-page-utils';
import { CrudAction } from '@/types/crud';
import { TFunction } from 'i18next';

/**
 * Stats card configurations for users
 */
export function getUserStatsConfigs(t: TFunction): StatsCardConfig[] {
  return [
    {
      title: t('stats.totalUsers'),
      getValue: (stats) => stats.totalUsers,
      icon: <Users />,
    },
    {
      title: t('stats.activeUsers'),
      getValue: (stats) => stats.activeUsers,
      icon: <div className="h-4 w-4 bg-green-500 rounded-full" />,
      getDescription: (stats) =>
        stats.totalUsers > 0
          ? `${Math.round((stats.activeUsers / stats.totalUsers) * 100)}% of total`
          : undefined,
      valueClassName: 'text-green-600',
    },
    {
      title: t('stats.inactiveUsers'),
      getValue: (stats) => stats.inactiveUsers,
      icon: <div className="h-4 w-4 bg-gray-500 rounded-full" />,
      getDescription: (stats) =>
        stats.totalUsers > 0
          ? `${Math.round((stats.inactiveUsers / stats.totalUsers) * 100)}% of total`
          : undefined,
      valueClassName: 'text-gray-600',
    },
    {
      title: t('stats.superAdmins'),
      getValue: (stats) => stats.superAdmins,
      icon: <Shield className="h-4 w-4" />,
      iconClassName: 'text-blue-500',
      getDescription: () => t('stats.fullSystemAccess'),
      valueClassName: 'text-blue-600',
    },
  ];
}

/**
 * Simple filter configurations for users
 */
export function getUserSimpleFilters(t: TFunction, stats: any): SimpleFilterConfig[] {
  return [
    {
      key: 'active',
      label: t('filters.active'),
      value: 'active',
      badge: stats?.activeUsers || 0,
      filterParams: { is_active: 'true' }
    },
    {
      key: 'inactive',
      label: t('filters.inactive'),
      value: 'inactive',
      badge: stats?.inactiveUsers || 0,
      filterParams: { is_active: 'false' }
    },
    {
      key: 'super_admins',
      label: t('filters.superAdmins'),
      value: 'super_admin',
      badge: stats?.superAdmins || 0,
      filterParams: { level_min: '90' }
    },
  ];
}

/**
 * Page action configurations for users
 */
export function getUserPageActions(
  t: TFunction,
  permissions: any,
  handlers: {
    onImport: () => void;
    onExport: () => void;
    onRefresh: () => void;
  }
): PageActionConfig[] {
  return [
    {
      key: 'import',
      label: t('actions.importUsers'),
      icon: <Upload className="h-4 w-4" />,
      handler: handlers.onImport,
      permission: permissions.canManage,
    },
    {
      key: 'export',
      label: t('actions.exportUsers'),
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
}


/*
 * User action handlers using the generic resource handler
 */
export const userActionHandlers = {
  activate: (t: TFunction) => (id: number) =>
    handleResourceAction('users', id, 'activate', {
      successMessage: t('toast.activated'),
    }),

  deactivate: (t: TFunction) => (id: number) =>
    handleResourceAction('users', id, 'deactivate', {
      successMessage: t('toast.deactivated'),
    }),

  resetPassword: (t: TFunction) => (id: number) =>
    handleResourceAction('users', id, 'reset-password', {
      confirmMessage: t('confirm.resetPassword'),
      successMessage: t('toast.passwordResetSent'),
    }),

  impersonate: (t: TFunction) => (id: number) =>
    handleResourceAction('users', id, 'impersonate', {
      confirmMessage: t('confirm.impersonate'),
    }).then(() => {
      window.location.href = '/admin/dashboard';
    }),

  sendWelcomeEmail: (t: TFunction) => (id: number) =>
    handleResourceAction('users', id, 'send-welcome', {
      successMessage: t('toast.welcomeEmailSent'),
    }),
};