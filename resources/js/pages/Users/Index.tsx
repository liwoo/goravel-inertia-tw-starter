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
  userColumns,
  userColumnsMobile,
  userFilters,
  createUserAdditionalActions,
  userStatsConfigs,
  userSimpleFilters,
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
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';

interface Sme {
  id: number;
  name: string;
  usme_number: string;
}

export default function UsersIndex({ 
  data, 
  filters, 
  stats,
  roles,
  permissions 
}: UserIndexProps) {
  const isMobile = false; // useIsMobile();

  // Dialog states
  const [showImportDialog, setShowImportDialog] = useState(false);
  const [showExportDialog, setShowExportDialog] = useState(false);
  const [selectedUsers, setSelectedUsers] = useState<User[]>([]);

  // Assign to SME dialog states
  const [showAssignSmeDialog, setShowAssignSmeDialog] = useState(false);
  const [selectedUserForSme, setSelectedUserForSme] = useState<User | null>(null);
  const [selectedSmeId, setSelectedSmeId] = useState<string>('');
  const [smes, setSmes] = useState<Sme[]>([]);
  const [isLoadingSmes, setIsLoadingSmes] = useState(false);
  const [isAssigning, setIsAssigning] = useState(false);

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

  // Fetch SMEs for the assign dialog
  const fetchSmes = async () => {
    setIsLoadingSmes(true);
    try {
      const response = await fetch('/api/smes', {
        headers: {
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
      });

      if (response.ok) {
        const data = await response.json();
        setSmes(data.data?.data || []);
      } else {
        toast.error('Failed to load MSMEs');
      }
    } catch (error) {
      console.error('Error fetching MSMEs:', error);
      toast.error('Failed to load MSMEs');
    } finally {
      setIsLoadingSmes(false);
    }
  };

  // Fetch user's currently assigned SME
  const fetchUserAssignedSme = async (userEmail: string) => {
    try {
      const response = await fetch(`/api/primary-business-owners?filter[email]=${encodeURIComponent(userEmail)}&pageSize=1`, {
        headers: {
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
      });

      if (response.ok) {
        const data = await response.json();
        const owners = data.data?.data || [];
        if (owners.length > 0 && owners[0].sme_id) {
          setSelectedSmeId(owners[0].sme_id.toString());
        }
      }
    } catch (error) {
      console.error('Error fetching user assigned MSME:', error);
    }
  };

  // Handle assign to SME click
  const handleAssignToSmeClick = async (user: User) => {
    setSelectedUserForSme(user);
    setSelectedSmeId(''); // Reset first
    setShowAssignSmeDialog(true);
    await fetchSmes();
    // Pre-select user's currently assigned SME if any
    await fetchUserAssignedSme(user.email);
  };

  // Handle assign to SME submission
  const handleAssignToSme = async () => {
    if (!selectedUserForSme || !selectedSmeId) return;

    setIsAssigning(true);
    try {
      const response = await fetch(`/api/users/${selectedUserForSme.id}/assign-sme`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify({ sme_id: parseInt(selectedSmeId) }),
      });

      const data = await response.json();

      if (response.ok) {
        toast.success(data.message || 'User assigned to MSME successfully');
        setShowAssignSmeDialog(false);
        setSelectedUserForSme(null);
        setSelectedSmeId('');
        router.reload({ only: ['data'] });
      } else {
        toast.error(data.message || 'Failed to assign user to MSME');
      }
    } catch (error) {
      console.error('Error assigning user to MSME:', error);
      toast.error('Failed to assign user to MSME');
    } finally {
      setIsAssigning(false);
    }
  };

  // Create additional actions using the extracted handlers
  const additionalActions = createUserAdditionalActions({
    onActivate: permissions.canEdit ? userActionHandlers.activate : undefined,
    onDeactivate: permissions.canEdit ? userActionHandlers.deactivate : undefined,
    onResetPassword: permissions.canEdit ? userActionHandlers.resetPassword : undefined,
    onImpersonate: permissions.canManage ? userActionHandlers.impersonate : undefined,
    onSendWelcomeEmail: permissions.canEdit ? userActionHandlers.sendWelcomeEmail : undefined,
    onAssignToSme: permissions.canEdit ? handleAssignToSmeClick : undefined,
  });

  // Custom form wrappers to include roles
  const CreateFormWithRoles = forwardRef((props: any, ref) => (
    <UserCreateForm {...props} roles={roles} ref={ref} />
  ));

  const EditFormWithRoles = forwardRef((props: any, ref) => (
    <UserEditForm {...props} roles={roles} ref={ref} />
  ));

  // Use extracted configurations
  const simpleFilters = createSimpleFilters(userSimpleFilters(stats));
  
  const pageActions = createPageActions(
    getUserPageActions(permissions, {
      onImport: () => setShowImportDialog(true),
      onExport: () => setShowExportDialog(true),
      onRefresh: handleRefresh,
    })
  );

  return (
    <Admin title="User Management">
      <Head title="Users - Management" />
      
      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        {/* Statistics Cards */}
        {renderStatsCards(stats, userStatsConfigs)}



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

      {/* Assign to MSME Dialog */}
      <Dialog open={showAssignSmeDialog} onOpenChange={setShowAssignSmeDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Assign User to MSME</DialogTitle>
            <DialogDescription>
              Select an MSME to assign {selectedUserForSme?.name} to. This will link the user account to the selected MSME.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="sme">Select MSME *</Label>
              <Select value={selectedSmeId} onValueChange={setSelectedSmeId} disabled={isLoadingSmes}>
                <SelectTrigger>
                  <SelectValue placeholder={isLoadingSmes ? "Loading MSMEs..." : "Select an MSME"} />
                </SelectTrigger>
                <SelectContent>
                  {smes.map((sme) => (
                    <SelectItem key={sme.id} value={sme.id.toString()}>
                      {sme.name} ({sme.usme_number})
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowAssignSmeDialog(false)} disabled={isAssigning}>
              Cancel
            </Button>
            <Button onClick={handleAssignToSme} disabled={isAssigning || !selectedSmeId}>
              {isAssigning ? 'Assigning...' : 'Assign to MSME'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Admin>
  );
}