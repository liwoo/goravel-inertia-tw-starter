import React, { useState, forwardRef, useImperativeHandle } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { CrudFormProps } from '@/types/crud';
import { RoleFormData } from '@/types/permissions';
import { Shield, Crown, Users, Hash, FileText } from 'lucide-react';
import { Separator } from '@/components/ui/separator';

interface RoleCreateFormProps extends CrudFormProps {
  roles?: Array<{ id: number; name: string; level: number }>;
  setIsSaving?: (saving: boolean) => void;
}

export const RoleCreateForm = forwardRef<any, RoleCreateFormProps>(({ 
  onSuccess,
  onError,
  onCancel, 
  isLoading = false,
  roles = [],
  setIsSaving
}, ref) => {
  const [formData, setFormData] = useState<RoleFormData>({
    name: '',
    slug: '',
    description: '',
    level: 1,
    is_active: true,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  // Auto-generate slug from name
  const handleNameChange = (name: string) => {
    setFormData({
      ...formData,
      name,
      slug: name.toLowerCase().replace(/\s+/g, '-').replace(/[^a-z0-9-]/g, '')
    });
  };

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
      const response = await fetch('/api/roles', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
        },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        onSuccess('Role created successfully');
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
      {/* Role Icon */}
      <div className="flex justify-center">
        <div className="p-6 rounded-full bg-teal-100 dark:bg-teal-900/30">
          <Shield className="h-16 w-16 text-teal-600 dark:text-teal-400" />
        </div>
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
                  onChange={(e) => handleNameChange(e.target.value)}
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
                  URL-friendly identifier (auto-generated from name)
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
          </div>
        </div>
      </div>

    </form>
  );
});

RoleCreateForm.displayName = 'RoleCreateForm';