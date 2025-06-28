import React, { useState } from 'react';
import { Head, router } from '@inertiajs/react';
import Admin from '@/layouts/Admin';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { Checkbox } from '@/components/ui/checkbox';
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

interface RoleCreateProps {
  auth: any;
  allPermissions: Permission[];
  services: Service[];
  actions: Action[];
}

export default function RoleCreate({ auth, allPermissions, services, actions }: RoleCreateProps) {
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    level: 1,
  });
  const [selectedPermissions, setSelectedPermissions] = useState<Set<string>>(new Set());
  const [loading, setLoading] = useState(false);

  const handleBack = () => {
    router.visit('/admin/permissions');
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.name.trim()) {
      toast.error('Role name is required');
      return;
    }

    setLoading(true);
    try {
      const response = await fetch('/api/roles', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
          'Accept': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          name: formData.name.trim(),
          description: formData.description.trim(),
          level: formData.level,
          permissions: Array.from(selectedPermissions),
        }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to create role');
      }

      toast.success('Role created successfully');
      router.visit('/admin/permissions');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to create role');
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
    <Admin title="Create Role">
      <Head title="Create Role" />
      
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
                <h1 className="text-2xl font-bold">Create a role</h1>
                <p className="text-muted-foreground text-sm mt-1">Define the rights given to the role</p>
              </div>
            </div>
            <Button
              onClick={handleSubmit}
              disabled={loading || !formData.name.trim()}
              className="bg-primary hover:bg-primary/90 text-primary-foreground"
            >
              {loading ? 'Creating...' : 'Create'}
            </Button>
          </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Role Info */}
            <div className="bg-card rounded-lg p-6 border">
              <h2 className="text-xl font-semibold mb-4">Role Information</h2>
              
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="name" className="text-foreground">Name*</Label>
                  <Input
                    id="name"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    className="mt-1 bg-background border-border text-foreground"
                    placeholder="e.g., Editor, Manager"
                    required
                  />
                </div>
                <div>
                  <Label htmlFor="level" className="text-foreground">Level</Label>
                  <Input
                    id="level"
                    type="number"
                    min="1"
                    max="100"
                    value={formData.level}
                    onChange={(e) => setFormData({ ...formData, level: parseInt(e.target.value) || 1 })}
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
                    placeholder="Describe what this role can do..."
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
          </form>
        </div>
      </div>
    </Admin>
  );
}