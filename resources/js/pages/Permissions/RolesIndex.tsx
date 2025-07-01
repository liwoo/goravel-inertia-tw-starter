import React from 'react';
import { Head, router } from '@inertiajs/react';
import { Shield } from 'lucide-react';
import { Role } from '@/types/permissions';
import { CrudPage } from '@/components/Crud/CrudPage';
import { PageAction, SimpleFilter } from '@/types/crud';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { 
  RoleCreateForm, 
  RoleEditForm, 
  RoleDetailView,
  roleColumns, 
  roleColumnsMobile, 
  roleFilters,
  roleQuickFilters,
  createRoleAdditionalActions
} from './sections';
import Admin from '@/layouts/Admin';

interface RoleListResponse {
  data: Role[];
  total: number;
  perPage: number;
  currentPage: number;
  lastPage: number;
  from: number | null;
  to: number | null;
}

interface RolesIndexProps {
  auth: any;
  data: RoleListResponse;
  filters: any;
  stats?: {
    total_roles: number;
    active_roles: number;
    inactive_roles: number;
    total_users_with_roles: number;
  };
  permissions?: {
    canCreate: boolean;
    canEdit: boolean;
    canDelete: boolean;
    canManage: boolean;
  };
}

export default function RolesIndex({ 
  data, 
  filters = {}, 
  stats,
  permissions = {
    canCreate: true,
    canEdit: true,
    canDelete: true,
    canManage: true,
  },
  allPermissions = [],
  services = [],
  actions = []
}: RolesIndexProps & {
  allPermissions?: any[];
  services?: any[];
  actions?: any[];
}) {
  const isMobile = false; // Could use useIsMobile hook if available
  
  const handleRefresh = () => {
    // Refresh logic handled by CrudPage
  };

  // Additional action handlers (beyond the default View/Edit/Delete)
  const handleActivateRole = async (id: number) => {
    try {
      const response = await fetch(`/api/roles/${id}/activate`, {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
      });
      
      if (response.ok) {
        handleRefresh();
      }
    } catch (error) {
      console.error('Activate error:', error);
    }
  };

  const handleDeactivateRole = async (id: number) => {
    try {
      const response = await fetch(`/api/roles/${id}/deactivate`, {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
      });
      
      if (response.ok) {
        handleRefresh();
      }
    } catch (error) {
      console.error('Deactivate error:', error);
    }
  };

  const handleDuplicateRole = async (id: number) => {
    if (confirm('Are you sure you want to duplicate this role?')) {
      try {
        const response = await fetch(`/api/roles/${id}/duplicate`, {
          method: 'POST',
          headers: {
            'Accept': 'application/json',
            'X-Requested-With': 'XMLHttpRequest',
          },
        });
        
        if (response.ok) {
          handleRefresh();
          alert('Role duplicated successfully');
        }
      } catch (error) {
        console.error('Duplicate error:', error);
      }
    }
  };

  const handleAssignUsers = (id: number) => {
    // Navigate to user assignment page
    window.location.href = `/admin/roles/${id}/assign-users`;
  };

  const handleManagePermissions = (id: number) => {
    // Navigate to permissions management page using Inertia
    router.visit(`/admin/roles/${id}/permissions`);
  };

  // Create additional actions (beyond the default View/Edit/Delete that CrudPage provides)
  const additionalActions = createRoleAdditionalActions({
    onActivate: permissions.canEdit ? handleActivateRole : undefined,
    onDeactivate: permissions.canEdit ? handleDeactivateRole : undefined,
    onDuplicate: permissions.canCreate ? handleDuplicateRole : undefined,
    onAssignUsers: permissions.canManage ? handleAssignUsers : undefined,
    onManagePermissions: permissions.canManage ? handleManagePermissions : undefined,
  });

  // Custom form wrappers to include additional data
  const CreateFormWithData = (props: any) => (
    <RoleCreateForm {...props} allPermissions={allPermissions} services={services} actions={actions} />
  );

  const EditFormWithData = (props: any) => (
    <RoleEditForm {...props} allPermissions={allPermissions} services={services} actions={actions} />
  );

  // Convert quick filters to SimpleFilter format
  const simpleFilters: SimpleFilter[] = [
    {
      key: 'active',
      label: 'Active',
      value: 'active',
      badge: stats?.active_roles || 0,
    },
    {
      key: 'inactive',
      label: 'Inactive',
      value: 'inactive',
      badge: stats?.inactive_roles || 0,
    },
    {
      key: 'super_admin',
      label: 'Super Admin',
      value: 'super_admin',
    },
    {
      key: 'admin',
      label: 'Admin',
      value: 'admin',
    },
    {
      key: 'user',
      label: 'User Roles',
      value: 'user',
    },
  ];

  // No page actions for permissions page in current implementation
  const pageActions: PageAction[] = [];

  return (
    <Admin title="Role Management">
      <Head title="Roles - Management" />
      
      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        {/* Statistics Cards - matching Books page */}
        {stats && (
          <div className="grid grid-cols-1 gap-4 px-4 lg:px-6 md:grid-cols-2 xl:grid-cols-4">
            <Card className="bg-gradient-to-br from-primary/5 to-card shadow-xs">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-base font-medium">Total Roles</CardTitle>
                <Shield className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.total_roles}</div>
              </CardContent>
            </Card>

            <Card className="bg-gradient-to-br from-primary/5 to-card shadow-xs">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-base font-medium">Active</CardTitle>
                <div className="h-4 w-4 bg-green-500 rounded-full" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-green-600">{stats.active_roles}</div>
                <p className="text-xs text-muted-foreground">
                  {stats.total_roles > 0 && `${Math.round((stats.active_roles / stats.total_roles) * 100)}% active`}
                </p>
              </CardContent>
            </Card>

            <Card className="bg-gradient-to-br from-primary/5 to-card shadow-xs">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-base font-medium">Inactive</CardTitle>
                <div className="h-4 w-4 bg-gray-500 rounded-full" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-gray-600">{stats.inactive_roles}</div>
                <p className="text-xs text-muted-foreground">
                  {stats.total_roles > 0 && `${Math.round((stats.inactive_roles / stats.total_roles) * 100)}% inactive`}
                </p>
              </CardContent>
            </Card>

            <Card className="bg-gradient-to-br from-primary/5 to-card shadow-xs">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-base font-medium">Users with Roles</CardTitle>
                <div className="h-4 w-4 bg-blue-500 rounded-full" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-blue-600">{stats.total_users_with_roles}</div>
                <p className="text-xs text-muted-foreground">
                  Total assigned
                </p>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Main CRUD Component */}
        <div className="px-0">
          <CrudPage<Role>
          data={data}
          filters={filters}
          title="My Roles"
          resourceName="roles"
          columns={isMobile ? roleColumnsMobile : roleColumns}
          actions={additionalActions}
          customFilters={roleFilters}
          simpleFilters={simpleFilters}
          pageActions={pageActions}
          createForm={CreateFormWithData}
          editForm={EditFormWithData}
          detailView={RoleDetailView}
          onRefresh={handleRefresh}
        />
        </div>
      </div>
    </Admin>
  );
}