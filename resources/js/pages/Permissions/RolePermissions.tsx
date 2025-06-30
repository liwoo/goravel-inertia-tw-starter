import React, { useState, useEffect } from 'react';
import { Head, router } from '@inertiajs/react';
import { Shield, CheckCircle, XCircle, Save, ArrowLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Switch } from '@/components/ui/switch';
import { toast } from 'sonner';
import Admin from '@/layouts/Admin';

interface Service {
  id: string;
  name: string;
  slug: string;
  actions: Record<string, boolean>;
}

interface Action {
  id: string;
  name: string;
  slug: string;
}

interface Role {
  id: number;
  name: string;
  slug: string;
  description: string;
  level: number;
  is_active: boolean;
  permissions: string[];
}

interface RolePermissionsProps {
  role: Role;
  allPermissions: any[];
  services: Service[];
  actions: Action[];
}

export default function RolePermissions({ 
  role, 
  allPermissions = [], 
  services = [], 
  actions = [] 
}: RolePermissionsProps) {
  const [selectedPermissions, setSelectedPermissions] = useState<Set<string>>(
    new Set(role.permissions || [])
  );
  const [isLoading, setIsLoading] = useState(false);
  const [hasChanges, setHasChanges] = useState(false);

  console.log('RolePermissions component props:', {
    roleId: role.id,
    roleName: role.name,
    rolePermissions: role.permissions,
    servicesCount: services.length,
    actionsCount: actions.length,
    services: services,
    actions: actions
  });

  // Sync selected permissions when role prop changes
  // Using role.id as dependency to ensure we reset when switching between roles
  useEffect(() => {
    console.log('Setting permissions from role data:', role.permissions);
    setSelectedPermissions(new Set(role.permissions || []));
    setHasChanges(false);
  }, [role.id, role.permissions]);

  useEffect(() => {
    const originalPermissions = new Set(role.permissions || []);
    const currentPermissions = selectedPermissions;
    
    // Check if there are any changes
    const hasChanges = originalPermissions.size !== currentPermissions.size ||
      [...originalPermissions].some(perm => !currentPermissions.has(perm)) ||
      [...currentPermissions].some(perm => !originalPermissions.has(perm));
    
    setHasChanges(hasChanges);
  }, [selectedPermissions, role.permissions]);

  const handlePermissionToggle = (serviceSlug: string, actionSlug: string) => {
    const permissionSlug = `${serviceSlug}_${actionSlug}`;
    const newPermissions = new Set(selectedPermissions);
    
    if (newPermissions.has(permissionSlug)) {
      newPermissions.delete(permissionSlug);
    } else {
      newPermissions.add(permissionSlug);
    }
    
    setSelectedPermissions(newPermissions);
  };

  const handleSelectAllForService = (serviceSlug: string, serviceActions: Record<string, boolean>) => {
    const newPermissions = new Set(selectedPermissions);
    const servicePermissions = Object.keys(serviceActions).map(action => `${serviceSlug}_${action}`);
    
    // Check if all permissions for this service are already selected
    const allSelected = servicePermissions.every(perm => newPermissions.has(perm));
    
    if (allSelected) {
      // Remove all permissions for this service
      servicePermissions.forEach(perm => newPermissions.delete(perm));
    } else {
      // Add all permissions for this service
      servicePermissions.forEach(perm => newPermissions.add(perm));
    }
    
    setSelectedPermissions(newPermissions);
  };

  const handleSave = async () => {
    setIsLoading(true);
    
    try {
      const response = await fetch(`/api/roles/${role.id}/permissions`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify({
          permissions: Array.from(selectedPermissions),
        }),
      });

      if (response.ok) {
        const responseData = await response.json();
        toast.success('Permissions updated successfully');
        setHasChanges(false);
        
        // Use Inertia's reload to refresh the page data from the server
        // This will get the updated permissions from the backend
        router.reload({ 
          only: ['role', 'allPermissions', 'services', 'actions'],
          onSuccess: () => {
            // The state will be updated via the useEffect hook when new props arrive
          }
        });
      } else {
        const errorData = await response.json().catch(() => ({}));
        toast.error(errorData.error || errorData.message || 'Failed to update permissions');
      }
    } catch (error) {
      console.error('Save error:', error);
      toast.error('An error occurred while saving');
    } finally {
      setIsLoading(false);
    }
  };

  const getPermissionCount = (serviceSlug: string, serviceActions: Record<string, boolean>) => {
    const servicePermissions = Object.keys(serviceActions);
    const selectedCount = servicePermissions.filter(action => 
      selectedPermissions.has(`${serviceSlug}_${action}`)
    ).length;
    return { selected: selectedCount, total: servicePermissions.length };
  };

  const totalPermissions = services.reduce((total, service) => 
    total + Object.keys(service.actions).length, 0
  );
  const selectedCount = selectedPermissions.size;

  return (
    <Admin title={`${role.name} - Permissions`}>
      <Head title={`${role.name} - Permissions`} />
      
      <div className="flex flex-col gap-6 py-6">
        {/* Header */}
        <div className="px-4 lg:px-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <Button
                variant="outline"
                size="sm"
                onClick={() => window.history.back()}
                className="flex items-center gap-2"
              >
                <ArrowLeft className="h-4 w-4" />
                Back
              </Button>
              <div>
                <h1 className="text-2xl font-bold text-foreground">Role Permissions</h1>
                <div className="text-sm text-muted-foreground">
                  Manage permissions for <Badge variant="secondary">{role.name}</Badge>
                </div>
              </div>
            </div>
            
            <div className="flex items-center gap-3">
              <div className="text-right">
                <p className="text-sm font-medium text-foreground">
                  {selectedCount} of {totalPermissions} permissions
                </p>
                <p className="text-xs text-muted-foreground">
                  {hasChanges ? 'You have unsaved changes' : 'All changes saved'}
                </p>
              </div>
              <Button
                onClick={handleSave}
                disabled={!hasChanges || isLoading}
                className="bg-teal-600 hover:bg-teal-700 text-white"
              >
                <Save className="h-4 w-4 mr-2" />
                {isLoading ? 'Saving...' : 'Save Changes'}
              </Button>
            </div>
          </div>
        </div>

        {/* Role Info Card */}
        <div className="px-4 lg:px-6">
          <Card className="bg-gradient-to-br from-teal-500/10 to-teal-500/5 border-teal-200 dark:border-teal-800">
            <CardContent className="p-6">
              <div className="flex items-center gap-4">
                <div className="p-3 rounded-full bg-teal-100 dark:bg-teal-900/30">
                  <Shield className="h-8 w-8 text-teal-600 dark:text-teal-400" />
                </div>
                <div className="flex-1">
                  <h3 className="text-lg font-semibold text-foreground">{role.name}</h3>
                  <p className="text-sm text-muted-foreground">{role.description || 'No description provided'}</p>
                  <div className="flex items-center gap-4 mt-2">
                    <Badge variant="outline">Level {role.level}</Badge>
                    <Badge variant={role.is_active ? "default" : "secondary"}>
                      {role.is_active ? 'Active' : 'Inactive'}
                    </Badge>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Permissions Matrix */}
        <div className="px-4 lg:px-6">
          <div className="space-y-4">
            {services.map((service) => {
              const permCount = getPermissionCount(service.slug, service.actions);
              const allSelected = permCount.selected === permCount.total;
              const someSelected = permCount.selected > 0 && permCount.selected < permCount.total;

              return (
                <Card key={service.id} className="overflow-hidden">
                  <CardHeader className="pb-3">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <CardTitle className="text-lg">{service.name}</CardTitle>
                        <Badge variant="outline">
                          {permCount.selected}/{permCount.total}
                        </Badge>
                      </div>
                      <div className="flex items-center gap-2">
                        <span className="text-sm text-muted-foreground">
                          {allSelected ? 'All selected' : someSelected ? 'Partial' : 'None selected'}
                        </span>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleSelectAllForService(service.slug, service.actions)}
                          className={`${allSelected ? 'bg-teal-50 border-teal-200 text-teal-700 dark:bg-teal-900/20' : ''}`}
                        >
                          {allSelected ? (
                            <>
                              <XCircle className="h-4 w-4 mr-1" />
                              Deselect All
                            </>
                          ) : (
                            <>
                              <CheckCircle className="h-4 w-4 mr-1" />
                              Select All
                            </>
                          )}
                        </Button>
                      </div>
                    </div>
                  </CardHeader>
                  <Separator />
                  <CardContent className="pt-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                      {Object.keys(service.actions).map((actionSlug) => {
                        const permissionSlug = `${service.slug}_${actionSlug}`;
                        const isSelected = selectedPermissions.has(permissionSlug);
                        const actionName = actions.find(a => a.slug === actionSlug)?.name || actionSlug;
                        
                        console.log(`Permission check: ${permissionSlug} - Selected: ${isSelected}`, {
                          serviceSlug: service.slug,
                          actionSlug: actionSlug,
                          permissionSlug: permissionSlug,
                          isSelected: isSelected,
                          selectedPermissions: Array.from(selectedPermissions)
                        });

                        return (
                          <div
                            key={actionSlug}
                            className={`
                              flex items-center justify-between p-3 rounded-lg border transition-colors
                              ${isSelected 
                                ? 'bg-teal-50 border-teal-200 dark:bg-teal-900/20 dark:border-teal-800' 
                                : 'bg-muted/30 border-border hover:bg-muted/50'
                              }
                            `}
                          >
                            <div className="flex items-center gap-2">
                              <div className={`
                                w-2 h-2 rounded-full 
                                ${isSelected ? 'bg-teal-500' : 'bg-gray-300'}
                              `} />
                              <span className="text-sm font-medium">
                                {actionName}
                              </span>
                            </div>
                            <Switch
                              checked={isSelected}
                              onCheckedChange={() => handlePermissionToggle(service.slug, actionSlug)}
                              size="sm"
                            />
                          </div>
                        );
                      })}
                    </div>
                  </CardContent>
                </Card>
              );
            })}
          </div>
        </div>

        {/* Save Notice */}
        {hasChanges && (
          <div className="px-4 lg:px-6">
            <Card className="bg-amber-50 border-amber-200 dark:bg-amber-900/20 dark:border-amber-800">
              <CardContent className="p-4">
                <div className="flex items-center gap-2">
                  <div className="w-2 h-2 rounded-full bg-amber-500 animate-pulse" />
                  <p className="text-sm text-amber-800 dark:text-amber-200">
                    You have unsaved changes. Don't forget to save your permission updates.
                  </p>
                </div>
              </CardContent>
            </Card>
          </div>
        )}
      </div>
    </Admin>
  );
}