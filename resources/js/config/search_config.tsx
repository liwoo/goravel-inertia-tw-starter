import {
  Users,
  FileText,
  Settings,
  Building2,
  Briefcase,
  Calendar,
  FileBox,
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

export type SearchEntityType = 'sme' | 'bdsp' | 'event' | 'procurement' | 'user' | 'config';

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
    type: 'sme',
    label: 'SMEs',
    icon: <Building2 className="h-4 w-4" />,
    permissionService: 'smes',
    permissionAction: 'read',
    colors: {
      light: 'bg-blue-100 text-blue-800',
      dark: 'dark:bg-blue-900/30 dark:text-blue-400',
    },
    urlPrefix: '/admin/smes',
  },
  {
    type: 'bdsp',
    label: 'BDSPs',
    icon: <Briefcase className="h-4 w-4" />,
    permissionService: 'bdsps',
    permissionAction: 'read',
    colors: {
      light: 'bg-green-100 text-green-800',
      dark: 'dark:bg-green-900/30 dark:text-green-400',
    },
    urlPrefix: '/admin/bdsps',
  },
  {
    type: 'event',
    label: 'Events',
    icon: <Calendar className="h-4 w-4" />,
    permissionService: 'events',
    permissionAction: 'read',
    colors: {
      light: 'bg-orange-100 text-orange-800',
      dark: 'dark:bg-orange-900/30 dark:text-orange-400',
    },
    urlPrefix: '/admin/events',
  },
  {
    type: 'procurement',
    label: 'Procurements',
    icon: <FileBox className="h-4 w-4" />,
    permissionService: 'procurement_notices',
    permissionAction: 'read',
    colors: {
      light: 'bg-purple-100 text-purple-800',
      dark: 'dark:bg-purple-900/30 dark:text-purple-400',
    },
    urlPrefix: '/admin/procurements',
  },
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
