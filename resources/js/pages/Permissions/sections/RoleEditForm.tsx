import React, { useState, forwardRef, useImperativeHandle } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { CrudEditFormProps } from '@/types/crud';
import { Role } from '@/types/permissions';
import { Shield, FileText, Hash, Users, Calendar } from 'lucide-react';
import { Switch } from '@/components/ui/switch';
import { Separator } from '@/components/ui/separator';
import { Badge } from '@/components/ui/badge';

interface RoleFormData {
  name: string;
  slug: string;
  description: string;
  level: number;
  is_active: boolean;
  parent_id?: number;
}

interface RoleEditFormProps extends CrudEditFormProps<Role> {
  setIsSaving?: (saving: boolean) => void;
}

export const RoleEditForm = forwardRef<any, RoleEditFormProps>(({ 
  item: role,
  onSuccess,
  onError,
  onCancel, 
  isLoading = false,
  setIsSaving
}, ref) => {
  const [formData, setFormData] = useState<RoleFormData>({
    name: role.name,
    slug: role.slug,
    description: role.description || '',
    level: role.level,
    is_active: role.is_active,
    parent_id: role.parent_id,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleSubmit = async () => {
    
    // Basic validation
    const newErrors: Record<string, string> = {};
    if (!formData.name.trim()) {
      newErrors.name = 'Name is required';
    }
    if (!formData.slug.trim()) {
      newErrors.slug = 'Slug is required';
    }
    if (formData.level < 1) {
      newErrors.level = 'Level must be at least 1';
    }
    
    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }
    
    setErrors({});
    setIsSaving?.(true);
    
    try {
      const response = await fetch(`/api/roles/${role.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        onSuccess('Role updated successfully');
      } else {
        const errorData = await response.json().catch(() => ({}));
        onError?.(errorData);
      }
    } catch (error) {
      onError?.(error);
    } finally {
      setIsSaving?.(false);
    }
  };

  // Expose handleSubmit to parent component
  useImperativeHandle(ref, () => ({
    handleSubmit
  }));

  return (
    <form onSubmit={(e) => e.preventDefault()} className="space-y-8">
      {/* Role Header */}
      <div className="flex justify-center">
        <div className={`p-6 rounded-full ${role.is_active ? 'bg-teal-100 dark:bg-teal-900/30' : 'bg-gray-100 dark:bg-gray-900/30'}`}>
          <Shield className={`h-16 w-16 ${role.is_active ? 'text-teal-600 dark:text-teal-400' : 'text-gray-600 dark:text-gray-400'}`} />
        </div>
      </div>
      <div className="text-center">
        <h2 className="text-2xl font-bold text-foreground">{role.name}</h2>
        <Badge variant="secondary" className="mt-2">ID: #{role.id}</Badge>
      </div>

      {/* Basic Information */}
      <div className="space-y-6">
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">Basic Information</h3>
          <div className="space-y-4">
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Shield className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="name">Role Name *</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="e.g., Content Editor"
                  className={errors.name ? 'border-destructive' : ''}
                />
                {errors.name && (
                  <p className="text-sm text-destructive">{errors.name}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Hash className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="slug">Slug *</Label>
                <Input
                  id="slug"
                  value={formData.slug}
                  onChange={(e) => setFormData({ ...formData, slug: e.target.value })}
                  placeholder="e.g., content-editor"
                  className={errors.slug ? 'border-destructive' : ''}
                />
                {errors.slug && (
                  <p className="text-sm text-destructive">{errors.slug}</p>
                )}
                <p className="text-xs text-muted-foreground">
                  URL-friendly identifier
                </p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <FileText className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="description">Description</Label>
                <Textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  placeholder="Brief description of the role's responsibilities"
                  rows={3}
                  className="resize-none"
                />
              </div>
            </div>
          </div>
        </div>

        <Separator className="my-6" />

        {/* Role Settings */}
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">Role Settings</h3>
          <div className="space-y-4">
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Users className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="level">Hierarchy Level *</Label>
                <Input
                  id="level"
                  type="number"
                  min="1"
                  max="100"
                  value={formData.level}
                  onChange={(e) => setFormData({ ...formData, level: parseInt(e.target.value) || 1 })}
                  className={errors.level ? 'border-destructive' : ''}
                />
                {errors.level && (
                  <p className="text-sm text-destructive">{errors.level}</p>
                )}
                <p className="text-xs text-muted-foreground">
                  Higher numbers indicate higher authority (1-100)
                </p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Shield className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <div className="flex items-center justify-between">
                  <Label htmlFor="is_active">Active Status</Label>
                  <Switch
                    id="is_active"
                    checked={formData.is_active}
                    onCheckedChange={(checked) => setFormData({ ...formData, is_active: checked })}
                  />
                </div>
                <p className="text-sm text-muted-foreground">
                  {formData.is_active ? 'Role can be assigned to users' : 'Role is disabled'}
                </p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Calendar className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-1">
                <p className="text-sm text-muted-foreground">Last Updated</p>
                <p className="font-medium text-foreground">
                  {new Date(role.updated_at).toLocaleDateString('en-US', {
                    month: 'long',
                    day: 'numeric',
                    year: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit'
                  })}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

    </form>
  );
});

RoleEditForm.displayName = 'RoleEditForm';