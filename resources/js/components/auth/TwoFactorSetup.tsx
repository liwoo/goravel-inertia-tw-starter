/**
 * TwoFactorSetup Component
 *
 * Multi-step dialog for enabling TOTP two-factor authentication.
 *
 * Steps:
 * 1. Password confirmation
 * 2. QR code display with manual secret
 * 3. TOTP code verification
 * 4. Backup codes display
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
import { Card, CardContent } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';
import {
  Shield,
  Key,
  Copy,
  Download,
  CheckCircle,
  Loader2,
  Eye,
  EyeOff,
  AlertCircle,
  Smartphone,
  QrCode,
} from 'lucide-react';
import { toast } from 'sonner';
import { setupTOTP, verifyTOTP } from '@/lib/api/two-factor';
import type {
  TOTPSetupData,
  TOTPFormErrors,
  ApiErrorResponse,
} from '@/types/account';
import axios from '@/lib/axios';

type SetupStep = 'password' | 'qr-code' | 'verify' | 'backup-codes';

interface TwoFactorSetupProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
}

export const TwoFactorSetup: React.FC<TwoFactorSetupProps> = ({
  open,
  onOpenChange,
  onSuccess,
}) => {
  const [step, setStep] = useState<SetupStep>('password');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [verificationCode, setVerificationCode] = useState('');
  const [setupData, setSetupData] = useState<TOTPSetupData | null>(null);
  const [backupCodes, setBackupCodes] = useState<string[]>([]);
  const [errors, setErrors] = useState<TOTPFormErrors>({});
  const [isLoading, setIsLoading] = useState(false);
  const [copiedSecret, setCopiedSecret] = useState(false);
  const [copiedCodes, setCopiedCodes] = useState(false);

  const resetState = () => {
    setStep('password');
    setPassword('');
    setShowPassword(false);
    setVerificationCode('');
    setSetupData(null);
    setBackupCodes([]);
    setErrors({});
    setIsLoading(false);
    setCopiedSecret(false);
    setCopiedCodes(false);
  };

  const handleClose = () => {
    if (step === 'backup-codes') {
      // Don't allow closing until they acknowledge backup codes
      return;
    }
    resetState();
    onOpenChange(false);
  };

  const handlePasswordSubmit = async () => {
    if (!password) {
      setErrors({ password: 'Password is required' });
      return;
    }

    setIsLoading(true);
    setErrors({});

    try {
      const response = await setupTOTP(password);
      // Backend wraps response in { success, message, data } structure
      setSetupData(response.data.data);
      setStep('qr-code');
    } catch (error: unknown) {
      if (axios.isAxiosError(error) && error.response) {
        const errorData = error.response.data as ApiErrorResponse;
        if (error.response.status === 401 || error.response.status === 403) {
          setErrors({ password: 'Incorrect password' });
        } else if (errorData.errors?.password) {
          setErrors({ password: errorData.errors.password[0] });
        } else {
          setErrors({ general: errorData.message || 'Failed to initiate 2FA setup' });
        }
      } else {
        setErrors({ general: 'An unexpected error occurred' });
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleVerifyCode = async () => {
    if (!verificationCode || verificationCode.length !== 6) {
      setErrors({ code: 'Please enter a 6-digit code' });
      return;
    }

    setIsLoading(true);
    setErrors({});

    try {
      const response = await verifyTOTP(verificationCode);
      // Backend wraps response in { success, message, data } structure
      setBackupCodes(response.data.data.backup_codes);
      setStep('backup-codes');
      toast.success('Two-factor authentication enabled successfully');
    } catch (error: unknown) {
      if (axios.isAxiosError(error) && error.response) {
        const errorData = error.response.data as ApiErrorResponse;
        if (errorData.errors?.code) {
          setErrors({ code: errorData.errors.code[0] });
        } else {
          setErrors({ code: errorData.message || 'Invalid verification code' });
        }
      } else {
        setErrors({ code: 'An unexpected error occurred' });
      }
    } finally {
      setIsLoading(false);
    }
  };

  const copyToClipboard = async (text: string, type: 'secret' | 'codes') => {
    try {
      await navigator.clipboard.writeText(text);
      if (type === 'secret') {
        setCopiedSecret(true);
        setTimeout(() => setCopiedSecret(false), 2000);
      } else {
        setCopiedCodes(true);
        setTimeout(() => setCopiedCodes(false), 2000);
      }
      toast.success(`${type === 'secret' ? 'Secret key' : 'Backup codes'} copied to clipboard`);
    } catch {
      toast.error('Failed to copy to clipboard');
    }
  };

  const downloadBackupCodes = () => {
    const content = [
      'SMEDI Database - Two-Factor Authentication Backup Codes',
      '========================================================',
      '',
      'Store these codes in a safe place. Each code can only be used once.',
      '',
      ...backupCodes.map((code, i) => `${i + 1}. ${code}`),
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

  const handleComplete = () => {
    resetState();
    onOpenChange(false);
    onSuccess();
  };

  const renderPasswordStep = () => (
    <>
      <DialogHeader>
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-lg bg-purple-100 dark:bg-purple-900">
            <Shield className="h-5 w-5 text-purple-600 dark:text-purple-400" />
          </div>
          <div>
            <DialogTitle>Enable Two-Factor Authentication</DialogTitle>
            <DialogDescription>
              Confirm your password to continue with 2FA setup
            </DialogDescription>
          </div>
        </div>
      </DialogHeader>

      <div className="space-y-4 py-4">
        <p className="text-sm text-muted-foreground">
          Two-factor authentication adds an extra layer of security to your account.
          You will need to enter a code from your authenticator app each time you sign in.
        </p>

        {errors.general && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{errors.general}</AlertDescription>
          </Alert>
        )}

        <div className="space-y-2">
          <Label htmlFor="setup-password">Current Password</Label>
          <div className="relative">
            <Input
              id="setup-password"
              type={showPassword ? 'text' : 'password'}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Enter your password"
              className={`pr-10 ${errors.password ? 'border-destructive' : ''}`}
              disabled={isLoading}
              onKeyDown={(e) => e.key === 'Enter' && handlePasswordSubmit()}
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
      </div>

      <DialogFooter>
        <Button variant="outline" onClick={handleClose} disabled={isLoading}>
          Cancel
        </Button>
        <Button onClick={handlePasswordSubmit} disabled={isLoading || !password}>
          {isLoading ? (
            <>
              <Loader2 className="h-4 w-4 mr-2 animate-spin" />
              Verifying...
            </>
          ) : (
            'Continue'
          )}
        </Button>
      </DialogFooter>
    </>
  );

  const renderQRCodeStep = () => (
    <>
      <DialogHeader>
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-lg bg-purple-100 dark:bg-purple-900">
            <QrCode className="h-5 w-5 text-purple-600 dark:text-purple-400" />
          </div>
          <div>
            <DialogTitle>Scan QR Code</DialogTitle>
            <DialogDescription>
              Use your authenticator app to scan this QR code
            </DialogDescription>
          </div>
        </div>
      </DialogHeader>

      <div className="space-y-4 py-4">
        <div className="flex items-start gap-2 text-sm text-muted-foreground">
          <Smartphone className="h-4 w-4 mt-0.5 flex-shrink-0" />
          <p>
            Open your authenticator app (like Google Authenticator, Authy, or 1Password)
            and scan the QR code below.
          </p>
        </div>

        {setupData && (
          <>
            <Card className="bg-white">
              <CardContent className="flex justify-center p-6">
                <img
                  src={setupData.qr_code}
                  alt="2FA QR Code"
                  className="w-48 h-48"
                />
              </CardContent>
            </Card>

            <div className="space-y-2">
              <Label className="text-xs text-muted-foreground">
                Can't scan? Enter this key manually:
              </Label>
              <div className="flex items-center gap-2">
                <code className="flex-1 p-3 bg-muted rounded-md font-mono text-sm break-all">
                  {setupData.secret}
                </code>
                <Button
                  variant="outline"
                  size="icon"
                  onClick={() => copyToClipboard(setupData.secret, 'secret')}
                  className="flex-shrink-0"
                >
                  {copiedSecret ? (
                    <CheckCircle className="h-4 w-4 text-green-600" />
                  ) : (
                    <Copy className="h-4 w-4" />
                  )}
                </Button>
              </div>
            </div>
          </>
        )}
      </div>

      <DialogFooter>
        <Button variant="outline" onClick={handleClose}>
          Cancel
        </Button>
        <Button onClick={() => setStep('verify')}>
          I've Scanned the Code
        </Button>
      </DialogFooter>
    </>
  );

  const renderVerifyStep = () => (
    <>
      <DialogHeader>
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-lg bg-purple-100 dark:bg-purple-900">
            <Key className="h-5 w-5 text-purple-600 dark:text-purple-400" />
          </div>
          <div>
            <DialogTitle>Verify Code</DialogTitle>
            <DialogDescription>
              Enter the 6-digit code from your authenticator app
            </DialogDescription>
          </div>
        </div>
      </DialogHeader>

      <div className="space-y-4 py-4">
        <p className="text-sm text-muted-foreground">
          Enter the verification code displayed in your authenticator app to complete setup.
        </p>

        <div className="space-y-2">
          <Label htmlFor="verification-code">Verification Code</Label>
          <Input
            id="verification-code"
            type="text"
            inputMode="numeric"
            pattern="[0-9]*"
            maxLength={6}
            value={verificationCode}
            onChange={(e) => {
              const value = e.target.value.replace(/\D/g, '');
              setVerificationCode(value);
            }}
            placeholder="000000"
            className={`text-center text-2xl font-mono tracking-widest ${
              errors.code ? 'border-destructive' : ''
            }`}
            disabled={isLoading}
            autoComplete="one-time-code"
            onKeyDown={(e) => e.key === 'Enter' && handleVerifyCode()}
          />
          {errors.code && (
            <p className="text-sm text-destructive flex items-center gap-1">
              <AlertCircle className="h-3 w-3" />
              {errors.code}
            </p>
          )}
        </div>
      </div>

      <DialogFooter>
        <Button variant="outline" onClick={() => setStep('qr-code')} disabled={isLoading}>
          Back
        </Button>
        <Button
          onClick={handleVerifyCode}
          disabled={isLoading || verificationCode.length !== 6}
        >
          {isLoading ? (
            <>
              <Loader2 className="h-4 w-4 mr-2 animate-spin" />
              Verifying...
            </>
          ) : (
            'Verify & Enable'
          )}
        </Button>
      </DialogFooter>
    </>
  );

  const renderBackupCodesStep = () => (
    <>
      <DialogHeader>
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-lg bg-green-100 dark:bg-green-900">
            <CheckCircle className="h-5 w-5 text-green-600 dark:text-green-400" />
          </div>
          <div>
            <DialogTitle>Save Your Backup Codes</DialogTitle>
            <DialogDescription>
              Store these codes safely - you'll need them if you lose access to your authenticator
            </DialogDescription>
          </div>
        </div>
      </DialogHeader>

      <div className="space-y-4 py-4">
        <Alert className="border-amber-500 bg-amber-50 dark:bg-amber-950">
          <AlertCircle className="h-4 w-4 text-amber-600" />
          <AlertDescription className="text-amber-700 dark:text-amber-300">
            <strong>Important:</strong> Each backup code can only be used once.
            Store them in a secure location like a password manager.
          </AlertDescription>
        </Alert>

        <Card>
          <CardContent className="p-4">
            <div className="grid grid-cols-2 gap-2">
              {backupCodes.map((code, index) => (
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
            onClick={() => copyToClipboard(backupCodes.join('\n'), 'codes')}
          >
            {copiedCodes ? (
              <>
                <CheckCircle className="h-4 w-4 mr-2 text-green-600" />
                Copied!
              </>
            ) : (
              <>
                <Copy className="h-4 w-4 mr-2" />
                Copy All
              </>
            )}
          </Button>
          <Button variant="outline" className="flex-1" onClick={downloadBackupCodes}>
            <Download className="h-4 w-4 mr-2" />
            Download
          </Button>
        </div>
      </div>

      <DialogFooter>
        <Button onClick={handleComplete} className="w-full">
          I've Saved My Backup Codes
        </Button>
      </DialogFooter>
    </>
  );

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent
        className="sm:max-w-md"
        showCloseButton={step !== 'backup-codes'}
      >
        {step === 'password' && renderPasswordStep()}
        {step === 'qr-code' && renderQRCodeStep()}
        {step === 'verify' && renderVerifyStep()}
        {step === 'backup-codes' && renderBackupCodesStep()}
      </DialogContent>
    </Dialog>
  );
};

export default TwoFactorSetup;
