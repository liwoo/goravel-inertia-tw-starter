import React, { useState, useCallback } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Checkbox } from '@/components/ui/checkbox';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { Separator } from '@/components/ui/separator';
import { 
  Users, 
  Shield, 
  Plus,
  Edit,
  Trash2,
  Save,
  X,
  CheckCircle2,
  AlertTriangle
} from 'lucide-react';

// Toast placeholder
const toast = (options: any) => {
  console.log('Toast:', options);
};

interface RoleData {
  id: number;
  name: string;
  slug: string;
  description: string;
  level: number;
  permissions: Record<string, boolean>;
  is_active: boolean;
}

interface ServiceData {
  id: string;
  name: string;
  slug: string;
  actions: string[];
}

interface ActionData {
  id: string;
  name: string;
  slug: string;
}

interface RoleManagerProps {
  initialData: {
    roles: RoleData[];
    services: ServiceData[];
    actions: ActionData[];
    stats: Record<string, number>;
  };
  onCreateRole: (roleData: Partial<RoleData>) => Promise<void>;
  onUpdateRole: (roleId: number, roleData: Partial<RoleData>) => Promise<void>;
  onDeleteRole: (roleId: number) => Promise<void>;
  onPermissionToggle: (roleId: number, serviceSlug: string, action: string, isAssigned: boolean) => Promise<void>;
  loading?: boolean;
}

