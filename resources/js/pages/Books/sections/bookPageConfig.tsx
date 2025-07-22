import React from 'react';
import { BookOpen, Upload, Download, BarChart3 } from 'lucide-react';
import { 
  StatsCardConfig, 
  PageActionConfig, 
  SimpleFilterConfig 
} from '@/lib/crud-page-utils';

/**
 * Stats card configurations for books
 */
export const bookStatsConfigs: StatsCardConfig[] = [
  {
    title: 'Total Books',
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
    title: 'Available',
    getValue: (stats) => stats.availableBooks,
    icon: <div className="h-4 w-4 bg-green-500 rounded-full" />,
    getDescription: (stats) => 
      stats.totalBooks > 0 
        ? `${Math.round((stats.availableBooks / stats.totalBooks) * 100)}% available` 
        : undefined,
    valueClassName: 'text-green-600',
  },
  {
    title: 'Borrowed',
    getValue: (stats) => stats.borrowedBooks,
    icon: <div className="h-4 w-4 bg-blue-500 rounded-full" />,
    getDescription: (stats) => 
      stats.totalBooks > 0 
        ? `${Math.round((stats.borrowedBooks / stats.totalBooks) * 100)}% borrowed` 
        : undefined,
    valueClassName: 'text-blue-600',
  },
  {
    title: 'Maintenance',
    getValue: (stats) => stats.maintenanceBooks,
    icon: <div className="h-4 w-4 bg-orange-500 rounded-full" />,
    getDescription: (stats) => `Avg. price $${stats.averagePrice.toFixed(2)}`,
    valueClassName: 'text-orange-600',
  },
];

/**
 * Simple filter configurations for books
 */
export const bookSimpleFilters = (stats: any): SimpleFilterConfig[] => [
  {
    key: 'status',
    label: 'Available',
    value: 'AVAILABLE',
    badge: stats?.availableBooks || 0,
    filterParams: { status: 'AVAILABLE' }
  },
  {
    key: 'status',
    label: 'Borrowed',
    value: 'BORROWED',
    badge: stats?.borrowedBooks || 0,
    filterParams: { status: 'BORROWED' }
  },
  {
    key: 'status',
    label: 'Maintenance',
    value: 'MAINTENANCE',
    badge: stats?.maintenanceBooks || 0,
    filterParams: { status: 'MAINTENANCE' }
  },
];

/**
 * Page action configurations for books
 */
export const getBookPageActions = (
  permissions: any,
  handlers: {
    onImport: () => void;
    onExport: () => void;
    onReports: () => void;
  }
): PageActionConfig[] => [
  {
    key: 'import',
    label: 'Import Books',
    icon: <Upload className="h-4 w-4" />,
    handler: handlers.onImport,
    permission: permissions.canManageLibrary,
  },
  {
    key: 'export',
    label: 'Export Books',
    icon: <Download className="h-4 w-4" />,
    handler: handlers.onExport,
    permission: permissions.canManageLibrary,
  },
  {
    key: 'reports',
    label: 'View Reports',
    icon: <BarChart3 className="h-4 w-4" />,
    handler: handlers.onReports,
    permission: permissions.canViewReports,
  },
];

/**
 * Bulk action configurations for books
 */
export const bookBulkActions = {
  handleBulkDelete: (bookIds: number[]) => {
    const confirmMessage = `Are you sure you want to delete ${bookIds.length} book(s)? This action cannot be undone.`;
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
    const format = prompt('Export format (csv, json, excel):') || 'csv';
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

// Import router for bulk actions
import { router } from '@inertiajs/react';