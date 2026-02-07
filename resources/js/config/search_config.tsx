import {
  Users,
  FileText,
  Settings,
} from 'lucide-react';

/**
 * Search Entity Configuration
 *
 * This file centralizes all searchable entity definitions for the global search (CMD+K).
 * To add a new entity to search:
 * 1. Add the entity type to SearchEntityType
 * 2. Add a configuration entry to SEARCH_ENTITIES
 * 3. Ensure the backend search controller handles the entity type
 */

export type SearchEntityType = 'user' | 'config' | 'application';

export type PermissionAction = 'create' | 'read' | 'update' | 'delete' | 'export' | 'bulk_update' | 'bulk_delete' | 'write' | 'manage';

export interface SearchEntityConfig {
  type: SearchEntityType;
  label: string;
  icon: React.ReactNode;
  iconClassName?: string;
  permissionService: string; // Service name for permission check
  permissionAction: PermissionAction;  // Action name for permission check
  colors: {
    light: string;
    dark: string;
  };
  urlPrefix: string; // URL prefix for navigation (e.g., '/admin/books')
}

/**
 * Centralized search entity configurations
 * Add new searchable entities here
 */
export const SEARCH_ENTITIES: SearchEntityConfig[] = [
  {
    type: 'user',
    label: 'Users',
    icon: <Users className="h-4 w-4" />,
    permissionService: 'users',
    permissionAction: 'read',
    colors: {
      light: 'bg-cyan-100 text-cyan-800',
      dark: 'dark:bg-cyan-900/30 dark:text-cyan-400',
    },
    urlPrefix: '/admin/users',
  },
  {
    type: 'config',
    label: 'Configs',
    icon: <Settings className="h-4 w-4" />,
    permissionService: 'config',
    permissionAction: 'read',
    colors: {
      light: 'bg-gray-100 text-gray-800',
      dark: 'dark:bg-gray-900/30 dark:text-gray-400',
    },
    urlPrefix: '/admin/configs',
  },
  {
    type: 'application',
    label: 'Applications',
    icon: <FileText className="h-4 w-4" />,
    permissionService: 'applications',
    permissionAction: 'read',
    colors: {
      light: 'bg-amber-100 text-amber-800',
      dark: 'dark:bg-amber-900/30 dark:text-amber-400',
    },
    urlPrefix: '/admin/applications',
  },
];

/**
 * Default/fallback configuration for unknown entity types
 */
export const DEFAULT_ENTITY_CONFIG: Omit<SearchEntityConfig, 'type' | 'label' | 'permissionService' | 'permissionAction' | 'urlPrefix'> = {
  icon: <FileText className="h-4 w-4" />,
  colors: {
    light: 'bg-gray-100 text-gray-800',
    dark: 'dark:bg-gray-900/30 dark:text-gray-400',
  },
};

/**
 * Helper function to get entity config by type
 */
export function getEntityConfig(type: SearchEntityType): SearchEntityConfig | undefined {
  return SEARCH_ENTITIES.find(entity => entity.type === type);
}

/**
 * Helper function to get entity colors
 */
export function getEntityColors(type: SearchEntityType): string {
  const config = getEntityConfig(type);
  if (!config) {
    return `${DEFAULT_ENTITY_CONFIG.colors.light} ${DEFAULT_ENTITY_CONFIG.colors.dark}`;
  }
  return `${config.colors.light} ${config.colors.dark}`;
}

/**
 * Helper function to get entity icon
 */
export function getEntityIcon(type: SearchEntityType): React.ReactNode {
  const config = getEntityConfig(type);
  return config?.icon ?? DEFAULT_ENTITY_CONFIG.icon;
}
