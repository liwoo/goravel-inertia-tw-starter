/**
 * TwoFactorDisable Component
 *
 * Dialog for disabling TOTP two-factor authentication.
 * Requires password and current TOTP code for verification.
 */

import React, { useState } from 'react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription } from '@/components/ui/alert';
import {
  ShieldOff,
  Eye,
  EyeOff,
  AlertCircle,
  Loader2,
  AlertTriangle,
} from 'lucide-react';
import { toast } from 'sonner';
import { disableTOTP } from '@/lib/api/two-factor';
import type { TOTPFormErrors, ApiErrorResponse } from '@/types/account';
import axios from '@/lib/axios';

interface TwoFactorDisableProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
}

export const TwoFactorDisable: React.FC<TwoFactorDisableProps> = ({
  open,
  onOpenChange,
  onSuccess,
}) => {
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [code, setCode] = useState('');
  const [errors, setErrors] = useState<TOTPFormErrors>({});
  const [isLoading, setIsLoading] = useState(false);

  const resetState = () => {
    setPassword('');
    setShowPassword(false);
    setCode('');
    setErrors({});
    setIsLoading(false);
  };

  const handleClose = () => {
    resetState();
    onOpenChange(false);
  };

  const handleSubmit = async () => {
    const newErrors: TOTPFormErrors = {};

    if (!password) {
      newErrors.password = 'Password is required';
    }

    if (!code) {
      newErrors.code = 'Verification code is required';
    } else if (code.length !== 6 && code.length !== 8) {
      // Allow 6-digit TOTP or 8-character backup code
      newErrors.code = 'Enter a 6-digit code or backup code';
    }

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setIsLoading(true);
    setErrors({});

    try {
      await disableTOTP(password, code);
      toast.success('Two-factor authentication has been disabled');
      resetState();
      onOpenChange(false);
      onSuccess();
    } catch (error: unknown) {
      if (axios.isAxiosError(error) && error.response) {
        const errorData = error.response.data as ApiErrorResponse;

        if (error.response.status === 401 || error.response.status === 403) {
          setErrors({ password: 'Incorrect password' });
        } else if (errorData.errors) {
          const apiErrors: TOTPFormErrors = {};
          if (errorData.errors.password) {
            apiErrors.password = errorData.errors.password[0];
          }
          if (errorData.errors.code) {
            apiErrors.code = errorData.errors.code[0];
          }
          if (Object.keys(apiErrors).length > 0) {
            setErrors(apiErrors);
          } else {
            setErrors({ general: errorData.message || 'Failed to disable 2FA' });
          }
        } else {
          setErrors({ general: errorData.message || 'Failed to disable 2FA' });
        }
      } else {
        setErrors({ general: 'An unexpected error occurred' });
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <div className="flex items-center gap-3">
            <div className="p-2 rounded-lg bg-red-100 dark:bg-red-900">
              <ShieldOff className="h-5 w-5 text-red-600 dark:text-red-400" />
            </div>
            <div>
              <DialogTitle>Disable Two-Factor Authentication</DialogTitle>
              <DialogDescription>
                This will reduce the security of your account
              </DialogDescription>
            </div>
          </div>
        </DialogHeader>

        <div className="space-y-4 py-4">
          <Alert variant="destructive">
            <AlertTriangle className="h-4 w-4" />
            <AlertDescription>
              Disabling 2FA will make your account less secure.
              Anyone with your password will be able to access your account.
            </AlertDescription>
          </Alert>

          {errors.general && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>{errors.general}</AlertDescription>
            </Alert>
          )}

          <div className="space-y-2">
            <Label htmlFor="disable-password">Current Password</Label>
            <div className="relative">
              <Input
                id="disable-password"
                type={showPassword ? 'text' : 'password'}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Enter your password"
                className={`pr-10 ${errors.password ? 'border-destructive' : ''}`}
                disabled={isLoading}
              />
              <Button
                type="button"
                variant="ghost"
                size="icon"
                className="absolute right-1 top-1/2 -translate-y-1/2 h-8 w-8"
                onClick={() => setShowPassword(!showPassword)}
                tabIndex={-1}
              >
                {showPassword ? (
                  <EyeOff className="h-4 w-4 text-muted-foreground" />
                ) : (
                  <Eye className="h-4 w-4 text-muted-foreground" />
                )}
              </Button>
            </div>
            {errors.password && (
              <p className="text-sm text-destructive flex items-center gap-1">
                <AlertCircle className="h-3 w-3" />
                {errors.password}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="disable-code">Authentication Code</Label>
            <Input
              id="disable-code"
              type="text"
              inputMode="numeric"
              pattern="[0-9A-Za-z]*"
              maxLength={8}
              value={code}
              onChange={(e) => {
                const value = e.target.value.toUpperCase();
                setCode(value);
              }}
              placeholder="Enter 6-digit code or backup code"
              className={`text-center font-mono tracking-widest ${
                errors.code ? 'border-destructive' : ''
              }`}
              disabled={isLoading}
              autoComplete="one-time-code"
              onKeyDown={(e) => e.key === 'Enter' && handleSubmit()}
            />
            {errors.code && (
              <p className="text-sm text-destructive flex items-center gap-1">
                <AlertCircle className="h-3 w-3" />
                {errors.code}
              </p>
            )}
            <p className="text-xs text-muted-foreground">
              Enter the code from your authenticator app or one of your backup codes.
            </p>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={handleClose} disabled={isLoading}>
            Cancel
          </Button>
          <Button
            variant="destructive"
            onClick={handleSubmit}
            disabled={isLoading || !password || !code}
          >
            {isLoading ? (
              <>
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                Disabling...
              </>
            ) : (
              'Disable 2FA'
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default TwoFactorDisable;
