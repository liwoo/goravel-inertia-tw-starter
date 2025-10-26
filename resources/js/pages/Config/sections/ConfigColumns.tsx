import React from 'react';
import { Config } from '@/types/config';
import { CrudColumn, CrudFilter } from '@/types/crud';
import { Badge } from '@/components/ui/badge';

// Config Type color configuration
const CONFIG_TYPE_COLORS: Record<string, string> = {
  'Financing': 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
  'Improvement Aspects': 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400',
  'Business Categories': 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400',
  'Industries': 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400',
  'Sectors': 'bg-cyan-100 text-cyan-800 dark:bg-cyan-900/30 dark:text-cyan-400',
  'Registration Status': 'bg-pink-100 text-pink-800 dark:bg-pink-900/30 dark:text-pink-400',
  'Development Partners': 'bg-indigo-100 text-indigo-800 dark:bg-indigo-900/30 dark:text-indigo-400',
};

/**
 * Config table columns configuration
 */
export const configColumns: CrudColumn<Config>[] = [
  {
    key: 'name',
    label: 'Name',
    sortable: true,
    className: 'min-w-[250px]',
    render: (config) => (
      <div className="space-y-1">
        <div className="font-medium text-foreground">{config.name}</div>
        <div className="text-xs text-muted-foreground">
          {config.description || <span className="italic opacity-60">No Description</span>}
        </div>
      </div>
    ),
  },
  {
    key: 'configType',
    label: 'Config Type',
    sortable: true,
    className: 'min-w-[180px]',
    render: (config) => {
      const configType = config.configType || config.config_type;
      const colorClass = CONFIG_TYPE_COLORS[configType] || 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400';

      return (
        <Badge className={colorClass}>
          {configType}
        </Badge>
      );
    },
  },
  {
    key: 'createdAt',
    label: 'Created',
    sortable: true,
    className: 'w-32',
    render: (config) => {
      const dateValue = config.createdAt || config.created_at;
      if (!dateValue) return <span className="text-sm text-muted-foreground">-</span>;

      const date = new Date(dateValue);
      return (
        <div className="text-sm text-muted-foreground">
          {date.toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric'
          })}
        </div>
      );
    },
  },
];

/**
 * Compact columns for mobile/smaller screens
 */
export const configColumnsMobile: CrudColumn<Config>[] = [
  {
    key: 'combined',
    label: 'Config',
    sortable: false,
    render: (config) => {
      const configType = config.configType || config.config_type;
      const colorClass = CONFIG_TYPE_COLORS[configType] || 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400';

      return (
        <div className="space-y-2">
          <div className="font-medium text-foreground">{config.name}</div>
          <div className="flex items-center gap-2">
            <Badge className={`${colorClass} text-xs`}>
              {configType}
            </Badge>
          </div>
          {config.description && (
            <div className="text-xs text-muted-foreground line-clamp-1">
              {config.description}
            </div>
          )}
          <div className="text-xs text-muted-foreground">
            {new Date(config.createdAt || config.created_at || '').toLocaleDateString()}
          </div>
        </div>
      );
    },
  },
];

/**
 * Config filters configuration
 */
export const configFilters: CrudFilter[] = [
  {
    key: 'name',
    label: 'Name',
    type: 'text',
    placeholder: 'Enter name',
  },
  {
    key: 'configType',
    label: 'Config Type',
    type: 'text',
    placeholder: 'Enter config type',
  },
  {
    key: 'description',
    label: 'Description',
    type: 'text',
    placeholder: 'Enter description',
  }
];
