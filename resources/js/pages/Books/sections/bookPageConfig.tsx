import React from 'react';
import { BookOpen, Upload, Download, BarChart3 } from 'lucide-react';
import { TFunction } from 'i18next';
import {
  StatsCardConfig,
  PageActionConfig,
  SimpleFilterConfig
} from '@/lib/crud-page-utils';

/**
 * Stats card configurations for books
 */
export function getBookStatsConfigs(t: TFunction): StatsCardConfig[] {
  return [
    {
      title: t('stats.totalBooks'),
      getValue: (stats) => stats.totalBooks,
      icon: <BookOpen />,
      getDescription: (stats) =>
        stats.totalValue > 0
          ? `Worth ${new Intl.NumberFormat('en-US', {
              style: 'currency',
              currency: 'USD',
            }).format(stats.totalValue)}`
          : undefined,
    },
    {
      title: t('stats.available'),
      getValue: (stats) => stats.availableBooks,
      icon: <div className="h-4 w-4 bg-green-500 rounded-full" />,
      getDescription: (stats) =>
        stats.totalBooks > 0
          ? t('stats.percentAvailable', { percent: Math.round((stats.availableBooks / stats.totalBooks) * 100) })
          : undefined,
      valueClassName: 'text-green-600',
    },
    {
      title: t('stats.borrowed'),
      getValue: (stats) => stats.borrowedBooks,
      icon: <div className="h-4 w-4 bg-blue-500 rounded-full" />,
      getDescription: (stats) =>
        stats.totalBooks > 0
          ? t('stats.percentBorrowed', { percent: Math.round((stats.borrowedBooks / stats.totalBooks) * 100) })
          : undefined,
      valueClassName: 'text-blue-600',
    },
    {
      title: t('stats.maintenance'),
      getValue: (stats) => stats.maintenanceBooks,
      icon: <div className="h-4 w-4 bg-orange-500 rounded-full" />,
      getDescription: (stats) => t('stats.avgPrice', { price: stats.averagePrice.toFixed(2) }),
      valueClassName: 'text-orange-600',
    },
  ];
}

/**
 * Simple filter configurations for books
 */
export function getBookSimpleFilters(t: TFunction, stats: any): SimpleFilterConfig[] {
  return [
    {
      key: 'status-available',
      label: t('status.available'),
      value: 'AVAILABLE',
      badge: stats?.availableBooks || 0,
      filterParams: { status: 'AVAILABLE' }
    },
    {
      key: 'status-borrowed',
      label: t('status.borrowed'),
      value: 'BORROWED',
      badge: stats?.borrowedBooks || 0,
      filterParams: { status: 'BORROWED' }
    },
    {
      key: 'status-maintenance',
      label: t('status.maintenance'),
      value: 'MAINTENANCE',
      badge: stats?.maintenanceBooks || 0,
      filterParams: { status: 'MAINTENANCE' }
    },
  ];
}

/**
 * Page action configurations for books
 */
export function getBookPageActions(
  t: TFunction,
  permissions: any,
  handlers: {
    onImport: () => void;
    onExport: () => void;
    onReports: () => void;
  }
): PageActionConfig[] {
  return [
    {
      key: 'import',
      label: t('actions.importBooks'),
      icon: <Upload className="h-4 w-4" />,
      handler: handlers.onImport,
      permission: permissions.canManageLibrary,
    },
    {
      key: 'export',
      label: t('actions.exportBooks'),
      icon: <Download className="h-4 w-4" />,
      handler: handlers.onExport,
      permission: permissions.canManageLibrary,
    },
    {
      key: 'reports',
      label: t('actions.viewReports'),
      icon: <BarChart3 className="h-4 w-4" />,
      handler: handlers.onReports,
      permission: permissions.canViewReports,
    },
  ];
}

/**
 * Bulk action configurations for books
 */
export function getBookBulkActions(t: TFunction) {
  return {
    handleBulkDelete: (bookIds: number[]) => {
      const confirmMessage = t('confirm.bulkDelete', { count: bookIds.length });
      if (confirm(confirmMessage)) {
        router.delete('/api/books/bulk', {
          data: { bookIds },
        });
      }
    },

    handleBulkStatusUpdate: (bookIds: number[], status: string) => {
      router.put('/api/books/bulk/status', {
        bookIds,
        status,
      });
    },

    handleBulkExport: (bookIds: number[], filters: any) => {
      const format = prompt(t('confirm.exportFormat')) || 'csv';
      const params = new URLSearchParams({
        format: format,
        bookIds: bookIds.join(','),
        ...Object.fromEntries(
          Object.entries(filters || {}).map(([key, value]) => [key, String(value)])
        ),
      });

      window.open(`/api/books/export?${params.toString()}`);
    },

    handleBulkAddTags: (bookIds: number[], tags: string[]) => {
      router.put('/api/books/bulk/tags', {
        bookIds,
        tags,
        action: 'add',
      });
    },
  };
}

// Import router for bulk actions
import { router } from '@inertiajs/react';