import React, { ReactNode } from 'react';
import { usePermissions } from '@/contexts/PermissionsContext';

interface PermissionGateProps {
  children: ReactNode;
  permission?: string;
  service?: string;
  action?: 'create' | 'read' | 'update' | 'delete' | 'export' | 'bulk_update' | 'bulk_delete';
  requireSuperAdmin?: boolean;
  requireAdmin?: boolean;
  fallback?: ReactNode;
  any?: boolean; // If true, any of the conditions can be true (OR logic)
}

/**
 * PermissionGate component that conditionally renders children based on user permissions
 * 
 * Examples:
 * <PermissionGate permission="books_create">
 *   <CreateButton />
 * </PermissionGate>
 * 
 * <PermissionGate service="books" action="create">
 *   <CreateButton />
 * </PermissionGate>
 * 
 * <PermissionGate requireSuperAdmin>
 *   <AdminPanel />
 * </PermissionGate>
 */
export function PermissionGate({
  children,
  permission,
  service,
  action,
  requireSuperAdmin = false,
  requireAdmin = false,
  fallback = null,
  any = false,
}: PermissionGateProps) {
  const { hasPermission, canPerformAction, isSuperAdmin, isAdmin } = usePermissions();

  const checks: boolean[] = [];

  // Check specific permission
  if (permission) {
    checks.push(hasPermission(permission));
  }

  // Check service + action combination
  if (service && action) {
    checks.push(canPerformAction(service, action));
  }

  // Check super admin requirement
  if (requireSuperAdmin) {
    checks.push(isSuperAdmin());
  }

  // Check admin requirement
  if (requireAdmin) {
    checks.push(isAdmin());
  }

  // If no checks specified, show children
  if (checks.length === 0) {
    return <>{children}</>;
  }

  // Determine if user has access
  const hasAccess = any ? checks.some(Boolean) : checks.every(Boolean);

  return hasAccess ? <>{children}</> : <>{fallback}</>;
}

interface ConditionalRenderProps {
  condition: boolean;
  children: ReactNode;
  fallback?: ReactNode;
}

/**
 * Simple conditional render component
 */
export function ConditionalRender({ condition, children, fallback = null }: ConditionalRenderProps) {
  return condition ? <>{children}</> : <>{fallback}</>;
}

interface PermissionButtonProps {
  service: string;
  action: 'create' | 'read' | 'update' | 'delete' | 'export' | 'bulk_update' | 'bulk_delete';
  children: ReactNode;
  fallback?: ReactNode;
  className?: string;
  onClick?: () => void;
  disabled?: boolean;
  [key: string]: any; // Allow other button props
}

/**
 * Button component that's automatically disabled/hidden based on permissions
 */
export function PermissionButton({
  service,
  action,
  children,
  fallback = null,
  onClick,
  disabled = false,
  ...props
}: PermissionButtonProps) {
  const { canPerformAction } = usePermissions();
  
  const canPerform = canPerformAction(service, action);
  
  if (!canPerform) {
    return <>{fallback}</>;
  }

  return (
    <button
      onClick={onClick}
      disabled={disabled}
      {...props}
    >
      {children}
    </button>
  );
}