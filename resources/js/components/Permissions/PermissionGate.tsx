import React, { ReactNode } from 'react';
import { usePermissions } from '@/contexts/PermissionsContext';

interface PermissionGateProps {
  children: ReactNode;
  permission?: string;
  service?: string;
  action?: 'create' | 'read' | 'update' | 'delete' | 'export' | 'bulk_update' | 'bulk_delete';
  requireSuperAdmin?: boolean;
  requireAdmin?: boolean;
  require2FA?: boolean; // Explicitly require 2FA for this gate
  fallback?: ReactNode;
  twoFactorFallback?: ReactNode; // Custom fallback when 2FA is required but not enabled
  any?: boolean; // If true, any of the conditions can be true (OR logic)
}

/**
 * PermissionGate component that conditionally renders children based on user permissions
 *
 * When global 2FA requirement is enabled (auth.require_2fa config), this gate will also
 * check if the user has 2FA enabled before granting access to permission-protected content.
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
 *
 * <PermissionGate require2FA twoFactorFallback={<Setup2FAPrompt />}>
 *   <SensitiveData />
 * </PermissionGate>
 */
export function PermissionGate({
  children,
  permission,
  service,
  action,
  requireSuperAdmin = false,
  requireAdmin = false,
  require2FA = false,
  fallback = null,
  twoFactorFallback,
  any = false,
}: PermissionGateProps) {
  const { hasPermission, canPerformAction, isSuperAdmin, isAdmin, needs2FASetup, is2FARequired } = usePermissions();

  // Check if this is a permission-protected gate (requires explicit permission)
  const isPermissionProtected = !!(permission || (service && action) || requireSuperAdmin || requireAdmin);

  // If 2FA is required globally and this is a permission-protected gate,
  // or if require2FA is explicitly set, check 2FA status first
  const shouldCheck2FA = require2FA || (is2FARequired && isPermissionProtected);

  if (shouldCheck2FA && needs2FASetup()) {
    // User needs to set up 2FA - show the 2FA fallback or regular fallback
    return <>{twoFactorFallback ?? fallback}</>;
  }

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
 * Also checks 2FA requirement when enabled globally
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
  const { canPerformAction, needs2FASetup, is2FARequired } = usePermissions();

  // Check 2FA first if required
  if (is2FARequired && needs2FASetup()) {
    return <>{fallback}</>;
  }

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

interface TwoFactorRequiredAlertProps {
  className?: string;
  onSetup?: () => void;
}

/**
 * Alert component to display when 2FA is required but not enabled
 * Can be used as a fallback in PermissionGate
 */
export function TwoFactorRequiredAlert({ className, onSetup }: TwoFactorRequiredAlertProps) {
  const { needs2FASetup } = usePermissions();

  if (!needs2FASetup()) {
    return null;
  }

  return (
    <div className={`rounded-lg border border-amber-500 bg-amber-50 dark:bg-amber-950 p-4 ${className || ''}`}>
      <div className="flex items-start gap-3">
        <svg
          className="h-5 w-5 text-amber-600 mt-0.5 flex-shrink-0"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          />
        </svg>
        <div className="flex-1">
          <h4 className="text-sm font-medium text-amber-800 dark:text-amber-200">
            Two-Factor Authentication Required
          </h4>
          <p className="mt-1 text-sm text-amber-700 dark:text-amber-300">
            You must enable two-factor authentication to access this feature.
            Please set up 2FA in your account security settings.
          </p>
          {onSetup && (
            <button
              onClick={onSetup}
              className="mt-3 inline-flex items-center px-3 py-1.5 text-sm font-medium rounded-md bg-amber-600 text-white hover:bg-amber-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-amber-500"
            >
              Set Up 2FA
            </button>
          )}
        </div>
      </div>
    </div>
  );
}

/**
 * Hook to check if user needs 2FA setup for permission-protected actions
 */
export function useRequires2FASetup(): boolean {
  const { needs2FASetup, is2FARequired } = usePermissions();
  return is2FARequired && needs2FASetup();
}