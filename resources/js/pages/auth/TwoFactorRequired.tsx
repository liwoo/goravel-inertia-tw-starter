/**
 * TwoFactorRequired Page
 *
 * Minimal 2FA setup flow for forced enrollment.
 */

import React, { useState } from 'react';
import AuthLayout from '@/layouts/Auth';
import { Head, router, useForm } from '@inertiajs/react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Copy,
  Download,
  CheckCircle,
  Loader2,
  Eye,
  EyeOff,
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

const TwoFactorRequired: React.FC = () => {
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

  const handlePasswordSubmit = async () => {
    if (!password) {
      setErrors({ password: 'Password is required' });
      return;
    }

    setIsLoading(true);
    setErrors({});

    try {
      const response = await setupTOTP(password);
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
      setBackupCodes(response.data.data.backup_codes);
      setStep('backup-codes');
      toast.success('2FA enabled');
    } catch (error: unknown) {
      if (axios.isAxiosError(error) && error.response) {
        const errorData = error.response.data as ApiErrorResponse;
        if (errorData.errors?.code) {
          setErrors({ code: errorData.errors.code[0] });
        } else {
          setErrors({ code: errorData.message || 'Invalid code' });
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
      toast.success('Copied');
    } catch {
      toast.error('Failed to copy');
    }
  };

  const downloadBackupCodes = () => {
    const content = [
      'Backup Codes - Store securely',
      '',
      ...backupCodes.map((code, i) => `${i + 1}. ${code}`),
      '',
      `Generated: ${new Date().toISOString()}`,
    ].join('\n');

    const blob = new Blob([content], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = '2fa-backup-codes.txt';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    toast.success('Downloaded');
  };

  const handleComplete = () => {
    router.visit('/dashboard');
  };

  const handleCancel = () => {
    router.post('/logout');
  };

  const renderPasswordStep = () => (
    <Card>
      <CardHeader className="text-center">
        <CardTitle>Set Up 2FA</CardTitle>
        <CardDescription>
          Two-factor authentication is required. Enter your password to continue.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {errors.general && (
          <p className="text-sm text-destructive text-center">{errors.general}</p>
        )}

        <div className="space-y-2">
          <Label htmlFor="setup-password">Password</Label>
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
            <p className="text-sm text-destructive">{errors.password}</p>
          )}
        </div>

        <Button onClick={handlePasswordSubmit} disabled={isLoading || !password} className="w-full">
          {isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : 'Continue'}
        </Button>

        <Button variant="ghost" onClick={handleCancel} className="w-full text-muted-foreground">
          Cancel and sign out
        </Button>
      </CardContent>
    </Card>
  );

  const renderQRCodeStep = () => (
    <Card>
      <CardHeader className="text-center">
        <CardTitle>Scan QR Code</CardTitle>
        <CardDescription>Scan with your authenticator app</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {setupData && (
          <>
            <div className="flex justify-center p-4 bg-white rounded-lg">
              <img src={setupData.qr_code} alt="QR Code" className="w-40 h-40" />
            </div>

            <div className="space-y-2">
              <Label className="text-xs text-muted-foreground">Manual entry key</Label>
              <div className="flex items-center gap-2">
                <code className="flex-1 p-2 bg-muted rounded text-xs font-mono break-all">
                  {setupData.secret}
                </code>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => copyToClipboard(setupData.secret, 'secret')}
                  className="h-8 w-8 shrink-0"
                >
                  {copiedSecret ? (
                    <CheckCircle className="h-4 w-4 text-green-600" />
                  ) : (
                    <Copy className="h-4 w-4" />
                  )}
                </Button>
              </div>
            </div>

            <Button onClick={() => setStep('verify')} className="w-full">
              Next
            </Button>
          </>
        )}
      </CardContent>
    </Card>
  );

  const renderVerifyStep = () => (
    <Card>
      <CardHeader className="text-center">
        <CardTitle>Enter Code</CardTitle>
        <CardDescription>Enter the 6-digit code from your app</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <Input
            type="text"
            inputMode="numeric"
            pattern="[0-9]*"
            maxLength={6}
            value={verificationCode}
            onChange={(e) => setVerificationCode(e.target.value.replace(/\D/g, ''))}
            placeholder="000000"
            className={`text-center text-2xl font-mono tracking-widest ${errors.code ? 'border-destructive' : ''}`}
            disabled={isLoading}
            autoComplete="one-time-code"
            onKeyDown={(e) => e.key === 'Enter' && handleVerifyCode()}
          />
          {errors.code && (
            <p className="text-sm text-destructive text-center">{errors.code}</p>
          )}
        </div>

        <div className="flex gap-2">
          <Button variant="outline" onClick={() => setStep('qr-code')} disabled={isLoading} className="flex-1">
            Back
          </Button>
          <Button onClick={handleVerifyCode} disabled={isLoading || verificationCode.length !== 6} className="flex-1">
            {isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : 'Verify'}
          </Button>
        </div>
      </CardContent>
    </Card>
  );

  const renderBackupCodesStep = () => (
    <Card>
      <CardHeader className="text-center">
        <CardTitle>Backup Codes</CardTitle>
        <CardDescription>Save these codes securely. Each can only be used once.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-2 gap-2 p-3 bg-muted rounded-lg">
          {backupCodes.map((code, index) => (
            <code key={index} className="p-1.5 bg-background rounded text-center font-mono text-sm">
              {code}
            </code>
          ))}
        </div>

        <div className="flex gap-2">
          <Button variant="outline" className="flex-1" onClick={() => copyToClipboard(backupCodes.join('\n'), 'codes')}>
            {copiedCodes ? <CheckCircle className="h-4 w-4 mr-1.5 text-green-600" /> : <Copy className="h-4 w-4 mr-1.5" />}
            Copy
          </Button>
          <Button variant="outline" className="flex-1" onClick={downloadBackupCodes}>
            <Download className="h-4 w-4 mr-1.5" />
            Download
          </Button>
        </div>

        <Button onClick={handleComplete} className="w-full">
          Done
        </Button>
      </CardContent>
    </Card>
  );

  return (
    <AuthLayout>
      <Head title="Set Up 2FA" />
      <div className="w-full">
        {step === 'password' && renderPasswordStep()}
        {step === 'qr-code' && renderQRCodeStep()}
        {step === 'verify' && renderVerifyStep()}
        {step === 'backup-codes' && renderBackupCodesStep()}

        {/* Step indicator */}
        <div className="flex justify-center mt-6 gap-1.5">
          {['password', 'qr-code', 'verify', 'backup-codes'].map((s, i) => (
            <div
              key={s}
              className={`h-1.5 w-1.5 rounded-full transition-colors ${
                ['password', 'qr-code', 'verify', 'backup-codes'].indexOf(step) >= i
                  ? 'bg-primary'
                  : 'bg-muted'
              }`}
            />
          ))}
        </div>
      </div>
    </AuthLayout>
  );
};

export default TwoFactorRequired;
