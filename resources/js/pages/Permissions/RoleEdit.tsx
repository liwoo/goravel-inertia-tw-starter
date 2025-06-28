import React, { useState } from 'react';
import { Head, router } from '@inertiajs/react';
import Admin from '@/layouts/Admin';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { Checkbox } from '@/components/ui/checkbox';
import { Badge } from '@/components/ui/badge';
import { ArrowLeft } from 'lucide-react';
import { toast } from 'sonner';

interface Permission {
  id: number;
  name: string;
  slug: string;
  category: string;
  resource: string;
  action: string;
}

interface Role {
  id: number;
  name: string;
  slug: string;
  description: string;
  level: number;
  is_active: boolean;
  users_count: number;
  permissions: string[]; // Array of permission slugs
}

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

interface RoleEditProps {
  auth: any;
  role: Role;
  allPermissions: Permission[];
  services: Service[];
  actions: Action[];
}

export default function RoleEdit({ auth, role, allPermissions, services, actions }: RoleEditProps) {
  const [formData, setFormData] = useState({
    name: role.name,
    description: role.description,
  });
  const [selectedPermissions, setSelectedPermissions] = useState<Set<string>>(
    new Set(role.permissions)
  );
  const [loading, setLoading] = useState(false);

  const handleBack = () => {
    router.visit('/admin/permissions');
  };

  const handleSave = async () => {
    setLoading(true);
    try {
      const response = await fetch(`/api/roles/${role.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
          'Accept': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          name: formData.name,
          description: formData.description,
          permissions: Array.from(selectedPermissions),
        }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to update role');
      }

      toast.success('Role updated successfully');
      router.visit('/admin/permissions');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to update role');
    } finally {
      setLoading(false);
    }
  };

  const togglePermission = (service: string, action: string) => {
    const permissionSlug = `${service}_${action}`;
    const newPermissions = new Set(selectedPermissions);
    
    if (newPermissions.has(permissionSlug)) {
      newPermissions.delete(permissionSlug);
    } else {
      newPermissions.add(permissionSlug);
    }
    
    setSelectedPermissions(newPermissions);
  };

  const toggleAllPermissions = (service: string, value: boolean) => {
    const serviceData = services.find(s => s.slug === service);
    if (!serviceData) return;
    
    const newPermissions = new Set(selectedPermissions);
    
    Object.keys(serviceData.actions).forEach(action => {
      if (serviceData.actions[action]) {
        const permissionSlug = `${service}_${action}`;
        if (value) {
          newPermissions.add(permissionSlug);
        } else {
          newPermissions.delete(permissionSlug);
        }
      }
    });
    
    setSelectedPermissions(newPermissions);
  };

  const renderPermissionMatrix = () => {
    return (
      <div className="rounded-md border overflow-hidden">
        <div className="overflow-x-auto max-h-[600px] overflow-y-auto">
          <table className="w-full min-w-full">
            {/* Header */}
            <thead className="bg-muted/50 sticky top-0 z-10">
              <tr>
                <th className="px-4 py-3 text-left font-medium text-muted-foreground min-w-[200px] bg-muted/50 border-b">
                  SERVICE
                </th>
                {actions.map((action) => (
                  <th 
                    key={action.slug} 
                    className="px-2 py-3 text-center font-medium text-muted-foreground min-w-[80px] bg-muted/50 border-b border-l"
                  >
                    <div className="text-xs uppercase whitespace-nowrap">
                      {action.name}
                    </div>
                  </th>
                ))}
              </tr>
            </thead>
            
            {/* Body */}
            <tbody>
              {services.map((service) => {
                const serviceActions = Object.keys(service.actions).filter(action => service.actions[action]);
                const allChecked = serviceActions.every(action => selectedPermissions.has(`${service.slug}_${action}`));
                const someChecked = serviceActions.some(action => selectedPermissions.has(`${service.slug}_${action}`));

                return (
                  <tr key={service.slug} className="hover:bg-muted/50 border-b">
                    <td className="px-4 py-3 min-w-[200px]">
                      <div className="flex items-center space-x-3">
                        <Checkbox
                          checked={allChecked}
                          ref={(el) => {
                            if (el) el.indeterminate = someChecked && !allChecked;
                          }}
                          onCheckedChange={(checked) => toggleAllPermissions(service.slug, checked as boolean)}
                          className="data-[state=checked]:bg-primary data-[state=checked]:border-primary"
                        />
                        <span className="text-foreground font-medium whitespace-nowrap">
                          {service.name}
                        </span>
                      </div>
                    </td>
                    
                    {actions.map((action) => {
                      const isSupported = service.actions[action.slug];
                      return (
                        <td key={action.slug} className="px-2 py-3 text-center min-w-[80px] border-l">
                          <div className="flex justify-center">
                            <Checkbox
                              checked={selectedPermissions.has(`${service.slug}_${action.slug}`)}
                              onCheckedChange={() => togglePermission(service.slug, action.slug)}
                              className="data-[state=checked]:bg-primary data-[state=checked]:border-primary"
                              disabled={!isSupported}
                            />
                          </div>
                        </td>
                      );
                    })}
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>
    );
  };

  return (
    <Admin title={`Edit ${role.name} Role`}>
      <Head title={`Edit ${role.name} Role`} />
      
      <div className="space-y-6 min-w-0 overflow-hidden">
        <div className="max-w-5xl mx-auto">
          {/* Header */}
          <div className="flex items-center justify-between mb-8">
            <div className="flex items-center space-x-4">
              <Button
                variant="ghost"
                onClick={handleBack}
                className="text-muted-foreground hover:text-foreground"
              >
                <ArrowLeft className="w-4 h-4 mr-2" />
                Back
              </Button>
              <div>
                <h1 className="text-2xl font-bold">Edit a role</h1>
                <p className="text-muted-foreground text-sm mt-1">Define the rights given to the role</p>
              </div>
            </div>
            <Button
              onClick={handleSave}
              disabled={loading}
              className="bg-primary hover:bg-primary/90 text-primary-foreground"
            >
              {loading ? 'Saving...' : 'Save'}
            </Button>
          </div>

          {/* Role Info */}
          <div className="bg-card rounded-lg p-6 mb-6 border">
            <div className="flex items-center justify-between mb-4">
              <div>
                <h2 className="text-xl font-semibold">{role.name}</h2>
                <p className="text-muted-foreground text-sm mt-1">{role.description}</p>
              </div>
              <Badge variant="secondary" className="bg-muted text-muted-foreground">
                {role.users_count} users with this role
              </Badge>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label htmlFor="name" className="text-foreground">Name*</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="mt-1 bg-background border-border text-foreground"
                />
              </div>
              <div className="col-span-2">
                <Label htmlFor="description" className="text-foreground">Description</Label>
                <Textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  className="mt-1 bg-background border-border text-foreground"
                  rows={3}
                />
              </div>
            </div>
          </div>

          {/* Permissions Matrix */}
          <div className="bg-card rounded-lg p-6 border">
            <h2 className="text-xl font-semibold mb-4">Permissions</h2>
            <p className="text-sm text-muted-foreground mb-6">
              Select the permissions this role should have for each service. Disabled checkboxes indicate actions not supported by that service.
            </p>
            
            {renderPermissionMatrix()}
          </div>
        </div>
      </div>
    </Admin>
  );
}