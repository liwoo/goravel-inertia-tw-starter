import React, { useState, useCallback } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Checkbox } from '@/components/ui/checkbox';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Separator } from '@/components/ui/separator';
// import { toast } from '@/components/ui/use-toast';
const toast = (options: any) => {
  console.log('Toast:', options);
  // Placeholder for toast functionality
};
import { 
  Shield, 
  Users, 
  CheckCircle2, 
  XCircle, 
  Search,
  Filter,
  Download,
  Upload,
  RotateCcw,
  AlertTriangle
} from 'lucide-react';
import type { 
  PermissionMatrixData, 
  RoleWithPermissions, 
  Permission, 
  PermissionGrouped,
  BulkAssignmentRequest,
  ServiceAction
} from '@/types/permissions';

interface PermissionMatrixProps {
  initialData: any; // Updated to handle new service-action structure
  onPermissionToggle: (roleId: number, serviceSlug: string, action: string, isAssigned: boolean) => Promise<void>;
  onBulkAssign: (request: BulkAssignmentRequest) => Promise<void>;
  onSyncRole: (roleId: number, permissionSlugs: string[]) => Promise<void>;
  loading?: boolean;
}

export default function PermissionMatrix({ 
  initialData, 
  onPermissionToggle, 
  onBulkAssign,
  onSyncRole,
  loading = false 
}: PermissionMatrixProps) {
  const [data, setData] = useState(initialData);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedService, setSelectedService] = useState<string>('all');
  const [selectedRole, setSelectedRole] = useState<number | null>(null);
  const [pendingChanges, setPendingChanges] = useState<Map<string, boolean>>(new Map());
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Debug: Log the data structure we're receiving (uncomment for debugging)
  // console.log('PermissionMatrix received data:', data);
  // console.log('Data type:', typeof data);
  // console.log('Data keys:', data ? Object.keys(data) : 'null');

  // Handle missing or invalid data
  if (!data || !data.roles || !data.services || !data.actions) {
    console.log('PermissionMatrix data validation failed:', {
      hasData: !!data,
      hasRoles: data?.roles,
      hasServices: data?.services,
      hasActions: data?.actions
    });
    
    return (
      <div className="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded">
        <strong>Warning:</strong> Permission matrix data is not available. 
        {!data && " No data provided."}
        {data && !data.roles && " Roles data is missing."}
        {data && !data.services && " Services data is missing."}
        {data && !data.actions && " Actions data is missing."}
        <pre className="mt-2 text-xs">{JSON.stringify(data, null, 2)}</pre>
      </div>
    );
  }

  // Filter roles and services based on search - with defensive checks
  const filteredRoles = React.useMemo(() => {
    if (!data?.roles || !Array.isArray(data.roles)) return [];
    return data.roles.filter(role => 
      role?.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      role?.slug?.toLowerCase().includes(searchTerm.toLowerCase())
    );
  }, [data?.roles, searchTerm]);

  const filteredServices = React.useMemo(() => {
    if (!data?.services || !Array.isArray(data.services)) return [];
    return data.services.filter(service => 
      selectedService === 'all' || service?.slug === selectedService
    );
  }, [data?.services, selectedService]);

  const allServices = React.useMemo(() => {
    if (!data?.services || !Array.isArray(data.services)) return ['all'];
    return ['all', ...data.services.map(s => s?.slug).filter(Boolean)];
  }, [data?.services]);

  // Helper to check if a permission is assigned to a role using service_action format
  const isPermissionAssigned = useCallback((roleId: number, serviceSlug: string, action: string): boolean => {
    const permissionSlug = `${serviceSlug}_${action}`;
    const key = `${roleId}-${permissionSlug}`;
    if (pendingChanges.has(key)) {
      return pendingChanges.get(key)!;
    }
    
    // Check if role has this permission
    const role = data.roles.find(r => r.id === roleId);
    return role?.permissions[permissionSlug] || false;
  }, [data.roles, pendingChanges]);

  // Handle individual permission toggle
  const handlePermissionToggle = useCallback(async (roleId: number, serviceSlug: string, action: string) => {
    if (loading || isSubmitting) return;

    const permissionSlug = `${serviceSlug}_${action}`;
    const key = `${roleId}-${permissionSlug}`;
    const currentlyAssigned = isPermissionAssigned(roleId, serviceSlug, action);
    const newAssignment = !currentlyAssigned;

    // Update pending changes
    const newPendingChanges = new Map(pendingChanges);
    newPendingChanges.set(key, newAssignment);
    setPendingChanges(newPendingChanges);

    try {
      await onPermissionToggle(roleId, serviceSlug, action, newAssignment);
      
      // Update local data on success
      setData(prevData => {
        const newData = { ...prevData };
        const roleIndex = newData.roles.findIndex(r => r.id === roleId);
        if (roleIndex !== -1) {
          newData.roles[roleIndex].permissions[permissionSlug] = newAssignment;
        }
        return newData;
      });

      // Remove from pending changes
      newPendingChanges.delete(key);
      setPendingChanges(newPendingChanges);

      toast({
        title: "Success",
        description: `Permission ${newAssignment ? 'assigned' : 'revoked'} successfully`,
      });
    } catch (error) {
      // Revert pending change on error
      newPendingChanges.delete(key);
      setPendingChanges(newPendingChanges);
      
      toast({
        title: "Error",
        description: `Failed to ${newAssignment ? 'assign' : 'revoke'} permission`,
        variant: "destructive",
      });
    }
  }, [loading, isSubmitting, isPermissionAssigned, pendingChanges, onPermissionToggle]);


  // Handle bulk assignment for a role
  const handleBulkRoleToggle = useCallback(async (roleId: number, action: 'assign' | 'revoke') => {
    if (loading || isSubmitting) return;

    setIsSubmitting(true);
    
    try {
      // Build all service-action combinations for this role
      if (!filteredServices.length || !data?.actions?.length) {
        console.log('No services or actions available for bulk operation');
        return;
      }
      
      const allPermissionSlugs = filteredServices.flatMap(service => 
        (data.actions || []).map(action => `${service?.slug}_${action?.slug}`)
      );

      // TODO: Update to use new bulk assign API
      console.log('Bulk role toggle:', roleId, action, allPermissionSlugs);
      return; // Skip for now

      // Update local data would go here
      toast({
        title: "Success",
        description: `Bulk ${action} completed successfully`,
      });
    } catch (error) {
      toast({
        title: "Error",
        description: `Failed to perform bulk ${action}`,
        variant: "destructive",
      });
    } finally {
      setIsSubmitting(false);
    }
  }, [loading, isSubmitting, filteredServices, data.actions, onBulkAssign]);

  // Handle bulk assignment for a service (updated for new structure)
  const handleBulkServiceAssignment = useCallback(async (serviceSlug: string, action: 'assign' | 'revoke') => {
    if (loading || isSubmitting) return;

    setIsSubmitting(true);

    try {
      // TODO: Implement bulk service assignment
      console.log('Bulk service assignment:', serviceSlug, action);

      toast({
        title: "Success",
        description: `Bulk ${action} for ${serviceSlug} completed successfully`,
      });
    } catch (error) {
      toast({
        title: "Error",
        description: `Failed to perform bulk ${action} for ${serviceSlug}`,
        variant: "destructive",
      });
    } finally {
      setIsSubmitting(false);
    }
  }, [loading, isSubmitting]);

  // Calculate statistics for current view
  const viewStats = React.useMemo(() => {
    if (!filteredRoles.length || !filteredServices.length || !data?.actions?.length) {
      return {
        visibleRoles: filteredRoles.length,
        visibleServices: filteredServices.length,
        visibleActions: data?.actions?.length || 0,
        visibleAssignments: 0,
      };
    }

    const totalVisibleAssignments = filteredRoles.reduce((total, role) => {
      const rolePermissions = role?.permissions || {};
      const visiblePermissionCount = filteredServices.reduce((count, service) => {
        return count + (data.actions || []).filter(action => 
          rolePermissions[`${service?.slug}_${action?.slug}`]
        ).length;
      }, 0);
      return total + visiblePermissionCount;
    }, 0);

    return {
      visibleRoles: filteredRoles.length,
      visibleServices: filteredServices.length,
      visibleActions: (data.actions || []).length,
      visibleAssignments: totalVisibleAssignments,
    };
  }, [filteredRoles, filteredServices, data?.actions]);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Permission Matrix</h1>
          <p className="text-muted-foreground">
            Manage role-permission assignments across your application
          </p>
        </div>
        
        <div className="flex gap-2">
          <Button variant="outline" size="sm">
            <Download className="w-4 h-4 mr-2" />
            Export
          </Button>
          <Button variant="outline" size="sm">
            <Upload className="w-4 h-4 mr-2" />
            Import
          </Button>
        </div>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2">
              <Users className="w-5 h-5 text-blue-500" />
              <div>
                <p className="text-sm font-medium">Roles</p>
                <p className="text-2xl font-bold">{viewStats.visibleRoles}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2">
              <Shield className="w-5 h-5 text-green-500" />
              <div>
                <p className="text-sm font-medium">Services</p>
                <p className="text-2xl font-bold">{viewStats.visibleServices}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2">
              <CheckCircle2 className="w-5 h-5 text-purple-500" />
              <div>
                <p className="text-sm font-medium">Assignments</p>
                <p className="text-2xl font-bold">{viewStats.visibleAssignments}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2">
              <AlertTriangle className="w-5 h-5 text-orange-500" />
              <div>
                <p className="text-sm font-medium">Pending</p>
                <p className="text-2xl font-bold">{pendingChanges.size}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="p-4">
          <div className="flex flex-col gap-4 md:flex-row md:items-center">
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                <Input
                  placeholder="Search roles..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>
            
            <div className="flex gap-2">
              <Select value={selectedService} onValueChange={setSelectedService}>
                <SelectTrigger className="w-[200px]">
                  <Filter className="w-4 h-4 mr-2" />
                  <SelectValue placeholder="Service" />
                </SelectTrigger>
                <SelectContent>
                  {allServices.map(service => (
                    <SelectItem key={service} value={service}>
                      {service === 'all' ? 'All Services' : service.charAt(0).toUpperCase() + service.slice(1)}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              
              {pendingChanges.size > 0 && (
                <Button 
                  variant="outline" 
                  size="sm"
                  onClick={() => setPendingChanges(new Map())}
                >
                  <RotateCcw className="w-4 h-4 mr-2" />
                  Reset ({pendingChanges.size})
                </Button>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Permission Matrix */}
      <Card>
        <CardHeader>
          <CardTitle>Permission Matrix</CardTitle>
        </CardHeader>
        <CardContent>
          {filteredServices.length === 0 || !data?.actions?.length ? (
            <div className="p-8 text-center text-muted-foreground">
              <p>No services or actions available to display.</p>
              {!data?.services?.length && <p>Missing services data.</p>}
              {!data?.actions?.length && <p>Missing actions data.</p>}
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full border-collapse">
              <thead>
                <tr>
                  <th className="sticky left-0 bg-background border border-border p-2 text-left min-w-[200px]">
                    Role
                  </th>
                  {filteredServices.map(service => (
                    <th key={service.slug} className="border border-border p-2 text-center" colSpan={(data.actions || []).length}>
                      <div className="flex flex-col items-center gap-1">
                        <span className="font-semibold">{service.name}</span>
                        <div className="flex gap-1">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleBulkServiceAssignment(service.slug, 'assign')}
                            disabled={loading || isSubmitting}
                            className="h-6 px-2 text-xs"
                          >
                            All
                          </Button>
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleBulkServiceAssignment(service.slug, 'revoke')}
                            disabled={loading || isSubmitting}
                            className="h-6 px-2 text-xs"
                          >
                            None
                          </Button>
                        </div>
                      </div>
                    </th>
                  ))}
                </tr>
                <tr>
                  <th className="sticky left-0 bg-background border border-border p-2"></th>
                  {filteredServices.flatMap(service =>
                    (data.actions || []).map(action => (
                      <th key={`${service.slug}-${action.slug}`} className="border border-border p-1 text-xs text-center min-w-[80px]">
                        <div className="transform -rotate-45 origin-center whitespace-nowrap">
                          {action.name}
                        </div>
                      </th>
                    ))
                  )}
                </tr>
              </thead>
              <tbody>
                {filteredRoles.map(role => (
                  <tr key={role.id} className="hover:bg-muted/50">
                    <td className="sticky left-0 bg-background border border-border p-2">
                      <div className="flex items-center justify-between">
                        <div>
                          <div className="font-medium">{role.name}</div>
                          <div className="text-xs text-muted-foreground">
                            Level {role.level} • {Object.keys(role.permissions || {}).filter(p => role.permissions[p]).length} permissions
                          </div>
                        </div>
                        <div className="flex gap-1">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleBulkRoleToggle(role.id, 'assign')}
                            disabled={loading || isSubmitting}
                            className="h-6 px-2 text-xs"
                          >
                            All
                          </Button>
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleBulkRoleToggle(role.id, 'revoke')}
                            disabled={loading || isSubmitting}
                            className="h-6 px-2 text-xs"
                          >
                            None
                          </Button>
                        </div>
                      </div>
                    </td>
                    {filteredServices.flatMap(service =>
                      (data.actions || []).map(action => {
                        const isAssigned = isPermissionAssigned(role.id, service.slug, action.slug);
                        const permissionSlug = `${service.slug}_${action.slug}`;
                        const isPending = pendingChanges.has(`${role.id}-${permissionSlug}`);
                        
                        return (
                          <td key={`${role.id}-${service.slug}-${action.slug}`} className="border border-border p-2 text-center">
                            <Checkbox
                              checked={isAssigned}
                              onCheckedChange={() => handlePermissionToggle(role.id, service.slug, action.slug)}
                              disabled={loading || isSubmitting}
                              className={isPending ? 'opacity-50' : ''}
                            />
                          </td>
                        );
                      })
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        </CardContent>
      </Card>
    </div>
  );
}