import React from 'react';
// @ts-ignore
import { Head, router } from '@inertiajs/react';
import { Shield } from 'lucide-react';
import { Role } from '@/types/permissions';
import { CrudPage } from '@/components/Crud/CrudPage';
import { Button } from '@/components/ui/button';
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

interface PermissionsIndexProps {
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
  allPermissions?: any[];
  services?: any[];
  actions?: any[];
  title: string;
  subtitle: string;
}

export default function PermissionsIndex({ 
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
}: PermissionsIndexProps) {
  const isMobile = false; // Could use useIsMobile hook if available
  
  const handleRefresh = () => {
    router.reload({ only: ['data', 'filters', 'stats'] });
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

  return (
    <Admin title="Role Management">
      <Head title="Roles - Management" />
      
      <div className="flex flex-col">
        {/* Statistics Cards - Clean design matching Books page */}
        {stats && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 p-6 border-b">
            <div className="bg-card rounded-lg border p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Total Roles</p>
                  <p className="text-3xl font-bold mt-1">{stats.total_roles}</p>
                </div>
                <div className="w-12 h-12 rounded-lg bg-muted flex items-center justify-center">
                  <Shield className="h-6 w-6 text-muted-foreground" />
                </div>
              </div>
            </div>

            <div className="bg-card rounded-lg border p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Active</p>
                  <p className="text-3xl font-bold mt-1 text-green-600">{stats.active_roles}</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    {stats.total_roles > 0 && `${Math.round((stats.active_roles / stats.total_roles) * 100)}% active`}
                  </p>
                </div>
                <div className="w-12 h-12 rounded-full bg-green-500 opacity-20" />
              </div>
            </div>

            <div className="bg-card rounded-lg border p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Inactive</p>
                  <p className="text-3xl font-bold mt-1 text-gray-600">{stats.inactive_roles}</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    {stats.total_roles > 0 && `${Math.round((stats.inactive_roles / stats.total_roles) * 100)}% inactive`}
                  </p>
                </div>
                <div className="w-12 h-12 rounded-full bg-gray-500 opacity-20" />
              </div>
            </div>

            <div className="bg-card rounded-lg border p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-muted-foreground">Users with Roles</p>
                  <p className="text-3xl font-bold mt-1 text-blue-600">{stats.total_users_with_roles}</p>
                  <p className="text-xs text-muted-foreground mt-1">Total assigned</p>
                </div>
                <div className="w-12 h-12 rounded-full bg-blue-500 opacity-20" />
              </div>
            </div>
          </div>
        )}

        {/* Quick Filters - Clean horizontal layout like Books */}
        <div className="flex items-center gap-2 p-6 border-b">
          {roleQuickFilters.map((filter) => (
            <Button
              key={filter.key}
              variant={JSON.stringify(filters) === JSON.stringify(filter.filters) ? 'default' : 'outline'}
              size="sm"
              onClick={() => {
                router.get('/admin/permissions', filter.filters, {
                  preserveState: true,
                  preserveScroll: true,
                  only: ['data', 'filters', 'stats'],
                });
              }}
              className="flex items-center gap-2"
            >
              {filter.icon}
              <span>{filter.label}</span>
            </Button>
          ))}
        </div>

        {/* Main CRUD Component - No extra padding */}
        <CrudPage<Role>
          data={data}
          filters={filters}
          title="My Roles"
          resourceName="roles"
          columns={isMobile ? roleColumnsMobile : roleColumns}
          actions={additionalActions}
          customFilters={roleFilters}
          createForm={CreateFormWithData}
          editForm={EditFormWithData}
          detailView={RoleDetailView}
          onRefresh={handleRefresh}
        />
      </div>
    </Admin>
  );
}