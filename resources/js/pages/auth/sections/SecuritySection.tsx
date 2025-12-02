import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';
import {
  Lock,
  Eye,
  EyeOff,
  Loader2,
  CheckCircle2,
  AlertCircle,
  ShieldCheck
} from 'lucide-react';
import { toast } from 'sonner';
import axios from '@/lib/axios';
import type {
  ChangePasswordRequest,
  PasswordFormErrors,
  ApiErrorResponse
} from '@/types/account';

interface SecuritySectionProps {
  userId?: number;
}

export const SecuritySection: React.FC<SecuritySectionProps> = () => {
  const [formData, setFormData] = useState<ChangePasswordRequest>({
    current_password: '',
    new_password: '',
    confirm_password: '',
  });
  const [errors, setErrors] = useState<PasswordFormErrors>({});
  const [isSaving, setIsSaving] = useState(false);
  const [showCurrentPassword, setShowCurrentPassword] = useState(false);
  const [showNewPassword, setShowNewPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const validateForm = (): boolean => {
    const newErrors: PasswordFormErrors = {};

    if (!formData.current_password) {
      newErrors.current_password = 'Current password is required';
    }

    if (!formData.new_password) {
      newErrors.new_password = 'New password is required';
    } else if (formData.new_password.length < 8) {
      newErrors.new_password = 'Password must be at least 8 characters';
    } else if (formData.new_password.length > 128) {
      newErrors.new_password = 'Password must not exceed 128 characters';
    } else if (!/[A-Z]/.test(formData.new_password)) {
      newErrors.new_password = 'Password must contain at least one uppercase letter';
    } else if (!/[a-z]/.test(formData.new_password)) {
      newErrors.new_password = 'Password must contain at least one lowercase letter';
    } else if (!/[0-9]/.test(formData.new_password)) {
      newErrors.new_password = 'Password must contain at least one number';
    }

    if (!formData.confirm_password) {
      newErrors.confirm_password = 'Please confirm your new password';
    } else if (formData.new_password !== formData.confirm_password) {
      newErrors.confirm_password = 'Passwords do not match';
    }

    if (formData.current_password && formData.new_password &&
        formData.current_password === formData.new_password) {
      newErrors.new_password = 'New password must be different from current password';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSuccessMessage(null);

    if (!validateForm()) {
      return;
    }

    setIsSaving(true);
    try {
      await axios.put('/api/account/password', formData);

      // Clear form on success
      setFormData({
        current_password: '',
        new_password: '',
        confirm_password: '',
      });
      setErrors({});
      setSuccessMessage('Your password has been changed successfully.');
      toast.success('Password changed successfully');
    } catch (error: unknown) {
      console.error('Failed to change password:', error);
      setSuccessMessage(null);

      if (axios.isAxiosError(error) && error.response) {
        const errorData = error.response.data as ApiErrorResponse;

        if (errorData.errors) {
          const apiErrors: PasswordFormErrors = {};
          if (errorData.errors.current_password) {
            apiErrors.current_password = errorData.errors.current_password[0];
          }
          if (errorData.errors.new_password) {
            apiErrors.new_password = errorData.errors.new_password[0];
          }
          if (errorData.errors.confirm_password) {
            apiErrors.confirm_password = errorData.errors.confirm_password[0];
          }
          setErrors(apiErrors);
        }

        // Handle specific error messages
        if (error.response.status === 401 || error.response.status === 403) {
          setErrors({ current_password: 'Current password is incorrect' });
          toast.error('Current password is incorrect');
        } else {
          toast.error(errorData.message || 'Failed to change password');
        }
      } else {
        toast.error('An unexpected error occurred');
      }
    } finally {
      setIsSaving(false);
    }
  };

  const getPasswordStrength = (password: string): { level: number; label: string; color: string } => {
    if (!password) return { level: 0, label: '', color: '' };

    let score = 0;
    if (password.length >= 8) score++;
    if (password.length >= 12) score++;
    if (/[A-Z]/.test(password)) score++;
    if (/[a-z]/.test(password)) score++;
    if (/[0-9]/.test(password)) score++;
    if (/[^A-Za-z0-9]/.test(password)) score++;

    if (score <= 2) return { level: 1, label: 'Weak', color: 'bg-red-500' };
    if (score <= 4) return { level: 2, label: 'Fair', color: 'bg-yellow-500' };
    if (score <= 5) return { level: 3, label: 'Good', color: 'bg-blue-500' };
    return { level: 4, label: 'Strong', color: 'bg-green-500' };
  };

  const passwordStrength = getPasswordStrength(formData.new_password);

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-lg bg-muted">
            <ShieldCheck className="h-5 w-5 text-muted-foreground" />
          </div>
          <div>
            <CardTitle>Change Password</CardTitle>
            <CardDescription>
              Update your password to keep your account secure
            </CardDescription>
          </div>
        </div>
      </CardHeader>
      <form onSubmit={handleSubmit}>
        <CardContent className="space-y-6">
          {successMessage && (
            <Alert className="border-green-500 bg-green-50 dark:bg-green-950">
              <CheckCircle2 className="h-4 w-4 text-green-600" />
              <AlertDescription className="text-green-700 dark:text-green-300">
                {successMessage}
              </AlertDescription>
            </Alert>
          )}

          {/* Current Password */}
          <div className="space-y-2">
            <Label htmlFor="current_password">Current Password</Label>
            <div className="relative">
              <Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                id="current_password"
                type={showCurrentPassword ? 'text' : 'password'}
                value={formData.current_password}
                onChange={(e) => setFormData({ ...formData, current_password: e.target.value })}
                placeholder="Enter your current password"
                className={`pl-10 pr-10 ${errors.current_password ? 'border-destructive' : ''}`}
                disabled={isSaving}
              />
              <Button
                type="button"
                variant="ghost"
                size="icon"
                className="absolute right-1 top-1/2 -translate-y-1/2 h-8 w-8"
                onClick={() => setShowCurrentPassword(!showCurrentPassword)}
                tabIndex={-1}
              >
                {showCurrentPassword ? (
                  <EyeOff className="h-4 w-4 text-muted-foreground" />
                ) : (
                  <Eye className="h-4 w-4 text-muted-foreground" />
                )}
              </Button>
            </div>
            {errors.current_password && (
              <p className="text-sm text-destructive flex items-center gap-1">
                <AlertCircle className="h-3 w-3" />
                {errors.current_password}
              </p>
            )}
          </div>

          {/* New Password */}
          <div className="space-y-2">
            <Label htmlFor="new_password">New Password</Label>
            <div className="relative">
              <Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                id="new_password"
                type={showNewPassword ? 'text' : 'password'}
                value={formData.new_password}
                onChange={(e) => setFormData({ ...formData, new_password: e.target.value })}
                placeholder="Enter your new password"
                className={`pl-10 pr-10 ${errors.new_password ? 'border-destructive' : ''}`}
                disabled={isSaving}
              />
              <Button
                type="button"
                variant="ghost"
                size="icon"
                className="absolute right-1 top-1/2 -translate-y-1/2 h-8 w-8"
                onClick={() => setShowNewPassword(!showNewPassword)}
                tabIndex={-1}
              >
                {showNewPassword ? (
                  <EyeOff className="h-4 w-4 text-muted-foreground" />
                ) : (
                  <Eye className="h-4 w-4 text-muted-foreground" />
                )}
              </Button>
            </div>
            {formData.new_password && (
              <div className="space-y-2">
                <div className="flex gap-1">
                  {[1, 2, 3, 4].map((level) => (
                    <div
                      key={level}
                      className={`h-1 flex-1 rounded-full transition-colors ${
                        passwordStrength.level >= level
                          ? passwordStrength.color
                          : 'bg-muted'
                      }`}
                    />
                  ))}
                </div>
                <p className={`text-xs ${
                  passwordStrength.level <= 1 ? 'text-red-500' :
                  passwordStrength.level === 2 ? 'text-yellow-600' :
                  passwordStrength.level === 3 ? 'text-blue-500' :
                  'text-green-500'
                }`}>
                  Password strength: {passwordStrength.label}
                </p>
              </div>
            )}
            {errors.new_password && (
              <p className="text-sm text-destructive flex items-center gap-1">
                <AlertCircle className="h-3 w-3" />
                {errors.new_password}
              </p>
            )}
            <p className="text-xs text-muted-foreground">
              Password must be at least 8 characters and contain uppercase, lowercase, and numbers.
            </p>
          </div>

          {/* Confirm Password */}
          <div className="space-y-2">
            <Label htmlFor="confirm_password">Confirm New Password</Label>
            <div className="relative">
              <Lock className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                id="confirm_password"
                type={showConfirmPassword ? 'text' : 'password'}
                value={formData.confirm_password}
                onChange={(e) => setFormData({ ...formData, confirm_password: e.target.value })}
                placeholder="Confirm your new password"
                className={`pl-10 pr-10 ${errors.confirm_password ? 'border-destructive' : ''}`}
                disabled={isSaving}
              />
              <Button
                type="button"
                variant="ghost"
                size="icon"
                className="absolute right-1 top-1/2 -translate-y-1/2 h-8 w-8"
                onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                tabIndex={-1}
              >
                {showConfirmPassword ? (
                  <EyeOff className="h-4 w-4 text-muted-foreground" />
                ) : (
                  <Eye className="h-4 w-4 text-muted-foreground" />
                )}
              </Button>
            </div>
            {errors.confirm_password && (
              <p className="text-sm text-destructive flex items-center gap-1">
                <AlertCircle className="h-3 w-3" />
                {errors.confirm_password}
              </p>
            )}
            {formData.confirm_password && formData.new_password === formData.confirm_password && !errors.confirm_password && (
              <p className="text-sm text-green-600 flex items-center gap-1">
                <CheckCircle2 className="h-3 w-3" />
                Passwords match
              </p>
            )}
          </div>
        </CardContent>
        <CardFooter className="border-t pt-6">
          <Button
            type="submit"
            disabled={isSaving || !formData.current_password || !formData.new_password || !formData.confirm_password}
            className="ml-auto"
          >
            {isSaving ? (
              <>
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                Changing Password...
              </>
            ) : (
              <>
                <Lock className="h-4 w-4 mr-2" />
                Change Password
              </>
            )}
          </Button>
        </CardFooter>
      </form>
    </Card>
  );
};

export default SecuritySection;
