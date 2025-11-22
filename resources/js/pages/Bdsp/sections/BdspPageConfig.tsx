import React from 'react';
import { Download, Upload, BarChart3 } from 'lucide-react';
import {
  StatsCardConfig,
  PageActionConfig,
  SimpleFilterConfig
} from '@/lib/crud-page-utils';
import { BdspStats } from '@/types/bdsp';
import { router } from '@inertiajs/react';

/**
 * Stats card configurations for bdsps
 */
export const bdspStatsConfigs: StatsCardConfig[] = [
  {
    title: 'Total BDSPs',
    getValue: (stats: BdspStats) => stats.totalBdsps,
    icon: <BarChart3 className="h-4 w-4 text-muted-foreground" />,
    getDescription: (stats: BdspStats) => 'Total registered BDSPs',
  },
  {
    title: 'Active',
    getValue: (stats: BdspStats) => stats.activeBdsps,
    icon: <BarChart3 className="h-4 w-4 text-green-600" />,
    getDescription: (stats: BdspStats) => 'Confirmed registrations',
    valueClassName: 'text-green-600',
  },
  {
    title: 'Pending',
    getValue: (stats: BdspStats) => stats.pendingBdsps,
    icon: <BarChart3 className="h-4 w-4 text-yellow-600" />,
    getDescription: (stats: BdspStats) => 'Awaiting approval',
    valueClassName: 'text-yellow-600',
  },
  {
    title: 'Rejected/Suspended',
    getValue: (stats: BdspStats) => (stats.rejectedBdsps || 0) + (stats.suspendedBdsps || 0),
    icon: <BarChart3 className="h-4 w-4 text-red-600" />,
    getDescription: (stats: BdspStats) => 'Rejected or Suspended',
    valueClassName: 'text-red-600',
  },
];

/**
 * Simple filter configurations for bdsps
 */
export const bdspSimpleFilters = (stats: any): SimpleFilterConfig[] => [
  // TODO: Configure your simple filters here
];

/**
 * Page action configurations for bdsps
 */
export const getBdspPageActions = (
  permissions: any,
  handlers: {
    onImport?: () => void;
    onExport: () => void;
  }
): PageActionConfig[] => [
    {
      key: 'export',
      label: 'Export',
      icon: <Download className="h-4 w-4" />,
      handler: handlers.onExport,
    },
  ];

/**
 * Bulk action configurations for bdsps
 */
export const bdspBulkActions = {
  handleBulkDelete: (bdspIds: number[]) => {
    const confirmMessage = `Are you sure you want to delete ${bdspIds.length} item(s)?`;
    if (confirm(confirmMessage)) {
      // TODO: Implement bulk delete
    }
  },
};
