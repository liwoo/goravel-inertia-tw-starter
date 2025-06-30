import React, { useState } from 'react';
import { Head, router } from '@inertiajs/react';
import { Download, Upload, Plus, RefreshCw, Shield, Users } from 'lucide-react';
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
  userColumns, 
  userColumnsMobile, 
  userFilters, 
  userQuickFilters, 
  createUserAdditionalActions
} from './sections';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
// import { useIsMobile } from '@/hooks/use-mobile';
import Admin from '@/layouts/Admin';

export default function UsersIndex({ 
  data, 
  filters, 
  stats,
  roles,
  permissions 
}: UserIndexProps) {
  const isMobile = false; // useIsMobile();
  
  // Debug logging
  console.log('UsersIndex - data:', data);
  console.log('UsersIndex - filters:', filters);
  console.log('UsersIndex - stats:', stats);
  console.log('UsersIndex - roles:', roles);
  console.log('UsersIndex - permissions:', permissions);
  
  // Dialog states
  const [showImportDialog, setShowImportDialog] = useState(false);
  const [showExportDialog, setShowExportDialog] = useState(false);
  const [selectedUsers, setSelectedUsers] = useState<User[]>([]);

  // Update user filters with roles
  const updatedUserFilters = userFilters.map(filter => {
    if (filter.key === 'role' && roles) {
      return {
        ...filter,
        options: [
          { label: 'All', value: '' },
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

  // Additional action handlers (beyond the default View/Edit/Delete)
  const handleActivateUser = async (id: number) => {
    try {
      const response = await fetch(`/api/users/${id}/activate`, {
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

  const handleDeactivateUser = async (id: number) => {
    try {
      const response = await fetch(`/api/users/${id}/deactivate`, {
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

  const handleResetPassword = async (id: number) => {
    if (confirm('Are you sure you want to reset this user\'s password?')) {
      try {
        const response = await fetch(`/api/users/${id}/reset-password`, {
          method: 'POST',
          headers: {
            'Accept': 'application/json',
            'X-Requested-With': 'XMLHttpRequest',
          },
        });
        
        if (response.ok) {
          alert('Password reset email sent successfully');
        }
      } catch (error) {
        console.error('Reset password error:', error);
      }
    }
  };

  const handleImpersonateUser = async (id: number) => {
    if (confirm('Are you sure you want to impersonate this user?')) {
      try {
        const response = await fetch(`/api/users/${id}/impersonate`, {
          method: 'POST',
          headers: {
            'Accept': 'application/json',
            'X-Requested-With': 'XMLHttpRequest',
          },
        });
        
        if (response.ok) {
          // Redirect to dashboard as the impersonated user
          window.location.href = '/admin/dashboard';
        }
      } catch (error) {
        console.error('Impersonate error:', error);
      }
    }
  };

  const handleSendWelcomeEmail = async (id: number) => {
    try {
      const response = await fetch(`/api/users/${id}/send-welcome`, {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
      });
      
      if (response.ok) {
        alert('Welcome email sent successfully');
      }
    } catch (error) {
      console.error('Send welcome email error:', error);
    }
  };

  // Create additional actions (beyond the default View/Edit/Delete that CrudPage provides)
  const additionalActions = createUserAdditionalActions({
    onActivate: permissions.canEdit ? handleActivateUser : undefined,
    onDeactivate: permissions.canEdit ? handleDeactivateUser : undefined,
    onResetPassword: permissions.canEdit ? handleResetPassword : undefined,
    onImpersonate: permissions.canManage ? handleImpersonateUser : undefined,
    onSendWelcomeEmail: permissions.canEdit ? handleSendWelcomeEmail : undefined,
  });

  // Custom form wrappers to include roles
  const CreateFormWithRoles = (props: any) => (
    <UserCreateForm {...props} roles={roles} />
  );

  const EditFormWithRoles = (props: any) => (
    <UserEditForm {...props} roles={roles} />
  );

  return (
    <Admin title="User Management">
      <Head title="Users - Management" />
      
      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        {/* Admin Notice */}
        <div className="px-4 lg:px-6">
          <Card className="bg-gradient-to-t from-primary/5 to-card shadow-xs">
            <CardContent className="p-4">
              <div className="flex items-center gap-3">
                <div className="flex-shrink-0">
                  <Shield className="h-5 w-5 text-blue-600 dark:text-blue-400" />
                </div>
                <div className="min-w-0 flex-1">
                  <p className="text-sm font-medium text-blue-800 dark:text-blue-200">
                    Super Admin Access
                  </p>
                  <p className="text-sm text-blue-700 dark:text-blue-300">
                    This page is only accessible to super administrators.
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Statistics Cards */}
        {stats && (
          <div className="grid grid-cols-1 gap-4 px-4 lg:px-6 md:grid-cols-2 xl:grid-cols-4">
            <Card className="bg-gradient-to-br from-primary/5 to-card shadow-xs">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-base font-medium">Total Users</CardTitle>
                <Users className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.totalUsers}</div>
              </CardContent>
            </Card>

            <Card className="bg-gradient-to-br from-primary/5 to-card shadow-xs">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-base font-medium">Active Users</CardTitle>
                <div className="h-4 w-4 bg-green-500 rounded-full" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-green-600">{stats.activeUsers}</div>
                <p className="text-xs text-muted-foreground">
                  {stats.totalUsers > 0 && `${Math.round((stats.activeUsers / stats.totalUsers) * 100)}% of total`}
                </p>
              </CardContent>
            </Card>

            <Card className="bg-gradient-to-br from-primary/5 to-card shadow-xs">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-base font-medium">Inactive Users</CardTitle>
                <div className="h-4 w-4 bg-gray-500 rounded-full" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-gray-600">{stats.inactiveUsers}</div>
                <p className="text-xs text-muted-foreground">
                  {stats.totalUsers > 0 && `${Math.round((stats.inactiveUsers / stats.totalUsers) * 100)}% of total`}
                </p>
              </CardContent>
            </Card>

            <Card className="bg-gradient-to-br from-primary/5 to-card shadow-xs">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-base font-medium">Super Admins</CardTitle>
                <Shield className="h-4 w-4 text-blue-500" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-blue-600">{stats.superAdmins}</div>
                <p className="text-xs text-muted-foreground">
                  Full system access
                </p>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Quick Filter Buttons */}
        <div className="flex flex-wrap gap-2 px-4 lg:px-6">
          {userQuickFilters.map((filter) => (
            <Button
              key={filter.key}
              variant={JSON.stringify(filters) === JSON.stringify(filter.filters) ? 'default' : 'outline'}
              size="sm"
              onClick={() => {
                router.get('/admin/users', filter.filters, {
                  preserveState: true,
                  preserveScroll: true,
                  only: ['data', 'filters', 'stats'],
                });
              }}
              className="flex items-center space-x-2"
            >
              {filter.icon}
              <span>{filter.label}</span>
            </Button>
          ))}
        </div>

        {/* Management Actions */}
        {permissions.canManage && (
          <div className="flex flex-wrap gap-2 px-4 lg:px-6">
            <Button 
              variant="outline" 
              size="sm"
              onClick={() => setShowImportDialog(true)}
            >
              <Upload className="h-4 w-4 mr-2" />
              Import Users
            </Button>
            <Button 
              variant="outline" 
              size="sm" 
              onClick={() => setShowExportDialog(true)}
            >
              <Download className="h-4 w-4 mr-2" />
              Export Users
            </Button>
            <Button 
              variant="outline" 
              size="sm"
              onClick={handleRefresh}
            >
              <RefreshCw className="h-4 w-4 mr-2" />
              Refresh
            </Button>
          </div>
        )}

        {/* Main CRUD Component */}
        <div className="px-0">
          <CrudPage<User>
          data={data}
          filters={filters}
          title="Users"
          resourceName="users"
          columns={isMobile ? userColumnsMobile : userColumns}
          actions={additionalActions}
          customFilters={updatedUserFilters}
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