export default function RoleManager({ 
  initialData, 
  onCreateRole,
  onUpdateRole,
  onDeleteRole,
  onPermissionToggle,
  loading = false 
}: RoleManagerProps) {
  const [data, setData] = useState(initialData);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingRole, setEditingRole] = useState<RoleData | null>(null);
  const [pendingChanges, setPendingChanges] = useState<Map<string, boolean>>(new Map());
  const [isSubmitting, setIsSubmitting] = useState(false);

  // New role form data
  const [newRole, setNewRole] = useState({
    name: '',
    description: '',
    level: 1,
  });

  // Handle missing or invalid data
  if (!data || !data.roles || !data.services || !data.actions) {
    return (
      <div className="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded">
        <strong>Warning:</strong> Role management data is not available.
        {!data && " No data provided."}
        {data && !data.roles && " Roles data is missing."}
        {data && !data.services && " Services data is missing."}
        {data && !data.actions && " Actions data is missing."}
      </div>
    );
  }

  // Helper to check if a permission is assigned to a role
  const isPermissionAssigned = useCallback((roleId: number, serviceSlug: string, action: string): boolean => {
    const permissionSlug = `${serviceSlug}_${action}`;
    const key = `${roleId}-${permissionSlug}`;
    if (pendingChanges.has(key)) {
      return pendingChanges.get(key)!;
    }
    
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

  // Handle role creation
  const handleCreateRole = async () => {
    if (!newRole.name.trim()) return;

    setIsSubmitting(true);
    try {
      await onCreateRole({
        name: newRole.name,
        description: newRole.description,
        level: newRole.level,
      });

      setNewRole({ name: '', description: '', level: 1 });
      setIsCreateDialogOpen(false);
      
      toast({
        title: "Success",
        description: "Role created successfully",
      });
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to create role",
        variant: "destructive",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  // Handle role update
  const handleUpdateRole = async (roleId: number, updates: Partial<RoleData>) => {
    setIsSubmitting(true);
    try {
      await onUpdateRole(roleId, updates);
      setEditingRole(null);
      
      toast({
        title: "Success",
        description: "Role updated successfully",
      });
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to update role",
        variant: "destructive",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  // Handle role deletion
  const handleDeleteRole = async (roleId: number) => {
    if (!confirm('Are you sure you want to delete this role?')) return;

    setIsSubmitting(true);
    try {
      await onDeleteRole(roleId);
      
      toast({
        title: "Success",
        description: "Role deleted successfully",
      });
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to delete role",
        variant: "destructive",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  // Calculate role statistics
  const getRoleStats = (role: RoleData) => {
    const totalPermissions = Object.keys(role.permissions || {}).length;
    const assignedPermissions = Object.values(role.permissions || {}).filter(Boolean).length;
    const percentage = totalPermissions > 0 ? Math.round((assignedPermissions / totalPermissions) * 100) : 0;
    
    return { totalPermissions, assignedPermissions, percentage };
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Role Management</h1>
          <p className="text-muted-foreground">
            Create and manage user roles with granular permissions
          </p>
        </div>
        
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="w-4 h-4 mr-2" />
              Create Role
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>Create New Role</DialogTitle>
              <DialogDescription>
                Define a new role with a name, description, and permission level.
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div>
                <Label htmlFor="name">Role Name</Label>
                <Input
                  id="name"
                  value={newRole.name}
                  onChange={(e) => setNewRole(prev => ({ ...prev, name: e.target.value }))}
                  placeholder="e.g., Editor, Manager"
                />
              </div>
              <div>
                <Label htmlFor="description">Description</Label>
                <Textarea
                  id="description"
                  value={newRole.description}
                  onChange={(e) => setNewRole(prev => ({ ...prev, description: e.target.value }))}
                  placeholder="Describe what this role can do..."
                  rows={3}
                />
              </div>
              <div>
                <Label htmlFor="level">Permission Level</Label>
                <Input
                  id="level"
                  type="number"
                  min="1"
                  max="100"
                  value={newRole.level}
                  onChange={(e) => setNewRole(prev => ({ ...prev, level: parseInt(e.target.value) || 1 }))}
                />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsCreateDialogOpen(false)}>
                Cancel
              </Button>
              <Button onClick={handleCreateRole} disabled={!newRole.name.trim() || isSubmitting}>
                {isSubmitting ? "Creating..." : "Create Role"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2">
              <Users className="w-5 h-5 text-blue-500" />
              <div>
                <p className="text-sm font-medium">Total Roles</p>
                <p className="text-2xl font-bold">{data.roles.length}</p>
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
                <p className="text-2xl font-bold">{data.services.length}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2">
              <CheckCircle2 className="w-5 h-5 text-purple-500" />
              <div>
                <p className="text-sm font-medium">Actions</p>
                <p className="text-2xl font-bold">{data.actions.length}</p>
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

      {/* Roles Accordion */}
      <Card>
        <CardHeader>
          <CardTitle>Roles & Permissions</CardTitle>
        </CardHeader>
        <CardContent>
          <Accordion type="single" collapsible className="w-full">
            {data.roles.map((role) => {
              const stats = getRoleStats(role);
              
              return (
                <AccordionItem key={role.id} value={`role-${role.id}`}>
                  <AccordionTrigger className="hover:no-underline">
                    <div className="flex items-center justify-between w-full pr-4">
                      <div className="flex items-center gap-3">
                        <div>
                          <div className="flex items-center gap-2">
                            <h3 className="font-semibold text-left">{role.name}</h3>
                            <Badge variant={role.is_active ? "default" : "secondary"}>
                              Level {role.level}
                            </Badge>
                          </div>
                          <p className="text-sm text-muted-foreground text-left">
                            {role.description || "No description"}
                          </p>
                        </div>
                      </div>
                      <div className="flex items-center gap-4">
                        <div className="text-right">
                          <div className="text-sm font-medium">
                            {stats.assignedPermissions}/{stats.totalPermissions} permissions
                          </div>
                          <div className="text-xs text-muted-foreground">
                            {stats.percentage}% configured
                          </div>
                        </div>
                        <div className="flex gap-1">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={(e) => {
                              e.stopPropagation();
                              setEditingRole(role);
                            }}
                          >
                            <Edit className="w-4 h-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={(e) => {
                              e.stopPropagation();
                              handleDeleteRole(role.id);
                            }}
                          >
                            <Trash2 className="w-4 h-4" />
                          </Button>
                        </div>
                      </div>
                    </div>
                  </AccordionTrigger>
                  <AccordionContent>
                    <div className="pt-4">
                      <div className="grid gap-4">
                        {data.services.map((service) => (
                          <Card key={service.slug} className="border-muted">
                            <CardHeader className="pb-3">
                              <CardTitle className="text-lg">{service.name}</CardTitle>
                            </CardHeader>
                            <CardContent>
                              <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3">
                                {data.actions.map((action) => {
                                  const isAssigned = isPermissionAssigned(role.id, service.slug, action.slug);
                                  const permissionSlug = `${service.slug}_${action.slug}`;
                                  const isPending = pendingChanges.has(`${role.id}-${permissionSlug}`);
                                  
                                  return (
                                    <div
                                      key={`${service.slug}-${action.slug}`}
                                      className="flex items-center space-x-2 p-2 rounded border"
                                    >
                                      <Checkbox
                                        id={`${role.id}-${service.slug}-${action.slug}`}
                                        checked={isAssigned}
                                        onCheckedChange={() => handlePermissionToggle(role.id, service.slug, action.slug)}
                                        disabled={loading || isSubmitting}
                                        className={isPending ? 'opacity-50' : ''}
                                      />
                                      <Label
                                        htmlFor={`${role.id}-${service.slug}-${action.slug}`}
                                        className="text-sm font-medium cursor-pointer"
                                      >
                                        {action.name}
                                      </Label>
                                    </div>
                                  );
                                })}
                              </div>
                            </CardContent>
                          </Card>
                        ))}
                      </div>
                    </div>
                  </AccordionContent>
                </AccordionItem>
              );
            })}
          </Accordion>
        </CardContent>
      </Card>

      {/* Edit Role Dialog */}
      {editingRole && (
        <Dialog open={!!editingRole} onOpenChange={() => setEditingRole(null)}>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>Edit Role</DialogTitle>
              <DialogDescription>
                Update the role name, description, and permission level.
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div>
                <Label htmlFor="edit-name">Role Name</Label>
                <Input
                  id="edit-name"
                  value={editingRole.name}
                  onChange={(e) => setEditingRole(prev => prev ? { ...prev, name: e.target.value } : null)}
                />
              </div>
              <div>
                <Label htmlFor="edit-description">Description</Label>
                <Textarea
                  id="edit-description"
                  value={editingRole.description}
                  onChange={(e) => setEditingRole(prev => prev ? { ...prev, description: e.target.value } : null)}
                  rows={3}
                />
              </div>
              <div>
                <Label htmlFor="edit-level">Permission Level</Label>
                <Input
                  id="edit-level"
                  type="number"
                  min="1"
                  max="100"
                  value={editingRole.level}
                  onChange={(e) => setEditingRole(prev => prev ? { ...prev, level: parseInt(e.target.value) || 1 } : null)}
                />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setEditingRole(null)}>
                Cancel
              </Button>
              <Button onClick={() => editingRole && handleUpdateRole(editingRole.id, editingRole)} disabled={isSubmitting}>
                {isSubmitting ? "Saving..." : "Save Changes"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      )}
    </div>
  );
}