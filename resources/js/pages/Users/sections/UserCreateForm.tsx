import React, { useState, forwardRef, useImperativeHandle } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { CrudFormProps } from '@/types/crud';
import { UserFormData, Role } from '@/types/user';
import { User, Mail, Shield, Lock, UserCheck } from 'lucide-react';
import { Switch } from '@/components/ui/switch';
import { Separator } from '@/components/ui/separator';
import { useTranslation } from 'react-i18next';

interface UserCreateFormProps extends CrudFormProps {
  roles?: Role[];
  setIsSaving?: (saving: boolean) => void;
}

export const UserCreateForm = forwardRef<any, UserCreateFormProps>(({
  onSuccess,
  onError,
  onCancel,
  isLoading = false,
  roles = [],
  setIsSaving
}, ref) => {
  const { t } = useTranslation('users');
  const [formData, setFormData] = useState<UserFormData>({
    name: '',
    email: '',
    password: '',
    is_active: true,
    is_super_admin: false,
    role_id: undefined,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleSubmit = async () => {

    // Basic validation
    const newErrors: Record<string, string> = {};
    if (!formData.name.trim()) {
      newErrors.name = t('validation.nameRequired');
    }
    if (!formData.email.trim()) {
      newErrors.email = t('validation.emailRequired');
    }
    if (!formData.password.trim()) {
      newErrors.password = t('validation.passwordRequired');
    } else if (formData.password.length < 8) {
      newErrors.password = t('validation.passwordMinLength');
    }

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setErrors({});
    setIsSaving?.(true);

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
        onSuccess(t('toast.created'));
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
      {/* Profile Picture Section */}
      <div className="flex flex-col items-center space-y-4">
        <Avatar className="h-32 w-32 ring-4 ring-background shadow-xl">
          <AvatarImage src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${formData.email || 'default'}`} />
          <AvatarFallback className="text-2xl bg-gradient-to-br from-teal-500 to-teal-600 text-white">
            {formData.name ? formData.name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2) : 'U'}
          </AvatarFallback>
        </Avatar>
        <p className="text-xs text-muted-foreground text-center max-w-xs">
          {t('form.avatarHint')}
        </p>
      </div>

      {/* User Information */}
      <div className="space-y-6">
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">{t('form.userInfo')}</h3>
          <div className="space-y-4">
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <User className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="name">{t('form.fullName')}</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder={t('form.enterFullName')}
                  className={errors.name ? 'border-destructive' : ''}
                />
                {errors.name && (
                  <p className="text-sm text-destructive">{errors.name}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Mail className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="email">{t('form.emailAddress')}</Label>
                <Input
                  id="email"
                  type="email"
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  placeholder={t('form.emailPlaceholder')}
                  className={errors.email ? 'border-destructive' : ''}
                />
                {errors.email && (
                  <p className="text-sm text-destructive">{errors.email}</p>
                )}
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Lock className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="password">{t('form.password')}</Label>
                <Input
                  id="password"
                  type="password"
                  value={formData.password}
                  onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                  placeholder={t('form.enterPassword')}
                  className={errors.password ? 'border-destructive' : ''}
                />
                {errors.password && (
                  <p className="text-sm text-destructive">{errors.password}</p>
                )}
              </div>
            </div>
          </div>
        </div>

        <Separator className="my-6" />

        {/* Account Settings */}
        <div>
          <h3 className="text-lg font-semibold mb-4 text-foreground">{t('form.accountSettings')}</h3>
          <div className="space-y-4">
            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Shield className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <Label htmlFor="role">{t('form.role')}</Label>
                <Select
                  value={formData.role_id?.toString()}
                  onValueChange={(value) => setFormData({ ...formData, role_id: value ? parseInt(value) : undefined })}
                >
                  <SelectTrigger>
                    <SelectValue placeholder={t('form.selectRole')} />
                  </SelectTrigger>
                  <SelectContent>
                    {roles.map((role) => (
                      <SelectItem key={role.id} value={role.id.toString()}>
                        {role.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <UserCheck className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <div className="flex items-center justify-between">
                  <Label htmlFor="is_active">{t('form.accountStatus')}</Label>
                  <Switch
                    id="is_active"
                    checked={formData.is_active}
                    onCheckedChange={(checked) => setFormData({ ...formData, is_active: checked })}
                  />
                </div>
                <p className="text-sm text-muted-foreground">
                  {formData.is_active ? t('form.canLogin') : t('form.cannotLogin')}
                </p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="p-2 rounded-lg bg-muted">
                <Shield className="h-4 w-4 text-muted-foreground" />
              </div>
              <div className="flex-1 space-y-2">
                <div className="flex items-center justify-between">
                  <Label htmlFor="is_super_admin">{t('form.superAdministrator')}</Label>
                  <Switch
                    id="is_super_admin"
                    checked={formData.is_super_admin}
                    onCheckedChange={(checked) => setFormData({ ...formData, is_super_admin: checked })}
                  />
                </div>
                <p className="text-sm text-muted-foreground">
                  {formData.is_super_admin ? t('form.hasFullAccess') : t('form.limitedPermissions')}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

    </form>
  );
});

UserCreateForm.displayName = 'UserCreateForm';