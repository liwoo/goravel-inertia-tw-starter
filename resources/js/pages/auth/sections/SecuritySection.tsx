import React, { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import {
  Lock,
  Eye,
  EyeOff,
  Loader2,
  CheckCircle2,
  AlertCircle,
  ShieldCheck,
  Shield,
  ShieldOff,
  Key,
  RefreshCw,
} from 'lucide-react';
import { toast } from 'sonner';
import axios from '@/lib/axios';
import { getTOTPStatus, regenerateBackupCodes } from '@/lib/api/two-factor';
import { TwoFactorSetup } from '@/components/auth/TwoFactorSetup';
import { TwoFactorDisable } from '@/components/auth/TwoFactorDisable';
import type {
  ChangePasswordRequest,
  PasswordFormErrors,
  ApiErrorResponse,
  TOTPStatusData,
} from '@/types/account';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';

interface SecuritySectionProps {
  userId?: number;
}

export const SecuritySection: React.FC<SecuritySectionProps> = () => {
  // Password change state
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

  // 2FA state
  const [totpStatus, setTotpStatus] = useState<TOTPStatusData | null>(null);
  const [isLoadingStatus, setIsLoadingStatus] = useState(true);
  const [showSetupDialog, setShowSetupDialog] = useState(false);
  const [showDisableDialog, setShowDisableDialog] = useState(false);
  const [showRegenerateDialog, setShowRegenerateDialog] = useState(false);
  const [regeneratePassword, setRegeneratePassword] = useState('');
  const [showRegeneratePassword, setShowRegeneratePassword] = useState(false);
  const [isRegenerating, setIsRegenerating] = useState(false);
  const [regeneratedCodes, setRegeneratedCodes] = useState<string[] | null>(null);
  const [regenerateError, setRegenerateError] = useState<string | null>(null);

  // Fetch 2FA status on mount
  useEffect(() => {
    fetchTOTPStatus();
  }, []);

  const fetchTOTPStatus = async () => {
    setIsLoadingStatus(true);
    try {
      const response = await getTOTPStatus();
      // Backend wraps response in { success, message, data } structure
      setTotpStatus(response.data.data);
    } catch (error) {
      console.error('Failed to fetch 2FA status:', error);
      // Set default status on error
      setTotpStatus({ enabled: false, verified_at: null, backup_codes_remaining: 0 });
    } finally {
      setIsLoadingStatus(false);
    }
  };

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

  const handleRegenerateBackupCodes = async () => {
    if (!regeneratePassword) {
      setRegenerateError('Password is required');
      return;
    }

    setIsRegenerating(true);
    setRegenerateError(null);

    try {
      const response = await regenerateBackupCodes(regeneratePassword);
      // Backend wraps response in { success, message, data } structure
      setRegeneratedCodes(response.data.data.backup_codes);
      toast.success('Backup codes regenerated successfully');
      // Refresh status to get new count
      fetchTOTPStatus();
    } catch (error: unknown) {
      if (axios.isAxiosError(error) && error.response) {
        const errorData = error.response.data as ApiErrorResponse;
        if (error.response.status === 401 || error.response.status === 403) {
          setRegenerateError('Incorrect password');
        } else {
          setRegenerateError(errorData.message || 'Failed to regenerate backup codes');
        }
      } else {
        setRegenerateError('An unexpected error occurred');
      }
    } finally {
      setIsRegenerating(false);
    }
  };

  const handleCloseRegenerateDialog = () => {
    setShowRegenerateDialog(false);
    setRegeneratePassword('');
    setShowRegeneratePassword(false);
    setRegeneratedCodes(null);
    setRegenerateError(null);
  };

  const copyBackupCodes = async (codes: string[]) => {
    try {
      await navigator.clipboard.writeText(codes.join('\n'));
      toast.success('Backup codes copied to clipboard');
    } catch {
      toast.error('Failed to copy to clipboard');
    }
  };

  const downloadBackupCodes = (codes: string[]) => {
    const content = [
      'SMEDI Database - Two-Factor Authentication Backup Codes',
      '========================================================',
      '',
      'Store these codes in a safe place. Each code can only be used once.',
      '',
      ...codes.map((code, i) => `${i + 1}. ${code}`),
      '',
      `Generated: ${new Date().toISOString()}`,
    ].join('\n');

    const blob = new Blob([content], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'smedi-2fa-backup-codes.txt';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    toast.success('Backup codes downloaded');
  };

  const passwordStrength = getPasswordStrength(formData.new_password);

  return (
    <div className="space-y-6">
      {/* Password Change Card */}
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

      {/* Two-Factor Authentication Card */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className={`p-2 rounded-lg ${totpStatus?.enabled ? 'bg-purple-100 dark:bg-purple-900' : 'bg-muted'}`}>
                <Shield className={`h-5 w-5 ${totpStatus?.enabled ? 'text-purple-600 dark:text-purple-400' : 'text-muted-foreground'}`} />
              </div>
              <div>
                <CardTitle className="flex items-center gap-2">
                  Two-Factor Authentication
                  {isLoadingStatus ? (
                    <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
                  ) : (
                    <Badge
                      variant={totpStatus?.enabled ? 'default' : 'secondary'}
                      className={totpStatus?.enabled ? 'bg-green-600 hover:bg-green-600' : ''}
                    >
                      {totpStatus?.enabled ? 'Enabled' : 'Disabled'}
                    </Badge>
                  )}
                </CardTitle>
                <CardDescription>
                  Add an extra layer of security to your account using a TOTP authenticator app
                </CardDescription>
              </div>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {!isLoadingStatus && totpStatus?.enabled && (
            <>
              <div className="flex items-center justify-between p-4 bg-muted/50 rounded-lg">
                <div className="flex items-center gap-3">
                  <Key className="h-5 w-5 text-muted-foreground" />
                  <div>
                    <p className="text-sm font-medium">Backup Codes</p>
                    <p className="text-xs text-muted-foreground">
                      {totpStatus.backup_codes_remaining} codes remaining
                    </p>
                  </div>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setShowRegenerateDialog(true)}
                >
                  <RefreshCw className="h-4 w-4 mr-2" />
                  Regenerate
                </Button>
              </div>

              {totpStatus.backup_codes_remaining <= 3 && totpStatus.backup_codes_remaining > 0 && (
                <Alert className="border-amber-500 bg-amber-50 dark:bg-amber-950">
                  <AlertCircle className="h-4 w-4 text-amber-600" />
                  <AlertDescription className="text-amber-700 dark:text-amber-300">
                    You have only {totpStatus.backup_codes_remaining} backup code{totpStatus.backup_codes_remaining !== 1 ? 's' : ''} remaining.
                    Consider regenerating your backup codes.
                  </AlertDescription>
                </Alert>
              )}

              {totpStatus.backup_codes_remaining === 0 && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>
                    You have no backup codes remaining. Regenerate new codes immediately to avoid being locked out.
                  </AlertDescription>
                </Alert>
              )}

              {totpStatus.verified_at && (
                <p className="text-xs text-muted-foreground">
                  Enabled on {new Date(totpStatus.verified_at).toLocaleDateString('en-US', {
                    year: 'numeric',
                    month: 'long',
                    day: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit'
                  })}
                </p>
              )}
            </>
          )}

          {!isLoadingStatus && !totpStatus?.enabled && (
            <div className="text-sm text-muted-foreground">
              <p>
                Two-factor authentication adds an extra layer of security to your account.
                When enabled, you will need to enter a code from your authenticator app
                each time you sign in.
              </p>
              <p className="mt-2">
                We recommend using apps like Google Authenticator, Authy, or 1Password.
              </p>
            </div>
          )}
        </CardContent>
        <CardFooter className="border-t pt-6">
          {!isLoadingStatus && (
            totpStatus?.enabled ? (
              <Button
                variant="destructive"
                onClick={() => setShowDisableDialog(true)}
                className="ml-auto"
              >
                <ShieldOff className="h-4 w-4 mr-2" />
                Disable 2FA
              </Button>
            ) : (
              <Button
                onClick={() => setShowSetupDialog(true)}
                className="ml-auto"
              >
                <Shield className="h-4 w-4 mr-2" />
                Enable 2FA
              </Button>
            )
          )}
        </CardFooter>
      </Card>

      {/* 2FA Setup Dialog */}
      <TwoFactorSetup
        open={showSetupDialog}
        onOpenChange={setShowSetupDialog}
        onSuccess={fetchTOTPStatus}
      />

      {/* 2FA Disable Dialog */}
      <TwoFactorDisable
        open={showDisableDialog}
        onOpenChange={setShowDisableDialog}
        onSuccess={fetchTOTPStatus}
      />

      {/* Regenerate Backup Codes Dialog */}
      <Dialog open={showRegenerateDialog} onOpenChange={handleCloseRegenerateDialog}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <div className="flex items-center gap-3">
              <div className="p-2 rounded-lg bg-purple-100 dark:bg-purple-900">
                <RefreshCw className="h-5 w-5 text-purple-600 dark:text-purple-400" />
              </div>
              <div>
                <DialogTitle>
                  {regeneratedCodes ? 'New Backup Codes' : 'Regenerate Backup Codes'}
                </DialogTitle>
                <DialogDescription>
                  {regeneratedCodes
                    ? 'Store these codes safely - your old codes are now invalid'
                    : 'This will invalidate all existing backup codes'}
                </DialogDescription>
              </div>
            </div>
          </DialogHeader>

          {!regeneratedCodes ? (
            <>
              <div className="space-y-4 py-4">
                <Alert className="border-amber-500 bg-amber-50 dark:bg-amber-950">
                  <AlertCircle className="h-4 w-4 text-amber-600" />
                  <AlertDescription className="text-amber-700 dark:text-amber-300">
                    All existing backup codes will be permanently deleted and replaced with new ones.
                  </AlertDescription>
                </Alert>

                {regenerateError && (
                  <Alert variant="destructive">
                    <AlertCircle className="h-4 w-4" />
                    <AlertDescription>{regenerateError}</AlertDescription>
                  </Alert>
                )}

                <div className="space-y-2">
                  <Label htmlFor="regenerate-password">Confirm Password</Label>
                  <div className="relative">
                    <Input
                      id="regenerate-password"
                      type={showRegeneratePassword ? 'text' : 'password'}
                      value={regeneratePassword}
                      onChange={(e) => setRegeneratePassword(e.target.value)}
                      placeholder="Enter your password"
                      className="pr-10"
                      disabled={isRegenerating}
                      onKeyDown={(e) => e.key === 'Enter' && handleRegenerateBackupCodes()}
                    />
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      className="absolute right-1 top-1/2 -translate-y-1/2 h-8 w-8"
                      onClick={() => setShowRegeneratePassword(!showRegeneratePassword)}
                      tabIndex={-1}
                    >
                      {showRegeneratePassword ? (
                        <EyeOff className="h-4 w-4 text-muted-foreground" />
                      ) : (
                        <Eye className="h-4 w-4 text-muted-foreground" />
                      )}
                    </Button>
                  </div>
                </div>
              </div>

              <DialogFooter>
                <Button variant="outline" onClick={handleCloseRegenerateDialog} disabled={isRegenerating}>
                  Cancel
                </Button>
                <Button onClick={handleRegenerateBackupCodes} disabled={isRegenerating || !regeneratePassword}>
                  {isRegenerating ? (
                    <>
                      <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                      Regenerating...
                    </>
                  ) : (
                    'Regenerate Codes'
                  )}
                </Button>
              </DialogFooter>
            </>
          ) : (
            <>
              <div className="space-y-4 py-4">
                <Alert className="border-amber-500 bg-amber-50 dark:bg-amber-950">
                  <AlertCircle className="h-4 w-4 text-amber-600" />
                  <AlertDescription className="text-amber-700 dark:text-amber-300">
                    <strong>Important:</strong> Each backup code can only be used once.
                    Store them in a secure location.
                  </AlertDescription>
                </Alert>

                <Card>
                  <CardContent className="p-4">
                    <div className="grid grid-cols-2 gap-2">
                      {regeneratedCodes.map((code, index) => (
                        <code
                          key={index}
                          className="p-2 bg-muted rounded text-center font-mono text-sm"
                        >
                          {code}
                        </code>
                      ))}
                    </div>
                  </CardContent>
                </Card>

                <div className="flex gap-2">
                  <Button
                    variant="outline"
                    className="flex-1"
                    onClick={() => copyBackupCodes(regeneratedCodes)}
                  >
                    Copy All
                  </Button>
                  <Button
                    variant="outline"
                    className="flex-1"
                    onClick={() => downloadBackupCodes(regeneratedCodes)}
                  >
                    Download
                  </Button>
                </div>
              </div>

              <DialogFooter>
                <Button onClick={handleCloseRegenerateDialog} className="w-full">
                  I've Saved My Backup Codes
                </Button>
              </DialogFooter>
            </>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default SecuritySection;
