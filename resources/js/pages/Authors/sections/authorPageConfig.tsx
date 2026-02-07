import React from 'react';
import { User } from 'lucide-react';
import { TFunction } from 'i18next';
import {
  StatsCardConfig,
  SimpleFilterConfig
} from '@/lib/crud-page-utils';

/**
 * Stats card configurations for authors
 */
export function getAuthorStatsConfigs(t: TFunction): StatsCardConfig[] {
  return [
    {
      title: t('stats.totalAuthors'),
      getValue: (stats) => stats.totalAuthors,
      icon: <User />,
      getDescription: () => undefined,
    },
    {
      title: t('stats.active'),
      getValue: (stats) => stats.activeAuthors,
      icon: <div className="h-4 w-4 bg-green-500 rounded-full" />,
      getDescription: (stats) =>
        stats.totalAuthors > 0
          ? t('stats.percentActive', { percent: Math.round((stats.activeAuthors / stats.totalAuthors) * 100) })
          : undefined,
      valueClassName: 'text-green-600',
    },
    {
      title: t('stats.inactive'),
      getValue: (stats) => stats.inactiveAuthors,
      icon: <div className="h-4 w-4 bg-gray-500 rounded-full" />,
      getDescription: (stats) =>
        stats.totalAuthors > 0
          ? t('stats.percentInactive', { percent: Math.round((stats.inactiveAuthors / stats.totalAuthors) * 100) })
          : undefined,
      valueClassName: 'text-gray-600',
    },
  ];
}

/**
 * Simple filter configurations for authors
 */
export function getAuthorSimpleFilters(t: TFunction, stats: any): SimpleFilterConfig[] {
  return [
    {
      key: 'status-active',
      label: t('status.active'),
      value: 'ACTIVE',
      badge: stats?.activeAuthors || 0,
      filterParams: { status: 'ACTIVE' }
    },
    {
      key: 'status-inactive',
      label: t('status.inactive'),
      value: 'INACTIVE',
      badge: stats?.inactiveAuthors || 0,
      filterParams: { status: 'INACTIVE' }
    },
  ];
}
