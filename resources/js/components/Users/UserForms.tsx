import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { CrudFormProps, CrudEditFormProps, CrudDetailViewProps } from '@/types/crud';
import { User, UserFormData, Role } from '@/types/user';
import { EyeIcon, EyeOffIcon, Shield, UserIcon } from 'lucide-react';
import { router } from '@inertiajs/react';

// Create form component
export function UserCreateForm({ 
  onSuccess,
  onError,
  onCancel, 
  isLoading = false,
  roles = []
}: CrudFormProps & {
  roles?: Role[];
}) {
  const [formData, setFormData] = useState<UserFormData>({
    name: '',
    email: '',
    password: '',
    is_active: true,
    is_super_admin: false,
    role_id: undefined,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const [showPassword, setShowPassword] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    // Basic validation
    const newErrors: Record<string, string> = {};
    if (!formData.name.trim()) {
      newErrors.name = 'Name is required';
    }
    if (!formData.email.trim()) {
      newErrors.email = 'Email is required';
    }
    if (!formData.password.trim()) {
      newErrors.password = 'Password is required';
    } else if (formData.password.length < 8) {
      newErrors.password = 'Password must be at least 8 characters';
    }
    
    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }
    
    setErrors({});
    
    try {
      const response = await fetch('/api/users', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
          'X-Inertia': 'true',
          'X-Inertia-Version': '1.0.0',
        },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        onSuccess('User created successfully');
      } else {
        const errorData = await response.json().catch(() => ({}));
        onError?.(errorData);
      }
    } catch (error) {
      onError?.(error);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="name">Full Name *</Label>
        <Input
          id="name"
          value={formData.name}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          placeholder="Enter user's full name"
          className={errors.name ? 'border-destructive' : ''}
        />
        {errors.name && (
          <p className="text-sm text-destructive">{errors.name}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="email">Email Address *</Label>
        <Input
          id="email"
          type="email"
          value={formData.email}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          placeholder="user@example.com"
          className={errors.email ? 'border-destructive' : ''}
        />
        {errors.email && (
          <p className="text-sm text-destructive">{errors.email}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="password">Password *</Label>
        <div className="relative">
          <Input
            id="password"
            type={showPassword ? 'text' : 'password'}
            value={formData.password}
            onChange={(e) => setFormData({ ...formData, password: e.target.value })}
            placeholder="Enter password (min 8 characters)"
            className={errors.password ? 'border-destructive' : ''}
          />
          <Button
            type="button"
            variant="ghost"
            size="sm"
            className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
            onClick={() => setShowPassword(!showPassword)}
          >
            {showPassword ? (
              <EyeOffIcon className="h-4 w-4" />
            ) : (
              <EyeIcon className="h-4 w-4" />
            )}
          </Button>
        </div>
        {errors.password && (
          <p className="text-sm text-destructive">{errors.password}</p>
        )}
      </div>

      {roles.length > 0 && (
        <div className="space-y-2">
          <Label htmlFor="role">Role</Label>
          <Select
            value={formData.role_id?.toString() || 'none'}
            onValueChange={(value) => setFormData({ ...formData, role_id: value && value !== 'none' ? parseInt(value) : undefined })}
          >
            <SelectTrigger>
              <SelectValue placeholder="Select a role" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="none">No role</SelectItem>
              {roles.map((role) => (
                <SelectItem key={role.id} value={role.id.toString()}>
                  {role.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      )}

      <div className="space-y-4">
        <div className="flex items-center space-x-2">
          <Switch
            id="is_active"
            checked={formData.is_active}
            onCheckedChange={(checked) => setFormData({ ...formData, is_active: checked })}
          />
          <Label htmlFor="is_active">Active</Label>
        </div>

        <div className="flex items-center space-x-2">
          <Switch
            id="is_super_admin"
            checked={formData.is_super_admin}
            onCheckedChange={(checked) => setFormData({ ...formData, is_super_admin: checked })}
          />
          <Label htmlFor="is_super_admin">Super Admin</Label>
          <Shield className="h-4 w-4 text-blue-500" />
        </div>
      </div>

      <div className="flex justify-end space-x-2">
        <Button type="button" variant="outline" onClick={onCancel}>
          Cancel
        </Button>
        <Button type="submit" disabled={isLoading}>
          {isLoading ? 'Creating...' : 'Create User'}
        </Button>
      </div>
    </form>
  );
}

// Edit form component
export function UserEditForm({ 
  item: user,
  onSuccess,
  onError,
  onCancel, 
  isLoading = false,
  roles = []
}: CrudEditFormProps<User> & {
  roles?: Role[];
}) {
  const [formData, setFormData] = useState<UserFormData>({
    name: user.name,
    email: user.email,
    password: '', // Empty for updates
    is_active: user.is_active,
    is_super_admin: user.is_super_admin,
    role_id: user.roles && user.roles.length > 0 ? user.roles[0].id : undefined,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const [showPassword, setShowPassword] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    // Basic validation
    const newErrors: Record<string, string> = {};
    if (!formData.name.trim()) {
      newErrors.name = 'Name is required';
    }
    if (!formData.email.trim()) {
      newErrors.email = 'Email is required';
    }
    if (formData.password && formData.password.length < 8) {
      newErrors.password = 'Password must be at least 8 characters';
    }
    
    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }
    
    setErrors({});
    
    try {
      const response = await fetch(`/api/users/${user.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
          'X-Requested-With': 'XMLHttpRequest',
          'X-Inertia': 'true',
          'X-Inertia-Version': '1.0.0',
        },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        onSuccess('User updated successfully');
      } else {
        const errorData = await response.json().catch(() => ({}));
        onError?.(errorData);
      }
    } catch (error) {
      onError?.(error);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="name">Full Name *</Label>
        <Input
          id="name"
          value={formData.name}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          placeholder="Enter user's full name"
          className={errors.name ? 'border-destructive' : ''}
        />
        {errors.name && (
          <p className="text-sm text-destructive">{errors.name}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="email">Email Address *</Label>
        <Input
          id="email"
          type="email"
          value={formData.email}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          placeholder="user@example.com"
          className={errors.email ? 'border-destructive' : ''}
        />
        {errors.email && (
          <p className="text-sm text-destructive">{errors.email}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="password">New Password (leave blank to keep current)</Label>
        <div className="relative">
          <Input
            id="password"
            type={showPassword ? 'text' : 'password'}
            value={formData.password}
            onChange={(e) => setFormData({ ...formData, password: e.target.value })}
            placeholder="Enter new password (optional)"
            className={errors.password ? 'border-destructive' : ''}
          />
          <Button
            type="button"
            variant="ghost"
            size="sm"
            className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
            onClick={() => setShowPassword(!showPassword)}
          >
            {showPassword ? (
              <EyeOffIcon className="h-4 w-4" />
            ) : (
              <EyeIcon className="h-4 w-4" />
            )}
          </Button>
        </div>
        {errors.password && (
          <p className="text-sm text-destructive">{errors.password}</p>
        )}
      </div>

      {roles.length > 0 && (
        <div className="space-y-2">
          <Label htmlFor="role">Role</Label>
          <Select
            value={formData.role_id?.toString() || 'none'}
            onValueChange={(value) => setFormData({ ...formData, role_id: value && value !== 'none' ? parseInt(value) : undefined })}
          >
            <SelectTrigger>
              <SelectValue placeholder="Select a role" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="none">No role</SelectItem>
              {roles.map((role) => (
                <SelectItem key={role.id} value={role.id.toString()}>
                  {role.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      )}

      <div className="space-y-4">
        <div className="flex items-center space-x-2">
          <Switch
            id="is_active"
            checked={formData.is_active}
            onCheckedChange={(checked) => setFormData({ ...formData, is_active: checked })}
          />
          <Label htmlFor="is_active">Active</Label>
        </div>

        <div className="flex items-center space-x-2">
          <Switch
            id="is_super_admin"
            checked={formData.is_super_admin}
            onCheckedChange={(checked) => setFormData({ ...formData, is_super_admin: checked })}
          />
          <Label htmlFor="is_super_admin">Super Admin</Label>
          <Shield className="h-4 w-4 text-blue-500" />
        </div>
      </div>

      <div className="flex justify-end space-x-2">
        <Button type="button" variant="outline" onClick={onCancel}>
          Cancel
        </Button>
        <Button type="submit" disabled={isLoading}>
          {isLoading ? 'Updating...' : 'Update User'}
        </Button>
      </div>
    </form>
  );
}

// Detail view component
export function UserDetailView({ 
  item: user,
  onEdit,
  onClose,
  canEdit
}: CrudDetailViewProps<User>) {
  return (
    <Card>
      <CardHeader>
        <div className="flex items-center space-x-2">
          {user.is_super_admin ? (
            <Shield className="h-5 w-5 text-blue-500" />
          ) : (
            <UserIcon className="h-5 w-5 text-gray-400" />
          )}
          <CardTitle>{user.name}</CardTitle>
        </div>
        <CardDescription>
          User ID: #{user.id}
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div>
          <Label className="text-sm font-medium">Email</Label>
          <p className="text-sm text-muted-foreground mt-1">
            {user.email}
          </p>
        </div>

        <div>
          <Label className="text-sm font-medium">Role</Label>
          <div className="flex flex-wrap gap-1 mt-1">
            {user.roles && user.roles.length > 0 ? (
              user.roles.map((role) => (
                <Badge key={role.id} variant="secondary">
                  {role.name}
                </Badge>
              ))
            ) : (
              <span className="text-sm text-muted-foreground">No role assigned</span>
            )}
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <Label className="text-sm font-medium">Status</Label>
            <p className="text-sm text-muted-foreground mt-1">
              {user.is_active ? 'Active' : 'Inactive'}
            </p>
          </div>
          <div>
            <Label className="text-sm font-medium">Admin Type</Label>
            <p className="text-sm text-muted-foreground mt-1">
              {user.is_super_admin ? 'Super Admin' : 'Regular User'}
            </p>
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <Label className="text-sm font-medium">Created</Label>
            <p className="text-sm text-muted-foreground mt-1">
              {new Date(user.created_at).toLocaleDateString()}
            </p>
          </div>
          <div>
            <Label className="text-sm font-medium">Updated</Label>
            <p className="text-sm text-muted-foreground mt-1">
              {new Date(user.updated_at).toLocaleDateString()}
            </p>
          </div>
        </div>

        <div className="flex justify-end space-x-2 pt-4">
          {canEdit && onEdit && (
            <Button onClick={onEdit}>
              Edit User
            </Button>
          )}
          <Button variant="outline" onClick={onClose}>
            Close
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}