import React from 'react';
import { PageAction, SimpleFilter } from '@/types/crud';
import { StatsCard } from '@/components/ui/stats-card';

/**
 * Configuration for stats cards
 */
export interface StatsCardConfig {
  title: string;
  getValue: (stats: any) => number | string;
  icon?: React.ReactNode;
  getDescription?: (stats: any) => string | undefined;
  valueClassName?: string;
  iconClassName?: string;
}

/**
 * Renders a grid of stats cards
 */
export function renderStatsCards(stats: any, configs: StatsCardConfig[]) {
  if (!stats) return null;

  return (
    <div className="grid grid-cols-1 gap-4 px-4 lg:px-6 md:grid-cols-2 xl:grid-cols-4">
      {configs.map((config, index) => (
        <StatsCard
          key={index}
          title={config.title}
          value={config.getValue(stats)}
          icon={config.icon}
          description={config.getDescription?.(stats)}
          valueClassName={config.valueClassName}
          iconClassName={config.iconClassName}
        />
      ))}
    </div>
  );
}

/**
 * Common bulk action handlers
 */
export interface BulkActionHandlers {
  onDelete?: (ids: number[]) => void;
  onStatusUpdate?: (ids: number[], status: any) => void;
  onExport?: (ids: number[]) => void;
  [key: string]: ((ids: number[], ...args: any[]) => void) | undefined;
}

export function createBulkActionHandler(
  action: string,
  selectedIds: number[],
  handlers: BulkActionHandlers,
  getSelectedItems: (ids: number[]) => any[]
) {
  if (selectedIds.length === 0) return;

  const handler = handlers[action];
  if (handler) {
    if (action === 'delete' && handlers.onDelete) {
      handlers.onDelete(selectedIds);
    } else {
      handler(selectedIds);
    }
  }
}

/**
 * Common page action configurations
 */
export interface PageActionConfig {
  key: string;
  label: string;
  icon?: React.ReactNode;
  handler: () => void | Promise<void>;
  permission?: boolean;
}

export function createPageActions(configs: PageActionConfig[]): PageAction[] {
  return configs
    .filter(config => config.permission !== false)
    .map(({ permission, ...action }) => action);
}

/**
 * Common simple filter configurations
 */
export interface SimpleFilterConfig {
  key: string;
  label: string;
  value: string | number | boolean;
  badge?: number | string;
  icon?: React.ReactNode;
  filterParams?: Record<string, any>;
}

export function createSimpleFilters(configs: SimpleFilterConfig[]): SimpleFilter[] {
  return configs;
}

/**
 * Generic resource action handlers
 */
export async function handleResourceAction(
  resourceName: string,
  id: number,
  action: string,
  options?: {
    confirmMessage?: string;
    successMessage?: string;
    method?: string;
    body?: any;
  }
) {
  if (options?.confirmMessage && !confirm(options.confirmMessage)) {
    return;
  }

  try {
    const csrfToken = document.querySelector('meta[name="csrf-token"]')?.getAttribute('content');
    
    const response = await fetch(`/api/${resourceName}/${id}/${action}`, {
      method: options?.method || 'POST',
      headers: {
        'Accept': 'application/json',
        'X-Requested-With': 'XMLHttpRequest',
        ...(csrfToken && { 'X-CSRF-TOKEN': csrfToken }),
        ...(options?.body && { 'Content-Type': 'application/json' }),
      },
      ...(options?.body && { body: JSON.stringify(options.body) }),
    });
    
    if (response.ok) {
      if (options?.successMessage) {
        alert(options.successMessage);
      }
      // Trigger refresh
      window.location.reload();
    } else {
      throw new Error(`Action failed with status ${response.status}`);
    }
  } catch (error) {
    console.error(`${action} error:`, error);
    alert(`Failed to ${action} ${resourceName}`);
  }
}

/**
 * Common dialog/modal state management
 */
export interface DialogState<T = any> {
  isOpen: boolean;
  type?: string;
  data?: T;
}

export function useDialogState<T = any>(initialState?: Partial<DialogState<T>>) {
  const [state, setState] = React.useState<DialogState<T>>({
    isOpen: false,
    ...initialState,
  });

  const open = (type?: string, data?: T) => {
    setState({ isOpen: true, type, data });
  };

  const close = () => {
    setState({ isOpen: false, type: undefined, data: undefined });
  };

  return { state, open, close };
}