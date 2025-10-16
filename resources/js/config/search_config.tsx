import {
  BookOpen,
  Users,
  Shield,
  FileText,
  Landmark,
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

export type SearchEntityType = 'book' | 'user' | 'lender';

export interface SearchEntityConfig {
  type: SearchEntityType;
  label: string;
  icon: React.ReactNode;
  iconClassName?: string;
  permissionService: string; // Service name for permission check
  permissionAction: string;  // Action name for permission check
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
    type: 'book',
    label: 'Books',
    icon: <BookOpen className="h-4 w-4" />,
    permissionService: 'books',
    permissionAction: 'read' as const,
    colors: {
      light: 'bg-blue-100 text-blue-800',
      dark: 'dark:bg-blue-900/30 dark:text-blue-400',
    },
    urlPrefix: '/admin/books',
  },
  {
    type: 'user',
    label: 'Users',
    icon: <Users className="h-4 w-4" />,
    permissionService: 'users',
    permissionAction: 'read' as const,
    colors: {
      light: 'bg-green-100 text-green-800',
      dark: 'dark:bg-green-900/30 dark:text-green-400',
    },
    urlPrefix: '/admin/users',
  },
  {
    type: 'lender',
    label: 'Lenders',
    icon: <Landmark className="h-4 w-4" />,
    permissionService: 'lenders',
    permissionAction: 'read' as const,
    colors: {
      light: 'bg-orange-100 text-orange-800',
      dark: 'dark:bg-orange-900/30 dark:text-orange-400',
    },
    urlPrefix: '/admin/lenders',
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
