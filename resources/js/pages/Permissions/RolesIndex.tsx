import React from 'react';
import { Head } from '@inertiajs/react';
import { Role } from '@/types/permissions';
import { CrudPage } from '@/components/Crud/CrudPage';
import { renderStatsCards } from '@/lib/crud-page-utils';
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
import { 
  roleStatsConfigs, 
  roleSimpleFilters, 
  rolePageActions,
  roleActionHandlers
} from './sections/rolePageConfig';
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

  // Create additional actions (beyond the default View/Edit/Delete that CrudPage provides)
  const additionalActions = createRoleAdditionalActions({
    onActivate: permissions.canEdit ? roleActionHandlers.handleActivateRole : undefined,
    onDeactivate: permissions.canEdit ? roleActionHandlers.handleDeactivateRole : undefined,
    onDuplicate: permissions.canCreate ? roleActionHandlers.handleDuplicateRole : undefined,
    onAssignUsers: permissions.canManage ? roleActionHandlers.handleAssignUsers : undefined,
    onManagePermissions: permissions.canManage ? roleActionHandlers.handleManagePermissions : undefined,
  });

  // Custom form wrappers to include additional data
  const CreateFormWithData = (props: any) => (
    <RoleCreateForm {...props} allPermissions={allPermissions} services={services} actions={actions} />
  );

  const EditFormWithData = (props: any) => (
    <RoleEditForm {...props} allPermissions={allPermissions} services={services} actions={actions} />
  );

  // Configure filters and actions
  const simpleFilters = roleSimpleFilters(stats);
  const pageActions = rolePageActions(permissions);

  return (
    <Admin title="Role Management">
      <Head title="Roles - Management" />
      
      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        {/* Statistics Cards */}
        {renderStatsCards(stats, roleStatsConfigs)}

        {/* Main CRUD Component */}
        <div className="px-0">
          <CrudPage<Role>
          data={data}
          filters={filters}
          title="My Roles"
          resourceName="roles"
          route="/admin/permissions"
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