import React, { useState, useEffect } from 'react';
import { Head } from '@inertiajs/react';
import { router } from '@inertiajs/react';
import Admin from '@/layouts/Admin';
import RoleManager from '@/components/Roles/RoleManager';
// import { toast } from '@/components/ui/use-toast';
const toast = (options: any) => {
  console.log('Toast:', options);
  // Placeholder for toast functionality
};
import type { 
  PermissionMatrixProps,
  BulkAssignmentRequest,
  PermissionAssignmentRequest,
  SyncPermissionsRequest 
} from '@/types/permissions';

interface PermissionsIndexProps {
  auth: any;
  breadcrumbs: any[];
  matrixData: any;
  permissions: any;
  stats: any;
  subtitle: string;
  title: string;
  user: any;
}

export default function PermissionsIndex(props: PermissionsIndexProps) {
  const [loading, setLoading] = useState(false);

  // Debug: Log all props to see what we're receiving
  console.log('Permissions page props:', props);

  // Handle role creation
  const handleCreateRole = async (roleData: any) => {
    setLoading(true);
    
    try {
      const response = await fetch('/api/roles', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
          'Accept': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(roleData)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to create role');
      }

      const result = await response.json();
      
      toast({
        title: "Success",
        description: result.message || 'Role created successfully',
      });

      // Refresh data
      router.reload({ only: ['matrixData', 'stats'] });

    } catch (error) {
      console.error('Create role error:', error);
      
      toast({
        title: "Error",
        description: error instanceof Error ? error.message : 'Failed to create role',
        variant: "destructive",
      });
      
      throw error;
    } finally {
      setLoading(false);
    }
  };

  // Handle role update
  const handleUpdateRole = async (roleId: number, roleData: any) => {
    setLoading(true);
    
    try {
      const response = await fetch(`/api/roles/${roleId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
          'Accept': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(roleData)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to update role');
      }

      const result = await response.json();
      
      toast({
        title: "Success",
        description: result.message || 'Role updated successfully',
      });

      // Refresh data
      router.reload({ only: ['matrixData', 'stats'] });

    } catch (error) {
      console.error('Update role error:', error);
      
      toast({
        title: "Error",
        description: error instanceof Error ? error.message : 'Failed to update role',
        variant: "destructive",
      });
      
      throw error;
    } finally {
      setLoading(false);
    }
  };

  // Handle role deletion
  const handleDeleteRole = async (roleId: number) => {
    setLoading(true);
    
    try {
      const response = await fetch(`/api/roles/${roleId}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
          'Accept': 'application/json',
        },
        credentials: 'include'
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to delete role');
      }

      const result = await response.json();
      
      toast({
        title: "Success",
        description: result.message || 'Role deleted successfully',
      });

      // Refresh data
      router.reload({ only: ['matrixData', 'stats'] });

    } catch (error) {
      console.error('Delete role error:', error);
      
      toast({
        title: "Error",
        description: error instanceof Error ? error.message : 'Failed to delete role',
        variant: "destructive",
      });
      
      throw error;
    } finally {
      setLoading(false);
    }
  };

  // Handle individual permission toggle
  const handlePermissionToggle = async (roleId: number, serviceSlug: string, action: string, isAssigned: boolean) => {
    setLoading(true);
    
    try {
      const url = isAssigned ? '/api/permissions/assign' : '/api/permissions/revoke';
      const data = {
        role_id: roleId,
        service: serviceSlug,
        action: action
      };

      const response = await fetch(url, {
        method: isAssigned ? 'POST' : 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
          'Accept': 'application/json',
        },
        credentials: 'include', // Include cookies for JWT
        body: JSON.stringify(data)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || `Failed to ${isAssigned ? 'assign' : 'revoke'} permission`);
      }

      const result = await response.json();
      
      // Show success message
      toast({
        title: "Success",
        description: result.message || `Permission ${isAssigned ? 'assigned' : 'revoked'} successfully`,
      });

    } catch (error) {
      console.error('Permission toggle error:', error);
      
      // Show error message
      toast({
        title: "Error",
        description: error instanceof Error ? error.message : `Failed to ${isAssigned ? 'assign' : 'revoke'} permission`,
        variant: "destructive",
      });
      
      // Re-throw to let the component handle the reversion
      throw error;
    } finally {
      setLoading(false);
    }
  };

  // Handle bulk permission assignment
  const handleBulkAssign = async (request: BulkAssignmentRequest) => {
    setLoading(true);
    
    try {
      const response = await fetch('/api/permissions/bulk', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
          'Accept': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(request)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || `Failed to perform bulk ${request.action}`);
      }

      const result = await response.json();
      
      toast({
        title: "Success",
        description: result.message || `Bulk ${request.action} completed successfully`,
      });

    } catch (error) {
      console.error('Bulk assign error:', error);
      
      toast({
        title: "Error",
        description: error instanceof Error ? error.message : `Failed to perform bulk ${request.action}`,
        variant: "destructive",
      });
      
      throw error;
    } finally {
      setLoading(false);
    }
  };

  // Handle role permission sync
  const handleSyncRole = async (roleId: number, permissionSlugs: string[]) => {
    setLoading(true);
    
    try {
      const data = {
        permission_slugs: permissionSlugs
      };

      const response = await fetch(`/api/permissions/sync/${roleId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
          'Accept': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(data)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to sync role permissions');
      }

      const result = await response.json();
      
      toast({
        title: "Success",
        description: result.message || 'Role permissions synchronized successfully',
      });

    } catch (error) {
      console.error('Sync role error:', error);
      
      toast({
        title: "Error",
        description: error instanceof Error ? error.message : 'Failed to sync role permissions',
        variant: "destructive",
      });
      
      throw error;
    } finally {
      setLoading(false);
    }
  };

  // Handle page refresh after major changes
  const refreshData = () => {
    router.reload({ only: ['matrixData', 'stats'] });
  };

  // If no props or missing essential data, show error state
  if (!props || !props.title) {
    return (
      <Admin>
        <Head title="Permissions" />
        <div className="p-6">
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
            <strong>Error:</strong> Failed to load permissions data. Props: {JSON.stringify(props, null, 2)}
          </div>
        </div>
      </Admin>
    );
  }

  return (
    <Admin title="Role Management">
      <Head title="Role Management" />
      
      <div className="space-y-6 min-w-0 overflow-hidden">
        {/* Breadcrumbs */}
        <nav className="flex" aria-label="Breadcrumb">
          <ol className="inline-flex items-center space-x-1 md:space-x-3">
            {(props.breadcrumbs || []).map((breadcrumb, index) => (
              <li key={index} className="inline-flex items-center">
                {index > 0 && (
                  <svg className="w-3 h-3 text-gray-400 mx-1" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 6 10">
                    <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="m1 9 4-4-4-4"/>
                  </svg>
                )}
                {breadcrumb.active ? (
                  <span className="text-sm font-medium text-gray-500 dark:text-gray-400">
                    {breadcrumb.label}
                  </span>
                ) : (
                  <a 
                    href={breadcrumb.href}
                    className="text-sm font-medium text-gray-700 hover:text-blue-600 dark:text-gray-400 dark:hover:text-white"
                  >
                    {breadcrumb.label}
                  </a>
                )}
              </li>
            ))}
          </ol>
        </nav>

        {/* Page Header */}
        <div className="border-b border-gray-200 dark:border-gray-700 pb-4">
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            Role Management
          </h1>
          <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
            Create and manage user roles with granular permissions
          </p>
        </div>

        {/* Role Manager Component */}
        {props.matrixData ? (
          <RoleManager
            initialData={props.matrixData}
            onCreateRole={handleCreateRole}
            onUpdateRole={handleUpdateRole}
            onDeleteRole={handleDeleteRole}
            onPermissionToggle={handlePermissionToggle}
            loading={loading}
          />
        ) : (
          <div className="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded">
            <strong>Warning:</strong> No role data available. Loading...
          </div>
        )}

        {/* Access Control Notice */}
        <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
          <div className="flex">
            <div className="flex-shrink-0">
              <svg className="w-5 h-5 text-blue-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
              </svg>
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-blue-800 dark:text-blue-200">
                Super Admin Access
              </h3>
              <div className="mt-2 text-sm text-blue-700 dark:text-blue-300">
                <p>
                  This permission matrix is only accessible to super administrators. 
                  Changes made here will immediately affect user access across the application.
                  Use caution when modifying role permissions.
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Quick Stats Summary */}
        <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
          <h3 className="text-sm font-medium text-gray-900 dark:text-white mb-3">
            Current System Overview
          </h3>
          <div className="grid grid-cols-2 md:grid-cols-5 gap-4 text-sm">
            <div>
              <span className="text-gray-500 dark:text-gray-400">Total Roles:</span>
              <span className="ml-2 font-semibold text-gray-900 dark:text-white">
                {props.stats?.total_roles || 0}
              </span>
            </div>
            <div>
              <span className="text-gray-500 dark:text-gray-400">Total Permissions:</span>
              <span className="ml-2 font-semibold text-gray-900 dark:text-white">
                {props.stats?.total_permissions || 0}
              </span>
            </div>
            <div>
              <span className="text-gray-500 dark:text-gray-400">Active Roles:</span>
              <span className="ml-2 font-semibold text-green-600 dark:text-green-400">
                {props.stats?.active_roles || 0}
              </span>
            </div>
            <div>
              <span className="text-gray-500 dark:text-gray-400">Active Permissions:</span>
              <span className="ml-2 font-semibold text-green-600 dark:text-green-400">
                {props.stats?.active_permissions || 0}
              </span>
            </div>
            <div>
              <span className="text-gray-500 dark:text-gray-400">Total Assignments:</span>
              <span className="ml-2 font-semibold text-blue-600 dark:text-blue-400">
                {props.stats?.total_assignments || 0}
              </span>
            </div>
          </div>
        </div>
      </div>
    </Admin>
  );
}

// Set the layout for this page
PermissionsIndex.layout = (page: React.ReactElement) => page;