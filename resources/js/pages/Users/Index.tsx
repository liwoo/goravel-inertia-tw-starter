import React, { useState, forwardRef } from 'react';
import { Head, router } from '@inertiajs/react';
import {
  User,
  UserIndexProps,
  UserFormData
} from '@/types/user';
import { CrudPage } from '@/components/Crud/CrudPage';
import {
  UserCreateForm,
  UserEditForm,
  UserDetailView,
  getUserColumns,
  getUserColumnsMobile,
  getUserFilters,
  createUserAdditionalActions,
  getUserStatsConfigs,
  getUserSimpleFilters,
  getUserPageActions,
  userActionHandlers
} from './sections';
import {
  renderStatsCards,
  createPageActions,
  createSimpleFilters
} from '@/lib/crud-page-utils';
// import { useIsMobile } from '@/hooks/use-mobile';
import Admin from '@/layouts/Admin';
import { toast } from 'sonner';
import { useTranslation } from 'react-i18next';
export default function UsersIndex({
  data,
  filters,
  stats,
  roles,
  permissions
}: UserIndexProps) {
  const { t } = useTranslation('users');
  const isMobile = false; // useIsMobile();

  // Dialog states
  const [showImportDialog, setShowImportDialog] = useState(false);
  const [showExportDialog, setShowExportDialog] = useState(false);
  const [selectedUsers, setSelectedUsers] = useState<User[]>([]);

  // Update user filters with roles
  const userFilters = getUserFilters(t);
  const updatedUserFilters = userFilters.map(filter => {
    if (filter.key === 'role' && roles) {
      return {
        ...filter,
        options: [
          { label: t('filters.allRoles'), value: '' },
          ...roles.map(role => ({ label: role.name, value: role.slug }))
        ]
      };
    }
    return filter;
  });

  // Handle bulk operations
  const handleBulkAction = async (action: string, selectedIds: number[]) => {
    if (selectedIds.length === 0) return;

    // Get selected user objects
    const selected = data.data.filter(user => selectedIds.includes(user.id));
    setSelectedUsers(selected);

    const operations: Record<string, () => void> = {
      delete: () => handleBulkDelete(selectedIds),
      activate: () => handleBulkStatusUpdate(selectedIds, true),
      deactivate: () => handleBulkStatusUpdate(selectedIds, false),
      export: () => handleBulkExport(selectedIds),
    };

    const operation = operations[action];
    if (operation) {
      operation();
    }
  };

  const handleBulkDelete = (userIds: number[]) => {
    const confirmMessage = `Are you sure you want to delete ${userIds.length} user(s)? This action cannot be undone.`;
    if (confirm(confirmMessage)) {
      router.delete('/api/users/bulk', {
        data: { userIds },
        onSuccess: () => {
          // Refresh will be handled by the parent
        },
      });
    }
  };

  const handleBulkStatusUpdate = (userIds: number[], isActive: boolean) => {
    router.put('/api/users/bulk/status', {
      userIds,
      is_active: isActive,
    });
  };

  const handleBulkExport = (userIds: number[]) => {
    const format = prompt('Export format (csv, json, excel):') || 'csv';
    const options = {
      format: format as any,
      filters: { ...filters, userIds },
    };
    
    // Trigger download
    window.open(`/api/users/export?${new URLSearchParams(options as any).toString()}`);
  };

  const handleRefresh = () => {
    router.reload({ only: ['data', 'stats'] });
  };

  // Create additional actions using the extracted handlers
  const additionalActions = createUserAdditionalActions(t, {
    onActivate: permissions.canEdit ? userActionHandlers.activate(t) : undefined,
    onDeactivate: permissions.canEdit ? userActionHandlers.deactivate(t) : undefined,
    onResetPassword: permissions.canEdit ? userActionHandlers.resetPassword(t) : undefined,
    onImpersonate: permissions.canManage ? userActionHandlers.impersonate(t) : undefined,
    onSendWelcomeEmail: permissions.canEdit ? userActionHandlers.sendWelcomeEmail(t) : undefined,
  });

  // Custom form wrappers to include roles
  const CreateFormWithRoles = forwardRef((props: any, ref) => (
    <UserCreateForm {...props} roles={roles} ref={ref} />
  ));

  const EditFormWithRoles = forwardRef((props: any, ref) => (
    <UserEditForm {...props} roles={roles} ref={ref} />
  ));

  // Use extracted configurations
  const simpleFilters = createSimpleFilters(getUserSimpleFilters(t, stats));

  const pageActions = createPageActions(
    getUserPageActions(t, permissions, {
      onImport: () => setShowImportDialog(true),
      onExport: () => setShowExportDialog(true),
      onRefresh: handleRefresh,
    })
  );

  return (
    <Admin title={t('page.title')}>
      <Head title={t('page.headTitle')} />

      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        {/* Statistics Cards */}
        {renderStatsCards(stats, getUserStatsConfigs(t))}



        {/* Main CRUD Component */}
        <div className="px-0">
          <CrudPage<User>
          data={data}
          filters={filters}
          title={t('page.crudTitle')}
          resourceName="users"
          columns={isMobile ? getUserColumnsMobile(t) : getUserColumns(t)}
          actions={additionalActions}
          customFilters={updatedUserFilters}
          simpleFilters={simpleFilters}
          pageActions={pageActions}
          createForm={CreateFormWithRoles}
          editForm={EditFormWithRoles}
          detailView={UserDetailView}
          onBulkAction={handleBulkAction}
          onRefresh={handleRefresh}
          canCreate={permissions.canCreate}
          canEdit={permissions.canEdit}
          canDelete={permissions.canDelete}
          canView={true}
          />
        </div>
      </div>

    </Admin>
  );
}