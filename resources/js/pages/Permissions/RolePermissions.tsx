import React, { useState, useEffect } from 'react';
import { Head, router } from '@inertiajs/react';
import { Shield, CheckCircle, XCircle, Save, ArrowLeft, Settings2, Eye, Edit } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { toast } from 'sonner';
import Admin from '@/layouts/Admin';
import { PermissionScopeSelector, PermissionScope, scopeConfig } from '@/components/Permissions/PermissionScopeSelector';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

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
  currentPermissions: Record<string, boolean>;
}

export default function RolePermissions({ 
  role, 
  allPermissions = [], 
  services = [], 
  actions = [],
  currentPermissions = {}
}: RolePermissionsProps) {
  // Store permissions as a map of service_action -> scope
  const [permissions, setPermissions] = useState<Record<string, PermissionScope>>({});
  const [isLoading, setIsLoading] = useState(false);
  const [hasChanges, setHasChanges] = useState(false);
  const [bulkService, setBulkService] = useState<string>('all');
  const [bulkAction, setBulkAction] = useState<string>('all');
  const [bulkScope, setBulkScope] = useState<PermissionScope>('none');

  // Initialize permissions from role data
  useEffect(() => {
    const perms: Record<string, PermissionScope> = {};
    
    // Process each service/action combination
    services.forEach(service => {
      Object.keys(service.actions).forEach(actionSlug => {
        const key = `${service.slug}_${actionSlug}`;
        
        // Check which scope the role has
        const byAllPerm = `${service.slug}_${actionSlug}_by_all`;
        const byRolePerm = `${service.slug}_${actionSlug}_by_my_role`;
        const byMePerm = `${service.slug}_${actionSlug}_by_me`;
        
        if (currentPermissions[byAllPerm]) {
          perms[key] = 'by_all';
        } else if (currentPermissions[byRolePerm]) {
          perms[key] = 'by_my_role';
        } else if (currentPermissions[byMePerm]) {
          perms[key] = 'by_me';
        } else if (currentPermissions[key]) {
          // Handle old-style permissions (treat as by_all)
          perms[key] = 'by_all';
        } else {
          perms[key] = 'none';
        }
      });
    });
    
    setPermissions(perms);
    setHasChanges(false);
  }, [role.id, currentPermissions, services]);

  // Check for changes
  // Note: key is already "service_action" (e.g., "procurement_notices_create")
  useEffect(() => {
    let changed = false;

    Object.entries(permissions).forEach(([key, scope]) => {
      // key is already "service_action", don't split it incorrectly
      // Check if current state differs from original
      const hasAny = ['by_all', 'by_my_role', 'by_me'].some(s =>
        currentPermissions[`${key}_${s}`] ||
        (s === 'by_all' && currentPermissions[key])
      );

      if (scope === 'none' && hasAny) {
        changed = true;
      } else if (scope !== 'none') {
        const expectedPerm = scope === 'by_all' ?
          [`${key}_by_all`, key] :
          [`${key}_${scope}`];

        const hasExpected = expectedPerm.some(p => currentPermissions[p]);
        if (!hasExpected) {
          changed = true;
        }
      }
    });

    setHasChanges(changed);
  }, [permissions, currentPermissions]);

  const handlePermissionChange = (serviceSlug: string, actionSlug: string, scope: PermissionScope) => {
    const key = `${serviceSlug}_${actionSlug}`;
    setPermissions(prev => ({
      ...prev,
      [key]: scope
    }));
  };

  const handleSelectAllForService = (serviceSlug: string, serviceActions: Record<string, boolean>, scope: PermissionScope) => {
    const updates: Record<string, PermissionScope> = {};
    
    Object.keys(serviceActions).forEach(actionSlug => {
      updates[`${serviceSlug}_${actionSlug}`] = scope;
    });
    
    setPermissions(prev => ({
      ...prev,
      ...updates
    }));
  };

  const applyBulkAction = () => {
    const updates: Record<string, PermissionScope> = {};
    
    services.forEach(service => {
      if (bulkService !== 'all' && service.slug !== bulkService) return;
      
      Object.keys(service.actions).forEach(actionSlug => {
        if (bulkAction !== 'all' && actionSlug !== bulkAction) return;
        
        updates[`${service.slug}_${actionSlug}`] = bulkScope;
      });
    });
    
    setPermissions(prev => ({
      ...prev,
      ...updates
    }));
    
    toast.success('Bulk permissions applied');
  };

  const applyPreset = (preset: 'admin' | 'editor' | 'viewer') => {
    const updates: Record<string, PermissionScope> = {};
    
    services.forEach(service => {
      Object.keys(service.actions).forEach(actionSlug => {
        const key = `${service.slug}_${actionSlug}`;
        
        switch (preset) {
          case 'admin':
            updates[key] = 'by_all';
            break;
          case 'editor':
            if (actionSlug === 'read') {
              updates[key] = 'by_all';
            } else if (actionSlug === 'create' || actionSlug === 'update') {
              updates[key] = 'by_my_role';
            } else if (actionSlug === 'delete') {
              updates[key] = 'by_me';
            } else {
              updates[key] = 'none';
            }
            break;
          case 'viewer':
            updates[key] = actionSlug === 'read' ? 'by_all' : 'none';
            break;
        }
      });
    });
    
    setPermissions(updates);
    toast.success(`${preset.charAt(0).toUpperCase() + preset.slice(1)} preset applied`);
  };

  const handleSave = async () => {
    setIsLoading(true);
    
    try {
      // Convert permissions to the format expected by the API
      // key is already "service_action" (e.g., "procurement_notices_create")
      // We just append the scope to form "service_action_scope"
      const permissionsToSave: string[] = [];

      Object.entries(permissions).forEach(([key, scope]) => {
        if (scope !== 'none') {
          // key is already in format "service_action", just append scope
          permissionsToSave.push(`${key}_${scope}`);
        }
      });
      
      const response = await fetch(`/api/roles/${role.id}/permissions`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify({
          permissions: permissionsToSave,
        }),
      });

      if (response.ok) {
        const responseData = await response.json();
        toast.success('Permissions updated successfully');
        setHasChanges(false);
        
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
    const counts = { by_all: 0, by_my_role: 0, by_me: 0, none: 0, total: 0 };
    
    Object.keys(serviceActions).forEach(actionSlug => {
      const scope = permissions[`${serviceSlug}_${actionSlug}`] || 'none';
      counts[scope]++;
      counts.total++;
    });
    
    return counts;
  };

  const totalPermissions = services.reduce((total, service) => 
    total + Object.keys(service.actions).length, 0
  );
  const permissionCounts = { by_all: 0, by_my_role: 0, by_me: 0, none: 0 };
  
  Object.values(permissions).forEach(scope => {
    permissionCounts[scope]++;
  });
  
  const activePermissions = totalPermissions - permissionCounts.none;

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
                  {activePermissions} of {totalPermissions} permissions
                </p>
                <div className="flex gap-1 text-xs text-muted-foreground">
                  <Badge variant="outline" className="text-xs px-1 py-0">
                    🌍 {permissionCounts.by_all}
                  </Badge>
                  <Badge variant="outline" className="text-xs px-1 py-0">
                    👥 {permissionCounts.by_my_role}
                  </Badge>
                  <Badge variant="outline" className="text-xs px-1 py-0">
                    👤 {permissionCounts.by_me}
                  </Badge>
                </div>
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

        {/* Presets and Bulk Actions */}
        <div className="px-4 lg:px-6">
          <Card className="bg-muted/30">
            <CardContent className="p-4">
              <div className="space-y-4">
                {/* Quick Presets */}
                <div>
                  <h4 className="text-sm font-medium mb-2 flex items-center gap-2">
                    <Settings2 className="h-4 w-4" />
                    Quick Presets
                  </h4>
                  <div className="flex gap-2">
                    <Button variant="outline" size="sm" onClick={() => applyPreset('admin')}>
                      <Shield className="w-4 h-4 mr-2" />
                      Admin Preset
                    </Button>
                    <Button variant="outline" size="sm" onClick={() => applyPreset('editor')}>
                      <Edit className="w-4 h-4 mr-2" />
                      Editor Preset
                    </Button>
                    <Button variant="outline" size="sm" onClick={() => applyPreset('viewer')}>
                      <Eye className="w-4 h-4 mr-2" />
                      Viewer Preset
                    </Button>
                  </div>
                </div>
                
                <Separator />
                
                {/* Bulk Actions */}
                <div>
                  <h4 className="text-sm font-medium mb-2">Bulk Actions</h4>
                  <div className="flex gap-2 items-end">
                    <div className="flex-1">
                      <label className="text-xs text-muted-foreground">Service</label>
                      <Select value={bulkService} onValueChange={setBulkService}>
                        <SelectTrigger className="w-full">
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="all">All Services</SelectItem>
                          {services.map(service => (
                            <SelectItem key={service.slug} value={service.slug}>
                              {service.name}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="flex-1">
                      <label className="text-xs text-muted-foreground">Action</label>
                      <Select value={bulkAction} onValueChange={setBulkAction}>
                        <SelectTrigger className="w-full">
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="all">All Actions</SelectItem>
                          {actions.map(action => (
                            <SelectItem key={action.slug} value={action.slug}>
                              {action.name}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="flex-1">
                      <label className="text-xs text-muted-foreground">Scope</label>
                      <Select value={bulkScope} onValueChange={(value) => setBulkScope(value as PermissionScope)}>
                        <SelectTrigger className="w-full">
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {Object.entries(scopeConfig).map(([scope, config]) => (
                            <SelectItem key={scope} value={scope}>
                              <div className="flex items-center gap-2">
                                <span className={config.color}>{config.icon}</span>
                                <span>{config.label}</span>
                              </div>
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <Button onClick={applyBulkAction} size="sm">
                      Apply
                    </Button>
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
              const counts = getPermissionCount(service.slug, service.actions);
              const hasAnyPermissions = counts.by_all + counts.by_my_role + counts.by_me > 0;
              const allMaxScope = counts.by_all === counts.total;

              return (
                <Card key={service.id} className="overflow-hidden">
                  <CardHeader className="pb-3">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <CardTitle className="text-lg">{service.name}</CardTitle>
                        <div className="flex gap-1">
                          {counts.by_all > 0 && (
                            <Badge variant="outline" className="text-xs px-2 py-1 text-green-600">
                              🌍 {counts.by_all}
                            </Badge>
                          )}
                          {counts.by_my_role > 0 && (
                            <Badge variant="outline" className="text-xs px-2 py-1 text-blue-600">
                              👥 {counts.by_my_role}
                            </Badge>
                          )}
                          {counts.by_me > 0 && (
                            <Badge variant="outline" className="text-xs px-2 py-1 text-orange-600">
                              👤 {counts.by_me}
                            </Badge>
                          )}
                          {counts.none === counts.total && (
                            <Badge variant="outline" className="text-xs px-2 py-1 text-gray-500">
                              No access
                            </Badge>
                          )}
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        <div className="flex gap-1">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleSelectAllForService(service.slug, service.actions, 'by_all')}
                            className="text-xs"
                          >
                            All Access
                          </Button>
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleSelectAllForService(service.slug, service.actions, 'by_my_role')}
                            className="text-xs"
                          >
                            Role Access
                          </Button>
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleSelectAllForService(service.slug, service.actions, 'none')}
                            className="text-xs"
                          >
                            <XCircle className="h-3 w-3 mr-1" />
                            Clear
                          </Button>
                        </div>
                      </div>
                    </div>
                  </CardHeader>
                  <Separator />
                  <CardContent className="pt-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                      {Object.keys(service.actions).map((actionSlug) => {
                        const key = `${service.slug}_${actionSlug}`;
                        const currentScope = permissions[key] || 'none';
                        const actionName = actions.find(a => a.slug === actionSlug)?.name || actionSlug;
                        const scopeColor = scopeConfig[currentScope].color;

                        return (
                          <div
                            key={actionSlug}
                            className={`
                              flex items-center justify-between p-3 rounded-lg border transition-colors
                              ${currentScope !== 'none' 
                                ? 'bg-primary/5 border-primary/20 dark:bg-primary/10' 
                                : 'bg-muted/30 border-border hover:bg-muted/50'
                              }
                            `}
                          >
                            <div className="flex items-center gap-2 min-w-0 flex-1">
                              <div className={`w-2 h-2 rounded-full flex-shrink-0 ${scopeColor.replace('text-', 'bg-')}`} />
                              <span className="text-sm font-medium truncate">
                                {actionName}
                              </span>
                            </div>
                            <div className="ml-2">
                              <PermissionScopeSelector
                                service={service.slug}
                                action={actionSlug}
                                currentScope={currentScope}
                                onChange={(scope) => handlePermissionChange(service.slug, actionSlug, scope)}
                              />
                            </div>
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
                    You have unsaved scoped permission changes. Don't forget to save your updates.
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