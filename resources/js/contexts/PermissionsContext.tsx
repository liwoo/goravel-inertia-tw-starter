import React, { createContext, useContext, ReactNode } from 'react';
import { usePage } from '@inertiajs/react';

// Types for our permission system
export interface UserPermissions {
  id: number;
  name: string;
  email: string;
  role: string;
  roles: Array<{
    id: number;
    name: string;
    slug: string;
    description: string;
    is_active: boolean;
  }>;
  permissions: string[];
  isSuperAdmin: boolean;
  isAdmin: boolean;
}

export interface ServicePermissions {
  canView: boolean;
  canCreate: boolean;
  canEdit: boolean;
  canDelete: boolean;
  canManage: boolean;
  canExport: boolean;
  canBulkUpdate: boolean;
  canBulkDelete: boolean;
  canViewReports: boolean;
  isAdmin: boolean;
  isSuperAdmin: boolean;
}

export interface PermissionsContextType {
  user: UserPermissions | null;
  permissions: Record<string, ServicePermissions>;
  hasPermission: (permission: string) => boolean;
  hasServicePermission: (service: string, action: string) => boolean;
  canPerformAction: (service: string, action: 'create' | 'read' | 'update' | 'delete' | 'export' | 'bulk_update' | 'bulk_delete' | 'write' | 'manage') => boolean;
  isSuperAdmin: () => boolean;
  isAdmin: () => boolean;
}

const PermissionsContext = createContext<PermissionsContextType | undefined>(undefined);

export function PermissionsProvider({ children }: { children: ReactNode }) {
  const { props } = usePage();
  const auth = props.auth as { user: UserPermissions | null; permissions?: Record<string, ServicePermissions> };

  const hasPermission = (permission: string): boolean => {
    if (!auth?.user) return false;
    if (auth.user.isSuperAdmin) return true;
    return auth.user.permissions?.includes(permission) || false;
  };

  const hasServicePermission = (service: string, action: string): boolean => {
    if (!auth?.user) return false;
    if (auth.user.isSuperAdmin) return true;
    
    // Check specific permission in format "service_action"
    const permissionSlug = `${service}_${action}`;
    return auth.user.permissions?.includes(permissionSlug) || false;
  };

  const canPerformAction = (service: string, action: 'create' | 'read' | 'update' | 'delete' | 'export' | 'bulk_update' | 'bulk_delete' | 'write' | 'manage'): boolean => {
    if (!auth?.user) {
      return false;
    }
    if (auth.user.isSuperAdmin) {
      return true;
    }

    // Check the permissions object for this service
    const servicePerms = auth?.permissions?.[service];
    
    if (!servicePerms) {
      // Fallback to checking user's permission array
      const permissionSlug = `${service}_${action}`;
      const hasPermission = auth.user.permissions?.includes(permissionSlug) || false;
      return hasPermission;
    }

    // Map actions to permission keys
    const actionMap: Record<string, keyof ServicePermissions> = {
      create: 'canCreate',
      read: 'canView',
      update: 'canEdit',
      delete: 'canDelete',
      export: 'canExport',
      bulk_update: 'canBulkUpdate',
      bulk_delete: 'canBulkDelete',
      write: 'canEdit', // Map 'write' to 'canEdit'
      manage: 'canManage',
    };

    const permissionKey = actionMap[action];
    const result = permissionKey ? servicePerms[permissionKey] : false;
    
    return result;
  };

  const isSuperAdmin = (): boolean => {
    return auth.user?.isSuperAdmin || false;
  };

  const isAdmin = (): boolean => {
    return auth.user?.isAdmin || false;
  };

  const contextValue: PermissionsContextType = {
    user: auth?.user || null,
    permissions: auth?.permissions || {},
    hasPermission,
    hasServicePermission,
    canPerformAction,
    isSuperAdmin,
    isAdmin,
  };

  return (
    <PermissionsContext.Provider value={contextValue}>
      {children}
    </PermissionsContext.Provider>
  );
}

export function usePermissions(): PermissionsContextType {
  const context = useContext(PermissionsContext);
  if (context === undefined) {
    throw new Error('usePermissions must be used within a PermissionsProvider');
  }
  return context;
}

// Convenience hooks for common permission checks
export function useCanPerform(service: string, action: 'create' | 'read' | 'update' | 'delete' | 'export' | 'bulk_update' | 'bulk_delete' | 'write' | 'manage'): boolean {
  const { canPerformAction } = usePermissions();
  return canPerformAction(service, action);
}

export function useServicePermissions(service: string): ServicePermissions | null {
  const { permissions } = usePermissions();
  return permissions[service] || null;
}